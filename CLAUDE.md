# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Smileyan Backend — a Go/Gin REST API for a personal blog system. Supports posts, pages, categories, tags, comments (with email reply notifications), email-based passwordless login, subscriptions, JWT auth, and Bleve full-text search.

Module path: `github.com/smileyan/backend`
Go version: 1.25

## Common Commands

### Run / Build
```bash
go mod tidy                          # install/update deps
go run main.go                       # run dev server on :8080
go build -o smileyan-backend .       # build binary
```

### Local lifecycle scripts (macOS dev)
```bash
./bin/start          # loads .env, frees :8080, runs `go run main.go`, writes PID to logs/app.pid
./bin/stop           # stops via logs/app.pid or port lookup
./bin/restart        # stop + start
```

### Production packaging
```bash
./bin/build-and-package.sh
# Outputs to ./output/:
#   smileyan-backend-dev.tar.gz
#   smileyan-backend-prod.tar.gz
#   smileyan-backend-all.tar.gz
# Each contains a self-contained folder (binary + config.yaml + start.sh/stop.sh/restart.sh + .env.example).
# config.yaml's `mode` is patched to match the env (dev|prod).
```

### Utility commands
```bash
go run cmd/genjwt         # prints two JWTs (a "dev" token and an admin token for i@ccyan.cn).
                          # The dev token only works if the server is started with SMILEYAN_BACKEND_DEV_TOKEN set
                          # (default fallback: "dev-skip-auth-token-2024"). The admin token works against the default
                          # admin_emails list in config.yaml. Use as: curl -H "Authorization: Bearer <token>" ...
go run scripts/rebuild_html.go   # re-renders HTMLContent for every Post/Page from its Markdown Content.
                                  # Run this after changing the Markdown→HTML pipeline (utils/markdown.go).
```

### Tests
There are no `*_test.go` files in the repo. No test framework is wired up; add one before writing unit tests.

## Configuration

Two layers, both required:

1. **`config.yaml`** — non-sensitive defaults (server port/mode, DB host/user/dbname, Redis host/port/db, email host/port/username, upload paths, `admin_emails` list, JWT expire hours). Values here are fallbacks.
2. **`.env`** (gitignored) — secret overrides. Required vars: `SMILEYAN_BACKEND_DB_PASSWORD`, `SMILEYAN_BACKEND_REDIS_PASSWORD`, `SMILEYAN_BACKEND_EMAIL_PASSWORD`, `SMILEYAN_BACKEND_JWT_SECRET`. Optional: `SMILEYAN_BACKEND_DB_HOST/USER/NAME`, `SMILEYAN_BACKEND_REDIS_HOST/USERNAME`, `SMILEYAN_BACKEND_DEV_TOKEN`.

Viper is configured with `SetEnvPrefix("SMILEYAN_BACKEND")`, and `config/config.go::GetConfig()` then **explicitly re-reads** the env vars listed above and overrides the unmarshalled config. That's why prefix matching isn't enough — the override list in `GetConfig()` is the source of truth for which env vars win.

`config.AdminEmails` (drives `cfg.IsAdminEmail`) controls:
- Who is granted `RoleAdmin` on first login (`controllers/user.go::Login`).
- Who bypasses the per-email rate limit on `POST /api/send-code`.

`models.AutoMigrate()` runs on every startup and creates a default `root@smileyan.cn` admin user **only if the `users` table is empty**.

## High-Level Architecture

The app follows a classic layered structure with no service layer between controllers and GORM. Entry point is `main.go`, which wires everything up in this order: `config.InitConfig` → `utils.InitLogger` → `config.InitDatabase` → `config.InitRedis` → `models.AutoMigrate` → `services.InitSearchService` → seed posts into the in-memory Bleve index → `routes.SetupRoutes` → `r.Run`. The `createUploadDirs()` function in `main.go` is a no-op; uploads rely on the process having `./uploads/avatars/` present (or the upload call failing).

```
main.go
  └─ config/      Viper + GORM + go-redis bootstrap. Exposes package-level globals: config.DB, config.Redis, config.Config.
  └─ models/      GORM models (User, Post, Page, Category, Tag, Comment, Subscription) and AutoMigrate.
  └─ utils/       Logger (zap), JWT, email (go-mail), Markdown→HTML (goldmark + post-processing), random strings, verification codes.
  └─ middleware/  AuthMiddleware (JWT bearer, with dev-token bypass for user_id=1), AdminMiddleware, GetCurrentUser, LoggerMiddleware.
  └─ controllers/ Thin Gin handlers. Mostly direct GORM via config.DB. Search delegates to services.
  └─ services/    SearchService wrapping a bleve.NewMemOnly() index, protected by a sync.RWMutex.
  └─ routes/      Three groups: public /api, authenticated /api, /api/admin (auth + admin middleware).
  └─ cmd/genjwt   Standalone CLI that mints JWTs.
  └─ scripts/     rebuild_html.go — one-shot DB re-render tool.
  └─ bin/         Local dev start/stop/restart and the cross-build packager.
```

