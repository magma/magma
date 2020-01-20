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

#include <grpcpp/create_channel.h>
#include <grpcpp/security/credentials.h>
#include <stdint.h>
#include <string.h>
#include <string>
#include <memory>

#include "MobilityClient.h"
#include "rpc_client.h"

namespace grpc {
class Channel;
}  // namespace grpc

// TODO: MobilityService IP:port config (t14002037)
#define MOBILITYD_ENDPOINT "localhost:60051"

using grpc::Channel;
using grpc::ChannelCredentials;
using grpc::CreateChannel;
using grpc::InsecureChannelCredentials;
using magma::MobilityServiceClient;

int get_assigned_ipv4_block(
  int index,
  struct in_addr *netaddr,
  uint32_t *netmask)
{
  const std::shared_ptr<ChannelCredentials> cred = InsecureChannelCredentials();
  const std::shared_ptr<Channel> channel =
    CreateChannel(MOBILITYD_ENDPOINT, cred);
  MobilityServiceClient client(channel);
  int status = client.GetAssignedIPv4Block(index, netaddr, netmask);
  return status;
}

int allocate_ipv4_address(const char *subscriber_id, const char *apn,
                          struct in_addr *addr)
{
  const std::shared_ptr<ChannelCredentials> cred = InsecureChannelCredentials();
  const std::shared_ptr<Channel> channel =
    CreateChannel(MOBILITYD_ENDPOINT, cred);
  MobilityServiceClient client(channel);
  int status = client.AllocateIPv4Address(subscriber_id, apn, addr);
  return status;
}

int release_ipv4_address(const char *subscriber_id, const char *apn,
                         const struct in_addr *addr)
{
  const std::shared_ptr<ChannelCredentials> cred = InsecureChannelCredentials();
  const std::shared_ptr<Channel> channel =
    CreateChannel(MOBILITYD_ENDPOINT, cred);
  MobilityServiceClient client(channel);
  int status = client.ReleaseIPv4Address(subscriber_id, apn, *addr);
  return status;
}

int get_ipv4_address_for_subscriber(
  const char *subscriber_id,
  const char *apn,
  struct in_addr *addr)
{
  const std::shared_ptr<ChannelCredentials> cred = InsecureChannelCredentials();
  const std::shared_ptr<Channel> channel =
    CreateChannel(MOBILITYD_ENDPOINT, cred);
  MobilityServiceClient client(channel);
  int status = client.GetIPv4AddressForSubscriber(subscriber_id, apn, addr);
  return status;
}

int get_subscriber_id_from_ipv4(
  const struct in_addr *addr,
  char **subscriber_id)
{
  const std::shared_ptr<ChannelCredentials> cred = InsecureChannelCredentials();
  const std::shared_ptr<Channel> channel =
    CreateChannel(MOBILITYD_ENDPOINT, cred);
  MobilityServiceClient client(channel);
  std::string subscriber_id_str;
  int status = client.GetSubscriberIDFromIPv4(*addr, &subscriber_id_str);
  if (!subscriber_id_str.empty()) {
    *subscriber_id = strdup(subscriber_id_str.c_str());
  }
  return status;
}
