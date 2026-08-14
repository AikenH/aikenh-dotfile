#!/bin/bash

POPUP_SHOW="sketchybar --set calendar_bracket popup.drawing=on"
POPUP_HIDE="sketchybar --set calendar_bracket popup.drawing=off"
POPUP_TOGGLE="sketchybar --set calendar_bracket popup.drawing=toggle"

sketchybar --add item calendar_label right \
           --set calendar_label \
               update_freq=10 \
               script="$PLUGIN_DIR/calendar.sh" \
               click_script="$POPUP_TOGGLE" \
               mouse.entered="$POPUP_SHOW" \
               mouse.exited="$POPUP_HIDE" \
               label.font="GoMono Nerd Font:Bold:13.0" \
               label.color=$WHITE \
               label.padding_left=4 \
               label.padding_right=10 \
               background.drawing=off \
               padding_left=0 \
               padding_right=0

sketchybar --subscribe calendar_label system_woke mouse.entered mouse.exited

sketchybar --add item calendar_icon right \
           --set calendar_icon \
               icon="󰥔" \
               icon.font="GoMono Nerd Font:Bold:15.0" \
               icon.color=$BAR_SOLID \
               icon.padding_left=6 \
               icon.padding_right=6 \
               label.drawing=off \
               background.drawing=on \
               background.color=$MAGENTA \
               background.height=26 \
               background.corner_radius=6 \
               padding_left=0 \
               padding_right=0 \
               click_script="$POPUP_TOGGLE" \
               mouse.entered="$POPUP_SHOW" \
               mouse.exited="$POPUP_HIDE"

sketchybar --subscribe calendar_icon mouse.entered mouse.exited

sketchybar --add bracket calendar_bracket calendar_icon calendar_label \
           --set calendar_bracket \
               background.drawing=on \
               background.color=$BAR_SOLID \
               background.height=26 \
               background.corner_radius=8 \
               background.border_width=2 \
               background.border_color=$MAGENTA \
               padding_left=2 \
               padding_right=2 \
               click_script="$POPUP_TOGGLE" \
               mouse.entered="$POPUP_SHOW" \
               mouse.exited="$POPUP_HIDE" \
               popup.background.color=$POPUP_BACKGROUND_COLOR \
               popup.background.border_color=$MAGENTA \
               popup.background.border_width=2 \
               popup.background.corner_radius=10 \
               popup.align=right \
               popup.y_offset=6

sketchybar --subscribe calendar_bracket mouse.entered mouse.exited

# Popup Month Grid items (exact 13pt font for 100% alignment)
sketchybar --add item calendar.header popup.calendar_bracket \
           --set calendar.header \
               label.font="GoMono Nerd Font:Bold:13.0" \
               label.color=$MAGENTA \
               label.padding_left=16 \
               label.padding_right=16 \
               background.drawing=off \
               mouse.entered="$POPUP_SHOW" \
               mouse.exited="$POPUP_HIDE"

sketchybar --subscribe calendar.header mouse.entered mouse.exited

sketchybar --add item calendar.week popup.calendar_bracket \
           --set calendar.week \
               label.font="GoMono Nerd Font:Bold:13.0" \
               label.color=$BLUE \
               label.padding_left=16 \
               label.padding_right=16 \
               background.drawing=off \
               mouse.entered="$POPUP_SHOW" \
               mouse.exited="$POPUP_HIDE"

sketchybar --subscribe calendar.week mouse.entered mouse.exited

for r in {1..6}; do
  sketchybar --add item "calendar.row$r" popup.calendar_bracket \
             --set "calendar.row$r" \
                 label.font="GoMono Nerd Font:Bold:13.0" \
                 label.color=$WHITE \
                 label.padding_left=16 \
                 label.padding_right=16 \
                 background.drawing=off \
                 mouse.entered="$POPUP_SHOW" \
                 mouse.exited="$POPUP_HIDE"

  sketchybar --subscribe "calendar.row$r" mouse.entered mouse.exited
done
