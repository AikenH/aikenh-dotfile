#!/bin/bash
# zen.sh — toggle Zen mode: keep left nav (workspace badges + front app),
# hide right distractions (volume/battery/calendar/apple).

ZEN_STATE_FILE="$HOME/.cache/sketchybar_zen"
if [ -f "$ZEN_STATE_FILE" ] && [ "$(cat "$ZEN_STATE_FILE")" = "on" ]; then
  ZEN="off"
else
  ZEN="on"
fi
echo "$ZEN" > "$ZEN_STATE_FILE"

if [ "$ZEN" = "on" ]; then
  # Enter Zen: hide apple + right-side brackets, keep workspace nav + front_app
  sketchybar --set apple.logo drawing=off \
             --set volume_bracket drawing=off \
             --set battery_bracket drawing=off \
             --set calendar_bracket drawing=off \
             --set zen_button icon="󰝥" \
             --set theme_button icon="󰔎"
  # Keep left nav (workspace badges + front app) visible
  bash "$HOME/.config/sketchybar/plugins/spaces.sh"
  bash "$HOME/.config/sketchybar/plugins/front_app.sh"
else
  # Exit Zen: restore everything
  sketchybar --set apple.logo drawing=on \
             --set volume_bracket drawing=on \
             --set battery_bracket drawing=on \
             --set calendar_bracket drawing=on \
             --set zen_button icon="󰨇" \
             --set theme_button icon="󰖔"
  bash "$HOME/.config/sketchybar/plugins/spaces.sh"
  bash "$HOME/.config/sketchybar/plugins/front_app.sh"
fi
