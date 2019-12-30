/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ReaderFragment } from 'relay-runtime';
type EquipmentBreadcrumbs_equipment$ref = any;
export type FutureState = "INSTALL" | "REMOVE" | "%future added value";
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "equipment" | "float" | "gps_location" | "int" | "location" | "range" | "service" | "string" | "%future added value";
export type WorkOrderStatus = "DONE" | "PENDING" | "PLANNED" | "%future added value";
import type { FragmentReference } from "relay-runtime";
declare export opaque type EquipmentPortsTable_equipment$ref: FragmentReference;
declare export opaque type EquipmentPortsTable_equipment$fragmentType: EquipmentPortsTable_equipment$ref;
export type EquipmentPortsTable_equipment = {|
  +id: string,
  +name: string,
  +equipmentType: {|
    +id: string,
    +name: string,
    +portDefinitions: $ReadOnlyArray<?{|
      +id: string,
      +name: string,
      +index: ?number,
      +visibleLabel: ?string,
      +portType: ?{|
        +id: string,
        +name: string,
        +propertyTypes: $ReadOnlyArray<?{|
          +id: string,
          +name: string,
          +type: PropertyKind,
          +index: ?number,
          +stringValue: ?string,
          +intValue: ?number,
          +booleanValue: ?boolean,
          +floatValue: ?number,
          +latitudeValue: ?number,
          +longitudeValue: ?number,
          +rangeFromValue: ?number,
          +rangeToValue: ?number,
          +isEditable: ?boolean,
          +isInstanceProperty: ?boolean,
          +isMandatory: ?boolean,
        |}>,
        +linkPropertyTypes: $ReadOnlyArray<?{|
          +id: string,
          +name: string,
          +type: PropertyKind,
          +index: ?number,
          +stringValue: ?string,
          +intValue: ?number,
          +booleanValue: ?boolean,
          +floatValue: ?number,
          +latitudeValue: ?number,
          +longitudeValue: ?number,
          +rangeFromValue: ?number,
          +rangeToValue: ?number,
          +isEditable: ?boolean,
          +isInstanceProperty: ?boolean,
          +isMandatory: ?boolean,
        |}>,
      |},
    |}>,
  |},
  +ports: $ReadOnlyArray<?{|
    +id: string,
    +definition: {|
      +id: string,
      +name: string,
      +index: ?number,
      +visibleLabel: ?string,
      +portType: ?{|
        +id: string,
        +name: string,
        +propertyTypes: $ReadOnlyArray<?{|
          +id: string,
          +name: string,
          +type: PropertyKind,
          +index: ?number,
          +stringValue: ?string,
          +intValue: ?number,
          +booleanValue: ?boolean,
          +floatValue: ?number,
          +latitudeValue: ?number,
          +longitudeValue: ?number,
          +rangeFromValue: ?number,
          +rangeToValue: ?number,
          +isEditable: ?boolean,
          +isInstanceProperty: ?boolean,
          +isMandatory: ?boolean,
        |}>,
        +linkPropertyTypes: $ReadOnlyArray<?{|
          +id: string,
          +name: string,
          +type: PropertyKind,
          +index: ?number,
          +stringValue: ?string,
          +intValue: ?number,
          +booleanValue: ?boolean,
          +floatValue: ?number,
          +latitudeValue: ?number,
          +longitudeValue: ?number,
          +rangeFromValue: ?number,
          +rangeToValue: ?number,
          +isEditable: ?boolean,
          +isInstanceProperty: ?boolean,
          +isMandatory: ?boolean,
        |}>,
      |},
    |},
    +parentEquipment: {|
      +id: string,
      +name: string,
      +equipmentType: {|
        +id: string,
        +name: string,
        +portDefinitions: $ReadOnlyArray<?{|
          +id: string,
          +name: string,
          +visibleLabel: ?string,
          +portType: ?{|
            +id: string,
            +name: string,
          |},
          +bandwidth: ?string,
        |}>,
      |},
    |},
    +link: ?{|
      +id: string,
      +futureState: ?FutureState,
      +ports: $ReadOnlyArray<?{|
        +id: string,
        +definition: {|
          +id: string,
          +name: string,
          +visibleLabel: ?string,
          +portType: ?{|
            +linkPropertyTypes: $ReadOnlyArray<?{|
              +id: string,
              +name: string,
              +type: PropertyKind,
              +index: ?number,
              +stringValue: ?string,
              +intValue: ?number,
              +booleanValue: ?boolean,
              +floatValue: ?number,
              +latitudeValue: ?number,
              +longitudeValue: ?number,
              +rangeFromValue: ?number,
              +rangeToValue: ?number,
              +isEditable: ?boolean,
              +isInstanceProperty: ?boolean,
              +isMandatory: ?boolean,
            |}>
          |},
        |},
        +parentEquipment: {|
          +id: string,
          +name: string,
          +futureState: ?FutureState,
          +equipmentType: {|
            +id: string,
            +name: string,
            +portDefinitions: $ReadOnlyArray<?{|
              +id: string,
              +name: string,
              +visibleLabel: ?string,
              +bandwidth: ?string,
              +portType: ?{|
                +id: string,
                +name: string,
              |},
            |}>,
          |},
          +$fragmentRefs: EquipmentBreadcrumbs_equipment$ref,
        |},
      |}>,
      +workOrder: ?{|
        +id: string,
        +status: WorkOrderStatus,
      |},
      +properties: $ReadOnlyArray<?{|
        +id: string,
        +propertyType: {|
          +id: string,
          +name: string,
          +type: PropertyKind,
          +isEditable: ?boolean,
          +isMandatory: ?boolean,
          +isInstanceProperty: ?boolean,
          +stringValue: ?string,
        |},
        +stringValue: ?string,
        +intValue: ?number,
        +floatValue: ?number,
        +booleanValue: ?boolean,
        +latitudeValue: ?number,
        +longitudeValue: ?number,
        +rangeFromValue: ?number,
        +rangeToValue: ?number,
        +equipmentValue: ?{|
          +id: string,
          +name: string,
        |},
        +locationValue: ?{|
          +id: string,
          +name: string,
        |},
      |}>,
      +services: $ReadOnlyArray<?{|
        +id: string,
        +name: string,
      |}>,
    |},
    +properties: $ReadOnlyArray<{|
      +id: string,
      +propertyType: {|
        +id: string,
        +name: string,
        +type: PropertyKind,
        +isEditable: ?boolean,
        +isMandatory: ?boolean,
        +isInstanceProperty: ?boolean,
        +stringValue: ?string,
      |},
      +stringValue: ?string,
      +intValue: ?number,
      +floatValue: ?number,
      +booleanValue: ?boolean,
      +latitudeValue: ?number,
      +longitudeValue: ?number,
      +rangeFromValue: ?number,
      +rangeToValue: ?number,
      +equipmentValue: ?{|
        +id: string,
        +name: string,
      |},
      +locationValue: ?{|
        +id: string,
        +name: string,
      |},
    |}>,
  |}>,
  +positions: $ReadOnlyArray<?{|
    +attachedEquipment: ?{|
      +id: string,
      +name: string,
      +ports: $ReadOnlyArray<?{|
        +id: string,
        +definition: {|
          +id: string,
          +name: string,
          +index: ?number,
          +visibleLabel: ?string,
          +portType: ?{|
            +id: string,
            +name: string,
            +propertyTypes: $ReadOnlyArray<?{|
              +id: string,
              +name: string,
              +type: PropertyKind,
              +index: ?number,
              +stringValue: ?string,
              +intValue: ?number,
              +booleanValue: ?boolean,
              +floatValue: ?number,
              +latitudeValue: ?number,
              +longitudeValue: ?number,
              +rangeFromValue: ?number,
              +rangeToValue: ?number,
              +isEditable: ?boolean,
              +isInstanceProperty: ?boolean,
              +isMandatory: ?boolean,
            |}>,
            +linkPropertyTypes: $ReadOnlyArray<?{|
              +id: string,
              +name: string,
              +type: PropertyKind,
              +index: ?number,
              +stringValue: ?string,
              +intValue: ?number,
              +booleanValue: ?boolean,
              +floatValue: ?number,
              +latitudeValue: ?number,
              +longitudeValue: ?number,
              +rangeFromValue: ?number,
              +rangeToValue: ?number,
              +isEditable: ?boolean,
              +isInstanceProperty: ?boolean,
              +isMandatory: ?boolean,
            |}>,
          |},
        |},
        +parentEquipment: {|
          +id: string,
          +name: string,
          +equipmentType: {|
            +id: string,
            +name: string,
            +portDefinitions: $ReadOnlyArray<?{|
              +id: string,
              +name: string,
              +visibleLabel: ?string,
              +portType: ?{|
                +id: string,
                +name: string,
              |},
              +bandwidth: ?string,
            |}>,
          |},
        |},
        +link: ?{|
          +id: string,
          +futureState: ?FutureState,
          +ports: $ReadOnlyArray<?{|
            +id: string,
            +definition: {|
              +id: string,
              +name: string,
              +visibleLabel: ?string,
              +portType: ?{|
                +linkPropertyTypes: $ReadOnlyArray<?{|
                  +id: string,
                  +name: string,
                  +type: PropertyKind,
                  +index: ?number,
                  +stringValue: ?string,
                  +intValue: ?number,
                  +booleanValue: ?boolean,
                  +floatValue: ?number,
                  +latitudeValue: ?number,
                  +longitudeValue: ?number,
                  +rangeFromValue: ?number,
                  +rangeToValue: ?number,
                  +isEditable: ?boolean,
                  +isInstanceProperty: ?boolean,
                  +isMandatory: ?boolean,
                |}>
              |},
            |},
            +parentEquipment: {|
              +id: string,
              +name: string,
              +futureState: ?FutureState,
              +equipmentType: {|
                +id: string,
                +name: string,
                +portDefinitions: $ReadOnlyArray<?{|
                  +id: string,
                  +name: string,
                  +visibleLabel: ?string,
                  +bandwidth: ?string,
                  +portType: ?{|
                    +id: string,
                    +name: string,
                  |},
                |}>,
              |},
              +$fragmentRefs: EquipmentBreadcrumbs_equipment$ref,
            |},
          |}>,
          +workOrder: ?{|
            +id: string,
            +status: WorkOrderStatus,
          |},
          +properties: $ReadOnlyArray<?{|
            +id: string,
            +propertyType: {|
              +id: string,
              +name: string,
              +type: PropertyKind,
              +isEditable: ?boolean,
              +isMandatory: ?boolean,
              +isInstanceProperty: ?boolean,
              +stringValue: ?string,
            |},
            +stringValue: ?string,
            +intValue: ?number,
            +floatValue: ?number,
            +booleanValue: ?boolean,
            +latitudeValue: ?number,
            +longitudeValue: ?number,
            +rangeFromValue: ?number,
            +rangeToValue: ?number,
            +equipmentValue: ?{|
              +id: string,
              +name: string,
            |},
            +locationValue: ?{|
              +id: string,
              +name: string,
            |},
          |}>,
          +services: $ReadOnlyArray<?{|
            +id: string,
            +name: string,
          |}>,
        |},
        +properties: $ReadOnlyArray<{|
          +id: string,
          +propertyType: {|
            +id: string,
            +name: string,
            +type: PropertyKind,
            +isEditable: ?boolean,
            +isMandatory: ?boolean,
            +isInstanceProperty: ?boolean,
            +stringValue: ?string,
          |},
          +stringValue: ?string,
          +intValue: ?number,
          +floatValue: ?number,
          +booleanValue: ?boolean,
          +latitudeValue: ?number,
          +longitudeValue: ?number,
          +rangeFromValue: ?number,
          +rangeToValue: ?number,
          +equipmentValue: ?{|
            +id: string,
            +name: string,
          |},
          +locationValue: ?{|
            +id: string,
            +name: string,
          |},
        |}>,
      |}>,
      +equipmentType: {|
        +portDefinitions: $ReadOnlyArray<?{|
          +id: string,
          +name: string,
          +visibleLabel: ?string,
          +bandwidth: ?string,
        |}>
      |},
      +positions: $ReadOnlyArray<?{|
        +attachedEquipment: ?{|
          +id: string,
          +name: string,
          +ports: $ReadOnlyArray<?{|
            +id: string,
            +definition: {|
              +id: string,
              +name: string,
              +index: ?number,
              +visibleLabel: ?string,
              +portType: ?{|
                +id: string,
                +name: string,
                +propertyTypes: $ReadOnlyArray<?{|
                  +id: string,
                  +name: string,
                  +type: PropertyKind,
                  +index: ?number,
                  +stringValue: ?string,
                  +intValue: ?number,
                  +booleanValue: ?boolean,
                  +floatValue: ?number,
                  +latitudeValue: ?number,
                  +longitudeValue: ?number,
                  +rangeFromValue: ?number,
                  +rangeToValue: ?number,
                  +isEditable: ?boolean,
                  +isInstanceProperty: ?boolean,
                  +isMandatory: ?boolean,
                |}>,
                +linkPropertyTypes: $ReadOnlyArray<?{|
                  +id: string,
                  +name: string,
                  +type: PropertyKind,
                  +index: ?number,
                  +stringValue: ?string,
                  +intValue: ?number,
                  +booleanValue: ?boolean,
                  +floatValue: ?number,
                  +latitudeValue: ?number,
                  +longitudeValue: ?number,
                  +rangeFromValue: ?number,
                  +rangeToValue: ?number,
                  +isEditable: ?boolean,
                  +isInstanceProperty: ?boolean,
                  +isMandatory: ?boolean,
                |}>,
              |},
            |},
            +parentEquipment: {|
              +id: string,
              +name: string,
              +equipmentType: {|
                +id: string,
                +name: string,
                +portDefinitions: $ReadOnlyArray<?{|
                  +id: string,
                  +name: string,
                  +visibleLabel: ?string,
                  +portType: ?{|
                    +id: string,
                    +name: string,
                  |},
                  +bandwidth: ?string,
                |}>,
              |},
            |},
            +link: ?{|
              +id: string,
              +futureState: ?FutureState,
              +ports: $ReadOnlyArray<?{|
                +id: string,
                +definition: {|
                  +id: string,
                  +name: string,
                  +visibleLabel: ?string,
                  +portType: ?{|
                    +linkPropertyTypes: $ReadOnlyArray<?{|
                      +id: string,
                      +name: string,
                      +type: PropertyKind,
                      +index: ?number,
                      +stringValue: ?string,
                      +intValue: ?number,
                      +booleanValue: ?boolean,
                      +floatValue: ?number,
                      +latitudeValue: ?number,
                      +longitudeValue: ?number,
                      +rangeFromValue: ?number,
                      +rangeToValue: ?number,
                      +isEditable: ?boolean,
                      +isInstanceProperty: ?boolean,
                      +isMandatory: ?boolean,
                    |}>
                  |},
                |},
                +parentEquipment: {|
                  +id: string,
                  +name: string,
                  +futureState: ?FutureState,
                  +equipmentType: {|
                    +id: string,
                    +name: string,
                    +portDefinitions: $ReadOnlyArray<?{|
                      +id: string,
                      +name: string,
                      +visibleLabel: ?string,
                      +bandwidth: ?string,
                      +portType: ?{|
                        +id: string,
                        +name: string,
                      |},
                    |}>,
                  |},
                  +$fragmentRefs: EquipmentBreadcrumbs_equipment$ref,
                |},
              |}>,
              +workOrder: ?{|
                +id: string,
                +status: WorkOrderStatus,
              |},
              +properties: $ReadOnlyArray<?{|
                +id: string,
                +propertyType: {|
                  +id: string,
                  +name: string,
                  +type: PropertyKind,
                  +isEditable: ?boolean,
                  +isMandatory: ?boolean,
                  +isInstanceProperty: ?boolean,
                  +stringValue: ?string,
                |},
                +stringValue: ?string,
                +intValue: ?number,
                +floatValue: ?number,
                +booleanValue: ?boolean,
                +latitudeValue: ?number,
                +longitudeValue: ?number,
                +rangeFromValue: ?number,
                +rangeToValue: ?number,
                +equipmentValue: ?{|
                  +id: string,
                  +name: string,
                |},
                +locationValue: ?{|
                  +id: string,
                  +name: string,
                |},
              |}>,
              +services: $ReadOnlyArray<?{|
                +id: string,
                +name: string,
              |}>,
            |},
            +properties: $ReadOnlyArray<{|
              +id: string,
              +propertyType: {|
                +id: string,
                +name: string,
                +type: PropertyKind,
                +isEditable: ?boolean,
                +isMandatory: ?boolean,
                +isInstanceProperty: ?boolean,
                +stringValue: ?string,
              |},
              +stringValue: ?string,
              +intValue: ?number,
              +floatValue: ?number,
              +booleanValue: ?boolean,
              +latitudeValue: ?number,
              +longitudeValue: ?number,
              +rangeFromValue: ?number,
              +rangeToValue: ?number,
              +equipmentValue: ?{|
                +id: string,
                +name: string,
              |},
              +locationValue: ?{|
                +id: string,
                +name: string,
              |},
            |}>,
          |}>,
          +equipmentType: {|
            +portDefinitions: $ReadOnlyArray<?{|
              +id: string,
              +name: string,
              +visibleLabel: ?string,
              +bandwidth: ?string,
            |}>
          |},
          +positions: $ReadOnlyArray<?{|
            +attachedEquipment: ?{|
              +id: string,
              +name: string,
              +ports: $ReadOnlyArray<?{|
                +id: string,
                +definition: {|
                  +id: string,
                  +name: string,
                  +index: ?number,
                  +visibleLabel: ?string,
                  +portType: ?{|
                    +id: string,
                    +name: string,
                    +propertyTypes: $ReadOnlyArray<?{|
                      +id: string,
                      +name: string,
                      +type: PropertyKind,
                      +index: ?number,
                      +stringValue: ?string,
                      +intValue: ?number,
                      +booleanValue: ?boolean,
                      +floatValue: ?number,
                      +latitudeValue: ?number,
                      +longitudeValue: ?number,
                      +rangeFromValue: ?number,
                      +rangeToValue: ?number,
                      +isEditable: ?boolean,
                      +isInstanceProperty: ?boolean,
                      +isMandatory: ?boolean,
                    |}>,
                    +linkPropertyTypes: $ReadOnlyArray<?{|
                      +id: string,
                      +name: string,
                      +type: PropertyKind,
                      +index: ?number,
                      +stringValue: ?string,
                      +intValue: ?number,
                      +booleanValue: ?boolean,
                      +floatValue: ?number,
                      +latitudeValue: ?number,
                      +longitudeValue: ?number,
                      +rangeFromValue: ?number,
                      +rangeToValue: ?number,
                      +isEditable: ?boolean,
                      +isInstanceProperty: ?boolean,
                      +isMandatory: ?boolean,
                    |}>,
                  |},
                |},
                +parentEquipment: {|
                  +id: string,
                  +name: string,
                  +equipmentType: {|
                    +id: string,
                    +name: string,
                    +portDefinitions: $ReadOnlyArray<?{|
                      +id: string,
                      +name: string,
                      +visibleLabel: ?string,
                      +portType: ?{|
                        +id: string,
                        +name: string,
                      |},
                      +bandwidth: ?string,
                    |}>,
                  |},
                |},
                +link: ?{|
                  +id: string,
                  +futureState: ?FutureState,
                  +ports: $ReadOnlyArray<?{|
                    +id: string,
                    +definition: {|
                      +id: string,
                      +name: string,
                      +visibleLabel: ?string,
                      +portType: ?{|
                        +linkPropertyTypes: $ReadOnlyArray<?{|
                          +id: string,
                          +name: string,
                          +type: PropertyKind,
                          +index: ?number,
                          +stringValue: ?string,
                          +intValue: ?number,
                          +booleanValue: ?boolean,
                          +floatValue: ?number,
                          +latitudeValue: ?number,
                          +longitudeValue: ?number,
                          +rangeFromValue: ?number,
                          +rangeToValue: ?number,
                          +isEditable: ?boolean,
                          +isInstanceProperty: ?boolean,
                          +isMandatory: ?boolean,
                        |}>
                      |},
                    |},
                    +parentEquipment: {|
                      +id: string,
                      +name: string,
                      +futureState: ?FutureState,
                      +equipmentType: {|
                        +id: string,
                        +name: string,
                        +portDefinitions: $ReadOnlyArray<?{|
                          +id: string,
                          +name: string,
                          +visibleLabel: ?string,
                          +bandwidth: ?string,
                          +portType: ?{|
                            +id: string,
                            +name: string,
                          |},
                        |}>,
                      |},
                      +$fragmentRefs: EquipmentBreadcrumbs_equipment$ref,
                    |},
                  |}>,
                  +workOrder: ?{|
                    +id: string,
                    +status: WorkOrderStatus,
                  |},
                  +properties: $ReadOnlyArray<?{|
                    +id: string,
                    +propertyType: {|
                      +id: string,
                      +name: string,
                      +type: PropertyKind,
                      +isEditable: ?boolean,
                      +isMandatory: ?boolean,
                      +isInstanceProperty: ?boolean,
                      +stringValue: ?string,
                    |},
                    +stringValue: ?string,
                    +intValue: ?number,
                    +floatValue: ?number,
                    +booleanValue: ?boolean,
                    +latitudeValue: ?number,
                    +longitudeValue: ?number,
                    +rangeFromValue: ?number,
                    +rangeToValue: ?number,
                    +equipmentValue: ?{|
                      +id: string,
                      +name: string,
                    |},
                    +locationValue: ?{|
                      +id: string,
                      +name: string,
                    |},
                  |}>,
                  +services: $ReadOnlyArray<?{|
                    +id: string,
                    +name: string,
                  |}>,
                |},
                +properties: $ReadOnlyArray<{|
                  +id: string,
                  +propertyType: {|
                    +id: string,
                    +name: string,
                    +type: PropertyKind,
                    +isEditable: ?boolean,
                    +isMandatory: ?boolean,
                    +isInstanceProperty: ?boolean,
                    +stringValue: ?string,
                  |},
                  +stringValue: ?string,
                  +intValue: ?number,
                  +floatValue: ?number,
                  +booleanValue: ?boolean,
                  +latitudeValue: ?number,
                  +longitudeValue: ?number,
                  +rangeFromValue: ?number,
                  +rangeToValue: ?number,
                  +equipmentValue: ?{|
                    +id: string,
                    +name: string,
                  |},
                  +locationValue: ?{|
                    +id: string,
                    +name: string,
                  |},
                |}>,
              |}>,
              +equipmentType: {|
                +portDefinitions: $ReadOnlyArray<?{|
                  +id: string,
                  +name: string,
                  +visibleLabel: ?string,
                  +bandwidth: ?string,
                |}>
              |},
              +positions: $ReadOnlyArray<?{|
                +attachedEquipment: ?{|
                  +id: string,
                  +name: string,
                  +ports: $ReadOnlyArray<?{|
                    +id: string,
                    +definition: {|
                      +id: string,
                      +name: string,
                      +index: ?number,
                      +visibleLabel: ?string,
                      +portType: ?{|
                        +id: string,
                        +name: string,
                        +propertyTypes: $ReadOnlyArray<?{|
                          +id: string,
                          +name: string,
                          +type: PropertyKind,
                          +index: ?number,
                          +stringValue: ?string,
                          +intValue: ?number,
                          +booleanValue: ?boolean,
                          +floatValue: ?number,
                          +latitudeValue: ?number,
                          +longitudeValue: ?number,
                          +rangeFromValue: ?number,
                          +rangeToValue: ?number,
                          +isEditable: ?boolean,
                          +isInstanceProperty: ?boolean,
                          +isMandatory: ?boolean,
                        |}>,
                        +linkPropertyTypes: $ReadOnlyArray<?{|
                          +id: string,
                          +name: string,
                          +type: PropertyKind,
                          +index: ?number,
                          +stringValue: ?string,
                          +intValue: ?number,
                          +booleanValue: ?boolean,
                          +floatValue: ?number,
                          +latitudeValue: ?number,
                          +longitudeValue: ?number,
                          +rangeFromValue: ?number,
                          +rangeToValue: ?number,
                          +isEditable: ?boolean,
                          +isInstanceProperty: ?boolean,
                          +isMandatory: ?boolean,
                        |}>,
                      |},
                    |},
                    +parentEquipment: {|
                      +id: string,
                      +name: string,
                      +equipmentType: {|
                        +id: string,
                        +name: string,
                        +portDefinitions: $ReadOnlyArray<?{|
                          +id: string,
                          +name: string,
                          +visibleLabel: ?string,
                          +portType: ?{|
                            +id: string,
                            +name: string,
                          |},
                          +bandwidth: ?string,
                        |}>,
                      |},
                    |},
                    +link: ?{|
                      +id: string,
                      +futureState: ?FutureState,
                      +ports: $ReadOnlyArray<?{|
                        +id: string,
                        +definition: {|
                          +id: string,
                          +name: string,
                          +visibleLabel: ?string,
                          +portType: ?{|
                            +linkPropertyTypes: $ReadOnlyArray<?{|
                              +id: string,
                              +name: string,
                              +type: PropertyKind,
                              +index: ?number,
                              +stringValue: ?string,
                              +intValue: ?number,
                              +booleanValue: ?boolean,
                              +floatValue: ?number,
                              +latitudeValue: ?number,
                              +longitudeValue: ?number,
                              +rangeFromValue: ?number,
                              +rangeToValue: ?number,
                              +isEditable: ?boolean,
                              +isInstanceProperty: ?boolean,
                              +isMandatory: ?boolean,
                            |}>
                          |},
                        |},
                        +parentEquipment: {|
                          +id: string,
                          +name: string,
                          +futureState: ?FutureState,
                          +equipmentType: {|
                            +id: string,
                            +name: string,
                            +portDefinitions: $ReadOnlyArray<?{|
                              +id: string,
                              +name: string,
                              +visibleLabel: ?string,
                              +bandwidth: ?string,
                              +portType: ?{|
                                +id: string,
                                +name: string,
                              |},
                            |}>,
                          |},
                          +$fragmentRefs: EquipmentBreadcrumbs_equipment$ref,
                        |},
                      |}>,
                      +workOrder: ?{|
                        +id: string,
                        +status: WorkOrderStatus,
                      |},
                      +properties: $ReadOnlyArray<?{|
                        +id: string,
                        +propertyType: {|
                          +id: string,
                          +name: string,
                          +type: PropertyKind,
                          +isEditable: ?boolean,
                          +isMandatory: ?boolean,
                          +isInstanceProperty: ?boolean,
                          +stringValue: ?string,
                        |},
                        +stringValue: ?string,
                        +intValue: ?number,
                        +floatValue: ?number,
                        +booleanValue: ?boolean,
                        +latitudeValue: ?number,
                        +longitudeValue: ?number,
                        +rangeFromValue: ?number,
                        +rangeToValue: ?number,
                        +equipmentValue: ?{|
                          +id: string,
                          +name: string,
                        |},
                        +locationValue: ?{|
                          +id: string,
                          +name: string,
                        |},
                      |}>,
                      +services: $ReadOnlyArray<?{|
                        +id: string,
                        +name: string,
                      |}>,
                    |},
                    +properties: $ReadOnlyArray<{|
                      +id: string,
                      +propertyType: {|
                        +id: string,
                        +name: string,
                        +type: PropertyKind,
                        +isEditable: ?boolean,
                        +isMandatory: ?boolean,
                        +isInstanceProperty: ?boolean,
                        +stringValue: ?string,
                      |},
                      +stringValue: ?string,
                      +intValue: ?number,
                      +floatValue: ?number,
                      +booleanValue: ?boolean,
                      +latitudeValue: ?number,
                      +longitudeValue: ?number,
                      +rangeFromValue: ?number,
                      +rangeToValue: ?number,
                      +equipmentValue: ?{|
                        +id: string,
                        +name: string,
                      |},
                      +locationValue: ?{|
                        +id: string,
                        +name: string,
                      |},
                    |}>,
                  |}>,
                  +equipmentType: {|
                    +portDefinitions: $ReadOnlyArray<?{|
                      +id: string,
                      +name: string,
                      +visibleLabel: ?string,
                      +bandwidth: ?string,
                    |}>
                  |},
                |}
              |}>,
            |}
          |}>,
        |}
      |}>,
    |}
  |}>,
  +$refType: EquipmentPortsTable_equipment$ref,
|};
export type EquipmentPortsTable_equipment$data = EquipmentPortsTable_equipment;
export type EquipmentPortsTable_equipment$key = {
  +$data?: EquipmentPortsTable_equipment$data,
  +$fragmentRefs: EquipmentPortsTable_equipment$ref,
};
*/


