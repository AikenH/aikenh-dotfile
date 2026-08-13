#!/bin/bash

# Update calendar label (YYYY-MM-DD HH:MM)
sketchybar --set calendar_label label="$(date +'%Y-%m-%d %H:%M')"

# Update monthly calendar popup grid (exact 3-char columns, today with combining underline)
python3 - << 'PY'
import subprocess, sys, datetime

try:
    today_day = datetime.datetime.now().day
    res = subprocess.run(['cal', '-h'], capture_output=True, text=True)
    lines = [l for l in res.stdout.splitlines() if l.strip()]
except Exception:
    sys.exit(0)

if not lines:
    sys.exit(0)

header = lines[0].strip()
weekdays = lines[1] if len(lines) > 1 else ""
rows = lines[2:] if len(lines) > 2 else []

cmds = [
    '--set', 'calendar.header', 'label=' + header,
    '--set', 'calendar.week', 'label=' + weekdays
]

for i in range(1, 7):
    row_idx = i - 1
    if row_idx < len(rows):
        r = rows[row_idx].ljust(21)
        chunks = [r[j:j+3] for j in range(0, 21, 3)]
        formatted_chunks = []
        for c in chunks:
            val_str = c.strip()
            if not val_str:
                formatted_chunks.append('   ')
            elif val_str.isdigit() and int(val_str) == today_day:
                underlined = ''.join([ch + '\u0332' for ch in val_str])
                if len(val_str) == 1:
                    formatted_chunks.append(f' {underlined} ')
                else:
                    formatted_chunks.append(f'{underlined} ')
            else:
                formatted_chunks.append(c)
        formatted_row = ''.join(formatted_chunks)
        cmds.extend(['--set', f'calendar.row{i}', 'drawing=on', 'label=' + formatted_row])
    else:
        cmds.extend(['--set', f'calendar.row{i}', 'drawing=off'])

if cmds:
    subprocess.run(['sketchybar'] + cmds)
PY
