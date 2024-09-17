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
import logging
import sys
import warnings

import requests
from base import credentials
from delete_prompush import delete_prompush
from get_ports import get
from requests.auth import HTTPBasicAuth


def delete_item(library_id, test_name, tas_ip, auth, item):
    logger = logging.getLogger("delete_test")
    url = (
        "http://" + tas_ip + ":8080/api/libraries/" + str(library_id) + item + test_name
    )
    r = requests.request("DELETE", url, auth=auth, data={})
    logger.info(r.text)


def list_of_test_cases(tas_ip, library_id, auth):
    logger = logging.getLogger("list_of_test_cases")
    tests_to_be_deleted = []
    url = (
        "http://" + tas_ip + ":8080/api/libraries/" + str(library_id) + "/testSessions/"
    )
    r = get(url, auth)
    for sessions in r["testSessions"]:
        try:
            if sessions["keywords"] == "DELETE_ME ":
                tests_to_be_deleted.append(sessions["name"])
        except KeyError:
            logger.warning("no keywords")
    return tests_to_be_deleted


def list_of_tdfs(tas_ip, library_id, auth):
    logger = logging.getLogger("list_of_tdfs")
    tdfs_to_be_deleted = []
    url = "http://" + tas_ip + ":8080/api/libraries/" + str(library_id) + "/tdfs/"
    r = get(url, auth)
    for tdfs in r["tdfs"]:
        try:
            tdfs_to_be_deleted.append(tdfs["name"])
        except KeyError:
            logger.warning("no tdfs")
    return tdfs_to_be_deleted


def get_all_libraries(tas_ip, auth):
    # This function finds the appropriate library_id for the said library name!
    url = "http://" + tas_ip + ":8080/api/libraries/"
    libraries = []
    r = get(url, auth)
    for n in r["libraries"]:
        if int(n["id"]) > 0:
            libraries.append(n["id"])
    return libraries


def delete_artifacts(tas_ip, auth):
    logger = logging.getLogger("delete_artifacts")
    test_cases = {}
    tdfs = {}
    library_list = get_all_libraries(tas_ip, auth)
    for n in library_list:
        test_cases[n] = list_of_test_cases(tas_ip, n, auth)
        tdfs[n] = list_of_tdfs(tas_ip, n, auth)

    for n in test_cases:
        if (
            len(test_cases[n]) > 0
        ):  # delete this case only if we identified anything to be deleted.
            for test_name in test_cases[n]:
                delete_item(str(n), test_name, tas_ip, auth, "/testSessions/")
                logging.info("Deleting TC {name}".format(name=test_name))
        else:
            continue
    for n in tdfs:
        if len(tdfs[n]) > 0:  # go through deleting tdfs if something in the library id
            for tdf in tdfs[n]:
                delete_item(str(n), tdf, tas_ip, auth, "/tdfs/")
                logging.info("Deleting TDF {name}".format(name=tdf))
        else:
            continue

    logger.info("All deleted your royal highness; nothing more to delete!")


def main():
    logging.basicConfig(level=logging.WARNING)
    data = credentials("creds.json")
    auth = HTTPBasicAuth(data["username"], data["password"])
    tas_ip = data["tas_ip"]
    delete_artifacts(tas_ip, auth)
    delete_prompush()


if __name__ == "__main__":
    main()
