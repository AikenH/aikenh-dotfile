# ─── FZF + fd + bat integration ─────────────────────────────────────────
# Requires: fzf, fd, bat
# Append guard: # FZF Integration Flag

# Resolve bat/batcat and fd/fdfind naming differences (Ubuntu vs Homebrew)
if command -v bat >/dev/null 2>&1; then
  alias b='bat'
elif command -v batcat >/dev/null 2>&1; then
  alias bat='batcat'
  alias b='batcat'
fi

if command -v fdfind >/dev/null 2>&1 && ! command -v fd >/dev/null 2>&1; then
  alias fd='fdfind'
fi

if command -v fzf >/dev/null 2>&1; then
  eval "$(fzf --zsh 2>/dev/null || fzf --bash 2>/dev/null || true)"
  export FZF_DEFAULT_COMMAND='fd --type f --hidden --follow --exclude .git 2>/dev/null || fdfind --type f --hidden --follow --exclude .git'
  export FZF_CTRL_T_COMMAND="$FZF_DEFAULT_COMMAND"
  export FZF_ALT_C_COMMAND='fd --type d --hidden --follow --exclude .git 2>/dev/null || fdfind --type d --hidden --follow --exclude .git'
  export FZF_DEFAULT_OPTS=$'--height 40%\n--layout=reverse\n--border\n--preview "if [ -d \"{}\" ]; then fd --max-results 100 --type f . \"{}\" 2>/dev/null || fdfind --max-results 100 --type f . \"{}\"; else bat --color=always --style=numbers --line-range=:200 -- \"{}\" 2>/dev/null || batcat --color=always --style=numbers --line-range=:200 -- \"{}\"; fi"'
fi

alias fdf='fd -t f'
alias fdd='fd -t d'
