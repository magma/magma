

import logging
import time
from ebpf_utils import read_bpf_map, write_bpf_map, delete_bpf_map_entry

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='[%(asctime)s] %(levelname)s: %(message)s'
)

class GtpSession:
    """Represents a GTP session entry"""
    def __init__(self, imsi, teid, ue_ip, state):
        self.imsi = imsi
        self.teid = teid
        self.ue_ip = ue_ip
        self.state = state

    def __repr__(self):
        return f"GtpSession(imsi={self.imsi}, teid={self.teid}, ue_ip={self.ue_ip}, state={self.state})"


def get_kernel_gtp_sessions(map_name="/sys/fs/bpf/gtp_session_map"):
    """
    Read GTP sessions from pinned BPF map in kernel
    Returns: list of GtpSession
    """
    sessions = []
    kernel_entries = read_bpf_map(map_name)
    for entry in kernel_entries:
        # Example entry: {'imsi': '001010123456789', 'teid': 12345, 'ue_ip': '10.0.0.5', 'state': 1}
        session = GtpSession(
            imsi=entry['imsi'],
            teid=entry['teid'],
            ue_ip=entry['ue_ip'],
            state=entry['state']
        )
        sessions.append(session)
    return sessions


def sync_gtp_sessions():
    """
    Synchronize GTP sessions between kernel BPF map and PipelineD userspace
    """
    logging.info("Syncing GTP sessions")
    kernel_sessions = get_kernel_gtp_sessions()

    # Placeholder: get userspace sessions (from DB or PipelineD)
    userspace_sessions = get_userspace_sessions()

    # Sync kernel -> userspace
    for ksession in kernel_sessions:
        if ksession.imsi not in userspace_sessions:
            logging.info("Adding new session to userspace: %s", ksession)
            add_userspace_session(ksession)
        else:
            # Optional: update state if changed
            pass

    # Sync userspace -> kernel (remove stale entries)
    kernel_imsi_set = {s.imsi for s in kernel_sessions}
    for usession in userspace_sessions.values():
        if usession.imsi not in kernel_imsi_set:
            logging.info("Removing stale kernel session: %s", usession)
            delete_bpf_map_entry("/sys/fs/bpf/gtp_session_map", usession.teid)


def get_userspace_sessions():
    """
    Retrieve GTP sessions from userspace database / PipelineD
    Returns: dict of imsi -> GtpSession
    """
    # Placeholder for real DB/API call
    return {}


def add_userspace_session(session: GtpSession):
    """
    Add a new session to userspace
    """
    # Placeholder for real DB/API call
    logging.info("Session added to userspace: %s", session)


if __name__ == "__main__":
    while True:
        try:
            sync_gtp_sessions()
            time.sleep(10)  # configurable
        except KeyboardInterrupt:
            logging.info("GTP session manager stopped")
            break
