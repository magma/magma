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
from __future__ import annotations

import operator
from contextlib import contextmanager
from datetime import datetime, timezone
from http import HTTPStatus
from threading import Event, Thread
from time import sleep
from typing import Any, Dict, List, Optional
from uuid import uuid4

import pytest
import requests
from dp.protos.enodebd_dp_pb2 import CBSDRequest, CBSDStateResult, LteChannel
from magma.test_runner.config import TestConfig
from magma.test_runner.tests.integration_testcase import (
    DomainProxyIntegrationTestCase,
)
from retrying import retry

config = TestConfig()

DP_HTTP_PREFIX = 'magma/v1/dp'
NETWORK = 'some_network'
SOME_FCC_ID = "some_fcc_id"
OTHER_FCC_ID = "other_fcc_id"
USER_ID = "some_user_id"
SAS = 'SAS'
DP = 'DP'
DATETIME_WAY_BACK = '2010-03-28T09:13:25.407877399+00:00'


@pytest.mark.orc8r
class DomainProxyOrc8rTestCase(DomainProxyIntegrationTestCase):
    def setUp(self) -> None:
        self.serial_number = self._testMethodName + '_' + str(uuid4())

    def test_cbsd_sas_flow(self):
        builder = CbsdAPIDataBuilder().with_serial_number(self.serial_number)
        cbsd_id = self.given_cbsd_provisioned(builder)

        with self.while_cbsd_is_active():
            when_elastic_indexes_data()

            logs = self.when_logs_are_fetched(get_current_sas_filters(self.serial_number))
            self.then_logs_are(logs, self.get_sas_provision_messages())

            filters = get_filters_for_request_type('heartbeat', self.serial_number)
            self.then_message_is_eventually_sent(filters)

        self.delete_cbsd(cbsd_id)

    def test_cbsd_unregistered_when_requested_by_desired_state(self):
        builder = CbsdAPIDataBuilder().with_serial_number(self.serial_number)
        cbsd_id = self.given_cbsd_provisioned(builder)

        with self.while_cbsd_is_active():
            filters = get_filters_for_request_type('deregistration', self.serial_number)

            builder = builder.with_desired_state('unregistered')
            self.when_cbsd_is_updated(cbsd_id, builder.build_post_data())

            # TODO maybe asking for state (cbsd api instead of log api) would be better
            self.then_message_is_eventually_sent(filters)

    def test_sas_flow_restarted_when_user_requested_deregistration(self):
        builder = CbsdAPIDataBuilder().with_serial_number(self.serial_number)
        cbsd_id = self.given_cbsd_provisioned(builder)

        with self.while_cbsd_is_active():
            filters = get_filters_for_request_type('deregistration', self.serial_number)

            self.when_cbsd_is_deregistered(cbsd_id)

            self.then_message_is_eventually_sent(filters)

            self.then_state_is_eventually(builder.build_grant_state_data())

    def test_sas_flow_restarted_for_updated_cbsd(self):
        builder = CbsdAPIDataBuilder().with_serial_number(self.serial_number)
        cbsd_id = self.given_cbsd_provisioned(builder)

        with self.while_cbsd_is_active():
            builder = builder.with_fcc_id(OTHER_FCC_ID)
            self.when_cbsd_is_updated(cbsd_id, builder.build_post_data())

            filters = get_filters_for_request_type('deregistration', self.serial_number)
            self.then_message_is_eventually_sent(filters)

            self.then_state_is_eventually(builder.build_grant_state_data())

            cbsd = self.when_cbsd_is_fetched(builder.serial_number)
            self.then_cbsd_is(cbsd, builder.build_registered_active_data())

        self.delete_cbsd(cbsd_id)

    def test_activity_status(self):
        builder = CbsdAPIDataBuilder().with_serial_number(self.serial_number)
        cbsd_id = self.given_cbsd_provisioned(builder)

        cbsd = self.when_cbsd_is_fetched(builder.serial_number)
        self.then_cbsd_is(cbsd, builder.build_registered_active_data())

        self.when_cbsd_is_inactive()
        cbsd = self.when_cbsd_is_fetched(builder.serial_number)
        self.then_cbsd_is(cbsd, builder.build_registered_inactive_data())

    def test_frequency_preferences(self):
        builder = CbsdAPIDataBuilder(). \
            with_serial_number(self.serial_number). \
            with_frequency_preferences(5, [3625]). \
            with_expected_grant(5, 3625, 31)
        cbsd_id = self.given_cbsd_provisioned(builder)

        self.delete_cbsd(cbsd_id)

    def test_creating_cbsd_with_the_same_unique_fields_returns_409(self):
        builder = CbsdAPIDataBuilder().with_serial_number(self.serial_number)

        self.when_cbsd_is_created(builder.build_post_data())
        self.when_cbsd_is_created(builder.build_post_data(), expected_status=HTTPStatus.CONFLICT)

    def test_updating_cbsd_returns_409_when_setting_existing_serial_num(self):
        builder = CbsdAPIDataBuilder()

        cbsd1_serial = self.serial_number + "_foo"
        cbsd2_serial = self.serial_number + "_bar"
        cbsd1_payload = builder.with_serial_number(cbsd1_serial).build_post_data()
        cbsd2_payload = builder.with_serial_number(cbsd2_serial).build_post_data()
        self.when_cbsd_is_created(cbsd1_payload)
        self.when_cbsd_is_created(cbsd2_payload)
        cbsd2 = self.when_cbsd_is_fetched(serial_number=cbsd2_serial)
        self.when_cbsd_is_updated(
            cbsd_id=cbsd2.get("id"),
            data=cbsd1_payload,
            expected_status=HTTPStatus.CONFLICT,
        )

    def test_fetch_cbsds_filtered_by_serial_number(self):
        cbsd1_serial = self.serial_number + "_foo"
        cbsd2_serial = self.serial_number + "_bar"

        builder1 = CbsdAPIDataBuilder().with_serial_number(cbsd1_serial)
        builder2 = CbsdAPIDataBuilder().with_serial_number(cbsd2_serial)

        self.when_cbsd_is_created(builder1.build_post_data())
        self.when_cbsd_is_created(builder2.build_post_data())

        cbsd1 = self.when_cbsd_is_fetched(serial_number=cbsd1_serial)
        cbsd2 = self.when_cbsd_is_fetched(serial_number=cbsd2_serial)

        self.then_cbsd_is(cbsd1, builder1.build_unregistered_data())
        self.then_cbsd_is(cbsd2, builder2.build_unregistered_data())

    def test_fetching_logs_with_custom_filters(self):
        builder = CbsdAPIDataBuilder().with_serial_number(self.serial_number)

        sas_to_dp_end_date_only = {
            'serial_number': self.serial_number,
            'from': SAS,
            'to': DP,
            'end': now(),
        }
        sas_to_dp_begin_date_only = {
            'serial_number': self.serial_number,
            'from': SAS,
            'to': DP,
            'begin': DATETIME_WAY_BACK,
        }
        sas_to_dp_end_date_too_early = {
            'serial_number': self.serial_number,
            'from': SAS,
            'to': DP,
            'end': DATETIME_WAY_BACK,
        }
        dp_to_sas = {
            'serial_number': self.serial_number,
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
            (sas_to_dp_end_date_only, operator.eq, 0),
            (sas_to_dp_begin_date_only, operator.gt, 3),
            (sas_to_dp_end_date_too_early, operator.eq, 0),
            (dp_to_sas, operator.gt, 0),
            (dp_to_sas_incorrect_serial_number, operator.eq, 0),
            (sas_to_dp_with_limit, operator.gt, 3),
            (sas_to_dp_with_limit_and_too_large_offset, operator.eq, 0),
        ]

        self.given_cbsd_provisioned(builder)
        with self.while_cbsd_is_active():
            when_elastic_indexes_data()

            for params in scenarios:
                self._verify_logs_count(params)

    def given_cbsd_provisioned(self, builder: CbsdAPIDataBuilder) -> int:
        self.when_cbsd_is_created(builder.build_post_data())
        cbsd = self.when_cbsd_is_fetched(builder.serial_number)
        self.then_cbsd_is(cbsd, builder.build_unregistered_data())

        state = self.when_cbsd_asks_for_state()
        self.then_state_is(state, get_empty_state())

        self.then_state_is_eventually(builder.build_grant_state_data())

        cbsd = self.when_cbsd_is_fetched(builder.serial_number)
        self.then_cbsd_is(cbsd, builder.build_registered_active_data())

        return cbsd['id']

    def when_cbsd_is_created(self, data: Dict[str, Any], expected_status: int = HTTPStatus.CREATED):
        r = send_request_to_backend('post', 'cbsds', json=data)
        self.assertEqual(r.status_code, expected_status)

    def when_cbsd_is_fetched(self, serial_number: str = None) -> Dict[str, Any]:
        return self._check_for_cbsd(serial_number=serial_number)

    def when_logs_are_fetched(self, params: Dict[str, Any]) -> Dict[str, Any]:
        r = send_request_to_backend('get', 'logs', params=params)
        self.assertEqual(r.status_code, HTTPStatus.OK)
        data = r.json()
        return data

    def when_cbsd_is_deleted(self, cbsd_id: int):
        r = send_request_to_backend('delete', f'cbsds/{cbsd_id}')
        self.assertEqual(r.status_code, HTTPStatus.NO_CONTENT)

    def when_cbsd_is_updated(self, cbsd_id: int, data: Dict[str, Any], expected_status: int = HTTPStatus.NO_CONTENT):
        r = send_request_to_backend('put', f'cbsds/{cbsd_id}', json=data)
        self.assertEqual(r.status_code, expected_status)

    def when_cbsd_is_deregistered(self, cbsd_id: int):
        r = send_request_to_backend('post', f'cbsds/{cbsd_id}/deregister')
        self.assertEqual(r.status_code, HTTPStatus.NO_CONTENT)

    def when_cbsd_asks_for_state(self) -> CBSDStateResult:
        return self.dp_client.GetCBSDState(get_cbsd_request(self.serial_number))

    @staticmethod
    def when_cbsd_is_inactive():
        inactivity = 3
        polling = 1
        delta = 3
        total_wait_time = inactivity + 2 * polling + delta
        sleep(total_wait_time)

    @contextmanager
    def while_cbsd_is_active(self):
        done = Event()

        def keep_asking_for_state():
            while not done.wait(timeout=1):
                self.when_cbsd_asks_for_state()

        t = Thread(target=keep_asking_for_state)
        try:
            t.start()
            yield
        finally:
            done.set()
            t.join()

    def then_cbsd_is(self, actual: Dict[str, Any], expected: Dict[str, Any]):
        actual = actual.copy()
        del actual['id']
        grant = actual.get('grant')
        if grant:
            del grant['grant_expire_time']
            del grant['transmit_expire_time']
        self.assertEqual(actual, expected)

    def then_cbsd_is_deleted(self, serial_number: str):
        self._check_for_cbsd(serial_number=serial_number, should_exist=False)

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
    def then_message_is_eventually_sent(self, filters: Dict[str, Any]):
        logs = self.when_logs_are_fetched(filters)
        self.assertEqual(logs["total_count"], 1)

    def delete_cbsd(self, cbsd_id: int):
        filters = get_filters_for_request_type('deregistration', self.serial_number)

        self.when_cbsd_is_deleted(cbsd_id)
        self.then_cbsd_is_deleted(self.serial_number)

        state = self.when_cbsd_asks_for_state()
        self.then_state_is(state, get_empty_state())

        self.then_message_is_eventually_sent(filters)

    @staticmethod
    def get_sas_provision_messages() -> List[str]:
        names = ['heartbeat', 'grant', 'spectrumInquiry', 'registration']
        return [f'{x}Response' for x in names]

    def _verify_logs_count(self, params):
        using_filters, _operator, expected_count = params
        logs = self.when_logs_are_fetched(using_filters)
        logs_len = len(logs["logs"])
        comparison = _operator(logs_len, expected_count)
        self.assertTrue(comparison)

    def _check_for_cbsd(self, serial_number: str, should_exist: bool = True) -> Optional[Dict[str, Any]]:
        params = {'serial_number': serial_number}
        expected_count = 1 if should_exist else 0
        r = send_request_to_backend('get', 'cbsds', params=params)
        self.assertEqual(r.status_code, HTTPStatus.OK)
        data = r.json()
        total_count = data.get('total_count')
        self.assertEqual(total_count, expected_count)
        cbsds = data.get('cbsds', [])
        self.assertEqual(len(cbsds), expected_count)
        if should_exist:
            return cbsds[0]


