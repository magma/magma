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
import os
import re
import sys

from generate_html import (
    add_compilation_details,
    add_compile_rows,
    add_target_image_gen_row,
    add_target_image_size_row,
    generate_build_footer,
    generate_build_header,
    generate_footer,
    generate_git_summary,
    generate_header,
)

REPORT_NAME = 'build_results_magma_oai_mme.html'


def main() -> None:
    """Provide command-line options to generate a HTML report"""
    args = _parse_args()

    if args.git_merge_request:
        if args.git_target_branch is None:
            sys.exit('git_target_branch: Missing Parameter')
        if args.git_target_commit is None:
            sys.exit('git_target_commit: Missing Parameter')

    if args.mode == 'Build':
        generate_build_report(args)

    if args.mode == 'TestWithDsTest':
        append_build_summary(args, 'test_results_magma_oai_epc.html')
        append_build_summary(args, 'test_results_magma_epc_u18.html')

    if args.mode == 'RHEL8SanityCheck':
        append_build_summary(args, 'test_results_magma_epc_rhel8.html')


def _parse_args() -> argparse.Namespace:
    """Parse the command line args

    Returns:
        argparse.Namespace: the created parser
    """
    parser = argparse.ArgumentParser(description='OAI HTML report generator')

    # Jenkins Job parameters
    parser.add_argument(
        '--job_name', '-jn',
        action='store',
        required=True,
        help='Jenkins Job Name',
    )
    parser.add_argument(
        '--job_id', '-id',
        action='store',
        required=True,
        help='Jenkins Job Build ID',
    )
    parser.add_argument(
        '--job_url', '-url',
        action='store',
        required=True,
        help='Jenkins Job Build URL',
    )

    # Git Parameters
    parser.add_argument(
        '--git_url',
        action='store',
        required=True,
        help='Git Repository URL',
    )
    parser.add_argument(
        '--git_src_branch',
        action='store',
        required=True,
        help='Git Source Branch',
    )
    parser.add_argument(
        '--git_src_commit',
        action='store',
        required=True,
        help='Git Source Commit (SHA-ONE)',
    )
    parser.add_argument(
        '--git_src_commit_msg',
        action='store',
        help='Git Source Commit Message',
    )

    # Pull Request Parameters
    parser.add_argument(
        '--git_merge_request',
        action='store_true',
        default=False,
        help='Git Pull Request Active',
    )
    parser.add_argument(
        '--git_target_branch',
        action='store',
        help='Git Target Branch',
    )
    parser.add_argument(
        '--git_target_commit',
        action='store',
        help='Git Target Commit (SHA-ONE)',
    )

    # Mode
    parser.add_argument(
        '--mode',
        action='store',
        required=True,
        choices=['Build', 'TestWithDsTest', 'RHEL8SanityCheck'],
        help='HTML Generation Mode',
    )

    return parser.parse_args()


def generate_build_report(args):
    """
    Create the BUILD HTML report.

    Args:
        args: results from argument parser
    """
    cwd = os.getcwd()
    with open(os.path.join(cwd, REPORT_NAME), 'w') as wfile:
        wfile.write(generate_header(args))
        wfile.write(generate_git_summary(args))
        wfile.write(generate_build_header())
        wfile.write(add_compile_rows())
        wfile.write(add_target_image_gen_row())
        wfile.write(add_target_image_size_row(args))
        wfile.write(generate_build_footer())
        wfile.write(add_compilation_details())
        wfile.write(generate_footer())


def append_build_summary(args, filename):
    """
    Append the GIT summary to test report.

    Args:
        args: results from argument parser
        filename: file to append to
    """
    cwd = os.getcwd()
    if not os.path.isfile(cwd + '/' + filename):
        return
    report = ''
    org_file = os.path.join(cwd, filename)
    with open(org_file, 'r') as org_f:
        report = org_f.read()

    build_summary_to_be_done = True
    with open(org_file, 'w') as org_f:
        for line in report.split('\n'):
            my_res = re.search('Deployment Summary', line)
            if (my_res is not None) and build_summary_to_be_done:
                summary = generate_git_summary(args)
                org_f.write(summary)
                build_summary_to_be_done = False
            org_f.write(line + '\n')


if __name__ == '__main__':
    main()
