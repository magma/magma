#include "lte/gateway/c/session_manager/AmfServiceClient.hpp"

#include <glog/logging.h>
#include <grpcpp/channel.h>
#include <grpcpp/impl/codegen/status.h>
#include <lte/protos/session_manager.grpc.pb.h>
#include <lte/protos/session_manager.pb.h>
#include <functional>
#include <ostream>
#include <utility>

#include "orc8r/gateway/c/common/service_registry/ServiceRegistrySingleton.hpp"
#include "orc8r/gateway/c/common/logging/magma_logging.hpp"
#include "lte/gateway/c/session_manager/GrpcMagmaUtils.hpp"

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
  MLOG(MDEBUG) << "Sending Set SM Session Response from SMF ";
  PrintGrpcMessage(static_cast<const google::protobuf::Message&>(response));
  auto local_resp = new AsyncLocalResponse<SmContextVoid>(std::move(callback),
                                                          RESPONSE_TIMEOUT);
  local_resp->set_response_reader(stub_->AsyncSetSmfSessionContext(
      local_resp->get_context(), response, &queue_));
  return true;
}

bool AsyncAmfServiceClient::handle_notification_to_access(
    const magma::SetSmNotificationContext& notif) {
  MLOG(MDEBUG) << "Sending Set SM Session Notification from SMF ";
  auto local_resp = new AsyncLocalResponse<SmContextVoid>(std::move(callback),
                                                          RESPONSE_TIMEOUT);
  local_resp->set_response_reader(stub_->AsyncSetAmfNotification(
      local_resp->get_context(), notif, &queue_));
  return true;
}

}  // namespace magma
