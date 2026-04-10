#include <bpf/libbpf.h>
#include <bpf/bpf.h>
#include <linux/perf_event.h>
#include <signal.h>
#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

static volatile bool exiting = false;

struct event {
    __u32 session_id;
    __u32 ul_bytes;
    __u32 dl_bytes;
    char info[64];
};

// Signal handler for Ctrl+C
static void sig_handler(int sig) {
    exiting = true;
}

// Callback for perf buffer events
static int handle_event(void *ctx, void *data, size_t data_sz) {
    struct event *e = data;
    printf("Event: session_id=%u, UL=%u bytes, DL=%u bytes, info=%s\n",
           e->session_id, e->ul_bytes, e->dl_bytes, e->info);
    return 0;
}

// Callback for lost events
static void handle_lost_events(void *ctx, int lost) {
    fprintf(stderr, "Lost %d events\n", lost);
}

int main(int argc, char **argv) {
    struct perf_buffer *pb = NULL;
    struct bpf_map *events_map;
    int map_fd;
    int err;

    signal(SIGINT, sig_handler);
    signal(SIGTERM, sig_handler);

    // Open BPF object file
    struct bpf_object *obj = bpf_object__open_file("tc_core.o", NULL);
    if (libbpf_get_error(obj)) {
        fprintf(stderr, "Failed to open BPF object\n");
        return 1;
    }

    if ((err = bpf_object__load(obj))) {
        fprintf(stderr, "Failed to load BPF object\n");
        return 1;
    }

    // Open the perf event map created in BPF program
    map_fd = bpf_object__find_map_fd_by_name(obj, "events");
    if (map_fd < 0) {
        fprintf(stderr, "Failed to find events map\n");
        return 1;
    }

    // Setup perf buffer
    struct perf_buffer_opts pb_opts = {};
    pb_opts.sample_cb = handle_event;
    pb_opts.lost_cb = handle_lost_events;

    pb = perf_buffer__new(map_fd, 8, &pb_opts);
    if (libbpf_get_error(pb)) {
        fprintf(stderr, "Failed to open perf buffer\n");
        return 1;
    }

    printf("Listening for events... Press Ctrl+C to exit.\n");

    while (!exiting) {
        err = perf_buffer__poll(pb, 100 /* timeout ms */);
        if (err < 0 && err != -EINTR) {
            fprintf(stderr, "Error polling perf buffer: %d\n", err);
        }
    }

    perf_buffer__free(pb);
    bpf_object__close(obj);

    printf("Exiting...\n");
    return 0;
}
