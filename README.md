# vim-mdlink

Vim plugin to convert URL to Markdown style with title.

## Install

Checkout into your plugin directory. Or use plugin manager.

## Requirements

* Vim 8.0+ (with +channel +job)
* [Go](https://golang.org/)
* [kana/vim-textobj-user](https://github.com/kana/vim-textobj-user)
* [mattn/vim-textobj-url](https://github.com/mattn/vim-textobj-url)

## Settings

1. Create $HOME/.vim-mdlink file like the following. And input your github personal token.

        let g:vim_mdlink = {
          \ 'github_token':   'your personal token'
        \}

1. If using GHE, add optional settings about GHE.

        let g:vim_mdlink = {
          \ 'github_token':   'your personal token',
          \ 'ghe_url':        'https://ghe.domain.name/',
          \ 'ghe_api_url':    'https://ghe.domain.name/api/v3/repos/',
          \ 'ghe_token':      'your personal token'
        \}

## Usage

### MarkdownLink

1. Put the cursor on the row existing URL. Or select multi rows by visual mode.
1. Enter command `:MarkdownLink`.

### MarkdownLinkOnlyOnCursor

1. Put the cursor on the URL.
1. Enter command `:MarkdownLinkOnlyOnCursor`.

### Convert all URL

1. Enter command `:%Markdownlink`.

### Map example

```
nnoremap <silent> ml :MarkdownLink<CR>
vnoremap <silent> ml :MarkdownLink<CR>
nnoremap <silent> mo :MarkdownLinkOnlyOnCursor<CR>
```

## License

MIT

## Authors

[genkiroid](https://github.com/genkiroid)

