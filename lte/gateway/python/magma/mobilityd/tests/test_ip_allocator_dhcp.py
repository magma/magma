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
from magma.mobilityd.mac import sid_to_mac
from magma.mobilityd.mobility_store import AssignedIpBlocksSet, ip_states, defaultdict_key

SID = "IMSI123456789"
MAC = sid_to_mac(SID).lower()
MAC2 = "01:23:45:67:89:ab"
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
def ip_allocator_fixture(dhcp_desc_fixture: DHCPDescriptor) -> IPAllocatorDHCP:
    store = MagicMock()
    store.dhcp_store = MagicMock()
    store.dhcp_store.values.return_value = [dhcp_desc_fixture]

    ip_allocator = IPAllocatorDHCP(
        store=store,
        lease_renew_wait_min=4,
        start=False,
    )
    yield ip_allocator

    ip_allocator._monitor_thread_event.set()
    ip_allocator._monitor_thread.join()


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


def get_mock_run() -> MagicMock:
    m = MagicMock()
    m.returncode = 0
    m.stdout = """{"lease_expiration_time": 4}"""
    return m


def setup_allocate_test(
        frozen_datetime: Any, ip_allocator_fixture: IPAllocatorDHCP,
        mock_run: MagicMock, freeze_time: float) -> None:
    mock_run.return_value = get_mock_run()
    ip_allocator_fixture._monitor_thread_event.set()
    frozen_datetime.tick(timedelta(seconds=freeze_time))
    ip_allocator_fixture.start_monitor_thread()
    ip_allocator_fixture._monitor_thread.join()


@patch("subprocess.run")
def test_no_renewal_of_ip(mock_run: MagicMock, ip_allocator_fixture: IPAllocatorDHCP) -> None:
    with freezegun.freeze_time(FROZEN_TEST_TIME) as frozen_datetime:
        setup_allocate_test(
            frozen_datetime=frozen_datetime,
            ip_allocator_fixture=ip_allocator_fixture,
            mock_run=mock_run,
            freeze_time=1
        )
        mock_run.assert_not_called()


@patch("subprocess.run")
def test_renewal_of_ip(
        mock_run: MagicMock, ip_allocator_fixture: IPAllocatorDHCP,
        dhcp_desc_fixture: DHCPDescriptor) -> None:
    with freezegun.freeze_time(FROZEN_TEST_TIME) as frozen_datetime:
        setup_allocate_test(
            frozen_datetime=frozen_datetime,
            ip_allocator_fixture=ip_allocator_fixture,
            mock_run=mock_run,
            freeze_time=3
        )
        mock_run.assert_called_once()
        mock_run.assert_called_with([
            DHCP_CLI_HELPER_PATH,
            "--mac", str(dhcp_desc_fixture.mac),
            "--vlan", str(dhcp_desc_fixture.vlan),
            "--interface", ip_allocator_fixture._iface,
            "--json",
            "renew",
            "--ip", str(dhcp_desc_fixture.ip),
            "--server-ip", str(dhcp_desc_fixture.server_ip),
        ],
            capture_output=True
        )


@patch("subprocess.run")
def test_allocate_ip_after_expiry(
        mock_run: MagicMock, ip_allocator_fixture: IPAllocatorDHCP,
        dhcp_desc_fixture: DHCPDescriptor) -> None:
    with freezegun.freeze_time(FROZEN_TEST_TIME) as frozen_datetime:
        setup_allocate_test(
            frozen_datetime=frozen_datetime,
            ip_allocator_fixture=ip_allocator_fixture,
            mock_run=mock_run,
            freeze_time=5
        )
        mock_run.assert_called_once()
        mock_run.assert_called_with([
            DHCP_CLI_HELPER_PATH,
            "--mac", str(dhcp_desc_fixture.mac),
            "--vlan", str(dhcp_desc_fixture.vlan),
            "--interface", ip_allocator_fixture._iface,
            "--json",
            "allocate",
        ],
            capture_output=True
        )


