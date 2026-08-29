#!/bin/sh
set -e

# Named volumes are created as root. Ensure the non-root appuser can write
# to the Go module and build caches (and the air tmp dir).
for dir in \
    /home/appuser/.cache \
    /home/appuser/.cache/go-build \
    /go/pkg/mod \
    /workspace/tmp
do
    mkdir -p "$dir"
    chown -R appuser:appgroup "$dir" 2>/dev/null || true
done

# Drop privileges and run the original command (air by default)
exec gosu appuser "$@"
