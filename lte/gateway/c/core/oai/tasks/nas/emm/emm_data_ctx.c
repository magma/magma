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
#include <string.h>
#include <stdlib.h>
#include <stdio.h>
#include <sys/types.h>

#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#include "lte/gateway/c/core/oai/common/assertions.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.301.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.008.h"
#include "lte/gateway/c/core/oai/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/include/mme_app_ue_context.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/NasSecurityAlgorithms.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/emm_data.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/esm_data.h"
#include "lte/gateway/c/core/oai/common/security_types.h"
#include "lte/gateway/c/core/oai/lib/secu/secu_defs.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_23.003.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_33.401.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_36.401.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/EpsBearerContextStatus.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/MobileStationClassmark2.h"
#include "lte/gateway/c/core/oai/include/TrackingAreaIdentity.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/TrackingAreaIdentityList.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/sap/emm_fsm.h"
#include "lte/gateway/c/core/oai/lib/hashtable/hashtable.h"
#include "lte/gateway/c/core/oai/tasks/nas/api/mme/mme_api.h"
#include "lte/gateway/c/core/oai/include/mme_app_state.h"
#include "lte/gateway/c/core/oai/tasks/nas/nas_procedures.h"
#include "lte/gateway/c/core/oai/tasks/nas/util/nas_timer.h"
#include "lte/gateway/c/core/oai/include/nas/securityDef.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_timer.h"

//------------------------------------------------------------------------------
mme_ue_s1ap_id_t emm_ctx_get_new_ue_id(const emm_context_t* const ctxt) {
  return (mme_ue_s1ap_id_t)((uint)((uintptr_t) ctxt) >> 4);
}

//------------------------------------------------------------------------------
inline void emm_ctx_set_attribute_present(
    emm_context_t* const ctxt, const int attribute_bit_pos) {
  ctxt->member_present_mask |= attribute_bit_pos;
}

inline void emm_ctx_clear_attribute_present(
    emm_context_t* const ctxt, const int attribute_bit_pos) {
  ctxt->member_present_mask &= ~attribute_bit_pos;
  ctxt->member_valid_mask &= ~attribute_bit_pos;
}

inline void emm_ctx_set_attribute_valid(
    emm_context_t* const ctxt, const int attribute_bit_pos) {
  ctxt->member_present_mask |= attribute_bit_pos;
  ctxt->member_valid_mask |= attribute_bit_pos;
}

inline void emm_ctx_clear_attribute_valid(
    emm_context_t* const ctxt, const int attribute_bit_pos) {
  ctxt->member_valid_mask &= ~attribute_bit_pos;
}

