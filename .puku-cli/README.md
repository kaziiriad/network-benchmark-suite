# Session Configuration (canonical, in-repo)

This directory holds the **tracked, canonical copy** of the session configuration
that the Puku CLI harness reads at session start. The live, working copy lives
one level up at `/home/poridhian/code/.puku-cli/` and is **machine-local** —
not tracked, not pushed, regenerated on each new machine.

## Files

| File | Role | Tracked? |
|---|---|---|
| `settings.local.json` | Harness config: output style, plans directory | yes |
| `sessionstart-hook-3.sh` | Runs on session start, prints restored context | yes |
| `README.session.md` | Migration procedure (Claude Code / portable setup) | yes |
| `README.md` | This file — explains the in-repo vs live split | yes |

The **transcript** (`c7c59082-.../c7c59082-...jsonl`, ~3.7 MB) lives only in the
live `.puku-cli/` directory and is **deliberately not tracked**. See the plan
file (`PLAN.md`) and git history for the durable record of decisions.

## How the live and in-repo copies relate

```
/home/poridhian/code/
├── network-benchmark-suite/        ← git repo (this directory)
│   ├── PLAN.md                     ← tracked plan
│   ├── .puku-cli/                  ← tracked config (canonical)
│   │   ├── settings.local.json
│   │   ├── sessionstart-hook-3.sh
│   │   ├── README.session.md
│   │   └── README.md
│   ├── CLAUDE.md
│   └── .gitignore
└── .puku-cli/                      ← machine-local, not tracked
    ├── settings.local.json         ← active session reads this
    ├── sessionstart-hook-3.sh
    ├── c7c59082-.../               ← active session directory
    │   ├── ethereal-sparking-wind.md
    │   └── c7c59082-...jsonl       ← transcript, never tracked
    └── ...
```

For the **current session**, the live `.puku-cli/` is authoritative and must
not be overwritten — the harness is reading from it. On a **new machine**,
copy the three files in this directory into the live `.puku-cli/` location
to reproduce the session setup.
