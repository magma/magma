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

# pylint: disable=W0223

# We cannot create instances directly for transport we need them to be
# created by event loop so hacking this.


class MockTransport(asyncio.Transport):
    def __init__(self, extra=None):
        self.sent = []
        self.open = True
        super().__init__()

    def write(self, data):
        if self.open:
            self.sent.append(data)

    def close(self):
        self.open = False

    def flush(self):
        while self.sent:
            self.sent.pop()

    def get_extra_info(self, name, default=None):
        return self.extra.get(name, default)

    def is_closing(self):
        return self.open is False
