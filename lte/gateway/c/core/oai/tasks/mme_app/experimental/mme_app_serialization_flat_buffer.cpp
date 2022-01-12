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

#include "lte/gateway/c/core/oai/tasks/mme_app/experimental/mme_app_serialization.h"
// --C system includes --------------------------------------------------------
#include <mcheck.h>
#include <sys/time.h>      // rusage()
#include <sys/resource.h>  //rusage()
//--C++ includes ---------------------------------------------------------------
#include <chrono>
#include <cmath>
#include <tuple>
#include <vector>
// --C includes ---------------------------------------------------------------
extern "C" {
#include "lte/gateway/c/core/oai/tasks/nas/emm/emm_data.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/emm_proc.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/esm_proc.h"
#include "lte/gateway/c/core/oai/tasks/nas/nas_procedures.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#include "lte/gateway/c/core/oai/lib/itti/itti_types.h"
#include "lte/gateway/c/core/oai/include/mme_app_state.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_23.003.h"
}
#include "lte/flat/oai/experimental/mme_nas_state_generated.h"
//--Other includes -------------------------------------------------------------
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_state_manager.h"
extern task_zmq_ctx_t main_zmq_ctx;

using magma::lte::MmeNasStateManager;
using namespace magma::lte::test_flat_buffer;

extern void log_rusage_diff(
    struct rusage& ru_first, struct rusage& ru_last, const char* context);
//------------------------------------------------------------------------------
std::vector<std::tuple<uint8_t*, flatbuffers::uoffset_t, UeMmContext*>>
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

Imsi imsi2fb(const imsi_t& imsi) {
  Imsi imsi_fb(imsi.u.value, imsi.length);
  return imsi_fb;
}

