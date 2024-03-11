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

#include "lte/gateway/c/core/oai/tasks/amf/amf_app_state_converter.hpp"
#include <vector>
#include <memory>
extern "C" {
#include "lte/gateway/c/core/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/lib/message_utils/bytes_to_ie.h"
#include "lte/gateway/c/core/oai/lib/message_utils/ie_to_bytes.h"
}

using magma::lte::oai::EmmContext;
using magma::lte::oai::EmmSecurityContext;
using magma::lte::oai::MmeNasState;
namespace magma5g {

AmfNasStateConverter::AmfNasStateConverter() = default;
AmfNasStateConverter::~AmfNasStateConverter() = default;

void AmfNasStateConverter::plmn_to_chars(const plmn_t& state_plmn,
                                         char* plmn_array) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  plmn_array[0] = static_cast<char>(state_plmn.mcc_digit1 + ASCII_ZERO);
  plmn_array[1] = static_cast<char>(state_plmn.mcc_digit2 + ASCII_ZERO);
  plmn_array[2] = static_cast<char>(state_plmn.mcc_digit3 + ASCII_ZERO);
  plmn_array[3] = static_cast<char>(state_plmn.mnc_digit1 + ASCII_ZERO);
  plmn_array[4] = static_cast<char>(state_plmn.mnc_digit2 + ASCII_ZERO);
  plmn_array[5] = static_cast<char>(state_plmn.mnc_digit3 + ASCII_ZERO);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void AmfNasStateConverter::chars_to_plmn(const char* plmn_array,
                                         plmn_t* state_plmn) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  state_plmn->mcc_digit1 = static_cast<int>(plmn_array[0]) - ASCII_ZERO;
  state_plmn->mcc_digit2 = static_cast<int>(plmn_array[1]) - ASCII_ZERO;
  state_plmn->mcc_digit3 = static_cast<int>(plmn_array[2]) - ASCII_ZERO;
  state_plmn->mnc_digit1 = static_cast<int>(plmn_array[3]) - ASCII_ZERO;
  state_plmn->mnc_digit2 = static_cast<int>(plmn_array[4]) - ASCII_ZERO;
  state_plmn->mnc_digit3 = static_cast<int>(plmn_array[5]) - ASCII_ZERO;
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

// HelperFunction: Converts guti_m5_t to std::string
std::string AmfNasStateConverter::amf_app_convert_guti_m5_to_string(
    const guti_m5_t& guti) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
#define GUTI_M5_STRING_LEN 25
  char* temp_str =
      reinterpret_cast<char*>(calloc(1, sizeof(char) * GUTI_M5_STRING_LEN));
  snprintf(temp_str, GUTI_M5_STRING_LEN, "%x%x%x%x%x%x%02x%04x%04x%08x",
           guti.guamfi.plmn.mcc_digit1, guti.guamfi.plmn.mcc_digit2,
           guti.guamfi.plmn.mcc_digit3, guti.guamfi.plmn.mnc_digit1,
           guti.guamfi.plmn.mnc_digit2, guti.guamfi.plmn.mnc_digit3,
           guti.guamfi.amf_regionid, guti.guamfi.amf_set_id,
           guti.guamfi.amf_pointer, guti.m_tmsi);
  std::string guti_str(temp_str);
  free(temp_str);
  OAILOG_FUNC_RETURN(LOG_AMF_APP, guti_str);
}

// HelperFunction: Converts std:: string back to guti_m5_t
void AmfNasStateConverter::amf_app_convert_string_to_guti_m5(
    const std::string& guti_str, guti_m5_t* guti_m5_p) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  int idx = 0;
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
  chars_to_read = 2;
  guti_m5_p->guamfi.amf_regionid = std::stoul(
      guti_str.substr(idx, chars_to_read), &chars_to_read, HEX_BASE_VAL);
  idx += chars_to_read;
  chars_to_read = 4;
  guti_m5_p->guamfi.amf_set_id = std::stoul(guti_str.substr(idx, chars_to_read),
                                            &chars_to_read, HEX_BASE_VAL);
  idx += chars_to_read;
  chars_to_read = 4;
  guti_m5_p->guamfi.amf_pointer = std::stoul(
      guti_str.substr(idx, chars_to_read), &chars_to_read, HEX_BASE_VAL);
  idx += chars_to_read;
  chars_to_read = 8;
  guti_m5_p->m_tmsi = std::stoul(guti_str.substr(idx, chars_to_read),
                                 &chars_to_read, HEX_BASE_VAL);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}
// Converts Map<guti_m5_t,uint64_t> to proto
void AmfNasStateConverter::map_guti_uint64_to_proto(
    const map_guti_m5_uint64_t guti_map,
    google::protobuf::Map<std::string, uint64_t>* proto_map) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  std::string guti_str;
  for (const auto& elm : guti_map.umap) {
    guti_str = amf_app_convert_guti_m5_to_string(elm.first);
    (*proto_map)[guti_str] = elm.second;
  }
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

