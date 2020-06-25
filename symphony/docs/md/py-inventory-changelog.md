---
id: py-inventory-release-notes
title: Python API Release Notes
---

<!--
***
This is template for release notes
# new version number
 Features
 Changes
 Deprecated
 Removed
 Bug fixes
***
-->


***
## 3.0.0 - release date 10.6.2020
### Features
- ServiceType
    - `add_service_type`
    - `get_service_type`
    - `edit_service_type`
    - `delete_service_type`
    - `delete_service_type_with_services`
- Service
    - `add_service`
    - `add_service_endpoint_definition`
    - `add_service_endpoint`
    - `add_service_link`
    - `get_service`
### Breaking Changes
- Customer
    - `Customer` DataClass attribute `externalId` renamed to `external_id`
    - `add_customer` changed attribute `externalId` to `external_id`
- Document
    - `Document` DataClass attribute `parentId` renamed to `parent_id`
    - `Document` DataClass attribute `parentEntity` renamed to `parent_entity`
- EquipmentType
    - `add_equipment_type` - the type of `properties` attribute `Sequence[Tuple[str, str, Optional[PropertyValue], Optional[bool]]]` changed to `List[PropertyDefinition]`
- LocationType
    - `add_location_type` - the type of `properties` attribute `Sequence[Tuple[str, str, Optional[PropertyValue], Optional[bool]]]` changed to `List[PropertyDefinition]`
- Location
    - `Location` DataClass attribute `externalId` renamed to `external_id`
    - `Location` DataClass attribute `locationTypeName` renamed to `location_type_name`
    - `add_location` changed attribute `externalID` to `external_id`
### Removed
- `get_locations_by_external_id` - deprecated in 2.4.0
### Bug fixes
***


***
## 2.6.1 - release date 23.04.2020
### Changes
- `get_location` and `get_location_by_external_id` performance is improved (`get_location_by_external_id` had 3X time improvement from 0.9 seconds to 0.3 seconds)
### Bug fixes
- Fixed a server breaking change introduced on 15.4.2020. The changes breaks all APIs that add or edit properties for all pyinventory versions.
***


***
## 2.6.0 - release date 14.04.2020
### Features
- Equipment functionality:
    - `get_equipments_by_type`
    - `get_equipments_by_location`
    - `get_equipment_by_external_id`
- EquipmentType functionality:
    - `get_equipment_type_property_type`
    - `get_equipment_type_property_type_by_external_id`
    - `edit_equipment_type_property_type`
- PropertyType functionality:
    - `get_property_type_id`
- PropertyDefinition:
    - `is_mandatory` value added
### Changes
- Equipment functionality:
    - `external_id` variable added to functions
        - `add_equipment`
        - `add_equipment_to_position`
        - `copy_equipment_in_position`
        - `copy_equipment`
        - `get_or_create_equipment`
        - `get_or_create_equipment_in_position`
### Bug fixes
***


***
## 2.5.0 - release date 23.03.2020
### Features

- User functionality:
    - `add_user`
    - `get_user`
    - `edit_user`
    - `deactivate_user`
    - `activate_user`
    - `get_users`
    - `get_active_users`

- Port type functionality:
    - `add_equipment_port_type`
    - `get_equipment_port_type`
    - `edit_equipment_port_type`
    - `delete_equipment_port_type`

- Port functionality:
    - `get_port`
    - `edit_port_properties`
    - `edit_link_properties`

- Equipment functionality:
    - `edit equipment`
    - `get_equipment_properties`
### Changes
- use BasicAuth login to graphql server which improves first connection performance

- functions now raise warning if they query against deprecated Graphql Endpoints. If you get such warning you are adviced to upgrade to newer version of API that will call different graphql endpoints instead
### Bug fixes
***


***
## 2.4.0 - release date 22.02.2020
### Features
- `edit_location function` - edit location properties.
- `edit_equipment_type` - edit/add port types.

### Changes
- `get_locations_by_external_id` raises exception on not found location
### Deprecated
- `get_locations_by_external_id` function is deprecated by `get_location_by_external_id` function
### Bug fixes
***
