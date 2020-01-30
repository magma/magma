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

#include <grpc++/grpc++.h>
#include <grpcpp/impl/codegen/status.h>

#include "lte/protos/spgw_service.grpc.pb.h"
#include "lte/protos/policydb.pb.h"

extern "C" {
#include "spgw_service_handler.h"
#include "log.h"
}

namespace grpc {
class ServerContext;
} // namespace grpc
namespace magma {
namespace lte {
class CreateBearerRequest;
class CreateBearerResult;
class DeleteBearerRequest;
class DeleteBearerResult;
} // namespace lte
} // namespace magma

using grpc::ServerContext;
using magma::lte::CreateBearerRequest;
using magma::lte::CreateBearerResult;
using magma::lte::DeleteBearerRequest;
using magma::lte::DeleteBearerResult;
using magma::lte::SpgwService;

namespace magma {
using namespace lte;

class SpgwServiceImpl final : public SpgwService::Service {
 public:
  SpgwServiceImpl();

  /*
       * CreateBearerRequest.
       *
       * @param context: the grpc Server context
       * @param request: createBearerRequest
       * @param response (out): the CreateBearerResult that contains
                                err message.
       * @return grpc Status instance
       */
  grpc::Status CreateBearer(
    ServerContext* context,
    const CreateBearerRequest* request,
    CreateBearerResult* response) override;

  /*
       * DeleteBearerRequest.
       *
       * @param context: the grpc Server context
       * @param request: DeleteBearerRequest
       * @param response (out): the DeleteBearerResult that contains
                                err message.
       * @return grpc Status instance
       */
  grpc::Status DeleteBearer(
    ServerContext* context,
    const DeleteBearerRequest* request,
    DeleteBearerResult* response) override;

 private:
  /*
    * Fill up the packet filter contents such as flags and flow tuple fields
    * @param pf_content: packet filter content to be filled
    * @param flow_match_rule: pf_content is filled based on flow match rule
    * @return void
    */
  void fillUpPacketFilterContents(
    packet_filter_contents_t* pf_content,
    const FlowMatch* flow_match_rule);

  /*
    * Fill up the ipv4 remote address field in packet filter
    * @param pf_content: packet filter object to be filled
    * @param ipv4addr: IPv4 address in string form (e.g, "172.12.0.1")
    * @return void
    */
  void fillIpv4(
    packet_filter_contents_t* pf_content,
    const std::string ipv4addr);
};

} // namespace magma
