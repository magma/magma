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

import json
from collections import namedtuple
from typing import Any, Optional, Union

from lte.protos.mconfig import mconfigs_pb2
from magma.common.misc_utils import get_ip_from_if
from magma.configuration.exceptions import LoadConfigError
from magma.configuration.mconfig_managers import load_service_mconfig_as_json
from magma.enodebd.data_models.data_model import DataModel
from magma.enodebd.data_models.data_model_parameters import ParameterName
from magma.enodebd.device_config.enodeb_config_postprocessor import (
    EnodebConfigurationPostProcessor,
)
from magma.enodebd.device_config.enodeb_configuration import EnodebConfiguration
from magma.enodebd.exceptions import ConfigurationError
from magma.enodebd.logger import EnodebdLogger as logger
from magma.enodebd.lte_utils import (
    DuplexMode,
    map_earfcndl_to_band_earfcnul_mode,
    map_earfcndl_to_duplex_mode,
)

# LTE constants
DEFAULT_S1_PORT = 36412
# This is a known working value for supported eNB devices.
# Cell Identity is a 28 bit number, but not all values are supported.
DEFAULT_CELL_IDENTITY = 138777000

SingleEnodebConfig = namedtuple(
    'SingleEnodebConfig',
    [
        'earfcndl', 'subframe_assignment',
        'special_subframe_pattern',
        'pci', 'plmnid_list', 'tac',
        'bandwidth_mhz', 'cell_id',
        'allow_enodeb_transmit',
        'mme_address', 'mme_port',
        'a1_threshold_rsrp', 'lte_a1_threshold_rsrq',
        'hysteresis', 'time_to_trigger',
        'a2_threshold_rsrp', 'lte_a2_threshold_rsrq',
        'lte_a2_threshold_rsrp_irat_volte', 'lte_a2_threshold_rsrq_irat_volte',
        'a3_offset', 'a3_offset_anr',
        'a4_threshold_rsrp', 'lte_intra_a5_threshold_1_rsrp',
        'lte_intra_a5_threshold_2_rsrp', 'b2_threshold1_rsrp',
        'b2_threshold2_rsrp', 'b2_geran_irat_threshold',
        'qrxlevmin_sib1', 'qrxlevminoffset',
        's_intrasearch', 's_nonintrasearch',
        'qrxlevmin_sib3', 'reselection_priority',
        'threshservinglow', 'x2_enable_disable',
        'ciphering_algorithm', 'integrity_algorithm',
    ],
)


def config_assert(condition: bool, message: str = None) -> None:
    """ To be used in place of 'assert' so that ConfigurationError is raised
        for all config-related exceptions. """
    if not condition:
        raise ConfigurationError(message)


