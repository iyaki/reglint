#!/usr/bin/env bash

set -euo pipefail

if ! command -v go >/dev/null 2>&1; then
	echo "Go is required to install tooling." >&2
	exit 1
fi

export GOBIN="${GOBIN:-$(go env GOPATH)/bin}"
mkdir -p "$GOBIN"

curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.10.1
go install golang.org/x/vuln/cmd/govulncheck@latest
curl -sSfL -o $GOPATH/bin/nancy https://github.com/sonatype-nexus-community/nancy/releases/download/v1.2.0/nancy-v1.2.0-linux-amd64
chmod +x $GOPATH/bin/nancy
go install github.com/fe3dback/go-arch-lint@latest
# go install -v github.com/avito-tech/go-mutesting/...
go install github.com/evilmartians/lefthook/v2@v2.1.2
