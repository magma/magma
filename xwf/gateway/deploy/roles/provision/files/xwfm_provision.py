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
from typing import List

import jsonpickle
import requests
import urllib3

NETWORK_TYPE = "carrier_wifi_network"

admin_cert = (
    "/var/opt/magma/certs/rest_admin.crt",
    "/var/opt/magma/certs/rest_admin.key",
)

# Disable warnings about SSL verification since its a local VM
urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)


class AllowedGREPeers:
    def __init__(self, ip: str, key: int):
        self.ip = ip
        self.key = key


class CarrierWiFiConfig:
    def __init__(self, grePeers: List[AllowedGREPeers]):
        self.allowed_gre_peers = grePeers


class NetworkDNSConfig:
    def __init__(self, enable_caching: bool, local_ttl: int):
        self.enable_caching = enable_caching
        self.local_ttl = local_ttl


class XwFMNetwork:
    def __init__(self, id: str, name: str, description: str):
        self.id = id
        self.name = name
        self.description = description
        self.type = NETWORK_TYPE
        self.dns = NetworkDNSConfig(enable_caching=False, local_ttl=60)


class TierImage:
    def __init__(self, name: str, order: int):
        self.name = name
        self.order = order


class Tier:
    def __init__(
        self,
        id: str,
        name: str,
        version: str,
        images: List[TierImage],
        gateways: List[str],
    ):
        self.id = id
        self.name = name
        self.images = images
        self.version = version
        self.gateways = gateways


class MagmadGatewayConfigs:
    def __init__(
        self,
        autoupgrade_enabled: bool,
        autoupgrade_poll_interval: int,
        checkin_interval: int,
        checkin_timeout: int,
    ):
        self.autoupgrade_enabled = autoupgrade_enabled
        self.autoupgrade_poll_interval = autoupgrade_poll_interval
        self.checkin_interval = checkin_interval
        self.checkin_timeout = checkin_timeout


class ChallengeKey:
    def __init__(self, key_type: str):
        self.key_type = key_type


class GatewayDevice:
    def __init__(self, hardware_id: str, key: ChallengeKey):
        self.hardware_id = hardware_id
        self.key = key


class Gateway:
    def __init__(
        self,
        id: str,
        name: str,
        description: str,
        magmad: MagmadGatewayConfigs,
        device: GatewayDevice,
        tier: str,
        carrier_wifi: CarrierWiFiConfig,
    ):
        self.id = id
        self.name, self.description = name, description
        self.magmad = magmad
        self.device = device
        self.tier = tier
        self.carrier_wifi = carrier_wifi


def cloud_get(url: str):
    resp = requests.get(url, verify=False, cert=admin_cert)
    if resp.status_code != 200:
        raise Exception("Received a %d response: %s" % (resp.status_code, resp.text))
        return
    return resp.json()


def cloud_post(url: str, data: str):
    resp = requests.post(
        url,
        data=data,
        headers={"content-type": "application/json"},
        verify=False,
        cert=admin_cert,
    )

    if resp.status_code not in [200, 201, 204]:
        raise Exception("Received a %d response: %s" % (resp.status_code, resp.text))


def create_network_if_not_exists(url: str, network_id: str):
    values = cloud_get(url + "/networks")
    if network_id in values:
        print(f"NMS XWF-M Network exists already - {network_id}")
    else:
        data = XwFMNetwork(
            id=network_id, name="XWFM Network", description="XWFM Network",
        )
        cloud_post(url + "/networks", jsonpickle.pickler.encode(data))

        # create tier
        tier_payload = Tier(
            id="default", name="default", version="0.0.0-0", images=[], gateways=[],
        )
        cloud_post(
            url + f"/networks/{network_id}/tiers",
            jsonpickle.pickler.encode(tier_payload),
        )
        print(f"{network_id} NMS XWF-M Network created successfully")


def get_next_gateway_id(url: str, network_id: str, hw_id: str) -> (bool, str):
    gateways = cloud_get(url + f"/cwf/{network_id}/gateways")
    for gw in gateways.values():
        if gw['device']['hardware_id'] == hw_id:
            return True, gw['id']
    nbr = len(gateways) + 1
    return False, str(nbr)


def register_gateway(url: str, network_id: str, hardware_id: str, tier_id: str):
    """
    Register XwF-M Gateway in the requested network.
    """
    found, gid = get_next_gateway_id(url, network_id, hardware_id)
    if found:
        print(f"XWF-M Gateway exists already - {hardware_id}")
    else:
        grePeer = AllowedGREPeers(ip="192.168.128.2", key=100)
        data = Gateway(
	    name=socket.gethostname().strip(),
	    description=f"XWFM Gateway {gid}",
	    tier="default",
            id=f"fbc_gw_{gid}",
            device=GatewayDevice(
                hardware_id=hardware_id, key=ChallengeKey(key_type="ECHO"),
            ),
	    magmad=MagmadGatewayConfigs(
	        autoupgrade_enabled=True,
	        autoupgrade_poll_interval=60,
	        checkin_interval=60,
	        checkin_timeout=30,
	    ),
	    carrier_wifi=CarrierWiFiConfig(grePeers=[grePeer]),
        )
        cloud_post(url + f"/cwf/{network_id}/gateways", jsonpickle.pickler.encode(data))


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        description="Provision XwF-M Gateway",
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
    return parser


def main():
    parser = create_parser()
    args = parser.parse_args()
    if not (args.hwid and args.url and args.partner):
        parser.print_usage()
        exit(1)

    # Create XwF-M Network
    partner = args.partner.strip()
    create_network_if_not_exists(args.url, partner)
    register_gateway(args.url, partner, args.hwid, "default")


if __name__ == "__main__":
    main()