def build_desired_config(
        mconfig: Any,
        service_config: Any,
        device_config: EnodebConfiguration,
        data_model: DataModel,
        post_processor: EnodebConfigurationPostProcessor,
) -> EnodebConfiguration:
    """
    Factory for initializing DESIRED data model configuration.

    When working with the configuration of an eNodeB, we track the
    current state of configuration for that device, as well as what
    configuration we want to set on the device.
    Args:
        mconfig: Managed configuration, eNodeB protobuf message
        service_config:
    Returns:
        Desired data model configuration for the device
    """
    cfg_desired = EnodebConfiguration(data_model)

    # Determine configuration parameters
    _set_management_server(cfg_desired)

    # Attempt to load device configuration from YANG before service mconfig
    enb_config = _get_enb_yang_config(device_config) or \
                 _get_enb_config(mconfig, device_config)
    # Ho Algorithm
    _set_algorithm_a1_threshold_rsrp(cfg_desired, enb_config.a1_threshold_rsrp)
    # _set_algorithm_lte_a1_threshold_rsrq(cfg_desired, enb_config.lte_a1_threshold_rsrq)
    # _set_algorithm_hysteresis(cfg_desired, enb_config.hysteresis)
    # _set_algorithm_time_to_trigger(cfg_desired, enb_config.time_to_trigger)
    _set_algorithm_a2_threshold_rsrp(cfg_desired, enb_config.a2_threshold_rsrp)
    # _set_algorithm_lte_a2_threshold_rsrq(cfg_desired, enb_config.lte_a2_threshold_rsrq)
    # _set_algorithm_lte_a2_threshold_rsrp_irat_volte(cfg_desired, enb_config.lte_a2_threshold_rsrp_irat_volte)
    # _set_algorithm_lte_a2_threshold_rsrq_irat_volte(cfg_desired, enb_config.lte_a2_threshold_rsrq_irat_volte)
    _set_algorithm_a3_offset(cfg_desired, enb_config.a3_offset)
    _set_algorithm_a3_offset_anr(cfg_desired, enb_config.a3_offset_anr)
    _set_algorithm_a4_threshold_rsrp(cfg_desired, enb_config.a4_threshold_rsrp)
    _set_algorithm_lte_intra_a5_threshold_1_rsrp(cfg_desired, enb_config.lte_intra_a5_threshold_1_rsrp)
    _set_algorithm_lte_intra_a5_threshold_2_rsrp(cfg_desired, enb_config.lte_intra_a5_threshold_2_rsrp)
    _set_algorithm_b2_threshold1_rsrp(cfg_desired, enb_config.b2_threshold1_rsrp)
    _set_algorithm_b2_threshold2_rsrp(cfg_desired, enb_config.b2_threshold2_rsrp)
    _set_algorithm_b2_geran_irat_threshold(cfg_desired, enb_config.b2_geran_irat_threshold)
    _set_algorithm_qrxlevmin_sib1(cfg_desired, enb_config.qrxlevmin_sib1)
    _set_algorithm_qrxlevminoffset(cfg_desired, enb_config.qrxlevminoffset)
    _set_algorithm_s_intrasearch(cfg_desired, enb_config.s_intrasearch)
    _set_algorithm_s_nonintrasearch(cfg_desired, enb_config.s_nonintrasearch)
    _set_algorithm_qrxlevmin_sib3(cfg_desired, enb_config.qrxlevmin_sib3)
    _set_algorithm_reselection_priority(cfg_desired, enb_config.reselection_priority)
    _set_algorithm_threshservinglow(cfg_desired, enb_config.threshservinglow)
    _set_algorithm_x2_enable_disable(cfg_desired, enb_config.x2_enable_disable)
    _set_algorithm_ciphering_algorithm(cfg_desired, enb_config.ciphering_algorithm)
    _set_algorithm_integrity_algorithm(cfg_desired, enb_config.integrity_algorithm)

    _set_earfcn_freq_band_mode(
        device_config, cfg_desired, data_model,
        enb_config.earfcndl,
    )
    if enb_config.subframe_assignment is not None:
        _set_tdd_subframe_config(
            device_config, cfg_desired,
            enb_config.subframe_assignment,
            enb_config.special_subframe_pattern,
        )
    _set_pci(cfg_desired, enb_config.pci)
    _set_plmnids_tac(cfg_desired, enb_config.plmnid_list, enb_config.tac)
    _set_bandwidth(cfg_desired, data_model, enb_config.bandwidth_mhz)
    _set_cell_id(cfg_desired, enb_config.cell_id)
    _set_perf_mgmt(
        cfg_desired,
        get_ip_from_if(service_config['tr069']['interface']),
        service_config['tr069']['perf_mgmt_port'],
    )
    _set_misc_static_params(device_config, cfg_desired, data_model)
    if enb_config.mme_address is not None and enb_config.mme_port is not None:
        _set_s1_connection(
            cfg_desired,
            enb_config.mme_address,
            enb_config.mme_port,
        )
    else:
        _set_s1_connection(
            cfg_desired, get_ip_from_if(service_config['s1_interface']),
        )

    # Enable LTE if we should
    cfg_desired.set_parameter(
        ParameterName.ADMIN_STATE,
        enb_config.allow_enodeb_transmit,
    )

    post_processor.postprocess(mconfig, service_config, cfg_desired)
    return cfg_desired


def _get_enb_yang_config(
        device_config: EnodebConfiguration,
) -> Optional[SingleEnodebConfig]:
    """"
    Proof of concept configuration function to load eNB configs from YANG
    data model. Attempts to load configuration from YANG for the eNodeB if
    an entry exists with a matching serial number.
    Args:
        device_config: eNodeB device configuration
    Returns:
        None or a SingleEnodebConfig from YANG with matching serial number
    """
    enb = []
    mme_list = []
    mme_address = None
    mme_port = None
    try:
        enb_serial = \
            device_config.get_parameter(ParameterName.SERIAL_NUMBER)
        config = json.loads(
            load_service_mconfig_as_json('yang').get('value', '{}'),
        )
        enb.extend(
            filter(
                lambda entry: entry['serial'] == enb_serial,
                config.get('cellular', {}).get('enodeb', []),
            ),
        )
    except (ValueError, KeyError, LoadConfigError):
        return None
    if len(enb) == 0:
        return None
    enb_config = enb[0].get('config', {})
    mme_list.extend(enb_config.get('mme', []))
    if len(mme_list) > 0:
        mme_address = mme_list[0].get('host')
        mme_port = mme_list[0].get('port')
    single_enodeb_config = SingleEnodebConfig(
        earfcndl=enb_config.get('earfcndl'),
        subframe_assignment=enb_config.get('subframe_assignment'),
        special_subframe_pattern=enb_config.get('special_subframe_pattern'),
        pci=enb_config.get('pci'),
        plmnid_list=",".join(enb_config.get('plmnid', [])),
        tac=enb_config.get('tac'),
        bandwidth_mhz=enb_config.get('bandwidth_mhz'),
        cell_id=enb_config.get('cell_id'),
        allow_enodeb_transmit=enb_config.get('transmit_enabled'),
        mme_address=mme_address,
        mme_port=mme_port,
    )
    return single_enodeb_config


