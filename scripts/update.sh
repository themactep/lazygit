#!/usr/bin/env bash
set -euo pipefail

SELF_DIR="$(cd "$(dirname "$0")"/.. && pwd)"
cd "$SELF_DIR"

echo "=== Fetching upstream ==="
git fetch upstream

echo "=== Rebasing current branch onto upstream/master ==="
git rebase upstream/master

echo "=== Building ==="
make build

echo "=== Installing ==="
go install

echo "=== Done ==="
