

set -euo pipefail

# Interfaces to detach eBPF from
INTERFACES=("br-magma" "gtp0" "sgi0")  # Update if needed

echo "Unloading TC eBPF programs..."

for IFACE in "${INTERFACES[@]}"; do
    echo "Removing eBPF programs from interface: $IFACE"

    # Delete ingress & egress filters by removing clsact
    sudo tc qdisc del dev "$IFACE" clsact 2>/dev/null || true

    echo "Detached eBPF program from $IFACE"
done

echo "All TC eBPF programs unloaded successfully!"
