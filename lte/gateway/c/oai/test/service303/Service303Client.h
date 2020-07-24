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

#ifndef SERVICE303_CLIENT_H
#define SERVICE303_CLIENT_H

#include <grpc++/grpc++.h>

#include "orc8r/protos/service303.grpc.pb.h"

using grpc::Channel;
using grpc::ClientContext;
using grpc::Status;
using magma::orc8r::MetricsContainer;
using magma::orc8r::Service303;
using magma::orc8r::ServiceInfo;

namespace magma {

/**
 * gRPC client for Service303
 */
class Service303Client {
 public:
  explicit Service303Client(const std::shared_ptr<Channel>& channel);

  /**
   * Get Service303 Info
   *
   * @param response: a pointer to the ServiceInfo object to populate
   * @return 0 on success, -1 on failure
   */
  int GetServiceInfo(ServiceInfo* response);

  /**
   * Get Metrics from server
   *
   * @param response: the MetricsContainer instance to populate
   * @return 0 on success, -1 on failure
   */
  int GetMetrics(MetricsContainer* response);

 private:
  std::shared_ptr<Service303::Stub> stub_;
};

}  // namespace magma
#endif  // SERVICE303_CLIENT_H
