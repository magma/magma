#!/usr/bin/env python3

"""
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import argparse
import logging
import os
import re
import xml.etree.ElementTree as ET
from datetime import datetime
from pathlib import Path


def merge_all_report(working_dir, list_xml_report_paths, output_path):
    """ 
    Given the paths of all xml report files, merge them all into a single report files
    Args:
        working_dir          : path of folder contains all the report files
        list_xml_report_paths: list of all the report xml files
        output_path          : path of output report file in xml format 
    """
    testsuites_node = ET.Element("testsuites")

    testsuites_node.set("name", "AllTests")
    num_all_failures = 0
    num_all_disabled = 0
    num_all_errors = 0
    total_time = 0
    num_all_tests = 0
    init_time = datetime.max
    c = 0

    # iterate through files
    for xml_file_path in list_xml_report_paths:
        xml_file_path = working_dir + "/" + xml_file_path
        test_result = ET.parse(xml_file_path)
        test_suites_data = test_result.getroot()

        num_all_failures += int(test_suites_data.attrib['failures'])
        num_all_tests += int(test_suites_data.attrib['tests'])
        total_time += float(test_suites_data.attrib['time'])
        num_all_errors += int(test_suites_data.attrib['errors'])
        num_all_disabled += int(test_suites_data.attrib['disabled'])
        init_time = min(init_time, datetime.fromisoformat(test_suites_data.attrib['timestamp']))

        for single_test_suite_data in test_suites_data:
            testsuites_node.append(single_test_suite_data)

    # adding node here
    testsuites_node.set("timestamp", str(init_time))
    testsuites_node.set("time", str(total_time))
    testsuites_node.set("errors", str(num_all_errors))
    testsuites_node.set("disabled", str(num_all_disabled))
    testsuites_node.set("failures", str(num_all_failures))
    testsuites_node.set("tests", str(num_all_tests))

    tree = ET.ElementTree()
    tree._setroot(testsuites_node)
    tree.write(output_path)


if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)

    magma_root_path = os.environ.get('MAGMA_ROOT')
    default_output_path = ((magma_root_path + "/report/merged_report/") if magma_root_path else "/tmp/magma_oai/report/") +\
        f"report_all_tests.xml"

    parser = argparse.ArgumentParser(description="Merging all the .xml unittest reports into a single file")
    parser.add_argument('-w', '--working_dir', default=magma_root_path, help='Path of folder contains all report files', required=False)
    parser.add_argument('-i', '--input', help='<Required> regex for relative paths of xml files', required=True)
    parser.add_argument('-o', '--output', help='Output of xml report', default=default_output_path, required=False)

    args = parser.parse_args()
    path_working_dir = Path(args.working_dir)
    paths_all_report_files = filter(
        re.compile(args.input).match,
        [
            str(file_or_dir.relative_to(path_working_dir))
            for file_or_dir in path_working_dir.rglob("*")
            if file_or_dir.is_file()
        ],
    )
    Path(args.output).parent.mkdir(parents=True, exist_ok=True)
    merge_all_report(args.working_dir, paths_all_report_files, args.output)
    logging.info(f"Final report is generated at: {Path(args.output).resolve()}")
