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
#include <mcheck.h>
//--C++ includes ---------------------------------------------------------------
#include <chrono>
#include <cmath>
//--Other includes -------------------------------------------------------------
#include "mme_app_state_manager.h"
extern task_zmq_ctx_t main_zmq_ctx;

using magma::lte::MmeNasStateManager;

//------------------------------------------------------------------------------
void mme_app_schedule_test_protobuf_serialization(uint num_loops) {
  MessageDef* message_p =
      itti_alloc_new_message(TASK_UNKNOWN, MME_APP_TEST_PROTOBUF_SERIALIZATION);
  MME_APP_TEST_PROTOBUF_SERIALIZATION(message_p).num_loops = num_loops;
  send_msg_to_task(&main_zmq_ctx, TASK_MME_APP, message_p);
  return;
}
//------------------------------------------------------------------------------
void mme_app_test_protobuf_serialization(
    mme_app_desc_t* mme_app_desc, uint num_loops) {

  srand (time(NULL));
  ue_mm_context_t* ue_mm_contexts[num_loops][2] = {};
  enb_s1ap_id_key_t enb_s1ap_id_key = INVALID_ENB_UE_S1AP_ID_KEY;
  enb_ue_s1ap_id_t enb_ue_s1ap_id = rand() & 0X00FFFFFF;
  mme_ue_s1ap_id_t mme_ue_s1ap_id = rand();



  for (int i = 0; i < num_loops; i++) {
    ue_mm_context_t* ue_mm_context = mme_create_new_ue_context();
    emm_context_t* emm_ctx         = &ue_mm_context->emm_context;
    esm_context_t* esm_ctx         = &emm_ctx->esm_ctx;

    ue_mm_contexts[i][0] = ue_mm_context;
    enb_ue_s1ap_id++;
    mme_ue_s1ap_id++;

    imsi64_t imsi64       = 1010000000000 + i;
    imsi_t imsi           = {};
    imsi.u.num.digit1     = (uint8_t)((imsi64_t)(imsi64/std::pow(10,14))%10);
    imsi.u.num.digit2     = (uint8_t)((imsi64_t)(imsi64/std::pow(10,13))%10);
    imsi.u.num.digit3     = (uint8_t)((imsi64_t)(imsi64/std::pow(10,12))%10);
    imsi.u.num.digit4     = (uint8_t)((imsi64_t)(imsi64/std::pow(10,11))%10);
    imsi.u.num.digit5     = (uint8_t)((imsi64_t)(imsi64/std::pow(10,10))%10);
    imsi.u.num.digit6     = (uint8_t)((imsi64_t)(imsi64/std::pow(10,9))%10);
    imsi.u.num.digit7     = (uint8_t)((imsi64_t)(imsi64/std::pow(10,8))%10);
    imsi.u.num.digit8     = (uint8_t)((imsi64_t)(imsi64/std::pow(10,7))%10);
    imsi.u.num.digit9     = (uint8_t)((imsi64_t)(imsi64/std::pow(10,6))%10);
    imsi.u.num.digit10    = (uint8_t)((imsi64_t)(imsi64/std::pow(10,5))%10);
    imsi.u.num.digit11    = (uint8_t)((imsi64_t)(imsi64/std::pow(10,4))%10);
    imsi.u.num.digit12    = (uint8_t)((imsi64_t)(imsi64/std::pow(10,3))%10);
    imsi.u.num.digit13    = (uint8_t)((imsi64_t)(imsi64/std::pow(10,2))%10);
    imsi.u.num.digit14    = (uint8_t)((imsi64_t)(imsi64/std::pow(10,1))%10);
    imsi.u.num.digit15    = (uint8_t)((imsi64_t)(imsi64/std::pow(10,0))%10);
    imsi.u.num.parity     = 0xF;
    emm_ctx->saved_imsi64 = imsi64;

    guti_t guti                     = {};
    guti.gummei.plmn.mcc_digit2     = 0;
    guti.gummei.plmn.mcc_digit1     = 0;
    guti.gummei.plmn.mnc_digit3     = 1;
    guti.gummei.plmn.mcc_digit3     = 0xF;
    guti.gummei.plmn.mnc_digit2     = 1;
    guti.gummei.plmn.mnc_digit1     = 0;
    guti.gummei.mme_gid             = 1;
    guti.gummei.mme_code            = 1;
    guti.m_tmsi                     = 2106150532 + i;
    guti_t old_guti                 = {};
    old_guti.gummei.plmn.mcc_digit2 = 0;
    old_guti.gummei.plmn.mcc_digit1 = 0;
    old_guti.gummei.plmn.mnc_digit3 = 0;
    old_guti.gummei.plmn.mcc_digit3 = 0;
    old_guti.gummei.plmn.mnc_digit2 = 0;
    old_guti.gummei.plmn.mnc_digit1 = 0;
    old_guti.gummei.mme_gid         = 0;
    old_guti.gummei.mme_code        = 0;
    old_guti.m_tmsi                 = 429496729 + i;

    emm_ctx_set_valid_imsi(emm_ctx, &imsi, imsi64);
    emm_ctx_set_valid_guti(emm_ctx, &guti);
    emm_ctx_set_valid_old_guti(emm_ctx, &old_guti);

    emm_ctx->emm_cause      = -1;
    emm_ctx->_emm_fsm_state = EMM_REGISTERED;

    esm_ctx->n_active_ebrs = 2;
    esm_ctx->esm_proc_data =
        (struct esm_proc_data_s*) calloc(1, sizeof(struct esm_proc_data_s));
    struct esm_proc_data_s* esm_proc_data = esm_ctx->esm_proc_data;
    esm_proc_data->pti                    = 1;
    esm_proc_data->request_type           = 1;
    esm_proc_data->apn                    = bfromcstr("ims");
    esm_proc_data->pdn_cid                = 1;
    esm_proc_data->pdn_type               = ESM_PDN_TYPE_IPV4;
    esm_proc_data->bearer_qos.pci         = 1;
    esm_proc_data->bearer_qos.pl          = 15;
    esm_proc_data->bearer_qos.qci         = 5;
    // no bearer_qos.gbr, bearer_qos.mbr
    // no pco


    MME_APP_ENB_S1AP_ID_KEY(
        ue_mm_context->enb_s1ap_id_key, rand() & 0X0000FFFF, enb_ue_s1ap_id);
    ue_mm_context->enb_ue_s1ap_id = enb_ue_s1ap_id;
    ue_mm_context->mme_ue_s1ap_id = mme_ue_s1ap_id;

    if (mme_insert_ue_context(&mme_app_desc->mme_ue_contexts, ue_mm_context) !=
        RETURNok) {
      OAILOG_ERROR_UE(
          LOG_MME_APP, imsi64,
          "Failed to insert UE contxt, MME UE S1AP Id: " MME_UE_S1AP_ID_FMT "\n",
          ue_mm_context->mme_ue_s1ap_id);
      //return;
    }

    auto start_ctxt_to_proto = std::chrono::high_resolution_clock::now();
    //mtrace();

    put_mme_ue_state(mme_app_desc, imsi64, true);

    auto end_ctxt_to_proto = std::chrono::high_resolution_clock::now();
    auto duration_ctxt_to_proto =
        std::chrono::duration_cast<std::chrono::microseconds>(
            end_ctxt_to_proto - start_ctxt_to_proto);
    std::cout << "Time taken by context to proto conversion : "
              << duration_ctxt_to_proto.count() << " microseconds" << std::endl;
    OAILOG_INFO_UE(
        LOG_MME_APP, imsi64,
        "Time taken by context to proto conversion : %ld µs\n",
        duration_ctxt_to_proto.count());

    ue_mm_contexts[i][1] = mme_create_new_ue_context();
    emm_context_t* emm_ctx2         = &ue_mm_contexts[i][1]->emm_context;

    auto start_proto_to_ctxt = std::chrono::high_resolution_clock::now();

    mme_app_desc_t* mme_app_desc2 = get_mme_nas_state(true);
    MmeNasStateManager::getInstance().read_ue_state_from_db();
    //muntrace();

    auto end_proto_to_ctxt = std::chrono::high_resolution_clock::now();
    auto duration_proto_to_ctxt =
        std::chrono::duration_cast<std::chrono::microseconds>(
            end_proto_to_ctxt - start_proto_to_ctxt);
    std::cout << "Time taken by proto to context conversion : "
              << duration_proto_to_ctxt.count() << " microseconds" << std::endl;
    OAILOG_INFO_UE(
        LOG_MME_APP, imsi64,
        "Time taken by proto to context conversion : %ld µs\n",
        duration_proto_to_ctxt.count());
    /*auto imsi_str = MmeNasStateManager::getInstance().get_imsi_str(imsi64);
    MmeNasStateManager::getInstance().write_ue_state_to_db(
        ue_context, imsi_str);
    put_mme_ue_state(mme_app_desc_p, imsi64, force_ue_write);
    put_mme_nas_state(); */
  }
  mme_app_desc_t* mme_app_desc_p = get_mme_nas_state(false);
  if (!mme_app_desc_p) {
    OAILOG_ERROR(LOG_MME_APP, "Failed to fetch mme_app_desc_p \n");
    return;
  }
  for (int r = 0; r < num_loops; r++) {
    mme_remove_ue_context(&mme_app_desc_p->mme_ue_contexts, ue_mm_contexts[r][0]);
  }

  send_terminate_message_fatal(&main_zmq_ctx);
  return;
}
