# Changelog: np-common

## v0.1.2 — 2026-11-03

- New `httplog` package: one stderr line per request, `http <method> <path> -> <status> <duration>ms`.
  Never the query string, headers, cookies, tokens, passwords or bodies; `/healthz` is skipped; the
  status writer keeps `http.Flusher`, `http.Hijacker` and `http.ResponseController` working.
  `httplog.Command` gives non-HTTP handlers the same shape (`cmd <name> -> <result> <duration>ms`).
- `nphttp.Logging` now uses it (it used slog without a duration); the PS3 status header is appended
  as `ps3_status=<n>` when the handler sets one.
