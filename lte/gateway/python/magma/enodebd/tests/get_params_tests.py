from parameterized import parameterized

from magma.enodebd.data_models.data_model_parameters import ParameterName
from magma.enodebd.device_config.configuration_init import build_desired_config
from magma.enodebd.devices.baicells_436Q import Baicells436QTrDataModel
from magma.enodebd.devices.device_utils import EnodebDeviceName
from magma.enodebd.state_machines.acs_state_utils import get_object_params_to_get
from magma.enodebd.tests.test_utils.enb_acs_builder import EnodebAcsStateMachineBuilder
from magma.enodebd.tests.test_utils.enodeb_handler import EnodebHandlerTestCase


class GetParamsTestCase(EnodebHandlerTestCase):
    @parameterized.expand([
        (True, True, 4),
        (True, False, 4),
        (False, True, 0),
        (False, False, 4),
    ])
    def test_get_object_params_to_get(
            self, request_all_params: bool, with_desired_config: bool, expected_object_names_list_len: int,
    ):
        acs_state_machine = self._prepare_sm()
        data_model = Baicells436QTrDataModel()

        if with_desired_config:
            acs_state_machine.desired_cfg = self._prepare_desired_cfg_for_sm(acs_state_machine)

        obj_names = get_object_params_to_get(
            desired_cfg=acs_state_machine.desired_cfg,
            device_cfg=acs_state_machine.device_cfg,
            data_model=data_model,
            request_all_params=request_all_params,
        )

        self.assertEqual(expected_object_names_list_len, len(obj_names))

    @staticmethod
    def _prepare_sm():
        sm = EnodebAcsStateMachineBuilder.build_acs_state_machine(EnodebDeviceName.BAICELLS_436Q)
        sm.device_cfg.set_parameter(ParameterName.IP_SEC_ENABLE, False)
        sm.device_cfg.set_parameter(ParameterName.NUM_PLMNS, 1)
        return sm

    @staticmethod
    def _prepare_desired_cfg_for_sm(sm):
        return build_desired_config(
            sm.mconfig,
            sm.service_config,
            sm.device_cfg,
            sm.data_model,
            sm.config_postprocessor,
        )