def _get_enb_config(
        mconfig: mconfigs_pb2.EnodebD,
        device_config: EnodebConfiguration,
) -> SingleEnodebConfig:
    # For fields that are specified per eNB
    if mconfig.enb_configs_by_serial is not None and \
            len(mconfig.enb_configs_by_serial) > 0:
        enb_serial = \
            device_config.get_parameter(ParameterName.SERIAL_NUMBER)
        if enb_serial in mconfig.enb_configs_by_serial:
            enb_config = mconfig.enb_configs_by_serial[enb_serial]
            earfcndl = enb_config.earfcndl
            pci = enb_config.pci
            allow_enodeb_transmit = enb_config.transmit_enabled
            tac = enb_config.tac
            bandwidth_mhz = enb_config.bandwidth_mhz
            cell_id = enb_config.cell_id
            duplex_mode = map_earfcndl_to_duplex_mode(earfcndl)
            subframe_assignment = None
            special_subframe_pattern = None
            if duplex_mode == DuplexMode.TDD:
                subframe_assignment = enb_config.subframe_assignment
                special_subframe_pattern = \
                    enb_config.special_subframe_pattern

            if enb_config.a1_threshold_rsrp is not None:
                a1_threshold_rsrp = enb_config.a1_threshold_rsrp
            if enb_config.lte_a1_threshold_rsrq is not None:
                lte_a1_threshold_rsrq = enb_config.lte_a1_threshold_rsrq
            if enb_config.hysteresis is not None:
                hysteresis = enb_config.hysteresis
            if enb_config.time_to_trigger is not None:
                time_to_trigger = enb_config.time_to_trigger
            if enb_config.a2_threshold_rsrp is not None:
                a2_threshold_rsrp = enb_config.a2_threshold_rsrp
            if enb_config.lte_a2_threshold_rsrq is not None:
                lte_a2_threshold_rsrq = enb_config.lte_a2_threshold_rsrq
            if enb_config.lte_a2_threshold_rsrp_irat_volte is not None:
                lte_a2_threshold_rsrp_irat_volte = enb_config.lte_a2_threshold_rsrp_irat_volte
            if enb_config.lte_a2_threshold_rsrq_irat_volte is not None:
                lte_a2_threshold_rsrq_irat_volte = enb_config.lte_a2_threshold_rsrq_irat_volte

            a3_offset = enb_config.a3_offset
            a3_offset_anr = enb_config.a3_offset_anr
            a4_threshold_rsrp = enb_config.a4_threshold_rsrp
            lte_intra_a5_threshold_1_rsrp = enb_config.lte_intra_a5_threshold_1_rsrp
            lte_intra_a5_threshold_2_rsrp = enb_config.lte_intra_a5_threshold_2_rsrp
            b2_threshold1_rsrp = enb_config.b2_threshold1_rsrp
            b2_threshold2_rsrp = enb_config.b2_threshold2_rsrp
            b2_geran_irat_threshold = enb_config.b2_geran_irat_threshold
            qrxlevmin_sib1 = enb_config.qrxlevmin_sib1
            qrxlevminoffset = enb_config.qrxlevminoffset
            s_intrasearch = enb_config.s_intrasearch
            s_nonintrasearch = enb_config.s_nonintrasearch
            qrxlevmin_sib3 = enb_config.qrxlevmin_sib3
            reselection_priority = enb_config.reselection_priority
            threshservinglow = enb_config.threshservinglow
            x2_enable_disable = enb_config.x2_enable_disable
            ciphering_algorithm = enb_config.ciphering_algorithm
            integrity_algorithm = enb_config.integrity_algorithm

        else:
            raise ConfigurationError(
                'Could not construct desired config '
                'for eNB',
            )
    else:
        pci = mconfig.pci
        allow_enodeb_transmit = mconfig.allow_enodeb_transmit
        tac = mconfig.tac
        bandwidth_mhz = mconfig.bandwidth_mhz
        cell_id = DEFAULT_CELL_IDENTITY
        if mconfig.tdd_config is not None and str(mconfig.tdd_config) != '':
            earfcndl = mconfig.tdd_config.earfcndl
            subframe_assignment = mconfig.tdd_config.subframe_assignment
            special_subframe_pattern = \
                mconfig.tdd_config.special_subframe_pattern
        elif mconfig.fdd_config is not None and str(mconfig.fdd_config) != '':
            earfcndl = mconfig.fdd_config.earfcndl
            subframe_assignment = None
            special_subframe_pattern = None
        else:
            earfcndl = mconfig.earfcndl
            subframe_assignment = mconfig.subframe_assignment
            special_subframe_pattern = mconfig.special_subframe_pattern

    # And now the rest of the fields
    plmnid_list = mconfig.plmnid_list

    single_enodeb_config = SingleEnodebConfig(
        earfcndl=earfcndl,
        subframe_assignment=subframe_assignment,
        special_subframe_pattern=special_subframe_pattern,
        pci=pci,
        plmnid_list=plmnid_list,
        tac=tac,
        bandwidth_mhz=bandwidth_mhz,
        cell_id=cell_id,
        allow_enodeb_transmit=allow_enodeb_transmit,
        mme_address=None,
        mme_port=None,
        a1_threshold_rsrp=a1_threshold_rsrp,
        lte_a1_threshold_rsrq=lte_a1_threshold_rsrq,
        hysteresis=hysteresis,
        time_to_trigger=time_to_trigger,
        a2_threshold_rsrp=a2_threshold_rsrp,
        lte_a2_threshold_rsrq=lte_a2_threshold_rsrq,
        lte_a2_threshold_rsrp_irat_volte=lte_a2_threshold_rsrp_irat_volte,
        lte_a2_threshold_rsrq_irat_volte=lte_a2_threshold_rsrq_irat_volte,
        a3_offset=a3_offset,
        a3_offset_anr=a3_offset_anr,
        a4_threshold_rsrp=a4_threshold_rsrp,
        lte_intra_a5_threshold_1_rsrp=lte_intra_a5_threshold_1_rsrp,
        lte_intra_a5_threshold_2_rsrp=lte_intra_a5_threshold_2_rsrp,
        b2_threshold1_rsrp=b2_threshold1_rsrp,
        b2_threshold2_rsrp=b2_threshold2_rsrp,
        b2_geran_irat_threshold=b2_geran_irat_threshold,
        qrxlevmin_sib1=qrxlevmin_sib1,
        qrxlevminoffset=qrxlevminoffset,
        s_intrasearch=s_intrasearch,
        s_nonintrasearch=s_nonintrasearch,
        qrxlevmin_sib3=qrxlevmin_sib3,
        reselection_priority=reselection_priority,
        threshservinglow=threshservinglow,
        x2_enable_disable=x2_enable_disable,
        ciphering_algorithm=ciphering_algorithm,
        integrity_algorithm=integrity_algorithm,
    )
    return single_enodeb_config


