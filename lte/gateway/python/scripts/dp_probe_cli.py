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
    parser.add_argument("-i", "--imsi", required=True, help="IMSI of the subscriber")
    parser.add_argument(
        "-d",
        "--direction",
        required=True,
        choices=["DL", "UL"],
        help="Direction - DL/UL",
    )
    parser.add_argument(
        "-I", "--ip", nargs="?", const="8.8.8.8", default="8.8.8.8", help="External IP"
    )
    parser.add_argument(
        "-P", "--port", nargs="?", const="80", default="80", help="External Port"
    )
    parser.add_argument(
        "-UP", "--ue_port", nargs="?", const="3372", default="3372", help="UE Port"
    )
    parser.add_argument(
        "-p",
        "--protocol",
        choices=["tcp", "udp", "icmp"],
        nargs="?",
        const="tcp",
        default="tcp",
        help="Portocol (i.e. tcp, udp, icmp)",
    )

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
    return


def output_datapath_actions(
    imsi: str,
    direction: str,
    ue_ip: str,
    external_ip: str,
    external_port: str,
    ue_port: str,
    protocol: str,
):
    """
    Returns the Output of Datapath Actions based as per the supplied UE IP
    """
    service = MagmaService("pipelined", mconfigs_pb2.PipelineD())
    cmd = ["sudo", "ovs-appctl", "ofproto/trace", "gtp_br0"]
    if service.mconfig.nat_enabled:
        in_port = "local"
    else:
        in_port = "patch-up"

    if direction == "DL":

        cmd_append = (
            protocol
            + ",in_port="
            + in_port
            + ",ip_dst="
            + ue_ip
            + ",ip_src="
            + external_ip
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

    elif direction == "UL":
        ingress_tun_id = get_ingress_tunid(ue_ip, in_port)
        if not ingress_tun_id:
            print("Ingress Tunnel not Found")
            exit(1)

        data = get_egress_tunid_and_port(ue_ip, ingress_tun_id)
        if not data:
            print("Egress Tunnel not Found")
            exit(1)

        cmd_append = (
            protocol
            + ",in_port="
            + data["in_port"]
            + ",tun_id="
            + data["tun_id"]
            + ",ip_dst="
            + external_ip
            + ",ip_src="
            + ue_ip
        )

        if protocol != "icmp":
            cmd_append += (
                ","
                + protocol
                + "_src="
                + ue_port
                + ","
                + protocol
                + "_dst="
                + external_port
            )

        cmd.append(cmd_append)
    else:
        return

    print("Running: " + " ".join(cmd))
    output = subprocess.check_output(cmd)
    output_str = str(output, "utf-8").strip()
    pattern = "Datapath\sactions:(.*)"
    match = re.search(pattern, output_str)
    if match:
        return match.group(1).strip()
    else:
        return

def get_ingress_tunid(ue_ip: str, in_port: str):
    cmd = ["sudo", "ovs-ofctl", "dump-flows", "gtp_br0", "table=0"]
    output = subprocess.check_output(cmd)
    output = str(output, "utf-8").strip()
    pattern = (
        ".*?in_port="
        + in_port
        + ",nw_dst="
        + ue_ip
        + ".*?load:(0x(?:[a-z]|[0-9]){1,})->OXM_OF_METADATA.*?"
    )
    match = re.findall(pattern, output, re.IGNORECASE)
    if match:
        return match[0]
    return

def get_egress_tunid_and_port(ue_ip: str, ingress_tun: str):
    cmd = ["sudo", "ovs-ofctl", "dump-flows", "gtp_br0", "table=0"]
    output = subprocess.check_output(cmd)
    output = str(output, "utf-8").strip()
    pattern = pattern = (
        "tun_id=(0x(?:[a-z]|[0-9]){1,}),in_port=([a-z]|[0-9]).*?load:"
        + ingress_tun
        + "->OXM_OF_METADATA.*?\n"
    )
    match = re.findall(pattern, output)
    if match:
        return {"tun_id": match[0][0], "in_port": match[0][1]}
    return

def main():
    parser = create_parser()
    # Parse the args
    args = parser.parse_args()
    ue_ip = find_ue_ip(args.imsi)
    if not ue_ip:
        print("UE is not connected")
        exit(1)

    print("IMSI: " + args.imsi + ", IP: " + ue_ip)

    dp_actions = output_datapath_actions(
        args.imsi,
        args.direction,
        ue_ip,
        args.ip,
        args.port,
        args.ue_port,
        args.protocol,
    )
    if not dp_actions:
        print("Cannot find Datapath Actions for the UE")

    print("Datapath Actions: " + dp_actions)


if __name__ == "__main__":
    main()
