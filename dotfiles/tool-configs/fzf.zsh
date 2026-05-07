# fzf shell integration
# fzf Config Import Flag
if command -v fzf >/dev/null 2>&1; then
  # Load fzf key bindings and completion (Ctrl-T, Ctrl-R, Alt-C)
  # Priority 1: fzf >= 0.48 supports --zsh flag directly (most reliable)
  if fzf --zsh >/dev/null 2>&1; then
    source <(fzf --zsh)
  # Priority 2: common system package locations
  elif [[ -f /usr/share/fzf/shell/key-bindings.zsh ]]; then
    source /usr/share/fzf/shell/key-bindings.zsh
    [[ -f /usr/share/fzf/shell/completion.zsh ]] && source /usr/share/fzf/shell/completion.zsh
  elif [[ -f /usr/share/fzf/key-bindings.zsh ]]; then
    source /usr/share/fzf/key-bindings.zsh
  elif [[ -f /usr/share/doc/fzf/examples/key-bindings.zsh ]]; then
    source /usr/share/doc/fzf/examples/key-bindings.zsh
  elif command -v brew >/dev/null 2>&1; then
    _fzf_prefix="$(brew --prefix fzf 2>/dev/null)"
    [[ -f "${_fzf_prefix}/shell/key-bindings.zsh" ]] && source "${_fzf_prefix}/shell/key-bindings.zsh"
    [[ -f "${_fzf_prefix}/shell/completion.zsh" ]]   && source "${_fzf_prefix}/shell/completion.zsh"
    unset _fzf_prefix
  fi

  # Use fd as default find command if available
  if command -v fd >/dev/null 2>&1; then
    export FZF_DEFAULT_COMMAND='fd --type f --hidden --follow --exclude .git'
    export FZF_CTRL_T_COMMAND="$FZF_DEFAULT_COMMAND"
    export FZF_ALT_C_COMMAND='fd --type d --hidden --follow --exclude .git'
  fi

  # Use bat for preview if available
  if command -v bat >/dev/null 2>&1; then
    export FZF_CTRL_T_OPTS="--preview 'bat --color=always --line-range :200 {}'"
  fi
fi
