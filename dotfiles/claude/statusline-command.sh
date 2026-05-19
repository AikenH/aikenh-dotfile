#!/bin/bash

input=$(cat)

if ! command -v jq >/dev/null 2>&1; then
    printf '%s\n' 'statusline error: jq not found'
    printf '%s\n' 'install jq or update the parser'
    exit 0
fi

json_get() {
    printf '%s' "$input" | jq -r "$1 // empty" 2>/dev/null
}

MODEL=$(json_get '.model.display_name')
[ -z "$MODEL" ] && MODEL=$(json_get '.model.id')
DIR=$(json_get '.workspace.current_dir')
[ -z "$DIR" ] && DIR=$(json_get '.cwd')
DIR_NAME=${DIR##*/}
[ -z "$DIR_NAME" ] && DIR_NAME="$DIR"

GREEN='\033[32m'
YELLOW='\033[33m'
RESET='\033[0m'

LINE1="[$MODEL] 📁 ${DIR_NAME}"

if [ -n "$DIR" ] && git --no-optional-locks -C "$DIR" rev-parse --git-dir >/dev/null 2>&1; then
    BRANCH=$(git --no-optional-locks -C "$DIR" branch --show-current 2>/dev/null)
    [ -z "$BRANCH" ] && BRANCH=$(git --no-optional-locks -C "$DIR" rev-parse --short HEAD 2>/dev/null)

    STAGED=$(git --no-optional-locks -C "$DIR" diff --cached --name-only 2>/dev/null | wc -l | tr -d ' ')
    MODIFIED=$(git --no-optional-locks -C "$DIR" diff --name-only 2>/dev/null | wc -l | tr -d ' ')

    GIT_STATUS=""
    [ -n "$STAGED" ] && [ "$STAGED" -gt 0 ] && GIT_STATUS="${GIT_STATUS} ${GREEN}+${STAGED}${RESET}"
    [ -n "$MODIFIED" ] && [ "$MODIFIED" -gt 0 ] && GIT_STATUS="${GIT_STATUS} ${YELLOW}~${MODIFIED}${RESET}"

    LINE1="${LINE1} |  ${BRANCH}${GIT_STATUS}"
fi

USED=$(json_get '.context_window.used_percentage')
COST=$(json_get '.cost.total_cost_usd')
DURATION=$(json_get '.cost.total_duration_ms')

USED_FMT=0
if [ -n "$USED" ]; then
    USED_FMT=$(printf '%.0f' "$USED" 2>/dev/null)
    [ -z "$USED_FMT" ] && USED_FMT=0
fi
[ "$USED_FMT" -lt 0 ] && USED_FMT=0
[ "$USED_FMT" -gt 100 ] && USED_FMT=100

BAR_WIDTH=10
FILLED=$((USED_FMT * BAR_WIDTH / 100))
EMPTY=$((BAR_WIDTH - FILLED))
printf -v FILL '%*s' "$FILLED" ''
printf -v PAD '%*s' "$EMPTY" ''
BAR="${FILL// /▓}${PAD// /░}"
LINE2="${BAR} ${USED_FMT}%"

COST_FMT="$0.00"
if [ -n "$COST" ]; then
    COST_FMT=$(printf '$%.2f' "$COST" 2>/dev/null)
    [ -z "$COST_FMT" ] && COST_FMT="$0.00"
fi

MINS=0
SECS=0
if [ -n "$DURATION" ]; then
    DURATION_MS=$(printf '%.0f' "$DURATION" 2>/dev/null)
    if [ -n "$DURATION_MS" ]; then
        TOTAL_SECS=$((DURATION_MS / 1000))
        MINS=$((TOTAL_SECS / 60))
        SECS=$((TOTAL_SECS % 60))
    fi
fi

LINE2="${LINE2} | ${COST_FMT} | ${MINS}m ${SECS}s"

EFFORT=$(json_get '.effort.level')
[ -n "$EFFORT" ] && [ "$EFFORT" != "null" ] && LINE2="${LINE2} | effort:${EFFORT}"

printf '%b\n' "$LINE1"
printf '%s\n' "$LINE2"
