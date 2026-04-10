#!/usr/bin/env python3


import logging
from ebpf_utils import write_bpf_map, delete_bpf_map_entry, read_bpf_map
from ebpf_gtp_manager import GtpSession

logging.basicConfig(
    level=logging.INFO,
    format='[%(asctime)s] %(levelname)s: %(message)s'
)

# Pinned map locations
UL_MAP = "/sys/fs/bpf/gtp_ul_map"
DL_MAP = "/sys/fs/bpf/gtp_dl_map"
SESS_MAP = "/sys/fs/bpf/gtp_session_map"


# -------------------------------------------------------------------
# PipelineD → eBPF: Session Creation
# -------------------------------------------------------------------

def pipelined_create_session(imsi: str, ue_ip: str, ul_teid: int, dl_teid: int):
    """
    Write UL/DL TEIDs + session metadata to eBPF kernel maps.

    Called when PipelineD installs a new bearer/session.
    """
    logging.info(f"[PipelineD → eBPF] Creating session: IMSI={imsi}, UE={ue_ip}, "
                 f"UL_TEID={ul_teid}, DL_TEID={dl_teid}")

    # Session metadata for gtp_session_map
    session_value = {
        "imsi": imsi,
        "ue_ip": ue_ip,
        "ul_teid": ul_teid,
        "dl_teid": dl_teid,
        "state": 1,
    }

    # Write UL entry
    write_bpf_map(UL_MAP, key=str(ul_teid), value={
        "imsi": imsi,
        "ue_ip": ue_ip
    })

    # Write DL entry
    write_bpf_map(DL_MAP, key=str(dl_teid), value={
        "imsi": imsi,
        "ue_ip": ue_ip
    })

    # Write full session entry
    write_bpf_map(SESS_MAP, key=str(dl_teid), value=session_value)

    return True


# -------------------------------------------------------------------
# PipelineD → eBPF: Session Deletion
# -------------------------------------------------------------------

def pipelined_delete_session(imsi: str, dl_teid: int, ul_teid: int):
    """
    Remove UL, DL and session metadata entries from eBPF maps.
    """
    logging.info(f"[PipelineD → eBPF] Deleting session: IMSI={imsi}, "
                 f"UL_TEID={ul_teid}, DL_TEID={dl_teid}")

    delete_bpf_map_entry(UL_MAP, str(ul_teid))
    delete_bpf_map_entry(DL_MAP, str(dl_teid))
    delete_bpf_map_entry(SESS_MAP, str(dl_teid))

    return True


# -------------------------------------------------------------------
# PipelineD → eBPF: Update (Optional)
# -------------------------------------------------------------------

def pipelined_update_session(imsi: str, ue_ip: str = None, state: int = None):
    """
    Update existing session metadata (e.g., state, IP change).

    PipelineD uses this when bearer transitions or lifecycle changes occur.
    """
    logging.info(f"[PipelineD → eBPF] Updating session IMSI={imsi}")

    # Read entire session map
    session_entries = read_bpf_map(SESS_MAP)

    # Locate entry by IMSI
    target = None
    for entry in session_entries:
        if entry.get("imsi") == imsi:
            target = entry
            break

    if not target:
        logging.warning(f"Session for IMSI {imsi} not found")
        return False

    # Modify fields
    if ue_ip:
        target["ue_ip"] = ue_ip
    if state:
        target["state"] = state

    # Write back
    dl_teid = target.get("dl_teid")
    write_bpf_map(SESS_MAP, key=str(dl_teid), value=target)

    return True


# -------------------------------------------------------------------
# PipelineD Stats → eBPF
# -------------------------------------------------------------------

def pipelined_get_stats():
    """
    Return eBPF byte/packet counters (if your maps include stats).

    Called when PipelineD queries flow stats.
    """
    logging.info("[PipelineD → eBPF] Stats requested")

    # Example: reading UL map stats
    ul_entries = read_bpf_map(UL_MAP)
    dl_entries = read_bpf_map(DL_MAP)

    return {
        "uplink": ul_entries,
        "downlink": dl_entries,
    }


# -------------------------------------------------------------------
# PipelineD Startup Sync → eBPF
# -------------------------------------------------------------------

def pipelined_initial_sync(pipelined_sessions):
    """
    Called at PipelineD startup to push all existing session rules into eBPF.

    pipelined_sessions: list of dicts:
        {
            "imsi": "...",
            "ue_ip": "...",
            "ul_teid": int,
            "dl_teid": int,
        }
    """
    logging.info("[PipelineD → eBPF] Performing initial sync...")

    for session in pipelined_sessions:
        pipelined_create_session(
            session["imsi"],
            session["ue_ip"],
            session["ul_teid"],
            session["dl_teid"],
        )

    logging.info(f"[PipelineD → eBPF] Initial sync complete ({len(pipelined_sessions)} sessions).")


# -------------------------------------------------------------------
# End of File
# -------------------------------------------------------------------

if __name__ == "__main__":
    logging.info("eBPF PipelineD integration module loaded.")
