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

import _thread
import socket
from wsgiref.simple_server import (
    ServerHandler,
    WSGIRequestHandler,
    WSGIServer,
    make_server,
)

from magma.common.misc_utils import get_ip_from_if
from magma.configuration.service_configs import load_service_config
from magma.enodebd.logger import EnodebdLogger as logger
from magma.enodebd.state_machines.enb_acs_manager import StateMachineManager
from spyne.server.wsgi import WsgiApplication

from .models import CWMP_NS
from .rpc_methods import AutoConfigServer
from .spyne_mods import Tr069Application, Tr069Soap11

# Socket timeout in seconds. Should be set larger than the longest TR-069
# response time (typically for a GetParameterValues of the entire data model),
# measured at 168secs. Should also be set smaller than ENB_CONNECTION_TIMEOUT,
# to avoid incorrectly detecting eNodeB timeout.
SOCKET_TIMEOUT = 240


class tr069_WSGIRequestHandler(WSGIRequestHandler):
    timeout = 10
    # pylint: disable=attribute-defined-outside-init

    def handle_single(self):
        """Handle a single HTTP request"""
        self.raw_requestline = self.rfile.readline(65537)
        if len(self.raw_requestline) > 65536:
            self.requestline = ''
            self.request_version = ''
            self.command = ''
            self.close_connection = 1
            self.send_error(414)
            return

        if not self.parse_request():  # An error code has been sent, just exit
            return

        handler = ServerHandler(
            self.rfile, self.wfile, self.get_stderr(), self.get_environ(),
        )
        handler.http_version = "1.1"
        handler.request_handler = self  # backpointer for logging

        # eNodeB will sometimes close connection to enodebd.
        # The cause of this is unknown, but we can safely ignore the
        # closed connection, and continue as normal otherwise.
        #
        # While this throws a BrokenPipe exception in wsgi server,
        # it also causes an AttributeError to be raised because of a
        # bug in the wsgi server.
        # https://bugs.python.org/issue27682
        try:
            handler.run(self.server.get_app())
        except BrokenPipeError:
            self.log_error(
                "eNodeB has unexpectedly closed the TCP connection.",
            )

    def handle(self):
        self.protocol_version = "HTTP/1.1"
        self.close_connection = 0

        try:
            while not self.close_connection:
                self.handle_single()

        except socket.timeout as e:
            self.log_error("tr069 WSGI Server Socket Timeout: %r", e)
            self.close_connection = 1
            return

        except socket.error as e:
            self.log_error("tr069 WSGI Server Socket Error: %r", e)
            self.close_connection = 1
            return

    # Disable pylint warning because we are using same parameter name as built-in
    # pylint: disable=redefined-builtin
    def log_message(self, format, *args):
        """ Overwrite message logging to use python logging framework rather
            than stderr """
        logger.debug("%s - %s", self.client_address[0], format % args)

    # Disable pylint warning because we are using same parameter name as built-in
    # pylint: disable=redefined-builtin
    def log_error(self, format, *args):
        """ Overwrite message logging to use python logging framework rather
            than stderr """
        logger.warning("%s - %s", self.client_address[0], format % args)


def tr069_server(state_machine_manager: StateMachineManager) -> None:
    """
    TR-069 server
    Inputs:
        - acs_to_cpe_queue = instance of Queue
            containing messages from parent process/thread to be sent to CPE
        - cpe_to_acs_queue = instance of Queue
            containing messages from CPE to be sent to parent process/thread
    """
    config = load_service_config("enodebd")

    AutoConfigServer.set_state_machine_manager(state_machine_manager)

    app = Tr069Application(
        [AutoConfigServer], CWMP_NS,
        in_protocol=Tr069Soap11(validator='soft'),
        out_protocol=Tr069Soap11(),
    )
    wsgi_app = WsgiApplication(app)

    try:
        ip_address = get_ip_from_if(config['tr069']['interface'])
    except (ValueError, KeyError) as e:
        # Interrupt main thread since process should not continue without TR-069
        _thread.interrupt_main()
        raise e

    socket.setdefaulttimeout(SOCKET_TIMEOUT)
    logger.info(
        'Starting TR-069 server on %s:%s',
        ip_address, config['tr069']['port'],
    )
    server = make_server(
        ip_address,
        config['tr069']['port'], wsgi_app,
        WSGIServer, tr069_WSGIRequestHandler,
    )

    # Note: use single-thread server, to avoid state contention
    try:
        server.serve_forever()
    finally:
        # Log error and interrupt main thread, to ensure that entire process
        # is restarted if this thread exits
        logger.error('Hit error in TR-069 thread. Interrupting main thread.')
        _thread.interrupt_main()
