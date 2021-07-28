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
import argparse
import socket
import subprocess

import jsonpickle
import requests
import urllib3

admin_cert = (
    "/var/opt/magma/certs/rest_admin.crt",
    "/var/opt/magma/certs/rest_admin.key",
)

# Disable warnings about SSL verification since its a local VM
urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)


def cloud_get(url: str):
    resp = requests.get(url, verify=False, cert=admin_cert)
    if resp.status_code != 200:
        return
    return resp.json()


def get_gateway_id(url: str, network_id: str, hw_id: str):
    """
    Retrieve the unique server identifier for this machine.
    """
    gateways = cloud_get(url + f"/cwf/{network_id}/gateways")
    if gateways is None:
        return str(1)

    if hw_id != "":
        for gw in gateways.values():
            if hw_id == gw['device']['hardware_id']:
                gw_id = gw['id'].split("_")[2]
                return gw_id
    return str(len(gateways) + 1)


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        description="Set unique hostname for XwF-M Gateway",
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )
    parser.add_argument(
        "--partner", dest="partner", action="store", help="Partner Short Name",
    )
    parser.add_argument(
        "--hwid", dest="hwid", action="store", help="Gateway Hardware ID",
    )
    parser.add_argument(
        "--url", dest="url", action="store", help="Orchestrator URL Address",
    )
    parser.add_argument(
        "--env", dest="env", action="store", help="Environment type: QA/Prod",
    )
    return parser


def main():
    parser = create_parser()
    args = parser.parse_args()
    if not (args.url and args.partner and args.env):
        parser.print_usage()
        exit(1)

    # Create XwF-M Network
    partner = args.partner.strip()
    server_id = get_gateway_id(args.url, partner, args.hwid)
    hostname = 'xwfm.' + args.partner + '.' + args.env + '.' + server_id
    with open('/etc/hostname', 'w') as config:
        config.write(hostname)
    subprocess.call(['hostname', hostname])


if __name__ == "__main__":
    main()
