/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the Apache License, Version 2.0  (the "License"); you may not use this file
 * except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

extern "C" {
#include "sctpd_downlink_client.h"

#include <arpa/inet.h>

#include "assertions.h"
#include "log.h"

#include "sctp_defs.h"
}

#include <memory.h>

#include <grpcpp/grpcpp.h>

#include <lte/protos/sctpd.grpc.pb.h>

namespace magma {
namespace lte {

using grpc::Channel;
using grpc::ClientContext;

using magma::sctpd::InitReq;
using magma::sctpd::InitRes;
using magma::sctpd::SctpdDownlink;
using magma::sctpd::SendDlReq;
using magma::sctpd::SendDlRes;

class SctpdDownlinkClient {
 public:
  explicit SctpdDownlinkClient(
    const std::shared_ptr<Channel>& channel,
    bool force_restart);

  int init(InitReq& req, InitRes* res);
  int sendDl(SendDlReq& req, SendDlRes* res);

  bool should_force_restart = false;

 private:
  std::unique_ptr<SctpdDownlink::Stub> _stub;
};

SctpdDownlinkClient::SctpdDownlinkClient(
  const std::shared_ptr<Channel>& channel,
  bool force_restart)
{
  _stub = SctpdDownlink::NewStub(channel);
  should_force_restart = force_restart;
}

int SctpdDownlinkClient::init(InitReq& req, InitRes* res)
{
  assert(res != nullptr);

  ClientContext context;

  auto status = _stub->Init(&context, req, res);

  if (!status.ok()) {
    OAILOG_ERROR(
      LOG_SCTP, "sctpdl.init error = %s\n", status.error_message().c_str());
  }

  return status.ok() ? 0 : -1;
}

int SctpdDownlinkClient::sendDl(SendDlReq& req, SendDlRes* res)
{
  assert(res != nullptr);

  ClientContext context;

  auto status = _stub->SendDl(&context, req, res);

  if (!status.ok()) {
    OAILOG_ERROR(
      LOG_SCTP, "sctpdl.senddl error = %s\n", status.error_message().c_str());
  }

  return status.ok() ? 0 : -1;
}

} // namespace lte
} // namespace magma

using magma::lte::SctpdDownlinkClient;
using magma::sctpd::InitReq;
using magma::sctpd::InitRes;
using magma::sctpd::SendDlReq;
using magma::sctpd::SendDlRes;

std::unique_ptr<SctpdDownlinkClient> _client = nullptr;

int init_sctpd_downlink_client(bool force_restart)
{
  auto channel =
    grpc::CreateChannel(DOWNSTREAM_SOCK, grpc::InsecureChannelCredentials());
  _client = std::make_unique<SctpdDownlinkClient>(channel, force_restart);
}

// init
int sctpd_init(sctp_init_t* init)
{
  assert(init != nullptr);

  int i;
  InitReq req;
  InitRes res;
  char ipv4_str[INET_ADDRSTRLEN];
  char ipv6_str[INET6_ADDRSTRLEN];

  req.set_use_ipv4(init->ipv4);
  req.set_use_ipv6(init->ipv6);

  for (i = 0; i < init->nb_ipv4_addr; i++) {
    auto ipv4_addr = init->ipv4_address[i];
    if (inet_ntop(AF_INET, &ipv4_addr, ipv4_str, INET_ADDRSTRLEN) < 0) {
      Fatal("failed to convert ipv4 addr\n");
      return -1;
    }
    req.add_ipv4_addrs(ipv4_str);
  }

  for (i = 0; i < init->nb_ipv6_addr; i++) {
    auto ipv6_addr = init->ipv6_address[i];
    if (inet_ntop(AF_INET6, &ipv6_addr, ipv6_str, INET6_ADDRSTRLEN) < 0) {
      Fatal("failed to convert ipv6 addr\n");
      return -1;
    }
    req.add_ipv6_addrs(ipv6_str);
  }

  req.set_port(init->port);
  req.set_ppid(init->ppid);

  req.set_force_restart(_client->should_force_restart);

  auto rc = _client->init(req, &res);
  auto init_ok = res.result() == InitRes::INIT_OK;

  return (rc == 0) && init_ok ? 0 : -1;
}

// sendDl
int sctpd_send_dl(uint32_t assoc_id, uint16_t stream, bstring payload)
{
  SendDlReq req;
  SendDlRes res;

  req.set_assoc_id(assoc_id);
  req.set_stream(stream);
  req.set_payload(bdata(payload), blength(payload));

  auto rc = _client->sendDl(req, &res);

  OAILOG_ERROR(LOG_SCTP, "rc = %d\n", rc);

  return rc == 0 && res.result() == SendDlRes::SEND_DL_OK ? 0 : -1;
}