def _set_pci(
        cfg: EnodebConfiguration,
        pci: Any,
) -> None:
    """
    Set the following parameters:
     - PCI
    """
    if pci not in range(0, 504 + 1):
        raise ConfigurationError('Invalid PCI (%d)' % pci)
    cfg.set_parameter(ParameterName.PCI, pci)


def _set_bandwidth(
        cfg: EnodebConfiguration,
        data_model: DataModel,
        bandwidth_mhz: Any,
) -> None:
    """
    Set the following parameters:
     - DL bandwidth
     - UL bandwidth
    """
    _set_param_if_present(
        cfg, data_model, ParameterName.DL_BANDWIDTH,
        bandwidth_mhz,
    )
    _set_param_if_present(
        cfg, data_model, ParameterName.UL_BANDWIDTH,
        bandwidth_mhz,
    )


def _set_cell_id(
        cfg: EnodebConfiguration,
        cell_id: int,
) -> None:
    config_assert(
        cell_id in range(0, 268435456),
        'Cell Identity should be from 0 - (2^28 - 1)',
    )
    cfg.set_parameter(ParameterName.CELL_ID, cell_id)


def _set_tdd_subframe_config(
        device_cfg: EnodebConfiguration,
        cfg: EnodebConfiguration,
        subframe_assignment: Any,
        special_subframe_pattern: Any,
) -> None:
    """
    Set the following parameters:
     - Subframe assignment
     - Special subframe pattern
    """
    # Don't try to set if this is not TDD mode
    if (
        device_cfg.has_parameter(ParameterName.DUPLEX_MODE_CAPABILITY)
            and device_cfg.get_parameter(ParameterName.DUPLEX_MODE_CAPABILITY)
            != 'TDDMode'
    ):
        return

    config_assert(
        subframe_assignment in range(0, 6 + 1),
        'Invalid TDD subframe assignment (%d)' % subframe_assignment,
    )
    config_assert(
        special_subframe_pattern in range(0, 9 + 1),
        'Invalid TDD special subframe pattern (%d)'
        % special_subframe_pattern,
    )

    cfg.set_parameter(
        ParameterName.SUBFRAME_ASSIGNMENT,
        subframe_assignment,
    )
    cfg.set_parameter(
        ParameterName.SPECIAL_SUBFRAME_PATTERN,
        special_subframe_pattern,
    )


