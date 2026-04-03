# DEPLOY GUIDE

## OVERVIEW
`deploy/` is the operational subtree for Docker Compose, binary install, systemd units, reverse proxy config, and runtime configuration examples.

## STRUCTURE
```text
deploy/
├── README.md                     # canonical deployment guide
├── docker-compose.yml            # named-volume compose
├── docker-compose.local.yml      # local-directory compose (recommended for easy migration)
├── docker-compose.dev.yml        # dev compose
├── docker-compose.standalone.yml # app with external DB/Redis
├── docker-deploy.sh              # one-click Docker prep
├── install.sh                    # one-click binary/systemd install
├── install-datamanagementd.sh    # companion daemon installer
├── sub2api.service               # hardened systemd unit
├── sub2api-datamanagementd.service # companion daemon systemd unit
├── Caddyfile                     # sample reverse proxy config
├── config.example.yaml           # full runtime config example
└── .env.example                  # Docker env template
```

## WHERE TO LOOK
| Task | Location | Notes |
|------|----------|-------|
| Choose deployment mode | `README.md` | Docker vs binary install guidance lives here |
| Fast Docker setup | `docker-deploy.sh`, `docker-compose.local.yml` | Local-dir variant is the repo's recommended path |
| Binary/systemd install | `install.sh`, `sub2api.service` | Includes service hardening and upgrade flow |
| Runtime env vars | `.env.example` | Docker-oriented environment template |
| Full app config | `config.example.yaml` | Authoritative config shape and comments |
| Reverse proxy | `Caddyfile` | Helpful baseline proxy config |
| Release pipeline context | `../.github/workflows/release.yml`, `../.goreleaser*.yaml` | Artifact production happens outside this subtree |

## CONVENTIONS
- `README.md` is the canonical human runbook; keep AGENTS notes here short and point back to it.
- Prefer `docker-compose.local.yml` when documenting migration/backup-friendly setups; it uses local directories instead of named volumes.
- Docker flows and binary-install flows are both first-class. Keep instructions explicit about which path a change affects.
- This subtree also covers `datamanagementd`; its install script and service file are separate from the main app.
- `datamanagementd` integration assumes the fixed socket path `/tmp/sub2api-datamanagement.sock`; Docker setups must bind-mount that host socket to the same in-container path.

## COMMANDS
```bash
cd deploy && docker compose -f docker-compose.local.yml up -d
cd deploy && docker compose -f docker-compose.local.yml logs -f sub2api
cd deploy && docker compose -f docker-compose.local.yml down

curl -sSL https://raw.githubusercontent.com/Wei-Shaw/sub2api/main/deploy/install.sh | sudo bash
```

## ANTI-PATTERNS
- Do not commit real secrets or filled-in production env files.
- Do not change install/deploy behavior without updating `deploy/README.md` if operator steps changed.
- Do not assume named-volume Compose and local-directory Compose are interchangeable; document which one a change targets.
- Do not treat release automation files as runtime docs; release packaging lives in repo root workflows/configs.

## NOTES
- `install.sh` is large and operationally important; review it carefully before changing upgrade/install semantics.
- `sub2api.service` uses systemd hardening options; keep those constraints in mind when changing file paths or runtime permissions.
- `DATAMANAGEMENTD_CN.md` documents the companion socket-based daemon path and operational coupling.
