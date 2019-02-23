"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import logging
import os
from collections import namedtuple
from typing import Any, Dict
from magma.common import serialization_utils
from magma.enodebd import metrics
from magma.enodebd.data_models.data_model_parameters import ParameterName
from magma.enodebd.exceptions import ConfigurationError


# There are 2 levels of caching for GPS coordinates from the enodeB: module
# variables (in-memory) and on disk. In the event the enodeB stops reporting
# GPS, we will continue to report the cached coordinates from the in-memory
# cached coordinates. If enodebd is restarted, this in-memory cache will be
# populated by the file
from magma.enodebd.state_machines.enb_acs import EnodebAcsStateMachine

CACHED_GPS_COORD_FILE_PATH = os.path.join(
    '/var/opt/magma/enodebd',
    'gps_coords.txt',
)

# Cache GPS coordinates in memory so we don't write to the file cache if the
# coordinates have not changed. We can read directly from here instead of the
# file cache when the enodeB goes down unless these are unintialized.
_gps_lat_cached = None
_gps_lon_cached = None

EnodebStatus = namedtuple('EnodebStatus',
                          ['enodeb_configured', 'gps_latitude',
                           'gps_longitude', 'enodeb_connected',
                           'opstate_enabled', 'rf_tx_on', 'gps_connected',
                           'ptp_connected', 'mme_connected', 'enodeb_state'])


def update_status_metrics(status: EnodebStatus) -> None:
    """ Update metrics for eNodeB status """
    # Call every second
    metrics_by_stat_key = {
        'enodeb_connected': metrics.STAT_ENODEB_CONNECTED,
        'enodeb_configured': metrics.STAT_ENODEB_CONFIGURED,
        'opstate_enabled': metrics.STAT_OPSTATE_ENABLED,
        'rf_tx_on': metrics.STAT_RF_TX_ENABLED,
        'gps_connected': metrics.STAT_GPS_CONNECTED,
        'ptp_connected': metrics.STAT_PTP_CONNECTED,
        'mme_connected': metrics.STAT_MME_CONNECTED,
    }

    def get_metric_value(enodeb_status, key):
        # Metrics are "sticky" when synced to the cloud - if we don't
        # receive a status update from enodeb, set the metric to 0
        # to explicitly indicate that it was not received, otherwise the
        # metrics collector will continue to report the last value
        if key not in enodeb_status:
            return 0

        try:
            return int(enodeb_status[key])
        except ValueError:
            logging.error('Could not cast metric value %s to int',
                          enodeb_status[key])
            return 0

    for stat_key, metric in metrics_by_stat_key.items():
        metric.set(get_metric_value(status, stat_key))


def get_enodeb_status(enodeb: EnodebAcsStateMachine) -> Dict[str, Any]:
    """
    Returns a dict representing the status of an enodeb

    The returned dictionary will be a subset of the following keys:
        - enodeb_connected
        - enodeb_configured
        - opstate_enabled
        - rf_tx_on
        - gps_connected
        - ptp_connected
        - mme_connected
        - gps_latitude
        - gps_longitude

    The set of keys returned will depend on the connection status of the
    enodeb. A missing key indicates that the value is unknown.

    Returns:
        Status dictionary for the enodeb state
    """
    enodeb_configured = '1' if enodeb.is_enodeb_configured() else '0'

    # We cache GPS coordinates so try to read them before the early return
    # if the enB is not connected
    gps_lat, gps_lon = _get_and_cache_gps_coords(enodeb)

    enodeb_connected = '1' if enodeb.is_enodeb_connected() else '0'
    opstate_enabled = _parse_param_as_bool(enodeb, ParameterName.OP_STATE)
    rf_tx_on = _parse_param_as_bool(enodeb, ParameterName.RF_TX_STATUS)
    mme_connected = _parse_param_as_bool(enodeb, ParameterName.MME_STATUS)

    try:
        param = enodeb.get_parameter(ParameterName.GPS_STATUS)
        pval = param.lower().strip()
        if pval == '0' or pval == '1':
            gps_connected = pval
        elif pval == '2':
            # 2 = GPS locking
            gps_connected = '0'
        else:
            logging.warning(
                'GPS status parameter not understood (%s)', param)
            gps_connected = '0'
    except (KeyError, ConfigurationError):
        gps_connected = '0'

    try:
        param = enodeb.get_parameter(ParameterName.PTP_STATUS)
        pval = param.lower().strip()
        if pval == '0' or pval == '1':
            ptp_connected = pval
        else:
            logging.warning(
                'PTP status parameter not understood (%s)', param)
            ptp_connected = '0'
    except (KeyError, ConfigurationError):
        ptp_connected = '0'

    return EnodebStatus(enodeb_configured=enodeb_configured,
                        gps_latitude=gps_lat,
                        gps_longitude=gps_lon,
                        enodeb_connected=enodeb_connected,
                        opstate_enabled=opstate_enabled,
                        rf_tx_on=rf_tx_on,
                        gps_connected=gps_connected,
                        ptp_connected=ptp_connected,
                        mme_connected=mme_connected,
                        enodeb_state=enodeb.get_state())._asdict()


def _parse_param_as_bool(
    enodeb: EnodebAcsStateMachine,
    param_name: ParameterName
) -> str:
    """
    Returns '1' for true, and '0' for false
    """
    try:
        param = enodeb.get_parameter(param_name)
        pval = param.lower().strip()
        if pval in {'true', '1'}:
            return '1'
        elif pval in {'false', '0'}:
            return '0'
        else:
            logging.warning(
                '%s parameter not understood (%s)', param_name, param)
            return '0'
    except (KeyError, ConfigurationError):
        return '0'


def _get_and_cache_gps_coords(enodeb: EnodebAcsStateMachine) -> tuple:
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
        logging.warning('GPS lat/long not understood (%s/%s)', lat, lon)
        return '0', '0'


def _get_cached_gps_coords():
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
                logging.warning('Expected to find 2 lines in GPS '
                                'coordinate file but only found %d',
                                len(lines))
                return '0', '0'
            return tuple(map(lambda l: l.strip(), lines))
    except OSError:
        logging.warning('Could not open cached GPS coordinate file')
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
