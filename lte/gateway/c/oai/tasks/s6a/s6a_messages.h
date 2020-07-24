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

/*! \file s6a_messages.h
  \brief
  \author Sebastien ROUX
  \company Eurecom
*/

#ifndef S6A_MESSAGES_H_
#define S6A_MESSAGES_H_

#include "s6a_messages_types.h"

#include <freeDiameter/freeDiameter-host.h>
#include <freeDiameter/libfdproto.h>

int s6a_generate_update_location(s6a_update_location_req_t* ulr_p);
int s6a_generate_authentication_info_req(s6a_auth_info_req_t* uar_p);
int s6a_send_cancel_location_ans(s6a_cancel_location_ans_t* cla_pP);
int s6a_generate_purge_ue_req(const char* imsi);

int s6a_ula_cb(
    struct msg** msg, struct avp* paramavp, struct session* sess, void* opaque,
    enum disp_action* act);
int s6a_aia_cb(
    struct msg** msg, struct avp* paramavp, struct session* sess, void* opaque,
    enum disp_action* act);

int s6a_clr_cb(
    struct msg** msg, struct avp* paramavp, struct session* sess, void* opaque,
    enum disp_action* act);

int s6a_pua_cb(
    struct msg** msg, struct avp* paramavp, struct session* sess, void* opaque,
    enum disp_action* act);
int s6a_rsr_cb(
    struct msg** msg, struct avp* paramavp, struct session* sess, void* opaque,
    enum disp_action* act);

int s6a_parse_subscription_data(
    struct avp* avp_subscription_data, subscription_data_t* subscription_data);

#endif /* S6A_MESSAGES_H_ */
