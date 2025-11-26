#ifndef BPF_LOADER_H
#define BPF_LOADER_H

#include <linux/if_link.h>   // For TC hook constants
#include <bpf/libbpf.h>

// Structure to hold loader context
struct bpf_loader {
    struct bpf_object *obj;
    int ifindex_ingress;
    int ifindex_egress;
};

// Initialize the loader context
int bpf_loader_init(struct bpf_loader *loader);

// Load a .o eBPF file
int bpf_loader_load(struct bpf_loader *loader, const char *filename);

// Attach the loaded eBPF program to TC hooks
int bpf_loader_attach_tc(struct bpf_loader *loader, const char *ifname);

// Detach eBPF program from TC hooks
int bpf_loader_detach_tc(struct bpf_loader *loader, const char *ifname);

// Cleanup loader context
void bpf_loader_cleanup(struct bpf_loader *loader);

#endif // BPF_LOADER_H
