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
#include <sstream>
#include <thread>
namespace magma5g {

// Authentication related procedure
typedef struct nas5g_amf_auth_proc_s {
  nas_amf_common_proc_t amf_com_proc;
  nas5g_timer_t T3560; /* Authentication timer         */
  unsigned int retransmission_count;
  amf_ue_ngap_id_t ue_id;
  bool is_cause_is_registered;  //  could also be done by seeking parent
                                //  procedure
  ksi_t ksi;
  uint8_t rand[AUTH_RAND_SIZE]; /* Random challenge number  */
  uint8_t autn[AUTH_AUTN_SIZE]; /* Authentication token     */
  int amf_cause;
} nas5g_amf_auth_proc_t;

typedef struct nas5g_auth_info_proc_s {
  nas5g_cn_proc_t cn_proc;
  bool request_sent;
  amf_ue_ngap_id_t ue_id;
  bool resync;  // Indicates whether the authentication information is requested
                // due to sync failure
} nas5g_auth_info_proc_t;

nas5g_auth_info_proc_t* nas5g_cn_auth_info_procedure(
    amf_context_t* const amf_context);

nas5g_auth_info_proc_t* get_nas5g_cn_procedure_auth_info(
    const amf_context_t* ctxt);

int amf_proc_authentication_ksi(
    amf_context_t* amf_context,
    nas_amf_specific_proc_t* const amf_specific_proc, ksi_t ksi,
    const uint8_t* const rand, const uint8_t* const autn, success_cb_t success,
    failure_cb_t failure);
int amf_proc_authentication(
    amf_context_t* amf_context,
    nas_amf_specific_proc_t* const amf_specific_proc, success_cb_t success,
    failure_cb_t failure);
int amf_proc_authentication_complete(
    amf_ue_ngap_id_t ue_id, AuthenticationResponseMsg* msg, int amf_cause,
    const unsigned char* res);
int amf_proc_authentication_failure(
    amf_ue_ngap_id_t ue_id, AuthenticationFailureMsg* msg, int amf_cause);
int amf_registration_security(amf_context_t* amf_context);
int amf_send_authentication_request(
    amf_context_t* amf_context, nas5g_amf_auth_proc_t* auth_proc);
}  // namespace magma5g
