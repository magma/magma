---
id: py-inventory-release-notes
title: Python API Release Notes
---

<!---
***
This is template for release notes
# new version number
## Features:
## Changes:
## Deprecated:
## Removed:
## Bug fixes
***
--->

<!---
***
# new version number
## Features:
- `add_equipment_port_type` - create new equipment port type.
- `get_equipment_port_type` - get existing equipment port type.
- `edit_equipment_port_type` - edit existing equipment port type.
- `delete_equipment_port_type` - delete existing equipment port.
## Changes:
- functions now raise warning if they query against deprecated Graphql Endpoints. If you get such warning you are adviced to upgrade to newer version of API that will call different graphql endpoints instead
- use BasicAuth login to graphql server which improves first connection performance
## Deprecated:
## Deprecated:
## Bug fixes
***
---> 


***
# 2.4.0
## Features:
- `edit_location function` - edit location properties.
- `edit_equipment_type` - edit/add port types.

## Changes:
- `get_locations_by_external_id` raises exception on not found location
## Deprecated:
- `get_locations_by_external_id` function is deprecated by `get_location_by_external_id` function
## Bug fixes
***
