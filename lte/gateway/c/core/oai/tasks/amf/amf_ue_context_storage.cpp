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

#include "lte/gateway/c/core/oai/tasks/amf/include/amf_ue_context_storage.h"

namespace magma5g {

/****************************************************************************
 **                                                                        **
 ** Name:    amf_create_new_ue_context()                                   **
 **                                                                        **
 ** Description: Creates new UE context                                    **
 **                                                                        **
 ***************************************************************************/
std::shared_ptr<ue_m5gmm_context_t>
AmfUeContextStorage::amf_create_new_ue_context(void) {
  std::shared_ptr<ue_m5gmm_context_t> new_p =
      std::make_shared<ue_m5gmm_context_t>();

  if (!new_p) {
    OAILOG_ERROR(LOG_AMF_APP, "Failed to allocate memory for UE context \n");
    return NULL;
  }

  new_p->amf_ue_ngap_id  = generate_amf_ue_ngap_id();
  new_p->gnb_ngap_id_key = INVALID_GNB_UE_NGAP_ID_KEY;
  new_p->gnb_ue_ngap_id  = INVALID_GNB_UE_NGAP_ID;

  // Initialize timers to INVALID IDs
  new_p->m5_mobile_reachability_timer.id    = AMF_APP_TIMER_INACTIVE_ID;
  new_p->m5_implicit_detach_timer.id        = AMF_APP_TIMER_INACTIVE_ID;
  new_p->m5_initial_context_setup_rsp_timer = (amf_app_timer_t){
      AMF_APP_TIMER_INACTIVE_ID, AMF_APP_INITIAL_CONTEXT_SETUP_RSP_TIMER_VALUE};
  new_p->m5_ulr_response_timer = (amf_app_timer_t){
      AMF_APP_TIMER_INACTIVE_ID, AMF_APP_ULR_RESPONSE_TIMER_VALUE};
  new_p->m5_ue_context_modification_timer = (amf_app_timer_t){
      AMF_APP_TIMER_INACTIVE_ID, AMF_APP_UE_CONTEXT_MODIFICATION_TIMER_VALUE};
  new_p->mm_state = DEREGISTERED;

  new_p->amf_context._security.eksi = KSI_NO_KEY_AVAILABLE;
  new_p->mm_state                   = DEREGISTERED;
  new_p->amf_context.ue_context     = new_p;

  return new_p;
}

// id-AMF-UE-NGAP-ID <--> ue_m5gmm_context_t map
bool AmfUeContextStorage::amf_insert_into_amfid_ue_context_map(
    std::shared_ptr<ue_m5gmm_context_t> pContext) {
  if (!pContext && (INVALID_AMF_UE_NGAP_ID == pContext->amf_ue_ngap_id)) {
    return false;
  }
  amfid_ue_context_map.insert_or_update(pContext->amf_ue_ngap_id, pContext);
  return true;
}
bool AmfUeContextStorage::amf_remove_from_amfid_ue_context_map(
    amf_ue_ngap_id_t ue_amf_id) {
  magma::map_rc_t status = amfid_ue_context_map.remove(ue_amf_id);

  if (magma::MAP_OK != status) {
    return false;
  }
  return true;
}

std::shared_ptr<ue_m5gmm_context_t>
AmfUeContextStorage::amf_get_from_amfid_ue_context_map(
    amf_ue_ngap_id_t ue_amf_id) {
  std::shared_ptr<ue_m5gmm_context_t> ctx;
  magma::map_rc_t status = amfid_ue_context_map.get(ue_amf_id, ctx);
  return ctx;
}

// id-GNB-UE-NGAP-ID <---> ue_m5gmm_context_t map
bool AmfUeContextStorage::amf_insert_into_gnbid_ue_context_map(
    std::shared_ptr<ue_m5gmm_context_t> pContext) {
  if (!pContext && (INVALID_GNB_UE_NGAP_ID_KEY == pContext->gnb_ngap_id_key)) {
    return false;
  }
  gnbid_ue_context_map.insert_or_update(pContext->gnb_ngap_id_key, pContext);
  return true;
}
bool AmfUeContextStorage::amf_remove_from_gnbid_ue_context_map(
    gnb_ue_ngap_id_t ue_gnb_id) {
  magma::map_rc_t status = gnbid_ue_context_map.remove(ue_gnb_id);

  if (magma::MAP_OK != status) {
    return false;
  }
  return true;
}
std::shared_ptr<ue_m5gmm_context_t>
AmfUeContextStorage::amf_get_from_gnbid_ue_context_map(
    gnb_ue_ngap_id_t ue_gnb_id) {
  std::shared_ptr<ue_m5gmm_context_t> ctx;
  magma::map_rc_t status = gnbid_ue_context_map.get(ue_gnb_id, ctx);
  return ctx;
}

// GUTI <---> ue_m5gmm_context_t map
bool AmfUeContextStorage::amf_insert_into_guti_ue_context_map(
    std::shared_ptr<ue_m5gmm_context_t> pContext) {
  if (!pContext) {
    return false;
  }
  guti_ue_context_map.insert_or_update(pContext->amf_context.m5_guti, pContext);
  return true;
}
bool AmfUeContextStorage::amf_remove_from_guti_ue_context_map(guti_m5_t guti) {
  magma::map_rc_t status = guti_ue_context_map.remove(guti);
  if (magma::MAP_OK != status) {
    return false;
  }
  return true;
}
std::shared_ptr<ue_m5gmm_context_t>
AmfUeContextStorage::amf_get_from_guti_ue_context_map(guti_m5_t guti) {
  std::shared_ptr<ue_m5gmm_context_t> ctx;
  magma::map_rc_t status = guti_ue_context_map.get(guti, ctx);
  return ctx;
}

// SUPI  <---> ue_m5gmm_context_t map
bool AmfUeContextStorage::amf_insert_into_supi_ue_context_map(
    std::shared_ptr<ue_m5gmm_context_t> pContext) {
  if (!pContext) {
    return false;
  }
  supi_ue_context_map.insert_or_update(pContext->amf_context.imsi64, pContext);
  return true;
}
bool AmfUeContextStorage::amf_remove_from_supi_ue_context_map(imsi64_t supi) {
  magma::map_rc_t status = supi_ue_context_map.remove(supi);
  if (magma::MAP_OK != status) {
    return false;
  }
  return true;
}
std::shared_ptr<ue_m5gmm_context_t>
AmfUeContextStorage::amf_get_from_supi_ue_context_map(imsi64_t supi) {
  std::shared_ptr<ue_m5gmm_context_t> ctx;
  magma::map_rc_t status = supi_ue_context_map.get(supi, ctx);
  return ctx;
}

bool AmfUeContextStorage::amf_remove_ue_context_from_cache(
    amf_ue_ngap_id_t ue_amf_id) {
  std::shared_ptr<ue_m5gmm_context_t> ctx =
      amf_get_from_amfid_ue_context_map(ue_amf_id);
  if (ctx) {
    amf_remove_from_amfid_ue_context_map(ue_amf_id);
    amf_remove_from_gnbid_ue_context_map(ctx->gnb_ngap_id_key);
    amf_remove_from_guti_ue_context_map(ctx->amf_context.m5_guti);
    amf_remove_from_supi_ue_context_map(ctx->amf_context.imsi64);
    return true;
  }
  return false;
}

bool AmfUeContextStorage::amf_add_ue_context_in_cache(
    std::shared_ptr<ue_m5gmm_context_t> pContext) {
    if(!pContext) {
        return false;
    }
    amf_insert_into_amfid_ue_context_map(pContext->amf_ue_ngap_id);
    amf_insert_into_gnbid_ue_context_map(pContext->gnb_ngap_id_key);
    amf_insert_into_guti_ue_context_map(pContext->amf_context.m5_guti);
    amf_insert_into_supi_ue_context_map(pContext->amf_context.imsi64);
    return true;
}

}  // namespace magma5g
