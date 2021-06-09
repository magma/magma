/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#ifdef __cplusplus
extern "C" {
#endif
#include "log.h"
#include "conversions.h"
#include "3gpp_38.401.h"
#ifdef __cplusplus
}
#endif
#include "common_defs.h"
#include <cstdint>
#include <cstring>
#include <string>
#include "MobilityServiceClient.h"
#include <unistd.h>
#include <thread>
#include "SmfServiceClient.h"
#include "amf_smfDefs.h"
#include "conversions.h"
#include "lte/protos/session_manager.pb.h"
#define VERSION_0 0

using grpc::Channel;
using grpc::ChannelCredentials;
using grpc::CreateChannel;
using grpc::InsecureChannelCredentials;
using grpc::Status;
using magma::lte::AllocateIPAddressResponse;
using magma::lte::IPAddress;
using magma::lte::MobilityServiceClient;
using magma::lte::SetSMSessionContext;
using magma::lte::TeidSet;
using magma5g::AsyncSmfServiceClient;

/***************************************************************************
**                                                                        **
** Name:    grpc_prep_estab_req_to_smf()                                  **
**                                                                        **
** Description: get IP from mobilityd and send session request to SMF     **
**                                                                        **
**                                                                        **
***************************************************************************/
void grpc_prep_estab_req_to_smf(magma::lte::SetSMSessionContext req) {
  auto smf_srv_client = std::make_shared<magma5g::AsyncSmfServiceClient>();
  std::thread smf_srv_client_response_handling_thread(
      [&]() { smf_srv_client->rpc_response_loop(); });
  smf_srv_client_response_handling_thread.detach();
  MobilityServiceClient::getInstance().AllocateIPv4AddressAsync(
      req.common_context().sid().id().c_str(),
      req.common_context().apn().c_str(),
      [&req, &smf_srv_client](
          const Status& status, AllocateIPAddressResponse ip_msg) {
        struct in_addr ip_addr = {0};
        char ip_str[INET_ADDRSTRLEN];
        unsigned char buff_ip[4];
        memcpy(
            &(ip_addr), ip_msg.mutable_ip_list(0)->mutable_address()->c_str(),
            sizeof(in_addr));
        memset(ip_str, '\0', sizeof(ip_str));
        //        uint32_t ip_int = ntohl(ip_addr.s_addr);
        //        INT32_TO_BUFFER(ip_int, buff_ip);
        //        memcpy(ip_str, buff_ip, sizeof(buff_ip));
        inet_ntop(AF_INET, &(ip_addr.s_addr), ip_str, INET_ADDRSTRLEN);
        req.mutable_rat_specific_context()
            ->mutable_m5gsm_session_context()
            ->mutable_pdu_address()
            ->set_redirect_server_address((char*) ip_str);
        OAILOG_DEBUG(
            LOG_AMF_APP, "Sending PDU session Establishment Request to SMF");

        smf_srv_client->set_smf_session(req);
      });

  std::this_thread::sleep_for(std::chrono::milliseconds(
      20));  // TODO remove this blocking call without which the gRPC call
             // set_smf_session() doesn't initiate, as per the way its
             // implemeted now
}

