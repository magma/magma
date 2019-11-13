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
#include <string>

#include "lte/protos/spgw_service.pb.h"

extern "C" {
#include "spgw_service_handler.h"
#include "log.h"
}
#include "SpgwServiceImpl.h"

namespace grpc {
class ServerContext;
} // namespace grpc

using grpc::ServerContext;
using grpc::Status;
using magma::CreateBearerRequest;
using magma::CreateBearerResult;
using magma::DeleteBearerRequest;
using magma::DeleteBearerResult;
using magma::SpgwService;

namespace magma {
using namespace lte;

SpgwServiceImpl::SpgwServiceImpl() {}

Status SpgwServiceImpl::CreateBearer(
  ServerContext *context,
  const CreateBearerRequest *request,
  CreateBearerResult *response)
{
  OAILOG_INFO(LOG_UTIL,"Received CreateBearer GRPC request\n");
  itti_pgw_nw_init_actv_bearer_request_t itti_msg;
  itti_msg.imsi_length = request->sid().id().size();
  strcpy(itti_msg.imsi, request->sid().id().c_str());
  itti_msg.lbi = request->link_bearer_id();

  //TODO: figure tfts out
  memset(&itti_msg.ul_tft,0,sizeof(traffic_flow_template_t));
  memset(&itti_msg.dl_tft,0,sizeof(traffic_flow_template_t));

  bearer_qos_t *qos = &itti_msg.eps_bearer_qos;
  for (const auto &policy_rule : request->policy_rules()) {
    qos->pci = policy_rule.qos().arp().pre_capability();
    qos->pl = policy_rule.qos().arp().priority_level();
    qos->pvi = policy_rule.qos().arp().pre_vulnerability();
    qos->qci = policy_rule.qos().qci();
    qos->gbr.br_ul = policy_rule.qos().gbr_ul();
    qos->gbr.br_dl = policy_rule.qos().gbr_dl();
    qos->mbr.br_ul = policy_rule.qos().max_req_bw_ul();
    qos->mbr.br_dl = policy_rule.qos().max_req_bw_dl();
    send_activate_bearer_request_itti(&itti_msg);
  }

  return Status::OK;
}

Status SpgwServiceImpl::DeleteBearer(
  ServerContext *context,
  const DeleteBearerRequest *request,
  DeleteBearerResult *response)
{
  OAILOG_INFO(LOG_UTIL,"Received DeleteBearer GRPC request\n");
  itti_pgw_nw_init_deactv_bearer_request_t itti_msg;
  itti_msg.imsi_length = request->sid().id().size();
  strcpy(itti_msg.imsi, request->sid().id().c_str());
  itti_msg.lbi = request->link_bearer_id();
  itti_msg.no_of_bearers = request->eps_bearer_ids_size();
  for (int i = 0; i < request->eps_bearer_ids_size() && i < BEARERS_PER_UE;
       i++) {
    itti_msg.ebi[i] = request->eps_bearer_ids(i);
  }
  send_deactivate_bearer_request_itti(&itti_msg);
  return Status::OK;
}
} // namespace magma
