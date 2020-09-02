#pragma once

#include <grpc++/grpc++.h>
#include <stdint.h>
#include <functional>
#include <memory>

#include "lte/protos/session_manager.grpc.pb.h"
#include "GRPCReceiver.h"


namespace grpc {
class Status;
}  // namespace grpc
namespace magma {
namespace lte {

class SetSMSessionContext;
class SmContextVoid;
}
}

using grpc::Status;

void amf_create_session();
namespace magma {
using namespace lte;

class AMFClient : public GRPCReceiver {
 public:

 static void amf_create_session_final(
  const SetSMSessionContext& request,
  std::function<void(Status, SmContextVoid)> callback);

 public:
  AMFClient(AMFClient const&) = delete;
  void operator=(AMFClient const&) = delete;

 private:
  AMFClient();
  static AMFClient& get_instance();
  std::unique_ptr<AmfPduSessionSmContext::Stub> stub_;
  static const uint32_t RESPONSE_TIMEOUT = 10;  // seconds
 };
 }
