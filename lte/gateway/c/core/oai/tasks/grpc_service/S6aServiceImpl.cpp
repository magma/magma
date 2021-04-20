/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */
#include <string>

#include "lte/protos/s6a_service.pb.h"

extern "C" {
#include "s6a_service_handler.h"
#include "log.h"
}
#include "S6aServiceImpl.h"

namespace grpc {
class Channel;
class ServerContext;
}  // namespace grpc

using grpc::Channel;
using grpc::ServerContext;
using grpc::Status;
using magma::DeleteSubscriberRequest;
using magma::DeleteSubscriberResponse;
using magma::S6aService;

namespace magma {
using namespace lte;

S6aServiceImpl::S6aServiceImpl() {}

Status S6aServiceImpl::DeleteSubscriber(
    ServerContext* context, const DeleteSubscriberRequest* request,
    DeleteSubscriberResponse* response) {
  auto imsi_size = request->imsi_list_size();
  for (int i = 0; i < imsi_size; i++) {
    auto imsi = request->imsi_list(i);
    OAILOG_INFO(
        LOG_MME_APP, "Sending deleting subscriber %s request\n ", imsi.c_str());
    auto imsi_len = imsi.length();
    delete_subscriber_request(imsi.c_str(), imsi_len);
  }
  return Status::OK;
}
}  // namespace magma
