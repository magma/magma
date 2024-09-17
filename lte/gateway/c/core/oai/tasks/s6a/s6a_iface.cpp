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
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/include/s6a_messages_types.hpp"
#include "lte/gateway/c/core/oai/tasks/s6a/s6a_c_iface.hpp"
#if S6A_OVER_GRPC
#include "lte/gateway/c/core/oai/tasks/s6a/s6a_grpc_iface.hpp"
#else
#include "lte/gateway/c/core/oai/tasks/s6a/s6a_fd_iface.hpp"
#endif

#include <new>
#include <exception>

S6aViface* s6a_interface = nullptr;

//------------------------------------------------------------------------------
bool s6a_viface_open(const s6a_config_t* config) {
  if (!s6a_interface) {
#if S6A_OVER_GRPC
    s6a_interface = new S6aGrpcIface();
#else
    s6a_interface = new S6aFdIface(config);
#endif
  }
  return true;
}

//------------------------------------------------------------------------------
void s6a_viface_close() {
  if (s6a_interface) {
    delete s6a_interface;
    s6a_interface = nullptr;
  }
}

//------------------------------------------------------------------------------
bool s6a_viface_update_location_req(s6a_update_location_req_t* ulr_p) {
  if (s6a_interface) {
    return s6a_interface->update_location_req(ulr_p);
  }
  return false;
}

//------------------------------------------------------------------------------
bool s6a_viface_authentication_info_req(s6a_auth_info_req_t* air_p) {
  if (s6a_interface) {
    return s6a_interface->authentication_info_req(air_p);
  }
  return false;
}

//------------------------------------------------------------------------------
bool s6a_viface_send_cancel_location_ans(s6a_cancel_location_ans_t* cla_pP) {
  if (s6a_interface) {
    return s6a_interface->send_cancel_location_ans(cla_pP);
  }
  return false;
}
//------------------------------------------------------------------------------
bool s6a_viface_purge_ue(const char* imsi) {
  if (s6a_interface) {
    return s6a_interface->purge_ue(imsi);
  }
  return false;
}
