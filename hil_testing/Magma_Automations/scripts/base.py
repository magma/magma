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
from ipaddress import ip_address

import requests
from get_ports import get
from netaddr import *


def input_validation(input_dict):
    for arg in input_dict:
        if type(input_dict[arg]) == str or input_dict[arg] is None:
            continue
        else:
            if input_dict[arg] < 0:
                raise ValueError("{arg_name} cannot be negative!".format(arg_name=arg))


def mac_offset(addr, offset):
    mac = EUI(addr)
    new_mac = EUI(int(mac) + offset)
    new_mac.dialect = mac_unix_expanded
    return str(new_mac)


def save_test(url, payload, auth):
    response = requests.request("POST", url, data=json.dumps(payload), auth=auth)
    return response


def run_test(payload, tas_ip, auth):
    url = "http://" + tas_ip + ":8080/api/runningTests"
    response = requests.request("POST", url, data=json.dumps(payload), auth=auth)
    return json.loads(response.text)


def credentials(filename):
    # Read-in creds to access the spirent TS
    with open(filename) as f:
        data = json.load(f)
    return data


def get_library_id(name, tas_ip, auth):
    # This function finds the appropriate library_id for the said library name!
    url = "http://" + tas_ip + ":8080/api/libraries"
    output = get(url, auth)
    sorry = "no valid library ID found"
    try:
        for n in output["libraries"]:
            if name in n["name"]:
                library_id = n["id"]
        return str(library_id)
    except:
        return sorry


def imsi_calc(imsi_start, invalid_sub):
    output = "00" + str(int(imsi_start) - invalid_sub)
    return output
