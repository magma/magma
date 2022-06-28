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
    load_service_config,
    save_override_config,
)

return_codes = Enum(
    "return_codes", "IPV4_ENABLED IPV6_ENABLED INVALID", start=0,
)
IPV6_IFACE_CONFIGS = [
    ("pipelined", "nat_iface"),
    ("pipelined", "uplink_eth_port_name"),
]


def check_ipv6_iface_config(service, config_name, config_value):
    """Check if eth3 interface is configured"""
    service_config = load_service_config(service)
    if service_config.get(config_name) == config_value:
        print("IPV6_ENABLED\t%s -> %s" % (service, config_name))
        return return_codes.IPV6_ENABLED
    return return_codes.IPV4_ENABLED


def check_ipv6_iface():
    """Check the interface configured"""
    ipv6_enabled = 0
    for service, config in IPV6_IFACE_CONFIGS:
        if (
            check_ipv6_iface_config(service, config, "eth3")
            == return_codes.IPV6_ENABLED
        ):
            ipv6_enabled += 1

    if ipv6_enabled == 0:
        res = return_codes.IPV4_ENABLED
    elif ipv6_enabled == len(IPV6_IFACE_CONFIGS):
        res = return_codes.IPV6_ENABLED
    else:
        res = return_codes.INVALID

    print("Check returning", res)
    return res


def enable_eth3_iface():
    """Enable eth3 interface as nat_iface"""
    if check_ipv6_iface() == return_codes.IPV6_ENABLED:
        print("eth3 interface is already enabled")
        sys.exit(return_codes.IPV6_ENABLED.value)
    for service, config in IPV6_IFACE_CONFIGS:
        cfg = load_override_config(service) or {}
        cfg[config] = "eth3"
        save_override_config(service, cfg)


def disable_eth3_iface():
    """Disable eth3 interface as nat_iface"""
    if check_ipv6_iface() == return_codes.IPV4_ENABLED:
        print("IPv4 is already enabled")
        sys.exit(return_codes.IPV4_ENABLED.value)
    for service, config in IPV6_IFACE_CONFIGS:
        cfg = load_override_config(service) or {}
        cfg[config] = "eth2"
        save_override_config(service, cfg)


ETH3_IFACE_FUNC_DICT = {
    "enable": enable_eth3_iface,
    "disable": disable_eth3_iface,
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
