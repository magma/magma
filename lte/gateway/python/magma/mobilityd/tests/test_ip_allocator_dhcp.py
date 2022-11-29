"""
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""
from datetime import timedelta
from ipaddress import IPv4Network, IPv4Address
from typing import Any, List
from unittest.mock import patch, MagicMock

import fakeredis
import freezegun
import pytest
from magma.mobilityd.dhcp_desc import DHCPDescriptor, DHCPState
from magma.mobilityd.ip_allocator_dhcp import IPAllocatorDHCP, DHCP_CLI_HELPER_PATH
from magma.mobilityd.ip_descriptor import IPState, IPDesc, IPType
from magma.mobilityd.ip_descriptor_map import IpDescriptorMap
from magma.mobilityd.mac import sid_to_mac, MacAddress
from magma.mobilityd.mobility_store import AssignedIpBlocksSet, ip_states, defaultdict_key, MobilityStore

SID = "IMSI123456789"
MAC = MacAddress(sid_to_mac(SID).lower())
MAC2 = MacAddress("01:23:45:67:89:ab")
IP = "1.2.3.4"
SUBNET = "24"
IP_NETWORK = "1.2.3.0/" + SUBNET
IP_NETWORK_2 = "1.2.4.0/" + SUBNET
VLAN = "0"
LEASE_EXPIRATION_TIME = 10
FROZEN_TEST_TIME = "2021-01-01"


@pytest.fixture
def dhcp_desc_fixture() -> DHCPDescriptor:
    with freezegun.freeze_time(FROZEN_TEST_TIME):
        return DHCPDescriptor(
            mac=MAC,
            ip=IP_NETWORK,
            vlan=VLAN,
            state=DHCPState.ACK,
            state_requested=DHCPState.REQUEST,
            lease_expiration_time=4,
        )


@pytest.fixture
def ip_allocator_fixture() -> IPAllocatorDHCP:
    client = fakeredis.FakeStrictRedis()

    ip_allocator = IPAllocatorDHCP(
        store=MobilityStore(client),
        lease_renew_wait_min=4,
        start=False,
    )

    yield ip_allocator

    ip_allocator._monitor_thread_event.set()
    ip_allocator._monitor_thread.join()


@pytest.fixture
def ip_allocator_dhcp_fixture(ip_allocator_fixture, dhcp_desc_fixture: DHCPDescriptor) -> IPAllocatorDHCP:
    ip_allocator_fixture._store.dhcp_store[
        dhcp_desc_fixture.mac.as_redis_key(dhcp_desc_fixture.vlan)] = dhcp_desc_fixture

    yield ip_allocator_fixture


@pytest.fixture
def ip_desc_fixture() -> IPDesc:
    return IPDesc(
        ip=IPv4Address(IP),
        state=IPState.ALLOCATED,
        vlan_id=int(VLAN),
        ip_block=IPv4Network(IP_NETWORK),
        ip_type=IPType.DHCP,
        sid=SID,
    )


def create_subprocess_mock() -> MagicMock:
    m = MagicMock()
    m.returncode = 0
    m.stdout = """{"lease_expiration_time": 4}"""
    return m


def run_dhcp_allocator_thread(
        frozen_datetime: Any,
        ip_allocator_dhcp_fixture: IPAllocatorDHCP,
        freeze_time: float) -> None:
    ip_allocator_dhcp_fixture._monitor_thread_event.set()
    frozen_datetime.tick(timedelta(seconds=freeze_time))
    ip_allocator_dhcp_fixture.start_monitor_thread()
    ip_allocator_dhcp_fixture._monitor_thread.join()


def test_no_renewal_of_ip(ip_allocator_dhcp_fixture: IPAllocatorDHCP) -> None:
    with freezegun.freeze_time(FROZEN_TEST_TIME) as frozen_datetime, \
            patch("subprocess.run", return_value=create_subprocess_mock()) as subprocess_mock:
        run_dhcp_allocator_thread(
            frozen_datetime=frozen_datetime,
            ip_allocator_dhcp_fixture=ip_allocator_dhcp_fixture,
            freeze_time=1
        )

        subprocess_mock.assert_not_called()


def test_renewal_of_ip(
        ip_allocator_dhcp_fixture: IPAllocatorDHCP,
        dhcp_desc_fixture: DHCPDescriptor) -> None:
    with freezegun.freeze_time(FROZEN_TEST_TIME) as frozen_datetime, \
            patch("subprocess.run", return_value=create_subprocess_mock()) as subprocess_mock:
        run_dhcp_allocator_thread(
            frozen_datetime=frozen_datetime,
            ip_allocator_dhcp_fixture=ip_allocator_dhcp_fixture,
            freeze_time=3
        )

        subprocess_mock.assert_called_once()
        subprocess_mock.assert_called_with([
            DHCP_CLI_HELPER_PATH,
            "--mac", str(dhcp_desc_fixture.mac),
            "--vlan", str(dhcp_desc_fixture.vlan),
            "--interface", ip_allocator_dhcp_fixture._iface,
            "--json",
            "renew",
            "--ip", str(dhcp_desc_fixture.ip),
            "--server-ip", str(dhcp_desc_fixture.server_ip),
        ],
            capture_output=True
        )


def test_allocate_ip_after_expiry(ip_allocator_dhcp_fixture: IPAllocatorDHCP,
                                  dhcp_desc_fixture: DHCPDescriptor) -> None:
    with freezegun.freeze_time(FROZEN_TEST_TIME) as frozen_datetime, \
            patch("subprocess.run", return_value=create_subprocess_mock()) as subprocess_mock:
        run_dhcp_allocator_thread(
            frozen_datetime=frozen_datetime,
            ip_allocator_dhcp_fixture=ip_allocator_dhcp_fixture,
            freeze_time=5
        )
        subprocess_mock.assert_called_once()
        subprocess_mock.assert_called_with([
            DHCP_CLI_HELPER_PATH,
            "--mac", str(dhcp_desc_fixture.mac),
            "--vlan", str(dhcp_desc_fixture.vlan),
            "--interface", ip_allocator_dhcp_fixture._iface,
            "--json",
            "allocate",
        ],
            capture_output=True
        )


@pytest.fixture
def ip_allocator_block_fixture(ip_allocator_fixture):
    networks = [IPv4Network(IP_NETWORK), IPv4Network(IP_NETWORK_2)]
    for network in networks:
        ip_allocator_fixture._store.assigned_ip_blocks.add(network)

    return ip_allocator_fixture


@patch("subprocess.run", MagicMock())
def test_remove_ip_block(ip_allocator_block_fixture: IPAllocatorDHCP) -> None:
    ip_allocator_block_fixture.start_monitor_thread()

    actual_removed = ip_allocator_block_fixture.remove_ip_blocks(
        ipblocks=[IPv4Network(IP_NETWORK)],
        force=False,
    )

    actual_remain = ip_allocator_block_fixture._store.assigned_ip_blocks

    assert set(actual_removed) == set([IPv4Network(IP_NETWORK)])
    assert set(actual_remain) == set([IPv4Network(IP_NETWORK_2)])


@patch("subprocess.run", MagicMock())
def test_keep_ip_block_with_allocated_ip(
        ip_allocator_block_fixture: IPAllocatorDHCP,
        ip_desc_fixture: IPDesc,
) -> None:
    ip_allocator_block_fixture._store.ip_state_map.add_ip_to_state(
        ip=IPv4Address(IP),
        ip_desc=ip_desc_fixture,
        state=IPState.ALLOCATED,
    )
    ip_allocator_block_fixture.start_monitor_thread()

    actual_removed = ip_allocator_block_fixture.remove_ip_blocks(
        ipblocks=[IPv4Network(IP_NETWORK)],
        force=False,
    )

    actual_remain = ip_allocator_block_fixture._store.assigned_ip_blocks

    assert set(actual_removed) == set()
    assert set(actual_remain) == set([IPv4Network(IP_NETWORK), IPv4Network(IP_NETWORK_2)])


@patch("subprocess.run", MagicMock())
def test_force_remove_ip_block_with_allocated_ip(
        ip_allocator_block_fixture: IPAllocatorDHCP,
        ip_desc_fixture: IPDesc,
) -> None:
    ip_allocator_block_fixture._store.ip_state_map.add_ip_to_state(
        ip=IPv4Address(IP),
        ip_desc=ip_desc_fixture,
        state=IPState.ALLOCATED,
    )
    ip_allocator_block_fixture.start_monitor_thread()

    removed_blocks = ip_allocator_block_fixture.remove_ip_blocks(
        ipblocks=[IPv4Network(IP_NETWORK)],
        force=True,
    )
    remaining_blocks = ip_allocator_block_fixture._store.assigned_ip_blocks

    assert set(removed_blocks) == set([IPv4Network(IP_NETWORK)])
    assert set(remaining_blocks) == set([IPv4Network(IP_NETWORK_2)])


def create_subprocess_mock_dhcp_return() -> MagicMock:
    m = MagicMock()
    m.returncode = 0
    m.stdout = """{"ip": "%s","subnet": "%s","server_ip": "5.6.7.8","lease_expiration_time": "4"}""" % (IP, IP_NETWORK)
    return m


def test_allocate_ip_address(
        ip_allocator_fixture: IPAllocatorDHCP,
        ip_desc_fixture: IPDesc,
        dhcp_desc_fixture: DHCPDescriptor,
) -> None:
    ip_allocator_fixture.start_monitor_thread()

    with patch("subprocess.run", return_value=create_subprocess_mock_dhcp_return()) as subprocess_mock:
        actual_ip_desc = ip_allocator_fixture.alloc_ip_address(
            sid=SID,
            vlan=int(VLAN),
        )

        subprocess_mock.assert_called_once()
        subprocess_mock.assert_called_with([
            DHCP_CLI_HELPER_PATH,
            "--mac", str(dhcp_desc_fixture.mac),
            "--vlan", str(dhcp_desc_fixture.vlan),
            "--interface", ip_allocator_fixture._iface,
            "--json",
            "allocate",
        ],
            capture_output=True
        )

    assert actual_ip_desc == ip_desc_fixture

    actual_added_ip_blocks = ip_allocator_fixture.list_added_ip_blocks()
    assert actual_added_ip_blocks == [IPv4Network(IP_NETWORK)]

    actual_allocated_ips = ip_allocator_fixture.list_allocated_ips(ipblock=IPv4Network(IP_NETWORK))
    assert actual_allocated_ips == [actual_ip_desc.ip]
