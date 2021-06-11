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
import os
from collections import namedtuple
from typing import Any, Dict, List, NamedTuple, Optional, Tuple, Union

from lte.protos.enodebd_pb2 import SingleEnodebStatus
from lte.protos.mconfig import mconfigs_pb2
from magma.common import serialization_utils
from magma.enodebd import metrics
from magma.enodebd.data_models.data_model_parameters import ParameterName
from magma.enodebd.device_config.configuration_util import (
    find_enb_by_cell_id,
    get_enb_rf_tx_desired,
)
from magma.enodebd.exceptions import ConfigurationError
from magma.enodebd.logger import EnodebdLogger as logger
from magma.enodebd.s1ap_client import get_all_enb_state
from magma.enodebd.state_machines.enb_acs import EnodebAcsStateMachine
from magma.enodebd.state_machines.enb_acs_manager import StateMachineManager
from orc8r.protos.service303_pb2 import State

# There are 2 levels of caching for GPS coordinates from the enodeB: module
# variables (in-memory) and on disk. In the event the enodeB stops reporting
# GPS, we will continue to report the cached coordinates from the in-memory
# cached coordinates. If enodebd is restarted, this in-memory cache will be
# populated by the file

CACHED_GPS_COORD_FILE_PATH = os.path.join(
    '/var/opt/magma/enodebd',
    'gps_coords.txt',
)

# Cache GPS coordinates in memory so we don't write to the file cache if the
# coordinates have not changed. We can read directly from here instead of the
# file cache when the enodeB goes down unless these are unintialized.
_gps_lat_cached = None
_gps_lon_cached = None

EnodebStatus = NamedTuple(
    'EnodebStatus',
    [
        ('enodeb_configured', bool),
        ('gps_latitude', str),
        ('gps_longitude', str),
        ('enodeb_connected', bool),
        ('opstate_enabled', bool),
        ('rf_tx_on', bool),
        ('rf_tx_desired', bool),
        ('gps_connected', bool),
        ('ptp_connected', bool),
        ('mme_connected', bool),
        ('fsm_state', str),
        ('cell_id', int),
    ],
)

# TODO: Remove after checkins support multiple eNB status
MagmaOldEnodebdStatus = namedtuple(
    'MagmaOldEnodebdStatus',
    [
        'enodeb_serial',
        'enodeb_configured',
        'gps_latitude',
        'gps_longitude',
        'enodeb_connected',
        'opstate_enabled',
        'rf_tx_on',
        'rf_tx_desired',
        'gps_connected',
        'ptp_connected',
        'mme_connected',
        'enodeb_state',
    ],
)

MagmaEnodebdStatus = NamedTuple(
    'MagmaEnodebdStatus',
    [
        ('n_enodeb_connected', str),
        ('all_enodeb_configured', str),
        ('all_enodeb_opstate_enabled', str),
        ('all_enodeb_rf_tx_configured', str),
        ('any_enodeb_gps_connected', str),
        ('all_enodeb_ptp_connected', str),
        ('all_enodeb_mme_connected', str),
        ('gateway_gps_longitude', str),
        ('gateway_gps_latitude', str),
    ],
)


def update_status_metrics(status: EnodebStatus) -> None:
    """ Update metrics for eNodeB status """
    # Call every second
    metrics_by_stat_key = {
        'enodeb_connected': metrics.STAT_ENODEB_CONNECTED,
        'enodeb_configured': metrics.STAT_ENODEB_CONFIGURED,
        'opstate_enabled': metrics.STAT_OPSTATE_ENABLED,
        'rf_tx_on': metrics.STAT_RF_TX_ENABLED,
        'rf_tx_desired': metrics.STAT_RF_TX_DESIRED,
        'gps_connected': metrics.STAT_GPS_CONNECTED,
        'ptp_connected': metrics.STAT_PTP_CONNECTED,
        'mme_connected': metrics.STAT_MME_CONNECTED,
    }

    def get_metric_value(enodeb_status: Dict[str, str], key: str):
        # Metrics are "sticky" when synced to the cloud - if we don't
        # receive a status update from enodeb, set the metric to 0
        # to explicitly indicate that it was not received, otherwise the
        # metrics collector will continue to report the last value
        val = enodeb_status.get(key, None)
        if val is None:
            return 0
        if type(val) is not bool:
            logger.error('Could not cast metric value %s to int', val)
            return 0
        return int(val)  # val should be either True or False

    for stat_key, metric in metrics_by_stat_key.items():
        metric.set(get_metric_value(status._asdict(), stat_key))


