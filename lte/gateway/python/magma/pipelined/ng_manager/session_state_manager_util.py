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
from typing import NamedTuple
from lte.protos.pipelined_pb2 import ActivateFlowsRequest, \
            DeactivateFlowsRequest

FARRuleEntry =  NamedTuple(
                   'FARRuleEntry',
                    [('apply_action', int),
                     ('o_teid', int),
                     ('gnb_ip_addr', str)])

PDRRuleEntry = NamedTuple(
                'PDRRuleEntry',
                [('pdr_id', int),
                 ('pdr_version', int),
                 ('pdr_state', int),
                 ('precedence', int),
                 ('local_f_teid', int),
                 ('ue_ip_addr', str),
                 ('del_qos_enforce_rule', DeactivateFlowsRequest),
                 ('add_qos_enforce_rule', ActivateFlowsRequest),
                 ('far_action', FARRuleEntry)])

# Create the Named tuple for the FAR entry
def far_create_rule_entry(far_entry) -> FARRuleEntry:
    o_teid = 0
    fwd_gnb_ip_addr = None

    if far_entry.fwd_parm.HasField('outr_head_cr'):
        o_teid = far_entry.fwd_parm.outr_head_cr.o_teid
        fwd_gnb_ip_addr = far_entry.fwd_parm.outr_head_cr.gnb_ipv4_adr

    far_rule = FARRuleEntry(far_entry.far_action_to_apply[0],
                            o_teid, fwd_gnb_ip_addr)

    return far_rule

# Create the Named tuple for the PDR entry
def pdr_create_rule_entry(pdr_entry) -> PDRRuleEntry:
    local_f_teid = 0
    ue_ip_addr = None
    far_entry = None
    deactivate_flow_req = None
    activate_flow_req = None

    # get local teid
    if pdr_entry.pdi.local_f_teid:
        local_f_teid = pdr_entry.pdi.local_f_teid

    # Get UE IP address
    if len(pdr_entry.pdi.ue_ip_adr):
        ue_ip_addr = pdr_entry.pdi.ue_ip_adr

    if len(pdr_entry.set_gr_far.ListFields()):
        far_entry = far_create_rule_entry(pdr_entry.set_gr_far)

    if pdr_entry.HasField('deactivate_flow_req') == True:
        deactivate_flow_req = pdr_entry.deactivate_flow_req

    if pdr_entry.HasField('activate_flow_req') == True:
        activate_flow_req = pdr_entry.activate_flow_req

    pdr_rule = PDRRuleEntry(pdr_entry.pdr_id, pdr_entry.pdr_version,
                            pdr_entry.pdr_state, pdr_entry.precedence,
                            local_f_teid, ue_ip_addr,
                            deactivate_flow_req, activate_flow_req,
                            far_entry)
    return pdr_rule

