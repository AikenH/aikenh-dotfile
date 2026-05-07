# bat shell integration
# bat Config Import Flag
if command -v bat >/dev/null 2>&1; then
  export BAT_THEME="Catppuccin Mocha"
  alias cat='bat --paging=never'

  # If fzf is also installed, enhance its preview (bat handles this side)
  if command -v fzf >/dev/null 2>&1; then
    export FZF_CTRL_T_OPTS="--preview 'bat --color=always --line-range :200 {}'"
  fi
fi
