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

import xml.etree.ElementTree as ET
from pprint import pprint
from pathlib import Path
import argparse
import os
import glob

def parsing_and_display(report_list):
    header = ['Testsuites', 'Testsuite', 'Testcase', 'Failed', 'Runtime']
    reformatted_report = [header,]
    for report_file_name in report_list:
        test_result = ET.parse(report_file_name)
        test_suites_data = test_result.getroot()
        testuites = [Path(report_file_name).stem, '', '', test_suites_data.attrib['failures'] + "/" + test_suites_data.attrib['tests'], test_suites_data.attrib['time']]
        reformatted_report.append(testuites)

        for single_test_suite_data in test_suites_data:
            test_suite_name = single_test_suite_data.attrib['name']
            number_failed_case = single_test_suite_data.attrib['failures'] + "/" + single_test_suite_data.attrib['tests']
            reformatted_report.append(['', test_suite_name, '', number_failed_case, single_test_suite_data.attrib['time']])
            for test_case in single_test_suite_data:
                failed = '0/1' if test_case.attrib['result'] == 'completed' else '1/1'
                reformatted_report.append(['', '', test_case.attrib['name'], failed, test_case.attrib['time']]) 

    max_len_per_col = [max(len(row[i]) for row in reformatted_report) for i in range(len(reformatted_report[0]))]

    format_row = ''.join(["{:>" + str(x + 3)+ "}"  for x in max_len_per_col])
    print(' ')
    print("RUNTIME REPORT:")
    for row in reformatted_report:
        print(format_row.format(*row))
    print(" ")

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("xml_report")
    args = parser.parse_args()

    print(f"Value of env variable GTEST_OUTPUT = {os.environ['GTEST_OUTPUT']}")
    print(f"Processing test reports at : {args.xml_report}")

    print("=" * 50)
    print("NOTE: THIS REPORT IS ONLY FOR FINISHED TESTS. FOR THE ABORTED TESTS, WE DO NOT HAVE DATA TO REPORT!")
    files = glob.glob(args.xml_report + "/*.xml")
    if len(files) < 1:
        print("No report is generated")
    else:
        parsing_and_display(files)
    # parsing_and_display(args.xml_report)