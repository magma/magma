"""
Copyright 2021 The Magma Authors.

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
import shlex
import subprocess
from typing import List  # noqa

from .tc_ops import TcOpsBase

LOG = logging.getLogger('pipelined.qos.tc_cmd')


# this code can run in either a docker container(CWAG) or as a native
# python process(AG). When we are running as a root there is no need for
# using sudo. (related to task T63499189 where tc commands failed since
# sudo wasn't available in the docker container)
def argSplit(cmd: str) -> List[str]:
    args = [] if os.geteuid() == 0 else ["sudo"]
    args.extend(shlex.split(cmd))
    return args


def run_cmd(cmd_list, show_error=True) -> int:
    err = 0
    for cmd in cmd_list:
        LOG.debug("running %s", cmd)
        try:
            args = argSplit(cmd)
            subprocess.check_call(args)
        except subprocess.CalledProcessError as e:
            err = -1
            if show_error:
                LOG.error("error: %s running: %s", str(e.returncode), cmd)
    return err


class TcOpsCmd(TcOpsBase):
    def __init__(self):
        LOG.info("initialized")

    def create_htb(self, iface: str, qid: str, max_bw: str, rate:str,
                   parent_qid: str = None) -> int:
        tc_cmd = "tc class add dev {intf} parent {parent_qid} "
        tc_cmd += "classid 1:{qid} htb rate {rate} ceil {maxbw} prio 2"
        tc_cmd = tc_cmd.format(intf=iface, parent_qid=parent_qid,
                               qid=qid, rate=rate, maxbw=max_bw)

        return run_cmd([tc_cmd], True)

    def del_htb(self, iface: str, qid: str) -> int:
        del_cmd = "tc class del dev {intf} classid 1:{qid}".format(intf=iface, qid=qid)
        return run_cmd([del_cmd], True)

    def create_filter(self, iface: str, mark: str, qid: str, proto: int = 3) -> int:
        filter_cmd = "tc filter add dev {intf} protocol ip parent 1: prio 1 "
        filter_cmd += "handle {mark} fw flowid 1:{qid}"
        filter_cmd = filter_cmd.format(intf=iface, mark=mark, qid=qid)
        return run_cmd([filter_cmd], True)

    def del_filter(self, iface: str, mark: str, qid: str, proto: int = 3) -> int:
        filter_cmd = "tc filter del dev {intf} protocol ip parent 1: prio 1 "
        filter_cmd += "handle {mark} fw flowid 1:{qid}"
        filter_cmd = filter_cmd.format(intf=iface, mark=mark, qid=qid)
        return run_cmd([filter_cmd], True)
