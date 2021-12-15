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

#include "lte/gateway/c/core/oai/tasks/amf/amf_app_state_converter.h"
#include <vector>
#include <memory>
extern "C" {
#include "lte/gateway/c/core/oai/lib/message_utils/bytes_to_ie.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/lib/message_utils/ie_to_bytes.h"
#include "lte/gateway/c/core/oai/common/log.h"
}

using magma::lte::oai::MmeNasState;
namespace magma5g {

AmfNasStateConverter::AmfNasStateConverter()  = default;
AmfNasStateConverter::~AmfNasStateConverter() = default;

// HelperFunction: Converts guti_m5_t to std::string
std::string AmfNasStateConverter::amf_app_convert_guti_m5_to_string(
    const guti_m5_t& guti) {
#define GUTI_STRING_LEN 25
  char* temp_str =
      reinterpret_cast<char*>(calloc(1, sizeof(char) * GUTI_STRING_LEN));
  snprintf(
      temp_str, GUTI_STRING_LEN, "%x%x%x%x%x%x%02x%04x%04x%08x",
      guti.guamfi.plmn.mcc_digit1, guti.guamfi.plmn.mcc_digit2,
      guti.guamfi.plmn.mcc_digit3, guti.guamfi.plmn.mnc_digit1,
      guti.guamfi.plmn.mnc_digit2, guti.guamfi.plmn.mnc_digit3,
      guti.guamfi.amf_regionid, guti.guamfi.amf_set_id, guti.guamfi.amf_pointer,
      guti.m_tmsi);
  std::string guti_str(temp_str);
  free(temp_str);
  return guti_str;
}

// HelperFunction: Converts std:: string back to guti_m5_t
void AmfNasStateConverter::amf_app_convert_string_to_guti_m5(
    const std::string& guti_str, guti_m5_t* guti_m5_p) {
  int idx                   = 0;
  std::size_t chars_to_read = 1;
#define HEX_BASE_VAL 16
  guti_m5_p->guamfi.plmn.mcc_digit1 = std::stoul(
      guti_str.substr(idx++, chars_to_read), &chars_to_read, HEX_BASE_VAL);
  guti_m5_p->guamfi.plmn.mcc_digit2 = std::stoul(
      guti_str.substr(idx++, chars_to_read), &chars_to_read, HEX_BASE_VAL);
  guti_m5_p->guamfi.plmn.mcc_digit3 = std::stoul(
      guti_str.substr(idx++, chars_to_read), &chars_to_read, HEX_BASE_VAL);
  guti_m5_p->guamfi.plmn.mnc_digit1 = std::stoul(
      guti_str.substr(idx++, chars_to_read), &chars_to_read, HEX_BASE_VAL);
  guti_m5_p->guamfi.plmn.mnc_digit2 = std::stoul(
      guti_str.substr(idx++, chars_to_read), &chars_to_read, HEX_BASE_VAL);
  guti_m5_p->guamfi.plmn.mnc_digit3 = std::stoul(
      guti_str.substr(idx++, chars_to_read), &chars_to_read, HEX_BASE_VAL);
  chars_to_read                  = 2;
  guti_m5_p->guamfi.amf_regionid = std::stoul(
      guti_str.substr(idx, chars_to_read), &chars_to_read, HEX_BASE_VAL);
  idx += chars_to_read;
  chars_to_read                = 4;
  guti_m5_p->guamfi.amf_set_id = std::stoul(
      guti_str.substr(idx, chars_to_read), &chars_to_read, HEX_BASE_VAL);
  idx += chars_to_read;
  chars_to_read                 = 4;
  guti_m5_p->guamfi.amf_pointer = std::stoul(
      guti_str.substr(idx, chars_to_read), &chars_to_read, HEX_BASE_VAL);
  idx += chars_to_read;
  chars_to_read     = 8;
  guti_m5_p->m_tmsi = std::stoul(
      guti_str.substr(idx, chars_to_read), &chars_to_read, HEX_BASE_VAL);
}
// Converts Map<guti_m5_t,uint64_t> to proto
void AmfNasStateConverter::map_guti_uint64_to_proto(
    const map_guti_m5_uint64_t guti_map,
    google::protobuf::Map<std::string, uint64_t>* proto_map) {
  std::string guti_str;
  for (const auto& elm : guti_map.umap) {
    guti_str               = amf_app_convert_guti_m5_to_string(elm.first);
    (*proto_map)[guti_str] = elm.second;
  }
}

// Converts Proto to Map<guti_m5_t,uint64_t>
void AmfNasStateConverter::proto_to_guti_map(
    const google::protobuf::Map<std::string, uint64_t>& proto_map,
    map_guti_m5_uint64_t* guti_map) {
  for (auto const& kv : proto_map) {
    amf_ue_ngap_id_t amf_ue_ngap_id = kv.second;
    std::unique_ptr<guti_m5_t> guti = std::make_unique<guti_m5_t>();
    memset(guti.get(), 0, sizeof(guti_m5_t));
    // Converts guti to string.
    amf_app_convert_string_to_guti_m5(kv.first, guti.get());

    guti_m5_t guti_received = *guti.get();
    magma::map_rc_t m_rc    = guti_map->insert(guti_received, amf_ue_ngap_id);
    if (m_rc != magma::MAP_OK) {
      OAILOG_ERROR(
          LOG_AMF_APP,
          "Failed to insert amf_ue_ngap_id %lu in GUTI table, error: %s\n",
          amf_ue_ngap_id, map_rc_code2string(m_rc).c_str());
    }
  }
}

/*********************************************************
 *                AMF app state<-> Proto                  *
 * Functions to serialize/desearialize AMF app state      *
 * The caller is responsible for all memory management    *
 **********************************************************/

void AmfNasStateConverter::state_to_proto(
    const amf_app_desc_t* amf_nas_state_p, MmeNasState* state_proto) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  state_proto->set_mme_app_ue_s1ap_id_generator(
      amf_nas_state_p->amf_app_ue_ngap_id_generator);

