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
import shlex
import subprocess
import sys


class TrafficError(Exception):
    pass


def interface_exists(interface: str) -> None:
    cmd = 'nmcli device wifi list'
    args = shlex.split(cmd)
    nmcli = subprocess.Popen(args, stdout=subprocess.PIPE)
    cmd = 'grep {}'.format(interface)
    args = shlex.split(cmd)
    proc = subprocess.check_output(args, stdin=nmcli.stdout)
    print(proc.decode('utf-8'))


def connect_to_wifi(interface: str, password: str) -> None:
    cmd = 'nmcli device wifi connect {} password {}'.format(interface, password)
    args = shlex.split(cmd)
    nmcli = subprocess.Popen(args, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    errorstr = 'Error:'
    stdout = nmcli.communicate()[0].decode('utf-8')
    print(stdout)
    if errorstr in stdout:
        raise TrafficError('Error could not connect to {}'.format(interface))


def send_traffic(endpt: str) -> int:
    # 30 seconds until timeout
    cmd = 'curl -m 30 {}'.format(endpt)
    args = shlex.split(cmd)
    proc = subprocess.Popen(args, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
    print(proc.communicate()[0])
    return proc.returncode


def gen_traffic_handler(args):
    iface = args.iface
    pw = args.password
    endpt = args.endpoint
    try:
        interface_exists(iface)
    except Exception:
        print('Interface {} not found'.format(iface))
        sys.exit(1)
    connect_to_wifi(iface, pw)
    if send_traffic(endpt) != 0:
        print('Error pinging {}'.format(endpt))
        sys.exit(1)


def main():
    parser = argparse.ArgumentParser(
        description='Traffic CLI that generates traffic to an endpoint',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )

    # Add subcommands
    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')

    # gen_traffic
    subparser = subparsers.add_parser('gen_traffic', help='Generate traffic')
    subparser.add_argument('iface')
    subparser.add_argument('password')
    subparser.add_argument('endpoint')
    subparser.set_defaults(func=gen_traffic_handler)

    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        sys.exit(1)
    args.func(args)


if __name__ == '__main__':
    main()
