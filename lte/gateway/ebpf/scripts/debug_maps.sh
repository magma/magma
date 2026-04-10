

set -euo pipefail

# Path to the pinned maps (adjust if needed)
MAP_DIR="/sys/fs/bpf/magma"

# List of maps to debug (update based on your eBPF maps)
MAPS=("session_map" "stats_map" "metadata_map")

echo "Debugging BPF maps in $MAP_DIR..."

for MAP in "${MAPS[@]}"; do
    MAP_PATH="$MAP_DIR/$MAP"

    if [[ ! -e "$MAP_PATH" ]]; then
        echo "Map $MAP not found at $MAP_PATH"
        continue
    fi

    echo "=== Dumping $MAP ==="
    sudo bpftool map dump pinned "$MAP_PATH"
    echo ""
done

echo "All maps dumped successfully!"
