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

import asyncio
import datetime
from concurrent import futures
from unittest import TestCase
from unittest.mock import ANY, MagicMock, call, patch

import grpc
import magma.magmad.bootstrap_manager as bm
from cryptography import x509
from cryptography.exceptions import InternalError
from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives import hashes, serialization
from cryptography.hazmat.primitives.asymmetric import ec, rsa
from cryptography.hazmat.primitives.asymmetric.utils import encode_dss_signature
from google.protobuf.timestamp_pb2 import Timestamp
from orc8r.protos import bootstrapper_pb2_grpc
from orc8r.protos.bootstrapper_pb2 import Challenge, ChallengeKey
from orc8r.protos.certifier_pb2 import CSR, Certificate

# Allow access to protected variables for unit testing
# pylint: disable=protected-access

BM = 'magma.magmad.bootstrap_manager'


# https://stackoverflow.com/questions/32480108/mocking-async-call-in-python-3-5
def AsyncMock():
    coro = MagicMock(name="CoroutineResult")
    corofunc = MagicMock(
        name="CoroutineFunction",
        side_effect=asyncio.coroutine(coro),
    )
    corofunc.coro = coro
    return corofunc


def make_awaitable(func):
    future = asyncio.Future()
    future.set_result(None)
    func.return_value = future


class DummpyBootstrapperServer(bootstrapper_pb2_grpc.BootstrapperServicer):
    def __init__(self):
        pass

    def add_to_server(self, server):
        bootstrapper_pb2_grpc.add_BootstrapperServicer_to_server(self, server)

    def GetChallenge(self, request, context):
        challenge = Challenge(
            challenge=b'simple_challenge',
            key_type=ChallengeKey.ECHO,
        )
        return challenge

    def RequestSign(self, request, context):
        return create_cert_message()


