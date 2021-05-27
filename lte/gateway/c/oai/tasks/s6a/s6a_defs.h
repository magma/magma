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

/*! \file s6a_defs.h
  \brief
  \author Sebastien ROUX
  \company Eurecom
*/
#ifndef S6A_DEFS_H_
#define S6A_DEFS_H_

#if HAVE_CONFIG_H
#include "config.h"
#endif

#include <freeDiameter/freeDiameter-host.h>
#include <freeDiameter/libfdcore.h>

#include "mme_config.h"
#include "queue.h"
#include "intertask_interface.h"

extern task_zmq_ctx_t s6a_task_zmq_ctx;

#define VENDOR_3GPP (10415)
#define APP_S6A (16777251)

/* Errors that fall within the Permanent Failures category shall be used to
 * inform the peer that the request has failed, and should not be attempted
 * again. The Result-Code AVP values defined in Diameter Base Protocol RFC 3588
 * shall be applied. When one of the result codes defined here is included in a
 * response, it shall be inside an Experimental-Result AVP and the Result-Code
 * AVP shall be absent.
 */
#define DIAMETER_ERROR_USER_UNKNOWN (5001)
#define DIAMETER_ERROR_ROAMING_NOT_ALLOWED (5004)
#define DIAMETER_ERROR_UNKNOWN_EPS_SUBSCRIPTION (5420)
#define DIAMETER_ERROR_RAT_NOT_ALLOWED (5421)
#define DIAMETER_ERROR_EQUIPMENT_UNKNOWN (5422)
#define DIAMETER_ERROR_UNKOWN_SERVING_NODE (5423)

/* Result codes that fall within the transient failures category shall be used
 * to inform a peer that the request could not be satisfied at the time it was
 * received, but may be able to satisfy the request in the future. The
 * Result-Code AVP values defined in Diameter Base Protocol RFC 3588 shall be
 * applied. When one of the result codes defined here is included in a response,
 * it shall be inside an Experimental-Result AVP and the Result-Code AVP shall
 * be absent.
 */
#define DIAMETER_AUTHENTICATION_DATA_UNAVAILABLE (4181)

#define DIAMETER_ERROR_IS_VENDOR(x)                                            \
  ((x == DIAMETER_ERROR_USER_UNKNOWN) ||                                       \
   (x == DIAMETER_ERROR_ROAMING_NOT_ALLOWED) ||                                \
   (x == DIAMETER_ERROR_UNKNOWN_EPS_SUBSCRIPTION) ||                           \
   (x == DIAMETER_ERROR_RAT_NOT_ALLOWED) ||                                    \
   (x == DIAMETER_ERROR_EQUIPMENT_UNKNOWN) ||                                  \
   (x == DIAMETER_AUTHENTICATION_DATA_UNAVAILABLE) ||                          \
   (x == DIAMETER_ERROR_UNKOWN_SERVING_NODE))

typedef struct {
  struct dict_object* dataobj_s6a_vendor; /* s6a vendor object */
  struct dict_object* dataobj_s6a_app;    /* s6a application object */

  /* Commands */
  struct dict_object* dataobj_s6a_air; /* s6a authentication request */
  struct dict_object* dataobj_s6a_aia; /* s6a authentication answer */
  struct dict_object* dataobj_s6a_ulr; /* s6a update location request */
  struct dict_object* dataobj_s6a_ula; /* s6a update location asnwer */
  struct dict_object* dataobj_s6a_pur; /* s6a purge ue request */
  struct dict_object* dataobj_s6a_pua; /* s6a purge ue answer */
  struct dict_object* dataobj_s6a_clr; /* s6a Cancel Location req */
  struct dict_object* dataobj_s6a_cla; /* s6a Cancel Location ans */
  struct dict_object* dataobj_s6a_rsr; /* s6a Reset req */
  struct dict_object* dataobj_s6a_rsa; /* s6a Reset ans */

  /* Some standard basic AVPs */
  struct dict_object* dataobj_s6a_origin_host;
  struct dict_object* dataobj_s6a_origin_realm;
  struct dict_object* dataobj_s6a_destination_host;
  struct dict_object* dataobj_s6a_destination_realm;
  struct dict_object* dataobj_s6a_user_name;
  struct dict_object* dataobj_s6a_session_id;
  struct dict_object* dataobj_s6a_auth_session_state;
  struct dict_object* dataobj_s6a_result_code;
  struct dict_object* dataobj_s6a_experimental_result;
  struct dict_object* dataobj_s6a_vendor_id;
  struct dict_object* dataobj_s6a_experimental_result_code;

  /* S6A specific AVPs */
  struct dict_object* dataobj_s6a_visited_plmn_id;
  struct dict_object* dataobj_s6a_rat_type;
  struct dict_object* dataobj_s6a_ulr_flags;
  struct dict_object* dataobj_s6a_ula_flags;
  struct dict_object* dataobj_s6a_subscription_data;
  struct dict_object* dataobj_s6a_req_eutran_auth_info;
  struct dict_object* dataobj_s6a_number_of_requested_vectors;
  struct dict_object* dataobj_s6a_immediate_response_pref;
  struct dict_object* dataobj_s6a_authentication_info;
  struct dict_object* dataobj_s6a_re_synchronization_info;
  struct dict_object* dataobj_s6a_service_selection;
  struct dict_object* dataobj_s6a_ue_srvcc_cap;
  struct dict_object* dataobj_s6a_cancellation_type;
  struct dict_object* dataobj_s6a_pua_flags;
  struct dict_object* dataobj_s6a_supported_features;

  /* Handlers */
  struct disp_hdl* aia_hdl; /* Authentication Information Answer Handle */
  struct disp_hdl* ula_hdl; /* Update Location Answer Handle */
  struct disp_hdl* pua_hdl; /* Purge UE Answer Handle */
  struct disp_hdl* clr_hdl; /* Cancel Location Request Handle */
  struct disp_hdl* rsr_hdl; /* Hss Reset Request Handle */
} s6a_fd_cnf_t;

