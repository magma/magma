"""
Packet Capture Module

Utilities for capturing and analyzing packets during tests.
"""

import time
import threading
from typing import List, Optional, Callable

from scapy.all import Packet, sniff, AsyncSniffer, conf

# Disable scapy verbosity
conf.verb = 0


class CaptureError(Exception):
    """Packet capture error"""
    pass


class PacketCapture:
    """
    Async packet capture with filtering.

    Usage:
        cap = PacketCapture("eth0", filter="udp port 2152")
        cap.start()
        # ... send packets ...
        packets = cap.stop()
    """

    def __init__(
        self,
        iface: str,
        filter: str = "",
        timeout: float = 10.0,
        count: int = 0
    ):
        """
        Initialize packet capture.

        Args:
            iface: Interface to capture on
            filter: BPF filter string (e.g., "udp port 2152")
            timeout: Maximum capture time in seconds
            count: Stop after capturing this many packets (0 = no limit)
        """
        self.iface = iface
        self.filter = filter
        self.timeout = timeout
        self.count = count

        self._sniffer: Optional[AsyncSniffer] = None
        self._packets: List[Packet] = []
        self._lock = threading.Lock()

    def _packet_handler(self, pkt: Packet) -> None:
        """Handle captured packet"""
        with self._lock:
            self._packets.append(pkt)

    def start(self) -> None:
        """Start capturing packets asynchronously"""
        with self._lock:
            self._packets = []

        self._sniffer = AsyncSniffer(
            iface=self.iface,
            filter=self.filter if self.filter else None,
            prn=self._packet_handler,
            store=False,
            count=self.count if self.count > 0 else 0
        )
        self._sniffer.start()

        # Give sniffer time to initialize
        time.sleep(0.1)

    def stop(self) -> List[Packet]:
        """
        Stop capturing and return captured packets.

        Returns:
            List of captured Scapy Packet objects
        """
        if self._sniffer:
            self._sniffer.stop()
            time.sleep(0.1)  # Allow final packets to be processed
            self._sniffer = None

        with self._lock:
            return list(self._packets)

    def get_packets(self) -> List[Packet]:
        """
        Get currently captured packets without stopping.

        Returns:
            List of captured packets so far
        """
        with self._lock:
            return list(self._packets)

    def wait_for_packet(
        self,
        match: Callable[[Packet], bool],
        timeout: Optional[float] = None
    ) -> Optional[Packet]:
        """
        Wait for a packet matching the given condition.

        Args:
            match: Function that returns True for matching packet
            timeout: Timeout in seconds (default: self.timeout)

        Returns:
            Matching packet or None if timeout
        """
        timeout = timeout or self.timeout
        start_time = time.time()

        while time.time() - start_time < timeout:
            with self._lock:
                for pkt in self._packets:
                    if match(pkt):
                        return pkt
            time.sleep(0.05)

        return None

    def clear(self) -> None:
        """Clear captured packets"""
        with self._lock:
            self._packets = []


def capture_packets(
    iface: str,
    duration: float,
    filter: str = "",
    count: int = 0
) -> List[Packet]:
    """
    Capture packets synchronously for a duration.

    Args:
        iface: Interface to capture on
        duration: Capture duration in seconds
        filter: BPF filter string
        count: Stop after this many packets (0 = no limit)

    Returns:
        List of captured packets
    """
    return sniff(
        iface=iface,
        filter=filter if filter else None,
        timeout=duration,
        count=count if count > 0 else 0
    )


def find_packet(
    packets: List[Packet],
    match: Callable[[Packet], bool]
) -> Optional[Packet]:
    """
    Find first packet matching condition.

    Args:
        packets: List of packets to search
        match: Function that returns True for matching packet

    Returns:
        First matching packet or None
    """
    for pkt in packets:
        if match(pkt):
            return pkt
    return None


def find_all_packets(
    packets: List[Packet],
    match: Callable[[Packet], bool]
) -> List[Packet]:
    """
    Find all packets matching condition.

    Args:
        packets: List of packets to search
        match: Function that returns True for matching packets

    Returns:
        List of matching packets
    """
    return [pkt for pkt in packets if match(pkt)]


def packet_contains_payload(pkt: Packet, payload: bytes) -> bool:
    """
    Check if packet contains the given payload.

    Args:
        pkt: Scapy packet
        payload: Payload bytes to search for

    Returns:
        True if payload found in packet
    """
    try:
        pkt_bytes = bytes(pkt)
        return payload in pkt_bytes
    except:
        return False


def packets_summary(packets: List[Packet]) -> str:
    """
    Generate summary of captured packets.

    Args:
        packets: List of packets

    Returns:
        Summary string
    """
    from scapy.all import IP, UDP, TCP

    lines = [f"Captured {len(packets)} packets:"]

    for i, pkt in enumerate(packets):
        summary = f"  [{i+1}] "

        if IP in pkt:
            ip = pkt[IP]
            summary += f"{ip.src} -> {ip.dst}"

            if UDP in pkt:
                udp = pkt[UDP]
                summary += f" UDP {udp.sport}->{udp.dport}"
            elif TCP in pkt:
                tcp = pkt[TCP]
                summary += f" TCP {tcp.sport}->{tcp.dport}"
            else:
                summary += f" proto={ip.proto}"

            summary += f" len={len(pkt)}"
        else:
            summary += f"Non-IP packet, len={len(pkt)}"

        lines.append(summary)

    return "\n".join(lines)
