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

#include "MobilityServiceClient.h"

#include <cassert>
#include <grpcpp/channel.h>
#include <grpcpp/impl/codegen/client_context.h>
#include <grpcpp/impl/codegen/status.h>
#include <netinet/in.h>
#include <cstring>
#include <iostream>
#include <memory>
#include <string>

#include "lte/protos/mobilityd.grpc.pb.h"
#include "lte/protos/mobilityd.pb.h"
#include "orc8r/protos/common.pb.h"
#include "lte/protos/subscriberdb.pb.h"

using grpc::Channel;
using grpc::ClientContext;
using grpc::Status;
using grpc::ChannelCredentials;
using grpc::InsecureChannelCredentials;
using magma::orc8r::Void;

// TODO: MobilityService IP:port config (t14002037)
#define MOBILITYD_ENDPOINT "localhost:60051"

namespace magma {
namespace lte {

int MobilityServiceClient::AllocateIPv4Address(
  const std::string& imsi,
  const std::string& apn,
  struct in_addr* addr)
{
  AllocateIPRequest request = AllocateIPRequest();
  request.set_version(AllocateIPRequest::IPV4);

  SubscriberID* sid = request.mutable_sid();
  sid->set_id(imsi);
  sid->set_type(SubscriberID::IMSI);

  request.set_apn(apn);

  ClientContext context;
  IPAddress ip_msg;
  // TODO: Add AllocateIPv4Address response handler here
  Status status = stub_->AllocateIPAddress(&context, request, &ip_msg);
  if (!status.ok()) {
    // TODO: use logging
    std::cout << "AllocateIPAddress fails with code " << status.error_code()
              << ", msg: " << status.error_message() << std::endl;
    return status.error_code();
  }
  memcpy(addr, ip_msg.mutable_address()->c_str(), sizeof(in_addr));
  return 0;
}

int MobilityServiceClient::ReleaseIPv4Address(
  const std::string& imsi,
  const std::string& apn,
  const struct in_addr& addr)
{
  ReleaseIPRequest request = ReleaseIPRequest();
  SubscriberID* sid = request.mutable_sid();
  sid->set_id(imsi);
  sid->set_type(SubscriberID::IMSI);

  request.set_apn(apn);

  IPAddress* ip = request.mutable_ip();
  ip->set_version(IPAddress::IPV4);
  ip->set_address(&addr, sizeof(struct in_addr));

  ReleaseIPv4AddressRPC(request, [](const Status& status, Void resp) {
    if (!status.ok()) {
      // TODO: use logging
      std::cout << "ReleaseIPAddress fails with code " << status.error_code()
                << ", msg: " << status.error_message() << std::endl;
    }
  });
  return 0;
}

// More than one IP can be assigned due to multiple PDNs (one per PDN)
// Get PDN specific IP address
int MobilityServiceClient::GetIPv4AddressForSubscriber(
  const std::string& imsi,
  const std::string& apn,
  struct in_addr* addr)
{
  IPLookupRequest request = IPLookupRequest();
  SubscriberID* sid = request.mutable_sid();
  sid->set_id(imsi);
  sid->set_type(SubscriberID::IMSI);

  request.set_apn(apn);

  IPAddress ip_msg;

  ClientContext context;

  Status status = stub_->GetIPForSubscriber(&context, request, &ip_msg);
  if (!status.ok()) {
    std::cout << "GetIPv4AddressForSubscriber fails with code "
              << status.error_code() << ", msg: " << status.error_message()
              << std::endl;
    return status.error_code();
  }
  memcpy(addr, ip_msg.mutable_address()->c_str(), sizeof(in_addr));
  return 0;
}

int MobilityServiceClient::GetAssignedIPv4Block(
  int index,
  struct in_addr* netaddr,
  uint32_t* netmask)
{
  ClientContext context;
  Void request;
  ListAddedIPBlocksResponse response;
  uint32_t prefix_len = 0;

  assert(index == 0 && "Only one IP block is supported currently");

  Status status = stub_->ListAddedIPv4Blocks(&context, request, &response);
  if (!status.ok()) {
    // TODO: use logging
    std::cout << "GetAssignedIPBlock fails with code " << status.error_code()
              << ", msg: " << status.error_message() << std::endl;
    return status.error_code();
  }

  memcpy(
    netaddr,
    response.mutable_ip_block_list(index)->mutable_net_address()->c_str(),
    sizeof(in_addr));
  *netmask = response.mutable_ip_block_list(index)->prefix_len();
  return 0;
}

int MobilityServiceClient::GetSubscriberIDFromIPv4(
  const struct in_addr& addr,
  std::string* imsi)
{
  IPAddress ip_addr = IPAddress();
  ip_addr.set_version(IPAddress::IPV4);
  ip_addr.set_address(&addr, sizeof(struct in_addr));

  SubscriberID match;

  ClientContext context;
  Void resp;
  Status status = stub_->GetSubscriberIDFromIP(&context, ip_addr, &match);
  if (!status.ok()) {
    std::cout << "GetSubscriberIDFromIPv4 fails with code "
              << status.error_code() << ", msg: " << status.error_message()
              << std::endl;
    return status.error_code();
  }
  imsi->assign(match.id());
  return 0;
}

void MobilityServiceClient::AllocateIPv4AddressRPC(
  const AllocateIPRequest& request,
  const std::function<void(Status, IPAddress)>& callback)
{
  auto localResp =
    new AsyncLocalResponse<IPAddress>(std::move(callback), RESPONSE_TIMEOUT);
  localResp->set_response_reader(std::move(
    stub_->AsyncAllocateIPAddress(localResp->get_context(), request, &queue_)));
}

void MobilityServiceClient::ReleaseIPv4AddressRPC(
  const ReleaseIPRequest& request,
  const std::function<void(grpc::Status, magma::orc8r::Void)>& callback)
{
  auto localResp =
    new AsyncLocalResponse<Void>(callback, RESPONSE_TIMEOUT);
  localResp->set_response_reader(std::move(
    stub_->AsyncReleaseIPAddress(localResp->get_context(), request, &queue_)));
}

MobilityServiceClient::MobilityServiceClient()
{
  const std::shared_ptr<ChannelCredentials> cred = InsecureChannelCredentials();
  const std::shared_ptr<Channel> channel =
    CreateChannel(MOBILITYD_ENDPOINT, cred);
  stub_ = MobilityService::NewStub(channel);
}

MobilityServiceClient& MobilityServiceClient::getInstance()
{
  static MobilityServiceClient instance;
  return instance;
}

} // namespace lte
} // namespace magma
