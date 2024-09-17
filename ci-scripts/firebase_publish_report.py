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
import json
import os
import sys
import time
from typing import Optional

from firebase_admin import credentials, db, initialize_app


def publish_report(
    worker_id: str,
    build_id: str,
    verdict: str,
    report: str,
):
    """Publish report to Firebase realtime database"""
    # Read Firebase service account config from envirenment
    firebase_config = os.environ["FIREBASE_SERVICE_CONFIG"]
    config = json.loads(firebase_config)

    cred = credentials.Certificate(config)
    initialize_app(
        cred, {
            'databaseURL': 'https://magma-ci-default-rtdb.firebaseio.com/',
        },
    )

    report_dict = {
        'report': report,
        'timestamp': int(time.time()),
        'verdict': verdict,
    }

    ref = db.reference('/workers')
    reports_ref = ref.child(worker_id).child('reports').child(build_id)
    reports_ref.set(report_dict)


def url_to_html_redirect(run_id: str, url: Optional[str]):
    """Convert URL into a redirecting HTML page"""
    report_url = url
    if not url:
        report_url = f'https://github.com/magma/magma/actions/runs/{run_id}'

    return (
        f'<script>'
        f'  window.location.href = "{report_url}";'
        f'</script>'
    )


def debian_lte_integ_test(args):
    """Prepare and publish LTE Integ Test report"""
    prepare_and_publish('debian_lte_integ_test', args, 'test_status.txt')


def feg_integ_test(args):
    """Prepare and publish FEG Integ Test report"""
    prepare_and_publish('feg_integ_test', args, 'test_status.txt')


def cwf_integ_test(args):
    """Prepare and publish CWF Integ Test report"""
    prepare_and_publish('cwf_integ_test', args)


def sudo_python_tests(args):
    """Prepare and publish Sudo Python Test report"""
    prepare_and_publish('sudo_python_tests', args)


def containerized_lte_integ_test(args):
    """Prepare and publish containerized LTE Integ Test report"""
    prepare_and_publish('containerized_lte_integ_test', args, 'test_status.txt')


def prepare_and_publish(test_type: str, args, path: Optional[str] = None):
    """Prepare and publish test report"""
    report = url_to_html_redirect(args.run_id, args.url)
    # Possible args.verdict values are success, failure, or inconclusive
    verdict = 'inconclusive'

    if path and os.path.exists(path):
        # As per the recent change, CI process runs all integ tests ignoring
        # the failing test cases, because of which CI report always shows lte
        # integ test as success. Here we read the CI status from file for more
        # accurate lte integ test execution status
        with open(path, 'r') as file:
            status_file_content = file.read().rstrip()
            expected_verdict_list = ["pass", "fail"]
            if status_file_content in expected_verdict_list:
                verdict = status_file_content
    else:
        if args.verdict.lower() == 'success':
            verdict = 'pass'
        elif args.verdict.lower() == 'failure':
            verdict = 'fail'
    publish_report(test_type, args.build_id, verdict, report)


# Create the top-level parser
parser = argparse.ArgumentParser(
    description='Traffic CLI that generates traffic to an endpoint',
    formatter_class=argparse.ArgumentDefaultsHelpFormatter,
)


# Add arguments
parser.add_argument("--build_id", "-id", required=True, help="build ID")
parser.add_argument("--verdict", required=True, help="Test verdict")
parser.add_argument("--run_id", default="none", help="Github Actions Run ID")

# Add subcommands
subparsers = parser.add_subparsers(title='subcommands', dest='cmd')

tests = {
    'feg': feg_integ_test,
    'cwf': cwf_integ_test,
    'sudo_python_tests': sudo_python_tests,
    'debian_lte_integ_test': debian_lte_integ_test,
    'containerized_lte': containerized_lte_integ_test,
}

for key, value in tests.items():
    test_parser = subparsers.add_parser(key)
    test_parser.add_argument(
        "--url", default="none", help="Report URL", nargs='?',
    )
    test_parser.set_defaults(func=value)

# Read arguments from the command line
args = parser.parse_args()
if not args.cmd:
    parser.print_usage()
    sys.exit(1)
args.func(args)
