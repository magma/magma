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
import importlib
import json
import logging
import math
import os
import random
import re
import statistics
import subprocess
import sys
import time
from datetime import datetime, timedelta
from ipaddress import ip_address, ip_network
from multiprocessing import Process, Queue
from typing import Dict, List, Union

import attr
import config
import get_ports
import pyshark
import requests
from jinja2 import Template
from netaddr import *
from TS.kpi_lib import cruncher


@attr.s
class TSBase(object):
    template = attr.ib()
    library_id = attr.ib()
    auth = attr.ib()

    def __attrs_post_init__(self):
        self.override_url = (
            "http://"
            + config.TAS.get("tas_ip")
            + ":8080/api/libraries/"
            + self.library_id
            + "/testSessions/"
            + self.template
            + "?action=overrideAndSaveAs"
        )
        self.url = "http://" + config.TAS.get("tas_ip") + ":8080/api/runningTests"
        self.timeout = config.TAS.get("timeout")

    def get_base_constants(self):
        pass

    def imsi_calc(self, invalid_sub: int, Imsi: str) -> str:
        output = "00" + str(int(Imsi) - invalid_sub)
        return output

    def save_test(self, payload: Dict[str, any]):
        response = requests.request(
            "POST", self.override_url, data=json.dumps(payload), auth=self.auth,
        )
        return response

    def get_test(self, saved_test_name: str):
        response = requests.request(
            "GET",
            "http://"
            + config.TAS.get("tas_ip")
            + ":8080/api/libraries/"
            + self.library_id
            + "/testSessions/"
            + saved_test_name,
            auth=self.auth,
        )
        return response

    def run_test(self, saved_test_name: str = None):
        payload = {
            "library": self.library_id,
            "name": saved_test_name or self.save_name,
        }
        try:
            response = requests.request(
                "POST", self.url, data=json.dumps(payload), auth=self.auth,
            )
            return json.loads(response.text)
        except:
            try:
                time.sleep(10)
                response = requests.request(
                    "POST", self.url, data=json.dumps(payload), auth=self.auth,
                )
                return json.loads(response.text)
            except Exception as e:
                logging.error(
                    f"Error Running test case on spirent {e} response {response}",
                )

    def check_verify_agw_imsi_state(self, gateway: str = "") -> bool:
        _ue_state_verification = {}
        expected_state = config.MAGMA_AGW.get("UE_STATE", {})
        current_state = self.SutObject.sut_magma_state_check(role="agw-state-check")

        logging.info(
            f"Current SUT UE state for its processes is "
            f"{current_state}\n"
            f"Expected state is {expected_state}\n",
        )
        for k, v in expected_state.items():
            if current_state[k] != v:
                _ue_state_verification[k] = False
        return _ue_state_verification

    def check_test_summary(self, **kwargs: Union[str, int, float]) -> Dict[str, any]:
        """check test summary with expected result and gives verdict"""

        def _s1ap_check(s1ap: Dict[str, any]) -> Dict[str, bool]:
            logging.debug(
                f"{kwargs['tc_name']} S1-AP Setup Requests - {s1ap['Setup Requests Sent']}",
            )
            logging.debug(
                f"{kwargs['tc_name']} S1-AP Setup Responses - {s1ap['Setup Responses Received']}",
            )
            logging.debug(
                f"{kwargs['tc_name']} S1-AP Init Context Setup Failures - {s1ap['Init Context Setup Failures']}",
            )
            _s1ap_decision = {}
            tolerance = config.MAGMA.get("control_procedure_tolerance_pct")

            _s1ap_decision["S1_Setup_Procedure"] = (
                (
                    s1ap["Setup Requests Sent"] == str(kwargs["enbs"])
                    and s1ap["Setup Responses Received"] == str(kwargs["enbs"])
                )
                and True
                or False
            )
            _s1ap_decision["S1_Release_Timeout"] = (
                (
                    int(s1ap["S1 Release Timeouts"])
                    <= int(s1ap["S1 Release Commands"]) * (tolerance / 100)
                )
                and True
                or False
            )
            _s1ap_decision["Init_Context_Setup_Failure"] = (
                (
                    int(s1ap["Init Context Setup Failures"])
                    <= int(s1ap["Init Context Setup Requests"]) * (tolerance / 100)
                )
                and True
                or False
            )

            return _s1ap_decision

        def _handover_check(
            s1ho: Dict[str, any], test_summary: Dict[str, any],
        ) -> Dict[str, bool]:
            logging.debug(
                f"{kwargs['tc_name']} S1-AP Setup Requests - {s1ho['Setup Requests Sent']}",
            )
            logging.debug(
                f"{kwargs['tc_name']} S1-AP Setup Responses - {s1ho['Setup Responses Received']}",
            )
            logging.debug(
                f"{kwargs['tc_name']} S1 Handover Completed - {s1ho['S1 Handover Count']}",
            )
            logging.debug(
                f"{kwargs['tc_name']} X2 Handover Completed - {s1ho['X2 Handover Count']}",
            )
            logging.debug(
                f"{kwargs['tc_name']} Total Handover Attempts - {test_summary['Handoff Attempts']}",
            )
            logging.debug(
                f"{kwargs['tc_name']} Total Handover Successes - {test_summary['Handoff Successes']}",
            )
            _ho_decision = {}

            _ho_decision["S1_Setup_Procedure_serving"] = (
                (
                    int(s1ho["Setup Requests Sent"]) // 2 == kwargs["enbs"]
                    and int(s1ho["Setup Responses Received"]) // 2 == kwargs["enbs"]
                )
                and True
                or False
            )
            _ho_decision["X2-Handover"] = (
                (
                    int(s1ho["Path Switch Requests"])
                    == int(s1ho["Path Switch Request Acks"])
                )
                and True
                or False
            )

            _ho_decision["S1-Handover"] = (
                (int(s1ho["Handover Required"]) == int(s1ho["Handover Notifications"]))
                and True
                or False
            )

            return _ho_decision

        def _emm_check(emm: Dict[str, any]) -> Dict[str, bool]:
            logging.info(
                f"{kwargs['tc_name']} EMM Attach Req - {emm['Attach Requests']}",
            )
            logging.info(
                f"{kwargs['tc_name']} EMM Attach Accepts - {emm['Attach Accepts']}",
            )
            logging.info(f"{kwargs['tc_name']} EMM Timeouts - {emm['Timeouts']}")
            _emm_decision = {}
            tolerance = config.MAGMA.get("control_procedure_tolerance_pct")
            _emm_decision["Attach_Procedure"] = (
                (
                    int(emm["Attach Requests"]) >= t_subs
                    and int(emm["Attach Accepts"]) >= t_subs * 0.99
                )
                and True
                or False
            )  # 99 percentile can be PASS

            _emm_decision["Service_Rejects"] = (
                (
                    int(emm["Service Rejects"])
                    <= int(emm["Service Requests"]) * (tolerance / 100)
                )
                and True
                or False
            )  # SR failure and Timeout will be based on a 0.1% of total service requests
            _emm_decision["Service_Request_Timeout"] = (
                (
                    int(emm["Service Request Timeouts"])
                    <= int(emm["Service Requests"]) * (tolerance / 100)
                )
                and True
                or False
            )

            return _emm_decision

        def _esm_check(esm: Dict[str, any]) -> Dict[str, bool]:
            logging.info(
                f"{kwargs['tc_name']} ESM PDN Req - {esm['PDN Connectivity Requests']}",
            )
            logging.info(
                f"{kwargs['tc_name']} ESM PDN Connectivity Rejects - {esm['PDN Connectivity Rejects']}",
            )
            logging.info(
                f"{kwargs['tc_name']} ESM PDN Connectivity Success Count - {esm['PDN Connectivity Success Count']}",
            )
            _esm_decision = {}

            _esm_decision["PDN_Con_Procedure"] = (
                (
                    int(esm["PDN Connectivity Requests"]) >= t_subs
                    and int(esm["PDN Connectivity Success Count"]) >= t_subs * 0.99
                )
                and True
                or False
            )  # 99 percentile can be PASS

            return _esm_decision

        def _check_active_idle(l3_client: Dict[str, any]) -> Dict[str, bool]:
            logging.info(f"Total Pings sent: {l3_client['Total Pings Sent']}")
            logging.info(
                f"Total Responses recvd: {l3_client['Total Ping Replies Received']}",
            )
            _active_idle_decision = {}
            try:
                _active_idle_decision["Active_Idle_Data_Check"] = (
                    (
                        (
                            int(l3_client["Total Ping Replies Received"])
                            / int(l3_client["Total Pings Sent"])
                        )
                        >= config.MAGMA.get("pass_percentile")
                    )
                    and True
                    or False
                )
            except Exception as e:
                logging.error(f"Maybe no pings were sent? msg: {e} ")

            return _active_idle_decision

        def _calculate_availability() -> Dict[str, bool]:
            """This func would gets reports for mme nodal test 1 runs 1 UE per second and do send 1 ping
            it calculate Availability form test run duration and report that in output
            """

            url = self.url + "/" + str(kwargs["id"]) + "/measurements?Ts=0&Tc=2"
            _res = get_ports.get(url, self.auth)
            _res = _res["testservers"][0]["testcases"][
                2
            ]  # get ts=0 and tc=2 results as mme nodal 1 is on this location

            logging.info(
                f"{kwargs['tc_name']} Attempted Session Connects - {_res['tabs']['Test Summary']['Attempted Session Connects']}",
            )
            logging.info(
                f"{kwargs['tc_name']} Actual Session Connects - {_res['tabs']['Test Summary']['Actual Session Connects']}",
            )
            logging.info(
                f"{kwargs['tc_name']} Attempted Session Disconnects - {_res['tabs']['Test Summary']['Attempted Session Disconnects']}",
            )
            logging.info(
                f"{kwargs['tc_name']} Actual Session Disconnects - {_res['tabs']['Test Summary']['Actual Session Disconnects']}",
            )
            logging.info(
                f"{kwargs['tc_name']} Session Errors - {_res['tabs']['Test Summary']['Session Errors']}",
            )
            logging.info(
                f"{kwargs['tc_name']} Data Verify Attempts - {_res['tabs']['L3 Client']['Data Verify Attempts']}",
            )
            logging.info(
                f"{kwargs['tc_name']} Data Verify Successes - {_res['tabs']['L3 Client']['Data Verify Successes']}",
            )
            logging.info(
                f"{kwargs['tc_name']} Data Verify Failures - {_res['tabs']['L3 Client']['Data Verify Failures']}",
            )
            availability_calculation = {}

            Total_Test_Duration = int(
                _res["tabs"]["Test Summary"]["Attempted Session Connects"],
            )
            Total_Sessions_Connects = int(
                _res["tabs"]["Test Summary"]["Actual Session Connects"],
            )
            Total_Data_Traffic_Verified = int(
                _res["tabs"]["L3 Client"]["Data Verify Successes"],
            )

            Total_Unavailable_Time_Control = (
                Total_Test_Duration - Total_Sessions_Connects
            )
            Total_Unavailable_Time_Data = (
                Total_Sessions_Connects - Total_Data_Traffic_Verified
            )

            Total_Unavailable_Time = (
                Total_Unavailable_Time_Control + Total_Unavailable_Time_Data
            )

            System_Availability = (
                (Total_Test_Duration - Total_Unavailable_Time) / Total_Test_Duration
            ) * 100

            logging.info(
                f"{kwargs['tc_name']} Total test time - {Total_Test_Duration} seconds, Total Unavailable time - {Total_Unavailable_Time} seconds, System Availability - {System_Availability}%",
            )

            availability_calculation["System_Availability"] = System_Availability

            return availability_calculation

        def _check_pcap(proc: str) -> Dict[str, bool]:
            """
            proc == HEADER_ENRICHMENT or IPFIX
            1. Check if the pcap stopped
            2. Move pcap from SUT to local temp directory on script execution machine
            3. Parse through script using pyshark
            """
            # Let's delete the trace of the pcaps!
            def _del_pcap(proc: str):
                cap_vars = {
                    "output_file": getattr(config, proc).get(
                        "output_file", "/tmp/test.pcapng",
                    ),
                }
                self.pcap_operations(
                    cmdline="--tags pcap_cleanup",
                    cap_vars=cap_vars,
                    sut=kwargs["gateway"],
                    **kwargs,
                )
                subprocess.call([f"rm -r /tmp/{dest_name}"], shell=True)

            dest_name = os.path.split(getattr(config, proc).get("output_file"))[1]
            res = self.pcap_operations(
                cmdline="--tags pcap_check",
                cap_vars={},
                sut=kwargs["gateway"],
                **kwargs,
            )

            tshark_pid = res.get_fact_cache(kwargs["gateway"])["tshark_pid"]
            tshark_pid_post_run = res.get_fact_cache(kwargs["gateway"])[
                "tshark_pid_post_run"
            ]
            logging.info(f"PCAP procedure is {proc}.")

            if proc == "HEADER_ENRICHMENT":
                logging.debug("Header Enrichment validation in progress.")
                _he_decision = {}

                if tshark_pid in tshark_pid_post_run:
                    logging.info("PCAP did not finish! Test failed")
                    _he_decision["HE_wo_ciphering"] = False
                    logging.info(
                        f"Force terminating the pcap on SUT with pid {tshark_pid}",
                    )
                    self.pcap_operations(
                        cmdline="--tags pcap_kill",
                        cap_vars={"tshark_id": str(tshark_pid)},
                        sut=kwargs["gateway"],
                        **kwargs,
                    )
                    _del_pcap(proc)
                    return _he_decision

                self.SutObject.get_sut_file(
                    file_path=getattr(config, proc).get("output_file"),
                    dest_location="/tmp",
                    dest_name=dest_name,
                )

                enriched = 0
                not_enriched = 0
                cap = pyshark.FileCapture(f"/tmp/{dest_name}")
                for n in cap:
                    try:
                        if n.http.get_field_by_showname("imsi"):
                            enriched += 1
                        else:
                            not_enriched += 1
                    except Exception as e:
                        logging.error(f"pcap parsing faild: {e}")

                enriched_pct = enriched / (enriched + not_enriched)
                logging.info(
                    f"{enriched_pct*100}% of all HTTP GET requests were enriched!",
                )

                _he_decision[
                    "HE_wo_ciphering"
                ] = True and enriched_pct > config.MAGMA.get("pass_percentile")
                _del_pcap(proc)

                return _he_decision

            elif proc == "IPFIX":
                logging.debug("IPFIX validation in progress.")
                _ipfix_decision = {}
                if tshark_pid in tshark_pid_post_run:
                    logging.warning("PCAP did not finish! Terminating and proceeding.")
                    logging.warning(
                        f"Force terminating the pcap on SUT with pid {tshark_pid}",
                    )
                    self.pcap_operations(
                        cmdline="--tags pcap_kill",
                        cap_vars={"tshark_id": str(tshark_pid)},
                        sut=kwargs["gateway"],
                        **kwargs,
                    )
                self.SutObject.get_sut_file(
                    file_path=getattr(config, proc).get("output_file"),
                    dest_location="/tmp",
                    dest_name=dest_name,
                )
                cap = pyshark.FileCapture(
                    f"/tmp/{dest_name}", decode_as={"udp.port==65010": "cflow"},
                )
                hosts = []
                for n in cap:
                    try:
                        if "no_template_found" in n.cflow.field_names:
                            _ipfix_decision["IPFIX_Template_Not_Found"] = False
                            logging.error(
                                "No template file found in this run. Sorry :( Run again after restarting pipelined",
                            )
                            return _ipfix_decision
                        elif "template_id" in n.cflow.field_names:
                            continue
                        else:
                            addr = str(n.cflow.srcaddr)
                            if addr not in hosts:
                                hosts.append(str(n.cflow.srcaddr))
                            else:
                                continue
                    except AttributeError:
                        continue
                    except Exception as e:
                        logging.error(e)

                ipfix_pct = len(hosts) / t_subs
                _ipfix_decision["IPFIX_Export"] = True and ipfix_pct > config.MAGMA.get(
                    "pass_percentile",
                )
                logging.info(f"IPFIX received pct is: {ipfix_pct}")
                _del_pcap(proc)
                return _ipfix_decision

            else:
                logging.error("Sorry, this feature has not been implemented yet.")

        def _enb_check(enb: Dict[str, any]) -> Dict[str, bool]:
            _enb_decision = {}
            _enb_decision["gtpu_echo"] = (
                (
                    int(enb["Echo Requests Sent"]) > 0
                    and int(enb["Echo Requests Sent"])
                    == int(enb["Echo Responses Received"])
                )
                and True
                or False
            )
            logging.debug(f"GTP-U Echo requests sent {enb['Echo Requests Sent']}")
            logging.debug(
                f"GTP-U Echo responses received {enb['Echo Responses Received']}",
            )

            return _enb_decision

        def _data_traffic_check(rid: str, test_suite: str) -> Dict[str, bool]:
            def _data_traffic_get(filename: str):
                logging.debug("Reading the excel output from TAS")
                test_data = cruncher()
                values = test_data.read(
                    filename=filename,
                    sheets=["L3 Client", "L3 Server"],
                    columns=[
                        "Total Bits Received/Sec  (P-I)",
                        "Total Bits Sent/Sec  (P-I)",
                    ],
                )
                logging.debug(f'L3 Server sent bits/sec {list(values["L3 Server"])}')
                logging.debug(
                    f'L3 Client received bits/sec {list(values["L3 Client"])}',
                )

                l3_client_bitrate = test_data.quant(values["L3 Client"])
                l3_server_bitrate = test_data.quant(values["L3 Server"])
                logging.debug("Computed l3_server_bitrate and l3_client_bitrate")
                return (l3_server_bitrate, l3_client_bitrate)

            if file_dl_check:
                l3_server_bitrate, l3_client_bitrate = _data_traffic_get(
                    tc_res_filepath,
                )
            else:
                logging.error(
                    "Something went wrong with getting the TAS excel output for this run",
                )
                l3_server_bitrate = [0]
                l3_client_bitrate = [0]

            logging.info(
                f"{kwargs['tc_name']} L3 Server Total Bits Sent/Sec - {l3_server_bitrate}",
            )
            logging.info(
                f"{kwargs['tc_name']} L3 Client Total Bits Received/Sec - {l3_client_bitrate}",
            )
            data_traffic_decision = {}
            tput_multiplier = config.MAGMA.get("pass_percentile")
            datarate = kwargs.get("expected_qos_flow_datarate") or (
                t_subs
                * int(re.search(r"\d+", kwargs["dmf"]).group())
                * 1000
                * tput_multiplier
                if "K_" in kwargs["dmf"]
                else t_subs
                * int(re.search(r"\d+", kwargs["dmf"]).group())
                * 1000000
                * tput_multiplier
            )
            logging.info(f"{kwargs['tc_name']} Expected Received/Sec - {datarate}")
            if test_suite == "SANITY":
                data_traffic_decision["L3_Data_Traffic"] = (
                    (
                        l3_server_bitrate[-1] >= datarate
                        and l3_client_bitrate[-1] >= datarate
                    )
                    and True
                    or False
                )
            else:
                data_traffic_decision["L3_Data_Traffic"] = (
                    (
                        l3_server_bitrate[1] >= datarate
                        and l3_client_bitrate[1] >= datarate
                    )
                    and True
                    or False
                )
            return data_traffic_decision

        # get test detailed report file
        file_dl_check, tc_res_filepath = self.get_tc_result_file(
            kwargs["tc_name"], kwargs["gateway"],
        )

        t_subs = kwargs.get("verify_UE_online") or (
            kwargs.get("total_subs", 0)
            + kwargs.get("nodal2_total_subs", 0)
            + kwargs.get("nodal3_total_subs", 0)
        )
        logging.info(f"Expected UEs to attach {t_subs}")
        url = self.url + "/" + str(kwargs["id"]) + "/measurements"
        time.sleep(15)  # some reports updates after test get completed
        _res = get_ports.get(url, self.auth)
        categories = kwargs["categories"]
        tsum = {}
        cat = ""

        try:
            tsum["Timeout"] = kwargs["timeout"]
            for cat in categories:
                if cat == "S1-AP":
                    tsum[cat] = _s1ap_check(_res["tabs"][cat])
                elif cat == "ESM":
                    tsum[cat] = _esm_check(_res["tabs"][cat])
                elif cat == "Data_Traffic":
                    tsum[cat] = _data_traffic_check(kwargs["id"], kwargs["test_suite"])
                elif cat == "Handover":
                    tsum[cat] = _handover_check(
                        _res["tabs"]["S1-AP"], _res["tabs"]["Test Summary"],
                    )
                elif cat == "AVAILABILITY":
                    tsum[cat] = _calculate_availability()
                elif cat == "UE_STATE_CHECK":
                    tsum[cat] = self.check_verify_agw_imsi_state(kwargs["gateway"])
                elif cat == "HE":
                    tsum[cat] = _check_pcap("HEADER_ENRICHMENT")
                elif cat == "IPFIX":
                    tsum[cat] = _check_pcap("IPFIX")
                elif cat == "active_idle":
                    tsum[cat] = _check_active_idle(_res["tabs"]["L3 Client"])
                elif cat == "ENB":
                    tsum[cat] = _enb_check(_res["tabs"]["eNodeB Node"])
                else:
                    tsum[cat] = _emm_check(_res["tabs"][cat])
        except Exception as e:
            logging.error(f"Encountered {e}; check in category {cat}.")

        return tsum, _res

    def pcap_operations(self, **kwargs: Union[int, str, float]) -> Dict[str, any]:
        if kwargs.get("he", False) or kwargs.get("ipfix", False):
            """Run PCAP collection Playbook"""
            if not kwargs.get("cap_vars"):
                if kwargs.get("he", False):
                    logging.info("Configuring pcap for header enrichment")
                    cap_vars = config.HEADER_ENRICHMENT
                    cap_vars.update({"packet_count": str(kwargs.get("total_subs", 1))})
                elif kwargs.get("ipfix", False):
                    logging.info("Configuring pcap for IPFIX")
                    cap_vars = config.IPFIX
                    cap_vars.update(
                        {"packet_count": str(int(kwargs.get("total_subs", 1) * 10))},
                    )  # PCAP will be stopped by the validator
                else:
                    cap_vars = {}
            else:
                cap_vars = kwargs["cap_vars"]

            res = self.SutObject.pcap_collection(
                role="pcap", cmdline=kwargs["cmdline"], cap_vars=cap_vars,
            )
            return res
        else:
            logging.error("Please run a valid PCAP using TC")

    def check_and_wait_for_tc(self, **kwargs: Union[int, str, float]) -> Dict[str, any]:
        """Method to wait for test to complete"""
        url = self.url + "/" + str(kwargs["id"])
        self.test_run_id = kwargs["id"]
        start_time = datetime.now()
        delta_time = timedelta(
            seconds=self.timeout
            if self.timeout > kwargs["tc_dpruntime"]
            else kwargs["tc_dpruntime"] + self.timeout,
        )
        test_status = {"timeout": False, "l3_server_bitrate": 0, "l3_client_bitrate": 0}

        sleep_profile = config.SPIRENT_SLEEP_PROFILE
        tc_dpruntime = kwargs.get("tc_dpruntime", 60)
        sleep_time = sleep_profile[
            min(
                range(len(sleep_profile)),
                key=lambda i: abs(sleep_profile[i] - tc_dpruntime),
            )
        ]

        logging.info(f"kwargs {kwargs}")
        if kwargs.get("pcap", False):
            logging.debug("Found a pcap TC; dialing pcap ops!")
            self.pcap_operations(cmdline="--tags pcap_start", **kwargs)

        while True:
            if datetime.now() > start_time + delta_time and not test_status["timeout"]:
                logging.error(f"Test id {str(kwargs['id'])} timed out!")
                self.force_stop_test(kwargs["id"])
                test_status["timeout"] = True

            tc_res = get_ports.get(url, self.auth)
            if tc_res["testStateOrStep"] == "COMPLETE":
                logging.info(f"Test id {str(kwargs['id'])} is completed")
                if kwargs.get("sequencer_mode", False):
                    kwargs["queue"].put(test_status)  # Only needed for sequencer mode
                return test_status
            else:
                logging.debug(f"Test id {str(kwargs['id'])} on going")
                logging.debug(f"Sleeping for {sleep_time} seconds")
                time.sleep(sleep_time)

    def force_stop_test(self, id: int) -> None:
        logging.debug(f"Stop running test with id {str(id)}, as it seems stuck")
        url = self.url + "/" + str(id) + "?action=stop"
        payload = {}
        response = requests.request(
            "POST", url, data=json.dumps(payload), auth=self.auth,
        )

    def get_tc_steps(
        self, *args: None, **kwargs: Union[int, str, float]
    ) -> List[Dict[str, any]]:
        """get TC steps if required, TC can pass its own function to override this"""

        return config.SPIRENT.get("1nodal_1nw_steps", [])(
            run_time=kwargs["tc_dpruntime"],
        )

    def get_tc_mme_nodal(self, **kwargs: Union[str, int, float]) -> Dict[str, any]:
        """return mme nodal config as required by TC"""

        def _sequencer_config(mme_nodal: Dict[str, any]) -> Dict[str, any]:
            if kwargs.get("sequencer_mode", False):
                inter_nodal_offset = str(
                    random.randint(5, 15),
                )  # Some random value to separate each nodal further. This can be made into a wait; but that's more API calls!
                detach_delay = (
                    "1"
                    if kwargs["total_subs"] < kwargs["rate"]
                    else str(int(kwargs["total_subs"] / kwargs["rate"]))
                )

                command_template = Template(
                    config.SPIRENT_SEQUENCER_COMMAND.get("active_idle"),
                )
                mme_nodal["parameters"]["CommandSequence"] = command_template.render(
                    rate=kwargs["rate"],
                    total_subs=kwargs["total_subs"],
                    inter_nodal_offset=inter_nodal_offset,
                    iterations=kwargs["iterations"],
                    detach_delay=detach_delay,
                )
                mme_nodal["parameters"]["MobMmeSut"] = {
                    "class": "Sut",
                    "name": kwargs["gateway"],
                }
            else:
                mme_nodal["parameters"]["StartRate"] = kwargs["rate"]
                mme_nodal["parameters"]["DisconnectRate"] = kwargs["rate"]
                mme_nodal["parameters"]["HoldTime"] = str(kwargs.get("HoldTime", "610"))
                mme_nodal["parameters"]["IdleTime"] = str(kwargs.get("IdleTime", "30"))

            return mme_nodal

        def _dmf_config(mme_nodal: Dict[str, any]) -> Dict[str, any]:
            if kwargs.get("dmf_conf", False) and not kwargs.get(
                "sequencer_mode", False,
            ):
                mme_nodal["parameters"]["Dmf"] = {
                    "_replace": "true",
                    "mainflows": [{"library": 342, "name": kwargs["dmf"]}],
                    "instanceGroups": [
                        {
                            "mainflowIdx": 0,
                            "mixType": "",
                            "rate": 0.0,
                            "rows": [
                                {
                                    "clientPort": 0,
                                    "context": 0,
                                    "node": 0,
                                    "ratingGroup": 0,
                                    "role": 0,
                                    "serviceId": 0,
                                    "transport": kwargs.get("dmf_transport", "Any"),
                                },
                            ],
                        },
                    ],
                }

            return mme_nodal

        def _nw_host_add_remote_config(mme_nodal: Dict[str, any]) -> Dict[str, any]:
            if kwargs.get("NetworkHostAddrRemote", False):
                mme_nodal["parameters"]["NetworkHostAddrRemote"] = {
                    "class": "Array",
                    "array": [kwargs["start_ip_for_network_host"]],
                }
            return mme_nodal

        def _EnbUserAddr_config(mme_nodal: Dict[str, any]) -> Dict[str, any]:
            if kwargs.get("EnbUserAddr", False):
                mme_nodal["parameters"]["EnbUserAddr"] = {
                    "ip": kwargs["start_ip_for_user_enb"],
                    "mac": kwargs["start_mac_user_enb"],
                    "nextHop": str(kwargs["nodal_network"][1]),
                    "numLinksOrNodes": kwargs["enbs"],
                    "numVlan": 1,
                    "phy": kwargs["port_name"],
                    "vlanId": kwargs["vlan_enb"],
                }
            return mme_nodal

        def _MobEnbControlAddr(mme_nodal: Dict[str, any]) -> Dict[str, any]:
            if kwargs.get("MobEnbControlAddr", False):
                mme_nodal["parameters"]["MobEnbControlAddr"] = {
                    "ip": kwargs["start_ip_for_mob_control_enb"],
                    "mac": kwargs["start_mac_mob_control_enb"],
                    "nextHop": str(kwargs["nodal_network"][1]),
                    "numLinksOrNodes": kwargs["MobEnbControlAddr_enb"]
                    if kwargs.get("MobEnbControlAddr_enb")
                    else kwargs.get("enbs"),
                    "numVlan": 1,
                    "phy": kwargs["port_name"],
                    "vlanId": kwargs["vlan_enb"],
                }

            return mme_nodal

        def _s1_handover(mme_nodal: Dict[str, any]) -> Dict[str, any]:
            if kwargs.get("s1_handover", False):
                mme_nodal["parameters"]["MobMmeSut"] = {
                    "class": "Sut",
                    "name": kwargs["gateway"],
                }
            return mme_nodal

        def _mme_nodal_nw_host_dual_stck(mme_nodal: Dict[str, any]) -> Dict[str, any]:
            if kwargs.get("mme_nodal_nw_host_dual_stack", False):
                mme_nodal["parameters"]["Ipv4NetworkHostAddrLocal"] = {
                    "ip": str(kwargs["start_ip_for_network_host"]),
                    "mac": str(kwargs["start_mac_network_host"]),
                    "nextHop": str(kwargs["network_host"][1]),
                    "numLinksOrNodes": 1,
                    "phy": str(kwargs["port_name"]) + "_1",
                    "vlanId": kwargs["vlan_network_host"],
                }
                mme_nodal["parameters"]["Ipv6NetworkHostAddrLocal"] = {
                    "ip": str(kwargs["start_ip_for_network_host_v6"]),
                    "mac": str(kwargs["start_mac_network_host"]),
                    "nextHop": str(kwargs["network_host_v6"][1]),
                    "numLinksOrNodes": 1,
                    "phy": str(kwargs["port_name"]) + "v6_1",
                }
            return mme_nodal

        mme_nodal = {
            "name": "",
            "type": "MME Nodal",
            "parameters": {
                "EnbControlAddr": {
                    "ip": kwargs["start_ip_for_control_enb"],
                    "mac": kwargs["start_mac_control_enb"],
                    "nextHop": str(kwargs["nodal_network"][1]),
                    "numLinksOrNodes": kwargs["enbs"],
                    "numVlan": 1,
                    "phy": kwargs["port_name_v6"],
                    "vlanId": kwargs["vlan_enb"],
                },
                "Imei": kwargs.get("mme_nodal_Imei", "11111111111111"),
                "HomeAddrType": kwargs.get("pdn_type", "46"),
                "Imsi": kwargs["imsi_start"],
                "EnbMcc": kwargs.get("mcc", "001"),
                "EnbMnc": kwargs.get("mnc", "01"),
                "Sessions": kwargs["total_subs"],
                "MmeSut": {"class": "Sut", "name": kwargs["gateway"]},
            },
        }
        return _mme_nodal_nw_host_dual_stck(
            _s1_handover(
                _MobEnbControlAddr(
                    _EnbUserAddr_config(
                        _nw_host_add_remote_config(
                            _dmf_config(_sequencer_config(mme_nodal)),
                        ),
                    ),
                ),
            ),
        )

    def get_tc_extra_phy_subnet(
        self, **kwargs: Union[str, int, float]
    ) -> Dict[str, any]:
        """add phy subnet as required by TC"""

        return {
            "base": str(kwargs["start_ip_for_network_host"]),
            "mask": str(kwargs["cidr_mask_network_host"]),
            "name": str(kwargs["port_name"]) + "_1",
            "numIps": kwargs["num_of_ip_network_host"] - 5,
        }

    def get_tc_ipv6_subnet(self, **kwargs: Union[str, int, float]) -> Dict[str, any]:
        """add phy subnet as required by TC"""

        return {
            "base": str(kwargs["start_ip_for_network_host_v6"]),
            "mask": str(kwargs["cidr_mask_network_host_v6"]),
            "name": str(kwargs["port_name"]) + "v6_1",
            "numIps": kwargs["num_of_ip_network_host_v6"],
        }

        # Making this dule stack for all

    def get_tc_nw_host(self, **kwargs: Union[str, int, float]) -> Dict[str, any]:
        """return network host config as required by TC"""

        def _dmf_config(nw_host: Dict[str, any]) -> Dict[str, any]:
            if kwargs.get("dmf_conf", False):
                nw_host["parameters"]["Dmf"] = {
                    "mainflows": [{"library": 342, "name": kwargs["dmf"]}],
                }
            return nw_host

        nw_host = {
            "name": "",
            "type": "Network Host",
            "parameters": {
                "NetworkHostAddrLocal": {
                    "ip": kwargs["start_ip_for_network_host"],
                    "mac": kwargs["start_mac_network_host"],
                    "nextHop": str(kwargs["network_host"][1]),
                    "numLinksOrNodes": 1,
                    "numVlan": 1,
                    "phy": kwargs["port_name"] + "_1",
                    "vlanId": kwargs["vlan_network_host"],
                },
            },
        }
        return _dmf_config(nw_host)

    def common_TC_body(
        self, *args: None, **kwargs: Union[str, int, float]
    ) -> Dict[str, any]:
        """method to create template body as required by TC options"""
        self.save_name = kwargs["tc_name"] + "_" + kwargs["gateway"]
        port_info = kwargs["port_info"]

        kwargs["port_name"] = port_name = next(iter(port_info))
        # Setup SUT object based on epc type store as class var
        try:
            sut_mod = importlib.import_module(
                config.HIL.get("TEST_LIB_PATH") + "." + kwargs["epc"],
            )
            self.SutObject = getattr(sut_mod, kwargs["epc"])(kwargs["gateway"])
        except ModuleNotFoundError:
            from TS.agw import SUT_AGW

            self.SutObject = SUT_AGW(kwargs["gateway"])
        except Exception as e:
            logging.error(f"Error while importing SUT class {e}")
            sys.exit(1)

        if kwargs.get("ipv6", False):
            kwargs["port_name_v6"] = kwargs["port_name"] + "v6"
            # v6 ports have a special designation in Spirent
            kwargs["nodal_network"] = nodal_network = ip_network(
                port_info[port_name]["nodal_v6"],
            )
            kwargs["num_of_ip_nodal"] = 512
            kwargs.update({"gateway": kwargs.get("gateway") + "-v6"})

        else:
            kwargs["nodal_network"] = nodal_network = ip_network(
                port_info[port_name]["nodal"],
            )
            kwargs["num_of_ip_nodal"] = nodal_network.num_addresses
            kwargs["port_name_v6"] = kwargs["port_name"]

        kwargs["start_ip_for_control_enb"] = str(nodal_network[4])
        kwargs["network_host"] = network_host = ip_network(
            port_info[port_name]["network_host"],
        )
        kwargs["network_host_v6"] = network_host_v6 = ip_network(
            port_info[port_name]["network_host_v6"],
        )
        kwargs["num_of_ip_network_host"] = network_host.num_addresses
        kwargs["num_of_ip_network_host_v6"] = 512
        kwargs["start_ip_for_network_host_v6"] = str(network_host_v6[4])
        kwargs["cidr_mask_network_host_v6"] = "/" + str(network_host_v6.prefixlen)
        kwargs["start_ip_for_mob_control_enb"] = str(nodal_network[4 + kwargs["enbs"]])
        kwargs["start_ip_for_nodal3_control_enb"] = str(
            nodal_network[4 + kwargs["enbs"] + 4],
        )
        kwargs["cidr_mask_enb"] = "/" + str(nodal_network.prefixlen)
        kwargs["start_ip_for_network_host"] = str(network_host[4])
        kwargs["start_ip_for_network_host2"] = str(network_host[4 + 1])
        kwargs["cidr_mask_network_host"] = "/" + str(network_host.prefixlen)
        kwargs["start_mac_control_enb"] = port_info[port_name]["nodal_mac"]
        kwargs["start_mac_mob_control_enb"] = port_info[port_name]["nodal_mob_mac"]
        kwargs["start_mac_nodal3_control_enb"] = port_info[port_name]["nodal3_mac"]
        kwargs["start_mac_network_host"] = port_info[port_name]["network_host_mac"]
        kwargs["start_mac_network_host2"] = port_info[port_name]["network_host2_mac"]

        if kwargs.get(
            "nw_host3", False,
        ):  # If we do not need a nw_host nodal, use the same one as the first one
            kwargs["start_mac_network_host3"] = port_info[port_name][
                "network_host3_mac"
            ]
            kwargs["start_ip_for_network_host3"] = str(network_host[4 + 2])
        else:
            kwargs["start_ip_for_network_host3"] = kwargs["start_ip_for_network_host"]
            kwargs["start_mac_network_host3"] = kwargs["start_mac_network_host"]

        kwargs["vlan_enb"] = port_info[port_name]["nodal_vlan"]
        kwargs["vlan_network_host"] = port_info[port_name]["network_host_vlan"]

        kwargs["imsi_start"] = self.imsi_calc(
            kwargs["invalid_subs"], kwargs.get("mme_nodal_Imsi", "001011234560000"),
        )
        body = {
            "library": self.library_id,
            "name": self.save_name,
            "keywords": "DELETE_ME",
            "reservations": [
                {
                    "tsId": 2,
                    "phySubnets": [
                        {
                            "base": kwargs["start_ip_for_control_enb"],
                            "mask": kwargs["cidr_mask_enb"],
                            "name": str(kwargs["port_name_v6"]),
                            "numIps": kwargs["num_of_ip_nodal"] - 5,
                        },
                    ],
                },
            ],
            "tsGroups": [{"tsId": 2, "testCases": []}],
        }
        body["iterations"] = (
            kwargs.get("sequencer_mode", False) and 1 or kwargs["iterations"]
        )  # Num of iterations are managed by CommandSequencer
        if kwargs.get("steps", False):
            body["steps"] = self.get_tc_steps(**kwargs)
        if kwargs.get("nodal", False):
            body["tsGroups"][0]["testCases"].append(self.get_tc_mme_nodal(**kwargs))
        if kwargs.get("nw_host", False):
            body["tsGroups"][0]["testCases"].append(self.get_tc_nw_host(**kwargs))
        if kwargs.get("extra_phy_subnets", False):
            body["reservations"][0]["phySubnets"].append(
                self.get_tc_extra_phy_subnet(**kwargs),
            )
        if kwargs.get("mme_nodal_nw_host_dual_stack", False):
            body["reservations"][0]["phySubnets"].append(
                self.get_tc_ipv6_subnet(**kwargs),
            )
        # keep this at end as we change input dict
        if kwargs.get("availability_nodal", False):

            kwargs["start_ip_for_control_enb"] = kwargs["start_ip_for_mob_control_enb"]
            kwargs["start_mac_control_enb"] = kwargs["start_mac_mob_control_enb"]
            kwargs["start_ip_for_network_host"] = kwargs["start_ip_for_network_host2"]
            kwargs["start_mac_network_host"] = kwargs["start_mac_network_host2"]
            kwargs.update(
                {
                    # default values are for availability test
                    "HoldTime": kwargs.get("nodal2_HoldTime", 30),
                    "IdleTime": kwargs.get("nodal2_IdleTime", 30),
                    "total_subs": kwargs.get("nodal2_total_subs", 60),
                    "rate": kwargs.get("nodal2_rate", 1),
                    "mme_nodal_Imei": kwargs.get("nodal2_Imei", "11111111111111"),
                    "imsi_start": kwargs.get("nodal2_Imsi", "001011234560611"),
                    "dmf_conf": False,
                },
            )
            body["tsGroups"][0]["testCases"].append(self.get_tc_mme_nodal(**kwargs))
        if kwargs.get("availability_nw_host", False):
            body["tsGroups"][0]["testCases"].append(self.get_tc_nw_host(**kwargs))

        if kwargs.get("nodal3", False):

            kwargs["start_ip_for_control_enb"] = kwargs[
                "start_ip_for_nodal3_control_enb"
            ]
            kwargs["start_mac_control_enb"] = kwargs["start_mac_nodal3_control_enb"]
            kwargs["start_ip_for_network_host"] = kwargs["start_ip_for_network_host3"]
            kwargs["start_mac_network_host"] = kwargs["start_mac_network_host3"]

            mob_nodal_cnt = (
                kwargs["MobEnbControlAddr_enb"]
                if kwargs.get("MobEnbControlAddr_enb", False)
                else kwargs["enbs"]
            )

            kwargs["start_mac_for_mob_control_enb"] = str(
                EUI(int(EUI(port_info[port_name]["nodal_mob_mac"])) + mob_nodal_cnt),
            )  # Create offset off of initial nodal mob mac
            kwargs["start_ip_for_mob_control_enb"] = str(
                ip_address(kwargs["start_ip_for_nodal3_control_enb"])
                + kwargs["enbs"]
                + 1,
            )  # Create an IP offset from nodal IPs

            kwargs.update(
                {
                    "HoldTime": kwargs.get("nodal3_HoldTime", 30),
                    "IdleTime": kwargs.get("nodal3_IdleTime", 30),
                    "total_subs": kwargs.get("nodal3_total_subs", 60),
                    "rate": kwargs.get("nodal3_rate", 1),
                    "mme_nodal_Imei": kwargs.get("nodal3_Imei", "11111111111111"),
                    "imsi_start": kwargs.get("nodal3_Imsi", "001011234560611"),
                    "dmf_conf": kwargs.get("nodal3_dmf_conf", True),
                    "MobEnbControlAddr": kwargs.get("MobEnbControlAddr_nodal3", False),
                    "MobEnbControlAddr_enb": kwargs.get(
                        "MobEnbControlAddr_nodal3_enb"
                        if kwargs.get("MobEnbControlAddr_nodal3_enb")
                        else kwargs.get("enbs"),
                    ),
                },
            )

            body["tsGroups"][0]["testCases"].append(self.get_tc_mme_nodal(**kwargs))
        if kwargs.get("nw_host3", False):
            body["tsGroups"][0]["testCases"].append(self.get_tc_nw_host(**kwargs))
        return body

    def continue_func(self, **kwargs: Union[str, int, float]) -> Dict[str, any]:
        # -> None
        """
        This method sends as continue command to the sequencer at random times
        so as to simulate UEs being idle for a duration of time.

        IdleEntryTime is set to 4 seconds (default for Airspan)
        Consequently, UEs will remain idle anywhere from 5-10 seconds forcing
        UEs to transition to RRC_IDLE in each period of inactivity.

        """
        self.test_run_id = kwargs["id"]
        try:
            processes = []
            for tc in kwargs["tc_with_wait"]:
                kwargs["tc"] = tc
                kwargs["ts"] = 0
                proc = Process(target=self.send_continue, kwargs=kwargs)

                processes.append(proc)
                proc.start()

            kwargs["queue"] = Queue()
            data_measure = Process(target=self.check_and_wait_for_tc, kwargs=kwargs)
            processes.append(data_measure)
            data_measure.start()
            test_status = kwargs["queue"].get()

            for process in processes:
                process.join()

        except Exception as e:
            logging.error(f"Error: {e}")
            test_status = 0

        finally:
            return test_status

    def random_wait(self, profile_list: Dict[str, List[int]]) -> int:
        # Function to randomly calculate a wait time based on the constraints of the step.
        return random.randint(profile_list[0], profile_list[1])

    def get_sequencer_status(self, run_id: str) -> Dict[int, str]:
        # Function to find which testcase is sitting in the waiting state.
        # Called on-demand from individual processes.
        try:
            url = (
                "http://"
                + config.TAS.get("tas_ip")
                + ":8080/api/runningTests/"
                + run_id
            )
            r = get_ports.get(url, self.auth)
            if "COMPLETE" in r["testStateOrStep"]:
                logging.info("Test is complete; time to quit.")
                return {0: "COMPLETE", 2: "COMPLETE"}
            else:
                output = r["tsGroups"]
                values = []
                tc_step = {}
                for n in output[0]["testCases"]:
                    values.append(n["fullState"])
                for n in range(len(values)):
                    if "waiting" in values[n].lower():
                        step = values[n][
                            values[n].find(":") + 1: values[n].find("#")
                        ]  # Get's the actual step #.
                        tc_step[str(n)] = step

                        """output will be of the format {0:'8',2:'5'}
                        """
                    else:
                        logging.debug("TC is not in waiting state.")

                logging.debug(tc_step)
                return tc_step

        except Exception as e:
            logging.error(f"Error: {e}")
            sys.exit(1)

    def get_tc_result_file(self, tc_name: str, gateway: str) -> (bool, str):
        """helper func to get tc result file from TAS"""
        subprocess.call(
            [f"rm -f {config.TAS.get('test_report_path')}*{tc_name}_{gateway}.xls"],
            shell=True,
        )
        subprocess.call(
            [f"rm -f {config.TAS.get('test_report_path')}*{tc_name}_{gateway}.csv"],
            shell=True,
        )
        url = self.url + "/" + str(self.test_run_id)
        _res = get_ports.get(url, self.auth)
        output_file_list = _res["resultFilesList"]
        file_url = [name for name in output_file_list if "xls" in name][0]
        filename = file_url.split("/")[-1]
        filepath = config.TAS.get("test_report_path") + filename
        logging.info(f"Downloading file {filename} from the TAS server for data check")
        try:
            data = requests.get(file_url)
            with open(filepath, "wb") as f:
                f.write(data.content)
        except Exception as e:
            logging.error(f"Error while getting TC result file {e}")
            subprocess.call([f"rm -f {filepath}"], shell=True)
            return False, ""

        return True, filepath

    def send_continue(self, **kwargs: Union[str, int, float]) -> None:
        """
        Function spanawed for individual processes;
        each process maintains it's own lifecycle for the # of iterations.
        """

        num_of_iterations = (kwargs["iterations"] + 1) * 2

        for n in range(num_of_iterations):
            while True:
                tc_step = self.get_sequencer_status(kwargs["id"])
                if "COMPLETE" in tc_step.values():
                    return 1  # In the off chance that we get here when the TC is complete; continue with join() of the processes
                try:
                    tc = kwargs["tc"]
                    wait_time = self.random_wait(kwargs["wait_profile"][tc_step[tc]])
                    mode = kwargs["mode_lut"][tc_step[tc]]
                    logging.info(
                        f"Keeping UE/eNB group {tc} in {mode} for {wait_time} seconds.",
                    )
                    logging.info(
                        f"Starting Sequencer iteration {n} for UE/eNB group {tc}",
                    )
                    time.sleep(wait_time)  # Wait before "clicking" continue.
                    url = (
                        "http://"
                        + config.TAS.get("tas_ip")
                        + ":8080/api/runningTests/"
                        + kwargs["id"]
                        + "?action=sendTcCommand&ts=0"
                        + "&tc="
                        + tc
                        + "&cmd=continueSequencer"
                    )
                    continue_press = requests.request(
                        "POST", url, data=json.dumps(""), auth=self.auth,
                    )

                    if continue_press.status_code == 200:
                        logging.info(
                            f"UE group {tc} leaving {mode} state after {wait_time} seconds.",
                        )
                    time.sleep(
                        5,
                    )  # Buffer time so as to account for transport / processing delays.
                    break
                except KeyError:
                    logging.debug(f"Waiting for nodal {tc} to start sending traffic.")
                    time.sleep(5)
                    pass
                except KeyboardInterrupt:
                    url = (
                        "http://"
                        + config.TAS.get("tas_ip")
                        + ":8080/api/runningTests/"
                        + run_id
                        + "?action=abort"
                    )
                    abort = requests.request(
                        "POST", url, data=json.dumps(""), auth=self.auth,
                    )
                    if abort.status_code == 200:
                        logging.info("Python shell will exit now; test in TAS aborted")
                        sys.exit(0)
