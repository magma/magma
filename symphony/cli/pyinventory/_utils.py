#!/usr/bin/env python3
# pyre-strict

from datetime import date, datetime
from typing import Any, Dict, List, Tuple, Union


PROPERTY_TYPE_TO_FIELD_NAME = {
    "date": "stringValue",
    "float": "floatValue",
    "int": "intValue",
    "email": "stringValue",
    "string": "stringValue",
    "bool": "booleanValue",
}

PropertyValue = Union[date, float, int, str, bool, Tuple[float, float]]


def _get_properties_to_add(
    property_types: List[Dict[str, Any]], properties_dict: Dict[str, PropertyValue]
) -> List[Dict[str, PropertyValue]]:
    properties = []
    for property_type in property_types:
        property_type_name = property_type["name"]
        property_type_id = property_type["id"]
        if property_type_name in properties_dict:
            type = property_type["type"]
            value = properties_dict[property_type_name]
            assert property_type[
                "isInstanceProperty"
            ], "property {} is not instance property".format(property_type_name)
            if type == "date":
                assert isinstance(
                    value, date
                ), "property {} is not of type datetime.date".format(property_type_name)
                properties.append(
                    {"propertyTypeID": property_type_id, "stringValue": str(value)}
                )
            elif type == "float":
                assert isinstance(
                    value, float
                ), "property {} is not of type float".format(property_type_name)
                properties.append(
                    {"propertyTypeID": property_type_id, "floatValue": value}
                )
            elif type == "int":
                assert isinstance(value, int), "property {} is not of type int".format(
                    property_type_name
                )
                properties.append(
                    {"propertyTypeID": property_type_id, "intValue": value}
                )
            elif type == "email":
                assert isinstance(
                    value, str
                ), "property {} is not of type string".format(property_type_name)
                properties.append(
                    {"propertyTypeID": property_type_id, "stringValue": value}
                )
            elif type == "string":
                assert isinstance(value, str) or isinstance(
                    value, bytes
                ), "property {} is not of type string".format(property_type_name)
                properties.append(
                    {"propertyTypeID": property_type_id, "stringValue": value}
                )
            elif type == "bool":
                assert isinstance(
                    value, bool
                ), "property {} is not of type bool".format(property_type_name)
                properties.append(
                    {"propertyTypeID": property_type_id, "booleanValue": value}
                )
            elif type == "gps_location":
                assert (
                    isinstance(value, tuple)
                    and len(value) == 2
                    and isinstance(value[0], float)
                    and isinstance(value[1], float)
                ), "property {} is not of type tuple(float, float)".format(
                    property_type_name
                )
                properties.append(
                    {
                        "propertyTypeID": property_type_id,
                        "latitudeValue": value[0],
                        "longitudeValue": value[1],
                    }
                )
            else:
                raise Exception(
                    "property type {} has not supported type {}".format(
                        property_type_name, type
                    )
                )
    return properties


def _get_property_value(
    property_type: Dict[str, Any], property: Dict[str, Any]
) -> Union[date, float, int, str, bool, Tuple[float, float]]:
    if property_type["type"] == "gps_location":
        return (property["latitudeValue"], property["longitudeValue"])
    else:
        for property_type_name, field_name in PROPERTY_TYPE_TO_FIELD_NAME.items():
            if property_type["type"] == property_type_name:
                value = property[field_name]
                if property_type_name == "date":
                    value = datetime.strptime(value, "%Y-%m-%d").date()
                # pyre-fixme[7]: Expected `Union[bool, date, float, int, str,
                #  Tuple[float, float]]` but got implicit return value of `None`.
                return value


def _get_property_default_value(
    name: str, type: str, value: PropertyValue
) -> Dict[str, Any]:
    if value is None:
        return {}
    if type == "date":
        assert isinstance(
            value, date
        ), "property {} is not of type datetime.date".format(name)
        return {"stringValue": str(value)}
    elif type == "float":
        assert isinstance(value, float), "property {} is not of type float".format(name)
        return {"floatValue": value}
    elif type == "int":
        assert isinstance(value, int), "property {} is not of type int".format(name)
        return {"intValue": value}
    elif type == "email":
        assert isinstance(value, str), "property {} is not of type string".format(name)
        return {"stringValue": value}
    elif type == "string":
        assert isinstance(value, str) or isinstance(
            value, bytes
        ), "property {} is not of type string".format(name)
        return {"stringValue": value}
    elif type == "bool":
        assert isinstance(value, bool), "property {} is not of type bool".format(name)
        return {"booleanValue": value}
    elif type == "gps_location":
        assert (
            isinstance(value, tuple)
            and len(value) == 2
            and isinstance(value[0], float)
            and isinstance(value[1], float)
        ), "property {} is not of type tuple(float, float)".format(name)
        return {"latitudeValue": value[0], "longitudeValue": value[1]}
    else:
        raise Exception("property type {} has not supported type {}".format(name, type))


def _make_property_types(
    properties: List[Tuple[str, str, PropertyValue, bool]]
) -> List[Dict[str, Any]]:
    property_types = [
        {
            "name": arg[0],
            "type": arg[1],
            "index": i,
            **_get_property_default_value(arg[0], arg[1], arg[2]),
            "isInstanceProperty": arg[3],
        }
        for i, arg in enumerate(properties)
    ]
    return property_types
