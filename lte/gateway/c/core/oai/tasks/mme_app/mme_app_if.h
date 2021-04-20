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

#ifndef FILE_MME_APP_IF_SEEN
#define FILE_MME_APP_IF_SEEN

/*
 * IF method called by lower layers (S1AP) delivering the content of initial NAS
 * UE message
 */
int itf_mme_app_nas_initial_ue_message(
    const sctp_assoc_id_t assoc_id, const enb_ue_s1ap_id_t enb_ue_s1ap_id,
    const mme_ue_s1ap_id_t mme_ue_s1ap_id, const uint8_t* const nas_msg,
    const size_t nas_msg_length, const tai_t const* tai,
    const ecgi_t const* cgi, const long rrc_cause,
    const as_stmsi_t const* opt_s_tmsi, const void const* opt_csg_id,
    const gummei_t const* opt_gummei, const void const* opt_cell_access_mode,
    const void const* opt_cell_gw_transport_address,
    const void const* opt_relay_node_indicator);

#endif /* FILE_MME_APP_IF_SEEN */
