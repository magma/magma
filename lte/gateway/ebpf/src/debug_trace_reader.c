// src/debug_trace_reader.c
#define _GNU_SOURCE
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <fcntl.h>
#include <errno.h>
#include <poll.h>
#include <signal.h>
#include <bpf/libbpf.h>
#include <bpf/bpf.h>

static volatile bool exiting = false;

static void handle_sig(int sig) { exiting = true; }

/* Simple tail -f on trace_pipe */
static int read_trace_pipe() {
    const char *trace_path = "/sys/kernel/debug/tracing/trace_pipe";
    int fd = open(trace_path, O_RDONLY | O_NONBLOCK);
    if (fd < 0) {
        perror("open trace_pipe");
        return -1;
    }

    struct pollfd pfd = { .fd = fd, .events = POLLIN };
    char buf[4096];

    printf("Reading kernel trace from %s (Ctrl+C to stop)\n", trace_path);

    while (!exiting) {
        int ret = poll(&pfd, 1, 500);
        if (ret < 0) {
            if (errno == EINTR) continue;
            perror("poll");
            break;
        } else if (ret == 0) {
            continue; // timeout, loop
        }

        if (pfd.revents & POLLIN) {
            ssize_t n = read(fd, buf, sizeof(buf)-1);
            if (n > 0) {
                buf[n] = '\0';
                printf("%s", buf); // trace_pipe returns newline-terminated lines
                fflush(stdout);
            } else if (n == 0) {
                // EOF unlikely for trace_pipe
                usleep(100000);
            } else if (errno != EAGAIN) {
                perror("read");
                break;
            }
        }
    }

    close(fd);
    return 0;
}

/* Alternative: read events from a pinned perf map fd (like in event_handler) */
static int read_perf_map(const char *pinned_map_path) {
    int map_fd = bpf_obj_get(pinned_map_path);
    if (map_fd < 0) {
        fprintf(stderr, "bpf_obj_get(%s) failed: %s\n", pinned_map_path, strerror(errno));
        return -1;
    }

    struct perf_buffer *pb = NULL;
    struct perf_buffer_opts opts = {};

    // Minimal callbacks
    opts.sample_cb = [](void *ctx, void *data, size_t size) -> int {
        // data is opaque; print raw hex (user should cast to correct struct)
        const unsigned char *d = data;
        printf("Perf sample (%zu bytes): ", size);
        for (size_t i = 0; i < size; ++i) printf("%02x", d[i]);
        printf("\n");
        return 0;
    };
    opts.lost_cb = [](void *ctx, int lost) {
        fprintf(stderr, "Lost %d events\n", lost);
    };

    // NOTE: libbpf's perf_buffer__new expects C-style function pointers.
    // Use wrapper functions or compile with C++17 lambdas with extern "C" if needed.
    // For simplicity we won't implement perf_map read here in full.
    // The event_handler.c is a more complete example for perf map reading.

    close(map_fd);
    fprintf(stderr, "read_perf_map: not fully implemented here, use event_handler.c for perf maps\n");
    return -1;
}

int main(int argc, char **argv) {
    signal(SIGINT, handle_sig);
    signal(SIGTERM, handle_sig);

    if (argc > 1 && strcmp(argv[1], "--perf") == 0 && argc > 2) {
        return read_perf_map(argv[2]);
    }

    return read_trace_pipe();
}
