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
import re
import subprocess

from lte.protos.mconfig import mconfigs_pb2
from magma.common.service import MagmaService


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        description="CLI wrapper around ovs-appctl ofproto/trace.\n"
        "To display the Datapath actions of the supplied IMSI",
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )
    parser.add_argument("-i", "--imsi", help="IMSI of the subscriber")
    parser.add_argument("-I", "--ip", help="External IP")
    parser.add_argument("-P", "--port", help="External Port")
    parser.add_argument("-UP", "--ue_port", help="UE Port")
    parser.add_argument("-p", "--protocol", help="Portocol (i.e. tcp, udp, icmp)")

    return parser


def find_ue_ip(imsi: str):
    """
    Finds the UE IP address corresponding to the IMSI
    """
    cmd = ["mobility_cli.py", "get_subscriber_table"]
    output = subprocess.check_output(cmd)
    output_str = str(output, "utf-8").strip()
    pattern = "IMSI.*?" + imsi + ".*?([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3})"
    match = re.search(pattern, output_str)
    if match:
        return match.group(1)
    else:
        return "NA"


def output_datapath_actions(
    ue_ip: str, external_ip: str, external_port: str, ue_port: str, protocol: str
):
    """
    Returns the Output of Datapath Actions based as per the supplied UE IP
    """
    service = MagmaService("pipelined", mconfigs_pb2.PipelineD())
    if service.mconfig.nat_enabled:
        in_port = "local"
    else:
        in_port = "patch-up"

    cmd = ["sudo", "ovs-appctl", "ofproto/trace", "gtp_br0"]

    cmd_append = (
        protocol + ",in_port=" + in_port + ",ip_dst=" + ue_ip + ",ip_src=" + external_ip
    )

    if protocol != "icmp":
        cmd_append += (
            ","
            + protocol
            + "_src="
            + external_port
            + ","
            + protocol
            + "_dst="
            + ue_port
        )

    cmd.append(cmd_append)

    print("Running: " + " ".join(cmd))
    output = subprocess.check_output(cmd)
    output_str = str(output, "utf-8").strip()
    pattern = "Datapath\sactions:(.*)"
    match = re.search(pattern, output_str)
    if match:
        return match.group(1).strip()
    else:
        return "NA"


def get_options(args):
    external_ip = args.ip if args.ip else "8.8.8.8"
    external_port = args.port if args.port else "80"
    ue_port = args.ue_port if args.ue_port else "3372"
    protocol = args.protocol if args.protocol else "tcp"

    return {
        "external_ip": external_ip,
        "external_port": external_port,
        "ue_port": ue_port,
        "protocol": protocol,
    }


def main():
    parser = create_parser()
    # Parse the args
    args = parser.parse_args()
    if not args.imsi:
        parser.print_usage()
        exit(1)
    ue_ip = find_ue_ip(args.imsi)
    if ue_ip == "NA":
        print("UE is not connected")
        exit(1)

    print("IMSI: " + args.imsi + ", IP: " + ue_ip)

    input_options = get_options(args)
    dp_actions = output_datapath_actions(
        ue_ip,
        input_options["external_ip"],
        input_options["external_port"],
        input_options["ue_port"],
        input_options["protocol"],
    )

    print("Datapath Actions: " + dp_actions)


if __name__ == "__main__":
    main()
