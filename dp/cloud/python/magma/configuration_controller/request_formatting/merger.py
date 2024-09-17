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

from collections import defaultdict
from typing import Dict, List

from magma.configuration_controller.custom_types.custom_types import (
    DBRequest,
    MergedRequests,
)


def merge_requests(request_map: Dict[str, List[DBRequest]]) -> MergedRequests:
    """
    Merge request_map objects into one JSON object with request names as keys
    and sub-arrays of payload values

    Parameters:
        request_map: Map containing list of request objects at request type key

    Returns:
        merged_requests (MergedRequests): JSON serialized object with a subarray of payloads
    """

    merged_requests = defaultdict(list)

    for request_type in request_map.keys():
        payloads = [req.payload for req in request_map[request_type]]
        merged_requests[request_type].extend(payloads)

    return merged_requests
