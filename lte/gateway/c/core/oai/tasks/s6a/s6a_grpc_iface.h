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

#ifndef S6A_GRPC_IFACE_H_SEEN
#define S6A_GRPC_IFACE_H_SEEN

#include "s6a_viface.h"

class S6aGrpcIface : public S6aViface {
 public:
  S6aGrpcIface();
  bool update_location_req(s6a_update_location_req_t* const ulr_p);
  bool authentication_info_req(s6a_auth_info_req_t* const air_p);
  bool send_cancel_location_ans(s6a_cancel_location_ans_t* cla_pP);
  bool purge_ue(const char* imsi);

  ~S6aGrpcIface(){};
};

#endif /* S6A_GRPC_IFACE_H_SEEN */
