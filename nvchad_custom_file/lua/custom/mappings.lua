---@type MappingsTable
local M = {}

M.general = {
  n = {
    [";"] = { ":", "enter command mode", opts = { nowait = true } },
  },
}

M.indent = {
  v = {
    ["<tab>"] = { ">gv", "tab this line in visual mode" },
    ["<S-tab>"] = { "<gv", "del tab for this line in visual mode"},
  },
}

function ToggleNumber()
  if vim.wo.relativenumber == true then
    vim.opt.relativenumber = false
    vim.opt.number = true
  else
    vim.opt.relativenumber  = true
  end
end

function CloseNumber()
  if vim.wo.relativenumber == true or vim.wo.number == true then
    vim.wo.relativenumber = false
    vim.wo.number = false
  else
    vim.wo.number = true
  end
end

function ToggleWrap()
  if vim.wo.wrap == true then
    vim.wo.wrap = false
  else
    vim.wo.wrap = true
  end
end

-- may need figure it out whethere we still set paste in neovim?
function TogglePaste()
  if vim.o.paste == false then
    vim.opt.paste = true
    print("PASTE change to", vim.o.paste)
  else
    vim.opt.paste = false
    print("PASTE change to", vim.o.paste)
  end
end

M.line = {
  n = {
    ["<leader>lr"] = { ToggleNumber, "ChangeNumberMode" },
    ["<leader>ln"] = { CloseNumber, "CloseNumber" },
    ["<leader>lw"] = { ToggleWrap, "Toggle Line Wrap"},
    ["<leader>lp"] = { TogglePaste, "Toggle paste mode"},
  },
}

M.format = {
  n = {
    ["<leader>rt"] = { "<cmd>retab!<CR>", "retab! this file "},
    ["<leader>re"] = { "<cmd>%s/\\s\\+$//g<CR>", "del the space/tab in the end of line."},
  }
}

M.pagemove = {
  n = {
    ["<leader>pl"] = { "<PageUp>", "pageup" },
    ["<leader>pp"] = {"<PageDown>", "pagedown"},
  }
}
-- more keybinds!
return M
