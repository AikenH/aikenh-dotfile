local wezterm = require 'wezterm'

local function normalize_path(path)
  if not path or path == '' then
    return ''
  end

  local home = os.getenv('HOME') or ''
  if home ~= '' then
    path = path:gsub('^' .. home, '~')
  end

  return path
end

local function get_path_segments(path, num_segments)
  path = normalize_path(path)
  if path == '' then
    return ''
  end

  local segments = {}
  for segment in path:gmatch('[^/]+') do
    table.insert(segments, segment)
  end

  local start_idx = math.max(1, #segments - num_segments + 1)
  local result = {}
  for i = start_idx, #segments do
    table.insert(result, segments[i])
  end

  return table.concat(result, '/')
end

local function trim_tail(text, max_len)
  if max_len <= 0 then
    return ''
  end
  if #text <= max_len then
    return text
  end
  if max_len == 1 then
    return '…'
  end
  return '…' .. text:sub(-(max_len - 1))
end

local function fit_title_to_budget(title, budget)
  if budget <= 0 then
    return ''
  end
  if #title <= budget then
    return title
  end

  local segments = {}
  for segment in title:gmatch('[^/]+') do
    table.insert(segments, segment)
  end

  if #segments <= 1 then
    return trim_tail(title, budget)
  end

  local leaf = segments[#segments]
  if #leaf >= budget then
    return trim_tail(leaf, budget)
  end

  local result = { leaf }
  local used = #leaf
  local hidden_prefix = #segments > 1

  for i = #segments - 1, 1, -1 do
    local seg = segments[i]
    local extra = #seg + 1
    local ellipsis = i > 1 and 2 or 0
    if used + extra + ellipsis <= budget then
      table.insert(result, 1, seg)
      used = used + extra
      hidden_prefix = i > 1
    else
      break
    end
  end

  local collapsed = table.concat(result, '/')
  if hidden_prefix then
    if #collapsed + 2 <= budget then
      return '…/' .. collapsed
    end
    return trim_tail(leaf, budget)
  end

  return collapsed
end

local function split_title_parts(title)
  local slash_idx = title:match('^.*()/')
  if not slash_idx then
    return '', title
  end

  return title:sub(1, slash_idx), title:sub(slash_idx + 1)
end

wezterm.on('format-tab-title', function(tab, tabs, panes, effective_config, hover, max_width)
  local KAKU = {
    INDEX_BG = '#2d3d52',
    INDEX_FG = '#9db4cc',
    TITLE_ACTIVE_BG = '#2a2940',
    TITLE_INACTIVE_BG = '#17151d',
    TITLE_HOVER_BG = '#201e28',
    PREFIX_ACTIVE_FG = '#9ea6ba',
    PREFIX_INACTIVE_FG = '#646878',
    LEAF_ACTIVE_FG = '#f4f7fb',
    LEAF_INACTIVE_FG = '#aeb4c2',
    LEAF_HOVER_FG = '#cdd3e2',
  }

  local active = tab.is_active
  local hovered = hover and not active

  local title_bg = active and KAKU.TITLE_ACTIVE_BG
    or (hovered and KAKU.TITLE_HOVER_BG or KAKU.TITLE_INACTIVE_BG)
  local prefix_fg = active and KAKU.PREFIX_ACTIVE_FG
    or (hovered and '#7a8099' or KAKU.PREFIX_INACTIVE_FG)
  local leaf_fg = active and KAKU.LEAF_ACTIVE_FG
    or (hovered and KAKU.LEAF_HOVER_FG or KAKU.LEAF_INACTIVE_FG)

  local cwd = ''
  local pane = tab.active_pane
  if pane and pane.current_working_dir then
    local cwd_uri = pane.current_working_dir
    if type(cwd_uri) == 'string' then
      cwd = cwd_uri:gsub('file://[^/]*', '')
    elseif cwd_uri.file_path then
      cwd = cwd_uri.file_path
    end
  end

  local title = tab.tab_title
  if title == '' or not title then
    if cwd ~= '' then
      title = get_path_segments(cwd, 3)
    else
      title = pane and pane.title or '?'
    end
  end

  if pane and pane.is_zoomed then
    title = title .. ' [Z]'
  end

  local index_text = ' ' .. tostring(tab.tab_index + 1) .. ' '
  local title_budget = math.max(4, max_width - #index_text - 2)
  title = fit_title_to_budget(title, title_budget)
  local prefix, leaf = split_title_parts(title)

  local items = {
    { Background = { Color = KAKU.INDEX_BG } },
    { Foreground = { Color = KAKU.INDEX_FG } },
    { Attribute = { Intensity = 'Normal' } },
    { Text = index_text },
    { Background = { Color = title_bg } },
    { Text = ' ' },
  }

  if prefix ~= '' then
    table.insert(items, { Foreground = { Color = prefix_fg } })
    table.insert(items, { Attribute = { Intensity = 'Normal' } })
    table.insert(items, { Text = prefix })
  end

  table.insert(items, { Foreground = { Color = leaf_fg } })
  table.insert(items, { Attribute = { Intensity = active and 'Bold' or 'Normal' } })
  table.insert(items, { Text = leaf })
  table.insert(items, { Background = { Color = title_bg } })
  table.insert(items, { Text = ' ' })

  return items
end)

return {}
