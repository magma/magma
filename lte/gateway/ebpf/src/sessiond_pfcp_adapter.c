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
static const char *PINNED_MAP_PATH = "/sys/fs/bpf/gtp/ue_session_map"; // expected pinned map path

static void sigint(int sig) { exiting = true; }

/*
 * NOTE: the session map key/value must match your BPF program.
 * This is an example layout:
 */
typedef struct {
    __be32 ue_ip;   // UE IPv4 (network byte order)
} ue_session_key_t;

typedef struct {
    __be32 enb_ip;

    __u32 teid_ul_in;
    __u32 teid_ul_out;
    __u32 teid_dl_in;
    __u32 teid_dl_out;

    __u32 s1u_ifindex;
    __u32 sgi_ifindex;
    __u32 ovs_ifindex;

    __u8  ul_mac_src[6];
    __u8  ul_mac_dst[6];

    __u32 qos_mark;
    __u32 bearer_id;

    __u64 ul_bytes;
    __u64 dl_bytes;
    __u64 ul_packets;
    __u64 dl_packets;

    __u64 last_seen;
    __u32 session_flags;

    __u8  imsi[16];
    __u32 imsi_len;
    __u64 encoded_imsi;

    __u8  qfi;
    __u32 tunnel_id;
    __be32 tun_ipv4_dst;

    __u8  tun_flags;
    __u8  direction;
    __u32 original_port;
    __u8  reserved[3];

    __u32 metadata_mark;
} ue_session_info_t;

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
static int session_map_add(
    int map_fd,
    const char *ue_ip_str,
    const char *enb_ip_str,
    __u32 teid_ul_in,
    __u32 teid_dl_out,
    const char *ifname
) {
    ue_session_key_t key = {};
    ue_session_info_t val = {};

    if (parse_ipv4(ue_ip_str, &key.ue_ip) != 0)
        return -1;

    if (parse_ipv4(enb_ip_str, &val.enb_ip) != 0)
        return -1;

    val.teid_ul_in = teid_ul_in;
    val.teid_dl_out = teid_dl_out;

    int ifidx = bpf_utils_ifindex(ifname);
    if (ifidx <= 0)
        return -1;

    val.s1u_ifindex = ifidx;
    val.ovs_ifindex = ifidx;

    val.session_flags = 1;  // ACTIVE

    if (bpf_map_update_elem(map_fd, &key, &val, BPF_ANY) != 0) {
        perror("bpf_map_update_elem");
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
         *   ADD <ue_ip> <enb_ip> <teid_ul> <teid_dl> <ifname>\n
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
