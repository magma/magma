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
#pragma once
#include <grpc++/grpc++.h>
#include <lte/protos/abort_session.grpc.pb.h>
#include <lte/protos/session_manager.grpc.pb.h>

#include <memory>
#include <utility>

#include "LocalSessionManagerHandler.h"
#include "SessionProxyResponderHandler.h"
#include "SetMessageManagerHandler.h"
#include "UpfMsgManageHandler.h"

using grpc::ServerCompletionQueue;
using grpc::ServerContext;
using grpc::Status;
using namespace std;
namespace magma {

/**
 * General async gRPC service. It runs a loop pulling requests off a server
 * queue and then processing them. In general, the C++ implementation of async
 * gRPC servers is quite confusing and requires a lot of groundwork. Hopefully
 * this documentation can clear up some of the "magic".
 *
 * Every request is represented by a CallData, with one implementation per
 * RPC call. The CallData has a state machine of the life cycle of the request.
 * A request appears on the queue when it needs to be *processed* for the first
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

  /**
   * Start the server, blocks
   */
  void wait_for_requests();

  /**
   * Stop the server and shutdown the completion queue
   */
  void stop();

 protected:
  /**
   * Initialize request handlers for report rule stats, end/create session
   */
  virtual void init_call_data() = 0;

 protected:
  std::unique_ptr<ServerCompletionQueue> cq_;
  std::atomic<bool> running_;
};

/**
 * LocalSessionManagerAsyncService handles gRPC calls to LocalSessionManager
 * through a completion queue where requests are processed and returned async
 */
class LocalSessionManagerAsyncService final
    : public AsyncService,
      public LocalSessionManager::AsyncService {
 public:
  LocalSessionManagerAsyncService(
      std::unique_ptr<ServerCompletionQueue> cq,
      std::unique_ptr<LocalSessionManagerHandler> handler);

 protected:
  void init_call_data();

 private:
  std::unique_ptr<LocalSessionManagerHandler> handler_;
};

/*AmfPduSessionSmContextToSmf  Set RPC service object for 5G */
class AmfPduSessionSmContextAsyncService final
    : public AsyncService,
      public AmfPduSessionSmContext::AsyncService {
 public:
  AmfPduSessionSmContextAsyncService(
      std::unique_ptr<ServerCompletionQueue> cq,
      std::unique_ptr<SetMessageManager> handler);

 protected:
  void init_call_data();

 private:
  std::unique_ptr<SetMessageManager> handler_;
};

/* UpfContextAsysnService  Set RPC service object for 5G */
class SetInterfaceForUserPlaneAsyncService final
    : public AsyncService,
      public SetInterfaceForUserPlane::AsyncService {
 public:
  SetInterfaceForUserPlaneAsyncService(
      std::unique_ptr<ServerCompletionQueue> cq,
      std::unique_ptr<UpfMsgManageHandler> handler);

 protected:
  void init_call_data();

 private:
  std::unique_ptr<UpfMsgManageHandler> handler_;
};
/**
 * SessionProxyResponderAsyncService handles gRPC calls to SessionProxyResponder
 * through a completion queue where requests are processed and returned async
 */
class SessionProxyResponderAsyncService final
    : public AsyncService,
      public SessionProxyResponder::AsyncService {
 public:
  SessionProxyResponderAsyncService(
      std::unique_ptr<ServerCompletionQueue> cq,
      std::shared_ptr<SessionProxyResponderHandler> handler);

 protected:
  void init_call_data();

 private:
  std::shared_ptr<SessionProxyResponderHandler> handler_;
};

class AbortSessionResponderAsyncService final
    : public AsyncService,
      public AbortSessionResponder::AsyncService {
 public:
  AbortSessionResponderAsyncService(
      std::unique_ptr<ServerCompletionQueue> cq,
      std::shared_ptr<SessionProxyResponderHandler> handler);

 protected:
  void init_call_data();

 private:
  std::shared_ptr<SessionProxyResponderHandler> handler_;
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
  virtual ~CallData()    = default;
};

/**
 * AsyncGRPCRequest represents a GRPC call through the lifetime of its call.
 * On construction, the state machine is started and the call data is
 * created. When a request is made, the call data moves to processing. When
 * finished, the call data destroys itself
 */
template<class GRPCService, class RequestType, class ResponseType>
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
  std::function<void(Status, ResponseType)> get_finish_callback();
};

/**
 * Class to handle ReportRuleStats requests
 */
