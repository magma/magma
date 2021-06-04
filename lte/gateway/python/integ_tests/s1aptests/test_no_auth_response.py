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

import unittest

import orc8r.protos.metricsd_pb2 as metricsd
import s1ap_types
import s1ap_wrapper
from python.integ_tests.common.service303_utils import (
    MetricValue,
    verify_gateway_metrics,
)


class TestNoAuthResponse(unittest.TestCase):

    TEST_METRICS = [
        MetricValue(
            service="mme",
            name=str(metricsd.ue_attach),
            labels={
                str(metricsd.result): "failure",
                str(metricsd.cause): "no_response_for_auth_request",
            },
            value=1,
        ),
        MetricValue(
            service="mme",
            name=str(metricsd.ue_attach),
            labels={str(metricsd.action): "attach_accept_sent"},
            value=0,
        ),
        MetricValue(
            service="mme",
            name=str(metricsd.ue_detach),
            labels={str(metricsd.cause): "implicit_detach"},
            value=1,
        ),
        MetricValue(
            service="mme",
            name=str(metricsd.nas_auth_rsp_timer_expired),
            labels={},
            value=1,
        ),
        MetricValue(
            service="mme",
            name=str(metricsd.spgw_create_session),
            labels={str(metricsd.result): "success"},
            value=0,
        ),
    ]

    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()
        self.gateway_services = self._s1ap_wrapper.get_gateway_services_util()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    @verify_gateway_metrics
    def test_no_auth_response(self):
        # Ground work.
        self._s1ap_wrapper.configIpBlock()
        self._s1ap_wrapper.configUEDevice(1)

        req = self._s1ap_wrapper.ue_req
        print(
            "************************* Running attach no auth response \
            timer expiry test",
        )

        attach_req = s1ap_types.ueAttachRequest_t()
        attach_req.ue_Id = req.ue_id
        sec_ctxt = s1ap_types.TFW_CREATE_NEW_SECURITY_CONTEXT
        id_type = s1ap_types.TFW_MID_TYPE_IMSI
        eps_type = s1ap_types.TFW_EPS_ATTACH_TYPE_EPS_ATTACH
        attach_req.mIdType = id_type
        attach_req.epsAttachType = eps_type
        attach_req.useOldSecCtxt = sec_ctxt

        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ATTACH_REQUEST, attach_req,
        )
        # Wait for timer expiry 5 times, until context is released
        for i in range(5):
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_AUTH_REQ_IND.value,
            )
            print("************************* Timeout", i + 1)

        print("************************* Timeouts complete")
        # Attach Reject
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_ATTACH_REJECT_IND.value,
        )

        # Context release
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_CTX_REL_IND.value,
        )
        print("************************* Context released")


if __name__ == "__main__":
    unittest.main()
