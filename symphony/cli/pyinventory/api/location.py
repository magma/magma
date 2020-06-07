#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from typing import Dict, List, Optional, Sequence, Tuple

from pysymphony import SymphonyClient

from .._utils import deprecated, get_graphql_property_inputs
from ..common.cache import LOCATION_TYPES
from ..common.constant import LOCATION_PAGINATION_STEP, LOCATIONS_TO_SEARCH
from ..common.data_class import Document, ImageEntity, Location, PropertyValue
from ..common.data_enum import Entity
from ..exceptions import (
    EntityNotFoundError,
    LocationCannotBeDeletedWithDependency,
    LocationIsNotUniqueException,
)
from ..graphql.enum.filter_operator import FilterOperator
from ..graphql.fragment.location import LocationFragment
from ..graphql.input.add_location import AddLocationInput
from ..graphql.input.edit_location import EditLocationInput
from ..graphql.input.location_filter import LocationFilterInput, LocationFilterType
from ..graphql.mutation.add_location import AddLocationMutation
from ..graphql.mutation.edit_location import EditLocationMutation
from ..graphql.mutation.move_location import MoveLocationMutation
from ..graphql.mutation.remove_location import RemoveLocationMutation
from ..graphql.query.get_locations import GetLocationsQuery
from ..graphql.query.location_children import LocationChildrenQuery
from ..graphql.query.location_deps import LocationDepsQuery
from ..graphql.query.location_documents import LocationDocumentsQuery
from ..graphql.query.location_search import LocationSearchQuery


def _get_locations_by_name_and_type(
    client: SymphonyClient, location_name: str, location_type_name: str
) -> Sequence[LocationFragment]:
    location_filters = [
        LocationFilterInput(
            filterType=LocationFilterType.LOCATION_TYPE,
            operator=FilterOperator.IS_ONE_OF,
            idSet=[LOCATION_TYPES[location_type_name].id],
            stringSet=[],
        ),
        LocationFilterInput(
            filterType=LocationFilterType.LOCATION_INST_NAME,
            operator=FilterOperator.IS,
            stringValue=location_name,
            idSet=[],
            stringSet=[],
        ),
    ]

    result = LocationSearchQuery.execute(
        client, filters=location_filters, limit=LOCATIONS_TO_SEARCH
    )

    return result.locations


def _get_location_children_by_name_and_type(
    client: SymphonyClient,
    location_name: str,
    location_type_name: str,
    current_location: Optional[LocationFragment],
) -> Sequence[LocationFragment]:
    if current_location is None:
        return _get_locations_by_name_and_type(
            client=client,
            location_name=location_name,
            location_type_name=location_type_name,
        )

    else:
        location_id = current_location.id
        location_with_children = LocationChildrenQuery.execute(client, id=location_id)
        if location_with_children is None:
            raise EntityNotFoundError(entity=Entity.Location, entity_id=location_id)

        return [
            location
            for location in location_with_children.children
            if location.locationType.name == location_type_name
            and location.name == location_name
        ]


