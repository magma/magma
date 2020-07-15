#!/usr/bin/env python3
# Copyright (c) 2004-present Facebook All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from numbers import Number
from typing import Dict, Iterator, List, Optional, Sequence, Tuple, cast

from pysymphony import SymphonyClient

from .._utils import get_graphql_property_inputs
from ..common.cache import LOCATION_TYPES
from ..common.constant import LOCATIONS_TO_SEARCH, PAGINATION_STEP
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
    locations = []
    if result is not None:
        for edge in result.edges:
            node = edge.node
            if node is not None:
                locations.append(node)
    return locations


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
    external_id: Optional[str] = None,
) -> Location:
    """Create a new location of a specific type with a specific name.
        It will also get the requested location specifiers for hirerchy
        leading to it and will create all the hirerchy.
        However the `lat`,`long` and `properties_dict` would only apply for the last location in the chain.
        If a location with its name in this place already exists, then existing location is returned

        :param location_hirerchy: Locations hierarchy
        :type location_hirerchy: List[Tuple[str, str]]

            * str - location type name
            * str - location name

        :param properties_dict: Dictionary of property name to property value

            * str - property name
            * PropertyValue - new value of the same type for this property

        :type properties_dict: Dict[str, PropertyValue]
        :param lat: Latitude
        :type lat: float, optional
        :param long: Longitude
        :type long: float, optional
        :param external_id: Location external ID
        :type external_id: str, optional

        :raises:
            * LocationIsNotUniqueException: There is more than one location to return
              in the chain and it is not clear where to create or what to return
            * FailedOperationException: Internal inventory error
            * :class:`~pyinventory.exceptions.EntityNotFoundError`: Parent location in the chain does not exist

        :return: Location object
        :rtype: :class:`~pyinventory.common.data_class.Location`

        **Example**

        .. code-block:: python

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
                external_id=None,
            )
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
            lat_val = cast(Number, lat)
            long_val = cast(Number, long)

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
                    externalID=external_id,
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
        external_id=last_location.externalId,
        location_type_name=last_location.locationType.name,
        properties=last_location.properties,
    )


def get_location(
    client: SymphonyClient, location_hirerchy: List[Tuple[str, str]]
) -> Location:
    """This function returns a location of a specific type with a specific name.
        It can get only the requested location specifiers or the hirerchy leading to it

        :param location_hirerchy: Locations hierarchy
        :type location_hirerchy: List[Tuple[str, str]]

            * str - location type name
            * str - location name

        :raises:
            * LocationIsNotUniqueException: There is more than one location to return
              in the chain and it is not clear where to create or what to return
            * LocationNotFoundException: Location was not found
            * FailedOperationException: Internal inventory error
            * :class:`~pyinventory.exceptions.EntityNotFoundError`: Location in the chain does not exist

        :return: Location object
        :rtype: :class:`~pyinventory.common.data_class.Location`

        **Example**

        .. code-block:: python

            location = client.get_location(
                location_hirerchy=[
                    ("Country", "England"),
                    ("City", "Milton Keynes"),
                    ("Site", "Bletchley Park")
                ])

        .. code-block:: python

            # this call will fail if there is Bletchley Park in two cities
            location = client.get_location(location_hirerchy=[("Site", "Bletchley Park")])
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
        external_id=last_location.externalId,
        location_type_name=last_location.locationType.name,
        properties=last_location.properties,
    )


def get_locations(client: SymphonyClient) -> Iterator[Location]:
    """This function returns all existing locations

        :return: Locations Iterator
        :rtype: Iterator[ :class:`~pyinventory.common.data_class.Location` ]

        **Example**

        .. code-block:: python

            all_locations = client.get_locations()
    """

    def generate_pages(
        client: SymphonyClient,
    ) -> Iterator[GetLocationsQuery.GetLocationsQueryData.LocationConnection]:
        locations = GetLocationsQuery.execute(client, first=PAGINATION_STEP)
        if locations:
            yield locations
        while locations is not None and locations.pageInfo.hasNextPage:
            locations = GetLocationsQuery.execute(
                client, after=locations.pageInfo.endCursor, first=PAGINATION_STEP
            )
            if locations is not None:
                yield locations

    for page in generate_pages(client):
        for edge in page.edges:
            node = edge.node
            if node is not None:
                yield Location(
                    name=node.name,
                    id=node.id,
                    latitude=node.latitude,
                    longitude=node.longitude,
                    external_id=node.externalId,
                    location_type_name=node.locationType.name,
                    properties=node.properties,
                )


