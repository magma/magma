"""
Copyright (c) 2018-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from typing import Any as TAny, Dict

# Import to register all mconfig protobufs in symbol database
from google.protobuf import symbol_database
from google.protobuf.internal.well_known_types import Any

from magma.configuration import service_configs
from magma.configuration.exceptions import LoadConfigError


def filter_configs_by_key(configs_by_key: Dict[str, TAny]) -> Dict[str, TAny]:
    """
    Given a JSON-deserialized map of mconfig protobuf Any's keyed by service
    name, filter out any entires without a corresponding service or which have
    values that aren't registered in the protobuf symbol database yet.

    Args:
        configs_by_key:
            JSON-deserialized service mconfigs keyed by service name

    Returns:
        The input map without any services which currently don't exist or have
        types which are not in the protobuf type registry.
    """
    services = service_configs.get_service_config_value(
        'magmad',
        'magma_services', [],
    )
    services.append('magmad')
    services = set(services)

    filtered_configs_by_key = {}
    for srv, cfg in configs_by_key.items():
        if srv not in services:
            continue

        try:
            type_name = cfg['@type'].split('/')[-1]
            symbol_database.Default().GetSymbol(type_name)
        except KeyError:
            continue
        filtered_configs_by_key[srv] = cfg
    return filtered_configs_by_key


def unpack_mconfig_any(mconfig_any: Any) -> TAny:
    """
    Unpack a protobuf Any type into its concrete protobuf type.

    Args:
        mconfig_any: protobuf Any type to unpack

    Returns: Concrete protobuf object that the provided Any wraps
    """
    type_name = mconfig_any.TypeName()
    try:
        msg = symbol_database.Default().GetSymbol(type_name)()
    except KeyError as e:
        raise LoadConfigError(
            'Mconfig proto type %s not found' % type_name,
        ) from e
    mconfig_any.Unpack(msg)
    return msg
