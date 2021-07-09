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

class DpProbe:
    """Class for dp_probe cli"""

    def __init__(self) -> None:
        """Initialize common ue network parameters and ovs commands"""

        # Dump flow table zero
        self.output_flow_table_zero = self.get_ovs_dump_flow_table_zero()

        # Parse the args
        parser = self.create_parser()
        self.args = parser.parse_args()

        # Get ue network parameters
        self.mtr_port = self.get_mtr0_port()
        self.ue_ip = self.find_ue_ip(self.args.imsi)

        service = MagmaService("pipelined", mconfigs_pb2.PipelineD())
        if service.mconfig.nat_enabled:
            self.ingress_in_port = "local"
        else:
            self.ingress_in_port = "patch-up"

        if self.args.direction == "UL":
            self.ingress_tun_id = self.get_ingress_tunid(self.ue_ip, self.ingress_in_port)

            if not self.ingress_tun_id:
                print("Ingress Tunnel not Found")
                exit(1)

            data = self.get_egress_tunid_and_port(self.ue_ip, self.ingress_tun_id)
            if  data:
                self.egress_tun_id = data["tun_id"]
                self.egress_in_port = data["in_port"]
            else:
                print("Egress Tunnel not Found")
                exit(1)

    def create_parser(self):
        """
        Creates the argparse parser with all the arguments.
        """
        parser = argparse.ArgumentParser(
            description="CLI wrapper around ovs-appctl ofproto/trace.\n"
            "To display the Datapath actions of the supplied IMSI",
            formatter_class=argparse.ArgumentDefaultsHelpFormatter,
        )
        subparsers = parser.add_subparsers(dest="subcmd")
        parser.add_argument(
            "-i", "--imsi", required=True,
            help="IMSI of the subscriber",
        )
        parser.add_argument(
            "-d",
            "--direction",
            required=True,
            choices=["DL", "UL"],
            help="Direction - DL/UL",
        )
        parser.add_argument(
            "-I",
            "--ip",
            nargs="?",
            const="8.8.8.8",
            default="8.8.8.8",
            help="External IP",
        )
        parser.add_argument(
            "-P",
            "--port",
            nargs="?",
            const="80",
            default="80",
            help="External Port",
        )
        parser.add_argument(
            "-UP",
            "--ue_port",
            nargs="?",
            const="3372",
            default="3372",
            help="UE Port",
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
        parser_list_rules = subparsers.add_parser(
            'list_rules', help="List uplink or downlink enforced rules",
        )
        parser_list_rules.set_defaults(func=self.get_enforced_rules)
        parser_get_stats = subparsers.add_parser(
            'stats', help="Output uplink or downlink stats",
        )
        parser_get_stats.set_defaults(func=self.get_stats)
        return parser


    def find_ue_ip(self,imsi: str):
        """Find the UE IP address corresponding to the imsi"""
        cmd = ["mobility_cli.py", "get_subscriber_table"]
        output = subprocess.check_output(cmd)
        output_str = str(output, "utf-8").strip()
        pattern = "IMSI.*?" + imsi + \
            ".*?([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3})"
        match = re.search(pattern, output_str)
        if match:
            return match.group(1)
        return


    def output_datapath_actions(
        self,
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

        cmd = ["sudo", "ovs-appctl", "ofproto/trace", "gtp_br0"]

        if direction == "DL":

            cmd_append = (
                protocol
                + ",in_port="
                + self.ingress_in_port
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
            self.get_ingress_stats(ue_ip)


        elif direction == "UL":

            cmd_append = (
                protocol
                + ",in_port="
                + self.egress_in_port
                + ",tun_id="
                + self.egress_tun_id
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
            self.get_egress_stats(self.egress_tun_id)
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


    def get_ingress_tunid(self, ue_ip: str, in_port: str):
        pattern = (
            ".*?in_port="
            + in_port
            + ",nw_dst="
            + ue_ip
            + ".*?load:(0x(?:[a-z]|[0-9]){1,})->OXM_OF_METADATA.*?"
        )
        match = re.findall(pattern, self.output_flow_table_zero, re.IGNORECASE)
        if match:
            return match[0]
        return


    def get_egress_tunid_and_port(self, ue_ip: str, ingress_tun: str):
        pattern = pattern = (
            "tun_id=(0x(?:[a-z]|[0-9]){1,}),in_port=([a-z]|[0-9]).*?load:"
            + ingress_tun
            + "->OXM_OF_METADATA.*?\n"
        )
        match = re.findall(pattern, self.output_flow_table_zero)
        if match:
            return {"tun_id": match[0][0], "in_port": match[0][1]}
        return


    def get_enforced_rules(self, args):
        """Output enforced rules in pipelined using ue_ip obtained from args"""
        cmd = ["sudo", "pipelined_cli.py", "enforcement", "display_flows"]
        output = subprocess.check_output(cmd, stderr=subprocess.DEVNULL)
        output = str(output, "utf-8").strip()

        if args.direction == "DL":
            pattern = ".*table=enforcement\(main_table\).*nw_dst=" + \
                self.ue_ip + " actions=note:b\'(.*)\',load.*"
            rules_dl = re.findall(pattern, output)
            if len(rules_dl) > 0:
                print("Downlink rules: " + '\n'.join(map(str, rules_dl)))
            else:
                print("No downlink rules found for UE")
        elif args.direction == "UL":
            pattern = ".*table=enforcement\(main_table\).*nw_src=" + \
                self.ue_ip + " actions=note:b\'(.*)\',load.*"
            rules_ul = re.findall(pattern, output)
            if len(rules_ul) > 0:
                print("Uplink rules: " + '\n'.join(map(str, rules_ul)))
            else:
                print("No uplink rules found for UE")


    def get_stats(self, args):
        """Output dowlink or uplink stats"""

        if args.direction == "DL":
            return self.get_ingress_stats(self.ue_ip)
        elif args.direction == "UL":
            return self.get_egress_stats(self.egress_tun_id)


    def get_ingress_stats(self, ue_ip: str):
        """Output ingress Downlink stats using ue_ip obtained from args"""

        pattern_local = "(n_packets=[0-9]{1,}.*n_bytes=[0-9]{1,}).*in_port=LOCAL,nw_dst=" + ue_ip  +  ".*actions.*"
        match_local = re.findall(pattern_local, self.output_flow_table_zero)

        pattern_mtr0 =  "(n_packets=[0-9]{1,}.*n_bytes=[0-9]{1,}).*in_port="+self.mtr_port+",nw_dst=" + ue_ip  +  ".*actions.*"
        match_mtr0 = re.findall(pattern_mtr0, self.output_flow_table_zero)

        if match_local:
            print("LOCAL: " + match_local[0])

        if match_mtr0:
            print("mtr0: " + match_mtr0[0])


    def get_egress_stats(self, tun_id: str):
        """Output egress Uplink stats using ue_ip obtained from args"""

        pattern = "(n_packets=[0-9]{1,}.*n_bytes=[0-9]{1,}).*tun_id="+ tun_id +",in_port.*"
        match = re.findall(pattern, self.output_flow_table_zero)
        print("UL stats: " + match[0])


    def get_mtr0_port(self):
        """Output mtr0 port number"""
        cmd = ["ovs-ofctl", "show", "gtp_br0"]
        output = subprocess.check_output(cmd)
        output = str(output, "utf-8").strip()
        pattern = "([0-9]{1,})\(mtr0\)"
        match_mtr0 = re.findall(pattern, output)
        if match_mtr0:
            return match_mtr0[0]


    def get_ovs_dump_flow_table_zero(self):
        """Output ovs dump flow table zero"""
        cmd = ["sudo", "ovs-ofctl", "dump-flows", "gtp_br0", "table=0"]
        output = subprocess.check_output(cmd)
        output_str = str(output, "utf-8").strip()
        return output_str


def main():
    dp_probe = DpProbe()
    if not dp_probe.ue_ip:
        print("UE is not connected")
        exit(1)
    print("IMSI: " + dp_probe.args.imsi + ", IP: " + dp_probe.ue_ip)

    # If there are subcomands, execute subcommands only
    if dp_probe.args.subcmd:
        dp_probe.args.func(dp_probe.args)
        exit(1)

    dp_actions = dp_probe.output_datapath_actions(
        dp_probe.args.imsi,
        dp_probe.args.direction,
        dp_probe.ue_ip,
        dp_probe.args.ip,
        dp_probe.args.port,
        dp_probe.args.ue_port,
        dp_probe.args.protocol,
    )
    if not dp_actions:
        print("Cannot find Datapath Actions for the UE")

    print("Datapath Actions: " + dp_actions)
    dp_probe.get_enforced_rules(dp_probe.args)


if __name__ == "__main__":
    main()
