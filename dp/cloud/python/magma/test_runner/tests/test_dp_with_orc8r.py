"""
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""
import operator
from datetime import datetime, timezone
from http import HTTPStatus
from time import sleep
from typing import Any, Dict, List, Optional

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

DP_HTTP_PREFIX = 'magma/v1/dp'
NETWORK = 'some_network'
SERIAL_NUMBER = "some_serial_number"
SOME_FCC_ID = "some_fcc_id"
OTHER_FCC_ID = "other_fcc_id"
USER_ID = "some_user_id"
SAS = 'SAS'
DP = 'DP'
DATETIME_WAY_BACK = '2010-03-28T09:13:25.407877399+00:00'


@pytest.mark.orc8r
class DomainProxyOrc8rTestCase(DBTestCase):
    def setUp(self):
        super().setUp()
        DBInitializer(SessionManager(self.engine)).initialize()
        grpc_channel = grpc.insecure_channel(
            f"{config.GRPC_SERVICE}:{config.GRPC_PORT}",
        )
        self.dp_client = DPServiceStub(grpc_channel)
        # Indexing from previous test may not have been finished yet
        sleep(5)
        self._delete_dp_elasticsearch_indices()

    def tearDown(self):
        super().tearDown()
        self._delete_dp_elasticsearch_indices()

    # retrying is needed because of a possible deadlock
    # with cc locking request table
    @retry(stop_max_attempt_number=5, wait_fixed=100)
    def drop_all(self):
        super().drop_all()

    def test_cbsd_sas_flow(self):
        cbsd_id = self.given_cbsd_provisioned()
        # Giving ElasticSearch time to index logs
        sleep(5)

        logs = self.when_logs_are_fetched(self.get_current_sas_filters())
        self.then_logs_are(logs, self.get_sas_provision_messages())

        filters = self.get_filters_for_request_type('heartbeat')
        self.then_message_is_eventually_sent(filters, keep_alive=True)

        self.delete_cbsd(cbsd_id)

    def test_sas_flow_restarted_for_updated_cbsd(self):
        cbsd_id = self.given_cbsd_provisioned()

        self.when_cbsd_is_updated(cbsd_id, OTHER_FCC_ID)

        filters = self.get_filters_for_request_type('deregistration')
        self.then_message_is_eventually_sent(filters, keep_alive=True)

        self.then_state_is_eventually(self.get_state_with_grant())

        cbsd = self.when_cbsd_is_fetched()
        self.then_cbsd_is(cbsd, self.get_cbsd_data_with_grant(OTHER_FCC_ID))

        self.delete_cbsd(cbsd_id)

    def test_activity_status(self):
        cbsd_id = self.given_cbsd_provisioned()

        cbsd = self.when_cbsd_is_fetched()
        self.then_cbsd_is(cbsd, self.get_cbsd_data_with_grant())

        self.when_cbsd_is_inactive()
        cbsd = self.when_cbsd_is_fetched()
        self.then_cbsd_is(cbsd, self.get_registered_cbsd_data())

        self.delete_cbsd(cbsd_id)

    def test_fetching_logs_with_custom_filters(self):
        self.given_cbsd_provisioned()
        # Giving ElasticSearch time to index logs
        sleep(5)

        sas_to_dp_end_date_only = {
            'serial_number': SERIAL_NUMBER,
            'from': SAS,
            'to': DP,
            'end': self.now(),
        }
        sas_to_dp_begin_date_only = {
            'serial_number': SERIAL_NUMBER,
            'from': SAS,
            'to': DP,
            'begin': DATETIME_WAY_BACK,
        }
        sas_to_dp_end_date_too_early = {
            'serial_number': SERIAL_NUMBER,
            'from': SAS,
            'to': DP,
            'end': DATETIME_WAY_BACK,
        }
        dp_to_sas = {
            'serial_number': SERIAL_NUMBER,
            'from': DP,
            'to': SAS,
        }
        dp_to_sas_incorrect_serial_number = {
            'serial_number': 'incorrect_serial_number',
            'from': DP,
            'to': SAS,
        }
        sas_to_dp_with_limit = {
            'limit': '100',
            'from': SAS,
            'to': DP,
        }
        sas_to_dp_with_limit_and_too_large_offset = {
            'limit': '100',
            'offset': '100',
            'from': SAS,
            'to': DP,
        }
        scenarios = [
            (sas_to_dp_end_date_only, operator.eq, 4, HTTPStatus.OK),
            (sas_to_dp_begin_date_only, operator.eq, 4, HTTPStatus.OK),
            (sas_to_dp_end_date_too_early, operator.eq, 0, HTTPStatus.OK),
            (dp_to_sas, operator.gt, 0, HTTPStatus.OK),
            (dp_to_sas_incorrect_serial_number, operator.eq, 0, HTTPStatus.OK),
            (sas_to_dp_with_limit, operator.eq, 4, HTTPStatus.OK),
            (sas_to_dp_with_limit_and_too_large_offset, operator.eq, 0, HTTPStatus.OK),
        ]

        for params in scenarios:
            self._verify_logs_count(params)

    def given_cbsd_provisioned(self) -> int:
        self.when_cbsd_is_created()
        cbsd = self.when_cbsd_is_fetched()
        self.then_cbsd_is(cbsd, self.get_unregistered_cbsd_data())

        state = self.when_cbsd_asks_for_state()
        self.then_state_is(state, self.get_empty_state())

        self.then_state_is_eventually(self.get_state_with_grant())

        cbsd = self.when_cbsd_is_fetched()
        self.then_cbsd_is(cbsd, self.get_cbsd_data_with_grant())

        return cbsd['id']

    def delete_cbsd(self, cbsd_id: int):
        filters = self.get_filters_for_request_type('deregistration')

        self.when_cbsd_is_deleted(cbsd_id)
        self.then_cbsd_is_deleted()

        state = self.when_cbsd_asks_for_state()
        self.then_state_is(state, self.get_empty_state())

        self.then_message_is_eventually_sent(filters, keep_alive=False)

    def when_cbsd_is_created(self):
        r = send_request_to_backend(
            'post', 'cbsds', json=self.get_cbsd_post_data(),
        )
        self.assertEqual(r.status_code, HTTPStatus.CREATED)

    def when_cbsd_is_fetched(self) -> Dict[str, Any]:
        r = send_request_to_backend('get', 'cbsds')
        self.assertEqual(r.status_code, HTTPStatus.OK)
        data = r.json()
        self.assertEqual(data.get('total_count'), 1)
        cbsds = data.get('cbsds', [])
        self.assertEqual(len(cbsds), 1)
        return cbsds[0]

    def when_logs_are_fetched(self, params: Dict[str, Any]) -> Dict[str, Any]:
        r = send_request_to_backend('get', 'logs', params=params)
        self.assertEqual(r.status_code, HTTPStatus.OK)
        data = r.json()
        return data

    def when_cbsd_is_deleted(self, cbsd_id: int):
        r = send_request_to_backend('delete', f'cbsds/{cbsd_id}')
        self.assertEqual(r.status_code, HTTPStatus.NO_CONTENT)

    def when_cbsd_is_updated(self, cbsd_id: int, fcc_id: str):
        r = send_request_to_backend(
            'put', f'cbsds/{cbsd_id}', json=self.get_cbsd_post_data(fcc_id=fcc_id),
        )
        self.assertEqual(r.status_code, HTTPStatus.NO_CONTENT)

    def when_cbsd_asks_for_state(self) -> CBSDStateResult:
        return self.dp_client.GetCBSDState(self.get_cbsd_request())

    @staticmethod
    def when_cbsd_is_inactive():
        inactivity = 3
        polling = 1
        delta = 3
        total_wait_time = inactivity + 2 * polling + delta
        sleep(total_wait_time)

    def then_cbsd_is(self, actual: Dict[str, Any], expected: Dict[str, Any]):
        actual = actual.copy()
        del actual['id']
        grant = actual.get('grant')
        if grant:
            del grant['grant_expire_time']
            del grant['transmit_expire_time']
        self.assertEqual(actual, expected)

    def then_cbsd_is_deleted(self):
        r = send_request_to_backend('get', 'cbsds')
        self.assertEqual(r.status_code, HTTPStatus.OK)
        data = r.json()
        self.assertFalse(data.get('total_count', True))

    def then_state_is(self, actual: CBSDStateResult, expected: CBSDStateResult):
        self.assertEqual(actual, expected)

    @retry(stop_max_attempt_number=30, wait_fixed=1000)
    def then_state_is_eventually(self, expected):
        actual = self.when_cbsd_asks_for_state()
        self.then_state_is(actual, expected)

    def then_logs_are(self, actual: Dict[str, Any], expected: List[str]):
        actual = [x['type'] for x in actual['logs']]
        self.assertEqual(actual, expected)

    @retry(stop_max_attempt_number=60, wait_fixed=1000)
    def then_message_is_eventually_sent(self, filters: Dict[str, Any], keep_alive):
        if keep_alive:
            self.when_cbsd_asks_for_state()
        logs = self.when_logs_are_fetched(filters)
        self.assertEqual(logs["total_count"], 1)

    @staticmethod
    def get_cbsd_request() -> CBSDRequest:
        return CBSDRequest(serial_number=SERIAL_NUMBER)

    @staticmethod
    def get_empty_state() -> CBSDStateResult:
        return CBSDStateResult(radio_enabled=False)

    @staticmethod
    def get_state_with_grant() -> CBSDStateResult:
        return CBSDStateResult(
            radio_enabled=True,
            channel=LteChannel(
                low_frequency_hz=3620_000_000,
                high_frequency_hz=3630_000_000,
                max_eirp_dbm_mhz=28.0,
            ),
        )

    @staticmethod
    def get_cbsd_post_data(fcc_id: str = SOME_FCC_ID) -> Dict[str, Any]:
        return {
            "capabilities": {
                "antenna_gain": 15,
                "max_power": 20,
                "min_power": 0,
                "number_of_antennas": 2,
            },
            "fcc_id": fcc_id,
            "serial_number": SERIAL_NUMBER,
            "user_id": USER_ID,
        }

    def get_unregistered_cbsd_data(self) -> Dict[str, Any]:
        data = self.get_cbsd_post_data()
        data.update({
            'is_active': False,
            'state': 'unregistered',
        })
        return data

    def get_registered_cbsd_data(self, fcc_id: str = SOME_FCC_ID) -> Dict[str, Any]:
        data = self.get_cbsd_post_data(fcc_id)
        data.update({
            'cbsd_id': f'{fcc_id}/{SERIAL_NUMBER}',
            'is_active': False,
            'state': 'registered',
        })
        return data

    def get_cbsd_data_with_grant(self, fcc_id: str = SOME_FCC_ID) -> Dict[str, Any]:
        data = self.get_registered_cbsd_data(fcc_id)
        data.update({
            'is_active': True,
            'grant': {
                'bandwidth_mhz': 10,
                'frequency_mhz': 3625,
                'max_eirp': 28,
                'state': 'authorized',
            },
        })
        return data

    def get_current_sas_filters(self) -> Dict[str, Any]:
        return {
            'serial_number': SERIAL_NUMBER,
            'from': SAS,
            'to': DP,
            'end': self.now(),
        }

    def get_filters_for_request_type(self, request_type: str) -> Dict[str, Any]:
        return {
            'serial_number': SERIAL_NUMBER,
            'type': f'{request_type}Response',
            'begin': self.now(),
        }

    @staticmethod
    def get_sas_provision_messages() -> List[str]:
        names = ['heartbeat', 'grant', 'spectrumInquiry', 'registration']
        return [f'{x}Response' for x in names]

    def now(self):
        return datetime.now(timezone.utc).isoformat()

    def _delete_dp_elasticsearch_indices(self):
        requests.delete(f"{config.ELASTICSEARCH_URL}/{config.ELASTICSEARCH_INDEX}*")

    def _verify_logs_count(self, params):
        using_filters, _operator, expected_count, expected_status = params
        logs = self.when_logs_are_fetched(using_filters)

        comparison = _operator(len(logs["logs"]), expected_count)
        self.assertTrue(comparison)


def send_request_to_backend(
    method: str, url_suffix: str, params: Optional[Dict[str, Any]] = None,
    json: Optional[Dict[str, Any]] = None,
) -> requests.Response:
    return requests.request(
        method,
        f'{config.HTTP_SERVER}/{DP_HTTP_PREFIX}/{NETWORK}/{url_suffix}',
        cert=(config.DP_CERT_PATH, config.DP_SSL_KEY_PATH),
        verify=False,  # noqa: S501
        params=params,
        json=json,
    )
