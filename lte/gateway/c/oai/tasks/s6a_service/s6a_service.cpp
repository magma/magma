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
#include "s6a_service.h"
#include "S6aServiceImpl.h"
#include "S6aProxyImpl.h"
extern "C" {
#include "log.h"
}

using grpc::InsecureServerCredentials;
using grpc::Server;
using grpc::ServerBuilder;
using magma::S6aProxyImpl;
using magma::S6aServiceImpl;

static S6aServiceImpl s6a_service;
static S6aProxyImpl s6a_proxy;
static std::unique_ptr<Server> server;

void start_s6a_service_server(bstring server_address)
{
  OAILOG_INFO(
    LOG_MME_APP,
    "Starting s6a grpc service in service_server at : %s\n ",
    bdata(server_address));
  ServerBuilder builder;
  builder.AddListeningPort(
    bdata(server_address), grpc::InsecureServerCredentials());
  builder.RegisterService(&s6a_proxy);
  builder.RegisterService(&s6a_service);
  server = builder.BuildAndStart();
}

void stop_s6a_service_server(void)
{
  server->Shutdown();
  server->Wait(); // Blocks until server finishes shutting down
}
