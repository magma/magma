/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#pragma once

#include <mutex>

#include <lte/protos/policydb.pb.h>
#include <lte/protos/spgw_service.grpc.pb.h>

#include "GRPCReceiver.h"

using grpc::Status;

namespace magma {
using namespace lte;

/**
 * SpgwServiceClient is the base class for sending dedicated bearer
 * create/delete to PGW
 */
class SpgwServiceClient {
public:
  /**
   * Delete a default bearer (all session bearers)
   * @param imsi - msi to identify a UE
   * @param apn_ip_addr - imsi and apn_ip_addrs identify a default bearer
   * @param linked_bearer_id - identifier for link bearer
   * @return true if the operation was successful
   */
  virtual bool delete_default_bearer(const std::string &imsi,
                                     const std::string &apn_ip_addr,
                                     const uint32_t linked_bearer_id) = 0;

  /**
   * Delete a dedicated bearer
   * @param imsi - msi to identify a UE
   * @param apn_ip_addr - imsi and apn_ip_addrs identify a default bearer
   * @param linked_bearer_id - identifier for link bearer
   * @param eps_bearer_ids - ids of bearers to delete
   * @return true if the operation was successful
   */
  virtual bool
  delete_dedicated_bearer(const std::string &imsi,
                          const std::string &apn_ip_addr,
                          const uint32_t linked_bearer_id,
                          const std::vector<uint32_t> &eps_bearer_ids) = 0;

  /**
   * Create a dedicated bearer
   * @param imsi - msi to identify a UE
   * @param apn_ip_addr - imsi and apn_ip_addrs identify a default bearer
   * @param linked_bearer_id - identifier for link bearer
   * @param flows - flow information required for a dedicated bearer
   * @return true if the operation was successful
   */
  virtual bool
  create_dedicated_bearer(const std::string &imsi,
                          const std::string &apn_ip_addr,
                          const uint32_t linked_bearer_id,
                          const std::vector<PolicyRule> &flows) = 0;
};

/**
 * AsyncSpgwServiceClient implements SpgwServiceClient but sends calls
 * asynchronously to PGW.
 */
class AsyncSpgwServiceClient : public GRPCReceiver, public SpgwServiceClient {
public:
  AsyncSpgwServiceClient();

  AsyncSpgwServiceClient(std::shared_ptr<grpc::Channel> pgw_channel);
  /**
   * Delete a default bearer (all session bearers)
   * @param imsi - msi to identify a UE
   * @param apn_ip_addr - imsi and apn_ip_addrs identify a default bearer
   * @param linked_bearer_id - identifier for link bearer
   * @return true if the operation was successful
   */
  bool delete_default_bearer(const std::string &imsi,
                             const std::string &apn_ip_addr,
                             const uint32_t linked_bearer_id);

  /**
   * Delete a dedicated bearer
   * @param imsi - msi to identify a UE
   * @param apn_ip_addr - imsi and apn_ip_addrs identify a default bearer
   * @param linked_bearer_id - identifier for link bearer
   * @param flows - flow information required for a dedicated bearer
   * @return true if the operation was successful
   */
  bool delete_dedicated_bearer(const std::string &imsi,
                               const std::string &apn_ip_addr,
                               const uint32_t linked_bearer_id,
                               const std::vector<uint32_t> &eps_bearer_ids);

  /**
   * Create a dedicated bearer
   * @param imsi - msi to identify a UE
   * @param apn_ip_addr - imsi and apn_ip_addrs identify a default bearer
   * @param linked_bearer_id - identifier for link bearer
   * @param flows - flow information required for a dedicated bearer
   * @return true if the operation was successful
   */
  bool create_dedicated_bearer(const std::string &imsi,
                               const std::string &apn_ip_addr,
                               const uint32_t linked_bearer_id,
                               const std::vector<PolicyRule> &flows);

private:
  static const uint32_t RESPONSE_TIMEOUT = 6; // seconds
  std::unique_ptr<SpgwService::Stub> stub_;

private:
  bool delete_bearer(const std::string &imsi, const std::string &apn_ip_addr,
                     const uint32_t linked_bearer_id,
                     const std::vector<uint32_t> &eps_bearer_ids);

  void
  delete_bearer_rpc(const DeleteBearerRequest &request,
                    std::function<void(Status, DeleteBearerResult)> callback);

  void create_dedicated_bearer_rpc(
      const CreateBearerRequest &request,
      std::function<void(Status, CreateBearerResult)> callback);
};

} // namespace magma
