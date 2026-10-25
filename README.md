# PS3 Common Library

OpenPak's PS3 Common Library.

| | |
| --- | --- |
| Network | PlayStation 3 |
| Language | Go module `openpak/np-common` |
| Licence | AGPL-3.0-only |

## Packages

| Package | What |
| --- | --- |
| `ticket` | Unchanged firmware envelope, Ed25519 verification, avatar capabilities, and current-account/link verification |
| `coreclient`, `accountpb` | Account core gRPC client and generated protocol |
| `identity`, `onlineid` | Consistent NP link provisioning and numeric identity |
| `events` | Per-process standing cache and core invalidation consumer |
| `tlsfront`, `nphttp` | PS3 firmware TLS and bounded HTTP listener lifecycle |
| `resource` | Captured firmware resource signature framing |
| `fakecore` | In-process core for protocol/integration tests and isolated labs |

## Use

```go
import "openpak/np-common/ticket"
```

Each service uses `replace openpak/np-common => ../np-common` and commits `vendor/`; details in [deploy](docs/deploy.md).

## Docs

- [Protocol](docs/protocol.md)
- [Evidence](docs/evidence.md)
- [Deploy](docs/deploy.md)
- [Notes](docs/notes.md)
- [Wire provenance](docs/provenance.md)

## Licence

AGPL-3.0-only. See [LICENSE](LICENSE).
