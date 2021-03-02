/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#include "util.h"

#include <netinet/sctp.h>
#include <unistd.h>

#include "sctpd.h"

namespace magma {
namespace sctpd {

// Set up basic sctp options of socket "sd"
int set_sctp_opts(
    const int sd, const uint16_t instreams, const uint16_t outstreams,
    const uint16_t max_attempts, const uint16_t init_timeout);

// Convert address specified in InitReq into struct sockaddr for sctp setup
int convert_addrs(const InitReq* req, struct sockaddr** addrs, int* num_addrs);

int create_sctp_sock(const InitReq& req) {
  int num_addrs;
  struct sockaddr* addrs;
  int sd;

  if (!req.use_ipv4() && !req.use_ipv6()) return -1;

  if (req.ipv4_addrs_size() == 0 && req.ipv6_addrs_size() == 0) return -1;

  sd = socket(AF_INET6, SOCK_STREAM, IPPROTO_SCTP);
  if (sd < 0) {
    MLOG_perror("socket");
    return -1;
  }

  if (set_sctp_opts(sd, SCTP_IN_STREAMS, SCTP_OUT_STREAMS, 0, 0) < 0) {
    goto fail;
  }

  if (convert_addrs(&req, &addrs, &num_addrs) < 0) goto fail;

  if (sctp_bindx(sd, addrs, num_addrs, SCTP_BINDX_ADD_ADDR)) {
    MLOG_perror("sctp_bindx ADD error");
    if (sctp_bindx(sd, addrs, num_addrs, SCTP_BINDX_REM_ADDR)) {
      MLOG_perror("sctp_bindx REM error");
    }
    goto fail;
  }

  if (listen(sd, 5) < 0) {
    MLOG_perror("listen");
    goto fail;
  }

  free(addrs);

  return sd;

fail:
  close(sd);
  free(addrs);
  return -1;
}

int set_sctp_opts(
    const int sd, const uint16_t instreams, const uint16_t outstreams,
    const uint16_t max_attempts, const uint16_t init_timeout) {
  struct sctp_initmsg init;
  init.sinit_num_ostreams   = outstreams;
  init.sinit_max_instreams  = instreams;
  init.sinit_max_attempts   = max_attempts;
  init.sinit_max_init_timeo = init_timeout;

  if (setsockopt(sd, IPPROTO_SCTP, SCTP_INITMSG, &init, sizeof(init)) < 0) {
    MLOG_perror("setsockopt sctp");
    return -1;
  }

  int on = 1;

  struct linger sctp_linger;
  sctp_linger.l_onoff  = on;
  sctp_linger.l_linger = 0;  // send an ABORT
  if (setsockopt(sd, SOL_SOCKET, SO_LINGER, &sctp_linger, sizeof(sctp_linger)) <
      0) {
    MLOG_perror("setsockopt linger");
    return -1;
  }

  struct sctp_event_subscribe event;
  event.sctp_association_event = on;
  event.sctp_shutdown_event    = on;
  event.sctp_data_io_event     = on;

  if (setsockopt(sd, IPPROTO_SCTP, SCTP_EVENTS, &event, sizeof(event)) < 0) {
    MLOG_perror("setsockopt");
    return -1;
  }

  return 0;
}

int convert_addrs(const InitReq* req, struct sockaddr** addrs, int* num_addrs) {
  int i;
  struct sockaddr_in* ipv4_addr;
  struct sockaddr_in6* ipv6_addr;

  auto num_ipv4_addrs = req->ipv4_addrs_size();
  auto num_ipv6_addrs = req->ipv6_addrs_size();
  *num_addrs          = num_ipv4_addrs + num_ipv6_addrs;

  *addrs = (struct sockaddr*) calloc(*num_addrs, sizeof(struct sockaddr*));
  if (*addrs == NULL) return -1;

  for (i = 0; i < num_ipv4_addrs; i++) {
    ipv4_addr = (struct sockaddr_in*) &(*addrs)[i];

    ipv4_addr->sin_family = AF_INET;
    ipv4_addr->sin_port   = htons(req->port());
    if (inet_aton(req->ipv4_addrs(i).c_str(), &ipv4_addr->sin_addr) < 0) {
      return -1;
    }
  }

  for (i = 0; i < num_ipv6_addrs; i++) {
    ipv6_addr = (struct sockaddr_in6*) &(*addrs)[i + num_ipv4_addrs];

    ipv6_addr->sin6_family = AF_INET6;
    ipv6_addr->sin6_port   = htons(req->port());
    if (inet_pton(AF_INET6, req->ipv6_addrs(i).c_str(), &ipv6_addr->sin6_addr) <
        0) {
      return -1;
    }
  }

  return 0;
}

int pull_peer_ipaddr(
    const int sd, const uint32_t assoc_id, std::string& ran_cp_ipaddr) {
  int n_remote_addr             = -1;
  struct sockaddr* remote_addrs = NULL;
  n_remote_addr                 = sctp_getpaddrs(sd, assoc_id, &remote_addrs);

  // Since socket is opened as AF_INET6, remote address comes as IPv6 formatted
  // for both IPv4 and IPv6 end points
  const uint8_t* remote_addr_ipv6_bytes =
      ((const struct sockaddr_in6*) &remote_addrs[0])->sin6_addr.s6_addr;
  const char* fromaddr = NULL;
  if (n_remote_addr >= 1) {
    // Picking the first address only.
    // Check if remote_addrs[0] is IPv6 formatted IPv4 address
    if (IN6_IS_ADDR_V4MAPPED(
            &((struct sockaddr_in6*) &remote_addrs[0])->sin6_addr)) {
      // First 12 bytes are ::FFFF for IPv4-mapped-IPv6
      fromaddr      = (const char*) remote_addr_ipv6_bytes + 12;
      ran_cp_ipaddr = std::string(fromaddr, 4);
    } else {
      fromaddr      = (const char*) remote_addr_ipv6_bytes;
      ran_cp_ipaddr = std::string(fromaddr, 16);
    }
  }

  sctp_freepaddrs(remote_addrs);
  return 0;
}

}  // namespace sctpd
}  // namespace magma
