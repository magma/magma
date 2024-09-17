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

from unittest import mock

from spyne.server.wsgi import WsgiMethodContext


def get_spyne_context_with_ip(
    req_ip: str = "192.168.60.145",
) -> WsgiMethodContext:
    with mock.patch('spyne.server.wsgi.WsgiApplication') as MockTransport:
        MockTransport.req_env = {"REMOTE_ADDR": req_ip}
        with mock.patch('spyne.server.wsgi.WsgiMethodContext') as MockContext:
            MockContext.transport = MockTransport
            return MockContext
