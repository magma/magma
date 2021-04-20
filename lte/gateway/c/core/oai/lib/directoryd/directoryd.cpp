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

#include <grpcpp/impl/codegen/status.h>
#include <string>
#include <iostream>

#include "GatewayDirectorydClient.h"
#include "directoryd.h"
#include "orc8r/protos/common.pb.h"
#include "orc8r/protos/directoryd.pb.h"

static void directoryd_rpc_call_done(const grpc::Status& status);

bool directoryd_report_location(char* imsi) {
  // Actual GW_ID will be filled in the cloud
  magma::GatewayDirectoryServiceClient::UpdateRecord(
      "IMSI" + std::string(imsi), std::string("GW_ID"),
      [&](grpc::Status status, magma::Void response) {
        directoryd_rpc_call_done(status);
      });
  return true;
}

bool directoryd_remove_location(char* imsi) {
  magma::GatewayDirectoryServiceClient::DeleteRecord(
      "IMSI" + std::string(imsi),
      [&](grpc::Status status, magma::Void response) {
        directoryd_rpc_call_done(status);
      });
  return true;
}

bool directoryd_update_location(char* imsi, char* location) {
  magma::GatewayDirectoryServiceClient::UpdateRecord(
      "IMSI" + std::string(imsi), std::string(location),
      [&](grpc::Status status, magma::Void response) {
        directoryd_rpc_call_done(status);
      });
  return true;
}

bool directoryd_update_record_field(char* imsi, char* key, char* value) {
  // Actual GW_ID will be filled in the cloud
  magma::GatewayDirectoryServiceClient::UpdateRecordField(
      "IMSI" + std::string(imsi), std::string(key), std::string(value),
      [&](grpc::Status status, magma::Void response) {
        directoryd_rpc_call_done(status);
      });
  return true;
}

void directoryd_rpc_call_done(const grpc::Status& status) {
  if (!status.ok()) {
    std::cerr << "Directoryd RPC failed with code " << status.error_code()
              << ", msg: " << status.error_message() << std::endl;
  }
}
