/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the Apache License, Version 2.0  (the "License"); you may not use this file
 * except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
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
}

#include "nas_state_converter.h"
#include "spgw_state_converter.h"

namespace magma {
namespace lte {

#define PLMN_BYTES 6

NasStateConverter::NasStateConverter() = default;
NasStateConverter::~NasStateConverter() = default;

/*************************************************/
/*        Common Types <-> Proto                 */
/*************************************************/

void NasStateConverter::proto_to_guti(
  const Guti& guti_proto,
  guti_t* state_guti)
{
  strncpy(
    (char*) &state_guti->gummei.plmn,
    (guti_proto.plmn()).c_str(),
    sizeof(plmn_t));
  state_guti->gummei.mme_gid = guti_proto.mme_gid();
  state_guti->gummei.mme_code = guti_proto.mme_code();
  state_guti->m_tmsi = (tmsi_t) guti_proto.m_tmsi();
}

void NasStateConverter::proto_to_ecgi(
  const Ecgi& ecgi_proto,
  ecgi_t* state_ecgi)
{
  strcpy((char*) &state_ecgi->plmn, ecgi_proto.plmn().c_str());
  state_ecgi->cell_identity.enb_id = ecgi_proto.enb_id();
  state_ecgi->cell_identity.cell_id = ecgi_proto.cell_id();
  state_ecgi->cell_identity.empty = ecgi_proto.empty();
}

void NasStateConverter::partial_tai_list_to_proto(
  const partial_tai_list_t* state_partial_tai_list,
  PartialTaiList* partial_tai_list_proto)
{
  partial_tai_list_proto->set_type_of_list(state_partial_tai_list->typeoflist);
  partial_tai_list_proto->set_number_of_elements(
    state_partial_tai_list->numberofelements);
  // TODO
}

void NasStateConverter::tai_list_to_proto(
  const tai_list_t* state_tai_list,
  TaiList* tai_list_proto)
{
  tai_list_proto->set_numberoflists(state_tai_list->numberoflists);
  //TODO
}

void NasStateConverter::proto_to_tai_list(
  const TaiList& tai_list_proto,
  tai_list_t* state_tai_list)
{
  state_tai_list->numberoflists = tai_list_proto.numberoflists();
  // TODO
}

void NasStateConverter::tai_to_proto(const tai_t* state_tai, Tai* tai_proto)
{
  char plmn_array[PLMN_BYTES];
  plmn_array[0] = (char) state_tai->mcc_digit2;
  plmn_array[1] = (char) state_tai->mcc_digit1;
  plmn_array[2] = (char) state_tai->mnc_digit3;
  plmn_array[3] = (char) state_tai->mcc_digit3;
  plmn_array[4] = (char) state_tai->mnc_digit2;
  plmn_array[5] = (char) state_tai->mnc_digit1;

  tai_proto->set_mcc_mnc(plmn_array);
  tai_proto->set_tac(state_tai->tac);
}

void NasStateConverter::proto_to_tai(const Tai& tai_proto, tai_t* state_tai)
{
  state_tai->mcc_digit2 = tai_proto.mcc_mnc()[0];
  state_tai->mcc_digit1 = tai_proto.mcc_mnc()[1];
  state_tai->mnc_digit3 = tai_proto.mcc_mnc()[2];
  state_tai->mcc_digit3 = tai_proto.mcc_mnc()[3];
  state_tai->mnc_digit2 = tai_proto.mcc_mnc()[4];
  state_tai->mnc_digit1 = tai_proto.mcc_mnc()[5];
  state_tai->tac = tai_proto.tac();
}

/*************************************************/
/*        ESM State <-> Proto                  */
/*************************************************/
void NasStateConverter::ambr_to_proto(const ambr_t& state_ambr, Ambr* ambr_proto)
{
  ambr_proto->set_br_ul(state_ambr.br_ul);
  ambr_proto->set_br_dl(state_ambr.br_dl);
}

void NasStateConverter::proto_to_ambr(
  const Ambr& ambr_proto,
  ambr_t* state_ambr)
{
  state_ambr->br_ul = ambr_proto.br_ul();
  state_ambr->br_dl = ambr_proto.br_dl();
}

void NasStateConverter::bearer_qos_to_proto(
  const bearer_qos_t& state_bearer_qos,
  BearerQos* bearer_qos_proto)
{
  bearer_qos_proto->set_pci(state_bearer_qos.pci);
  bearer_qos_proto->set_pl(state_bearer_qos.pl);
  bearer_qos_proto->set_pvi(state_bearer_qos.pvi);
  bearer_qos_proto->set_qci(state_bearer_qos.qci);
  ambr_to_proto(state_bearer_qos.gbr, bearer_qos_proto->mutable_gbr());
  ambr_to_proto(state_bearer_qos.mbr, bearer_qos_proto->mutable_mbr());
}

void NasStateConverter::proto_to_bearer_qos(
  const BearerQos& bearer_qos_proto,
  bearer_qos_t* state_bearer_qos)
{
  state_bearer_qos->pci = bearer_qos_proto.pci();
  state_bearer_qos->pl = bearer_qos_proto.pl();
  state_bearer_qos->pvi = bearer_qos_proto.pvi();
  state_bearer_qos->qci = bearer_qos_proto.qci();
  proto_to_ambr(bearer_qos_proto.gbr(), &state_bearer_qos->gbr);
  proto_to_ambr(bearer_qos_proto.mbr(), &state_bearer_qos->mbr);
}

void NasStateConverter::pco_protocol_or_container_id_to_proto(
  const protocol_configuration_options_t& state_protocol_configuration_options,
  ProtocolConfigurationOptions* protocol_configuration_options_proto)
{
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
  const ProtocolConfigurationOptions& protocol_configuration_options_proto,
  protocol_configuration_options_t* state_protocol_configuration_options)
{
  auto proto_pco_ids =
    protocol_configuration_options_proto.proto_or_container_id();
  int i = 0;
  for (auto ptr = proto_pco_ids.begin(); ptr < proto_pco_ids.end(); ptr++) {
    pco_protocol_or_container_id_t state_pco_protocol_or_container_id =
      state_protocol_configuration_options->protocol_or_container_ids[i];
    state_pco_protocol_or_container_id.id = ptr->id();
    state_pco_protocol_or_container_id.length = ptr->length();
    state_pco_protocol_or_container_id.contents =
      bfromcstr(ptr->contents().c_str());
    i++;
  }
}

void NasStateConverter::protocol_configuration_options_to_proto(
  const protocol_configuration_options_t& state_protocol_configuration_options,
  ProtocolConfigurationOptions* protocol_configuration_options_proto)
{
  protocol_configuration_options_proto->set_ext(
    state_protocol_configuration_options.ext);
  protocol_configuration_options_proto->set_spare(
    state_protocol_configuration_options.spare);
  protocol_configuration_options_proto->set_config_protocol(
    state_protocol_configuration_options.configuration_protocol);
  protocol_configuration_options_proto->set_num_protocol_or_container_id(
    state_protocol_configuration_options.num_protocol_or_container_id);
  pco_protocol_or_container_id_to_proto(
    state_protocol_configuration_options, protocol_configuration_options_proto);
}

void NasStateConverter::proto_to_protocol_configuration_options(
  const ProtocolConfigurationOptions& protocol_configuration_options_proto,
  protocol_configuration_options_t* state_protocol_configuration_options)
{
  state_protocol_configuration_options->ext =
    protocol_configuration_options_proto.ext();
  state_protocol_configuration_options->spare =
    protocol_configuration_options_proto.spare();
  state_protocol_configuration_options->configuration_protocol =
    protocol_configuration_options_proto.config_protocol();
  state_protocol_configuration_options->num_protocol_or_container_id =
    protocol_configuration_options_proto.num_protocol_or_container_id();
  proto_to_pco_protocol_or_container_id(
    protocol_configuration_options_proto, state_protocol_configuration_options);
}

void NasStateConverter::esm_proc_data_to_proto(
  const esm_proc_data_t* state_esm_proc_data,
  EsmProcData* esm_proc_data_proto)
{
  esm_proc_data_proto->set_pti(state_esm_proc_data->pti);
  esm_proc_data_proto->set_request_type(state_esm_proc_data->request_type);
  if (state_esm_proc_data->apn) {
    esm_proc_data_proto->set_apn(
      (char*) state_esm_proc_data->apn->data, state_esm_proc_data->apn->slen);
  }
  esm_proc_data_proto->set_pdn_cid(state_esm_proc_data->pdn_cid);
  esm_proc_data_proto->set_pdn_type(state_esm_proc_data->pdn_type);
  if (state_esm_proc_data->pdn_addr) {
    esm_proc_data_proto->set_pdn_addr(
      (char*) state_esm_proc_data->pdn_addr->data,
      state_esm_proc_data->pdn_addr->slen);
  }
  bearer_qos_to_proto(
    state_esm_proc_data->bearer_qos, esm_proc_data_proto->mutable_bearer_qos());
  /* TODO: SpgwStateConverter::traffic_flow_template_to_proto(
    &state_esm_proc_data->tft, esm_proc_data_proto->mutable_tft());*/
  protocol_configuration_options_to_proto(
    state_esm_proc_data->pco, esm_proc_data_proto->mutable_pco());
}

void NasStateConverter::proto_to_esm_proc_data(
  const EsmProcData& esm_proc_data_proto,
  esm_proc_data_t* state_esm_proc_data)
{
  state_esm_proc_data->pti = esm_proc_data_proto.pti();
  state_esm_proc_data->request_type = esm_proc_data_proto.request_type();
  state_esm_proc_data->apn = bfromcstr(esm_proc_data_proto.apn().c_str());
  state_esm_proc_data->pdn_cid = esm_proc_data_proto.pdn_cid();
  state_esm_proc_data->pdn_type =
    (esm_proc_pdn_type_t) esm_proc_data_proto.pdn_type();
  state_esm_proc_data->pdn_addr =
    bfromcstr(esm_proc_data_proto.pdn_addr().c_str());
  proto_to_bearer_qos(
    esm_proc_data_proto.bearer_qos(), &state_esm_proc_data->bearer_qos);
  /* TODO: SpgwStateConverter::proto_to_traffic_flow_template(
    esm_proc_data_proto.tft(), &state_esm_proc_data->tft);*/
  proto_to_protocol_configuration_options(
    esm_proc_data_proto.pco(), &state_esm_proc_data->pco);
}

void NasStateConverter::esm_context_to_proto(
  const esm_context_t* state_esm_context,
  EsmContext* esm_context_proto)
{
  esm_context_proto->set_n_active_ebrs(state_esm_context->n_active_ebrs);
  esm_context_proto->set_n_active_pdns(state_esm_context->n_active_pdns);
  esm_context_proto->set_n_pdns(state_esm_context->n_pdns);
  esm_context_proto->set_is_emergency(state_esm_context->is_emergency);
  if (state_esm_context->esm_proc_data) {
    esm_proc_data_to_proto(
      state_esm_context->esm_proc_data,
      esm_context_proto->mutable_esm_proc_data());
  }
  nas_timer_to_proto(
    state_esm_context->T3489, esm_context_proto->mutable_t3489());
}

void NasStateConverter::proto_to_esm_context(
  const EsmContext& esm_context_proto,
  esm_context_t* state_esm_context)
{
  state_esm_context->n_active_ebrs = esm_context_proto.n_active_ebrs();
  state_esm_context->n_active_pdns = esm_context_proto.n_active_pdns();
  state_esm_context->n_pdns = esm_context_proto.n_pdns();
  state_esm_context->is_emergency = esm_context_proto.is_emergency();
  if (esm_context_proto.has_esm_proc_data()) {
    state_esm_context->esm_proc_data =
      (esm_proc_data_t*) calloc(1, sizeof(state_esm_context->esm_proc_data));
    proto_to_esm_proc_data(
      esm_context_proto.esm_proc_data(), state_esm_context->esm_proc_data);
  }
  proto_to_nas_timer(esm_context_proto.t3489(), &state_esm_context->T3489);
}

void NasStateConverter::esm_ebr_context_to_proto(
  const esm_ebr_context_t& state_esm_ebr_context,
  EsmEbrContext* esm_ebr_context_proto)
{
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
  nas_timer_to_proto(
    state_esm_ebr_context.timer, esm_ebr_context_proto->mutable_timer());
}

void NasStateConverter::proto_to_esm_ebr_context(
  const EsmEbrContext& esm_ebr_context_proto,
  esm_ebr_context_t* state_esm_ebr_context)
{
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
  proto_to_nas_timer(
    esm_ebr_context_proto.timer(), &state_esm_ebr_context->timer);
}

/*************************************************/
/*        EMM State <-> Proto                  */
/*************************************************/
void NasStateConverter::nas_timer_to_proto(
  const nas_timer_t& state_nas_timer,
  Timer* timer_proto)
{
  timer_proto->set_id(state_nas_timer.id);
  timer_proto->set_sec(state_nas_timer.sec);
}

void NasStateConverter::proto_to_nas_timer(
  const Timer& timer_proto,
  nas_timer_t* state_nas_timer)
{
  state_nas_timer->id = timer_proto.id();
  state_nas_timer->sec = timer_proto.sec();
}

void NasStateConverter::ue_network_capability_to_proto(
  const ue_network_capability_t* state_ue_network_capability,
  UeNetworkCapability* ue_network_capability_proto)
{
  ue_network_capability_proto->set_eea(state_ue_network_capability->eea);
  ue_network_capability_proto->set_eia(state_ue_network_capability->eia);
  ue_network_capability_proto->set_uea(state_ue_network_capability->uea);
  ue_network_capability_proto->set_ucs2(state_ue_network_capability->ucs2);
  ue_network_capability_proto->set_uia(state_ue_network_capability->uia);
  ue_network_capability_proto->set_spare(state_ue_network_capability->spare);
  ue_network_capability_proto->set_csfb(state_ue_network_capability->csfb);
  ue_network_capability_proto->set_lpp(state_ue_network_capability->lpp);
  ue_network_capability_proto->set_lcs(state_ue_network_capability->lcs);
  ue_network_capability_proto->set_srvcc(state_ue_network_capability->srvcc);
  ue_network_capability_proto->set_nf(state_ue_network_capability->nf);
  ue_network_capability_proto->set_umts_present(
    state_ue_network_capability->umts_present);
  ue_network_capability_proto->set_misc_present(
    state_ue_network_capability->misc_present);
}

void NasStateConverter::proto_to_ue_network_capability(
  const UeNetworkCapability& ue_network_capability_proto,
  ue_network_capability_t* state_ue_network_capability)
{
  state_ue_network_capability->eea = ue_network_capability_proto.eea();
  state_ue_network_capability->eia = ue_network_capability_proto.eia();
  state_ue_network_capability->uea = ue_network_capability_proto.uea();
  state_ue_network_capability->ucs2 = ue_network_capability_proto.ucs2();
  state_ue_network_capability->uia = ue_network_capability_proto.uia();
  state_ue_network_capability->spare = ue_network_capability_proto.spare();
  state_ue_network_capability->csfb = ue_network_capability_proto.csfb();
  state_ue_network_capability->lpp = ue_network_capability_proto.lpp();
  state_ue_network_capability->lcs = ue_network_capability_proto.lcs();
  state_ue_network_capability->srvcc = ue_network_capability_proto.srvcc();
  state_ue_network_capability->nf = ue_network_capability_proto.nf();
  state_ue_network_capability->umts_present =
    ue_network_capability_proto.umts_present();
  state_ue_network_capability->misc_present =
    ue_network_capability_proto.misc_present();
}

void NasStateConverter::classmark2_to_proto(
  const MobileStationClassmark2* state_MobileStationClassmark,
  MobileStaClassmark2* mobile_station_classmark2_proto)
{
  // TODO
}

void NasStateConverter::proto_to_classmark2(
  const MobileStaClassmark2& mobile_sta_classmark2_proto,
  MobileStationClassmark2* state_MobileStationClassmar)
{
  // TODO
}

void NasStateConverter::voice_preference_to_proto(
  const voice_domain_preference_and_ue_usage_setting_t*
    state_voice_domain_preference_and_ue_usage_setting,
  VoicePreference* voice_preference_proto)
{
  // TODO
}

void NasStateConverter::proto_to_voice_preference(
  const VoicePreference& voice_preference_proto,
  voice_domain_preference_and_ue_usage_setting_t*
    state_voice_domain_preference_and_ue_usage_setting)
{
  // TODO
}

void NasStateConverter::nas_message_decode_status_to_proto(
  const nas_message_decode_status_t* state_nas_message_decode_status,
  NasMsgDecodeStatus* nas_msg_decode_status_proto)
{
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
  const NasMsgDecodeStatus& nas_msg_decode_status_proto,
  nas_message_decode_status_t* state_nas_message_decode_status)
{
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
  AttachRequestIes* attach_request_ies_proto)
{
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
      attach_request_ies_proto->mutable_imsi(),
      IMSI_BCD8_SIZE);
  }

