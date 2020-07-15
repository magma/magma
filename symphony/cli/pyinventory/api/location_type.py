#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from typing import List

from pysymphony import SymphonyClient

from ..common.cache import LOCATION_TYPES
from ..common.data_class import Location, LocationType, PropertyDefinition
from ..common.data_enum import Entity
from ..common.data_format import (
    format_to_property_definitions,
    format_to_property_type_inputs,
)
from ..exceptions import EntityNotFoundError
from ..graphql.input.add_location_type import AddLocationTypeInput
from ..graphql.mutation.add_location_type import AddLocationTypeMutation
from ..graphql.mutation.remove_location_type import RemoveLocationTypeMutation
from ..graphql.query.location_type_locations import LocationTypeLocationsQuery
from ..graphql.query.location_types import LocationTypesQuery
from .location import delete_location


def _populate_location_types(client: SymphonyClient) -> None:
    location_types = LocationTypesQuery.execute(client)
    if not location_types:
        return
    edges = location_types.edges
    for edge in edges:
        node = edge.node
        if node:
            LOCATION_TYPES[node.name] = LocationType(
                name=node.name,
                id=node.id,
                property_types=format_to_property_definitions(node.propertyTypes),
            )


def add_location_type(
    client: SymphonyClient,
    name: str,
    properties: List[PropertyDefinition],
    map_zoom_level: int = 8,
) -> LocationType:
    """This function creates new location type.

        :param name: Location type name
        :type name: str
        :param properties: List of property definitions
        :type properties: List[ :class:`~pyinventory.common.data_class.PropertyDefinition` ]
        :param map_zoom_level: Map zoom level
        :type map_zoom_level: int

        :raises:
            FailedOperationException: Internal inventory error

        :return: LocationType object
        :rtype: :class:`~pyinventory.common.data_class.LocationType`

        **Example**

        .. code-block:: python

            location_type = client.add_location_type(
                name="city",
                properties=[
                    PropertyDefinition(
                        property_name="Contact",
                        property_kind=PropertyKind.string,
                        default_raw_value=None,
                        is_fixed=True
                    )
                ],
                map_zoom_level=5,
            )
    """
    new_property_types = format_to_property_type_inputs(data=properties)
    result = AddLocationTypeMutation.execute(
        client,
        AddLocationTypeInput(
            name=name,
            mapZoomLevel=map_zoom_level,
            properties=new_property_types,
            surveyTemplateCategories=[],
        ),
    )

    location_type = LocationType(
        name=result.name,
        id=result.id,
        property_types=format_to_property_definitions(result.propertyTypes),
    )
    LOCATION_TYPES[result.name] = location_type
    return location_type


def delete_locations_by_location_type(
    client: SymphonyClient, location_type: LocationType
) -> None:
    """Delete locatons by location type.

        :param location_type: LocationType object
        :type location_type: :class:`~pyinventory.common.data_class.LocationType`

        :raises:
            `pyinventory.exceptions.EntityNotFoundError`: `location_type` does not exist

        :rtype: None

        **Example**

        .. code-block:: python

            client.delete_locations_by_location_type(location_type=location_type)
    """
    location_type_with_locations = LocationTypeLocationsQuery.execute(
        client, id=location_type.id
    )
    if location_type_with_locations is None:
        raise EntityNotFoundError(
            entity=Entity.LocationType, entity_id=location_type.id
        )
    locations = location_type_with_locations.locations
    if locations is None:
        return
    for location in locations.edges:
        node = location.node
        if node:
            delete_location(
                client,
                Location(
                    id=node.id,
                    name=node.name,
                    latitude=node.latitude,
                    longitude=node.longitude,
                    external_id=node.externalId,
                    location_type_name=node.locationType.name,
                    properties=node.properties,
                ),
            )


def delete_location_type_with_locations(
    client: SymphonyClient, location_type: LocationType
) -> None:
    """Delete locaton type with existing locations.

        :param location_type: LocationType object
        :type location_type: :class:`~pyinventory.common.data_class.LocationType`

        :raises:
            `pyinventory.exceptions.EntityNotFoundError`: `location_type` does not exist

        :rtype: None

        **Example**

        .. code-block:: python

            client.delete_location_type_with_locations(location_type=location_type)
    """
    delete_locations_by_location_type(client, location_type)
    RemoveLocationTypeMutation.execute(client, id=location_type.id)