# TODO: Remove after checkins support multiple eNB status
def get_service_status_old(
        enb_acs_manager: StateMachineManager,
) -> Dict[str, Any]:
    """ Get service status compatible with older controller """
    enb_status_by_serial = get_all_enb_status(enb_acs_manager)
    # Since we only expect users to plug in a single eNB, generate service
    # status with the first one we find that is connected
    for enb_serial, enb_status in enb_status_by_serial.items():
        if enb_status.enodeb_connected:
            return MagmaOldEnodebdStatus(
                enodeb_serial=enb_serial,
                enodeb_configured=_bool_to_str(enb_status.enodeb_configured),
                gps_latitude=enb_status.gps_latitude,
                gps_longitude=enb_status.gps_longitude,
                enodeb_connected=_bool_to_str(enb_status.enodeb_connected),
                opstate_enabled=_bool_to_str(enb_status.opstate_enabled),
                rf_tx_on=_bool_to_str(enb_status.rf_tx_on),
                rf_tx_desired=_bool_to_str(enb_status.rf_tx_desired),
                gps_connected=_bool_to_str(enb_status.gps_connected),
                ptp_connected=_bool_to_str(enb_status.ptp_connected),
                mme_connected=_bool_to_str(enb_status.mme_connected),
                enodeb_state=enb_status.fsm_state,
            )._asdict()
    return MagmaOldEnodebdStatus(
        enodeb_serial='N/A',
        enodeb_configured='0',
        gps_latitude='0.0',
        gps_longitude='0.0',
        enodeb_connected='0',
        opstate_enabled='0',
        rf_tx_on='0',
        rf_tx_desired='N/A',
        gps_connected='0',
        ptp_connected='0',
        mme_connected='0',
        enodeb_state='N/A',
    )._asdict()


def get_service_status(enb_acs_manager: StateMachineManager) -> Dict[str, Any]:
    enodebd_status = _get_enodebd_status(enb_acs_manager)
    return enodebd_status._asdict()


def _get_enodebd_status(
        enb_acs_manager: StateMachineManager,
) -> MagmaEnodebdStatus:
    enb_status_by_serial = get_all_enb_status(enb_acs_manager)
    # Start from default values for enodebd status
    n_enodeb_connected = 0
    all_enodeb_configured = False
    all_enodeb_opstate_enabled = False
    all_enodeb_rf_tx_configured = False
    any_enodeb_gps_connected = False
    all_enodeb_ptp_connected = False
    all_enodeb_mme_connected = False
    gateway_gps_longitude = '0.0'
    gateway_gps_latitude = '0.0'

    def _is_rf_tx_configured(enb_status: EnodebStatus) -> bool:
        return enb_status.rf_tx_on == enb_status.rf_tx_desired

    if enb_status_by_serial:
        enb_status_list = list(enb_status_by_serial.values())
        # Aggregate all eNB status for enodebd status, repetitive but
        # clearer for output purposes.
        n_enodeb_connected = sum(
            enb_status.enodeb_connected for enb_status in enb_status_list
        )
        all_enodeb_configured = all(
            enb_status.enodeb_configured for enb_status in enb_status_list
        )
        all_enodeb_mme_connected = all(
            enb_status.mme_connected for enb_status in enb_status_list
        )
        all_enodeb_opstate_enabled = all(
            enb_status.opstate_enabled for enb_status in enb_status_list
        )
        all_enodeb_ptp_connected = all(
            enb_status.ptp_connected for enb_status in enb_status_list
        )
        any_enodeb_gps_connected = any(
            enb_status.gps_connected for enb_status in enb_status_list
        )
        all_enodeb_rf_tx_configured = all(
            _is_rf_tx_configured(enb_status) for enb_status in enb_status_list
        )
        if n_enodeb_connected:
            gateway_gps_longitude = enb_status_list[0].gps_longitude
            gateway_gps_latitude = enb_status_list[0].gps_latitude

    return MagmaEnodebdStatus(
        n_enodeb_connected=str(n_enodeb_connected),
        all_enodeb_configured=str(all_enodeb_configured),
        all_enodeb_opstate_enabled=str(all_enodeb_opstate_enabled),
        all_enodeb_rf_tx_configured=str(all_enodeb_rf_tx_configured),
        any_enodeb_gps_connected=str(any_enodeb_gps_connected),
        all_enodeb_ptp_connected=str(all_enodeb_ptp_connected),
        all_enodeb_mme_connected=str(all_enodeb_mme_connected),
        gateway_gps_longitude=str(gateway_gps_longitude),
        gateway_gps_latitude=str(gateway_gps_latitude),
    )


