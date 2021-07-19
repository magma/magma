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
 *-----------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */
extern "C" {
#include "log.h"
#include "dynamic_memory_check.h"
}

#include "nas_state_converter.h"
//#include "spgw_state_converter.h"

namespace magma {
namespace lte {

NasStateConverter::NasStateConverter()  = default;
NasStateConverter::~NasStateConverter() = default;

/*************************************************/
/*        Common Types <-> Proto                 */
/*************************************************/

void NasStateConverter::proto_to_guti(
    const oai::Guti& guti_proto, guti_t* state_guti) {
  state_guti->gummei.plmn.mcc_digit1 =
      ((int) guti_proto.plmn()[0]) - ASCII_ZERO;
  state_guti->gummei.plmn.mcc_digit2 =
      ((int) guti_proto.plmn()[1]) - ASCII_ZERO;
  state_guti->gummei.plmn.mcc_digit3 =
      ((int) guti_proto.plmn()[2]) - ASCII_ZERO;
  state_guti->gummei.plmn.mnc_digit1 =
      ((int) guti_proto.plmn()[3]) - ASCII_ZERO;
  state_guti->gummei.plmn.mnc_digit2 =
      ((int) guti_proto.plmn()[4]) - ASCII_ZERO;
  state_guti->gummei.plmn.mnc_digit3 =
      ((int) guti_proto.plmn()[5]) - ASCII_ZERO;

  state_guti->gummei.mme_gid  = guti_proto.mme_gid();
  state_guti->gummei.mme_code = guti_proto.mme_code();
  state_guti->m_tmsi          = (tmsi_t) guti_proto.m_tmsi();
}

void NasStateConverter::proto_to_ecgi(
    const oai::Ecgi& ecgi_proto, ecgi_t* state_ecgi) {
  state_ecgi->plmn.mcc_digit1 = (int) (ecgi_proto.plmn()[0]) - ASCII_ZERO;
  state_ecgi->plmn.mcc_digit2 = (int) (ecgi_proto.plmn()[1]) - ASCII_ZERO;
  state_ecgi->plmn.mcc_digit3 = (int) (ecgi_proto.plmn()[2]) - ASCII_ZERO;
  state_ecgi->plmn.mnc_digit1 = (int) (ecgi_proto.plmn()[3]) - ASCII_ZERO;
  state_ecgi->plmn.mnc_digit2 = (int) (ecgi_proto.plmn()[4]) - ASCII_ZERO;
  state_ecgi->plmn.mnc_digit3 = (int) (ecgi_proto.plmn()[5]) - ASCII_ZERO;

  state_ecgi->cell_identity.enb_id  = ecgi_proto.enb_id();
  state_ecgi->cell_identity.cell_id = ecgi_proto.cell_id();
  state_ecgi->cell_identity.empty   = ecgi_proto.empty();
}

void NasStateConverter::partial_tai_list_to_proto(
    const partial_tai_list_t* state_partial_tai_list,
    oai::PartialTaiList* partial_tai_list_proto) {
  partial_tai_list_proto->set_type_of_list(state_partial_tai_list->typeoflist);
  partial_tai_list_proto->set_number_of_elements(
      state_partial_tai_list->numberofelements);
  switch (state_partial_tai_list->typeoflist) {
    case TRACKING_AREA_IDENTITY_LIST_MANY_PLMNS: {
      for (int idx = 0; idx < TRACKING_AREA_IDENTITY_LIST_MAXIMUM_NUM_TAI;
           idx++) {
        oai::Tai* proto_many_plmn = partial_tai_list_proto->add_tai_many_plmn();
        tai_to_proto(
            &state_partial_tai_list->u.tai_many_plmn[idx], proto_many_plmn);
      }
    } break;
    case TRACKING_AREA_IDENTITY_LIST_ONE_PLMN_CONSECUTIVE_TACS: {
      tai_to_proto(
          &state_partial_tai_list->u.tai_one_plmn_consecutive_tacs,
          partial_tai_list_proto->mutable_tai_one_plmn_consecutive_tacs());
    } break;
    case TRACKING_AREA_IDENTITY_LIST_ONE_PLMN_NON_CONSECUTIVE_TACS: {
      char plmn_array[PLMN_BYTES];
      plmn_array[0] =
          (char) (state_partial_tai_list->u.tai_one_plmn_non_consecutive_tacs.plmn.mcc_digit1 + ASCII_ZERO);
      plmn_array[1] =
          (char) (state_partial_tai_list->u.tai_one_plmn_non_consecutive_tacs.plmn.mcc_digit2 + ASCII_ZERO);
      plmn_array[2] =
          (char) (state_partial_tai_list->u.tai_one_plmn_non_consecutive_tacs.plmn.mcc_digit3 + ASCII_ZERO);
      plmn_array[3] =
          (char) (state_partial_tai_list->u.tai_one_plmn_non_consecutive_tacs.plmn.mnc_digit1 + ASCII_ZERO);
      plmn_array[4] =
          (char) (state_partial_tai_list->u.tai_one_plmn_non_consecutive_tacs.plmn.mnc_digit2 + ASCII_ZERO);
      plmn_array[5] =
          (char) (state_partial_tai_list->u.tai_one_plmn_non_consecutive_tacs.plmn.mnc_digit3 + ASCII_ZERO);
      partial_tai_list_proto->set_plmn(plmn_array);
      for (int idx = 0; idx < TRACKING_AREA_IDENTITY_LIST_MAXIMUM_NUM_TAI;
           idx++) {
        partial_tai_list_proto->add_tac(
            state_partial_tai_list->u.tai_one_plmn_non_consecutive_tacs
                .tac[idx]);
      }
    } break;
  }
}

void NasStateConverter::proto_to_partial_tai_list(
    const oai::PartialTaiList& partial_tai_list_proto,
    partial_tai_list_t* state_partial_tai_list) {
  state_partial_tai_list->typeoflist = partial_tai_list_proto.type_of_list();
  state_partial_tai_list->numberofelements =
      partial_tai_list_proto.number_of_elements();
  switch (state_partial_tai_list->typeoflist) {
    case TRACKING_AREA_IDENTITY_LIST_MANY_PLMNS: {
      for (int idx = 0; idx < TRACKING_AREA_IDENTITY_LIST_MAXIMUM_NUM_TAI;
           idx++) {
        proto_to_tai(
            partial_tai_list_proto.tai_many_plmn(idx),
            &state_partial_tai_list->u.tai_many_plmn[idx]);
      }
    } break;
    case TRACKING_AREA_IDENTITY_LIST_ONE_PLMN_CONSECUTIVE_TACS: {
      proto_to_tai(
          partial_tai_list_proto.tai_one_plmn_consecutive_tacs(),
          &state_partial_tai_list->u.tai_one_plmn_consecutive_tacs);
    } break;
    case TRACKING_AREA_IDENTITY_LIST_ONE_PLMN_NON_CONSECUTIVE_TACS: {
      state_partial_tai_list->u.tai_one_plmn_non_consecutive_tacs.plmn
          .mcc_digit1 = (int) (partial_tai_list_proto.plmn()[0]) - ASCII_ZERO;
      state_partial_tai_list->u.tai_one_plmn_non_consecutive_tacs.plmn
          .mcc_digit2 = (int) (partial_tai_list_proto.plmn()[1]) - ASCII_ZERO;
      state_partial_tai_list->u.tai_one_plmn_non_consecutive_tacs.plmn
          .mcc_digit3 = (int) (partial_tai_list_proto.plmn()[2]) - ASCII_ZERO;
      state_partial_tai_list->u.tai_one_plmn_non_consecutive_tacs.plmn
          .mnc_digit1 = (int) (partial_tai_list_proto.plmn()[3]) - ASCII_ZERO;
      state_partial_tai_list->u.tai_one_plmn_non_consecutive_tacs.plmn
          .mnc_digit2 = (int) (partial_tai_list_proto.plmn()[4]) - ASCII_ZERO;
      state_partial_tai_list->u.tai_one_plmn_non_consecutive_tacs.plmn
          .mnc_digit3 = (int) (partial_tai_list_proto.plmn()[5]) - ASCII_ZERO;
      for (int idx = 0; idx < TRACKING_AREA_IDENTITY_LIST_MAXIMUM_NUM_TAI;
           idx++) {
        state_partial_tai_list->u.tai_one_plmn_non_consecutive_tacs.tac[idx] =
            partial_tai_list_proto.tac(idx);
      }
    } break;
  }
}

void NasStateConverter::tai_list_to_proto(
    const tai_list_t* state_tai_list, oai::TaiList* tai_list_proto) {
  tai_list_proto->set_numberoflists(state_tai_list->numberoflists);
  for (int idx = 0; idx < state_tai_list->numberoflists; idx++) {
    oai::PartialTaiList* partial_tai_list =
        tai_list_proto->add_partial_tai_lists();
    partial_tai_list_to_proto(
        &state_tai_list->partial_tai_list[idx], partial_tai_list);
  }
}

void NasStateConverter::proto_to_tai_list(
    const oai::TaiList& tai_list_proto, tai_list_t* state_tai_list) {
  state_tai_list->numberoflists = tai_list_proto.numberoflists();
  for (int idx = 0; idx < state_tai_list->numberoflists; idx++) {
    proto_to_partial_tai_list(
        tai_list_proto.partial_tai_lists(idx),
        &state_tai_list->partial_tai_list[idx]);
  }
}

void NasStateConverter::tai_to_proto(
    const tai_t* state_tai, oai::Tai* tai_proto) {
  OAILOG_DEBUG(
      LOG_MME_APP, "State PLMN " PLMN_FMT "to proto",
      PLMN_ARG(&state_tai->plmn));
  char plmn_array[PLMN_BYTES];
  plmn_array[0] = (char) (state_tai->plmn.mcc_digit1 + ASCII_ZERO);
  plmn_array[1] = (char) (state_tai->plmn.mcc_digit2 + ASCII_ZERO);
  plmn_array[2] = (char) (state_tai->plmn.mcc_digit3 + ASCII_ZERO);
  plmn_array[3] = (char) (state_tai->plmn.mnc_digit1 + ASCII_ZERO);
  plmn_array[4] = (char) (state_tai->plmn.mnc_digit2 + ASCII_ZERO);
  plmn_array[5] = (char) (state_tai->plmn.mnc_digit3 + ASCII_ZERO);
  tai_proto->set_mcc_mnc(plmn_array);
  tai_proto->set_tac(state_tai->tac);
}

void NasStateConverter::proto_to_tai(
    const oai::Tai& tai_proto, tai_t* state_tai) {
  state_tai->plmn.mcc_digit1 = (int) (tai_proto.mcc_mnc()[0]) - ASCII_ZERO;
  state_tai->plmn.mcc_digit2 = (int) (tai_proto.mcc_mnc()[1]) - ASCII_ZERO;
  state_tai->plmn.mcc_digit3 = (int) (tai_proto.mcc_mnc()[2]) - ASCII_ZERO;
  state_tai->plmn.mnc_digit1 = (int) (tai_proto.mcc_mnc()[3]) - ASCII_ZERO;
  state_tai->plmn.mnc_digit2 = (int) (tai_proto.mcc_mnc()[4]) - ASCII_ZERO;
  state_tai->plmn.mnc_digit3 = (int) (tai_proto.mcc_mnc()[5]) - ASCII_ZERO;
  state_tai->tac             = tai_proto.tac();
  OAILOG_DEBUG(
      LOG_MME_APP, "State PLMN " PLMN_FMT "from proto",
      PLMN_ARG(&state_tai->plmn));
}

/*************************************************/
/*        ESM State <-> Proto                  */
/*************************************************/
void NasStateConverter::ambr_to_proto(
    const ambr_t& state_ambr, oai::Ambr* ambr_proto) {
  ambr_proto->set_br_ul(state_ambr.br_ul);
  ambr_proto->set_br_dl(state_ambr.br_dl);
  ambr_proto->set_br_unit(
      static_cast<magma::lte::oai::Ambr::BitrateUnitsAMBR>(state_ambr.br_unit));
}

void NasStateConverter::proto_to_ambr(
    const oai::Ambr& ambr_proto, ambr_t* state_ambr) {
  state_ambr->br_ul   = ambr_proto.br_ul();
  state_ambr->br_dl   = ambr_proto.br_dl();
  state_ambr->br_unit = (apn_ambr_bitrate_unit_t) ambr_proto.br_unit();
}

void NasStateConverter::bearer_qos_to_proto(
    const bearer_qos_t& state_bearer_qos, oai::BearerQos* bearer_qos_proto) {
  bearer_qos_proto->set_pci(state_bearer_qos.pci);
  bearer_qos_proto->set_pl(state_bearer_qos.pl);
  bearer_qos_proto->set_pvi(state_bearer_qos.pvi);
  bearer_qos_proto->set_qci(state_bearer_qos.qci);
  ambr_to_proto(state_bearer_qos.gbr, bearer_qos_proto->mutable_gbr());
  ambr_to_proto(state_bearer_qos.mbr, bearer_qos_proto->mutable_mbr());
}

void NasStateConverter::proto_to_bearer_qos(
    const oai::BearerQos& bearer_qos_proto, bearer_qos_t* state_bearer_qos) {
  state_bearer_qos->pci = bearer_qos_proto.pci();
  state_bearer_qos->pl  = bearer_qos_proto.pl();
  state_bearer_qos->pvi = bearer_qos_proto.pvi();
  state_bearer_qos->qci = bearer_qos_proto.qci();
  proto_to_ambr(bearer_qos_proto.gbr(), &state_bearer_qos->gbr);
  proto_to_ambr(bearer_qos_proto.mbr(), &state_bearer_qos->mbr);
}

void NasStateConverter::pco_protocol_or_container_id_to_proto(
    const protocol_configuration_options_t&
        state_protocol_configuration_options,
    oai::ProtocolConfigurationOptions* protocol_configuration_options_proto) {
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
      BSTRING_TO_STRING(
          state_pco_protocol_or_container_id.contents,
          pco_protocol_or_container_id_proto->mutable_contents());
    }
  }
}

