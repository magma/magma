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
import os
from datetime import datetime, timedelta
from ipaddress import IPv4Address, IPv4Network
from typing import Any, List
from unittest.mock import MagicMock, patch

import fakeredis
import freezegun
import pytest
from magma.mobilityd.dhcp_desc import DHCPDescriptor, DHCPState
from magma.mobilityd.ip_allocator_dhcp import DHCP_HELPER_CLI, IPAllocatorDHCP
from magma.mobilityd.ip_descriptor import IPDesc, IPState, IPType
from magma.mobilityd.mac import MacAddress, sid_to_mac
from magma.mobilityd.mobility_store import MobilityStore

SID = "IMSI123456789"
MAC = MacAddress(sid_to_mac(SID).lower())
MAC2 = MacAddress("01:23:45:67:89:ab")
IP = "1.2.3.4"
SERVER_IP = "5.6.7.8"
ROUTER_IP = "11.22.33.44"
SUBNET = "24"
IP_NETWORK = "1.2.3.0/" + SUBNET
IP_NETWORK_2 = "1.2.4.0/" + SUBNET
VLAN = "0"
LEASE_EXPIRATION_TIME = 4
FROZEN_TEST_TIME = "2021-01-01"
TMP_FILE = "/tmp/tmpfile"


@pytest.fixture
def dhcp_desc_fixture() -> DHCPDescriptor:
    with freezegun.freeze_time(FROZEN_TEST_TIME):
        return DHCPDescriptor(
            mac=MAC,
            ip=IP_NETWORK,
            vlan=VLAN,
            state=DHCPState.ACK,
            state_requested=DHCPState.REQUEST,
            lease_expiration_time=LEASE_EXPIRATION_TIME,
        )


@pytest.fixture
def ip_allocator_fixture() -> IPAllocatorDHCP:
    client = fakeredis.FakeStrictRedis()

    ip_allocator = IPAllocatorDHCP(
        store=MobilityStore(client),
        lease_renew_wait_min=LEASE_EXPIRATION_TIME,
        start=False,
    )

    yield ip_allocator

    ip_allocator._monitor_thread_event.set()
    ip_allocator._monitor_thread.join()


@pytest.fixture
def ip_allocator_dhcp_fixture(
    ip_allocator_fixture: IPAllocatorDHCP,
    dhcp_desc_fixture: DHCPDescriptor,
) -> IPAllocatorDHCP:
    ip_allocator_fixture._store.dhcp_store[
        dhcp_desc_fixture.mac.as_redis_key(dhcp_desc_fixture.vlan)
    ] = dhcp_desc_fixture

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


def run_dhcp_allocator_thread(
    frozen_datetime: Any,
    ip_allocator_dhcp_fixture: IPAllocatorDHCP,
    freeze_time: float,
) -> None:
    ip_allocator_dhcp_fixture._monitor_thread_event.set()
    frozen_datetime.tick(timedelta(seconds=freeze_time))
    ip_allocator_dhcp_fixture.start_monitor_thread()
    ip_allocator_dhcp_fixture._monitor_thread.join()


@patch("tempfile.NamedTemporaryFile")
def test_allocate_ip_address(
    mock_tempfile: MagicMock,
    ip_allocator_fixture: IPAllocatorDHCP,
    ip_desc_fixture: IPDesc,
    dhcp_desc_fixture: DHCPDescriptor,
) -> None:
    mock_tempfile.return_value.__enter__.return_value.name = TMP_FILE
    mock_tempfile.return_value.__exit__.side_effect = lambda *args: os.remove(TMP_FILE)
    ip_allocator_fixture.start_monitor_thread()
    call_args = [
        DHCP_HELPER_CLI,
        "--mac", str(dhcp_desc_fixture.mac),
        "--vlan", str(dhcp_desc_fixture.vlan),
        "--interface", ip_allocator_fixture._iface,
        "--save-file", TMP_FILE,
        "--json",
        "allocate",
    ]

    with freezegun.freeze_time(FROZEN_TEST_TIME):
        with patch(
            "subprocess.run",
            return_value=create_subprocess_mock_dhcp_return(),
            side_effect=create_subprocess_mock_json_file,
        ) as subprocess_mock:
            reference_time = datetime.now()
            actual_ip_desc = ip_allocator_fixture.alloc_ip_address(
                sid=SID,
                vlan_id=int(VLAN),
            )
            _assert_calls_and_deadlines(
                advance_time=0,
                call_args=call_args,
                ip_allocator=ip_allocator_fixture,
                reference_time=reference_time,
                subprocess_mock=subprocess_mock,
            )

    assert actual_ip_desc == ip_desc_fixture

    actual_added_ip_blocks = ip_allocator_fixture.list_added_ip_blocks()
    assert actual_added_ip_blocks == [IPv4Network(IP_NETWORK)]

    actual_allocated_ips = ip_allocator_fixture.list_allocated_ips(ipblock=IPv4Network(IP_NETWORK))
    assert actual_allocated_ips == [actual_ip_desc.ip]


def test_no_renewal_of_ip(ip_allocator_dhcp_fixture: IPAllocatorDHCP) -> None:
    advance_time = 1
    with freezegun.freeze_time(FROZEN_TEST_TIME) as frozen_datetime:
        with patch("subprocess.run", return_value=create_subprocess_mock_dhcp_return()) as subprocess_mock:
            run_dhcp_allocator_thread(
                frozen_datetime=frozen_datetime,
                ip_allocator_dhcp_fixture=ip_allocator_dhcp_fixture,
                freeze_time=advance_time,
            )
            subprocess_mock.assert_not_called()


