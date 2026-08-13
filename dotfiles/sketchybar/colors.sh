#!/bin/bash

# ── Catppuccin Mocha (Dark) — frosted glass variants ──
dark_black=0xff181926
dark_white=0xffcad3f5
dark_red=0xfff38ba8
dark_green=0xffa6e3a1
dark_blue=0xff89b4fa
dark_yellow=0xfff9e2af
dark_orange=0xfffab387
dark_magenta=0xffcba6f7
dark_grey=0xff6c7086
dark_bar=0xcc1e1e2e          # 80% (frosted)
dark_bg1=0xb3313244          # 70%
dark_bg2=0xb345475a          # 70%
dark_popup_bg=0xe61e1e2e     # 90%

# ── Vintage Paper (Light) — from AFA vintage theme, glass/frosted variants ──
# 语义：WHITE=前景文字(深棕)，BLACK=高亮槽文字(深棕)，BG1/BG2=浅色槽底
# 玻璃效果：bar/popup 半透明 + blur（见 sketchybarrc blur_radius）
light_black=0xff493f35        # foreground 深棕墨（高亮文字）
light_white=0xff493f35        # foreground 深棕墨（普通文字）
light_red=0xffcb624d          # chart-1 red
light_green=0xff799a4a        # chart-2 green
light_blue=0xff4a6c9a         # chart-3 blue
light_yellow=0xffd6aa41       # chart-4 yellow
light_orange=0xffcc6632       # chart-10 orange
light_magenta=0xff904cb2      # chart-6 magenta
light_grey=0xff7d6b56         # muted-foreground
light_bar=0xccf4f0e5          # background @80% (frosted)
light_bg1=0x99ebe4d7          # muted @60%
light_bg2=0xb3e1d7c2          # secondary @70%
light_popup_bg=0xe6fffbf4     # card @90%

# ── 当前主题状态（默认 dark）──
export CURRENT_THEME="${CURRENT_THEME:-dark}"

apply_palette() {
  local t="$1"
  if [ "$t" = "light" ]; then
    export BLACK=$light_black WHITE=$light_white RED=$light_red GREEN=$light_green
    export BLUE=$light_blue YELLOW=$light_yellow ORANGE=$light_orange
    export MAGENTA=$light_magenta GREY=$light_grey
    export BAR_COLOR=$light_bar BACKGROUND_1=$light_bg1 BACKGROUND_2=$light_bg2
    export POPUP_BACKGROUND_COLOR=$light_popup_bg
  else
    export BLACK=$dark_black WHITE=$dark_white RED=$dark_red GREEN=$dark_green
    export BLUE=$dark_blue YELLOW=$dark_yellow ORANGE=$dark_orange
    export MAGENTA=$dark_magenta GREY=$dark_grey
    export BAR_COLOR=$dark_bar BACKGROUND_1=$dark_bg1 BACKGROUND_2=$dark_bg2
    export POPUP_BACKGROUND_COLOR=$dark_popup_bg
  fi
  export ICON_COLOR=$WHITE
  export LABEL_COLOR=$WHITE
  export POPUP_BORDER_COLOR=$WHITE
  export SHADOW_COLOR=$BLACK
  export TRANSPARENT=0x00000000
}

apply_palette "$CURRENT_THEME"
