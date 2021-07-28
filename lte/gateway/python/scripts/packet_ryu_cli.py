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
import logging

from integ_tests.s1aptests.ovs.rest_api import (
    add_flowentry,
    delete_flowentry,
    get_datapath,
    get_flows,
)
from scapy.all import IP, Ether, sendp

logging.getLogger("scapy.runtime").setLevel(logging.ERROR)

DEFAULT_PKT_MAC_SRC = "00:00:00:00:00:01"
DEFAULT_PKT_MAC_DST = "12:99:cc:97:47:4e"


def _simple_remove(args):
    del_old = {"dpid": datapath, "priority": args.priority}
    if args.cookie:
        del_old["cookie"] = args.cookie
    if args.table_id:
        del_old["table_id"] = args.table_id
    print(del_old)
    delete_flowentry(del_old)


def _simple_get(args):
    query = {"table_id": args.table_id}
    flows = get_flows(datapath, query)
    print("FlowEntry Match : captured packets")
    for flowentry in flows:
        print("Prior ", flowentry["priority"], end=' - ')
        print(flowentry["match"], flowentry["packet_count"], sep=', pkts: ')


def _simple_add(args):
        fields = {
            "dpid": datapath, "table_id": args.table_start,
            "priority": args.priority, "instructions":
            [{"type": "GOTO_TABLE", "table_id": args.table_end}],
        }
        if args.cookie:
            fields["cookie"] = args.cookie
        if args.reg1:
            reg1 = int(args.reg1, 0)
            fields["instructions"].append({
                "type": "APPLY_ACTIONS", "actions":
                [{
                    "type": "SET_FIELD", "field":
                    "reg1", "value": reg1,
                }],
            })
        add_flowentry(fields)


def _simple_send(args):
    eth = Ether(dst=DEFAULT_PKT_MAC_DST, src=DEFAULT_PKT_MAC_SRC)
    ip = IP(proto=1, src=args.ipv4_src, dst=args.ipv4_dst)
    pkt = eth / ip
    print(pkt.show())
    sendp(pkt, iface=args.iface, count=args.num)


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        description='CLI for testing packet movement through pipelined,\
                     using RYU REST API & Scapy',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )

    # Add subcommands
    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')
    parser_dump = subparsers.add_parser('dump', help='Dump packet stats')
    parser_dump.add_argument('table_id', help='table id to print', type=int)

    parser_send = subparsers.add_parser('send', help='Send packets')
    parser_send.add_argument('iface', help='iface to send to')
    parser_send.add_argument('-ipd', '--ipv4_dst', help='ipv4 dst for pkt')
    parser_send.add_argument('-ips', '--ipv4_src', help='ipv4 src for pkt')
    parser_send.add_argument(
        '-n', '--num', help='number of packets to send',
        default=5, type=int,
    )

    parser_skip = subparsers.add_parser('skip', help='Add flowentry')
    parser_skip.add_argument(
        'table_start', type=int,
        help='table to insert flowentry',
    )
    parser_skip.add_argument(
        'table_end', type=int,
        help='table to forward to',
    )
    parser_skip.add_argument(
        '-c', '--cookie', default=0, type=int,
        help='flowentry cookie value',
    )
    parser_skip.add_argument('-r1', '--reg1', help='flowentry reg1 value')
    parser_skip.add_argument(
        '-p', '--priority', help='flowentry priority',
        type=int, default=65535,
    )

    parser_rem = subparsers.add_parser('rem', help='Remove flowentry')
    parser_rem.add_argument(
        '-tid', '--table_id', type=int,
        help='table to remove flowentry from',
    )
    parser_rem.add_argument(
        '-p', '--priority', default=65535, type=int,
        help='rm flowentry matching priority value',
    )
    parser_rem.add_argument(
        '-c', '--cookie',
        help='rm flowentry matching cookie value',
    )

    # Add function callbacks
    parser_dump.set_defaults(func=_simple_get)
    parser_send.set_defaults(func=_simple_send)
    parser_skip.set_defaults(func=_simple_add)
    parser_rem.set_defaults(func=_simple_remove)

    return parser


def main():
    global datapath
    datapath = get_datapath()
    if not datapath:
        print("Coudn't get datapath")
        exit(1)

    parser = create_parser()
    # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        exit(1)

    # Execute the subcommand function
    args.func(args)


if __name__ == "__main__":
    main()
