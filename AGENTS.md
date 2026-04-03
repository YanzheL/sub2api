# PROJECT KNOWLEDGE BASE

**Generated:** 2026-04-04
**Commit:** 7c8923d6
**Branch:** next

## OVERVIEW
Sub2API is an AI API gateway platform with a Go backend, a Vue 3/Vite frontend, and a deployment subtree for Docker and systemd installs. The repo is effectively split into three working domains: `backend/`, `frontend/`, and `deploy/`.

## STRUCTURE
```text
sub2api/
├── backend/      # Go service, Ent schema/codegen, SQL migrations
├── frontend/     # Vue 3 admin + user UI, built into backend/internal/web/dist
├── deploy/       # Docker/systemd/install artifacts and ops docs
├── .github/      # CI, release, security scan workflows
├── docs/         # Additional documentation
└── tools/        # Repo utilities such as secret/audit helpers
```

## CHILD GUIDES
| Path | Scope |
|------|-------|
| `backend/AGENTS.md` | Go architecture, backend commands, layer boundaries |
| `backend/ent/AGENTS.md` | Generated Ent boundary and schema workflow |
| `frontend/AGENTS.md` | Vue app architecture, pnpm workflow, local docs |
| `deploy/AGENTS.md` | Docker/systemd deployment files and ops-safe changes |

## WHERE TO LOOK
| Task | Location | Notes |
|------|----------|-------|
| Backend bootstrap | `backend/cmd/server/main.go` | Setup mode, auto-setup, normal server startup |
| Route registration | `backend/internal/server/router.go`, `backend/internal/server/routes/` | HTTP entrypoint after bootstrap |
| Gateway behavior | `backend/internal/service/gateway_service.go` and provider-specific gateway files | Core request routing and upstream logic |
| Usage/billing persistence | `backend/internal/repository/usage_log_repo.go` | Hot path for accounting data |
| Data model changes | `backend/ent/schema/`, `backend/migrations/` | Schema is handwritten; migrations are forward-only |
| Frontend bootstrap | `frontend/src/main.ts`, `frontend/src/App.vue` | Theme, config injection, i18n, router readiness |
| Frontend routing/auth | `frontend/src/router/`, `frontend/src/stores/auth.ts` | Route guards and auth state |
| Admin UI | `frontend/src/views/admin/`, `frontend/src/components/admin/`, `frontend/src/api/admin/` | View/API/component split |
| Shared UI | `frontend/src/components/common/`, `frontend/src/components/layout/` | Reusable tables, dialogs, layout shell |
| Deploy/runtime config | `deploy/README.md`, `deploy/config.example.yaml`, `deploy/.env.example` | Canonical ops docs live here |

## CONVENTIONS
- Frontend uses **pnpm**, not npm/yarn. CI installs with `pnpm install --frozen-lockfile`.
- Frontend production build writes to `backend/internal/web/dist/`; that output is only served by Go when using an embed-enabled build path.
- Backend codegen is part of normal workflow: `cd backend && make generate` runs Ent generation and Wire generation.
- Go tests are split by build tags: `unit`, `integration`, `e2e`.
- Backend CI pins exact Go `1.26.1` and frontend CI runs lint, typecheck, Vitest, and Playwright.
- `deploy/docker-compose.local.yml` is the migration-friendly Docker variant; named-volume compose is the simpler alternative.

## ANTI-PATTERNS (THIS PROJECT)
- Do not hand-edit generated Ent files under `backend/ent/` outside `backend/ent/schema/`.
- Do not modify already-applied SQL migrations; add a new migration instead.
- Do not edit `backend/internal/web/dist/`; edit `frontend/src/` and rebuild.

## COMMANDS
```bash
make build
make test
make secret-scan
```

## NOTES
- `RUN_MODE=simple` disables billing/quota enforcement; treat it as a special deployment mode, not the default product path.
- A plain backend build does not automatically embed the freshly built frontend; check child guides when you need a single binary that serves UI assets.
- Sora-related support is documented as temporarily unavailable in the main README; avoid treating it as a stable production path.
- When proxying through Nginx for CLI traffic, the repo README explicitly calls out `underscores_in_headers on;`.
- Antigravity and Anthropic mixed-context behavior has repo-specific caveats; verify group/account routing before changing that area.