class BootstrapManagerTest(TestCase):
    @patch('magma.common.cert_utils.write_key')
    @patch('%s.BootstrapManager._bootstrap_check' % BM)
    @patch('%s.snowflake.snowflake' % BM)
    @patch('%s.load_service_config' % BM)
    # Pylint doesn't handle decorators correctly
    # pylint: disable=arguments-differ, unused-argument
    def setUp(
        self,
        load_service_config_mock,
        snowflake_mock,
        bootstrap_check_mock,
        write_key_mock,
    ):

        self.gateway_key_file = '__test_gw.key'
        self.gateway_cert_file = '__test_hw_cert'
        self.hw_id = 'hwid_test'

        load_service_config_mock.return_value = {
            'gateway_key': self.gateway_key_file,
            'gateway_cert': self.gateway_cert_file,
        }
        snowflake_mock.return_value = self.hw_id

        self.loop = asyncio.new_event_loop()
        service = MagicMock()
        service.loop = self.loop
        asyncio.set_event_loop(self.loop)
        service.config = {
            'bootstrap_config': {
                'challenge_key': '__test_challenge.key',
            },
        }

        bootstrap_success_cb = MagicMock()

        self.manager = bm.BootstrapManager(service, bootstrap_success_cb)
        self.manager._bootstrap_success_cb = bootstrap_success_cb
        self.manager.start_bootstrap_manager()
        write_key_mock.assert_has_calls(
            [call(ANY, service.config['bootstrap_config']['challenge_key'])],
        )

        # Bind the rpc server to a free port
        self._rpc_server = grpc.server(
            futures.ThreadPoolExecutor(max_workers=10),
        )
        port = self._rpc_server.add_insecure_port('0.0.0.0:0')
        # Add the servicer
        self._servicer = DummpyBootstrapperServer()
        self._servicer.add_to_server(self._rpc_server)
        self._rpc_server.start()
        # Create a rpc stub
        self.channel = grpc.insecure_channel('0.0.0.0:{}'.format(port))

        self.manager.SHORT_BOOTSTRAP_RETRY_INTERVAL = datetime.timedelta(
            seconds=0,
        )
        self.manager.LONG_BOOTSTRAP_RETRY_INTERVAL = datetime.timedelta(
            seconds=0,
        )

    def tearDown(self):
        self._rpc_server.stop(None)
        self.manager.stop_bootstrap_manager()
        self.loop.close()

    @patch('magma.common.cert_utils.load_cert')
    @patch('%s.BootstrapManager._bootstrap_now' % BM)
    @patch('%s.BootstrapManager._schedule_next_bootstrap_check' % BM)
    def test__bootstrap_check(
        self,
        schedule_bootstrap_check_mock,
        bootstrap_now_mock,
        load_cert_mock,
    ):
        async def test():
            make_awaitable(self.manager._bootstrap_now)
            make_awaitable(self.manager._bootstrap_success_cb)

            # cannot load cert
            load_cert_mock.side_effect = IOError
            await self.manager._bootstrap_check()
            load_cert_mock.assert_has_calls([call(self.gateway_cert_file)])
            bootstrap_now_mock.assert_has_calls([call()])

            # invalid not_before
            load_cert_mock.reset_mock()
            load_cert_mock.side_effect = None  # clear IOError side effect
            not_before = datetime.datetime.utcnow() + datetime.timedelta(
                days=3,
            )
            not_after = not_before + datetime.timedelta(days=3)
            load_cert_mock.return_value = create_cert(not_before, not_after)
            await self.manager._bootstrap_check()
            bootstrap_now_mock.assert_has_calls([call()])

            # invalid not_after
            load_cert_mock.reset_mock()
            not_before = datetime.datetime.utcnow()
            not_after = not_before + datetime.timedelta(hours=1)
            load_cert_mock.return_value = create_cert(not_before, not_after)
            await self.manager._bootstrap_check()
            bootstrap_now_mock.assert_has_calls([call()])

            # cert is present and valid,
            load_cert_mock.reset_mock()
            not_before = datetime.datetime.utcnow()
            not_after = not_before + datetime.timedelta(days=10)
            load_cert_mock.return_value = create_cert(not_before, not_after)
            await self.manager._bootstrap_check()
            schedule_bootstrap_check_mock.assert_has_calls([call()])

        # Cancel the loop so that there's no periodic bootstrap/bootstrap_check
        self.manager._periodic_task.cancel()
        self.loop.run_until_complete(test())

    @patch('%s.BootstrapManager._schedule_next_bootstrap_check' % BM)
    @patch('%s.ServiceRegistry.get_bootstrap_rpc_channel' % BM)
    @patch('%s.cert_utils.write_cert' % BM)
    @patch('magma.common.cert_utils.write_key')
    def test__bootstrap_now(
        self,
        write_key_mock,
        write_cert_mock,
        bootstrap_channel_mock,
        schedule_mock,
    ):
        def fake_schedule():
            self.manager._loop.stop()

        async def test():
            make_awaitable(self.manager._bootstrap_success_cb)
            bootstrap_channel_mock.return_value = self.channel
            schedule_mock.side_effect = fake_schedule

            await self.manager._bootstrap_now()
            write_key_mock.assert_has_calls(
                [call(ANY, self.manager._gateway_key_file)],
            )
            write_cert_mock.assert_has_calls(
                [call(ANY, self.manager._gateway_cert_file)],
            )
            self.assertIs(self.manager._state, bm.BootstrapState.BOOTSTRAPPING)
            self.manager._bootstrap_success_cb.assert_has_calls([call(True)])

        # Cancel the loop so that there's no periodic bootstrap/bootstrap_check
        self.manager._periodic_task.cancel()
        self.loop.run_until_complete(test())

    @patch('%s.BootstrapManager._schedule_next_bootstrap' % BM)
    @patch('%s.ServiceRegistry.get_bootstrap_rpc_channel' % BM)
    def test__bootstrap_fail(
        self,
        bootstrap_channel_mock,
        schedule_next_bootstrap_mock,
    ):
        async def test():
            # test fail to get channel
            bootstrap_channel_mock.side_effect = ValueError
            await self.manager._bootstrap_now()
            schedule_next_bootstrap_mock.assert_has_calls(
                [call(hard_failure=False)],
            )
            # because retry is mocked, state should still be bootstrapping
            self.assertIs(self.manager._state, bm.BootstrapState.BOOTSTRAPPING)

        # Cancel the loop so that there's no periodic bootstrap/bootstrap_check
        self.manager._periodic_task.cancel()
        self.loop.run_until_complete(test())

    @patch('%s.ec.generate_private_key' % BM)
    @patch('%s.BootstrapManager._schedule_next_bootstrap' % BM)
    def test__get_challenge_done_pk_exception(
            self,
            schedule_next_bootstrap_mock,
            generate_pk_mock,
    ):
        async def test():
            future = MagicMock()
            # Private key generation returns error
            generate_pk_mock.side_effect = InternalError("", 0)
            await self.manager._get_challenge_done_success(future.result)
            schedule_next_bootstrap_mock.assert_has_calls(
                [call(hard_failure=True)],
            )
        # Cancel the loop so that there's no periodic bootstrap/bootstrap_check
        self.manager._periodic_task.cancel()
        self.loop.run_until_complete(test())

    @patch('%s.BootstrapManager._schedule_next_bootstrap' % BM)
    @patch('%s.BootstrapManager._request_sign' % BM)
    def test__get_challenge_done(
            self,
            request_sign_mock,
            schedule_next_bootstrap_mock,
    ):
        async def test():
            future = MagicMock()

            # GetChallenge returns error
            self.manager._get_challenge_done_fail(future.exception)
            schedule_next_bootstrap_mock.assert_has_calls(
                [call(hard_failure=False)],
            )

            # Fail to construct response
            schedule_next_bootstrap_mock.reset_mock()
            future.exception = lambda: None
            await self.manager._get_challenge_done_success(future.result)
            schedule_next_bootstrap_mock.assert_has_calls(
                [call(hard_failure=True)],
            )

            # No error
            schedule_next_bootstrap_mock.reset_mock()
            self.manager._loop = MagicMock()
            future.result = lambda: Challenge(
                challenge=b'simple_challenge',
                key_type=ChallengeKey.ECHO,
            )

            make_awaitable(request_sign_mock)

            await self.manager._get_challenge_done_success(future.result())
            schedule_next_bootstrap_mock.assert_not_called()
            request_sign_mock.assert_has_calls([call(ANY)])
        # Cancel the loop so that there's no periodic bootstrap/bootstrap_check
        self.manager._periodic_task.cancel()
        self.loop.run_until_complete(test())

    @patch('%s.BootstrapManager._schedule_next_bootstrap' % BM)
    @patch('%s.ServiceRegistry.get_bootstrap_rpc_channel' % BM)
    def test__request_sign_fail(
            self,
            bootstrap_channel_mock,
            schedule_next_bootstrap_mock,
    ):
        async def test():
            challenge = Challenge(
                challenge=b'simple_challenge',
                key_type=ChallengeKey.ECHO,
            )
            self.manager._gateway_key = ec.generate_private_key(
                ec.SECP256R1(), default_backend(),
            )
            csr = self.manager._create_csr()
            response = self.manager._construct_response(challenge, csr)

            # test fail to get channel
            bootstrap_channel_mock.side_effect = ValueError
            await self.manager._request_sign(response)
            schedule_next_bootstrap_mock.assert_has_calls(
                [call(hard_failure=False)],
            )

        # Cancel the loop so that there's no periodic bootstrap/bootstrap_check
        self.manager._periodic_task.cancel()
        self.loop.run_until_complete(test())

    @patch('%s.BootstrapManager._schedule_next_bootstrap' % BM)
    @patch('%s.ServiceRegistry.get_bootstrap_rpc_channel' % BM)
    def test__request_sign(
            self,
            bootstrap_channel_mock,
            schedule_next_bootstrap_mock,
    ):
        async def test():
            challenge = Challenge(
                challenge=b'simple_challenge',
                key_type=ChallengeKey.ECHO,
            )
            self.manager._gateway_key = ec.generate_private_key(
                ec.SECP256R1(), default_backend(),
            )
            csr = self.manager._create_csr()
            response = self.manager._construct_response(challenge, csr)
            # test no error
            schedule_next_bootstrap_mock.reset_mock()
            bootstrap_channel_mock.reset_mock()
            bootstrap_channel_mock.side_effect = None
            bootstrap_channel_mock.return_value = self.channel
            make_awaitable(self.manager._bootstrap_success_cb)
            await self.manager._request_sign(response)
            schedule_next_bootstrap_mock.assert_not_called()
        # Cancel the loop so that there's no periodic bootstrap/bootstrap_check
        self.manager._periodic_task.cancel()
        self.manager._loop.run_until_complete(test())

    @patch('%s.cert_utils.write_cert' % BM)
    @patch('%s.cert_utils.write_key' % BM)
    @patch('%s.BootstrapManager._schedule_next_bootstrap_check' % BM)
    @patch('%s.BootstrapManager._schedule_next_bootstrap' % BM)
    def test__request_sign_done(
        self,
        schedule_next_bootstrap_mock,
        schedule_bootstrap_check_mock,
        write_key_mock,
        write_cert_mock,
    ):
        async def test():
            future = MagicMock()

            # RequestSign returns error
            self.manager._request_sign_done_fail(future.exception)
            schedule_next_bootstrap_mock.assert_has_calls(
                [call(hard_failure=False)],
            )
            # certificate is invalid
            schedule_next_bootstrap_mock.reset_mock()
            not_before = \
                datetime.datetime.utcnow() + datetime.timedelta(hours=1)
            invalid_cert = create_cert_message(not_before=not_before)
            await self.manager._request_sign_done_success(invalid_cert)
            schedule_next_bootstrap_mock.assert_has_calls(
                [call(hard_failure=True)],
            )
            # certificate is valid
            schedule_next_bootstrap_mock.reset_mock()
            make_awaitable(self.manager._bootstrap_success_cb)
            valid_cert = create_cert_message()
            await self.manager._request_sign_done_success(valid_cert)
            self.manager._bootstrap_success_cb.assert_has_calls([call(True)])
            schedule_next_bootstrap_mock.assert_not_called()
            write_key_mock.assert_has_calls(
                [call(ANY, self.manager._gateway_key_file)],
            )
            write_cert_mock.assert_has_calls(
                [call(ANY, self.manager._gateway_cert_file)],
            )

            schedule_bootstrap_check_mock.assert_has_calls([call()])
        # Cancel the loop so that there's no periodic bootstrap/bootstrap_check
        self.manager._periodic_task.cancel()
        self.manager._loop.run_until_complete(test())

    def test__schedule_next_bootstrap(self):
        self.manager._loop = MagicMock()
        self.manager.LONG_BOOTSTRAP_RETRY_INTERVAL = datetime.timedelta(
            seconds=1,
        )
        self.manager.SHORT_BOOTSTRAP_RETRY_INTERVAL = datetime.timedelta(
            seconds=0,
        )
        self.manager._state = bm.BootstrapState.BOOTSTRAPPING

        self.manager._schedule_next_bootstrap(False)
        self.assertAlmostEqual(
            self.manager._interval,
            self.manager.SHORT_BOOTSTRAP_RETRY_INTERVAL.total_seconds(),
        )
        self.assertIs(
            self.manager._state,
            bm.BootstrapState.SCHEDULED_BOOTSTRAP,
        )

        self.manager._state = bm.BootstrapState.BOOTSTRAPPING
        self.manager._schedule_next_bootstrap(True)
        self.assertAlmostEqual(
            self.manager._interval,
            self.manager.LONG_BOOTSTRAP_RETRY_INTERVAL.total_seconds(),
        )
        self.assertIs(
            self.manager._state,
            bm.BootstrapState.SCHEDULED_BOOTSTRAP,
        )

    def test__schedule_next_bootstrap_check(self):
        self.manager._loop = MagicMock()
        self.manager._state = bm.BootstrapState.BOOTSTRAPPING
        self.manager._schedule_next_bootstrap_check()
        self.assertAlmostEqual(
            self.manager._interval,
            self.manager.PERIODIC_BOOTSTRAP_CHECK_INTERVAL.total_seconds(),
        )
        self.assertIs(self.manager._state, bm.BootstrapState.SCHEDULED_CHECK)

    def test__create_csr(self):
        self.manager._gateway_key = ec.generate_private_key(
            ec.SECP256R1(), default_backend(),
        )
        csr_msg = self.manager._create_csr()
        self.assertEqual(csr_msg.id.gateway.hardware_id, self.hw_id)

    @patch('magma.common.cert_utils.load_key')
    def test__construct_response(self, load_key_mock):
        ecdsa_key = ec.generate_private_key(ec.SECP256R1(), default_backend())

        key_types = {
            ChallengeKey.ECHO: None,
            ChallengeKey.SOFTWARE_ECDSA_SHA256: ecdsa_key,
        }
        for key_type, key in key_types.items():
            load_key_mock.return_value = key
            challenge = Challenge(key_type=key_type, challenge=b'challenge')
            response = self.manager._construct_response(challenge, CSR())
            self.assertEqual(response.hw_id.id, self.hw_id)
            self.assertEqual(response.challenge, challenge.challenge)

        challenge = Challenge(key_type=5, challenge=b'crap challenge')
        with self.assertRaises(
                bm.BootstrapError,
                msg='Unknown key type: %s' % challenge.key_type,
        ):
            self.manager._construct_response(challenge, CSR())

    @patch('magma.common.cert_utils.load_key')
    def test__ecdsa_sha256_response(self, load_key_mock):
        challenge = b'challenge'

        # success case
        private_key = ec.generate_private_key(
            ec.SECP256R1(), default_backend(),
        )
        load_key_mock.return_value = private_key
        r, s = self.manager._ecdsa_sha256_response(challenge)
        r = int.from_bytes(r, 'big')
        s = int.from_bytes(s, 'big')
        signature = encode_dss_signature(r, s)
        private_key.public_key().verify(
            signature, challenge, ec.ECDSA(hashes.SHA256()),
        )

        # no key found
        load_key_mock.reset_mock()
        load_key_mock.side_effect = IOError
        with self.assertRaises(bm.BootstrapError):
            self.manager._ecdsa_sha256_response(challenge)

        # wrong type of key, e.g. rsa
        load_key_mock.reset_mock()
        load_key_mock.return_value = rsa.generate_private_key(
            65537, 2048, default_backend(),
        )
        with self.assertRaises(
                bm.BootstrapError,
                msg='Challenge key cannot be used for ECDSA signature',
        ):
            self.manager._ecdsa_sha256_response(challenge)

    def test__is_valid_certificate(self):
        # not-yet-valid
        not_before = datetime.datetime.utcnow() + datetime.timedelta(hours=1)
        cert = create_cert_message(not_before=not_before)
        is_valid = self.manager._is_valid_certificate(cert)
        self.assertFalse(is_valid)

        # expiring soon
        with self.assertLogs() as log:
            not_before = datetime.datetime.utcnow()
            not_after = not_before + datetime.timedelta(hours=1)
            cert = create_cert_message(
                not_before=not_before,
                not_after=not_after,
            )
            is_valid = self.manager._is_valid_certificate(cert)
            self.assertTrue(is_valid)
            self.assertEqual(
                log.output,
                ['WARNING:root:Received a 1.0-hour certificate'],
            )

        # correct
        cert = create_cert_message()
        is_valid = self.manager._is_valid_certificate(cert)
        self.assertTrue(is_valid)


