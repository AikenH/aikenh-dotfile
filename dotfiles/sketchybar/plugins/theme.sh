#!/bin/bash
# theme.sh — toggle light/dark palette and re-apply colors to bar + items.

CONFIG_DIR="$HOME/.config/sketchybar"
COLORS="$CONFIG_DIR/colors.sh"
STATE_FILE="$HOME/.cache/sketchybar_theme"

# Read current theme from state (default dark)
CUR=$(cat "$STATE_FILE" 2>/dev/null || echo "dark")
if [ "$CUR" = "dark" ]; then
  NEW="light"
else
  NEW="dark"
fi

# Persist and apply
echo "$NEW" > "$STATE_FILE"
export CURRENT_THEME="$NEW"
# shellcheck disable=SC1090
source "$COLORS"

# Apply to bar
sketchybar --bar color=$BAR_COLOR

# Global defaults
sketchybar --default icon.color=$ICON_COLOR \
                     label.color=$LABEL_COLOR

# Space badges + app pills (spaces.sh reads current palette vars)
bash "$CONFIG_DIR/plugins/spaces.sh"

# Theme button icon (sun/moon) + color
if [ "$NEW" = "light" ]; then
  sketchybar --set theme_button icon="󰖙" icon.color=$ORANGE
else
  sketchybar --set theme_button icon="󰖔" icon.color=$YELLOW
fi

# Right-side brackets: pill body = bar color (matching theme bg),
# icon chip keeps its accent color, text = foreground.
sketchybar --set volume_icon background.color=$BLUE icon.color=$BAR_COLOR
sketchybar --set battery_icon background.color=$GREEN icon.color=$BAR_COLOR
sketchybar --set calendar_icon background.color=$MAGENTA icon.color=$BAR_COLOR
sketchybar --set volume_bracket background.border_color=$BLUE background.color=$BAR_COLOR
sketchybar --set battery_bracket background.border_color=$GREEN background.color=$BAR_COLOR
sketchybar --set calendar_bracket background.border_color=$MAGENTA background.color=$BAR_COLOR

# Text labels follow the theme foreground (fix: explicit per-item, not just default)
sketchybar --set volume_label label.color=$WHITE
sketchybar --set battery_label label.color=$WHITE
sketchybar --set calendar_label label.color=$WHITE
sketchybar --set front_app label.color=$WHITE

# Apple logo + its popup
sketchybar --set apple.logo icon.color=$GREEN popup.background.color=$POPUP_BACKGROUND_COLOR popup.background.border_color=$GREEN
# Apple popup items: text/icon follow theme foreground
sketchybar --set apple.prefs label.color=$WHITE icon.color=$GREEN
sketchybar --set apple.activity label.color=$WHITE icon.color=$GREEN
sketchybar --set apple.lock label.color=$WHITE icon.color=$GREEN

# Calendar popup background + month grid items follow theme
sketchybar --set calendar_bracket popup.background.color=$POPUP_BACKGROUND_COLOR popup.background.border_color=$MAGENTA
sketchybar --set calendar.header label.color=$MAGENTA
sketchybar --set calendar.week label.color=$BLUE
for r in 1 2 3 4 5 6; do
  sketchybar --set "calendar.row$r" label.color=$WHITE
done

# Zen + theme buttons: body matches bar, icons stay accent
sketchybar --set zen_button background.color=$BAR_COLOR icon.color=$MAGENTA
sketchybar --set theme_button background.color=$BAR_COLOR

echo "theme -> $NEW"
