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
#include <iostream>
#include <string>

#include "sms_orc8r_client_api.h"
#include "SMSOrc8rClient.h"
#include "orc8r/protos/common.pb.h"

extern "C" {
#include "log.h"
}

void void_callback(grpc::Status status, magma::orc8r::Void void_response) {
  return;
}

void send_smo_uplink_unitdata(const itti_sgsap_uplink_unitdata_t* msg) {
  OAILOG_DEBUG(LOG_SMS_ORC8R, "Sending UPLINK_UNITDATA\n");
  magma::SMSOrc8rClient::send_uplink_unitdata(msg, void_callback);
  return;
}
