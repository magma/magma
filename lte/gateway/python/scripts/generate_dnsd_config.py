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

import ipaddress
import logging

from generate_service_config import generate_template_config
from lte.protos.mconfig.mconfigs_pb2 import DnsD
from magma.common.misc_utils import get_ip_from_if_cidr
from magma.configuration.exceptions import LoadConfigError
from magma.configuration.mconfig_managers import load_service_mconfig
from magma.configuration.service_configs import load_service_config

CONFIG_OVERRIDE_DIR = '/var/opt/magma/tmp'


def _get_addresses(cfg, mconfig):
    """
    Return list of addresses mapping domain to a type of record.

    EG: [{'ip': '192.88.99.142', 'domain': 'baiomc.cloudapp.net'}]}
    Currently uses record types A record, AAAA record and CNAME record.
    """
    # Start with list of addresses from YML
    addresses = cfg.get('addresses', [])
    for record in mconfig.records:
        # Unpack each record type into list NOTE: list concat doesn't work here
        domain_records = [
            *record.a_record,
            *record.aaaa_record,
            #  *record.cname_record, TODO: Figure out how to repr CNAME
        ]
        for ip in domain_records:
            addresses.append({'domain': record.domain, 'ip': ip})

    return addresses


def get_context():
    """
    Provide context to pass to Jinja2 for templating.
    """
    context = {}
    cfg = load_service_config("dnsd")
    try:
        mconfig = load_service_mconfig('dnsd', DnsD())
    except LoadConfigError as err:
        logging.warning("Error! Using default config because: %s", err)
        mconfig = DnsD()
    ip = get_ip_from_if_cidr(cfg['enodeb_interface'])
    if int(ip.split('/')[1]) < 16:
        logging.fatal(
            "Large interface netmasks hang dnsmasq, consider using a "
            "netmask in range /16 - /24",
        )
        raise Exception(
            "Interface %s netmask is to large."
            % cfg['enodeb_interface'],
        )

    dhcp_block_size = cfg['dhcp_block_size']
    available_hosts = list(ipaddress.IPv4Interface(ip).network.hosts())

    context['dhcp_server_enabled'] = mconfig.dhcp_server_enabled
    if dhcp_block_size < len(available_hosts):
        context['dhcp_range'] = {
            "lower": str(available_hosts[-dhcp_block_size]),
            "upper": str(available_hosts[-1]),
        }
    else:
        logging.fatal(
            "Not enough available hosts to allocate a DHCP block of \
            %d addresses." % (dhcp_block_size),
        )

    context['addresses'] = _get_addresses(cfg, mconfig)
    return context


def main():
    logging.basicConfig(
        level=logging.INFO,
        format='[%(asctime)s %(levelname)s %(name)s] %(message)s',
    )

    generate_template_config(
        'dnsd', 'dnsd',
        CONFIG_OVERRIDE_DIR, get_context(),
    )


if __name__ == '__main__':
    main()