void NasStateConverter::proto_to_pco_protocol_or_container_id(
    const oai::ProtocolConfigurationOptions&
        protocol_configuration_options_proto,
    protocol_configuration_options_t* state_protocol_configuration_options) {
  auto proto_pco_ids =
      protocol_configuration_options_proto.proto_or_container_id();
  int i = 0;
  for (auto ptr = proto_pco_ids.begin(); ptr < proto_pco_ids.end(); ptr++) {
    pco_protocol_or_container_id_t state_pco_protocol_or_container_id =
        state_protocol_configuration_options->protocol_or_container_ids[i];
    state_pco_protocol_or_container_id.id     = ptr->id();
    state_pco_protocol_or_container_id.length = ptr->length();
    if (ptr->contents().length()) {
      state_pco_protocol_or_container_id.contents = bfromcstr_with_str_len(
          ptr->contents().c_str(), ptr->contents().length());
    }
    i++;
  }
}

void NasStateConverter::protocol_configuration_options_to_proto(
    const protocol_configuration_options_t&
        state_protocol_configuration_options,
    oai::ProtocolConfigurationOptions* protocol_configuration_options_proto) {
  protocol_configuration_options_proto->set_ext(
      state_protocol_configuration_options.ext);
  protocol_configuration_options_proto->set_spare(
      state_protocol_configuration_options.spare);
  protocol_configuration_options_proto->set_config_protocol(
      state_protocol_configuration_options.configuration_protocol);
  protocol_configuration_options_proto->set_num_protocol_or_container_id(
      state_protocol_configuration_options.num_protocol_or_container_id);
  pco_protocol_or_container_id_to_proto(
      state_protocol_configuration_options,
      protocol_configuration_options_proto);
}

void NasStateConverter::proto_to_protocol_configuration_options(
    const oai::ProtocolConfigurationOptions&
        protocol_configuration_options_proto,
    protocol_configuration_options_t* state_protocol_configuration_options) {
  state_protocol_configuration_options->ext =
      protocol_configuration_options_proto.ext();
  state_protocol_configuration_options->spare =
      protocol_configuration_options_proto.spare();
  state_protocol_configuration_options->configuration_protocol =
      protocol_configuration_options_proto.config_protocol();
  state_protocol_configuration_options->num_protocol_or_container_id =
      protocol_configuration_options_proto.num_protocol_or_container_id();
  proto_to_pco_protocol_or_container_id(
      protocol_configuration_options_proto,
      state_protocol_configuration_options);
}

void NasStateConverter::esm_ebr_timer_data_to_proto(
    const esm_ebr_timer_data_t& state_esm_ebr_timer_data,
    oai::EsmEbrTimerData* proto_esm_ebr_timer_data) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  proto_esm_ebr_timer_data->set_ue_id(state_esm_ebr_timer_data.ue_id);
  proto_esm_ebr_timer_data->set_ebi(state_esm_ebr_timer_data.ebi);
  proto_esm_ebr_timer_data->set_count(state_esm_ebr_timer_data.count);
  if (state_esm_ebr_timer_data.msg) {
    BSTRING_TO_STRING(
        state_esm_ebr_timer_data.msg,
        proto_esm_ebr_timer_data->mutable_esm_msg());
  }
  OAILOG_FUNC_OUT(LOG_NAS_ESM);
}

void NasStateConverter::proto_to_esm_ebr_timer_data(
    const oai::EsmEbrTimerData& proto_esm_ebr_timer_data,
    esm_ebr_timer_data_t** state_esm_ebr_timer_data) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  *state_esm_ebr_timer_data =
      (esm_ebr_timer_data_t*) calloc(1, sizeof(esm_ebr_timer_data_t));
  if (*state_esm_ebr_timer_data) {
    (*state_esm_ebr_timer_data)->ue_id = proto_esm_ebr_timer_data.ue_id();
    (*state_esm_ebr_timer_data)->ebi   = proto_esm_ebr_timer_data.ebi();
    (*state_esm_ebr_timer_data)->count = proto_esm_ebr_timer_data.count();
    if (!proto_esm_ebr_timer_data.esm_msg().empty()) {
      (*state_esm_ebr_timer_data)->msg = bfromcstr_with_str_len(
          proto_esm_ebr_timer_data.esm_msg().c_str(),
          proto_esm_ebr_timer_data.esm_msg().length());
    }
  }
  OAILOG_FUNC_OUT(LOG_NAS_ESM);
}

void NasStateConverter::esm_proc_data_to_proto(
    const esm_proc_data_t* state_esm_proc_data,
    oai::EsmProcData* esm_proc_data_proto) {
  OAILOG_DEBUG(LOG_NAS_ESM, "Writing esm proc data to proto");
  esm_proc_data_proto->set_pti(state_esm_proc_data->pti);
  esm_proc_data_proto->set_request_type(state_esm_proc_data->request_type);
  if (state_esm_proc_data->apn) {
    BSTRING_TO_STRING(
        state_esm_proc_data->apn, esm_proc_data_proto->mutable_apn());
  }
  esm_proc_data_proto->set_pdn_cid(state_esm_proc_data->pdn_cid);
  esm_proc_data_proto->set_pdn_type(state_esm_proc_data->pdn_type);
  if (state_esm_proc_data->pdn_addr) {
    BSTRING_TO_STRING(
        state_esm_proc_data->pdn_addr, esm_proc_data_proto->mutable_pdn_addr());
  }
  bearer_qos_to_proto(
      state_esm_proc_data->bearer_qos,
      esm_proc_data_proto->mutable_bearer_qos());
  protocol_configuration_options_to_proto(
      state_esm_proc_data->pco, esm_proc_data_proto->mutable_pco());
}

void NasStateConverter::proto_to_esm_proc_data(
    const oai::EsmProcData& esm_proc_data_proto,
    esm_proc_data_t* state_esm_proc_data) {
  OAILOG_DEBUG(LOG_NAS_ESM, "Reading esm proc data from proto");
  state_esm_proc_data->pti          = esm_proc_data_proto.pti();
  state_esm_proc_data->request_type = esm_proc_data_proto.request_type();
  if (!esm_proc_data_proto.apn().empty()) {
    state_esm_proc_data->apn = bfromcstr_with_str_len(
        esm_proc_data_proto.apn().c_str(), esm_proc_data_proto.apn().length());
  }
  state_esm_proc_data->pdn_cid = esm_proc_data_proto.pdn_cid();
  state_esm_proc_data->pdn_type =
      (esm_proc_pdn_type_t) esm_proc_data_proto.pdn_type();
  if (!esm_proc_data_proto.pdn_addr().empty()) {
    state_esm_proc_data->pdn_addr = bfromcstr_with_str_len(
        esm_proc_data_proto.pdn_addr().c_str(),
        esm_proc_data_proto.pdn_addr().length());
  }
  proto_to_bearer_qos(
      esm_proc_data_proto.bearer_qos(), &state_esm_proc_data->bearer_qos);
  proto_to_protocol_configuration_options(
      esm_proc_data_proto.pco(), &state_esm_proc_data->pco);
}

void NasStateConverter::esm_context_to_proto(
    const esm_context_t* state_esm_context,
    oai::EsmContext* esm_context_proto) {
  esm_context_proto->set_n_active_ebrs(state_esm_context->n_active_ebrs);
  esm_context_proto->set_is_emergency(state_esm_context->is_emergency);
  esm_context_proto->set_pending_standalone(
      state_esm_context->pending_standalone);
  esm_context_proto->set_is_pdn_disconnect(
      state_esm_context->is_pdn_disconnect);
  if (state_esm_context->esm_proc_data) {
    esm_proc_data_to_proto(
        state_esm_context->esm_proc_data,
        esm_context_proto->mutable_esm_proc_data());
  }
}

void NasStateConverter::proto_to_esm_context(
    const oai::EsmContext& esm_context_proto,
    esm_context_t* state_esm_context) {
  OAILOG_DEBUG(LOG_NAS_ESM, "Reading esm context from proto");
  state_esm_context->n_active_ebrs     = esm_context_proto.n_active_ebrs();
  state_esm_context->is_emergency      = esm_context_proto.is_emergency();
  state_esm_context->is_pdn_disconnect = esm_context_proto.is_pdn_disconnect();
  state_esm_context->pending_standalone =
      esm_context_proto.pending_standalone();
  if (esm_context_proto.has_esm_proc_data()) {
    state_esm_context->esm_proc_data =
        (esm_proc_data_t*) calloc(1, sizeof(*state_esm_context->esm_proc_data));
    proto_to_esm_proc_data(
        esm_context_proto.esm_proc_data(), state_esm_context->esm_proc_data);
  }
}

void NasStateConverter::esm_ebr_context_to_proto(
    const esm_ebr_context_t& state_esm_ebr_context,
    oai::EsmEbrContext* esm_ebr_context_proto) {
  esm_ebr_context_proto->set_status(state_esm_ebr_context.status);
  esm_ebr_context_proto->set_gbr_dl(state_esm_ebr_context.gbr_dl);
  esm_ebr_context_proto->set_gbr_ul(state_esm_ebr_context.gbr_ul);
  esm_ebr_context_proto->set_mbr_dl(state_esm_ebr_context.mbr_dl);
  esm_ebr_context_proto->set_mbr_ul(state_esm_ebr_context.mbr_ul);
  /*TODO
   * SpgwStateConverter::traffic_flow_template_to_proto(state_esm_ebr_context.tft,
   * esm_ebr_context_proto->mutable_tft());*/
  if (state_esm_ebr_context.pco != nullptr) {
    protocol_configuration_options_to_proto(
        *state_esm_ebr_context.pco, esm_ebr_context_proto->mutable_pco());
  }
  if (state_esm_ebr_context.args != nullptr) {
    esm_ebr_timer_data_to_proto(
        *state_esm_ebr_context.args,
        esm_ebr_context_proto->mutable_esm_ebr_timer_data());
  }
}

void NasStateConverter::proto_to_esm_ebr_context(
    const oai::EsmEbrContext& esm_ebr_context_proto,
    esm_ebr_context_t* state_esm_ebr_context) {
  state_esm_ebr_context->status =
      (esm_ebr_state) esm_ebr_context_proto.status();
  state_esm_ebr_context->gbr_dl = esm_ebr_context_proto.gbr_dl();
  state_esm_ebr_context->gbr_ul = esm_ebr_context_proto.gbr_ul();
  state_esm_ebr_context->mbr_dl = esm_ebr_context_proto.mbr_dl();
  state_esm_ebr_context->mbr_ul = esm_ebr_context_proto.mbr_ul();
  /* SpgwStateConverter::proto_to_traffic_flow_template(esm_ebr_context_proto.tft(),
   * &state_esm_ebr_context->tft);*/
  if (esm_ebr_context_proto.has_pco()) {
    proto_to_protocol_configuration_options(
        esm_ebr_context_proto.pco(), state_esm_ebr_context->pco);
  }
  if (esm_ebr_context_proto.has_esm_ebr_timer_data()) {
    proto_to_esm_ebr_timer_data(
        esm_ebr_context_proto.esm_ebr_timer_data(),
        &state_esm_ebr_context->args);
  }
}

