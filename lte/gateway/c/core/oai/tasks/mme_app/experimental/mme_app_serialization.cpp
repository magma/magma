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
// --C includes ---------------------------------------------------------------
#include "lte/gateway/c/core/oai/tasks/mme_app/experimental/mme_app_serialization.hpp"
#include <mcheck.h>
#include <sys/time.h>      // rusage()
#include <sys/resource.h>  // rusage()
// --C++ includes -------------------------------------------------------------
#include <chrono>
#include <cmath>
#include <cstdlib>
#include <vector>
// --Other includes -----------------------------------------------------------
extern "C" {
#include "lte/gateway/c/core/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/include/mme_app_state.hpp"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_23.003.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#include "lte/gateway/c/core/oai/lib/itti/itti_types.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/emm_data.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/emm_proc.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/esm_proc.h"
#include "lte/gateway/c/core/oai/tasks/nas/nas_procedures.h"
}

#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_state_manager.hpp"
extern task_zmq_ctx_t main_zmq_ctx;

using magma::lte::MmeNasStateManager;

uint64_t kFirstImsi = 1010000000000;

std::vector<ue_mm_context_t*> mme_app_allocate_ues(uint num_ues);
void mme_app_deallocate_ues(mme_app_desc_t* mme_app_desc,
                            std::vector<ue_mm_context_t*>* contexts);
void mme_app_insert_ues(mme_app_desc_t* mme_app_desc,
                        const std::vector<ue_mm_context_t*>& contexts);
void mme_app_serialize_ues(const std::vector<ue_mm_context_t*>& contexts);
void mme_app_deserialize_ues(void);

void mme_app_schedule_test_protobuf_serialization(uint num_ues) {
  MessageDef* message_p =
      itti_alloc_new_message(TASK_UNKNOWN, MME_APP_TEST_PROTOBUF_SERIALIZATION);
  MME_APP_TEST_PROTOBUF_SERIALIZATION(message_p).num_ues = num_ues;
  send_msg_to_task(&main_zmq_ctx, TASK_MME_APP, message_p);
  return;
}

