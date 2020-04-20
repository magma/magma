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

/*
 * Sends a SGS_LOCATION_UPDATE_ACCEPT message to MME.
 */
int handle_sgs_location_update_accept(
    const itti_sgsap_location_update_acc_t* itti_sgsap_location_update_acc_p);

/*
 * Sends a SGS_LOCATION_UPDATE_REJECT message to MME.
 */
int handle_sgs_location_update_reject(
    const itti_sgsap_location_update_rej_t* itti_sgsap_location_update_rej_p);

/*
 * Sends a SGS_EPS_DETACH_ACK message to MME.
 */
int handle_sgs_eps_detach_ack(
    const itti_sgsap_eps_detach_ack_t* sgs_eps_detach_ack_p);

/*
 * Sends a SGS_IMSI_DETACH_ACK message to MME.
 */
int handle_sgs_imsi_detach_ack(
    const itti_sgsap_imsi_detach_ack_t* sgs_imsi_detach_ack_p);

/*
 * Sends a SGS_MM_INFORMATION_REQUEST message to MME.
 */

int handle_sgs_mm_information_request(
    const itti_sgsap_mm_information_req_t* mm_information_req_pP);

/*
 * Sends a SGS_PAGING_REQUEST message to MME.
 */
int handle_sgs_paging_request(
    const itti_sgsap_paging_request_t* const sgs_paging_req_pP);

/*
 * Sends a SGS_VLR_RESET_INDICATION message to MME.
 */
int handle_sgs_vlr_reset_indication(
    const itti_sgsap_vlr_reset_indication_t* const sgs_vlr_reset_ind_pP);

/*
 * Sends a SGS_DOWNLINK_UNITDATA message to NAS.
 */
int handle_sgs_downlink_unitdata(
    const itti_sgsap_downlink_unitdata_t* sgs_dl_unitdata_p);
/*
 * Sends a SGS_RELEASE_REQUEST message to NAS.
 */
int handle_sgs_release_req(const itti_sgsap_release_req_t* sgs_release_req_p);

int handle_sgsap_alert_request(
    const itti_sgsap_alert_request_t* const sgsap_alert_request);

int handle_sgs_service_abort_req(
    const itti_sgsap_service_abort_req_t* const itti_sgsap_service_abort_req_p);
