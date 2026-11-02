# Protocol: np-common

Nothing recorded yet.

## Packages

- `ticket`: unchanged firmware envelope, Ed25519 verification, avatar capabilities, and current-account/link verification.
- `coreclient`: account core gRPC client (protocol from `github.com/openpak/account-api`).
- `identity`, `onlineid`: consistent NP link provisioning and numeric identity.
- `events`: per-process standing cache and core invalidation consumer.
- `tlsfront`, `nphttp`: PS3 firmware TLS and bounded HTTP listener lifecycle.
- `resource`: captured firmware resource signature framing.
- `fakecore`: in-process core for protocol/integration tests and isolated labs.
