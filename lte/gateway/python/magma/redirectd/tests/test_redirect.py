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

import unittest

from lte.protos.policydb_pb2 import RedirectInformation
from magma.redirectd.redirect_server import (
    HTTP_NOT_FOUND,
    HTTP_REDIRECT,
    NOT_FOUND_HTML,
    ServerResponse,
    setup_flask_server,
)


class RedirectdTest(unittest.TestCase):
    def setUp(self):
        """
        Sets up a test version of the redirect server, mocks url_dict
        """
        app = setup_flask_server()
        app.config['TESTING'] = True

        test_dict = {
            '192.5.82.1':
                RedirectInformation(
                    support=1,
                    address_type=2,
                    server_address='http://www.example.com/',
                ),
        }

        def get_resp(src_ip):
            if src_ip not in test_dict:
                return ServerResponse(NOT_FOUND_HTML, HTTP_NOT_FOUND)
            return ServerResponse(
                test_dict[src_ip].server_address, HTTP_REDIRECT,
            )
        # Replaces all url_dict polls with a mocked dict (for all url rules)
        # pylint: disable=protected-access
        for rule in app.url_map._rules:
            if rule is not None and rule.defaults is not None:
                rule.defaults['get_redirect_response'] = get_resp
        self.client = app.test_client()

    def test_302_homepage(self):
        """
        Assert 302 http response, proper reponse headers with new dest url
        """
        resp = self.client.get('/', environ_base={'REMOTE_ADDR': '192.5.82.1'})

        self.assertEqual(resp.status_code, HTTP_REDIRECT)
        self.assertEqual(resp.headers['Location'], 'http://www.example.com/')

    def test_302_with_path(self):
        """
        Assert 302 http response, proper reponse headers with new dest url
        """
        resp = self.client.get(
            '/generate_204',
            environ_base={'REMOTE_ADDR': '192.5.82.1'},
        )

        self.assertEqual(resp.status_code, HTTP_REDIRECT)
        self.assertEqual(resp.headers['Location'], 'http://www.example.com/')

    def test_404(self):
        """
        Assert 404 http response
        """
        resp = self.client.get('/', environ_base={'REMOTE_ADDR': '127.0.0.1'})

        self.assertEqual(resp.status_code, HTTP_NOT_FOUND)
