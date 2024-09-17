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
from TS.aws_lib import AWSbase

ANSIBLE_LOC = os.path.join(
    os.path.dirname(os.path.abspath(sys.path[0])), "Magma_Ansible",
)

LOGS_LOC = os.path.join(os.path.abspath(sys.path[0]), "logs")


@attr.s
class SUT_AGW(object):
    sut = attr.ib()

    def __attrs_post_init__(self):
        agw_pass = config.MAGMA.get("password")
        self.e_vars = {
            "ansible_ssh_pass": agw_pass,
            "ansible_sudo_pass": agw_pass,
            "ansible_user": config.MAGMA.get("username"),
            "ansible_become": "yes",
            "ANSIBLE_HOST_KEY_CHECKING": False,
        }

    def upgrade(
        self, book: str = None, magma_rel: str = "ci", magma_build: str = "latest",
    ) -> bool:
        """this is dummy mathod, EPC specific class mathod should be used"""
        logging.info(f"Please use EPC required mathod to upgrade")
        return True

    def _run_ansible_role(self, **kwargs: Union[str, int, float]) -> Union[bool, any]:
        """Worker func to run role"""

        try:
            r = ansible_runner.run(
                private_data_dir=ANSIBLE_LOC,
                limit=self.sut,
                role=kwargs["role"],
                rotate_artifacts=1,
                directory_isolation_base_path="/tmp/runner",
                extravars=kwargs["extra_vars"],
                cmdline=kwargs.get("cmdline", "--tags all"),
            )
            subprocess.call(["rm", "-f", ANSIBLE_LOC + "/project/main.json"])
            subprocess.call(["rm", "-f", ANSIBLE_LOC + "/env/extravars"])
        except Exception as e:
            logging.error(f"Ansible role run got error - {e}")
            # clean up
            subprocess.call(["rm", "-f", ANSIBLE_LOC + "/project/main.json"])
            subprocess.call(["rm", "-f", ANSIBLE_LOC + "/env/extravars"])
            return False
        if r.status == "successful" and r.rc == 0:
            return r
        else:
            return False

    def pcap_collection(self, **kwargs: Union[str, int, float]) -> Union[bool, str]:
        """This method will start the pcap collection on the SUT"""
        kwargs["cap_vars"].update(self.e_vars)
        res = self._run_ansible_role(
            role=kwargs["role"],
            extra_vars=kwargs["cap_vars"],
            cmdline=kwargs["cmdline"],
        )
        if res:
            logging.info(
                f"Successfully executed {kwargs['cmdline']} in role {kwargs['role']}",
            )

        else:
            logging.info(
                f"Execution failed while executing {kwargs['cmdline']} in role {kwargs['role']}",
            )

        return res

    def get_build_info(self, role: str = None) -> Union[bool, str]:
        """Use EPC specific class mathod"""
        logging.warn(f"Please use EPC specific mathod to retrive SW version")
        return True

    def sut_magma_state_check(self, role: str = None) -> Dict[str, int]:
        """
        Use EPC specific health check
        """
        logging.info(f"EPC specific health check needs to be used")
        return config.MAGMA_AGW.get("UE_STATE", {})

    def sut_check(self, sut_res: any = None, key: str = None) -> Dict[str, any]:
        """EPC specific checks mathod should be used to retrive memory and pid info"""
        mem_dict = {"epc": 1}
        pid_dict = {"epc": 1}

        return mem_dict, pid_dict

    def sut_file_cleanup(self, file_path: str) -> None:
        """EPC specific cleanup mathod should be used"""
        pass
        # self._run_ansible_role(role="test-cleanup", extra_vars=self.e_vars)

    def sut_magma_checks(self, role: str = None) -> Union[bool, Dict[str, any]]:
        """EPC specific health check roles/mathods should be used"""
        pass
        # res = self._run_ansible_role(role=role, extra_vars=self.e_vars)
        mem_dict, pid_dict = self.sut_check()
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
        This is pure magma specific, EPC AGW specific mathod should be used for non-magma
        """
        pass
        return []

    def upload_file_s3(self, file_path: str, aws_obj: any):
        """upload file to s3 with aws object"""

        aws_obj.upload_file(file_path, ts="cores", c_type="application/octet-stream")

    def reboot(self, role: str = None) -> bool:
        """This method would reboot the SUT"""

        res = self._run_ansible_role(role=role, extra_vars=self.e_vars)
        if res:
            logging.info(f"SUT restarted Successfully")
            return True
        else:
            return False

    def get_sut_file(self, file_path: str, dest_location: str, dest_name: str) -> bool:
        """helper function to retrieve file from SUT"""
        copy_flag = True
        extra_vars = {
            "dest_location": dest_location + "/" + dest_name,
            "file_path": file_path,
        }
        extra_vars.update(self.e_vars)
        res = self._run_ansible_role(role="get-file", extra_vars=extra_vars)
        if res:
            logging.info(f" SUT {file_path} retrieved Successfully")
        else:
            logging.info(f" SUT {file_path} retrieval failed")
            copy_flag = False

        return copy_flag

    def get_logs(self, **kwargs: Union[str, int, float]) -> bool:
        """EPC specific Func to get log"""
        return True
