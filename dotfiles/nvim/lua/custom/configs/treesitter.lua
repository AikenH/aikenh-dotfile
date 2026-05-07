require("nvim-treesitter.install").command_extra_args = {
  curl = { "--proxy", "http://172.30.240.1:7890" },
}

local options = {
  ensure_installed = {
    -- defaults
    "vim",
    "vimdoc",
    "lua",
    "c",
    "css",
    "html",
    "javascript",
    "markdown",
    "markdown_inline",
    "query",
    "tsx",
    "typescript",

    -- personal dev backend
    "python",
    "cpp",
    "go",
    "bash",

    -- frontent dev
    "vue",

    -- markup langs
    "yaml",
    "json",

    -- database
    "sql",
  },
}

return options
