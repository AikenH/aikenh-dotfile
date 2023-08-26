local on_attach = require("plugins.configs.lspconfig").on_attach
local capabilities = require("plugins.configs.lspconfig").capabilities

local lspconfig = require "lspconfig"

-- if you just want default config for the servers then put them in a table
-- local servers = { "html", "cssls", "tsserver", "clangd", "bashls", "cmake", "lua_ls", "pyright", "sqlls"}
-- FIXME: we delete lua_ls here because lua_ls have setted in the default lspconfig and add some useful setting.
local servers = { "html", "cssls", "tsserver", "clangd", "bashls", "cmake", "pyright", "sqlls"}

for _, lsp in ipairs(servers) do
  lspconfig[lsp].setup {
    on_attach = on_attach,
    capabilities = capabilities,
  }
end

-- 
-- lspconfig.pyright.setup { blabla}