// Converts Proto to Map<guti_m5_t,uint64_t>
void AmfNasStateConverter::proto_to_guti_map(
    const google::protobuf::Map<std::string, uint64_t>& proto_map,
    map_guti_m5_uint64_t* guti_map) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  for (auto const& kv : proto_map) {
    amf_ue_ngap_id_t amf_ue_ngap_id = kv.second;
    std::unique_ptr<guti_m5_t> guti = std::make_unique<guti_m5_t>();
    memset(guti.get(), 0, sizeof(guti_m5_t));
    // Converts guti to string.
    amf_app_convert_string_to_guti_m5(kv.first, guti.get());

    guti_m5_t guti_received = *guti.get();
    magma::map_rc_t m_rc = guti_map->insert(guti_received, amf_ue_ngap_id);
    if (m_rc != magma::MAP_OK) {
      OAILOG_ERROR(
          LOG_AMF_APP,
          "Failed to insert amf_ue_ngap_id %lu in GUTI table, error: %s\n",
          amf_ue_ngap_id, map_rc_code2string(m_rc).c_str());
    }
  }
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

/*********************************************************
 *                AMF app state<-> Proto                  *
 * Functions to serialize/desearialize AMF app state      *
 * The caller is responsible for all memory management    *
 **********************************************************/

void AmfNasStateConverter::state_to_proto(const amf_app_desc_t* amf_nas_state_p,
                                          MmeNasState* state_proto) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  state_proto->set_nb_ue_connected(amf_nas_state_p->nb_ue_connected);
  state_proto->set_nb_ue_attached(amf_nas_state_p->nb_ue_attached);
  state_proto->set_nb_ue_idle(amf_nas_state_p->nb_ue_idle);
  state_proto->set_nb_pdu_sessions(amf_nas_state_p->nb_pdu_sessions);
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

void AmfNasStateConverter::proto_to_state(const MmeNasState& state_proto,
                                          amf_app_desc_t* amf_nas_state_p) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  amf_nas_state_p->nb_ue_connected = state_proto.nb_ue_connected();
  amf_nas_state_p->nb_ue_attached = state_proto.nb_ue_attached();
  amf_nas_state_p->nb_ue_idle = state_proto.nb_ue_idle();
  amf_nas_state_p->nb_pdu_sessions = state_proto.nb_pdu_sessions();

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
  proto_to_map_uint64_uint64(amf_ue_ctxts_proto.imsi_ue_id_htbl(),
                             &amf_ue_ctxt_state->imsi_amf_ue_id_htbl);
  proto_to_map_uint64_uint64(amf_ue_ctxts_proto.tun11_ue_id_htbl(),
                             &amf_ue_ctxt_state->tun11_ue_context_htbl);
  proto_to_map_uint64_uint64(
      amf_ue_ctxts_proto.enb_ue_id_ue_id_htbl(),
      &amf_ue_ctxt_state->gnb_ue_ngap_id_ue_context_htbl);

  proto_to_guti_map(amf_ue_ctxts_proto.guti_ue_id_htbl(),
                    &amf_ue_ctxt_state->guti_ue_context_htbl);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void AmfNasStateConverter::ue_to_proto(
    const ue_m5gmm_context_t* ue_ctxt,
    magma::lte::oai::UeContext* ue_ctxt_proto) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  ue_m5gmm_context_to_proto(ue_ctxt, ue_ctxt_proto);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void AmfNasStateConverter::proto_to_ue(
    const magma::lte::oai::UeContext& ue_ctxt_proto,
    ue_m5gmm_context_t* ue_ctxt) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  proto_to_ue_m5gmm_context(ue_ctxt_proto, ue_ctxt);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

/*********************************************************
 *                UE Context <-> Proto                    *
 * Functions to serialize/desearialize UE context         *
 * The caller needs to acquire a lock on UE context       *
 **********************************************************/