def setup_remove_ip_block_test(
        ip_allocator_fixture: IPAllocatorDHCP, mock_run: MagicMock,
        ip_networks: List[IPv4Network]) -> None:
    mock_run.return_value = get_mock_run()
    client = fakeredis.FakeStrictRedis()
    ip_allocator_fixture._store.assigned_ip_blocks = AssignedIpBlocksSet(client)
    for network in ip_networks:
        ip_allocator_fixture._store.assigned_ip_blocks.add(network)
    ip_allocator_fixture._store.ip_state_map = IpDescriptorMap(
        defaultdict_key(lambda key: ip_states(client, key)),  # type: ignore[arg-type]
    )


@patch("subprocess.run")
def test_remove_ip_block(mock_run: MagicMock, ip_allocator_fixture: IPAllocatorDHCP) -> None:
    networks = [IPv4Network(IP_NETWORK), IPv4Network(IP_NETWORK_2)]
    setup_remove_ip_block_test(
        ip_allocator_fixture=ip_allocator_fixture,
        mock_run=mock_run,
        ip_networks=networks)
    ip_allocator_fixture.start_monitor_thread()

    actual_removed = ip_allocator_fixture.remove_ip_blocks(
        ipblocks=[IPv4Network(IP_NETWORK)],
        force=False,
    )
    expected_removed = [IPv4Network(IP_NETWORK)]

    actual_remain = ip_allocator_fixture._store.assigned_ip_blocks
    expected_remain = AssignedIpBlocksSet(fakeredis.FakeStrictRedis())
    expected_remain.add(IPv4Network(IP_NETWORK_2))

    assert_remove_blocks(actual_remain=actual_remain, expected_remain=expected_remain,
                         actual_removed=actual_removed, expected_removed=expected_removed)


@patch("subprocess.run")
def test_keep_ip_block_with_allocated_ip(
        mock_run: MagicMock, ip_allocator_fixture: IPAllocatorDHCP,
        ip_desc_fixture: IPDesc,
) -> None:
    networks = [IPv4Network(IP_NETWORK)]
    setup_remove_ip_block_test(
        ip_allocator_fixture=ip_allocator_fixture,
        mock_run=mock_run,
        ip_networks=networks)
    ip_allocator_fixture._store.ip_state_map.add_ip_to_state(
        ip=IPv4Address(IP),
        ip_desc=ip_desc_fixture,
        state=IPState.ALLOCATED,
    )
    ip_allocator_fixture.start_monitor_thread()

    actual_removed = ip_allocator_fixture.remove_ip_blocks(
        ipblocks=[IPv4Network(IP_NETWORK)],
        force=False,
    )
    expected_removed = []

    actual_remain = ip_allocator_fixture._store.assigned_ip_blocks
    expected_remain = AssignedIpBlocksSet(fakeredis.FakeStrictRedis())
    expected_remain.add(IPv4Network(IP_NETWORK))

    assert_remove_blocks(actual_remain=actual_remain, expected_remain=expected_remain,
                         actual_removed=actual_removed, expected_removed=expected_removed)


@patch("subprocess.run")
def test_force_remove_ip_block_with_allocated_ip(
        mock_run: MagicMock, ip_allocator_fixture: IPAllocatorDHCP,
        ip_desc_fixture: IPDesc,
) -> None:
    networks = [IPv4Network(IP_NETWORK), IPv4Network(IP_NETWORK_2)]
    setup_remove_ip_block_test(
        ip_allocator_fixture=ip_allocator_fixture,
        mock_run=mock_run,
        ip_networks=networks)
    ip_allocator_fixture._store.ip_state_map.add_ip_to_state(
        ip=IPv4Address(IP),
        ip_desc=ip_desc_fixture,
        state=IPState.ALLOCATED,
    )
    ip_allocator_fixture.start_monitor_thread()

    actual_removed = ip_allocator_fixture.remove_ip_blocks(
        ipblocks=[IPv4Network(IP_NETWORK)],
        force=True,
    )
    expected_removed = [IPv4Network(IP_NETWORK)]

    actual_remain = ip_allocator_fixture._store.assigned_ip_blocks
    expected_remain = AssignedIpBlocksSet(fakeredis.FakeStrictRedis())
    expected_remain.add(IPv4Network(IP_NETWORK_2))

    assert_remove_blocks(actual_remain=actual_remain, expected_remain=expected_remain,
                         actual_removed=actual_removed, expected_removed=expected_removed)


