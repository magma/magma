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
from typing import List, Optional

from dp.protos.requests_pb2 import RequestPayload
from magma.protocol_controller.grpc_client.grpc_client import GrpcClient


def upload_requests(client: GrpcClient, payload: str) -> Optional[List[int]]:
    """
    Upload requests to radio controller

    Parameters:
        client: a gRPC client
        payload: JSON string with request payload

    Returns:
        List[int]: IDs of uploaded requests
    """
    grpc_req = RequestPayload(payload=payload)
    grpc_response = client.UploadRequests(grpc_req)
    req_db_ids = grpc_response.ids
    logging.info(
        f'Executed upload_requests action. GRPC response: {req_db_ids}',
    )
    return req_db_ids
