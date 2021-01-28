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

/*! \file s6a_hss_reset.c
   \brief Handle a hss reset message and create the answer.
   \date 2018
   \version 0.1
*/

#include <stdio.h>

#include "assertions.h"
#include "intertask_interface.h"
#include "s6a_defs.h"
#include "s6a_messages.h"
#include "log.h"
#include "common_types.h"
#include "intertask_interface_types.h"
#include "itti_types.h"

struct avp;
struct msg;
struct session;

int s6a_rsr_cb(
    struct msg** msg_p, struct avp* paramavp_p, struct session* sess_p,
    void* opaque_p, enum disp_action* act_p) {
  struct msg* ans_p = NULL;
  struct msg* qry_p = NULL;
  struct avp *avp_p, *origin_host_p, *origin_realm_p;
  struct avp* failed_avp_p = NULL;
  struct avp_hdr* hdr_p    = NULL;
  struct avp_hdr *origin_host_hdr, *origin_realm_hdr;
  int result_code       = ER_DIAMETER_SUCCESS;
  int experimental      = 0;
  MessageDef* message_p = NULL;
  // s6a_reset_req_t                        *s6a_reset_req_p = NULL;

  DevAssert(msg_p);
  OAILOG_DEBUG(LOG_S6A, "Sending S6A_RESET_REQ to task MME_APP\n");

  qry_p = *msg_p;
  /*
   * Create the answer
   */
  CHECK_FCT(fd_msg_new_answer_from_req(fd_g_config->cnf_dict, msg_p, 0));
  ans_p = *msg_p;

  /*
   * Retrieving Origin host AVP
   */
  CHECK_FCT(fd_msg_search_avp(
      qry_p, s6a_fd_cnf.dataobj_s6a_origin_host, &origin_host_p));

  if (!origin_host_p) {
    OAILOG_ERROR(LOG_S6A, "origin_host ER_DIAMETER_MISSING_AVP\n");
    result_code = ER_DIAMETER_MISSING_AVP;
    goto out;
  }

  /*
   * Retrieving Origin realm AVP
   */
  CHECK_FCT(fd_msg_search_avp(
      qry_p, s6a_fd_cnf.dataobj_s6a_origin_realm, &origin_realm_p));

  if (!origin_realm_p) {
    OAILOG_ERROR(LOG_S6A, "origin_realm ER_DIAMETER_MISSING_AVP\n");
    result_code = ER_DIAMETER_MISSING_AVP;
    goto out;
  }

  /*
   * Retrieve the header from origin host and realm avps
   */
  CHECK_FCT(fd_msg_avp_hdr(origin_host_p, &origin_host_hdr));
  CHECK_FCT(fd_msg_avp_hdr(origin_realm_p, &origin_realm_hdr));
  // Send it to MME module for further processing
  message_p = itti_alloc_new_message(TASK_S6A, S6A_RESET_REQ);
  // s6a_reset_req_p->msg_rsa_p = msg_p;
  send_msg_to_task(&s6a_task_zmq_ctx, TASK_MME_APP, message_p);
  OAILOG_DEBUG(LOG_S6A, "Sending S6A_RESET_REQ to task MME_APP\n");
  result_code = DIAMETER_SUCCESS;
  return 0;
out:
  /*
   * Add the Auth-Session-State AVP
   */
  CHECK_FCT(fd_msg_search_avp(
      qry_p, s6a_fd_cnf.dataobj_s6a_auth_session_state, &avp_p));
  CHECK_FCT(fd_msg_avp_hdr(avp_p, &hdr_p));
  CHECK_FCT(
      fd_msg_avp_new(s6a_fd_cnf.dataobj_s6a_auth_session_state, 0, &avp_p));
  CHECK_FCT(fd_msg_avp_setvalue(avp_p, hdr_p->avp_value));
  CHECK_FCT(fd_msg_avp_add(ans_p, MSG_BRW_LAST_CHILD, avp_p));
  /*
   * Append the result code to the answer
   */
  CHECK_FCT(
      s6a_add_result_code(ans_p, failed_avp_p, result_code, experimental));
  CHECK_FCT(fd_msg_send(msg_p, NULL, NULL));
  return 0;
}
