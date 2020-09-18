#include "AmfServiceClient.h"
#include "ServiceRegistrySingleton.h"
#include "magma_logging.h"

// using google::protobuf::RepeatedPtrField;
using grpc::Status;

namespace {

void call_back_void_amf(grpc::Status, magma::SmContextVoid respvoid)
{
  //do nothinf but to only passing call back
  //cout <<" Only for testing call back" << endl;
}
std::function<void(grpc::Status, magma::SmContextVoid)>callback  = call_back_void_amf;

#if 0
magma::SetSMSessionContextAccess create_set_smsess_access_req(
    const std::string& imsi, const std::string& apn_ip_addr,
    const uint32_t linked_bearer_id,
    const std::vector<magma::PolicyRule>& flows) {
  magma::SetSMSessionContextAccess sreq;
  
  auto *req =  sreq.mutable_rat_specific_context().mutable_m5g_session_context_rsp();

  req.set_pdu_session_id();
  req.set_pdu_session_type(magma::PduSessionType::IPV4);
  req.set_selected_ssc_mode(magma::SscMode::SSC_MODE_1);
  auto req_auth_qos_rules = req.mutable_authorized_qos_rules();
  for (const auto& flow : flows) {
    req_auth_qos_rules->Add()->CopyFrom(flow);
  }
  .
  .
  .

 return req;
}

#endif
}

namespace magma {

AsyncAmfServiceClient::AsyncAmfServiceClient(
    std::shared_ptr<grpc::Channel> channel)
    : stub_(SmfPduSessionSmContext::NewStub(channel)) {}

AsyncAmfServiceClient::AsyncAmfServiceClient()
    : AsyncAmfServiceClient(
          ServiceRegistrySingleton::Instance()->GetGrpcChannel(
              "spgw_service", ServiceRegistrySingleton::LOCAL)) {}


bool AsyncAmfServiceClient::handle_response_to_access(
    const magma::SetSMSessionContextAccess& response) {
  
  MLOG(MINFO) << "Sending Set SM Session Response from SMF ";
  #if 0
  handle_response_to_access_rpc(
      response, [](Status status, SmContextVoid resp) {
        if (!status.ok()) {
          MLOG(MERROR) << "Could not send Set SM Session Response from SMF" << imsi
                       << apn_ip_addr << ": " << status.error_message();
        }
      });
  #endif
  handle_response_to_access_rpc(response, callback);
  return true;
}


void AsyncAmfServiceClient::handle_response_to_access_rpc(
    const SetSMSessionContextAccess& response,
    std::function<void(Status, SmContextVoid)> callback) {
	
  std::cerr <<__LINE__ << " " << __FUNCTION__ 
	         << " calling RPC of AMF SetSMFSessionContext and sending message" << "\n";
  auto local_resp = new AsyncLocalResponse<SmContextVoid>(
      std::move(callback), RESPONSE_TIMEOUT);
  local_resp->set_response_reader(std::move(
      stub_->AsyncSetSmfSessionContext(local_resp->get_context(), response, &queue_)));
}

}
