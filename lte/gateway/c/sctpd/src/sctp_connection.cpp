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

#include "lte/gateway/c/sctpd/src/sctp_connection.h"

#include <arpa/inet.h>
#include <assert.h>
#include <errno.h>
#include <netinet/sctp.h>
#include <stdlib.h>
#include <string.h>
#include <sys/epoll.h>
#include <sys/select.h>
#include <sys/socket.h>
#include <unistd.h>

#include "lte/gateway/c/sctpd/src/sctpd.h"
#include "lte/gateway/c/sctpd/src/util.h"

namespace magma {
namespace sctpd {

const int NUM_EPOLL_EVENTS = 10;

SctpConnection::SctpConnection(const InitReq& req, SctpEventHandler& handler)
    : _done(false),
      _handler(handler),
      _ppid(req.ppid()),
      _sctp_desc(0),
      _thread(nullptr) {
  int sock = create_sctp_sock(req);
  if (sock < 0) throw std::exception();

  _sctp_desc = SctpDesc(sock);
}

void SctpConnection::Start() {
  assert(_done == false);
  assert(_thread == nullptr);

  _thread = std::make_unique<std::thread>(&SctpConnection::Listen, this);
}

void SctpConnection::Close() {
  assert(_done == false);
  assert(_thread != nullptr);

  _done = true;
  _thread->join();

  for (auto kv : _sctp_desc) {
    auto assoc = kv.second;
    shutdown(assoc.sd, SHUT_RDWR);
    close(assoc.sd);
  }
  close(_sctp_desc.sd());
}

void SctpConnection::Send(
    uint32_t assoc_id, uint32_t stream, const std::string& msg) {
  assert(_thread != nullptr);

  auto assoc = _sctp_desc.getAssoc(assoc_id);
  assert(assoc.sd >= 0);

  auto buf = msg.c_str();
  auto n   = msg.size();
  auto rc  = sctp_sendmsg(
      assoc.sd, buf, n, NULL, 0, htonl(assoc.ppid), 0, stream, 0, 0);

  if (rc < 0) {
    MLOG_perror("sctp_sendmsg");
    throw std::exception();
  }
}

void SctpConnection::Listen() {
  int server_fd = _sctp_desc.sd();
  MLOG(MINFO) << "starting sctp connection listener sd = "
              << std::to_string(server_fd);

  int epoll_fd = epoll_create(1);
  if (epoll_fd < 0) {
    MLOG_perror("epoll_create");
    std::terminate();
  }

  struct epoll_event event;
  event.events  = EPOLLIN;
  event.data.fd = server_fd;

  if (epoll_ctl(epoll_fd, EPOLL_CTL_ADD, server_fd, &event) < 0) {
    MLOG_perror("epoll_ctl");
    std::terminate();
  }

  struct epoll_event events[NUM_EPOLL_EVENTS];

  while (!_done) {
    int timeout    = 100;  // milliseconds = .1s
    int num_events = epoll_wait(epoll_fd, events, NUM_EPOLL_EVENTS, timeout);

    switch (num_events) {
      case -1: {  // errored
        if (errno == EINTR) continue;
        MLOG_perror("epoll_wait");
        std::terminate();
      }
      case 0: {  // timed out
        continue;
      }
      default: { break; }
    }

    for (int i = 0; i < num_events; i++) {
      if (events[i].data.fd == server_fd) {
        // new connection
        int client_sd = accept(server_fd, NULL, NULL);
        if (client_sd < 0) {
          if (errno == ECONNABORTED || errno == EINTR) continue;
          MLOG_perror("accept");
          std::terminate();
        }

        struct epoll_event event;
        event.events  = EPOLLIN;
        event.data.fd = client_sd;

        if (epoll_ctl(epoll_fd, EPOLL_CTL_ADD, client_sd, &event) < 0) {
          MLOG_perror("epoll_ctl");
          std::terminate();
        }
      } else {
        int client_sd = events[i].data.fd;

        auto status = HandleClientSock(client_sd);

        if ((status == SctpStatus::DISCONNECT) ||
            (status == SctpStatus::NEW_ASSOC_NOTIF_FAILED)) {
          if (epoll_ctl(epoll_fd, EPOLL_CTL_DEL, client_sd, nullptr) < 0) {
            MLOG_perror("epoll_ctl");
            std::terminate();
          }
        }

        if (status == SctpStatus::NEW_ASSOC_NOTIF_FAILED) {
          shutdown(client_sd, 0);
        }
      }
    }
  }
}

SctpStatus SctpConnection::HandleClientSock(int sd) {
  assert(sd >= 0);

  MLOG(MDEBUG) << "HandleClientSock sd = " << std::to_string(sd);

  char msg[SCTP_RECV_BUFFER_SIZE];
  struct sctp_sndrcvinfo sinfo;
  int flags;

  int n = sctp_recvmsg(sd, msg, sizeof(msg), nullptr, nullptr, &sinfo, &flags);

  if (n < 0) {
    MLOG_perror("sctp_recvmsg");
    return SctpStatus::FAILURE;
  }

  if (flags & MSG_NOTIFICATION) {
    auto notif = (union sctp_notification*) msg;

    switch (notif->sn_header.sn_type) {
      case SCTP_SHUTDOWN_EVENT: {
        MLOG(MDEBUG) << "SCTP_SHUTDOWN_EVENT received";
        return HandleComDown(notif->sn_shutdown_event.sse_assoc_id);
      }
      case SCTP_ASSOC_CHANGE: {
        MLOG(MDEBUG) << "SCTP association change event received";
        return HandleAssocChange(sd, &notif->sn_assoc_change);
      }
      default: {
        MLOG(MWARNING) << "Unhandled notification type "
                       << std::to_string(notif->sn_header.sn_type);
        return SctpStatus::OK;
      }
    }
  } else {
    // Data payload received
    SctpAssoc&& assoc = SctpAssoc();
    try {
      assoc = _sctp_desc.getAssoc(sinfo.sinfo_assoc_id);
    } catch (const std::out_of_range&) {
      MLOG(MERROR) << "Received sctp msg for untracked assoc: "
                   << std::to_string(sinfo.sinfo_assoc_id);
      // TODO: handle this case
      return SctpStatus::FAILURE;
    }

    assoc.messages_recv++;

    if (ntohl(sinfo.sinfo_ppid) != assoc.ppid) {
      // may have received unsollicited traffic from stack other than S1AP.
      MLOG(MERROR) << "Received data from peer with unsollicited PPID "
                   << std::to_string(ntohl(sinfo.sinfo_ppid)) << ", expecting "
                   << std::to_string(assoc.ppid);
      return SctpStatus::FAILURE;
    }

    MLOG(MDEBUG) << "[sd:" << std::to_string(sd) << "] msg of len "
                 << std::to_string(n) << " on "
                 << std::to_string(sinfo.sinfo_assoc_id) << ":"
                 << std::to_string(sinfo.sinfo_stream);

    _handler.HandleRecv(
        ntohl(sinfo.sinfo_ppid), sinfo.sinfo_assoc_id, sinfo.sinfo_stream,
        std::string(msg, n));

    return SctpStatus::OK;
  }
}

SctpStatus SctpConnection::HandleAssocChange(
    int sd, struct sctp_assoc_change* change) {
  switch (change->sac_state) {
    case SCTP_COMM_UP: {
      return HandleComUp(sd, change);
    }
    case SCTP_RESTART: {
      return HandleReset(change->sac_assoc_id);
    }
    case SCTP_COMM_LOST:
    case SCTP_SHUTDOWN_COMP:
    case SCTP_CANT_STR_ASSOC: {
      return HandleComDown(change->sac_assoc_id);
    }
    default:
      MLOG(MWARNING) << "Unhandled sctp message "
                     << std::to_string(change->sac_state);
      return SctpStatus::FAILURE;
  }
}

SctpStatus SctpConnection::HandleComUp(
    int sd, struct sctp_assoc_change* change) {
  SctpAssoc assoc;

  assoc.sd         = sd;
  assoc.ppid       = _ppid;
  assoc.assoc_id   = change->sac_assoc_id;
  assoc.instreams  = change->sac_inbound_streams;
  assoc.outstreams = change->sac_outbound_streams;

  _sctp_desc.addAssoc(assoc);

  std::string ran_cp_ipaddr;
  pull_peer_ipaddr(sd, change->sac_assoc_id, ran_cp_ipaddr);

  if (_handler.HandleNewAssoc(
          assoc.ppid, change->sac_assoc_id, change->sac_inbound_streams,
          change->sac_outbound_streams, ran_cp_ipaddr) < 0) {
    _sctp_desc.delAssoc(assoc.assoc_id);
    return SctpStatus::NEW_ASSOC_NOTIF_FAILED;
  }

  return SctpStatus::OK;
}

SctpStatus SctpConnection::HandleComDown(uint32_t assoc_id) {
  MLOG(MDEBUG) << "Sending close connection for assoc_id "
               << std::to_string(assoc_id);

  _sctp_desc.delAssoc(assoc_id);

  _handler.HandleCloseAssoc(_ppid, assoc_id, false);

  return SctpStatus::DISCONNECT;
}

SctpStatus SctpConnection::HandleReset(uint32_t assoc_id) {
  MLOG(MDEBUG) << "Handling sctp reset";

  _handler.HandleCloseAssoc(_ppid, assoc_id, true);

  return SctpStatus::OK;
}

}  // namespace sctpd
}  // namespace magma