  if (state_emm_attach_request_ies->imei) {
    identity_tuple_to_proto<imei_t>(
      state_emm_attach_request_ies->imei,
      attach_request_ies_proto->mutable_imei(),
      IMEI_BCD8_SIZE);
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
}

void NasStateConverter::proto_to_emm_attach_request_ies(
  const AttachRequestIes& attach_request_ies_proto,
  emm_attach_request_ies_t* state_emm_attach_request_ies)
{
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
      attach_request_ies_proto.imsi(),
      state_emm_attach_request_ies->imsi,
      IMSI_BCD8_SIZE);
  }
  if (attach_request_ies_proto.has_imei()) {
    state_emm_attach_request_ies->imei = (imei_t*) calloc(1, sizeof(imei_t));
    proto_to_identity_tuple<imei_t>(
      attach_request_ies_proto.imei(),
      state_emm_attach_request_ies->imei,
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
}

void NasStateConverter::nas_attach_proc_to_proto(
  const nas_emm_attach_proc_t* state_nas_attach_proc,
  AttachProc* attach_proc_proto)
{
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
  nas_timer_to_proto(
    state_nas_attach_proc->T3450, attach_proc_proto->mutable_t3450());
}

void NasStateConverter::proto_to_nas_emm_attach_proc(
  const AttachProc& attach_proc_proto,
  nas_emm_attach_proc_t* state_nas_emm_attach_proc)
{
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
  state_nas_emm_attach_proc->esm_msg_out =
    bfromcstr(attach_proc_proto.esm_msg_out().c_str());
  if (attach_proc_proto.has_ies()) {
    state_nas_emm_attach_proc->ies = (emm_attach_request_ies_t*) calloc(
      1, sizeof(*(state_nas_emm_attach_proc->ies)));
    proto_to_emm_attach_request_ies(
      attach_proc_proto.ies(), state_nas_emm_attach_proc->ies);
  }
  state_nas_emm_attach_proc->ue_id = attach_proc_proto.ue_id();
  state_nas_emm_attach_proc->ksi = attach_proc_proto.ksi();
  state_nas_emm_attach_proc->emm_cause = attach_proc_proto.emm_cause();
  proto_to_nas_timer(
    attach_proc_proto.t3450(), &state_nas_emm_attach_proc->T3450);
  set_callbacks_for_attach_proc(state_nas_emm_attach_proc);
}

void NasStateConverter::emm_detach_request_ies_to_proto(
  const emm_detach_request_ies_t* state_emm_detach_request_ies,
  DetachRequestIes* detach_request_ies_proto)
{
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
    detach_request_ies_proto->mutable_imsi(),
    IMSI_BCD8_SIZE);
  identity_tuple_to_proto<imei_t>(
    state_emm_detach_request_ies->imei,
    detach_request_ies_proto->mutable_imei(),
    IMEI_BCD8_SIZE);
  nas_message_decode_status_to_proto(
    &state_emm_detach_request_ies->decode_status,
    detach_request_ies_proto->mutable_decode_status());
}

