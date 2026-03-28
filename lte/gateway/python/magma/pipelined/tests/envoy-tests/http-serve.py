#!/usr/bin/env python

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


import SimpleHTTPServer
import SocketServer

PORT = 80


class GetHandler(SimpleHTTPServer.SimpleHTTPRequestHandler):

    def do_GET(self):
        self.send_head()
        for h in self.headers:
            self.send_header(h, self.headers[h])
        self.end_headers()
        self.send_response(200, "")


Handler = GetHandler
httpd = SocketServer.TCPServer(("192.168.128.1", PORT), Handler)
print("starting server")
httpd.serve_forever()
