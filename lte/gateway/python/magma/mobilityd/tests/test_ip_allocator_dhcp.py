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
import time
from datetime import datetime, timedelta
from ipaddress import IPv4Network, IPv4Address
from unittest.mock import patch, MagicMock
import unittest
import fakeredis
import pytest
import freezegun

from magma.mobilityd.dhcp_desc import DHCPDescriptor, DHCPState
from magma.mobilityd.ip_allocator_dhcp import IPAllocatorDHCP, DHCP_CLI_HELPER_PATH
from magma.mobilityd.ip_descriptor import IPState, IPDesc, IPType
from magma.mobilityd.ip_descriptor_map import IpDescriptorMap
from magma.mobilityd.mobility_store import AssignedIpBlocksSet, ip_states, defaultdict_key
from magma.mobilityd.mac import sid_to_mac

SID = "IMSI123456789"
MAC = sid_to_mac(SID)
# MAC = "01:23:45:67:89:ab"
MAC2 = "01:23:45:67:89:cd"
IP = "1.2.3.4"
IP_NETWORK = "1.2.3.0/24"
IP_NETWORK_2 = "1.2.4.0/24"
VLAN = "0"
LEASE_EXPIRATION_TIME = 10
FROZEN_TEST_TIME = "2021-01-01"


@pytest.fixture
def dhcp_desc_fixture():
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
def ipallocator_fixture(dhcp_desc_fixture):
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
def ip_desc_fixture():
    return IPDesc(
        ip=IPv4Address(IP),
        state=IPState.ALLOCATED,
        vlan_id=int(VLAN),
        ip_block=IPv4Network(IP_NETWORK),
        ip_type=IPType.DHCP,
        sid=SID,
    )


def get_mock_run():
    m = MagicMock()
    m.returncode = 0
    m.stdout = """{"lease_expiration_time": 4}"""
    return m


@patch("subprocess.run")
def test_no_renewal_of_ip(mock_run, ipallocator_fixture):
    with freezegun.freeze_time(FROZEN_TEST_TIME) as frozen_datetime:
        mock_run.return_value = get_mock_run()
        ipallocator_fixture._monitor_thread_event.set()
        mock_run.assert_not_called()
        frozen_datetime.tick(timedelta(seconds=1))
        ipallocator_fixture.start_monitor_thread()
        ipallocator_fixture._monitor_thread.join()
        mock_run.assert_not_called()


@patch("subprocess.run")
def test_renewal_of_ip(mock_run, ipallocator_fixture, dhcp_desc_fixture):
    with freezegun.freeze_time(FROZEN_TEST_TIME) as frozen_datetime:
        mock_run.return_value = get_mock_run()
        ipallocator_fixture._monitor_thread_event.set()
        mock_run.assert_not_called()
        frozen_datetime.tick(timedelta(seconds=3))
        ipallocator_fixture.start_monitor_thread()
        ipallocator_fixture._monitor_thread.join()
        mock_run.assert_called_once()
        mock_run.assert_called_with([
            DHCP_CLI_HELPER_PATH,
            "--mac", str(dhcp_desc_fixture.mac),
            "--vlan", str(dhcp_desc_fixture.vlan),
            "--interface", ipallocator_fixture._iface,
            "--json",
            "renew",
            "--ip", str(dhcp_desc_fixture.ip),
            "--server-ip", str(dhcp_desc_fixture.server_ip),
        ],
            capture_output=True
        )


@patch("subprocess.run")
def test_allocate_ip_after_expiry(mock_run, ipallocator_fixture, dhcp_desc_fixture):
    with freezegun.freeze_time(FROZEN_TEST_TIME) as frozen_datetime:
        mock_run.return_value = get_mock_run()
        ipallocator_fixture._monitor_thread_event.set()
        mock_run.assert_not_called()
        frozen_datetime.tick(timedelta(seconds=5))
        ipallocator_fixture.start_monitor_thread()
        ipallocator_fixture._monitor_thread.join()
        mock_run.assert_called_once()
        mock_run.assert_called_with([
            DHCP_CLI_HELPER_PATH,
            "--mac", str(dhcp_desc_fixture.mac),
            "--vlan", str(dhcp_desc_fixture.vlan),
            "--interface", ipallocator_fixture._iface,
            "--json",
            "allocate",
        ],
            capture_output=True
        )


