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

# Prevent flakiness due to prometheus library import
sys.modules["magma.pipelined.metrics"] = MagicMock()

from magma.pipelined.ifaces import get_mac_address


@patch("magma.pipelined.ifaces.get_mac_address")
def test_get_mac_address(get_mac_address_mock):
    get_mac_address_mock.return_value = "00:11:22:33:44:55"
    assert get_mac_address(interface="eth0") == "00:11:22:33:44:55"


@patch("magma.pipelined.ifaces.get_mac_address")
def test_get_mac_address_invalid(get_mac_address_mock):
    get_mac_address_mock.return_value = None
    with pytest.raises(ValueError):
        get_mac_address(interface="eth0")