void NasStateConverter::proto_to_emm_detach_request_ies(
  const DetachRequestIes& detach_request_ies_proto,
  emm_detach_request_ies_t* state_emm_detach_request_ies)
{
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
    detach_request_ies_proto.imsi(),
    state_emm_detach_request_ies->imsi,
    IMSI_BCD8_SIZE);
  state_emm_detach_request_ies->imei = (imei_t*) calloc(1, sizeof(imei_t));
  proto_to_identity_tuple<imei_t>(
    detach_request_ies_proto.imei(),
    state_emm_detach_request_ies->imei,
    IMEI_BCD8_SIZE);

  proto_to_nas_message_decode_status(
    detach_request_ies_proto.decode_status(),
    &state_emm_detach_request_ies->decode_status);
}

void NasStateConverter::emm_tau_request_ies_to_proto(
  const emm_tau_request_ies_t* state_emm_tau_request_ies,
  TauRequestIes* tau_request_ies_proto)
{
  // TODO
}

void NasStateConverter::proto_to_emm_tau_request_ies(
  const TauRequestIes& tau_request_ies_proto,
  emm_tau_request_ies_t* state_emm_tau_request_ies)
{
  // TODO
}

void NasStateConverter::nas_emm_tau_proc_to_proto(
  const nas_emm_tau_proc_t* state_nas_emm_tau_proc,
  NasTauProc* nas_tau_proc_proto)
{
  // TODO
}

