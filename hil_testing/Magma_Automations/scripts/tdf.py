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
import csv
import math
import subprocess

import requests


def create_nas_tdf(*args, **kwargs):
    local_subs = math.ceil(kwargs["total_subs"] * (1 - kwargs["pct_roaming"]))
    remote_subs = kwargs["total_subs"] - local_subs
    local_imsi_start = kwargs["local_imsi_start"]
    remote_imsi_start = kwargs["remote_imsi_start"]
    local_imei_seed = kwargs["local_imei_start"]
    remote_imei_seed = kwargs["remote_imei_start"]
    final_list = []
    fields = ["IMSI", "IMEI", "Mobile Network Code"]
    filename = "_".join(
        [
            "inbound_roaming",
            str(kwargs["total_subs"]),
            "subs",
            str(kwargs["pct_roaming"]),
            "pct_roaming.csv",
        ],
    )

    # create the user list
    for n in range(local_subs):
        imsi = str((int(local_imsi_start) + n)).zfill(15)
        imei = str((int(local_imei_seed) + n)).zfill(14)
        mnc = kwargs["local_mnc"]
        final_list.append([imsi, imei, mnc])

    for n in range(remote_subs):
        imsi = str((int(remote_imsi_start) + n)).zfill(15)
        imei = str((int(remote_imei_seed) + n)).zfill(14)
        mnc = kwargs["remote_mnc"]
        final_list.append([imsi, imei, mnc])
    # write to csv
    with open("/tmp/" + filename, "w") as csvfile:
        csvwriter = csv.writer(csvfile)
        csvwriter.writerow(fields)
        csvwriter.writerows(final_list)
    # read from csv
    payload = open("/tmp/" + filename, "rb").read()
    # url = "http://" + kwargs["tas_ip"] + ":8080/api/libraries/342/tdfs/" + filename
    url = "".join(
        [
            "http://",
            kwargs["tas_ip"],
            ":8080/api/libraries/",
            kwargs["library_id"],
            "/tdfs/",
            filename,
        ],
    )
    headers = {
        "Content-Type": "application/binary",
    }
    response = requests.request(
        "POST", url, headers=headers, data=payload, auth=kwargs["auth"],
    )
    if response.status_code == 201:
        # delete if saving it to TAS is good
        subprocess.call(["rm /tmp/" + filename], shell=True)

    return filename
