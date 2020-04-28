#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from typing import Optional

from .common.data_enum import Entity


class CustomException(Exception):
    pass


def get_location_id_msg(
    location_name: Optional[str] = None,
    location_type: Optional[str] = None,
    external_id: Optional[str] = None,
) -> str:
    msg = "location"
    if location_name:
        msg = msg + f" {location_name}"
    if location_type:
        msg = msg + f" of type {location_type}"
    if external_id:
        msg = msg + f" with external id {external_id}"
    return msg


class EntityNotFoundError(Exception):
    def __init__(
        self,
        entity: Entity,
        entity_id: Optional[str] = None,
        entity_name: Optional[str] = None,
        msg: Optional[str] = None,
    ) -> None:
        self.entity: Entity = entity
        self.entity_id: str = entity_id if entity_id is not None else ""
        self.entity_name: str = entity_name if entity_name is not None else ""
        self.msg: str = msg if msg is not None else ""

    def __str__(self) -> str:
        message = f"{self.entity} entity"
        if self.msg:
            message = "\n".join([self.msg, f"{self.entity} entity"])
        if self.entity_name:
            message = " ".join([message, f"<enitity_name: {self.entity_name}>"])
        if self.entity_id:
            message = " ".join([message, f"<entity_id: {self.entity_id}>"])
        return " ".join([message, "does not found in inventory"])


class LocationNotFoundException(CustomException):
    def __init__(
        self,
        location_name: Optional[str] = None,
        location_type: Optional[str] = None,
        external_id: Optional[str] = None,
    ) -> None:
        self.location_name: Optional[str] = location_name
        self.location_type: Optional[str] = location_type
        self.external_id: Optional[str] = external_id
        msg = get_location_id_msg(location_name, location_type, external_id)
        msg = msg + " does not exist in inventory"
        super(LocationNotFoundException, self).__init__(msg)


class LocationIsNotUniqueException(CustomException):
    def __init__(
        self,
        location_name: Optional[str] = None,
        location_type: Optional[str] = None,
        external_id: Optional[str] = None,
    ) -> None:
        self.location_name: Optional[str] = location_name
        self.location_type: Optional[str] = location_type
        self.external_id: Optional[str] = external_id
        msg = get_location_id_msg(location_name, location_type, external_id)
        msg = msg + " has more than one result in inventory"
        super(LocationIsNotUniqueException, self).__init__(msg)


class EquipmentTypeNotFoundException(CustomException):
    def __init__(self, equipment_type_name: str) -> None:
        self.equipmentTypeName: str = equipment_type_name
        super(EquipmentTypeNotFoundException, self).__init__(
            f"Equipment type {equipment_type_name} does not exist in inventory"
        )


class EquipmentNotFoundException(CustomException):
    def __init__(
        self,
        equipment_name: Optional[str] = None,
        parent_equipment_name: Optional[str] = None,
        parent_position_name: Optional[str] = None,
    ) -> None:
        self.equipment_name: Optional[str] = equipment_name
        self.parent_equipment_name: Optional[str] = parent_equipment_name
        self.parent_position_name: Optional[str] = parent_position_name
        if equipment_name:
            super(EquipmentNotFoundException, self).__init__(
                f"equipment {equipment_name} does not exist in inventory"
            )
        else:
            super(EquipmentNotFoundException, self).__init__(
                f"Position {parent_equipment_name} in equipment "
                f"{parent_position_name} is not occupied"
            )


class EquipmentIsNotUniqueException(CustomException):
    def __init__(
        self,
        equipment_name: Optional[str] = None,
        parent_equipment_name: Optional[str] = None,
        parent_position_name: Optional[str] = None,
    ) -> None:
        self.equipment_name: Optional[str] = equipment_name
        self.parent_equipment_name: Optional[str] = parent_equipment_name
        self.parent_position_name: Optional[str] = parent_position_name
        if equipment_name:
            super(EquipmentIsNotUniqueException, self).__init__(
                f"equipment {equipment_name} has more than one result in inventory"
            )
        else:
            super(EquipmentIsNotUniqueException, self).__init__(
                f"More than one installed equipment in equipment "
                f"{parent_equipment_name} at position {parent_position_name}"
            )


class EquipmentPositionNotFoundException(CustomException):
    def __init__(self, parent_equipment_name: str, parent_position_name: str) -> None:
        self.parent_equipment_name: str = parent_equipment_name
        self.parent_position_name: str = parent_position_name
        super(EquipmentPositionNotFoundException, self).__init__(
            f"Equipment {parent_equipment_name} has no position "
            f"{parent_position_name}"
        )


class EquipmentPositionIsNotUniqueException(CustomException):
    def __init__(self, parent_equipment_name: str, parent_position_name: str) -> None:
        self.parent_equipment_name: str = parent_equipment_name
        self.parent_position_name: str = parent_position_name
        super(EquipmentPositionIsNotUniqueException, self).__init__(
            f"Equipment {parent_equipment_name} has more than one position "
            f"{parent_position_name}"
        )


class LinkNotFoundException(CustomException):
    def __init__(self, equipment_name: str, port_name: str) -> None:
        self.equipment_name: str = equipment_name
        self.port_name: str = port_name
        super(LinkNotFoundException, self).__init__(
            f"No link in port {port_name} in equipment {equipment_name}"
        )


class PortAlreadyOccupiedException(CustomException):
    def __init__(self, equipment_name: str, port_name: str) -> None:
        self.equipment_name: str = equipment_name
        self.port_name: str = port_name
        super(PortAlreadyOccupiedException, self).__init__(
            f"Port {port_name} in equipment {equipment_name} is already occupied"
        )


class EquipmentPortNotFoundException(CustomException):
    def __init__(self, equipment_name: str, port_name: str) -> None:
        self.equipment_name: str = equipment_name
        self.port_name: str = port_name
        super(EquipmentPortNotFoundException, self).__init__(
            f"Equipment {equipment_name} has no port {port_name}"
        )


class EquipmentPortIsNotUniqueException(CustomException):
    def __init__(self, equipment_name: str, port_name: str) -> None:
        self.equipment_name: str = equipment_name
        self.port_name: str = port_name
        super(EquipmentPortIsNotUniqueException, self).__init__(
            f"Equipment {equipment_name} has more than one port {port_name}"
        )


class LocationCannotBeDeletedWithDependency(CustomException):
    def __init__(self, location_name: str, dependency: str) -> None:
        self.location_name: str = location_name
        self.dependency: str = dependency
        super(LocationCannotBeDeletedWithDependency, self).__init__(
            f"Location {location_name} cannot be deleted with {dependency}"
        )