def get_all_enb_status(
        enb_acs_manager: StateMachineManager,
) -> Dict[str, EnodebStatus]:
    enb_status_by_serial = {}
    serial_list = enb_acs_manager.get_connected_serial_id_list()
    for enb_serial in serial_list:
        handler = enb_acs_manager.get_handler_by_serial(enb_serial)
        status = get_enb_status(handler)
        enb_status_by_serial[enb_serial] = status

    return enb_status_by_serial


def get_enb_status(enodeb: EnodebAcsStateMachine) -> EnodebStatus:
    """
    Returns a dict representing the status of an enodeb

    The returned dictionary will be a subset of the following keys:
        - enodeb_connected
        - enodeb_configured
        - opstate_enabled
        - rf_tx_on
        - rf_tx_desired
        - gps_connected
        - ptp_connected
        - mme_connected
        - gps_latitude
        - gps_longitude
        - ip_address
        - cell_id

    The set of keys returned will depend on the connection status of the
    enodeb. A missing key indicates that the value is unknown.

    Returns:
        Status dictionary for the enodeb state
    """
    enodeb_configured = enodeb.is_enodeb_configured()

    # We cache GPS coordinates so try to read them before the early return
    # if the enB is not connected
    gps_lat, gps_lon = _get_and_cache_gps_coords(enodeb)

    enodeb_connected = enodeb.is_enodeb_connected()
    opstate_enabled = _parse_param_as_bool(enodeb, ParameterName.OP_STATE)
    rf_tx_on = _parse_param_as_bool(enodeb, ParameterName.RF_TX_STATUS)
    rf_tx_on = rf_tx_on and enodeb_connected
    try:
        enb_serial = \
            enodeb.device_cfg.get_parameter(ParameterName.SERIAL_NUMBER)
        enb_cell_id = int(
            enodeb.device_cfg.get_parameter(ParameterName.CELL_ID),
        )
        rf_tx_desired = get_enb_rf_tx_desired(enodeb.mconfig, enb_serial)
    except (KeyError, ConfigurationError):
        rf_tx_desired = False
        enb_cell_id = 0

    mme_connected = _parse_param_as_bool(enodeb, ParameterName.MME_STATUS)
    gps_connected = _get_gps_status_as_bool(enodeb)
    try:
        ptp_connected = _parse_param_as_bool(enodeb, ParameterName.PTP_STATUS)
    except ConfigurationError:
        ptp_connected = False

    return EnodebStatus(
        enodeb_configured=enodeb_configured,
        gps_latitude=gps_lat,
        gps_longitude=gps_lon,
        enodeb_connected=enodeb_connected,
        opstate_enabled=opstate_enabled,
        rf_tx_on=rf_tx_on,
        rf_tx_desired=rf_tx_desired,
        gps_connected=gps_connected,
        ptp_connected=ptp_connected,
        mme_connected=mme_connected,
        fsm_state=enodeb.get_state(),
        cell_id=enb_cell_id,
    )


