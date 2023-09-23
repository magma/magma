#!/usr/bin/env python3

"""
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Script to trigger pre and post start commands for the Sctpd systemd unit
"""

import argparse
import sys
from enum import Enum

from magma.configuration.service_configs import (
    load_override_config,
    save_override_config,
)

PIPELINED_SERVICE = "pipelined"

IPV6_IFACE_CONFIGS = [
    "nat_iface",
    "uplink_eth_port_name",
]


def _enable_eth3_iface():
    """Enable eth3 interface as nat_iface"""
    _change_iface("eth3")


def _disable_eth3_iface():
    """Disable eth3 interface as nat_iface"""
    _change_iface("eth2")


def _change_iface(iface):
    for config in IPV6_IFACE_CONFIGS:
        cfg = load_override_config(PIPELINED_SERVICE) or {}
        cfg[config] = iface
        save_override_config(PIPELINED_SERVICE, cfg)


ETH3_IFACE_FUNC_DICT = {
    "enable": _enable_eth3_iface,
    "disable": _disable_eth3_iface,
}


def main():
    """Script to add eth3 iface as nat_iface in pipelined.yml"""
    parser = argparse.ArgumentParser()
    parser.add_argument("command", choices=ETH3_IFACE_FUNC_DICT.keys())

    args = parser.parse_args()

    func = ETH3_IFACE_FUNC_DICT[args.command]
    func()


if __name__ == "__main__":
    main()
