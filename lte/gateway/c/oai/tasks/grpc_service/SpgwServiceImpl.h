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
    ServerContext *context,
    const CreateBearerRequest *request,
    CreateBearerResult *response) override;

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
    ServerContext *context,
    const DeleteBearerRequest *request,
    DeleteBearerResult *response) override;
};

} // namespace magma