void AmfNasStateConverter::ue_m5gmm_context_to_proto(
    const ue_m5gmm_context_t* state_ue_m5gmm_context,
    magma::lte::oai::UeContext* ue_context_proto) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  ue_context_proto->set_amf_ue_ngap_id(state_ue_m5gmm_context->amf_ue_ngap_id);
  ue_context_proto->set_rel_cause(state_ue_m5gmm_context->ue_context_rel_cause);

  ue_context_proto->set_mm_state(state_ue_m5gmm_context->mm_state);
  ue_context_proto->set_ecm_state(state_ue_m5gmm_context->cm_state);

  EmmContext* emm_ctx = ue_context_proto->mutable_emm_context();
  AmfNasStateConverter::amf_context_to_proto(
      &state_ue_m5gmm_context->amf_context, emm_ctx);
  ue_context_proto->set_sctp_assoc_id_key(
      state_ue_m5gmm_context->sctp_assoc_id_key);
  ue_context_proto->set_gnb_ue_ngap_id(state_ue_m5gmm_context->gnb_ue_ngap_id);
  ue_context_proto->set_gnb_ngap_id_key(
      state_ue_m5gmm_context->gnb_ngap_id_key);

  StateConverter::apn_config_profile_to_proto(
      state_ue_m5gmm_context->amf_context.apn_config_profile,
      ue_context_proto->mutable_apn_config());

  ue_context_proto->set_amf_teid_n11(state_ue_m5gmm_context->amf_teid_n11);
  StateConverter::ambr_to_proto(
      state_ue_m5gmm_context->amf_context.subscribed_ue_ambr,
      ue_context_proto->mutable_subscribed_ue_ambr());
  ue_context_proto->set_paging_retx_count(
      state_ue_m5gmm_context->paging_context.paging_retx_count);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void AmfNasStateConverter::proto_to_ue_m5gmm_context(
    const magma::lte::oai::UeContext& ue_context_proto,
    ue_m5gmm_context_t* state_ue_m5gmm_context) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  state_ue_m5gmm_context->amf_ue_ngap_id = ue_context_proto.amf_ue_ngap_id();
  state_ue_m5gmm_context->ue_context_rel_cause =
      static_cast<n2cause_e>(ue_context_proto.rel_cause());

  state_ue_m5gmm_context->mm_state =
      static_cast<m5gmm_state_t>(ue_context_proto.mm_state());
  state_ue_m5gmm_context->cm_state =
      static_cast<m5gcm_state_t>(ue_context_proto.ecm_state());

  AmfNasStateConverter::proto_to_amf_context(
      ue_context_proto.emm_context(), &state_ue_m5gmm_context->amf_context);
  state_ue_m5gmm_context->sctp_assoc_id_key =
      ue_context_proto.sctp_assoc_id_key();
  state_ue_m5gmm_context->gnb_ue_ngap_id = ue_context_proto.gnb_ue_ngap_id();
  state_ue_m5gmm_context->gnb_ngap_id_key = ue_context_proto.gnb_ngap_id_key();

  StateConverter::proto_to_apn_config_profile(
      ue_context_proto.apn_config(),
      &state_ue_m5gmm_context->amf_context.apn_config_profile);
  state_ue_m5gmm_context->amf_teid_n11 = ue_context_proto.amf_teid_n11();
  StateConverter::proto_to_ambr(
      ue_context_proto.subscribed_ue_ambr(),
      &state_ue_m5gmm_context->amf_context.subscribed_ue_ambr);

  // Initialize timers to INVALID IDs
  state_ue_m5gmm_context->m5_mobile_reachability_timer.id =
      AMF_APP_TIMER_INACTIVE_ID;
  state_ue_m5gmm_context->m5_implicit_deregistration_timer.id =
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
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void AmfNasStateConverter::tai_to_proto(const tai_t* state_tai,
                                        magma::lte::oai::Tai* tai_proto) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  OAILOG_DEBUG(LOG_AMF_APP, "State PLMN " PLMN_FMT "to proto",
               PLMN_ARG(&state_tai->plmn));
  char plmn_array[PLMN_BYTES] = {0};
  plmn_array[0] = static_cast<char>(state_tai->plmn.mcc_digit1 + ASCII_ZERO);
  plmn_array[1] = static_cast<char>(state_tai->plmn.mcc_digit2 + ASCII_ZERO);
  plmn_array[2] = static_cast<char>(state_tai->plmn.mcc_digit3 + ASCII_ZERO);
  plmn_array[3] = static_cast<char>(state_tai->plmn.mnc_digit1 + ASCII_ZERO);
  plmn_array[4] = static_cast<char>(state_tai->plmn.mnc_digit2 + ASCII_ZERO);
  plmn_array[5] = static_cast<char>(state_tai->plmn.mnc_digit3 + ASCII_ZERO);
  tai_proto->set_mcc_mnc(plmn_array);
  tai_proto->set_tac(state_tai->tac);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void AmfNasStateConverter::proto_to_tai(const magma::lte::oai::Tai& tai_proto,
                                        tai_t* state_tai) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  state_tai->plmn.mcc_digit1 =
      static_cast<int>(tai_proto.mcc_mnc()[0]) - ASCII_ZERO;
  state_tai->plmn.mcc_digit2 =
      static_cast<int>(tai_proto.mcc_mnc()[1]) - ASCII_ZERO;
  state_tai->plmn.mcc_digit3 =
      static_cast<int>(tai_proto.mcc_mnc()[2]) - ASCII_ZERO;
  state_tai->plmn.mnc_digit1 =
      static_cast<int>(tai_proto.mcc_mnc()[3]) - ASCII_ZERO;
  state_tai->plmn.mnc_digit2 =
      static_cast<int>(tai_proto.mcc_mnc()[4]) - ASCII_ZERO;
  state_tai->plmn.mnc_digit3 =
      static_cast<int>(tai_proto.mcc_mnc()[5]) - ASCII_ZERO;
  state_tai->tac = tai_proto.tac();
  OAILOG_DEBUG(LOG_AMF_APP, "State PLMN " PLMN_FMT "from proto",
               PLMN_ARG(&state_tai->plmn));
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void AmfNasStateConverter::guti_m5_to_proto(
    const guti_m5_t& state_guti_m5, magma::lte::oai::Guti_m5* guti_m5_proto) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  guti_m5_proto->Clear();
  char plmn_array[PLMN_BYTES] = {0};
  AmfNasStateConverter::plmn_to_chars(state_guti_m5.guamfi.plmn, plmn_array);
  guti_m5_proto->set_plmn(plmn_array);
  guti_m5_proto->set_amf_regionid(state_guti_m5.guamfi.amf_regionid);
  guti_m5_proto->set_amf_set_id(state_guti_m5.guamfi.amf_set_id);
  guti_m5_proto->set_amf_pointer(state_guti_m5.guamfi.amf_pointer);
  guti_m5_proto->set_m_tmsi(state_guti_m5.m_tmsi);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void AmfNasStateConverter::proto_to_guti_m5(
    const magma::lte::oai::Guti_m5& guti_m5_proto, guti_m5_t* state_guti_m5) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  AmfNasStateConverter::chars_to_plmn(guti_m5_proto.plmn().c_str(),
                                      &state_guti_m5->guamfi.plmn);
  state_guti_m5->guamfi.amf_regionid = guti_m5_proto.amf_regionid();
  state_guti_m5->guamfi.amf_set_id = guti_m5_proto.amf_set_id();
  state_guti_m5->guamfi.amf_pointer = guti_m5_proto.amf_pointer();
  state_guti_m5->m_tmsi = (tmsi_t)guti_m5_proto.m_tmsi();
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void AmfNasStateConverter::amf_context_to_proto(const amf_context_t* amf_ctx,
                                                EmmContext* emm_context_proto) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  emm_context_proto->set_imsi64(amf_ctx->imsi64);
  identity_tuple_to_proto<imsi_t>(
      &amf_ctx->imsi, emm_context_proto->mutable_imsi(), IMSI_BCD8_SIZE);
  emm_context_proto->set_saved_imsi64(amf_ctx->saved_imsi64);
  identity_tuple_to_proto<imei_t>(
      &amf_ctx->imei, emm_context_proto->mutable_imei(), IMEI_BCD8_SIZE);
  identity_tuple_to_proto<imeisv_t>(
      &amf_ctx->imeisv, emm_context_proto->mutable_imeisv(), IMEISV_BCD8_SIZE);

  AmfNasStateConverter::amf_security_context_to_proto(
      &amf_ctx->_security, emm_context_proto->mutable_security());
  emm_context_proto->set_emm_cause(amf_ctx->amf_cause);
  emm_context_proto->set_emm_fsm_state(amf_ctx->amf_fsm_state);
  emm_context_proto->set_attach_type(amf_ctx->m5gsregistrationtype);

  emm_context_proto->set_member_present_mask(amf_ctx->member_present_mask);
  emm_context_proto->set_member_valid_mask(amf_ctx->member_valid_mask);
  emm_context_proto->set_is_dynamic(amf_ctx->is_dynamic);
  emm_context_proto->set_is_attached(amf_ctx->is_registered);
  emm_context_proto->set_is_initial_identity_imsi(
      amf_ctx->is_initial_identity_imsi);
  emm_context_proto->set_is_guti_based_attach(
      amf_ctx->is_guti_based_registered);
  emm_context_proto->set_is_imsi_only_detach(amf_ctx->is_imsi_only_detach);
  tai_to_proto(&amf_ctx->originating_tai,
               emm_context_proto->mutable_originating_tai());
  emm_context_proto->set_ksi(amf_ctx->ksi);
  AmfNasStateConverter::smf_context_map_to_proto(
      amf_ctx->smf_ctxt_map,
      emm_context_proto->mutable_pdu_session_id_smf_context_map());
  AmfNasStateConverter::guti_m5_to_proto(amf_ctx->m5_guti,
                                         emm_context_proto->mutable_m5_guti());
  AmfNasStateConverter::guti_m5_to_proto(
      amf_ctx->m5_old_guti, emm_context_proto->mutable_m5_old_guti());
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void AmfNasStateConverter::proto_to_amf_context(
    const EmmContext& emm_context_proto, amf_context_t* amf_ctx) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  amf_ctx->imsi64 = emm_context_proto.imsi64();
  proto_to_identity_tuple<imsi_t>(emm_context_proto.imsi(), &amf_ctx->imsi,
                                  IMSI_BCD8_SIZE);
  amf_ctx->saved_imsi64 = emm_context_proto.saved_imsi64();

  proto_to_identity_tuple<imei_t>(emm_context_proto.imei(), &amf_ctx->imei,
                                  IMEI_BCD8_SIZE);
  proto_to_identity_tuple<imeisv_t>(emm_context_proto.imeisv(),
                                    &amf_ctx->imeisv, IMEISV_BCD8_SIZE);

  AmfNasStateConverter::proto_to_amf_security_context(
      emm_context_proto.security(), &amf_ctx->_security);
  amf_ctx->amf_cause = emm_context_proto.emm_cause();
  amf_ctx->amf_fsm_state = (amf_fsm_state_t)emm_context_proto.emm_fsm_state();
  amf_ctx->m5gsregistrationtype = emm_context_proto.attach_type();
  amf_ctx->member_present_mask = emm_context_proto.member_present_mask();
  amf_ctx->member_valid_mask = emm_context_proto.member_valid_mask();
  amf_ctx->is_dynamic = emm_context_proto.is_dynamic();
  amf_ctx->is_registered = emm_context_proto.is_attached();
  amf_ctx->is_initial_identity_imsi =
      emm_context_proto.is_initial_identity_imsi();
  amf_ctx->is_guti_based_registered = emm_context_proto.is_guti_based_attach();
  amf_ctx->is_imsi_only_detach = emm_context_proto.is_imsi_only_detach();
  proto_to_tai(emm_context_proto.originating_tai(), &amf_ctx->originating_tai);
  amf_ctx->ksi = emm_context_proto.ksi();
  AmfNasStateConverter::proto_to_smf_context_map(
      emm_context_proto.pdu_session_id_smf_context_map(),
      &amf_ctx->smf_ctxt_map);
  AmfNasStateConverter::proto_to_guti_m5(emm_context_proto.m5_guti(),
                                         &amf_ctx->m5_guti);
  AmfNasStateConverter::proto_to_guti_m5(emm_context_proto.m5_old_guti(),
                                         &amf_ctx->m5_old_guti);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}
