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

/*! \file s6a_cancel_loc.c
   \brief Handle an cancel location message and create the answer.
   \date 2017
   \version 0.1
*/

#include <stdio.h>
#include <string.h>
#include <conversions.h>

#include "assertions.h"
#include "intertask_interface.h"
#include "s6a_defs.h"
#include "s6a_messages.h"
#include "s6a_messages_types.h"
#include "log.h"
#include "common_types.h"
#include "intertask_interface_types.h"
#include "itti_types.h"

struct avp;
struct msg;
struct session;

#define IMSI_LENGTH 15
int s6a_clr_cb(
    struct msg** msg_p, struct avp* paramavp_p, struct session* sess_p,
    void* opaque_p, enum disp_action* act_p) {
  struct msg* ans_p = NULL;
  struct msg* qry_p = NULL;
  struct avp *avp_p, *origin_host_p, *origin_realm_p;
  struct avp* failed_avp_p = NULL;
  struct avp_hdr* hdr_p    = NULL;
  struct avp_hdr *origin_host_hdr, *origin_realm_hdr;
  int result_code                                      = ER_DIAMETER_SUCCESS;
  int experimental                                     = 0;
  MessageDef* message_p                                = NULL;
  s6a_cancel_location_req_t* s6a_cancel_location_req_p = NULL;
  int imsi_len                                         = 0;
  char imsi_str[IMSI_LENGTH + 1];

  DevAssert(msg_p);
  OAILOG_DEBUG(LOG_S6A, "Sending S6A_CANCEL_LOCATION_REQ to task MME_APP\n");

  qry_p = *msg_p;
  /*
   * Create the answer
   */
  CHECK_FCT(fd_msg_new_answer_from_req(fd_g_config->cnf_dict, msg_p, 0));
  ans_p = *msg_p;
  /*
   * Retrieving IMSI AVP
   */
  CHECK_FCT(fd_msg_search_avp(qry_p, s6a_fd_cnf.dataobj_s6a_user_name, &avp_p));

  if (avp_p) {
    CHECK_FCT(fd_msg_avp_hdr(avp_p, &hdr_p));

    if (hdr_p->avp_value->os.len > IMSI_LENGTH) {
      OAILOG_ERROR(
          LOG_S6A, "Received s6a clr for imsi=%*s\n",
          (int) hdr_p->avp_value->os.len, hdr_p->avp_value->os.data);
      OAILOG_ERROR(LOG_S6A, "IMSI_LENGTH ER_DIAMETER_INVALID_AVP_VALUE\n");
      result_code = ER_DIAMETER_INVALID_AVP_VALUE;
      goto out;
    }
    memcpy(imsi_str, hdr_p->avp_value->os.data, hdr_p->avp_value->os.len);
    imsi_len = hdr_p->avp_value->os.len;
  } else {
    OAILOG_ERROR(LOG_S6A, "Cannot get IMSI AVP which is mandatory\n");
    result_code = ER_DIAMETER_MISSING_AVP;
    goto out;
  }

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
  /*
   * Retrieving Cancellation type AVP
   */
  CHECK_FCT(fd_msg_search_avp(
      qry_p, s6a_fd_cnf.dataobj_s6a_cancellation_type, &avp_p));

  if (avp_p) {
    CHECK_FCT(fd_msg_avp_hdr(avp_p, &hdr_p));
    if (hdr_p->avp_value->u32 == SUBSCRIPTION_WITHDRAWL) {
      // Send it to MME module for further processing
      message_p = itti_alloc_new_message(TASK_S6A, S6A_CANCEL_LOCATION_REQ);
      s6a_cancel_location_req_p = &message_p->ittiMsg.s6a_cancel_location_req;
      memcpy(s6a_cancel_location_req_p->imsi, imsi_str, imsi_len);
      s6a_cancel_location_req_p->imsi[imsi_len]    = '\0';
      s6a_cancel_location_req_p->imsi_length       = imsi_len;
      s6a_cancel_location_req_p->cancellation_type = SUBSCRIPTION_WITHDRAWL;
      s6a_cancel_location_req_p->msg_cla_p         = msg_p;
      IMSI_STRING_TO_IMSI64(
          (char*) s6a_cancel_location_req_p->imsi,
          &message_p->ittiMsgHeader.imsi);
      send_msg_to_task(&s6a_task_zmq_ctx, TASK_MME_APP, message_p);
      OAILOG_DEBUG(
          LOG_S6A, "Sending S6A_CANCEL_LOCATION_REQ to task MME_APP\n");
      result_code = DIAMETER_SUCCESS;
    } else {
      OAILOG_ERROR(
          LOG_S6A, "S6A_CANCEL_LOCATION cancellation type %d not supported \n",
          hdr_p->avp_value->u32);
      result_code = ER_DIAMETER_INVALID_AVP_VALUE;
      goto out;
    }
  } else {
    OAILOG_ERROR(LOG_S6A, "S6A_CANCEL_LOCATION cancellation type missing \n");
    result_code = ER_DIAMETER_MISSING_AVP;
  }
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
int s6a_add_result_code(
    struct msg* ans, struct avp* failed_avp, int result_code,
    int experimental) {
  struct avp* avp;
  union avp_value value;

  if (DIAMETER_ERROR_IS_VENDOR(result_code) && experimental != 0) {
    struct avp* experimental_result;

    CHECK_FCT(fd_msg_avp_new(
        s6a_fd_cnf.dataobj_s6a_experimental_result, 0, &experimental_result));
    CHECK_FCT(fd_msg_avp_new(s6a_fd_cnf.dataobj_s6a_vendor_id, 0, &avp));
    value.u32 = VENDOR_3GPP;
    CHECK_FCT(fd_msg_avp_setvalue(avp, &value));
    CHECK_FCT(fd_msg_avp_add(experimental_result, MSG_BRW_LAST_CHILD, avp));
    CHECK_FCT(fd_msg_avp_new(
        s6a_fd_cnf.dataobj_s6a_experimental_result_code, 0, &avp));
    value.u32 = result_code;
    CHECK_FCT(fd_msg_avp_setvalue(avp, &value));
    CHECK_FCT(fd_msg_avp_add(experimental_result, MSG_BRW_LAST_CHILD, avp));
    CHECK_FCT(fd_msg_avp_add(ans, MSG_BRW_LAST_CHILD, experimental_result));
    /*
     * Add Origin_Host & Origin_Realm AVPs
     */
    CHECK_FCT(fd_msg_add_origin(ans, 0));
  } else {
    /*
     * This is a code defined in the base protocol: result-code AVP should
     * * * * be used.
     */
    CHECK_FCT(fd_msg_rescode_set(
        ans, retcode_2_string(result_code), NULL, failed_avp, 1));
  }

  return 0;
}

int s6a_send_cancel_location_ans(s6a_cancel_location_ans_t* cla_pP) {
  struct msg** msg_p       = NULL;
  struct msg* ans_p        = NULL;
  struct avp* failed_avp_p = NULL;
  int result_code          = 0;
  int experimental         = 0;

  DevAssert(cla_pP);

  OAILOG_DEBUG(LOG_S6A, "Received S6A_CANCEL_LOCATION_ANS from task MME_APP\n");

  result_code = cla_pP->result;

  msg_p = (struct msg**) cla_pP->msg_cla_p;
  if (msg_p == NULL) {
    return -1;
  }
  ans_p = *msg_p; /* Get the received CLA */
  /*
   * Append the result code to the answer
   */
  CHECK_FCT(
      s6a_add_result_code(ans_p, failed_avp_p, result_code, experimental));
  CHECK_FCT(fd_msg_send(msg_p, NULL, NULL));
  return 0;
}
