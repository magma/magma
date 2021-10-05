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
from collections import defaultdict
from pathlib import Path


class UserTrace(object):
    """Class for user_trace cli"""

    SESSIOND_ERRORS_LIST = [
        {
            "log": "cause=99",
            "error": "IE is not implemented in Magma",
            "suggestion": "Identify the parameter not supported by Magma, verify with the UE/CPE vendor if you can disable it. If you can't disable, file a feature request",
        },
        {
            "log": "No suitable APN found",
            "error": "Wrong APN configuration",
            "suggestion": "Verify the APN has been provisioned correctly in the phone/CPE and/or Orc8r/HSS",
        },
    ]

    def __init__(self) -> None:
        """Initialize common parameters for user tracing"""

        # Parse the args
        parser = self.create_parser()
        self.args = parser.parse_args()

        with open(self.args.path, "r") as f:
            self.lines = f.readlines()

        self.pattern_time = r"(\d{1,}) \w* (\w* \w* \d{2}:\d{2}:\d{2} \d{4}) "

    def create_parser(self):
        """
        Create the argparse parser with all the arguments.
        """
        parser = argparse.ArgumentParser(
            description="CLI outputs user control signaling. for a \n"
            "supplied IMSI and mme user id",
            formatter_class=argparse.ArgumentDefaultsHelpFormatter,
        )

        parser.add_argument(
            "-p",
            "--path",
            type=Path,
            default=Path(__file__).absolute().parent / "/var/log/mme.log",
            help="Path to the data directory",
        )

        subparsers = parser.add_subparsers(dest="subcmd")

        # list imsi and number of ocurrences in logs
        list_imsi = subparsers.add_parser(
            'list_imsi', help="List imsi and number of ocurrences in logs",
        )
        list_imsi.set_defaults(func=self.get_list_imsi)

        # list mme and enb user id pairs
        list_ue_id = subparsers.add_parser(
            'list_ue_id', help="List mme and enb ue id pairs for a given imsi",
        )
        list_ue_id.add_argument(
            "-i", "--imsi", required=True,
            help="IMSI of the subscriber",
        )
        list_ue_id.set_defaults(func=self.get_list_ue_id)

        # output session trace based on imsi,  mme id
        session_trace = subparsers.add_parser(
            'session_trace', help="Dump session trace for given imsi and user id",
        )
        session_trace.add_argument(
            "-m", "--mme_id", required=True,
            help="mme id of the user", type=int,
        )
        session_trace.add_argument(
            "-i", "--imsi", required=True,
            help="IMSI of the subscriber",
        )
        session_trace.set_defaults(func=self.get_session_trace)

        return parser

    def get_list_imsi(self, args):
        """List imsi and number of ocurrences in logs"""

        imsi_pattern = r"\b\d{15}\b"

        imsi_list = defaultdict(int)
        lines = self.lines
        for line in lines:
            match = re.findall(imsi_pattern, line)
            if not match:
                continue
            imsi = match[0]
            imsi_list[imsi] += 1

        sorted_imsi_list = sorted(imsi_list.items(), key=lambda x: x[1], reverse=True)

        print("\n{0:<15} {1:<15} ".format('IMSI', 'Occurrences'))
        print('-' * 30)

        for v in sorted_imsi_list:
            imsi, number = v
            print("{0:<15} {1:<15}".format(imsi, number))

        print("\nNumber of unique IMSI found in logs: " + str(len(imsi_list)))

    def get_list_ue_id(self, args):
        """
        Provide list of mme id and enodeb id pairs from given an IMSI.
        Some sessions may not have the pair(missing enodeb id). Getting mme_id from Attach Req

        Args:
            args: Filtered logs for a given imsi and mme id

        Returns:
            ue_id_list: list of mme id and enodeb id pairs from given an IMSI
        """

        hex_pattern = "0X.*"
        pattern_mme_enb_ue_id = r".*MME_UE_S1AP_ID = (\w+)\b eNB_UE_S1AP_ID = (\w+)\b"
        ue_id_list = []
        mme_id_list = []
        lines = self.lines
        imsi = self.args.imsi
        for line in lines:
            ue_id_dict = {}
            match = re.findall(imsi + pattern_mme_enb_ue_id, line, re.IGNORECASE)

            if not match:
                continue

            ue_id_dict["mme_ue_id"] = match[0][0]
            ue_id_dict["enb_ue_id"] = match[0][1]

            if re.search(hex_pattern, match[0][0], re.IGNORECASE) and re.search(hex_pattern, match[0][1], re.IGNORECASE):
                ue_id_dict["mme_ue_id"] = int(ue_id_dict["mme_ue_id"], 16)
                ue_id_dict["enb_ue_id"] = int(ue_id_dict["enb_ue_id"], 16)

            if ue_id_dict and ue_id_dict not in ue_id_list:
                ue_id_list.append(ue_id_dict)
                mme_id_list.append(ue_id_dict["mme_ue_id"])

        # Some sessions may not have enodeb id, getting the mme_id from Attach request

        attach_pattern = r"ATTACH REQ \(ue_id = (\w+)\) \(IMSI = " + imsi + r"\)"

        for line in lines:
            match = re.findall(attach_pattern, line)
            if not match:
                continue
            mme_id = int(match[0], 16)
            if mme_id not in mme_id_list:
                mme_id_list.append(mme_id)
                ue_id_list.append({"mme_ue_id": mme_id})

        print("\nIMSI: " + imsi + "\n")
        print("{0:<15} {1:<15} ".format('mme_id', 'enodeb_id'))
        print('-' * 25)

        for v in ue_id_list:
            mme_ue_id = v.get('mme_ue_id')
            enodeb_ue_id = v.get('enb_ue_id')
            if enodeb_ue_id:
                print("{0} 0x{1:<10} {2} 0x{3:<15}".format(mme_ue_id, format(int(mme_ue_id), 'x'), enodeb_ue_id, format(int(enodeb_ue_id), 'x')))
            else:
                print("{0} 0x{1:<10} {2}".format(mme_ue_id, format(int(mme_ue_id), 'x'), enodeb_ue_id))

        return ue_id_list

    def get_session_trace(self, args):
        """Provide a session trace given and IMSI and mme id(ue_id).
        1. Get timestamps for ocurrences where the mme id (ue_id) was found. This will be "session_timestamps"
        2. Filter out logs where the timestamp doesn't match the session_timestamps. This will be "time_filtered_logs"
        3. Filter out logs where line doesn't match imsi or ue_id
        4. From filtered logs, find and show if there are any known session errors.

        Args:
            args: args obtained from cli

        Returns:
            sorted_log: returned session trace sorted on timestamp
        """

        imsi = self.args.imsi
        log_filter = {}
        session_timestamps = self.get_session_timestamp(args)
        time_filtered_logs = self.get_log_time_filtered(session_timestamps)

        self.set_ue_id_pattern(args)
        for line in time_filtered_logs:
            match = re.findall(self.pattern_str, line, re.IGNORECASE)
            match_imsi = re.findall(imsi, line)
            if match or match_imsi:
                match_time = re.findall(self.pattern_time, line)
                log_filter[match_time[0][0]] = line

        sorted_log = sorted(log_filter.items(), key=lambda kv: kv[1])

        [print(value) for key, value in sorted_log]
        self.get_session_error(sorted_log)
        return sorted_log

    def get_session_error(self, session_trace):
        """
        From the session trace, find if any of the logs match any know errors.
        Provide the error found and suggestion to fix it.

        Args:
            session_trace: Filtered logs for a given imsi and mme id
        """

        for _key, value in session_trace:
            for error in self.SESSIOND_ERRORS_LIST:
                if re.search(error["log"], value, re.IGNORECASE):
                    print('*' * 100)
                    print("\nError: " + error["error"])
                    print("\nSuggestion: " + error["suggestion"])
                    print("\nLog: " + value)
                    print('*' * 100)
                    break

    def set_ue_id_pattern(self, args):
        """Generate regex pattern using provided user_id"""
        ue_id = args.mme_id

        pattern_list = [
            r"\bUE id  " + str(ue_id) + r"\b",
            r"\bue 0x0*" + format(ue_id, 'x') + r"\b",
            r"\b\(?ue(_|-)id ?(=|:) ? ?\(?(0x0*" + format(ue_id, 'x') + "|" + str(ue_id) + r")\b\)?",
            "MME_UE_S1AP_ID ?=? 0x0*" + format(ue_id, 'x') + r"\b",
            r"mme_ue_s1ap_id = \(" + str(ue_id) + r"\)",
        ]
        self.pattern_str = "|".join(pattern_list)

    def get_session_timestamp(self, args):
        """Get timestamps for ocurrences where the mme id (ue_id) was found.

        Args:
            args: args obtained from cli

        Returns:
            session_timestamps: List timestamps for ocurrences where the mme id (ue_id) was found.
        """
        session_timestamps = []
        lines = self.lines
        self.set_ue_id_pattern(args)
        for line in lines:
            match = re.findall(self.pattern_str, line, re.IGNORECASE)
            if match:
                match_time = re.findall(self.pattern_time, line)
                if match_time:
                    session_timestamps.append(match_time[0][1])
        return session_timestamps

    def get_log_time_filtered(self, timestamps):
        """Filter out logs where the timestamp doesn't match the session_timestamps

        Args:
            timestamps: session timestamps

        Returns:
            filtered_lines: List of logs where timestamp match session_timestamps.
        """
        lines = self.lines
        pattern_timestamps = "|".join(timestamps)
        filtered_lines = []
        for line in lines:
            if re.search(pattern_timestamps, line):
                filtered_lines.append(line)
        return filtered_lines


def main():
    """ Initialize user trace class and run subparses functions"""
    user_trace = UserTrace()

    if user_trace.args.subcmd:
        user_trace.args.func(user_trace.args)
    exit(1)


if __name__ == "__main__":
    main()