namespace magma5g {

/***************************************************************************
**                                                                        **
** Name:   create_session_grpc_req_on_gnb_setup_rsp()                     **
**                                                                        **
** Description:                                                           **
** Creating session grpc req upon receiving resource setup response       **
** from gNB which have gNB IP and TEID, and sends to UPF through SMF      **
**                                                                        **
**                                                                        **
***************************************************************************/
int create_session_grpc_req_on_gnb_setup_rsp(
    amf_smf_establish_t* message, char* imsi, uint32_t version) {
  int rc = RETURNerror;
  magma::lte::SetSMSessionContext req;

  auto smf_srv_client = std::make_shared<magma5g::AsyncSmfServiceClient>();
  std::thread smf_srv_client_response_handling_thread(
      [&]() { smf_srv_client->rpc_response_loop(); });
  smf_srv_client_response_handling_thread.detach();

  auto* req_common = req.mutable_common_context();
  auto* req_rat_specific =
      req.mutable_rat_specific_context()->mutable_m5gsm_session_context();
  // IMSI retrieved from amf context
  req_common->mutable_sid()->mutable_id()->assign(imsi);  // string id
  req_common->mutable_sid()->set_type(
      magma::lte::SubscriberID_IDType::SubscriberID_IDType_IMSI);
  req_common->set_rat_type(magma::lte::RATType::TGPP_NR);
  // PDU session state to CREATING
  req_common->set_sm_session_state(magma::lte::SMSessionFSMState::CREATING_0);
  // Version Added
  req_common->set_sm_session_version(version);
  req_rat_specific->set_pdu_session_id((message->pdu_session_id));
  req_rat_specific->set_rquest_type(magma::lte::RequestType::INITIAL_REQUEST);

  TeidSet* gnode_endpoint = req_rat_specific->mutable_gnode_endpoint();
  gnode_endpoint->set_teid(&message->gnb_gtp_teid, 4);

  req_rat_specific->mutable_gnode_endpoint()->mutable_end_ipv4_addr()->assign(
      (char*) message->gnb_gtp_teid_ip_addr);

  OAILOG_DEBUG(
      LOG_AMF_APP, "Sending PDU session Establishment 2nd Request to SMF");
  smf_srv_client->set_smf_session(req);
  OAILOG_DEBUG(LOG_AMF_APP, "sent Establish Request 2nd time to SMF");

  return rc;
}

/***************************************************************************
**                                                                        **
** Name:    create_session_grpc_req()                                     **
**                                                                        **
** Description: fill session establishment gRPC request to SMF            **
**                                                                        **
**                                                                        **
***************************************************************************/
int create_session_grpc_req(amf_smf_establish_t* message, char* imsi) {
  magma::lte::SetSMSessionContext req;
  auto* req_common = req.mutable_common_context();
  req_common->mutable_sid()->mutable_id()->assign(imsi);
  req_common->mutable_sid()->set_type(
      magma::lte::SubscriberID_IDType::SubscriberID_IDType_IMSI);
  req_common->set_apn("internet");  // TODO upcoming PR this value as default
  req_common->set_rat_type(magma::lte::RATType::TGPP_NR);
  req_common->set_sm_session_state(magma::lte::SMSessionFSMState::CREATING_0);
  req_common->set_sm_session_version(VERSION_0);
  auto* req_rat_specific =
      req.mutable_rat_specific_context()->mutable_m5gsm_session_context();
  req.mutable_rat_specific_context()->mutable_m5gsm_session_context();
  req_rat_specific->set_pdu_session_id(message->pdu_session_id);
  req_rat_specific->set_rquest_type(magma::lte::RequestType::INITIAL_REQUEST);
  req_rat_specific->mutable_pdu_address()->set_redirect_address_type(
      magma::lte::RedirectServer::IPV4);
  req_rat_specific->set_pdu_session_type(magma::lte::PduSessionType::IPV4);
  req_rat_specific->mutable_gnode_endpoint()->mutable_teid()->assign(
      (char*) message->gnb_gtp_teid);
  req_rat_specific->mutable_gnode_endpoint()->mutable_end_ipv4_addr()->assign(
      (char*) message->gnb_gtp_teid_ip_addr);
  req_rat_specific->set_procedure_trans_identity(
      (const char*) (&(message->pti)));
  grpc_prep_estab_req_to_smf(req);

  return RETURNok;
}

/***************************************************************************
**                                                                        **
** Name:    release_session_gprc_req()                                    **
**                                                                        **
** Description: Sends session release request to SMF                      **
**                                                                        **
**                                                                        **
***************************************************************************/
int release_session_gprc_req(amf_smf_release_t* message, char* imsi) {
  magma::lte::SetSMSessionContext req;
  auto* req_common = req.mutable_common_context();
  req_common->mutable_sid()->mutable_id()->assign(imsi);
  req_common->mutable_sid()->set_type(
      magma::lte::SubscriberID_IDType::SubscriberID_IDType_IMSI);
  req_common->set_sm_session_state(magma::lte::SMSessionFSMState::RELEASED_4);
  req_common->set_sm_session_version(1);  // uint32
  auto* req_rat_specific =
      req.mutable_rat_specific_context()->mutable_m5gsm_session_context();
  req_rat_specific->set_pdu_session_id(message->pdu_session_id);
  req_rat_specific->set_procedure_trans_identity(
      (const char*) (&(message->pti)));
  auto smf_srv_client = std::make_shared<magma5g::AsyncSmfServiceClient>();
  std::thread smf_srv_client_response_handling_thread(
      [&]() { smf_srv_client->rpc_response_loop(); });
  smf_srv_client_response_handling_thread.detach();
  smf_srv_client->set_smf_session(req);

  return RETURNok;
}
}  // namespace magma5g
