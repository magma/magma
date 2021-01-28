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

#include <stdint.h>
#include <stdbool.h>
#include <stdlib.h>
#include <folly/Format.h>
#include <folly/json.h>
#include <folly/dynamic.h>
#include <grpcpp/support/status.h>

#include "mme_events.h"
#include "conversions.h"
#include "bstrlib.h"
#include "log.h"
#include "dynamic_memory_check.h"
#include "assertions.h"
#include "common_types.h"
#include "amf_data.h"
#include "amf_app_ue_context_and_proc.h"
#include "amf_app_defs.h"
#include "amf_sap.h"
//#include "smf_sap.h"
#include "service303.h"
//#include "M5GDeRegistrationRequest.h"
#include "common_defs.h"
#include "amf_asDefs.h"
#include "amf_fsm.h"
//#include "smf_sapDef.h"
//#include "smf_data.h"
//#include "amf_api.h"
//#include "nas_procedures.h"

namespace magma5g {
constexpr char AMF_STREAM_NAME[]        = "amf";
constexpr char DEREGISTRATION_SUCCESS[] = "deregistration_success";
/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

/* String representation of the deregistration type */
static const char* _amf_deregistration_type_str[] = {
    "3GPP",
    "NON3GPP",
    "3GPP/NON3GPP",
    "RE-REGISTRATION REQUIRED",
    "RE-REGISTRATION NOT REQUIRED",
    "RESERVED"};

typedef enum {
  AMF_UE_INITIATED_3GPP_DEREGISTRATION = 0,
  AMF_UE_INITIATED_NON3GPP_DEREGISTRATION,
  AMF_UE_INITIATED_3GPP_NON3GPP_DEREGISTRATION,
} amf_proc_access_type_t;

typedef struct amf_deregistration_request_ies_s {
  amf_proc_access_type_t type;
  bool switch_off;
  bool is_native_sc;
  ksi_t ksi;
  guti_m5_t* guti;
  imsi_t* imsi;
} amf_deregistration_request_ies_t;

static int report_event(
    folly::dynamic& event_value, const std::string& event_type,
    const std::string& stream_name, const std::string& event_tag) {
  Event event_request = Event();
  event_request.set_event_type(event_type);
  event_request.set_stream_name(stream_name);

  std::string event_value_string = folly::toJson(event_value);
  event_request.set_value(event_value_string);
  event_request.set_tag(event_tag);
  int rc = log_event(event_request);
  return rc;
}

int deregistration_success_event(imsi64_t imsi64, const char* action) {
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(imsi64, (char*) imsi_str, IMSI_BCD_DIGITS_MAX);

  folly::dynamic event_value = folly::dynamic::object;
  event_value["imsi"]        = imsi_str;
  event_value["action"]      = action;

  return report_event(
      event_value, DEREGISTRATION_SUCCESS, AMF_STREAM_NAME, imsi_str);
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_proc_deregistration_request()                             **
 **                                                                        **
 ***************************************************************************/

int amf_proc_deregistration_request(
    amf_ue_ngap_id_t ue_id, amf_deregistration_request_ies_t* params) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc;

  // OAILOG_INFO(
  //    LOG_NAS_AMF,
  //    "AMF-PROC  - De_Registration type = %s (%d) requested"
  //    " (ue_id=" AMF_UE_NGAP_ID_FMT ")\n",
  //    _amf_deregistration_type_str[params->type], params->type, ue_id);

  amf_context_t* amf_ctx = amf_context_get(&_amf_data, ue_id);
  if (params->switch_off) {
    increment_counter("ue_deregistration", 1, 1, "result", "success");
    increment_counter(
        "ue_deregistration", 1, 1, "action", "deregistration_accept_not_sent");
    deregistration_success_event(
        amf_ctx->_imsi64, "deregistration_accept_not_sent");
    rc = RETURNok;
  }

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_cn_ue_initiated_deregistration_ue()                       **
 **                                                                        **
 ***************************************************************************/

static int amf_cn_ue_initiated_deregistration_ue(const uint32_t ue_id) {
  int rc                   = RETURNok;
  amf_context_t* amf_ctx_p = NULL;

  OAILOG_FUNC_IN(LOG_NAS_AMF);

  // OAILOG_DEBUG(
  //    LOG_NAS_AMF,
  //    "AMF-PROC Implicit DeRegistration UE" AMF_UE_NGAP_ID_FMT "\n", ue_id);

  amf_deregistration_request_ies_t params = {0};
  // params.decode_status
  // params.guti = NULL;
  // params.imei = NULL;
  // params.imsi = NULL;
  params.is_native_sc = true;
  params.ksi          = 0;
  params.switch_off   = true;
  params.type         = AMF_UE_INITIATED_3GPP_DEREGISTRATION;

  amf_ctx_p = amf_context_get(&_amf_data, ue_id);

  if (amf_ctx_p && (params.type == AMF_UE_INITIATED_3GPP_DEREGISTRATION)) {
    rc = amf_proc_deregistration_request(ue_id, &params);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}
}  // namespace magma5g
