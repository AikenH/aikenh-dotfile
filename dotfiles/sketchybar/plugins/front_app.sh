#!/bin/bash

ACTIVE_JSON=$(paneru query active --json 2>/dev/null)

if [ -n "$ACTIVE_JSON" ]; then
  APP_NAME=$(echo "$ACTIVE_JSON" | jq -r '.focused_app_name // ""' 2>/dev/null)
  WIN_TITLE=$(echo "$ACTIVE_JSON" | jq -r '.focused_window_title // ""' 2>/dev/null)

  if [ "$APP_NAME" = "null" ]; then APP_NAME=""; fi
  if [ "$WIN_TITLE" = "null" ]; then WIN_TITLE=""; fi

  if [ -n "$WIN_TITLE" ] && [ "$WIN_TITLE" != "$APP_NAME" ]; then
    if [ ${#WIN_TITLE} -gt 45 ]; then
      WIN_TITLE="${WIN_TITLE:0:42}..."
    fi
    LABEL="$APP_NAME — $WIN_TITLE"
  elif [ -n "$APP_NAME" ]; then
    LABEL="$APP_NAME"
  else
    LABEL="Desktop"
  fi

  sketchybar --set "$NAME" label="$LABEL" drawing=on
fi