class ReportRuleStatsCallData
    : public AsyncGRPCRequest<
          LocalSessionManager::AsyncService, RuleRecordTable, Void> {
 public:
  ReportRuleStatsCallData(
      ServerCompletionQueue* cq, LocalSessionManager::AsyncService& service,
      LocalSessionManagerHandler& handler)
      : AsyncGRPCRequest(cq, service), handler_(handler) {
    // By calling RequestReportRuleStats, any RPC to ReportRuleStats will get
    // added to the request queue cq_ with the tag being the memory address
    // of this instance. When the request is completed, it will be added to
    // cq_ again to be finished
    service_.RequestReportRuleStats(
        &ctx_, &request_, &responder_, cq_, cq_, (void*) this);
  }

 protected:
  void clone() override {
    // When processing a request, create a new ReportRuleStatsCallData to
    // process another request if it comes in.
    new ReportRuleStatsCallData(cq_, service_, handler_);
  }

  void process() override {
    // Get a response from a handler and call the finish callback
    handler_.ReportRuleStats(&ctx_, &request_, get_finish_callback());
  }

 private:
  LocalSessionManagerHandler& handler_;
};
/*Set RPC calldata to invoke first first function of landing object for 5G */
// AmfPduSessionSmContextToSmf
class SetAmfSessionContextCallData : public AsyncGRPCRequest<
                                         AmfPduSessionSmContext::AsyncService,
                                         SetSMSessionContext, SmContextVoid> {
 public:
  SetAmfSessionContextCallData(
      ServerCompletionQueue* cq, AmfPduSessionSmContext::AsyncService& service,
      SetMessageManager& handler)
      : AsyncGRPCRequest(cq, service), handler_(handler) {
    service_.RequestSetAmfSessionContext(
        &ctx_, &request_, &responder_, cq_, cq_, (void*) this);
  }

 protected:
  void clone() override {
    new SetAmfSessionContextCallData(cq_, service_, handler_);
  }

  void process() override {
    handler_.SetAmfSessionContext(&ctx_, &request_, get_finish_callback());
  }

 private:
  SetMessageManager& handler_;
};

/*
 *  Class to handle SetUPFNodeStateCallData
 */
class SetUPFNodeStateCallData
    : public AsyncGRPCRequest<
          SetInterfaceForUserPlane::AsyncService, UPFNodeState, SmContextVoid> {
 public:
  SetUPFNodeStateCallData(
      ServerCompletionQueue* cq,
      SetInterfaceForUserPlane::AsyncService& service,
      UpfMsgManageHandler& handler)
      : AsyncGRPCRequest(cq, service), handler_(handler) {
    service_.RequestSetUPFNodeState(
        &ctx_, &request_, &responder_, cq_, cq_, reinterpret_cast<void*>(this));
  }

 protected:
  void clone() override {
    new SetUPFNodeStateCallData(cq_, service_, handler_);
  }

  void process() override {
    handler_.SetUPFNodeState(&ctx_, &request_, get_finish_callback());
  }

 private:
  UpfMsgManageHandler& handler_;
};

/**
 * Class to handle CreateSession requests
 */
class CreateSessionCallData
    : public AsyncGRPCRequest<
          LocalSessionManager::AsyncService, LocalCreateSessionRequest,
          LocalCreateSessionResponse> {
 public:
  CreateSessionCallData(
      ServerCompletionQueue* cq, LocalSessionManager::AsyncService& service,
      LocalSessionManagerHandler& handler)
      : AsyncGRPCRequest(cq, service), handler_(handler) {
    service_.RequestCreateSession(
        &ctx_, &request_, &responder_, cq_, cq_, (void*) this);
  }

 protected:
  void clone() override { new CreateSessionCallData(cq_, service_, handler_); }

  void process() override {
    handler_.CreateSession(&ctx_, &request_, get_finish_callback());
  }

 private:
  LocalSessionManagerHandler& handler_;
};

/**
 * Class to handle EndSession requests
 */
class EndSessionCallData
    : public AsyncGRPCRequest<
          LocalSessionManager::AsyncService, LocalEndSessionRequest,
          LocalEndSessionResponse> {
 public:
  EndSessionCallData(
      ServerCompletionQueue* cq, LocalSessionManager::AsyncService& service,
      LocalSessionManagerHandler& handler)
      : AsyncGRPCRequest(cq, service), handler_(handler) {
    service_.RequestEndSession(
        &ctx_, &request_, &responder_, cq_, cq_, (void*) this);
  }

 protected:
  void clone() override { new EndSessionCallData(cq_, service_, handler_); }

  void process() override {
    handler_.EndSession(&ctx_, &request_, get_finish_callback());
  }

 private:
  LocalSessionManagerHandler& handler_;
};

/**
 * Class to handle BindPolicy2Bearer requests
 */
class BindPolicy2BearerCallData
    : public AsyncGRPCRequest<
          LocalSessionManager::AsyncService, PolicyBearerBindingRequest,
          PolicyBearerBindingResponse> {
 public:
  BindPolicy2BearerCallData(
      ServerCompletionQueue* cq, LocalSessionManager::AsyncService& service,
      LocalSessionManagerHandler& handler)
      : AsyncGRPCRequest(cq, service), handler_(handler) {
    service_.RequestBindPolicy2Bearer(
        &ctx_, &request_, &responder_, cq_, cq_, (void*) this);
  }

 protected:
  void clone() override {
    new BindPolicy2BearerCallData(cq_, service_, handler_);
  }

  void process() override {
    handler_.BindPolicy2Bearer(&ctx_, &request_, get_finish_callback());
  }

 private:
  LocalSessionManagerHandler& handler_;
};

