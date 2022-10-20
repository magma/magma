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

#define grpc_async_service
#define grpc_async_service_TASK_C

#include <thread>
#include "orc8r/gateway/c/common/service303/MagmaService.hpp"
#include "lte/gateway/c/core/oai/tasks/async_grpc_service/grpc_async_service_task.hpp"
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/include/s6a_service_handler.hpp"

static void grpc_async_service_exit(void);
task_zmq_ctx_t grpc_async_service_task_zmq_ctx;

namespace magma {
#define S6A_ASYNC_PROXY_SERVICE "s6a_async_service"
#define S6A_ASYNC_PROXY_VERSION "1.0"

magma::S6aProxyResponderAsyncService*
    S6aProxyResponderAsyncService::s6a_async_service = 0;

magma::service303::MagmaService s6a_async_grpc_server(S6A_ASYNC_PROXY_SERVICE,
                                                      S6A_ASYNC_PROXY_VERSION);

void stop_async_s6a_grpc_server(void) {
  magma::S6aProxyResponderAsyncService* s6a_proxy_async_instance =
      magma::S6aProxyResponderAsyncService::getInstance();
  magma::s6a_async_grpc_server.Stop();  // Stop async grpc server
  s6a_proxy_async_instance->stop();     // stop queue after server shuts down
  s6a_proxy_async_instance->async_grpc_message_receiver_thread.join();
  delete s6a_proxy_async_instance;
}

AsyncService::AsyncService(std::unique_ptr<ServerCompletionQueue> cq)
    : cq_(std::move(cq)) {}

void AsyncService::wait_for_requests() {
  init_call_data();
  void* tag;
  bool ok;
  running_ = true;
  while (running_) {
    if (cq_ == nullptr) {
      OAILOG_CRITICAL(LOG_S6A, "Completion queue on s6a interface is null \n");
      return;
    }
    if (!(cq_->Next(&tag, &ok))) {
      OAILOG_ERROR(LOG_S6A,
                   "Completion queue on s6a interface is shutting down \n");
      return;
    }
    if (!ok) {
      OAILOG_ERROR(LOG_S6A, "encountered error while processing request");
      delete static_cast<CallData*>(tag);
      continue;
    }
    static_cast<CallData*>(tag)->proceed();
  }
  return;
}

void AsyncService::stop() {
  running_ = false;
  cq_->Shutdown();
  // Pop all items in the queue until it is empty
  void* tag;
  bool ok;
  while (cq_->Next(&tag, &ok)) {
    OAILOG_INFO(LOG_S6A,
                "CQ is shutdown, pop out all messages from CQ and release the "
                "messages ");
    delete static_cast<CallData*>(tag);
  }
}

S6aProxyResponderAsyncService::S6aProxyResponderAsyncService(
    std::unique_ptr<ServerCompletionQueue> cq,
    std::shared_ptr<S6aProxyAsyncResponderHandler> handler)
    : AsyncService(std::move(cq)), handler_(handler) {}

void S6aProxyResponderAsyncService::init_call_data(void) {
  new CancelLocationCallData(cq_.get(), *this, *handler_);
  new ResetCallData(cq_.get(), *this, *handler_);
}

void S6aProxyResponderAsyncService::set_callback(
    std::unique_ptr<ServerCompletionQueue> cq,
    std::shared_ptr<S6aProxyAsyncResponderHandler> handler) {
  this->cq_ = std::move(cq);
  this->handler_ = handler;
}

template <class GRPCService, class RequestType, class ResponseType>
AsyncGRPCRequest<GRPCService, RequestType, ResponseType>::AsyncGRPCRequest(
    ServerCompletionQueue* cq, GRPCService& service)
    : cq_(cq), status_(PROCESS), responder_(&ctx_), service_(service) {}

// Internal state management logic for AsyncGRPCRequest
// Once a request has started processing, create a new AsyncGRPCRequest to
// standby for new requests. After a request has finished processing, delete the
// object.
template <class GRPCService, class RequestType, class ResponseType>
void AsyncGRPCRequest<GRPCService, RequestType, ResponseType>::proceed() {
  if (status_ == PROCESS) {
    clone();  // Create another stand by CallData
    process();
    status_ = FINISH;
  } else {
    GPR_ASSERT(status_ == FINISH);
    delete this;
  }
}

template <class GRPCService, class RequestType, class ResponseType>
std::function<void(grpc::Status, ResponseType)> AsyncGRPCRequest<
    GRPCService, RequestType, ResponseType>::get_finish_callback() {
  return [this](grpc::Status status, ResponseType response) {
    responder_.Finish(response, status, reinterpret_cast<void*>(this));
  };
}

void S6aProxyAsyncResponderHandler::CancelLocation(
    ServerContext* context, const CancelLocationRequest* request,
    std::function<void(grpc::Status, CancelLocationAnswer)> response_callback) {
  CancelLocationAnswer ans;
  Status status;
  status = Status::OK;
  ans.set_error_code(ErrorCode::SUCCESS);
  if (response_callback) {
    response_callback(status, ans);
  }
  auto imsi = request->user_name();
  auto cancellation_type = request->cancellation_type();
  OAILOG_INFO(LOG_S6A, "Received CLR for %s of type %d\n ", imsi.c_str(),
              cancellation_type);
  if (cancellation_type == CancelLocationRequest::SUBSCRIPTION_WITHDRAWAL) {
    auto imsi_len = imsi.length();
    delete_subscriber_request(imsi.c_str(), imsi_len);
  }
  return;
}

void S6aProxyAsyncResponderHandler::Reset(
    ServerContext* context, const ResetRequest* request,
    std::function<void(grpc::Status, ResetAnswer)> response_callback) {
  ResetAnswer ans;
  Status status;
  status = Status::OK;
  ans.set_error_code(ErrorCode::SUCCESS);
  if (response_callback) {
    response_callback(status, ans);
  }
  OAILOG_INFO(LOG_S6A, "Received S6a-Reset message \n");
  // Send message to MME_APP for furher processing
  handle_reset_request();
  return;
}

}  // namespace magma

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  MessageDef* received_message_p = receive_msg(reader);

  switch (ITTI_MSG_ID(received_message_p)) {
    case TERMINATE_MESSAGE:
      free(received_message_p);
      grpc_async_service_exit();
      break;
    default:
      OAILOG_DEBUG(LOG_UTIL, "Unknown message ID %d: %s\n",
                   ITTI_MSG_ID(received_message_p),
                   ITTI_MSG_NAME(received_message_p));
      break;
  }
  free(received_message_p);
  return 0;
}