Guti guti2fb(const guti_t& guti) {
  Guti guti_fb(
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
std::tuple<uint8_t*, flatbuffers::uoffset_t, UeMmContext*>
build_ue_mm_context() {
  flatbuffers::FlatBufferBuilder builder(20000);
  // Force all fields you set to actually be written. This, of course,
  // increases the size of the buffer somewhat, but this may be
  // acceptable for a mutable buffer.
  builder.ForceDefaults(true);
  UeMmContextBuilder ue_mm_context_builder(builder);

  magma::lte::test_flat_buffer::Msisdn msisdn                         = {};
  magma::lte::test_flat_buffer::S1Cause ue_context_rel_cause          = {};
  magma::lte::test_flat_buffer::MmState mm_state                      = {};
  magma::lte::test_flat_buffer::EcmState ecm_state                    = {};
  magma::lte::test_flat_buffer::EmmContext emm_context                = {};
  uint32_t sctp_assoc_id_key                                          = 0;
  uint32_t enb_ue_s1ap_id                                             = 0;
  uint32_t enb_s1ap_id_key                                            = 0;
  uint32_t mme_ue_s1ap_id                                             = 0;
  uint8_t attach_type                                                 = 0;
  uint8_t sgs_detach_type                                             = 0;
  magma::lte::test_flat_buffer::Ecgi e_utran_cgi                      = {};
  uint64_t cell_age                                                   = 0;
  magma::lte::test_flat_buffer::Lai lai                               = {};
  magma::lte::test_flat_buffer::ApnConfigProfile apn_config_profile   = {};
  magma::lte::test_flat_buffer::SubscriberStatus subscriber_status    = {};
  magma::lte::test_flat_buffer::NetworkAccessMode network_access_mode = {};
  uint32_t access_restriction_data                                    = 0;
  magma::lte::test_flat_buffer::ApnOi apn_oi_replacement              = {};
  uint32_t mme_teid_s11                                               = 0;
  magma::lte::test_flat_buffer::Ambr subscribed_ue_ambr               = {};
  magma::lte::test_flat_buffer::Ambr used_ue_ambr                     = {};
  uint8_t nb_active_pdn_contexts                                      = 0;
  magma::lte::test_flat_buffer::PdnContextArray pdn_contexts          = {};
  magma::lte::test_flat_buffer::BearerContextArray bearer_contexts    = {};
  magma::lte::test_flat_buffer::UeRadioCapability ue_radio_capability = {};
  bool send_ue_purge_request                                          = false;
  bool hss_initiated_detach                                           = false;
  bool location_info_confirmed_in_hss                                 = false;
  bool ppf                                                            = false;
  bool subscription_known                                             = false;
  bool path_switch_req                                                = false;
  magma::lte::test_flat_buffer::GrantedService granted_service        = {};
  uint8_t num_reg_sub                                                 = 0;
  magma::lte::test_flat_buffer::RegionalSubscriptionArray reg_sub     = {};
  int32_t cs_fallback_indicator                                       = 0;
  magma::lte::test_flat_buffer::SgsContext sgs_context                = {};
  uint32_t rau_tau_timer                                              = 0;
  uint32_t time_mobile_reachability_timer_started                     = 0;
  uint32_t time_implicit_detach_timer_started                         = 0;
  uint32_t time_paging_response_timer_started                         = 0;
  uint8_t paging_retx_count                                           = 0;
  uint32_t time_ics_rsp_timer_started                                 = 0;

  ue_mm_context_builder.add_msisdn(&msisdn);
  ue_mm_context_builder.add_ue_context_rel_cause(ue_context_rel_cause);
  ue_mm_context_builder.add_mm_state(mm_state);
  ue_mm_context_builder.add_ecm_state(ecm_state);
  ue_mm_context_builder.add_emm_context(&emm_context);
  ue_mm_context_builder.add_sctp_assoc_id_key(sctp_assoc_id_key);
  ue_mm_context_builder.add_enb_ue_s1ap_id(enb_ue_s1ap_id);
  ue_mm_context_builder.add_enb_s1ap_id_key(enb_s1ap_id_key);
  ue_mm_context_builder.add_mme_ue_s1ap_id(mme_ue_s1ap_id);
  ue_mm_context_builder.add_attach_type(attach_type);
  ue_mm_context_builder.add_sgs_detach_type(sgs_detach_type);
  ue_mm_context_builder.add_e_utran_cgi(&e_utran_cgi);
  ue_mm_context_builder.add_cell_age(cell_age);
  ue_mm_context_builder.add_lai(&lai);
  ue_mm_context_builder.add_apn_config_profile(&apn_config_profile);
  ue_mm_context_builder.add_subscriber_status(subscriber_status);
  ue_mm_context_builder.add_network_access_mode(network_access_mode);
  ue_mm_context_builder.add_access_restriction_data(access_restriction_data);
  ue_mm_context_builder.add_apn_oi_replacement(&apn_oi_replacement);
  ue_mm_context_builder.add_mme_teid_s11(mme_teid_s11);
  ue_mm_context_builder.add_subscribed_ue_ambr(&subscribed_ue_ambr);
  ue_mm_context_builder.add_used_ue_ambr(&used_ue_ambr);
  ue_mm_context_builder.add_nb_active_pdn_contexts(nb_active_pdn_contexts);
  ue_mm_context_builder.add_pdn_contexts(&pdn_contexts);
  ue_mm_context_builder.add_bearer_contexts(&bearer_contexts);
  ue_mm_context_builder.add_ue_radio_capability(&ue_radio_capability);
  ue_mm_context_builder.add_send_ue_purge_request(send_ue_purge_request);
  ue_mm_context_builder.add_hss_initiated_detach(hss_initiated_detach);
  ue_mm_context_builder.add_location_info_confirmed_in_hss(
      location_info_confirmed_in_hss);
  ue_mm_context_builder.add_ppf(ppf);
  ue_mm_context_builder.add_subscription_known(subscription_known);
  ue_mm_context_builder.add_path_switch_req(path_switch_req);
  ue_mm_context_builder.add_granted_service(granted_service);
  ue_mm_context_builder.add_num_reg_sub(num_reg_sub);
  ue_mm_context_builder.add_reg_sub(&reg_sub);
  ue_mm_context_builder.add_cs_fallback_indicator(cs_fallback_indicator);
  ue_mm_context_builder.add_sgs_context(&sgs_context);
  ue_mm_context_builder.add_rau_tau_timer(rau_tau_timer);
  ue_mm_context_builder.add_time_mobile_reachability_timer_started(
      time_mobile_reachability_timer_started);
  ue_mm_context_builder.add_time_implicit_detach_timer_started(
      time_implicit_detach_timer_started);
  ue_mm_context_builder.add_time_paging_response_timer_started(
      time_paging_response_timer_started);
  ue_mm_context_builder.add_paging_retx_count(paging_retx_count);
  ue_mm_context_builder.add_time_ics_rsp_timer_started(
      time_ics_rsp_timer_started);

  flatbuffers::Offset<UeMmContext> o_ue_mm_context =
      ue_mm_context_builder.Finish();
  builder.Finish(o_ue_mm_context);
  uint8_t* buf                 = builder.GetBufferPointer();
  flatbuffers::uoffset_t size  = builder.GetSize();
  flatbuffers::uoffset_t align = 64;
  // TODO optimize '+1'
  flatbuffers::uoffset_t size_aligned = ((size / align) + 1) * align;
  // std::cout << "builder.GetSize() = " << builder.GetSize() << std::endl;
  // recopy buffer in another buffer we can manage
  uint8_t* buf_cp            = (uint8_t*) aligned_alloc(align, size_aligned);
  UeMmContext* ue_mm_context = nullptr;
  if (buf_cp) {
    memcpy((void*) buf_cp, buf, size);
    ue_mm_context = GetMutableUeMmContext(buf_cp);
  }
  // no need to copy size_aligned
  // but need to read with size_aligned
  std::tuple<uint8_t*, flatbuffers::uoffset_t, UeMmContext*> ret(
      buf_cp, size, ue_mm_context);
  return ret;
}

//------------------------------------------------------------------------------
std::vector<std::tuple<uint8_t*, flatbuffers::uoffset_t, UeMmContext*>>
mme_app_fb_allocate_ues(uint num_ues) {
  enb_s1ap_id_key_t enb_s1ap_id_key = INVALID_ENB_UE_S1AP_ID_KEY;
  enb_ue_s1ap_id_t enb_ue_s1ap_id   = rand() & 0X00FFFFFF;
  mme_ue_s1ap_id_t mme_ue_s1ap_id   = rand();
  std::vector<std::tuple<uint8_t*, flatbuffers::uoffset_t, UeMmContext*>>
      contexts;
  uint64_t kFirstImsi = 1010000000000;

  contexts.reserve(num_ues);

  for (int i = 0; i < num_ues; i++) {
    // UeMmContext
    std::tuple<uint8_t*, flatbuffers::uoffset_t, UeMmContext*> ue_mm_fb_ctx =
        build_ue_mm_context();
    // std::cout << "Built " << i << "'th UeMmContext context @"
    //          << (void*) std::get<2>(ue_mm_fb_ctx) << " buffer @"
    //          << (void*) std::get<0>(ue_mm_fb_ctx) << std::endl;

    enb_ue_s1ap_id++;
    std::get<2>(ue_mm_fb_ctx)->mutate_enb_ue_s1ap_id(enb_ue_s1ap_id);

    MME_APP_ENB_S1AP_ID_KEY(
        enb_s1ap_id_key, rand() & 0X0000FFFF, enb_ue_s1ap_id);
    std::get<2>(ue_mm_fb_ctx)->mutate_enb_s1ap_id_key(enb_s1ap_id_key);

    mme_ue_s1ap_id++;
    std::get<2>(ue_mm_fb_ctx)->mutate_mme_ue_s1ap_id(mme_ue_s1ap_id);

    // ESM context
    std::get<2>(ue_mm_fb_ctx)
        ->mutable_emm_context()
        ->mutable_esm_ctx()
        .mutate_n_active_ebrs(2);

    magma::lte::test_flat_buffer::EsmProcData& esm_proc_data =
        std::get<2>(ue_mm_fb_ctx)
            ->mutable_emm_context()
            ->mutable_esm_ctx()
            .mutable_esm_proc_data();

    esm_proc_data.mutable_bearer_qos().mutate_pci(true);
    esm_proc_data.mutable_bearer_qos().mutate_pl(15);
    esm_proc_data.mutable_bearer_qos().mutate_qci(5);
    int size = snprintf(
        (char*) esm_proc_data.mutable_apn().mutable_bytes()->data(),
        esm_proc_data.apn().bytes()->size(), "ims");
    if (size > 0) {
      esm_proc_data.mutable_apn().mutate_length(size);
    }
    esm_proc_data.mutate_pdn_cid(1);
    esm_proc_data.mutate_pdn_type(PdnTypeValue_IPv4);
    esm_proc_data.mutate_request_type(1);
    esm_proc_data.mutate_pti(1);

    // EMM context
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
    Imsi imsi_fb       = imsi2fb(imsi);
    Plmn plmn(0, 0, 1, 0, 1, 0xF);
    Gummei gummei(plmn, 1, 1);
    Guti guti(gummei, 2106150532 + i);
    Guti old_guti;
    old_guti.mutate_m_tmsi(429496729 + i);

    std::get<2>(ue_mm_fb_ctx)->mutable_emm_context()->mutate__imsi64(imsi64);
    std::get<2>(ue_mm_fb_ctx)->mutable_emm_context()->mutate__imsi64(imsi64);
    std::get<2>(ue_mm_fb_ctx)->mutable_emm_context()->mutable__imsi() = imsi_fb;
    std::get<2>(ue_mm_fb_ctx)->mutable_emm_context()->mutable__guti() = guti;
    std::get<2>(ue_mm_fb_ctx)->mutable_emm_context()->mutable__old_guti() =
        old_guti;
    std::get<2>(ue_mm_fb_ctx)
        ->mutable_emm_context()
        ->mutate_emm_cause(UINT32_MAX);
    std::get<2>(ue_mm_fb_ctx)
        ->mutable_emm_context()
        ->mutate__emm_fsm_state(EmmFsmState_EMM_REGISTERED);

    contexts.push_back(ue_mm_fb_ctx);
  }
  return contexts;
}
//------------------------------------------------------------------------------
void mme_app_fb_deallocate_ues(
    std::vector<std::tuple<uint8_t*, flatbuffers::uoffset_t, UeMmContext*>>&
        contexts) {
  // TODO
  contexts.clear();
}

//------------------------------------------------------------------------------
void mme_app_fb_serialize_ues(
    mme_app_desc_t* mme_app_desc,
    std::vector<std::tuple<uint8_t*, flatbuffers::uoffset_t, UeMmContext*>>&
        contexts,
    std::vector<uint64_t>& durations) {
  for (std::vector<std::tuple<uint8_t*, flatbuffers::uoffset_t, UeMmContext*>>::
           iterator it = contexts.begin();
       it != contexts.end(); ++it) {
    auto imsi_str = MmeNasStateManager::getInstance().get_imsi_str(
        std::get<2>(*it)->emm_context()->_imsi64());
    auto start = std::chrono::high_resolution_clock::now();
    MmeNasStateManager::getInstance().write_ue_state_to_db(
        reinterpret_cast<uint8_t*>(std::get<0>(*it)), std::get<1>(*it),
        imsi_str);
    auto stop = std::chrono::high_resolution_clock::now();
    auto duration =
        std::chrono::duration_cast<std::chrono::nanoseconds>(stop - start);
    durations.push_back(duration.count());
  }
}
//------------------------------------------------------------------------------
void mme_app_fb_insert_ues(
    std::vector<std::tuple<uint8_t*, flatbuffers::uoffset_t, UeMmContext*>>&
        contexts) {
  for (std::vector<std::tuple<uint8_t*, flatbuffers::uoffset_t, UeMmContext*>>::
           iterator it = contexts.begin();
       it != contexts.end(); ++it) {
    // TODO
  }
}

//------------------------------------------------------------------------------
void mme_app_fb_deserialize_ues(void) {
  // TODO
  // MmeNasStateManager::getInstance().read_fb_ue_state_from_db();
}

//------------------------------------------------------------------------------
void mme_app_fb_test_serialization(mme_app_desc_t* mme_app_desc, uint num_ues) {
  struct rusage ru_start_ctxt_to_proto, ru_end_ctxt_to_proto;

  std::cout << "sizeof(UeMmContext) : " << sizeof(struct UeMmContext)
            << std::endl;
  std::cout << "sizeof(EmmContext) : " << sizeof(struct EmmContext)
            << std::endl;
  std::cout << "sizeof(EsmContext) : " << sizeof(struct EsmContext)
            << std::endl;

  std::cout << "sizeof(ApnConfigProfile) : " << sizeof(struct ApnConfigProfile)
            << std::endl;
  std::cout << "sizeof(PdnContextArray) : " << sizeof(struct PdnContextArray)
            << std::endl;
  std::cout << "sizeof(BearerContextArray) : "
            << sizeof(struct BearerContextArray) << std::endl;
  std::cout << "sizeof(BearerContext) : " << sizeof(struct BearerContext)
            << std::endl;
  std::cout << "sizeof(EsmEbrContext) : " << sizeof(struct EsmEbrContext)
            << std::endl;
  std::cout << "sizeof(UeRadioCapability) : "
            << sizeof(struct UeRadioCapability) << std::endl;
  std::cout << "sizeof(RegionalSubscriptionArray) : "
            << sizeof(struct RegionalSubscriptionArray) << std::endl;
  std::cout << "sizeof(SgsContext) : " << sizeof(struct SgsContext)
            << std::endl;

  std::vector<uint64_t> durations;
  durations.reserve(num_ues);
  std::vector<std::tuple<uint8_t*, flatbuffers::uoffset_t, UeMmContext*>>
      contexts = mme_app_fb_allocate_ues(num_ues);

  mme_app_fb_insert_ues(contexts);

  getrusage(RUSAGE_SELF, &ru_start_ctxt_to_proto);
  auto start_ctxt_to_proto = std::chrono::high_resolution_clock::now();
  mme_app_fb_serialize_ues(mme_app_desc, contexts, durations);
  auto end_ctxt_to_proto = std::chrono::high_resolution_clock::now();
  getrusage(RUSAGE_SELF, &ru_end_ctxt_to_proto);
  log_rusage_diff(
      ru_start_ctxt_to_proto, ru_end_ctxt_to_proto,
      "RUSAGE Contexts serialization");
  auto duration_ctxt_to_proto =
      std::chrono::duration_cast<std::chrono::nanoseconds>(
          end_ctxt_to_proto - start_ctxt_to_proto);
  std::cout << "Time taken to serialize contexts: "
            << duration_ctxt_to_proto.count() << " nanoseconds" << std::endl;
  OAILOG_INFO(
      LOG_MME_APP, "Time taken to serialize contexts: %ld µs\n",
      duration_ctxt_to_proto.count());

  auto start_proto_to_ctxt = std::chrono::high_resolution_clock::now();
  mme_app_fb_deserialize_ues();
  auto end_proto_to_ctxt = std::chrono::high_resolution_clock::now();
  auto duration_proto_to_ctxt =
      std::chrono::duration_cast<std::chrono::nanoseconds>(
          end_proto_to_ctxt - start_proto_to_ctxt);
  std::cout << "Time taken to deserialize contexts: "
            << duration_proto_to_ctxt.count() << " nanoseconds" << std::endl;
  OAILOG_INFO(
      LOG_MME_APP, "Time taken to deserialize contexts:  %ld ns\n",
      duration_proto_to_ctxt.count());

  // dump context serialization history one by one
  std::cout << "CONTEXT SERIALIZATION HISTORY ";
  for (std::vector<uint64_t>::iterator it = durations.begin();
       it != durations.end(); ++it) {
    std::cout << (*it) << " ";
  }
  std::cout << std::endl;
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
  sleep(1);
  exit(0);
  return;
}
