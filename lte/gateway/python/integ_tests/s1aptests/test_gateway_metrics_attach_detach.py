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


class TestGatewayMetricsAttachDetach(unittest.TestCase):

    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()
        self._gateway_service = self._s1ap_wrapper.get_gateway_services_util()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def _getMetricValueGivenLabel(self, metric_name, label_values):
        service = self._gateway_service.get_mme_service_util()
        return service.get_metric_value(
            metric_name,
            label_values,
            default=0,
        )

    def test_gateway_metrics_attach_detach(self):
        """ Basic gateway metrics with attach/detach for a single UE """

        label_values_success = {"result": "success"}
        mme_new_association_before = self._getMetricValueGivenLabel(
            "mme_new_association",
            label_values_success,
        )

        num_ues = 2
        detach_type = [
            s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value,
            s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
        ]
        wait_for_s1 = [True, False]
        self._s1ap_wrapper.configUEDevice(num_ues)

        label_values_ue_attach_result = \
            {"result": "attach_proc_successful"}

        for i in range(num_ues):
            v_ue_attach = self._getMetricValueGivenLabel(
                "ue_attach",
                label_values_ue_attach_result,
            )
            v_ue_detach = self._getMetricValueGivenLabel(
                "ue_detach",
                label_values_success,
            )
            v_spgw_create_session = self._getMetricValueGivenLabel(
                "spgw_create_session",
                label_values_success,
            )
            v_spgw_delete_session = self._getMetricValueGivenLabel(
                "spgw_delete_session",
                label_values_success,
            )

            req = self._s1ap_wrapper.ue_req
            print(
                "************************* Running End to End attach for ",
                "UE id ", req.ue_id,
            )
            # Now actually complete the attach
            self._s1ap_wrapper._s1_util.attach(
                req.ue_id, s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
                s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
                s1ap_types.ueAttachAccept_t,
            )

            # Wait on EMM Information from MME
            self._s1ap_wrapper._s1_util.receive_emm_info()

            # waits until the metrics get updated
            time.sleep(0.5)
            val = self._getMetricValueGivenLabel(
                "ue_attach",
                label_values_ue_attach_result,
            )
            assert (val == v_ue_attach + 1)

            val = self._getMetricValueGivenLabel(
                "spgw_create_session",
                label_values_success,
            )
            assert (val == v_spgw_create_session + 1)

            val = self._getMetricValueGivenLabel(
                "ue_detach",
                label_values_success,
            )
            assert (val == v_ue_detach)

            val = self._getMetricValueGivenLabel(
                "spgw_delete_session",
                label_values_success,
            )
            assert (val == v_spgw_delete_session)

            print(
                "************************* Running UE detach for UE id ",
                req.ue_id,
            )
            # Now detach the UE
            self._s1ap_wrapper.s1_util.detach(
                req.ue_id, detach_type[i], wait_for_s1[i],
            )

            # waits so that metrics have time to be updated
            time.sleep(0.5)
            val = self._getMetricValueGivenLabel(
                "ue_detach",
                label_values_success,
            )
            assert(val == v_ue_detach + 1)

            val = self._getMetricValueGivenLabel(
                "spgw_delete_session",
                label_values_success,
            )
            assert (val == v_spgw_delete_session + 1)

        mme_new_association = self._getMetricValueGivenLabel(
            "mme_new_association",
            label_values_success,
        )
        assert (mme_new_association == mme_new_association_before + 1)


if __name__ == "__main__":
    unittest.main()
