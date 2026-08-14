#!/bin/bash

sketchybar --add item volume_label right \
           --set volume_label \
               script="$PLUGIN_DIR/volume.sh" \
               label.font="GoMono Nerd Font:Bold:13.0" \
               label.color=$WHITE \
               label.padding_left=4 \
               label.padding_right=10 \
               background.drawing=off \
               padding_left=0 \
               padding_right=0

sketchybar --subscribe volume_label volume_change

sketchybar --add item volume_icon right \
           --set volume_icon \
               icon="󰕾" \
               icon.font="GoMono Nerd Font:Bold:15.0" \
               icon.color=$BAR_SOLID \
               icon.padding_left=6 \
               icon.padding_right=6 \
               label.drawing=off \
               background.drawing=on \
               background.color=$BLUE \
               background.height=26 \
               background.corner_radius=6 \
               padding_left=0 \
               padding_right=0

sketchybar --add bracket volume_bracket volume_icon volume_label \
           --set volume_bracket \
               background.drawing=on \
               background.color=$BAR_SOLID \
               background.height=26 \
               background.corner_radius=8 \
               background.border_width=2 \
               background.border_color=$BLUE \
               padding_left=2 \
               padding_right=2
