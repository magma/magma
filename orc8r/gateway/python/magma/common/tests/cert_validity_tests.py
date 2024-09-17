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
import errno
import ssl
from unittest import TestCase
from unittest.mock import MagicMock, patch

import magma.common.cert_validity as cv


# https://stackoverflow.com/questions/32480108/mocking-async-call-in-python-3-5
def AsyncMock():
    coro = MagicMock(name="CoroutineResult")
    corofunc = MagicMock(
        name="CoroutineFunction",
        side_effect=asyncio.coroutine(coro),
    )
    corofunc.coro = coro
    return corofunc


class CertValidityTests(TestCase):
    def setUp(self):
        self.host = 'localhost'
        self.port = 8080
        self.certfile = 'certfile'
        self.keyfile = 'keyfile'

        asyncio.set_event_loop(None)
        self.loop = asyncio.new_event_loop()

    def test_tcp_connection(self):
        """
        Test that loop.create_connection called with the correct TCP args.
        """
        self.loop.create_connection = MagicMock()

        @asyncio.coroutine
        def go():
            yield from cv.create_tcp_connection(
                self.host,
                self.port,
                self.loop,
            )
        self.loop.run_until_complete(go())

        self.loop.create_connection.assert_called_once_with(
            cv.TCPClientProtocol,
            self.host,
            self.port,
        )

    @patch('magma.common.cert_validity.ssl.SSLContext')
    def test_ssl_connection(self, mock_ssl):
        """
        Test that ssl.SSLContext and loop.create_connection are called with the
        correct SSL args.
        """
        self.loop.create_connection = MagicMock()

        @asyncio.coroutine
        def go():
            yield from cv.create_ssl_connection(
                self.host,
                self.port,
                self.certfile,
                self.keyfile,
                self.loop,
            )
        self.loop.run_until_complete(go())

        mock_context = mock_ssl.return_value

        mock_ssl.assert_called_once_with(ssl.PROTOCOL_SSLv23)
        mock_context.load_cert_chain.assert_called_once_with(
            self.certfile,
            keyfile=self.keyfile,
        )

        self.loop.create_connection.assert_called_once_with(
            cv.TCPClientProtocol,
            self.host,
            self.port,
            ssl=mock_context,
        )

    @patch(
        'magma.common.cert_validity.create_ssl_connection',
        new_callable=AsyncMock,
    )
    @patch(
        'magma.common.cert_validity.create_tcp_connection',
        new_callable=AsyncMock,
    )
    def test_cert_is_invalid_both_ok(self, mock_create_tcp, mock_create_ssl):
        """
        Test the appropriate calls and return value for cert_is_invalid()
        cert_is_invalid() == False when TCP and SSL succeed
        """

        @asyncio.coroutine
        def go():
            return (
                yield from cv.cert_is_invalid(
                    self.host,
                    self.port,
                    self.certfile,
                    self.keyfile,
                    self.loop,
                )
            )
        ret_val = self.loop.run_until_complete(go())

        mock_create_tcp.assert_called_once_with(
            self.host,
            self.port,
            self.loop,
        )
        mock_create_ssl.assert_called_once_with(
            self.host,
            self.port,
            self.certfile,
            self.keyfile,
            self.loop,
        )
        self.assertEqual(ret_val, False)

    @patch(
        'magma.common.cert_validity.create_ssl_connection',
        new_callable=AsyncMock,
    )
    @patch('magma.common.cert_validity.create_tcp_connection', AsyncMock())
    def test_cert_is_invalid_ssl_fail(self, mock_create_ssl):
        """
        Test cert_is_invalid() == True when TCP succeeds and SSL fails
        """

        mock_err = TimeoutError()
        mock_err.errno = errno.ETIMEDOUT
        mock_create_ssl.coro.side_effect = mock_err

        @asyncio.coroutine
        def go():
            return (
                yield from cv.cert_is_invalid(
                    self.host,
                    self.port,
                    self.certfile,
                    self.keyfile,
                    self.loop,
                )
            )
        ret_val = self.loop.run_until_complete(go())
        self.assertEqual(ret_val, True)

    @patch(
        'magma.common.cert_validity.create_ssl_connection',
        new_callable=AsyncMock,
    )
    @patch('magma.common.cert_validity.create_tcp_connection', AsyncMock())
    def test_cert_is_invalid_ssl_fail_none_errno(self, mock_create_ssl):
        """
        Test cert_is_invalid() == True when TCP succeeds and SSL fails w/o error number
        """

        mock_err = TimeoutError()
        mock_err.errno = None
        mock_create_ssl.coro.side_effect = mock_err

        @asyncio.coroutine
        def go():
            return (
                yield from cv.cert_is_invalid(
                    self.host,
                    self.port,
                    self.certfile,
                    self.keyfile,
                    self.loop,
                )
            )
        ret_val = self.loop.run_until_complete(go())
        self.assertEqual(ret_val, True)

    @patch('magma.common.cert_validity.create_ssl_connection', AsyncMock())
    @patch(
        'magma.common.cert_validity.create_tcp_connection',
        new_callable=AsyncMock,
    )
    def test_cert_is_invalid_tcp_fail_none_errno(self, mock_create_tcp):
        """
        Test cert_is_invalid() == False when TCP fails w/o errno and SSL succeeds
        """

        mock_err = TimeoutError()
        mock_err.errno = None
        mock_create_tcp.coro.side_effect = mock_err

        @asyncio.coroutine
        def go():
            return (
                yield from cv.cert_is_invalid(
                    self.host,
                    self.port,
                    self.certfile,
                    self.keyfile,
                    self.loop,
                )
            )
        ret_val = self.loop.run_until_complete(go())
        self.assertEqual(ret_val, False)

    @patch('magma.common.cert_validity.create_ssl_connection', AsyncMock())
    @patch(
        'magma.common.cert_validity.create_tcp_connection',
        new_callable=AsyncMock,
    )
    def test_cert_is_invalid_tcp_fail(self, mock_create_tcp):
        """
        Test cert_is_invalid() == False when TCP fails and SSL succeeds
        """

        mock_err = TimeoutError()
        mock_err.errno = errno.ETIMEDOUT
        mock_create_tcp.coro.side_effect = mock_err

        @asyncio.coroutine
        def go():
            return (
                yield from cv.cert_is_invalid(
                    self.host,
                    self.port,
                    self.certfile,
                    self.keyfile,
                    self.loop,
                )
            )
        ret_val = self.loop.run_until_complete(go())
        self.assertEqual(ret_val, False)

    @patch(
        'magma.common.cert_validity.create_ssl_connection',
        new_callable=AsyncMock,
    )
    @patch(
        'magma.common.cert_validity.create_tcp_connection',
        new_callable=AsyncMock,
    )
    def test_cert_is_invalid_both_fail(self, mock_create_tcp, mock_create_ssl):
        """
        Test cert_is_invalid() == False when TCP and SSL fail
        """

        mock_tcp_err = TimeoutError()
        mock_tcp_err.errno = errno.ETIMEDOUT
        mock_create_tcp.coro.side_effect = mock_tcp_err

        mock_ssl_err = TimeoutError()
        mock_ssl_err.errno = errno.ETIMEDOUT
        mock_create_ssl.coro.side_effect = mock_ssl_err

        @asyncio.coroutine
        def go():
            return (
                yield from cv.cert_is_invalid(
                    self.host,
                    self.port,
                    self.certfile,
                    self.keyfile,
                    self.loop,
                )
            )
        ret_val = self.loop.run_until_complete(go())
        self.assertEqual(ret_val, False)
