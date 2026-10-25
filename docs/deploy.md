# Deploy: np-common

Services use `replace openpak/np-common => ../np-common` and commit `vendor/`,
matching the Switch pattern. Each service image builds from its own repository.
After changing this module refresh and test every consumer's vendor tree.

`GOWORK=off go test ./...` and `GOWORK=off go vet ./...` verify this module.
See [wire provenance](provenance.md). AGPL-3.0-only.