/*************************************************/
/*        EMM State <-> Proto                  */
/*************************************************/
void NasStateConverter::ue_network_capability_to_proto(
    const ue_network_capability_t* state_ue_network_capability,
    oai::UeNetworkCapability* ue_network_capability_proto) {
  ue_network_capability_proto->set_eea(state_ue_network_capability->eea);
  ue_network_capability_proto->set_eia(state_ue_network_capability->eia);
  ue_network_capability_proto->set_uea(state_ue_network_capability->uea);
  ue_network_capability_proto->set_ucs2(state_ue_network_capability->ucs2);
  ue_network_capability_proto->set_uia(state_ue_network_capability->uia);
  ue_network_capability_proto->set_prosedd(
      state_ue_network_capability->prosedd);
  ue_network_capability_proto->set_prose(state_ue_network_capability->prose);
  ue_network_capability_proto->set_h245ash(
      state_ue_network_capability->h245ash);
  ue_network_capability_proto->set_csfb(state_ue_network_capability->csfb);
  ue_network_capability_proto->set_lpp(state_ue_network_capability->lpp);
  ue_network_capability_proto->set_lcs(state_ue_network_capability->lcs);
  ue_network_capability_proto->set_srvcc(state_ue_network_capability->srvcc);
  ue_network_capability_proto->set_nf(state_ue_network_capability->nf);
  ue_network_capability_proto->set_epco(state_ue_network_capability->epco);
  ue_network_capability_proto->set_hccpciot(
      state_ue_network_capability->hccpciot);
  ue_network_capability_proto->set_erwfopdn(
      state_ue_network_capability->erwfopdn);
  ue_network_capability_proto->set_s1udata(
      state_ue_network_capability->s1udata);
  ue_network_capability_proto->set_upciot(state_ue_network_capability->upciot);
  ue_network_capability_proto->set_cpciot(state_ue_network_capability->cpciot);
  ue_network_capability_proto->set_proserelay(
      state_ue_network_capability->proserelay);
  ue_network_capability_proto->set_prosedc(
      state_ue_network_capability->prosedc);
  ue_network_capability_proto->set_bearer(state_ue_network_capability->bearer);
  ue_network_capability_proto->set_sgc(state_ue_network_capability->sgc);
  ue_network_capability_proto->set_n1mod(state_ue_network_capability->n1mod);
  ue_network_capability_proto->set_dcnr(state_ue_network_capability->dcnr);
  ue_network_capability_proto->set_cpbackoff(
      state_ue_network_capability->cpbackoff);
  ue_network_capability_proto->set_restrictec(
      state_ue_network_capability->restrictec);
  ue_network_capability_proto->set_v2xpc5(state_ue_network_capability->v2xpc5);
  ue_network_capability_proto->set_multipledrb(
      state_ue_network_capability->multipledrb);
  ue_network_capability_proto->set_umts_present(
      state_ue_network_capability->umts_present);
  ue_network_capability_proto->set_length(state_ue_network_capability->length);
}

void NasStateConverter::proto_to_ue_network_capability(
    const oai::UeNetworkCapability& ue_network_capability_proto,
    ue_network_capability_t* state_ue_network_capability) {
  state_ue_network_capability->eea     = ue_network_capability_proto.eea();
  state_ue_network_capability->eia     = ue_network_capability_proto.eia();
  state_ue_network_capability->uea     = ue_network_capability_proto.uea();
  state_ue_network_capability->ucs2    = ue_network_capability_proto.ucs2();
  state_ue_network_capability->uia     = ue_network_capability_proto.uia();
  state_ue_network_capability->prosedd = ue_network_capability_proto.prosedd();
  state_ue_network_capability->prose   = ue_network_capability_proto.prose();
  state_ue_network_capability->h245ash = ue_network_capability_proto.h245ash();
  state_ue_network_capability->csfb    = ue_network_capability_proto.csfb();
  state_ue_network_capability->lpp     = ue_network_capability_proto.lpp();
  state_ue_network_capability->lcs     = ue_network_capability_proto.lcs();
  state_ue_network_capability->srvcc   = ue_network_capability_proto.srvcc();
  state_ue_network_capability->nf      = ue_network_capability_proto.nf();
  state_ue_network_capability->epco    = ue_network_capability_proto.epco();
  state_ue_network_capability->hccpciot =
      ue_network_capability_proto.hccpciot();
  state_ue_network_capability->erwfopdn =
      ue_network_capability_proto.erwfopdn();
  state_ue_network_capability->s1udata = ue_network_capability_proto.s1udata();
  state_ue_network_capability->upciot  = ue_network_capability_proto.upciot();
  state_ue_network_capability->cpciot  = ue_network_capability_proto.cpciot();
  state_ue_network_capability->proserelay =
      ue_network_capability_proto.proserelay();
  state_ue_network_capability->prosedc = ue_network_capability_proto.prosedc();
  state_ue_network_capability->bearer  = ue_network_capability_proto.bearer();
  state_ue_network_capability->sgc     = ue_network_capability_proto.sgc();
  state_ue_network_capability->n1mod   = ue_network_capability_proto.n1mod();
  state_ue_network_capability->dcnr    = ue_network_capability_proto.dcnr();
  state_ue_network_capability->cpbackoff =
      ue_network_capability_proto.cpbackoff();
  state_ue_network_capability->restrictec =
      ue_network_capability_proto.restrictec();
  state_ue_network_capability->v2xpc5 = ue_network_capability_proto.v2xpc5();
  state_ue_network_capability->multipledrb =
      ue_network_capability_proto.multipledrb();
  state_ue_network_capability->umts_present =
      ue_network_capability_proto.umts_present();
  state_ue_network_capability->length = ue_network_capability_proto.length();
}

void NasStateConverter::classmark2_to_proto(
    const MobileStationClassmark2* state_MobileStationClassmark,
    oai::MobileStaClassmark2* mobile_station_classmark2_proto) {
  // TODO
}

void NasStateConverter::proto_to_classmark2(
    const oai::MobileStaClassmark2& mobile_sta_classmark2_proto,
    MobileStationClassmark2* state_MobileStationClassmar) {
  // TODO
}

void NasStateConverter::voice_preference_to_proto(
    const voice_domain_preference_and_ue_usage_setting_t*
        state_voice_domain_preference_and_ue_usage_setting,
    oai::VoicePreference* voice_preference_proto) {
  // TODO
}

void NasStateConverter::proto_to_voice_preference(
    const oai::VoicePreference& voice_preference_proto,
    voice_domain_preference_and_ue_usage_setting_t*
        state_voice_domain_preference_and_ue_usage_setting) {
  // TODO
}

void NasStateConverter::ue_additional_security_capability_to_proto(
    const ue_additional_security_capability_t*
        state_ue_additional_security_capability,
    oai::UeAdditionalSecurityCapability*
        ue_additional_security_capability_proto) {
  ue_additional_security_capability_proto->set_ea(
      state_ue_additional_security_capability->_5g_ea);
  ue_additional_security_capability_proto->set_ia(
      state_ue_additional_security_capability->_5g_ia);
}

void NasStateConverter::proto_to_ue_additional_security_capability(
    const oai::UeAdditionalSecurityCapability&
        ue_additional_security_capability_proto,
    ue_additional_security_capability_t*
        state_ue_additional_security_capability) {
  state_ue_additional_security_capability->_5g_ea =
      ue_additional_security_capability_proto.ea();
  state_ue_additional_security_capability->_5g_ia =
      ue_additional_security_capability_proto.ia();
}

void NasStateConverter::nas_message_decode_status_to_proto(
    const nas_message_decode_status_t* state_nas_message_decode_status,
    oai::NasMsgDecodeStatus* nas_msg_decode_status_proto) {
  nas_msg_decode_status_proto->set_integrity_protected_message(
      state_nas_message_decode_status->integrity_protected_message);
  nas_msg_decode_status_proto->set_ciphered_message(
      state_nas_message_decode_status->ciphered_message);
  nas_msg_decode_status_proto->set_mac_matched(
      state_nas_message_decode_status->mac_matched);
  nas_msg_decode_status_proto->set_security_context_available(
      state_nas_message_decode_status->security_context_available);
  nas_msg_decode_status_proto->set_emm_cause(
      state_nas_message_decode_status->emm_cause);
}

void NasStateConverter::proto_to_nas_message_decode_status(
    const oai::NasMsgDecodeStatus& nas_msg_decode_status_proto,
    nas_message_decode_status_t* state_nas_message_decode_status) {
  state_nas_message_decode_status->integrity_protected_message =
      nas_msg_decode_status_proto.integrity_protected_message();
  state_nas_message_decode_status->ciphered_message =
      nas_msg_decode_status_proto.ciphered_message();
  state_nas_message_decode_status->mac_matched =
      nas_msg_decode_status_proto.mac_matched();
  state_nas_message_decode_status->security_context_available =
      nas_msg_decode_status_proto.security_context_available();
  state_nas_message_decode_status->emm_cause =
      nas_msg_decode_status_proto.emm_cause();
}

void NasStateConverter::emm_attach_request_ies_to_proto(
    const emm_attach_request_ies_t* state_emm_attach_request_ies,
    oai::AttachRequestIes* attach_request_ies_proto) {
  attach_request_ies_proto->set_is_initial(
      state_emm_attach_request_ies->is_initial);
  attach_request_ies_proto->set_type(state_emm_attach_request_ies->type);
  attach_request_ies_proto->set_additional_update_type(
      state_emm_attach_request_ies->additional_update_type);
  attach_request_ies_proto->set_is_native_sc(
      state_emm_attach_request_ies->is_native_sc);
  attach_request_ies_proto->set_ksi(state_emm_attach_request_ies->ksi);
  attach_request_ies_proto->set_is_native_guti(
      state_emm_attach_request_ies->is_native_guti);
  if (state_emm_attach_request_ies->guti) {
    guti_to_proto(
        *state_emm_attach_request_ies->guti,
        attach_request_ies_proto->mutable_guti());
  }
  if (state_emm_attach_request_ies->imsi) {
    identity_tuple_to_proto<imsi_t>(
        state_emm_attach_request_ies->imsi,
        attach_request_ies_proto->mutable_imsi(), IMSI_BCD8_SIZE);
  }

  if (state_emm_attach_request_ies->imei) {
    identity_tuple_to_proto<imei_t>(
        state_emm_attach_request_ies->imei,
        attach_request_ies_proto->mutable_imei(), IMEI_BCD8_SIZE);
  }

  if (state_emm_attach_request_ies->last_visited_registered_tai) {
    tai_to_proto(
        state_emm_attach_request_ies->last_visited_registered_tai,
        attach_request_ies_proto->mutable_last_visited_tai());
  }

  if (state_emm_attach_request_ies->originating_tai) {
    tai_to_proto(
        state_emm_attach_request_ies->originating_tai,
        attach_request_ies_proto->mutable_origin_tai());
  }
  if (state_emm_attach_request_ies->originating_ecgi) {
    ecgi_to_proto(
        *state_emm_attach_request_ies->originating_ecgi,
        attach_request_ies_proto->mutable_origin_ecgi());
  }
  ue_network_capability_to_proto(
      &state_emm_attach_request_ies->ue_network_capability,
      attach_request_ies_proto->mutable_ue_nw_capability());
  if (state_emm_attach_request_ies->esm_msg) {
    BSTRING_TO_STRING(
        state_emm_attach_request_ies->esm_msg,
        attach_request_ies_proto->mutable_esm_msg());
  }
  nas_message_decode_status_to_proto(
      &state_emm_attach_request_ies->decode_status,
      attach_request_ies_proto->mutable_decode_status());
  classmark2_to_proto(
      state_emm_attach_request_ies->mob_st_clsMark2,
      attach_request_ies_proto->mutable_classmark2());
  voice_preference_to_proto(
      state_emm_attach_request_ies->voicedomainpreferenceandueusagesetting,
      attach_request_ies_proto->mutable_voice_preference());
  if (state_emm_attach_request_ies->ueadditionalsecuritycapability) {
    ue_additional_security_capability_to_proto(
        state_emm_attach_request_ies->ueadditionalsecuritycapability,
        attach_request_ies_proto->mutable_ue_additional_security_capability());
  }
}

