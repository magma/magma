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

#include <gtest/gtest.h>

extern "C" {
#include "lte/gateway/c/core/oai/common/log.h"
}

#include "lte/gateway/c/core/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/emm_data.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/emm_proc.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/nas_state_converter.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/nas_procedures.hpp"

namespace magma {
namespace lte {

TEST(NasStateConverterTest, TestEmmContextConversion) {
  emm_context_t emm_context;

  emm_init_context(&emm_context, true);

  emm_context._imsi64 = 310150123456789;
  emm_context._imsi.u.num.digit1 = 3;
  emm_context._imsi.u.num.digit2 = 1;
  emm_context._imsi.u.num.digit3 = 0;
  emm_context._imsi.u.num.digit4 = 1;
  emm_context._imsi.u.num.digit5 = 5;
  emm_context._imsi.u.num.digit6 = 0;
  emm_context._imsi.u.num.digit7 = 1;
  emm_context._imsi.u.num.digit8 = 2;
  emm_context._imsi.u.num.digit9 = 3;
  emm_context._imsi.u.num.digit10 = 4;
  emm_context._imsi.u.num.digit11 = 5;
  emm_context._imsi.u.num.digit12 = 6;
  emm_context._imsi.u.num.digit13 = 7;
  emm_context._imsi.u.num.digit14 = 8;
  emm_context._imsi.u.num.digit15 = 9;
  emm_context.saved_imsi64 = 310150123456789;

  // Initialize EMM procedures
  nas_emm_attach_proc_t* attach_proc = nas_new_attach_procedure(&emm_context);
  nas_emm_auth_proc_t* auth_proc =
      nas_new_authentication_procedure(&emm_context);
  nas_emm_smc_proc_t* smc_proc = nas_new_smc_procedure(&emm_context);
  nas_emm_ident_proc_t* ident_proc =
      nas_new_identification_procedure(&emm_context);
  // TODO (ssanadhya): Add state for auth_info_proc

  // Initialize authentication vectors
  emm_context.remaining_vectors = 0;
  memset(emm_context._vector, sizeof(auth_vector_t), 1);

  emm_context.esm_ctx.esm_proc_data =
      (esm_proc_data_t*)calloc(1, sizeof(*emm_context.esm_ctx.esm_proc_data));
  emm_context.esm_ctx.esm_proc_data->pti = 5;
  bstring bstr = bfromcstr_with_str_len("192.168.0.1", 11);
  emm_context.esm_ctx.esm_proc_data->pdn_addr = bstr;
  bstring bstr_apn = bfromcstr_with_str_len("magma", 5);
  emm_context.esm_ctx.esm_proc_data->apn = bstr_apn;
  emm_context.esm_ctx.T3489.id = NAS_TIMER_INACTIVE_ID;

  emm_context._tai_list.numberoflists = 0;

  emm_context.new_attach_info =
      (new_attach_info_t*)calloc(1, sizeof(new_attach_info_t));
  emm_context.new_attach_info->mme_ue_s1ap_id = 1;
  emm_context.new_attach_info->is_mm_ctx_new = true;
  emm_context.new_attach_info->ies = (emm_attach_request_ies_t*)calloc(
      1, sizeof(*(emm_context.new_attach_info->ies)));
  ;
  emm_context.new_attach_info->ies->is_initial = true;
  emm_context.new_attach_info->ies->type = EMM_ATTACH_TYPE_EPS;

  oai::EmmContext proto_state;
  NasStateConverter::emm_context_to_proto(&emm_context, &proto_state);

  emm_context_t final_state;
  NasStateConverter::proto_to_emm_context(proto_state, &final_state);

  EXPECT_EQ(emm_context._imsi64, final_state._imsi64);

  EXPECT_STREQ((char*)emm_context.esm_ctx.esm_proc_data->pdn_addr->data,
               (char*)final_state.esm_ctx.esm_proc_data->pdn_addr->data);
  EXPECT_STREQ((char*)emm_context.esm_ctx.esm_proc_data->apn,
               (char*)final_state.esm_ctx.esm_proc_data->apn);

  EXPECT_TRUE(final_state.new_attach_info->ies->is_initial);
  EXPECT_EQ(final_state.new_attach_info->ies->type, EMM_ATTACH_TYPE_EPS);

  EXPECT_EQ(final_state.T3422.id, NAS_TIMER_INACTIVE_ID);

  // check that all procedures from initial state are in the final state
  EXPECT_TRUE(is_nas_specific_procedure_attach_running(&final_state));
  EXPECT_TRUE(is_nas_common_procedure_authentication_running(&final_state));
  EXPECT_TRUE(is_nas_common_procedure_smc_running(&final_state));
  // TODO (ssanadhya): Add check for Identification procedure, once state
  // conversion is implemented for it

  free_wrapper((void**)&emm_context.new_attach_info->ies);
  free_wrapper((void**)&emm_context.new_attach_info);
  clear_emm_ctxt(&emm_context);
  bdestroy_wrapper(&bstr);

  bdestroy_wrapper(&final_state.esm_ctx.esm_proc_data->pdn_addr);
  bdestroy_wrapper(&final_state.esm_ctx.esm_proc_data->apn);
  free_wrapper((void**)&final_state.new_attach_info->ies);
  free_wrapper((void**)&final_state.new_attach_info);
  free_wrapper((void**)&final_state.esm_ctx.esm_proc_data);
  clear_emm_ctxt(&final_state);
}

TEST(NasStateConverterTest, TestEsmContextSetInactiveT3489) {
  oai::EsmContext esm_context_proto;
  esm_context_t state_esm_context;

  esm_context_proto.clear_esm_proc_data();
  NasStateConverter::proto_to_esm_context(esm_context_proto,
                                          &state_esm_context);

  EXPECT_EQ(state_esm_context.T3489.id, NAS_TIMER_INACTIVE_ID);
}

TEST(NasStateConverterTest, TestEsmEbrContextSetInactiveTimer) {
  oai::EsmEbrContext esm_context_proto;
  esm_ebr_context_t state_esm_ebr_context;

  esm_context_proto.clear_pco();
  esm_context_proto.clear_esm_ebr_timer_data();
  NasStateConverter::proto_to_esm_ebr_context(esm_context_proto,
                                              &state_esm_ebr_context);

  EXPECT_EQ(state_esm_ebr_context.timer.id, NAS_TIMER_INACTIVE_ID);
}

TEST(NasStateConverterTest, TestNasEmmAttachProcSetInactiveT3450) {
  oai::AttachProc attach_proc_proto;
  nas_emm_attach_proc_t state_nas_emm_attach_proc;

  attach_proc_proto.clear_ies();
  attach_proc_proto.clear_esm_msg_out();
  magma::lte::NasStateConverter::proto_to_nas_emm_attach_proc(
      attach_proc_proto, &state_nas_emm_attach_proc);

  EXPECT_EQ(state_nas_emm_attach_proc.T3450.id, NAS_TIMER_INACTIVE_ID);
}

TEST(NasStateConverterTest, TestNasEmmAuthProcSetInactiveT3460) {
  oai::AuthProc auth_proc_proto;
  nas_emm_auth_proc_t state_nas_emm_auth_proc;

  NasStateConverter::proto_to_nas_emm_auth_proc(auth_proc_proto,
                                                &state_nas_emm_auth_proc);

  EXPECT_EQ(state_nas_emm_auth_proc.T3460.id, NAS_TIMER_INACTIVE_ID);
}

TEST(NasStateConverterTest, TestNasEmmSmcProcSetInactiveT3460) {
  oai::SmcProc smc_proc_proto;
  nas_emm_smc_proc_t state_nas_emm_smc_proc;

  magma::lte::NasStateConverter::proto_to_nas_emm_smc_proc(
      smc_proc_proto, &state_nas_emm_smc_proc);

  EXPECT_EQ(state_nas_emm_smc_proc.T3460.id, NAS_TIMER_INACTIVE_ID);
}

}  // namespace lte
}  // namespace magma

// Note: This is necessary for setting up a log thread (Might be addressed by
// #11736)
int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  OAILOG_INIT("MME", OAILOG_LEVEL_INFO, MAX_LOG_PROTOS);
  return RUN_ALL_TESTS();
}