std::vector<ue_mm_context_t*> mme_app_allocate_ues(uint num_ues) {
  enb_s1ap_id_key_t enb_s1ap_id_key = INVALID_ENB_UE_S1AP_ID_KEY;
  unsigned int seed = 0;
  enb_ue_s1ap_id_t enb_ue_s1ap_id = rand_r(&seed) & 0X00FFFFFF;
  mme_ue_s1ap_id_t mme_ue_s1ap_id = rand_r(&seed);
  std::vector<ue_mm_context_t*> contexts;

  contexts.reserve(num_ues);

  for (int i = 0; i < num_ues; i++) {
    ue_mm_context_t* ue_mm_context = mme_create_new_ue_context();
    emm_context_t* emm_ctx = &ue_mm_context->emm_context;
    esm_context_t* esm_ctx = &emm_ctx->esm_ctx;

    enb_ue_s1ap_id++;
    mme_ue_s1ap_id++;

    imsi64_t imsi64 = kFirstImsi + i;
    imsi_t imsi = {};
    imsi.u.num.digit1 = (uint8_t)((imsi64_t)(imsi64 / std::pow(10, 14)) % 10);
    imsi.u.num.digit2 = (uint8_t)((imsi64_t)(imsi64 / std::pow(10, 13)) % 10);
    imsi.u.num.digit3 = (uint8_t)((imsi64_t)(imsi64 / std::pow(10, 12)) % 10);
    imsi.u.num.digit4 = (uint8_t)((imsi64_t)(imsi64 / std::pow(10, 11)) % 10);
    imsi.u.num.digit5 = (uint8_t)((imsi64_t)(imsi64 / std::pow(10, 10)) % 10);
    imsi.u.num.digit6 = (uint8_t)((imsi64_t)(imsi64 / std::pow(10, 9)) % 10);
    imsi.u.num.digit7 = (uint8_t)((imsi64_t)(imsi64 / std::pow(10, 8)) % 10);
    imsi.u.num.digit8 = (uint8_t)((imsi64_t)(imsi64 / std::pow(10, 7)) % 10);
    imsi.u.num.digit9 = (uint8_t)((imsi64_t)(imsi64 / std::pow(10, 6)) % 10);
    imsi.u.num.digit10 = (uint8_t)((imsi64_t)(imsi64 / std::pow(10, 5)) % 10);
    imsi.u.num.digit11 = (uint8_t)((imsi64_t)(imsi64 / std::pow(10, 4)) % 10);
    imsi.u.num.digit12 = (uint8_t)((imsi64_t)(imsi64 / std::pow(10, 3)) % 10);
    imsi.u.num.digit13 = (uint8_t)((imsi64_t)(imsi64 / std::pow(10, 2)) % 10);
    imsi.u.num.digit14 = (uint8_t)((imsi64_t)(imsi64 / std::pow(10, 1)) % 10);
    imsi.u.num.digit15 = (uint8_t)((imsi64_t)(imsi64 / std::pow(10, 0)) % 10);
    imsi.u.num.parity = 0xF;
    emm_ctx->saved_imsi64 = imsi64;

    guti_t guti = {};
    guti.gummei.plmn.mcc_digit2 = 0;
    guti.gummei.plmn.mcc_digit1 = 0;
    guti.gummei.plmn.mnc_digit3 = 1;
    guti.gummei.plmn.mcc_digit3 = 0xF;
    guti.gummei.plmn.mnc_digit2 = 1;
    guti.gummei.plmn.mnc_digit1 = 0;
    guti.gummei.mme_gid = 1;
    guti.gummei.mme_code = 1;
    guti.m_tmsi = 2106150532 + i;
    guti_t old_guti = {};
    old_guti.gummei.plmn.mcc_digit2 = 0;
    old_guti.gummei.plmn.mcc_digit1 = 0;
    old_guti.gummei.plmn.mnc_digit3 = 0;
    old_guti.gummei.plmn.mcc_digit3 = 0;
    old_guti.gummei.plmn.mnc_digit2 = 0;
    old_guti.gummei.plmn.mnc_digit1 = 0;
    old_guti.gummei.mme_gid = 0;
    old_guti.gummei.mme_code = 0;
    old_guti.m_tmsi = 429496729 + i;

    emm_ctx_set_valid_imsi(emm_ctx, &imsi, imsi64);
    emm_ctx_set_valid_guti(emm_ctx, &guti);
    emm_ctx_set_valid_old_guti(emm_ctx, &old_guti);

    emm_ctx->emm_cause = -1;
    emm_ctx->_emm_fsm_state = EMM_REGISTERED;

    esm_ctx->n_active_ebrs = 2;
    esm_ctx->esm_proc_data =
        (struct esm_proc_data_s*)calloc(1, sizeof(struct esm_proc_data_s));
    struct esm_proc_data_s* esm_proc_data = esm_ctx->esm_proc_data;
    esm_proc_data->pti = 1;
    esm_proc_data->request_type = 1;
    esm_proc_data->apn = bfromcstr("ims");
    esm_proc_data->pdn_cid = 1;
    esm_proc_data->pdn_type = ESM_PDN_TYPE_IPV4;
    esm_proc_data->bearer_qos.pci = 1;
    esm_proc_data->bearer_qos.pl = 15;
    esm_proc_data->bearer_qos.qci = 5;
    // no bearer_qos.gbr, bearer_qos.mbr
    // no pco

    MME_APP_ENB_S1AP_ID_KEY(ue_mm_context->enb_s1ap_id_key, rand() & 0X0000FFFF,
                            enb_ue_s1ap_id);
    ue_mm_context->enb_ue_s1ap_id = enb_ue_s1ap_id;
    ue_mm_context->mme_ue_s1ap_id = mme_ue_s1ap_id;

    contexts.push_back(ue_mm_context);
  }
  return contexts;
}

void mme_app_deallocate_ues(mme_app_desc_t* mme_app_desc,
                            std::vector<ue_mm_context_t*>* contexts) {
  while (!contexts.empty()) {
    mme_remove_ue_context(&mme_app_desc->mme_ue_contexts, contexts.back());
    contexts.pop_back();
  }
}

void mme_app_serialize_ues(mme_app_desc_t* mme_app_desc,
                           const std::vector<ue_mm_context_t*>& contexts) {
  for (auto it = contexts.begin(); it != contexts.end(); ++it) {
    put_mme_ue_state(mme_app_desc, (*it)->emm_context.saved_imsi64, true);
  }
}

void mme_app_insert_ues(mme_app_desc_t* mme_app_desc,
                        const std::vector<ue_mm_context_t*>& contexts) {
  for (auto it = contexts.begin(); it != contexts.end(); ++it) {
    if (mme_insert_ue_context(&mme_app_desc->mme_ue_contexts, *it) !=
        RETURNok) {
      OAILOG_ERROR_UE(
          LOG_MME_APP, (*it)->emm_context.saved_imsi64,
          "Failed to insert UE contxt, MME UE S1AP Id: " MME_UE_S1AP_ID_FMT
          "\n",
          (*it)->mme_ue_s1ap_id);
      return;
    }
  }
}

void mme_app_deserialize_ues(void) {
  mme_app_desc_t* mme_app_desc2 = get_mme_nas_state(true);
  MmeNasStateManager::getInstance().read_ue_state_from_db();
}

