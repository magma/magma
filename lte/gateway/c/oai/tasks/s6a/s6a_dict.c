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

/*! \file s6a_dict.c
  \brief
  \author Sebastien ROUX
  \company Eurecom
*/

#if HAVE_CONFIG_H
#include "config.h"
#endif
#include <string.h>
#include <errno.h>

#include "common_defs.h"
#include "s6a_defs.h"
#include "s6a_messages.h"
#include "assertions.h"

#define CHECK_FD_FCT(fCT) DevAssert(fCT == 0);

/*! \file s6a_dict.c
   \brief Initialize s6a dictionnary and setup callbacks for procedures
   \author Sebastien ROUX <sebastien.roux@eurecom.fr>
   \date 2013
   \version 0.1
*/

int s6a_fd_init_dict_objs(void) {
  struct disp_when when;
  vendor_id_t vendor_3gpp  = VENDOR_3GPP;
  application_id_t app_s6a = APP_S6A;

  /*
   * Pre-loading vendor object
   */
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_VENDOR, VENDOR_BY_ID, (void*) &vendor_3gpp,
      &s6a_fd_cnf.dataobj_s6a_vendor, ENOENT));
  /*
   * Pre-loading application object
   */
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_APPLICATION, APPLICATION_BY_ID,
      (void*) &app_s6a, &s6a_fd_cnf.dataobj_s6a_app, ENOENT));
  /*
   * Pre-loading commands objects
   */
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_COMMAND, CMD_BY_NAME,
      "Authentication-Information-Request", &s6a_fd_cnf.dataobj_s6a_air,
      ENOENT));
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_COMMAND, CMD_BY_NAME,
      "Authentication-Information-Answer", &s6a_fd_cnf.dataobj_s6a_aia,
      ENOENT));
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_COMMAND, CMD_BY_NAME,
      "Update-Location-Request", &s6a_fd_cnf.dataobj_s6a_ulr, ENOENT));
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_COMMAND, CMD_BY_NAME,
      "Update-Location-Answer", &s6a_fd_cnf.dataobj_s6a_ula, ENOENT));
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_COMMAND, CMD_BY_NAME, "Purge-UE-Request",
      &s6a_fd_cnf.dataobj_s6a_pur, ENOENT));
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_COMMAND, CMD_BY_NAME, "Purge-UE-Answer",
      &s6a_fd_cnf.dataobj_s6a_pua, ENOENT));
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_COMMAND, CMD_BY_NAME,
      "Cancel-Location-Request", &s6a_fd_cnf.dataobj_s6a_clr, ENOENT));
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_COMMAND, CMD_BY_NAME,
      "Cancel-Location-Answer", &s6a_fd_cnf.dataobj_s6a_cla, ENOENT));
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_COMMAND, CMD_BY_NAME, "Reset-Request",
      &s6a_fd_cnf.dataobj_s6a_rsr, ENOENT));
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_COMMAND, CMD_BY_NAME, "Reset-Answer",
      &s6a_fd_cnf.dataobj_s6a_rsa, ENOENT));
  /*
   * Pre-loading base avps
   */
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_AVP, AVP_BY_NAME, "Destination-Host",
      &s6a_fd_cnf.dataobj_s6a_destination_host, ENOENT));
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_AVP, AVP_BY_NAME, "Destination-Realm",
      &s6a_fd_cnf.dataobj_s6a_destination_realm, ENOENT));
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_AVP, AVP_BY_NAME, "User-Name",
      &s6a_fd_cnf.dataobj_s6a_user_name, ENOENT));
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_AVP, AVP_BY_NAME, "Session-Id",
      &s6a_fd_cnf.dataobj_s6a_session_id, ENOENT));
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_AVP, AVP_BY_NAME, "Auth-Session-State",
      &s6a_fd_cnf.dataobj_s6a_auth_session_state, ENOENT));
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_AVP, AVP_BY_NAME, "Result-Code",
      &s6a_fd_cnf.dataobj_s6a_result_code, ENOENT));
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_AVP, AVP_BY_NAME, "Experimental-Result",
      &s6a_fd_cnf.dataobj_s6a_experimental_result, ENOENT));
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_AVP, AVP_BY_NAME, "Vendor-Id",
      &s6a_fd_cnf.dataobj_s6a_vendor_id, ENOENT));
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_AVP, AVP_BY_NAME, "Experimental-Result-Code",
      &s6a_fd_cnf.dataobj_s6a_experimental_result_code, ENOENT));
  /*
   * Pre-loading S6A specifics AVPs
   */
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_AVP, AVP_BY_NAME_ALL_VENDORS,
      "Visited-PLMN-Id", &s6a_fd_cnf.dataobj_s6a_visited_plmn_id, ENOENT));
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_AVP, AVP_BY_NAME_ALL_VENDORS, "RAT-Type",
      &s6a_fd_cnf.dataobj_s6a_rat_type, ENOENT));
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_AVP, AVP_BY_NAME_ALL_VENDORS, "ULR-Flags",
      &s6a_fd_cnf.dataobj_s6a_ulr_flags, ENOENT));
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_AVP, AVP_BY_NAME_ALL_VENDORS, "ULA-Flags",
      &s6a_fd_cnf.dataobj_s6a_ula_flags, ENOENT));
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_AVP, AVP_BY_NAME_ALL_VENDORS,
      "Subscription-Data", &s6a_fd_cnf.dataobj_s6a_subscription_data, ENOENT));
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_AVP, AVP_BY_NAME_ALL_VENDORS,
      "Requested-EUTRAN-Authentication-Info",
      &s6a_fd_cnf.dataobj_s6a_req_eutran_auth_info, ENOENT));
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_AVP, AVP_BY_NAME_ALL_VENDORS,
      "Number-Of-Requested-Vectors",
      &s6a_fd_cnf.dataobj_s6a_number_of_requested_vectors, ENOENT));
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_AVP, AVP_BY_NAME_ALL_VENDORS,
      "Immediate-Response-Preferred",
      &s6a_fd_cnf.dataobj_s6a_immediate_response_pref, ENOENT));
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_AVP, AVP_BY_NAME_ALL_VENDORS,
      "Authentication-Info", &s6a_fd_cnf.dataobj_s6a_authentication_info,
      ENOENT));
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_AVP, AVP_BY_NAME_ALL_VENDORS,
      "Re-Synchronization-Info",
      &s6a_fd_cnf.dataobj_s6a_re_synchronization_info, ENOENT));
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_AVP, AVP_BY_NAME_ALL_VENDORS,
      "Service-Selection", &s6a_fd_cnf.dataobj_s6a_service_selection, ENOENT));
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_AVP, AVP_BY_NAME_ALL_VENDORS,
      "UE-SRVCC-Capability", &s6a_fd_cnf.dataobj_s6a_ue_srvcc_cap, ENOENT));
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_AVP, AVP_BY_NAME_ALL_VENDORS,
      "Cancellation-Type", &s6a_fd_cnf.dataobj_s6a_cancellation_type, ENOENT));
  CHECK_FD_FCT(fd_dict_search(
      fd_g_config->cnf_dict, DICT_AVP, AVP_BY_NAME_ALL_VENDORS, "PUA-Flags",
      &s6a_fd_cnf.dataobj_s6a_pua_flags, ENOENT));
  /*
   * Register callbacks
   */
  memset(&when, 0, sizeof(when));
  when.command = s6a_fd_cnf.dataobj_s6a_ula;
  when.app     = s6a_fd_cnf.dataobj_s6a_app;
  /*
   * Register the callback for Update Location Answer S6A Application
   */
  CHECK_FD_FCT(fd_disp_register(
      s6a_ula_cb, DISP_HOW_CC, &when, NULL, &s6a_fd_cnf.ula_hdl));
  DevAssert(s6a_fd_cnf.ula_hdl);
  when.command = s6a_fd_cnf.dataobj_s6a_aia;
  when.app     = s6a_fd_cnf.dataobj_s6a_app;
  /*
   * Register the callback for Authentication Information Answer S6A Application
   */
  CHECK_FD_FCT(fd_disp_register(
      s6a_aia_cb, DISP_HOW_CC, &when, NULL, &s6a_fd_cnf.aia_hdl));
  DevAssert(s6a_fd_cnf.aia_hdl);

  /*
   * Register the callback for cancel Location Request S6A Application
   */
  when.command = s6a_fd_cnf.dataobj_s6a_clr;
  when.app     = s6a_fd_cnf.dataobj_s6a_app;
  CHECK_FD_FCT(fd_disp_register(
      s6a_clr_cb, DISP_HOW_CC, &when, NULL, &s6a_fd_cnf.clr_hdl));
  DevAssert(s6a_fd_cnf.clr_hdl);

  /*
   * Register the callback for Purge UE Answer S6A Application
   */
  when.command = s6a_fd_cnf.dataobj_s6a_pua;
  when.app     = s6a_fd_cnf.dataobj_s6a_app;
  CHECK_FD_FCT(fd_disp_register(
      s6a_pua_cb, DISP_HOW_CC, &when, NULL, &s6a_fd_cnf.pua_hdl));
  DevAssert(s6a_fd_cnf.pua_hdl);

  /*
   * Register the callback for hss Reset Request S6A Application
   */
  when.command = s6a_fd_cnf.dataobj_s6a_rsr;
  when.app     = s6a_fd_cnf.dataobj_s6a_app;
  CHECK_FD_FCT(fd_disp_register(
      s6a_rsr_cb, DISP_HOW_CC, &when, NULL, &s6a_fd_cnf.rsr_hdl));
  DevAssert(s6a_fd_cnf.rsr_hdl);
  /*

   * Advertise the support for the test application in the peer
   */
  CHECK_FD_FCT(fd_disp_app_support(
      s6a_fd_cnf.dataobj_s6a_app, s6a_fd_cnf.dataobj_s6a_vendor, 1, 0));
  return RETURNok;
}
