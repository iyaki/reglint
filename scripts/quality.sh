#!/usr/bin/env bash

set -euo pipefail

mode="${1:-all}"

require_cmd() {
	if ! command -v "$1" >/dev/null 2>&1; then
		echo "Missing required tool: $1" >&2
		exit 1
	fi
}

has_go_files() {
	[ -n "$(git ls-files '*.go')" ]
}

run_go_fmt() {
	require_cmd gofmt
	git ls-files -z '*.go' >/tmp/quality-go-files
	if [ -s /tmp/quality-go-files ]; then
		xargs -0 gofmt -w </tmp/quality-go-files
	fi
}

run_golangci() {
	require_cmd golangci-lint
	if ! has_go_files; then
		echo "No Go files found; skipping golangci-lint."
		return 0
	fi
	golangci-lint run
}

run_govulncheck() {
	require_cmd govulncheck
	if ! has_go_files; then
		echo "No Go files found; skipping govulncheck."
		return 0
	fi
	govulncheck ./...
}

run_go_arch_lint() {
	require_cmd go-arch-lint
	if ! has_go_files; then
		echo "No Go files found; skipping go-arch-lint."
		return 0
	fi
	go-arch-lint check
}

run_nancy() {
	require_cmd nancy
	go list -m all | nancy sleuth
}

case "$mode" in
gofmt)
	run_go_fmt
	;;
golangci | golangci-lint)
	run_golangci
	;;
govulncheck)
	run_govulncheck
	;;
nancy)
	run_nancy
	;;
go-arch-lint | arch)
	run_go_arch_lint
	;;
format)
	run_go_fmt
	;;
lint)
	run_golangci
	;;
security)
	run_govulncheck
	run_nancy
	;;
all)
	run_go_fmt
	run_golangci
	run_govulncheck
	run_nancy
	run_go_arch_lint
	;;
*)
	echo "Usage: $0 [all|format|lint|security|arch|gofmt|golangci|govulncheck|nancy|go-arch-lint]" >&2
	exit 1
	;;
esac