extern s6a_fd_cnf_t s6a_fd_cnf;

#define ULR_SINGLE_REGISTRATION_IND (1U)
#define ULR_S6A_S6D_INDICATOR (1U << 1)
#define ULR_SKIP_SUBSCRIBER_DATA (1U << 2)
#define ULR_GPRS_SUBSCRIPTION_DATA_IND (1U << 3)
#define ULR_NODE_TYPE_IND (1U << 4)
#define ULR_INITIAL_ATTACH_IND (1U << 5)
#define ULR_PS_LCS_SUPPORTED_BY_UE (1U << 6)
#define ULR_DUAL_REGIS_5G_IND (1U << 8)

#define ULA_SEPARATION_IND (1U)

#define PUA_FREEZE_M_TMSI (1U)
#define PUA_FREEZE_P_TMSI (1U << 1)

#define FLAG_IS_SET(x, flag) ((x) & (flag))

#define FLAGS_SET(x, flags) ((x) |= (flags))

#define FLAGS_CLEAR(x, flags) ((x) = (x) & ~(flags))

/* IANA defined IP address type */
#define IANA_IPV4 (0x1)
#define IANA_IPV6 (0x2)

#define AVP_CODE_3GPP_CHARGING_CHARACTERISTICS (13)
#define AVP_CODE_VENDOR_ID (266)
#define AVP_CODE_EXPERIMENTAL_RESULT (297)
#define AVP_CODE_EXPERIMENTAL_RESULT_CODE (298)
#define AVP_CODE_MIP_HOME_AGENT_ADDRESS (334)
#define AVP_CODE_MIP6_AGENT_INFO (486)
#define AVP_CODE_SERVICE_SELECTION (493)
#define AVP_CODE_BANDWIDTH_UL (516)
#define AVP_CODE_BANDWIDTH_DL (515)
#define AVP_CODE_SUPPORTED_FEATURES (628)
#define AVP_CODE_MSISDN (701)
#define AVP_CODE_SERVED_PARTY_IP_ADDRESS (848)
#define AVP_CODE_QCI (1028)
#define AVP_CODE_ALLOCATION_RETENTION_PRIORITY (1034)
#define AVP_CODE_PRIORITY_LEVEL (1046)
#define AVP_CODE_PRE_EMPTION_CAPABILITY (1047)
#define AVP_CODE_PRE_EMPTION_VULNERABILITY (1048)
#define AVP_CODE_SUBSCRIPTION_DATA (1400)
#define AVP_CODE_AUTHENTICATION_INFO (1413)
#define AVP_CODE_E_UTRAN_VECTOR (1414)
#define AVP_CODE_NETWORK_ACCESS_MODE (1417)
#define AVP_CODE_CANCELLATION_TYPE (1420)
#define AVP_CODE_CONTEXT_IDENTIFIER (1423)
#define AVP_CODE_SUBSCRIBER_STATUS (1424)
#define AVP_CODE_ACCESS_RESTRICTION_DATA (1426)
#define AVP_CODE_APN_OI_REPLACEMENT (1427)
#define AVP_CODE_ALL_APN_CONFIG_INC_IND (1428)
#define AVP_CODE_APN_CONFIGURATION_PROFILE (1429)
#define AVP_CODE_APN_CONFIGURATION (1430)
#define AVP_CODE_EPS_SUBSCRIBED_QOS_PROFILE (1431)
#define AVP_CODE_VPLMN_DYNAMIC_ADDRESS_ALLOWED (1432)
#define AVP_CODE_AMBR (1435)
#define AVP_CODE_PDN_GW_ALLOCATION_TYPE (1438)
#define AVP_CODE_REGIONAL_SUBSCRIPTION_ZONE_CODE (1446)
#define AVP_CODE_RAND (1447)
#define AVP_CODE_XRES (1448)
#define AVP_CODE_AUTN (1449)
#define AVP_CODE_KASME (1450)
#define AVP_CODE_PDN_TYPE (1456)
#define AVP_CODE_SUBSCRIBED_PERIODIC_RAU_TAU_TIMER (1619)

int s6a_init(const mme_config_t* mme_config);

int s6a_fd_new_peer(void);

void s6a_peer_connected_cb(struct peer_info* info, void* arg);

int s6a_fd_init_dict_objs(void);

int s6a_parse_subscription_data(
    struct avp* avp_subscription_data, subscription_data_t* subscription_data);

int s6a_parse_experimental_result(
    struct avp* avp, s6a_experimental_result_t* ptr);
char* experimental_retcode_2_string(uint32_t ret_code);
char* retcode_2_string(uint32_t ret_code);

int s6a_add_result_code(
    struct msg* ans, struct avp* failed_avp, int result_code, int experimental);

void send_activate_messages(void);

#endif /* S6A_DEFS_H_ */
