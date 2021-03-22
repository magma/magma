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
import ipaddress
import socket
import subprocess
from collections import namedtuple

from lte.protos.mobilityd_pb2 import IPAddress
from magma.pipelined.app.base import ControllerType, MagmaController
from magma.pipelined.app.inout import INGRESS
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import TUN_PORT_REG, Direction
from magma.pipelined.policy_converters import get_eth_type, get_ue_ip_match_args
from magma.pipelined.utils import Utils
from ryu.lib.packet import ether_types, packet, ipv4
from ryu.ofproto.ofproto_v1_4 import OFPP_LOCAL
from ryu.controller import ofp_event
from ryu.controller.handler import MAIN_DISPATCHER, set_ev_cls
from lte.protos.session_manager_pb2 import UPFPagingInfo
from magma.pipelined.app.enforcement_stats import _get_sid, _get_ipv4, _get_tunnel_id
from ryu.ofproto.ofproto_v1_4 import OFPMPF_REPLY_MORE
from ryu.ofproto import ofproto_v1_0_parser
from magma.pipelined.openflow import messages
from magma.pipelined.openflow.exceptions import MagmaOFError
from lte.protos.pipelined_pb2 import PdrState

GTP_PORT_MAC = "02:00:00:00:00:01"
TUNNEL_OAM_FLAG = 1

