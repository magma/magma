#!/usr/bin/env bash
set -euo pipefail

# Output path for vmlinux.h
OUT_FILE="./vmlinux.h"

echo "Generating vmlinux.h..."

# Check if bpftool exists
if ! command -v bpftool &> /dev/null; then
    echo "bpftool not found. Please install linux-tools (bpftool)."
    exit 1
fi

# Generate BTF from running kernel
sudo bpftool btf dump file /sys/kernel/btf/vmlinux format c > "$OUT_FILE"

echo "vmlinux.h generated at $OUT_FILE"
