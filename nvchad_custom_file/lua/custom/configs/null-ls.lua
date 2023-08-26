local null_ls = require "null-ls"

local b = null_ls.builtins

local sources = {

  ------------------------------FORMATTER-----------------------
  -- webdev stuff
  b.formatting.deno_fmt, -- choosed deno for ts/js files cuz its very fast!
  b.formatting.prettier.with { filetypes = { "html", "markdown", "css" } }, -- so prettier works only on these filetypes
  b.formatting.stylua, --lua
  b.formatting.clang_format, --cpp
  b.formatting.shfmt, --bash
  b.formatting.black, --python

  -------------------------------LINTER------------------------
  b.diagnostics.flake8, --python
  b.diagnostics.shellcheck, --shell
}

null_ls.setup {
  debug = true,
  sources = sources,
}
