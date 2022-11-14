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

import sys
from unittest.mock import MagicMock, patch

import pytest

from magma.pipelined.ifaces import get_mac_address

# Prevent flakiness due to prometheus library import
sys.modules["magma.pipelined.metrics"] = MagicMock()


@patch("magma.pipelined.ifaces.netifaces")
def test_get_mac_address(netifaces_mock):
    netifaces_mock.AF_LINK = 13
    netifaces_mock.ifaddresses.return_value = {netifaces_mock.AF_LINK: [{"addr": "00:11:22:33:44:55"}]}

    assert get_mac_address(interface="eth0") == "00:11:22:33:44:55"


@patch("magma.pipelined.ifaces.netifaces")
def test_get_mac_address_invalid(netifaces_mock):
    netifaces_mock.AF_LINK = 13
    netifaces_mock.ifaddresses.return_value = {}
    with pytest.raises(ValueError):
        get_mac_address(interface="eth0")


@patch("magma.pipelined.ifaces.netifaces")
def test_get_mac_address_from_ip4(netifaces_mock):
    netifaces_mock.AF_LINK = 13
    netifaces_mock.AF_INET = 3
    netifaces_mock.interfaces.return_value = ["eth0"]
    netifaces_mock.ifaddresses.return_value = {
        netifaces_mock.AF_LINK: [{"addr": "00:11:22:33:44:55"}],
        netifaces_mock.AF_INET: [{"addr": "10.0.2.15"}],
    }

    assert get_mac_address(ip4="10.0.2.15") == "00:11:22:33:44:55"


@patch("magma.pipelined.ifaces.netifaces")
def test_get_mac_address_from_ip4_invalid(netifaces_mock):
    netifaces_mock.AF_LINK = 13
    netifaces_mock.AF_INET = 3
    netifaces_mock.interfaces.return_value = ["eth0"]
    netifaces_mock.ifaddresses.return_value = {}
    with pytest.raises(ValueError):
        get_mac_address(ip4="10.0.2.15")


@patch("magma.pipelined.ifaces.netifaces")
def test_get_mac_address_from_ip4_no_mac(netifaces_mock):
    netifaces_mock.AF_LINK = 13
    netifaces_mock.AF_INET = 3
    netifaces_mock.interfaces.return_value = ["eth0"]
    netifaces_mock.ifaddresses.return_value = {
        netifaces_mock.AF_INET: [{"addr": "10.0.2.15"}],
    }
    with pytest.raises(ValueError):
        get_mac_address(ip4="10.0.2.15")


@patch("magma.pipelined.ifaces.netifaces")
def test_get_mac_address_from_ip6(netifaces_mock):
    netifaces_mock.AF_LINK = 13
    netifaces_mock.AF_INET6 = 7
    netifaces_mock.interfaces.return_value = ["eth0"]
    netifaces_mock.ifaddresses.return_value = {
        netifaces_mock.AF_LINK: [{"addr": "00:11:22:33:44:55"}],
        netifaces_mock.AF_INET6: [{"addr": "fe80::5054:ff:fe12:3456%eth0"}],
    }

    assert get_mac_address(ip6="fe80::5054:ff:fe12:3456") == "00:11:22:33:44:55"


@patch("magma.pipelined.ifaces.netifaces")
def test_get_mac_address_from_ip6_invalid(netifaces_mock):
    netifaces_mock.AF_LINK = 13
    netifaces_mock.AF_INET6 = 7
    netifaces_mock.interfaces.return_value = ["eth0"]
    netifaces_mock.ifaddresses.return_value = {}
    with pytest.raises(ValueError):
        get_mac_address(ip6="fe80::64ba:b0ff:fe23:87f0")


@patch("magma.pipelined.ifaces.netifaces")
def test_get_mac_address_from_ip6_no_mac(netifaces_mock):
    netifaces_mock.AF_LINK = 13
    netifaces_mock.AF_INET6 = 7
    netifaces_mock.interfaces.return_value = ["eth0"]
    netifaces_mock.ifaddresses.return_value = {
        netifaces_mock.AF_INET6: [{"addr": "fe80::5054:ff:fe12:3456%eth0"}],
    }
    with pytest.raises(ValueError):
        get_mac_address(ip6="fe80::5054:ff:fe12:3456")
