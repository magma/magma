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
# pylint: disable=broad-except

import datetime
import enum
import logging
import os

import grpc
import magma.common.cert_utils as cert_utils
import snowflake
from cryptography.exceptions import InternalError
from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives import hashes, serialization
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.hazmat.primitives.asymmetric.utils import decode_dss_signature
from google.protobuf.duration_pb2 import Duration
from magma.common.rpc_utils import grpc_async_wrapper
from magma.common.sdwatchdog import SDWatchdogTask
from magma.common.service_registry import ServiceRegistry
from magma.configuration.service_configs import load_service_config
from magma.magmad.metrics import BOOTSTRAP_EXCEPTION
from orc8r.protos.bootstrapper_pb2 import ChallengeKey, Response
from orc8r.protos.bootstrapper_pb2_grpc import BootstrapperStub
from orc8r.protos.certifier_pb2 import CSR
from orc8r.protos.identity_pb2 import AccessGatewayID, Identity


class BootstrapError(Exception):
    pass


@enum.unique
class BootstrapState(enum.Enum):
    INITIAL = 0
    BOOTSTRAPPING = 1
    SCHEDULED_BOOTSTRAP = 2
    SCHEDULED_CHECK = 3
    IDLE = 4


class BootstrapManager(SDWatchdogTask):
    """
    Bootstrap the gateway by contacting the controller.

    Bootstrap manager responds to the challenge from the controller to
    verify the device. As a result of the bootstrap process, the
    gateways' session certs would be written to /var/opt/magma/certs.
    Before the session certs expire, bootstrap would make sure we
    fetch new certs by maintaining a timer internally.
    """
    # delay in asyncio should not exceed one day
    PERIODIC_BOOTSTRAP_CHECK_INTERVAL = datetime.timedelta(hours=1)
    PREEXPIRY_BOOTSTRAP_INTERVAL = datetime.timedelta(hours=20)
    SHORT_BOOTSTRAP_RETRY_INTERVAL = datetime.timedelta(seconds=30)
    LONG_BOOTSTRAP_RETRY_INTERVAL = datetime.timedelta(minutes=1)

    def __init__(self, service, bootstrap_success_cb):
        super().__init__(
            self.PERIODIC_BOOTSTRAP_CHECK_INTERVAL.total_seconds(),
            service.loop,
        )

        control_proxy_config = load_service_config('control_proxy')

        self._challenge_key_file \
            = service.config['bootstrap_config']['challenge_key']
        self._hw_id = snowflake.snowflake()
        self._gateway_key_file = control_proxy_config['gateway_key']
        self._gateway_cert_file = control_proxy_config['gateway_cert']
        self._gateway_key = None
        self._state = BootstrapState.INITIAL
        self._bootstrap_success_cb = bootstrap_success_cb

        # give some margin on watchdog check interval
        self.set_timeout(self._interval * 1.1)

    def start_bootstrap_manager(self):
        self.start()
        self._maybe_create_challenge_key()

    def stop_bootstrap_manager(self):
        self._state = BootstrapState.IDLE
        self.stop()

    async def _run(self):
        if self._state == BootstrapState.INITIAL:
            await self._bootstrap_check()
        elif self._state == BootstrapState.BOOTSTRAPPING:
            pass
        elif self._state == BootstrapState.SCHEDULED_BOOTSTRAP:
            await self._bootstrap_now()
        elif self._state == BootstrapState.SCHEDULED_CHECK:
            await self._bootstrap_check()
        elif self._state == BootstrapState.IDLE:
            pass

    async def schedule_bootstrap_now(self):
        """Public Interface to start a bootstrap

        1. If the device is already bootstrapping, do nothing
        2. If it is waiting for a next bootstrap check or bootstrap, wake up
           and do it now.
        """
        if self._state is BootstrapState.BOOTSTRAPPING:
            return
        await self.wake_up()

    def _maybe_create_challenge_key(self):
        """Generate key the first time it runs if key does not exist"""
        if not os.path.exists(self._challenge_key_file):
            logging.info(
                'Generating challenge key and written into %s',
                self._challenge_key_file,
            )
            challenge_key = ec.generate_private_key(
                ec.SECP384R1(), default_backend(),
            )
            cert_utils.write_key(challenge_key, self._challenge_key_file)

    async def _bootstrap_check(self):
        """Check whether bootstrap is need

        Check whether cert is present and still valid
        If so, a future _bootstrap_check will be scheduled.
        Otherwise _bootstrap_now will be called immediately
        """
        # flag to ensure the loop is still running, successfully or not
        self.heartbeat()

        try:
            cert = cert_utils.load_cert(self._gateway_cert_file)
        except (IOError, ValueError):
            logging.info('Cannot load a proper cert, start bootstrapping')
            await self._bootstrap_now()
            return

        now = datetime.datetime.utcnow()
        if now + self.PREEXPIRY_BOOTSTRAP_INTERVAL > cert.not_valid_after:
            logging.info(
                'Certificate is expiring soon at %s, start bootstrapping',
                cert.not_valid_after,
            )
            await self._bootstrap_now()
            return
        if now < cert.not_valid_before:
            logging.error(
                'Certificate is not valid until %s', cert.not_valid_before,
            )
            await self._bootstrap_now()
            return

        # no need to restart control_proxy
        await self._bootstrap_success_cb(False)
        self._schedule_next_bootstrap_check()

    async def _bootstrap_now(self):
        """Main entrance to bootstrapping

        1. set self._state to BOOTSTRAPPING
        2. set up a gPRC channel and get a challenge (async)
        3. call _get_challenge_done_success  to deal with the response
        If any steps fails, a new _bootstrap_now call will be scheduled.
        """
        assert self._state != BootstrapState.BOOTSTRAPPING, \
            'At most one bootstrap is happening'
        self._state = BootstrapState.BOOTSTRAPPING

        try:
            chan = ServiceRegistry.get_bootstrap_rpc_channel()
        except ValueError as exp:
            logging.error('Failed to get rpc channel: %s', exp)
            self._schedule_next_bootstrap(hard_failure=False)
            return

        client = BootstrapperStub(chan)
        try:
            result = await grpc_async_wrapper(
                client.GetChallenge.future(AccessGatewayID(id=self._hw_id)),
                self._loop,
            )
            await self._get_challenge_done_success(result)

        except grpc.RpcError as err:
            self._get_challenge_done_fail(err)

    async def _get_challenge_done_success(self, challenge):
        # create key
        try:
            # GRPC python client only supports P256 elliptic curve cipher
            # See https://github.com/grpc/grpc/issues/23235
            # Behind the nghttpx control_proxy this isn't a problem because
            # nghttpx handles the handshake, but if you have a P384 cert and
            # don't proxy your cloud connections, every authenticated Python
            # GRPC call will fail.
            self._gateway_key = ec.generate_private_key(
                ec.SECP256R1(),
                default_backend(),
            )
        except InternalError as exp:
            logging.error('Fail to generate private key: %s', exp)
            BOOTSTRAP_EXCEPTION.labels(
                cause='GetChallengeDonePrivateKey',
            ).inc()
            self._schedule_next_bootstrap(hard_failure=True)
            return
        # create csr and send for signing
        try:
            csr = self._create_csr()
        except Exception as exp:
            logging.error('Fail to create csr: %s', exp)
            BOOTSTRAP_EXCEPTION.labels(
                cause='GetChallengeDoneCreateCSR:%s' % type(
                    exp,
                ).__name__,
            ).inc()

        try:
            response = self._construct_response(challenge, csr)
        except BootstrapError as exp:
            logging.error('Fail to create response: %s', exp)
            BOOTSTRAP_EXCEPTION.labels(
                cause='GetChallengeDoneCreateResponse',
            ).inc()
            self._schedule_next_bootstrap(hard_failure=True)
            return
        await self._request_sign(response)

    def _get_challenge_done_fail(self, err):
        err = 'GetChallenge error! [%s] %s' % (err.code(), err.details())
        logging.error(err)
        BOOTSTRAP_EXCEPTION.labels(cause='GetChallengeResp').inc()
        self._schedule_next_bootstrap(hard_failure=False)

    async def _request_sign(self, response):
        """Request a signed certificate

        set up a gPRC channel and set the response

        If it fails, schedule the next bootstrap,
        Otherwise _request_sign_done callback is called
        """
        try:
            chan = ServiceRegistry.get_bootstrap_rpc_channel()
        except ValueError as exp:
            logging.error('Failed to get rpc channel: %s', exp)
            BOOTSTRAP_EXCEPTION.labels(cause='RequestSignGetRPC').inc()
            self._schedule_next_bootstrap(hard_failure=False)
            return

        try:
            client = BootstrapperStub(chan)
            result = await grpc_async_wrapper(
                client.RequestSign.future(response),
                self._loop,
            )
            await self._request_sign_done_success(result)

        except grpc.RpcError as err:
            self._request_sign_done_fail(err)

    async def _request_sign_done_success(self, cert):
        if not self._is_valid_certificate(cert):
            BOOTSTRAP_EXCEPTION.labels(
                cause='RequestSignDoneInvalidCert',
            ).inc()
            self._schedule_next_bootstrap(hard_failure=True)
            return
        try:
            cert_utils.write_key(self._gateway_key, self._gateway_key_file)
            cert_utils.write_cert(cert.cert_der, self._gateway_cert_file)
        except Exception as exp:
            BOOTSTRAP_EXCEPTION.labels(
                cause='RequestSignDoneWriteCert:%s' % type(exp).__name__,
            ).inc()
            logging.error('Failed to write cert: %s', exp)

        # need to restart control_proxy
        await self._bootstrap_success_cb(True)
        self._gateway_key = None
        self._schedule_next_bootstrap_check()
        logging.info("Bootstrapped Successfully!")

    def _request_sign_done_fail(self, err):
        err = 'RequestSign error! [%s], %s' % (err.code(), err.details())
        BOOTSTRAP_EXCEPTION.labels(cause='RequestSignDoneResp').inc()
        logging.error(err)
        self._schedule_next_bootstrap(hard_failure=False)

    def _schedule_next_bootstrap(self, hard_failure):
        """Schedule a bootstrap

        Args:
            hard_failure: bool. If set, the time to next retry will be longer
        """
        if hard_failure:
            interval = self.LONG_BOOTSTRAP_RETRY_INTERVAL.total_seconds()
        else:
            interval = self.SHORT_BOOTSTRAP_RETRY_INTERVAL.total_seconds()
        logging.info('Retrying bootstrap in %d seconds', interval)
        self.set_interval(int(interval))
        self._state = BootstrapState.SCHEDULED_BOOTSTRAP

    def _schedule_next_bootstrap_check(self):
        """Schedule a bootstrap_check"""
        self.set_interval(
            int(self.PERIODIC_BOOTSTRAP_CHECK_INTERVAL.total_seconds()),
        )
        self._state = BootstrapState.SCHEDULED_CHECK

    def _create_csr(self):
        """Create CSR protobuf

        Returns:
             CSR protobuf object
        """
        csr = cert_utils.create_csr(self._gateway_key, self._hw_id)
        duration = Duration()
        duration.FromTimedelta(datetime.timedelta(days=4))
        csr = CSR(
            id=Identity(gateway=Identity.Gateway(hardware_id=self._hw_id)),
            valid_time=duration,
            csr_der=csr.public_bytes(serialization.Encoding.DER),
        )
        return csr

    def _construct_response(self, challenge, csr):
        """Construct response message given challenge and csr message

        Args:
            challenge: Challenge(key_type, challenge)
            csr: CSR object returned by create_csr

        Returns:
             protobuf Response object

        Raises:
            BootstrapError: Unknown key type, cannot load challenge key,
             or wrong type of challenge key
        """
        if challenge.key_type == ChallengeKey.ECHO:
            echo_resp = Response.Echo(
                response=challenge.challenge,
            )
            response = Response(
                hw_id=AccessGatewayID(id=self._hw_id),
                challenge=challenge.challenge,
                echo_response=echo_resp,
                csr=csr,
            )
        elif challenge.key_type == ChallengeKey.SOFTWARE_ECDSA_SHA256:
            r_bytes, s_bytes = self._ecdsa_sha256_response(challenge.challenge)
            ecdsa_resp = Response.ECDSA(r=r_bytes, s=s_bytes)
            response = Response(
                hw_id=AccessGatewayID(id=self._hw_id),
                challenge=challenge.challenge,
                ecdsa_response=ecdsa_resp,
                csr=csr,
            )
        else:
            raise BootstrapError('Unknown key type: %s' % challenge.key_type)
        return response

    def _is_valid_certificate(self, cert):
        """Check whether certificate is usable

        Args:
            cert: Certificate object returned by RequestSign gRPC call

        Returns:
            err: error message, None if no error
        """
        now = datetime.datetime.utcnow()
        not_before = cert.not_before.ToDatetime()
        if now < not_before:
            logging.error(
                'Current system time indicates certificate received is not yet valid (notBefore: %s). Consider checking NTP.', not_before,
            )
            return False

        not_after = cert.not_after.ToDatetime()
        valid_time = not_after - now
        # log a warning if the cert is short-lived
        if valid_time < self.PREEXPIRY_BOOTSTRAP_INTERVAL:
            valid_hours = valid_time.total_seconds() / 3600
            logging.warning('Received a %.1f-hour certificate', valid_hours)

        return True

    def _ecdsa_sha256_response(self, challenge):
        """Compute the ecdsa signature

        Args:
            challenge: content of challenge in bytes

        Returns:
            r_bytes, s_bytes: ecdsa signature R, S in bytes

        Raises:
            BootstrapError: if the gateway cannot be properly loaded
        """
        try:
            challenge_key = cert_utils.load_key(self._challenge_key_file)
        except (IOError, ValueError, TypeError) as e:
            raise BootstrapError(
                'Gateway does not have a proper challenge key: %s' % e,
            )

        try:
            signature = challenge_key.sign(
                challenge, ec.ECDSA(hashes.SHA256()),
            )
        except TypeError:
            raise BootstrapError(
                'Challenge key cannot be used for ECDSA signature',
            )

        r_int, s_int = decode_dss_signature(signature)
        r_bytes = r_int.to_bytes((r_int.bit_length() + 7) // 8, 'big')
        s_bytes = s_int.to_bytes((s_int.bit_length() + 7) // 8, 'big')
        return r_bytes, s_bytes