static void* grpc_async_service_thread(__attribute__((unused)) void* args) {
  itti_mark_task_ready(TASK_ASYNC_GRPC_SERVICE);
  const task_id_t tasks[] = {TASK_MME_APP, TASK_S6A};

  init_task_context(TASK_ASYNC_GRPC_SERVICE, tasks, 2, handle_message,
                    &grpc_async_service_task_zmq_ctx);
  magma::S6aProxyResponderAsyncService* s6a_proxy_async_instance =
      magma::S6aProxyResponderAsyncService::getInstance();
  auto async_service_handler =
      std::make_shared<magma::S6aProxyAsyncResponderHandler>();
  s6a_proxy_async_instance->set_callback(
      std::move(magma::s6a_async_grpc_server.GetNewCompletionQueue()),
      async_service_handler);
  magma::s6a_async_grpc_server.AddServiceToServer(s6a_proxy_async_instance);

  magma::s6a_async_grpc_server.Start();
  OAILOG_INFO(LOG_S6A, "Started async grpc server for s6a interface \n");
  s6a_proxy_async_instance->async_grpc_message_receiver_thread =
      std::thread([&]() {
        s6a_proxy_async_instance
            ->wait_for_requests();  // block here instead of on server
      });
  zloop_start(grpc_async_service_task_zmq_ctx.event_loop);
  AssertFatal(
      0,
      "Asserting as async grpc_service_thread should not be exiting on its "
      "own!");
  return NULL;
}

extern "C" int grpc_async_service_init(void) {
  OAILOG_DEBUG(LOG_UTIL, "Initializing async_grpc_service task interface\n");

  if (itti_create_task(TASK_ASYNC_GRPC_SERVICE, &grpc_async_service_thread,
                       NULL) < 0) {
    OAILOG_ALERT(LOG_UTIL, "Initializing async_grpc_service: ERROR\n");
    return RETURNerror;
  }
  return RETURNok;
}

static void grpc_async_service_exit(void) {
  destroy_task_context(&grpc_async_service_task_zmq_ctx);
  magma::stop_async_s6a_grpc_server();
  OAI_FPRINTF_INFO("TASK_ASYNC_GRPC_SERVICE terminated\n");
  pthread_exit(NULL);
}
