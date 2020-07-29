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

# This script allows supervisord processes to log to stdout and stderr, by
# prefixing their process names.
# See http://supervisord.org/events.html for more info.

import sys


def write_stdout(s):
    # only eventlistener protocol messages may be sent to stdout
    sys.stdout.write(s)
    sys.stdout.flush()


def write_stderr(s):
    sys.stderr.write(s)
    sys.stderr.flush()


def main():
    while 1:
        # transition from ACKNOWLEDGED to READY
        write_stdout('READY\n')

        # read header line
        line = sys.stdin.readline()

        # read event payload
        headers = dict(x.split(':') for x in line.split())
        data = sys.stdin.read(int(headers['len']))

        # transition from READY to ACKNOWLEDGED
        write_stdout('RESULT %s\n%s' % (len(data), data))


def result_handler(event, response):
    # Parse the headers
    line, data = response.split('\n', 1)
    headers = dict(x.split(':') for x in line.split())

    # Get the log lines and prefix the process name and stdout/stderr
    lines = data.rstrip().split('\n')
    prefix = '%s %s | ' % (headers['processname'], headers['channel'])
    print('\n'.join([prefix + l for l in lines]))


if __name__ == '__main__':
    main()