def assert_remove_blocks(
        actual_remain: AssignedIpBlocksSet, expected_remain: AssignedIpBlocksSet,
        actual_removed: List[IPv4Network], expected_removed: List[IPv4Network]) -> None:
    assert expected_removed == actual_removed
    assert set(expected_remain) == set(actual_remain)


def get_mock_run_dhcp_return() -> MagicMock:
    m = MagicMock()
    m.returncode = 0
    m.stdout = """{"ip": %s,"subnet": %s,"server_ip": "5.6.7.8","lease_expiration_time": "4"}""" % (f""" "{IP}" """, f""" "{IP_NETWORK}" """)
    return m


def setup_allocate_ip_test(
        ip_allocator_fixture: IPAllocatorDHCP, mock_run: MagicMock) -> None:
    mock_run.return_value = get_mock_run_dhcp_return()
    client = fakeredis.FakeStrictRedis()
    ip_allocator_fixture._store.assigned_ip_blocks = AssignedIpBlocksSet(client)


@patch("subprocess.run")
def test_allocate_ip_address(
        mock_run: MagicMock,
        ip_allocator_fixture: IPAllocatorDHCP,
        ip_desc_fixture: IPDesc,
        dhcp_desc_fixture: DHCPDescriptor,
) -> None:
    setup_allocate_ip_test(
        ip_allocator_fixture=ip_allocator_fixture, mock_run=mock_run
    )

    actual_ip_desc = ip_allocator_fixture.alloc_ip_address(
        sid=SID,
        vlan=int(VLAN),
    )
    expected_ip_desc = ip_desc_fixture

    assert_ip_desc(actual_ip_desc=actual_ip_desc, expected_ip_desc=expected_ip_desc)
    mock_run.assert_called_once()
    mock_run.assert_called_with([
        DHCP_CLI_HELPER_PATH,
        "--mac", str(dhcp_desc_fixture.mac),
        "--vlan", str(dhcp_desc_fixture.vlan),
        "--interface", ip_allocator_fixture._iface,
        "--json",
        "allocate",
    ],
        capture_output=True
    )
    actual_added_ip_blocks = ip_allocator_fixture.list_added_ip_blocks()
    expected_added_ip_blocks = [IPv4Network(IP_NETWORK)]
    assert actual_added_ip_blocks == expected_added_ip_blocks
    actual_allocated_ips = ip_allocator_fixture.list_allocated_ips(ipblock=IPv4Network(IP_NETWORK))
    expected_allocated_ips = actual_ip_desc.ip
    print(IPv4Network(IP_NETWORK))
    for ip in ip_allocator_fixture._store.ip_state_map.list_ips(IPState.ALLOCATED):
        print(ip)
    assert actual_allocated_ips == expected_allocated_ips


def assert_ip_desc(actual_ip_desc, expected_ip_desc):
    assert actual_ip_desc.ip == expected_ip_desc.ip
    assert actual_ip_desc.state == expected_ip_desc.state
    assert actual_ip_desc.sid == expected_ip_desc.sid
    assert actual_ip_desc.ip_block == expected_ip_desc.ip_block
    assert actual_ip_desc.type == expected_ip_desc.type
    assert actual_ip_desc.vlan_id == expected_ip_desc.vlan_id
