#!/bin/bash
# zen_theme.sh — Theme (light/dark) toggle button + Zen toggle button.
# Placed on the right side, LEFT of the calendar (time), RIGHT of battery.
# sketchybar right: first added = rightmost. Add theme first, then zen,
# so (left->right) order is: zen | theme.

# ── Theme (light/dark) toggle button ──
sketchybar --add item theme_button right \
           --set theme_button \
               icon="󰖔" \
               icon.font="GoMono Nerd Font:Bold:14.0" \
               icon.color=$YELLOW \
               icon.padding_left=6 \
               icon.padding_right=6 \
               label.drawing=off \
               background.drawing=on \
               background.color=$BAR_COLOR \
               background.corner_radius=6 \
               background.height=26 \
               padding_left=2 \
               padding_right=2 \
               click_script="$PLUGIN_DIR/theme.sh"

# ── Zen toggle button ──
sketchybar --add item zen_button right \
           --set zen_button \
               icon="󰨇" \
               icon.font="GoMono Nerd Font:Bold:14.0" \
               icon.color=$MAGENTA \
               icon.padding_left=6 \
               icon.padding_right=6 \
               label.drawing=off \
               background.drawing=on \
               background.color=$BAR_COLOR \
               background.corner_radius=6 \
               background.height=26 \
               padding_left=2 \
               padding_right=2 \
               click_script="$PLUGIN_DIR/zen.sh"
