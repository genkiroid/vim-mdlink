package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:11111")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

type Body struct {
	Hash         string  `json:"hash"`
	URL          string  `json:"url"`
	APIEndpoint  string  `json:"api_endpoint"`
	Token        string  `json:"token"`
	Start        float64 `json:"start"`
	End          float64 `json:"end"`
	MarkdownLink string  `json:"markdown_link"`
}

type Message struct {
	ID   float64
	Body Body
}

func handleConn(c net.Conn) {
	defer c.Close()
	s := bufio.NewScanner(c)
	for s.Scan() {
		var m Message
		err := json.Unmarshal(s.Bytes(), &m)
		if err != nil {
			log.Println(err)
			return
		}
		go func(m *Message) {
			err := m.createMarkdownLink()
			if err != nil {
				log.Println(err)
				return
			}
			err = json.NewEncoder(c).Encode(m)
			if err != nil {
				log.Println(err)
				return
			}
		}(&m)
	}
}

func (m *Message) UnmarshalJSON(b []byte) error {
	var v [2]interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	m.ID = v[0].(float64)
	m.Body = Body{
		Hash:         v[1].(map[string]interface{})["hash"].(string),
		URL:          v[1].(map[string]interface{})["url"].(string),
		APIEndpoint:  v[1].(map[string]interface{})["api_endpoint"].(string),
		Token:        v[1].(map[string]interface{})["token"].(string),
		Start:        v[1].(map[string]interface{})["start"].(float64),
		End:          v[1].(map[string]interface{})["end"].(float64),
		MarkdownLink: "",
	}
	return nil
}

func (m *Message) MarshalJSON() ([]byte, error) {
	var v [2]interface{}
	v[0] = m.ID
	v[1] = m.Body
	return json.Marshal(v)
}

func (m *Message) createMarkdownLink() error {
	var title = ""
	var err error

	if m.Body.Token == "" {
		title, err = m.pageTitle()
		if err != nil {
			return err
		}
	} else {
		title, err = m.issueTitle()
		if err != nil {
			return err
		}
	}

	m.Body.MarkdownLink = fmt.Sprintf("[%s](%s)", title, m.Body.URL)

	return nil
}

func (m *Message) pageTitle() (string, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	res, err := client.Get(m.Body.URL)
	if err != nil {
		return "", err
	}
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", err
	}
	doc, err := html.Parse(strings.NewReader(string(content)))
	if err != nil {
		return "", err
	}
	var title = ""
	var f func(n *html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" && n.Parent.Data == "head" {
			title = n.FirstChild.Data
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return title, nil
}

func (m *Message) issueTitle() (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", m.Body.APIEndpoint, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "token "+m.Body.Token)
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", err
	}
	var v interface{}
	err = json.Unmarshal(content, &v)
	if err != nil {
		return "", err
	}
	title := v.(map[string]interface{})["title"]
	if title == nil {
		return m.Body.URL, nil
	}
	return title.(string), nil
}
