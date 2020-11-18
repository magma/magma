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

import abc
import copy

from lte.protos.mobilityd_pb2 import IPAddress
from magma.pipelined.imsi import encode_imsi
from magma.pipelined.policy_converters import convert_ip_str_to_ip_proto
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import DIRECTION_REG, Direction, \
    IMSI_REG
from integ_tests.s1aptests.ovs import LOCALHOST
from integ_tests.s1aptests.ovs.rest_api import get_datapath,\
    delete_flowentry, add_flowentry

from ryu.lib.packet import ether_types
from ryu.lib import hub


class TableIsolator(abc.ABC):
    """
    Interface for table isolation

    Table isolation forwards all packets directly to the destination table,
    it also sets the register values specified
    TableIsolator implemntations are used as contexts, rules will be added to
    ovs and then terminated when the context is no longer used
    """

    def __enter__(self):
        """
        Used for running 'with' (isolates the destination table)
        """
        self._activate_flow_rules()

    def __exit__(self, type, value, traceback):
        """
        Clean up after using 'with' (cleans up everything that was added)
        """
        self._deactivate_flow_rules()

    def _activate_flow_rules(self):
        """
        Sets up the flows to forward packets to destination table

        If the ip is set then 2 flows are created, uplink and downlink,
        these rules will set the direction register(reg1)'s value.

        Otherwise use a single flow rule independent of all packets and set
        the register values configured by calling set_reg_value
        """
        raise NotImplementedError()

    def _deactivate_flow_rules(self):
        """
        Removes all flow rules that were added from ovs
        """
        raise NotImplementedError()


