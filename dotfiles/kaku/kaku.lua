local wezterm = require 'wezterm'

local function resolve_bundled_config()
  local resource_dir = wezterm.executable_dir:gsub('MacOS/?$', 'Resources')
  local bundled = resource_dir .. '/kaku.lua'
  local f = io.open(bundled, 'r')
  if f then
    f:close()
    return bundled
  end

  local dev_bundled = wezterm.executable_dir .. '/../../assets/macos/Kaku.app/Contents/Resources/kaku.lua'
  f = io.open(dev_bundled, 'r')
  if f then
    f:close()
    return dev_bundled
  end

  local app_bundled = '/Applications/Kaku.app/Contents/Resources/kaku.lua'
  f = io.open(app_bundled, 'r')
  if f then
    f:close()
    return app_bundled
  end

  local home = os.getenv('HOME') or ''
  local home_bundled = home .. '/Applications/Kaku.app/Contents/Resources/kaku.lua'
  f = io.open(home_bundled, 'r')
  if f then
    f:close()
    return home_bundled
  end

  return nil
end

local config = {}
local bundled = resolve_bundled_config()

if bundled then
  local ok, loaded = pcall(dofile, bundled)
  if ok and type(loaded) == 'table' then
    config = loaded
  else
    wezterm.log_error('Kaku: failed to load bundled defaults from ' .. bundled)
  end
else
  wezterm.log_error('Kaku: bundled defaults not found')
end

-- User overrides:
-- Kaku intentionally keeps WezTerm-compatible Lua API names
-- for maximum compatibility, so `wezterm.*` here is expected.
-- Full API docs: https://wezfurlong.org/wezterm/config/lua/
--
-- 1) Font family and size
-- config.font = wezterm.font('JetBrains Mono')
-- config.font_size = 16.0
-- config.line_height = 1.2
--
-- 2) Color scheme
-- config.color_scheme = 'Catppuccin Mocha'
--
-- 3) Window size and padding
-- config.initial_cols = 120
-- config.initial_rows = 30
-- config.window_padding = { left = '24px', right = '24px', top = '40px', bottom = '20px' }
--
-- 4) Window transparency and blur
-- config.window_background_opacity = 0.95
-- config.macos_window_background_blur = 20
--
-- 5) Copy on select
-- config.copy_on_select = false
--
-- 6) Default shell/program
-- config.default_prog = { '/bin/zsh', '-l' }
--
-- 7) Cursor and scrollback
-- config.default_cursor_style = 'BlinkingBar'
-- config.cursor_blink_rate = 500
-- config.scrollback_lines = 20000
--
-- 8) Tab bar
-- config.hide_tab_bar_if_only_one_tab = true
-- config.tab_bar_at_bottom = true
-- config.tab_title_show_basename_only = true
--
-- 9) Working directory inheritance
-- config.window_inherit_working_directory = true
-- config.tab_inherit_working_directory = true
-- config.split_pane_inherit_working_directory = true
--
-- 10) Split pane
-- config.split_pane_gap = 2
-- config.inactive_pane_hsb = { saturation = 1.0, brightness = 0.9 }
--
-- 11) Add or override a key binding
-- table.insert(config.keys, {
--   key = 'Enter',
--   mods = 'CMD|SHIFT',
--   action = wezterm.action.TogglePaneZoomState,
-- })

-- ═══════════════════════════════════════════════════════════════
-- Font: English primary + Chinese fallback (mirrors Ghostty setup)
-- ═══════════════════════════════════════════════════════════════
config.font = wezterm.font_with_fallback({
    { family = "Inconsolata LGC Nerd Font Mono", weight = "Regular" },
    "TsangerJinKai03",       -- Chinese fallback
})
config.font_size = 14
config.line_height = 1.1
-- Ligatures: calt(contextual alternates) + liga(standard) + dlig(discretionary)
config.harfbuzz_features = { "calt=1", "liga=1", "dlig=1" }

-- ═══════════════════════════════════════════════════════════════
-- Window: opacity, blur, padding
-- ═══════════════════════════════════════════════════════════════
config.window_background_opacity = 0.7
config.macos_window_background_blur = 20
config.window_padding = {
    left = 10, right = 10, top = 40, bottom = 0,
}
config.window_close_confirmation = 'NeverPrompt'

-- ═══════════════════════════════════════════════════════════════
-- Tab bar
-- ═══════════════════════════════════════════════════════════════
config.tab_bar_at_bottom = true
config.hide_tab_bar_if_only_one_tab = false
config.use_fancy_tab_bar = false
config.tab_max_width = 25
config.show_tab_index_in_tab_bar = true

