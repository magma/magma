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

from typing import List

from lte.protos.policydb_pb2 import FlowDescription, FlowMatch, PolicyRule


def get_allow_all_policy_rule(
    subscriber_id: str,
    apn: str,
) -> PolicyRule:
    """
    This builds a PolicyRule used as a default to allow all traffic
    for an attached subscriber.
    """
    policy_id = _get_allow_all_rule_id(subscriber_id, apn)
    return PolicyRule(
        # Don't set the rating group
        # Don't set the monitoring key
        # Don't set the hard timeout
        id=policy_id,
        priority=2,
        flow_list=_get_allow_all_flows(),
        tracking_type=PolicyRule.TrackingType.Value("NO_TRACKING"),
    )


def _get_allow_all_rule_id(subscriber_id: str, apn: str) -> str:
    rule_id_info = {'sid': subscriber_id, 'apn-id': apn}
    return "allowlist_sid-{sid}-{apn-id}".format(**rule_id_info)


def _get_allow_all_flows() -> List[FlowDescription]:
    """
    Get flows for allowing all traffic
    Returns:
        Two flows, for outgoing and incoming traffic
    """
    return [
        # Set flow match for all packets
        # Don't set the app_name field
        FlowDescription(  # uplink flow
            match=FlowMatch(
                direction=FlowMatch.Direction.Value("UPLINK"),
            ),
            action=FlowDescription.Action.Value("PERMIT"),
        ),
        FlowDescription(  # downlink flow
            match=FlowMatch(
                direction=FlowMatch.Direction.Value("DOWNLINK"),
            ),
            action=FlowDescription.Action.Value("PERMIT"),
        ),
    ]
