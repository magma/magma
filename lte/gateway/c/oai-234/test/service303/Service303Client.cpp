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

#include "orc8r/protos/service303.grpc.pb.h"
#include "orc8r/protos/metricsd.pb.h"
#include "orc8r/protos/common.pb.h"

#include "Service303Client.h"

using grpc::Channel;
using grpc::ClientContext;
using grpc::Status;
using magma::Service303Client;
using magma::orc8r::MetricsContainer;
using magma::orc8r::Service303;
using magma::orc8r::ServiceInfo;
using magma::orc8r::Void;

Service303Client::Service303Client(const std::shared_ptr<Channel>& channel)
    : stub_(Service303::NewStub(channel)) {}

int Service303Client::GetServiceInfo(ServiceInfo* response) {
  Void request;
  ClientContext context;
  Status status = stub_->GetServiceInfo(&context, request, response);
  if (!status.ok()) {
    std::cout << "GetServiceInfo fails with code " << status.error_code()
              << ", msg: " << status.error_message() << std::endl;
    return -1;
  }
  return 0;
}

int Service303Client::GetMetrics(MetricsContainer* response) {
  ClientContext context;
  Void request;
  Status status = stub_->GetMetrics(&context, request, response);
  if (!status.ok()) {
    std::cout << "GetMetrics fails with code " << status.error_code()
              << ", msg: " << status.error_message() << std::endl;
    return -1;
  }
  return 0;
}
