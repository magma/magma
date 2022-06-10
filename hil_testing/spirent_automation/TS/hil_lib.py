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
import glob
import importlib
import json
import logging
import os
import pickle
import pprint
import subprocess
import sys
import time
from datetime import datetime
from typing import Dict, List, Union

import ansible_runner
import attr
import config
import jinja2
import pytz
import texttable
from anybadge import Badge
from pytz import timezone
from requests.auth import HTTPBasicAuth
from TS.agw import SUT_AGW
from TS.aws_lib import AWSbase
from TS.kpi_lib import cruncher
from TS.slack import SlackSender

LOGS_LOC = os.path.join(os.path.abspath(sys.path[0]), "logs")
sys.path.append(
    os.path.join(
        os.path.dirname(os.path.abspath(sys.path[0])), "Magma_Automations/scripts",
    ),
)
import base
import get_ports

# TODO we should create separate library for CI automation, this would not change
LIBRARY_NAME = "sms/AGW Scale"


@attr.s
class hil_lib(object):
    args = attr.ib()
    now_pt = attr.ib()

    def __attrs_post_init__(self):
        # get all var setup
        pass

    def store_data(self, data):
        """this function is to store result data to pass to CI
        data in dict form
        """
        dbfile = open("/tmp/test_res_pickle", "ab")
        pickle.dump(data, dbfile)
        dbfile.close()

    def hil_list(self):
        # trigger listing feature for all test cases
        print(f"\033[4m    {self.args.only_list:45} : \033[0m \n")
        for name in glob.glob("TC/*" + self.args.only_list + "*"):
            __tc = os.path.splitext(name)[0]
            tc_name = __tc.split("/")[1]
            print(f"    {tc_name:100} ")

    def run_cmd_worker(self, cmd: str) -> None:
        """
        This func can be used to send cmd to worker node to perform certain function
        """
        subprocess.call([cmd], shell=True)

    def hil_run(self):
        subprocess.call(["rm -f /tmp/test_res_pickle"], shell=True)
        if self.args.credentials_file:
            try:
                with open(self.args.credentials_file, "r") as f:
                    creds = json.load(f)
            except Exception as e:
                logging.error(f"Error reading credentials file {e}")
                sys.exit(1)
            if "MAGMA_PASSWORD" in creds:
                config.TAS["password"] = creds["TAS_PASSWORD"]
                config.AWS["secret_key"] = creds["AWS_SECRET_KEY"]
                config.MAGMA["password"] = creds["MAGMA_PASSWORD"]
                config.RDS["db_host"] = creds["RDS_HOST"]
                config.RDS["db_pass"] = creds["RDS_PASS"]
                config.SLACK["slack_webhook_path"] = creds["SLACK_WEBHOOK_PATH"]
        if config.TAS.get("password"):
            auth = HTTPBasicAuth(config.TAS.get("username"), config.TAS.get("password"))

        # get required class instance stored as class vars
        self.kpi_data = cruncher()
        self.a_client = AWSbase()

        # TODO create pool of gateway and use other gateway if one is busy
        r = get_ports.check_sut(self.args.gateway, config.TAS, auth)
        if r:
            logging.error(
                f"SUT {self.args.gateway} IS BUSY RUNNING OTHER TEST, SELECT OTHER FROM POOL {r}",
            )
            sys.exit(1)

        library_id = base.get_library_id(LIBRARY_NAME, config.TAS.get("tas_ip"), auth)
        test_summary = {}
        test_run_info = {}
        initial_magma_checks = {}
        try:
            sut_mod = importlib.import_module(
                config.HIL.get("TEST_LIB_PATH") + "." + self.args.epc,
            )
            logging.info(f"importing {self.args.epc} specific module ")
            sut = getattr(sut_mod, self.args.epc)(self.args.gateway)
        except ModuleNotFoundError:
            from TS.agw import SUT_AGW

            sut = SUT_AGW(self.args.gateway)

        except Exception as e:
            logging.error(f"Error while importing SUT class {e}")
            sys.exit(1)

        if self.args.upgrade or self.args.build != "latest":
            if not sut.upgrade(
                book="upgrade.yaml",
                magma_rel=self.args.rel,
                magma_build=self.args.build,
            ):
                logging.info(f" SUT SW upgrade Failed or no new build available")
                if not self.args.build_check:
                    sys.exit(0)

        cur_sw_ver = sut.get_build_info(role="get-sw-info")

        if self.args.reboot:
            if not sut.reboot(role="reboot"):
                logging.error(f" SUT SW reboot timeout - Exiting!!!")
                sys.exit(1)

        if cur_sw_ver:
            logging.info(
                f"Current running sw version on {self.args.gateway} is {cur_sw_ver} - Healthy",
            )
        else:
            logging.error(f"No Magma SW running on {self.args.gateway} exiting!!!")
            sys.exit(1)
        # get all test cases for mentioned test suite and run_test

        initial_magma_memory, initial_magma_pid = sut.sut_magma_checks(
            role="agw-health-checks",
        )
        run_only_tests = None
        run_test_list = []
        tests_set = sorted(glob.glob("TC/*" + self.args.test_suite + "*"))
        if self.args.only_run:
            run_only_tests = set(self.args.only_run)
            for name in tests_set:
                __tc = os.path.splitext(name)[0].replace("/", ".")
                if not any(_t in __tc for _t in run_only_tests):
                    continue
                run_test_list.append(__tc)
        else:
            run_test_list = [
                os.path.splitext(x)[0].replace("/", ".") for x in tests_set
            ]
        total_tests = len(run_test_list)
        for name in run_test_list:
            _from_time = round(time.time() * 1000)
            _tc = importlib.import_module(name)
            # TODO port ip assignment config in get_next_avail_port should be part of config
            # Reserve port per test.
            port_info = get_ports.get_next_avail_port(config.TAS, auth)
            logging.info(f"Reserving port {port_info} for test {name}")
            test_run_id, tc_result, test_metrics = getattr(_tc, "run_test")(
                gateway=self.args.gateway,
                port_info=port_info,
                library_id=library_id,
                auth=auth,
                test_suite=self.args.test_suite,
                epc=self.args.epc,
            )
            # Add SUT check here before running next test
            # get SUT info
            post_magma_memory, post_magma_pid = sut.sut_magma_checks(
                role="agw-health-checks",
            )

            proc_restarted_list = sut.detect_failed_procedure(
                initial_magma_pid,
                post_magma_pid,
                cur_sw_ver,
                initial_magma_memory,
                post_magma_memory,
                config.MAGMA.get("memory_delta_failed_pct"),
            )
            # store per test pre and post mem utilization per process
            test_metrics["pre_test_mem"] = initial_magma_memory
            test_metrics["post_test_mem"] = post_magma_memory

            total_tests -= 1
            time.sleep(60) if total_tests else time.sleep(0)
            # TODO we can multitask this with multiple sut and executors
            _to_time = round(time.time() * 1000)
            _name = name.split(".")[1]
            test_summary[_name] = tc_result
            # get detailed test result file for each test
            check, csv_filename = self.check_and_get_file(_name)
            logging.info(f"Got test result csv file {csv_filename} for test {_name}")
            if check:
                test_run_info[_name] = {"test_results_file": csv_filename}

                self.a_client.upload_file(
                    f"{config.TAS.get('test_report_path')}{csv_filename}",
                    ts=self.args.test_suite.lower(),
                    c_type="application/octet-stream",
                )
            else:
                test_run_info[_name] = {"test_results_file": None}

            test_run_info[_name].update(
                {
                    "tc_run_id": test_run_id,
                    "tc_run_sut": self.args.gateway,
                    "tc_start_time": _from_time,
                    "tc_stop_time": _to_time,
                    "failed_proc": proc_restarted_list,
                    "test_metrics": json.dumps(test_metrics),
                    "test_suite": self.args.test_suite.lower(),
                },
            )

            for p_names in post_magma_pid.keys():
                if not "core_files" in p_names:
                    """
                    initial core_files will always be "NO_CORE"; detect_failed_procedure resets
                    after every TC is run.
                    """
                    initial_magma_pid[p_names] = post_magma_pid[p_names]

        # Get log from sut for all required processes
        test_suite = self.args.test_suite.lower()
        log_result = sut.get_logs(
            role="get-logs",
            start=self.now_pt,
            processes=config.MAGMA.get("processes", []),
            test_suite=test_suite,
        )
        if self.args.output_text:
            self.print_formatted_text(test_summary, test_run_info)

        # upload test summary to aws s3
        if config.AWS.get("secret_key") and self.args.output_s3:
            run_log = LOGS_LOC + "/" + self.now_pt + "_" + test_suite + "_magma.log"
            test_results = self.adjudicator(test_summary, test_run_info)
            res_file = self.create_test_summary_html(
                cur_sw_ver,
                run_log,
                test_results,
                self.now_pt,
            )

            self.store_data({"verdict": self.verdict(test_results), "report": res_file})

            if test_suite == "sanity":
                sanity_pass_badge = config.MAGMA.get("sanity_pass_badge", None)
                sanity_res_badge = config.MAGMA.get("sanity_res_badge", None)
                if self.verdict(test_results):
                    logging.info(f"Generating Badge for sanity pass")
                    self.create_badge(
                        cur_sw_ver, sanity_pass_badge, "HIL AGW tests stable", "green",
                    )
                    self.create_badge(
                        "passing", sanity_res_badge, "HIL AGW tests", "green",
                    )
                    self.a_client.upload_file(
                        sanity_pass_badge, ts=test_suite, c_type="image/svg+xml",
                    )
                    self.a_client.upload_file(
                        sanity_res_badge, ts=test_suite, c_type="image/svg+xml",
                    )
                else:
                    logging.info(f"Generating Badge for sanity failed")
                    self.create_badge(
                        "failing", sanity_res_badge, "HIL AGW tests", "red",
                    )
                    self.a_client.upload_file(
                        sanity_res_badge, ts=test_suite, c_type="image/svg+xml",
                    )

                self.run_cmd_worker("rm -f *.svg")

            for test, res in test_results.items():
                _test_result = True if res["result"] == "PASS" else False
                if test_suite == "availability" and not _test_result:
                    continue
                self.a_client.db_connect_insert(
                    release=cur_sw_ver.split("-")[0],
                    testname=test,
                    testresult=_test_result,
                    runtime=self.now_pt,
                    build=cur_sw_ver,
                    testsuite=test_suite,
                    testnote=str(res["fail_procs"]),
                    sut=self.args.gateway,
                    SystemAvailability=res["System_Availability"],
                    TestKPI=test_run_info[test]["test_metrics"],
                )
            # update DB with test execution results
            try:
                self.a_client.upload_file(run_log, c_type="application/octet-stream")
                self.a_client.upload_file(res_file, ts=test_suite)
                latest_file = "/tmp/" + test_suite + str(os.getpid()) + "_latest.html"
                subprocess.call(["cp", res_file, latest_file])
                self.a_client.upload_file(latest_file, ts=test_suite)
            except Exception as e:
                logging.error(f" failed to push log/report file to S3 Error: {e}")
            # push report to slack hil-test channel
            slack_config = config.SLACK
            if slack_config.get("slack_webhook_path"):
                slacker = SlackSender()
                slacker.send_report(
                    test_results,
                    cur_sw_ver,
                    test_suite,
                    link=slack_config.get("dashboard"),
                )
        self.cleanup()
        logging.info(f"Completed Test Execution check test summary page")

    def check_and_get_file(self, tc_name: str) -> (bool, str):
        """
        helper func to check and convert to xls to csv
        """
        try:
            filepath = glob.glob(
                f"{config.TAS.get('test_report_path')}*{tc_name}_{self.args.gateway}.xls",
            )[0]
            print(
                glob.glob(
                    f"{config.TAS.get('test_report_path')}*{tc_name}_{self.args.gateway}.xls",
                ),
            )
            if os.path.exists(filepath):
                check, csvfile = self.kpi_data.xls_to_csv(filepath)
                if check and os.path.exists(
                    config.TAS.get("test_report_path") + csvfile,
                ):
                    return True, csvfile
            return False, ""
        except Exception as e:
            logging.error(f"could not get xls/csv file for test {e}")
            return False, ""

    def cleanup(self):
        """
        func to cleanup
        """
        pass

    def print_formatted_text(
        self,
        test_summary: Dict[str, Dict[str, any]],
        test_run_info: Dict[str, Dict[str, any]] = None,
    ):
        """
        This function would print formatted result table on screen
        """
        table = texttable.Texttable(max_width=100)
        table.set_cols_align(["l", "c", "l"])
        rows = [["TEST_NAME", "TEST_RESULT", "TEST_RESULTS_ANALYTICS"]]
        results = self.adjudicator(test_summary, test_run_info)
        for test, res in results.items():
            rows.append([test, res["result"], res["fail_procs"]])
        table.add_rows(rows, header=True)
        print(table.draw())

    def verdict(self, results: Dict[str, Dict[str, any]]) -> Dict[str, bool]:
        """This Functions is to call out perticular test suite passed or not"""
        test_res_list = []
        for k, v in results.items():
            test_res_list.append(True if v.get("result", "FAIL") == "PASS" else False)
        print(test_res_list)
        return all(test_res_list)

    def create_badge(
        self, ver: str, badge_file: str, badge_label: str, badge_color: str,
    ) -> None:
        """This functions will create badge as per run results"""
        badge = Badge(label=badge_label, value=ver, default_color=badge_color)
        badge.write_badge(badge_file, overwrite=True)

    def adjudicator(
        self,
        test_summary: Dict[str, Dict[str, any]],
        tc_run_info: Dict[str, Dict[str, any]] = None,
    ) -> Dict[str, Dict[str, any]]:
        """
        This Function would update results based on result summary received
        """
        results = {}
        excluded_cat_list = ["Timeout", "AVAILABILITY"]

        for k, v in test_summary.items():
            fail_list = []
            try:
                for cat, proc in v.items():
                    if cat not in excluded_cat_list:
                        for n in proc:
                            if proc[n] == False:
                                fail_list.append(n)
            except Exception as e:
                logging.error(f" failed to run test suite Error: {e} on object {v}")
                return {}
            if v["Timeout"]:
                fmt_val = "TIMEOUT"
            elif len(fail_list) > 0:
                fmt_val = "FAIL"
                fail_list.insert(0, "Failed Procedure -> ")
            else:
                fmt_val = "PASS"

            if tc_run_info:
                System_Availability = v.get("AVAILABILITY", {}).get(
                    "System_Availability", None,
                )
                results[k] = {
                    "result": fmt_val,
                    "fail_procs": fail_list
                    + tc_run_info[k]["failed_proc"]
                    + (
                        ["System Availability - " + str(System_Availability)]
                        if System_Availability
                        else []
                    ),
                    "tc_run_id": tc_run_info[k]["tc_run_id"],
                    "tc_run_sut": tc_run_info[k]["tc_run_sut"],
                    "tc_start_time": tc_run_info[k]["tc_start_time"],
                    "tc_stop_time": tc_run_info[k]["tc_stop_time"],
                    "test_results_file": tc_run_info[k]["test_results_file"],
                    "test_suite": tc_run_info[k]["test_suite"],
                    "System_Availability": System_Availability or 0.0,
                }
            else:
                results[k] = {"result": fmt_val, "fail_procs": fail_list}
        return results

    def create_test_summary_html(
        self,
        sw: str,
        log_file: str,
        test_results: Dict[str, Dict[str, any]],
        timestr: str,
    ) -> str:
        """
        This Function would render webpage for test summary for s3 web portal
        """
        env = jinja2.Environment(
            loader=jinja2.FileSystemLoader(searchpath="templates/"),
        )
        template = env.get_template("test_case_results.html")
        _out = template.render(
            sw=sw, result=test_results, log_file=os.path.basename(log_file),
        )
        # File name to push to S3 includes build version + timestamp
        # to match dashborad so that we can create link
        res_filename = "/tmp/test_result_" + sw + "_" + timestr + ".html"
        with open(res_filename, "w") as hf:
            hf.write(_out)
        return res_filename
