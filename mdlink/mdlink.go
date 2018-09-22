package main

import (
	"bufio"
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

type message struct {
	Hash         string  `json:"hash"`
	URL          string  `json:"url"`
	APIEndpoint  string  `json:"api_endpoint"`
	Token        string  `json:"token"`
	Start        float64 `json:"start"`
	End          float64 `json:"end"`
	MarkdownLink string  `json:"markdown_link"`
}

func handleConn(c net.Conn) {
	defer c.Close()
	s := bufio.NewScanner(c)
	for s.Scan() {
		var v [2]interface{}
		err := json.Unmarshal(s.Bytes(), &v)
		if err != nil {
			log.Println(err)
			return
		}
		m := &message{
			Hash:        v[1].(map[string]interface{})["hash"].(string),
			URL:         v[1].(map[string]interface{})["url"].(string),
			APIEndpoint: v[1].(map[string]interface{})["api_endpoint"].(string),
			Token:       v[1].(map[string]interface{})["token"].(string),
			Start:       v[1].(map[string]interface{})["start"].(float64),
			End:         v[1].(map[string]interface{})["end"].(float64),
		}
		go func(m *message) {
			err := m.createMarkdownLink()
			if err != nil {
				log.Println(err)
				return
			}
			v[1] = m
			err = json.NewEncoder(c).Encode(v)
			if err != nil {
				log.Println(err)
				return
			}
		}(m)
	}
}

func (m *message) createMarkdownLink() error {
	var title = ""
	var err error

	if m.Token == "" {
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

	m.MarkdownLink = fmt.Sprintf("[%s](%s)", title, m.URL)

	return nil
}

func (m *message) pageTitle() (string, error) {
	res, err := http.Get(m.URL)
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

func (m *message) issueTitle() (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", m.APIEndpoint, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "token "+m.Token)
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
		return m.URL, nil
	}
	return title.(string), nil
}