def _set_management_server(cfg: EnodebConfiguration) -> None:
    """
    Set the following parameters:
     - Periodic inform enable
     - Periodic inform interval (hard-coded)
    """
    cfg.set_parameter(ParameterName.PERIODIC_INFORM_ENABLE, True)
    # In seconds
    cfg.set_parameter(ParameterName.PERIODIC_INFORM_INTERVAL, 5)


def _set_s1_connection(
        cfg: EnodebConfiguration,
        mme_ip: Any,
        mme_port: Any = DEFAULT_S1_PORT,
) -> None:
    """
    Set the following parameters:
     - MME IP
     - MME port (defalts to 36412 as per TR-196 recommendation)
    """
    config_assert(type(mme_ip) == str, 'Invalid MME IP type')
    config_assert(type(mme_port) == int, 'Invalid MME Port type')
    cfg.set_parameter(ParameterName.MME_IP, mme_ip)
    cfg.set_parameter(ParameterName.MME_PORT, mme_port)


def _set_perf_mgmt(
        cfg: EnodebConfiguration,
        perf_mgmt_ip: str,
        perf_mgmt_port: int,
) -> None:
    """
    Set the following parameters:
     - Perf mgmt enable
     - Perf mgmt upload interval
     - Perf mgmt upload URL
    """
    cfg.set_parameter(ParameterName.PERF_MGMT_ENABLE, True)
    # Upload interval supported values (in secs):
    # [60, 300, 900, 1800, 3600]
    # Note: eNodeB crashes have been experienced with 60-sec interval.
    # Hence using 300sec
    cfg.set_parameter(
        ParameterName.PERF_MGMT_UPLOAD_INTERVAL,
        300,
    )
    cfg.set_parameter(
        ParameterName.PERF_MGMT_UPLOAD_URL,
        'http://%s:%d/' % (perf_mgmt_ip, perf_mgmt_port),
    )


def _set_misc_static_params(
        device_cfg: EnodebConfiguration,
        cfg: EnodebConfiguration,
        data_model: DataModel,
) -> None:
    """
    Set the following parameters:
     - Local gateway enable
     - GPS enable
    """
    _set_param_if_present(
        cfg, data_model, ParameterName.LOCAL_GATEWAY_ENABLE,
        0,
    )
    _set_param_if_present(cfg, data_model, ParameterName.GPS_ENABLE, True)
    # For BaiCells eNodeBs, IPSec enable may be either integer or bool.
    # Set to false/0 depending on the current type
    if data_model.is_parameter_present(ParameterName.IP_SEC_ENABLE):
        try:
            int(device_cfg.get_parameter(ParameterName.IP_SEC_ENABLE))
            cfg.set_parameter(ParameterName.IP_SEC_ENABLE, value=0)
        except ValueError:
            cfg.set_parameter(ParameterName.IP_SEC_ENABLE, value=False)

    _set_param_if_present(cfg, data_model, ParameterName.CELL_RESERVED, False)
    _set_param_if_present(
        cfg, data_model, ParameterName.MME_POOL_ENABLE,
        False,
    )