void NasStateConverter::proto_to_emm_attach_request_ies(
    const oai::AttachRequestIes& attach_request_ies_proto,
    emm_attach_request_ies_t* state_emm_attach_request_ies) {
  state_emm_attach_request_ies->is_initial =
      attach_request_ies_proto.is_initial();
  state_emm_attach_request_ies->type =
      (emm_proc_attach_type_t) attach_request_ies_proto.type();
  state_emm_attach_request_ies->additional_update_type =
      (additional_update_type_t)
          attach_request_ies_proto.additional_update_type();
  state_emm_attach_request_ies->is_native_sc =
      attach_request_ies_proto.is_native_sc();
  state_emm_attach_request_ies->ksi = attach_request_ies_proto.ksi();
  state_emm_attach_request_ies->is_native_guti =
      attach_request_ies_proto.is_native_guti();
  if (attach_request_ies_proto.has_guti()) {
    state_emm_attach_request_ies->guti = (guti_t*) calloc(1, sizeof(guti_t));
    proto_to_guti(
        attach_request_ies_proto.guti(), state_emm_attach_request_ies->guti);
  }
  if (attach_request_ies_proto.has_imsi()) {
    state_emm_attach_request_ies->imsi = (imsi_t*) calloc(1, sizeof(imsi_t));
    proto_to_identity_tuple<imsi_t>(
        attach_request_ies_proto.imsi(), state_emm_attach_request_ies->imsi,
        IMSI_BCD8_SIZE);
  }
  if (attach_request_ies_proto.has_imei()) {
    state_emm_attach_request_ies->imei = (imei_t*) calloc(1, sizeof(imei_t));
    proto_to_identity_tuple<imei_t>(
        attach_request_ies_proto.imei(), state_emm_attach_request_ies->imei,
        IMEI_BCD8_SIZE);
  }
  if (attach_request_ies_proto.has_last_visited_tai()) {
    state_emm_attach_request_ies->last_visited_registered_tai =
        (tai_t*) calloc(1, sizeof(tai_t));
    proto_to_tai(
        attach_request_ies_proto.last_visited_tai(),
        state_emm_attach_request_ies->last_visited_registered_tai);
  }
  if (attach_request_ies_proto.has_origin_tai()) {
    state_emm_attach_request_ies->originating_tai =
        (tai_t*) calloc(1, sizeof(tai_t));
    proto_to_tai(
        attach_request_ies_proto.origin_tai(),
        state_emm_attach_request_ies->originating_tai);
  }
  if (attach_request_ies_proto.has_origin_ecgi()) {
    state_emm_attach_request_ies->originating_ecgi =
        (ecgi_t*) calloc(1, sizeof(ecgi_t));
    proto_to_ecgi(
        attach_request_ies_proto.origin_ecgi(),
        state_emm_attach_request_ies->originating_ecgi);
  }
  proto_to_ue_network_capability(
      attach_request_ies_proto.ue_nw_capability(),
      &state_emm_attach_request_ies->ue_network_capability);
  if (attach_request_ies_proto.esm_msg().length() > 0) {
    STRING_TO_BSTRING(
        attach_request_ies_proto.esm_msg(),
        state_emm_attach_request_ies->esm_msg);
  }
  proto_to_nas_message_decode_status(
      attach_request_ies_proto.decode_status(),
      &state_emm_attach_request_ies->decode_status);
  proto_to_classmark2(
      attach_request_ies_proto.classmark2(),
      state_emm_attach_request_ies->mob_st_clsMark2);
  proto_to_voice_preference(
      attach_request_ies_proto.voice_preference(),
      state_emm_attach_request_ies->voicedomainpreferenceandueusagesetting);
  if (attach_request_ies_proto.has_ue_additional_security_capability()) {
    state_emm_attach_request_ies->ueadditionalsecuritycapability =
        (ue_additional_security_capability_t*) calloc(
            1, sizeof(ue_additional_security_capability_t));
    proto_to_ue_additional_security_capability(
        attach_request_ies_proto.ue_additional_security_capability(),
        state_emm_attach_request_ies->ueadditionalsecuritycapability);
  }
}

void NasStateConverter::nas_attach_proc_to_proto(
    const nas_emm_attach_proc_t* state_nas_attach_proc,
    oai::AttachProc* attach_proc_proto) {
  attach_proc_proto->set_attach_accept_sent(
      state_nas_attach_proc->attach_accept_sent);
  attach_proc_proto->set_attach_reject_sent(
      state_nas_attach_proc->attach_reject_sent);
  attach_proc_proto->set_attach_complete_received(
      state_nas_attach_proc->attach_complete_received);
  guti_to_proto(state_nas_attach_proc->guti, attach_proc_proto->mutable_guti());

  char* esm_msg_buffer =
      bstr2cstr(state_nas_attach_proc->esm_msg_out, (char) '?');
  if (esm_msg_buffer) {
    attach_proc_proto->set_esm_msg_out(esm_msg_buffer);
    bcstrfree(esm_msg_buffer);
  } else {
    attach_proc_proto->set_esm_msg_out("");
  }

  if (state_nas_attach_proc->ies) {
    emm_attach_request_ies_to_proto(
        state_nas_attach_proc->ies, attach_proc_proto->mutable_ies());
  }
  attach_proc_proto->set_ue_id(state_nas_attach_proc->ue_id);
  attach_proc_proto->set_ksi(state_nas_attach_proc->ksi);
  attach_proc_proto->set_emm_cause(state_nas_attach_proc->emm_cause);
}

void NasStateConverter::proto_to_nas_emm_attach_proc(
    const oai::AttachProc& attach_proc_proto,
    nas_emm_attach_proc_t* state_nas_emm_attach_proc) {
  state_nas_emm_attach_proc->emm_spec_proc.emm_proc.base_proc.type =
      NAS_PROC_TYPE_EMM;
  state_nas_emm_attach_proc->emm_spec_proc.emm_proc.type =
      NAS_EMM_PROC_TYPE_CONN_MNGT;
  state_nas_emm_attach_proc->emm_spec_proc.type = EMM_SPEC_PROC_TYPE_ATTACH;
  state_nas_emm_attach_proc->attach_accept_sent =
      attach_proc_proto.attach_accept_sent();
  state_nas_emm_attach_proc->attach_reject_sent =
      attach_proc_proto.attach_reject_sent();
  state_nas_emm_attach_proc->attach_complete_received =
      attach_proc_proto.attach_complete_received();
  proto_to_guti(attach_proc_proto.guti(), &state_nas_emm_attach_proc->guti);
  if (attach_proc_proto.esm_msg_out().length() > 0) {
    state_nas_emm_attach_proc->esm_msg_out = bfromcstr_with_str_len(
        attach_proc_proto.esm_msg_out().c_str(),
        attach_proc_proto.esm_msg_out().length());
  }
  if (attach_proc_proto.has_ies()) {
    state_nas_emm_attach_proc->ies = (emm_attach_request_ies_t*) calloc(
        1, sizeof(*(state_nas_emm_attach_proc->ies)));
    proto_to_emm_attach_request_ies(
        attach_proc_proto.ies(), state_nas_emm_attach_proc->ies);
  }
  state_nas_emm_attach_proc->ue_id     = attach_proc_proto.ue_id();
  state_nas_emm_attach_proc->ksi       = attach_proc_proto.ksi();
  state_nas_emm_attach_proc->emm_cause = attach_proc_proto.emm_cause();
  state_nas_emm_attach_proc->T3450.sec = T3450_DEFAULT_VALUE;
  set_callbacks_for_attach_proc(state_nas_emm_attach_proc);
}

void NasStateConverter::emm_detach_request_ies_to_proto(
    const emm_detach_request_ies_t* state_emm_detach_request_ies,
    oai::DetachRequestIes* detach_request_ies_proto) {
  detach_request_ies_proto->set_type(state_emm_detach_request_ies->type);
  detach_request_ies_proto->set_switch_off(
      state_emm_detach_request_ies->switch_off);
  detach_request_ies_proto->set_is_native_sc(
      state_emm_detach_request_ies->is_native_sc);
  detach_request_ies_proto->set_ksi(state_emm_detach_request_ies->ksi);
  guti_to_proto(
      *state_emm_detach_request_ies->guti,
      detach_request_ies_proto->mutable_guti());
  identity_tuple_to_proto<imsi_t>(
      state_emm_detach_request_ies->imsi,
      detach_request_ies_proto->mutable_imsi(), IMSI_BCD8_SIZE);
  identity_tuple_to_proto<imei_t>(
      state_emm_detach_request_ies->imei,
      detach_request_ies_proto->mutable_imei(), IMEI_BCD8_SIZE);
  nas_message_decode_status_to_proto(
      &state_emm_detach_request_ies->decode_status,
      detach_request_ies_proto->mutable_decode_status());
}

void NasStateConverter::proto_to_emm_detach_request_ies(
    const oai::DetachRequestIes& detach_request_ies_proto,
    emm_detach_request_ies_t* state_emm_detach_request_ies) {
  state_emm_detach_request_ies->type =
      (emm_proc_detach_type_t) detach_request_ies_proto.type();
  state_emm_detach_request_ies->switch_off =
      detach_request_ies_proto.switch_off();
  state_emm_detach_request_ies->is_native_sc =
      detach_request_ies_proto.is_native_sc();
  state_emm_detach_request_ies->ksi = detach_request_ies_proto.ksi();
  proto_to_guti(
      detach_request_ies_proto.guti(), state_emm_detach_request_ies->guti);

  state_emm_detach_request_ies->imsi = (imsi_t*) calloc(1, sizeof(imsi_t));
  proto_to_identity_tuple<imsi_t>(
      detach_request_ies_proto.imsi(), state_emm_detach_request_ies->imsi,
      IMSI_BCD8_SIZE);
  state_emm_detach_request_ies->imei = (imei_t*) calloc(1, sizeof(imei_t));
  proto_to_identity_tuple<imei_t>(
      detach_request_ies_proto.imei(), state_emm_detach_request_ies->imei,
      IMEI_BCD8_SIZE);

  proto_to_nas_message_decode_status(
      detach_request_ies_proto.decode_status(),
      &state_emm_detach_request_ies->decode_status);
}

void NasStateConverter::emm_tau_request_ies_to_proto(
    const emm_tau_request_ies_t* state_emm_tau_request_ies,
    oai::TauRequestIes* tau_request_ies_proto) {
  // TODO
}

void NasStateConverter::proto_to_emm_tau_request_ies(
    const oai::TauRequestIes& tau_request_ies_proto,
    emm_tau_request_ies_t* state_emm_tau_request_ies) {
  // TODO
}

void NasStateConverter::nas_emm_tau_proc_to_proto(
    const nas_emm_tau_proc_t* state_nas_emm_tau_proc,
    oai::NasTauProc* nas_tau_proc_proto) {
  // TODO
}

void NasStateConverter::proto_to_nas_emm_tau_proc(
    const oai::NasTauProc& nas_tau_proc_proto,
    nas_emm_tau_proc_t* state_nas_emm_tau_proc) {
  // TODO
}

void NasStateConverter::nas_emm_auth_proc_to_proto(
    const nas_emm_auth_proc_t* state_nas_emm_auth_proc,
    oai::AuthProc* auth_proc_proto) {
  OAILOG_DEBUG(LOG_MME_APP, "Writing auth proc to proto");
  auth_proc_proto->Clear();
  auth_proc_proto->set_retransmission_count(
      state_nas_emm_auth_proc->retransmission_count);
  auth_proc_proto->set_sync_fail_count(
      state_nas_emm_auth_proc->sync_fail_count);
  auth_proc_proto->set_mac_fail_count(state_nas_emm_auth_proc->mac_fail_count);
  auth_proc_proto->set_ue_id(state_nas_emm_auth_proc->ue_id);
  auth_proc_proto->set_is_cause_is_attach(
      state_nas_emm_auth_proc->is_cause_is_attach);
  auth_proc_proto->set_ksi(state_nas_emm_auth_proc->ksi);
  auth_proc_proto->set_rand(
      (char*) state_nas_emm_auth_proc->rand, AUTH_RAND_SIZE);
  auth_proc_proto->set_autn(
      (char*) state_nas_emm_auth_proc->autn, AUTH_AUTN_SIZE);
  if (state_nas_emm_auth_proc->unchecked_imsi) {
    identity_tuple_to_proto<imsi_t>(
        state_nas_emm_auth_proc->unchecked_imsi,
        auth_proc_proto->mutable_unchecked_imsi(), IMSI_BCD8_SIZE);
  }
  auth_proc_proto->set_emm_cause(state_nas_emm_auth_proc->emm_cause);
}

