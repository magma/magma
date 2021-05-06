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
import logging

from magma.pipelined.openflow import flows
from ryu import cfg
from ryu.lib.ovs import bridge


class Utils:
    # Packet drop priority
    DROP_PRIORITY = flows.MINIMUM_PRIORITY + 1
    # For allowing unlcassified flows for app/service type rules.
    UNCLASSIFIED_ALLOW_PRIORITY = DROP_PRIORITY + 1
    # Should not overlap with the drop flow as drop matches all packets.
    MIN_PROGRAMMED_PRIORITY = UNCLASSIFIED_ALLOW_PRIORITY + 1
    MAX_PROGRAMMED_PRIORITY = flows.MAXIMUM_PRIORITY
    # Effectively range is 3 -> 65535
    APP_PRIORITY_RANGE = MAX_PROGRAMMED_PRIORITY - MIN_PROGRAMMED_PRIORITY

    # Resume tunnel flows
    RESUME_RULE_PRIORITY = flows.DEFAULT_PRIORITY + 1
    # Discard tunnel flows
    DISCARD_RULE_PRIORITY = RESUME_RULE_PRIORITY
    # Paging tunnel flows
    PAGING_RULE_PRIORITY = 5
    PAGING_RULE_DROP_PRIORITY = PAGING_RULE_PRIORITY + 1

    OVSDB_PORT = 6640  # The IANA registered port for OVSDB [RFC7047]
    CONF = cfg.CONF
    # OVSBridge instance instantiated later
    ovs = None

    @classmethod
    def get_of_priority(cls, precedence:int):
        """
        Lower the precedence higher the importance of the flow in 3GPP.
        Higher the priority higher the importance of the flow in openflow.
        Convert precedence to priority:
        1 - Flows with precedence > 65534 will have min priority which is the
        min priority for a programmed flow = (default drop + 1)
        2 - Flows in the precedence range 0-65534 will have priority 65535 -
        Precedence
        :param precedence:
        :return:
        """
        if precedence >= cls.APP_PRIORITY_RANGE:
            logging.warning(
                "Flow precedence is higher than OF range using min priority %d",
                cls.MIN_PROGRAMMED_PRIORITY)
            return cls.MIN_PROGRAMMED_PRIORITY
        return cls.MAX_PROGRAMMED_PRIORITY - precedence

    @classmethod
    def get_ovs_bridge(cls, datapath):
        dpid = datapath.id
        ovsdb_addr = 'tcp:%s:%d' % (datapath.address[0], cls.OVSDB_PORT)

        if (cls.ovs is not None
                and cls.ovs.datapath_id == dpid
                and cls.ovs.vsctl.remote == ovsdb_addr):
            return cls.ovs

        try:
            cls.ovs = bridge.OVSBridge(
                CONF=cls.CONF,
                datapath_id=dpid,
                ovsdb_addr=ovsdb_addr)
            cls.ovs.init()
        except (ValueError, KeyError) as e:
            logging.warning('Cannot initiate OVSDB connection: %s', e)
            return None
        return cls.ovs

