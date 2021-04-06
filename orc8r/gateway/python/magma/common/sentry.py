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


def sentry_init():
    """
    Initialize connection and start piping errors to sentry.io
    """
    sentry_url = get_service_config_value('control_proxy', 'sentry_url', default="")
    if sentry_url:
        sentry_sample_rate = get_service_config_value('control_proxy', 'sentry_sample_rate', default=1.0)
        sentry_sdk.init(
            dsn=sentry_url, 
            release=os.environ['COMMIT_HASH'],
            traces_sample_rate=sentry_sample_rate, 
        )
        sentry_sdk.set_tag("hwid", snowflake.snowflake())
