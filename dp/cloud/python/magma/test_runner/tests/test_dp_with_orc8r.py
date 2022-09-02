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

import grpc
import pytest
import requests
from dp.protos.cbsd_pb2 import CBSDStateResult, EnodebdUpdateCbsdRequest
from magma.test_runner.config import TestConfig
from magma.test_runner.tests.api_data_builder import CbsdAPIDataBuilder
from magma.test_runner.tests.integration_testcase import (
    Orc8rIntegrationTestCase,
)
from parameterized import parameterized
from retrying import retry

config = TestConfig()

DP_HTTP_PREFIX = 'magma/v1/dp'
NETWORK = 'some_network'
SOME_FCC_ID = "some_fcc_id"
OTHER_FCC_ID = "other_fcc_id"
USER_ID = "some_user_id"
SAS = 'SAS'
DP = 'DP'
UNREGISTERED = "unregistered"
DATETIME_WAY_BACK = '2010-03-28T09:13:25.407877399+00:00'
MAGMA_CLIENT_CERT_SERIAL_KEY = 'x-magma-client-cert-serial'
MAGMA_CLIENT_CERT_SERIAL_VALUE = '7ZZXAF7CAETF241KL22B8YRR7B5UF401'


@pytest.mark.orc8r
class DomainProxyOrc8rTestCase(Orc8rIntegrationTestCase):
    def setUp(self) -> None:
        self.now = now()
        self.serial_number = self._testMethodName + '_' + str(uuid4())
        # TODO why do we do that? Setup was supposed to be done by deployment...
        self.prometheus_url = f'http://{config.PROMETHEUS_SERVICE_HOST}:{config.PROMETHEUS_SERVICE_PORT}'
        requests.post(f"{self.prometheus_url}/api/v1/admin/tsdb/delete_series")

    def test_cbsd_sas_flow(self):
        builder = CbsdAPIDataBuilder() \
            .with_single_step_enabled(True) \
            .with_serial_number(self.serial_number)

        update_request = builder.build_enodebd_update_request(indoor_deployment=True)
        cbsd_id = self.given_cbsd_provisioned(builder, update_request)

        with self.while_cbsd_is_active(update_request):
            self.then_provision_logs_are_sent()
            self.then_logs_contain({'serial_number': self.serial_number}, ["EnodebdUpdateCbsd", "CbsdStateResponse"])

        self.delete_cbsd(cbsd_id)

    def test_cbsd_unregistered_when_requested_by_desired_state(self):
        builder = CbsdAPIDataBuilder() \
            .with_serial_number(self.serial_number)

        update_request = builder.build_enodebd_update_request()
        cbsd_id = self.given_cbsd_provisioned(builder, update_request)

        with self.while_cbsd_is_active(update_request):
            filters = self._get_filters_for_request_type('deregistration', self.serial_number)

            builder = builder.with_desired_state(UNREGISTERED)
            self.when_cbsd_is_updated_by_user(cbsd_id, builder.payload)

            # TODO maybe asking for state (cbsd api instead of log api) would be better
            self.then_message_is_eventually_sent(filters)

    @parameterized.expand([
        (10.6, 11.6, True),
        (10.500001, 11.500001, False),
    ])
    def test_cbsd_unregistered_when_enodebd_changes_coordinates(self, lat, lon, should_deregister):
        builder = CbsdAPIDataBuilder() \
            .with_serial_number(self.serial_number) \
            .with_single_step_enabled(True)

        update_request = builder.build_enodebd_update_request(indoor_deployment=True)

        self.given_cbsd_provisioned(builder, update_request)

        update_request.installation_param.latitude_deg.value = lat
        update_request.installation_param.longitude_deg.value = lon

        with self.while_cbsd_is_active(update_request):

            filters = self._get_filters_for_request_type('deregistration', self.serial_number)
            if should_deregister:
                self.then_message_is_eventually_sent(filters)
            else:
                self.then_message_is_never_sent(filters)

    def test_sas_flow_restarted_when_user_requested_deregistration(self):
        builder = CbsdAPIDataBuilder() \
            .with_serial_number(self.serial_number)

        update_request = builder.build_enodebd_update_request()
        cbsd_id = self.given_cbsd_provisioned(builder, update_request)

        with self.while_cbsd_is_active(update_request):
            filters = self._get_filters_for_request_type('deregistration', self.serial_number)

            self.when_cbsd_is_deregistered(cbsd_id)

            self.then_message_is_eventually_sent(filters)

            self.then_state_is_eventually(builder.build_grant_state_data(), update_request)

    def test_grants_relinquished_when_user_requested_relinquish(self):
        builder = CbsdAPIDataBuilder() \
            .with_serial_number(self.serial_number)

        update_request = builder.build_enodebd_update_request()
        cbsd_id = self.given_cbsd_provisioned(builder, update_request)

        with self.while_cbsd_is_active(update_request):
            self.when_cbsd_is_relinquished(cbsd_id)

            self.then_state_is_eventually(get_empty_state(), update_request)

            self.then_state_is_eventually(builder.build_grant_state_data(), update_request)

    def test_sas_flow_restarted_for_updated_cbsd(self):
        builder = CbsdAPIDataBuilder() \
            .with_serial_number(self.serial_number)

        update_request = builder.build_enodebd_update_request()
        cbsd_id = self.given_cbsd_provisioned(builder, update_request)

        with self.while_cbsd_is_active(update_request):
            builder = builder.with_fcc_id(OTHER_FCC_ID)
            self.when_cbsd_is_updated_by_user(cbsd_id, builder.payload)

            filters = self._get_filters_for_request_type('deregistration', self.serial_number)
            self.then_message_is_eventually_sent(filters)

            self.then_state_is_eventually(builder.build_grant_state_data(), update_request)

            cbsd = self.when_cbsd_is_fetched(builder.payload["serial_number"])
            self.then_cbsd_is(
                cbsd,
                builder
                .with_cbsd_id(f"{OTHER_FCC_ID}/{self.serial_number}")
                .with_is_active(True)
                .payload,
            )

    def test_activity_status(self):
        builder = CbsdAPIDataBuilder() \
            .with_serial_number(self.serial_number)
        update_request = builder.build_enodebd_update_request()
        self.given_cbsd_provisioned(builder, update_request)

        cbsd = self.when_cbsd_is_fetched(builder.payload["serial_number"])
        self.then_cbsd_is(cbsd, builder.with_is_active(True).payload)

        self.when_cbsd_is_inactive()
        cbsd = self.when_cbsd_is_fetched(builder.payload["serial_number"])
        # TODO this is bad builders shouldn't be used like that
        self.then_cbsd_is(cbsd, builder.without_grants().with_is_active(False).payload)

    def test_frequency_preferences(self):
        builder = CbsdAPIDataBuilder() \
            .with_serial_number(self.serial_number) \
            .with_frequency_preferences(5, [3625])

        self.when_cbsd_is_created(builder.payload)

        sn = builder.payload['serial_number']
        fcc_id = builder.payload['fcc_id']
        cbsd_id = f'{fcc_id}/{sn}'

        builder \
            .with_state('registered') \
            .with_is_active(True) \
            .with_full_installation_param(indoor_deployment=False) \
            .with_cbsd_id(cbsd_id) \
            .with_grant(bandwidth_mhz=5, frequency_mhz=3625, max_eirp=31)

        update_request = builder.build_enodebd_update_request()

        with self.while_cbsd_is_active(update_request):
            self.then_state_is_eventually(builder.build_grant_state_data(), update_request)
            cbsd = self.when_cbsd_is_fetched(sn)
            self.then_cbsd_is(cbsd, builder.payload)

    def test_carrier_aggregation(self):
        builder = CbsdAPIDataBuilder() \
            .with_serial_number(self.serial_number) \
            .with_frequency_preferences(20, [3600]) \
            .with_carrier_aggregation()

        self.when_cbsd_is_created(builder.payload)

        sn = builder.payload['serial_number']
        fcc_id = builder.payload['fcc_id']
        cbsd_id = f'{fcc_id}/{sn}'

        builder \
            .with_state('registered') \
            .with_is_active(True) \
            .with_full_installation_param(indoor_deployment=False) \
            .with_cbsd_id(cbsd_id) \
            .with_grant(bandwidth_mhz=20, frequency_mhz=3580, max_eirp=25) \
            .with_grant(bandwidth_mhz=20, frequency_mhz=3600, max_eirp=25)

        update_request = builder.build_enodebd_update_request()

        with self.while_cbsd_is_active(update_request):
            self.then_state_is_eventually(builder.build_grant_state_data(), update_request)
            cbsd = self.when_cbsd_is_fetched(sn)
            self.then_cbsd_is(cbsd, builder.payload)

    def test_creating_cbsd_with_the_same_unique_fields_returns_409(self):
        builder = CbsdAPIDataBuilder() \
            .with_serial_number(self.serial_number)

        self.when_cbsd_is_created(builder.payload)
        self.when_cbsd_is_created(builder.payload, expected_status=HTTPStatus.CONFLICT)

    def test_updating_cbsd_returns_409_when_setting_existing_serial_num(self):
        builder = CbsdAPIDataBuilder() \

        cbsd1_serial = self.serial_number + "_foo"
        cbsd2_serial = self.serial_number + "_bar"
        self.when_cbsd_is_created(builder.with_serial_number(cbsd1_serial).payload)
        self.when_cbsd_is_created(builder.with_serial_number(cbsd2_serial).payload)
        cbsd2 = self.when_cbsd_is_fetched(serial_number=cbsd2_serial)
        self.when_cbsd_is_updated_by_user(
            cbsd_id=cbsd2.get("id"),
            data=builder.with_serial_number(cbsd1_serial).payload,
            expected_status=HTTPStatus.CONFLICT,
        )

    def test_fetch_cbsds_filtered_by_serial_number(self):
        cbsd1_serial = self.serial_number + "_foo"
        cbsd2_serial = self.serial_number + "_bar"

        builder1 = CbsdAPIDataBuilder() \
            .with_serial_number(cbsd1_serial)
        builder2 = CbsdAPIDataBuilder() \
            .with_serial_number(cbsd2_serial)

        self.when_cbsd_is_created(builder1.payload)
        self.when_cbsd_is_created(builder2.payload)

        cbsd1 = self.when_cbsd_is_fetched(serial_number=cbsd1_serial)
        cbsd2 = self.when_cbsd_is_fetched(serial_number=cbsd2_serial)

        self.then_cbsd_is(
            cbsd1,
            builder1
            .with_state(UNREGISTERED)
            .with_indoor_deployment(False)
            .with_is_active(False)
            .payload,
        )
        self.then_cbsd_is(
            cbsd2,
            builder2
            .with_state(UNREGISTERED)
            .with_indoor_deployment(False)
            .with_is_active(False)
            .payload,
        )

    def test_fetching_logs_with_custom_filters(self):
        builder = CbsdAPIDataBuilder() \
            .with_serial_number(self.serial_number)

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

        update_request = builder.build_enodebd_update_request()

        self.given_cbsd_provisioned(builder, update_request)
        with self.while_cbsd_is_active(update_request):
            self._verify_logs_count(scenarios)

    def then_metrics_are_in_prometheus(self):
        metrics = [
            'dp_rc_grpc_request_processing_seconds_bucket',
            'dp_rc_grpc_request_processing_seconds_count',
            'dp_rc_grpc_request_processing_seconds_sum',
            'dp_cc_pending_requests_fetching_seconds_bucket',
            'dp_cc_pending_requests_fetching_seconds_count',
            'dp_cc_pending_requests_fetching_seconds_sum',
            'dp_cc_request_processing_seconds_bucket',
            'dp_cc_request_processing_seconds_count',
            'dp_cc_request_processing_seconds_sum',
            'dp_cc_response_processing_seconds_bucket',
            'dp_cc_response_processing_seconds_count',
            'dp_cc_response_processing_seconds_sum',
        ]
        for m in metrics:
            resp = requests.get(self.prometheus_url + '/api/v1/query', params={'query': m})
            data = resp.json().get('data')
            self.assertIsNotNone(data)
            result = data.get('result')
            self.assertIsNotNone(result)
            self.assertGreater(len(result), 0)

    def given_cbsd_provisioned(self, builder: CbsdAPIDataBuilder, request: EnodebdUpdateCbsdRequest) -> int:
        self.when_cbsd_is_created(builder.payload)
        cbsd = self.when_cbsd_is_fetched(builder.payload["serial_number"])
        self.then_cbsd_is(
            cbsd,
            builder
            .with_state(UNREGISTERED)
            .with_is_active(False)
            .with_indoor_deployment(False)
            .payload,
        )

        state = self.when_cbsd_is_updated_by_enodebd(request)
        self.then_state_is(state, get_empty_state())

        cbsd = self._check_cbsd_successfully_provisioned(builder, request)

        return cbsd['id']

    def _check_cbsd_successfully_provisioned(
            self, builder: CbsdAPIDataBuilder, request: EnodebdUpdateCbsdRequest,
    ) -> Dict[str, Any]:
        sn = builder.payload["serial_number"]
        fcc_id = builder.payload["fcc_id"]
        cbsd_id = f"{fcc_id}/{sn}"
        # TODO this is very unexpected, builder modifies state
        expected_cbsd = builder \
            .with_is_active(True) \
            .with_grant() \
            .with_state("registered") \
            .with_cbsd_id(cbsd_id) \
            .with_full_installation_param(
                latitude_deg=request.installation_param.latitude_deg.value,
                longitude_deg=request.installation_param.longitude_deg.value,
                indoor_deployment=request.installation_param.indoor_deployment.value,
                height_type=request.installation_param.height_type.value,
                height_m=request.installation_param.height_m.value,
            ) \
            .payload

        self.then_state_is_eventually(builder.build_grant_state_data(), request)

        cbsd = self.when_cbsd_is_fetched(sn)

        self.then_cbsd_is(cbsd, expected_cbsd)

        return cbsd

    def then_provision_logs_are_sent(self):
        self.then_logs_are_in_only_one_network()

        filters = self._get_filters_for_request_type('heartbeat', self.serial_number)
        self.then_message_is_eventually_sent(filters)

    def then_logs_are_in_only_one_network(self):
        self.then_logs_are(_get_current_sas_filters(self.serial_number), self.get_sas_provision_messages())
        self.then_logs_are(_get_current_sas_filters(self.serial_number), [], network="someOtherNetworkId")

    def when_cbsd_is_created(self, data: Dict[str, Any], expected_status: int = HTTPStatus.CREATED):
        r = send_request_to_backend('post', 'cbsds', json=data)
        self.assertEqual(r.status_code, expected_status)

    def when_cbsd_is_fetched(self, serial_number: str = None) -> Dict[str, Any]:
        return self._check_for_cbsd(serial_number=serial_number)

    def when_logs_are_fetched(self, params: Dict[str, Any], network: Optional[str] = None) -> Dict[str, Any]:
        r = send_request_to_backend('get', 'logs', params=params, network=network)
        self.assertEqual(r.status_code, HTTPStatus.OK)
        data = r.json()
        return data

    def when_cbsd_is_deleted(self, cbsd_id: int):
        r = send_request_to_backend('delete', f'cbsds/{cbsd_id}')
        self.assertEqual(r.status_code, HTTPStatus.NO_CONTENT)

    def when_cbsd_is_updated_by_user(
            self, cbsd_id: int, data: Dict[str, Any],
            expected_status: int = HTTPStatus.NO_CONTENT,
    ):
        r = send_request_to_backend('put', f'cbsds/{cbsd_id}', json=data)
        self.assertEqual(r.status_code, expected_status)

    def when_cbsd_is_updated_by_enodebd(self, req: EnodebdUpdateCbsdRequest) -> CBSDStateResult:
        return self.orc8r_dp_client.EnodebdUpdateCbsd(
            request=req, metadata=(
                (MAGMA_CLIENT_CERT_SERIAL_KEY, MAGMA_CLIENT_CERT_SERIAL_VALUE),
            ),
        )

    def when_cbsd_is_deregistered(self, cbsd_id: int):
        r = send_request_to_backend('post', f'cbsds/{cbsd_id}/deregister')
        self.assertEqual(r.status_code, HTTPStatus.NO_CONTENT)

    def when_cbsd_is_relinquished(self, cbsd_id: int):
        r = send_request_to_backend('post', f'cbsds/{cbsd_id}/relinquish')
        self.assertEqual(r.status_code, HTTPStatus.NO_CONTENT)

    @staticmethod
    def when_cbsd_is_inactive():
        inactivity = 3
        polling = 1
        delta = 3
        total_wait_time = inactivity + 2 * polling + delta
        sleep(total_wait_time)

    @contextmanager
    def while_cbsd_is_active(self, update_request):
        done = Event()

        def keep_asking_for_state():
            while not done.wait(timeout=1):
                self.when_cbsd_is_updated_by_enodebd(update_request)

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
        for grant in actual.get('grants', []):
            del grant['grant_expire_time']
            del grant['transmit_expire_time']
        self.assertEqual(expected, actual)

    def then_cbsd_is_deleted(self, serial_number: str):
        self._check_for_cbsd(serial_number=serial_number, should_exist=False)

    def then_state_is(self, actual: CBSDStateResult, expected: CBSDStateResult):
        self.assertEqual(actual, expected)

    @retry(stop_max_attempt_number=30, wait_fixed=1000)
    def then_state_is_eventually(self, expected, request):
        actual = self.when_cbsd_is_updated_by_enodebd(request)
        self.then_state_is(actual, expected)

    @retry(stop_max_attempt_number=30, wait_fixed=1000)
    def then_logs_are(self, filters: Dict[str, Any], expected: List[str], network: Optional[str] = None):
        actual = self._get_log_types(filters, network)
        self.assertCountEqual(expected, actual)

    @retry(stop_max_attempt_number=30, wait_fixed=1000)
    def then_logs_contain(self, filters: Dict[str, Any], expected: List[str], network: Optional[str] = None):
        actual = self._get_log_types(filters, network)
        for log_type in expected:
            self.assertIn(log_type, actual)

    def _get_log_types(self, filters: Dict[str, Any], network: Optional[str] = None):
        actual = self.when_logs_are_fetched(filters, network)
        return [x['type'] for x in actual['logs']]

    @retry(stop_max_attempt_number=60, wait_fixed=1000)
    def then_message_is_eventually_sent(self, filters: Dict[str, Any]):
        logs = self.when_logs_are_fetched(filters)
        self.assertEqual(logs["total_count"], 1)

    def then_message_is_never_sent(self, filters: Dict[str, Any]):
        sleep(5)
        logs = self.when_logs_are_fetched(filters)
        self.assertEqual(logs["total_count"], 0)

    def delete_cbsd(self, cbsd_id: int):
        filters = self._get_filters_for_request_type('deregistration', self.serial_number)

        self.when_cbsd_is_deleted(cbsd_id)
        self.then_cbsd_is_deleted(self.serial_number)

        try:
            self.when_cbsd_is_updated_by_enodebd(EnodebdUpdateCbsdRequest(serial_number=self.serial_number))
        except grpc.RpcError as e:
            self.assertEqual(grpc.StatusCode.NOT_FOUND, e.code())

        self.then_message_is_eventually_sent(filters)

    @staticmethod
    def get_sas_provision_messages() -> List[str]:
        names = ['heartbeat', 'grant', 'spectrumInquiry', 'registration']
        return [f'{x}Response' for x in names]

    @retry(stop_max_attempt_number=30, wait_fixed=1000)
    def _verify_logs_count(self, scenarios):
        for params in scenarios:
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

    def _get_filters_for_request_type(self, request_type: str, serial_number: str) -> Dict[str, Any]:
        return {
            'serial_number': serial_number,
            'type': f'{request_type}Response',
            'begin': self.now,
        }


def _get_current_sas_filters(serial_number: str) -> Dict[str, Any]:
    return {
        'serial_number': serial_number,
        'from': SAS,
        'to': DP,
        'end': now(),
    }


def get_empty_state() -> CBSDStateResult:
    return CBSDStateResult(radio_enabled=False)


def now() -> str:
    return datetime.now(timezone.utc).isoformat()


def send_request_to_backend(
        method: str, url_suffix: str, params: Optional[Dict[str, Any]] = None,
        json: Optional[Dict[str, Any]] = None,
        network: Optional[str] = None,
) -> requests.Response:
    return requests.request(
        method,
        f'{config.HTTP_SERVER}/{DP_HTTP_PREFIX}/{network or NETWORK}/{url_suffix}',
        cert=(config.DP_CERT_PATH, config.DP_SSL_KEY_PATH),
        verify=False,  # noqa: S501
        params=params,
        json=json,
    )
