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
extern "C" {}

namespace magma {
namespace lte {

void send_sctp_mme_server_initialized();

void send_activate_message_to_mme_app();

void send_mme_app_initial_ue_msg(
    const uint8_t* nas_msg, uint8_t nas_msg_length, const plmn_t& plmn);

void send_mme_app_uplink_data_ind(
    const uint8_t* nas_msg, uint8_t nas_msg_length, const plmn_t& plmn);

void send_authentication_info_resp(const std::string& imsi);

void send_s6a_ula(const std::string& imsi);

void send_create_session_resp();

void send_delete_session_resp();

void send_ics_response();

void send_ue_ctx_release_complete();

void send_ue_capabilities_ind();

}  // namespace lte
}  // namespace magma
