/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#pragma once

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/include/amf_app_messages_types.h"
#include "lte/protos/subscriberauth.pb.h"
#include "lte/protos/subscriberdb.pb.h"

using magma::lte::M5GAuthenticationInformationAnswer;
using magma::lte::M5GSUCIRegistrationAnswer;

namespace magma5g {

void convert_proto_msg_to_itti_m5g_auth_info_ans(
    M5GAuthenticationInformationAnswer msg,
    itti_amf_subs_auth_info_ans_t* itti_msg);

void convert_proto_msg_to_itti_amf_decrypted_imsi_info_ans(
    M5GSUCIRegistrationAnswer response,
    itti_amf_decrypted_imsi_info_ans_t* amf_app_decrypted_imsi_info_resp);

}  // namespace magma5g
