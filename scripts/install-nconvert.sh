#!/bin/sh
set -eu

SRC="${1:-./third_party/nconvert/nconvert}"
DST="${2:-/usr/local/bin/nconvert}"

if [ ! -f "$SRC" ]; then
  echo "nconvert bundle not found at: $SRC" >&2
  exit 1
fi

install -d "$(dirname "$DST")"
install -m 0755 "$SRC" "$DST"

echo "installed nconvert -> $DST"
