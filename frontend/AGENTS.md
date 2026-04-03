# FRONTEND GUIDE

## OVERVIEW
`frontend/` is a Vue 3 + Vite + Pinia application. It builds into `backend/internal/web/dist/`, so frontend changes often affect backend embed output as well.

## STRUCTURE
```text
frontend/
├── package.json         # pnpm scripts
├── vite.config.ts       # Vite config, proxy, build output
├── vitest.config.ts     # unit/integration test config
├── playwright.config.ts # E2E config
└── src/
    ├── main.ts          # bootstrap
    ├── App.vue          # root component
    ├── router/          # route table + guards
    ├── stores/          # Pinia stores
    ├── api/             # Axios client + API modules
    ├── views/           # route targets
    ├── components/      # UI building blocks
    ├── composables/     # reusable state/behavior hooks
    ├── i18n/            # locale setup + messages
    ├── utils/           # formatting and helper logic
    └── __tests__/       # shared test setup + integration tests
```

## WHERE TO LOOK
| Task | Location | Notes |
|------|----------|-------|
| App bootstrap | `src/main.ts`, `src/App.vue` | Theme/config/i18n/router init sequence |
| Auth flow | `src/stores/auth.ts`, `src/api/auth.ts`, `src/views/auth/` | Router guards depend on this stack |
| Route behavior | `src/router/` | Lazy-loaded routes, auth/admin gating, title handling |
| Shared app state | `src/stores/` | `app.ts`, `adminSettings.ts`, `subscriptions.ts`, etc. |
| Backend API calls | `src/api/client.ts`, `src/api/`, `src/api/admin/` | Reuse the existing client/interceptors |
| Shared UI | `src/components/common/`, `src/components/layout/` | Tables, dialogs, badges, app shell |
| Admin pages | `src/views/admin/`, `src/components/admin/`, `src/api/admin/` | Admin view/API/component split |
| Reusable behavior | `src/composables/` | `use*` helpers with co-located tests |
| Local docs | `src/router/README.md`, `src/stores/README.md`, `src/components/common/README.md`, `src/components/layout/README.md`, `src/views/auth/README.md` | Helpful context; verify against current code before trusting details |

## CONVENTIONS
- Use **pnpm**. CI and release automation assume the pnpm lockfile.
- `src/main.ts` initializes config/theme/i18n before mount and waits for `router.isReady()`; preserve that startup ordering.
- Route files rely on route meta (`requiresAuth`, `requiresAdmin`, `titleKey`, etc.) rather than ad-hoc guard logic inside views.
- API calls should normally flow through `src/api/client.ts` and the existing module barrels, not raw Axios instances created in random components.
- Many narrow frontend areas already have local README files; use them as supplemental context, but verify behavior against the current code before copying conventions.
- Tests are mostly co-located in `__tests__/` directories, with top-level integration tests under `src/__tests__/integration/`.

## COMMANDS
```bash
cd frontend && pnpm install
cd frontend && pnpm run dev
cd frontend && pnpm run build
cd frontend && pnpm run lint:check
cd frontend && pnpm run typecheck
cd frontend && pnpm run test:run
cd frontend && pnpm run test:e2e
```

## ANTI-PATTERNS
- Do not edit `../backend/internal/web/dist/` directly.
- Do not bypass `src/api/client.ts` unless you are intentionally changing HTTP client behavior.
- Do not change `package.json` without updating `pnpm-lock.yaml`.
- Do not duplicate conventions already documented in local frontend README files.

## NOTES
- `src/router/index.ts` is large and central; route changes often imply store/API/title impacts.
- `src/views/admin/` and `src/components/admin/` hold most of the admin-dashboard complexity.
- `src/i18n/locales/en.ts` and `zh.ts` are translation datasets, not logic hotspots.
