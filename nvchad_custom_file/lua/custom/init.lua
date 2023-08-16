-- local autocmd = vim.api.nvim_create_autocmd

-- Auto resize panes when resizing nvim window
-- autocmd("VimResized", {
--   pattern = "*",
--   command = "tabdo wincmd =",
-- })

-- merge vimrc configuration here.
vim.opt.scrolloff = 10

-- search option.
vim.opt.ignorecase = true
vim.opt.incsearch = true
vim.opt.showmatch = true

-- recover the pos we edit last time we leave
vim.cmd([[autocmd BufReadPost * if line("'\"") > 0 && line("'\"") <= line("$") | exe "normal! g`\"" | endif]])

-- incase can not del the special chara
vim.opt.backspace = "indent,eol,start"