def get_location_children(
    client: SymphonyClient, location_id: str
) -> Iterator[Location]:
    """This function returns all children locations of the given location

        :param location_id: Parent location ID
        :type location_id: str

        :raises:
            :class:`~pyinventory.exceptions.EntityNotFoundError`: Location does not exist

        :return: Locations Iterator
        :rtype: Iterator[ :class:`~pyinventory.common.data_class.Location` ]

        **Example**

        .. code-block:: python

            client.add_location(
                [
                    ("Country", "England"),
                    ("City", "Milton Keynes"),
                ],
                {},
            )
            client.add_location(
                [
                    ("Country", "England"),
                    ("City", "London"),
                ],
                {},
            )
            parent_location = client.get_location(
                location_hirerchy=[
                    ("Country", "England"),
                ],
            )
            children_locations = client.get_location_children(location_id=parent_location.id)
            # This call will return a list with 2 locations: "Milton Keynes" and "London"
    """
    location_with_children = LocationChildrenQuery.execute(client, id=location_id)
    if not location_with_children:
        raise EntityNotFoundError(entity=Entity.Location, entity_id=location_id)

    for location in location_with_children.children:
        yield Location(
            name=location.name,
            id=location.id,
            latitude=location.latitude,
            longitude=location.longitude,
            external_id=location.externalId,
            location_type_name=location.locationType.name,
            properties=location.properties,
        )


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

        :param location: Location object
        :type location: :class:`~pyinventory.common.data_class.Location`
        :param new_name: Location new name
        :type new_name: str, optional
        :param new_lat: Location new latitude
        :type new_lat: float, optional
        :param new_long: Location new longitude
        :type new_long: float, optional
        :param new_external_id: Location new external ID
        :type new_external_id: str, optional
        :param new_properties: Dictionary of property name to property value

            * str - property name
            * PropertyValue - new value of the same type for this property

        :type new_properties: Dict[str, PropertyValue], optional

        :raises:
            FailedOperationException: Internal inventory error

        :return: Location object
        :rtype: :class:`~pyinventory.common.data_class.Location`

        **Example**

        .. code-block:: python

            # this call will fail incase the 'Bletchley Park' in two cities
            location = client.get_location(location_hirerchy=[("Site", "Bletchley Park")])
            edited_location = client.edit_location(
                location=location,
                new_name="New Bletchley Park",
                new_lat=10,
                new_long=20,
                new_external_id=None,
                new_properties={"Contact": "new_contact@info.com"},
            )
    """
    properties = []
    location_type = location.location_type_name
    property_types = LOCATION_TYPES[location_type].property_types
    if new_properties:
        properties = get_graphql_property_inputs(property_types, new_properties)
    if new_external_id is None:
        new_external_id = location.external_id
    edit_location_input = EditLocationInput(
        id=location.id,
        name=new_name if new_name is not None else location.name,
        latitude=cast(Number, new_lat) if new_lat is not None else location.latitude,
        longitude=cast(Number, new_long)
        if new_long is not None
        else location.longitude,
        properties=properties,
        externalID=new_external_id,
    )
    result = EditLocationMutation.execute(client, edit_location_input)
    return Location(
        name=result.name,
        id=result.id,
        latitude=result.latitude,
        longitude=result.longitude,
        external_id=result.externalId,
        location_type_name=result.locationType.name,
        properties=result.properties,
    )


def delete_location(client: SymphonyClient, location: Location) -> None:
    """This delete existing location.

        :param location: Location object
        :type location: :class:`~pyinventory.common.data_class.Location`

        :raises:
            * LocationCannotBeDeletedWithDependency: Location has dependencies in one or more
              ["files", "images", "children", "surveys", "equipment"]
            * FailedOperationException: Internal inventory error
            * :class:`~pyinventory.exceptions.EntityNotFoundError`: Location does not exist

        :rtype: None

        **Example**

        .. code-block:: python

            location = client.get_location(location_hirerchy=[('Site', 'Bletchley Park')])
            client.delete_location(location=location)
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

        :param location_id: Existing location ID to be moved
        :type location_id: str
        :param new_parent_id: New existing parent location ID
        :type new_parent_id: str, optional

        :raises:
            FailedOperationException: Internal inventory error

        :return: Location object
        :rtype: :class:`~pyinventory.common.data_class.Location`

        **Example**

        .. code-block:: python

            # this call will fail if there is Bletchley Park in two cities
            location = client.get_location(location_hirerchy=[("Site", "Bletchley Park")])
            moved_location = client.move_locatoin(
                location_id=location.id,
                new_parent_id="12345"
            )
    """
    result = MoveLocationMutation.execute(
        client, locationID=location_id, parentLocationID=new_parent_id
    )
    return Location(
        name=result.name,
        id=result.id,
        latitude=result.latitude,
        longitude=result.longitude,
        external_id=result.externalId,
        location_type_name=result.locationType.name,
        properties=result.properties,
    )


def get_location_by_external_id(client: SymphonyClient, external_id: str) -> Location:
    """This function returns location by external ID.

        :param external_id: Location external ID
        :type external_id: str

        :raises:
            * LocationNotFoundException: Location with this external ID does not exists
            * FailedOperationException: Internal inventory error
            * :class:`~pyinventory.exceptions.EntityNotFoundError`: Location does not exist

        :return: Location object
        :rtype: :class:`~pyinventory.common.data_class.Location`

        **Example**

        .. code-block:: python

            location = client.get_location_by_external_id(external_id="12345")
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

    if not location_search_result or location_search_result.totalCount == 0:
        raise EntityNotFoundError(
            entity=Entity.Location, msg=f"<external_id: {external_id}"
        )
    if location_search_result.totalCount > 1:
        raise LocationIsNotUniqueException(external_id=external_id)

    for edge in location_search_result.edges:
        node = edge.node
        if node is not None:
            return Location(
                name=node.name,
                id=node.id,
                latitude=node.latitude,
                longitude=node.longitude,
                external_id=node.externalId,
                location_type_name=node.locationType.name,
                properties=node.properties,
            )
    raise EntityNotFoundError(
        entity=Entity.Location, msg=f"<external_id: {external_id}"
    )


def get_location_documents(
    client: SymphonyClient, location: Location
) -> List[Document]:
    """This function returns locations documents.

        :param location: Location object
        :type location: :class:`~pyinventory.common.data_class.Location`

        :raises:
            * FailedOperationException: Internal inventory error
            * :class:`~pyinventory.exceptions.EntityNotFoundError`: Location does not exist

        :return: Documents List
        :rtype: List[ :class:`~pyinventory.common.data_class.Document` ]

        **Example**

        .. code-block:: python

            # this call will fail if there is Bletchley Park in two cities
            location = client.get_location(location_hirerchy=[("Site", "Bletchley Park")])
            location = client.get_location_documents(location=location)
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
