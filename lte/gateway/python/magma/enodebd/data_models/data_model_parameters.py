"""
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
"""
class ParameterName():
    # Top-level objects
    DEVICE = 'Device'
    FAP_SERVICE = 'FAPService'

    # Device info parameters
    GPS_STATUS = 'GPS status'
    PTP_STATUS = 'PTP status'
    MME_STATUS = 'MME status'
    REM_STATUS = 'REM status'

    LOCAL_GATEWAY_ENABLE = 'Local gateway enable'
    GPS_ENABLE = 'GPS enable'
    GPS_LAT = 'GPS lat'
    GPS_LONG = 'GPS long'
    SW_VERSION = 'SW version'

    # Capabilities
    DUPLEX_MODE_CAPABILITY = 'Duplex mode capability'
    BAND_CAPABILITY = 'Band capability'

    # RF-related parameters
    EARFCNDL = 'EARFCNDL'
    EARFCNUL = 'EARFCNUL'
    BAND = 'Band'
    PCI = 'PCI'
    DL_BANDWIDTH = 'DL bandwidth'
    UL_BANDWIDTH = 'UL bandwidth'
    SUBFRAME_ASSIGNMENT = 'Subframe assignment'
    SPECIAL_SUBFRAME_PATTERN = 'Special subframe pattern'

    # Other LTE parameters
    ADMIN_STATE = 'Admin state'
    OP_STATE = 'Opstate'
    RF_TX_STATUS = 'RF TX status'

    # RAN parameters
    CELL_RESERVED = 'Cell reserved'
    CELL_BARRED = 'Cell barred'

    # Core network parameters
    MME_IP = 'MME IP'
    MME_PORT = 'MME port'
    NUM_PLMNS = 'Num PLMNs'
    PLMN = 'PLMN'

    # PLMN parameters
    PLMN_N = 'PLMN %d'
    PLMN_N_CELL_RESERVED = 'PLMN %d cell reserved'
    PLMN_N_ENABLE = 'PLMN %d enable'
    PLMN_N_PRIMARY = 'PLMN %d primary'
    PLMN_N_PLMNID = 'PLMN %d PLMNID'

    # PLMN arrays are added below
    TAC = 'TAC'
    IP_SEC_ENABLE = 'IPSec enable'
    MME_POOL_ENABLE = 'MME pool enable'

    # Management server parameters
    PERIODIC_INFORM_ENABLE = 'Periodic inform enable'
    PERIODIC_INFORM_INTERVAL = 'Periodic inform interval'

    # Performance management parameters
    PERF_MGMT_ENABLE = 'Perf mgmt enable'
    PERF_MGMT_UPLOAD_INTERVAL = 'Perf mgmt upload interval'
    PERF_MGMT_UPLOAD_URL = 'Perf mgmt upload URL'


class TrParameterType():
    BOOLEAN = 'boolean'
    STRING = 'string'
    INT = 'int'
    UNSIGNED_INT = 'unsignedInt'
    OBJECT = 'object'