void log_rusage(const struct rusage& ru, const char* context) {
  std::cout << context
            << "\tCpu user/system: " << static_cast<int>(ru.ru_utime.tv_sec)
            << "." << static_cast<int>(ru.ru_utime.tv_usec) << " / "
            << static_cast<int>(ru.ru_stime.tv_sec) << "."
            << static_cast<int>(ru.ru_stime.tv_usec) << std::endl;
  std::cout << context << "\tMaximum resident set size: " << ru.ru_maxrss
            << std::endl;
  std::cout << context
            << "\tPage reclaims (soft/hard page faults): " << ru.ru_minflt
            << " / " << ru.ru_majflt << std::endl;
  std::cout << context << "\tBlock operations (input/output): " << ru.ru_inblock
            << " / " << ru.ru_oublock << std::endl;
  std::cout << context
            << "\tContext switches (voluntary/involuntary): " << ru.ru_nvcsw
            << " / " << ru.ru_nivcsw << std::endl;
}

void log_rusage_diff(const struct rusage& ru_first,
                     const struct rusage& ru_last, const char* context) {
  struct rusage ru_diff = {0};
  ru_diff.ru_utime.tv_sec = ru_last.ru_utime.tv_sec - ru_first.ru_utime.tv_sec;
  ru_diff.ru_utime.tv_usec =
      ru_last.ru_utime.tv_usec - ru_first.ru_utime.tv_usec;
  ru_diff.ru_stime.tv_sec = ru_last.ru_stime.tv_sec - ru_first.ru_stime.tv_sec;
  ru_diff.ru_stime.tv_usec =
      ru_last.ru_stime.tv_usec - ru_first.ru_stime.tv_usec;

  ru_diff.ru_maxrss = ru_last.ru_maxrss - ru_first.ru_maxrss;
  ru_diff.ru_minflt = ru_last.ru_minflt - ru_first.ru_minflt;
  ru_diff.ru_majflt = ru_last.ru_majflt - ru_first.ru_majflt;
  ru_diff.ru_inblock = ru_last.ru_inblock - ru_first.ru_inblock;
  ru_diff.ru_oublock = ru_last.ru_oublock - ru_first.ru_oublock;
  ru_diff.ru_nvcsw = ru_last.ru_nvcsw - ru_first.ru_nvcsw;
  ru_diff.ru_nivcsw = ru_last.ru_nivcsw - ru_first.ru_nivcsw;
  std::string str3(context);
  log_rusage(ru_diff, str3.c_str());
}

void mme_app_test_protobuf_serialization(mme_app_desc_t* mme_app_desc,
                                         uint num_ues) {
  srand(time(NULL));
  ue_mm_context_t* ue_mm_contexts[num_ues][2] = {};

  std::vector<ue_mm_context_t*> contexts = mme_app_allocate_ues(num_ues);

  mme_app_insert_ues(mme_app_desc, contexts);

  struct rusage ru_start_ctxt_to_proto, ru_end_ctxt_to_proto;

  getrusage(RUSAGE_SELF, &ru_start_ctxt_to_proto);
  auto start_ctxt_to_proto = std::chrono::high_resolution_clock::now();
  mme_app_serialize_ues(mme_app_desc, contexts);
  auto end_ctxt_to_proto = std::chrono::high_resolution_clock::now();
  getrusage(RUSAGE_SELF, &ru_end_ctxt_to_proto);
  log_rusage_diff(ru_start_ctxt_to_proto, ru_end_ctxt_to_proto,
                  "RUSAGE Contexts serialization");
  auto duration_ctxt_to_proto =
      std::chrono::duration_cast<std::chrono::nanoseconds>(end_ctxt_to_proto -
                                                           start_ctxt_to_proto);
  std::cout << "Time taken to serialize contexts " << num_ues
            << " UEs : " << duration_ctxt_to_proto.count() << " nanoseconds"
            << std::endl;
  OAILOG_INFO(LOG_MME_APP, "Time taken to serialize contexts %d UEs: %ld µs\n",
              num_ues, duration_ctxt_to_proto.count());

  struct rusage ru_start_proto_to_ctxt, ru_end_proto_to_ctxt;
  getrusage(RUSAGE_SELF, &ru_start_proto_to_ctxt);
  auto start_proto_to_ctxt = std::chrono::high_resolution_clock::now();
  mme_app_deserialize_ues();
  auto end_proto_to_ctxt = std::chrono::high_resolution_clock::now();
  getrusage(RUSAGE_SELF, &ru_end_proto_to_ctxt);
  log_rusage_diff(ru_start_proto_to_ctxt, ru_end_proto_to_ctxt,
                  "RUSAGE Contexts deserialization");
  auto duration_proto_to_ctxt =
      std::chrono::duration_cast<std::chrono::nanoseconds>(end_proto_to_ctxt -
                                                           start_proto_to_ctxt);
  std::cout << "Time taken to deserialize contexts " << num_ues
            << " UEs : " << duration_proto_to_ctxt.count() << " nanoseconds"
            << std::endl;
  OAILOG_INFO(LOG_MME_APP, "Time taken to serialize contexts %d UEs : %ld ns\n",
              num_ues, duration_proto_to_ctxt.count());
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
  mme_app_deallocate_ues(mme_app_desc_p, &contexts);

  send_terminate_message_fatal(&main_zmq_ctx);
  sleep(1);
  exit(0);
  return;
}