/**
 * Class to handle UpdateTunnelIds requests
 */

class UpdateTunnelIdsCallData
    : public AsyncGRPCRequest<
          LocalSessionManager::AsyncService, UpdateTunnelIdsRequest,
          UpdateTunnelIdsResponse> {
 public:
  UpdateTunnelIdsCallData(
      ServerCompletionQueue* cq, LocalSessionManager::AsyncService& service,
      LocalSessionManagerHandler& handler)
      : AsyncGRPCRequest(cq, service), handler_(handler) {
    service_.RequestUpdateTunnelIds(
        &ctx_, &request_, &responder_, cq_, cq_, (void*) this);
  }

 protected:
  void clone() override {
    new UpdateTunnelIdsCallData(cq_, service_, handler_);
  }

  void process() override {
    handler_.UpdateTunnelIds(&ctx_, &request_, get_finish_callback());
  }

 private:
  LocalSessionManagerHandler& handler_;
};

/**
 * Class to handle SetSessionRules requests
 */
class SetSessionRulesCallData
    : public AsyncGRPCRequest<
          LocalSessionManager::AsyncService, SessionRules, Void> {
 public:
  SetSessionRulesCallData(
      ServerCompletionQueue* cq, LocalSessionManager::AsyncService& service,
      LocalSessionManagerHandler& handler)
      : AsyncGRPCRequest(cq, service), handler_(handler) {
    service_.RequestSetSessionRules(
        &ctx_, &request_, &responder_, cq_, cq_, (void*) this);
  }

 protected:
  void clone() override {
    new SetSessionRulesCallData(cq_, service_, handler_);
  }

  void process() override {
    handler_.SetSessionRules(&ctx_, &request_, get_finish_callback());
  }

 private:
  LocalSessionManagerHandler& handler_;
};

/**
 * Class to handle AbortSessionRequest requests
 */
class AbortSessionCallData : public AsyncGRPCRequest<
                                 AbortSessionResponder::AsyncService,
                                 AbortSessionRequest, AbortSessionResult> {
 public:
  AbortSessionCallData(
      ServerCompletionQueue* cq, AbortSessionResponder::AsyncService& service,
      SessionProxyResponderHandler& handler)
      : AsyncGRPCRequest(cq, service), handler_(handler) {
    service_.RequestAbortSession(
        &ctx_, &request_, &responder_, cq_, cq_, (void*) this);
  }

 protected:
  void clone() override { new AbortSessionCallData(cq_, service_, handler_); }

  void process() override {
    handler_.AbortSession(&ctx_, &request_, get_finish_callback());
  }

 private:
  SessionProxyResponderHandler& handler_;
};

/**
 * Class to handle ChargingReauth requests
 */
class ChargingReAuthCallData
    : public AsyncGRPCRequest<
          SessionProxyResponder::AsyncService, ChargingReAuthRequest,
          ChargingReAuthAnswer> {
 public:
  ChargingReAuthCallData(
      ServerCompletionQueue* cq, SessionProxyResponder::AsyncService& service,
      SessionProxyResponderHandler& handler)
      : AsyncGRPCRequest(cq, service), handler_(handler) {
    service_.RequestChargingReAuth(
        &ctx_, &request_, &responder_, cq_, cq_, (void*) this);
  }

 protected:
  void clone() override { new ChargingReAuthCallData(cq_, service_, handler_); }

  void process() override {
    handler_.ChargingReAuth(&ctx_, &request_, get_finish_callback());
  }

 private:
  SessionProxyResponderHandler& handler_;
};

/**
 * Class to handle PolicyReauth requests
 */
class PolicyReAuthCallData : public AsyncGRPCRequest<
                                 SessionProxyResponder::AsyncService,
                                 PolicyReAuthRequest, PolicyReAuthAnswer> {
 public:
  PolicyReAuthCallData(
      ServerCompletionQueue* cq, SessionProxyResponder::AsyncService& service,
      SessionProxyResponderHandler& handler)
      : AsyncGRPCRequest(cq, service), handler_(handler) {
    service_.RequestPolicyReAuth(
        &ctx_, &request_, &responder_, cq_, cq_, (void*) this);
  }

 protected:
  void clone() override { new PolicyReAuthCallData(cq_, service_, handler_); }

  void process() override {
    handler_.PolicyReAuth(&ctx_, &request_, get_finish_callback());
  }

 private:
  SessionProxyResponderHandler& handler_;
};

}  // namespace magma
