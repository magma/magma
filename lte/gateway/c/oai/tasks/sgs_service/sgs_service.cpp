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
#include "sgs_service.h"

#include <grpcpp/grpcpp.h>
#include <grpcpp/security/server_credentials.h>
#include <memory>

#include "CSFBGatewayServiceImpl.h"

extern "C" {
#include "log.h"
}

using grpc::InsecureServerCredentials;
using grpc::Server;
using grpc::ServerBuilder;
using magma::CSFBGatewayServiceImpl;

static CSFBGatewayServiceImpl sgs_service;
static std::unique_ptr<Server> server;

void start_sgs_service_server(bstring server_address)
{
  OAILOG_INFO(
    LOG_MME_APP,
    (char *) "Starting sgs grpc service in service_server at : %s\n ",
    bdata(server_address));
  ServerBuilder builder;
  builder.AddListeningPort(bdata(server_address), InsecureServerCredentials());
  builder.RegisterService(&sgs_service);
  server = builder.BuildAndStart();
  server->Wait(); // Blocking call
}

void stop_sgs_service_server(void)
{
  server->Shutdown();
}
