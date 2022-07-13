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
import pprint
import sys
import time
from typing import Dict, Tuple, Union

import attr
import config
from TS.ts_lib import TSBase

tc = {}
##########################################
# DO NOT CHANGE THIS, ADD NEW TEST INSTEAD

tc["total_subs"] = 3
tc["invalid_subs"] = 0
tc["enbs"] = 1
tc["rate"] = 3
tc["template"] = "TEMPLATE_DATA_MIX_200M"
tc["dmf"] = "UDP_DL_300M_Per_UE"
# tc["dmf"] = "HTTP_DL_80M_Per_UE_82"
tc["iterations"] = 1
tc["tc_dpruntime"] = 60  # seconds
tc["tc_name"] = os.path.basename(__file__).split(".")[0]
tc["steps"] = True
tc["nodal"] = True
tc["nw_host"] = True
tc["extra_phy_subnets"] = True
tc["dmf_conf"] = True
tc["NetworkHostAddrRemote"] = True
tc["EnbUserAddr"] = False
tc["categories"] = ["S1-AP", "ESM", "EMM", "Data_Traffic", "UE_STATE_CHECK"]
# TODO read APN AMBR using NMS API
tc["expected_qos_flow_datarate"] = (
    tc["total_subs"] * config.MAGMA["nms_apn_ambr"] * config.MAGMA["pass_percentile"]
)
tc["high_bw_limit"] = True
##########################################


@attr.s
class TCBase(TSBase):
    pass


def run_test(
    *args: str, **kwargs: Union[str, int, float]
) -> Tuple[int, Dict[str, any]]:
    ts = TCBase(tc["template"], kwargs["library_id"], kwargs["auth"])
    body = ts.common_TC_body(**tc, **kwargs)
    res = ts.save_test(body)
    try:
        if res.status_code == 200:
            run_info = ts.run_test()
            # For Automation no need to get continuous update
            logging.info(
                f"Started execution for {tc['tc_name']} with run id as {run_info['id']}",
            )
            test_status = ts.check_and_wait_for_tc(id=run_info["id"], **tc)
            logging.info(f"{tc['tc_name']} with run id {run_info['id']} is completed")
            if test_status:
                results, test_metrics = ts.check_test_summary(
                    id=run_info["id"],
                    **test_status,
                    **tc,
                    **kwargs,
                )
                return run_info["id"], results, test_metrics
            else:
                return 0, {"Timeout": False, "Generic": {"Spirent_system": False}}, {}
        else:
            logging.error(f"Error while saving test config to TAS...{res.text}")
    except Exception as e:
        logging.error(f"Error working with TAS API run info: {run_info} error: {e}")
    return (
        0,
        {"Timeout": False, "Generic": {"Spirent_system": False}},
        {},
    )  # all fail case
