// src/sessiond_pfcp_adapter.c
#define _GNU_SOURCE
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <errno.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <stdint.h>
#include <signal.h>

#include <bpf/libbpf.h>
#include <bpf/bpf.h>

#include "bpf_utils.h"

static volatile bool exiting = false;
static const char *SOCKET_PATH = "/var/run/magma/sessiond_pfcp.sock";
static const char *PINNED_MAP_PATH = "/sys/fs/bpf/magma_sessions"; // expected pinned map path

static void sigint(int sig) { exiting = true; }

/*
 * NOTE: the session map key/value must match your BPF program.
 * This is an example layout:
 */
typedef struct {
    __u32 seid;      // session id
} map_key_t;

typedef struct {
    __u32 ipv4;      // network order IPv4 addr
    __u32 teid;      // GTP TEID
    __u32 ifindex;   // outgoing interface index
} map_val_t;

/* Helper: open pinned map */
static int open_pinned_map(const char *path) {
    int fd = bpf_obj_get(path);
    if (fd < 0) {
        fprintf(stderr, "bpf_obj_get(%s) failed: %s\n", path, strerror(errno));
    }
    return fd;
}

/* Parse an IPv4 dotted string to uint32 (network order) */
static int parse_ipv4(const char *s, __u32 *out) {
    struct in_addr a;
    if (inet_pton(AF_INET, s, &a) != 1) return -1;
    *out = a.s_addr;
    return 0;
}

/* Update session map */
static int session_map_add(int map_fd, __u32 seid, const char *ipv4_s, __u32 teid, const char *ifname) {
    map_key_t key = { .seid = seid };
    map_val_t val;
    memset(&val, 0, sizeof(val));

    if (parse_ipv4(ipv4_s, &val.ipv4) != 0) {
        fprintf(stderr, "Invalid IPv4: %s\n", ipv4_s);
        return -1;
    }
    val.teid = teid;
    int ifidx = bpf_utils_ifindex(ifname);
    if (ifidx <= 0) {
        fprintf(stderr, "Invalid interface: %s\n", ifname);
        return -1;
    }
    val.ifindex = (__u32)ifidx;

    if (bpf_map_update_elem(map_fd, &key, &val, BPF_ANY) != 0) {
        fprintf(stderr, "bpf_map_update_elem failed: %s\n", strerror(errno));
        return -1;
    }
    return 0;
}

/* Delete session map entry */
static int session_map_del(int map_fd, __u32 seid) {
    map_key_t key = { .seid = seid };
    if (bpf_map_delete_elem(map_fd, &key) != 0) {
        fprintf(stderr, "bpf_map_delete_elem failed: %s\n", strerror(errno));
        return -1;
    }
    return 0;
}

/* Minimal unix socket server to accept simple textual PFCP commands */
int main(int argc, char **argv) {
    int server_fd, client_fd;
    struct sockaddr_un addr;
    char buf[256];
    int map_fd = -1;

    signal(SIGINT, sigint);
    signal(SIGTERM, sigint);

    /* open pinned map */
    map_fd = open_pinned_map(PINNED_MAP_PATH);
    if (map_fd < 0) {
        fprintf(stderr, "Failed to open pinned map at %s\n", PINNED_MAP_PATH);
        /* still run if you want to debug via socket, but map ops will fail */
    }

    /* create socket */
    if ((server_fd = socket(AF_UNIX, SOCK_STREAM, 0)) < 0) {
        perror("socket");
        return 1;
    }

    unlink(SOCKET_PATH);
    memset(&addr, 0, sizeof(addr));
    addr.sun_family = AF_UNIX;
    strncpy(addr.sun_path, SOCKET_PATH, sizeof(addr.sun_path) - 1);

    if (bind(server_fd, (struct sockaddr*)&addr, sizeof(addr)) < 0) {
        perror("bind");
        close(server_fd);
        return 1;
    }

    if (listen(server_fd, 5) < 0) {
        perror("listen");
        close(server_fd);
        return 1;
    }

    printf("sessiond_pfcp_adapter listening on %s\n", SOCKET_PATH);

    while (!exiting) {
        client_fd = accept(server_fd, NULL, NULL);
        if (client_fd < 0) {
            if (errno == EINTR) continue;
            perror("accept");
            break;
        }

        ssize_t n = read(client_fd, buf, sizeof(buf)-1);
        if (n <= 0) {
            close(client_fd);
            continue;
        }
        buf[n] = '\0';

        /* Expect commands of the form:
         *   ADD <seid> <ipv4> <teid> <ifname>\n
         *   DEL <seid>\n
         */
        char *save = NULL;
        char *tok = strtok_r(buf, " \t\r\n", &save);
        if (!tok) { close(client_fd); continue; }

        if (strcmp(tok, "ADD") == 0) {
            char *seid_s = strtok_r(NULL, " \t\r\n", &save);
            char *ipv4 = strtok_r(NULL, " \t\r\n", &save);
            char *teid_s = strtok_r(NULL, " \t\r\n", &save);
            char *ifname = strtok_r(NULL, " \t\r\n", &save);

            if (!seid_s || !ipv4 || !teid_s || !ifname) {
                dprintf(client_fd, "ERR missing args\n");
            } else {
                __u32 seid = (uint32_t)strtoul(seid_s, NULL, 0);
                __u32 teid = (uint32_t)strtoul(teid_s, NULL, 0);
                if (map_fd >= 0) {
                    if (session_map_add(map_fd, seid, ipv4, teid, ifname) == 0)
                        dprintf(client_fd, "OK\n");
                    else
                        dprintf(client_fd, "ERR map update failed\n");
                } else {
                    dprintf(client_fd, "ERR no map\n");
                }
            }
        } else if (strcmp(tok, "DEL") == 0) {
            char *seid_s = strtok_r(NULL, " \t\r\n", &save);
            if (!seid_s) {
                dprintf(client_fd, "ERR missing seid\n");
            } else {
                __u32 seid = (uint32_t)strtoul(seid_s, NULL, 0);
                if (map_fd >= 0) {
                    if (session_map_del(map_fd, seid) == 0)
                        dprintf(client_fd, "OK\n");
                    else
                        dprintf(client_fd, "ERR map delete failed\n");
                } else {
                    dprintf(client_fd, "ERR no map\n");
                }
            }
        } else {
            dprintf(client_fd, "ERR unknown cmd\n");
        }

        close(client_fd);
    }

    if (map_fd >= 0) close(map_fd);
    close(server_fd);
    unlink(SOCKET_PATH);
    printf("sessiond_pfcp_adapter exiting\n");
    return 0;
}
