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
#include "grpc_service.h"

#include <grpcpp/grpcpp.h>
#include <grpcpp/security/server_credentials.h>
#include <memory>

#include "CSFBGatewayServiceImpl.h"
#include "S6aGatewayImpl.h"
#include "S6aServiceImpl.h"
#include "SpgwServiceImpl.h"

extern "C" {
#include "log.h"
}

using grpc::InsecureServerCredentials;
using grpc::Server;
using grpc::ServerBuilder;
using magma::SpgwServiceImpl;
using magma::S6aGatewayImpl;
using magma::S6aServiceImpl;
using magma::CSFBGatewayServiceImpl;

static SpgwServiceImpl spgw_service;
static S6aServiceImpl s6a_service;
static S6aGatewayImpl s6a_proxy;
static CSFBGatewayServiceImpl sgs_service;
static std::unique_ptr<Server> server;

// TODO Candidate: GRPC service may be evolved into a
// MagmaService, which implements Service303::Service as the
// base service and can add other services on top.
void start_grpc_service(bstring server_address)
{
  OAILOG_INFO(
    LOG_SPGW_APP,
    "Starting service at : %s\n ",
    bdata(server_address));
  ServerBuilder builder;
  builder.AddListeningPort(
    bdata(server_address), grpc::InsecureServerCredentials());
#if SPGW_ENABLE_SESSIOND_AND_MOBILITYD
  builder.RegisterService(&spgw_service);
#endif
  builder.RegisterService(&s6a_proxy);
  builder.RegisterService(&s6a_service);
  builder.RegisterService(&sgs_service);
  server = builder.BuildAndStart();
}

void stop_grpc_service(void)
{
  server->Shutdown();
  server->Wait(); // Blocks until server finishes shutting down
}
