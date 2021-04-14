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

import asyncio
from unittest import mock

from magma.common.service import MagmaService
from magma.enodebd.devices.device_map import get_device_handler_from_name
from magma.enodebd.devices.device_utils import EnodebDeviceName
from magma.enodebd.state_machines.enb_acs import EnodebAcsStateMachine
from magma.enodebd.state_machines.enb_acs_manager import StateMachineManager
from magma.enodebd.tests.test_utils.config_builder import EnodebConfigBuilder


class EnodebAcsStateMachineBuilder:
    @classmethod
    def build_acs_manager(
        cls,
        device: EnodebDeviceName = EnodebDeviceName.BAICELLS,
    ) -> StateMachineManager:
        service = cls.build_magma_service(device)
        return StateMachineManager(service)

    @classmethod
    def build_multi_enb_acs_manager(
        cls,
    ) -> StateMachineManager:
        service = cls.build_multi_enb_magma_service()
        return StateMachineManager(service)

    @classmethod
    def build_multi_enb_acs_state_machine(
        cls,
        device: EnodebDeviceName = EnodebDeviceName.BAICELLS,
    ) -> EnodebAcsStateMachine:
        # Build the state_machine
        service = cls.build_multi_enb_magma_service()
        handler_class = get_device_handler_from_name(device)
        acs_state_machine = handler_class(service)
        return acs_state_machine

    @classmethod
    def build_acs_state_machine(
        cls,
        device: EnodebDeviceName = EnodebDeviceName.BAICELLS,
    ) -> EnodebAcsStateMachine:
        # Build the state_machine
        service = cls.build_magma_service(device)
        handler_class = get_device_handler_from_name(device)
        acs_state_machine = handler_class(service)
        return acs_state_machine

    @classmethod
    def build_magma_service(
        cls,
        device: EnodebDeviceName = EnodebDeviceName.BAICELLS,
    ) -> MagmaService:
        event_loop = asyncio.get_event_loop()
        mconfig = EnodebConfigBuilder.get_mconfig(device)
        service_config = EnodebConfigBuilder.get_service_config()
        with mock.patch('magma.common.service.MagmaService') as MockService:
            MockService.config = service_config
            MockService.mconfig = mconfig
            MockService.loop = event_loop
            return MockService

    @classmethod
    def build_multi_enb_magma_service(cls) -> MagmaService:
        event_loop = asyncio.get_event_loop()
        mconfig = EnodebConfigBuilder.get_multi_enb_mconfig()
        service_config = EnodebConfigBuilder.get_service_config()
        with mock.patch('magma.common.service.MagmaService') as MockService:
            MockService.config = service_config
            MockService.mconfig = mconfig
            MockService.loop = event_loop
            return MockService
