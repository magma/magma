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
import datetime

from freezegun import freeze_time
from magma.configuration_controller.custom_types.custom_types import DBResponse
from magma.db_service.models import (
    DBCbsd,
    DBCbsdState,
    DBRequest,
    DBRequestType,
)
from magma.db_service.tests.local_db_test_case import LocalDBTestCase
from magma.fluentd_client.client import DPLog
from magma.fluentd_client.dp_logs import make_dp_log, now
from parameterized import parameterized

DP = 'DP'
CBSD = 'CBSD'
SAS = 'SAS'
SOME_FCC_ID = 'some_fcc_id'
HEARTBEAT_REQUEST = 'heartbeatRequest'
SOME_SERIAL_NUMBER = 'some_serial_number'
SOME_NETWORK_ID = 'some_network_id'
SOME_MESSAGE = 'some_message'
SOME_DATE = datetime.datetime.now(datetime.timezone.utc)
SOME_TIMESTAMP = int(SOME_DATE.timestamp())


class IncorrectDPLog(object):
    ...


@freeze_time(SOME_DATE)
class DPLogsTestCase(LocalDBTestCase):

    def setUp(self):
        super().setUp()
        cbsd_state = DBCbsdState(name='some_cbsd_state')
        cbsd = DBCbsd(
            state=cbsd_state,
            desired_state=cbsd_state,
            fcc_id=SOME_FCC_ID,
            cbsd_serial_number=SOME_SERIAL_NUMBER,
            network_id=SOME_NETWORK_ID,
        )
        req_type = DBRequestType(name=HEARTBEAT_REQUEST)
        self.session.add_all([cbsd, req_type])
        self.session.commit()

    @parameterized.expand([
        (False, '', '', ''),
        (True, SOME_SERIAL_NUMBER, SOME_FCC_ID, SOME_NETWORK_ID),
    ])
    def test_dp_log_created_from_db_request(self, with_cbsd, serial_num, fcc_id, network_id):
        # Given
        req_type = self.session.query(DBRequestType).first()
        cbsd = None
        if with_cbsd:
            cbsd = self.session.query(DBCbsd).first()
        request = DBRequest(type=req_type, cbsd=cbsd, payload=SOME_MESSAGE)

        # When
        actual_log = make_dp_log(request)

        # Then
        expected_log = DPLog(
            event_timestamp=SOME_TIMESTAMP,
            cbsd_serial_number=serial_num,
            fcc_id=fcc_id,
            log_from=DP,
            log_message=SOME_MESSAGE,
            log_name=HEARTBEAT_REQUEST,
            log_to=SAS,
            network_id=network_id,
            response_code=None,
        )
        self.assertEqual(expected_log, actual_log)

    @parameterized.expand([
        (False, '', '', ''),
        (True, SOME_SERIAL_NUMBER, SOME_FCC_ID, SOME_NETWORK_ID),
    ])
    def test_dp_log_created_from_db_response(self, with_cbsd, serial_num, fcc_id, network_id):
        # Given
        req_type = self.session.query(DBRequestType).first()
        cbsd = None
        if with_cbsd:
            cbsd = self.session.query(DBCbsd).first()
        request = DBRequest(type=req_type, cbsd=cbsd, payload='some_request_message')
        resp_payload = {"response": {"responseCode": "0"}}
        response = DBResponse(request=request, response_code=200, payload=resp_payload)

        # When
        actual_log = make_dp_log(response)

        # Then
        expected_log = DPLog(
            event_timestamp=SOME_TIMESTAMP,
            cbsd_serial_number=serial_num,
            fcc_id=fcc_id,
            log_from=SAS,
            log_message="{'response': {'responseCode': '0'}}",
            log_name='heartbeatResponse',
            log_to=DP,
            network_id=network_id,
            response_code="0",
        )
        self.assertEqual(expected_log, actual_log)

    def test_make_dp_log_returns_type_error_for_unknown_message_type(self):
        with self.assertRaises(TypeError):
            make_dp_log(IncorrectDPLog())

    def test_datetime_now(self):
        self.assertEqual(SOME_TIMESTAMP, now())
