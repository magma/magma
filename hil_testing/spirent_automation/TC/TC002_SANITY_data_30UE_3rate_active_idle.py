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
from jinja2 import Template
from TS.ts_lib import TSBase

tc = {}
##########################################
# DO NOT CHANGE THIS, ADD NEW TEST INSTEAD

tc["total_subs"] = 15
tc["invalid_subs"] = 0
tc["enbs"] = 6  # This gets replicated for the second nodal also; so 12 eNBs total.
tc["rate"] = tc["nodal3_rate"] = 1.5
tc["template"] = "TEMPLATE_DATA_HTTP_SESSION_LOADING_SEQ"
tc["dmf"] = "HTTP_DL_500K_Per_UE"
tc[
    "iterations"
] = 3  # This is the sequencer mode iterations; n iterations == n+1 executions.
tc["tc_name"] = os.path.basename(__file__).split(".")[0]
tc["steps"] = False
tc["tc_dpruntime"] = 400

# Sequencer Mode configuration
tc["sequencer_mode"] = True
tc["mode_lut"] = {
    "7": "Active mode",
    "10": "Idle Mode",
}  # `7` and `10` refer to the steps in command mode sequencer. The meaning of those steps is here purely for log readability!
tc["wait_profile"] = {
    "7": config.SPIRENT_SEQUENCER_DELAY_PROFILE.get("medium"),  # Active mode
    "10": config.SPIRENT_SEQUENCER_DELAY_PROFILE.get("short"),  # Idle mode
}
tc["tc_with_wait"] = [
    "0",
    "2",
]  # This is used to spawn independent processes for each nodal. One `tc` here refers to the nodal # in Spirent Landslide.

# First nodal
tc["nodal"] = True
tc["mme_nodal_Imei"] = "21111111111111"
tc["mme_nodal_Imsi"] = "001011234560000"
tc["MobEnbControlAddr"] = True
tc[
    "MobEnbControlAddr_enb"
] = 1  # This is necessary for sequencer mode to work; it caters for potential mobility (not used)

# Second nodal
tc["nodal3"] = True
tc["nodal3_HoldTime"] = 1200
tc["nodal3_Imei"] = "31111111111111"
tc["nodal3_Imsi"] = "001011234560300"
tc["nodal3_total_subs"] = 15
tc["MobEnbControlAddr_nodal3"] = True
tc[
    "MobEnbControlAddr_nodal3_enb"
] = 1  # This is necessary for sequencer mode to work; it caters for potential mobility (not used)

tc["nw_host"] = True
tc["extra_phy_subnets"] = True
tc["dmf_conf"] = False
tc["NetworkHostAddrRemote"] = True
tc["EnbUserAddr"] = False
tc["categories"] = [
    "ESM",
    "EMM",
    "UE_STATE_CHECK",
    "active_idle",
]
##########################################


@attr.s
class TCBase(TSBase):
    pass


def run_test(
    *args: str, **kwargs: Union[str, int, float]
) -> Tuple[int, Dict[str, any]]:
    ts = TCBase(
        tc["template"],
        kwargs["library_id"],
        kwargs["auth"],
    )
    body = ts.common_TC_body(**tc, **kwargs)
    res = ts.save_test(body)
    try:
        if res.status_code == 200:
            run_info = ts.run_test()  # Run the test here.
            # For Automation no need to get continuous update
            logging.info(
                f"Started execution for {tc['tc_name']} with run id as {run_info['id']}",
            )

            test_status = ts.continue_func(
                id=run_info["id"], **tc, **kwargs,
            )  # function call to manage active/idle transitions.

            # test_status = ts.check_and_wait_for_tc(id=run_info["id"], **tc)
            # logging.info(f"{tc['tc_name']} with run id {run_info['id']} is completed")
            print(f"Test_Status is: {test_status} ")  # Delete
            if test_status:
                results, test_metrics = ts.check_test_summary(
                    id=run_info["id"],
                    **test_status,
                    **tc,
                    **kwargs,
                )
                return run_info["id"], results, test_metrics
            else:
                return (
                    0,
                    {
                        "Timeout": False,
                        "Generic": {"Spirent_system": False},
                    },
                    {},
                )
        else:
            logging.error(f"Error while saving test config to TAS...{res.text}")
    except Exception as e:
        logging.error(f"Error working with TAS API run info: {run_info} error: {e}")
    return (
        0,
        {
            "Timeout": False,
            "Generic": {"Spirent_system": False},
        },
        {},
    )  # all fail case
