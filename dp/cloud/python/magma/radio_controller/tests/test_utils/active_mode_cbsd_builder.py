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
from __future__ import annotations

from typing import List, Optional

from dp.protos.active_mode_pb2 import (
    Cbsd,
    CbsdState,
    Channel,
    DatabaseCbsd,
    EirpCapabilities,
    FrequencyPreferences,
    Grant,
    GrantSettings,
    GrantState,
    InstallationParams,
    SasSettings,
)
from google.protobuf.wrappers_pb2 import FloatValue


class ActiveModeCbsdBuilder:
    def __init__(self):
        self.desired_state = None
        self.cbsd_id = None
        self.state = None
        self.grants = []
        self.channels = []
        self.last_seen_timestamp = None
        self.eirp_capabilities = EirpCapabilities()
        self.sas_settings = SasSettings()
        self.preferences = FrequencyPreferences()
        self.installation = InstallationParams()
        self.db_data = DatabaseCbsd()
        self.grant_settings = GrantSettings()

    def build(self) -> Cbsd:
        return Cbsd(
            cbsd_id=self.cbsd_id,
            sas_settings=self.sas_settings,
            state=self.state,
            desired_state=self.desired_state,
            grants=self.grants,
            channels=self.channels,
            last_seen_timestamp=self.last_seen_timestamp,
            eirp_capabilities=self.eirp_capabilities,
            db_data=self.db_data,
            preferences=self.preferences,
            installation_params=self.installation,
            grant_settings=self.grant_settings,
        )

    def with_single_step_enabled(self) -> ActiveModeCbsdBuilder:
        self.sas_settings.single_step_enabled = True
        return self

    def with_category(self, category: str) -> ActiveModeCbsdBuilder:
        self.sas_settings.cbsd_category = category
        return self

    def deleted(self) -> ActiveModeCbsdBuilder:
        self.db_data.is_deleted = True
        return self

    def updated(self) -> ActiveModeCbsdBuilder:
        self.db_data.should_deregister = True
        return self

    def relinquished(self) -> ActiveModeCbsdBuilder:
        self.db_data.should_relinquish = True
        return self

    def with_id(self, db_id: int) -> ActiveModeCbsdBuilder:
        self.db_data.id = db_id
        return self

    def with_desired_state(self, state: CbsdState) -> ActiveModeCbsdBuilder:
        self.desired_state = state
        return self

    def with_state(self, state: CbsdState) -> ActiveModeCbsdBuilder:
        self.state = state
        return self

    def with_grant_settings(self, grant_settings: GrantSettings):
        self.grant_settings = grant_settings
        return self

    def with_registration(self, prefix: str) -> ActiveModeCbsdBuilder:
        self.cbsd_id = f'{prefix}_cbsd_id'
        self.sas_settings.fcc_id = f'{prefix}_fcc_id'
        self.sas_settings.user_id = f'{prefix}_user_id'
        self.sas_settings.serial_number = f'{prefix}_serial_number'
        return self

    def with_eirp_capabilities(
        self,
        min_power: float, max_power: float, no_ports: int,
    ) -> ActiveModeCbsdBuilder:
        eirp_capabilities = EirpCapabilities(
            min_power=min_power,
            max_power=max_power,
            number_of_ports=no_ports,
        )
        self.eirp_capabilities = eirp_capabilities
        return self

    def with_antenna_gain(self, antenna_gain_dbi: float) -> ActiveModeCbsdBuilder:
        self.installation.antenna_gain_dbi = antenna_gain_dbi
        return self

    def with_grant(
        self,
        grant_id: str, state: GrantState,
        hb_interval_sec: int, last_hb_ts: int,
        low_frequency_hz: int = 3500,
        high_frequency_hz: int = 3700,
    ) -> ActiveModeCbsdBuilder:
        grant = Grant(
            id=grant_id,
            state=state,
            heartbeat_interval_sec=hb_interval_sec,
            last_heartbeat_timestamp=last_hb_ts,
            low_frequency_hz=low_frequency_hz,
            high_frequency_hz=high_frequency_hz,
        )
        self.grants.append(grant)
        return self

    def with_preferences(self, bandwidth_mhz: int, frequencies_mhz: List[int]) -> ActiveModeCbsdBuilder:
        self.preferences = FrequencyPreferences(
            bandwidth_mhz=bandwidth_mhz,
            frequencies_mhz=frequencies_mhz,
        )
        return self

    def with_channel(
        self,
        low: int, high: int,
        max_eirp: Optional[float] = None,
    ) -> ActiveModeCbsdBuilder:
        channel = Channel(
            low_frequency_hz=low,
            high_frequency_hz=high,
            max_eirp=self.make_optional_float(max_eirp),
        )
        self.channels.append(channel)
        return self

    def with_installation_params(
        self,
        latitude_deg: float,
        longitude_deg: float,
        height_m: float,
        height_type: str,
        indoor_deployment: bool,
    ) -> ActiveModeCbsdBuilder:
        self.installation.latitude_deg = latitude_deg
        self.installation.longitude_deg = longitude_deg
        self.installation.height_m = height_m
        self.installation.height_type = height_type
        self.installation.indoor_deployment = indoor_deployment
        return self

    @staticmethod
    def make_optional_float(value: Optional[float] = None) -> FloatValue:
        return FloatValue(value=value) if value is not None else None

    def with_last_seen(self, last_seen_timestamp: int) -> ActiveModeCbsdBuilder:
        self.last_seen_timestamp = last_seen_timestamp
        return self
