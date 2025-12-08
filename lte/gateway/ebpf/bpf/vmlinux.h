#!/bin/bash
# gen_vmlinux.sh
# Generates vmlinux.h for BPF programs using kernel BTF

set -euo pipefail

OUTPUT_FILE="./vmlinux.h"

echo "[*] Generating vmlinux.h for kernel $(uname -r)..."

if [ ! -r /sys/kernel/btf/vmlinux ]; then
    echo "[!] Error: /sys/kernel/btf/vmlinux not found or not readable."
    echo "    Make sure your kernel has BTF enabled and bpftool is installed."
    exit 1
fi

sudo bpftool btf dump file /sys/kernel/btf/vmlinux format c > "$OUTPUT_FILE"

echo "[+] Generated $OUTPUT_FILE successfully."
