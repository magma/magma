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

    id1[LOCATION TYPE___]:::types <==> id2((LOCATION__))
    id3[EQUIPMENT TYPE____]:::types <==> id4((EQUIPMENT___))
    id3[EQUIPMENT TYPE____]:::types <==> id3.1[POSITION DEFINITION____]:::types
    id3[EQUIPMENT TYPE____]:::types <==> id3.2[PORT DEFINITION____]:::types
    id3.2[PORT DEFINITION____]:::types <--> id15[PORT TYPE___]:::types

    id4((EQUIPMENT___)) ==> id2{LOCATION__}
    id3.1[POSITION DEFINITION_____]:::types <==> id5(EQUIPMENT POSITION_____)
    id3.2[PORT DEFINITION____]:::types <==> id6(EQUIPMENT PORT_____)
    id4((EQUIPMENT___)) <==> id5(EQUIPMENT POSITION_____)
    id4((EQUIPMENT___)) <==> id6(EQUIPMENT PORT_____)
    id5(EQUIPMENT POSITION_____)<==>|_attachment____| id7((EQUIPMENT___))
    id6(EQUIPMENT PORT____)<==>|side a__| id8(LINK__)
    id9(EQUIPMENT PORT_____)<==>|side b__| id8(LINK__)
    
    id10[PROPERTY TYPE____]:::properties -.-> id1[LOCATION TYPE___]:::types
    id11[PROPERTY TYPE____]:::properties -.-> id3[EQUIPMENT TYPE____]:::types
    id14[PROPERTY TYPE____]:::properties -.-> id15[PORT TYPE___]:::types

    id10[PROPERTY TYPE____]:::properties -.-> id12(PROPERTY___):::properties
    id12(PROPERTY___):::properties -.-> id2((LOCATION__))

    id11[PROPERTY TYPE____]:::properties -.-> id13(PROPERTY___):::properties
    id13(PROPERTY___):::properties -.-> id4((EQUIPMENT___))

    id14[PROPERTY TYPE____]:::properties -.-> id16(PROPERTY___):::properties
    id16(PROPERTY___):::properties -.-> id6(EQUIPMENT PORT_____)

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