void NasStateConverter::proto_to_nas_emm_tau_proc(
  const NasTauProc& nas_tau_proc_proto,
  nas_emm_tau_proc_t* state_nas_emm_tau_proc)
{
  // TODO
}

void NasStateConverter::nas_emm_auth_proc_to_proto(
  const nas_emm_auth_proc_t* state_nas_emm_auth_proc,
  AuthProc* auth_proc_proto)
{
  OAILOG_INFO(LOG_MME_APP, "Writing auth proc to proto");
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
      auth_proc_proto->mutable_unchecked_imsi(),
      IMSI_BCD8_SIZE);
  }
  auth_proc_proto->set_emm_cause(state_nas_emm_auth_proc->emm_cause);
  nas_timer_to_proto(
    state_nas_emm_auth_proc->T3460, auth_proc_proto->mutable_t3460());
}

void NasStateConverter::proto_to_nas_emm_auth_proc(
  const AuthProc& auth_proc_proto,
  nas_emm_auth_proc_t* state_nas_emm_auth_proc)
{
  OAILOG_INFO(LOG_MME_APP, "Reading auth proc from proto");
  state_nas_emm_auth_proc->emm_com_proc.emm_proc.base_proc.type =
    NAS_PROC_TYPE_EMM;
  state_nas_emm_auth_proc->emm_com_proc.emm_proc.type =
    NAS_EMM_PROC_TYPE_COMMON;
  state_nas_emm_auth_proc->emm_com_proc.type = EMM_COMM_PROC_AUTH;
  state_nas_emm_auth_proc->retransmission_count =
    auth_proc_proto.retransmission_count();
  state_nas_emm_auth_proc->sync_fail_count = auth_proc_proto.sync_fail_count();
  state_nas_emm_auth_proc->mac_fail_count = auth_proc_proto.mac_fail_count();
  state_nas_emm_auth_proc->ue_id = auth_proc_proto.ue_id();
  state_nas_emm_auth_proc->is_cause_is_attach =
    auth_proc_proto.is_cause_is_attach();
  state_nas_emm_auth_proc->ksi = auth_proc_proto.ksi();
  strncpy(
    (char*) state_nas_emm_auth_proc->rand,
    auth_proc_proto.rand().c_str(),
    AUTH_RAND_SIZE);
  strncpy(
    (char*) state_nas_emm_auth_proc->autn,
    auth_proc_proto.autn().c_str(),
    AUTH_AUTN_SIZE);

  if (auth_proc_proto.has_unchecked_imsi()) {
    state_nas_emm_auth_proc->unchecked_imsi =
      (imsi_t*) calloc(1, sizeof(imsi_t));
    proto_to_identity_tuple<imsi_t>(
      auth_proc_proto.unchecked_imsi(),
      state_nas_emm_auth_proc->unchecked_imsi,
      IMSI_BCD8_SIZE);
  }

  state_nas_emm_auth_proc->emm_cause = auth_proc_proto.emm_cause();
  proto_to_nas_timer(auth_proc_proto.t3460(), &state_nas_emm_auth_proc->T3460);
  // update callback functions for auth proc
  set_callbacks_for_auth_proc(state_nas_emm_auth_proc);
  set_notif_callbacks_for_auth_proc(state_nas_emm_auth_proc);
}

void NasStateConverter::nas_emm_smc_proc_to_proto(
  const nas_emm_smc_proc_t* state_nas_emm_smc_proc,
  SmcProc* smc_proc_proto)
{
  OAILOG_INFO(LOG_MME_APP, "Writing smc proc to proto");
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
  const SmcProc& smc_proc_proto,
  nas_emm_smc_proc_t* state_nas_emm_smc_proc)
{
  OAILOG_INFO(LOG_MME_APP, "Reading smc proc from proto");
  state_nas_emm_smc_proc->emm_com_proc.emm_proc.base_proc.type =
    NAS_PROC_TYPE_EMM;
  state_nas_emm_smc_proc->emm_com_proc.emm_proc.type = NAS_EMM_PROC_TYPE_COMMON;
  state_nas_emm_smc_proc->emm_com_proc.type = EMM_COMM_PROC_SMC;
  state_nas_emm_smc_proc->ue_id = smc_proc_proto.ue_id();
  state_nas_emm_smc_proc->retransmission_count =
    smc_proc_proto.retransmission_count();
  state_nas_emm_smc_proc->ksi = smc_proc_proto.ksi();
  state_nas_emm_smc_proc->eea = smc_proc_proto.eea();
  state_nas_emm_smc_proc->eia = smc_proc_proto.eia();
  state_nas_emm_smc_proc->ucs2 = smc_proc_proto.ucs2();
  state_nas_emm_smc_proc->uea = smc_proc_proto.uea();
  state_nas_emm_smc_proc->uia = smc_proc_proto.uia();
  state_nas_emm_smc_proc->gea = smc_proc_proto.gea();
  state_nas_emm_smc_proc->umts_present = smc_proc_proto.umts_present();
  state_nas_emm_smc_proc->gprs_present = smc_proc_proto.gprs_present();
  state_nas_emm_smc_proc->selected_eea = smc_proc_proto.selected_eea();
  state_nas_emm_smc_proc->selected_eia = smc_proc_proto.selected_eia();
  state_nas_emm_smc_proc->saved_selected_eea =
    smc_proc_proto.saved_selected_eea();
  state_nas_emm_smc_proc->saved_selected_eia =
    smc_proc_proto.saved_selected_eia();
  state_nas_emm_smc_proc->saved_eksi = smc_proc_proto.saved_eksi();
  state_nas_emm_smc_proc->saved_overflow = smc_proc_proto.saved_overflow();
  state_nas_emm_smc_proc->saved_seq_num = smc_proc_proto.saved_seq_num();
  state_nas_emm_smc_proc->saved_sc_type =
    (emm_sc_type_t) smc_proc_proto.saved_sc_type();
  state_nas_emm_smc_proc->notify_failure = smc_proc_proto.notify_failure();
  state_nas_emm_smc_proc->is_new = smc_proc_proto.is_new();
  state_nas_emm_smc_proc->imeisv_request = smc_proc_proto.imeisv_request();
  set_notif_callbacks_for_smc_proc(state_nas_emm_smc_proc);
  set_callbacks_for_smc_proc(state_nas_emm_smc_proc);
}

