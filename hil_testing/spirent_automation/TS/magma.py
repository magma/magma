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
import logging
import os
import re
import subprocess
import sys
from typing import Dict, List, Tuple, Union

import ansible_runner
import attr
import config
from TS.agw import SUT_AGW
from TS.aws_lib import AWSbase

ANSIBLE_LOC = os.path.join(
    os.path.dirname(os.path.abspath(sys.path[0])), "Magma_Ansible",
)

LOGS_LOC = os.path.join(os.path.abspath(sys.path[0]), "logs")


@attr.s
class magma(SUT_AGW):
    sut = attr.ib()

    def upgrade(
        self, book: str = None, magma_rel: str = "ci", magma_build: str = "latest",
    ) -> bool:
        """this methond would upgrade given sut with latest software build for given Magma Release"""
        self.e_vars["magma_rel"] = magma_rel
        if magma_build:
            try:
                sctpd_build = config.MAGMA.get("magma_sctpd_pkg_name").join(
                    magma_build.rsplit(config.MAGMA.get("magma_pkg_name"), 1),
                )
                if sctpd_build:
                    self.e_vars["sctpd_build"] = sctpd_build
                    self.e_vars["Magma_build"] = magma_build
                else:
                    logging.error(f"Magma Build format is not proper")
                    return False
            except Exception as e:
                logging.error(f"Error while creating sctpd build url {e}")
                return False

        r = ansible_runner.run(
            private_data_dir=ANSIBLE_LOC,
            playbook=book,
            limit=self.sut,
            rotate_artifacts=1,
            directory_isolation_base_path="/tmp/runner",
            extravars=self.e_vars,
        )

        if r.status == "successful" and r.rc == 0:
            upgrade_done = r.get_fact_cache(self.sut)["need_upgrade"]
            logging.info(f" SUT upgrade done. Was it needed? - {upgrade_done}")
            return upgrade_done if upgrade_done else False
        else:
            logging.error(f" SUT upgrade Failed with error - {r.status}")
            return False

    def get_build_info(self, role: str = None) -> Union[bool, str]:
        """This method would retrive running magma sw build information"""
        res = self._run_ansible_role(role=role, extra_vars=self.e_vars)
        if res:
            logging.info(f"MAGMA SW retrieved Successfully")
            _sw_ver = res.get_fact_cache(self.sut)["agw_magma_build"].strip()
            return _sw_ver if _sw_ver else False
        else:
            logging.info(f"MAGMA sw retrieval failed")
            return False

    def sut_magma_state_check(self, role: str = None) -> Dict[str, int]:
        """This method would fetch and check UE state across magma agw processes
        on output it should give true or false if state is expected or not
        """
        current_state = {}
        res = self._run_ansible_role(role=role, extra_vars=self.e_vars)

        for n in res.get_fact_cache(self.sut)["svc_state"]:
            current_state[f"{n['item'].lower()}_state"] = int(n["stdout"])

        for n in res.get_fact_cache(self.sut)["table_state"]:
            current_state[f"table{str(n['item'])}_flows"] = int(n["stdout"])

        current_state["mobility_state"] = int(
            res.get_fact_cache(self.sut)["mobility_state"].strip(),
        )
        return current_state

    def sut_check(self, sut_res: any = None, key: str = None) -> Dict[str, any]:
        """helper function to do sut checks as per config definition"""
        mem_dict = {}
        pid_dict = {}

        if sut_res:
            logging.info(f"MAGMA check in progress")
            sanitized = sut_res.get_fact_cache(self.sut)["magma_pid_mem_values"].split(
                "\n",
            )
            output = [n for n in sanitized if n != ""]
            output_dict = {}
            pattern = pattern = "magma@(.*.)\.service"
            for n in range(len(output)):
                if "magma" in output[n]:
                    key = re.search(pattern, output[n]).group(1)
                    mem_dict[key] = output[n - 1]
                    pid_dict[key] = output[n - 2]
            pid_dict["core_files"] = sut_res.get_fact_cache(self.sut)[
                "magma_core_files"
            ]
        return mem_dict, pid_dict

    def sut_file_cleanup(self, file_path: str) -> None:
        """helper function to clear files from SUT as per test needs"""
        self._run_ansible_role(role="test-cleanup", extra_vars=self.e_vars)

    def sut_magma_checks(self, role: str = None) -> Union[bool, Dict[str, any]]:
        """This method would check for SUT magma software is healthy or not"""
        res = self._run_ansible_role(role=role, extra_vars=self.e_vars)
        mem_dict, pid_dict = self.sut_check(res)

        return mem_dict, pid_dict if bool(mem_dict and pid_dict) else False

    def detect_failed_procedure(
        self,
        initial_process_dict: Dict[str, any],
        post_test_dict: Dict[str, any],
        sw_ver: str,
        initial_magma_memory: Dict[str, any],
        post_test_magma_memory: Dict[str, any],
        memory_delta_failed_pct: float,
    ) -> List[str]:
        """
        this method is used to detect and return which process failed
        and if mem usage exceeded a pre-defined threshold after the test
        """
        logging.info(f"MAGMA process pid dict before test {initial_process_dict}")
        logging.info(f"MAGMA process pid dict after test {post_test_dict}")
        logging.info(f"MAGMA process mem dict before test {initial_magma_memory}")
        logging.info(f"MAGMA process mem dict after test {post_test_magma_memory}")
        failed_procedures = [
            k
            for k in initial_process_dict
            if k in post_test_dict and initial_process_dict[k] != post_test_dict[k]
        ]
        for k in initial_magma_memory:
            try:
                pct_change = round(
                    (
                        (int(post_test_magma_memory[k]) - int(initial_magma_memory[k]))
                        / int(initial_magma_memory[k])
                    )
                    * 100,
                    2,
                )
                post_mem_mb = round(
                    int(post_test_magma_memory[k]) / 1048576, 2,
                )  # outputs in Mb
                logging.info(
                    f"MAGMA process {k} post test memory usage was {post_mem_mb}M",
                )
                if pct_change > memory_delta_failed_pct and int(
                    post_mem_mb > 275,
                ):  # Essentially check that mem usage is greater than 275M; Most services except mme, sessiond, and pipelined have limits greater than 300m, while the others are limited to 299
                    failed_procedures.append(k + "_memory_" + str(pct_change) + "pct")
                    logging.warning(
                        f"MAGMA process {k} exceeded memory by {pct_change}%",
                    )
            except Exception as e:
                logging.error(f"MAGMA Proceess {k} memory check failed with error {e}")
        if "core_files" in failed_procedures:
            # upload file to aws
            a_client = AWSbase()
            for core_f in post_test_dict["core_files"]:
                core_with_ver = (
                    os.path.splitext(core_f)[0]
                    + "-"
                    + sw_ver
                    + os.path.splitext(core_f)[1]
                )
                failed_procedures.append(core_with_ver)
                dest_name = os.path.basename(core_with_ver)
                self.get_sut_file(
                    file_path=core_f, dest_location=LOGS_LOC, dest_name=dest_name,
                )
                self.upload_file_s3(
                    file_path=LOGS_LOC + "/" + dest_name, aws_obj=a_client,
                )
            self.sut_file_cleanup(file_path=config.MAGMA.get("core_path"))

        return failed_procedures

    def get_logs(self, **kwargs: Union[str, int, float]) -> bool:
        """Function to get all specified process logs"""

        log_flag = True
        for process in kwargs["processes"]:
            extra_vars = {
                "start": kwargs["start"],
                "process": process,
                "process_log_loc": LOGS_LOC
                + "/"
                + kwargs["start"]
                + "_"
                + kwargs["test_suite"]
                + "_magma.log",
            }
            extra_vars.update(self.e_vars)
            extra_vars.update({"ansible_become": "no"})
            res = self._run_ansible_role(role=kwargs["role"], extra_vars=extra_vars)
            if res:
                logging.info(f"MAGMA {process} logs retrieved Successfully")
            else:
                logging.info(f"MAGMA {process} logs retrieval failed")
                log_flag = False

        return log_flag
