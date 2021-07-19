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
        raise Exception("Received a %d response: %s" % (resp.status_code, resp.text))
        return
    return resp.json()


def cloud_delete(url: str):
    resp = requests.delete(
        url,
        headers={"content-type": "application/json"},
        verify=False,
        cert=admin_cert,
    )

    if resp.status_code not in [200, 201, 204]:
        raise Exception("Received a %d response: %s" % (resp.status_code, resp.text))


def deregister_all_gateways(url: str, network_id: str):
    """
    Deprovision all XwF-M Gateways in the requested network.
    """
    gateways = cloud_get(url + f"/cwf/{network_id}/gateways")
    for gw in gateways.values():
        gateway_id = gw['id']
        cloud_delete(url + f"/cwf/{network_id}/gateways/{gateway_id}")
        print(f"XWF-M gateway {gateway_id} deprovisioned")


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        description="Deprovision XwF-M Gateway",
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )
    parser.add_argument(
        "--url", dest="url", action="store", help="Orchestrator URL Address",
    )
    parser.add_argument(
        "--partner", dest="partner", action="store", help="Partner Short Name",
    )
    return parser


def main():
    parser = create_parser()
    args = parser.parse_args()
    if not (args.url and args.partner):
        parser.print_usage()
        exit(1)

    partner = args.partner.strip()
    deregister_all_gateways(args.url, partner)


if __name__ == "__main__":
    main()