void AmfNasStateConverter::amf_security_context_to_proto(
    const amf_security_context_t* state_amf_security_context,
    EmmSecurityContext* emm_security_context_proto) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  emm_security_context_proto->set_sc_type(state_amf_security_context->sc_type);
  emm_security_context_proto->set_eksi(state_amf_security_context->eksi);
  emm_security_context_proto->set_vector_index(
      state_amf_security_context->vector_index);
  emm_security_context_proto->set_knas_enc(state_amf_security_context->knas_enc,
                                           AUTH_KNAS_ENC_SIZE);
  emm_security_context_proto->set_knas_int(state_amf_security_context->knas_int,
                                           AUTH_KNAS_INT_SIZE);

  // Count values
  auto* dl_count_proto = emm_security_context_proto->mutable_dl_count();
  dl_count_proto->set_overflow(state_amf_security_context->dl_count.overflow);
  dl_count_proto->set_seq_num(state_amf_security_context->dl_count.seq_num);
  auto* ul_count_proto = emm_security_context_proto->mutable_ul_count();
  ul_count_proto->set_overflow(state_amf_security_context->ul_count.overflow);
  ul_count_proto->set_seq_num(state_amf_security_context->ul_count.seq_num);
  auto* kenb_ul_count_proto =
      emm_security_context_proto->mutable_kenb_ul_count();
  kenb_ul_count_proto->set_overflow(
      state_amf_security_context->kenb_ul_count.overflow);
  kenb_ul_count_proto->set_seq_num(
      state_amf_security_context->kenb_ul_count.seq_num);

  // Security algorithm
  auto* selected_algorithms_proto =
      emm_security_context_proto->mutable_selected_algos();
  selected_algorithms_proto->set_encryption(
      state_amf_security_context->selected_algorithms.encryption);
  selected_algorithms_proto->set_integrity(
      state_amf_security_context->selected_algorithms.integrity);
  emm_security_context_proto->set_direction_encode(
      state_amf_security_context->direction_encode);
  emm_security_context_proto->set_direction_decode(
      state_amf_security_context->direction_decode);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void AmfNasStateConverter::proto_to_amf_security_context(
    const EmmSecurityContext& emm_security_context_proto,
    amf_security_context_t* state_amf_security_context) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  state_amf_security_context->sc_type =
      (amf_sc_type_t)emm_security_context_proto.sc_type();
  state_amf_security_context->eksi = emm_security_context_proto.eksi();
  state_amf_security_context->vector_index =
      emm_security_context_proto.vector_index();
  memcpy(state_amf_security_context->knas_enc,
         emm_security_context_proto.knas_enc().c_str(), AUTH_KNAS_ENC_SIZE);
  memcpy(state_amf_security_context->knas_int,
         emm_security_context_proto.knas_int().c_str(), AUTH_KNAS_INT_SIZE);

  // Count values
  const auto& dl_count_proto = emm_security_context_proto.dl_count();
  state_amf_security_context->dl_count.overflow = dl_count_proto.overflow();
  state_amf_security_context->dl_count.seq_num = dl_count_proto.seq_num();
  const auto& ul_count_proto = emm_security_context_proto.ul_count();
  state_amf_security_context->ul_count.overflow = ul_count_proto.overflow();
  state_amf_security_context->ul_count.seq_num = ul_count_proto.seq_num();
  const auto& kenb_ul_count_proto = emm_security_context_proto.kenb_ul_count();
  state_amf_security_context->kenb_ul_count.overflow =
      kenb_ul_count_proto.overflow();
  state_amf_security_context->kenb_ul_count.seq_num =
      kenb_ul_count_proto.seq_num();

  // Security algorithm
  const auto& selected_algorithms_proto =
      emm_security_context_proto.selected_algos();
  state_amf_security_context->selected_algorithms.encryption =
      selected_algorithms_proto.encryption();
  state_amf_security_context->selected_algorithms.integrity =
      selected_algorithms_proto.integrity();
  state_amf_security_context->direction_encode =
      emm_security_context_proto.direction_encode();
  state_amf_security_context->direction_decode =
      emm_security_context_proto.direction_decode();
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void AmfNasStateConverter::smf_proc_data_to_proto(
    const smf_proc_data_t* state_smf_proc_data,
    magma::lte::oai::Smf_Proc_Data* smf_proc_data_proto) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  smf_proc_data_proto->set_pdu_session_id(state_smf_proc_data->pdu_session_id);
  smf_proc_data_proto->set_pti(state_smf_proc_data->pti);
  smf_proc_data_proto->set_message_type(
      static_cast<uint32_t>(state_smf_proc_data->message_type));
  smf_proc_data_proto->set_max_uplink(state_smf_proc_data->max_uplink);
  smf_proc_data_proto->set_max_downlink(state_smf_proc_data->max_downlink);
  smf_proc_data_proto->set_pdu_session_type(
      (magma::lte::oai::M5GPduSessionType)
          state_smf_proc_data->pdu_session_type);
  smf_proc_data_proto->set_ssc_mode(state_smf_proc_data->ssc_mode);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}
