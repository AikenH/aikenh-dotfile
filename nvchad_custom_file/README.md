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

## How to config linter & format.



