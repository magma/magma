"""
Copyright (C) 2022  The Magma Authors

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; If not, see <http://www.gnu.org/licenses/>.
"""
from typing import Optional
from unittest.mock import patch

import pytest
from dhcp_helper_cli.dhcp_helper_cli import (
    DhcpHelperCli,
    DHCPState,
    MacAddress,
    release_arg_handler,
)
from scapy.layers.dhcp import BOOTP, DHCP
from scapy.layers.l2 import Dot1Q, Ether

MACSTRING = "12:34:56:78:90:AB"
MAC = MacAddress(MACSTRING)
VLAN = 0
IFACE = "dhcp0"
IPv4 = 0x0800

DHCP_OFFER_PKT = Ether(src="00:00:00:00:00:00", dst=MAC)
DHCP_OFFER_PKT /= BOOTP(yiaddr="1.2.3.4", chaddr=MAC.as_hex())
DHCP_OFFER_PKT /= DHCP(options=[("message-type", DHCPState.OFFER), "end"])

DHCP_ACK_PKT = Ether(src="00:00:00:00:00:00", dst=MAC)
DHCP_ACK_PKT /= BOOTP(yiaddr="1.2.3.4", chaddr=MAC.as_hex())
DHCP_ACK_PKT /= DHCP(options=[("message-type", DHCPState.ACK), "end"])

DHCP_RELEASE_PKT = Ether(src="00:00:00:00:00:00", dst=MAC)
DHCP_RELEASE_PKT /= BOOTP(yiaddr="1.2.3.4", chaddr=MAC.as_hex(), ciaddr="4.5.6.7")
DHCP_RELEASE_PKT /= DHCP(options=[("message-type", DHCPState.RELEASE), "end"])


@pytest.fixture()
def dhcp_helper_cli_fixture(mac: MacAddress = MAC, vlan: int = VLAN, iface=IFACE, ip: Optional[str] = None, server_ip: Optional[str] = None):
    with patch("dhcp_helper_cli.dhcp_helper_cli.AsyncSniffer"):
        return DhcpHelperCli(mac, vlan, iface, ip, server_ip)


@patch("dhcp_helper_cli.dhcp_helper_cli.sendp")
def test_send_dhcp_discover(sendp_mock, dhcp_helper_cli_fixture):
    dhcp_helper_cli_fixture.send_dhcp_discover()
    pkt = sendp_mock.call_args[0][0]

    assert pkt[Ether].type == IPv4
    assert pkt[Ether].src == MACSTRING.lower()
    assert Dot1Q not in pkt
    assert ('message-type', 'discover') in pkt[DHCP].options
    assert dhcp_helper_cli_fixture._state == DHCPState.DISCOVER


@patch("dhcp_helper_cli.dhcp_helper_cli.sendp")
def test_send_dhcp_request(sendp_mock, dhcp_helper_cli_fixture):
    dhcp_helper_cli_fixture._state = DHCPState.OFFER
    dhcp_helper_cli_fixture.send_dhcp_request()
    pkt = sendp_mock.call_args[0][0]

    assert pkt[Ether].type == IPv4
    assert pkt[Ether].src == MACSTRING.lower()
    assert Dot1Q not in pkt
    assert ('message-type', 'request') in pkt[DHCP].options
    assert dhcp_helper_cli_fixture._state == DHCPState.REQUEST


@patch("dhcp_helper_cli.dhcp_helper_cli.sendp")
def test_send_dhcp_release(sendp_mock, dhcp_helper_cli_fixture):
    dhcp_helper_cli_fixture._state = DHCPState.ACK
    dhcp_helper_cli_fixture.send_dhcp_release()
    pkt = sendp_mock.call_args[0][0]

    assert pkt[Ether].type == IPv4
    assert pkt[Ether].src == MACSTRING.lower()
    assert Dot1Q not in pkt
    assert ('message-type', 'release') in pkt[DHCP].options
    assert dhcp_helper_cli_fixture._state == DHCPState.RELEASE


def create_send_dhcp_pkt_mock(dhcp_helper_cli_fixture):
    def mocked_send_dhcp_pkt(dhcp_opts, ciaddr: Optional[str] = None) -> None:
        if ('message-type', 'discover') in dhcp_opts:
            dhcp_helper_cli_fixture._pkt_queue.put(DHCP_OFFER_PKT)
        if ('message-type', 'request') in dhcp_opts:
            dhcp_helper_cli_fixture._pkt_queue.put(DHCP_ACK_PKT)
        if ('message-type', 'release') in dhcp_opts:
            dhcp_helper_cli_fixture._pkt_queue.put(DHCP_RELEASE_PKT)

    return mocked_send_dhcp_pkt


def test_allocate(dhcp_helper_cli_fixture):
    with patch.object(
            dhcp_helper_cli_fixture,
            "send_dhcp_pkt",
            create_send_dhcp_pkt_mock(dhcp_helper_cli_fixture),
    ):
        dhcp_helper_cli_fixture.allocate()
        assert dhcp_helper_cli_fixture._state == DHCPState.ACK


def test_release(dhcp_helper_cli_fixture):
    with patch.object(
            dhcp_helper_cli_fixture,
            "send_dhcp_pkt",
            create_send_dhcp_pkt_mock(dhcp_helper_cli_fixture),
    ):
        dhcp_helper_cli_fixture.release()
        assert dhcp_helper_cli_fixture._state == DHCPState.RELEASE


def test_renew(dhcp_helper_cli_fixture):
    with patch.object(
            dhcp_helper_cli_fixture,
            "send_dhcp_pkt",
            create_send_dhcp_pkt_mock(dhcp_helper_cli_fixture),
    ):
        dhcp_helper_cli_fixture.renew()
        assert dhcp_helper_cli_fixture._state == DHCPState.ACK
