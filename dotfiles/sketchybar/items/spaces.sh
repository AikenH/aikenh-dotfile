#!/bin/bash

# Items for Paneru Virtual Workspaces & Individual Clickable [Icon + App Name] Badges

for sid in {1..5}; do
  # Workspace number badge
  sketchybar --add item "space.$sid" left \
             --set "space.$sid" \
                 icon="$sid" \
                 icon.font="GoMono Nerd Font:Bold:12.0" \
                 icon.padding_left=6 \
                 icon.padding_right=6 \
                 label.drawing=off \
                 background.drawing=on \
                 background.color=$BACKGROUND_2 \
                 background.corner_radius=5 \
                 background.height=22 \
                 padding_left=3 \
                 padding_right=1 \
                 script="$PLUGIN_DIR/spaces.sh" \
                 click_script="paneru send-cmd window virtualnum $sid"

  # Up to 5 individual [Icon + App Name] badges per workspace
  for aid in {1..5}; do
    sketchybar --add item "space.$sid.app.$aid" left \
               --set "space.$sid.app.$aid" \
                   icon.font="GoMono Nerd Font:Regular:14.0" \
                   icon.padding_left=6 \
                   icon.padding_right=2 \
                   label.font="GoMono Nerd Font:Bold:12.0" \
                   label.padding_left=2 \
                   label.padding_right=6 \
                   label.drawing=on \
                   background.drawing=on \
                   background.color=$BACKGROUND_1 \
                   background.corner_radius=5 \
                   background.height=22 \
                   padding_left=1 \
                   padding_right=1 \
                   drawing=off
  done
done

sketchybar --add item separator left \
           --set separator \
               icon=$SEPARATOR \
               icon.font="GoMono Nerd Font:Bold:13.0" \
               icon.color=$GREEN \
               icon.padding_left=6 \
               icon.padding_right=6 \
               label.drawing=off \
               background.drawing=off \
               click_script="paneru send-cmd window virtual south"

# Initial update trigger
sketchybar --add event paneru_update
sketchybar --subscribe "space.1" paneru_update front_app_switched space_change
