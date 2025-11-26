"""
ebpf_sessiond_integration.py

SessionD → eBPF integration layer.
This module listens for SessionD events (install base names, bearer setup,
QoS changes, and session termination) and synchronizes them with BPF maps.

Maps updated:
- session_map (key: ue_ipv4, value: session metadata)
- qos_map (key: session_id, value: QoS rules)
- stats_map (UL/DL byte counters - read-only from eBPF)
"""

import json
import logging
import threading
import time

from ebpf_utils import (
    load_bpf_map,
    bpf_map_update,
    bpf_map_delete,
    map_key_ipv4,
)

SESSION_MAP = "/sys/fs/bpf/magma/session_map"
QOS_MAP = "/sys/fs/bpf/magma/qos_map"
STATS_MAP = "/sys/fs/bpf/magma/stats_map"


class EBPFSessiondIntegration:
    """
    High-level session synchronization class.
    Receives session install/remove events from SessionD via RPC or local file/IP socket.
    """

    def __init__(self):
        self.session_map = None
        self.qos_map = None
        self.stats_map = None
        self.running = False

    # ------------------------------------------------------------
    # Initialization
    # ------------------------------------------------------------
    def initialize(self):
        """Load pinned maps created by the BPF loader."""
        logging.info("[SessionD] Loading pinned BPF maps…")

        self.session_map = load_bpf_map(SESSION_MAP)
        self.qos_map = load_bpf_map(QOS_MAP)
        self.stats_map = load_bpf_map(STATS_MAP)

        if not self.session_map or not self.qos_map:
            raise RuntimeError("BPF map load failed — SessionD integration cannot start.")

        logging.info("[SessionD] Maps loaded successfully.")

    # ------------------------------------------------------------
    # Session Install / Remove / Update
    # ------------------------------------------------------------
    def install_session(self, session):
        """
        Install a new user session into BPF maps.

        session = {
            "imsi": "001010000000001",
            "ue_ipv4": "192.168.128.10",
            "enb_teid": 0x12345,
            "ags_teid": 0x54321,
            "qci": 9,
            "ambr_ul": 200000,  # kbps
            "ambr_dl": 500000,
        }
        """

        ue_key = map_key_ipv4(session["ue_ipv4"])

        session_val = {
            "enb_teid": session["enb_teid"],
            "ags_teid": session["ags_teid"],
            "qci": session["qci"],
            "ambr_ul": session["ambr_ul"],
            "ambr_dl": session["ambr_dl"],
        }

        logging.info(f"[SessionD] Installing session for {session['ue_ipv4']}")

        bpf_map_update(self.session_map, ue_key, session_val)

        # QoS map entry
        qos_val = {
            "max_ul_bps": session["ambr_ul"] * 1000,
            "max_dl_bps": session["ambr_dl"] * 1000,
            "qci": session["qci"],
        }

        bpf_map_update(self.qos_map, ue_key, qos_val)

        logging.info("[SessionD] Session installed successfully.")

    def remove_session(self, ue_ipv4):
        """Remove a user session from BPF maps."""
        ue_key = map_key_ipv4(ue_ipv4)

        logging.info(f"[SessionD] Removing session for {ue_ipv4}")

        bpf_map_delete(self.session_map, ue_key)
        bpf_map_delete(self.qos_map, ue_key)

        logging.info("[SessionD] Session removed.")

    def update_qos(self, ue_ipv4, qos_params):
        """
        Update QoS parameters for an active session.
        qos_params = { "ambr_ul": 300000, "ambr_dl": 700000, "qci": 7 }
        """
        ue_key = map_key_ipv4(ue_ipv4)

        qos_val = {
            "max_ul_bps": qos_params["ambr_ul"] * 1000,
            "max_dl_bps": qos_params["ambr_dl"] * 1000,
            "qci": qos_params["qci"],
        }

        logging.info(f"[SessionD] Updating QoS for {ue_ipv4}")

        bpf_map_update(self.qos_map, ue_key, qos_val)

        logging.info("[SessionD] QoS updated.")

    # ------------------------------------------------------------
    # Stats Polling (Optional)
    # ------------------------------------------------------------
    def get_session_stats(self, ue_ipv4):
        """Read UL/DL stats from stats_map for Prometheus or Enforcement."""
        ue_key = map_key_ipv4(ue_ipv4)

        try:
            stats = self.stats_map[ue_key]
            return {
                "ul_bytes": stats["ul_bytes"],
                "dl_bytes": stats["dl_bytes"],
            }
        except KeyError:
            return {"ul_bytes": 0, "dl_bytes": 0}

    # ------------------------------------------------------------
    # Background Daemon Loop (optional)
    # ------------------------------------------------------------
    def start_event_listener(self, json_events_path="/var/run/magma/sessiond_events.json"):
        """
        Simple file-based watcher for new events.
        Real version should use gRPC from SessionD.
        """
        logging.info("[SessionD] Starting event listener thread…")
        self.running = True

        th = threading.Thread(target=self._event_loop, args=(json_events_path,))
        th.daemon = True
        th.start()

    def _event_loop(self, json_events_path):
        logging.info("[SessionD] Listener thread ready.")
        last_size = 0

        while self.running:
            try:
                size = os.path.getsize(json_events_path)
                if size > last_size:
                    with open(json_events_path) as f:
                        for line in f:
                            evt = json.loads(line)
                            self.handle_event(evt)
                    last_size = size
            except Exception as e:
                logging.error(f"[SessionD] Event loop error: {e}")

            time.sleep(1)

    # ------------------------------------------------------------
    # Event Handler
    # ------------------------------------------------------------
    def handle_event(self, evt):
        """
        Handle incoming events from SessionD.
        evt = { "type": "INSTALL", "session": {...} }
        """
        etype = evt.get("type")

        if etype == "INSTALL":
            self.install_session(evt["session"])
        elif etype == "REMOVE":
            self.remove_session(evt["ue_ipv4"])
        elif etype == "QOS_UPDATE":
            self.update_qos(evt["ue_ipv4"], evt["qos"])
        else:
            logging.warning(f"[SessionD] Unknown event type: {etype}")

