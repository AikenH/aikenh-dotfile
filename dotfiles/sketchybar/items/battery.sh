#!/bin/bash

sketchybar --add item battery_label right \
           --set battery_label \
               update_freq=60 \
               script="$PLUGIN_DIR/battery.sh" \
               label.font="GoMono Nerd Font:Bold:13.0" \
               label.color=$WHITE \
               label.padding_left=4 \
               label.padding_right=10 \
               background.drawing=off \
               padding_left=0 \
               padding_right=0

sketchybar --subscribe battery_label power_source_change system_woke

sketchybar --add item battery_icon right \
           --set battery_icon \
               icon="󰁹" \
               icon.font="GoMono Nerd Font:Bold:15.0" \
               icon.color=$BAR_SOLID \
               icon.padding_left=6 \
               icon.padding_right=6 \
               label.drawing=off \
               background.drawing=on \
               background.color=$GREEN \
               background.height=26 \
               background.corner_radius=6 \
               padding_left=0 \
               padding_right=0

sketchybar --add bracket battery_bracket battery_icon battery_label \
           --set battery_bracket \
               background.drawing=on \
               background.color=$BAR_SOLID \
               background.height=26 \
               background.corner_radius=8 \
               background.border_width=2 \
               background.border_color=$GREEN \
               padding_left=2 \
               padding_right=2