def create_cert(not_before, not_after):
    key = rsa.generate_private_key(65537, 2048, default_backend())

    subject = issuer = x509.Name([
        x509.NameAttribute(x509.oid.NameOID.COUNTRY_NAME, u"US"),
        x509.NameAttribute(x509.oid.NameOID.STATE_OR_PROVINCE_NAME, u"CA"),
        x509.NameAttribute(x509.oid.NameOID.LOCALITY_NAME, u"San Francisco"),
        x509.NameAttribute(x509.oid.NameOID.ORGANIZATION_NAME, u"My Company"),
        x509.NameAttribute(x509.oid.NameOID.COMMON_NAME, u"mysite.com"),
    ])

    cert = x509.CertificateBuilder().subject_name(
        subject,
    ).issuer_name(
        issuer,
    ).public_key(
        key.public_key(),
    ).serial_number(
        x509.random_serial_number(),
    ).not_valid_before(
        not_before,
    ).not_valid_after(
        not_after,
    ).sign(key, hashes.SHA256(), default_backend())

    return cert


def create_cert_message(not_before=None, not_after=None):
    if not_before is None:
        not_before = datetime.datetime.utcnow()
    if not_after is None:
        not_after = not_before + datetime.timedelta(days=10)

    cert = create_cert(not_before, not_after)

    not_before_stamp = Timestamp()
    not_before_stamp.FromDatetime(not_before)

    not_after_stamp = Timestamp()
    not_after_stamp.FromDatetime(not_after)

    dummy_cert = Certificate(
        cert_der=cert.public_bytes(serialization.Encoding.DER),
        not_before=not_before_stamp,
        not_after=not_after_stamp,
    )
    return dummy_cert
