#!/bin/bash
# gen_vmlinux.sh
# Generates vmlinux.h for BPF programs

OUTPUT_FILE="./vmlinux.h"

echo "[*] Generating vmlinux.h for kernel $(uname -r)..."
sudo bpftool btf dump file /sys/kernel/btf/vmlinux format c > "$OUTPUT_FILE"

echo "[+] Generated $OUTPUT_FILE"
