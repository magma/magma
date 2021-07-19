/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
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
#include <functional>
#include <memory>
#include <string>

#include "includes/GRPCReceiver.h"
#include "lte/protos/mobilityd.grpc.pb.h"

namespace grpc {
class Channel;
class ClientContext;
class Status;
}  // namespace grpc

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
      int index, struct in_addr* netaddr, uint32_t* netmask);

  /**
   * Allocate an IPv4 address from the free IP pool (non-blocking)
   * @param imsi: IMSI string
   * @param apn:  APN string
   * @param addr (out): contains the IP address allocated upon returning in
   * "network byte order"
   * @return status of gRPC call
   */
  int AllocateIPv4AddressAsync(
      const std::string& imsi, const std::string& apn,
      const std::function<void(grpc::Status, AllocateIPAddressResponse)>&
          callback);

  /**
   * Allocate an IPv6 address from the free IP pool (non-blocking)
   * @param imsi: IMSI string
   * @param apn:  APN string
   * @param addr (out): contains the IP address allocated upon returning
   * @return status of gRPC call
   */
  int AllocateIPv6AddressAsync(
      const std::string& imsi, const std::string& apn,
      const std::function<void(grpc::Status, AllocateIPAddressResponse)>&
          callback);

  /**
   * Allocate an IPv4v6 address from the free IP pool (non-blocking)
   * @param imsi: IMSI string
   * @param apn:  APN string
   * @param addr (out): contains the IP address allocated upon returning
   * @return status of gRPC call
   */
  int AllocateIPv4v6AddressAsync(
      const std::string& imsi, const std::string& apn,
      const std::function<void(grpc::Status, AllocateIPAddressResponse)>&
          callback);

  /**
   * Release an allocated IPv4 address. (non-blocking)
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
      const std::string& imsi, const std::string& apn,
      const struct in_addr& addr);

  /**
   * Release an allocated IPv6 address. (non-blocking)
   *
   * The released IP address is put into a tombstone state, and recycled
   * periodically.
   *
   * @param imsi: IMSI string
   * @param apn:  APN string
   * @param addr: IPv6 address to release
   * @return 0 on success
   * @return -RPC_STATUS_NOT_FOUND if the requested (SID, IP) pair is not found
   */
  int ReleaseIPv6Address(
      const std::string& imsi, const std::string& apn,
      const struct in6_addr& addr);

  /**
   * Release an allocated IPv4v6 address. (non-blocking)
   *
   * The released IP address is put into a tombstone state, and recycled
   * periodically.
   *
   * @param imsi: IMSI string
   * @param apn:  APN string
   * @param ipv4_addr: IPv4 address to release in "network byte order"
   * @param ipv6_addr: IPv6 address to release
   * @return 0 on success
   * @return -RPC_STATUS_NOT_FOUND if the requested (SID, IP) pair is not found
   */
  int ReleaseIPv4v6Address(
      const std::string& imsi, const std::string& apn,
      const struct in_addr& ipv4_addr, const struct in6_addr& ipv6_addr);

  /*
   * Get the allocated IPv4 address for a subscriber
   * @param imsi: IMSI string
   * @param addr (out): contains the allocated IPv4 address for the subscriber
   * @return 0 on success
   * @return -RPC_STATUS_NOT_FOUND if the SID is not found
   */
  int GetIPv4AddressForSubscriber(
      const std::string& imsi, const std::string& apn, struct in_addr* addr);

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
  static MobilityServiceClient& getInstance();

  MobilityServiceClient(MobilityServiceClient const&) = delete;
  void operator=(MobilityServiceClient const&) = delete;

 private:
  MobilityServiceClient();
  static const uint32_t RESPONSE_TIMEOUT = 10;  // seconds
  std::unique_ptr<MobilityService::Stub> stub_{};

  /**
   * Helper function to chain callback for gRPC response
   * @param request AllocateIP gRPC Request proto
   * @param callback std::function that returns Status and actual gRPC response
   */
  void AllocateIPAddressRPC(
      const AllocateIPRequest& request,
      const std::function<void(grpc::Status, AllocateIPAddressResponse)>&
          callback);

  /**
   * Helper function to chain callback for gRPC response
   * @param request ReleaseIP gRPC Request proto
   * @param callback std::function that returns Status and actual gRPC response
   */
  void ReleaseIPAddressRPC(
      const ReleaseIPRequest& request,
      const std::function<void(grpc::Status, magma::orc8r::Void)>& callback);
};

}  // namespace lte
}  // namespace magma
