/**
 * Copyright 2021 The Magma Authors.
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
//--C includes -----------------------------------------------------------------
extern "C" {
#include "emm_data.h"
#include "emm_proc.h"
#include "esm_proc.h"
#include "nas_procedures.h"
#include "log.h"
#include "dynamic_memory_check.h"
#include "intertask_interface.h"
#include "common_types.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "mme_app_state.h"
#include "3gpp_23.003.h"
}
#include "mme_app_test_serialization.h"
#include "mme_nas_state_generated.h"
#include <mcheck.h>
//--C++ includes ---------------------------------------------------------------
#include <chrono>
#include <cmath>
#include <vector>
//--Other includes -------------------------------------------------------------
#include "mme_app_state_manager.h"
extern task_zmq_ctx_t main_zmq_ctx;

using magma::lte::MmeNasStateManager;
using namespace magma::lte::test_flat_buffer;

//------------------------------------------------------------------------------
std::vector<std::pair<UeMmContextBuilder*, UeMmContext*>>
mme_app_fb_allocate_ues(flatbuffers::FlatBufferBuilder& builder, uint num_ues);
void mme_app_fb_deallocate_ues(
    std::vector<std::pair<UeMmContextBuilder*, UeMmContext*>>& contexts);
void mme_app_fb_insert_ues(
    std::vector<std::pair<UeMmContextBuilder*, UeMmContext*>>& contexts);
void mme_app_fb_serialize_ues(
    std::vector<std::pair<UeMmContextBuilder*, UeMmContext*>>& contexts);
void mme_app_fb_deserialize_ues(void);

//------------------------------------------------------------------------------
void mme_app_fb_schedule_test_serialization(uint num_ues) {
  MessageDef* message_p = itti_alloc_new_message(
      TASK_UNKNOWN, MME_APP_TEST_FLATBUFFER_SERIALIZATION);
  MME_APP_TEST_FLATBUFFER_SERIALIZATION(message_p).num_ues = num_ues;
  send_msg_to_task(&main_zmq_ctx, TASK_MME_APP, message_p);
  return;
}

Imsi* imsi2fb(const imsi_t& imsi) {
  Imsi* imsi_fb = new Imsi(imsi.u.value, imsi.length);
  return imsi_fb;
}
Guti* guti2fb(guti_t& guti) {
  Guti* guti_fb = new Guti(
      Gummei(
          Plmn(
              guti.gummei.plmn.mcc_digit1, guti.gummei.plmn.mcc_digit2,
              guti.gummei.plmn.mcc_digit3, guti.gummei.plmn.mnc_digit1,
              guti.gummei.plmn.mnc_digit2, guti.gummei.plmn.mnc_digit3),
          guti.gummei.mme_gid, guti.gummei.mme_code),
      guti.m_tmsi);
  return guti_fb;
}

//------------------------------------------------------------------------------
std::vector<std::pair<UeMmContextBuilder*, UeMmContext*>>
mme_app_fb_allocate_ues(uint num_ues) {
  enb_s1ap_id_key_t enb_s1ap_id_key = INVALID_ENB_UE_S1AP_ID_KEY;
  enb_ue_s1ap_id_t enb_ue_s1ap_id   = rand() & 0X00FFFFFF;
  mme_ue_s1ap_id_t mme_ue_s1ap_id   = rand();
  std::vector<std::pair<UeMmContextBuilder*, UeMmContext*>> contexts;
  uint64_t kFirstImsi = 1010000000000;

  contexts.reserve(num_ues);

  for (int i = 0; i < num_ues; i++) {
    // TODO find the right size to be allocated
    flatbuffers::FlatBufferBuilder builder(
        sizeof(UeMmContext) + sizeof(EmmContext) + sizeof(EsmContext));
    // Force all fields you set to actually be written. This, of course,
    // increases the size of the buffer somewhat, but this may be
    // acceptable for a mutable buffer.
    builder.ForceDefaults(true);
    UeMmContextBuilder* uemmcontext_builder = new UeMmContextBuilder(builder);
    EmmContextBuilder emmcontext_builder(builder);
    EsmContextBuilder esmcontext_builder(builder);

    enb_ue_s1ap_id++;
    mme_ue_s1ap_id++;

    imsi64_t imsi64    = kFirstImsi + i;
    imsi_t imsi        = {};
    imsi.u.num.digit1  = (uint8_t)((imsi64_t)(imsi64 / std::pow(10, 14)) % 10);
    imsi.u.num.digit2  = (uint8_t)((imsi64_t)(imsi64 / std::pow(10, 13)) % 10);
    imsi.u.num.digit3  = (uint8_t)((imsi64_t)(imsi64 / std::pow(10, 12)) % 10);
    imsi.u.num.digit4  = (uint8_t)((imsi64_t)(imsi64 / std::pow(10, 11)) % 10);
    imsi.u.num.digit5  = (uint8_t)((imsi64_t)(imsi64 / std::pow(10, 10)) % 10);
    imsi.u.num.digit6  = (uint8_t)((imsi64_t)(imsi64 / std::pow(10, 9)) % 10);
    imsi.u.num.digit7  = (uint8_t)((imsi64_t)(imsi64 / std::pow(10, 8)) % 10);
    imsi.u.num.digit8  = (uint8_t)((imsi64_t)(imsi64 / std::pow(10, 7)) % 10);
    imsi.u.num.digit9  = (uint8_t)((imsi64_t)(imsi64 / std::pow(10, 6)) % 10);
    imsi.u.num.digit10 = (uint8_t)((imsi64_t)(imsi64 / std::pow(10, 5)) % 10);
    imsi.u.num.digit11 = (uint8_t)((imsi64_t)(imsi64 / std::pow(10, 4)) % 10);
    imsi.u.num.digit12 = (uint8_t)((imsi64_t)(imsi64 / std::pow(10, 3)) % 10);
    imsi.u.num.digit13 = (uint8_t)((imsi64_t)(imsi64 / std::pow(10, 2)) % 10);
    imsi.u.num.digit14 = (uint8_t)((imsi64_t)(imsi64 / std::pow(10, 1)) % 10);
    imsi.u.num.digit15 = (uint8_t)((imsi64_t)(imsi64 / std::pow(10, 0)) % 10);
    imsi.u.num.parity  = 0xF;

    Imsi* imsi_fb = imsi2fb(imsi);

    Plmn plmn(0, 0, 1, 0, 1, 0xF);
    Gummei gummei(plmn, 1, 1);
    Guti* guti = new Guti(gummei, 2106150532 + i);

    Guti* old_guti = new Guti();
    old_guti->mutate_m_tmsi(429496729 + i);

    esmcontext_builder.add_n_active_ebrs(2);
    magma::lte::test_flat_buffer::BearerQos* bearer_qos = new BearerQos();
    bearer_qos->mutate_pci(true);
    bearer_qos->mutate_pl(15);
    bearer_qos->mutate_qci(5);
    EsmProcDataBuilder esm_proc_data_builder(builder);
    esm_proc_data_builder.add_bearer_qos(bearer_qos);
    esm_proc_data_builder.add_pdn_cid(1);
    auto apn = builder.CreateString("ims");
    esm_proc_data_builder.add_apn(apn);
    esm_proc_data_builder.add_pdn_type(PdnTypeValue_IPv4);
    esm_proc_data_builder.add_request_type(1);
    esm_proc_data_builder.add_pti(1);
    esmcontext_builder.add_esm_proc_data(esm_proc_data_builder.Finish());

    emmcontext_builder.add_saved_imsi64(imsi64);
    emmcontext_builder.add__imsi(imsi_fb);
    emmcontext_builder.add__guti(guti);
    emmcontext_builder.add__old_guti(old_guti);
    emmcontext_builder.add_emm_cause(UINT32_MAX);
    emmcontext_builder.add__emm_fsm_state(EmmFsmState_EMM_REGISTERED);
    emmcontext_builder.add_esm_ctx(esmcontext_builder.Finish());

    MME_APP_ENB_S1AP_ID_KEY(
        enb_s1ap_id_key, rand() & 0X0000FFFF, enb_ue_s1ap_id);

    // uemmcontext_builder->add_cell_age(cell_age);
    // uemmcontext_builder->add_time_ics_rsp_timer_started(time_ics_rsp_timer_started);
    // uemmcontext_builder->add_time_paging_response_timer_started(time_paging_response_timer_started);
    // uemmcontext_builder->add_time_implicit_detach_timer_started(time_implicit_detach_timer_started);
    // uemmcontext_builder->add_time_mobile_reachability_timer_started(time_mobile_reachability_timer_started);
    // uemmcontext_builder->add_rau_tau_timer(rau_tau_timer);
    // uemmcontext_builder->add_sgs_context(sgs_context);
    // uemmcontext_builder->add_cs_fallback_indicator(cs_fallback_indicator);
    // uemmcontext_builder->add_reg_sub(reg_sub);
    // uemmcontext_builder->add_ue_radio_capability(ue_radio_capability);
    // uemmcontext_builder->add_bearer_contexts(bearer_contexts);
    // uemmcontext_builder->add_pdn_contexts(pdn_contexts);
    // uemmcontext_builder->add_used_ue_ambr(used_ue_ambr);
    // uemmcontext_builder->add_subscribed_ue_ambr(subscribed_ue_ambr);
    // uemmcontext_builder->add_mme_teid_s11(mme_teid_s11);
    // uemmcontext_builder->add_apn_oi_replacement(apn_oi_replacement);
    // uemmcontext_builder->add_access_restriction_data(access_restriction_data);
    // uemmcontext_builder->add_apn_config_profile(apn_config_profile);
    // uemmcontext_builder->add_lai(lai);
    // uemmcontext_builder->add_e_utran_cgi(e_utran_cgi);
    uemmcontext_builder->add_mme_ue_s1ap_id(mme_ue_s1ap_id);
    uemmcontext_builder->add_enb_s1ap_id_key(enb_s1ap_id_key);
    uemmcontext_builder->add_enb_ue_s1ap_id(enb_ue_s1ap_id);
    // uemmcontext_builder->add_sctp_assoc_id_key(sctp_assoc_id_key);
    uemmcontext_builder->add_emm_context(emmcontext_builder.Finish());
    // uemmcontext_builder->add_ue_context_rel_cause(ue_context_rel_cause);
    // uemmcontext_builder->add_msisdn(msisdn);
    // uemmcontext_builder->add_paging_retx_count(paging_retx_count);
    // uemmcontext_builder->add_num_reg_sub(num_reg_sub);
    // uemmcontext_builder->add_granted_service(granted_service);
    // uemmcontext_builder->add_path_switch_req(path_switch_req);
    // uemmcontext_builder->add_subscription_known(subscription_known);
    // uemmcontext_builder->add_ppf(ppf);
    // uemmcontext_builder->add_location_info_confirmed_in_hss(location_info_confirmed_in_hss);
    // uemmcontext_builder->add_hss_initiated_detach(hss_initiated_detach);
    // uemmcontext_builder->add_send_ue_purge_request(send_ue_purge_request);
    // uemmcontext_builder->add_nb_active_pdn_contexts(nb_active_pdn_contexts);
    // uemmcontext_builder->add_network_access_mode(network_access_mode);
    // uemmcontext_builder->add_subscriber_status(subscriber_status);
    // uemmcontext_builder->add_sgs_detach_type(sgs_detach_type);
    // uemmcontext_builder->add_attach_type(attach_type);
    // uemmcontext_builder->add_ecm_state(ecm_state);
    // uemmcontext_builder->add_mm_state(mm_state);

    auto ue_ctxt = uemmcontext_builder->Finish();

    UeMmContext* ue_mm_context =
        GetMutableUeMmContext(uemmcontext_builder->fbb_.GetBufferPointer());

    contexts.push_back(std::pair<UeMmContextBuilder*, UeMmContext*>(
        uemmcontext_builder, ue_mm_context));
  }
  return contexts;
}
//------------------------------------------------------------------------------
void mme_app_fb_deallocate_ues(
    std::vector<std::pair<UeMmContextBuilder*, UeMmContext*>>& contexts) {
  // TODO
  contexts.clear();
}

//------------------------------------------------------------------------------
void mme_app_fb_serialize_ues(
    mme_app_desc_t* mme_app_desc,
    std::vector<std::pair<UeMmContextBuilder*, UeMmContext*>>& contexts) {
  for (std::vector<std::pair<UeMmContextBuilder*, UeMmContext*>>::iterator it =
           contexts.begin();
       it != contexts.end(); ++it) {
    auto imsi_str = MmeNasStateManager::getInstance().get_imsi_str(
        (*it).second->emm_context()->_imsi64());
    MmeNasStateManager::getInstance().write_ue_state_to_db(
        (*it).first, (*it).second, imsi_str);
  }
}
//------------------------------------------------------------------------------
void mme_app_fb_insert_ues(
    std::vector<std::pair<UeMmContextBuilder*, UeMmContext*>>& contexts) {
  for (std::vector<std::pair<UeMmContextBuilder*, UeMmContext*>>::iterator it =
           contexts.begin();
       it != contexts.end(); ++it) {
    // TODO
  }
}

//------------------------------------------------------------------------------
void mme_app_fb_deserialize_ues(void) {
  // TODO
  // mme_app_desc_t* mme_app_desc2 = get_mme_nas_state(true);
  // MmeNasStateManager::getInstance().read_ue_state_from_db();
}
//------------------------------------------------------------------------------
void mme_app_fb_test_serialization(mme_app_desc_t* mme_app_desc, uint num_ues) {
  std::vector<std::pair<UeMmContextBuilder*, UeMmContext*>> contexts =
      mme_app_fb_allocate_ues(num_ues);

  mme_app_fb_insert_ues(contexts);

  auto start_ctxt_to_proto = std::chrono::high_resolution_clock::now();
  mme_app_fb_serialize_ues(mme_app_desc, contexts);
  auto end_ctxt_to_proto = std::chrono::high_resolution_clock::now();
  auto duration_ctxt_to_proto =
      std::chrono::duration_cast<std::chrono::microseconds>(
          end_ctxt_to_proto - start_ctxt_to_proto);
  std::cout << "Time taken by context to proto conversion : "
            << duration_ctxt_to_proto.count() << " microseconds" << std::endl;
  OAILOG_INFO(
      LOG_MME_APP, "Time taken by context to proto conversion : %ld µs\n",
      duration_ctxt_to_proto.count());

  auto start_proto_to_ctxt = std::chrono::high_resolution_clock::now();
  mme_app_fb_deserialize_ues();
  auto end_proto_to_ctxt = std::chrono::high_resolution_clock::now();
  auto duration_proto_to_ctxt =
      std::chrono::duration_cast<std::chrono::microseconds>(
          end_proto_to_ctxt - start_proto_to_ctxt);
  std::cout << "Time taken by proto to context conversion : "
            << duration_proto_to_ctxt.count() << " microseconds" << std::endl;
  OAILOG_INFO(
      LOG_MME_APP, "Time taken by proto to context conversion : %ld µs\n",
      duration_proto_to_ctxt.count());
  /*auto imsi_str = MmeNasStateManager::getInstance().get_imsi_str(imsi64);
  MmeNasStateManager::getInstance().write_ue_state_to_db(
      ue_context, imsi_str);
  put_mme_ue_state(mme_app_desc_p, imsi64, force_ue_write);
  put_mme_nas_state(); */
  mme_app_desc_t* mme_app_desc_p = get_mme_nas_state(false);
  if (!mme_app_desc_p) {
    OAILOG_ERROR(LOG_MME_APP, "Failed to fetch mme_app_desc_p \n");
    return;
  }
  mme_app_fb_deallocate_ues(contexts);

  send_terminate_message_fatal(&main_zmq_ctx);
  return;
}