-- ═══════════════════════════════════════════════════════════════
-- Keybinding fixes: passthrough Alt+keys to Zellij
-- Kaku/WezTerm on macOS intercepts some Alt+key combos.
-- Explicitly send them as raw ESC sequences for Zellij.
-- ═══════════════════════════════════════════════════════════════
config.keys = config.keys or {}
local passthrough_alt_keys = { 'n', 'h', 'j', 'k', 'l', 'f', 'i', 'o', 'p', '[', ']', '+', '-', '=' }
for _, k in ipairs(passthrough_alt_keys) do
    table.insert(config.keys, {
        key = k,
        mods = 'ALT',
        action = wezterm.action.SendKey({ key = k, mods = 'ALT' }),
    })
end

-- ═══════════════════════════════════════════════════════════════
-- Cursor & animation
-- ═══════════════════════════════════════════════════════════════
config.default_cursor_style = 'BlinkingBlock'
config.cursor_blink_rate = 500
config.cursor_blink_ease_in = "EaseIn"
config.cursor_blink_ease_out = "EaseOut"
config.animation_fps = 60
config.max_fps = 60

-- ═══════════════════════════════════════════════════════════════
-- Colors (dark theme, auto switch)
-- ═══════════════════════════════════════════════════════════════
config.color_scheme = 'Fahrenheit'

-- ═══════════════════════════════════════════════════════════════
-- Misc
-- ═══════════════════════════════════════════════════════════════
config.audible_bell = 'Disabled'
config.bell_dock_badge = true
config.scrollback_lines = 10000

-- inactive pane dim
config.inactive_pane_hsb = {
    saturation = 0.9,
    brightness = 0.6,
}

-- clipboard (mirrors Ghostty clipboard-read/write=allow, copy-on-select)

-- ═══════════════════════════════════════════════════════════════
-- URL auto-detect & underline
-- ═══════════════════════════════════════════════════════════════
config.hyperlink_rules = wezterm.default_hyperlink_rules()

-- ═══════════════════════════════════════════════════════════════
-- Tab Bar: bar.wezterm plugin
-- https://github.com/adriankarlen/bar.wezterm
-- NOTE: must be called AFTER setting color_scheme
-- ═══════════════════════════════════════════════════════════════
local bar = wezterm.plugin.require("https://github.com/adriankarlen/bar.wezterm")
bar.apply_to_config(config, {
    position = "bottom",
    max_width = 36,
    modules = {
        tabs = {
            active_tab_fg = 4,
            inactive_tab_fg = 6,
        },
        workspace = { enabled = true, color = 8 },
        leader = { enabled = true, color = 2 },
        zoom = { enabled = true, color = 4 },
        pane = { enabled = true, color = 7 },
        username = { enabled = false },
        hostname = { enabled = false },
        clock = { enabled = false },
        cwd = { enabled = false },  -- we handle cwd on the right side
        spotify = { enabled = false },
    },
})

-- Override Kaku's bundled update-right-status which clears the right status.
-- Renders: mode | battery | cwd | clock
wezterm.on('update-right-status', function(window, pane)
    local nf = wezterm.nerdfonts
    local cells = {}
    local sep = "  "

    -- CWD
    local cwd_uri = pane:get_current_working_dir()
    local cwd = ""
    if cwd_uri then
        cwd = (cwd_uri.file_path or tostring(cwd_uri)):gsub("^/Users/aikenhong", "~")
    end

    -- Battery
    local battery_text = ""
    for _, b in ipairs(wezterm.battery_info()) do
        local icon
        local charge = b.state_of_charge * 100
        if b.state == "Charging" then
            icon = nf.md_battery_charging
        elseif charge > 75 then
            icon = nf.md_battery
        elseif charge > 50 then
            icon = nf.md_battery_70
        elseif charge > 25 then
            icon = nf.md_battery_40
        else
            icon = nf.md_battery_10
        end
        battery_text = icon .. " " .. string.format("%.0f%%", charge)
    end

    -- Mode indicator
    local process = pane:get_foreground_process_name() or ""
    local basename = process:gsub("(.*[/\\])(.*)", "%2")
    local mode = ""
    if basename == "vim" or basename == "nvim" then
        mode = " VIM"
    elseif basename == "zellij" then
        mode = "󰕮 ZELLIJ"
    end

    -- Build right status
    if mode ~= "" then
        table.insert(cells, { Foreground = { Color = "#f38ba8" } })
        table.insert(cells, { Text = mode .. sep })
    end

    if battery_text ~= "" then
        table.insert(cells, { Foreground = { Color = "#f9e2af" } })
        table.insert(cells, { Text = battery_text .. sep })
    end

    table.insert(cells, { Foreground = { Color = "#89b4fa" } })
    table.insert(cells, { Text = nf.oct_file_directory .. " " .. cwd .. sep })

    table.insert(cells, { Foreground = { Color = "#cba6f7" } })
    table.insert(cells, { Text = nf.md_calendar_clock .. " " .. wezterm.strftime("%H:%M") .. " " })

    window:set_right_status(wezterm.format(cells))
end)

return config
