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

#ifdef __cplusplus
extern "C" {
#endif

#include "common_defs.h"
#include "s6a_defs.h"
#include "s6a_messages_types.h"

#ifdef __cplusplus
}
#endif

#include "s6a_client_api.h"
#include "s6a_grpc_iface.h"

//------------------------------------------------------------------------------
S6aGrpcIface::S6aGrpcIface(void) {
  send_activate_messages();
  OAILOG_DEBUG(LOG_S6A, "Initializing S6a interface over gRPC: DONE\n");
}
//------------------------------------------------------------------------------
bool S6aGrpcIface::update_location_req(s6a_update_location_req_t* ulr_p) {
  return s6a_update_location_req(ulr_p);
}
//------------------------------------------------------------------------------
bool S6aGrpcIface::authentication_info_req(s6a_auth_info_req_t* air_p) {
  return s6a_authentication_info_req(air_p);
}
//------------------------------------------------------------------------------
bool S6aGrpcIface::send_cancel_location_ans(s6a_cancel_location_ans_t* cla_pP) {
  return false;
}
//------------------------------------------------------------------------------
bool S6aGrpcIface::purge_ue(const char* imsi) {
  return s6a_purge_ue(imsi);
}
