---
id: py-inventory-release-notes
title: Python API Release Notes
---

<!--
***
This is template for release notes
#3 new version number
### Features
### Changes
### Deprecated
### Removed
### Bug fixes
***
-->

<!--
***
## new version number
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
### Changes
- Equipment functionality:
    - `external_id` variable added to functions
        - 'add_equipment`
        - `add_equipment_to_position`
        - `copy_equipment_in_position`
        - `copy_equipment`
        - `get_or_create_equipment`
        - `get_or_create_equipment_in_position`
### Deprecated
### Removed
### Bug fixes
***
-->

***
## 2.5.0 -release date 23.03.2020
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
