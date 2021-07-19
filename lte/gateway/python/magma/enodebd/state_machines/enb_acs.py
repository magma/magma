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
from abc import ABC, abstractmethod
from asyncio import BaseEventLoop
from typing import Any, Type

from magma.common.service import MagmaService
from magma.enodebd.data_models.data_model import DataModel
from magma.enodebd.data_models.data_model_parameters import ParameterName
from magma.enodebd.device_config.enodeb_config_postprocessor import (
    EnodebConfigurationPostProcessor,
)
from magma.enodebd.device_config.enodeb_configuration import EnodebConfiguration
from magma.enodebd.devices.device_utils import EnodebDeviceName
from magma.enodebd.state_machines.acs_state_utils import are_tr069_params_equal


class EnodebAcsStateMachine(ABC):
    """
    Handles all TR-069 messages.
    Acts as the Auto Configuration Server (ACS), as specified by TR-069.
    A device/version specific ACS message handler.
    Different devices have various idiosyncrasies.
    Subclass BasicEnodebAcsStateMachine for a specific device/version
    implementation.

    This ACS class can only handle a single connected eNodeB device.
    Multiple connected eNodeB devices will lead to undefined behavior.

    This ABC is more of an interface definition.
    """

    def __init__(self) -> None:
        self._service = None
        self._desired_cfg = None
        self._device_cfg = None
        self._data_model = None
        self._are_invasive_changes_applied = True

    def has_parameter(self, param: ParameterName) -> bool:
        """
        Return True if the data model has the parameter

        Raise KeyError if the parameter is optional and we do not know yet
        if this eNodeB has the parameter
        """
        return self.data_model.is_parameter_present(param)

    def get_parameter(self, param: ParameterName) -> Any:
        """
        Return the value of the parameter
        """
        return self.device_cfg.get_parameter(param)

    def set_parameter_asap(self, param: ParameterName, value: Any) -> None:
        """
        Set the parameter to the suggested value ASAP
        """
        self.desired_cfg.set_parameter(param, value)

    def is_enodeb_configured(self) -> bool:
        """
        True if the desired configuration matches the device configuration
        """
        if self.desired_cfg is None:
            return False
        if not self.data_model.are_param_presences_known():
            return False
        desired = self.desired_cfg.get_parameter_names()

        for name in desired:
            val1 = self.desired_cfg.get_parameter(name)
            val2 = self.device_cfg.get_parameter(name)
            type_ = self.data_model.get_parameter(name).type
            if not are_tr069_params_equal(val1, val2, type_):
                return False

        for obj_name in self.desired_cfg.get_object_names():
            params = self.desired_cfg.get_parameter_names_for_object(obj_name)
            for name in params:
                val1 = self.device_cfg.get_parameter_for_object(name, obj_name)
                val2 = self.desired_cfg.get_parameter_for_object(
                    name,
                    obj_name,
                )
                type_ = self.data_model.get_parameter(name).type
                if not are_tr069_params_equal(val1, val2, type_):
                    return False
        return True

    @abstractmethod
    def get_state(self) -> str:
        """
        Get info about the state of the ACS
        """
        pass

    @abstractmethod
    def handle_tr069_message(self, message: Any) -> Any:
        """
        Given a TR-069 message sent from the hardware, return an
        appropriate response
        """
        pass

    @abstractmethod
    def transition(self, next_state: str) -> None:
        pass

    @property
    def service(self) -> MagmaService:
        return self._service

    @service.setter
    def service(self, service: MagmaService) -> None:
        self._service = service

    @property
    def event_loop(self) -> BaseEventLoop:
        return self._service.loop

    @property
    def mconfig(self) -> Any:
        return self._service.mconfig

    @property
    def service_config(self) -> Any:
        return self._service.config

    @property
    def desired_cfg(self) -> EnodebConfiguration:
        return self._desired_cfg

    @desired_cfg.setter
    def desired_cfg(self, val: EnodebConfiguration) -> None:
        self._desired_cfg = val

    @property
    def device_cfg(self) -> EnodebConfiguration:
        return self._device_cfg

    @device_cfg.setter
    def device_cfg(self, val: EnodebConfiguration) -> None:
        self._device_cfg = val

    @property
    def data_model(self) -> DataModel:
        return self._data_model

    @data_model.setter
    def data_model(self, data_model) -> None:
        self._data_model = data_model

    @property
    def are_invasive_changes_applied(self) -> bool:
        return self._are_invasive_changes_applied

    @are_invasive_changes_applied.setter
    def are_invasive_changes_applied(self, is_applied: bool) -> None:
        self._are_invasive_changes_applied = is_applied

    @property
    @abstractmethod
    def data_model_class(self) -> Type[DataModel]:
        pass

    @property
    @abstractmethod
    def device_name(self) -> EnodebDeviceName:
        pass

    @property
    @abstractmethod
    def config_postprocessor(self) -> EnodebConfigurationPostProcessor:
        pass

    @abstractmethod
    def reboot_asap(self) -> None:
        """
        Send a request to reboot the eNodeB ASAP
        """
        pass

    @abstractmethod
    def is_enodeb_connected(self) -> bool:
        pass

    @abstractmethod
    def stop_state_machine(self) -> None:
        pass
