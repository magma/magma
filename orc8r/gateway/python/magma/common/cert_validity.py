"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Util module to distinguish between the reasons checkins stop working: network
is down or cert is invalid.
"""

import asyncio
import logging
import os
import ssl


class TCPClientProtocol(asyncio.Protocol):
    """
    Implementation of TCP Protocol to create and immediately close the
    connection
    """

    def connection_made(self, transport):
        transport.close()


@asyncio.coroutine
def create_tcp_connection(host, port, loop):
    """
    Creates tcp connection
    """
    tcp_conn = yield from loop.create_connection(
        TCPClientProtocol,
        host,
        port,
    )
    return tcp_conn


@asyncio.coroutine
def create_ssl_connection(host, port, certfile, keyfile, loop):
    """
    Creates ssl connection.
    """
    context = ssl.SSLContext(ssl.PROTOCOL_SSLv23)
    context.load_cert_chain(
        certfile,
        keyfile=keyfile,
    )

    ssl_conn = yield from loop.create_connection(
        TCPClientProtocol,
        host,
        port,
        ssl=context,
    )
    return ssl_conn


@asyncio.coroutine
def cert_is_invalid(host, port, certfile, keyfile, loop):
    """
    Asynchronously test if both a TCP and SSL connection can be made to host
    on port. If the TCP connection is successful, but the SSL connection fails,
    we assume this is due to an invalid cert.

    Args:
        host: host to connect to
        port: port to connect to on host
        certfile: path to a PEM encoded certificate
        keyfile: path to the corresponding key to the certificate
        loop: asyncio event loop
    Returns:
        True if the cert is invalid
        False otherwise
    """
    # Create connections
    tcp_coro = create_tcp_connection(host, port, loop)
    ssl_coro = create_ssl_connection(host, port, certfile, keyfile, loop)

    coros = tcp_coro, ssl_coro
    asyncio.set_event_loop(loop)
    res = yield from asyncio.gather(*coros, return_exceptions=True)
    tcp_res, ssl_res = res

    if isinstance(tcp_res, Exception):
        logging.error(
            'Error making TCP connection: %s, %s',
            'errno==None' if tcp_res.errno is None
            else os.strerror(tcp_res.errno),
            tcp_res,
        )
        return False

    # Invalid cert only when tcp succeeds and ssl fails
    if isinstance(ssl_res, Exception):
        logging.error(
            'Error making SSL connection: %s, %s',
            'errno==None' if ssl_res.errno is None
            else os.strerror(ssl_res.errno),
            ssl_res,
        )
        return True

    return False
