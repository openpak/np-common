# Changelog: np-common

## v0.2.1 - 2026-11-04

- New `capture` package and `NP_CAPTURE_DIR` in `nphttp.Run`: full request capture to daily
  `capture-YYYY-MM-DD.jsonl` (host, URI, headers, body up to 16 KiB, status, response headers, TLS client
  hello SNI), kept 90 days, mode 0600. Credential headers, token/ticket/password-looking query and form
  values and the bodies of authentication paths are masked or skipped before writing. Off when the variable
  is unset. For card #130 P0 (PS3 Store protocol discovery).

## v0.2.0 - 2026-11-04

- `events.Consumer.Purge` (optional): called once per `account_deleted` event so a service can drop its own
  per-account rows (account erasure, card #120). Unset keeps the old behaviour (cut sessions only).

## v0.1.2 — 2026-11-03

- New `httplog` package: one stderr line per request, `http <method> <path> -> <status> <duration>ms`.
  Never the query string, headers, cookies, tokens, passwords or bodies; `/healthz` is skipped; the
  status writer keeps `http.Flusher`, `http.Hijacker` and `http.ResponseController` working.
  `httplog.Command` gives non-HTTP handlers the same shape (`cmd <name> -> <result> <duration>ms`).
- `nphttp.Logging` now uses it (it used slog without a duration); the PS3 status header is appended
  as `ps3_status=<n>` when the handler sets one.
