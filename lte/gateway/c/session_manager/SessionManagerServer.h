/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#pragma once
#include <grpc++/grpc++.h>
#include <lte/protos/session_manager.grpc.pb.h>

#include "LocalSessionManagerHandler.h"
#include "SessionProxyResponderHandler.h"

using grpc::ServerCompletionQueue;
using grpc::ServerContext;
using grpc::Status;

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
class LocalSessionManagerAsyncService final :
  public AsyncService,
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

/**
 * SessionProxyResponderAsyncService handles gRPC calls to SessionProxyResponder
 * through a completion queue where requests are processed and returned async
 */
class SessionProxyResponderAsyncService final :
  public AsyncService,
  public SessionProxyResponder::AsyncService {
 public:
  SessionProxyResponderAsyncService(
    std::unique_ptr<ServerCompletionQueue> cq,
    std::unique_ptr<SessionProxyResponderHandler> handler);

 protected:
  void init_call_data();

 private:
  std::unique_ptr<SessionProxyResponderHandler> handler_;
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
  AsyncGRPCRequest(ServerCompletionQueue *cq, GRPCService &service);

  /**
   * proceed moves to the next step in the state machine. If it's in processing,
   * the call is processed through the session manager handler. If it's finished
   * then this is destroyed
   */
  void proceed();

 protected:
  ServerCompletionQueue *cq_;
  ServerContext ctx_;

  enum CallStatus { PROCESS, FINISH };
  CallStatus status_;

  RequestType request_;
  grpc::ServerAsyncResponseWriter<ResponseType> responder_;

  GRPCService &service_;

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
class ReportRuleStatsCallData :
  public AsyncGRPCRequest<
    LocalSessionManager::AsyncService,
    RuleRecordTable,
    Void> {
 public:
  ReportRuleStatsCallData(
    ServerCompletionQueue *cq,
    LocalSessionManager::AsyncService &service,
    LocalSessionManagerHandler &handler):
    AsyncGRPCRequest(cq, service),
    handler_(handler)
  {
    // By calling RequestReportRuleStats, any RPC to ReportRuleStats will get
    // added to the request queue cq_ with the tag being the memory address
    // of this instance. When the request is completed, it will be added to
    // cq_ again to be finished
    service_.RequestReportRuleStats(
      &ctx_, &request_, &responder_, cq_, cq_, (void *) this);
  }

 protected:
  void clone() override
  {
    // When processing a request, create a new ReportRuleStatsCallData to
    // process another request if it comes in.
    new ReportRuleStatsCallData(cq_, service_, handler_);
  }

  void process() override
  {
    // Get a response from a handler and call the finish callback
    handler_.ReportRuleStats(&ctx_, &request_, get_finish_callback());
  }

 private:
  LocalSessionManagerHandler &handler_;
};

/**
 * Class to handle CreateSession requests
 */
class CreateSessionCallData :
  public AsyncGRPCRequest<
    LocalSessionManager::AsyncService,
    LocalCreateSessionRequest,
    LocalCreateSessionResponse> {
 public:
  CreateSessionCallData(
    ServerCompletionQueue *cq,
    LocalSessionManager::AsyncService &service,
    LocalSessionManagerHandler &handler):
    AsyncGRPCRequest(cq, service),
    handler_(handler)
  {
    service_.RequestCreateSession(
      &ctx_, &request_, &responder_, cq_, cq_, (void *) this);
  }

 protected:
  void clone() override { new CreateSessionCallData(cq_, service_, handler_); }

  void process() override
  {
    handler_.CreateSession(&ctx_, &request_, get_finish_callback());
  }

 private:
  LocalSessionManagerHandler &handler_;
};

/**
 * Class to handle EndSession requests
 */
class EndSessionCallData :
  public AsyncGRPCRequest<
    LocalSessionManager::AsyncService,
    LocalEndSessionRequest,
    LocalEndSessionResponse> {
 public:
  EndSessionCallData(
    ServerCompletionQueue *cq,
    LocalSessionManager::AsyncService &service,
    LocalSessionManagerHandler &handler):
    AsyncGRPCRequest(cq, service),
    handler_(handler)
  {
    service_.RequestEndSession(
      &ctx_, &request_, &responder_, cq_, cq_, (void *) this);
  }

 protected:
  void clone() override { new EndSessionCallData(cq_, service_, handler_); }

  void process() override
  {
    handler_.EndSession(&ctx_, &request_, get_finish_callback());
  }

 private:
  LocalSessionManagerHandler &handler_;
};

/**
 * Class to handle ChargingReauth requests
 */
class ChargingReAuthCallData :
  public AsyncGRPCRequest<
    SessionProxyResponder::AsyncService,
    ChargingReAuthRequest,
    ChargingReAuthAnswer> {
 public:
  ChargingReAuthCallData(
    ServerCompletionQueue *cq,
    SessionProxyResponder::AsyncService &service,
    SessionProxyResponderHandler &handler):
    AsyncGRPCRequest(cq, service),
    handler_(handler)
  {
    service_.RequestChargingReAuth(
      &ctx_, &request_, &responder_, cq_, cq_, (void *) this);
  }

 protected:
  void clone() override { new ChargingReAuthCallData(cq_, service_, handler_); }

  void process() override
  {
    handler_.ChargingReAuth(&ctx_, &request_, get_finish_callback());
  }

 private:
  SessionProxyResponderHandler &handler_;
};

/**
 * Class to handle PolicyReauth requests
 */
class PolicyReAuthCallData :
  public AsyncGRPCRequest<
    SessionProxyResponder::AsyncService,
    PolicyReAuthRequest,
    PolicyReAuthAnswer> {
 public:
  PolicyReAuthCallData(
    ServerCompletionQueue *cq,
    SessionProxyResponder::AsyncService &service,
    SessionProxyResponderHandler &handler):
    AsyncGRPCRequest(cq, service),
    handler_(handler)
  {
    service_.RequestPolicyReAuth(
      &ctx_, &request_, &responder_, cq_, cq_, (void *) this);
  }

 protected:
  void clone() override { new PolicyReAuthCallData(cq_, service_, handler_); }

  void process() override
  {
    handler_.PolicyReAuth(&ctx_, &request_, get_finish_callback());
  }

 private:
  SessionProxyResponderHandler &handler_;
};

} // namespace magma
