"""
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import logging

from dp.protos.requests_pb2_grpc import RadioControllerStub


class GrpcClient(RadioControllerStub):
    """
    Basic gRPC Client class
    """

    def init_app(self, app):
        """
        Initialize Flask application

        Parameters:
            app: Flask application
        """

        logging.info("Initializing GRPC Client")
        app.extensions = getattr(app, "extensions", {})
        app.extensions[self.__class__.__name__] = self
