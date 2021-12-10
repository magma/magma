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

#include <grpcpp/impl/codegen/client_context.h>    // for ClientContext
#include <grpcpp/impl/codegen/completion_queue.h>  // for CompletionQueue
#include <grpcpp/impl/codegen/status.h>            // for Status
#include <stdint.h>                                // for uint32_t
#include <atomic>                                  // for atomic
#include <chrono>                                  // for operator+, seconds
#include <functional>                              // for function
#include <memory>                                  // for unique_ptr
namespace grpc {
template <class R>
class ClientAsyncResponseReader;
}

namespace magma {

/**
 * GRPCReceiver is the base class for receiving responses asynchronously from
 * the cloud. It uses a completion queue to wait for new responses, and call
 * the virtual handle_response callback on them
 */
class GRPCReceiver {
 public:
  /**
   * Begin the receiver loop, blocks
   */
  void rpc_response_loop();

  /**
   * Stop the receiver loop
   */
  void stop();

 protected:
  grpc::CompletionQueue queue_;

 private:
  std::atomic<bool> running_;
};

/**
 * AsyncResponse is the base class that all tags in the completion queue will
 * be cast to.
 */
class AsyncResponse {
 public:
  virtual ~AsyncResponse() = default;
  /**
   * Override handle_response to be called when a response comes into the queue
   */
  virtual void handle_response() = 0;
};

/**
 * AsyncGRPCResponse is the templatized response to be used for any gRPC call.
 * Setting the response reader will automatically call Finish with 'this' as
 * the tag.
 * The only thing that needs to be overridden is handle_response.
 *
 * This class is implemented here because non-specialized templates
 * must be visible to a translation unit
 */
template <typename ResponseType>
class AsyncGRPCResponse : public AsyncResponse {
 public:
  AsyncGRPCResponse(std::function<void(grpc::Status, ResponseType)> callback,
                    uint32_t timeout_sec)
      : callback_(callback) {
    context_.set_deadline(std::chrono::system_clock::now() +
                          std::chrono::seconds(timeout_sec));
  }
  virtual ~AsyncGRPCResponse() = default;

  virtual void handle_response() {}

  /**
   * Set the response reader which waits for the response back from the gRPC
   * call
   */
  void set_response_reader(
      std::unique_ptr<grpc::ClientAsyncResponseReader<ResponseType>> reader) {
    response_reader_ = std::move(reader);
    response_reader_->Finish(&response_, &status_, this);
  }

  /**
   * Helper function to retrieve the client context
   */
  grpc::ClientContext* get_context() { return &context_; }

 protected:
  // callback on completion
  std::function<void(grpc::Status, ResponseType)> callback_;
  ResponseType response_;  // response from the cloud
  grpc::ClientContext context_;
  grpc::Status status_;
  std::unique_ptr<grpc::ClientAsyncResponseReader<ResponseType>>
      response_reader_;
};

/**
 * AsyncLocalResponse is an example provided response that takes the callback
 * passed and executes it directly in the response loop's thread
 * It is important that when using this, the callback can be executed quickly,
 * because it blocks the response queue.
 * Here is an example usage:
 * auto response = new AsyncLocalResponse<YourRPCResponseValue>(
 *   callback, RESPONSE_TIMEOUT);
 * auto response_reader = stub_->AsyncYourRPCCall(
 *   local_response->get_context(), request_val, &completion_queue);
 * local_response->set_response_reader(std::move(response_reader));
 */
template <typename ResponseType>
class AsyncLocalResponse : public AsyncGRPCResponse<ResponseType> {
 public:
  AsyncLocalResponse(std::function<void(grpc::Status, ResponseType)> callback,
                     uint32_t timeout_sec)
      : AsyncGRPCResponse<ResponseType>(callback, timeout_sec) {}

  void handle_response() {
    this->callback_(this->status_, this->response_);
    delete this;
  }
};

}  // namespace magma
