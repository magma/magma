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

#include "sgw_s8_utility.h"
extern "C" {
#include "pgw_procedures.h"
#include "common_types.h"
#include "mme_config.h"
#include "mme_app_embedded_spgw.h"
}
using ::testing::Test;
spgw_config_t spgw_config;

// TC validates updation of bearer context on reception of Create Session Rsp
TEST_F(SgwS8Config, check_dedicated_bearer_creation_request) {
  mme_sgw_tunnel_t sgw_s11_tunnel = {0};
  sgw_state_t* sgw_state          = create_ue_context(&sgw_s11_tunnel);
  sgw_eps_bearer_context_information_t* sgw_pdn_session = NULL;
  sgw_pdn_session = sgw_create_bearer_context_information_in_collection(
      sgw_s11_tunnel.local_teid);
  itti_s11_create_session_request_t session_req = {0};
  fill_itti_csreq(&session_req, default_eps_bearer_id);
  memcpy(session_req.apn, "internet", sizeof("internet"));
  sgw_s11_tunnel.remote_teid = session_req.sender_fteid_for_cp.teid;
  sgw_update_bearer_context_information_on_csreq(
      sgw_state, sgw_pdn_session, sgw_s11_tunnel, &session_req, imsi64);
  s8_create_bearer_request_t cb_req = {0};
  fill_create_bearer_request(
      &cb_req, sgw_s11_tunnel.local_teid, default_eps_bearer_id);

  itti_gx_nw_init_actv_bearer_request_t itti_bearer_req = {0};
  s8_bearer_context_t bc_cbreq = cb_req.bearer_context[0];

  itti_bearer_req.lbi = cb_req.linked_eps_bearer_id;

  memcpy(
      &itti_bearer_req.ul_tft, &bc_cbreq.tft, sizeof(traffic_flow_template_t));
  memcpy(
      &itti_bearer_req.dl_tft, &bc_cbreq.tft, sizeof(traffic_flow_template_t));
  memcpy(&itti_bearer_req.eps_bearer_qos, &bc_cbreq.qos, sizeof(bearer_qos_t));
  teid_t s1_u_sgw_fteid = sgw_get_new_s1u_teid(sgw_state);
  // Validates temporary bearer context is created
  EXPECT_EQ(
      create_temporary_dedicated_bearer_context(
          sgw_pdn_session, &itti_bearer_req,
          sgw_state->sgw_ip_address_S1u_S12_S4_up.s_addr, s1_u_sgw_fteid,
          cb_req.sequence_number, LOG_SGW_S8),
      RETURNok);
  EXPECT_EQ(
      update_pgw_info_to_temp_dedicated_bearer_context(
          sgw_pdn_session, s1_u_sgw_fteid, &bc_cbreq, sgw_state,
          cb_req.pgw_cp_address),
      RETURNok);
  // Validates sequence number matches with received create bearer request
  pgw_ni_cbr_proc_t* pgw_ni_cbr_proc =
      pgw_get_procedure_create_bearer(sgw_pdn_session);
  EXPECT_TRUE(pgw_ni_cbr_proc != nullptr);

  bool is_seq_number_updated                             = false;
  sgw_eps_bearer_entry_wrapper_t* sgw_eps_bearer_entry_p = nullptr;
  LIST_FOREACH(
      sgw_eps_bearer_entry_p, pgw_ni_cbr_proc->pending_eps_bearers, entries) {
    if ((sgw_eps_bearer_entry_p) &&
        (sgw_eps_bearer_entry_p->sgw_eps_bearer_entry->sgw_sequence_number ==
         cb_req.sequence_number)) {
      is_seq_number_updated = true;
      break;
    }
  }
  EXPECT_EQ(is_seq_number_updated, true);
  EXPECT_EQ(
      sgw_eps_bearer_entry_p->sgw_eps_bearer_entry->p_gw_teid_S5_S8_up,
      cb_req.bearer_context[0].pgw_s8_up.teid);
  free_wrapper((void**) &cb_req.pgw_cp_address);
  sgw_state_exit();
}

TEST_F(SgwS8Config, dedicated_bearer_invalid_lbi) {
  mme_sgw_tunnel_t sgw_s11_tunnel = {0};
  sgw_state_t* sgw_state          = create_ue_context(&sgw_s11_tunnel);
  sgw_eps_bearer_context_information_t* sgw_pdn_session = NULL;
  sgw_pdn_session = sgw_create_bearer_context_information_in_collection(
      sgw_s11_tunnel.local_teid);
  itti_s11_create_session_request_t session_req = {0};
  fill_itti_csreq(&session_req, default_eps_bearer_id);
  memcpy(session_req.apn, "internet", sizeof("internet"));
  sgw_s11_tunnel.remote_teid = session_req.sender_fteid_for_cp.teid;
  sgw_update_bearer_context_information_on_csreq(
      sgw_state, sgw_pdn_session, sgw_s11_tunnel, &session_req, imsi64);
  s8_create_bearer_request_t cb_req = {0};
  // send invalid default eps bearer id
  fill_create_bearer_request(
      &cb_req, sgw_s11_tunnel.local_teid, default_eps_bearer_id + 1);
  gtpv2c_cause_value_t cause_value = REQUEST_ACCEPTED;
  EXPECT_EQ(
      sgw_s8_handle_create_bearer_request(sgw_state, &cb_req, &cause_value),
      INVALID_IMSI64);
  free_wrapper((void**) &cb_req.pgw_cp_address);
  sgw_state_exit();
}

