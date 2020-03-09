#!/usr/bin/env python3

from typing import Any, Dict, List, Optional, Tuple

from dacite import Config, from_dict
from gql.gql.client import OperationException
from gql.gql.reporter import FailedOperationException

from .._utils import deprecated, get_graphql_property_inputs
from ..client import SymphonyClient
from ..consts import Document, Entity, ImageEntity, Location
from ..exceptions import (
    EntityNotFoundError,
    LocationCannotBeDeletedWithDependency,
    LocationIsNotUniqueException,
    LocationNotFoundException,
)
from ..graphql.add_location_input import AddLocationInput
from ..graphql.add_location_mutation import AddLocationMutation
from ..graphql.edit_location_input import EditLocationInput
from ..graphql.edit_location_mutation import EditLocationMutation
from ..graphql.location_children_query import LocationChildrenQuery
from ..graphql.location_deps_query import LocationDepsQuery
from ..graphql.location_details_query import LocationDetailsQuery
from ..graphql.location_documents_query import LocationDocumentsQuery
from ..graphql.move_location_mutation import MoveLocationMutation
from ..graphql.property_input import PropertyInput
from ..graphql.remove_location_mutation import RemoveLocationMutation
from ..graphql.search_query import SearchQuery


ADD_LOCATION_MUTATION_NAME = "addLocation"
EDIT_LOCATION_MUTATION_NAME = "editLocation"
MOVE_LOCATION_MUTATION_NAME = "moveLocation"


def add_location(
    client: SymphonyClient,
    location_hirerchy: List[Tuple[str, str]],
    properties_dict: Dict[str, Any],
    lat: Optional[float] = None,
    long: Optional[float] = None,
    externalID: Optional[str] = None,
) -> Location:
    """Create a new location of a specific type with a specific name.
        It will also get the requested location specifiers for hirerchy leading to it and will create all
        the hirerchy.
        However the lat,long and propertiesDict would only apply for the last location in the chain.
        If a location with his name in this place already exists the existing location is returned

        Args:
            location_hirerchy (List[Tuple[str, str]]): An hirerchy of locations.
                The first str is location type name. The second str is location name
            properties_dict: dict of property name to property value. the property value should match
                            the property type. Otherwise exception is raised
            lat (float): latitude
            long (float): longitude
            externalID (str): ID from external system

        Returns:
            pyinventory.consts.Location object

        Raises:
            LocationIsNotUniqueException: if there is two possible locations
                inside the chain and it is not clear where to create or what to return
            FailedOperationException: for internal inventory error
            `pyinventory.exceptions.EntityNotFoundError`: parent location in the chain does not exist

        Example:
        ```
        location = client.add_location(
            [
                ('Country', 'England'),
                ('City', 'Milton Keynes'),
                ('Site', 'Bletchley Park')
            ],
            {
                'Date Property ': date.today(),
                'Lat/Lng Property: ': (-1.23,9.232),
                'E-mail Property ': "user@fb.com",
                'Number Property ': 11,
                'String Property ': "aa",
                'Float Property': 1.23
            },
            -11.32,
            98.32,
            None)
        ```
    """

    last_location = None

    for i, location in enumerate(location_hirerchy):
        location_type = location[0]
        location_name = location[1]

        properties = []
        lat_val = None
        long_val = None
        if i == len(location_hirerchy) - 1:
            property_types = client.locationTypes[location_type].propertyTypes
            properties = get_graphql_property_inputs(property_types, properties_dict)
            lat_val = lat
            long_val = long

        if last_location is None:
            locations = SearchQuery.execute(
                client, name=location_name
            ).searchForEntity.edges

            locations = [
                location.node
                for location in locations
                # pyre-fixme[16]: `Optional` has no attribute `entityType`.
                if location.node.entityType == "location"
                # pyre-fixme[16]: `Optional` has no attribute `type`.
                and location.node.type == location_type
                # pyre-fixme[16]: `Optional` has no attribute `name`.
                and location.node.name == location_name
            ]
            if len(locations) > 1:
                raise LocationIsNotUniqueException(
                    location_name=location_name, location_type=location_type
                )
            if len(locations) == 1:
                location_details = LocationDetailsQuery.execute(
                    client,
                    # pyre-fixme[16]: `Optional` has no attribute `entityId`.
                    id=locations[0].entityId,
                ).location
                if location_details is None:
                    raise EntityNotFoundError(
                        entity=Entity.Location, entity_id=locations[0].entityId
                    )
                last_location = Location(
                    name=location_details.name,
                    id=location_details.id,
                    latitude=location_details.latitude,
                    longitude=location_details.longitude,
                    externalId=location_details.externalId,
                    locationTypeName=location_details.locationType.name,
                )
            else:
                add_location_input = AddLocationInput(
                    name=location_name,
                    type=client.locationTypes[location_type].id,
                    latitude=lat,
                    longitude=long,
                    properties=properties,
                    externalID=externalID,
                )

                try:
                    result = AddLocationMutation.execute(
                        client, add_location_input
                    ).__dict__[ADD_LOCATION_MUTATION_NAME]
                    client.reporter.log_successful_operation(
                        ADD_LOCATION_MUTATION_NAME, add_location_input.__dict__
                    )
                except OperationException as e:
                    raise FailedOperationException(
                        client.reporter,
                        e.err_msg,
                        e.err_id,
                        ADD_LOCATION_MUTATION_NAME,
                        add_location_input.__dict__,
                    )
                last_location = Location(
                    name=result.name,
                    id=result.id,
                    latitude=result.latitude,
                    longitude=result.longitude,
                    externalId=result.externalId,
                    locationTypeName=result.locationType.name,
                )
        else:
            location_id = last_location.id
            location_with_children = LocationChildrenQuery.execute(
                client, id=location_id
            ).location
            if location_with_children is None:
                raise EntityNotFoundError(entity=Entity.Location, entity_id=location_id)

            locations = [
                location
                for location in location_with_children.children
                if location.locationType.name == location_type
                and location.name == location_name
            ]
            if len(locations) > 1:
                raise LocationIsNotUniqueException(
                    location_name=location_name, location_type=location_type
                )
            if len(locations) == 1:
                last_location = Location(
                    name=locations[0].name,
                    id=locations[0].id,
                    latitude=locations[0].latitude,
                    longitude=locations[0].longitude,
                    externalId=locations[0].externalId,
                    locationTypeName=locations[0].locationType.name,
                )
            else:
                add_location_input = AddLocationInput(
                    name=location_name,
                    type=client.locationTypes[location_type].id,
                    latitude=lat_val,
                    longitude=long_val,
                    parent=location_id,
                    properties=properties,
                    externalID=externalID,
                )
                try:
                    result = AddLocationMutation.execute(
                        client, add_location_input
                    ).__dict__[ADD_LOCATION_MUTATION_NAME]
                    client.reporter.log_successful_operation(
                        ADD_LOCATION_MUTATION_NAME, add_location_input.__dict__
                    )
                except OperationException as e:
                    raise FailedOperationException(
                        client.reporter,
                        e.err_msg,
                        e.err_id,
                        ADD_LOCATION_MUTATION_NAME,
                        add_location_input.__dict__,
                    )
                last_location = Location(
                    name=result.name,
                    id=result.id,
                    latitude=result.latitude,
                    longitude=result.longitude,
                    externalId=result.externalId,
                    locationTypeName=result.locationType.name,
                )

    if last_location is None:
        raise LocationNotFoundException()
    return last_location


