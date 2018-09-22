let s:save_cpo = &cpo
set cpo&vim

let s:github_url = "https://github.com/"
let s:github_api_url = "https://api.github.com/repos/"
let s:ghe_options = ['ghe_api_url', 'ghe_token', 'ghe_url']

function! mdlink#make_markdown_link(is_only_on_cursor) range
  let s:done = 0
  let s:call_cnt = 0
  let messages = []

  if a:is_only_on_cursor
    let url = s:get_url()
    if url != ''
      let messages += [s:create_message(url, col("."), line("."), line("."))]
      call s:to_markdown(messages)
    endif
    return
  endif

  for row in range(a:firstline, a:lastline)
    let messages += s:create_messages(row, a:firstline, a:lastline)
  endfor
  call s:to_markdown(messages)
endfunction

function! s:has_ghe_options() abort
  let l:cnt = 0
  for opt in s:ghe_options
    if has_key(g:vim_mdlink, opt)
      let l:cnt += 1
    endif
  endfor
  if l:cnt == len(s:ghe_options)
    return 1
  endif
  return 0
endfunction

function! s:get_url() abort
  let tmp = @r
  execute 'normal "ryiu'
  let url = @r
  let @r = tmp
  return url
endfunction

function! s:create_message(url, col, start, end) abort
  let hash = sha256(a:url . string(a:col))
  call s:replace_to_hash(hash)
  return {"hash": hash, "url": a:url, "api_endpoint": s:api_endpoint(a:url), "token": s:token(a:url), "start": a:start, "end": a:end}
endfunction

function! s:api_endpoint(url) abort
  if s:is_github(a:url) || s:is_ghe(a:url)
    let apiUrl = substitute(a:url, s:github_url, s:github_api_url, 'g')
    if s:has_ghe_options()
      let apiUrl = substitute(apiUrl, g:vim_mdlink['ghe_url'], g:vim_mdlink['ghe_api_url'], 'g')
    endif
    let apiUrl = substitute(apiUrl, '/pull/', '/pulls/', 'g')
    return apiUrl
  endif
  return a:url
endfunction

function! s:is_github(url) abort
  if a:url =~# s:github_url
    return 1
  endif
  return 0
endfunction

function! s:is_ghe(url) abort
  if a:url =~# g:vim_mdlink['ghe_url']
    return 1
  endif
  return 0
endfunction

function! s:token(url) abort
  if s:is_github(a:url)
    return g:vim_mdlink['github_token']
  endif
  if s:is_ghe(a:url)
    return g:vim_mdlink['ghe_token']
  endif
  return ''
endfunction

function! s:create_messages(row, start, end) abort
  let messages = []
  let col = s:get_url_col_position(a:row)
  while col != -1
    call cursor(a:row, col)
    let url = s:get_url()
    if url != ''
      call add(messages, s:create_message(url, col, a:start, a:end))
    endif
    let col = s:get_url_col_position(a:row)
  endwhile
  return messages
endfunction

function! s:get_url_col_position(row) abort
  let pos = match(getline(a:row), '\v[[(]@<!https?', 0)
  if pos > 0
    let pos += 1
  endif
  return pos
endfunction

function! s:replace_to_hash(hash) abort
  let tmp = @r
  let @r = a:hash
  execute 'normal viu"rP'
  let @r = tmp
endfunction

function! s:get_title(channel, msg) abort
  let cmd = ':' . a:msg.start . ',' . a:msg.end . 's;' . a:msg.hash . ';' . a:msg.markdown_link . ';g'
  execute cmd
  let s:done += 1
endfunction


function! s:to_markdown(messages) abort
  let channel = ch_open("127.0.0.1:11111", {"mode": "json"})
  for message in a:messages
    call ch_sendexpr(channel, message, {"callback": function("s:get_title")})
    let s:call_cnt += 1
  endfor

  while 1
    if s:done == s:call_cnt
      break
    endif
    sleep 1ms
  endwhile
endfunction

let &cpo = s:save_cpo
unlet s:save_cpo
