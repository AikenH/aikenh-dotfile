#!/bin/bash

BATT_INFO=$(pmset -g batt 2>/dev/null)
PERCENT=$(echo "$BATT_INFO" | grep -o "[0-9]*%" | head -1 | tr -d '%')
CHARGING=$(echo "$BATT_INFO" | grep -c "AC Power")

if [ -z "$PERCENT" ] || [ "$PERCENT" = "null" ]; then
  PERCENT="100"
fi

if [ "$CHARGING" -gt 0 ]; then
  ICON="󰂆"
elif [ "$PERCENT" -gt 85 ]; then
  ICON="󰁹"
elif [ "$PERCENT" -gt 60 ]; then
  ICON="󰂁"
elif [ "$PERCENT" -gt 35 ]; then
  ICON="󰁽"
elif [ "$PERCENT" -gt 15 ]; then
  ICON="󰁻"
else
  ICON="󰂎"
fi

sketchybar --set battery_icon icon="$ICON"
sketchybar --set battery_label label="${PERCENT}%"
