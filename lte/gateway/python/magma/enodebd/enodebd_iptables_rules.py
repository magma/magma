#!/usr/bin/env python3

"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""


import asyncio
import re
import shlex
import subprocess
from typing import List

from magma.common.misc_utils import (
    IpPreference,
    get_if_ip_with_netmask,
    get_ip_from_if,
)
from magma.configuration.service_configs import load_service_config
from magma.enodebd.logger import EnodebdLogger as logger

IPTABLES_RULE_FMT = """sudo iptables -t nat
    -{add} PREROUTING
    -d {public_ip}
    -p tcp
    --dport {port}
    -j DNAT --to-destination {private_ip}"""

EXPECTED_IP4 = ('192.168.60.142', '10.0.2.1')
EXPECTED_MASK = '255.255.255.0'


def get_iptables_rule(port, enodebd_public_ip, private_ip, add=True):
    return IPTABLES_RULE_FMT.format(
        add='A' if add else 'D',
        public_ip=enodebd_public_ip,
        port=port,
        private_ip=private_ip,
    )


def does_iface_config_match_expected(ip: str, netmask: str) -> bool:
    return ip in EXPECTED_IP4 and netmask == EXPECTED_MASK


def _get_prerouting_rules(output: str) -> List[str]:
    prerouting_rules = output.split('\n\n')[0]
    prerouting_rules = prerouting_rules.split('\n')
    # Skipping the first two lines since it contains only column names
    prerouting_rules = prerouting_rules[2:]
    return prerouting_rules


async def check_and_apply_iptables_rules(
    port: str,
    enodebd_public_ip: str,
    enodebd_ip: str,
) -> None:
    command = 'sudo iptables -t nat -L'
    output = subprocess.run(
        command, shell=True,
        stdout=subprocess.PIPE, check=True,
    )
    command_output = output.stdout.decode('utf-8').strip()
    prerouting_rules = _get_prerouting_rules(command_output)
    if not prerouting_rules:
        logger.info('Configuring Iptables rule')
        await run(
            get_iptables_rule(
                port,
                enodebd_public_ip,
                enodebd_ip,
                add=True,
            ),
        )
    else:
        # Checks each rule in PREROUTING Chain
        check_rules(prerouting_rules, port, enodebd_public_ip, enodebd_ip)


def check_rules(
    prerouting_rules: List[str],
    port: str,
    enodebd_public_ip: str,
    private_ip: str,
) -> None:
    unexpected_rules = []
    pattern = r'DNAT\s+tcp\s+--\s+anywhere\s+{pub_ip}\s+tcp\s+dpt:{dport} to:{ip}'.format(
        pub_ip=enodebd_public_ip,
        dport=port,
        ip=private_ip,
    )
    for rule in prerouting_rules:
        match = re.search(pattern, rule)
        if not match:
            unexpected_rules.append(rule)
    if unexpected_rules:
        logger.warning('The following Prerouting rule(s) are unexpected')
        for rule in unexpected_rules:
            logger.warning(rule)


async def run(cmd):
    """Fork shell and run command NOTE: Popen is non-blocking"""
    cmd = shlex.split(cmd)
    proc = await asyncio.create_subprocess_shell(" ".join(cmd))
    await proc.communicate()
    if proc.returncode != 0:
        # This can happen because the NAT prerouting rule didn't exist
        logger.error(
            'Possible error running async subprocess: %s exited with '
            'return code [%d].', cmd, proc.returncode,
        )
    return proc.returncode


async def set_enodebd_iptables_rule():
    """
    Remove & Set iptable rules for exposing public IP
    for enobeb instead of private IP..
    """
    # Remove & Set iptable rules for exposing public ip
    # for enobeb instead of private
    cfg = load_service_config('enodebd')
    port, interface = cfg['tr069']['port'], cfg['tr069']['interface']
    enodebd_public_ip = cfg['tr069']['public_ip']
    # IPv4 only as iptables only works for IPv4. TODO: Investigate ip6tables?
    enodebd_ip = get_ip_from_if(interface, preference=IpPreference.IPV4_ONLY)
    # Incoming data from 192.88.99.142 -> enodebd address (eg 192.168.60.142)
    enodebd_netmask = get_if_ip_with_netmask(
        interface,
        preference=IpPreference.IPV4_ONLY,
    )[1]
    verify_config = does_iface_config_match_expected(
        enodebd_ip,
        enodebd_netmask,
    )
    if not verify_config:
        logger.warning(
            'The IP address of the %s interface is %s. The '
            'expected IP addresses are %s',
            interface, enodebd_ip, str(EXPECTED_IP4),
        )
    await check_and_apply_iptables_rules(
        port,
        enodebd_public_ip,
        enodebd_ip,
    )


if __name__ == '__main__':
    set_enodebd_iptables_rule()