void NasStateConverter::proto_to_nas_emm_auth_proc(
    const oai::AuthProc& auth_proc_proto,
    nas_emm_auth_proc_t* state_nas_emm_auth_proc) {
  OAILOG_DEBUG(LOG_MME_APP, "Reading auth proc from proto");
  state_nas_emm_auth_proc->emm_com_proc.emm_proc.base_proc.type =
      NAS_PROC_TYPE_EMM;
  state_nas_emm_auth_proc->emm_com_proc.emm_proc.type =
      NAS_EMM_PROC_TYPE_COMMON;
  state_nas_emm_auth_proc->emm_com_proc.type = EMM_COMM_PROC_AUTH;
  state_nas_emm_auth_proc->retransmission_count =
      auth_proc_proto.retransmission_count();
  state_nas_emm_auth_proc->sync_fail_count = auth_proc_proto.sync_fail_count();
  state_nas_emm_auth_proc->mac_fail_count  = auth_proc_proto.mac_fail_count();
  state_nas_emm_auth_proc->ue_id           = auth_proc_proto.ue_id();
  state_nas_emm_auth_proc->is_cause_is_attach =
      auth_proc_proto.is_cause_is_attach();
  state_nas_emm_auth_proc->ksi = auth_proc_proto.ksi();
  memcpy(
      state_nas_emm_auth_proc->rand, auth_proc_proto.rand().c_str(),
      AUTH_RAND_SIZE);
  memcpy(
      state_nas_emm_auth_proc->autn, auth_proc_proto.autn().c_str(),
      AUTH_AUTN_SIZE);

  if (auth_proc_proto.has_unchecked_imsi()) {
    state_nas_emm_auth_proc->unchecked_imsi =
        (imsi_t*) calloc(1, sizeof(imsi_t));
    proto_to_identity_tuple<imsi_t>(
        auth_proc_proto.unchecked_imsi(),
        state_nas_emm_auth_proc->unchecked_imsi, IMSI_BCD8_SIZE);
  }

  state_nas_emm_auth_proc->emm_cause = auth_proc_proto.emm_cause();
  state_nas_emm_auth_proc->T3460.sec = T3460_DEFAULT_VALUE;
  // update callback functions for auth proc
  set_callbacks_for_auth_proc(state_nas_emm_auth_proc);
  set_notif_callbacks_for_auth_proc(state_nas_emm_auth_proc);
}

void NasStateConverter::nas_emm_smc_proc_to_proto(
    const nas_emm_smc_proc_t* state_nas_emm_smc_proc,
    oai::SmcProc* smc_proc_proto) {
  OAILOG_DEBUG(LOG_MME_APP, "Writing smc proc to proto");
  smc_proc_proto->set_ue_id(state_nas_emm_smc_proc->ue_id);
  smc_proc_proto->set_retransmission_count(
      state_nas_emm_smc_proc->retransmission_count);
  smc_proc_proto->set_ksi(state_nas_emm_smc_proc->ksi);
  smc_proc_proto->set_eea(state_nas_emm_smc_proc->eea);
  smc_proc_proto->set_eia(state_nas_emm_smc_proc->eia);
  smc_proc_proto->set_ucs2(state_nas_emm_smc_proc->ucs2);
  smc_proc_proto->set_uea(state_nas_emm_smc_proc->uea);
  smc_proc_proto->set_uia(state_nas_emm_smc_proc->uia);
  smc_proc_proto->set_gea(state_nas_emm_smc_proc->gea);
  smc_proc_proto->set_umts_present(state_nas_emm_smc_proc->umts_present);
  smc_proc_proto->set_gprs_present(state_nas_emm_smc_proc->gprs_present);
  smc_proc_proto->set_selected_eea(state_nas_emm_smc_proc->selected_eea);
  smc_proc_proto->set_selected_eia(state_nas_emm_smc_proc->selected_eia);
  smc_proc_proto->set_saved_selected_eea(
      state_nas_emm_smc_proc->saved_selected_eea);
  smc_proc_proto->set_saved_selected_eia(
      state_nas_emm_smc_proc->saved_selected_eia);
  smc_proc_proto->set_saved_eksi(state_nas_emm_smc_proc->saved_eksi);
  smc_proc_proto->set_saved_overflow(state_nas_emm_smc_proc->saved_overflow);
  smc_proc_proto->set_saved_seq_num(state_nas_emm_smc_proc->saved_seq_num);
  smc_proc_proto->set_saved_sc_type(state_nas_emm_smc_proc->saved_sc_type);
  smc_proc_proto->set_notify_failure(state_nas_emm_smc_proc->notify_failure);
  smc_proc_proto->set_is_new(state_nas_emm_smc_proc->is_new);
  smc_proc_proto->set_imeisv_request(state_nas_emm_smc_proc->imeisv_request);
}

void NasStateConverter::proto_to_nas_emm_smc_proc(
    const oai::SmcProc& smc_proc_proto,
    nas_emm_smc_proc_t* state_nas_emm_smc_proc) {
  OAILOG_DEBUG(LOG_MME_APP, "Reading smc proc from proto");
  state_nas_emm_smc_proc->emm_com_proc.emm_proc.base_proc.type =
      NAS_PROC_TYPE_EMM;
  state_nas_emm_smc_proc->emm_com_proc.emm_proc.type = NAS_EMM_PROC_TYPE_COMMON;
  state_nas_emm_smc_proc->emm_com_proc.type          = EMM_COMM_PROC_SMC;
  state_nas_emm_smc_proc->ue_id                      = smc_proc_proto.ue_id();
  state_nas_emm_smc_proc->retransmission_count =
      smc_proc_proto.retransmission_count();
  state_nas_emm_smc_proc->ksi          = smc_proc_proto.ksi();
  state_nas_emm_smc_proc->eea          = smc_proc_proto.eea();
  state_nas_emm_smc_proc->eia          = smc_proc_proto.eia();
  state_nas_emm_smc_proc->ucs2         = smc_proc_proto.ucs2();
  state_nas_emm_smc_proc->uea          = smc_proc_proto.uea();
  state_nas_emm_smc_proc->uia          = smc_proc_proto.uia();
  state_nas_emm_smc_proc->gea          = smc_proc_proto.gea();
  state_nas_emm_smc_proc->umts_present = smc_proc_proto.umts_present();
  state_nas_emm_smc_proc->gprs_present = smc_proc_proto.gprs_present();
  state_nas_emm_smc_proc->selected_eea = smc_proc_proto.selected_eea();
  state_nas_emm_smc_proc->selected_eia = smc_proc_proto.selected_eia();
  state_nas_emm_smc_proc->saved_selected_eea =
      smc_proc_proto.saved_selected_eea();
  state_nas_emm_smc_proc->saved_selected_eia =
      smc_proc_proto.saved_selected_eia();
  state_nas_emm_smc_proc->saved_eksi     = smc_proc_proto.saved_eksi();
  state_nas_emm_smc_proc->saved_overflow = smc_proc_proto.saved_overflow();
  state_nas_emm_smc_proc->saved_seq_num  = smc_proc_proto.saved_seq_num();
  state_nas_emm_smc_proc->saved_sc_type =
      (emm_sc_type_t) smc_proc_proto.saved_sc_type();
  state_nas_emm_smc_proc->notify_failure = smc_proc_proto.notify_failure();
  state_nas_emm_smc_proc->is_new         = smc_proc_proto.is_new();
  state_nas_emm_smc_proc->imeisv_request = smc_proc_proto.imeisv_request();
  set_notif_callbacks_for_smc_proc(state_nas_emm_smc_proc);
  set_callbacks_for_smc_proc(state_nas_emm_smc_proc);
}

void NasStateConverter::nas_proc_mess_sign_to_proto(
    const nas_proc_mess_sign_t* state_nas_proc_mess_sign,
    oai::NasProcMessSign* nas_proc_mess_sign_proto) {
  nas_proc_mess_sign_proto->set_puid(state_nas_proc_mess_sign->puid);
  nas_proc_mess_sign_proto->set_digest(
      (void*) state_nas_proc_mess_sign->digest, NAS_MSG_DIGEST_SIZE);
  nas_proc_mess_sign_proto->set_digest_length(
      state_nas_proc_mess_sign->digest_length);
  nas_proc_mess_sign_proto->set_nas_msg_length(
      state_nas_proc_mess_sign->nas_msg_length);
}

void NasStateConverter::proto_to_nas_proc_mess_sign(
    const oai::NasProcMessSign& nas_proc_mess_sign_proto,
    nas_proc_mess_sign_t* state_nas_proc_mess_sign) {
  state_nas_proc_mess_sign->puid = nas_proc_mess_sign_proto.puid();
  memcpy(
      state_nas_proc_mess_sign->digest,
      nas_proc_mess_sign_proto.digest().c_str(), NAS_MSG_DIGEST_SIZE);
  state_nas_proc_mess_sign->digest_length =
      nas_proc_mess_sign_proto.digest_length();
  state_nas_proc_mess_sign->nas_msg_length =
      nas_proc_mess_sign_proto.nas_msg_length();
}

void NasStateConverter::nas_base_proc_to_proto(
    const nas_base_proc_t* base_proc_p, oai::NasBaseProc* base_proc_proto) {
  base_proc_proto->set_nas_puid(base_proc_p->nas_puid);
  base_proc_proto->set_type(base_proc_p->type);
}

void NasStateConverter::proto_to_nas_base_proc(
    const oai::NasBaseProc& nas_base_proc_proto,
    nas_base_proc_t* state_nas_base_proc) {
  state_nas_base_proc->nas_puid = nas_base_proc_proto.nas_puid();
  state_nas_base_proc->type = (nas_base_proc_type_t) nas_base_proc_proto.type();
  state_nas_base_proc->success_notif = NULL;
  state_nas_base_proc->failure_notif = NULL;
  state_nas_base_proc->abort         = NULL;
  state_nas_base_proc->fail_in       = NULL;
  state_nas_base_proc->fail_out      = NULL;
  state_nas_base_proc->time_out      = NULL;
}

void NasStateConverter::emm_proc_to_proto(
    const nas_emm_proc_t* emm_proc_p, oai::NasEmmProc* emm_proc_proto) {
  nas_base_proc_to_proto(
      &emm_proc_p->base_proc, emm_proc_proto->mutable_base_proc());
  emm_proc_proto->set_type(emm_proc_p->type);
  emm_proc_proto->set_previous_emm_fsm_state(
      emm_proc_p->previous_emm_fsm_state);
}

void NasStateConverter::proto_to_nas_emm_proc(
    const oai::NasEmmProc& nas_emm_proc_proto,
    nas_emm_proc_t* state_nas_emm_proc) {
  proto_to_nas_base_proc(
      nas_emm_proc_proto.base_proc(), &state_nas_emm_proc->base_proc);
  state_nas_emm_proc->type = (nas_emm_proc_type_t) nas_emm_proc_proto.type();
  state_nas_emm_proc->previous_emm_fsm_state =
      (emm_fsm_state_t) nas_emm_proc_proto.previous_emm_fsm_state();
  state_nas_emm_proc->delivered        = NULL;
  state_nas_emm_proc->not_delivered    = NULL;
  state_nas_emm_proc->not_delivered_ho = NULL;
}