def add_location(
    client: SymphonyClient,
    location_hirerchy: List[Tuple[str, str]],
    properties_dict: Dict[str, PropertyValue],
    lat: Optional[float] = None,
    long: Optional[float] = None,
    externalID: Optional[str] = None,
) -> Location:
    """Create a new location of a specific type with a specific name.
        It will also get the requested location specifiers for hirerchy
        leading to it and will create all the hirerchy.
        However the `lat`,`long` and `properties_dict` would only apply for the last location in the chain.
        If a location with its name in this place already exists, then existing location is returned

        Args:
            location_hirerchy (List[Tuple[str, str]]): hirerchy of locations.
            - str - location type name
            - str - location name

            properties_dict (Dict[str, PropertyValue]): dict of property name to property value. The property value should match
                            the property type, otherwise exception is raised
            - str - property name
            - PropertyValue - new value of the same type for this property

            lat (float): latitude
            long (float): longitude
            externalID (str): location external ID

        Returns:
            `pyinventory.common.data_class.Location` object

        Raises:
            LocationIsNotUniqueException: if there is two possible locations
                inside the chain and it is not clear where to create or what to return
            FailedOperationException: for internal inventory error
            `pyinventory.exceptions.EntityNotFoundError`: parent location in the chain does not exist

        Example:
            ```
            location = client.add_location(
                location_hirerchy=[
                    ("Country", "England"),
                    ("City", "Milton Keynes"),
                    ("Site", "Bletchley Park")
                ],
                properties_dict={
                    "Date Property": date.today(),
                    "Lat/Lng Property": (-1.23,9.232),
                    "E-mail Property": "user@fb.com",
                    "Number Property": 11,
                    "String Property": "aa",
                    "Float Property": 1.23
                },
                lat=-11.32,
                long=98.32,
                externalID=None)
            ```
    """

    last_location: Optional[LocationFragment] = None

    for i, location in enumerate(location_hirerchy):

        location_type = location[0]
        location_name = location[1]
        properties = []
        lat_val = None
        long_val = None

        if i == len(location_hirerchy) - 1:
            property_types = LOCATION_TYPES[location_type].property_types
            properties = get_graphql_property_inputs(property_types, properties_dict)
            lat_val = lat
            long_val = long

        location_search_result = _get_location_children_by_name_and_type(
            client=client,
            location_name=location_name,
            location_type_name=location_type,
            current_location=last_location,
        )
        if not location_search_result:
            last_location = AddLocationMutation.execute(
                client=client,
                input=AddLocationInput(
                    name=location_name,
                    type=LOCATION_TYPES[location_type].id,
                    latitude=lat_val,
                    longitude=long_val,
                    parent=last_location.id if last_location else None,
                    properties=properties,
                    externalID=externalID,
                ),
            )

        elif len(location_search_result) == 1:
            last_location = location_search_result[0]

        elif len(location_search_result) > 1:
            raise LocationIsNotUniqueException(
                location_name=location_name, location_type=location_type
            )
    if last_location is None:
        raise EntityNotFoundError(
            entity=Entity.Location, msg=f"<location_hierarchy: {location_hirerchy}>"
        )

    return Location(
        id=last_location.id,
        name=last_location.name,
        latitude=last_location.latitude,
        longitude=last_location.longitude,
        externalId=last_location.externalId,
        locationTypeName=last_location.locationType.name,
        properties=last_location.properties,
    )


def get_location(
    client: SymphonyClient, location_hirerchy: List[Tuple[str, str]]
) -> Location:
    """This function returns a location of a specific type with a specific name.
        It can get only the requested location specifiers or the hirerchy leading to it

        Args:
            location_hirerchy (List[Tuple[str, str]]): hirerchy of locations
            - str - location type name
            - str - location name

        Returns:
            `pyinventory.common.data_class.Location` object

        Raises:
            LocationIsNotUniqueException: if there is more than one correct
                location to return
            LocationNotFoundException: if no location was found
            `pyinventory.exceptions.EntityNotFoundError`: location in the chain does not exist
            FailedOperationException: for internal inventory error

        Example:
            ```
            location = client.get_location(
                location_hirerchy=[
                    ("Country", "England"),
                    ("City", "Milton Keynes"),
                    ("Site", "Bletchley Park")
                ])
            ```
            or
            ```
            # this call will fail if there is Bletchley Park in two cities
            location = client.get_location(location_hirerchy=[("Site", "Bletchley Park")])
            ```
    """

    last_location = None

    for location in location_hirerchy:
        location_type = location[0]
        location_name = location[1]

        location_search_result = _get_location_children_by_name_and_type(
            client=client,
            location_name=location_name,
            location_type_name=location_type,
            current_location=last_location,
        )

        if not location_search_result:
            raise EntityNotFoundError(entity=Entity.Location, entity_name=location_name)

        if len(location_search_result) > 1:
            raise LocationIsNotUniqueException(
                location_name=location_name, location_type=location_type
            )

        last_location = location_search_result[0]

    if last_location is None:
        raise EntityNotFoundError(
            entity=Entity.Location, msg=f"<location_hierarchy: {location_hirerchy}>"
        )
    return Location(
        id=last_location.id,
        name=last_location.name,
        latitude=last_location.latitude,
        longitude=last_location.longitude,
        externalId=last_location.externalId,
        locationTypeName=last_location.locationType.name,
        properties=last_location.properties,
    )


def get_locations(client: SymphonyClient) -> List[Location]:
    """This function returns all existing locations

        Returns:
            List[ `pyinventory.common.data_class.Location` ]

        Example:
            ```
            all_locations = client.get_locations()
            ```
    """
    locations = GetLocationsQuery.execute(client, first=LOCATION_PAGINATION_STEP)
    edges = locations.edges if locations else []
    while locations is not None and locations.pageInfo.hasNextPage:
        locations = GetLocationsQuery.execute(
            client, after=locations.pageInfo.endCursor, first=LOCATION_PAGINATION_STEP
        )
        if locations is not None:
            edges.extend(locations.edges)

    result = []
    for edge in edges:
        node = edge.node
        if node is not None:
            result.append(
                Location(
                    name=node.name,
                    id=node.id,
                    latitude=node.latitude,
                    longitude=node.longitude,
                    externalId=node.externalId,
                    locationTypeName=node.locationType.name,
                    properties=node.properties,
                )
            )

    return result


