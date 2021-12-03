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

from firebase_admin import credentials, db, initialize_app


def publish_report(worker_id, build_id, verdict, report):
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


def url_to_html_redirect(run_id, url):
    """Convert URL into a redirecting HTML page"""
    report_url = url
    if url is None:
        report_url = 'https://github.com/magma/magma/actions/runs/' + run_id

    return '<script>'\
           '  window.location.href = "' + report_url + '";'\
           '</script>'


def lte_integ_test(args):
    """Prepare and publish LTE Integ Test report"""
    report = url_to_html_redirect(args.run_id, args.url)
    # Possible args.verdict values are success, failure, or canceled
    verdict = 'inconclusive'
    if args.verdict.lower() == 'success':
        verdict = 'pass'
    elif args.verdict.lower() == 'failure':
        verdict = 'fail'
    publish_report('lte_integ_test', args.build_id, verdict, report)


def cwf_integ_test(args):
    """Prepare and publish CWF Integ Test report"""
    report = url_to_html_redirect(args.run_id, args.url)
    # Possible args.verdict values are success, failure, or canceled
    verdict = 'inconclusive'
    if args.verdict.lower() == 'success':
        verdict = 'pass'
    elif args.verdict.lower() == 'failure':
        verdict = 'fail'
    publish_report('cwf_integ_test', args.build_id, verdict, report)


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

# Create the parser for the "lte" command
parser_lte = subparsers.add_parser('lte')
parser_lte.add_argument("--url", default="none", help="Report URL", nargs='?')
parser_lte.set_defaults(func=lte_integ_test)

# Create the parser for the "cwf" command
parser_cwf = subparsers.add_parser('cwf')
parser_cwf.add_argument("--url", default="none", help="Report URL", nargs='?')
parser_cwf.set_defaults(func=cwf_integ_test)

# Read arguments from the command line
args = parser.parse_args()
if not args.cmd:
    parser.print_usage()
    sys.exit(1)
args.func(args)
