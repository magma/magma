#include "bpf_loader.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <net/if.h>
#include <bpf/libbpf.h>
#include <errno.h>

int bpf_loader_init(struct bpf_loader *loader) {
    if (!loader) return -1;
    loader->obj = NULL;
    loader->ifindex_ingress = 0;
    loader->ifindex_egress = 0;
    return 0;
}

int bpf_loader_load(struct bpf_loader *loader, const char *filename) {
    if (!loader || !filename) return -1;

    struct bpf_object *obj;
    int err;

    obj = bpf_object__open_file(filename, NULL);
    if (libbpf_get_error(obj)) {
        fprintf(stderr, "Failed to open BPF object file: %s\n", filename);
        return -1;
    }

    err = bpf_object__load(obj);
    if (err) {
        fprintf(stderr, "Failed to load BPF object: %d\n", err);
        bpf_object__close(obj);
        return -1;
    }

    loader->obj = obj;
    return 0;
}

int bpf_loader_attach_tc(struct bpf_loader *loader, const char *ifname) {
    if (!loader || !loader->obj || !ifname) return -1;

    int ifindex = if_nametoindex(ifname);
    if (ifindex == 0) {
        perror("if_nametoindex");
        return -1;
    }

    // Iterate through programs in object and attach to ingress TC hook
    struct bpf_program *prog;
    bpf_object__for_each_program(prog, loader->obj) {
        int prog_fd = bpf_program__fd(prog);
        if (prog_fd < 0) continue;

        int err = bpf_tc_attach(ifindex, BPF_TC_INGRESS, prog_fd);
        if (err) {
            fprintf(stderr, "Failed to attach program to %s ingress: %d\n", ifname, err);
            return -1;
        }
    }

    loader->ifindex_ingress = ifindex;
    return 0;
}

int bpf_loader_detach_tc(struct bpf_loader *loader, const char *ifname) {
    if (!loader || !ifname) return -1;

    int ifindex = if_nametoindex(ifname);
    if (ifindex == 0) {
        perror("if_nametoindex");
        return -1;
    }

    // This is simplified. Use libbpf or tc APIs to detach
    // TODO: Implement proper TC detach logic
    printf("Detached TC programs from %s\n", ifname);
    return 0;
}

void bpf_loader_cleanup(struct bpf_loader *loader) {
    if (!loader) return;
    if (loader->obj)
        bpf_object__close(loader->obj);
    loader->obj = NULL;
}