def get_location(
    client: SymphonyClient, location_hirerchy: List[Tuple[str, str]]
) -> Location:
    """This function returns a location of a specific type with a specific name.
        It can get only the requested location specifiers or the hirerchy leading to it

        Args:
            location_hirerchy (list of tuple(str, str)):
                the first str is location type name
                the second str is location name

        Returns: pyinventory.consts.Location object

        Raises: LocationIsNotUniqueException: if there is more than one correct
                location to return
                LocationNotFoundException: if no location was found
                `pyinventory.exceptions.EntityNotFoundError`: location in the chain does not exist

        Example:
        ```
        location = client.get_location([
            ('Country', 'England'),
            ('City', 'Milton Keynes'),
            ('Site', 'Bletchley Park')
        ])
        ```
        or
        ```
        # this call will fail if there is Bletchley Park in two cities in london
        location = client.get_location([('Site', 'Bletchley Park')])
        ```
    """

    last_location = None

    for location in location_hirerchy:
        location_type = location[0]
        location_name = location[1]

        if last_location is None:
            entities = SearchQuery.execute(
                client, name=location_name
            ).searchForEntity.edges
            nodes = [entity.node for entity in entities]

            locations = [
                node
                for node in nodes
                if node is not None
                and node.entityType == "location"
                and node.type == location_type
                and node.name == location_name
            ]
            if len(locations) == 0:
                raise LocationNotFoundException(
                    location_name=location_name, location_type=location_type
                )
            if len(locations) != 1:
                raise LocationIsNotUniqueException(
                    location_name=location_name, location_type=location_type
                )
            location_details = LocationDetailsQuery.execute(
                client, id=locations[0].entityId
            ).location
            if location_details is None:
                raise EntityNotFoundError(
                    entity=Entity.Location, entity_id=locations[0].entityId
                )
            last_location = Location(
                name=location_details.name,
                id=location_details.id,
                latitude=location_details.latitude,
                longitude=location_details.longitude,
                externalId=location_details.externalId,
                locationTypeName=location_details.locationType.name,
            )
        else:
            location_id = last_location.id

            location_with_children = LocationChildrenQuery.execute(
                client, id=location_id
            ).location
            if location_with_children is None:
                raise EntityNotFoundError(entity=Entity.Location, entity_id=location_id)
            locations = [
                location
                for location in location_with_children.children
                if location.locationType.name == location_type
                and location.name == location_name
            ]
            if len(locations) == 0:
                raise LocationNotFoundException(location_name=location_name)
            if len(locations) != 1:
                raise LocationIsNotUniqueException(
                    location_name=location_name, location_type=location_type
                )
            last_location = Location(
                name=locations[0].name,
                id=locations[0].id,
                latitude=locations[0].latitude,
                longitude=locations[0].longitude,
                externalId=locations[0].externalId,
                locationTypeName=locations[0].locationType.name,
            )

    if last_location is None:
        raise LocationNotFoundException()
    return last_location