### Auth flow
1. `POST /api/send-code` — admin emails skip rate limiting; everyone else is rate-limited via Redis (`rate_limit:email:<email>` 10min, `rate_limit:ip:<ip>` 100/hr). A 6-digit code is stored in Redis (`verification_code:<email>`, 10min TTL) and emailed via SMTP.
2. `POST /api/login` — verifies code, finds-or-creates a `User` (admin role granted by `IsAdminEmail`), returns `{token, user}`. JWT is HS256 with `UserID/Email/Role` claims and a 7-day default expiry.
3. `AuthMiddleware` decodes the bearer token. **Special case:** if the server was started with `SMILEYAN_BACKEND_DEV_TOKEN` set, a JWT whose `dev` claim matches that env var logs the request in as a synthetic `user_id=1` admin. This is what `cmd/genjwt`'s "方式1" token exercises. The `GetCurrentUser` middleware also short-circuits `user_id == 1` to return a hardcoded admin `User{}` without touching the DB.
4. `AdminMiddleware` then enforces `role == "admin"` for the `/api/admin/*` group.

### Post lifecycle (the most behavior-rich area)
- `CreatePost`/`UpdatePost` render Markdown → HTML via `utils.MarkdownToHTML` and write **both** `Content` (raw MD) and `HTMLContent` (rendered). The `HTMLContent` post-processing intentionally matches the frontend's `markdown-it` output (see `utils/markdown.go::postProcessHTML`) — element classes like `article-h1`, `article-figure`, `article-table` are stable contracts with the Vue frontend. Change the rendering only in coordination with the frontend, and re-run `go run scripts/rebuild_html.go` after.
- Slug generation (`generateSlug` in `controllers/post.go`) lowercases, replaces spaces with `-`, strips everything except `[a-z0-9-]` and CJK characters (`r >= 0x4E00`).
- `DeletePost` is **soft by default**: it sets `is_deleted=true`, `status=hidden`, and rewrites `slug` to an 8-char random hex (`utils.GenerateRandomString(8)`) so the URL slot is freed. Only when called on an already-hidden post does it hard-delete (`Unscoped().Delete`). `RestorePost` flips the flags back and re-adds to the search index.
- Slug uniqueness is enforced including soft-deleted rows (`Unscoped().Where("slug = ?", ...)`).
- Public `GetPost` returns 404 for non-published posts unless the caller is an admin **and** passes `?edit=true` (which also suppresses the `view_count` increment).
- Category/tag list endpoints include `post_count` for *published, non-deleted* posts via a per-row `Count`.

### Search
`services.InitSearchService` creates an **in-memory** `bleve.NewMemOnly` index — it does not persist across restarts. `main.go` reloads every existing post into it on startup. Every Create/Update/Restore/Delete goes through `services.AddToSearchIndex`/`UpdateSearchIndex`/`RemoveFromSearchIndex`. The search index is keyed by `post.Slug` and looks up the full record from MySQL after hitting Bleve. If you need a persistent index, this is the file to change.

### Middleware behavior worth remembering
- `AuthMiddleware` aborts with 401 on missing/malformed/expired bearer tokens, including JSON bodies, so handlers behind it can assume the context values are set.
- `AdminMiddleware` must come **after** `AuthMiddleware` in the chain (it relies on the `role` context key being present).
- `GetCurrentUser` returns `nil` for unknown IDs — every handler must nil-check, e.g. `controllers/post.go::CreatePost` does so explicitly even though the route is behind `AuthMiddleware`.

### Conventions / non-obvious details
- IDs in admin routes use the path segment `/_id/:id` (e.g. `PUT /api/admin/posts/_id/:id`) to avoid colliding with other slug-based routes. Match that pattern when adding admin endpoints.
- `models.Post` has both `gorm.DeletedAt` (used by GORM soft-delete via `db.Delete`) **and** an explicit `is_deleted bool` flag with its own index. The flag is what soft-delete and restore go through, not `gorm.DeletedAt` — most queries filter `is_deleted = false` rather than relying on GORM's auto-scoping.
- `controllers/page.go` and `controllers/subscription.go` roll their own `parseInt` instead of using `strconv.Atoi` (inconsistent with the rest of the codebase). Prefer `strconv.Atoi` in new code.
- `utils/jwt.go::Claims` is the claims struct used when *minting* tokens. The middleware in `middleware/auth.go` defines its own `Claims` struct (with an extra `Dev string` field) — they are intentionally separate. The dev-token bypass compares `claims.Dev` against the `SMILEYAN_BACKEND_DEV_TOKEN` env var.
- `utils.MarkdownToHTML` always returns the post-processed HTML; on a goldmark error it returns an empty `[]byte`, never an error.
- The frontend (Vue) reads from `https://smileyan.cn` / `https://bigbigpig.cn` — see hardcoded URLs in `controllers/subscription.go` and `cmd/genjwt` usage example.

## Key Files
- `main.go` — startup sequence and route registration
- `config/config.go` — Viper wiring and the env-var override list
- `models/models.go`, `models/migrate.go` — schema and `AutoMigrate` (incl. default admin seed)
- `middleware/auth.go` — JWT verification, dev-token bypass, role/admin gates
- `controllers/post.go` — most of the business logic (CRUD, slug, soft-delete, restore)
- `utils/markdown.go` — Markdown→HTML pipeline that must stay in sync with the frontend
- `services/search.go` — in-memory Bleve index
- `routes/routes.go` — the canonical list of endpoints and their middleware chains
- `docs/api.md` — the full request/response contract for the public API
- `docs/database.md` — schema details beyond what `models.go` shows
