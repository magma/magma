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
from collections import defaultdict
from time import sleep

import grpc
import pytest
import requests
from dp.protos.enodebd_dp_pb2 import CBSDRequest, CBSDStateResult, LteChannel
from dp.protos.enodebd_dp_pb2_grpc import DPServiceStub
from magma.db_service.db_initialize import DBInitializer
from magma.db_service.session_manager import SessionManager
from magma.db_service.tests.db_testcase import DBTestCase
from magma.test_runner.config import TestConfig
from retrying import retry

config = TestConfig()

SERIAL_NUMBER = "some_serial_number"
FCC_ID = "some_fcc_id"
USER_ID = "some_user_id"


@pytest.mark.local
class ActiveModeControllerTestCase(DBTestCase):
    @classmethod
    def setUpClass(cls):
        wait_for_elastic_to_start()

    def setUp(self):
        super().setUp()
        grpc_channel = grpc.insecure_channel(
            f"{config.GRPC_SERVICE}:{config.GRPC_PORT}",
        )
        self.dp_client = DPServiceStub(grpc_channel)
        DBInitializer(SessionManager(self.engine)).initialize()
        # Indexing from previous test may not have been finished yet
        sleep(5)
        self._delete_dp_elasticsearch_indices()

    # retrying is needed because of a possible deadlock
    # with cc locking request table
    @retry(stop_max_attempt_number=5, wait_fixed=100)
    def drop_all(self):
        super().drop_all()

    def tearDown(self):
        super().tearDown()
        self._delete_dp_elasticsearch_indices()

    def test_provision_cbsd_in_sas_requested_by_dp_client(self):
        self.given_cbsd_provisioned()

    def test_logs_are_written_to_elasticsearch(self):
        self.given_cbsd_provisioned()
        # Giving ElasticSearch time to index logs
        sleep(5)
        self.then_elasticsearch_contains_logs()

    def test_grant_relinquished_after_inactivity(self):
        self.given_cbsd_provisioned()
        self.when_cbsd_is_inactive()
        self.then_cbsd_has_no_grants_in_sas(self.dp_client)

    def test_last_used_max_eirp_stays_the_same_after_inactivity(self):
        self.given_cbsd_provisioned()
        self.when_cbsd_is_inactive()
        self.then_cbsd_has_no_grants_in_sas(self.dp_client)
        self.given_cbsd_provisioned()

    def given_cbsd_provisioned(self):
        self.given_cbsd_with_transmission_parameters()
        self.dp_client.GetCBSDState(self._build_cbsd_request())

        self.then_cbsd_is_eventually_provisioned_in_sas(self.dp_client)

    def then_elasticsearch_contains_logs(self):
        actual = requests.get(f"{config.ELASTICSEARCH_URL}/dp*/_search?size=100").json()
        log_field_names = [
            "log_from",
            "log_to",
            "log_name",
            "log_message",
            "cbsd_serial_number",
            "network_id",
            "fcc_id",
        ]
        actual_log_types = defaultdict(int)
        logs = actual["hits"]["hits"]
        for log in logs:
            actual_log_types[log["_source"]["log_name"]] += 1
            for fn in log_field_names:
                self.assertIn(fn, log["_source"].keys())

        self.assertEqual(1, actual_log_types["CBSDRegisterRequest"])
        self.assertEqual(1, actual_log_types["CBSDRegisterResponse"])
        self.assertEqual(1, actual_log_types["registrationRequest"])
        self.assertEqual(1, actual_log_types["registrationResponse"])
        self.assertEqual(1, actual_log_types["spectrumInquiryRequest"])
        self.assertEqual(1, actual_log_types["spectrumInquiryResponse"])
        self.assertEqual(1, actual_log_types["spectrumInquiryResponse"])
        self.assertEqual(1, actual_log_types["grantRequest"])
        self.assertEqual(1, actual_log_types["grantResponse"])
        self.assertEqual(1, actual_log_types["heartbeatRequest"])
        # The number of GetCBSDStateRequest and heartbeatResponse may differ between tests, so only asserting they have been logged
        self.assertGreater(actual_log_types["heartbeatResponse"], 0)
        self.assertGreater(actual_log_types["GetCBSDStateRequest"], 0)
        self.assertGreater(actual_log_types["GetCBSDStateResponse"], 0)

    # TODO change this when some API for domain proxy is introduced
    def given_cbsd_with_transmission_parameters(self):
        state = self.dp_client.CBSDRegister(
            self._build_cbsd_request(), wait_for_ready=True,
        )
        self.session.commit()
        self.assertEqual(self._build_empty_state_result(), state)

    @staticmethod
    def when_cbsd_is_inactive():
        inactivity = 3
        polling = 1
        delta = 3  # TODO investigate if such high delta is needed
        total_wait_time = inactivity + 2 * polling + delta
        sleep(total_wait_time)

    @retry(stop_max_attempt_number=30, wait_fixed=1000)
    def then_cbsd_is_eventually_provisioned_in_sas(self, dp_client: DPServiceStub):
        state = dp_client.GetCBSDState(self._build_cbsd_request())
        self.assertEqual(self._build_get_state_result(), state)

    def then_cbsd_has_no_grants_in_sas(self, dp_client: DPServiceStub):
        state = dp_client.GetCBSDState(self._build_cbsd_request())
        self.assertEqual(self._build_empty_state_result(), state)

    @staticmethod
    def _build_cbsd_request() -> CBSDRequest:
        return CBSDRequest(
            user_id=USER_ID,
            fcc_id=FCC_ID,
            serial_number=SERIAL_NUMBER,
            min_power=0,
            max_power=20,
            antenna_gain=15,
            number_of_ports=2,
        )

    @staticmethod
    def _build_get_state_result() -> CBSDStateResult:
        return CBSDStateResult(
            radio_enabled=True,
            channel=LteChannel(
                low_frequency_hz=3620_000_000,
                high_frequency_hz=3630_000_000,
                max_eirp_dbm_mhz=28.0,
            ),
        )

    @staticmethod
    def _build_empty_state_result() -> CBSDStateResult:
        return CBSDStateResult(radio_enabled=False)

    def _delete_dp_elasticsearch_indices(self):
        requests.delete(f"{config.ELASTICSEARCH_URL}/{config.ELASTICSEARCH_INDEX}*")


@retry(stop_max_attempt_number=30, wait_fixed=1000)
def wait_for_elastic_to_start() -> None:
    requests.get(f'{config.ELASTICSEARCH_URL}/_status')
