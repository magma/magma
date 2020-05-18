#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from typing import Dict, Iterator, MutableMapping, TypeVar

from ..exceptions import EntityNotFoundError
from .data_class import EquipmentPortType, EquipmentType, LocationType, ServiceType
from .data_enum import Entity


T = TypeVar("T")


class Cache(MutableMapping[str, T]):
    def __init__(self, entity: Entity) -> None:
        self.store: Dict[str, T] = {}
        self.entity = entity

    def __getitem__(self, key: str) -> T:
        if key not in self.store:
            raise EntityNotFoundError(entity=self.entity, entity_name=key)
        return self.store[self.__keytransform__(key)]

    def __setitem__(self, key: str, value: T) -> None:
        self.store[self.__keytransform__(key)] = value

    def __delitem__(self, key: str) -> None:
        del self.store[self.__keytransform__(key)]

    def __iter__(self) -> Iterator[str]:
        return iter(self.store)

    def __len__(self) -> int:
        return len(self.store)

    def __keytransform__(self, key: str) -> str:
        return key


# pyre-fixme[16]: `Cache` has no attribute `__getitem__`.
LOCATION_TYPES = Cache[LocationType](Entity.LocationType)
EQUIPMENT_TYPES = Cache[EquipmentType](Entity.EquipmentType)
SERVICE_TYPES = Cache[ServiceType](Entity.ServiceType)
PORT_TYPES = Cache[EquipmentPortType](Entity.EquipmentPortType)


def clear_types() -> None:
    LOCATION_TYPES.clear()
    EQUIPMENT_TYPES.clear()
    SERVICE_TYPES.clear()
    PORT_TYPES.clear()
