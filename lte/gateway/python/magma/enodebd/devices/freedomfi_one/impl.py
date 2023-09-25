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
from typing import Any, Dict

from magma.common.service import MagmaService
from magma.enodebd.data_models.data_model import DataModel
from magma.enodebd.data_models.data_model_parameters import ParameterName
from magma.enodebd.device_config.enodeb_config_postprocessor import (
    EnodebConfigurationPostProcessor,
)
from magma.enodebd.device_config.enodeb_configuration import EnodebConfiguration
from magma.enodebd.devices.device_utils import EnodebDeviceName
from magma.enodebd.devices.freedomfi_one.data_model import (
    FreedomFiOneTrDataModel,
)
from magma.enodebd.devices.freedomfi_one.params import (
    CarrierAggregationParameters,
    FreedomFiOneMiscParameters,
    SASParameters,
)
from magma.enodebd.devices.freedomfi_one.states import (
    FreedomFiOneEndSessionState,
    FreedomFiOneGetInitState,
    FreedomFiOneGetObjectParametersState,
    FreedomFiOneNotifyDPState,
    FreedomFiOneSendGetTransientParametersState,
)
from magma.enodebd.logger import EnodebdLogger
from magma.enodebd.state_machines.enb_acs import EnodebAcsStateMachine
from magma.enodebd.state_machines.enb_acs_impl import BasicEnodebAcsStateMachine
from magma.enodebd.state_machines.enb_acs_states import (
    AddObjectsState,
    CheckFirmwareUpgradeDownloadState,
    DeleteObjectsState,
    EnbSendRebootState,
    EnodebAcsState,
    ErrorState,
    FirmwareUpgradeDownloadState,
    GetParametersState,
    SetParameterValuesState,
    WaitForFirmwareUpgradeDownloadResponse,
    WaitGetParametersState,
    WaitInformMRebootState,
    WaitInformState,
    WaitRebootResponseState,
    WaitSetParameterValuesState,
)

SAS_KEY = 'sas'
WEB_UI_ENABLE_LIST_KEY = 'web_ui_enable_list'
DP_MODE_KEY = 'dp_mode'

RADIO_MIN_POWER = 0
RADIO_MAX_POWER = 24


class FreedomFiOneHandler(BasicEnodebAcsStateMachine):
    """
    FreedomFi One State Machine
    """

    def __init__(
            self,
            service: MagmaService,
    ) -> None:
        self._state_map: Dict[str, Any] = {}
        super().__init__(service=service, use_param_key=True)

    def reboot_asap(self) -> None:
        """
        Transition to 'reboot' state
        """
        self.transition('reboot')

    def is_enodeb_connected(self) -> bool:
        """
        Check if enodebd has received an Inform from the enodeb

        Returns:
            bool
        """
        return not isinstance(self.state, WaitInformState)

    def _init_state_map(self) -> None:
        self._state_map = {
            # Inform comes in -> Respond with InformResponse
            'wait_inform': WaitInformState(self, when_done='get_rpc_methods'),
            # If first inform after boot -> GetRpc request comes in, if not
            # empty request comes in => Transition
            'get_rpc_methods': FreedomFiOneGetInitState(
                self,
                when_done='check_fw_upgrade_download',
            ),

            # Download flow
            'check_fw_upgrade_download': CheckFirmwareUpgradeDownloadState(
                self,
                when_download='fw_upgrade_download',
                when_skip='get_transient_params',
            ),
            'fw_upgrade_download': FirmwareUpgradeDownloadState(
                self,
                when_done='wait_fw_upgrade_download_response',
            ),
            'wait_fw_upgrade_download_response': WaitForFirmwareUpgradeDownloadResponse(
                self,
                when_done='get_transient_params',
                when_skip='get_transient_params',
            ),
            # Download flow ends

            # Read transient readonly params.
            'get_transient_params': FreedomFiOneSendGetTransientParametersState(
                self,
                when_done='get_params',
            ),

            'get_params': FreedomFiOneGetObjectParametersState(
                self,
                when_delete='delete_objs',
                when_add='add_objs',
                when_set='set_params',
                when_skip='end_session',
            ),

            'delete_objs': DeleteObjectsState(
                self, when_add='add_objs',
                when_skip='set_params',
            ),
            'add_objs': AddObjectsState(self, when_done='set_params'),
            'set_params': SetParameterValuesState(
                self,
                when_done='wait_set_params',
            ),
            'wait_set_params': WaitSetParameterValuesState(
                self,
                when_done='check_get_params',
                when_apply_invasive='check_get_params',
                status_non_zero_allowed=True,
            ),
            'check_get_params': GetParametersState(
                self,
                when_done='check_wait_get_params',
                request_all_params=True,
            ),
            'check_wait_get_params': WaitGetParametersState(
                self,
                when_done='end_session',
            ),
            'end_session': FreedomFiOneEndSessionState(self, when_dp_mode='notify_dp'),
            'notify_dp': FreedomFiOneNotifyDPState(self, when_inform='wait_inform'),

            # These states are only entered through manual user intervention
            'reboot': EnbSendRebootState(self, when_done='wait_reboot'),
            'wait_reboot': WaitRebootResponseState(
                self,
                when_done='wait_post_reboot_inform',
            ),
            'wait_post_reboot_inform': WaitInformMRebootState(
                self,
                when_done='wait_inform',
                when_timeout='wait_inform',
            ),
            # The states below are entered when an unexpected message type is
            # received
            'unexpected_fault': ErrorState(
                self,
                inform_transition_target='wait_inform',
            ),
        }

    @property
    def device_name(self) -> str:
        """
        Return the device name

        Returns:
            device name
        """
        return EnodebDeviceName.FREEDOMFI_ONE

    @property
    def data_model_class(self) -> DataModel:
        """
        Return the class of the data model

        Returns:
            DataModel
        """
        return FreedomFiOneTrDataModel

    @property
    def config_postprocessor(self) -> EnodebConfigurationPostProcessor:
        """
        Return the instance of config postprocessor

        Returns:
            EnodebConfigurationPostProcessor
        """
        return FreedomFiOneConfigurationInitializer(self)

    @property
    def state_map(self) -> Dict[str, EnodebAcsState]:
        """
        Return the state map for the State Machine

        Returns:
            Dict[str, EnodebAcsState]
        """
        return self._state_map

    @property
    def disconnected_state_name(self) -> str:
        """
        Return the string representation of a disconnected state

        Returns:
            str
        """
        return 'wait_inform'

    @property
    def unexpected_fault_state_name(self) -> str:
        """
        Return the string representation of an unexpected fault state

        Returns:
            str
        """
        return 'unexpected_fault'


