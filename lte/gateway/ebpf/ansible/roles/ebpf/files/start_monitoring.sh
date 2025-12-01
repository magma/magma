#!/bin/bash
# start_monitoring.sh
# Starts Python eBPF observer for real-time monitoring

PYTHON_BIN=${1:-python3}

echo "Starting eBPF observer..."
$PYTHON_BIN python/ebpf_observer.py &

OBSERVER_PID=$!
echo "eBPF observer started with PID $OBSERVER_PID"

# Optional: redirect logs to a file
# $PYTHON_BIN python/ebpf_observer.py > logs/ebpf_observer.log 2>&1 &
