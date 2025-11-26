#!/usr/bin/env python3
"""
ebpf_observer.py

Observer for eBPF session maps, stats, and GTP flow events.
- Monitors pinned BPF maps
- Provides periodic logging
- Optional callback/event hooks for PipelineD or SessionD
"""

import logging
import threading
import time
from ebpf_utils import read_bpf_map

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='[%(asctime)s] %(levelname)s: %(message)s'
)

# Pinned maps
SESSION_MAP = "/sys/fs/bpf/gtp_session_map"
UL_MAP = "/sys/fs/bpf/gtp_ul_map"
DL_MAP = "/sys/fs/bpf/gtp_dl_map"
STATS_MAP = "/sys/fs/bpf/gtp_stats_map"


class EBPFObserver:
    """
    Observer class to monitor eBPF maps
    """

    def __init__(self, poll_interval: float = 5.0, callback=None):
        """
        :param poll_interval: How often to poll BPF maps (seconds)
        :param callback: Optional callback for changes: func(map_name:str, entries:list)
        """
        self.poll_interval = poll_interval
        self.callback = callback
        self.running = False

    # ------------------------------------------------------------
    # Start/Stop Observer
    # ------------------------------------------------------------
    def start(self):
        logging.info("[Observer] Starting eBPF observer thread...")
        self.running = True
        th = threading.Thread(target=self._observer_loop, daemon=True)
        th.start()

    def stop(self):
        logging.info("[Observer] Stopping observer...")
        self.running = False

    # ------------------------------------------------------------
    # Polling Loop
    # ------------------------------------------------------------
    def _observer_loop(self):
        logging.info("[Observer] Observer loop running.")
        while self.running:
            self.poll_maps()
            time.sleep(self.poll_interval)

    # ------------------------------------------------------------
    # Poll Maps
    # ------------------------------------------------------------
    def poll_maps(self):
        """
        Read all relevant maps and invoke callback if set
        """
        maps_to_poll = {
            "session_map": SESSION_MAP,
            "ul_map": UL_MAP,
            "dl_map": DL_MAP,
            "stats_map": STATS_MAP
        }

        for name, path in maps_to_poll.items():
            try:
                entries = read_bpf_map(path)
                logging.debug(f"[Observer] {name} entries: {entries}")

                if self.callback:
                    self.callback(name, entries)

            except Exception as e:
                logging.error(f"[Observer] Failed to read map {name}: {e}")


# ------------------------------------------------------------
# Example Usage
# ------------------------------------------------------------
def example_callback(map_name, entries):
    """
    Example callback function invoked when map is polled
    """
    logging.info(f"[Callback] Map {map_name} has {len(entries)} entries")


if __name__ == "__main__":
    observer = EBPFObserver(poll_interval=5, callback=example_callback)
    observer.start()

    try:
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        observer.stop()
        logging.info("eBPF observer stopped.")
