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
from typing import NamedTuple, Optional

from lte.protos.pipelined_pb2 import (
    ActivateFlowsRequest,
    DeactivateFlowsRequest,
)


class FARRuleEntry(NamedTuple):
    apply_action: int
    gnb_ip_addr: Optional[str]


class PDRRuleEntry(NamedTuple):
    pdr_id: int
    pdr_version: int
    pdr_state: int
    precedence: int
    local_f_teid: int
    ue_ip_addr: Optional[str]
    o_teid: int
    del_qos_enforce_rule: DeactivateFlowsRequest
    add_qos_enforce_rule: ActivateFlowsRequest
    far_action: Optional[FARRuleEntry]
    ue_ipv6_addr: Optional[str]
    session_qfi: Optional[int]


# Create the Named tuple for the FAR entry


def far_create_rule_entry(far_entry) -> FARRuleEntry:
    fwd_gnb_ip_addr = None

    if far_entry.fwd_parm.HasField('outr_head_cr'):
        fwd_gnb_ip_addr = far_entry.fwd_parm.outr_head_cr.gnb_ipv4_adr

    far_rule = FARRuleEntry(
        far_entry.far_action_to_apply[0],
        fwd_gnb_ip_addr,
    )

    return far_rule

# Create the Named tuple for the PDR entry


def pdr_create_rule_entry(pdr_entry) -> PDRRuleEntry:
    local_f_teid = 0
    ue_ip_addr = None
    far_entry = None
    deactivate_flow_req = None
    activate_flow_req = None
    ue_ipv6_addr = None
    session_qfi = 0
    o_teid = 0

    if pdr_entry.gnb_teid:
        o_teid = pdr_entry.gnb_teid
    # get local teid
    if pdr_entry.pdi.local_f_teid:
        local_f_teid = pdr_entry.pdi.local_f_teid

    # Get UE IP address
    if pdr_entry.pdi.ue_ipv4:
        ue_ip_addr = pdr_entry.pdi.ue_ipv4

    # Get UE IP address
    if pdr_entry.pdi.ue_ipv6:
        ue_ipv6_addr = pdr_entry.pdi.ue_ipv6

    if len(pdr_entry.set_gr_far.ListFields()):
        far_entry = far_create_rule_entry(pdr_entry.set_gr_far)

    if pdr_entry.HasField('deactivate_flow_req') == True:
        deactivate_flow_req = pdr_entry.deactivate_flow_req

    if pdr_entry.HasField('activate_flow_req') == True:
        activate_flow_req = pdr_entry.activate_flow_req

    if pdr_entry.session_qfi:
        session_qfi = pdr_entry.session_qfi

    pdr_rule = PDRRuleEntry(
        pdr_entry.pdr_id, pdr_entry.pdr_version,
        pdr_entry.pdr_state, pdr_entry.precedence,
        local_f_teid, ue_ip_addr, o_teid,
        deactivate_flow_req, activate_flow_req,
        far_entry, ue_ipv6_addr, session_qfi,
    )
    return pdr_rule
