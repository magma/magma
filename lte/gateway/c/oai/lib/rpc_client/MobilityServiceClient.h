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
#pragma once

#include <arpa/inet.h>
#include <grpc++/grpc++.h>
#include <cstdint>
#include <memory>
#include <string>

#include "lte/protos/mobilityd.grpc.pb.h"

#include "GRPCReceiver.h"

namespace grpc {
class Channel;
class ClientContext;
class Status;
} // namespace grpc

namespace magma {
namespace lte {
/*
 * gRPC client for MobilityService
 */
class MobilityServiceClient : public GRPCReceiver {
 public:
  virtual ~MobilityServiceClient() = default;
  /*
     * Get the address and netmask of an assigned IPv4 block
     *
     * @param index (in): index of the IP block requested, currently only ONE
     * IP block (index=0) is supported
     * @param netaddr (out): network address in "network byte order"
     * @param netmask (out): network address mask
     * @return 0 on success
     * @return -RPC_STATUS_INVALID_ARGUMENT if IPBlock is invalid
     * @return -RPC_STATUS_FAILED_PRECONDITION if IPBlock overlaps
     */
  int GetAssignedIPv4Block(
    int index,
    struct in_addr* netaddr,
    uint32_t* netmask);

  /*
     * Allocate an IPv4 address from the free IP pool
     *
     * @param imsi: IMSI string
     * @param apn:  APN string
     * @param addr (out): contains the IP address allocated upon returning in
     * "network byte order"
     * @return 0 on success
     * @return -RPC_STATUS_RESOURCE_EXHAUSTED if no free IP available
     * @return -RPC_STATUS_ALREADY_EXISTS if an IP has been allocated for the
     *         subscriber
     */
  int AllocateIPv4Address(
    const std::string& imsi,
    const std::string& apn,
    struct in_addr* addr);

  /*
     * Release an allocated IPv4 address.
     *
     * The released IP address is put into a tombstone state, and recycled
     * periodically.
     *
     * @param imsi: IMSI string
     * @param apn:  APN string
     * @param addr: IP address to release in "network byte order"
     * @return 0 on success
     * @return -RPC_STATUS_NOT_FOUND if the requested (SID, IP) pair is not found
     */
  int ReleaseIPv4Address(
    const std::string& imsi,
    const std::string& apn,
    const struct in_addr& addr);

  /*
     * Get the allocated IPv4 address for a subscriber
     * @param imsi: IMSI string
     * @param addr (out): contains the allocated IPv4 address for the subscriber
     * @return 0 on success
     * @return -RPC_STATUS_NOT_FOUND if the SID is not found
     */
  int GetIPv4AddressForSubscriber(
    const std::string& imsi,
    const std::string& apn,
    struct in_addr* addr);

  /*
     * Get the subscriber id given its allocated IPv4 address. If the address
     * isn't associated with a subscriber, then it returns an error
     * @param addr: ipv4 address of subscriber
     * @param imsi (out): contains the imsi of the associated subscriber if it
     *                    exists
     * @return 0 on success
     * @return -RPC_STATUS_NOT_FOUND if IPv4 address is not found
     */
  int GetSubscriberIDFromIPv4(const struct in_addr& addr, std::string* imsi);

 public:
  static MobilityServiceClient &getInstance();

  MobilityServiceClient(MobilityServiceClient const &) = delete;
  void operator=(MobilityServiceClient const &) = delete;

 private:
  MobilityServiceClient();
  static const uint32_t RESPONSE_TIMEOUT = 6; // seconds
  std::shared_ptr<MobilityService::Stub> stub_{};

  /**
  * Helper function to chain callback for gRPC response
  * @param request AllocateIP gRPC Request proto
  * @param callback std::function that returns Status and actual gRPC response
  */
  void AllocateIPv4AddressRPC(
    const AllocateIPRequest& request,
    const std::function<void(grpc::Status, IPAddress)>& callback);

  /**
 * Helper function to chain callback for gRPC response
 * @param request ReleaseIP gRPC Request proto
 * @param callback std::function that returns Status and actual gRPC response
 */
  void ReleaseIPv4AddressRPC(
    const ReleaseIPRequest& request,
    const std::function<void(grpc::Status, magma::orc8r::Void)>& callback);

};

} // namespace lte
} // namespace magma