void NasStateConverter::emm_specific_proc_to_proto(
    const nas_emm_specific_proc_t* state_emm_specific_proc,
    oai::NasEmmProcWithType* emm_proc_with_type) {
  OAILOG_DEBUG(LOG_MME_APP, "Writing specific procs to proto");
  emm_proc_to_proto(
      &state_emm_specific_proc->emm_proc,
      emm_proc_with_type->mutable_emm_proc());
  switch (state_emm_specific_proc->type) {
    case EMM_SPEC_PROC_TYPE_ATTACH: {
      OAILOG_DEBUG(LOG_MME_APP, "Writing attach proc to proto");
      nas_attach_proc_to_proto(
          (nas_emm_attach_proc_t*) state_emm_specific_proc,
          emm_proc_with_type->mutable_attach_proc());
      break;
    }
    case EMM_SPEC_PROC_TYPE_DETACH: {
      emm_detach_request_ies_to_proto(
          ((nas_emm_detach_proc_t*) state_emm_specific_proc)->ies,
          emm_proc_with_type->mutable_detach_proc());
      break;
    }
    default:
      break;
  }
}

// This function allocated memory for any specific procedure stored
void NasStateConverter::proto_to_emm_specific_proc(
    const oai::NasEmmProcWithType& proto_emm_proc_with_type,
    emm_procedures_t* state_emm_procedures) {
  OAILOG_DEBUG(LOG_MME_APP, "Reading specific procs from proto");
  // read attach or detach proc based on message type present
  switch (proto_emm_proc_with_type.MessageTypes_case()) {
    case oai::NasEmmProcWithType::kAttachProc: {
      OAILOG_DEBUG(LOG_MME_APP, "Reading attach proc from proto");
      state_emm_procedures
          ->emm_specific_proc = (nas_emm_specific_proc_t*) calloc(
          1,
          sizeof(
              nas_emm_attach_proc_t));  // NOLINT(clang-analyzer-unix.MallocSizeof)
      nas_emm_attach_proc_t* attach_proc =
          (nas_emm_attach_proc_t*) state_emm_procedures->emm_specific_proc;

      // read the emm proc content
      proto_to_nas_emm_proc(
          proto_emm_proc_with_type.emm_proc(),
          &attach_proc->emm_spec_proc.emm_proc);
      proto_to_nas_emm_attach_proc(
          proto_emm_proc_with_type.attach_proc(), attach_proc);
      break;
    }
    case oai::NasEmmProcWithType::kDetachProc: {
      state_emm_procedures
          ->emm_specific_proc = (nas_emm_specific_proc_t*) calloc(
          1,
          sizeof(
              nas_emm_detach_proc_t));  // NOLINT(clang-analyzer-unix.MallocSizeof)
      nas_emm_detach_proc_t* detach_proc =
          (nas_emm_detach_proc_t*) state_emm_procedures->emm_specific_proc;
      // read the emm proc content
      proto_to_nas_emm_proc(
          proto_emm_proc_with_type.emm_proc(),
          &detach_proc->emm_spec_proc.emm_proc);
      proto_to_emm_detach_request_ies(
          proto_emm_proc_with_type.detach_proc(), detach_proc->ies);
      break;
    }
    default:
      break;
  }
}

void NasStateConverter::emm_common_proc_to_proto(
    const emm_procedures_t* state_emm_procedures,
    oai::EmmProcedures* emm_procedures_proto) {
  OAILOG_DEBUG(LOG_MME_APP, "Writing common procs to proto");
  nas_emm_common_procedure_t* p1 =
      LIST_FIRST(&state_emm_procedures->emm_common_procs);
  while (p1) {
    oai::NasEmmProcWithType* nas_emm_proc_with_type_proto =
        emm_procedures_proto->add_emm_common_proc();
    emm_proc_to_proto(
        &p1->proc->emm_proc, nas_emm_proc_with_type_proto->mutable_emm_proc());
    switch (p1->proc->type) {
      case EMM_COMM_PROC_AUTH: {
        nas_emm_auth_proc_t* state_nas_emm_auth_proc =
            (nas_emm_auth_proc_t*) p1->proc;
        nas_emm_auth_proc_to_proto(
            state_nas_emm_auth_proc,
            nas_emm_proc_with_type_proto->mutable_auth_proc());
      } break;
      case EMM_COMM_PROC_SMC: {
        nas_emm_smc_proc_t* state_nas_emm_smc_proc =
            (nas_emm_smc_proc_t*) p1->proc;
        nas_emm_smc_proc_to_proto(
            state_nas_emm_smc_proc,
            nas_emm_proc_with_type_proto->mutable_smc_proc());
      } break;
      default:
        break;
    }
    p1 = LIST_NEXT(p1, entries);
  }
}

void NasStateConverter::insert_proc_into_emm_common_procs(
    emm_procedures_t* state_emm_procedures,
    nas_emm_common_proc_t* emm_com_proc) {
  nas_emm_common_procedure_t* wrapper =
      (nas_emm_common_procedure_t*) calloc(1, sizeof(*wrapper));
  if (!wrapper) return;

  wrapper->proc = emm_com_proc;
  LIST_INSERT_HEAD(&state_emm_procedures->emm_common_procs, wrapper, entries);
  OAILOG_DEBUG(LOG_NAS_EMM, "New COMMON PROC added from state\n");
}

void NasStateConverter::proto_to_emm_common_proc(
    const oai::EmmProcedures& emm_procedures_proto,
    emm_context_t* state_emm_context) {
  OAILOG_DEBUG(LOG_MME_APP, "Reading common procs from proto");
  auto proto_common_procs = emm_procedures_proto.emm_common_proc();
  for (auto ptr = proto_common_procs.begin(); ptr < proto_common_procs.end();
       ptr++) {
    switch (ptr->MessageTypes_case()) {
      case oai::NasEmmProcWithType::kAuthProc: {
        OAILOG_DEBUG(LOG_NAS_EMM, "Inserting AUTH PROC from state\n");
        nas_emm_auth_proc_t* state_auth_proc =
            (nas_emm_auth_proc_t*) calloc(1, sizeof(*state_auth_proc));
        proto_to_nas_emm_proc(
            ptr->emm_proc(), &state_auth_proc->emm_com_proc.emm_proc);
        proto_to_nas_emm_auth_proc(ptr->auth_proc(), state_auth_proc);
        nas_emm_attach_proc_t* attach_proc =
            get_nas_specific_procedure_attach(state_emm_context);
        ((nas_base_proc_t*) state_auth_proc)->parent =
            (nas_base_proc_t*) &attach_proc->emm_spec_proc;
        insert_proc_into_emm_common_procs(
            state_emm_context->emm_procedures, &state_auth_proc->emm_com_proc);
        break;
      }
      case oai::NasEmmProcWithType::kSmcProc: {
        OAILOG_DEBUG(LOG_NAS_EMM, "Inserting SMC PROC from state\n");
        nas_emm_smc_proc_t* state_smc_proc =
            (nas_emm_smc_proc_t*) calloc(1, sizeof(*state_smc_proc));
        proto_to_nas_emm_proc(
            ptr->emm_proc(), &state_smc_proc->emm_com_proc.emm_proc);
        proto_to_nas_emm_smc_proc(ptr->smc_proc(), state_smc_proc);
        nas_emm_attach_proc_t* attach_proc =
            get_nas_specific_procedure_attach(state_emm_context);
        ((nas_base_proc_t*) state_smc_proc)->parent =
            (nas_base_proc_t*) &attach_proc->emm_spec_proc;
        insert_proc_into_emm_common_procs(
            state_emm_context->emm_procedures, &state_smc_proc->emm_com_proc);
        break;
      }
      default:
        break;
    }
  }
}

void NasStateConverter::eutran_vectors_to_proto(
    eutran_vector_t** state_eutran_vector_array, uint8_t num_vectors,
    oai::AuthInfoProc* auth_info_proc_proto) {
  oai::AuthVector* eutran_vector_proto = nullptr;
  OAILOG_DEBUG(LOG_NAS_EMM, "Writing %d eutran vectors", num_vectors);
  for (int i = 0; i < num_vectors; i++) {
    eutran_vector_proto = auth_info_proc_proto->add_vector();
    memcpy(
        eutran_vector_proto->mutable_kasme(),
        state_eutran_vector_array[i]->kasme, KASME_LENGTH_OCTETS);
    memcpy(
        eutran_vector_proto->mutable_rand(), state_eutran_vector_array[i]->rand,
        RAND_LENGTH_OCTETS);
    memcpy(
        eutran_vector_proto->mutable_autn(), state_eutran_vector_array[i]->autn,
        AUTN_LENGTH_OCTETS);
    memcpy(
        eutran_vector_proto->mutable_xres(),
        state_eutran_vector_array[i]->xres.data,
        state_eutran_vector_array[i]->xres.size);
  }
}

void NasStateConverter::proto_to_eutran_vectors(
    const oai::AuthInfoProc& auth_info_proc_proto,
    nas_auth_info_proc_t* state_nas_auth_info_proc) {
  auto proto_vectors = auth_info_proc_proto.vector();
  int i              = 0;
  for (auto ptr = proto_vectors.begin(); ptr < proto_vectors.end(); ptr++) {
    eutran_vector_t* this_vector =
        (eutran_vector_t*) calloc(1, sizeof(eutran_vector_t));
    memcpy(this_vector->kasme, ptr->kasme().c_str(), AUTH_KASME_SIZE);
    memcpy(this_vector->rand, ptr->rand().c_str(), AUTH_RAND_SIZE);
    memcpy(this_vector->autn, ptr->autn().c_str(), AUTH_AUTN_SIZE);
    this_vector->xres.size = ptr->xres().length();
    memcpy(this_vector->xres.data, ptr->xres().c_str(), this_vector->xres.size);
    state_nas_auth_info_proc->vector[i] = this_vector;
    i++;
  }
  state_nas_auth_info_proc->nb_vectors = i;
  OAILOG_DEBUG(
      LOG_NAS_EMM, "Read %d eutran vectors",
      state_nas_auth_info_proc->nb_vectors);
}

void NasStateConverter::nas_auth_info_proc_to_proto(
    nas_auth_info_proc_t* state_nas_auth_info_proc,
    oai::AuthInfoProc* auth_info_proc_proto) {
  auth_info_proc_proto->set_request_sent(
      state_nas_auth_info_proc->request_sent);
  eutran_vectors_to_proto(
      state_nas_auth_info_proc->vector, state_nas_auth_info_proc->nb_vectors,
      auth_info_proc_proto);
  auth_info_proc_proto->set_nas_cause(state_nas_auth_info_proc->nas_cause);
  auth_info_proc_proto->set_ue_id(state_nas_auth_info_proc->ue_id);
  auth_info_proc_proto->set_resync(state_nas_auth_info_proc->resync);
}

void NasStateConverter::proto_to_nas_auth_info_proc(
    const oai::AuthInfoProc& auth_info_proc_proto,
    nas_auth_info_proc_t* state_nas_auth_info_proc) {
  state_nas_auth_info_proc->request_sent = auth_info_proc_proto.request_sent();
  proto_to_eutran_vectors(auth_info_proc_proto, state_nas_auth_info_proc);
  state_nas_auth_info_proc->nas_cause = auth_info_proc_proto.nas_cause();
  state_nas_auth_info_proc->ue_id     = auth_info_proc_proto.ue_id();
  state_nas_auth_info_proc->resync    = auth_info_proc_proto.resync();
  // update success_notif and failure_notif
  set_callbacks_for_auth_info_proc(state_nas_auth_info_proc);
}

void NasStateConverter::nas_cn_procs_to_proto(
    const emm_procedures_t* state_emm_procedures,
    oai::EmmProcedures* emm_procedures_proto) {
  OAILOG_DEBUG(LOG_NAS_EMM, "Writing cn procs to proto");
  nas_cn_procedure_t* p1 = LIST_FIRST(&state_emm_procedures->cn_procs);
  while (p1) {
    oai::NasCnProc* nas_cn_proc_proto = emm_procedures_proto->add_cn_proc();
    nas_base_proc_to_proto(
        &p1->proc->base_proc, nas_cn_proc_proto->mutable_base_proc());
    switch (p1->proc->type) {
      case CN_PROC_AUTH_INFO: {
        OAILOG_DEBUG(LOG_NAS_EMM, "Writing auth_info_proc to proto");
        nas_auth_info_proc_t* state_auth_info_proc =
            (nas_auth_info_proc_t*) p1->proc;
        nas_auth_info_proc_to_proto(
            state_auth_info_proc, nas_cn_proc_proto->mutable_auth_info_proc());
      } break;
      default:
        OAILOG_DEBUG(
            LOG_NAS,
            "EMM_CN: Unknown procedure type, cannot convert"
            "to proto");
        break;
    }
    p1 = LIST_NEXT(p1, entries);
  }
}

