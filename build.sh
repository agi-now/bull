#!/usr/bin/env bash
set -euo pipefail

VERSION=${VERSION:-"dev"}
LDFLAGS="-s -w -X main.version=${VERSION}"
OUT_DIR="bin"

mkdir -p "${OUT_DIR}"

echo "Building bull ${VERSION} ..."

# Current platform
go build -ldflags="${LDFLAGS}" -o "${OUT_DIR}/bull" ./cmd/bull/
echo "  -> ${OUT_DIR}/bull"

# Cross compile (optional, pass CROSS=1)
if [ "${CROSS:-0}" = "1" ]; then
  for pair in linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64; do
    os="${pair%%/*}"
    arch="${pair##*/}"
    ext=""
    [ "$os" = "windows" ] && ext=".exe"
    output="${OUT_DIR}/bull-${os}-${arch}${ext}"
    echo "  building ${os}/${arch} ..."
    GOOS=$os GOARCH=$arch go build -ldflags="${LDFLAGS}" -o "${output}" ./cmd/bull/
    echo "  -> ${output}"
  done
fi

echo "Done."
