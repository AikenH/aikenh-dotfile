-- Core Basic Setting.
local wezterm = require 'wezterm'
local c = {}
if wezterm.config_builder then
    c = wezterm.config_builder()
end

-- default windows size.
c.initial_cols = 100
c.initial_rows = 24
c.enable_scroll_bar = true
c.window_padding = {
    left = 10, right = 10, top = 10, bottom = 0
}

-- .
c.status_update_interval = 1000
c.automatically_reload_config = true

-- windows behavior
c.window_close_confirmation = 'NeverPrompt'
c.exit_behavior = 'CloseOnCleanExit' -- if the shell program exited with a successful status
c.exit_behavior_messaging = 'Verbose'
c.audible_bell = 'Disabled'
c.animation_fps = 60
c.max_fps = 60
-- c.front_end = "WebGpu"
-- c.webgpu_power_preference = "HighPerformance"
-- local gpu_adapters = require('utils.gpu_adapter')
-- c.webgpu_preferred_adapter = gpu_adapters:pick_best()

c.inactive_pane_hsb = {
    saturation = 0.9,
    brightness = 0.5,
}

-- font setting
c.font_size = 13
c.font = wezterm.font "Monaspace Radon"

-- opacity and light dark auto switch.
c.window_background_opacity = 0.7
c.text_background_opacity = 0.9

function get_appearance()
    if wezterm.gui then
        return wezterm.gui.get_appearance()
    end
    return 'Dark'
end

function scheme_for_appearance(appearance)
    if appearance:find 'Dark' then
        return "Tangoesque (terminal.sexy)"
        -- OneHalfDark
    else
        c.window_background_opacity = 1
        c.text_background_opacity = 0.8
     -- OneHalfLight
        return 'One Light (base16)'
    end
end
c.color_scheme = scheme_for_appearance(get_appearance())

-- tab_bar setting.
c.window_decorations = "INTEGRATED_BUTTONS|RESIZE"
c.use_fancy_tab_bar = false
c.enable_tab_bar = true
c.hide_tab_bar_if_only_one_tab = false
c.tab_bar_at_bottom = true
c.tab_max_width = 25
c.window_frame = {
    active_titlebar_bg = '#123456',
}
c.show_tab_index_in_tab_bar = true
c.switch_to_last_active_tab_when_closing_tab = true



require('tabbar.new-tab-button').setup()
require('tabbar.left-status').setup()
require('tabbar.right-status').setup()
require('tabbar.tab-title').setup()

-- KEYBIND
local platform = require("utils.platform")()
local act = wezterm.action

local mod = {}

if platform.is_mac then
  mod.SUPER = "SUPER"
  mod.SUPER_REV = "SUPER|CTRL"
elseif platform.is_win then
  mod.SUPER = "ALT" -- to not conflict with Windows key shortcuts
  mod.SUPER_REV = "ALT|CTRL"
end

-- launch_menu & default terminal setting.
-- which should be change by platform.
if platform.is_win then
    c.launch_menu = {
        {
            label = "WSL2",
            args = { 'wsl.exe', '-u', 'aikenhong', '--cd', '/home/aikenhong' },
        },
        {   label = "PWSH7", args = { "pwsh.exe" }   },
        {   label = "CMD",args = { "cmd.exe" }}
    }
    c.default_prog = { "wsl.exe ", "-u", "aikenhong", "--cd", "/home/aikenhong"  }
elseif platform.is_mac then
  c.tab_bar_at_bottom = false
end

