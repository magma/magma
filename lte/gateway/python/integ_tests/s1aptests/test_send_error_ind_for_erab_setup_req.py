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

import time
import unittest

import s1ap_types
from integ_tests.s1aptests import s1ap_wrapper
from integ_tests.s1aptests.s1ap_utils import SpgwUtil


class TestSendErrorIndForErabSetupReq(unittest.TestCase):
    """Test sending of error indication for erab setup request"""

    def setUp(self):
        """Initialize"""
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()
        self._spgw_util = SpgwUtil()

    def tearDown(self):
        """Cleanup"""
        self._s1ap_wrapper.cleanup()

    def test_send_error_ind_for_erab_setup_req(self):
        """attach/detach + erabsetup req + error ind with a single UE"""
        num_ues = 1
        detach_type = [
            s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value,
            s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
        ]
        wait_for_s1 = [True, False]
        self._s1ap_wrapper.configUEDevice(num_ues)

        for i in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            print(
                "********************** Running End to End attach for ",
                "UE id ",
                req.ue_id,
            )
            # attach
            attach_acc = self._s1ap_wrapper._s1_util.attach(
                req.ue_id,
                s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
                s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
                s1ap_types.ueAttachAccept_t,
            )
            # Wait on EMM Information from MME
            self._s1ap_wrapper._s1_util.receive_emm_info()

            # Send indication to drop erab setup req
            print("*** Sending indication to tfw to drop erab setup req***")
            drop_erab_setup_req = s1ap_types.DropErabSetupReq_t()
            drop_erab_setup_req.ue_Id = req.ue_id
            drop_erab_setup_req.flag = 1
            self._s1ap_wrapper._s1_util.issue_cmd(
                s1ap_types.tfwCmd.DROP_ERAB_SETUP_REQ, drop_erab_setup_req,
            )

            print("Sleeping for 5 seconds")
            time.sleep(5)
            print(
                "********************** Adding dedicated bearer to IMSI",
                "".join([str(i) for i in req.imsi]),
            )

            # Create default flow list
            flow_list = self._spgw_util.create_default_ipv4_flows()
            self._spgw_util.create_bearer(
                "IMSI" + "".join([str(i) for i in req.imsi]),
                attach_acc.esmInfo.epsBearerId,
                flow_list,
            )

            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_ACT_DED_BER_REQ.value,
            )
            act_ded_ber_ctxt_req = response.cast(
                s1ap_types.UeActDedBearCtxtReq_t,
            )
            self._s1ap_wrapper.sendActDedicatedBearerAccept(
                req.ue_id, act_ded_ber_ctxt_req.bearerId,
            )
            # Send error indication for erab setup req
            error_ind = s1ap_types.fwNbErrIndMsg_t()
            # isUeAssoc flag to include optional MME_UE_S1AP_ID and
            # eNB_UE_S1AP_ID
            error_ind.isUeAssoc = True
            error_ind.ue_Id = req.ue_id
            error_ind.cause.pres = True
            # Radio network causeType = 0
            error_ind.cause.causeType = 0
            # causeVal- Unknown-pair-ue-s1ap-id
            error_ind.cause.causeVal = 15
            print("*** Sending error indication ***")
            self._s1ap_wrapper._s1_util.issue_cmd(
                s1ap_types.tfwCmd.ENB_ERR_IND_MSG, error_ind,
            )
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_CTX_REL_IND.value,
            )
            print(
                "************************* Received UE_CTX_REL_IND for UE id ",
                req.ue_id,
            )
            print(
                "********************** Running UE detach for UE id ",
                req.ue_id,
            )
            # Now detach the UE
            self._s1ap_wrapper.s1_util.detach(
                req.ue_id, detach_type[i], wait_for_s1[i],
            )


if __name__ == "__main__":
    unittest.main()
