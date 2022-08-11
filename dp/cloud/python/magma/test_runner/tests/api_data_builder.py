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

from dp.protos.cbsd_pb2 import (
    CBSDStateResult,
    EnodebdUpdateCbsdRequest,
    InstallationParam,
    LteChannel,
)
from google.protobuf.wrappers_pb2 import BoolValue, DoubleValue, StringValue

SOME_FCC_ID = "some_fcc_id"
USER_ID = "some_user_id"
UNREGISTERED = "unregistered"


class CbsdAPIDataBuilder:
    def __init__(self):
        self.payload = {
            'fcc_id': SOME_FCC_ID,
            'serial_number': str(uuid4()),
            'user_id': USER_ID,
            'cbsd_category': 'a',
            'desired_state': 'registered',
            'single_step_enabled': False,
            'carrier_aggregation_enabled': False,
            'grant_redundancy': False,
            'grants': [],
            'capabilities': {
                'max_power': 20,
                'min_power': 0,
                'number_of_antennas': 2,
                'max_ibw_mhz': 150,
            },
            'frequency_preferences': {
                'bandwidth_mhz': 10,
                'frequencies_mhz': [3625],
            },
            'installation_param': {
                'antenna_gain': 15,
            },
        }

    def with_serial_number(self, serial_number: str) -> CbsdAPIDataBuilder:
        self.payload['serial_number'] = serial_number
        return self

    def with_fcc_id(self, fcc_id: str = SOME_FCC_ID) -> CbsdAPIDataBuilder:
        self.payload['fcc_id'] = fcc_id
        return self

    def with_cbsd_category(self, cbsd_category: str) -> CbsdAPIDataBuilder:
        self.payload['cbsd_category'] = cbsd_category
        return self

    def with_latitude_deg(self, latitude_deg: float = 10.5) -> CbsdAPIDataBuilder:
        self.payload['installation_param']['latitude_deg'] = latitude_deg
        return self

    def with_longitude_deg(self, longitude_deg: float = 11.5) -> CbsdAPIDataBuilder:
        self.payload['installation_param']['longitude_deg'] = longitude_deg
        return self

    def with_indoor_deployment(self, indoor_deployment: bool = False) -> CbsdAPIDataBuilder:
        self.payload['installation_param']["indoor_deployment"] = indoor_deployment
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

    def with_frequency_preferences(self, bandwidth_mhz: int, frequencies_mhz: List[int]) -> CbsdAPIDataBuilder:
        self.payload["frequency_preferences"] = {
            "bandwidth_mhz": bandwidth_mhz,
            "frequencies_mhz": frequencies_mhz,
        }
        return self

    def with_carrier_aggregation(self) -> CbsdAPIDataBuilder:
        self.payload['grant_redundancy'] = True
        self.payload['carrier_aggregation_enabled'] = True
        return self

    def with_desired_state(self, desired_state: str = "registered") -> CbsdAPIDataBuilder:
        self.payload["desired_state"] = desired_state
        return self

    def without_grants(self) -> CbsdAPIDataBuilder:
        self.payload['grants'] = []
        return self

    def with_grant(
            self, bandwidth_mhz: int = 10, frequency_mhz: int = 3625, max_eirp: int = 28,
    ) -> CbsdAPIDataBuilder:
        self.payload['grants'].append({
            'bandwidth_mhz': bandwidth_mhz,
            'frequency_mhz': frequency_mhz,
            'max_eirp': max_eirp,
            'state': 'authorized',
        })
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

    def build_grant_state_data(self) -> CBSDStateResult:
        # TODO rewrite builders to dataclasses
        grants = [_api_to_proto_grant(g) for g in self.payload['grants']]
        return CBSDStateResult(
            radio_enabled=True,
            carrier_aggregation_enabled=self.payload['carrier_aggregation_enabled'],
            channel=grants[0],
            channels=grants,
        )

    def build_enodebd_update_request(self, indoor_deployment=False, cbsd_category="a") -> EnodebdUpdateCbsdRequest:
        return EnodebdUpdateCbsdRequest(
            serial_number=self.payload["serial_number"],
            installation_param=InstallationParam(
                latitude_deg=DoubleValue(value=10.5),
                longitude_deg=DoubleValue(value=11.5),
                indoor_deployment=BoolValue(value=indoor_deployment),
                height_type=StringValue(value="agl"),
                height_m=DoubleValue(value=12.5),
            ),
            cbsd_category=cbsd_category,
        )


def _api_to_proto_grant(grant: dict[str, any]) -> LteChannel:
    frequency_mhz = grant['frequency_mhz']
    bandwidth_mhz = grant['bandwidth_mhz']
    max_eirp_dbm_mhz = grant['max_eirp']
    frequency_hz = 10**6 * frequency_mhz
    bandwidth_hz = 10**6 * bandwidth_mhz
    return LteChannel(
        low_frequency_hz=frequency_hz - bandwidth_hz // 2,
        high_frequency_hz=frequency_hz + bandwidth_hz // 2,
        max_eirp_dbm_mhz=max_eirp_dbm_mhz,
    )
