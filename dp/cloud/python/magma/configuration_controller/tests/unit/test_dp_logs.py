from dp.protos.enodebd_dp_pb2 import CBSDRequest, CBSDStateResult, LteChannel
from parameterized import parameterized

from magma.db_service.models import DBRequestType, DBRequest, DBRequestState, DBCbsd, DBCbsdState, DBResponse
from magma.db_service.tests.local_db_test_case import LocalDBTestCase
from magma.fluentd_client.client import DPLog
from magma.fluentd_client.dp_logs import make_dp_log


DP = 'DP'
CBSD = 'CBSD'
SAS = 'SAS'
SOME_FCC_ID = 'some_fcc_id'
HEARTBEAT_REQUEST = 'heartbeatRequest'
SOME_SERIAL_NUMBER = 'some_serial_number'
SOME_NETWORK_ID = 'some_network_id'
SOME_MESSAGE = 'some_message'
CBSD_REGISTER = 'CBSDRegister'


class IncorrectDPLog(object):
    ...


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
        req_type = DBRequestType(name=HEARTBEAT_REQUEST)
        req_state = DBRequestState(name='some_req_state')
        self.session.add_all([cbsd, req_type, req_state])
        self.session.commit()

    @parameterized.expand([
        (False, '', '', ''),
        (True, SOME_SERIAL_NUMBER, SOME_FCC_ID, SOME_NETWORK_ID),
    ])
    def test_dp_log_created_from_db_request(self, with_cbsd, serial_num, fcc_id, network_id):
        # Given
        req_type = self.session.query(DBRequestType).first()
        req_state = self.session.query(DBRequestState).first()
        cbsd = None
        if with_cbsd:
            cbsd = self.session.query(DBCbsd).first()
        request = DBRequest(type=req_type, state=req_state, cbsd=cbsd, payload=SOME_MESSAGE)
        
        # When
        actual_log = make_dp_log(request)
        
        # Then
        expected_log = DPLog(
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
    def test_dp_log_created_from_db_request(self, with_cbsd, serial_num, fcc_id, network_id):
        # Given
        req_type = self.session.query(DBRequestType).first()
        req_state = self.session.query(DBRequestState).first()
        cbsd = None
        if with_cbsd:
            cbsd = self.session.query(DBCbsd).first()
        request = DBRequest(type=req_type, state=req_state, cbsd=cbsd, payload='some_request_message')
        resp_payload = {"response": {"responseCode": "0"}}
        response = DBResponse(request=request, response_code=200, payload=resp_payload)
        
        # When
        actual_log = make_dp_log(response)
        
        # Then
        expected_log = DPLog(
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
    def test_dp_log_created_from_grpc_request(self, with_cbsd, fcc_id, network_id):
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
            radio_enabled=True
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

    def test_make_dp_log_returns_type_error_for_unknown_message_type(self):
        with self.assertRaises(TypeError):
            make_dp_log(IncorrectDPLog())
