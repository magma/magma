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

from prometheus_client import Histogram

# Metrics for current configuration controller status
SAS_REQUEST_PROCESSING_TIME = Histogram(
    'dp_cc_request_processing_seconds',
    'Time spent processing a SAS request',
)

SAS_RESPONSE_PROCESSING_TIME = Histogram(
    'dp_cc_response_processing_seconds',
    'Time spent processing a SAS response',
)

PENDING_REQUESTS_FETCH_TIME = Histogram(
    'dp_cc_pending_requests_fetching_seconds',
    'Time spent fetching pending requests from the database',
)
