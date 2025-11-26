
import os
import subprocess
import time
import logging
from ebpf_utils import run_command, check_tc_attached
from ebpf_gtp_manager import sync_gtp_sessions

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='[%(asctime)s] %(levelname)s: %(message)s'
)

class EbpfManager:
    def __init__(self, config_path="/etc/magma/magma_bpf.conf"):
        self.config_path = config_path
        self.load_config()

    def load_config(self):
        """Load configuration for eBPF (ports, maps, XDP/TC flags)"""
        # Example config parsing
        self.ovs_ports = ["br-magma", "gtp0", "sgi0"]
        self.enable_tc = True
        self.enable_xdp = False
        logging.info("Configuration loaded from %s", self.config_path)

    def load_bpf_programs(self):
        """Attach TC or XDP eBPF programs to OVS ports"""
        for port in self.ovs_ports:
            if self.enable_tc:
                logging.info("Attaching TC eBPF to port %s", port)
                run_command(f"./scripts/load_tc.sh {port}")
            elif self.enable_xdp:
                logging.info("Attaching XDP eBPF to port %s", port)
                run_command(f"./scripts/load_xdp.sh {port}")

    def unload_bpf_programs(self):
        """Detach eBPF programs"""
        for port in self.ovs_ports:
            logging.info("Detaching eBPF from port %s", port)
            run_command(f"./scripts/unload_tc.sh {port}")

    def monitor_events(self):
        """Monitor kernel perf buffer events"""
        logging.info("Starting perf buffer event monitoring")
        # Placeholder loop for event reading
        try:
            while True:
                # In reality, this would attach to perf buffer using bcc/libbpf
                time.sleep(5)
                logging.debug("Perf buffer poll cycle")
        except KeyboardInterrupt:
            logging.info("Event monitoring stopped")

    def sync_sessions_loop(self):
        """Continuously sync GTP sessions with Magma"""
        logging.info("Starting GTP session sync loop")
        try:
            while True:
                sync_gtp_sessions()
                time.sleep(10)  # configurable
        except KeyboardInterrupt:
            logging.info("Session sync loop stopped")

    def start(self):
        """Start the manager"""
        logging.info("Starting eBPF manager")
        self.load_bpf_programs()
        # Optionally run event monitoring and session sync in parallel threads
        self.sync_sessions_loop()

    def stop(self):
        """Stop manager and cleanup"""
        logging.info("Stopping eBPF manager")
        self.unload_bpf_programs()


if __name__ == "__main__":
    manager = EbpfManager()
    try:
        manager.start()
    except Exception as e:
        logging.error("Manager crashed: %s", e)
        manager.stop()
