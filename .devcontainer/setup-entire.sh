#!/usr/bin/env bash

# https://docs.entire.io/cli/installation#linux

# Bash strict mode - https://olivergondza.github.io/2019/10/01/bash-strict-mode.html
set -euo pipefail
# shellcheck disable=SC2154
trap 's=$?; echo "$0: Error on line "$LINENO": $BASH_COMMAND"; exit $s' ERR

mkdir entire

cd entire

# Download the latest release
curl -sSfLO https://github.com/entireio/cli/releases/latest/download/entire_linux_amd64.tar.gz

# Extract and install
tar -xzf entire_linux_amd64.tar.gz
sudo mv entire /usr/local/bin/
sudo mv completions/entire.bash /etc/bash_completion.d/entire

cd ..
rm -rf entire

entire enable --agent opencode
