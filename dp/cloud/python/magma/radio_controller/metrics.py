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
GRPC_REQUEST_PROCESSING_TIME = Histogram(
    'dp_rc_grpc_request_processing_seconds',
    'Time spent processing a GRPC request',
    ('name',),
)

GET_DB_STATE_PROCESSING_TIME = GRPC_REQUEST_PROCESSING_TIME.labels('get_database_state')
DELETE_CBSD_PROCESSING_TIME = GRPC_REQUEST_PROCESSING_TIME.labels('delete_cbsd')
ACKNOWLEDGE_UPDATE_PROCESSING_TIME = GRPC_REQUEST_PROCESSING_TIME.labels('acknowledge_cbsd_update')
ACKNOWLEDGE_RELINQUISH_PROCESSING_TIME = GRPC_REQUEST_PROCESSING_TIME.labels('acknowledge_cbsd_relinquish')
INSERT_TO_DB_PROCESSING_TIME = GRPC_REQUEST_PROCESSING_TIME.labels('insert_requests_to_db')
STORE_AVAILABLE_FREQUENCIES_PROCESSING_TIME = GRPC_REQUEST_PROCESSING_TIME.labels('store_available_frequencies_in_db')