@patch("subprocess.run")
def test_remove_ip_block(mock_run, ipallocator_fixture, dhcp_desc_fixture):
    mock_run.return_value = get_mock_run()
    client = fakeredis.FakeStrictRedis()
    ipallocator_fixture._store.assigned_ip_blocks = AssignedIpBlocksSet(client)
    ipallocator_fixture._store.assigned_ip_blocks.add(IPv4Network(IP_NETWORK))
    ipallocator_fixture._store.assigned_ip_blocks.add(IPv4Network(IP_NETWORK_2))

    def get_ip_states(key): return ip_states(client, key)
    ipallocator_fixture._store.ip_state_map = IpDescriptorMap(
        defaultdict_key(get_ip_states),  # type: ignore[arg-type]
    )
    ipallocator_fixture.start_monitor_thread()
    removed_block = ipallocator_fixture.remove_ip_blocks([IPv4Network(IP_NETWORK)])

    assert [IPv4Network(IP_NETWORK)] == removed_block
    expected = AssignedIpBlocksSet(fakeredis.FakeStrictRedis())
    expected.add(IPv4Network(IP_NETWORK_2))
    assert len(expected) == len(ipallocator_fixture._store.assigned_ip_blocks)
    for block in expected:
        assert block in ipallocator_fixture._store.assigned_ip_blocks
    print(expected)


@patch("subprocess.run")
def test_keep_ip_block_with_allocated_ip(
        mock_run, ipallocator_fixture,
        dhcp_desc_fixture, ip_desc_fixture,
):
    mock_run.return_value = get_mock_run()
    client = fakeredis.FakeStrictRedis()
    ipallocator_fixture._store.assigned_ip_blocks = AssignedIpBlocksSet(client)
    ipallocator_fixture._store.assigned_ip_blocks.add(IPv4Network(IP_NETWORK))

    def get_ip_states(key): return ip_states(client, key)
    ipallocator_fixture._store.ip_state_map = IpDescriptorMap(
        defaultdict_key(get_ip_states),  # type: ignore[arg-type]
    )
    ipallocator_fixture._store.ip_state_map.add_ip_to_state(
        ip=IPv4Address(IP),
        ip_desc=ip_desc_fixture,
        state=IPState.ALLOCATED,
    )
    ipallocator_fixture.start_monitor_thread()
    removed_block = ipallocator_fixture.remove_ip_blocks([IPv4Network(IP_NETWORK)])

    assert [] == removed_block
    expected = AssignedIpBlocksSet(fakeredis.FakeStrictRedis())
    expected.add(IPv4Network(IP_NETWORK))
    assert len(expected) == len(ipallocator_fixture._store.assigned_ip_blocks)
    for block in expected:
        assert block in ipallocator_fixture._store.assigned_ip_blocks
    print(expected)


@patch("subprocess.run")
def test_force_remove_ip_block_with_allocated_ip(
        mock_run, ipallocator_fixture,
        dhcp_desc_fixture, ip_desc_fixture,
):
    mock_run.return_value = get_mock_run()
    client = fakeredis.FakeStrictRedis()
    ipallocator_fixture._store.assigned_ip_blocks = AssignedIpBlocksSet(client)
    ipallocator_fixture._store.assigned_ip_blocks.add(IPv4Network(IP_NETWORK))
    ipallocator_fixture._store.assigned_ip_blocks.add(IPv4Network(IP_NETWORK_2))

    def get_ip_states(key): return ip_states(client, key)
    ipallocator_fixture._store.ip_state_map = IpDescriptorMap(
        defaultdict_key(get_ip_states),  # type: ignore[arg-type]
    )
    ipallocator_fixture._store.ip_state_map.add_ip_to_state(
        ip=IPv4Address(IP),
        ip_desc=ip_desc_fixture,
        state=IPState.ALLOCATED,
    )
    ipallocator_fixture.start_monitor_thread()
    removed_block = ipallocator_fixture.remove_ip_blocks(
        ipblocks=[IPv4Network(IP_NETWORK)],
        force=True,
    )

    assert [IPv4Network(IP_NETWORK)] == removed_block
    expected = AssignedIpBlocksSet(fakeredis.FakeStrictRedis())
    expected.add(IPv4Network(IP_NETWORK_2))
    assert len(expected) == len(ipallocator_fixture._store.assigned_ip_blocks)
    for block in expected:
        assert block in ipallocator_fixture._store.assigned_ip_blocks
    print(expected)
