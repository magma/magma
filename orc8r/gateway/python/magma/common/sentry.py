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

import os

import sentry_sdk
import snowflake
from magma.configuration.service_configs import get_service_config_value

CONTROL_PROXY = 'control_proxy'
SENTRY_URL = 'sentry_url_python'
SENTRY_SAMPLE_RATE = 'sentry_sample_rate'
COMMIT_HASH = 'COMMIT_HASH'
HWID = 'hwid'
SERVICE_NAME = 'service_name'


def sentry_init(service_name: str):
    """Initialize connection and start piping errors to sentry.io."""
    sentry_url = get_service_config_value(
        CONTROL_PROXY,
        SENTRY_URL,
        default='',
    )
    if not sentry_url:
        return

    sentry_sample_rate = get_service_config_value(
        CONTROL_PROXY,
        SENTRY_SAMPLE_RATE,
        default=1.0,
    )
    sentry_sdk.init(
        dsn=sentry_url,
        release=os.getenv(COMMIT_HASH),
        traces_sample_rate=sentry_sample_rate,
    )
    sentry_sdk.set_tag(HWID, snowflake.snowflake())
    sentry_sdk.set_tag(SERVICE_NAME, service_name)
