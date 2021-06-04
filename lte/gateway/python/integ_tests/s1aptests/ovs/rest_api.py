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
import logging
import time

import requests
from integ_tests.s1aptests.ovs import DEV_VM_URL, MAX_RETRIES, OF_REST_PORT


def get_flows(datapath, fields, ip=DEV_VM_URL):
    """
    Query OVS REST API for flows for a datapath given some filter criteria
    in `fields`.
    """
    url = "http://%s:%d/stats/flow/%s" % (ip, OF_REST_PORT, datapath)
    data = json.dumps(fields)
    return _ovs_api_request('POST', url, data=data)[datapath]


def add_flowentry(fields, ip=DEV_VM_URL):
    """
    Add a flowentry to OVS, flowentry info such as datapath, table_id, match,
    actions, etc., is stored in `fields`
    """
    url = "http://%s:%d/stats/flowentry/add" % (ip, OF_REST_PORT)
    data = json.dumps(fields)
    return _ovs_api_request('POST', url, data=data, return_json=False)


def delete_flowentry(fields, ip=DEV_VM_URL):
    """
    Deletes a flowentry by using OVS REST API, flowentry is matched based on
    the information in  `fields`
    """
    url = "http://%s:%d/stats/flowentry/delete_strict" % (ip, OF_REST_PORT)
    data = json.dumps(fields)
    return _ovs_api_request('POST', url, data=data, return_json=False)


def get_datapath(ip=DEV_VM_URL):
    """
    Get the first datapath object from the OVS REST API.

    For integ tests there should only ever be one datapath.
    """
    url = "http://%s:%d/stats/switches" % (ip, OF_REST_PORT)
    return str(_ovs_api_request('GET', url)[0])


def _ovs_api_request(
    method, url, data=None, max_retries=MAX_RETRIES,
    return_json=True,
):
    """
    Send generic OVS REST API request and retry if request fails. Returns json
    decoded message
    """
    for _ in range(MAX_RETRIES):
        response = requests.request(method, url, data=data)
        if response.status_code == 200:
            if return_json:
                return response.json()
            else:
                return response.status_code
        time.sleep(1)
    logging.error(
        "Could not send %s request to OVS REST API at %s with data %s",
        method, url, data,
    )
    return {}
