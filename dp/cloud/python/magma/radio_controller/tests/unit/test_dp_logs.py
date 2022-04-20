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

from dp.protos.enodebd_dp_pb2 import CBSDRequest, CBSDStateResult, LteChannel
from magma.db_service.models import DBCbsd, DBCbsdState
from magma.db_service.tests.local_db_test_case import LocalDBTestCase
from magma.fluentd_client.client import DPLog
from magma.fluentd_client.dp_logs import make_dp_log
from parameterized import parameterized

DP = 'DP'
CBSD = 'CBSD'
SAS = 'SAS'
SOME_FCC_ID = 'some_fcc_id'
HEARTBEAT_REQUEST = 'heartbeatRequest'
SOME_SERIAL_NUMBER = 'some_serial_number'
SOME_NETWORK_ID = 'some_network_id'
CBSD_REGISTER = 'CBSDRegister'


class DPLogsTestCase(LocalDBTestCase):

    def setUp(self):
        super().setUp()
        cbsd_state = DBCbsdState(name='some_cbsd_state')
        cbsd = DBCbsd(
            state=cbsd_state,
            fcc_id=SOME_FCC_ID,
            cbsd_serial_number=SOME_SERIAL_NUMBER,
            network_id=SOME_NETWORK_ID,
        )
        self.session.add(cbsd)
        self.session.commit()

    @parameterized.expand([
        (False, '', '', ''),
        (True, SOME_SERIAL_NUMBER, SOME_FCC_ID, SOME_NETWORK_ID),
    ])
    def test_dp_log_created_from_grpc_request(self, with_cbsd, serial_num, fcc_id, network_id):
        # Given
        cbsd = None
        if with_cbsd:
            cbsd = self.session.query(DBCbsd).first()
        message = CBSDRequest(
            user_id='some_user_id',
            fcc_id=fcc_id,
            serial_number=serial_num,
            min_power=2,
            max_power=3,
            antenna_gain=4,
            number_of_ports=5,
        )

        # When
        actual_log = make_dp_log(method_name=CBSD_REGISTER, message=message, cbsd=cbsd)

        # Then
        expected_log = DPLog(
            cbsd_serial_number=serial_num,
            fcc_id=fcc_id,
            log_from=CBSD,
            log_message=str(message),
            log_name='CBSDRegisterRequest',
            log_to=DP,
            network_id=network_id,
            response_code=None,
        )
        self.assertEqual(expected_log, actual_log)

    @parameterized.expand([
        (False, '', ''),
        (True, SOME_FCC_ID, SOME_NETWORK_ID),
    ])
    def test_dp_log_created_from_grpc_response(self, with_cbsd, fcc_id, network_id):
        # Given
        cbsd = None
        if with_cbsd:
            cbsd = self.session.query(DBCbsd).first()
        channel = LteChannel(
            low_frequency_hz=1,
            high_frequency_hz=2,
            max_eirp_dbm_mhz=3,
        )
        message = CBSDStateResult(
            channel=channel,
            radio_enabled=True,
        )

        # When
        actual_log = make_dp_log(
            method_name=CBSD_REGISTER,
            message=message,
            cbsd=cbsd,
            serial_number=SOME_SERIAL_NUMBER,
        )

        # Then
        expected_log = DPLog(
            cbsd_serial_number=SOME_SERIAL_NUMBER,
            fcc_id=fcc_id,
            log_from=DP,
            log_message=str(message),
            log_name='CBSDRegisterResponse',
            log_to=CBSD,
            network_id=network_id,
            response_code=None,
        )
        self.assertEqual(expected_log, actual_log)
