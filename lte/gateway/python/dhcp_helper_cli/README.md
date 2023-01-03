# DHCP Helper CLI

This is a CLI tool to allocate, renew, and release IPs from a DHCP server.

## Usage

To allocate:

```bash
dhcp_helper_cli.py [-h] --mac MAC [--json] [--vlan VLAN] [--interface INTERFACE] allocate
```

To renew or release a used IP:

```bash
dhcp_helper_cli.py [-h] --mac MAC [--json] [--vlan VLAN] [--interface INTERFACE] {renew,release}
 --ip IP --server-ip SERVER-IP
```

## Testing

The script can be tested by setting up a VM with a DHCP server running on it,
e.g. `udhcpd`. Then run the script, e.g. on the host, and set `INTERFACE` to
the one that connects to the VM. The allocated leases should appear on the VM
and can be inspected with `dumpleases`.