//------------------------------------------------------------------------------
/* Clear GUTI  */
inline void emm_ctx_clear_guti(emm_context_t* const ctxt) {
  clear_guti(&ctxt->_guti);
  emm_ctx_clear_attribute_present(ctxt, EMM_CTXT_MEMBER_GUTI);
  OAILOG_DEBUG(
      LOG_NAS_EMM, "ue_id=" MME_UE_S1AP_ID_FMT " GUTI cleared\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id);
}

/* Set GUTI */
inline void emm_ctx_set_guti(emm_context_t* const ctxt, guti_t* guti) {
  ctxt->_guti = *guti;
  emm_ctx_set_attribute_present(ctxt, EMM_CTXT_MEMBER_GUTI);
  OAILOG_DEBUG(
      LOG_NAS_EMM,
      "ue_id=" MME_UE_S1AP_ID_FMT " set GUTI " GUTI_FMT " (present)\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id,
      GUTI_ARG(&ctxt->_guti));
}

/* Set GUTI, mark it as valid */
inline void emm_ctx_set_valid_guti(emm_context_t* const ctxt, guti_t* guti) {
  ctxt->_guti = *guti;
  emm_ctx_set_attribute_valid(ctxt, EMM_CTXT_MEMBER_GUTI);
  OAILOG_DEBUG(
      LOG_NAS_EMM,
      "ue_id=" MME_UE_S1AP_ID_FMT " set GUTI " GUTI_FMT " (valid)\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id,
      GUTI_ARG(&ctxt->_guti));
}

//------------------------------------------------------------------------------
/* Clear old GUTI  */
inline void emm_ctx_clear_old_guti(emm_context_t* const ctxt) {
  clear_guti(&ctxt->_old_guti);
  emm_ctx_clear_attribute_present(ctxt, EMM_CTXT_MEMBER_OLD_GUTI);
  OAILOG_DEBUG(
      LOG_NAS_EMM, "ue_id=" MME_UE_S1AP_ID_FMT " old GUTI cleared\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id);
}

/* Set GUTI */
inline void emm_ctx_set_old_guti(emm_context_t* const ctxt, guti_t* guti) {
  ctxt->_old_guti = *guti;
  emm_ctx_set_attribute_present(ctxt, EMM_CTXT_MEMBER_OLD_GUTI);
  OAILOG_DEBUG(
      LOG_NAS_EMM,
      "ue_id=" MME_UE_S1AP_ID_FMT " set old GUTI " GUTI_FMT " (present)\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id,
      GUTI_ARG(&ctxt->_old_guti));
}

/* Set GUTI, mark it as valid */
inline void emm_ctx_set_valid_old_guti(
    emm_context_t* const ctxt, guti_t* guti) {
  ctxt->_old_guti = *guti;
  emm_ctx_set_attribute_valid(ctxt, EMM_CTXT_MEMBER_OLD_GUTI);
  OAILOG_DEBUG(
      LOG_NAS_EMM,
      "ue_id=" MME_UE_S1AP_ID_FMT " set old GUTI " GUTI_FMT " (valid)\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id,
      GUTI_ARG(&ctxt->_old_guti));
}

//------------------------------------------------------------------------------
/* Clear IMSI */
inline void emm_ctx_clear_imsi(emm_context_t* const ctxt) {
  clear_imsi(&ctxt->_imsi);
  ctxt->_imsi64 = INVALID_IMSI64;
  emm_ctx_clear_attribute_present(ctxt, EMM_CTXT_MEMBER_IMSI);
  OAILOG_DEBUG(
      LOG_NAS_EMM, "ue_id=" MME_UE_S1AP_ID_FMT " cleared IMSI\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id);
}

/* Set IMSI */
inline void emm_ctx_set_imsi(
    emm_context_t* const ctxt, imsi_t* imsi, imsi64_t imsi64) {
  ctxt->_imsi   = *imsi;
  ctxt->_imsi64 = imsi64;
  emm_ctx_set_attribute_present(ctxt, EMM_CTXT_MEMBER_IMSI);
#if DEBUG_IS_ON
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1] = {0};
  IMSI64_TO_STRING(ctxt->_imsi64, imsi_str, ctxt->_imsi.length);
  OAILOG_DEBUG(
      LOG_NAS_EMM, "ue_id=" MME_UE_S1AP_ID_FMT " set IMSI %s (valid)\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id,
      imsi_str);
#endif
}

/* Set IMSI, mark it as valid */
inline void emm_ctx_set_valid_imsi(
    emm_context_t* const ctxt, imsi_t* imsi, imsi64_t imsi64) {
  ctxt->_imsi   = *imsi;
  ctxt->_imsi64 = imsi64;
  emm_ctx_set_attribute_valid(ctxt, EMM_CTXT_MEMBER_IMSI);
#if DEBUG_IS_ON
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1] = {0};
  IMSI64_TO_STRING(ctxt->_imsi64, imsi_str, ctxt->_imsi.length);
  OAILOG_DEBUG(
      LOG_NAS_EMM, "ue_id=" MME_UE_S1AP_ID_FMT " set IMSI %s (valid)\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id,
      imsi_str);
#endif
  mme_api_notify_imsi(
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id,
      imsi64);
}

//------------------------------------------------------------------------------
/* Clear IMEI */
inline void emm_ctx_clear_imei(emm_context_t* const ctxt) {
  clear_imei(&ctxt->_imei);
  emm_ctx_clear_attribute_present(ctxt, EMM_CTXT_MEMBER_IMEI);
  OAILOG_DEBUG(
      LOG_NAS_EMM, "ue_id=" MME_UE_S1AP_ID_FMT " IMEI cleared\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id);
}

/* Set IMEI */
inline void emm_ctx_set_imei(emm_context_t* const ctxt, imei_t* imei) {
  ctxt->_imei = *imei;
  emm_ctx_set_attribute_present(ctxt, EMM_CTXT_MEMBER_IMEI);
#if DEBUG_IS_ON
  char imei_str[16];
  IMEI_TO_STRING(imei, imei_str, 16);
  OAILOG_DEBUG(
      LOG_NAS_EMM, "ue_id=" MME_UE_S1AP_ID_FMT " set IMEI %s (present)\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id,
      imei_str);
#endif
}

/* Set IMEI, mark it as valid */
inline void emm_ctx_set_valid_imei(emm_context_t* const ctxt, imei_t* imei) {
  ctxt->_imei = *imei;
  emm_ctx_set_attribute_valid(ctxt, EMM_CTXT_MEMBER_IMEI);
#if DEBUG_IS_ON
  char imei_str[16];
  IMEI_TO_STRING(imei, imei_str, 16);
  OAILOG_DEBUG(
      LOG_NAS_EMM, "ue_id=" MME_UE_S1AP_ID_FMT " set IMEI %s (valid)\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id,
      imei_str);
#endif
}

//------------------------------------------------------------------------------
/* Clear IMEI_SV */
inline void emm_ctx_clear_imeisv(emm_context_t* const ctxt) {
  clear_imeisv(&ctxt->_imeisv);
  emm_ctx_clear_attribute_present(ctxt, EMM_CTXT_MEMBER_IMEI_SV);
  OAILOG_DEBUG(
      LOG_NAS_EMM, "ue_id=" MME_UE_S1AP_ID_FMT " cleared IMEI_SV \n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id);
}

/* Set IMEI_SV */
inline void emm_ctx_set_imeisv(emm_context_t* const ctxt, imeisv_t* imeisv) {
  ctxt->_imeisv = *imeisv;
  emm_ctx_set_attribute_present(ctxt, EMM_CTXT_MEMBER_IMEI_SV);
  OAILOG_DEBUG(
      LOG_NAS_EMM, "ue_id=" MME_UE_S1AP_ID_FMT " set IMEI_SV (present)\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id);
}

/* Set IMEI_SV, mark it as valid */
inline void emm_ctx_set_valid_imeisv(
    emm_context_t* const ctxt, imeisv_t* imeisv) {
  ctxt->_imeisv = *imeisv;
  emm_ctx_set_attribute_valid(ctxt, EMM_CTXT_MEMBER_IMEI_SV);
  OAILOG_DEBUG(
      LOG_NAS_EMM, "ue_id=" MME_UE_S1AP_ID_FMT " set IMEI_SV (valid)\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id);
}

//------------------------------------------------------------------------------
/* Clear last_visited_registered_tai */
inline void emm_ctx_clear_lvr_tai(emm_context_t* const ctxt) {
  clear_tai(&ctxt->_lvr_tai);
  emm_ctx_clear_attribute_present(ctxt, EMM_CTXT_MEMBER_LVR_TAI);
  OAILOG_DEBUG(
      LOG_NAS_EMM,
      "ue_id=" MME_UE_S1AP_ID_FMT " cleared last visited registered TAI\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id);
}

/* Set last_visited_registered_tai, mark it as valid */
inline void emm_ctx_set_valid_lvr_tai(
    emm_context_t* const ctxt, tai_t* lvr_tai) {
  ctxt->_lvr_tai = *lvr_tai;
  emm_ctx_set_attribute_valid(ctxt, EMM_CTXT_MEMBER_LVR_TAI);
  // OAILOG_DEBUG (LOG_NAS_EMM, "ue_id="MME_UE_S1AP_ID_FMT" set last visited
  // registered TAI "TAI_FMT" (valid)\n", (PARENT_STRUCT(ctxt, struct
  // ue_mm_context_s, emm_context))->mme_ue_s1ap_id, TAI_ARG(&ctxt->_lvr_tai));
}

//------------------------------------------------------------------------------
/* Clear AUTH vectors  */
inline void emm_ctx_clear_auth_vectors(emm_context_t* const ctxt) {
  emm_ctx_clear_attribute_present(ctxt, EMM_CTXT_MEMBER_AUTH_VECTORS);
  for (int i = 0; i < MAX_EPS_AUTH_VECTORS; i++) {
    memset((void*) &ctxt->_vector[i], 0, sizeof(ctxt->_vector[i]));
    emm_ctx_clear_attribute_present(ctxt, EMM_CTXT_MEMBER_AUTH_VECTOR0 + i);
  }
  emm_ctx_clear_security_vector_index(ctxt);
  OAILOG_DEBUG(
      LOG_NAS_EMM, "ue_id=" MME_UE_S1AP_ID_FMT " cleared auth vectors \n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id);
}
//------------------------------------------------------------------------------
/* Clear AUTH vector  */
inline void emm_ctx_clear_auth_vector(emm_context_t* const ctxt, ksi_t eksi) {
  AssertFatal(eksi < MAX_EPS_AUTH_VECTORS, "Out of bounds eksi %d", eksi);
  memset(
      (void*) &ctxt->_vector[eksi % MAX_EPS_AUTH_VECTORS], 0,
      sizeof(ctxt->_vector[eksi % MAX_EPS_AUTH_VECTORS]));
  emm_ctx_clear_attribute_present(ctxt, EMM_CTXT_MEMBER_AUTH_VECTOR0 + eksi);
  int remaining_vectors = 0;
  for (int i = 0; i < MAX_EPS_AUTH_VECTORS; i++) {
    if (IS_EMM_CTXT_VALID_AUTH_VECTOR(ctxt, i)) {
      remaining_vectors += 1;
    }
  }
  ctxt->remaining_vectors = remaining_vectors;
  OAILOG_DEBUG(
      LOG_NAS_EMM, "ue_id=" MME_UE_S1AP_ID_FMT " cleared auth vector %u \n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id,
      eksi);
  if (!(remaining_vectors)) {
    emm_ctx_clear_attribute_present(ctxt, EMM_CTXT_MEMBER_AUTH_VECTORS);
    emm_ctx_clear_security_vector_index(ctxt);
    OAILOG_DEBUG(
        LOG_NAS_EMM, "ue_id=" MME_UE_S1AP_ID_FMT " cleared auth vectors\n",
        (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
            ->mme_ue_s1ap_id);
  }
}
//------------------------------------------------------------------------------
/* Clear security  */
inline void emm_ctx_clear_security(emm_context_t* const ctxt) {
  memset(&ctxt->_security, 0, sizeof(ctxt->_security));
  emm_ctx_set_security_type(ctxt, SECURITY_CTX_TYPE_NOT_AVAILABLE);
  emm_ctx_set_security_eksi(ctxt, KSI_NO_KEY_AVAILABLE);
  ctxt->_security.selected_algorithms.encryption = NAS_SECURITY_ALGORITHMS_EEA0;
  ctxt->_security.selected_algorithms.integrity  = NAS_SECURITY_ALGORITHMS_EIA0;
  emm_ctx_clear_attribute_present(ctxt, EMM_CTXT_MEMBER_SECURITY);
  ctxt->_security.direction_decode = SECU_DIRECTION_UPLINK;
  ctxt->_security.direction_encode = SECU_DIRECTION_DOWNLINK;
  OAILOG_DEBUG(
      LOG_NAS_EMM, "ue_id=" MME_UE_S1AP_ID_FMT " cleared security context \n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id);
}

//------------------------------------------------------------------------------
inline void emm_ctx_set_security_type(
    emm_context_t* const ctxt, emm_sc_type_t sc_type) {
  ctxt->_security.sc_type = sc_type;
  OAILOG_TRACE(
      LOG_NAS_EMM,
      "ue_id=" MME_UE_S1AP_ID_FMT " set security context security type %d\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id,
      sc_type);
}

//------------------------------------------------------------------------------
inline void emm_ctx_set_security_eksi(emm_context_t* const ctxt, ksi_t eksi) {
  ctxt->_security.eksi = eksi;
  OAILOG_TRACE(
      LOG_NAS_EMM,
      "ue_id=" MME_UE_S1AP_ID_FMT " set security context eksi %d\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id,
      eksi);
}

//------------------------------------------------------------------------------
inline void emm_ctx_clear_security_vector_index(emm_context_t* const ctxt) {
  ctxt->_security.vector_index = EMM_SECURITY_VECTOR_INDEX_INVALID;
  OAILOG_TRACE(
      LOG_NAS_EMM,
      "ue_id=" MME_UE_S1AP_ID_FMT " clear security context vector index\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id);
}
//------------------------------------------------------------------------------
inline void emm_ctx_set_security_vector_index(
    emm_context_t* const ctxt, int vector_index) {
  ctxt->_security.vector_index = vector_index;
  OAILOG_TRACE(
      LOG_NAS_EMM,
      "ue_id=" MME_UE_S1AP_ID_FMT " set security context vector index %d\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id,
      vector_index);
}

//------------------------------------------------------------------------------
/* Clear non current security  */
inline void emm_ctx_clear_non_current_security(emm_context_t* const ctxt) {
  memset(&ctxt->_non_current_security, 0, sizeof(ctxt->_non_current_security));
  ctxt->_non_current_security.sc_type = SECURITY_CTX_TYPE_NOT_AVAILABLE;
  ctxt->_non_current_security.eksi    = KSI_NO_KEY_AVAILABLE;
  ctxt->_non_current_security.selected_algorithms.encryption =
      NAS_SECURITY_ALGORITHMS_EEA0;
  ctxt->_non_current_security.selected_algorithms.integrity =
      NAS_SECURITY_ALGORITHMS_EIA0;
  emm_ctx_clear_attribute_present(ctxt, EMM_CTXT_MEMBER_NON_CURRENT_SECURITY);
  ctxt->_security.direction_decode = SECU_DIRECTION_UPLINK;
  ctxt->_security.direction_encode = SECU_DIRECTION_DOWNLINK;
  OAILOG_DEBUG(
      LOG_NAS_EMM,
      "ue_id=" MME_UE_S1AP_ID_FMT " cleared non current security context \n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id);
}

//------------------------------------------------------------------------------
/* Clear UE network capability IE   */
inline void emm_ctx_clear_ue_nw_cap(emm_context_t* const ctxt) {
  memset(
      &ctxt->_ue_network_capability, 0, sizeof(ctxt->_ue_network_capability));
  emm_ctx_clear_attribute_present(
      ctxt, EMM_CTXT_MEMBER_UE_NETWORK_CAPABILITY_IE);
  OAILOG_DEBUG(
      LOG_NAS_EMM,
      "ue_id=" MME_UE_S1AP_ID_FMT " cleared UE network capability IE\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id);
}

/* Set UE network capability IE */
inline void emm_ctx_set_ue_nw_cap(
    emm_context_t* const ctxt,
    const ue_network_capability_t* const ue_nw_cap_ie) {
  ctxt->_ue_network_capability = *ue_nw_cap_ie;
  emm_ctx_set_attribute_present(ctxt, EMM_CTXT_MEMBER_UE_NETWORK_CAPABILITY_IE);
  OAILOG_DEBUG(
      LOG_NAS_EMM,
      "ue_id=" MME_UE_S1AP_ID_FMT " set UE network capability IE (present)\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id);
}

/* Set UE network capability IE, mark it as valid */
inline void emm_ctx_set_valid_ue_nw_cap(
    emm_context_t* const ctxt,
    const ue_network_capability_t* const ue_nw_cap_ie) {
  ctxt->_ue_network_capability = *ue_nw_cap_ie;
  emm_ctx_set_attribute_valid(ctxt, EMM_CTXT_MEMBER_UE_NETWORK_CAPABILITY_IE);
  OAILOG_DEBUG(
      LOG_NAS_EMM,
      "ue_id=" MME_UE_S1AP_ID_FMT " set UE network capability IE (valid)\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id);
}

//------------------------------------------------------------------------------
/* Clear MS network capability IE   */
inline void emm_ctx_clear_ms_nw_cap(emm_context_t* const ctxt) {
  memset(
      &ctxt->_ms_network_capability, 0, sizeof(ctxt->_ms_network_capability));
  emm_ctx_clear_attribute_present(
      ctxt, EMM_CTXT_MEMBER_MS_NETWORK_CAPABILITY_IE);
  OAILOG_DEBUG(
      LOG_NAS_EMM,
      "ue_id=" MME_UE_S1AP_ID_FMT " cleared MS network capability IE\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id);
}

/* Set UE network capability IE */
inline void emm_ctx_set_ms_nw_cap(
    emm_context_t* const ctxt,
    const ms_network_capability_t* const ms_nw_cap_ie) {
  ctxt->_ms_network_capability = *ms_nw_cap_ie;
  emm_ctx_set_attribute_present(ctxt, EMM_CTXT_MEMBER_MS_NETWORK_CAPABILITY_IE);
  OAILOG_DEBUG(
      LOG_NAS_EMM,
      "ue_id=" MME_UE_S1AP_ID_FMT " set MS network capability IE (present)\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id);
}

/* Set UE network capability IE, mark it as valid */
inline void emm_ctx_set_valid_ms_nw_cap(
    emm_context_t* const ctxt,
    const ms_network_capability_t* const ms_nw_cap_ie) {
  ctxt->_ms_network_capability = *ms_nw_cap_ie;
  emm_ctx_set_attribute_valid(ctxt, EMM_CTXT_MEMBER_MS_NETWORK_CAPABILITY_IE);
  OAILOG_DEBUG(
      LOG_NAS_EMM,
      "ue_id=" MME_UE_S1AP_ID_FMT " set MS network capability IE (valid)\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id);
}

//------------------------------------------------------------------------------
/* Clear current DRX parameter   */
inline void emm_ctx_clear_drx_parameter(emm_context_t* const ctxt) {
  memset(&ctxt->_drx_parameter, 0, sizeof(drx_parameter_t));
  emm_ctx_clear_attribute_present(ctxt, EMM_CTXT_MEMBER_CURRENT_DRX_PARAMETER);
  OAILOG_DEBUG(
      LOG_NAS_EMM,
      "ue_id=" MME_UE_S1AP_ID_FMT " cleared current DRX parameter\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id);
}

/* Set current DRX parameter */
inline void emm_ctx_set_drx_parameter(
    emm_context_t* const ctxt, drx_parameter_t* drx) {
  memcpy(&ctxt->_drx_parameter, drx, sizeof(drx_parameter_t));
  emm_ctx_set_attribute_present(ctxt, EMM_CTXT_MEMBER_CURRENT_DRX_PARAMETER);
  OAILOG_DEBUG(
      LOG_NAS_EMM,
      "ue_id=" MME_UE_S1AP_ID_FMT " set current DRX parameter (present)\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id);
}

/* Set current DRX parameter, mark it as valid */
inline void emm_ctx_set_valid_drx_parameter(
    emm_context_t* const ctxt, drx_parameter_t* drx) {
  emm_ctx_set_drx_parameter(ctxt, drx);
  emm_ctx_set_attribute_valid(ctxt, EMM_CTXT_MEMBER_CURRENT_DRX_PARAMETER);
  OAILOG_DEBUG(
      LOG_NAS_EMM,
      "ue_id=" MME_UE_S1AP_ID_FMT " set current DRX parameter (valid)\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id);
}

//------------------------------------------------------------------------------
/* Clear UE additional security capability */
inline void emm_ctx_clear_ue_additional_security_capability(
    emm_context_t* const ctxt) {
  memset(
      &ctxt->ue_additional_security_capability, 0,
      sizeof(ue_additional_security_capability_t));
  emm_ctx_clear_attribute_present(
      ctxt, EMM_CTXT_MEMBER_UE_ADDITIONAL_SECURITY_CAPABILITY);
  OAILOG_DEBUG(
      LOG_NAS_EMM,
      "ue_id=" MME_UE_S1AP_ID_FMT
      " cleared ue additional security capability\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id);
}

/* Set UE additional security capability */
inline void emm_ctx_set_ue_additional_security_capability(
    emm_context_t* const ctxt, ue_additional_security_capability_t* uasc) {
  memcpy(
      &ctxt->ue_additional_security_capability, uasc,
      sizeof(ue_additional_security_capability_t));
  emm_ctx_set_attribute_present(
      ctxt, EMM_CTXT_MEMBER_UE_ADDITIONAL_SECURITY_CAPABILITY);
  OAILOG_DEBUG(
      LOG_NAS_EMM,
      "ue_id=" MME_UE_S1AP_ID_FMT
      " set ue additional security capability (present)\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id);
}

//------------------------------------------------------------------------------
/* Clear EPS bearer context status   */
inline void emm_ctx_clear_eps_bearer_context_status(emm_context_t* const ctxt) {
  memset(
      &ctxt->_eps_bearer_context_status, 0,
      sizeof(ctxt->_eps_bearer_context_status));
  emm_ctx_clear_attribute_present(
      ctxt, EMM_CTXT_MEMBER_EPS_BEARER_CONTEXT_STATUS);
  OAILOG_DEBUG(
      LOG_NAS_EMM,
      "ue_id=" MME_UE_S1AP_ID_FMT " cleared EPS bearer context status\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id);
}

/* Set current DRX parameter */
inline void emm_ctx_set_eps_bearer_context_status(
    emm_context_t* const ctxt, eps_bearer_context_status_t* status) {
  ctxt->_eps_bearer_context_status = *status;
  emm_ctx_set_attribute_present(
      ctxt, EMM_CTXT_MEMBER_EPS_BEARER_CONTEXT_STATUS);
  OAILOG_DEBUG(
      LOG_NAS_EMM,
      "ue_id=" MME_UE_S1AP_ID_FMT " set EPS bearer context status (present)\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id);
}

/* Set current DRX parameter, mark it as valid */
inline void emm_ctx_set_valid_eps_bearer_context_status(
    emm_context_t* const ctxt, eps_bearer_context_status_t* status) {
  ctxt->_eps_bearer_context_status = *status;
  emm_ctx_set_attribute_valid(ctxt, EMM_CTXT_MEMBER_EPS_BEARER_CONTEXT_STATUS);
  OAILOG_DEBUG(
      LOG_NAS_EMM,
      "ue_id=" MME_UE_S1AP_ID_FMT " set EPS bearer context status (valid)\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id);
}

//------------------------------------------------------------------------------
/* Clear mobile station class mark2 */
inline void emm_ctx_clear_mobile_station_clsMark2(emm_context_t* const ctxt) {
  memset(&ctxt->_mob_st_clsMark2, 0, sizeof(ctxt->_mob_st_clsMark2));
  emm_ctx_clear_attribute_present(ctxt, EMM_CTXT_MEMBER_MOB_STATION_CLSMARK2);
  OAILOG_DEBUG(
      LOG_NAS_EMM,
      "ue_id=" MME_UE_S1AP_ID_FMT " cleared mobile station classmark2\n",
      (PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context))
          ->mme_ue_s1ap_id);
}

/* Set mob_station_clsMark2 */
inline void emm_ctx_set_mobile_station_clsMark2(
    emm_context_t* const ctxt, MobileStationClassmark2* mob_st_clsMark2) {
  ctxt->_mob_st_clsMark2 = *mob_st_clsMark2;
  emm_ctx_set_attribute_present(ctxt, EMM_CTXT_MEMBER_MOB_STATION_CLSMARK2);
}

//------------------------------------------------------------------------------
/* Free dynamically allocated memory */
void free_emm_ctx_memory(
    emm_context_t* const ctxt, const mme_ue_s1ap_id_t ue_id) {
  OAILOG_DEBUG(
      LOG_NAS_EMM, "Freeing up emm_context for ue_id=" MME_UE_S1AP_ID_FMT,
      ue_id);
  if (!ctxt) {
    return;
  }
  if (ctxt->t3422_arg) {
    free_wrapper((void**) &ctxt->t3422_arg);
  }
  nas_delete_all_emm_procedures(ctxt);
  free_esm_context_content(&ctxt->esm_ctx);
  bdestroy_wrapper(&ctxt->esm_msg);
}

//------------------------------------------------------------------------------
struct emm_context_s* emm_context_get(
    emm_data_t* emm_data,  // TODO REMOVE
    const mme_ue_s1ap_id_t ue_id) {
  struct emm_context_s* emm_context_p = NULL;

  DevAssert(emm_data);
  if (INVALID_MME_UE_S1AP_ID != ue_id) {
    ue_mm_context_t* ue_mm_context =
        mme_ue_context_exists_mme_ue_s1ap_id(ue_id);
    if (ue_mm_context) {
      emm_context_p = &ue_mm_context->emm_context;
    }
    OAILOG_DEBUG(
        LOG_NAS_EMM, "EMM-CTX - get UE id " MME_UE_S1AP_ID_FMT " context %p\n",
        ue_id, emm_context_p);
  }
  return emm_context_p;
}

//------------------------------------------------------------------------------
struct emm_context_s* emm_context_get_by_imsi(
    emm_data_t* emm_data,  // TODO REMOVE
    imsi64_t imsi64) {
  struct emm_context_s* emm_context_p = NULL;

  mme_app_desc_t* mme_app_desc_p = get_mme_nas_state(false);
  ue_mm_context_t* ue_mm_context =
      mme_ue_context_exists_imsi(&mme_app_desc_p->mme_ue_contexts, imsi64);
  if (ue_mm_context) {
    emm_context_p = &ue_mm_context->emm_context;
  }

#if DEBUG_IS_ON
  if (emm_context_p) {
    OAILOG_DEBUG(
        LOG_NAS_EMM,
        "EMM-CTX - get UE id " MME_UE_S1AP_ID_FMT
        " context %p by imsi " IMSI_64_FMT "\n",
        ue_mm_context->mme_ue_s1ap_id, emm_context_p, imsi64);
  }
#endif
  return emm_context_p;
}

//------------------------------------------------------------------------------
struct emm_context_s* emm_context_get_by_guti(
    emm_data_t* emm_data, guti_t* guti) {
  struct emm_context_s* emm_context_p = NULL;

  mme_app_desc_t* mme_app_desc_p = get_mme_nas_state(false);
  ue_mm_context_t* ue_mm_context =
      mme_ue_context_exists_guti(&mme_app_desc_p->mme_ue_contexts, guti);
  if (ue_mm_context) {
    emm_context_p = &ue_mm_context->emm_context;
  }
#if DEBUG_IS_ON
  if (emm_context_p) {
    OAILOG_DEBUG(
        LOG_NAS_EMM,
        "EMM-CTX - get UE id " MME_UE_S1AP_ID_FMT
        " context %p by guti " GUTI_FMT "\n",
        ue_mm_context->mme_ue_s1ap_id, emm_context_p, GUTI_ARG(guti));
  }
#endif
  return emm_context_p;
}

//------------------------------------------------------------------------------

status_code_e emm_context_upsert_imsi(
    emm_data_t* emm_data, struct emm_context_s* elm) {
  hashtable_rc_t h_rc = HASH_TABLE_OK;
  mme_ue_s1ap_id_t ue_id =
      (PARENT_STRUCT(elm, struct ue_mm_context_s, emm_context))->mme_ue_s1ap_id;

  mme_app_desc_t* mme_app_desc_p = get_mme_nas_state(false);
  h_rc                           = hashtable_uint64_ts_remove(
      mme_app_desc_p->mme_ue_contexts.imsi_mme_ue_id_htbl,
      (const hash_key_t) elm->_imsi64);
  if (INVALID_MME_UE_S1AP_ID != ue_id) {
    h_rc = hashtable_uint64_ts_insert(
        mme_app_desc_p->mme_ue_contexts.imsi_mme_ue_id_htbl,
        (const hash_key_t) elm->_imsi64, ue_id);
  } else {
    h_rc = HASH_TABLE_KEY_NOT_EXISTS;
  }
  if (HASH_TABLE_OK != h_rc) {
    OAILOG_TRACE(
        LOG_MME_APP,
        "Error could not update this ue context "
        "mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT " imsi " IMSI_64_FMT ": %s\n",
        ue_id, elm->_imsi64, hashtable_rc_code2string(h_rc));
    return RETURNerror;
  }
  return RETURNok;
}

//------------------------------------------------------------------------------
void emm_init_context(
    struct emm_context_s* const emm_ctx, const bool init_esm_ctxt) {
  emm_ctx->_emm_fsm_state = EMM_DEREGISTERED;

  struct ue_mm_context_s* ue_mm_context =
      PARENT_STRUCT(emm_ctx, struct ue_mm_context_s, emm_context);
  OAILOG_DEBUG(
      LOG_NAS_EMM, "UE " MME_UE_S1AP_ID_FMT " Init EMM-CTX\n",
      ue_mm_context->mme_ue_s1ap_id);

  emm_ctx_clear_guti(emm_ctx);
  emm_ctx_clear_old_guti(emm_ctx);
  emm_ctx_clear_imsi(emm_ctx);
  emm_ctx_clear_imei(emm_ctx);
  emm_ctx_clear_imeisv(emm_ctx);
  emm_ctx_clear_lvr_tai(emm_ctx);
  emm_ctx_clear_security(emm_ctx);
  emm_ctx_clear_non_current_security(emm_ctx);
  emm_ctx_clear_auth_vectors(emm_ctx);
  emm_ctx_clear_ms_nw_cap(emm_ctx);
  emm_ctx_clear_ue_nw_cap(emm_ctx);
  emm_ctx_clear_drx_parameter(emm_ctx);
  emm_ctx_clear_mobile_station_clsMark2(emm_ctx);
  emm_ctx_clear_ue_additional_security_capability(emm_ctx);
  emm_ctx->T3422.id          = NAS_TIMER_INACTIVE_ID;
  emm_ctx->T3422.msec        = mme_config.nas_config.t3422_msec;
  emm_ctx->new_attach_info   = NULL;
  emm_ctx->emm_context_state = NEW_EMM_CONTEXT_NOT_CREATED;

  if (init_esm_ctxt) {
    esm_init_context(&emm_ctx->esm_ctx);
  }
  emm_ctx->emm_procedures = NULL;
  emm_ctx->esm_msg        = NULL;
}

//------------------------------------------------------------------------------
void nas_start_T3450(
    const mme_ue_s1ap_id_t ue_id, struct nas_timer_s* const T3450,
    time_out_t time_out_cb) {
  if ((T3450) && (T3450->id == NAS_TIMER_INACTIVE_ID)) {
    T3450->id =
        mme_app_start_timer(T3450->msec, TIMER_REPEAT_ONCE, time_out_cb, ue_id);
    if (NAS_TIMER_INACTIVE_ID != T3450->id) {
      OAILOG_DEBUG(
          LOG_NAS_EMM, "T3450 started UE " MME_UE_S1AP_ID_FMT "\n", ue_id);
    } else {
      OAILOG_ERROR(
          LOG_NAS_EMM, "Could not start T3450 UE " MME_UE_S1AP_ID_FMT " ",
          ue_id);
    }
  }
}
//------------------------------------------------------------------------------
void nas_start_T3460(
    const mme_ue_s1ap_id_t ue_id, struct nas_timer_s* const T3460,
    time_out_t time_out_cb) {
  if ((T3460) && (T3460->id == NAS_TIMER_INACTIVE_ID)) {
    T3460->id =
        mme_app_start_timer(T3460->msec, TIMER_REPEAT_ONCE, time_out_cb, ue_id);
    if (NAS_TIMER_INACTIVE_ID != T3460->id) {
      OAILOG_DEBUG(
          LOG_NAS_EMM, "T3460 started UE " MME_UE_S1AP_ID_FMT "\n", ue_id);
    } else {
      OAILOG_ERROR(
          LOG_NAS_EMM, "Could not start T3460 UE " MME_UE_S1AP_ID_FMT " ",
          ue_id);
    }
  }
}
//------------------------------------------------------------------------------
void nas_start_T3470(
    const mme_ue_s1ap_id_t ue_id, struct nas_timer_s* const T3470,
    time_out_t time_out_cb) {
  if ((T3470) && (T3470->id == NAS_TIMER_INACTIVE_ID)) {
    T3470->id =
        mme_app_start_timer(T3470->msec, TIMER_REPEAT_ONCE, time_out_cb, ue_id);
    if (NAS_TIMER_INACTIVE_ID != T3470->id) {
      OAILOG_DEBUG(
          LOG_NAS_EMM, "T3470 started UE " MME_UE_S1AP_ID_FMT "\n", ue_id);
    } else {
      OAILOG_ERROR(
          LOG_NAS_EMM, "Could not start T3470 UE " MME_UE_S1AP_ID_FMT " ",
          ue_id);
    }
  }
}
//------------------------------------------------------------------------------
void nas_start_T3422(
    const mme_ue_s1ap_id_t ue_id, struct nas_timer_s* const T3422,
    time_out_t time_out_cb) {
  if ((T3422) && (T3422->id == NAS_TIMER_INACTIVE_ID)) {
    T3422->id =
        mme_app_start_timer(T3422->msec, TIMER_REPEAT_ONCE, time_out_cb, ue_id);
    if (NAS_TIMER_INACTIVE_ID != T3422->id) {
      OAILOG_DEBUG(
          LOG_NAS_EMM, "T3422 started UE " MME_UE_S1AP_ID_FMT "\n", ue_id);
    } else {
      OAILOG_ERROR(
          LOG_NAS_EMM, "Could not start T3422 UE " MME_UE_S1AP_ID_FMT " ",
          ue_id);
    }
  }
}
//------------------------------------------------------------------------------
void nas_start_Ts6a_auth_info(
    const mme_ue_s1ap_id_t ue_id, struct nas_timer_s* const Ts6a_auth_info,
    time_out_t time_out_cb) {
  if ((Ts6a_auth_info) && (Ts6a_auth_info->id == NAS_TIMER_INACTIVE_ID)) {
    Ts6a_auth_info->id = mme_app_start_timer(
        Ts6a_auth_info->msec, TIMER_REPEAT_ONCE, time_out_cb, ue_id);
    if (NAS_TIMER_INACTIVE_ID != Ts6a_auth_info->id) {
      OAILOG_DEBUG(
          LOG_NAS_EMM, "Ts6a_auth_info started UE " MME_UE_S1AP_ID_FMT "\n",
          ue_id);
    } else {
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "Could not start Ts6a_auth_info UE " MME_UE_S1AP_ID_FMT " ", ue_id);
    }
  }
}
//------------------------------------------------------------------------------
void nas_stop_T3450(
    const mme_ue_s1ap_id_t ue_id, struct nas_timer_s* const T3450) {
  if ((T3450) && (T3450->id != NAS_TIMER_INACTIVE_ID)) {
    mme_app_stop_timer(T3450->id);
    T3450->id = NAS_TIMER_INACTIVE_ID;
    OAILOG_DEBUG(
        LOG_NAS_EMM, "T3450 stopped UE " MME_UE_S1AP_ID_FMT "\n", ue_id);
  }
}

//------------------------------------------------------------------------------
void nas_stop_T3460(
    const mme_ue_s1ap_id_t ue_id, struct nas_timer_s* const T3460) {
  if ((T3460) && (T3460->id != NAS_TIMER_INACTIVE_ID)) {
    mme_app_stop_timer(T3460->id);
    T3460->id = NAS_TIMER_INACTIVE_ID;
    OAILOG_DEBUG(
        LOG_NAS_EMM, "T3460 stopped UE " MME_UE_S1AP_ID_FMT "\n", ue_id);
  }
}

//------------------------------------------------------------------------------
void nas_stop_T3470(
    const mme_ue_s1ap_id_t ue_id, struct nas_timer_s* const T3470) {
  if ((T3470) && (T3470->id != NAS_TIMER_INACTIVE_ID)) {
    mme_app_stop_timer(T3470->id);
    T3470->id = NAS_TIMER_INACTIVE_ID;
    OAILOG_DEBUG(
        LOG_NAS_EMM, "T3470 stopped UE " MME_UE_S1AP_ID_FMT "\n", ue_id);
  }
}

//------------------------------------------------------------------------------
void nas_stop_T3422(const imsi64_t imsi64, struct nas_timer_s* const T3422) {
  if ((T3422) && (T3422->id != NAS_TIMER_INACTIVE_ID)) {
    mme_app_stop_timer(T3422->id);
    T3422->id = NAS_TIMER_INACTIVE_ID;
    OAILOG_DEBUG_UE(LOG_NAS_EMM, imsi64, "T3422 stopped ");
  }
}
