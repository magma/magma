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

from typing import List
from uuid import uuid4

from dp.protos.enodebd_dp_pb2 import CBSDStateResult, LteChannel

SOME_FCC_ID = "some_fcc_id"
USER_ID = "some_user_id"
UNREGISTERED = "unregistered"


class CbsdAPIDataBuilder:
    def __init__(self):
        self.frequency_mhz = 3625
        self.bandwidth_mhz = 10
        self.max_eirp = 28
        self.grant_state = "authorized"
        self.payload = {
            'fcc_id': SOME_FCC_ID,
            'serial_number': str(uuid4()),
            'user_id': USER_ID,
            'cbsd_category': 'b',
            'single_step_enabled': False,
        }

    def with_serial_number(self, serial_number: str) -> CbsdAPIDataBuilder:
        self.payload['serial_number'] = serial_number
        return self

    def with_fcc_id(self, fcc_id: str = SOME_FCC_ID) -> CbsdAPIDataBuilder:
        self.payload['fcc_id'] = fcc_id
        return self

    def with_cbsd_category(self, cbsd_category: str = "b") -> CbsdAPIDataBuilder:
        self.payload['cbsd_category'] = cbsd_category
        return self

    def with_latitude_deg(self, latitude_deg: float = 10.5) -> CbsdAPIDataBuilder:
        installation_param = self.payload.setdefault("installation_param", {})
        installation_param["latitude_deg"] = latitude_deg
        return self

    def with_longitude_deg(self, longitude_deg: float = 11.5) -> CbsdAPIDataBuilder:
        installation_param = self.payload.setdefault("installation_param", {})
        installation_param["longitude_deg"] = longitude_deg
        return self

    def with_antenna_gain(self, antenna_gain: int = 15) -> CbsdAPIDataBuilder:
        installation_param = self.payload.setdefault("installation_param", {})
        installation_param["antenna_gain"] = antenna_gain
        return self

    def with_indoor_deployment(self, indoor_deployment: bool = False) -> CbsdAPIDataBuilder:
        installation_param = self.payload.setdefault("installation_param", {})
        installation_param["indoor_deployment"] = indoor_deployment
        return self

    def with_full_installation_param(
            self,
            latitude_deg: float = 10.5,
            longitude_deg: float = 11.5,
            antenna_gain: int = 15,
            indoor_deployment: bool = True,
            height_m: float = 12.5,
            height_type: str = "agl",
    ) -> CbsdAPIDataBuilder:
        self.payload["installation_param"] = {
            "latitude_deg": latitude_deg,
            "longitude_deg": longitude_deg,
            "antenna_gain": antenna_gain,
            "indoor_deployment": indoor_deployment,
            "height_m": height_m,
            "height_type": height_type,
        }
        return self

    def with_frequency_preferences(
            self,
            bandwidth_mhz: int = 20,
            frequencies_mhz: List[int] = None,
    ) -> CbsdAPIDataBuilder:
        self.payload["frequency_preferences"] = {
            "bandwidth_mhz": bandwidth_mhz,
            "frequencies_mhz": frequencies_mhz or [],
        }
        return self

    def with_capabilities(self, max_power=20, min_power=0, number_of_antennas=2):
        self.payload['capabilities'] = {
            'max_power': max_power,
            'min_power': min_power,
            'number_of_antennas': number_of_antennas,
        }
        return self

    def with_desired_state(self, desired_state: str = "registered") -> CbsdAPIDataBuilder:
        self.payload["desired_state"] = desired_state
        return self

    def with_expected_grant(
        self, bandwidth_mhz: int = 10, frequency_mhz: int = 3625, max_eirp: int = 28,
        grant_state="authorized",
    ) -> CbsdAPIDataBuilder:
        self.bandwidth_mhz = bandwidth_mhz
        self.frequency_mhz = frequency_mhz
        self.max_eirp = max_eirp
        self.grant_state = grant_state
        return self

    def with_grant(
        self, bandwidth_mhz: int = None, frequency_mhz: int = None, max_eirp: int = None, grant_state=None,
    ) -> CbsdAPIDataBuilder:
        self.payload['grant'] = {
            'bandwidth_mhz': bandwidth_mhz or self.bandwidth_mhz,
            'frequency_mhz': frequency_mhz or self.frequency_mhz,
            'max_eirp': max_eirp or self.max_eirp,
            'state': grant_state or self.grant_state,
        }
        return self

    def with_max_eirp(self, max_eirp: int = 28) -> CbsdAPIDataBuilder:
        self.payload['max_eirp'] = max_eirp
        return self

    def with_state(self, state: str = UNREGISTERED) -> CbsdAPIDataBuilder:
        self.payload['state'] = state
        return self

    def with_cbsd_id(self, cbsd_id: str) -> CbsdAPIDataBuilder:
        self.payload['cbsd_id'] = cbsd_id
        return self

    def with_is_active(self, is_active: bool) -> CbsdAPIDataBuilder:
        self.payload['is_active'] = is_active
        return self

    def with_single_step_enabled(self, enabled: bool) -> CbsdAPIDataBuilder:
        self.payload['single_step_enabled'] = enabled
        return self

    def build_grant_state_data(self, frequenzy_mhz=None, bandwidth_mhz=None, max_eirp=None) -> CBSDStateResult:
        frequenzy_mhz = frequenzy_mhz or self.frequency_mhz
        bandwidth_mhz = bandwidth_mhz or self.bandwidth_mhz
        max_eirp = max_eirp or self.max_eirp
        frequency_hz = int(1e6) * frequenzy_mhz
        half_bandwidth_hz = int(5e5) * bandwidth_mhz
        return CBSDStateResult(
            radio_enabled=True,
            channel=LteChannel(
                low_frequency_hz=frequency_hz - half_bandwidth_hz,
                high_frequency_hz=frequency_hz + half_bandwidth_hz,
                max_eirp_dbm_mhz=max_eirp,
            ),
        )