def get_location_children(client: SymphonyClient, location_id: str) -> List[Location]:
    """This function returns all locations that are children of the given location

        Args:
            location_id (str):
                id of the parent location

        Returns: List of pyinventory.consts.Location objects

        Raises: `pyinventory.exceptions.EntityNotFoundError`: location does not exist

        Example:
        ```
        client.addLocation([('Country', 'England'), ('City', 'Milton Keynes')], {})
        client.addLocation([('Country', 'England'), ('City', 'London')], {})
        locations = client.get_location_children([('Country', 'England')])
        # This call will return a list with 2 locations: 'Milton Keynes' and 'London'
        ```
    """
    location_with_children = LocationChildrenQuery.execute(
        client, id=location_id
    ).location
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
    new_properties: Optional[Dict[str, Any]] = None,
) -> Location:

    properties = []
    location_type = location.locationTypeName
    property_types = client.locationTypes[location_type].propertyTypes
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

    try:
        result = EditLocationMutation.execute(client, edit_location_input).__dict__[
            EDIT_LOCATION_MUTATION_NAME
        ]
        client.reporter.log_successful_operation(
            EDIT_LOCATION_MUTATION_NAME, edit_location_input.__dict__
        )
        return Location(
            name=result.name,
            id=result.id,
            latitude=result.latitude,
            longitude=result.longitude,
            externalId=result.externalId,
            locationTypeName=result.locationType.name,
        )

    except OperationException as e:
        raise FailedOperationException(
            client.reporter,
            e.err_msg,
            e.err_id,
            EDIT_LOCATION_MUTATION_NAME,
            edit_location_input.__dict__,
        )
        return None


def delete_location(client: SymphonyClient, location: Location) -> None:
    location_with_deps = LocationDepsQuery.execute(client, id=location.id).location
    if location_with_deps is None:
        raise EntityNotFoundError(entity=Entity.Location, entity_id=location.id)
    if len(location_with_deps.files) > 0:
        raise LocationCannotBeDeletedWithDependency(location.name, "files")
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
    params = {"locationID": location_id, "parentLocationID": new_parent_id}
    try:
        result = MoveLocationMutation.execute(
            client, locationID=location_id, parentLocationID=new_parent_id
        ).__dict__[MOVE_LOCATION_MUTATION_NAME]
        client.reporter.log_successful_operation(MOVE_LOCATION_MUTATION_NAME, params)
        return Location(
            name=result.name,
            id=result.id,
            latitude=result.latitude,
            longitude=result.longitude,
            externalId=result.externalId,
            locationTypeName=result.locationType.name,
        )

    except OperationException as e:
        raise FailedOperationException(
            client.reporter, e.err_msg, e.err_id, MOVE_LOCATION_MUTATION_NAME, params
        )


@deprecated(deprecated_in="2.4.0", deprecated_by="get_location_by_external_id")
def get_locations_by_external_id(
    client: SymphonyClient, external_id: str
) -> List[Location]:

    locations = []
    locations.append(get_location_by_external_id(client, external_id))
    return locations


def get_location_by_external_id(client: SymphonyClient, external_id: str) -> Location:
    locations = SearchQuery.execute(client, name=external_id).searchForEntity.edges
    if not locations:
        raise LocationNotFoundException()

    location_details = None
    for location in locations:
        node = location.node
        if node is not None and node.entityType == "location":
            location_details = LocationDetailsQuery.execute(
                client, id=node.entityId
            ).location
            if location_details is None:
                raise EntityNotFoundError(
                    entity=Entity.Location, entity_id=node.entityId
                )
            if location_details.externalId == external_id:
                break
            else:
                location_details = None

    if not location_details:
        raise LocationNotFoundException()

    return Location(
        name=location_details.name,
        id=location_details.id,
        latitude=location_details.latitude,
        longitude=location_details.longitude,
        externalId=location_details.externalId,
        locationTypeName=location_details.locationType.name,
    )


def get_location_documents(
    client: SymphonyClient, location: Location
) -> List[Document]:
    location_with_documents = LocationDocumentsQuery.execute(
        client, id=location.id
    ).location
    if not location_with_documents:
        raise EntityNotFoundError(entity=Entity.Location, entity_id=location.id)
    return [
        Document(
            name=file.fileName,
            id=file.id,
            parentId=location.id,
            parentEntity=ImageEntity.LOCATION,
            category=file.category,
        )
        for file in location_with_documents.files
    ]
