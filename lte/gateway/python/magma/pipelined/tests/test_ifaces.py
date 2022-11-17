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
from magma.pipelined.ifaces import get_mac_address_from_iface

# Prevent flakiness due to prometheus library import
sys.modules["magma.pipelined.metrics"] = MagicMock()


@patch("magma.pipelined.ifaces.netifaces")
def test_get_mac_address_from_iface(netifaces_mock):
    netifaces_mock.AF_LINK = 13
    netifaces_mock.ifaddresses.return_value = {netifaces_mock.AF_LINK: [{"addr": "00:11:22:33:44:55"}]}

    assert get_mac_address_from_iface("eth0") == "00:11:22:33:44:55"


def test_get_mac_address_from_iface_invalid():
    with pytest.raises(ValueError):
        get_mac_address_from_iface("some_invalid_interface")
