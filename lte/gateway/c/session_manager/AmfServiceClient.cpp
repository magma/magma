#include "AmfServiceClient.h"
#include "includes/ServiceRegistrySingleton.h"
#include "magma_logging.h"

using grpc::Status;

namespace {

void call_back_void_amf(grpc::Status, magma::SmContextVoid respvoid) {
  // do nothinf but to only passing call back
}
std::function<void(grpc::Status, magma::SmContextVoid)> callback =
    call_back_void_amf;

}  // namespace

namespace magma {

AsyncAmfServiceClient::AsyncAmfServiceClient(
    std::shared_ptr<grpc::Channel> channel)
    : stub_(SmfPduSessionSmContext::NewStub(channel)) {}

AsyncAmfServiceClient::AsyncAmfServiceClient()
    : AsyncAmfServiceClient(
          ServiceRegistrySingleton::Instance()->GetGrpcChannel(
              "amf_service", ServiceRegistrySingleton::LOCAL)) {}

bool AsyncAmfServiceClient::handle_response_to_access(
    const magma::SetSMSessionContextAccess& response) {
  MLOG(MINFO) << "Sending Set SM Session Response from SMF ";
  auto local_resp = new AsyncLocalResponse<SmContextVoid>(
      std::move(callback), RESPONSE_TIMEOUT);
  local_resp->set_response_reader(std::move(stub_->AsyncSetSmfSessionContext(
      local_resp->get_context(), response, &queue_)));
  return true;
}

bool AsyncAmfServiceClient::handle_notification_to_access(
    const magma::SetSmNotificationContext& notif) {
  MLOG(MINFO) << "Sending Set SM Session Notification from SMF ";
  auto local_resp = new AsyncLocalResponse<SmContextVoid>(
      std::move(callback), RESPONSE_TIMEOUT);
  local_resp->set_response_reader(std::move(stub_->AsyncSetAmfNotification(
      local_resp->get_context(), notif, &queue_)));
  return true;
}

}  // namespace magma