@patch("tempfile.NamedTemporaryFile")
def test_renewal_of_ip(
    mock_tempfile: MagicMock,
    ip_allocator_dhcp_fixture: IPAllocatorDHCP,
) -> None:
    mock_tempfile.return_value.__enter__.return_value.name = TMP_FILE
    mock_tempfile.return_value.__exit__.side_effect = lambda *args: os.remove(TMP_FILE)

    dhcp_desc = list(ip_allocator_dhcp_fixture._store.dhcp_store.values())[0]
    call_args = [
        DHCP_HELPER_CLI,
        "--mac", str(dhcp_desc.mac),
        "--vlan", str(dhcp_desc.vlan),
        "--interface", ip_allocator_dhcp_fixture._iface,
        "--save-file", TMP_FILE,
        "--json",
        "renew",
        "--ip", str(dhcp_desc.ip),
        "--server-ip", str(dhcp_desc.server_ip),
    ]

    _run_allocator_and_assert(
        advance_time=3,
        call_args=call_args,
        ip_allocator_dhcp_fixture=ip_allocator_dhcp_fixture,
    )


@patch("tempfile.NamedTemporaryFile")
def test_allocate_ip_after_expiry(
    mock_tempfile: MagicMock,
    ip_allocator_dhcp_fixture: IPAllocatorDHCP,
) -> None:
    mock_tempfile.return_value.__enter__.return_value.name = TMP_FILE
    mock_tempfile.return_value.__exit__.side_effect = lambda *args: os.remove(TMP_FILE)

    dhcp_desc = list(ip_allocator_dhcp_fixture._store.dhcp_store.values())[0]
    call_args = [
        DHCP_HELPER_CLI,
        "--mac", str(dhcp_desc.mac),
        "--vlan", str(dhcp_desc.vlan),
        "--interface", ip_allocator_dhcp_fixture._iface,
        "--save-file", TMP_FILE,
        "--json",
        "allocate",
    ]
    _run_allocator_and_assert(
        advance_time=5,
        call_args=call_args,
        ip_allocator_dhcp_fixture=ip_allocator_dhcp_fixture,
    )


def _run_allocator_and_assert(
        advance_time: int, call_args: List[str], ip_allocator_dhcp_fixture: IPAllocatorDHCP,
) -> None:
    with freezegun.freeze_time(FROZEN_TEST_TIME) as frozen_datetime:
        with patch(
                "subprocess.run", return_value=create_subprocess_mock_dhcp_return(),
                side_effect=create_subprocess_mock_json_file,
        ) as subprocess_mock:
            reference_time = datetime.now()
            run_dhcp_allocator_thread(
                frozen_datetime=frozen_datetime,
                ip_allocator_dhcp_fixture=ip_allocator_dhcp_fixture,
                freeze_time=advance_time,
            )
            _assert_calls_and_deadlines(
                advance_time=advance_time,
                call_args=call_args,
                ip_allocator=ip_allocator_dhcp_fixture,
                reference_time=reference_time,
                subprocess_mock=subprocess_mock,
            )


def _assert_calls_and_deadlines(
        advance_time: int, call_args: List[str], ip_allocator: IPAllocatorDHCP,
        reference_time: datetime, subprocess_mock: MagicMock,
) -> None:
    subprocess_mock.assert_called_once()
    subprocess_mock.assert_called_with(
        call_args,
        capture_output=True,
        check=False,
    )
    dhcp_desc = list(ip_allocator._store.dhcp_store.values())[0]
    expected_lease_expiration_time = reference_time + timedelta(seconds=advance_time + LEASE_EXPIRATION_TIME)
    expected_lease_renew_deadline = reference_time + timedelta(seconds=advance_time + LEASE_EXPIRATION_TIME / 2)
    assert dhcp_desc.lease_expiration_time == expected_lease_expiration_time
    assert dhcp_desc.lease_renew_deadline == expected_lease_renew_deadline
    assert not os.path.exists(call_args[8])


@pytest.fixture
def ip_allocator_block_fixture(ip_allocator_fixture: IPAllocatorDHCP) -> IPAllocatorDHCP:
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

    assert set(actual_removed) == {IPv4Network(IP_NETWORK)}
    assert set(actual_remain) == {IPv4Network(IP_NETWORK_2)}


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
    assert set(actual_remain) == {IPv4Network(IP_NETWORK), IPv4Network(IP_NETWORK_2)}


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

    assert set(removed_blocks) == {IPv4Network(IP_NETWORK)}
    assert set(remaining_blocks) == {IPv4Network(IP_NETWORK_2)}


def create_subprocess_mock_dhcp_return() -> MagicMock:
    m = MagicMock()
    m.returncode = 0
    m.stdout = """{"ip": "%s","subnet": "%s","server_ip": "%s", "router_ip": "%s","lease_expiration_time": %s}""" % (
        IP, IP_NETWORK, SERVER_IP, ROUTER_IP, LEASE_EXPIRATION_TIME,
    )
    return m


def create_subprocess_mock_json_file(call_args: List[str], capture_output=True, check=False) -> MagicMock:
    with open(call_args[8], "w") as f:
        f.write(
            """{"ip": "%s","subnet": "%s","server_ip": "%s", "router_ip": "%s","lease_expiration_time": %s}"""
            % (IP, IP_NETWORK, SERVER_IP, ROUTER_IP, LEASE_EXPIRATION_TIME),
        )
    return create_subprocess_mock_dhcp_return()
