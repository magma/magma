#ifdef __cplusplus
extern "C" {
#endif
#include "log.h"
#include "conversions.h"
#include "3gpp_38.401.h"
#ifdef __cplusplus
}
#endif
#include <cstdint>
#include <cstring>
#include <string>
#include "MobilityServiceClient.h"
#include <unistd.h>
#include <thread>
#include "SmfServiceClient.h"
#include "amf_common_defs.h"
#include "amf_smfDefs.h"
#include "conversions.h"

using grpc::Channel;
using grpc::ChannelCredentials;
using grpc::CreateChannel;
using grpc::InsecureChannelCredentials;
using grpc::Status;
using magma::lte::AllocateIPAddressResponse;
using magma::lte::IPAddress;
using magma::lte::MobilityServiceClient;

using magma::lte::SetSMSessionContext;
using magma5g::AsyncSmfServiceClient;

namespace {
using namespace magma::lte;

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
        struct in_addr ip_addr;
        memcpy(
            //&(ip_addr), ip_msg.add_ip_list()->mutable_address()->c_str(),
            &(ip_addr), ip_msg.mutable_ip_list(0)->mutable_address()->c_str(),
            sizeof(in_addr));
        uint8_t ip_str[INET_ADDRSTRLEN];
        memset(ip_str, '\0', sizeof(ip_str));
        //  std::string ip_str;
        // inet_ntop(AF_INET, &(ip_addr.s_addr), ip_str, INET_ADDRSTRLEN);
        uint32_t ip_int = ntohl(ip_addr.s_addr);
        unsigned char buff_ip[4];
        INT32_TO_BUFFER(ip_int, buff_ip);
        memcpy(ip_str, buff_ip, 4);  // TODO fix string null_term for bytes
        // ip_str[4] = '\0';//TODO fix string null_term for bytes
        // ip_str.assign((char*)buff_ip);
        req.mutable_rat_specific_context()
            ->mutable_m5gsm_session_context()
            ->mutable_pdu_address()
            ->set_redirect_server_address((char*) ip_str);
        OAILOG_INFO(
            LOG_AMF_APP,
            "AMF_TEST: Sending PDU session Establishment Request to SMF");
        smf_srv_client->set_smf_session(req);
      });

  std::this_thread::sleep_for(std::chrono::milliseconds(
      20));  // TODO remove this blocking call without which the gRPC call
             // set_smf_session() doesn't initiate, as per the way its
             // implemeted now
}

}  // namespace

namespace magma5g {

int create_session_grpc_req(amf_smf_establish_t* message, char* imsi) {
  magma::lte::SetSMSessionContext req;
  auto* req_common = req.mutable_common_context();
  // req_common->mutable_sid()->set_id("imsi00000000001");  // string id
  //  req_common->mutable_sid()->set_id("310410100000001");  // string id
  req_common->mutable_sid()->mutable_id()->assign(imsi);  // string id
  req_common->mutable_sid()->set_type(
      magma::lte::SubscriberID_IDType::SubscriberID_IDType_IMSI);
  req_common->set_apn("blr");  // string apn
  req_common->set_rat_type(magma::lte::RATType::TGPP_NR);
  req_common->set_sm_session_state(magma::lte::SMSessionFSMState::CREATING_0);
  req_common->set_sm_session_version(0);  // uint32

  auto* req_rat_specific =
      req.mutable_rat_specific_context()->mutable_m5gsm_session_context();
  req.mutable_rat_specific_context()->mutable_m5gsm_session_context();

  req_rat_specific->set_pdu_session_id(
      (const char*) (&(message->pdu_session_id)));
  // req_rat_specific->set_pdu_session_id((message->pdu_session_id));

  //     (message->pdu_session_id));

  req_rat_specific->set_rquest_type(magma::lte::RequestType::INITIAL_REQUEST);
  req_rat_specific->mutable_pdu_address()->set_redirect_address_type(
      magma::lte::RedirectServer::IPV4);
  req_rat_specific->set_access_type(magma::lte::AccessType::M_3GPP_ACCESS_3GPP);
  req_rat_specific->set_pdu_session_type(
      (magma::lte::PduSessionType)(message->pdu_session_type - 1));
#if 0
  enum SscMode ssc_mode_var;
  switch(message->ssc_mode){
  case 1:
  ssc_mode_var = SSC_MODE_1;
  break;
  case 2:
  ssc_mode_var = SSC_MODE_2;
  break;
  case 3:
  ssc_mode_var = SSC_MODE_3;
  break;
  }
  req_rat_specific->set_ssc_mode(ssc_mode_var);//TODO fix mapping from NAS not covered in amf_smf_send
#endif
  req_rat_specific->mutable_gnode_endpoint()->mutable_teid()->assign(
      (char*) message->gnb_gtp_teid);
  req_rat_specific->mutable_gnode_endpoint()->mutable_end_ipv4_addr()->assign(
      (char*) message->gnb_gtp_teid_ip_addr);
  req_rat_specific->set_procedure_trans_identity(
      (const char*) (&(message->pti)));
  grpc_prep_estab_req_to_smf(req);

  return RETURNok;
}

int release_session_gprc_req(amf_smf_release_t* message, char* imsi) {
  magma::lte::SetSMSessionContext req;
  auto* req_common = req.mutable_common_context();
  req_common->mutable_sid()->mutable_id()->assign(imsi);
  req_common->mutable_sid()->set_type(
      magma::lte::SubscriberID_IDType::SubscriberID_IDType_IMSI);
  //  req_common->set_apn("blr");  // string apn
  //  req_common->set_rat_type(magma::lte::RATType::TGPP_NR);
  //  req_common->set_sm_session_version(0);  // uint32
  req_common->set_sm_session_state(magma::lte::SMSessionFSMState::RELEASED_4);

  auto* req_rat_specific =
      req.mutable_rat_specific_context()->mutable_m5gsm_session_context();

  req_rat_specific->set_pdu_session_id(
      (const char*) (&(message->pdu_session_id)));
  // req_rat_specific->set_pdu_session_id((message->pdu_session_id));
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
