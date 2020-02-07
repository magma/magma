---
id: inventory-intro
title: Intro
---

> Manage and track your equipment, deployments and connections.

### [Configure](/inventory/configure)
- Add or edit location types. These can be countries, cities, streets etc., their hierarchy in relation to other location types and properties unique to them.
- Add or edit equipment types. These can be splitters, cards, antennas etc., their ports, positions and properties.

### [Map](/inventory/map)
- See where your network is deployed. Location instances with latitude & longitude will be displayed here.

### [Inventory](/inventory/inventory)
- A tree view of the locations and equipment under each one.
- Add and edit location instances.
- Add, edit and link equipment instances.
- Data model:

```mermaid
graph LR

classDef properties fill:#f96;
classDef types fill:#aaa;

    id1[LOCATION TYPE]:::types <==> id2((LOCATION))
    id3[EQUIPMENT TYPE]:::types <==> id4((EQUIPMENT))
    id3[EQUIPMENT TYPE]:::types <==> id3.1[POSITION DEFINITION]:::types
    id3[EQUIPMENT TYPE]:::types <==> id3.2[PORT DEFINITION]:::types
    id3.2[PORT DEFINITION]:::types <--> id15[PORT TYPE]:::types

    id4((EQUIPMENT)) ==> id2{LOCATION}
    id3.1[POSITION DEFINITION]:::types <==> id5(EQUIPMENT POSITION)
    id3.2[PORT DEFINITION]:::types <==> id6(EQUIPMENT PORT)
    id4((EQUIPMENT)) <==> id5(EQUIPMENT POSITION)
    id4((EQUIPMENT)) <==> id6(EQUIPMENT PORT)
    id5(EQUIPMENT POSITION)<==>|attachment| id7((EQUIPMENT))
    id6(EQUIPMENT PORT)<==>|side a| id8(LINK)
    id9(EQUIPMENT PORT)<==>|side b| id8(LINK)

    id10[PROPERTY TYPE]:::properties -.-> id1[LOCATION TYPE]:::types
    id11[PROPERTY TYPE]:::properties -.-> id3[EQUIPMENT TYPE]:::types
    id14[PROPERTY TYPE]:::properties -.-> id15[PORT TYPE]:::types

    id10[PROPERTY TYPE]:::properties -.-> id12(PROPERTY):::properties
    id12(PROPERTY):::properties -.-> id2((LOCATION))

    id11[PROPERTY TYPE]:::properties -.-> id13(PROPERTY):::properties
    id13(PROPERTY):::properties -.-> id4((EQUIPMENT))

    id14[PROPERTY TYPE]:::properties -.-> id16(PROPERTY):::properties
    id16(PROPERTY):::properties -.-> id6(EQUIPMENT PORT)

```

### [Search](/inventory/search)
- Filter equipment by name/type/properties/locations etc..
- Export the list to CSV.

### [Services](/inventory/services)
- Filter service by name/type/properties etc..
- Export the list to CSV.
- Add or edit service instances

# Workforce Management
> Manage and track your projects and work orders.

### [Configure](/workorders/configure)
- Add or edit work order templates.
- Add or edit project templates. Project templates allow you to quickly create similar projects containing the same work order types.


### [Projects View](/workorders/projects)
- See all projects in either list or map view.
- Add or edit project instances.

### [Work Orders Search](/workorders/search)
- Filter work orders by type/location/name etc. View them in either list or map view.
- Add or edit work order instances.