// TC validates updation of bearer context on reception of Create Session Rsp
TEST_F(SgwS8Config, check_dedicated_bearer_creation_response) {
  mme_sgw_tunnel_t sgw_s11_tunnel   = {0};
  sgw_state_t* sgw_state            = create_ue_context(&sgw_s11_tunnel);
  int argc                          = 5;
  spgw_config.pgw_config.enable_nat = false;
  spgw_config.sgw_config.ovs_config.uplink_port_num = 2;

  spgw_config.sgw_config.ovs_config.uplink_mac =
      bfromcstr_with_str_len("ff:ff:ff:ff:ff:ff", strlen("ff:ff:ff:ff:ff:ff"));
  spgw_config.sgw_config.ovs_config.gtp_port_num                  = 32768;
  spgw_config.sgw_config.ovs_config.mtr_port_num                  = 15577;
  spgw_config.sgw_config.ovs_config.internal_sampling_port_num    = 15578;
  spgw_config.sgw_config.ovs_config.internal_sampling_fwd_tbl_num = 201;
  sgw_initialize_gtpv1u();
  bdestroy_wrapper(&spgw_config.sgw_config.ovs_config.uplink_mac);
#if 0
  sgw_eps_bearer_context_information_t* sgw_pdn_session = NULL;
  sgw_pdn_session = sgw_create_bearer_context_information_in_collection(
      sgw_s11_tunnel.local_teid);
  itti_s11_create_session_request_t session_req = {0};
  fill_itti_csreq(&session_req, default_eps_bearer_id);
  memcpy(session_req.apn, "internet", sizeof("internet"));
  sgw_s11_tunnel.remote_teid = session_req.sender_fteid_for_cp.teid;
  sgw_update_bearer_context_information_on_csreq(
      sgw_state, sgw_pdn_session, sgw_s11_tunnel, &session_req, imsi64);
  s8_create_bearer_request_t cb_req = {0};
  fill_create_bearer_request(
      &cb_req, sgw_s11_tunnel.local_teid, default_eps_bearer_id);

  itti_gx_nw_init_actv_bearer_request_t itti_bearer_req = {0};
  s8_bearer_context_t bc_cbreq = cb_req.bearer_context[0];

  itti_bearer_req.lbi = cb_req.linked_eps_bearer_id;

  memcpy(
      &itti_bearer_req.ul_tft, &bc_cbreq.tft, sizeof(traffic_flow_template_t));
  memcpy(
      &itti_bearer_req.dl_tft, &bc_cbreq.tft, sizeof(traffic_flow_template_t));
  memcpy(&itti_bearer_req.eps_bearer_qos, &bc_cbreq.qos, sizeof(bearer_qos_t));
  teid_t s1_u_sgw_fteid = sgw_get_new_s1u_teid(sgw_state);
  create_temporary_dedicated_bearer_context(
      sgw_pdn_session, &itti_bearer_req,
      sgw_state->sgw_ip_address_S1u_S12_S4_up.s_addr, s1_u_sgw_fteid,
      cb_req.sequence_number, LOG_SGW_S8);
  update_pgw_info_to_temp_dedicated_bearer_context(
      sgw_pdn_session, s1_u_sgw_fteid, &bc_cbreq, sgw_state,
      cb_req.pgw_cp_address);

  itti_s11_nw_init_actv_bearer_rsp_t s11_actv_bearer_rsp;
  memset(&s11_actv_bearer_rsp, 0, sizeof(itti_s11_nw_init_actv_bearer_rsp_t));
  fill_create_bearer_response(
      &s11_actv_bearer_rsp, sgw_s11_tunnel.local_teid, default_eps_bearer_id,
      s1_u_sgw_fteid);
  sgw_s8_proc_s11_create_bearer_rsp(
      sgw_pdn_session, &s11_actv_bearer_rsp.bearer_contexts.bearer_contexts[0],
      &s11_actv_bearer_rsp, imsi64, sgw_state);
  EXPECT_EQ(sgw_pdn_session->pending_procedures, nullptr);

  free_wrapper((void**) &cb_req.pgw_cp_address);
#endif
  sgw_uninitialize_gtpv1u();
  sgw_state_exit();
}