  // These Functions are to be removed as part of the stateless enhancement
  // maps to proto
  auto amf_ue_ctxts_proto = state_proto->mutable_mme_ue_contexts();
  OAILOG_DEBUG(LOG_AMF_APP, "IMSI table to proto");
  magma::lte::StateConverter::map_uint64_uint64_to_proto(
      amf_nas_state_p->amf_ue_contexts.imsi_amf_ue_id_htbl,
      amf_ue_ctxts_proto->mutable_imsi_ue_id_htbl());
  magma::lte::StateConverter::map_uint64_uint64_to_proto(
      amf_nas_state_p->amf_ue_contexts.tun11_ue_context_htbl,
      amf_ue_ctxts_proto->mutable_tun11_ue_id_htbl());
  magma::lte::StateConverter::map_uint64_uint64_to_proto(
      amf_nas_state_p->amf_ue_contexts.gnb_ue_ngap_id_ue_context_htbl,
      amf_ue_ctxts_proto->mutable_enb_ue_id_ue_id_htbl());
  map_guti_uint64_to_proto(
      amf_nas_state_p->amf_ue_contexts.guti_ue_context_htbl,
      amf_ue_ctxts_proto->mutable_guti_ue_id_htbl());
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void AmfNasStateConverter::proto_to_state(
    const MmeNasState& state_proto, amf_app_desc_t* amf_nas_state_p) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  amf_nas_state_p->amf_app_ue_ngap_id_generator =
      state_proto.mme_app_ue_s1ap_id_generator();

  if (amf_nas_state_p->amf_app_ue_ngap_id_generator == 0) {  // uninitialized
    amf_nas_state_p->amf_app_ue_ngap_id_generator = 1;
  }
  OAILOG_INFO(LOG_AMF_APP, "Done reading AMF statistics from data store");

  magma::lte::oai::MmeUeContext amf_ue_ctxts_proto =
      state_proto.mme_ue_contexts();

  amf_ue_context_t* amf_ue_ctxt_state = &amf_nas_state_p->amf_ue_contexts;

  // proto to maps
  OAILOG_INFO(LOG_AMF_APP, "Hashtable AMF UE ID => IMSI");
  proto_to_map_uint64_uint64(
      amf_ue_ctxts_proto.imsi_ue_id_htbl(),
      &amf_ue_ctxt_state->imsi_amf_ue_id_htbl);
  proto_to_map_uint64_uint64(
      amf_ue_ctxts_proto.tun11_ue_id_htbl(),
      &amf_ue_ctxt_state->tun11_ue_context_htbl);
  proto_to_map_uint64_uint64(
      amf_ue_ctxts_proto.enb_ue_id_ue_id_htbl(),
      &amf_ue_ctxt_state->gnb_ue_ngap_id_ue_context_htbl);

