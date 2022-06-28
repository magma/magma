/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

#pragma once

#include <grpc++/grpc++.h>
#include <grpcpp/impl/codegen/status.h>
#include "feg/protos/s6a_proxy.grpc.pb.h"

#include <memory>
#include <atomic>

namespace grpc {
class ServerContext;
}  // namespace grpc

namespace magma {
namespace feg {
class CancelLocationRequest;
class CancelLocationAnswer;
class ResetRequest;
class ResetAnswer;
}  // namespace feg
}  // namespace magma

using grpc::ServerCompletionQueue;
using grpc::ServerContext;
using grpc::Status;

namespace magma {

using namespace feg;

/**
 * General async gRPC service. It runs a loop pulling requests off a server
 * queue and then processing them. In general, the C++ implementation of async
 * gRPC servers is quite confusing and requires a lot of groundwork. Hopefully
 * this documentation can clear up some of the "magic".
 *
 * Every request is represented by a CallData, with one implementation per
 * RPC call. The CallData has a state machine of the life cycle of the request.
 * A request appears on the queue when it needs to be processed for the first
 * time and then when it needs to be finished (i.e. sent back as an answer).
 *
 * When the AsyncServer is first initialized, the CallData for each RPC must
 * be initialized. When a CallData is created for a particular RPC, the
 * server will mark an incoming request with a tag, and in this case the tag is
 * a pointer to the CallData itself. When the CallData is then processed, a new
 * CallData for that RPC is created and the next request will have that tag.
 */
class AsyncService {
 public:
  AsyncService(std::unique_ptr<ServerCompletionQueue> cq);
  virtual ~AsyncService() = default;

  // Start the server, blocks
  void wait_for_requests();

  // Stop the server and shutdown the completion queue
  void stop();

 protected:
  // Initialize request handlers for report rule stats, end/create session
  virtual void init_call_data() = 0;

 protected:
  std::unique_ptr<ServerCompletionQueue> cq_;
  std::atomic<bool> running_;
};

class S6aProxyAsyncResponderHandler {
 public:
  S6aProxyAsyncResponderHandler() {}

  //  Cancel Location Request is sent from HSS to delete subscriber
  void CancelLocation(ServerContext* context,
                      const CancelLocationRequest* request,
                      std::function<void(grpc::Status, CancelLocationAnswer)>
                          response_callback);
  //  Reset Request is sent from HSS to reset some or all subscribers
  void Reset(ServerContext* context, const ResetRequest* request,
             std::function<void(grpc::Status, ResetAnswer)> response_callback);
};

/*
 * S6aProxyResponderAsyncService handles gRPC calls to S6aGatewayService
 * through a completion queue where requests are processed and returns async
 * response
 */
class S6aProxyResponderAsyncService final
    : public AsyncService,
      public S6aGatewayService::AsyncService {
 private:
  static S6aProxyResponderAsyncService* s6a_async_service;

  S6aProxyResponderAsyncService(
      std::unique_ptr<ServerCompletionQueue> cq,
      std::shared_ptr<S6aProxyAsyncResponderHandler> handler);

 public:
  std::thread async_grpc_message_receiver_thread;
  static S6aProxyResponderAsyncService* getInstance() {
    if (!s6a_async_service) {
      s6a_async_service = new S6aProxyResponderAsyncService(nullptr, nullptr);
    }
    return s6a_async_service;
  }

  void set_callback(std::unique_ptr<ServerCompletionQueue> cq,
                    std::shared_ptr<S6aProxyAsyncResponderHandler> handler);

 protected:
  void init_call_data(void);

 private:
  std::shared_ptr<S6aProxyAsyncResponderHandler> handler_;
};

/**
 * Interface for capturing request call processing. This is required because
 * a templated class cannot be declared without knowing the type. Using this
 * class, we can arbitrarily keep a list of CallData objects of any kind and
 * call proceed on them
 */
class CallData {
 public:
  /**
   * proceed is called by the AsyncService loop when the CallData is found in
   * the queue
   */
  virtual void proceed() = 0;
  virtual ~CallData() = default;
};

/**
 * AsyncGRPCRequest represents a GRPC call through the lifetime of its call.
 * On construction, the state machine is started and the call data is
 * created. When a request is made, the call data moves to processing. When
 * finished, the call data destroys itself
 */
template <class GRPCService, class RequestType, class ResponseType>
class AsyncGRPCRequest : public CallData {
 public:
  AsyncGRPCRequest(ServerCompletionQueue* cq, GRPCService& service);
  virtual ~AsyncGRPCRequest() = default;

  /**
   * proceed moves to the next step in the state machine. If it's in processing,
   * the call is processed through the session manager handler. If it's finished
   * then this is destroyed
   */
  void proceed();

 protected:
  ServerCompletionQueue* cq_;
  ServerContext ctx_;

  enum CallStatus { PROCESS, FINISH };
  CallStatus status_;

  RequestType request_;
  grpc::ServerAsyncResponseWriter<ResponseType> responder_;

  GRPCService& service_;

 protected:
  /**
   * Create new CallData for the request so that more requests can be processed
   */
  virtual void clone() = 0;

  /**
   * Process the request and set a response
   */
  virtual void process() = 0;

  /**
   * Returns a callback function to pass to complete the request with the server
   */
  std::function<void(grpc::Status, ResponseType)> get_finish_callback();
};

/**
 * Class to handle Cancel Location requests
 */
class CancelLocationCallData
    : public AsyncGRPCRequest<S6aGatewayService::AsyncService,
                              CancelLocationRequest, CancelLocationAnswer> {
 public:
  CancelLocationCallData(ServerCompletionQueue* cq,
                         S6aGatewayService::AsyncService& service,
                         S6aProxyAsyncResponderHandler& handler)
      : AsyncGRPCRequest(cq, service), handler_(handler) {
    service_.RequestCancelLocation(&ctx_, &request_, &responder_, cq_, cq_,
                                   (void*)this);
  }

 protected:
  void clone() override { new CancelLocationCallData(cq_, service_, handler_); }

  void process() override {
    handler_.CancelLocation(&ctx_, &request_, get_finish_callback());
  }

 private:
  S6aProxyAsyncResponderHandler& handler_;
};

/**
 * Class to handle Reset requests
 */
class ResetCallData : public AsyncGRPCRequest<S6aGatewayService::AsyncService,
                                              ResetRequest, ResetAnswer> {
 public:
  ResetCallData(ServerCompletionQueue* cq,
                S6aGatewayService::AsyncService& service,
                S6aProxyAsyncResponderHandler& handler)
      : AsyncGRPCRequest(cq, service), handler_(handler) {
    service_.RequestReset(&ctx_, &request_, &responder_, cq_, cq_, (void*)this);
  }

 protected:
  void clone() override { new ResetCallData(cq_, service_, handler_); }

  void process() override {
    handler_.Reset(&ctx_, &request_, get_finish_callback());
  }

 private:
  S6aProxyAsyncResponderHandler& handler_;
};

}  // namespace magma
