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

import sentry_sdk

from magma.common.health.service_state_wrapper import ServiceStateWrapper
from magma.common.service import MagmaService
from magma.configuration.service_configs import load_service_config, get_service_config_value

from magma.health.state_recovery import StateRecoveryJob


def main():
    """
    Top-level function for health service
    """
    service = MagmaService('health', None)

    # Optionally pipe errors to Sentry
    sentry_url = get_service_config_value('control_proxy', 'sentry_url', default="")
    if sentry_url:
        sentry_sample_rate = get_service_config_value('control_proxy', 'sentry_sample_rate',default=1.0)
        sentry_sdk.init(dsn=sentry_url, traces_sample_rate=sentry_sample_rate)

    # Service state wrapper obj
    service_state = ServiceStateWrapper()

    # Load service YML config
    state_recovery_config = service.config["state_recovery"]
    services_check = state_recovery_config["services_check"]
    polling_interval = int(state_recovery_config["interval_check_mins"]) * 60
    restart_threshold = state_recovery_config["restart_threshold"]
    snapshots_dir = state_recovery_config["snapshots_dir"]

    redis_dump_src = load_service_config("redis").get("dir", "/var/opt/magma")

    state_recovery_job = StateRecoveryJob(service_state=service_state,
                                          polling_interval=polling_interval,
                                          services_check=services_check,
                                          restart_threshold=restart_threshold,
                                          redis_dump_src=redis_dump_src,
                                          snapshots_dir=snapshots_dir,
                                          service_loop=service.loop)
    state_recovery_job.start()

    # Run the service loop
    service.run()

    # Cleanup the service
    service.close()


if __name__ == "__main__":
    main()
