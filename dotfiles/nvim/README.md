# nvchad custom file

If u want to know how to transfer you .vimrc to .lua, read this [Everything you need to know to configure neovim using lua.](https://vonheikemen.github.io/devlog/tools/configuring-neovim-using-lua/).

Before using those config file, read [office doc](https://nvchad.com/docs/config/walkthrough) first, that way we only update those config in custom (unless necessary).

This dir contain(not finish yet).

- mappings add some keymap for daily use.
- init add some vim setting
- chadrc add some about the ui.

If dont need to modify the core config file, I'll continue update this or I will fork the repo.

TODO: 

- add proxy for treesitter update.

## How to config lsp server

Those requirements should be satisfied if we want lsp-server work well on neovim:

**INSTALL LSP-SERVER by MASON-PLUGIN**: edit your plugin.init, add those langs which you want get lsp-server on the mason config. like below:

```lua
    opts = {
      ensure_installed = {
        "cmake-language-server",
        "html-lsp",
        "powershell-editor-services",
        "pyright",
        "sqlls",
        "bash-language-server",
        "marksman",
      }
    }
```

- if u dont know what lsp-server there are, check it on this [website](https://github.com/neovim/nvim-lspconfig/blob/master/doc/server_configurations.md), but those name may not match which in Mason.
- linter and formatter can also install by mason.

**CONFIG LSP-SERVER**: edit custom/configs/lspconfig.lua, **NORMALLY** using the default setting is enough for us, so we can just easilly set uo the lsp here, now the name should match which on the website.(refer to same object, though the name is not excatly the same with mason).

check [this](https://nvchad.com/docs/config/lsp) out, now the part will be like:

```lua
local servers = { "html", "cssls", "tsserver", "clangd", "bashls", "cmake", "lua_ls", "pyright", "sqlls"}
```

after all those, we can open *.[filetyle] then type `:lsp_info` to checkout the server work or not.
Also the server will show up on right corner automaticlly.

## How to config format.

Just like config the lsp-server, config format function still have two step: 

ONE: **INSTALL FORMATTER**: using mason install formatter like black(python), pretter, clang_format(cpp), stylelua(lua) ...

STILL, add those on ensure_installed of mason, then type :MasonInstallAll and done.

TWO: **CONFIG NULLLS**: null-ls help those formatter work like 'server'(like lsp server), after config it, when open *.[filetype] , neovim will use it to ask server from from those formatter we had installed. 

> if you donnot know what formatter you need, you can open a *.[filetype] then type `:NullLsInfo` , then it'll show what this filetype need and whether this file has been supported or not.


NOW come to what we should edit: `custom/config/null-ls.lua`, add formatter we have installed like the bash one below:

```lua
local sources = {

  -- webdev stuff
  b.formatting.deno_fmt, -- choosed deno for ts/js files cuz its very fast!
  b.formatting.prettier.with { filetypes = { "html", "markdown", "css" } }, -- so prettier works only on these filetypes

  -- Lua
  b.formatting.stylua,

  -- cpp
  b.formatting.clang_format,

  -- bash
  b.formatting.shfmt,
}
```

this example show formatter for cpp,lua,bash, and using prettier to support html,markdown,css.


## How to config linter

Config add install linter is totally same with format.

1. **USING MASON 2 INSTALL THOSE LINTER WE NEED**. we dont talk it again.
2. **CONFIG IT in NULLLS**, there are a little bit different, the keyword we add linter in null-ls is not linter or formatting, it's **diagnostics**, so if we using flake8 for python and, shellcheck for bash, it'll be like:

```lua
local sources = {

  ------------------------------FORMATTER-----------------------
  b.formatting.black, --python

  -------------------------------LINTER------------------------
  b.diagnostics.flake8, --python
  b.diagnostics.shellcheck, --shell
}
```

done!

## How to add snippets

nvchad have config friendly-snippets already ising **luasnip** which support multi-style of snippets including vscode.

ref: we can learn about how to defined a snippets for self.

- [luasnip](https://github.com/L3MON4D3/LuaSnip/blob/master/DOC.md#loaders)
- [friendly-snippets](https://github.com/rafamadriz/friendly-snippets/tree/main): 
- [nvchad-snippets](https://nvchad.com/docs/config/snippets)

when we want to add personal snippets, there are some points we need to notics.

1. CONFIG SNIPPETS PATH IN **custom/init.lua**(intro vscode path only here), but remerber, here should be absolute path, you can write it like below:

```lua
-- snippets
-- 1. SHOULD BE ABSOLUTE PATH, CAN ALSO USE stdpath to CONCAT STRING
-- 2. snippets should have package.json to defined all the langs's snippets' path
vim.g.vscode_snippets_path = vim.fn.stdpath "config" .. "/lua/custom/snippets"
vim.g.vscode_snippets_path = "~/.config/nvim/lua/custom/snippets/"

```

one of those two is ok. other snippets path is same.

2. IN the snippets path, we should have package.json to index all the  snippet path of langs.
3. add legal snippets for langs you need.

EXAMPLE:

- `<...>/.config/nvim/lua/custom/snippets` I have vscode(dir) && package.json
- `<...>/custom/snippets/vscode` I have markdown.json && all.json


packages.json :

```json
{
	"name": "language",
	"contributes": {
		"snippets": [
			{
				"language": [
					"all"
				],
				"path": "./vscode/all.json"
			},
			{
				"language": [
					"markdown"
				],
				"path": "./vscode/markdown.json"
			}
		]
	}
}
```

markdown.json :

```json
{
  "docs title": {
    "prefix": "meta",
    "body": [
      "---",
      "title: $1",
      "catalog: true",
      "toc: true",
      "date: $2",
      "subtitle: $3",
      "lang: cn",
      "cover: /img/header_img/lml_bg2.jpg",
      "thumbnail: /img/header_img/lml_bg2.jpg",
      "tag:",
      "- $4",
      "categories:",
      "- $5",
      "mathjax: false",
      "---"
    ],
    "description": "docs title"
  }
}
```

about how to write snippets can checkout it's offial website, here provide some tool for write vscode snippets faster: [snippets-generator](https://snippet-generator.app/)

## FI
