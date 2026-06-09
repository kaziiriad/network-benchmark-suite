# Session Restore Guide

## What's in `.claude/`

```
.claude/
├── sessions/
│   └── c7c59082-733b-448a-968f-e80390c0e1a5/   ← current session
│       ├── ethereal-sparking-wind.md            ← active plan
│       └── c7c59082-...jsonl                   ← full transcript
├── sessionstart-hook-3.sh                      ← auto-runs on session start
└── settings.local.json                         ← project-level config
```

## Restoring on a New Machine

### Step 1: Copy project

```bash
rsync -av /path/to/npbs_project/ new-machine:/path/to/
```

### Step 2: Create project-level settings

On new machine, create `.claude/settings.local.json` inside the project:

```json
{
  "outputStyle": "Learning",
  "plansDirectory": ".claude/sessions/c7c59082-733b-448a-968f-e80390c0e1a5"
}
```

### Step 3: Configure SessionStart hook (global)

Add to `~/.claude/settings.json`:

```json
"hooks": {
  "PostToolUse": [{
    "matcher": "SessionStart",
    "hooks": [{
      "type": "command",
      "command": "bash /absolute/path/to/npbs_project/.claude/sessionstart-hook-3.sh"
    }]
  }]
}
```

### Step 4: Resume

Start Claude Code in the project directory. Hook runs, plan loads, session folder active.

## What Transfers

| File | Portable | Notes |
|------|----------|-------|
| `settings.local.json` | yes | Paste content into new machine's `.claude/settings.local.json` |
| `sessionstart-hook-3.sh` | yes | Hook activation is global — see Step 3 |
| `sessions/*/ethereal-sparking-wind.md` | yes | Plan auto-saves to session folder |
| `sessions/*/*.jsonl` | manual | Paste contents as initial context on new machine |
| `plan/` | partial | Old plans from before session folder setup |

## Session ID Note

Claude Code session IDs are machine-local. Moving `.claude/` to a new machine starts a fresh session. The `.jsonl` transcript is your恢复 manual reference — paste relevant context at session start if needed.