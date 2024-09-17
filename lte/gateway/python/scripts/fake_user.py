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

import argparse
import os
import sys

import netifaces
from lte.protos.session_manager_pb2 import LocalCreateSessionRequest
from lte.protos.session_manager_pb2_grpc import LocalSessionManagerStub
from magma.common.service_registry import ServiceRegistry
from magma.configuration import environment
from magma.pipelined.imsi import encode_imsi
from magma.subscriberdb.sid import SIDUtils


def sample_commands():
    return (
        'Example commands to run as fake_user:\n'
        'ping -I fake_user 8.8.8.8\n'
        'sudo curl --interface fake_user -vvv --ipv4 http://www.google.com'
        ' > /dev/null'
    )


def output(str, color=None):
    codes_map = {'blue': '\033[94m', 'green': '\033[92m', 'red': '\033[91m'}
    code = codes_map.get(color, '\x1b[0m')
    print(code + str + '\033[0m')


def mac_from_iface(iface_name):
    """ Returns the mac address of the specified interface """
    ifaddresses = netifaces.ifaddresses(iface_name)
    return ifaddresses[netifaces.AF_LINK][0]['addr']


def run(command, ignore_errors=False):
    """ Runs the shell command """
    output(command, 'blue')
    ret = os.system(command)
    if ret != 0 and not ignore_errors:
        output('Error!! Command returned: %d' % ret, 'red')
        sys.exit(1)


def add_flow(table, filter, actions, priority=300):
    """ Adds/modifies an OVS flow.
        We use '0xface0ff' as the cookie for all flows created by this tool """
    run(
        'sudo ovs-ofctl add-flow gtp_br0 "cookie=0xface0ff, '
        'table=%d, priority=%d,%s actions=%s"' %
        (table, priority, filter, actions),
    )


def pcef_create_session(args):
    chan = ServiceRegistry.get_rpc_channel(args.pcef, ServiceRegistry.LOCAL)
    client = LocalSessionManagerStub(chan)
    request = LocalCreateSessionRequest(
        sid=SIDUtils.to_pb(args.imsi),
        ue_ipv4=args.user_ip,
    )
    client.CreateSession(request)


def pcef_end_session(args):
    chan = ServiceRegistry.get_rpc_channel(args.pcef, ServiceRegistry.LOCAL)
    client = LocalSessionManagerStub(chan)
    client.EndSession(SIDUtils.to_pb(args.imsi))


def create_fake_user(args):
    output('Creating fake_user...')
    # Create the OVS port and interface
    run(
        'sudo ovs-vsctl add-port gtp_br0 fake_user -- '
        'set interface fake_user type=internal', ignore_errors=True,
    )
    run('sudo ifconfig fake_user up %s netmask 255.255.255.0' % args.iface_ip)

    if environment.is_dev_mode():
        # Make sure we can NAT to the internet
        run('sudo iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE')

    # Create table 0 flows to set the metadata, and to replace with iface_ip
    metadata = encode_imsi(args.imsi)
    add_flow(
        0, 'ip,ip_src=%s' % args.iface_ip,
        'set_field:%s->metadata,resubmit(,1)' % (metadata),
    )
    iface_mac = mac_from_iface('fake_user')
    add_flow(
        0, 'ip,ip_dst=%s' % args.user_ip,
        'set_field:%s->metadata,set_field:%s->ip_dst,'
        'set_field:%s->eth_dst,resubmit(,1)' %
        (metadata, args.iface_ip, iface_mac),
    )

    # Create a table 1 flow to set 0x1 as direction bit for pkts from fake_user
    add_flow(1, 'in_port=fake_user', 'set_field:0x1->reg1,resubmit(,2)')

    # Create a table 2 flow to respond to ARP requests from the fake_user
    # This doesn't happen with GTP, since the UEs by default send it to the
    # tunnel. Here we manually set the next hop as gtp_br0.
    gateway_mac = mac_from_iface('gtp_br0')
    add_flow(
        2, 'arp,reg1=0x1',
        'move:NXM_OF_ETH_SRC[]->NXM_OF_ETH_DST[],mod_dl_src:%s,'
        'load:0x2->NXM_OF_ARP_OP[],'
        'load:0x3243f1395144->NXM_NX_ARP_SHA[],'
        'move:NXM_NX_ARP_SHA[]->NXM_NX_ARP_THA[],'
        'move:NXM_OF_ARP_TPA[]->NXM_NX_REG0[],'
        'move:NXM_OF_ARP_SPA[]->NXM_OF_ARP_TPA[],'
        'move:NXM_NX_REG0[]->NXM_OF_ARP_SPA[],IN_PORT' % gateway_mac,
    )

    # Create a table 20 flow to output the pkts to face_user after swapping
    # back the destination ip.
    add_flow(
        20, 'ip,ip_src=%s' % args.iface_ip,
        'set_field:%s->ip_src,output:LOCAL' % (args.user_ip),
    )
    add_flow(20, 'ip,ip_dst=%s' % args.iface_ip, 'output:fake_user')

    # Create session with PCEF if needed
    if args.pcef:
        pcef_create_session(args)

    output('fake_user interface created successfully!!\n', 'green')
    output(sample_commands(), 'green')


def remove_fake_user(args):
    output('Removing fake_user...')
    run('sudo ovs-ofctl del-flows gtp_br0 "cookie=0xface0ff/0xfffffff"')
    run('sudo ovs-vsctl del-port fake_user', ignore_errors=True)

    # End session with PCEF if needed
    if args.pcef:
        pcef_end_session(args)

    output('All clear!!', 'green')


def main():
    parser = argparse.ArgumentParser(
        description='CLI for managing the fake_user interface:\n\n' +
                    sample_commands(),
        formatter_class=argparse.RawTextHelpFormatter,
    )
    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')

    # Create fake_user
    subparser = subparsers.add_parser('create', help='Creates the fake_user')
    subparser.add_argument(
        '--user_ip', help='IP address to use for the user',
        default='192.168.128.200',
    )
    subparser.add_argument(
        '--iface_ip', help='IP address of the interface',
        default='10.10.10.10',
    )
    subparser.add_argument('--imsi', help='IMSI the user', default='IMSI12345')
    subparser.add_argument('--pcef', help='PCEF to create session', default='')
    subparser.set_defaults(func=create_fake_user)

    # Remove fake_user
    subparser = subparsers.add_parser('remove', help='Removes the fake_user')
    subparser.add_argument('--imsi', help='IMSI of the user', default='IMSI12345')
    subparser.add_argument('--apn', help='APN of the session', default='oai.ipv4')
    subparser.add_argument('--pcef', help='PCEF to end session', default='')
    subparser.set_defaults(func=remove_fake_user)

    # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        sys.exit(1)

    # Execute the subcommand function
    args.func(args)


if __name__ == "__main__":
    main()