class FreedomFiOneConfigurationInitializer(EnodebConfigurationPostProcessor):
    """
    Overrides desired config on the State Machine
    """

    SAS_KEY = 'sas'
    WEB_UI_ENABLE_LIST_KEY = 'web_ui_enable_list'

    def __init__(self, acs: EnodebAcsStateMachine):
        super().__init__()
        self.acs = acs

    def postprocess(
            self, mconfig: Any, service_cfg: Any,
            desired_cfg: EnodebConfiguration,
    ) -> None:
        """
        Add some params to the desired config

        Args:
            mconfig (Any): mconfig
            service_cfg (Any): service config
            desired_cfg (EnodebConfiguration): desired config
        """
        desired_cfg.delete_parameter(ParameterName.BAND)
        desired_cfg.delete_parameter(ParameterName.EARFCNDL)
        desired_cfg.delete_parameter(ParameterName.DL_BANDWIDTH)
        desired_cfg.delete_parameter(ParameterName.UL_BANDWIDTH)

        self._set_default_params(desired_cfg)
        self._increment_param_version_key()
        self._verify_cell_reserved_param(desired_cfg)
        self._verify_ui_enable(service_cfg, desired_cfg)
        self._verify_sas_params(service_cfg, desired_cfg)
        self._set_misc_params_from_service_config(service_cfg, desired_cfg)

    def _set_default_params(self, desired_cfg):
        """Go through default params and set them in desired config"""
        defaults = {
            **FreedomFiOneMiscParameters.defaults,
            **SASParameters.defaults,
            **CarrierAggregationParameters.defaults,
        }
        for name, val in defaults.items():
            desired_cfg.set_parameter(param_name=name, value=val)

    def _increment_param_version_key(self):
        """Bump up the parameter key version"""
        self.acs.parameter_version_inc()

    def _verify_cell_reserved_param(self, desired_cfg):
        """
        Workaround a bug in Sercomm firmware in release 3920, 3921
        where the meaning of CellReservedForOperatorUse is wrong.
        Set to True to ensure the PLMN is not reserved

        Args:
            desired_cfg: desired eNB config
        """
        num_plmns = self.acs.data_model.get_num_plmns()
        for i in range(1, num_plmns + 1):
            object_name = ParameterName.PLMN_N % i
            desired_cfg.set_parameter_for_object(
                param_name=ParameterName.PLMN_N_CELL_RESERVED % i,
                value=True,
                object_name=object_name,
            )

    def _verify_ui_enable(self, service_cfg, desired_cfg):
        if WEB_UI_ENABLE_LIST_KEY in service_cfg:
            serial_nos = service_cfg.get(WEB_UI_ENABLE_LIST_KEY)
            if self.acs.device_cfg.has_parameter(
                    ParameterName.SERIAL_NUMBER,
            ):
                if self.acs.get_parameter(ParameterName.SERIAL_NUMBER) in \
                        serial_nos:
                    desired_cfg.set_parameter(
                        FreedomFiOneMiscParameters.WEB_UI_ENABLE,
                        value=True,
                    )
            else:
                # This should not happen
                EnodebdLogger.error("Serial number unknown for device")

    def _verify_sas_params(self, service_cfg, desired_cfg):
        sas_cfg = service_cfg.get(SAS_KEY)
        if not sas_cfg or sas_cfg[DP_MODE_KEY]:
            desired_cfg.set_parameter(SASParameters.SAS_METHOD, value=True)
            return

        sas_param_names = self.acs.data_model.get_sas_param_names()
        for name, val in sas_cfg.items():
            if name not in sas_param_names:
                EnodebdLogger.warning("Ignoring attribute %s", name)
                continue
            desired_cfg.set_parameter(name, val)

    def _set_misc_params_from_service_config(self, service_cfg, desired_cfg):
        prim_src_name = FreedomFiOneMiscParameters.PRIM_SOURCE
        prim_src_service_cfg_val = service_cfg.get(prim_src_name)
        desired_cfg.set_parameter(prim_src_name, prim_src_service_cfg_val)