void NasStateConverter::insert_proc_into_cn_procs(
    emm_procedures_t* state_emm_procedures, nas_cn_proc_t* cn_proc) {
  nas_cn_procedure_t* wrapper =
      (nas_cn_procedure_t*) calloc(1, sizeof(*wrapper));
  if (!wrapper) return;
  wrapper->proc = cn_proc;
  LIST_INSERT_HEAD(&state_emm_procedures->cn_procs, wrapper, entries);
  OAILOG_DEBUG(LOG_NAS_EMM, "New EMM_COMM_PROC_SMC\n");
}

void NasStateConverter::proto_to_nas_cn_proc(
    const oai::EmmProcedures& emm_procedures_proto,
    emm_procedures_t* state_emm_procedures) {
  OAILOG_DEBUG(LOG_NAS_EMM, "Reading cn procs from proto");
  auto proto_cn_procs = emm_procedures_proto.cn_proc();
  for (auto ptr = proto_cn_procs.begin(); ptr < proto_cn_procs.end(); ptr++) {
    switch (ptr->MessageTypes_case()) {
      case oai::NasCnProc::kAuthInfoProc: {
        nas_auth_info_proc_t* state_auth_info_proc =
            (nas_auth_info_proc_t*) calloc(1, sizeof(*state_auth_info_proc));
        OAILOG_DEBUG(LOG_NAS_EMM, "Inserting AUTH INFO PROC from state\n");
        proto_to_nas_base_proc(
            ptr->base_proc(), &state_auth_info_proc->cn_proc.base_proc);
        proto_to_nas_auth_info_proc(
            ptr->auth_info_proc(), state_auth_info_proc);
        state_auth_info_proc->cn_proc.type = CN_PROC_AUTH_INFO;
        insert_proc_into_cn_procs(
            state_emm_procedures, &state_auth_info_proc->cn_proc);
      }
      default:
        break;
    }
  }
}

void NasStateConverter::mess_sign_array_to_proto(
    const emm_procedures_t* state_emm_procedures,
    oai::EmmProcedures* emm_procedures_proto) {
  for (int i = 0; i < MAX_NAS_PROC_MESS_SIGN; i++) {
    oai::NasProcMessSign* nas_proc_mess_sign_proto =
        emm_procedures_proto->add_nas_proc_mess_sign();
    nas_proc_mess_sign_to_proto(
        &state_emm_procedures->nas_proc_mess_sign[i], nas_proc_mess_sign_proto);
  }
}

void NasStateConverter::proto_to_mess_sign_array(
    const oai::EmmProcedures& emm_procedures_proto,
    emm_procedures_t* state_emm_procedures) {
  int i                      = 0;
  auto proto_mess_sign_array = emm_procedures_proto.nas_proc_mess_sign();
  for (auto ptr = proto_mess_sign_array.begin();
       ptr < proto_mess_sign_array.end(); ptr++) {
    proto_to_nas_proc_mess_sign(
        *ptr, &state_emm_procedures->nas_proc_mess_sign[i]);
    i++;
  }
}

void NasStateConverter::emm_procedures_to_proto(
    const emm_procedures_t* state_emm_procedures,
    oai::EmmProcedures* emm_procedures_proto) {
  if (state_emm_procedures->emm_specific_proc) {
    emm_specific_proc_to_proto(
        state_emm_procedures->emm_specific_proc,
        emm_procedures_proto->mutable_emm_specific_proc());
  }
  emm_common_proc_to_proto(state_emm_procedures, emm_procedures_proto);

  // cn_procs
  nas_cn_procs_to_proto(state_emm_procedures, emm_procedures_proto);
  oai::NasEmmProcWithType* emm_proc_with_type =
      emm_procedures_proto->mutable_emm_con_mngt_proc();

  if (state_emm_procedures->emm_con_mngt_proc) {
    emm_proc_to_proto(
        &state_emm_procedures->emm_con_mngt_proc->emm_proc,
        emm_proc_with_type->mutable_emm_proc());
  }

  emm_procedures_proto->set_nas_proc_mess_sign_next_location(
      state_emm_procedures->nas_proc_mess_sign_next_location);

  mess_sign_array_to_proto(state_emm_procedures, emm_procedures_proto);
}

void NasStateConverter::proto_to_emm_procedures(
    const oai::EmmProcedures& emm_procedures_proto,
    emm_context_t* state_emm_context) {
  state_emm_context->emm_procedures =
      (emm_procedures_t*) calloc(1, sizeof(*state_emm_context->emm_procedures));
  emm_procedures_t* state_emm_procedures = state_emm_context->emm_procedures;
  LIST_INIT(&state_emm_procedures->emm_common_procs);
  proto_to_emm_specific_proc(
      emm_procedures_proto.emm_specific_proc(), state_emm_procedures);
  LIST_INIT(&state_emm_procedures->emm_common_procs);
  proto_to_emm_common_proc(emm_procedures_proto, state_emm_context);
  LIST_INIT(&state_emm_procedures->cn_procs);
  proto_to_nas_cn_proc(emm_procedures_proto, state_emm_procedures);
  if (emm_procedures_proto.has_emm_con_mngt_proc()) {
    state_emm_procedures->emm_con_mngt_proc =
        (nas_emm_con_mngt_proc_t*) calloc(1, sizeof(nas_emm_con_mngt_proc_t));
    proto_to_nas_emm_proc(
        emm_procedures_proto.emm_con_mngt_proc().emm_proc(),
        &state_emm_procedures->emm_con_mngt_proc->emm_proc);
  }
  state_emm_procedures->nas_proc_mess_sign_next_location =
      emm_procedures_proto.nas_proc_mess_sign_next_location();
  proto_to_mess_sign_array(emm_procedures_proto, state_emm_procedures);
}

void NasStateConverter::auth_vectors_to_proto(
    const auth_vector_t* state_auth_vector_array, int num_vectors,
    oai::EmmContext* emm_context_proto) {
  oai::AuthVector* auth_vector_proto = nullptr;
  for (int i = 0; i < num_vectors; i++) {
    auth_vector_proto = emm_context_proto->add_vector();
    auth_vector_proto->set_kasme(
        (const void*) state_auth_vector_array[i].kasme, AUTH_KASME_SIZE);
    auth_vector_proto->set_rand(
        (const void*) state_auth_vector_array[i].rand, AUTH_RAND_SIZE);
    auth_vector_proto->set_autn(
        (const void*) state_auth_vector_array[i].autn, AUTH_AUTN_SIZE);
    auth_vector_proto->set_xres(
        (const void*) state_auth_vector_array[i].xres,
        state_auth_vector_array[i].xres_size);
  }
}

int NasStateConverter::proto_to_auth_vectors(
    const oai::EmmContext& emm_context_proto,
    auth_vector_t* state_auth_vector) {
  auto proto_vectors = emm_context_proto.vector();
  int i              = 0;
  for (auto ptr = proto_vectors.begin(); ptr < proto_vectors.end(); ptr++) {
    memcpy(state_auth_vector[i].kasme, ptr->kasme().c_str(), AUTH_KASME_SIZE);
    memcpy(state_auth_vector[i].rand, ptr->rand().c_str(), AUTH_RAND_SIZE);
    memcpy(state_auth_vector[i].autn, ptr->autn().c_str(), AUTH_AUTN_SIZE);
    memcpy(
        state_auth_vector[i].xres, ptr->xres().c_str(), ptr->xres().length());
    i++;
  }
  return i;
}

void NasStateConverter::emm_security_context_to_proto(
    const emm_security_context_t* state_emm_security_context,
    oai::EmmSecurityContext* emm_security_context_proto) {
  emm_security_context_proto->set_sc_type(state_emm_security_context->sc_type);
  emm_security_context_proto->set_eksi(state_emm_security_context->eksi);
  emm_security_context_proto->set_vector_index(
      state_emm_security_context->vector_index);
  emm_security_context_proto->set_knas_enc(
      state_emm_security_context->knas_enc, AUTH_KNAS_ENC_SIZE);
  emm_security_context_proto->set_knas_int(
      state_emm_security_context->knas_int, AUTH_KNAS_INT_SIZE);

  // Count values
  oai::EmmSecurityContext_Count* dl_count_proto =
      emm_security_context_proto->mutable_dl_count();
  dl_count_proto->set_spare(state_emm_security_context->dl_count.spare);
  dl_count_proto->set_overflow(state_emm_security_context->dl_count.overflow);
  dl_count_proto->set_seq_num(state_emm_security_context->dl_count.seq_num);
  oai::EmmSecurityContext_Count* ul_count_proto =
      emm_security_context_proto->mutable_ul_count();
  ul_count_proto->set_spare(state_emm_security_context->ul_count.spare);
  ul_count_proto->set_overflow(state_emm_security_context->ul_count.overflow);
  ul_count_proto->set_seq_num(state_emm_security_context->ul_count.seq_num);
  oai::EmmSecurityContext_Count* kenb_ul_count_proto =
      emm_security_context_proto->mutable_kenb_ul_count();
  kenb_ul_count_proto->set_spare(
      state_emm_security_context->kenb_ul_count.spare);
  kenb_ul_count_proto->set_overflow(
      state_emm_security_context->kenb_ul_count.overflow);
  kenb_ul_count_proto->set_seq_num(
      state_emm_security_context->kenb_ul_count.seq_num);

  // TODO convert capability to proto

  // Security algorithm
  oai::EmmSecurityContext_SelectedAlgorithms* selected_algorithms_proto =
      emm_security_context_proto->mutable_selected_algos();
  selected_algorithms_proto->set_encryption(
      state_emm_security_context->selected_algorithms.encryption);
  selected_algorithms_proto->set_integrity(
      state_emm_security_context->selected_algorithms.integrity);
  emm_security_context_proto->set_activated(
      state_emm_security_context->activated);
  emm_security_context_proto->set_direction_encode(
      state_emm_security_context->direction_encode);
  emm_security_context_proto->set_direction_decode(
      state_emm_security_context->direction_decode);
  emm_security_context_proto->set_next_hop(
      state_emm_security_context->next_hop, AUTH_NEXT_HOP_SIZE);
  emm_security_context_proto->set_next_hop_chaining_count(
      state_emm_security_context->next_hop_chaining_count);
}

void NasStateConverter::proto_to_emm_security_context(
    const oai::EmmSecurityContext& emm_security_context_proto,
    emm_security_context_t* state_emm_security_context) {
  state_emm_security_context->sc_type =
      (emm_sc_type_t) emm_security_context_proto.sc_type();
  state_emm_security_context->eksi = emm_security_context_proto.eksi();
  state_emm_security_context->vector_index =
      emm_security_context_proto.vector_index();
  memcpy(
      state_emm_security_context->knas_enc,
      emm_security_context_proto.knas_enc().c_str(), AUTH_KNAS_ENC_SIZE);
  memcpy(
      state_emm_security_context->knas_int,
      emm_security_context_proto.knas_int().c_str(), AUTH_KNAS_INT_SIZE);

  // Count values
  const oai::EmmSecurityContext_Count& dl_count_proto =
      emm_security_context_proto.dl_count();
  state_emm_security_context->dl_count.spare    = dl_count_proto.spare();
  state_emm_security_context->dl_count.overflow = dl_count_proto.overflow();
  state_emm_security_context->dl_count.seq_num  = dl_count_proto.seq_num();
  const oai::EmmSecurityContext_Count& ul_count_proto =
      emm_security_context_proto.ul_count();
  state_emm_security_context->ul_count.spare    = ul_count_proto.spare();
  state_emm_security_context->ul_count.overflow = ul_count_proto.overflow();
  state_emm_security_context->ul_count.seq_num  = ul_count_proto.seq_num();
  const oai::EmmSecurityContext_Count& kenb_ul_count_proto =
      emm_security_context_proto.kenb_ul_count();
  state_emm_security_context->kenb_ul_count.spare = kenb_ul_count_proto.spare();
  state_emm_security_context->kenb_ul_count.overflow =
      kenb_ul_count_proto.overflow();
  state_emm_security_context->kenb_ul_count.seq_num =
      kenb_ul_count_proto.seq_num();

  // TODO read capability from proto

  // Security algorithm
  const oai::EmmSecurityContext_SelectedAlgorithms& selected_algorithms_proto =
      emm_security_context_proto.selected_algos();
  state_emm_security_context->selected_algorithms.encryption =
      selected_algorithms_proto.encryption();
  state_emm_security_context->selected_algorithms.integrity =
      selected_algorithms_proto.integrity();
  state_emm_security_context->activated =
      emm_security_context_proto.activated();
  state_emm_security_context->direction_encode =
      emm_security_context_proto.direction_encode();
  state_emm_security_context->direction_decode =
      emm_security_context_proto.direction_decode();
  memcpy(
      state_emm_security_context->next_hop,
      emm_security_context_proto.next_hop().c_str(), AUTH_NEXT_HOP_SIZE);
  state_emm_security_context->next_hop_chaining_count =
      emm_security_context_proto.next_hop_chaining_count();
}