void NasStateConverter::nas_proc_mess_sign_to_proto(
  const nas_proc_mess_sign_t* state_nas_proc_mess_sign,
  NasProcMessSign* nas_proc_mess_sign_proto)
{
  nas_proc_mess_sign_proto->set_puid(state_nas_proc_mess_sign->puid);
  nas_proc_mess_sign_proto->set_digest(
    (void*) state_nas_proc_mess_sign->digest, NAS_MSG_DIGEST_SIZE);
  nas_proc_mess_sign_proto->set_digest_length(
    state_nas_proc_mess_sign->digest_length);
  nas_proc_mess_sign_proto->set_nas_msg_length(
    state_nas_proc_mess_sign->nas_msg_length);
}

void NasStateConverter::proto_to_nas_proc_mess_sign(
  const NasProcMessSign& nas_proc_mess_sign_proto,
  nas_proc_mess_sign_t* state_nas_proc_mess_sign)
{
  state_nas_proc_mess_sign->puid = nas_proc_mess_sign_proto.puid();
  memcpy(
    state_nas_proc_mess_sign->digest,
    nas_proc_mess_sign_proto.digest().c_str(),
    NAS_MSG_DIGEST_SIZE);
  state_nas_proc_mess_sign->digest_length =
    nas_proc_mess_sign_proto.digest_length();
  state_nas_proc_mess_sign->nas_msg_length =
    nas_proc_mess_sign_proto.nas_msg_length();
}

void NasStateConverter::nas_base_proc_to_proto(
  const nas_base_proc_t* base_proc_p,
  NasBaseProc* base_proc_proto)
{
  base_proc_proto->set_nas_puid(base_proc_p->nas_puid);
  base_proc_proto->set_type(base_proc_p->type);
}

void NasStateConverter::proto_to_nas_base_proc(
  const NasBaseProc& nas_base_proc_proto,
  nas_base_proc_t* state_nas_base_proc)
{
  state_nas_base_proc->nas_puid = nas_base_proc_proto.nas_puid();
  state_nas_base_proc->type = (nas_base_proc_type_t) nas_base_proc_proto.type();
  state_nas_base_proc->success_notif = NULL;
  state_nas_base_proc->failure_notif = NULL;
  state_nas_base_proc->abort = NULL;
  state_nas_base_proc->fail_in = NULL;
  state_nas_base_proc->fail_out = NULL;
  state_nas_base_proc->time_out = NULL;
}

void NasStateConverter::emm_proc_to_proto(
  const nas_emm_proc_t* emm_proc_p,
  NasEmmProc* emm_proc_proto)
{
  nas_base_proc_to_proto(
    &emm_proc_p->base_proc, emm_proc_proto->mutable_base_proc());
  emm_proc_proto->set_type(emm_proc_p->type);
  emm_proc_proto->set_previous_emm_fsm_state(
    emm_proc_p->previous_emm_fsm_state);
}

void NasStateConverter::proto_to_nas_emm_proc(
  const NasEmmProc& nas_emm_proc_proto,
  nas_emm_proc_t* state_nas_emm_proc)
{
  proto_to_nas_base_proc(
    nas_emm_proc_proto.base_proc(), &state_nas_emm_proc->base_proc);
  state_nas_emm_proc->type = (nas_emm_proc_type_t) nas_emm_proc_proto.type();
  state_nas_emm_proc->previous_emm_fsm_state =
    (emm_fsm_state_t) nas_emm_proc_proto.previous_emm_fsm_state();
  state_nas_emm_proc->delivered = NULL;
  state_nas_emm_proc->not_delivered = NULL;
  state_nas_emm_proc->not_delivered_ho = NULL;
}

