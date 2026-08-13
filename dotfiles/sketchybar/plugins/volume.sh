#!/bin/bash

VOLUME=$(osascript -e "output volume of (get volume settings)" 2>/dev/null)
if [ -z "$VOLUME" ] || [ "$VOLUME" = "null" ]; then
  VOLUME="0"
fi

if [ "$VOLUME" -eq 0 ]; then
  ICON="󰝟"
elif [ "$VOLUME" -lt 33 ]; then
  ICON="󰕿"
elif [ "$VOLUME" -lt 66 ]; then
  ICON="󰖀"
else
  ICON="󰕾"
fi

sketchybar --set volume_icon icon="$ICON"
sketchybar --set volume_label label="${VOLUME}%"