class Classifier(MagmaController):
    """
    A controller that sets up an openflow pipeline for Magma.
    This controller is used for install the default tunnel entry,
    add tunnel flows for gtp_port/local_port/mtr_port and delete
    tunnel flows for gtp_port/local_port/mtr_port into OVS table 0.
    """
    APP_NAME = "classifier"
    APP_TYPE = ControllerType.SPECIAL
    SESSIOND_RPC_TIMEOUT = 5
    CLASSIFIER_CONTROLLER_ID = 5
    ClassifierConfig = namedtuple(
            'ClassifierConfig',
            ['gtp_port', 'mtr_ip', 'mtr_port', 'internal_sampling_port',
             'internal_sampling_fwd_tbl', 'multi_tunnel_flag',
             'internal_conntrack_port', 'internal_conntrack_fwd_tbl'],
    )

    def __init__(self, *args, **kwargs):
        super(Classifier, self).__init__(*args, **kwargs)
        self.config = self._get_config(kwargs['config'])
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.next_table = self._service_manager.get_table_num(INGRESS)
        self._uplink_port = OFPP_LOCAL
        self._datapath = None
        self._clean_restart = kwargs['config']['clean_restart']
        if self.config.multi_tunnel_flag:
            self._ovs_multi_tunnel_init()
        #Get SessionD Channel
        self._sessiond_setinterface = kwargs['rpc_stubs']['sessiond_setinterface']
        self.loop = kwargs['loop']
        self.pagingflag = 0

    def _get_config(self, config_dict):
        mtr_ip = None
        mtr_port = None
        
        if 'mtr_ip' in config_dict:
            self._mtr_service_enabled = True
            mtr_ip = config_dict['mtr_ip']
            mtr_port = config_dict['ovs_mtr_port_number']

        if 'ovs_multi_tunnel' in config_dict:
            multi_tunnel_flag = config_dict['ovs_multi_tunnel']
        else:
            multi_tunnel_flag = False

        return self.ClassifierConfig(
            gtp_port=config_dict['ovs_gtp_port_number'],
            mtr_ip=mtr_ip,
            mtr_port=mtr_port,
            internal_sampling_port=
                               config_dict['ovs_internal_sampling_port_number'],
            internal_sampling_fwd_tbl=
                               config_dict['ovs_internal_sampling_fwd_tbl_number'],
            multi_tunnel_flag=multi_tunnel_flag,
            internal_conntrack_port=
                               config_dict['ovs_internal_conntrack_port_number'],
            internal_conntrack_fwd_tbl=
                               config_dict['ovs_internal_conntrack_fwd_tbl_number'],
            paging_timeout = config_dict['paging_timeout'],

        )

    def _ovs_multi_tunnel_init(self):
        try:
            subprocess.check_output("sudo ovs-vsctl list Open_vSwitch | grep gtpu",
                                     shell=True,)
            self.ovs_gtp_type = "gtpu"
        except subprocess.CalledProcessError:
            self.ovs_gtp_type = "gtp"

    def _ip_addr_to_gtp_port_name(self, enodeb_ip_addr:str):
        ip_no = hex(socket.htonl(int(ipaddress.ip_address(enodeb_ip_addr))))
        buf = "g_{}".format(ip_no[2:])
        return buf

    def _get_ofport(self, port_name):
        ovs = Utils.get_ovs_bridge(self._datapath)
        if ovs is None:
            return None

        try:
            port_no = ovs.get_ofport(port_name)

        except AssertionError as error:
            self.logger.debug('Cannot get port number for %s: %s',
                               port_name, error)
            return None

        except Exception as e:     #pylint: disable=broad-except
            self.logger.debug('Cannot get port number for %s: %s',
                               port_name, e)
            return None

        return port_no

    def _add_gtp_port(self, gnb_ip):
        if not self.config.multi_tunnel_flag:
            return self.config.gtp_port

        port_name = self._ip_addr_to_gtp_port_name(gnb_ip)
        # If GTP port already exists, returns OFPort number
        gtpport = self._get_ofport(port_name)
        if gtpport is not None:
            return gtpport

        ovs = Utils.get_ovs_bridge(self._datapath)
        if ovs is None:
            return None

        port_name = self._ip_addr_to_gtp_port_name(gnb_ip)
        ovs.add_tunnel_port(port_name, self.ovs_gtp_type,
                            gnb_ip, key="flow")

        return self._get_ofport(port_name)

    def initialize_on_connect(self, datapath):
        self._datapath = datapath
        if self._clean_restart:
            self._delete_all_flows()

        self._set_classifier_controller_id()
        self._install_default_tunnel_flows()
        self._install_internal_pkt_fwd_flow()
        self._install_internal_conntrack_flow()

    def _delete_all_flows(self):
        flows.delete_all_flows_from_table(self._datapath, self.tbl_num)

    def cleanup_on_disconnect(self, datapath):
        if self._clean_restart:
            self._delete_all_flows()

    def _set_classifier_controller_id(self):
        req = ofproto_v1_0_parser.NXTSetControllerId(self._datapath,
                                                     controller_id=self.CLASSIFIER_CONTROLLER_ID)
        messages.send_msg(self._datapath, req)

    def _install_default_tunnel_flows(self):
        match = MagmaMatch()
        flows.add_resubmit_next_service_flow(self._datapath,self.tbl_num, match,
                                             priority=flows.MINIMUM_PRIORITY,
                                             reset_default_register=False,
                                             resubmit_table=self.next_table)

    def _install_internal_pkt_fwd_flow(self):
        match = MagmaMatch(in_port=self.config.internal_sampling_port)
        flows.add_resubmit_next_service_flow(self._datapath,self.tbl_num, match,
                                             priority=flows.MINIMUM_PRIORITY,
                                             reset_default_register=False,
                                             resubmit_table=self.config.internal_sampling_fwd_tbl)

    def _install_internal_conntrack_flow(self):
        match = MagmaMatch(in_port=self.config.internal_conntrack_port)
        flows.add_resubmit_next_service_flow(self._datapath,self.tbl_num, match, [],
                                             priority=flows.MINIMUM_PRIORITY,
                                             reset_default_register=False,
                                             resubmit_table=self.config.internal_conntrack_fwd_tbl)

    def _send_messsage_interface(self, ue_ip_add:str):
        """
        Build the message
        """
        ofproto, parser = self._datapath.ofproto, self._datapath.ofproto_parser
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP, ipv4_src=ue_ip_add)
        ryu_match = parser.OFPMatch(**match.ryu_match)
        req = parser.OFPFlowStatsRequest(self._datapath, table_id=14,
                out_group=ofproto.OFPG_ANY,
                out_port=ofproto.OFPP_ANY,match=ryu_match)
        try:
            messages.send_msg(self._datapath, req)
            self.pagingflag = 1
        except MagmaOFError as e:
            self.logger.warning("Couldn't poll datapath stats: %s", e)
            self.pagingflag = 0

    @set_ev_cls(ofp_event.EventOFPFlowStatsReply, MAIN_DISPATCHER)
    def _flow_info_reply_handler(self, ev):
        if self.pagingflag:
            flow_info = ev.msg.body
            self.pagingflag = 0
            if ev.msg.flags == OFPMPF_REPLY_MORE:
                # Wait for more multi-part responses thats received for the
                # single stats request.
                return
            if flow_info:
                self.loop.call_soon_threadsafe(
                                       self._handle_flow_info, flow_info)

    def _handle_flow_info(self, flow_info):
        for stat in flow_info:
            sid = _get_sid(stat)
            ipv4_addr = _get_ipv4(stat)
            tunnel_id = _get_tunnel_id(stat)

        if (sid != None and ipv4_addr !=None and tunnel_id):
            paging_message=UPFPagingInfo(subscriber_id = sid, local_f_teid=tunnel_id,
                                         ue_ip_addr=ipv4_addr)

            future = self._sessiond_setinterface.SendPagingReuest.future(
                            paging_message, self.SESSIOND_RPC_TIMEOUT)
            future.add_done_callback(
                              lambda future: self.loop.call_soon_threadsafe(
                              self._callback_method, future))

    def _callback_method(self, future):
        """
        Callback after sessiond RPC completion
        """
        err = future.exception()
        if err:
            self.logger.error('Couldnt send flow records to sessiond: %s', err)

    def add_tunnel_flows(self, precedence:int, i_teid:int,
                         o_teid:int, ue_ip_adr:IPAddress,
                         enodeb_ip_addr:str, sid:int = None) -> bool:

        parser = self._datapath.ofproto_parser
        priority = Utils.get_of_priority(precedence)
        # Set minimun priority of GTP flow as 10
        if priority < flows.DEFAULT_PRIORITY:
            priority = flows.DEFAULT_PRIORITY
        # Add flow for gtp port
        if enodeb_ip_addr:
            gtp_portno = self._add_gtp_port(enodeb_ip_addr)
        else:
            gtp_portno = self.config.gtp_port

        # Add flow for gtp port for Uplink Tunnel
        actions = []
        if i_teid:
            match = MagmaMatch(tunnel_id=i_teid, in_port=gtp_portno)

            actions = [parser.OFPActionSetField(eth_src=GTP_PORT_MAC),
                       parser.OFPActionSetField(eth_dst="ff:ff:ff:ff:ff:ff")]
            if sid:
                actions.append(parser.OFPActionSetField(metadata=sid))

            flows.add_resubmit_next_service_flow(self._datapath, self.tbl_num, match,
                                                 actions=actions, priority=priority,
                                                 reset_default_register=False,
                                                 resubmit_table=self.next_table)

        # Install Downlink Tunnel
        actions = []
        if not ue_ip_adr:
            self.logger.error("ue_ip_address is None")
            return
        else:
            # Add flow for LOCAL port
            ip_match_out = get_ue_ip_match_args(ue_ip_adr, Direction.IN)
            match = MagmaMatch(eth_type=get_eth_type(ue_ip_adr),
                               in_port=self._uplink_port, **ip_match_out)
            if o_teid and enodeb_ip_addr:

                actions = [parser.OFPActionSetField(tunnel_id=o_teid),
                           parser.OFPActionSetField(tun_ipv4_dst=enodeb_ip_addr),
                           parser.OFPActionSetField(tun_flags=TUNNEL_OAM_FLAG),
                           parser.NXActionRegLoad2(dst=TUN_PORT_REG, value=gtp_portno)]
                if sid:
                    actions.append(parser.OFPActionSetField(metadata=sid))

                flows.add_resubmit_next_service_flow(self._datapath, self.tbl_num, match,
                                                     actions=actions, priority=priority,
                                                     reset_default_register=False,
                                                     resubmit_table=self.next_table)

            # Add flow for mtr port
            match = MagmaMatch(eth_type=get_eth_type(ue_ip_adr),
                               in_port=self.config.mtr_port, **ip_match_out)

            flows.add_resubmit_next_service_flow(self._datapath, self.tbl_num, match,
                                                 actions=actions, priority=priority,
                                                 reset_default_register=False,
                                                 resubmit_table=self.next_table)
       
            # Add ARP flow for LOCAL port
            if ue_ip_adr.version == IPAddress.IPV4:
                match = MagmaMatch(eth_type=ether_types.ETH_TYPE_ARP,
                                   in_port=self._uplink_port,
                                   arp_tpa=ipaddress.IPv4Address(ue_ip_adr.address.decode('utf-8')))
            actions = []
            if sid:
                actions = [parser.OFPActionSetField(metadata=sid)]

            flows.add_resubmit_next_service_flow(self._datapath, self.tbl_num, match,
                                                 actions=actions, priority=priority,
                                                 reset_default_register=False,
                                                 resubmit_table=self.next_table)

            # Add ARP flow for mtr port
            if ue_ip_adr.version == IPAddress.IPV4:
                match = MagmaMatch(eth_type=ether_types.ETH_TYPE_ARP,
                                   in_port=self.config.mtr_port,
                                   arp_tpa=ipaddress.IPv4Address(ue_ip_adr.address.decode('utf-8')))

            flows.add_resubmit_next_service_flow(self._datapath, self.tbl_num, match,
                                                 actions=actions, priority=priority,
                                                 reset_default_register=False,
                                                 resubmit_table=self.next_table)

        return True

    def delete_tunnel_flows(self, i_teid:int, ue_ip_adr:IPAddress,
                                 enodeb_ip_addr:str = None) -> bool:

        # Delete flow for gtp port
        if enodeb_ip_addr:
            gtp_portno = self._add_gtp_port(enodeb_ip_addr)
        else:
            gtp_portno = self.config.gtp_port

        if i_teid:
            match = MagmaMatch(tunnel_id=i_teid, in_port=gtp_portno)

            flows.delete_flow(self._datapath, self.tbl_num, match)

        # Delete flow for LOCAL port
        if not ue_ip_adr:
            self.logger.error("ue_ip_address is None")
            return
        else:
            ip_match_out = get_ue_ip_match_args(ue_ip_adr, Direction.IN)
            match = MagmaMatch(eth_type=get_eth_type(ue_ip_adr),
                               in_port=self._uplink_port, **ip_match_out)
            flows.delete_flow(self._datapath, self.tbl_num, match)

            # Delete flow for mtr port
            match = MagmaMatch(eth_type=get_eth_type(ue_ip_adr),
                               in_port=self.config.mtr_port, **ip_match_out)

            flows.delete_flow(self._datapath, self.tbl_num, match)

            # Delete ARP flow for LOCAL port
            if ue_ip_adr.version == IPAddress.IPV4:
                match = MagmaMatch(eth_type=ether_types.ETH_TYPE_ARP,
                                   in_port=self._uplink_port,
                                   arp_tpa=ipaddress.IPv4Address(ue_ip_adr.address.decode('utf-8')))

            flows.delete_flow(self._datapath, self.tbl_num, match)

            # Delete ARP flow for mtr port
            if ue_ip_adr.version == IPAddress.IPV4:
                match = MagmaMatch(eth_type=ether_types.ETH_TYPE_ARP,
                                   in_port=self.config.mtr_port,
                                   arp_tpa=ipaddress.IPv4Address(ue_ip_adr.address.decode('utf-8')))

            flows.delete_flow(self._datapath, self.tbl_num, match)

        return True    


    def _resume_tunnel_flows(self, precedence:int, i_teid:int,
                              ue_ip_adr:IPAddress=None):
        priority = Utils.get_of_priority(precedence)
        # Set minimun priority of GTP flow as 10
        if priority < flows.DEFAULT_PRIORITY:
            priority = flows.DEFAULT_PRIORITY
        # discard uplink Tunnel
        match = MagmaMatch(tunnel_id=i_teid, in_port=self.config.gtp_port)

        flows.add_resubmit_next_service_flow(self._datapath, self.tbl_num, match,
                                             priority=priority + 1,
                                             reset_default_register=False,
                                             resubmit_table=self.next_table)

        # discard downlink Tunnel for LOCAL port
        if not ue_ip_adr:
            self.logger.error("ue_ip_address is None")
            return
        else:
            ip_match_out = get_ue_ip_match_args(ue_ip_adr, Direction.IN)
            match = MagmaMatch(eth_type=get_eth_type(ue_ip_adr),
                           in_port=self._uplink_port, **ip_match_out)

            flows.add_resubmit_next_service_flow(self._datapath, self.tbl_num, match,
                                                 priority=priority + 1,
                                                 reset_default_register=False,
                                                 resubmit_table=self.next_table)

            # discard downlink Tunnel for mtr port
            match = MagmaMatch(eth_type=get_eth_type(ue_ip_adr),
                               in_port=self.config.mtr_port, **ip_match_out)

            flows.add_resubmit_next_service_flow(self._datapath, self.tbl_num, match,
                                                 priority=priority + 1,
                                                 reset_default_register=False,
                                                 resubmit_table=self.next_table)

    def _discard_tunnel_flows(self, precedence:int, i_teid:int,
                             ue_ip_adr:IPAddress=None):

        priority = Utils.get_of_priority(precedence)
        # Set minimun priority of GTP flow as 10
        if priority < flows.DEFAULT_PRIORITY:
            priority = flows.DEFAULT_PRIORITY
        # Forward flow for gtp port
        match = MagmaMatch(tunnel_id=i_teid, in_port=self.config.gtp_port)

        flows.delete_flow(self._datapath, self.tbl_num, match,
                          priority=priority + 1)

        # Forward flow for LOCAL port
        if not ue_ip_adr:
            self.logger.error("ue_ip_address is None")
            return
        else:
            ip_match_out = get_ue_ip_match_args(ue_ip_adr, Direction.IN)
            match = MagmaMatch(eth_type=get_eth_type(ue_ip_adr),
                               in_port=self._uplink_port, **ip_match_out)

            flows.delete_flow(self._datapath, self.tbl_num, match,
                              priority=priority +1)

            match = MagmaMatch(eth_type=get_eth_type(ue_ip_adr),
                               in_port=self.config.mtr_port, **ip_match_out)

            flows.delete_flow(self._datapath, self.tbl_num, match,
                              priority=priority + 1)

    @set_ev_cls(ofp_event.EventOFPPacketIn, MAIN_DISPATCHER)
    def _packet_in_handler(self, ev):
        msg = ev.msg

        pkt = packet.Packet(msg.data)
        pkt_ipv4 = pkt.get_protocol(ipv4.ipv4)

        dst = pkt_ipv4.dst
        # For sending notification to SMF using GRPC
        self._send_messsage_interface(dst)
        # Add flow for paging with hard time.
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP, ipv4_dst=dst)

        flows.add_drop_flow(self._datapath, self.tbl_num, match, [],
                       priority=flows.PAGING_PRIORITY + 1,
                       hard_timeout= self.config.paging_timeout)

    def _install_paging_flow(self, ue_ip_addr:IPAddress, controller_id:int=0):
        ofproto = self._datapath.ofproto
        parser = self._datapath.ofproto_parser
        ip_match_out = get_ue_ip_match_args(ue_ip_addr, Direction.IN)
        # Add flow for paging.
        match = MagmaMatch(eth_type=get_eth_type(ue_ip_addr), **ip_match_out)

        # Pass Controller ID value as a ACTION
        actions = [parser.NXActionController(0, controller_id,
                                             ofproto.OFPR_ACTION_SET)]

        flows.add_output_flow(self._datapath, self.tbl_num,
                              match=match, actions=actions,
                              priority=flows.PAGING_PRIORITY,
                              output_port=ofproto.OFPP_CONTROLLER,
                              max_len=ofproto.OFPCML_NO_BUFFER)

    def _remove_paging_flow(self, ue_ip_addr:IPAddress):
        ip_match_out = get_ue_ip_match_args(ue_ip_addr, Direction.IN)
        match = MagmaMatch(eth_type=get_eth_type(ue_ip_addr), **ip_match_out)
        flows.delete_flow(self._datapath, self.tbl_num, match)

    def gtp_handler(self, pdr_state, precedence:int, local_f_teid:int,
                     o_teid:int, ue_ip_addr:IPAddress, gnb_ip_addr:str,
                     sid:int = None, controller_id:int = 0):

        if pdr_state == PdrState.Value('INSTALL'):
            self.remove_paging_flow(ue_ip_addr)
            self.add_tunnel_flows(precedence, local_f_teid,
                                  o_teid, ue_ip_addr,
                                  gnb_ip_addr, sid)

        elif pdr_state == PdrState.Value('IDLE'):
            self.delete_tunnel_flows(local_f_teid, ue_ip_addr)
            self._install_paging_flow(ue_ip_addr, controller_id)

        elif pdr_state == PdrState.Value('REMOVE'):
            self.delete_tunnel_flows(local_f_teid, ue_ip_addr)
            self.remove_paging_flow(ue_ip_addr)

        return True