void AmfNasStateConverter::proto_to_smf_proc_data(
    const magma::lte::oai::Smf_Proc_Data& smf_proc_data_proto,
    smf_proc_data_t* state_smf_proc_data) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  state_smf_proc_data->pdu_session_id = smf_proc_data_proto.pdu_session_id();
  state_smf_proc_data->pti = smf_proc_data_proto.pti();
  state_smf_proc_data->message_type =
      static_cast<M5GMessageType>(smf_proc_data_proto.message_type());
  state_smf_proc_data->max_uplink = smf_proc_data_proto.max_uplink();
  state_smf_proc_data->max_downlink = smf_proc_data_proto.max_downlink();
  state_smf_proc_data->pdu_session_type =
      static_cast<M5GPduSessionType>(smf_proc_data_proto.pdu_session_type());
  state_smf_proc_data->ssc_mode = smf_proc_data_proto.ssc_mode();
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void AmfNasStateConverter::s_nssai_to_proto(
    const s_nssai_t* state_s_nssai, magma::lte::oai::SNssai* snassi_proto) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  snassi_proto->set_sst(state_s_nssai->sst);
  snassi_proto->set_sd((char*)state_s_nssai->sd, SD_LENGTH);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}
void AmfNasStateConverter::proto_to_s_nssai(
    const magma::lte::oai::SNssai& snassi_proto, s_nssai_t* state_s_nssai) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  state_s_nssai->sst = snassi_proto.sst();
  memcpy((void*)state_s_nssai->sd, (void*)snassi_proto.sd().c_str(), SD_LENGTH);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void AmfNasStateConverter::pco_protocol_or_container_id_to_proto(
    const protocol_configuration_options_t&
        state_protocol_configuration_options,
    magma::lte::oai::ProtocolConfigurationOptions*
        protocol_configuration_options_proto) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  for (int i = 0;
       i < state_protocol_configuration_options.num_protocol_or_container_id;
       i++) {
    pco_protocol_or_container_id_t state_pco_protocol_or_container_id =
        state_protocol_configuration_options.protocol_or_container_ids[i];
    auto pco_protocol_or_container_id_proto =
        protocol_configuration_options_proto->add_proto_or_container_id();
    pco_protocol_or_container_id_proto->set_id(
        state_pco_protocol_or_container_id.id);
    pco_protocol_or_container_id_proto->set_length(
        state_pco_protocol_or_container_id.length);
    if (state_pco_protocol_or_container_id.contents) {
      BSTRING_TO_STRING(state_pco_protocol_or_container_id.contents,
                        pco_protocol_or_container_id_proto->mutable_contents());
    }
  }
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void AmfNasStateConverter::proto_to_pco_protocol_or_container_id(
    const magma::lte::oai::ProtocolConfigurationOptions&
        protocol_configuration_options_proto,
    protocol_configuration_options_t* state_protocol_configuration_options) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  auto proto_pco_ids =
      protocol_configuration_options_proto.proto_or_container_id();
  int i = 0;
  for (auto ptr = proto_pco_ids.begin(); ptr < proto_pco_ids.end(); ptr++) {
    pco_protocol_or_container_id_t* state_pco_protocol_or_container_id =
        &state_protocol_configuration_options->protocol_or_container_ids[i];
    state_pco_protocol_or_container_id->id = ptr->id();
    state_pco_protocol_or_container_id->length = ptr->length();
    if (ptr->contents().length()) {
      state_pco_protocol_or_container_id->contents = bfromcstr_with_str_len(
          ptr->contents().c_str(), ptr->contents().length());
    }
    i++;
  }
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void AmfNasStateConverter::protocol_configuration_options_to_proto(
    const protocol_configuration_options_t&
        state_protocol_configuration_options,
    magma::lte::oai::ProtocolConfigurationOptions*
        protocol_configuration_options_proto) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  protocol_configuration_options_proto->set_ext(
      state_protocol_configuration_options.ext);
  protocol_configuration_options_proto->set_spare(
      state_protocol_configuration_options.spare);
  protocol_configuration_options_proto->set_config_protocol(
      state_protocol_configuration_options.configuration_protocol);
  protocol_configuration_options_proto->set_num_protocol_or_container_id(
      state_protocol_configuration_options.num_protocol_or_container_id);

  AmfNasStateConverter::pco_protocol_or_container_id_to_proto(
      state_protocol_configuration_options,
      protocol_configuration_options_proto);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void AmfNasStateConverter::proto_to_protocol_configuration_options(
    const magma::lte::oai::ProtocolConfigurationOptions&
        protocol_configuration_options_proto,
    protocol_configuration_options_t* state_protocol_configuration_options) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  state_protocol_configuration_options->ext =
      protocol_configuration_options_proto.ext();
  state_protocol_configuration_options->spare =
      protocol_configuration_options_proto.spare();
  state_protocol_configuration_options->configuration_protocol =
      protocol_configuration_options_proto.config_protocol();
  state_protocol_configuration_options->num_protocol_or_container_id =
      protocol_configuration_options_proto.num_protocol_or_container_id();
  AmfNasStateConverter::proto_to_pco_protocol_or_container_id(
      protocol_configuration_options_proto,
      state_protocol_configuration_options);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void AmfNasStateConverter::session_ambr_to_proto(
    const session_ambr_t& state_session_ambr,
    magma::lte::oai::Ambr* ambr_proto) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  ambr_proto->set_br_ul(state_session_ambr.ul_session_ambr);
  ambr_proto->set_br_dl(state_session_ambr.dl_session_ambr);
  ambr_proto->set_br_unit(static_cast<magma::lte::oai::Ambr::BitrateUnitsAMBR>(
      state_session_ambr.dl_ambr_unit));
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}
void AmfNasStateConverter::proto_to_session_ambr(
    const magma::lte::oai::Ambr& ambr_proto,
    session_ambr_t* state_session_ambr) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  state_session_ambr->dl_ambr_unit =
      static_cast<M5GSessionAmbrUnit>(ambr_proto.br_unit());
  state_session_ambr->dl_session_ambr = ambr_proto.br_dl();
  state_session_ambr->ul_ambr_unit =
      static_cast<M5GSessionAmbrUnit>(ambr_proto.br_unit());
  state_session_ambr->ul_session_ambr = ambr_proto.br_ul();
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void AmfNasStateConverter::qos_flow_level_parameters_to_proto(
    const qos_flow_level_qos_parameters& state_qos_flow_parameters,
    magma::lte::oai::QosFlowParameters* qos_flow_parameters_proto) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  qos_flow_parameters_proto->set_fiveqi(
      state_qos_flow_parameters.qos_characteristic.non_dynamic_5QI_desc.fiveQI);
  qos_flow_parameters_proto->set_priority_level(
      state_qos_flow_parameters.alloc_reten_priority.priority_level);
  qos_flow_parameters_proto->set_preemption_vulnerability(
      state_qos_flow_parameters.alloc_reten_priority.pre_emption_vul);
  qos_flow_parameters_proto->set_preemption_capability(
      state_qos_flow_parameters.alloc_reten_priority.pre_emption_cap);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void AmfNasStateConverter::proto_to_qos_flow_level_parameters(
    const magma::lte::oai::QosFlowParameters& qos_flow_parameters_proto,
    qos_flow_level_qos_parameters* state_qos_flow_parameters) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  state_qos_flow_parameters->qos_characteristic.non_dynamic_5QI_desc.fiveQI =
      qos_flow_parameters_proto.fiveqi();
  state_qos_flow_parameters->alloc_reten_priority.priority_level =
      qos_flow_parameters_proto.priority_level();
  state_qos_flow_parameters->alloc_reten_priority.pre_emption_vul =
      static_cast<pre_emption_vulnerability>(
          qos_flow_parameters_proto.preemption_vulnerability());
  state_qos_flow_parameters->alloc_reten_priority.pre_emption_cap =
      static_cast<pre_emption_capability>(
          qos_flow_parameters_proto.preemption_capability());
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void AmfNasStateConverter::qos_flow_setup_request_item_to_proto(
    const qos_flow_setup_request_item& state_qos_flow_request_item,
    magma::lte::oai::M5GQosFlowItem* qos_flow_item_proto) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  qos_flow_item_proto->set_qfi(state_qos_flow_request_item.qos_flow_identifier);
  AmfNasStateConverter::qos_flow_level_parameters_to_proto(
      state_qos_flow_request_item.qos_flow_level_qos_param,
      qos_flow_item_proto->mutable_qos_flow_param());
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void AmfNasStateConverter::proto_to_qos_flow_setup_request_item(
    const magma::lte::oai::M5GQosFlowItem& qos_flow_item_proto,
    qos_flow_setup_request_item* state_qos_flow_request_item) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  state_qos_flow_request_item->qos_flow_identifier = qos_flow_item_proto.qfi();
  AmfNasStateConverter::proto_to_qos_flow_level_parameters(
      qos_flow_item_proto.qos_flow_param(),
      &state_qos_flow_request_item->qos_flow_level_qos_param);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void AmfNasStateConverter::smf_context_map_to_proto(
    const std::unordered_map<uint8_t, std::shared_ptr<smf_context_t>>&
        smf_ctxt_map,
    google::protobuf::Map<uint32_t, magma::lte::oai::SmfContext>* proto_map) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  for (const auto& it : smf_ctxt_map) {
    magma::lte::oai::SmfContext smf_context_proto =
        magma::lte::oai::SmfContext();
    AmfNasStateConverter::smf_context_to_proto(it.second.get(),
                                               &smf_context_proto);
    (*proto_map)[it.first] = smf_context_proto;
  }
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void AmfNasStateConverter::proto_to_smf_context_map(
    const google::protobuf::Map<uint32_t, magma::lte::oai::SmfContext>&
        proto_map,
    std::unordered_map<uint8_t, std::shared_ptr<smf_context_t>>* smf_ctxt_map) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  for (auto const& kv : proto_map) {
    smf_context_t smf_ctx;
    AmfNasStateConverter::proto_to_smf_context(kv.second, &smf_ctx);
    (*smf_ctxt_map)[kv.first] = std::make_shared<smf_context_t>(smf_ctx);
  }
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}
// smf_context to proto and proto to smf_context
void AmfNasStateConverter::smf_context_to_proto(
    const smf_context_t* state_smf_context,
    magma::lte::oai::SmfContext* smf_context_proto) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  smf_context_proto->set_sm_session_state(state_smf_context->pdu_session_state);
  smf_context_proto->set_pdu_session_version(
      state_smf_context->pdu_session_version);
  smf_context_proto->set_active_pdu_sessions(state_smf_context->n_active_pdus);
  smf_context_proto->set_is_emergency(state_smf_context->is_emergency);
  AmfNasStateConverter::session_ambr_to_proto(
      state_smf_context->selected_ambr,
      smf_context_proto->mutable_selected_ambr());

  smf_context_proto->set_gnb_gtp_teid(
      state_smf_context->gtp_tunnel_id.gnb_gtp_teid);

  char gnb_gtp_teid_ip_addr_str[16] = {0};
  inet_ntop(AF_INET, state_smf_context->gtp_tunnel_id.gnb_gtp_teid_ip_addr,
            gnb_gtp_teid_ip_addr_str, INET_ADDRSTRLEN);
  smf_context_proto->set_gnb_gtp_teid_ip_addr(gnb_gtp_teid_ip_addr_str);

  smf_context_proto->set_upf_gtp_teid(
      (char*)state_smf_context->gtp_tunnel_id.upf_gtp_teid, 4);

  char upf_gtp_teid_ip_addr_str[16] = {0};
  inet_ntop(AF_INET, state_smf_context->gtp_tunnel_id.upf_gtp_teid_ip_addr,
            upf_gtp_teid_ip_addr_str, INET_ADDRSTRLEN);
  smf_context_proto->set_upf_gtp_teid_ip_addr(upf_gtp_teid_ip_addr_str);

  bstring bstr_buffer = paa_to_bstring(&state_smf_context->pdu_address);
  BSTRING_TO_STRING(bstr_buffer, smf_context_proto->mutable_paa());
  bdestroy(bstr_buffer);

  StateConverter::ambr_to_proto(state_smf_context->apn_ambr,
                                smf_context_proto->mutable_apn_ambr());

  AmfNasStateConverter::smf_proc_data_to_proto(
      &state_smf_context->smf_proc_data,
      smf_context_proto->mutable_smf_proc_data());
  smf_context_proto->set_retransmission_count(
      state_smf_context->retransmission_count);
  AmfNasStateConverter::protocol_configuration_options_to_proto(
      state_smf_context->pco, smf_context_proto->mutable_pco());
  smf_context_proto->set_dnn_in_use(state_smf_context->dnn);

  AmfNasStateConverter::s_nssai_to_proto(
      &state_smf_context->requested_nssai,
      smf_context_proto->mutable_requested_nssai());

  AmfNasStateConverter::qos_flow_setup_request_item_to_proto(

      state_smf_context->smf_proc_data.qos_flow_list.item[0].qos_flow_req_item,
      smf_context_proto->mutable_smf_proc_data()->mutable_qos_flow_list());
}

