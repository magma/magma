---
id: inventory-intro
title: Inventory Management
---
## Overview

(See Figure 1)

The Inventory Management system enables the user to manage, configure and track equipment, location types, deployments, connections and services.
The Inventory Management system contains the following main tabs:

<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Figure+1.png' width=600> <br>
Figure 1: Inventory Management Tabs
<p/>

- Search (1)
- Locations (2)
- Map (3)
- Catalog (4)
- Services (5)

## Search Tab

(See Figure 2)

The Search tab enables the user to search and filter data stored in the Inventory management system and export reports based on this data.
The filtered data (1) can be:

- Equipment
- Links
- Port
- Locations

The data will be filtered according to different options (2) or according to previously saved searches (3).

The results of the filtered search will be displayed in a list (5) in the main display area and they can also be exported (4) to a CSV file.
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Figure+2.png' width=1000> <br>
Figure 2: Search Tab
</p>

## Locations Tab

(See Figure 3)

The Locations tab provides a hierarchical view of locations and equipment (1).

The Locations tab also enables the user to search (2) for locations and for equipment related to a specific location. The search can be done by using parameters such as: location/equipment name or external ID.

Using the Locations tab, it is possible to add new location instances (3) (see Figure 4).
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Figure+3.png' width=600> <br>
Figure 3: Locations Tab
</p>

<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Figure+4.png' width=600> <br>
Figure 4: Adding Locations
</p>

Once a location is added, it can be edited (4) and an equipment instance can be added to it (5) (see Figure 5).

<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Figure+5.png' width=1000> <br>
Figure 5: Editing Location/Adding Equipment
</p>

### Locations and Equipment Data Model

(See Figure 6)

  - Each location has its own Location Type, for example City or State.
  - Each Location Type has its own Property Types, for example, if the location type is “City”, “Mayor of the City” and “Postcode” may be useful property types.
  - Each Location will be created with properties according to the Property Types that were created.
  - Each Equipment has its own Equipment Type, for example Router or Card.
  - Each Equipment Type has its own Property Types, for example Operating System or Manufacturer.
  - Each Equipment has its own Position (one or more).
  - When creating Equipment, each Equipment Type will be created with properties according to the created Property Types and Positions.
  - A Location contains Equipment (one or more).

<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Figure+6.png' width=600> <br>
Figure 6: Locations and Equipment
</p>

## Map Tab

(See Figure 7)

The Map tab enables to view Location (1) in a map view, color coded by types and to search (2) for them. For example, a mobile provider would view mobile towers, offices and data centers that were uploaded to the system with geo-location.

Clicking a location on the map will display additional information about it (3).

The tab view can be switched between a map view to a satellite view (4).
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Figure+7.png' width=1000> <br>
Figure 7: Map Tab
</p>

## Catalog Tab

(See Figure 8)

The Inventory management system is used to create the schema of the system.

The system is used by many partners. Each partner is unique. For example, the usage of different equipment in a variety of different locations by each partner. When creating a system element such as equipment or location and adding their characteristics, a partner is not limited with the information that can be added due to the fact that the information is added as free text.

The Catalog tab enables the user to add or edit the following elements:

- EQUIPMENT (1) - add/edit equipment properties, positions and ports.
- LOCATIONS (2) - add/edit location properties and survey templates.
- PORTS (3) - add/edit port and link properties.
- SERVICES (4) - add/edit service endpoints and properties.

<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Figure+8.png' width=1000> <br>
Figure 8: Catalog Tab
</p>

## Services Tab

(See Figure 9)

The Services tab enables the user to add a service (2) describing the topology between endpoints and their links.

The Services tab also enables the user to filter services (1) by using parameters such as name, type or properties and to export a list of services to a CSV file.
<p align="center">
<img src='https://s3.amazonaws.com/purpleheadband.images/wiki/Inventory+management/Figure+9.png' width=1000> <br>
Figure 9: Services Tab
</p>

