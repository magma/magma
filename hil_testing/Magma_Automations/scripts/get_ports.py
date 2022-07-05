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
import json
import sys
import warnings

import requests

###### Following Code finds the Next Open Port for use ######


def get(url, auth):
    ls = requests.Session()
    get_response = ls.get(url, auth=auth)
    if get_response.ok:
        return json.loads(get_response.content.decode("utf-8"))
    else:
        warnings.warn(
            "Get Method failed for url: "
            + url
            + ". Response:\n"
            + get_response.content.decode("utf-8"),
        )
        print(f"ERROR-get_port-func-get {url} res {get_response}")
        return json.loads("{}")


def check_sut(sut_name, data, auth):
    url = "http://" + data["tas_ip"] + ":8080/api/runningTests"
    r = get(url, auth)
    for n in r["runningTests"]:
        if (
            sut_name == n["name"].split("_")[-1]
            and "COMPLETE" not in n["testStateOrStep"]
        ):
            return n["id"]
        else:
            continue


def get_next_avail_port(data, auth):

    url = "http://" + data["tas_ip"] + ":8080/api/"

    r = get(url + "testServers/2", auth)

    port_ip_assignment = {
        "eth2": {
            "nodal": "<ipv4_network/cidr>",
            "nodal_v6": "<ipv6_network/cidr>",
            "network_host_v6": "<ipv6_network/cidr>",
            "network_host": "<ipv4_network/cidr>",
            "nodal_mac": "<starting_mac_address>",
            "nodal3_mac": "<starting_mac_address_for_another_nodal>",
            "nodal_user_mac": "<eNB_userplane_mac_address>",
            "nodal_mob_mac": "<mobility_eNB_mac_address>",
            "network_host_mac": "<network_host_mac_address>",
            "network_host2_mac": "<2nd_network_host_mac_address>",
            "network_host3_mac": "<3rd_network_host_mac_address>",
            "nodal_vlan": 298,
            "network_host_vlan": 101,
        },
        "eth3": {
            "nodal": "<ipv4_network/cidr>",
            "nodal_v6": "<ipv6_network/cidr>",
            "network_host_v6": "<ipv6_network/cidr>",
            "network_host": "<ipv4_network/cidr>",
            "nodal_mac": "<starting_mac_address>",
            "nodal3_mac": "<starting_mac_address_for_another_nodal>",
            "nodal_user_mac": "<eNB_userplane_mac_address>",
            "nodal_mob_mac": "<mobility_eNB_mac_address>",
            "network_host_mac": "<network_host_mac_address>",
            "network_host2_mac": "<2nd_network_host_mac_address>",
            "network_host3_mac": "<3rd_network_host_mac_address>",
            "nodal_vlan": 299,
            "network_host_vlan": 102,
        },
        "eth4": {
            "nodal": "<ipv4_network/cidr>",
            "nodal_v6": "<ipv6_network/cidr>",
            "network_host_v6": "<ipv6_network/cidr>",
            "network_host": "<ipv4_network/cidr>",
            "nodal_mac": "<starting_mac_address>",
            "nodal3_mac": "<starting_mac_address_for_another_nodal>",
            "nodal_user_mac": "<eNB_userplane_mac_address>",
            "nodal_mob_mac": "<mobility_eNB_mac_address>",
            "network_host_mac": "<network_host_mac_address>",
            "network_host2_mac": "<2nd_network_host_mac_address>",
            "network_host3_mac": "<3rd_network_host_mac_address>",
            "nodal_vlan": 300,
            "network_host_vlan": 103,
        },
        "eth5": {
            "nodal": "<ipv4_network/cidr>",
            "nodal_v6": "<ipv6_network/cidr>",
            "network_host_v6": "<ipv6_network/cidr>",
            "network_host": "<ipv4_network/cidr>",
            "nodal_mac": "<starting_mac_address>",
            "nodal3_mac": "<starting_mac_address_for_another_nodal>",
            "nodal_user_mac": "<eNB_userplane_mac_address>",
            "nodal_mob_mac": "<mobility_eNB_mac_address>",
            "network_host_mac": "<network_host_mac_address>",
            "network_host2_mac": "<2nd_network_host_mac_address>",
            "network_host3_mac": "<3rd_network_host_mac_address>",
            "nodal_vlan": 301,
            "network_host_vlan": 104,
        },
        "eth6": {
            "nodal": "<ipv4_network/cidr>",
            "nodal_v6": "<ipv6_network/cidr>",
            "network_host_v6": "<ipv6_network/cidr>",
            "network_host": "<ipv4_network/cidr>",
            "nodal_mac": "<starting_mac_address>",
            "nodal3_mac": "<starting_mac_address_for_another_nodal>",
            "nodal_user_mac": "<eNB_userplane_mac_address>",
            "nodal_mob_mac": "<mobility_eNB_mac_address>",
            "network_host_mac": "<network_host_mac_address>",
            "network_host2_mac": "<2nd_network_host_mac_address>",
            "network_host3_mac": "<3rd_network_host_mac_address>",
            "nodal_vlan": 302,
            "network_host_vlan": 105,
        },
        "eth7": {
            "nodal": "<ipv4_network/cidr>",
            "nodal_v6": "<ipv6_network/cidr>",
            "network_host_v6": "<ipv6_network/cidr>",
            "network_host": "<ipv4_network/cidr>",
            "nodal_mac": "<starting_mac_address>",
            "nodal3_mac": "<starting_mac_address_for_another_nodal>",
            "nodal_user_mac": "<eNB_userplane_mac_address>",
            "nodal_mob_mac": "<mobility_eNB_mac_address>",
            "network_host_mac": "<network_host_mac_address>",
            "network_host2_mac": "<2nd_network_host_mac_address>",
            "network_host3_mac": "<3rd_network_host_mac_address>",
            "nodal_vlan": 303,
            "network_host_vlan": 106,
        },
        "eth10": {
            "nodal": "<ipv4_network/cidr>",
            "nodal_v6": "<ipv6_network/cidr>",
            "network_host_v6": "<ipv6_network/cidr>",
            "network_host": "<ipv4_network/cidr>",
            "nodal_mac": "<starting_mac_address>",
            "nodal3_mac": "<starting_mac_address_for_another_nodal>",
            "nodal_user_mac": "<eNB_userplane_mac_address>",
            "nodal_mob_mac": "<mobility_eNB_mac_address>",
            "network_host_mac": "<network_host_mac_address>",
            "network_host2_mac": "<2nd_network_host_mac_address>",
            "network_host3_mac": "<3rd_network_host_mac_address>",
            "nodal_vlan": 304,
            "network_host_vlan": 107,
        },
        "eth11": {
            "nodal": "<ipv4_network/cidr>",
            "nodal_v6": "<ipv6_network/cidr>",
            "network_host_v6": "<ipv6_network/cidr>",
            "network_host": "<ipv4_network/cidr>",
            "nodal_mac": "<starting_mac_address>",
            "nodal3_mac": "<starting_mac_address_for_another_nodal>",
            "nodal_user_mac": "<eNB_userplane_mac_address>",
            "nodal_mob_mac": "<mobility_eNB_mac_address>",
            "network_host_mac": "<network_host_mac_address>",
            "network_host2_mac": "<2nd_network_host_mac_address>",
            "network_host3_mac": "<3rd_network_host_mac_address>",
            "nodal_vlan": 305,
            "network_host_vlan": 108,
        },
        "eth12": {
            "nodal": "<ipv4_network/cidr>",
            "nodal_v6": "<ipv6_network/cidr>",
            "network_host_v6": "<ipv6_network/cidr>",
            "network_host": "<ipv4_network/cidr>",
            "nodal_mac": "<starting_mac_address>",
            "nodal3_mac": "<starting_mac_address_for_another_nodal>",
            "nodal_user_mac": "<eNB_userplane_mac_address>",
            "nodal_mob_mac": "<mobility_eNB_mac_address>",
            "network_host_mac": "<network_host_mac_address>",
            "network_host2_mac": "<2nd_network_host_mac_address>",
            "network_host3_mac": "<3rd_network_host_mac_address>",
            "nodal_vlan": 306,
            "network_host_vlan": 109,
        },
        "eth13": {
            "nodal": "10.22.101.192/26",
            "nodal_v6": "3000:11e:fff0:710::/64",
            "network_host_v6": "3000:11e:fff0:720::/64",
            "network_host": "10.22.24.16/28",
            "nodal_mac": "F4:15:63:A0:00:00",
            "nodal3_mac": "F4:15:63:A0:00:30",
            "nodal_mob_mac": "F4:15:63:A0:00:20",
            "nodal_user_mac": "F4:15:63:20:00:00",
            "network_host_mac": "F4:15:63:A5:00:00",
            "network_host2_mac": "F4:16:63:A5:00:00",
            "network_host3_mac": "F4:17:63:A5:00:00",
            "nodal_vlan": 307,
            "network_host_vlan": 110,
        },
    }

    max_sessions = 8
    current_sessions = 0
    acceptable_ports = [
        "eth4",
        "eth5",
        "eth6",
        "eth7",
        "eth10",
        "eth11",
        "eth12",
        "eth13",
    ]  # These are the 1G Ports we are currently reserving.
    sorry = "Sorry, no Spirent sessions available. Please try again later."
    try:
        for ports in r["ethInfos"]:
            if (
                not len(ports["runIds"]) > 0 and ports["name"] in acceptable_ports
            ):  # i.e. there are no run IDs associated with this port
                return {
                    ports["name"]: port_ip_assignment[ports["name"]],
                }  # Found the next avail port; we are done!

            if len(ports["runIds"]) > 0 and "v6" not in ports["name"]:
                current_sessions += 1

            if current_sessions == max_sessions:
                sys.exit(sorry)
    except KeyError:
        sys.exit("KeyError")