local keys = {
    -- misc/useful --
    { key = "F1", mods = "NONE", action = "ActivateCopyMode" },
    { key = "F2", mods = "NONE", action = act.ActivateCommandPalette },
    { key = "F3", mods = "NONE", action = act.ShowLauncher },
    { key = "F4", mods = "NONE", action = act.ShowTabNavigator },
    { key = "F11", mods = "NONE", action = act.ToggleFullScreen },
    { key = "F12", mods = "NONE", action = act.ShowDebugOverlay },
    { key = "f", mods = mod.SUPER, action = act.Search({ CaseInSensitiveString = "" }) },

    -- copy/paste --
    { key = "c", mods = "CTRL|SHIFT", action = act.CopyTo("Clipboard") },
    { key = "v", mods = "CTRL|SHIFT", action = act.PasteFrom("Clipboard") },

    -- tabs --
    -- tabs: spawn+close
    { key = "t", mods = mod.SUPER, action = act.SpawnTab("DefaultDomain") },
    { key = "t", mods = mod.SUPER_REV, action = act.SpawnTab({ DomainName = "WSL:Ubuntu" }) },
    { key = "w", mods = mod.SUPER_REV, action = act.CloseCurrentTab({ confirm = true }) },

    -- tabs: navigation
    { key = "[", mods = mod.SUPER, action = act.ActivateTabRelative(-1) },
    { key = "]", mods = mod.SUPER, action = act.ActivateTabRelative(1) },
    { key = "[", mods = mod.SUPER_REV, action = act.MoveTabRelative(-1) },
    { key = "]", mods = mod.SUPER_REV, action = act.MoveTabRelative(1) },

    -- window --
    -- spawn windows
    { key = "n", mods = mod.SUPER, action = act.SpawnWindow },

    -- panes --
    -- panes: split panes
    {
        key = [[/]],
        mods = mod.SUPER_REV,
        action = act.SplitVertical({ domain = "CurrentPaneDomain" }),
    },
    {
        key = [[\]],
        mods = mod.SUPER_REV,
        action = act.SplitHorizontal({ domain = "CurrentPaneDomain" }),
    },
    {
        key = [[-]],
        mods = mod.SUPER_REV,
        action = act.CloseCurrentPane({ confirm = true }),
    },

    -- panes: zoom+close pane
    { key = "z", mods = mod.SUPER_REV, action = act.TogglePaneZoomState },
    { key = "w", mods = mod.SUPER, action = act.CloseCurrentPane({ confirm = false }) },

    -- panes: navigation
    { key = "k", mods = mod.SUPER_REV, action = act.ActivatePaneDirection("Up") },
    { key = "j", mods = mod.SUPER_REV, action = act.ActivatePaneDirection("Down") },
    { key = "h", mods = mod.SUPER_REV, action = act.ActivatePaneDirection("Left") },
    { key = "l", mods = mod.SUPER_REV, action = act.ActivatePaneDirection("Right") },

    -- panes: resize
    { key = "UpArrow", mods = mod.SUPER_REV, action = act.AdjustPaneSize({ "Up", 1 }) },
    { key = "DownArrow", mods = mod.SUPER_REV, action = act.AdjustPaneSize({ "Down", 1 }) },
    { key = "LeftArrow", mods = mod.SUPER_REV, action = act.AdjustPaneSize({ "Left", 1 }) },
    { key = "RightArrow", mods = mod.SUPER_REV, action = act.AdjustPaneSize({ "Right", 1 }) },

    -- fonts --
    -- fonts: resize
    { key = "UpArrow", mods = mod.SUPER, action = act.IncreaseFontSize },
    { key = "DownArrow", mods = mod.SUPER, action = act.DecreaseFontSize },
    { key = "r", mods = mod.SUPER, action = act.ResetFontSize },

    -- key-tables --
    -- resizes fonts
    {
        key = "f",
        mods = "LEADER",
        action = act.ActivateKeyTable({
        name = "resize_font",
        one_shot = false,
        timemout_miliseconds = 1000,
        }),
    },
    -- resize panes
    {
        key = "p",
        mods = "LEADER",
        action = act.ActivateKeyTable({
        name = "resize_pane",
        one_shot = false,
        timemout_miliseconds = 1000,
        }),
    },
    -- rename tab bar
    {
        key = "R",
        mods = "CTRL|SHIFT",
        action = act.PromptInputLine({
        description = "Enter new name for tab",
        action = wezterm.action_callback(function(window, pane, line)
            -- line will be `nil` if they hit escape without entering anything
            -- An empty string if they just hit enter
            -- Or the actual line of text they wrote
            if line then
            window:active_tab():set_title(line)
            end
        end),
        }),
    },
}

local key_tables = {
    resize_font = {
        { key = "k", action = act.IncreaseFontSize },
        { key = "j", action = act.DecreaseFontSize },
        { key = "r", action = act.ResetFontSize },
        { key = "Escape", action = "PopKeyTable" },
        { key = "q", action = "PopKeyTable" },
    },
    resize_pane = {
        { key = "k", action = act.AdjustPaneSize({ "Up", 1 }) },
        { key = "j", action = act.AdjustPaneSize({ "Down", 1 }) },
        { key = "h", action = act.AdjustPaneSize({ "Left", 1 }) },
        { key = "l", action = act.AdjustPaneSize({ "Right", 1 }) },
        { key = "Escape", action = "PopKeyTable" },
        { key = "q", action = "PopKeyTable" },
    },
    }

    local mouse_bindings = {
    -- Ctrl-click will open the link under the mouse cursor
    {
        event = { Up = { streak = 1, button = "Left" } },
        mods = "CTRL",
        action = act.OpenLinkAtMouseCursor,
    },
    -- Move mouse will only select text and not copy text to clipboard
    {
        event = { Down = { streak = 1, button = "Left" } },
        mods = "NONE",
        action = act.SelectTextAtMouseCursor("Cell"),
    },
    {
        event = { Up = { streak = 1, button = "Left" } },
        mods = "NONE",
        action = act.ExtendSelectionToMouseCursor("Cell"),
    },
    {
        event = { Drag = { streak = 1, button = "Left" } },
        mods = "NONE",
        action = act.ExtendSelectionToMouseCursor("Cell"),
    },
    -- Triple Left click will select a line
    {
        event = { Down = { streak = 3, button = "Left" } },
        mods = "NONE",
        action = act.SelectTextAtMouseCursor("Line"),
    },
    {
        event = { Up = { streak = 3, button = "Left" } },
        mods = "NONE",
        action = act.SelectTextAtMouseCursor("Line"),
    },
    -- Double Left click will select a word
    {
        event = { Down = { streak = 2, button = "Left" } },
        mods = "NONE",
        action = act.SelectTextAtMouseCursor("Word"),
    },
    {
        event = { Up = { streak = 2, button = "Left" } },
        mods = "NONE",
        action = act.SelectTextAtMouseCursor("Word"),
    },
    -- Turn on the mouse wheel to scroll the screen
    {
        event = { Down = { streak = 1, button = { WheelUp = 1 } } },
        mods = "NONE",
        action = act.ScrollByCurrentEventWheelDelta,
    },
    {
        event = { Down = { streak = 1, button = { WheelDown = 1 } } },
        mods = "NONE",
        action = act.ScrollByCurrentEventWheelDelta,
    },
}

c.disable_default_key_bindings = true
c.disable_default_mouse_bindings = true
c.leader = { key = "Space", mods = "CTRL|SHIFT" }
c.keys = keys
c.key_tables = key_tables
c.mouse_bindings = mouse_bindings

return c
