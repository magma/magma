"""
Network Device Management Module

Utilities for creating and managing network interfaces for testing.
"""

import os
import subprocess
from contextlib import contextmanager
from dataclasses import dataclass
from typing import Optional, List, Generator


@dataclass
class VethPairInfo:
    """Information about a veth pair"""
    name_a: str
    name_b: str
    ifindex_a: int
    ifindex_b: int


class NetworkError(Exception):
    """Network operation error"""
    pass


def run_cmd(
    cmd: List[str],
    check: bool = True,
    capture: bool = True
) -> subprocess.CompletedProcess:
    """
    Run a shell command.

    Args:
        cmd: Command and arguments as list
        check: Raise exception on non-zero exit
        capture: Capture stdout/stderr

    Returns:
        CompletedProcess result
    """
    return subprocess.run(
        cmd,
        capture_output=capture,
        text=True,
        check=check
    )


def check_root() -> bool:
    """Check if running as root"""
    return os.geteuid() == 0


def require_root():
    """Raise exception if not running as root"""
    if not check_root():
        raise NetworkError("Root privileges required for network operations")


def interface_exists(iface: str) -> bool:
    """Check if network interface exists"""
    return os.path.exists(f'/sys/class/net/{iface}')


def get_ifindex(iface: str) -> int:
    """
    Get interface index.

    Args:
        iface: Interface name

    Returns:
        Interface index

    Raises:
        NetworkError: If interface doesn't exist
    """
    path = f'/sys/class/net/{iface}/ifindex'
    if not os.path.exists(path):
        raise NetworkError(f"Interface {iface} does not exist")

    with open(path, 'r') as f:
        return int(f.read().strip())


def get_mac_address(iface: str) -> str:
    """
    Get interface MAC address.

    Args:
        iface: Interface name

    Returns:
        MAC address as string (xx:xx:xx:xx:xx:xx)
    """
    path = f'/sys/class/net/{iface}/address'
    if not os.path.exists(path):
        raise NetworkError(f"Interface {iface} does not exist")

    with open(path, 'r') as f:
        return f.read().strip()


def create_veth_pair(name_a: str, name_b: str) -> VethPairInfo:
    """
    Create a veth pair.

    Args:
        name_a: First interface name
        name_b: Second interface name (peer)

    Returns:
        VethPairInfo with interface details

    Raises:
        NetworkError: If creation fails
    """
    require_root()

    # Delete if exists
    delete_interface(name_a)

    try:
        # Create veth pair
        run_cmd(['ip', 'link', 'add', name_a, 'type', 'veth', 'peer', 'name', name_b])

        # Bring both interfaces up
        run_cmd(['ip', 'link', 'set', name_a, 'up'])
        run_cmd(['ip', 'link', 'set', name_b, 'up'])

        # Disable checksum offloading (important for raw packet testing)
        run_cmd(['ethtool', '-K', name_a, 'tx', 'off', 'rx', 'off'], check=False)
        run_cmd(['ethtool', '-K', name_b, 'tx', 'off', 'rx', 'off'], check=False)

        return VethPairInfo(
            name_a=name_a,
            name_b=name_b,
            ifindex_a=get_ifindex(name_a),
            ifindex_b=get_ifindex(name_b)
        )

    except subprocess.CalledProcessError as e:
        raise NetworkError(f"Failed to create veth pair: {e.stderr}")


def delete_interface(iface: str) -> bool:
    """
    Delete a network interface.

    Args:
        iface: Interface name

    Returns:
        True if deleted, False if didn't exist
    """
    if not interface_exists(iface):
        return False

    try:
        run_cmd(['ip', 'link', 'del', iface], check=False)
        return True
    except:
        return False


def configure_interface(
    iface: str,
    ip_addr: Optional[str] = None,
    mtu: Optional[int] = None,
    up: bool = True
) -> None:
    """
    Configure interface parameters.

    Args:
        iface: Interface name
        ip_addr: IP address with prefix (e.g., "10.0.0.1/24")
        mtu: MTU value
        up: Bring interface up
    """
    require_root()

    if not interface_exists(iface):
        raise NetworkError(f"Interface {iface} does not exist")

    if ip_addr:
        # Flush existing addresses
        run_cmd(['ip', 'addr', 'flush', 'dev', iface], check=False)
        # Add new address
        run_cmd(['ip', 'addr', 'add', ip_addr, 'dev', iface])

    if mtu:
        run_cmd(['ip', 'link', 'set', iface, 'mtu', str(mtu)])

    if up:
        run_cmd(['ip', 'link', 'set', iface, 'up'])


def setup_tc_qdisc(iface: str) -> None:
    """
    Add clsact qdisc to interface for TC eBPF attachment.

    Args:
        iface: Interface name

    Raises:
        NetworkError: If setup fails
    """
    require_root()

    if not interface_exists(iface):
        raise NetworkError(f"Interface {iface} does not exist")

    # Remove existing clsact if present
    run_cmd(['tc', 'qdisc', 'del', 'dev', iface, 'clsact'], check=False)

    # Add clsact qdisc
    try:
        run_cmd(['tc', 'qdisc', 'add', 'dev', iface, 'clsact'])
    except subprocess.CalledProcessError as e:
        raise NetworkError(f"Failed to add clsact qdisc to {iface}: {e.stderr}")


def remove_tc_qdisc(iface: str) -> None:
    """
    Remove clsact qdisc from interface.

    Args:
        iface: Interface name
    """
    if interface_exists(iface):
        run_cmd(['tc', 'qdisc', 'del', 'dev', iface, 'clsact'], check=False)


def get_tc_filters(iface: str, direction: str = 'ingress') -> str:
    """
    Get TC filters attached to interface.

    Args:
        iface: Interface name
        direction: 'ingress' or 'egress'

    Returns:
        TC filter listing as string
    """
    if not interface_exists(iface):
        return ""

    result = run_cmd(['tc', 'filter', 'show', 'dev', iface, direction], check=False)
    return result.stdout


@contextmanager
def veth_pair(name_a: str, name_b: str) -> Generator[VethPairInfo, None, None]:
    """
    Context manager for veth pair with automatic cleanup.

    Usage:
        with veth_pair("test_a", "test_b") as veth:
            print(f"Created {veth.name_a} (idx={veth.ifindex_a})")
            # ... do tests ...
        # Automatically cleaned up

    Args:
        name_a: First interface name
        name_b: Second interface name

    Yields:
        VethPairInfo with interface details
    """
    veth_info = None
    try:
        veth_info = create_veth_pair(name_a, name_b)
        # Reload Scapy's interface cache so it sees the new ifindexes.
        # Without this, repeated create/delete of same-name veths causes
        # "OSError: No such device" because Scapy caches stale ifindexes.
        try:
            from scapy.config import conf
            conf.ifaces.reload()
        except Exception:
            pass
        yield veth_info
    finally:
        if veth_info:
            # Clean up TC first
            remove_tc_qdisc(name_a)
            remove_tc_qdisc(name_b)
            # Then delete interfaces
            delete_interface(name_a)


@contextmanager
def tc_qdisc(iface: str) -> Generator[None, None, None]:
    """
    Context manager for TC qdisc with automatic cleanup.

    Usage:
        with tc_qdisc("eth0"):
            # qdisc is set up
            pass
        # Automatically cleaned up

    Args:
        iface: Interface name
    """
    try:
        setup_tc_qdisc(iface)
        yield
    finally:
        remove_tc_qdisc(iface)
