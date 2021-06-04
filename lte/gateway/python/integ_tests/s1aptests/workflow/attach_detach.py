"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import s1ap_types


def attach_ue(ue, s1ap_wrapper):
    print(
        "************************* Running End to End attach for ",
        "UE id ", ue.ue_id,
    )
    s1ap_wrapper.s1_util.attach(
        ue.ue_id, s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
        s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
        s1ap_types.ueAttachAccept_t,
    )
    # Wait on EMM Information from MME
    s1ap_wrapper.s1_util.receive_emm_info()


def detach_ue(ue, s1ap_wrapper):
    print(
        "************************* Running UE detach for UE id ",
        ue.ue_id,
    )
    s1ap_wrapper.s1_util.detach(
        ue.ue_id, s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
        True,
    )