void AmfNasStateConverter::proto_to_smf_context(
    const magma::lte::oai::SmfContext& smf_context_proto,
    smf_context_t* state_smf_context) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  state_smf_context->pdu_session_state =
      (SMSessionFSMState)smf_context_proto.sm_session_state();
  state_smf_context->pdu_session_version =
      smf_context_proto.pdu_session_version();
  state_smf_context->n_active_pdus = smf_context_proto.active_pdu_sessions();
  state_smf_context->is_emergency = smf_context_proto.is_emergency();
  AmfNasStateConverter::proto_to_session_ambr(
      smf_context_proto.selected_ambr(), &state_smf_context->selected_ambr);
  state_smf_context->gtp_tunnel_id.gnb_gtp_teid =
      smf_context_proto.gnb_gtp_teid();

  memset(&state_smf_context->gtp_tunnel_id.gnb_gtp_teid_ip_addr, '\0',
         sizeof(state_smf_context->gtp_tunnel_id.gnb_gtp_teid_ip_addr));
  inet_pton(AF_INET, smf_context_proto.gnb_gtp_teid_ip_addr().c_str(),
            &(state_smf_context->gtp_tunnel_id.gnb_gtp_teid_ip_addr));

  memcpy((void*)state_smf_context->gtp_tunnel_id.upf_gtp_teid,
         (void*)smf_context_proto.upf_gtp_teid().c_str(), 4);

  memset(&state_smf_context->gtp_tunnel_id.upf_gtp_teid_ip_addr, '\0',
         sizeof(state_smf_context->gtp_tunnel_id.upf_gtp_teid_ip_addr));
  inet_pton(AF_INET, smf_context_proto.upf_gtp_teid_ip_addr().c_str(),
            &(state_smf_context->gtp_tunnel_id.upf_gtp_teid_ip_addr));

  bstring bstr_buffer;
  STRING_TO_BSTRING(smf_context_proto.paa(), bstr_buffer);
  bstring_to_paa(bstr_buffer, &state_smf_context->pdu_address);
  bdestroy(bstr_buffer);

  StateConverter::proto_to_ambr(smf_context_proto.apn_ambr(),
                                &state_smf_context->apn_ambr);

  AmfNasStateConverter::proto_to_smf_proc_data(
      smf_context_proto.smf_proc_data(), &state_smf_context->smf_proc_data);

  state_smf_context->retransmission_count =
      smf_context_proto.retransmission_count();

  AmfNasStateConverter::proto_to_protocol_configuration_options(
      smf_context_proto.pco(), &state_smf_context->pco);

  state_smf_context->dnn = smf_context_proto.dnn_in_use();

  AmfNasStateConverter::proto_to_s_nssai(smf_context_proto.requested_nssai(),
                                         &state_smf_context->requested_nssai);

  AmfNasStateConverter::proto_to_qos_flow_setup_request_item(
      smf_context_proto.smf_proc_data().qos_flow_list(),
      &state_smf_context->smf_proc_data.qos_flow_list.item[0]
           .qos_flow_req_item);
}
}  // namespace magma5g
