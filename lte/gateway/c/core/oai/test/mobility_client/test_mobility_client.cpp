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
#include "MobilityServiceClient.h"

#include <arpa/inet.h>
#include <grpcpp/impl/codegen/status.h>
#include <stdio.h>

#include "lte/protos/mobilityd.grpc.pb.h"
#include "lte/protos/mobilityd.pb.h"

using grpc::Status;
using magma::lte::AllocateIPAddressResponse;
using magma::lte::IPAddress;
using magma::lte::MobilityServiceClient;

int main(int argc, char** argv) {
  int status;
  struct in_addr tmp;
  struct in_addr ipv4_addr3;

  char apn[10] = "oai.ipv4";
  {
    char str[INET_ADDRSTRLEN];
    printf("Allocating IP address...\n");

    MobilityServiceClient::getInstance().AllocateIPv4AddressAsync(
        "0001", apn,
        [apn, &str, &tmp](
            const Status& status, AllocateIPAddressResponse ip_msg) {
          struct in_addr ipv4_addr1;
          memcpy(
              &ipv4_addr1,
              ip_msg.mutable_ip_list(0)->mutable_address()->c_str(),
              sizeof(in_addr));

          if (!status.ok()) {
            printf(
                "allocate_ipv4_address error %d for sid %s for apn %s\n",
                status.error_code(), "0001", apn);
            return -1;
          }
          tmp.s_addr = htonl(ipv4_addr1.s_addr);
          inet_ntop(AF_INET, &tmp, str, INET_ADDRSTRLEN);
          printf("IP allocated: %s\n", str);
          printf("Releasing IP address...\n");

          int release_status =
              MobilityServiceClient::getInstance().ReleaseIPv4Address(
                  "0001", apn, ipv4_addr1);
          if (release_status) {
            printf(
                "release_ipv4_address error %d for sid %s\n", release_status,
                "0001");
            return -1;
          }
          return 0;
        });

    MobilityServiceClient::getInstance().AllocateIPv4AddressAsync(
        "0002", apn,
        [apn, &str, &tmp](
            const Status& status, AllocateIPAddressResponse ip_msg) {
          struct in_addr ipv4_addr2;
          memcpy(
              &ipv4_addr2,
              ip_msg.mutable_ip_list(0)->mutable_address()->c_str(),
              sizeof(in_addr));
          if (!status.ok()) {
            printf(
                "allocate_ipv4_address error %d for sid %s for apn %s\n",
                status.error_code(), "0002", apn);
            return -1;
          }
          tmp.s_addr = htonl(ipv4_addr2.s_addr);
          inet_ntop(AF_INET, &tmp, str, INET_ADDRSTRLEN);
          printf("IP allocated: %s\n", str);
          int release_status2 =
              MobilityServiceClient::getInstance().ReleaseIPv4Address(
                  "0002", apn, ipv4_addr2);
          if (release_status2) {
            printf(
                "release_ipv4_address error %d for sid %s\n", release_status2,
                "0002");
            return -1;
          }
          return 0;
        });

    status = MobilityServiceClient::getInstance().GetIPv4AddressForSubscriber(
        "0002", apn, &ipv4_addr3);
    if (status) {
      printf(
          "get_ipv4_address_for_subscriber error %d for sid %s\n", status,
          "0002");
      return -1;
    }
    tmp.s_addr = htonl(ipv4_addr3.s_addr);
    inet_ntop(AF_INET, &tmp, str, INET_ADDRSTRLEN);
    printf("Retrieved ip address for user 0002: %s\n", str);
  }

  printf("All tests passed...\n");
}