def get_single_enb_status(
        device_serial: str,
        state_machine_manager: StateMachineManager,
) -> SingleEnodebStatus:
    try:
        handler = state_machine_manager.get_handler_by_serial(device_serial)
    except KeyError:
        return _empty_enb_status()

    # This namedtuple is missing IP and serial info
    status = get_enb_status(handler)

    # Get IP info
    ip = state_machine_manager.get_ip_of_serial(device_serial)

    def get_status_property(status: bool) -> SingleEnodebStatus.StatusProperty:
        if status:
            return SingleEnodebStatus.StatusProperty.Value('ON')
        return SingleEnodebStatus.StatusProperty.Value('OFF')

    # Build the message to return through gRPC
    enb_status = SingleEnodebStatus()
    enb_status.device_serial = device_serial
    enb_status.ip_address = ip
    enb_status.connected = get_status_property(status.enodeb_connected)
    enb_status.configured = get_status_property(status.enodeb_configured)
    enb_status.opstate_enabled = get_status_property(status.opstate_enabled)
    enb_status.rf_tx_on = get_status_property(status.rf_tx_on)
    enb_status.rf_tx_desired = get_status_property(status.rf_tx_desired)
    enb_status.gps_connected = get_status_property(status.gps_connected)
    enb_status.ptp_connected = get_status_property(status.ptp_connected)
    enb_status.mme_connected = get_status_property(status.mme_connected)
    enb_status.gps_longitude = status.gps_longitude
    enb_status.gps_latitude = status.gps_latitude
    enb_status.fsm_state = status.fsm_state
    return enb_status


def get_operational_states(
    enb_acs_manager: StateMachineManager,
    mconfig: mconfigs_pb2.EnodebD,
) -> List[State]:
    """
    Returns: A list of State with EnodebStatus encoded as JSON
    """
    states = []
    configured_serial_ids = []
    enb_status_by_serial = get_all_enb_status(enb_acs_manager)

    # Get S1 connected eNBs
    enb_statuses = get_all_enb_state()

    for serial_id in enb_status_by_serial:
        enb_status_dict = enb_status_by_serial[serial_id]._asdict()

        # Add IP address to state
        enb_status_dict['ip_address'] = enb_acs_manager.get_ip_of_serial(
            serial_id,
        )

        # Add num of UEs connected
        num_ue_connected = enb_statuses.get(enb_status_dict['cell_id'], 0)
        enb_status_dict['ues_connected'] = num_ue_connected

        serialized = json.dumps(enb_status_dict)
        state = State(
            type="single_enodeb",
            deviceID=serial_id,
            value=serialized.encode('utf-8'),
        )
        configured_serial_ids.append(serial_id)
        states.append(state)

    # Get state for externally configured enodebs
    s1_states = get_enb_s1_connected_states(
        enb_statuses,
        configured_serial_ids,
        mconfig,
    )
    states.extend(s1_states)

    return states


def get_enb_s1_connected_states(
    enb_s1_state_map, configured_serial_ids,
    mconfig,
) -> List[State]:
    states = []
    for enb_id in enb_s1_state_map:
        enb = find_enb_by_cell_id(mconfig, enb_id)
        if enb and enb.serial_num not in configured_serial_ids:
            status = EnodebStatus(
                enodeb_configured=False,
                gps_latitude='N/A',
                gps_longitude='N/A',
                enodeb_connected=True,
                opstate_enabled=False,
                rf_tx_on=False,
                rf_tx_desired=False,
                gps_connected=False,
                ptp_connected=False,
                mme_connected=True,
                fsm_state='N/A',
                cell_id=enb_id,
            )
            status_dict = status._asdict()

            # Add IP address to state
            status_dict['ip_address'] = enb.config.ip_address

            # Add num of UEs connected to state, use cellID from mconfig
            status_dict['ues_connected'] = enb_s1_state_map.get(enb_id, 0)

            serialized = json.dumps(status_dict)
            state = State(
                type="single_enodeb",
                deviceID=enb.serial_num,
                value=serialized.encode('utf-8'),
            )
            states.append(state)
    return states


def _empty_enb_status() -> SingleEnodebStatus:
    enb_status = SingleEnodebStatus()
    enb_status.device_serial = 'N/A'
    enb_status.ip_address = 'N/A'
    enb_status.connected = '0'
    enb_status.configured = '0'
    enb_status.opstate_enabled = '0'
    enb_status.rf_tx_on = '0'
    enb_status.rf_tx_desired = 'N/A'
    enb_status.gps_connected = '0'
    enb_status.ptp_connected = '0'
    enb_status.mme_connected = '0'
    enb_status.gps_longitude = '0.0'
    enb_status.gps_latitude = '0.0'
    enb_status.fsm_state = 'N/A'
    return enb_status