class RyuForwardFlowArgsBuilder():
    """
    Ryu Forward Flow Arguments Builder

    RyuForwardFlow is used to build ryu requests for table isolation, the flows
    will forward all traffic to given destination and set other optional args
    """

    """ Default values flow attribute values"""
    PRIORITY = 65535
    COOKIE = 9999

    def __init__(self, table_dest, table_start=0):
        self._ip = None
        self._teid_dlink_action = []
        self._teid_uplink_match = 0
        self._reg_sets = []
        self._request = {
            "table_id": table_start, "cookie": self.COOKIE,
            "priority": self.PRIORITY, "match": {},
            "instructions": [{"type": "GOTO_TABLE", "table_id": table_dest}],
        }
        self._match_kwargs = {}

    def set_reg_value(self, reg_name, value):
        """
        Set flow register action (flow will set reg_name to value)
        Args:
            reg_name (string): register name
            value (int): reg value
        Returns:
            Self
        """
        self._reg_sets.append(
            {"type": "SET_FIELD", "field": reg_name, "value": value}
        )
        return self

    def set_cookie(self, cookie):
        """
        Set flow cookie value
        Args:
            cookie (int): cookie value
        Returns:
            Self
        """
        self._request["cookie"] = cookie
        return self

    def set_priority(self, priority):
        """
        Set the flow priority
        Args:
            priority (int): priority value
        Returns:
            Self
        """
        self._request["priority"] = priority
        return self

    def set_ip(self, ip):
        """
        Set Match IPs and set register values for ovs flows:
            src_ip match for UPLINK    (DIRECTION_REG = 0x1   Direction.OUT)
            dst_ip match for DOWNLINK  (DIRECTION_REG = 0x10  Direction.IN)
        Args:
            ip (string): ip value
        Returns:
            Self
        """
        self._ip = ip
        self._ulink_action = {"type": "SET_FIELD", "field": DIRECTION_REG,
                              "value": Direction.OUT}

        self._dlink_action = {"type": "SET_FIELD", "field": DIRECTION_REG,
                              "value": Direction.IN}
        return self

    def set_tunnel_id(self, uplink_tunnel):
        self._teid_uplink_match = uplink_tunnel 

        if uplink_tunnel:
            self._teid_ulink_action = {"type": "SET_FIELD", "field": "tunnel_id",
                                        "value": uplink_tunnel}
        return self

    def set_eth_match(self, eth_src, eth_dst):
        self._match_kwargs['eth_src'] = eth_src
        self._match_kwargs['eth_dst'] = eth_dst
        return self

    def _update_ng_subscriber_tunnel_id_requests(self, uplink, downlink):
        if self._teid_dlink_action:
            for down_insts in downlink["instructions"]:
                if "actions" in down_insts.keys():
                    down_insts["actions"].append(self._teid_dlink_action)

        #Matching with tunnel id is challenging
        if self._teid_ulink_action:
            #uplink["match"].update({"tunnel_id": self._teid_uplink_match})
            for up_insts in uplink["instructions"]:
               if "actions" in up_insts.keys():
                  up_insts["actions"].append(self._teid_ulink_action)

        return [uplink, downlink]

    def _create_subscriber_ip_requests(self):
        """
        Generates ryu requests for subscriber flows

        Based on the provided ip create 2 flows, that are matched based on
        dst/src ip and set value of the direction register.

        Additional reg values are set from set_reg_value
        """

        uplink = copy.deepcopy(self._request)
        downlink = copy.deepcopy(self._request)

        uplink["instructions"].append({
            "type": "APPLY_ACTIONS",
            "actions": self._reg_sets + [self._ulink_action]
        })
        downlink["instructions"].append({
            "type": "APPLY_ACTIONS",
            "actions": self._reg_sets + [self._dlink_action]
        })

        ip_addr = convert_ip_str_to_ip_proto(self._ip)
        if ip_addr.version == IPAddress.IPV4:
            uplink["match"].update(
                {"ipv4_src": self._ip})
            downlink["match"].update(
                {"ipv4_dst": self._ip})
        else:
            uplink["match"].update(
                {"ipv6_src": self._ip})
            downlink["match"].update(
                {"ipv6_dst": self._ip})

        if self._teid_dlink_action or self._teid_uplink_match:
            self._update_ng_subscriber_tunnel_id_requests(uplink, downlink)

        return [uplink, downlink]

    def set_eth_type_arp(self):
        self._match_kwargs = {"eth_type": ether_types.ETH_TYPE_ARP}
        return self

    def _set_subscriber_match(self, sub_info, sub_ng_info=None):
        """ Sets up match/action for subscriber flows """
        if sub_info.ip.count(":") >= 2:
            self._match_kwargs = {"eth_type": ether_types.ETH_TYPE_IPV6}
        else:
            self._match_kwargs = {"eth_type": ether_types.ETH_TYPE_IP}

        if sub_ng_info and sub_ng_info.uplink_tunnel:
            return self.set_ip(sub_info.ip).\
                   set_reg_value(IMSI_REG, encode_imsi(sub_info.imsi)).\
                   set_tunnel_id(sub_ng_info.uplink_tunnel)

        return self.set_ip(sub_info.ip) \
            .set_reg_value(IMSI_REG, encode_imsi(sub_info.imsi))

    def build_requests(self):
        """
        From the set arguments generate ryu request(s)
        Returns:
            requests [dict]: generated ryu requests
        """
        self._request["match"] = MagmaMatch(**self._match_kwargs)
        if self._ip is not None:
            return self._create_subscriber_ip_requests()
        else:
            if self._reg_sets:
                self._request["instructions"].append(
                    {"type": "APPLY_ACTIONS", "actions": self._reg_sets}
                )
            return [self._request]

    @classmethod
    def from_subscriber(cls, sub_info, sub_ng_info=None):
        return cls(sub_info.table_id)._set_subscriber_match(sub_info, sub_ng_info)


# REST API is deprecated transition to RyuDirectTableIsolator
class RyuRestTableIsolator(TableIsolator):
    """
    RyuRestTableIsolator uses ryu REST api to isolate tables, sends the
    generated RyuForwardFlow requests as REST requests.
    """

    def __init__(self, requests, ovs_ip=LOCALHOST):
        self._requests = requests
        self._ovs_ip = ovs_ip
        self._datapath = get_datapath(ovs_ip)

    def _activate_flow_rules(self):
        """ Adds the flows to ovs, REST needs a dpid argument """
        for req in self._requests:
            req["dpid"] = self._datapath
            add_flowentry(req, self._ovs_ip)

    def _deactivate_flow_rules(self):
        """ Removes flows from ovs, REST needs a dpid argument """
        for req in self._requests:
            req["dpid"] = self._datapath
            delete_flowentry(req, self._ovs_ip)


class RyuDirectTableIsolator(TableIsolator):
    """
    RyuDirectTableIsolator uses ryu.hub and test_controller to isolate tables,
    sends the generated RyuForwardFlow requests to test_controller.
    """

    def __init__(self, requests, test_controller):
        self._tc = test_controller
        self._requests = requests

    def _activate_flow_rules(self):
        def insert_flow(req):
            self._tc.insert_flow(req)
        for req in self._requests:
            hub.joinall([hub.spawn(insert_flow, req)])

    def _deactivate_flow_rules(self):
        def delete_flow(req):
            self._tc.delete_flow(req)
        for req in self._requests:
            hub.joinall([hub.spawn(delete_flow, req)])
