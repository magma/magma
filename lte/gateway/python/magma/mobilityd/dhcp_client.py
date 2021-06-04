"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Allocates IP address as per DHCP server in the uplink network.
"""
import datetime
import logging
import threading
import time
from ipaddress import IPv4Network, ip_address
from threading import Condition
from typing import MutableMapping, Optional

from magma.mobilityd.dhcp_desc import DHCPDescriptor, DHCPState
from magma.mobilityd.mac import MacAddress, hex_to_mac
from magma.mobilityd.uplink_gw import UplinkGatewayInfo
from scapy.all import AsyncSniffer
from scapy.layers.dhcp import BOOTP, DHCP
from scapy.layers.inet import IP, UDP
from scapy.layers.l2 import Dot1Q, Ether
from scapy.sendrecv import sendp

LOG = logging.getLogger('mobilityd.dhcp.sniff')
DHCP_ACTIVE_STATES = [DHCPState.ACK, DHCPState.OFFER]


class DHCPClient:
    THREAD_YIELD_TIME = .1

    def __init__(
        self,
        dhcp_store: MutableMapping[str, DHCPDescriptor],
        gw_info: UplinkGatewayInfo,
        dhcp_wait: Condition,
        iface: str = "dhcp0",
        lease_renew_wait_min: int = 200,
    ):
        """
        Implement DHCP client to allocate IP for given Mac address.
        DHCP client state is maintained in user provided hash table.
        Args:
            dhcp_store: maintain DHCP transactions, key is mac address.
            gw_info_map: stores GW IP info from DHCP server
            dhcp_wait: notify users on new DHCP packet
            iface: DHCP egress and ingress interface.
        """
        self._sniffer = AsyncSniffer(
            iface=iface,
            filter="udp and (port 67 or 68)",
            store=False,
            prn=self._rx_dhcp_pkt,
        )

        self.dhcp_client_state = dhcp_store  # mac => DHCP_State
        self.dhcp_gw_info = gw_info
        self._dhcp_notify = dhcp_wait
        self._dhcp_interface = iface
        self._msg_xid = 0
        self._lease_renew_wait_min = lease_renew_wait_min
        self._monitor_thread = threading.Thread(
            target=self._monitor_dhcp_state,
        )
        self._monitor_thread.daemon = True
        self._monitor_thread_event = threading.Event()

    def run(self):
        """
        Start DHCP sniffer thread.
        This initializes state required for DHCP sniffer thread anf starts it.
        Returns: None
        """
        self._sniffer.start()
        LOG.info("DHCP sniffer started")
        # give it time to schedule the thread and start sniffing.
        time.sleep(self.THREAD_YIELD_TIME)
        self._monitor_thread.start()

    def stop(self):
        self._sniffer.stop()
        self._monitor_thread_event.set()

    def send_dhcp_packet(
        self, mac: MacAddress, vlan: str,
        state: DHCPState,
        dhcp_desc: DHCPDescriptor = None,
    ):
        """
        Send DHCP packet and record state in dhcp_client_state.

        Args:
            mac: MAC address of interface
            state: state of DHCP packet
            dhcp_desc: DHCP protocol state.
        Returns:
        """
        ciaddr = None

        # generate DHCP request packet
        if state == DHCPState.DISCOVER:
            dhcp_opts = [("message-type", "discover")]
            dhcp_desc = DHCPDescriptor(
                mac=mac, ip="", vlan=vlan,
                state_requested=DHCPState.DISCOVER,
            )
            self._msg_xid = self._msg_xid + 1
            pkt_xid = self._msg_xid
        elif state == DHCPState.REQUEST:
            dhcp_opts = [
                ("message-type", "request"),
                ("requested_addr", dhcp_desc.ip),
                ("server_id", dhcp_desc.server_ip),
            ]
            dhcp_desc.state_requested = DHCPState.REQUEST
            pkt_xid = dhcp_desc.xid
            ciaddr = dhcp_desc.ip
        elif state == DHCPState.RELEASE:
            dhcp_opts = [
                ("message-type", "release"),
                ("server_id", dhcp_desc.server_ip),
            ]
            dhcp_desc.state_requested = DHCPState.RELEASE
            self._msg_xid = self._msg_xid + 1
            pkt_xid = self._msg_xid
            ciaddr = dhcp_desc.ip
        else:
            LOG.warning(
                "Unknown egress request mac %s state %s",
                str(mac),
                state,
            )
            return

        dhcp_opts.append("end")
        dhcp_desc.xid = pkt_xid
        with self._dhcp_notify:
            self.dhcp_client_state[mac.as_redis_key(vlan)] = dhcp_desc

        pkt = Ether(src=str(mac), dst="ff:ff:ff:ff:ff:ff")
        if vlan and vlan != "0":
            pkt /= Dot1Q(vlan=int(vlan))
        pkt /= IP(src="0.0.0.0", dst="255.255.255.255")
        pkt /= UDP(sport=68, dport=67)
        pkt /= BOOTP(op=1, chaddr=mac.as_hex(), xid=pkt_xid, ciaddr=ciaddr)
        pkt /= DHCP(options=dhcp_opts)
        LOG.debug("DHCP pkt xmit %s", pkt.show(dump=True))

        sendp(pkt, iface=self._dhcp_interface, verbose=0)

    def get_dhcp_desc(
        self, mac: MacAddress,
        vlan: str,
    ) -> Optional[DHCPDescriptor]:
        """
                Get DHCP description for given MAC.
        Args:
            mac: Mac address of the client
            vlan: vlan id if the IP allocated in a VLAN

        Returns: Current DHCP info.
        """

        key = mac.as_redis_key(vlan)
        if key in self.dhcp_client_state:
            return self.dhcp_client_state[key]

        LOG.debug("lookup error for %s", str(key))
        return None

    def release_ip_address(self, mac: MacAddress, vlan: str):
        """
                Release DHCP allocated IP.
        Args:
            mac: MAC address of the IP allocated.
            vlan: vlan id if the IP allocated in a VLAN

        Returns: None
        """
        key = mac.as_redis_key(vlan)
        if key not in self.dhcp_client_state:
            LOG.error("Unallocated DHCP release for MAC: %s", key)
            return

        dhcp_desc = self.dhcp_client_state[key]
        self.send_dhcp_packet(
            mac,
            dhcp_desc.vlan,
            DHCPState.RELEASE,
            dhcp_desc,
        )
        del self.dhcp_client_state[key]

    def _monitor_dhcp_state(self):
        """
        monitor DHCP client state.
        """
        while True:
            wait_time = self._lease_renew_wait_min
            with self._dhcp_notify:
                for dhcp_record in self.dhcp_client_state.values():
                    logging.debug("monitor: %s", dhcp_record)
                    # Only process active records.
                    if dhcp_record.state not in DHCP_ACTIVE_STATES:
                        continue

                    now = datetime.datetime.now()
                    logging.debug("monitor time: %s", now)
                    request_state = DHCPState.REQUEST
                    # in case of lost DHCP lease rediscover it.
                    if now >= dhcp_record.lease_expiration_time:
                        request_state = DHCPState.DISCOVER

                    if now >= dhcp_record.lease_renew_deadline:
                        logging.debug("sending lease renewal")
                        self.send_dhcp_packet(
                            dhcp_record.mac, dhcp_record.vlan,
                            request_state, dhcp_record,
                        )
                    else:
                        # Find next renewal wait time.
                        time_to_renew = dhcp_record.lease_renew_deadline - now
                        wait_time = min(
                            wait_time, time_to_renew.total_seconds(),
                        )

            # default in wait is 30 sec
            wait_time = max(wait_time, self._lease_renew_wait_min)
            logging.debug("lease renewal check after: %s sec" % wait_time)
            self._monitor_thread_event.wait(wait_time)
            if self._monitor_thread_event.is_set():
                break

    @staticmethod
    def _get_option(packet, name):
        for opt in packet[DHCP].options:
            if opt[0] == name:
                return opt[1]
        return None

    def _process_dhcp_pkt(self, packet, state: DHCPState):
        LOG.debug("DHCP pkt recv %s", packet.show(dump=True))

        mac_addr = MacAddress(hex_to_mac(packet[BOOTP].chaddr.hex()[0:12]))
        vlan = ""
        if Dot1Q in packet:
            vlan = str(packet[Dot1Q].vlan)
        mac_addr_key = mac_addr.as_redis_key(vlan)

        with self._dhcp_notify:
            if mac_addr_key in self.dhcp_client_state:
                state_requested = self.dhcp_client_state[mac_addr_key].state_requested
                if BOOTP not in packet or packet[BOOTP].yiaddr is None:
                    LOG.error("no ip offered")
                    return

                ip_offered = packet[BOOTP].yiaddr
                subnet_mask = self._get_option(packet, "subnet_mask")
                if subnet_mask is not None:
                    ip_subnet = IPv4Network(
                        ip_offered + "/" + subnet_mask, strict=False,
                    )
                else:
                    ip_subnet = IPv4Network(
                        ip_offered + "/" + "32", strict=False,
                    )

                dhcp_server_ip = None
                if IP in packet:
                    dhcp_server_ip = packet[IP].src

                dhcp_router_opt = self._get_option(packet, "router")
                if dhcp_router_opt is not None:
                    router_ip_addr = ip_address(dhcp_router_opt)
                else:
                    # use DHCP as upstream router in case of missing Open 3.
                    router_ip_addr = dhcp_server_ip
                self.dhcp_gw_info.update_ip(router_ip_addr, vlan)

                lease_expiration_time = self._get_option(packet, "lease_time")
                dhcp_state = DHCPDescriptor(
                    mac=mac_addr,
                    ip=ip_offered,
                    state=state,
                    vlan=vlan,
                    state_requested=state_requested,
                    subnet=str(ip_subnet),
                    server_ip=dhcp_server_ip,
                    router_ip=router_ip_addr,
                    lease_expiration_time=lease_expiration_time,
                    xid=packet[BOOTP].xid,
                )
                LOG.info(
                    "Record DHCP for: %s state: %s",
                    mac_addr_key,
                    dhcp_state,
                )

                self.dhcp_client_state[mac_addr_key] = dhcp_state
                self._dhcp_notify.notifyAll()

                if state == DHCPState.OFFER:
                    # let other thread work on fulfilling IP allocation
                    # request.
                    threading.Event().wait(self.THREAD_YIELD_TIME)
                    self.send_dhcp_packet(
                        mac_addr, vlan, DHCPState.REQUEST, dhcp_state,
                    )
            else:
                LOG.debug("Unknown MAC: %s " % packet.summary())
                return

    # ref: https://fossies.org/linux/scapy/scapy/layers/dhcp.py
    def _rx_dhcp_pkt(self, packet):
        if DHCP not in packet:
            return

        # Match DHCP offer
        if packet[DHCP].options[0][1] == int(DHCPState.OFFER):
            self._process_dhcp_pkt(packet, DHCPState.OFFER)

        # Match DHCP ack
        elif packet[DHCP].options[0][1] == int(DHCPState.ACK):
            self._process_dhcp_pkt(packet, DHCPState.ACK)

        # TODO handle other DHCP protocol events.
