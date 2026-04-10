#ifndef BPF_UTILS_H
#define BPF_UTILS_H

#include <linux/if_link.h>
#include <bpf/libbpf.h>

#ifdef __cplusplus
extern "C" {
#endif

// Create clsact qdisc on interface
int bpf_utils_create_clsact(const char *ifname);

// Delete clsact qdisc
int bpf_utils_delete_clsact(const char *ifname);

// Returns interface index
int bpf_utils_ifindex(const char *ifname);

// Attach BPF program to TC ingress
int bpf_utils_tc_attach_ingress(int ifindex,
                                int prog_fd,
                                const char *section);

// Attach BPF program to TC egress
int bpf_utils_tc_attach_egress(int ifindex,
                               int prog_fd,
                               const char *section);

// Detach TC program (ingress + egress)
int bpf_utils_tc_detach(int ifindex);

// Pin a BPF map to bpffs
int bpf_utils_pin_map(int map_fd, const char *path);

// Unpin map
int bpf_utils_unpin_map(const char *path);

// Ensure bpffs is mounted
int bpf_utils_mount_bpffs(const char *path);

#ifdef __cplusplus
}
#endif

#endif // BPF_UTILS_H