def _set_plmnids_tac(
        cfg: EnodebConfiguration,
        plmnids: Union[int, str],
        tac: Any,
) -> None:
    """
    Set the following parameters:
     - PLMNID list (including all child parameters)

    Input 'plmnids' is comma-separated list of PLMNIDs
    """
    # Convert int PLMNID to string
    if type(plmnids) == int:
        plmnid_str = str(plmnids)
    else:
        config_assert(type(plmnids) == str, 'PLMNID must be string')
        plmnid_str = plmnids

    # Multiple PLMNIDs will be supported using comma-separated list.
    # Currently, just one supported
    for char in plmnid_str:
        config_assert(
            char in '0123456789, ',
            'Unhandled character (%s) in PLMNID' % char,
        )
    plmnid_list = plmnid_str.split(',')

    # TODO - add support for multiple PLMNIDs
    config_assert(
        len(plmnid_list) == 1,
        'Exactly one PLMNID must be configured',
    )

    # Validate PLMNIDs
    plmnid_list[0] = plmnid_list[0].strip()
    config_assert(
        len(plmnid_list[0]) <= 6,
        'PLMNID must be length <=6 (%s)' % plmnid_list[0],
    )

    # We just need one PLMN element in the config. Delete all others.
    for i in range(1, 2):  # data_model.get_num_plmns() + 1):
        object_name = ParameterName.PLMN_N % i
        enable_plmn = i == 1
        cfg.add_object(object_name)
        cfg.set_parameter_for_object(
            ParameterName.PLMN_N_ENABLE % i,
            enable_plmn,
            object_name,
        )
        if enable_plmn:
            cfg.set_parameter_for_object(
                ParameterName.PLMN_N_CELL_RESERVED % i,
                False, object_name,
            )
            cfg.set_parameter_for_object(
                ParameterName.PLMN_N_PRIMARY % i,
                enable_plmn,
                object_name,
            )
            cfg.set_parameter_for_object(
                ParameterName.PLMN_N_PLMNID % i,
                plmnid_list[i - 1],
                object_name,
            )
    cfg.set_parameter(ParameterName.TAC, tac)
def _set_algorithm_a1_threshold_rsrp(
    cfg: EnodebConfiguration,
    a1_threshold_rsrp: int,
) -> None:
    """
    Set the following parameters:
     - a1_threshold_rsrp
    """
    cfg.set_parameter(ParameterName.A1_THRESHOLD_RSRP, a1_threshold_rsrp)

def _set_algorithm_lte_a1_threshold_rsrq(
    cfg: EnodebConfiguration,
    lte_a1_threshold_rsrq: int,
) -> None:
    """
    Set the following parameters:
     - lte_a1_threshold_rsrq
    """
    cfg.set_parameter(ParameterName.LTE_A1_THRESHOLD_RSRQ, lte_a1_threshold_rsrq)

def _set_algorithm_hysteresis(
    cfg: EnodebConfiguration,
    hysteresis: int,
) -> None:
    """
    Set the following parameters:
     - hysteresis
    """
    cfg.set_parameter(ParameterName.HYSTERESIS, hysteresis)

def _set_algorithm_time_to_trigger(
    cfg: EnodebConfiguration,
    time_to_trigger: Any,
) -> None:
    """
    Set the following parameters:
     - time_to_trigger
    """
    cfg.set_parameter(ParameterName.TIME_TO_TRIGGER, time_to_trigger)

def _set_algorithm_a2_threshold_rsrp(
    cfg: EnodebConfiguration,
    a2_threshold_rsrp: int,
) -> None:
    """
    Set the following parameters:
     - a2_threshold_rsrp
    """
    cfg.set_parameter(ParameterName.A2_THRESHOLD_RSRP, a2_threshold_rsrp)

def _set_algorithm_lte_a2_threshold_rsrq(
    cfg: EnodebConfiguration,
    lte_a2_threshold_rsrq: int,
) -> None:
    """
    Set the following parameters:
     - lte_a2_threshold_rsrq
    """
    cfg.set_parameter(ParameterName.LTE_A2_THRESHOLD_RSRQ, lte_a2_threshold_rsrq)

def _set_algorithm_lte_a2_threshold_rsrp_irat_volte(
    cfg: EnodebConfiguration,
    lte_a2_threshold_rsrp_irat_volte: int,
) -> None:
    """
    Set the following parameters:
     - lte_a2_threshold_rsrp_irat_volte
    """
    cfg.set_parameter(ParameterName.LTE_A2_THRESHOLD_RSRP_IRAT_VOLTE, lte_a2_threshold_rsrp_irat_volte)

def _set_algorithm_lte_a2_threshold_rsrq_irat_volte(
    cfg: EnodebConfiguration,
    lte_a2_threshold_rsrq_irat_volte: int,
) -> None:
    """
    Set the following parameters:
     - lte_a2_threshold_rsrq_irat_volte
    """
    cfg.set_parameter(ParameterName.LTE_A2_THRESHOLD_RSRQ_IRAT_VOLTE, lte_a2_threshold_rsrq_irat_volte)

def _set_algorithm_a3_offset(
    cfg: EnodebConfiguration,
    a3_offset: int,
) -> None:
    """
    Set the following parameters:
     - a3_offset
    """
    cfg.set_parameter(ParameterName.A3_OFFSET, a3_offset)