def get_location_children(client: SymphonyClient, location_id: str) -> List[Location]:
    """This function returns all children locations of the given location

        Args:
            location_id (str): parent location ID

        Returns:
            List[ `pyinventory.common.data_class.Location` ]

        Raises:
            `pyinventory.exceptions.EntityNotFoundError`: location does not exist

        Example:
            ```
            client.add_location([("Country", "England"), ("City", "Milton Keynes")], {})
            client.add_location([("Country", "England"), ("City", "London")], {})
            parent_location = client.get_location(location_hirerchy=[("Country", "England")])
            children_locations = client.get_location_children(location_id=parent_location.id)
            # This call will return a list with 2 locations: "Milton Keynes" and "London"
            ```
    """
    location_with_children = LocationChildrenQuery.execute(client, id=location_id)
    if not location_with_children:
        raise EntityNotFoundError(entity=Entity.Location, entity_id=location_id)

    return [
        Location(
            name=location.name,
            id=location.id,
            latitude=location.latitude,
            longitude=location.longitude,
            externalId=location.externalId,
            locationTypeName=location.locationType.name,
            properties=location.properties,
        )
        for location in location_with_children.children
    ]


def edit_location(
    client: SymphonyClient,
    location: Location,
    new_name: Optional[str] = None,
    new_lat: Optional[float] = None,
    new_long: Optional[float] = None,
    new_external_id: Optional[str] = None,
    new_properties: Optional[Dict[str, PropertyValue]] = None,
) -> Location:
    """This function returns edited location.

        Args:
            location ( `pyinventory.common.data_class.Location` ): location object
            new_name (Optional[str]): location new name
            new_lat (Optional[float]): location new latitude
            new_long (Optional[float]): location new longitude
            new_external_id (Optional[float]): location new external ID
            new_properties (Optional[Dict[str, PropertyValue]]): dictionary of property name to property value
            - str - property name
            - PropertyValue - new value of the same type for this property

        Returns:
            `pyinventory.common.data_class.Location` object

        Raises:
            FailedOperationException: for internal inventory error

        Example:
            ```
            # this call will fail if there is Bletchley Park in two cities
            location = client.get_location(location_hirerchy=[("Site", "Bletchley Park")])
            edited_location = client.edit_location(
                location=location,
                new_name="New Bletchley Park",
                new_lat=10,
                new_long=20,
                new_external_id=None,
                new_properties={"Contact": "new_contact@info.com"},
            )
            ```
    """
    properties = []
    location_type = location.locationTypeName
    property_types = LOCATION_TYPES[location_type].property_types
    if new_properties:
        properties = get_graphql_property_inputs(property_types, new_properties)
    if new_external_id is None:
        new_external_id = location.externalId
    edit_location_input = EditLocationInput(
        id=location.id,
        name=new_name if new_name is not None else location.name,
        latitude=new_lat if new_lat is not None else location.latitude,
        longitude=new_long if new_long is not None else location.longitude,
        properties=properties,
        externalID=new_external_id,
    )
    result = EditLocationMutation.execute(client, edit_location_input)
    return Location(
        name=result.name,
        id=result.id,
        latitude=result.latitude,
        longitude=result.longitude,
        externalId=result.externalId,
        locationTypeName=result.locationType.name,
        properties=result.properties,
    )


def delete_location(client: SymphonyClient, location: Location) -> None:
    """This delete existing location.

        Args:
            location ( `pyinventory.common.data_class.Location` ): location object

        Raises:
            `pyinventory.exceptions.EntityNotFoundError`: if location does not exist
            FailedOperationException: for internal inventory error
            LocationCannotBeDeletedWithDependency: if there are dependencies in this location
            ["files", "images", "children", "surveys", "equipment"]

        Example:
            ```
            location = client.get_location(location_hirerchy=[('Site', 'Bletchley Park')])
            client.delete_location(location=location)
            ```
    """
    location_with_deps = LocationDepsQuery.execute(client, id=location.id)
    if location_with_deps is None:
        raise EntityNotFoundError(entity=Entity.Location, entity_id=location.id)
    if len(location_with_deps.files) > 0:
        raise LocationCannotBeDeletedWithDependency(location.name, "files")
    if len(location_with_deps.images) > 0:
        raise LocationCannotBeDeletedWithDependency(location.name, "images")
    if len(location_with_deps.children) > 0:
        raise LocationCannotBeDeletedWithDependency(location.name, "children")
    if len(location_with_deps.surveys) > 0:
        raise LocationCannotBeDeletedWithDependency(location.name, "surveys")
    if len(location_with_deps.equipments) > 0:
        raise LocationCannotBeDeletedWithDependency(location.name, "equipment")
    RemoveLocationMutation.execute(client, id=location.id)


