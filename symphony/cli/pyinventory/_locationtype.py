#!/usr/bin/env python3
# pyre-strict

from typing import List, Tuple, Union

from dacite import from_dict
from gql.gql.client import OperationException

from ._utils import PropertyValue, _make_property_types
from .api.location import delete_location
from .consts import Location, LocationType
from .graphql.add_location_type_mutation import (
    AddLocationTypeInput,
    AddLocationTypeMutation,
    PropertyKind,
)
from .graphql.location_type_locations_query import LocationTypeLocationsQuery
from .graphql.location_types_query import LocationTypesQuery
from .graphql.remove_location_type_mutation import RemoveLocationTypeMutation
from .graphql_client import GraphqlClient
from .reporter import FailedOperationException


ADD_LOCATION_TYPE_MUTATION_NAME = "addLocationType"


def _populate_location_types(client: GraphqlClient) -> None:
    edges = LocationTypesQuery.execute(client).locationTypes.edges
    for edge in edges:
        node = edge.node
        client.locationTypes[node.name] = LocationType(
            name=node.name,
            id=node.id,
            propertyTypes=[p.to_dict() for p in node.propertyTypes],
        )


def add_location_type(
    client: GraphqlClient,
    name: str,
    properties: List[Tuple[str, str, PropertyValue, bool]],
    map_zoom_level: int = 8,
) -> LocationType:
    property_types = _make_property_types(properties)

    def property_type_to_kind(
        key: str, value: Union[str, int, float, bool]
    ) -> Union[str, int, float, bool, PropertyKind]:
        return value if key != "type" else PropertyKind(value)

    new_property_types = [
        {k: property_type_to_kind(k, v) for k, v in property_type.items()}
        for property_type in property_types
    ]
    add_location_type_variables = {
        "name": name,
        "mapZoomLevel": map_zoom_level,
        "properties": new_property_types,
    }
    try:
        result = AddLocationTypeMutation.execute(
            client,
            AddLocationTypeInput(
                name=name,
                mapZoomLevel=map_zoom_level,
                properties=[
                    from_dict(data_class=AddLocationTypeInput.PropertyTypeInput, data=p)
                    for p in new_property_types
                ],
                surveyTemplateCategories=[],
            ),
        ).__dict__[ADD_LOCATION_TYPE_MUTATION_NAME]
        client.reporter.log_successful_operation(
            ADD_LOCATION_TYPE_MUTATION_NAME, add_location_type_variables
        )
    except OperationException as e:
        raise FailedOperationException(
            client.reporter,
            e.err_msg,
            e.err_id,
            ADD_LOCATION_TYPE_MUTATION_NAME,
            add_location_type_variables,
        )

    location_type = LocationType(
        name=result.name,
        id=result.id,
        propertyTypes=[p.to_dict() for p in result.propertyTypes],
    )
    client.locationTypes[result.name] = location_type
    return location_type


def delete_locations_by_location_type(
    client: GraphqlClient, location_type: LocationType
) -> None:
    locations = LocationTypeLocationsQuery.execute(
        client, id=location_type.id
    ).locationType.locations.edges
    for location in locations:
        delete_location(
            client,
            Location(
                id=location.node.id,
                name=location.node.name,
                latitude=location.node.latitude,
                longitude=location.node.longitude,
                externalId=location.node.externalId,
                locationTypeName=location.node.locationType.name,
            ),
        )


def delete_location_type_with_locations(
    client: GraphqlClient, location_type: LocationType
) -> None:
    delete_locations_by_location_type(client, location_type)
    RemoveLocationTypeMutation.execute(client, id=location_type.id)
