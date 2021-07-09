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
#include "M5GMobilityServiceClient.h"
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
using magma5g::AsyncM5GMobilityServiceClient;
using magma5g::AsyncSmfServiceClient;

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

  AsyncSmfServiceClient::getInstance().set_smf_session(req);
  OAILOG_DEBUG(LOG_AMF_APP, "sent Establish Request 2nd time to SMF");

  return rc;
}

/***************************************************************************
**                                                                        **
** Name:    amf_smf_create_ipv4_session_grpc_req()                        **
**                                                                        **
** Description: Fill session establishment gRPC request to SMF            **
**                                                                        **
**                                                                        **
***************************************************************************/
int amf_smf_create_ipv4_session_grpc_req(
    char* imsi, uint8_t* apn, uint32_t pdu_session_id,
    uint32_t pdu_session_type, uint8_t* gnb_gtp_teid, uint8_t pti,
    uint8_t* gnb_gtp_teid_ip_addr, char* ipv4_addr) {
  // Creating the request
  magma::lte::SetSMSessionContext req;

  auto* req_common = req.mutable_common_context();

  // Encode IMSI
  req_common->mutable_sid()->mutable_id()->assign(imsi);

  // Encode TYPE IMSI
  req_common->mutable_sid()->set_type(
      magma::lte::SubscriberID_IDType::SubscriberID_IDType_IMSI);

  // Encode APU
  req_common->set_apn((char*) apn);

  // Encode RAT TYPE
  req_common->set_rat_type(magma::lte::RATType::TGPP_NR);

  // Put in CREATING STATE
  req_common->set_sm_session_state(magma::lte::SMSessionFSMState::CREATING_0);

  // Create with Default Version
  req_common->set_sm_session_version(VERSION_0);

  auto* req_rat_specific =
      req.mutable_rat_specific_context()->mutable_m5gsm_session_context();

  // Set the Session ID
  req_rat_specific->set_pdu_session_id(pdu_session_id);

  // Set the Type of Request
  req_rat_specific->set_rquest_type(magma::lte::RequestType::INITIAL_REQUEST);

  // Set the Address type
  req_rat_specific->mutable_pdu_address()->set_redirect_address_type(
      magma::lte::RedirectServer::IPV4);

  // Type is IPv4
  req_rat_specific->set_pdu_session_type(magma::lte::PduSessionType::IPV4);

  // TEID of GNB
  req_rat_specific->mutable_gnode_endpoint()->mutable_teid()->assign(
      (char*) gnb_gtp_teid);

  // IP Address of GNB
  req_rat_specific->mutable_gnode_endpoint()->mutable_end_ipv4_addr()->assign(
      (char*) gnb_gtp_teid_ip_addr);

  // Set the PTI
  req_rat_specific->set_procedure_trans_identity((const char*) (&(pti)));

  // Set the PDU Address
  req_rat_specific->mutable_pdu_address()->set_redirect_server_address(
      (char*) ipv4_addr);

  AsyncSmfServiceClient::getInstance().set_smf_session(req);

  return RETURNok;
}

/***************************************************************************
**                                                                        **
** Name:    amf_smf_create_ipv4_session_grpc_req()                        **
**                                                                        **
** Description: Fill session establishment gRPC request to SMF            **
**                                                                        **
**                                                                        **
***************************************************************************/
int amf_smf_create_session_on_service_req(
    char* imsi, uint8_t* apn, uint32_t pdu_session_id,
    uint32_t pdu_session_type, uint8_t* gnb_gtp_teid, uint8_t pti,
    uint8_t* gnb_gtp_teid_ip_addr, char* ipv4_addr) {
  // Creating the request
  magma::lte::SetSMSessionContext req;

  auto* req_common = req.mutable_common_context();

  // Encode IMSI
  req_common->mutable_sid()->mutable_id()->assign(imsi);

  // Encode TYPE IMSI
  req_common->mutable_sid()->set_type(
      magma::lte::SubscriberID_IDType::SubscriberID_IDType_IMSI);

  // Encode APU
  req_common->set_apn((char*) apn);

  // Encode RAT TYPE
  req_common->set_rat_type(magma::lte::RATType::TGPP_NR);

  // Put in CREATING STATE
  req_common->set_sm_session_state(magma::lte::SMSessionFSMState::CREATING_0);

  // Create with Default Version
  req_common->set_sm_session_version(VERSION_0);

  auto* req_rat_specific =
      req.mutable_rat_specific_context()->mutable_m5gsm_session_context();

  // Set the Session ID
  req_rat_specific->set_pdu_session_id(pdu_session_id);

  // Set the Type of Request
  req_rat_specific->set_rquest_type(magma::lte::RequestType::INITIAL_REQUEST);

  // Set the Address type
  req_rat_specific->mutable_pdu_address()->set_redirect_address_type(
      magma::lte::RedirectServer::IPV4);

  // Type is IPv4
  req_rat_specific->set_pdu_session_type(magma::lte::PduSessionType::IPV4);

  // TEID of GNB
  req_rat_specific->mutable_gnode_endpoint()->mutable_teid()->assign(
      (char*) gnb_gtp_teid);

  // IP Address of GNB
  req_rat_specific->mutable_gnode_endpoint()->mutable_end_ipv4_addr()->assign(
      (char*) gnb_gtp_teid_ip_addr);

  // Set the PTI
  req_rat_specific->set_procedure_trans_identity((const char*) (&(pti)));

  // Set the PDU Address
  req_rat_specific->mutable_pdu_address()->set_redirect_server_address(
      (char*) ipv4_addr);

  AsyncSmfServiceClient::getInstance().set_smf_session(req);

  return RETURNok;
}
/***************************************************************************
 * **                                                                        **
 * ** Name:    amf_smf_create_pdu_session()                                  **
 * **                                                                        **
 * ** Description: Trigger PDU Session Creation in SMF                       **
 * **                                                                        **
 * **                                                                        **
 * ***************************************************************************/
int amf_smf_create_pdu_session(
    amf_smf_establish_t* message, char* imsi, uint32_t version) {
  AsyncM5GMobilityServiceClient::getInstance().allocate_ipv4_address(
      imsi, "internet", message->pdu_session_id, message->pti, AF_INET,
      message->gnb_gtp_teid, 4, message->gnb_gtp_teid_ip_addr, 4);

  return (RETURNok);
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

  AsyncSmfServiceClient::getInstance().set_smf_session(req);

  return RETURNok;
}
}  // namespace magma5g