def move_location(
    client: SymphonyClient, location_id: str, new_parent_id: Optional[str]
) -> Location:
    """This function moves existing location to another existing parent location.

        Args:
            location_id (str): existing location ID to be moved
            new_parent_id (Optional[str]): new existing parent location ID

        Returns:
            `pyinventory.common.data_class.Location` object

        Raises:
            FailedOperationException: for internal inventory error

        Example:
            ```
            # this call will fail if there is Bletchley Park in two cities
            location = client.get_location(location_hirerchy=[("Site", "Bletchley Park")])
            moved_location = client.move_locatoin(
                location_id=location.id,
                new_parent_id="12345"
            )
            ```
    """
    result = MoveLocationMutation.execute(
        client, locationID=location_id, parentLocationID=new_parent_id
    )
    return Location(
        name=result.name,
        id=result.id,
        latitude=result.latitude,
        longitude=result.longitude,
        externalId=result.externalId,
        locationTypeName=result.locationType.name,
        properties=result.properties,
    )


@deprecated(deprecated_in="2.4.0", deprecated_by="get_location_by_external_id")
def get_locations_by_external_id(
    client: SymphonyClient, external_id: str
) -> List[Location]:

    locations = []
    locations.append(get_location_by_external_id(client, external_id))
    return locations


def get_location_by_external_id(client: SymphonyClient, external_id: str) -> Location:
    """This function returns location by external ID.

        Args:
            external_id (str): location external ID

        Returns:
            `pyinventory.common.data_class.Location` object

        Raises:
            LocationNotFoundException: location with this external ID does not exists
            `pyinventory.exceptions.EntityNotFoundError`: location does not found
            FailedOperationException: for internal inventory error

        Example:
            ```
            location = client.get_location_by_external_id(external_id="12345")
            ```
    """
    location_filter = LocationFilterInput(
        filterType=LocationFilterType.LOCATION_INST_EXTERNAL_ID,
        operator=FilterOperator.IS,
        stringValue=external_id,
        idSet=[],
        stringSet=[],
    )

    location_search_result = LocationSearchQuery.execute(
        client, filters=[location_filter], limit=LOCATIONS_TO_SEARCH
    )

    if not location_search_result or location_search_result.count == 0:
        raise EntityNotFoundError(
            entity=Entity.Location, msg=f"<external_id: {external_id}"
        )
    if location_search_result.count > 1:
        raise LocationIsNotUniqueException(external_id=external_id)

    location_details = location_search_result.locations[0]

    return Location(
        name=location_details.name,
        id=location_details.id,
        latitude=location_details.latitude,
        longitude=location_details.longitude,
        externalId=location_details.externalId,
        locationTypeName=location_details.locationType.name,
        properties=location_details.properties,
    )


def get_location_documents(
    client: SymphonyClient, location: Location
) -> List[Document]:
    """This function returns locations documents.

        Args:
            location ( `pyinventory.common.data_class.Location` ): location object

        Returns:
            List[ `pyinventory.common.data_class.Document` ]

        Raises:
            `pyinventory.exceptions.EntityNotFoundError`: location does not exists
            FailedOperationException: for internal inventory error

        Example:
            ```
            # this call will fail if there is Bletchley Park in two cities
            location = client.get_location(location_hirerchy=[("Site", "Bletchley Park")])
            location = client.get_location_documents(location=location)
            ```
    """
    location_with_documents = LocationDocumentsQuery.execute(client, id=location.id)
    if not location_with_documents:
        raise EntityNotFoundError(entity=Entity.Location, entity_id=location.id)
    files = [
        Document(
            name=file.fileName,
            id=file.id,
            parent_id=location.id,
            parent_entity=ImageEntity.LOCATION,
            category=file.category,
        )
        for file in location_with_documents.files
    ]
    images = [
        Document(
            name=file.fileName,
            id=file.id,
            parent_id=location.id,
            parent_entity=ImageEntity.LOCATION,
            category=file.category,
        )
        for file in location_with_documents.images
    ]

    return files + images
