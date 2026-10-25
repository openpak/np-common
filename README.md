# np-common

Shared Go primitives for OpenPak's independent PS3 services. No daemon or image.

- `ticket`: unchanged firmware envelope, Ed25519 verification, avatar capabilities,
  and current-account/link verification.
- `coreclient`, `accountpb`: account core gRPC client and generated protocol.
- `identity`, `onlineid`: consistent NP link provisioning and numeric identity.
- `events`: per-process standing cache and core invalidation consumer.
- `tlsfront`, `nphttp`: PS3 firmware TLS and bounded HTTP listener lifecycle.
- `resource`: captured firmware resource signature framing.
- `fakecore`: in-process core for protocol/integration tests and isolated labs.

Services use `replace openpak/np-common => ../np-common` and commit `vendor/`,
matching the Switch pattern. Each service image builds from its own repository.
After changing this module refresh and test every consumer's vendor tree.

`GOWORK=off go test ./...` and `GOWORK=off go vet ./...` verify this module.
See [wire provenance](docs/provenance.md). AGPL-3.0-only.
