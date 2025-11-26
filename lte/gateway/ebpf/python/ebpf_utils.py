
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


def write_bpf_map(map_path: str, key: str, value: dict):
    """
    Write an entry to a BPF map
    """
    # This is a placeholder; in practice, use bpftool or libbpf APIs
    logging.info("Writing to BPF map %s: %s -> %s", map_path, key, value)
    # TODO: Implement actual write logic


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
