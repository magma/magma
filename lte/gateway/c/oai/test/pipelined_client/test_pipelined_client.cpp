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
#include "PipelinedServiceClient.h"

#include <arpa/inet.h>
#include <grpcpp/impl/codegen/status.h>
#include <stdio.h>

#include "lte/protos/pipelined.grpc.pb.h"
#include "lte/protos/pipelined.pb.h"

using grpc::Status;
using magma::lte::IPAddress;
using magma::lte::PipelinedServiceClient;
using magma::lte::UESessionContextResponse;
using magma::lte::UESessionSet;

int main(int argc, char** argv) {
  int status;
  struct in_addr tmp;
  struct in_addr ipv4_addr3;

  char apn[10] = "oai.ipv4";
  {
    char str[INET_ADDRSTRLEN];
    struct in_addr ue_ipv4_addr;
    struct in_addr enb_ipv4_addr;
    struct in6_addr ue_ipv6_addr;
    struct ip_flow_dl flow_dl;
    int ret = 0;

    ret = inet_pton(AF_INET, "192.168.128.11", &ue_ipv4_addr);
    if (ret < 0) {
      printf("Failed to allocate ue in_addr");
      return -1;
    }

    ret = inet_pton(AF_INET, "192.168.60.141", &enb_ipv4_addr);
    if (ret < 0) {
      printf("Failed to allocate enb in_addr");
      return -1;
    }

    printf("Default UE IPv4 without flow dl...\n");

    PipelinedServiceClient::get_instance().UpdateUEIPv4SessionSet(
        ue_ipv4_addr, 0, enb_ipv4_addr, 100, 200, "IMSI001222333", 10, apn,
        UE_SESSION_ACTIVE_STATE,
        [&](const Status& status, UESessionContextResponse response) {
          if (!status.ok()) {
            printf(
                "UpdateUEIPv4SessionSet error %d for sid %s for apn %s\n",
                status.error_code(), "IMSI001222333", apn);
            return -1;
          }
        });

    printf("UE IPv4, IPV6  with flow dl...\n");

    ret = inet_pton(AF_INET6, "2001::8", &ue_ipv6_addr);
    if (ret < 0) {
      printf("Failed to allocate ue ipv6 in6_addr");
      return -1;
    }

    memset(&flow_dl, 0, sizeof(struct ip_flow_dl));
    flow_dl.set_params   = 0x71;
    flow_dl.tcp_dst_port = 5002;
    flow_dl.tcp_src_port = 60;
    flow_dl.ip_proto     = 6;
    inet_pton(AF_INET, "192.168.128.9", &(flow_dl.dst_ip));

    PipelinedServiceClient::get_instance().UpdateUEIPv4v6SessionSetWithFlowdl(
        ue_ipv4_addr, ue_ipv6_addr, 0, enb_ipv4_addr, 100, 200, "IMSI001222333",
        flow_dl, 10, apn, UE_SESSION_ACTIVE_STATE,
        [&](const Status& status, UESessionContextResponse response) {
          if (!status.ok()) {
            printf(
                "UpdateUEIPv4v6SessionSetWithFlowdl error %d for sid %s for "
                "apn %s\n",
                status.error_code(), "IMSI001222333", apn);
            return -1;
          }
        });
    printf("All tests passed...\n");
  }
}
