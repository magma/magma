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

#pragma once

#include <gmp.h>

#include "feg/protos/s6a_proxy.pb.h"
#include "feg/protos/s6a_proxy.grpc.pb.h"
#include "s6a_messages_types.h"

extern "C" {
#include "intertask_interface.h"
}

namespace magma {
using namespace feg;

AuthenticationInformationRequest
convert_itti_s6a_authentication_info_req_to_proto_msg(
    const s6a_auth_info_req_t* const msg);

UpdateLocationRequest convert_itti_s6a_update_location_request_to_proto_msg(
    const s6a_update_location_req_t* const msg);

}  // namespace magma
