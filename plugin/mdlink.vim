if exists("g:loaded_vim_mdlink")
  finish
endif
let g:loaded_vim_mdlink = 1

try
  execute 'source $HOME/.vim-mdlink'
catch
  echoe "$HOME/.vim-mdlink not found. Check usage in README.md."
  unlet g:loaded_vim_mdlink
  finish
endtry

if !exists('g:vim_mdlink')
  echoe "g:vim_mdlink not defined. Check usage in README.md."
  unlet g:loaded_vim_mdlink
  finish
endif

let s:save_cpo = &cpo
set cpo&vim

let s:base = expand('<sfile>:h:h:gs?\\?/?')
let s:cmd = s:base . '/mdlink/mdlink' . (has('win32') ? '.exe' : '')
if !filereadable(s:cmd)
  execute(":cd " . s:base)
  call system("cd mdlink && go get -d && go build")
endif
let job = job_start(s:cmd)

command! MarkdownLinkOnlyOnCursor call mdlink#make_markdown_link(1)
command! -range MarkdownLink <line1>,<line2>call mdlink#make_markdown_link(0)

let &cpo = s:save_cpo
unlet s:save_cpo
