#include "bpf_utils.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <errno.h>
#include <sys/stat.h>
#include <sys/mount.h>

#include <linux/pkt_cls.h>
#include <bpf/libbpf.h>
#include <bpf/bpf.h>
#include <net/if.h>

// ------------------------------
// Helper: get interface index
// ------------------------------
int bpf_utils_ifindex(const char *ifname) {
    int idx = if_nametoindex(ifname);
    if (!idx)
        fprintf(stderr, "Error: interface %s not found\n", ifname);
    return idx;
}

// ------------------------------
// Ensure bpffs mounted
// ------------------------------
int bpf_utils_mount_bpffs(const char *path) {
    struct stat st = {};

    if (stat(path, &st) == 0)
        return 0;

    if (mkdir(path, 0755) && errno != EEXIST)
        return -errno;

    if (mount("bpffs", path, "bpf", 0, NULL) < 0 && errno != EBUSY) {
        perror("mount bpffs");
        return -errno;
    }

    return 0;
}

// ------------------------------
// Pin / Unpin Maps
// ------------------------------
int bpf_utils_pin_map(int map_fd, const char *path) {
    if (bpf_obj_pin(map_fd, path)) {
        perror("bpf_obj_pin");
        return -errno;
    }
    return 0;
}

int bpf_utils_unpin_map(const char *path) {
    if (unlink(path) < 0 && errno != ENOENT) {
        perror("unlink");
        return -errno;
    }
    return 0;
}

// ------------------------------
// Create/Delete clsact Qdisc
// ------------------------------
static int run_tc(const char *cmd) {
    printf("[TC] %s\n", cmd);
    return system(cmd);
}

int bpf_utils_create_clsact(const char *ifname) {
    char cmd[256];
    snprintf(cmd, sizeof(cmd),
             "tc qdisc replace dev %s clsact", ifname);
    return run_tc(cmd);
}

int bpf_utils_delete_clsact(const char *ifname) {
    char cmd[256];
    snprintf(cmd, sizeof(cmd),
             "tc qdisc del dev %s clsact", ifname);
    return run_tc(cmd);
}

// ------------------------------
// Attach Ingress / Egress
// ------------------------------
int bpf_utils_tc_attach_ingress(int ifindex, int prog_fd, const char *section) {
    struct bpf_tc_hook hook = {
        .ifindex = ifindex,
        .attach_point = BPF_TC_INGRESS,
    };

    struct bpf_tc_opts opts = {
        .prog_fd = prog_fd,
        .prog_id = 0,
        .attach_point = BPF_TC_INGRESS,
        .flags = BPF_TC_F_REPLACE,
    };

    int err = bpf_tc_hook_create(&hook);
    if (err && err != -EEXIST)
        return err;

    return bpf_tc_attach(&hook, &opts);
}

int bpf_utils_tc_attach_egress(int ifindex, int prog_fd, const char *section) {
    struct bpf_tc_hook hook = {
        .ifindex = ifindex,
        .attach_point = BPF_TC_EGRESS,
    };

    struct bpf_tc_opts opts = {
        .prog_fd = prog_fd,
        .attach_point = BPF_TC_EGRESS,
        .flags = BPF_TC_F_REPLACE,
    };

    int err = bpf_tc_hook_create(&hook);
    if (err && err != -EEXIST)
        return err;

    return bpf_tc_attach(&hook, &opts);
}

// ------------------------------
// Detach programs
// ------------------------------
int bpf_utils_tc_detach(int ifindex) {
    struct bpf_tc_hook hook = {
        .ifindex = ifindex,
        .attach_point = BPF_TC_INGRESS,
    };
    bpf_tc_detach(&hook, NULL);

    hook.attach_point = BPF_TC_EGRESS;
    bpf_tc_detach(&hook, NULL);

    return 0;
}
