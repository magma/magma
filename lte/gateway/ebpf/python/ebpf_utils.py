
import subprocess
import logging
import os
import json

logging.basicConfig(
    level=logging.INFO,
    format='[%(asctime)s] %(levelname)s: %(message)s'
)



def run_command(cmd: str, check=True) -> str:
    """
    Run a shell command and return output
    """
    logging.debug("Running command: %s", cmd)
    try:
        result = subprocess.run(
            cmd, shell=True, text=True, capture_output=True, check=check
        )
        return result.stdout.strip()
    except subprocess.CalledProcessError as e:
        logging.error("Command failed: %s\n%s", cmd, e.stderr)
        if check:
            raise e
        return e.stderr.strip()




def read_bpf_map(map_path: str):
    """
    Read a pinned BPF map using bpftool
    Returns a list of dict entries
    """
    if not os.path.exists(map_path):
        logging.warning("BPF map not found: %s", map_path)
        return []

    cmd = f"bpftool map dump pinned {map_path} -j"
    output = run_command(cmd)
    try:
        entries = json.loads(output)
        return entries
    except json.JSONDecodeError:
        logging.error("Failed to decode BPF map output for %s", map_path)
        return []


def write_bpf_map(map_path: str, key: bytes, value: bytes):
    """
    Write or update an entry in a pinned BPF map using bpftool.

    key   : raw bytes (already packed)
    value : raw bytes (already packed)
    """
    if not os.path.exists(map_path):
        raise FileNotFoundError(f"BPF map not found: {map_path}")

    # Convert bytes to hex format expected by bpftool
    key_hex = " ".join(f"{b:02x}" for b in key)
    value_hex = " ".join(f"{b:02x}" for b in value)

    cmd = (
        f"bpftool map update pinned {map_path} "
        f"key hex {key_hex} "
        f"value hex {value_hex}"
    )

    logging.debug("Executing: %s", cmd)
    run_command(cmd)


def delete_bpf_map_entry(map_path: str, key: str):
    """
    Delete an entry from a BPF map
    """
    logging.info("Deleting key %s from BPF map %s", key, map_path)
    cmd = f"bpftool map delete pinned {map_path} key {key}"
    run_command(cmd, check=False)



def check_tc_attached(iface: str) -> bool:
    """
    Check if a TC clsact is attached on the interface
    """
    cmd = f"tc qdisc show dev {iface}"
    output = run_command(cmd, check=False)
    return "clsact" in output


def mount_bpf_fs(mount_point="/sys/fs/bpf"):
    """
    Mount BPF filesystem if not mounted
    """
    if not os.path.exists(mount_point):
        os.makedirs(mount_point)
    cmd = f"mountpoint -q {mount_point} || mount -t bpf bpf {mount_point}"
    run_command(cmd)


def unmount_bpf_fs(mount_point="/sys/fs/bpf"):
    """
    Unmount BPF filesystem
    """
    cmd = f"umount {mount_point}"
    run_command(cmd, check=False)


if __name__ == "__main__":
    mount_bpf_fs()
    logging.info("BPF utils ready")