def _set_algorithm_a3_offset_anr(
    cfg: EnodebConfiguration,
    a3_offset_anr: int,
) -> None:
    """
    Set the following parameters:
     - a3_offset_anr
    """
    cfg.set_parameter(ParameterName.A3_OFFSET_ANR, a3_offset_anr)

def _set_algorithm_a4_threshold_rsrp(
    cfg: EnodebConfiguration,
    a4_threshold_rsrp: int,
) -> None:
    """
    Set the following parameters:
     - a4_threshold_rsrp
    """
    cfg.set_parameter(ParameterName.A4_THRESHOLD_RSRP, a4_threshold_rsrp)

def _set_algorithm_lte_intra_a5_threshold_1_rsrp(
    cfg: EnodebConfiguration,
    lte_intra_a5_threshold_1_rsrp: int,
) -> None:
    """
    Set the following parameters:
     - lte_intra_a5_threshold_1_rsrp
    """
    cfg.set_parameter(ParameterName.LTE_INTRA_A5_THRESHOLD_1_RSRP, lte_intra_a5_threshold_1_rsrp)

def _set_algorithm_lte_intra_a5_threshold_2_rsrp(
    cfg: EnodebConfiguration,
    lte_intra_a5_threshold_2_rsrp: int,
) -> None:
    """
    Set the following parameters:
     - lte_intra_a5_threshold_2_rsrp
    """
    cfg.set_parameter(ParameterName.LTE_INTRA_A5_THRESHOLD_2_RSRP, lte_intra_a5_threshold_2_rsrp)

def _set_algorithm_b2_threshold1_rsrp(
    cfg: EnodebConfiguration,
    b2_threshold1_rsrp: Any,
) -> None:
    """
    Set the following parameters:
     - b2_threshold1_rsrp
    """
    cfg.set_parameter(ParameterName.B2_THRESHOLD1_RSRP, b2_threshold1_rsrp)

def _set_algorithm_b2_threshold2_rsrp(
    cfg: EnodebConfiguration,
    b2_threshold2_rsrp: int,
) -> None:
    """
    Set the following parameters:
     - b2_threshold2_rsrp
    """
    cfg.set_parameter(ParameterName.B2_THRESHOLD2_RSRP, b2_threshold2_rsrp)

def _set_algorithm_b2_geran_irat_threshold(
    cfg: EnodebConfiguration,
    b2_geran_irat_threshold: Any,
) -> None:
    """
    Set the following parameters:
     - b2_geran_irat_threshold
    """
    cfg.set_parameter(ParameterName.B2_GERAN_IRAT_THRESHOLD, b2_geran_irat_threshold)

def _set_algorithm_qrxlevmin_sib1(
    cfg: EnodebConfiguration,
    qrxlevmin_sib1: int,
) -> None:
    """
    Set the following parameters:
     - qrxlevmin_sib1
    """
    cfg.set_parameter(ParameterName.QRXLEVMIN_SIB1, qrxlevmin_sib1)

def _set_algorithm_qrxlevminoffset(
    cfg: EnodebConfiguration,
    qrxlevminoffset: Any,
) -> None:
    """
    Set the following parameters:
     - qrxlevminoffset
    """
    cfg.set_parameter(ParameterName.QRXLEVMINOFFSET, qrxlevminoffset)

def _set_algorithm_s_intrasearch(
    cfg: EnodebConfiguration,
    s_intrasearch: Any,
) -> None:
    """
    Set the following parameters:
     - s_intrasearch
    """
    cfg.set_parameter(ParameterName.S_INTRASEARCH, s_intrasearch)

def _set_algorithm_s_nonintrasearch(
    cfg: EnodebConfiguration,
    s_nonintrasearch: Any,
) -> None:
    """
    Set the following parameters:
     - lte_a2_threshold_rsrq_irat_volte
    """
    cfg.set_parameter(ParameterName.S_NONINTRASEARCH, s_nonintrasearch)

def _set_algorithm_qrxlevmin_sib3(
    cfg: EnodebConfiguration,
    qrxlevmin_sib3: int,
) -> None:
    """
    Set the following parameters:
     - qrxlevmin_sib3
    """
    cfg.set_parameter(ParameterName.QRXLEVMIN_SIB3, qrxlevmin_sib3)

def _set_algorithm_reselection_priority(
    cfg: EnodebConfiguration,
    reselection_priority: Any,
) -> None:
    """
    Set the following parameters:
     - reselection_priority
    """
    cfg.set_parameter(ParameterName.RESELECTION_PRIORITY, reselection_priority)

