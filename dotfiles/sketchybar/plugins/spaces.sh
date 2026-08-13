#!/bin/bash
# spaces.sh — Scheme 1: show ONLY the currently focused screen's virtual workspaces.
# No cross-screen aggregation, no warp, no screen-concept mixing.
# Clicking is always same-screen: virtualnum N && focus M.

# Source current palette (respects persisted theme)
export CURRENT_THEME=$(cat "$HOME/.cache/sketchybar_theme" 2>/dev/null || echo "dark")
# shellcheck disable=SC1090
source "$HOME/.config/sketchybar/colors.sh"

ACTIVE_JSON=$(paneru query active --json 2>/dev/null)
WS_JSON=$(paneru query virtual-workspaces --json 2>/dev/null)

python3 - "$ACTIVE_JSON" "$WS_JSON" "$GREEN" "$BLACK" "$BACKGROUND_2" "$WHITE" "$BACKGROUND_1" << 'PY'
import json, sys, subprocess

active_raw = sys.argv[1] if len(sys.argv) > 1 else ""
ws_raw = sys.argv[2] if len(sys.argv) > 2 else ""
GREEN = sys.argv[3] if len(sys.argv) > 3 else "0xffa6e3a1"
BLACK = sys.argv[4] if len(sys.argv) > 4 else "0xff181926"
BG2   = sys.argv[5] if len(sys.argv) > 5 else "0xff45475a"
WHITE = sys.argv[6] if len(sys.argv) > 6 else "0xffcad3f5"
BG1   = sys.argv[7] if len(sys.argv) > 7 else "0xff313244"

APP_ICONS = {
    'Zed': '󰨞',
    'Ghostty': '󰆍',
    'WeTERM': '󰆍',
    'Terminal': '󰆍',
    'iTerm2': '󰆍',
    'Kaku': '󰆍',
    'Code': '󰨞',
    'VSCode': '󰨞',
    'Zen': '󰈹',
    'Safari': '󰈹',
    'Google Chrome': '󰈹',
    'Microsoft Edge': '󰖟',
    '微信': '󰭹',
    '企业微信': '󰵅',
    'Finder': '󰀵',
    'System Settings': '󰒓',
}

try:
    active_data = json.loads(active_raw) if active_raw else {}
    ws_data = json.loads(ws_raw) if ws_raw else []
except Exception:
    sys.exit(0)

active_native = active_data.get('native_workspace_id', 0)
active_num = active_data.get('virtual_workspace_number', 1)

# Only virtual workspaces belonging to the focused screen's native workspace.
native_vws = [vw for vw in ws_data if vw.get('native_workspace_id') == active_native]
vws_by_num = {vw.get('number', 1): vw for vw in native_vws}

MAX_SLOTS = 5
MAX_APPS = 5
cmds = ['--animate', 'tanh', '15']

for s_idx in range(1, MAX_SLOTS + 1):
    slot_name = f'space.{s_idx}'
    if s_idx in vws_by_num:
        vw = vws_by_num[s_idx]
        is_active = (s_idx == active_num)

        # Same-screen workspace badge click: direct virtualnum.
        click = f"paneru send-cmd window virtualnum {s_idx}"

        if is_active:
            cmds.extend([
                '--set', slot_name, 'drawing=on',
                'icon=' + str(s_idx),
                f'background.color={GREEN}',
                f'icon.color={BLACK}',
                'click_script=' + click,
            ])
        else:
            cmds.extend([
                '--set', slot_name, 'drawing=on',
                'icon=' + str(s_idx),
                f'background.color={BG2}',
                f'icon.color={WHITE}',
                'click_script=' + click,
            ])

        windows = vw.get('windows', [])
        for a_idx in range(1, MAX_APPS + 1):
            item_name = f'space.{s_idx}.app.{a_idx}'
            if a_idx <= len(windows):
                w = windows[a_idx - 1]
                app_name = w.get('app_name', '')
                is_app_focused = w.get('focused', False)
                icon = APP_ICONS.get(app_name, '󰈔')

                app_click = f"paneru send-cmd window virtualnum {s_idx} && paneru send-cmd window focus {a_idx}"

                if is_app_focused:
                    cmds.extend([
                        '--set', item_name, 'drawing=on',
                        'icon=' + icon,
                        'label=' + app_name,
                        'click_script=' + app_click,
                        f'background.color={GREEN}',
                        f'icon.color={BLACK}',
                        f'label.color={BLACK}',
                    ])
                else:
                    cmds.extend([
                        '--set', item_name, 'drawing=on',
                        'icon=' + icon,
                        'label=' + app_name,
                        'click_script=' + app_click,
                        f'background.color={BG1}',
                        f'icon.color={WHITE}',
                        f'label.color={WHITE}',
                    ])
            else:
                cmds.extend(['--set', item_name, 'drawing=off'])
    else:
        cmds.extend(['--set', slot_name, 'drawing=off'])
        for a_idx in range(1, MAX_APPS + 1):
            cmds.extend(['--set', f'space.{s_idx}.app.{a_idx}', 'drawing=off'])

if cmds:
    subprocess.run(['sketchybar'] + cmds)
PY