  proto_to_guti_map(
      amf_ue_ctxts_proto.guti_ue_id_htbl(),
      &amf_ue_ctxt_state->guti_ue_context_htbl);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void AmfNasStateConverter::ue_to_proto(
    const ue_m5gmm_context_t* ue_ctxt,
    magma::lte::oai::UeContext* ue_ctxt_proto) {
  ue_m5gmm_context_to_proto(ue_ctxt, ue_ctxt_proto);
}

void AmfNasStateConverter::proto_to_ue(
    const magma::lte::oai::UeContext& ue_ctxt_proto,
    ue_m5gmm_context_t* ue_ctxt) {
  proto_to_ue_m5gmm_context(ue_ctxt_proto, ue_ctxt);
}

/*********************************************************
 *                UE Context <-> Proto                    *
 * Functions to serialize/desearialize UE context         *
 * The caller needs to acquire a lock on UE context       *
 **********************************************************/

void AmfNasStateConverter::ue_m5gmm_context_to_proto(
    const ue_m5gmm_context_t* state_ue_m5gmm_context,
    magma::lte::oai::UeContext* ue_context_proto) {
  ue_context_proto->set_rel_cause(
      state_ue_m5gmm_context->ue_context_rel_cause.present);
  ue_context_proto->set_mm_state(state_ue_m5gmm_context->mm_state);
  ue_context_proto->set_ecm_state(state_ue_m5gmm_context->ecm_state);

  ue_context_proto->set_sctp_assoc_id_key(
      state_ue_m5gmm_context->sctp_assoc_id_key);
  ue_context_proto->set_enb_ue_s1ap_id(state_ue_m5gmm_context->gnb_ue_ngap_id);
  ue_context_proto->set_enb_s1ap_id_key(
      state_ue_m5gmm_context->gnb_ngap_id_key);

  StateConverter::apn_config_profile_to_proto(
      state_ue_m5gmm_context->amf_context.apn_config_profile,
      ue_context_proto->mutable_apn_config());

  ue_context_proto->set_mme_teid_s11(state_ue_m5gmm_context->amf_teid_n11);
  StateConverter::ambr_to_proto(
      state_ue_m5gmm_context->amf_context.subscribed_ue_ambr,
      ue_context_proto->mutable_subscribed_ue_ambr());
  ue_context_proto->set_paging_retx_count(
      state_ue_m5gmm_context->paging_context.paging_retx_count);
}

void AmfNasStateConverter::proto_to_ue_m5gmm_context(
    const magma::lte::oai::UeContext& ue_context_proto,
    ue_m5gmm_context_t* state_ue_m5gmm_context) {
  state_ue_m5gmm_context->ue_context_rel_cause.present =
      static_cast<ngap_Cause_PR>(ue_context_proto.rel_cause());
  state_ue_m5gmm_context->mm_state =
      static_cast<m5gmm_state_t>(ue_context_proto.mm_state());
  state_ue_m5gmm_context->ecm_state =
      static_cast<m5gcm_state_t>(ue_context_proto.ecm_state());

  state_ue_m5gmm_context->sctp_assoc_id_key =
      ue_context_proto.sctp_assoc_id_key();
  state_ue_m5gmm_context->gnb_ue_ngap_id  = ue_context_proto.enb_ue_s1ap_id();
  state_ue_m5gmm_context->gnb_ngap_id_key = ue_context_proto.enb_s1ap_id_key();

  StateConverter::proto_to_apn_config_profile(
      ue_context_proto.apn_config(),
      &state_ue_m5gmm_context->amf_context.apn_config_profile);
  state_ue_m5gmm_context->amf_teid_n11 = ue_context_proto.mme_teid_s11();
  StateConverter::proto_to_ambr(
      ue_context_proto.subscribed_ue_ambr(),
      &state_ue_m5gmm_context->amf_context.subscribed_ue_ambr);

  // Initialize timers to INVALID IDs
  state_ue_m5gmm_context->m5_mobile_reachability_timer.id =
      AMF_APP_TIMER_INACTIVE_ID;
  state_ue_m5gmm_context->m5_implicit_detach_timer.id =
      AMF_APP_TIMER_INACTIVE_ID;
  state_ue_m5gmm_context->m5_initial_context_setup_rsp_timer =
      (amf_app_timer_t){AMF_APP_TIMER_INACTIVE_ID,
                        AMF_APP_INITIAL_CONTEXT_SETUP_RSP_TIMER_VALUE};
  state_ue_m5gmm_context->paging_context.m5_paging_response_timer =
      (amf_app_timer_t){AMF_APP_TIMER_INACTIVE_ID,
                        AMF_APP_PAGING_RESPONSE_TIMER_VALUE};
  state_ue_m5gmm_context->m5_ulr_response_timer = (amf_app_timer_t){
      AMF_APP_TIMER_INACTIVE_ID, AMF_APP_ULR_RESPONSE_TIMER_VALUE};
  state_ue_m5gmm_context->m5_ue_context_modification_timer = (amf_app_timer_t){
      AMF_APP_TIMER_INACTIVE_ID, AMF_APP_UE_CONTEXT_MODIFICATION_TIMER_VALUE};

  state_ue_m5gmm_context->paging_context.paging_retx_count =
      ue_context_proto.paging_retx_count();
}

}  // namespace magma5g