void NasStateConverter::nw_detach_data_to_proto(
    nw_detach_data_t* detach_timer_arg,
    oai::NwDetachData* detach_timer_arg_proto) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  detach_timer_arg_proto->set_ue_id(detach_timer_arg->ue_id);
  detach_timer_arg_proto->set_retransmission_count(
      detach_timer_arg->retransmission_count);
  detach_timer_arg_proto->set_detach_type(detach_timer_arg->detach_type);
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

void NasStateConverter::proto_to_nw_detach_data(
    const oai::NwDetachData& detach_timer_arg_proto,
    nw_detach_data_t** detach_timer_arg) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  *detach_timer_arg = (nw_detach_data_t*) calloc(1, sizeof(nw_detach_data_t));
  (*detach_timer_arg)->ue_id = detach_timer_arg_proto.ue_id();
  (*detach_timer_arg)->retransmission_count =
      detach_timer_arg_proto.retransmission_count();
  (*detach_timer_arg)->detach_type = detach_timer_arg_proto.detach_type();
  OAILOG_FUNC_OUT(LOG_MME_APP);
}

void NasStateConverter::emm_context_to_proto(
    const emm_context_t* state_emm_context,
    oai::EmmContext* emm_context_proto) {
  emm_context_proto->set_imsi64(state_emm_context->_imsi64);
  identity_tuple_to_proto<imsi_t>(
      &state_emm_context->_imsi, emm_context_proto->mutable_imsi(),
      IMSI_BCD8_SIZE);
  emm_context_proto->set_saved_imsi64(state_emm_context->saved_imsi64);
  identity_tuple_to_proto<imei_t>(
      &state_emm_context->_imei, emm_context_proto->mutable_imei(),
      IMEI_BCD8_SIZE);
  identity_tuple_to_proto<imeisv_t>(
      &state_emm_context->_imeisv, emm_context_proto->mutable_imeisv(),
      IMEISV_BCD8_SIZE);
  emm_context_proto->set_emm_cause(state_emm_context->emm_cause);
  emm_context_proto->set_emm_fsm_state(state_emm_context->_emm_fsm_state);
  emm_context_proto->set_attach_type(state_emm_context->attach_type);

  if (state_emm_context->emm_procedures) {
    emm_procedures_to_proto(
        state_emm_context->emm_procedures,
        emm_context_proto->mutable_emm_procedures());
  }
  emm_context_proto->set_common_proc_mask(state_emm_context->common_proc_mask);
  esm_context_to_proto(
      &state_emm_context->esm_ctx, emm_context_proto->mutable_esm_ctx());
  emm_context_proto->set_member_present_mask(
      state_emm_context->member_present_mask);
  emm_context_proto->set_member_valid_mask(
      state_emm_context->member_valid_mask);
  int num_auth_vectors =
      MAX_EPS_AUTH_VECTORS - state_emm_context->remaining_vectors;
  auth_vectors_to_proto(
      state_emm_context->_vector, num_auth_vectors, emm_context_proto);
  emm_security_context_to_proto(
      &state_emm_context->_security, emm_context_proto->mutable_security());
  emm_security_context_to_proto(
      &state_emm_context->_non_current_security,
      emm_context_proto->mutable__non_current_security());
  emm_context_proto->set_is_dynamic(state_emm_context->is_dynamic);
  emm_context_proto->set_is_attached(state_emm_context->is_attached);
  emm_context_proto->set_is_initial_identity_imsi(
      state_emm_context->is_initial_identity_imsi);
  emm_context_proto->set_is_guti_based_attach(
      state_emm_context->is_guti_based_attach);
  emm_context_proto->set_is_imsi_only_detach(
      state_emm_context->is_imsi_only_detach);
  emm_context_proto->set_is_emergency(state_emm_context->is_emergency);
  emm_context_proto->set_additional_update_type(
      state_emm_context->additional_update_type);
  emm_context_proto->set_tau_updt_type(state_emm_context->tau_updt_type);
  emm_context_proto->set_num_attach_request(
      state_emm_context->num_attach_request);
  guti_to_proto(state_emm_context->_guti, emm_context_proto->mutable_guti());
  guti_to_proto(
      state_emm_context->_old_guti, emm_context_proto->mutable_old_guti());
  tai_list_to_proto(
      &state_emm_context->_tai_list, emm_context_proto->mutable_tai_list());
  tai_to_proto(
      &state_emm_context->_lvr_tai, emm_context_proto->mutable_lvr_tai());
  tai_to_proto(
      &state_emm_context->originating_tai,
      emm_context_proto->mutable_originating_tai());
  emm_context_proto->set_ksi(state_emm_context->ksi);
  ue_network_capability_to_proto(
      &state_emm_context->_ue_network_capability,
      emm_context_proto->mutable_ue_network_capability());
  ue_additional_security_capability_to_proto(
      &state_emm_context->ue_additional_security_capability,
      emm_context_proto->mutable_ue_additional_security_capability());
  if (state_emm_context->t3422_arg) {
    nw_detach_data_to_proto(
        (nw_detach_data_t*) state_emm_context->t3422_arg,
        emm_context_proto->mutable_nw_detach_data());
  }
  if (state_emm_context->new_attach_info) {
    new_attach_info_to_proto(
        state_emm_context->new_attach_info,
        emm_context_proto->mutable_new_attach_info());
  }
}

void NasStateConverter::proto_to_emm_context(
    const oai::EmmContext& emm_context_proto,
    emm_context_t* state_emm_context) {
  state_emm_context->_imsi64 = emm_context_proto.imsi64();
  proto_to_identity_tuple<imsi_t>(
      emm_context_proto.imsi(), &state_emm_context->_imsi, IMSI_BCD8_SIZE);
  state_emm_context->saved_imsi64 = emm_context_proto.saved_imsi64();

  proto_to_identity_tuple<imei_t>(
      emm_context_proto.imei(), &state_emm_context->_imei, IMEI_BCD8_SIZE);
  proto_to_identity_tuple<imeisv_t>(
      emm_context_proto.imeisv(), &state_emm_context->_imeisv,
      IMEISV_BCD8_SIZE);

  state_emm_context->emm_cause = emm_context_proto.emm_cause();
  state_emm_context->_emm_fsm_state =
      (emm_fsm_state_t) emm_context_proto.emm_fsm_state();
  state_emm_context->attach_type = emm_context_proto.attach_type();
  if (emm_context_proto.has_emm_procedures()) {
    proto_to_emm_procedures(
        emm_context_proto.emm_procedures(), state_emm_context);
  }
  nas_emm_auth_proc_t* auth_proc =
      get_nas_common_procedure_authentication(state_emm_context);
  if (auth_proc) {
    OAILOG_DEBUG(
        LOG_MME_APP, "Found non-null auth proc with RAND value " RAND_FORMAT,
        RAND_DISPLAY(auth_proc->rand));
  }
  nas_auth_info_proc_t* auth_info_proc =
      get_nas_cn_procedure_auth_info(state_emm_context);
  if (auth_info_proc) {
    OAILOG_DEBUG(LOG_MME_APP, "Found non-null auth info proc");
  }
  state_emm_context->common_proc_mask = emm_context_proto.common_proc_mask();
  proto_to_esm_context(
      emm_context_proto.esm_ctx(), &state_emm_context->esm_ctx);
  state_emm_context->member_present_mask =
      emm_context_proto.member_present_mask();
  state_emm_context->member_valid_mask = emm_context_proto.member_valid_mask();
  int num_vectors =
      proto_to_auth_vectors(emm_context_proto, state_emm_context->_vector);
  state_emm_context->remaining_vectors = MAX_EPS_AUTH_VECTORS - num_vectors;
  proto_to_emm_security_context(
      emm_context_proto.security(), &state_emm_context->_security);
  proto_to_emm_security_context(
      emm_context_proto._non_current_security(),
      &state_emm_context->_non_current_security);
  state_emm_context->is_dynamic  = emm_context_proto.is_dynamic();
  state_emm_context->is_attached = emm_context_proto.is_attached();
  state_emm_context->is_initial_identity_imsi =
      emm_context_proto.is_initial_identity_imsi();
  state_emm_context->is_guti_based_attach =
      emm_context_proto.is_guti_based_attach();
  state_emm_context->is_imsi_only_detach =
      emm_context_proto.is_imsi_only_detach();
  state_emm_context->is_emergency = emm_context_proto.is_emergency();
  state_emm_context->additional_update_type =
      (additional_update_type_t) emm_context_proto.additional_update_type();
  state_emm_context->tau_updt_type = emm_context_proto.tau_updt_type();
  state_emm_context->num_attach_request =
      emm_context_proto.num_attach_request();
  proto_to_guti(emm_context_proto.guti(), &state_emm_context->_guti);
  proto_to_guti(emm_context_proto.old_guti(), &state_emm_context->_old_guti);
  proto_to_tai_list(
      emm_context_proto.tai_list(), &state_emm_context->_tai_list);
  proto_to_tai(emm_context_proto.lvr_tai(), &state_emm_context->_lvr_tai);
  proto_to_tai(
      emm_context_proto.originating_tai(), &state_emm_context->originating_tai);
  state_emm_context->ksi = emm_context_proto.ksi();
  proto_to_ue_network_capability(
      emm_context_proto.ue_network_capability(),
      &state_emm_context->_ue_network_capability);
  proto_to_ue_additional_security_capability(
      emm_context_proto.ue_additional_security_capability(),
      &state_emm_context->ue_additional_security_capability);

  state_emm_context->T3422.id  = NAS_TIMER_INACTIVE_ID;
  state_emm_context->T3422.sec = T3422_DEFAULT_VALUE;
  if (emm_context_proto.has_nw_detach_data()) {
    proto_to_nw_detach_data(
        emm_context_proto.nw_detach_data(),
        (nw_detach_data_t**) &state_emm_context->t3422_arg);
  }
  if (emm_context_proto.has_new_attach_info()) {
    state_emm_context->new_attach_info =
        (new_attach_info_t*) calloc(1, sizeof(new_attach_info_t));
    proto_to_new_attach_info(
        emm_context_proto.new_attach_info(),
        state_emm_context->new_attach_info);
  }
}

void NasStateConverter::new_attach_info_to_proto(
    const new_attach_info_t* state_new_attach_info,
    oai::NewAttachInfo* proto_new_attach_info) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  proto_new_attach_info->set_mme_ue_s1ap_id(
      state_new_attach_info->mme_ue_s1ap_id);
  proto_new_attach_info->set_is_mm_ctx_new(
      state_new_attach_info->is_mm_ctx_new);

  if (state_new_attach_info->ies) {
    emm_attach_request_ies_to_proto(
        state_new_attach_info->ies, proto_new_attach_info->mutable_ies());
  }
  OAILOG_FUNC_OUT(LOG_NAS_EMM);
}

void NasStateConverter::proto_to_new_attach_info(
    const oai::NewAttachInfo& proto_new_attach_info,
    new_attach_info_t* state_new_attach_info) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  state_new_attach_info->mme_ue_s1ap_id =
      proto_new_attach_info.mme_ue_s1ap_id();
  state_new_attach_info->is_mm_ctx_new = proto_new_attach_info.is_mm_ctx_new();
  if (proto_new_attach_info.has_ies()) {
    state_new_attach_info->ies = (emm_attach_request_ies_t*) calloc(
        1, sizeof(*(state_new_attach_info->ies)));
    proto_to_emm_attach_request_ies(
        proto_new_attach_info.ies(), state_new_attach_info->ies);
  }
  OAILOG_FUNC_OUT(LOG_NAS_EMM);
}

}  // namespace lte
}  // namespace magma