def get_current_sas_filters(serial_number: str) -> Dict[str, Any]:
    return {
        'serial_number': serial_number,
        'from': SAS,
        'to': DP,
        'end': now(),
    }


def get_filters_for_request_type(request_type: str, serial_number: str) -> Dict[str, Any]:
    return {
        'serial_number': serial_number,
        'type': f'{request_type}Response',
        'begin': now(),
    }


def get_empty_state() -> CBSDStateResult:
    return CBSDStateResult(radio_enabled=False)


def get_cbsd_request(serial_number: str) -> CBSDRequest:
    return CBSDRequest(serial_number=serial_number)


def now() -> str:
    return datetime.now(timezone.utc).isoformat()


@retry(stop_max_attempt_number=30, wait_fixed=1000)
def wait_for_elastic_to_start() -> None:
    requests.get(f'{config.ELASTICSEARCH_URL}/_status')


def when_elastic_indexes_data():
    # TODO use retrying instead
    sleep(15)


def _delete_dp_elasticsearch_indices() -> None:
    requests.delete(f"{config.ELASTICSEARCH_URL}/{config.ELASTICSEARCH_INDEX}*")


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


class CbsdAPIDataBuilder:
    def __init__(self):
        self.serial_number = str(uuid4())
        self.fcc_id = SOME_FCC_ID
        self.preferred_bandwidth_mhz = 20
        self.preferred_frequencies_mhz = []
        self.frequency_mhz = 3625
        self.bandwidth_mhz = 10
        self.max_eirp = 28
        self.desired_state = 'registered'

    def with_serial_number(self, serial_number: str) -> CbsdAPIDataBuilder:
        self.serial_number = serial_number
        return self

    def with_fcc_id(self, fcc_id: str) -> CbsdAPIDataBuilder:
        self.fcc_id = fcc_id
        return self

    def with_frequency_preferences(self, bandwidth_mhz: int, frequencies_mhz: List[int]) -> CbsdAPIDataBuilder:
        self.preferred_bandwidth_mhz = bandwidth_mhz
        self.preferred_frequencies_mhz = frequencies_mhz
        return self

    def with_desired_state(self, desired_state: str) -> CbsdAPIDataBuilder:
        self.desired_state = desired_state
        return self

    def with_expected_grant(self, bandwidth_mhz: int, frequency_mhz: int, max_eirp: int) -> CbsdAPIDataBuilder:
        self.bandwidth_mhz = bandwidth_mhz
        self.frequency_mhz = frequency_mhz
        self.max_eirp = max_eirp
        return self

    def build_post_data(self) -> Dict[str, Any]:
        return {
            'capabilities': {
                'antenna_gain': 15,
                'max_power': 20,
                'min_power': 0,
                'number_of_antennas': 2,
            },
            'frequency_preferences': {
                'bandwidth_mhz': self.preferred_bandwidth_mhz,
                'frequencies_mhz': self.preferred_frequencies_mhz,
            },
            'fcc_id': self.fcc_id,
            'serial_number': self.serial_number,
            'user_id': USER_ID,
            'desired_state': self.desired_state,
        }

    def build_unregistered_data(self) -> Dict[str, Any]:
        data = self.build_post_data()
        data.update({
            'is_active': False,
            'state': 'unregistered',
        })
        return data

    def build_registered_inactive_data(self) -> Dict[str, Any]:
        data = self.build_post_data()
        data.update({
            'cbsd_id': f'{self.fcc_id}/{self.serial_number}',
            'is_active': False,
            'state': 'registered',
        })
        return data

    def build_registered_active_data(self) -> Dict[str, Any]:
        data = self.build_registered_inactive_data()
        data.update({
            'is_active': True,
            'grant': {
                'bandwidth_mhz': self.bandwidth_mhz,
                'frequency_mhz': self.frequency_mhz,
                'max_eirp': self.max_eirp,
                'state': 'authorized',
            },
        })
        return data

    def build_grant_state_data(self) -> CBSDStateResult:
        frequency_hz = int(1e6) * self.frequency_mhz
        half_bandwidth_hz = int(5e5) * self.bandwidth_mhz
        return CBSDStateResult(
            radio_enabled=True,
            channel=LteChannel(
                low_frequency_hz=frequency_hz - half_bandwidth_hz,
                high_frequency_hz=frequency_hz + half_bandwidth_hz,
                max_eirp_dbm_mhz=self.max_eirp,
            ),
        )
