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

from prometheus_client import Counter, Histogram

# Counters for Diameter/S6a application
S6A_AUTH_SUCCESS_TOTAL = Counter(
    's6a_auth_success',
    'Total successful S6a auth requests',
)
S6A_AUTH_FAILURE_TOTAL = Counter(
    's6a_auth_failure',
    'Total failed S6a auth requests with reason', ['code'],
)
S6A_LUR_TOTAL = Counter(
    's6a_location_update',
    'Total S6a location update requests',
)
M5G_AUTH_SUCCESS_TOTAL = Counter('m5g_auth_success',
                                 'Total successful M5G auth requests')
M5G_AUTH_FAILURE_TOTAL = Counter('m5g_auth_failure',
                                 'Total failed M5G auth requests with reason', ['code'])

DIAMETER_AUTHENTICATION_REJECTED = 4001
DIAMETER_ERROR_USER_UNKNOWN = 5001
DIAMETER_AUTHORIZATION_REJECTED = 5003
# Counters for Diameter base application
DIAMETER_CEX_TOTAL = Counter(
    'diameter_capabilities_exchange',
    'Total Diameter capabilities exchange requests',
)
DIAMETER_WATCHDOG_TOTAL = Counter(
    'diameter_watchdog',
    'Total Diameter watchdog requests',
)
DIAMETER_DISCONECT_TOTAL = Counter(
    'diameter_disconnect',
    'Total Diameter disconnect requests',
)

SUBSCRIBER_SYNC_LATENCY = Histogram(
    'subscriber_sync_latency_ms',
    'Latency syncing subscribers from cloud'
    'in milliseconds',
    buckets=[50, 100, 200, 500, 1000, 5000, 10000, 30000, 60000],
)

SUBSCRIBER_SYNC_SUCCESS_TOTAL = Counter(
    'subscriber_sync_success',
    'Total number of successful subscriber'
    'syncs with cloud',
)

SUBSCRIBER_SYNC_FAILURE_TOTAL = Counter(
    'subscriber_sync_failure',
    'Total number of failed subscriber'
    'syncs with cloud',
)
