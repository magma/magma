

set -euo pipefail

# Paths
BPF_OBJ="/path/to/bpf/tc_core.o"  # Path to compiled eBPF object
INTERFACES=("br-magma" "gtp0" "sgi0")  # OVS + SGi interfaces

echo "Loading TC eBPF program..."

for IFACE in "${INTERFACES[@]}"; do
    echo "Attaching to interface: $IFACE"

    # Delete existing clsact qdisc (if any)
    sudo tc qdisc del dev "$IFACE" clsact 2>/dev/null || true

    # Add clsact qdisc
    sudo tc qdisc add dev "$IFACE" clsact

    # Attach ingress program
    sudo tc filter add dev "$IFACE" ingress bpf \
        da obj "$BPF_OBJ" sec "tc_ingress"

    # Attach egress program
    sudo tc filter add dev "$IFACE" egress bpf \
        da obj "$BPF_OBJ" sec "tc_egress"

    echo "Attached eBPF program on $IFACE"
done

echo "TC eBPF programs loaded successfully!"
