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

/*! \file s11_mme_bearer_manager.h
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#ifndef FILE_S11_MME_BEARER_MANAGER_SEEN
#define FILE_S11_MME_BEARER_MANAGER_SEEN

/* @brief Create a new Release Access Bearers Request and send it to provided
 * S-GW. */
int s11_mme_release_access_bearers_request(
    nw_gtpv2c_stack_handle_t* stack_p,
    itti_s11_release_access_bearers_request_t* release_access_bearers_p);

/* @brief Handle a Release Access Bearer Response received from S-GW. */
int s11_mme_handle_release_access_bearer_response(
    nw_gtpv2c_stack_handle_t* stack_p, nw_gtpv2c_ulp_api_t* pUlpApi);

/* @brief Handle a Modify Bearer Response received from S-GW. */
int s11_mme_handle_modify_bearer_response(
    nw_gtpv2c_stack_handle_t* stack_p, nw_gtpv2c_ulp_api_t* pUlpApi);

/* @brief Create a new Delete Bearer Command and send it to provided S-GW. */
int s11_mme_delete_bearer_command(
    nw_gtpv2c_stack_handle_t* stack_p, itti_s11_delete_bearer_command_t* cmd_p);

/* @brief Handle a Create Bearer Request received from S-GW. */
int s11_mme_handle_create_bearer_request(
    nw_gtpv2c_stack_handle_t* stack_p, nw_gtpv2c_ulp_api_t* pUlpApi);

/* @brief Create a new Create Bearer Response and send it to provided S-GW. */
int s11_mme_create_bearer_response(
    nw_gtpv2c_stack_handle_t* stack_p,
    itti_s11_create_bearer_response_t* rsp_p);

/* @brief Handle a DeleteBearer Request received from S-GW. */
int s11_mme_handle_delete_bearer_request(
    nw_gtpv2c_stack_handle_t* stack_p, nw_gtpv2c_ulp_api_t* pUlpApi);

/* @brief Handle a Downlink Data Notification received from S-GW. */
int s11_mme_handle_downlink_data_notification(
    nw_gtpv2c_stack_handle_t* stack_p, nw_gtpv2c_ulp_api_t* pUlpApi);

/* @brief Handle a Downlink Data Notification acknowledge received from S-GW. */
int s11_mme_downlink_data_notification_acknowledge(
    nw_gtpv2c_stack_handle_t* stack_p,
    itti_s11_downlink_data_notification_acknowledge_t* ack_p);

#endif /* FILE_S11_MME_BEARER_MANAGER_SEEN */
