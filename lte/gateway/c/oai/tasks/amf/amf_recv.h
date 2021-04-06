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
/*****************************************************************************

  Source      amf_recv.h

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#ifndef AMF_RECV_SEEN
#define AMF_RECV_SEEN

#include <sstream>
//#include "amf_data.h" // optimized
//#include "amf_nas_common_defs.h"
//#include "amf_common_defs.h"
//#include "amf_nas_message.h"
#include "amf_app_ue_context_and_proc.h"  //includes "amf_common_defs.h"
#include "M5GRegistrationAccept.h"
#include "amf_asDefs.h"

using namespace std;

namespace magma5g {
class amf_procedure_handler {
 public:
  int amf_handle_registration_request(
      amf_ue_ngap_id_t ue_id, tai_t* originating_tai, ecgi_t* ecgi,
      RegistrationRequestMsg* msg, const bool is_initial,
      const bool is_amf_ctx_new, int amf_cause,
      amf_nas_message_decode_status_t decode_status);

  int amf_handle_identity_response(
      amf_ue_ngap_id_t ue_id, M5GSMobileIdentityMsg* msg, int amf_cause,
      amf_nas_message_decode_status_t decode_status);

  int amf_handle_authentication_response(
      amf_ue_ngap_id_t ue_id, AuthenticationResponseMsg* msg, int amf_cause,
      amf_nas_message_decode_status_t status);

  int amf_handle_securitycomplete_response(
      amf_ue_ngap_id_t ue_id, amf_nas_message_decode_status_t decode_status);

  int amf_handle_registrationcomplete_response(
      amf_ue_ngap_id_t ue_id, RegistrationCompleteMsg* msg, int amf_cause,
      amf_nas_message_decode_status_t decode_status);

  int amf_handle_deregistration_ue_origin_req(
      amf_ue_ngap_id_t ue_id, DeRegistrationRequestUEInitMsg* msg,
      int amf_cause, amf_nas_message_decode_status_t decode_status);

  int amf_smf_send(
      amf_ue_ngap_id_t ueid, ULNASTransportMsg* msg, int amf_cause);
};

class amf_registration_procedure : public amf_context_t,
                                   public amf_registration_request_ies_t {
 public:
  static int amf_proc_registration_request(
      amf_ue_ngap_id_t ue_id, const bool is_mm_ctx_new,
      amf_registration_request_ies_t* ies);
  static int amf_registration_run_procedure(amf_context_t* amf_context);
  static int amf_registration_success_identification_cb(
      amf_context_t* amf_context);
  static int amf_registration_success_authentication_cb(
      amf_context_t* amf_context);
  static int amf_registration_success_security_cb(amf_context_t* amf_context);
  static int amf_registration(amf_context_t* amf_context);
  int amf_proc_registration_reject(
      const amf_ue_ngap_id_t ue_id, amf_cause_t amf_cause);
  static int amf_registration_reject(
      amf_context_t* amf_context, nas_amf_registration_proc_t* nas_base_proc);
  static int amf_proc_registration_complete(
      amf_ue_ngap_id_t ue_id, bstring smf_msg_p, int amf_cause,
      const amf_nas_message_decode_status_t status);
  static int amf_send_registration_accept_dl_nas(
      const amf_as_data_t* msg, RegistrationAcceptMsg* amf_msg);
  static int amf_send_registration_accept(amf_context_t* amf_context);
};

}  // namespace magma5g
#endif