def _set_algorithm_threshservinglow(
    cfg: EnodebConfiguration,
    threshservinglow: Any,
) -> None:
    """
    Set the following parameters:
     - threshservinglow
    """
    cfg.set_parameter(ParameterName.THRESHSERVINGLOW, threshservinglow)

def _set_algorithm_x2_enable_disable(
    cfg: EnodebConfiguration,
    x2_enable_disable: Any,
) -> None:
    """
    Set the following parameters:
     - x2_enable_disable
    """
    cfg.set_parameter(ParameterName.X2_ENABLE_DISABLE, x2_enable_disable)

def _set_algorithm_ciphering_algorithm(
    cfg: EnodebConfiguration,
    ciphering_algorithm: Any,
) -> None:
    """
    Set the following parameters:
     - ciphering_algorithm
    """
    cfg.set_parameter(ParameterName.CIPHERING_ALGORITHM, ciphering_algorithm)

def _set_algorithm_integrity_algorithm(
    cfg: EnodebConfiguration,
    integrity_algorithm: Any,
) -> None:
    """
    Set the following parameters:
     - integrity_algorithm
    """
    cfg.set_parameter(ParameterName.INTEGRITY_ALGORITHM, integrity_algorithm)

def _set_earfcn_freq_band_mode(
        device_cfg: EnodebConfiguration,
        cfg: EnodebConfiguration,
        data_model: DataModel,
        earfcndl: int,
) -> None:
    """
    Set the following parameters:
     - EARFCNDL
     - EARFCNUL
     - Band
    """
    # Note: validation of EARFCNDL done by mapping function. If invalid
    # EARFCN, raise ConfigurationError
    try:
        band, duplex_mode, earfcnul = map_earfcndl_to_band_earfcnul_mode(
            earfcndl,
        )
    except ValueError as err:
        raise ConfigurationError(err)

    # Verify capabilities
    if device_cfg.has_parameter(ParameterName.DUPLEX_MODE_CAPABILITY):
        duplex_capability = \
            device_cfg.get_parameter(ParameterName.DUPLEX_MODE_CAPABILITY)
        if duplex_mode == DuplexMode.TDD and duplex_capability != 'TDDMode':
            raise ConfigurationError((
                'eNodeB duplex mode capability is <{0}>, '
                'but earfcndl is <{1}>, giving duplex '
                'mode <{2}> instead'
            ).format(
                duplex_capability, str(earfcndl), str(duplex_mode),
            ))
        elif duplex_mode == DuplexMode.FDD and duplex_capability != 'FDDMode':
            raise ConfigurationError((
                'eNodeB duplex mode capability is <{0}>, '
                'but earfcndl is <{1}>, giving duplex '
                'mode <{2}> instead'
            ).format(
                duplex_capability, str(earfcndl), str(duplex_mode),
            ))
        elif duplex_mode not in {DuplexMode.TDD, DuplexMode.FDD}:
            raise ConfigurationError(
                'Invalid duplex mode (%s)' % str(duplex_mode),
            )

    if device_cfg.has_parameter(ParameterName.BAND_CAPABILITY):
        # Baicells indicated that they no longer use the band capability list,
        # so it may not be populated correctly
        band_capability_list = device_cfg.get_parameter(
            ParameterName.BAND_CAPABILITY,
        )
        band_capabilities = band_capability_list.split(',')
        if str(band) not in band_capabilities:
            logger.warning(
                'Band %d not in capabilities list (%s). Continuing'
                ' with config because capabilities list may not be'
                ' correct', band, band_capabilities,
            )
    cfg.set_parameter(ParameterName.EARFCNDL, earfcndl)
    if duplex_mode == DuplexMode.FDD:
        _set_param_if_present(
            cfg, data_model, ParameterName.EARFCNUL,
            earfcnul,
        )
    else:
        logger.debug('Not setting EARFCNUL - duplex mode is not FDD')

    _set_param_if_present(cfg, data_model, ParameterName.BAND, band)

    if duplex_mode == DuplexMode.TDD:
        logger.debug('Set EARFCNDL=%d, Band=%d', earfcndl, band)
    elif duplex_mode == DuplexMode.FDD:
        logger.debug(
            'Set EARFCNDL=%d, EARFCNUL=%d, Band=%d',
            earfcndl, earfcnul, band,
        )


def _set_param_if_present(
        cfg: EnodebConfiguration,
        data_model: DataModel,
        param: ParameterName,
        value: Any,
) -> None:
    if data_model.is_parameter_present(param):
        cfg.set_parameter(param, value)