def _parse_param_as_bool(
        enodeb: EnodebAcsStateMachine,
        param_name: ParameterName,
) -> bool:
    try:
        return _format_as_bool(enodeb.get_parameter(param_name), param_name)
    except (KeyError, ConfigurationError):
        return False


def _format_as_bool(
        param_value: Union[bool, str, int],
        param_name: Optional[Union[ParameterName, str]] = None,
) -> bool:
    """ Returns '1' for true, and '0' for false """
    stripped_value = str(param_value).lower().strip()
    if stripped_value in {'true', '1'}:
        return True
    elif stripped_value in {'false', '0'}:
        return False
    else:
        logger.warning(
            '%s parameter not understood (%s)', param_name, param_value,
        )
        return False


def _get_gps_status_as_bool(enodeb: EnodebAcsStateMachine) -> bool:
    try:
        if not enodeb.has_parameter(ParameterName.GPS_STATUS):
            return False
        else:
            param = enodeb.get_parameter(ParameterName.GPS_STATUS)
            stripped_value = param.lower().strip()
            if stripped_value == '0' or stripped_value == '2':
                # 2 = GPS locking
                return False
            elif stripped_value == '1':
                return True
            else:
                logger.warning(
                    'GPS status parameter not understood (%s)', param,
                )
                return False
    except (KeyError, ConfigurationError):
        return False


def _get_and_cache_gps_coords(enodeb: EnodebAcsStateMachine) -> Tuple[
    str, str,
]:
    """
    Read the GPS coordinates of the enB from its configuration or the
    cached coordinate file if the preceding read fails. If reading from
    enB configuration succeeds, this method will cache the new coordinates.

    Returns:
        (str, str): GPS latitude, GPS longitude
    """
    lat, lon = '', ''
    try:
        lat = enodeb.get_parameter(ParameterName.GPS_LAT)
        lon = enodeb.get_parameter(ParameterName.GPS_LONG)

        if lat != _gps_lat_cached or lon != _gps_lon_cached:
            _cache_new_gps_coords(lat, lon)
        return lat, lon
    except (KeyError, ConfigurationError):
        return _get_cached_gps_coords()
    except ValueError:
        logger.warning('GPS lat/long not understood (%s/%s)', lat, lon)
        return '0', '0'


def _get_cached_gps_coords() -> Tuple[str, str]:
    """
    Returns cached GPS coordinates if enB is disconnected or otherwise not
    reporting coordinates.

    Returns:
        (str, str): (GPS lat, GPS lon)
    """
    # pylint: disable=global-statement
    global _gps_lat_cached, _gps_lon_cached
    if _gps_lat_cached is None or _gps_lon_cached is None:
        _gps_lat_cached, _gps_lon_cached = _read_gps_coords_from_file()
    return _gps_lat_cached, _gps_lon_cached


def _read_gps_coords_from_file():
    try:
        with open(CACHED_GPS_COORD_FILE_PATH) as f:
            lines = f.readlines()
            if len(lines) != 2:
                logger.warning(
                    'Expected to find 2 lines in GPS '
                    'coordinate file but only found %d',
                    len(lines),
                )
                return '0', '0'
            return tuple(map(lambda l: l.strip(), lines))
    except OSError:
        logger.warning('Could not open cached GPS coordinate file')
        return '0', '0'


def _cache_new_gps_coords(gps_lat, gps_lon):
    """
    Cache GPS coordinates in the module-level variables here and write them
    to a managed file on disk.

    Args:
        gps_lat (str): latitude as a string
        gps_lon (str): longitude as a string
    """
    # pylint: disable=global-statement
    global _gps_lat_cached, _gps_lon_cached
    _gps_lat_cached, _gps_lon_cached = gps_lat, gps_lon
    _write_gps_coords_to_file(gps_lat, gps_lon)


def _write_gps_coords_to_file(gps_lat, gps_lon):
    lines = '{lat}\n{lon}'.format(lat=gps_lat, lon=gps_lon)
    try:
        serialization_utils.write_to_file_atomically(
            CACHED_GPS_COORD_FILE_PATH,
            lines,
        )
    except OSError:
        pass


def _bool_to_str(b: bool) -> str:
    if b is True:
        return "1"
    return "0"
