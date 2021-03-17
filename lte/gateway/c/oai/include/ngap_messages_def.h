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
/*****************************************************************************

  Source      ngap_messages_def.h
  Date        2020/09/08
  Subsystem   Access and Mobility Management Function
  Description Defines Access and Mobility Management Messages

*****************************************************************************/
MESSAGE_DEF(
    NGAP_UE_CONTEXT_RELEASE_REQ, itti_ngap_ue_context_release_req_t,
    ngap_ue_context_release_req)
MESSAGE_DEF(
    NGAP_UE_CONTEXT_RELEASE_COMMAND, itti_ngap_ue_context_release_command_t,
    ngap_ue_context_release_command)
MESSAGE_DEF(
    NGAP_UE_CONTEXT_RELEASE_COMPLETE, itti_ngap_ue_context_release_complete_t,
    ngap_ue_context_release_complete)
MESSAGE_DEF(
    NGAP_NAS_DL_DATA_REQ, itti_ngap_nas_dl_data_req_t, ngap_nas_dl_data_req)
MESSAGE_DEF(
    NGAP_INITIAL_UE_MESSAGE, itti_ngap_initial_ue_message_t,
    ngap_initial_ue_message)
