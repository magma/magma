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

import threading
import unittest
from http.server import BaseHTTPRequestHandler, HTTPServer
from urllib import parse

from integ_tests.cloud.cloud_manager import CloudManager
from integ_tests.cloud.fixtures import GATEWAY_ID, NETWORK_ID
from integ_tests.gateway.rpc import get_gateway_hw_id


class TestHTTPServerRequestHandler(BaseHTTPRequestHandler):

    def do_POST(self):
        print("Received metric export POST request")
        # Get post body
        request_headers = self.headers
        length = int(request_headers.get_all('content-length')[0])
        post_body = self.rfile.read(length)
        post_body_dict = parse.parse_qs(
            parse.unquote(post_body.decode('utf-8')),
        )

        # Sanity check request, make sure it has the key 'datapoints'
        assert(len(post_body_dict['datapoints'][0]) > 0)
        print("Metrics export valid")

        # Send success response to cloud
        self.send_response(200)
        self.send_header('content-type', 'application/json')
        self.end_headers()
        self.wfile.write(bytes('"success"', 'utf-8'))
        return


class TestMetricsExport(unittest.TestCase):
    """
    Runs a test case which starts a mock metrics server (in this case, ODS)
    on a set IP and port and waits for gateway metrics to export via the cloud.
    """

    TEST_VM_IP = '192.168.60.141'
    TEST_VM_PORT = 8081
    METRIC_TIMEOUT = 120  # 2 minutes

    def setUp(self):
        self._cloud_manager = CloudManager()
        self._cloud_manager.delete_networks([NETWORK_ID])

        self._cloud_manager.create_network(NETWORK_ID)
        self._cloud_manager.register_gateway(
            NETWORK_ID, GATEWAY_ID,
            get_gateway_hw_id(),
        )

        self._test_server = HTTPServer(
            (self.TEST_VM_IP, self.TEST_VM_PORT),
            TestHTTPServerRequestHandler,
        )
        self._server_thread = threading.Thread(target=self.run_server)
        self._server_thread.daemon = True

    def tearDown(self):
        self._test_server.socket.close()
        self._cloud_manager.clean_up()

    def handle_timeout(self):
        self.assertTrue(
            False,
            "Metrics not received before timeout, test failed",
        )

    def run_server(self):
        self._test_server.timeout = self.METRIC_TIMEOUT
        self._test_server.handle_timeout = self.handle_timeout
        self._test_server.handle_request()

    def test_metrics_export(self):
        print("Starting test server, waiting for metrics export...")
        self._server_thread.start()
        self._server_thread.join()
        print("Metrics exported successfully")


if __name__ == "__main__":
    unittest.main()