void NasStateConverter::emm_specific_proc_to_proto(
  const nas_emm_specific_proc_t* state_emm_specific_proc,
  NasEmmProcWithType* emm_proc_with_type)
{
  OAILOG_INFO(LOG_MME_APP, "Writing specific procs to proto");
  emm_proc_to_proto(
    &state_emm_specific_proc->emm_proc, emm_proc_with_type->mutable_emm_proc());
  switch (state_emm_specific_proc->type) {
    case EMM_SPEC_PROC_TYPE_ATTACH: {
      OAILOG_INFO(LOG_MME_APP, "Writing attach proc to proto");
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
    default: break;
  }
}

// This function allocated memory for any specific procedure stored
void NasStateConverter::proto_to_emm_specific_proc(
  const NasEmmProcWithType& proto_emm_proc_with_type,
  emm_procedures_t* state_emm_procedures)
{
  OAILOG_INFO(LOG_MME_APP, "Reading specific procs from proto");
  // read attach or detach proc based on message type present
  switch (proto_emm_proc_with_type.MessageTypes_case()) {
    case NasEmmProcWithType::kAttachProc: {
      OAILOG_INFO(LOG_MME_APP, "Reading attach proc from proto");
      state_emm_procedures->emm_specific_proc =
        (nas_emm_specific_proc_t*) calloc(1, sizeof(nas_emm_attach_proc_t));
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
    case NasEmmProcWithType::kDetachProc: {
      state_emm_procedures->emm_specific_proc =
        (nas_emm_specific_proc_t*) calloc(1, sizeof(nas_emm_detach_proc_t));
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
    default: break;
  }
}

void NasStateConverter::emm_common_proc_to_proto(
  const emm_procedures_t* state_emm_procedures,
  EmmProcedures* emm_procedures_proto)
{
  OAILOG_INFO(LOG_MME_APP, "Writing common procs to proto");
  nas_emm_common_procedure_t* p1 =
    LIST_FIRST(&state_emm_procedures->emm_common_procs);
  while (p1) {
    NasEmmProcWithType* nas_emm_proc_with_type_proto =
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
      default: break;
    }
    p1 = LIST_NEXT(p1, entries);
  }
}

void NasStateConverter::insert_proc_into_emm_common_procs(
  emm_procedures_t* state_emm_procedures,
  nas_emm_common_proc_t* emm_com_proc)
{
  nas_emm_common_procedure_t* wrapper =
    (nas_emm_common_procedure_t*) calloc(1, sizeof(*wrapper));
  if (!wrapper) return;

  wrapper->proc = emm_com_proc;
  LIST_INSERT_HEAD(&state_emm_procedures->emm_common_procs, wrapper, entries);
  OAILOG_INFO(LOG_NAS_EMM, "New COMMON PROC added from state\n");
}

void NasStateConverter::proto_to_emm_common_proc(
  const EmmProcedures& emm_procedures_proto,
  emm_context_t* state_emm_context)
{
  OAILOG_INFO(LOG_MME_APP, "Reading common procs from proto");
  auto proto_common_procs = emm_procedures_proto.emm_common_proc();
  for (auto ptr = proto_common_procs.begin(); ptr < proto_common_procs.end();
       ptr++) {
    switch (ptr->MessageTypes_case()) {
      case NasEmmProcWithType::kAuthProc: {
        OAILOG_INFO(LOG_NAS_EMM, "Inserting AUTH PROC from state\n");
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
      case NasEmmProcWithType::kSmcProc: {
        OAILOG_INFO(LOG_NAS_EMM, "Inserting SMC PROC from state\n");
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
      default: break;
    }
  }
}

void NasStateConverter::eutran_vectors_to_proto(
  eutran_vector_t** state_eutran_vector_array,
  int num_vectors,
  AuthInfoProc* auth_info_proc_proto)
{
  AuthVector* eutran_vector_proto = nullptr;
  for (int i = 0; i < num_vectors; i++) {
    eutran_vector_proto = auth_info_proc_proto->add_vector();
    memcpy(
      eutran_vector_proto->mutable_kasme(),
      state_eutran_vector_array[i]->kasme,
      KASME_LENGTH_OCTETS);
    memcpy(
      eutran_vector_proto->mutable_rand(),
      state_eutran_vector_array[i]->rand,
      RAND_LENGTH_OCTETS);
    memcpy(
      eutran_vector_proto->mutable_autn(),
      state_eutran_vector_array[i]->autn,
      AUTN_LENGTH_OCTETS);
    memcpy(
      eutran_vector_proto->mutable_xres(),
      state_eutran_vector_array[i]->xres.data,
      state_eutran_vector_array[i]->xres.size);
  }
}

void NasStateConverter::proto_to_eutran_vectors(
  const AuthInfoProc& auth_info_proc_proto,
  nas_auth_info_proc_t* state_nas_auth_info_proc)
{
  auto proto_vectors = auth_info_proc_proto.vector();
  int i = 0;
  for (auto ptr = proto_vectors.begin(); ptr < proto_vectors.end(); ptr++) {
    eutran_vector_t* this_vector =
      (eutran_vector_t*) calloc(1, sizeof(eutran_vector_t));
    strncpy((char*) this_vector->kasme, ptr->kasme().c_str(), AUTH_KASME_SIZE);
    strncpy((char*) this_vector->rand, ptr->rand().c_str(), AUTH_RAND_SIZE);
    strncpy((char*) this_vector->autn, ptr->autn().c_str(), AUTH_AUTN_SIZE);
    this_vector->xres.size = ptr->xres().length();
    strncpy(
      (char*) this_vector->xres.data,
      ptr->xres().c_str(),
      this_vector->xres.size);
    state_nas_auth_info_proc->vector[i] = this_vector;
    i++;
  }
  state_nas_auth_info_proc->nb_vectors = i;
}

void NasStateConverter::nas_auth_info_proc_to_proto(
  nas_auth_info_proc_t* state_nas_auth_info_proc,
  AuthInfoProc* auth_info_proc_proto)
{
  auth_info_proc_proto->set_request_sent(
    state_nas_auth_info_proc->request_sent);
  eutran_vectors_to_proto(
    state_nas_auth_info_proc->vector,
    state_nas_auth_info_proc->nb_vectors,
    auth_info_proc_proto);
  auth_info_proc_proto->set_nas_cause(state_nas_auth_info_proc->nas_cause);
  nas_timer_to_proto(
    state_nas_auth_info_proc->timer_s6a,
    auth_info_proc_proto->mutable_timer_s6a());
  auth_info_proc_proto->set_ue_id(state_nas_auth_info_proc->ue_id);
  auth_info_proc_proto->set_resync(state_nas_auth_info_proc->resync);
}

void NasStateConverter::proto_to_nas_auth_info_proc(
  const AuthInfoProc& auth_info_proc_proto,
  nas_auth_info_proc_t* state_nas_auth_info_proc)
{
  state_nas_auth_info_proc->request_sent = auth_info_proc_proto.request_sent();
  proto_to_eutran_vectors(auth_info_proc_proto, state_nas_auth_info_proc);
  state_nas_auth_info_proc->nas_cause = auth_info_proc_proto.nas_cause();
  state_nas_auth_info_proc->ue_id = auth_info_proc_proto.ue_id();
  state_nas_auth_info_proc->resync = auth_info_proc_proto.resync();
  proto_to_nas_timer(
    auth_info_proc_proto.timer_s6a(), &state_nas_auth_info_proc->timer_s6a);
  // update success_notif and failure_notif
  set_callbacks_for_auth_info_proc(state_nas_auth_info_proc);
}

void NasStateConverter::nas_cn_procs_to_proto(
  const emm_procedures_t* state_emm_procedures,
  EmmProcedures* emm_procedures_proto)
{
  nas_cn_procedure_t* p1 = LIST_FIRST(&state_emm_procedures->cn_procs);
  while (p1) {
    NasCnProc* nas_cn_proc_proto = emm_procedures_proto->add_cn_proc();
    nas_base_proc_to_proto(
      &p1->proc->base_proc, nas_cn_proc_proto->mutable_base_proc());
    switch (p1->proc->type) {
      case CN_PROC_AUTH_INFO: {
        nas_auth_info_proc_t* state_auth_info_proc = (nas_auth_info_proc_t*) p1;
        nas_auth_info_proc_to_proto(
          state_auth_info_proc, nas_cn_proc_proto->mutable_auth_info_proc());
      } break;
      default:
        OAILOG_INFO(
          LOG_NAS,
          "EMM_CN: Unknown procedure type, cannot convert"
          "to proto");
        break;
    }
    p1 = LIST_NEXT(p1, entries);
  }
}

void NasStateConverter::insert_proc_into_cn_procs(
  emm_procedures_t* state_emm_procedures,
  nas_cn_proc_t* cn_proc)
{
  nas_cn_procedure_t* wrapper =
    (nas_cn_procedure_t*) calloc(1, sizeof(*wrapper));
  if (!wrapper) return;
  wrapper->proc = cn_proc;
  LIST_INSERT_HEAD(&state_emm_procedures->cn_procs, wrapper, entries);
  OAILOG_TRACE(LOG_NAS_EMM, "New EMM_COMM_PROC_SMC\n");
}

void NasStateConverter::proto_to_nas_cn_proc(
  const EmmProcedures& emm_procedures_proto,
  emm_procedures_t* state_emm_procedures)
{
  auto proto_cn_procs = emm_procedures_proto.cn_proc();
  for (auto ptr = proto_cn_procs.begin(); ptr < proto_cn_procs.end(); ptr++) {
    switch (ptr->MessageTypes_case()) {
      case NasCnProc::kAuthInfoProc: {
        nas_auth_info_proc_t* state_auth_info_proc =
          (nas_auth_info_proc_t*) calloc(1, sizeof(*state_auth_info_proc));
        OAILOG_INFO(LOG_NAS_EMM, "Inserting AUTH INFO PROC from state\n");
        proto_to_nas_base_proc(
          ptr->base_proc(), &state_auth_info_proc->cn_proc.base_proc);
        proto_to_nas_auth_info_proc(
          ptr->auth_info_proc(), state_auth_info_proc);
        insert_proc_into_cn_procs(
          state_emm_procedures, &state_auth_info_proc->cn_proc);
      }
      default: break;
    }
  }
}

void NasStateConverter::mess_sign_array_to_proto(
  const emm_procedures_t* state_emm_procedures,
  EmmProcedures* emm_procedures_proto)
{
  for (int i = 0; i < MAX_NAS_PROC_MESS_SIGN; i++) {
    NasProcMessSign* nas_proc_mess_sign_proto =
      emm_procedures_proto->add_nas_proc_mess_sign();
    nas_proc_mess_sign_to_proto(
      &state_emm_procedures->nas_proc_mess_sign[i], nas_proc_mess_sign_proto);
  }
}

void NasStateConverter::proto_to_mess_sign_array(
  const EmmProcedures& emm_procedures_proto,
  emm_procedures_t* state_emm_procedures)
{
  int i = 0;
  auto proto_mess_sign_array = emm_procedures_proto.nas_proc_mess_sign();
  for (auto ptr = proto_mess_sign_array.begin();
       ptr < proto_mess_sign_array.end();
       ptr++) {
    proto_to_nas_proc_mess_sign(
      *ptr, &state_emm_procedures->nas_proc_mess_sign[i]);
    i++;
  }
}

void NasStateConverter::emm_procedures_to_proto(
  const emm_procedures_t* state_emm_procedures,
  EmmProcedures* emm_procedures_proto)
{
  if (state_emm_procedures->emm_specific_proc) {
    emm_specific_proc_to_proto(
      state_emm_procedures->emm_specific_proc,
      emm_procedures_proto->mutable_emm_specific_proc());
  }
  emm_common_proc_to_proto(state_emm_procedures, emm_procedures_proto);

  // cn_procs
  nas_cn_procs_to_proto(state_emm_procedures, emm_procedures_proto);
  NasEmmProcWithType* emm_proc_with_type =
    emm_procedures_proto->mutable_emm_con_mngt_proc();

  if (state_emm_procedures->emm_con_mngt_proc) {
    emm_proc_to_proto(
      &state_emm_procedures->emm_con_mngt_proc->emm_proc,
      emm_proc_with_type->mutable_emm_proc());
  }

  emm_procedures_proto->set_nas_proc_mess_sign_next_location(
    state_emm_procedures->nas_proc_mess_sign_next_location);

  mess_sign_array_to_proto(state_emm_procedures, emm_procedures_proto);
  // temporarily storing the address in redis to make sure other state is
  // passing correctly
  emm_procedures_proto->set_pointer(
    reinterpret_cast<uintptr_t>(state_emm_procedures));
}

void NasStateConverter::proto_to_emm_procedures(
  const EmmProcedures& emm_procedures_proto,
  emm_context_t* state_emm_context)
{
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
  state_emm_procedures->emm_con_mngt_proc =
    (nas_emm_con_mngt_proc_t*) calloc(1, sizeof(nas_emm_con_mngt_proc_t));
  if (emm_procedures_proto.has_emm_con_mngt_proc()) {
    proto_to_nas_emm_proc(
      emm_procedures_proto.emm_con_mngt_proc().emm_proc(),
      &state_emm_procedures->emm_con_mngt_proc->emm_proc);
  }
  state_emm_procedures->nas_proc_mess_sign_next_location =
    emm_procedures_proto.nas_proc_mess_sign_next_location();
  proto_to_mess_sign_array(emm_procedures_proto, state_emm_procedures);
}

void NasStateConverter::auth_vectors_to_proto(
  const auth_vector_t* state_auth_vector_array,
  int num_vectors,
  EmmContext* emm_context_proto)
{
  AuthVector* auth_vector_proto = nullptr;
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
  const EmmContext& emm_context_proto,
  auth_vector_t* state_auth_vector)
{
  auto proto_vectors = emm_context_proto.vector();
  int i = 0;
  for (auto ptr = proto_vectors.begin(); ptr < proto_vectors.end(); ptr++) {
    strncpy(
      (char*) state_auth_vector[i].kasme,
      ptr->kasme().c_str(),
      AUTH_KASME_SIZE);
    strncpy(
      (char*) state_auth_vector[i].rand, ptr->rand().c_str(), AUTH_RAND_SIZE);
    strncpy(
      (char*) state_auth_vector[i].autn, ptr->autn().c_str(), AUTH_AUTN_SIZE);
    strncpy(
      (char*) state_auth_vector[i].xres,
      ptr->xres().c_str(),
      ptr->xres().length());
    i++;
  }
  return i;
}

void NasStateConverter::emm_security_context_to_proto(
  const emm_security_context_t* state_emm_security_context,
  EmmSecurityContext* emm_security_context_proto)
{
  emm_security_context_proto->set_sc_type(state_emm_security_context->sc_type);
  emm_security_context_proto->set_eksi(state_emm_security_context->eksi);
  emm_security_context_proto->set_vector_index(
    state_emm_security_context->vector_index);
  emm_security_context_proto->set_knas_enc(
    state_emm_security_context->knas_enc, AUTH_KNAS_ENC_SIZE);
  emm_security_context_proto->set_knas_int(
    state_emm_security_context->knas_int, AUTH_KNAS_INT_SIZE);

  // Count values
  EmmSecurityContext_Count* dl_count_proto =
    emm_security_context_proto->mutable_dl_count();
  dl_count_proto->set_spare(state_emm_security_context->dl_count.spare);
  dl_count_proto->set_overflow(state_emm_security_context->dl_count.overflow);
  dl_count_proto->set_seq_num(state_emm_security_context->dl_count.seq_num);
  EmmSecurityContext_Count* ul_count_proto =
    emm_security_context_proto->mutable_ul_count();
  ul_count_proto->set_spare(state_emm_security_context->ul_count.spare);
  ul_count_proto->set_overflow(state_emm_security_context->ul_count.overflow);
  ul_count_proto->set_seq_num(state_emm_security_context->ul_count.seq_num);
  EmmSecurityContext_Count* kenb_ul_count_proto =
    emm_security_context_proto->mutable_kenb_ul_count();
  kenb_ul_count_proto->set_spare(
    state_emm_security_context->kenb_ul_count.spare);
  kenb_ul_count_proto->set_overflow(
    state_emm_security_context->kenb_ul_count.overflow);
  kenb_ul_count_proto->set_seq_num(
    state_emm_security_context->kenb_ul_count.seq_num);

  // TODO convert capability to proto

  // Security algorithm
  EmmSecurityContext_SelectedAlgorithms* selected_algorithms_proto =
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
  const EmmSecurityContext& emm_security_context_proto,
  emm_security_context_t* state_emm_security_context)
{
  state_emm_security_context->sc_type =
    (emm_sc_type_t) emm_security_context_proto.sc_type();
  state_emm_security_context->eksi = emm_security_context_proto.eksi();
  state_emm_security_context->vector_index =
    emm_security_context_proto.vector_index();
  strcpy(
    (char*) state_emm_security_context->knas_enc,
    emm_security_context_proto.knas_enc().c_str());
  strcpy(
    (char*) state_emm_security_context->knas_int,
    emm_security_context_proto.knas_int().c_str());

  // Count values
  const EmmSecurityContext_Count& dl_count_proto =
    emm_security_context_proto.dl_count();
  state_emm_security_context->dl_count.spare = dl_count_proto.spare();
  state_emm_security_context->dl_count.overflow = dl_count_proto.overflow();
  state_emm_security_context->dl_count.seq_num = dl_count_proto.seq_num();
  const EmmSecurityContext_Count& ul_count_proto =
    emm_security_context_proto.ul_count();
  state_emm_security_context->ul_count.spare = ul_count_proto.spare();
  state_emm_security_context->ul_count.overflow = ul_count_proto.overflow();
  state_emm_security_context->ul_count.seq_num = ul_count_proto.seq_num();
  const EmmSecurityContext_Count& kenb_ul_count_proto =
    emm_security_context_proto.kenb_ul_count();
  state_emm_security_context->kenb_ul_count.spare = kenb_ul_count_proto.spare();
  state_emm_security_context->kenb_ul_count.overflow =
    kenb_ul_count_proto.overflow();
  state_emm_security_context->kenb_ul_count.seq_num =
    kenb_ul_count_proto.seq_num();

  // TODO read capability from proto

  // Security algorithm
  const EmmSecurityContext_SelectedAlgorithms& selected_algorithms_proto =
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
  strcpy(
    (char*) state_emm_security_context->next_hop,
    emm_security_context_proto.next_hop().c_str());
  state_emm_security_context->next_hop_chaining_count =
    emm_security_context_proto.next_hop_chaining_count();
}

void NasStateConverter::emm_context_to_proto(
  const emm_context_t* state_emm_context,
  EmmContext* emm_context_proto)
{
  emm_context_proto->set_imsi64(state_emm_context->_imsi64);
  identity_tuple_to_proto<imsi_t>(
    &state_emm_context->_imsi,
    emm_context_proto->mutable_imsi(),
    IMSI_BCD8_SIZE);
  emm_context_proto->set_saved_imsi64(state_emm_context->saved_imsi64);
  identity_tuple_to_proto<imei_t>(
    &state_emm_context->_imei,
    emm_context_proto->mutable_imei(),
    IMEI_BCD8_SIZE);
  identity_tuple_to_proto<imeisv_t>(
    &state_emm_context->_imeisv,
    emm_context_proto->mutable_imeisv(),
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
}

void NasStateConverter::proto_to_emm_context(
  const EmmContext& emm_context_proto,
  emm_context_t* state_emm_context)
{
  state_emm_context->_imsi64 = emm_context_proto.imsi64();
  proto_to_identity_tuple<imsi_t>(
    emm_context_proto.imsi(), &state_emm_context->_imsi, IMSI_BCD8_SIZE);
  state_emm_context->saved_imsi64 = emm_context_proto.saved_imsi64();

  proto_to_identity_tuple<imei_t>(
    emm_context_proto.imei(), &state_emm_context->_imei, IMEI_BCD8_SIZE);
  proto_to_identity_tuple<imeisv_t>(
    emm_context_proto.imeisv(), &state_emm_context->_imeisv, IMEISV_BCD8_SIZE);

  state_emm_context->emm_cause = emm_context_proto.emm_cause();
  state_emm_context->_emm_fsm_state =
    (emm_fsm_state_t) emm_context_proto.emm_fsm_state();
  state_emm_context->attach_type = emm_context_proto.attach_type();
  if (emm_context_proto.has_emm_procedures()) {
    state_emm_context->emm_procedures =
      (emm_procedures_t*) calloc(1, sizeof(*state_emm_context->emm_procedures));
    proto_to_emm_procedures(
      emm_context_proto.emm_procedures(), state_emm_context);
  }
  nas_emm_auth_proc_t* auth_proc =
    get_nas_common_procedure_authentication(state_emm_context);
  if (auth_proc) {
    OAILOG_INFO(
      LOG_MME_APP,
      "Found non-null auth proc with ue_id" RAND_FORMAT,
      RAND_DISPLAY(auth_proc->rand));
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
  state_emm_context->is_dynamic = emm_context_proto.is_dynamic();
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
}

} // namespace lte
} // namespace magma
