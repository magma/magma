"""
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""
from unittest import TestCase

import grpc
import requests
from dp.protos.cbsd_pb2_grpc import CbsdManagementStub
from magma.test_runner.config import TestConfig

config = TestConfig()


class Orc8rIntegrationTestCase(TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        super().setUpClass()
        cls.maxDiff = None
        grpc_channel = grpc.insecure_channel(
            f"{config.ORC8R_DP_GRPC_SERVICE}:{config.ORC8R_DP_GRPC_PORT}",
        )
        cls.orc8r_dp_client = CbsdManagementStub(grpc_channel)

    @classmethod
    def tearDownClass(cls) -> None:
        _delete_dp_elasticsearch_indices()


def _delete_dp_elasticsearch_indices() -> None:
    requests.delete(f"{config.ELASTICSEARCH_URL}/{config.ELASTICSEARCH_INDEX}*")
