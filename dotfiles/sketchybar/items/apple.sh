#!/bin/bash

POPUP_OFF="sketchybar --set apple.logo popup.drawing=off"
POPUP_CLICK_SCRIPT="sketchybar --set \$NAME popup.drawing=toggle"

apple_logo=(
  icon=$APPLE
  icon.font="GoMono Nerd Font:Bold:16.0"
  icon.color=$GREEN
  icon.padding_left=8
  icon.padding_right=8
  label.drawing=off
  click_script="$POPUP_CLICK_SCRIPT"
  popup.background.color=$POPUP_BACKGROUND_COLOR
  popup.background.border_color=$GREEN
  popup.background.border_width=2
  popup.background.corner_radius=10
  popup.blur_radius=0
  popup.y_offset=6
)

apple_prefs=(
  icon=$PREFERENCES
  icon.font="GoMono Nerd Font:Regular:14.0"
  icon.padding_left=12
  icon.padding_right=8
  label="System Settings"
  label.font="GoMono Nerd Font:Regular:13.0"
  label.padding_left=0
  label.padding_right=14
  click_script="open -a 'System Settings'; $POPUP_OFF"
)

apple_activity=(
  icon=$ACTIVITY
  icon.font="GoMono Nerd Font:Regular:14.0"
  icon.padding_left=12
  icon.padding_right=8
  label="Activity Monitor"
  label.font="GoMono Nerd Font:Regular:13.0"
  label.padding_left=0
  label.padding_right=14
  click_script="open -a 'Activity Monitor'; $POPUP_OFF"
)

apple_lock=(
  icon=$LOCK
  icon.font="GoMono Nerd Font:Regular:14.0"
  icon.padding_left=12
  icon.padding_right=8
  label="Lock Screen"
  label.font="GoMono Nerd Font:Regular:13.0"
  label.padding_left=0
  label.padding_right=14
  click_script="pmset displaysleepnow; $POPUP_OFF"
)

sketchybar --add item apple.logo left                  \
           --set apple.logo "${apple_logo[@]}"         \
                                                       \
           --add item apple.prefs popup.apple.logo     \
           --set apple.prefs "${apple_prefs[@]}"       \
                                                       \
           --add item apple.activity popup.apple.logo  \
           --set apple.activity "${apple_activity[@]}" \
                                                       \
           --add item apple.lock popup.apple.logo      \
           --set apple.lock "${apple_lock[@]}"