const node/*: ReaderFragment*/ = (function(){
var v0 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "id",
  "args": null,
  "storageKey": null
},
v1 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "name",
  "args": null,
  "storageKey": null
},
v2 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "index",
  "args": null,
  "storageKey": null
},
v3 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "visibleLabel",
  "args": null,
  "storageKey": null
},
v4 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "type",
  "args": null,
  "storageKey": null
},
v5 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "stringValue",
  "args": null,
  "storageKey": null
},
v6 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "intValue",
  "args": null,
  "storageKey": null
},
v7 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "booleanValue",
  "args": null,
  "storageKey": null
},
v8 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "floatValue",
  "args": null,
  "storageKey": null
},
v9 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "latitudeValue",
  "args": null,
  "storageKey": null
},
v10 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "longitudeValue",
  "args": null,
  "storageKey": null
},
v11 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeFromValue",
  "args": null,
  "storageKey": null
},
v12 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeToValue",
  "args": null,
  "storageKey": null
},
v13 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isEditable",
  "args": null,
  "storageKey": null
},
v14 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isInstanceProperty",
  "args": null,
  "storageKey": null
},
v15 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isMandatory",
  "args": null,
  "storageKey": null
},
v16 = [
  (v0/*: any*/),
  (v1/*: any*/),
  (v4/*: any*/),
  (v2/*: any*/),
  (v5/*: any*/),
  (v6/*: any*/),
  (v7/*: any*/),
  (v8/*: any*/),
  (v9/*: any*/),
  (v10/*: any*/),
  (v11/*: any*/),
  (v12/*: any*/),
  (v13/*: any*/),
  (v14/*: any*/),
  (v15/*: any*/)
],
v17 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "linkPropertyTypes",
  "storageKey": null,
  "args": null,
  "concreteType": "PropertyType",
  "plural": true,
  "selections": (v16/*: any*/)
},
v18 = [
  (v0/*: any*/),
  (v1/*: any*/),
  (v2/*: any*/),
  (v3/*: any*/),
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "portType",
    "storageKey": null,
    "args": null,
    "concreteType": "EquipmentPortType",
    "plural": false,
    "selections": [
      (v0/*: any*/),
      (v1/*: any*/),
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "propertyTypes",
        "storageKey": null,
        "args": null,
        "concreteType": "PropertyType",
        "plural": true,
        "selections": (v16/*: any*/)
      },
      (v17/*: any*/)
    ]
  }
],
v19 = [
  (v0/*: any*/),
  (v1/*: any*/)
],
v20 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "portType",
  "storageKey": null,
  "args": null,
  "concreteType": "EquipmentPortType",
  "plural": false,
  "selections": (v19/*: any*/)
},
v21 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "bandwidth",
  "args": null,
  "storageKey": null
},
v22 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "futureState",
  "args": null,
  "storageKey": null
},
v23 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "properties",
  "storageKey": null,
  "args": null,
  "concreteType": "Property",
  "plural": true,
  "selections": [
    (v0/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "propertyType",
      "storageKey": null,
      "args": null,
      "concreteType": "PropertyType",
      "plural": false,
      "selections": [
        (v0/*: any*/),
        (v1/*: any*/),
        (v4/*: any*/),
        (v13/*: any*/),
        (v15/*: any*/),
        (v14/*: any*/),
        (v5/*: any*/)
      ]
    },
    (v5/*: any*/),
    (v6/*: any*/),
    (v8/*: any*/),
    (v7/*: any*/),
    (v9/*: any*/),
    (v10/*: any*/),
    (v11/*: any*/),
    (v12/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "equipmentValue",
      "storageKey": null,
      "args": null,
      "concreteType": "Equipment",
      "plural": false,
      "selections": (v19/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "locationValue",
      "storageKey": null,
      "args": null,
      "concreteType": "Location",
      "plural": false,
      "selections": (v19/*: any*/)
    }
  ]
},
v24 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "ports",
  "storageKey": null,
  "args": null,
  "concreteType": "EquipmentPort",
  "plural": true,
  "selections": [
    (v0/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "definition",
      "storageKey": null,
      "args": null,
      "concreteType": "EquipmentPortDefinition",
      "plural": false,
      "selections": (v18/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "parentEquipment",
      "storageKey": null,
      "args": null,
      "concreteType": "Equipment",
      "plural": false,
      "selections": [
        (v0/*: any*/),
        (v1/*: any*/),
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "equipmentType",
          "storageKey": null,
          "args": null,
          "concreteType": "EquipmentType",
          "plural": false,
          "selections": [
            (v0/*: any*/),
            (v1/*: any*/),
            {
              "kind": "LinkedField",
              "alias": null,
              "name": "portDefinitions",
              "storageKey": null,
              "args": null,
              "concreteType": "EquipmentPortDefinition",
              "plural": true,
              "selections": [
                (v0/*: any*/),
                (v1/*: any*/),
                (v3/*: any*/),
                (v20/*: any*/),
                (v21/*: any*/)
              ]
            }
          ]
        }
      ]
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "link",
      "storageKey": null,
      "args": null,
      "concreteType": "Link",
      "plural": false,
      "selections": [
        (v0/*: any*/),
        (v22/*: any*/),
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "ports",
          "storageKey": null,
          "args": null,
          "concreteType": "EquipmentPort",
          "plural": true,
          "selections": [
            (v0/*: any*/),
            {
              "kind": "LinkedField",
              "alias": null,
              "name": "definition",
              "storageKey": null,
              "args": null,
              "concreteType": "EquipmentPortDefinition",
              "plural": false,
              "selections": [
                (v0/*: any*/),
                (v1/*: any*/),
                (v3/*: any*/),
                {
                  "kind": "LinkedField",
                  "alias": null,
                  "name": "portType",
                  "storageKey": null,
                  "args": null,
                  "concreteType": "EquipmentPortType",
                  "plural": false,
                  "selections": [
                    (v17/*: any*/)
                  ]
                }
              ]
            },
            {
              "kind": "LinkedField",
              "alias": null,
              "name": "parentEquipment",
              "storageKey": null,
              "args": null,
              "concreteType": "Equipment",
              "plural": false,
              "selections": [
                (v0/*: any*/),
                (v1/*: any*/),
                (v22/*: any*/),
                {
                  "kind": "LinkedField",
                  "alias": null,
                  "name": "equipmentType",
                  "storageKey": null,
                  "args": null,
                  "concreteType": "EquipmentType",
                  "plural": false,
                  "selections": [
                    (v0/*: any*/),
                    (v1/*: any*/),
                    {
                      "kind": "LinkedField",
                      "alias": null,
                      "name": "portDefinitions",
                      "storageKey": null,
                      "args": null,
                      "concreteType": "EquipmentPortDefinition",
                      "plural": true,
                      "selections": [
                        (v0/*: any*/),
                        (v1/*: any*/),
                        (v3/*: any*/),
                        (v21/*: any*/),
                        (v20/*: any*/)
                      ]
                    }
                  ]
                },
                {
                  "kind": "FragmentSpread",
                  "name": "EquipmentBreadcrumbs_equipment",
                  "args": null
                }
              ]
            }
          ]
        },
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "workOrder",
          "storageKey": null,
          "args": null,
          "concreteType": "WorkOrder",
          "plural": false,
          "selections": [
            (v0/*: any*/),
            {
              "kind": "ScalarField",
              "alias": null,
              "name": "status",
              "args": null,
              "storageKey": null
            }
          ]
        },
        (v23/*: any*/),
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "services",
          "storageKey": null,
          "args": null,
          "concreteType": "Service",
          "plural": true,
          "selections": (v19/*: any*/)
        }
      ]
    },
    (v23/*: any*/)
  ]
},
v25 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "equipmentType",
  "storageKey": null,
  "args": null,
  "concreteType": "EquipmentType",
  "plural": false,
  "selections": [
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "portDefinitions",
      "storageKey": null,
      "args": null,
      "concreteType": "EquipmentPortDefinition",
      "plural": true,
      "selections": [
        (v0/*: any*/),
        (v1/*: any*/),
        (v3/*: any*/),
        (v21/*: any*/)
      ]
    }
  ]
};
return {
  "kind": "Fragment",
  "name": "EquipmentPortsTable_equipment",
  "type": "Equipment",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
    (v1/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "equipmentType",
      "storageKey": null,
      "args": null,
      "concreteType": "EquipmentType",
      "plural": false,
      "selections": [
        (v0/*: any*/),
        (v1/*: any*/),
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "portDefinitions",
          "storageKey": null,
          "args": null,
          "concreteType": "EquipmentPortDefinition",
          "plural": true,
          "selections": (v18/*: any*/)
        }
      ]
    },
    (v24/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "positions",
      "storageKey": null,
      "args": null,
      "concreteType": "EquipmentPosition",
      "plural": true,
      "selections": [
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "attachedEquipment",
          "storageKey": null,
          "args": null,
          "concreteType": "Equipment",
          "plural": false,
          "selections": [
            (v0/*: any*/),
            (v1/*: any*/),
            (v24/*: any*/),
            (v25/*: any*/),
            {
              "kind": "LinkedField",
              "alias": null,
              "name": "positions",
              "storageKey": null,
              "args": null,
              "concreteType": "EquipmentPosition",
              "plural": true,
              "selections": [
                {
                  "kind": "LinkedField",
                  "alias": null,
                  "name": "attachedEquipment",
                  "storageKey": null,
                  "args": null,
                  "concreteType": "Equipment",
                  "plural": false,
                  "selections": [
                    (v0/*: any*/),
                    (v1/*: any*/),
                    (v24/*: any*/),
                    (v25/*: any*/),
                    {
                      "kind": "LinkedField",
                      "alias": null,
                      "name": "positions",
                      "storageKey": null,
                      "args": null,
                      "concreteType": "EquipmentPosition",
                      "plural": true,
                      "selections": [
                        {
                          "kind": "LinkedField",
                          "alias": null,
                          "name": "attachedEquipment",
                          "storageKey": null,
                          "args": null,
                          "concreteType": "Equipment",
                          "plural": false,
                          "selections": [
                            (v0/*: any*/),
                            (v1/*: any*/),
                            (v24/*: any*/),
                            (v25/*: any*/),
                            {
                              "kind": "LinkedField",
                              "alias": null,
                              "name": "positions",
                              "storageKey": null,
                              "args": null,
                              "concreteType": "EquipmentPosition",
                              "plural": true,
                              "selections": [
                                {
                                  "kind": "LinkedField",
                                  "alias": null,
                                  "name": "attachedEquipment",
                                  "storageKey": null,
                                  "args": null,
                                  "concreteType": "Equipment",
                                  "plural": false,
                                  "selections": [
                                    (v0/*: any*/),
                                    (v1/*: any*/),
                                    (v24/*: any*/),
                                    (v25/*: any*/)
                                  ]
                                }
                              ]
                            }
                          ]
                        }
                      ]
                    }
                  ]
                }
              ]
            }
          ]
        }
      ]
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = '3f850609f4f8209cb2e0623956b8386b';
module.exports = node;
