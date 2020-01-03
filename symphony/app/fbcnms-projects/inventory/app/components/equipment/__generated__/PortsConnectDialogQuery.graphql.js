/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 9fb63c327fbf6314cf44afc293dc814f
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type EquipmentBreadcrumbs_equipment$ref = any;
export type FutureState = "INSTALL" | "REMOVE" | "%future added value";
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "equipment" | "float" | "gps_location" | "int" | "location" | "range" | "service" | "string" | "%future added value";
export type ServiceEndpointRole = "CONSUMER" | "PROVIDER" | "%future added value";
export type WorkOrderStatus = "DONE" | "PENDING" | "PLANNED" | "%future added value";
export type PortsConnectDialogQueryVariables = {|
  equipmentId: string
|};
export type PortsConnectDialogQueryResponse = {|
  +equipment: ?{|
    +id?: string,
    +name?: string,
    +equipmentType?: {|
      +id: string,
      +name: string,
      +portDefinitions: $ReadOnlyArray<?{|
        +id: string,
        +name: string,
        +visibleLabel: ?string,
        +bandwidth: ?string,
      |}>,
    |},
    +positions?: $ReadOnlyArray<?{|
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
              +serviceValue: ?{|
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
            +serviceValue: ?{|
              +id: string,
              +name: string,
            |},
          |}>,
          +serviceEndpoints: $ReadOnlyArray<{|
            +role: ServiceEndpointRole,
            +service: {|
              +name: string
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
                  +serviceValue: ?{|
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
                +serviceValue: ?{|
                  +id: string,
                  +name: string,
                |},
              |}>,
              +serviceEndpoints: $ReadOnlyArray<{|
                +role: ServiceEndpointRole,
                +service: {|
                  +name: string
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
                      +serviceValue: ?{|
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
                    +serviceValue: ?{|
                      +id: string,
                      +name: string,
                    |},
                  |}>,
                  +serviceEndpoints: $ReadOnlyArray<{|
                    +role: ServiceEndpointRole,
                    +service: {|
                      +name: string
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
                          +serviceValue: ?{|
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
                        +serviceValue: ?{|
                          +id: string,
                          +name: string,
                        |},
                      |}>,
                      +serviceEndpoints: $ReadOnlyArray<{|
                        +role: ServiceEndpointRole,
                        +service: {|
                          +name: string
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
    +ports?: $ReadOnlyArray<?{|
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
          +serviceValue: ?{|
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
        +serviceValue: ?{|
          +id: string,
          +name: string,
        |},
      |}>,
      +serviceEndpoints: $ReadOnlyArray<{|
        +role: ServiceEndpointRole,
        +service: {|
          +name: string
        |},
      |}>,
    |}>,
  |}
|};
export type PortsConnectDialogQuery = {|
  variables: PortsConnectDialogQueryVariables,
  response: PortsConnectDialogQueryResponse,
|};
*/


/*
query PortsConnectDialogQuery(
  $equipmentId: ID!
) {
  equipment: node(id: $equipmentId) {
    __typename
    ... on Equipment {
      id
      name
      equipmentType {
        id
        name
        portDefinitions {
          id
          name
          visibleLabel
          bandwidth
        }
      }
      positions {
        attachedEquipment {
          id
          name
          ports {
            id
            definition {
              id
              name
              index
              visibleLabel
              portType {
                id
                name
                propertyTypes {
                  id
                  name
                  type
                  index
                  stringValue
                  intValue
                  booleanValue
                  floatValue
                  latitudeValue
                  longitudeValue
                  rangeFromValue
                  rangeToValue
                  isEditable
                  isInstanceProperty
                  isMandatory
                }
                linkPropertyTypes {
                  id
                  name
                  type
                  index
                  stringValue
                  intValue
                  booleanValue
                  floatValue
                  latitudeValue
                  longitudeValue
                  rangeFromValue
                  rangeToValue
                  isEditable
                  isInstanceProperty
                  isMandatory
                }
              }
            }
            parentEquipment {
              id
              name
              equipmentType {
                id
                name
                portDefinitions {
                  id
                  name
                  visibleLabel
                  portType {
                    id
                    name
                  }
                  bandwidth
                }
              }
            }
            link {
              id
              futureState
              ports {
                id
                definition {
                  id
                  name
                  visibleLabel
                  portType {
                    linkPropertyTypes {
                      id
                      name
                      type
                      index
                      stringValue
                      intValue
                      booleanValue
                      floatValue
                      latitudeValue
                      longitudeValue
                      rangeFromValue
                      rangeToValue
                      isEditable
                      isInstanceProperty
                      isMandatory
                    }
                    id
                  }
                }
                parentEquipment {
                  id
                  name
                  futureState
                  equipmentType {
                    id
                    name
                    portDefinitions {
                      id
                      name
                      visibleLabel
                      bandwidth
                      portType {
                        id
                        name
                      }
                    }
                  }
                  ...EquipmentBreadcrumbs_equipment
                }
              }
              workOrder {
                id
                status
              }
              properties {
                id
                propertyType {
                  id
                  name
                  type
                  isEditable
                  isMandatory
                  isInstanceProperty
                  stringValue
                }
                stringValue
                intValue
                floatValue
                booleanValue
                latitudeValue
                longitudeValue
                rangeFromValue
                rangeToValue
                equipmentValue {
                  id
                  name
                }
                locationValue {
                  id
                  name
                }
                serviceValue {
                  id
                  name
                }
              }
              services {
                id
                name
              }
            }
            properties {
              id
              propertyType {
                id
                name
                type
                isEditable
                isMandatory
                isInstanceProperty
                stringValue
              }
              stringValue
              intValue
              floatValue
              booleanValue
              latitudeValue
              longitudeValue
              rangeFromValue
              rangeToValue
              equipmentValue {
                id
                name
              }
              locationValue {
                id
                name
              }
              serviceValue {
                id
                name
              }
            }
            serviceEndpoints {
              role
              service {
                name
                id
              }
              id
            }
          }
          equipmentType {
            portDefinitions {
              id
              name
              visibleLabel
              bandwidth
            }
            id
          }
          positions {
            attachedEquipment {
              id
              name
              ports {
                id
                definition {
                  id
                  name
                  index
                  visibleLabel
                  portType {
                    id
                    name
                    propertyTypes {
                      id
                      name
                      type
                      index
                      stringValue
                      intValue
                      booleanValue
                      floatValue
                      latitudeValue
                      longitudeValue
                      rangeFromValue
                      rangeToValue
                      isEditable
                      isInstanceProperty
                      isMandatory
                    }
                    linkPropertyTypes {
                      id
                      name
                      type
                      index
                      stringValue
                      intValue
                      booleanValue
                      floatValue
                      latitudeValue
                      longitudeValue
                      rangeFromValue
                      rangeToValue
                      isEditable
                      isInstanceProperty
                      isMandatory
                    }
                  }
                }
                parentEquipment {
                  id
                  name
                  equipmentType {
                    id
                    name
                    portDefinitions {
                      id
                      name
                      visibleLabel
                      portType {
                        id
                        name
                      }
                      bandwidth
                    }
                  }
                }
                link {
                  id
                  futureState
                  ports {
                    id
                    definition {
                      id
                      name
                      visibleLabel
                      portType {
                        linkPropertyTypes {
                          id
                          name
                          type
                          index
                          stringValue
                          intValue
                          booleanValue
                          floatValue
                          latitudeValue
                          longitudeValue
                          rangeFromValue
                          rangeToValue
                          isEditable
                          isInstanceProperty
                          isMandatory
                        }
                        id
                      }
                    }
                    parentEquipment {
                      id
                      name
                      futureState
                      equipmentType {
                        id
                        name
                        portDefinitions {
                          id
                          name
                          visibleLabel
                          bandwidth
                          portType {
                            id
                            name
                          }
                        }
                      }
                      ...EquipmentBreadcrumbs_equipment
                    }
                  }
                  workOrder {
                    id
                    status
                  }
                  properties {
                    id
                    propertyType {
                      id
                      name
                      type
                      isEditable
                      isMandatory
                      isInstanceProperty
                      stringValue
                    }
                    stringValue
                    intValue
                    floatValue
                    booleanValue
                    latitudeValue
                    longitudeValue
                    rangeFromValue
                    rangeToValue
                    equipmentValue {
                      id
                      name
                    }
                    locationValue {
                      id
                      name
                    }
                    serviceValue {
                      id
                      name
                    }
                  }
                  services {
                    id
                    name
                  }
                }
                properties {
                  id
                  propertyType {
                    id
                    name
                    type
                    isEditable
                    isMandatory
                    isInstanceProperty
                    stringValue
                  }
                  stringValue
                  intValue
                  floatValue
                  booleanValue
                  latitudeValue
                  longitudeValue
                  rangeFromValue
                  rangeToValue
                  equipmentValue {
                    id
                    name
                  }
                  locationValue {
                    id
                    name
                  }
                  serviceValue {
                    id
                    name
                  }
                }
                serviceEndpoints {
                  role
                  service {
                    name
                    id
                  }
                  id
                }
              }
              equipmentType {
                portDefinitions {
                  id
                  name
                  visibleLabel
                  bandwidth
                }
                id
              }
              positions {
                attachedEquipment {
                  id
                  name
                  ports {
                    id
                    definition {
                      id
                      name
                      index
                      visibleLabel
                      portType {
                        id
                        name
                        propertyTypes {
                          id
                          name
                          type
                          index
                          stringValue
                          intValue
                          booleanValue
                          floatValue
                          latitudeValue
                          longitudeValue
                          rangeFromValue
                          rangeToValue
                          isEditable
                          isInstanceProperty
                          isMandatory
                        }
                        linkPropertyTypes {
                          id
                          name
                          type
                          index
                          stringValue
                          intValue
                          booleanValue
                          floatValue
                          latitudeValue
                          longitudeValue
                          rangeFromValue
                          rangeToValue
                          isEditable
                          isInstanceProperty
                          isMandatory
                        }
                      }
                    }
                    parentEquipment {
                      id
                      name
                      equipmentType {
                        id
                        name
                        portDefinitions {
                          id
                          name
                          visibleLabel
                          portType {
                            id
                            name
                          }
                          bandwidth
                        }
                      }
                    }
                    link {
                      id
                      futureState
                      ports {
                        id
                        definition {
                          id
                          name
                          visibleLabel
                          portType {
                            linkPropertyTypes {
                              id
                              name
                              type
                              index
                              stringValue
                              intValue
                              booleanValue
                              floatValue
                              latitudeValue
                              longitudeValue
                              rangeFromValue
                              rangeToValue
                              isEditable
                              isInstanceProperty
                              isMandatory
                            }
                            id
                          }
                        }
                        parentEquipment {
                          id
                          name
                          futureState
                          equipmentType {
                            id
                            name
                            portDefinitions {
                              id
                              name
                              visibleLabel
                              bandwidth
                              portType {
                                id
                                name
                              }
                            }
                          }
                          ...EquipmentBreadcrumbs_equipment
                        }
                      }
                      workOrder {
                        id
                        status
                      }
                      properties {
                        id
                        propertyType {
                          id
                          name
                          type
                          isEditable
                          isMandatory
                          isInstanceProperty
                          stringValue
                        }
                        stringValue
                        intValue
                        floatValue
                        booleanValue
                        latitudeValue
                        longitudeValue
                        rangeFromValue
                        rangeToValue
                        equipmentValue {
                          id
                          name
                        }
                        locationValue {
                          id
                          name
                        }
                        serviceValue {
                          id
                          name
                        }
                      }
                      services {
                        id
                        name
                      }
                    }
                    properties {
                      id
                      propertyType {
                        id
                        name
                        type
                        isEditable
                        isMandatory
                        isInstanceProperty
                        stringValue
                      }
                      stringValue
                      intValue
                      floatValue
                      booleanValue
                      latitudeValue
                      longitudeValue
                      rangeFromValue
                      rangeToValue
                      equipmentValue {
                        id
                        name
                      }
                      locationValue {
                        id
                        name
                      }
                      serviceValue {
                        id
                        name
                      }
                    }
                    serviceEndpoints {
                      role
                      service {
                        name
                        id
                      }
                      id
                    }
                  }
                  equipmentType {
                    portDefinitions {
                      id
                      name
                      visibleLabel
                      bandwidth
                    }
                    id
                  }
                  positions {
                    attachedEquipment {
                      id
                      name
                      ports {
                        id
                        definition {
                          id
                          name
                          index
                          visibleLabel
                          portType {
                            id
                            name
                            propertyTypes {
                              id
                              name
                              type
                              index
                              stringValue
                              intValue
                              booleanValue
                              floatValue
                              latitudeValue
                              longitudeValue
                              rangeFromValue
                              rangeToValue
                              isEditable
                              isInstanceProperty
                              isMandatory
                            }
                            linkPropertyTypes {
                              id
                              name
                              type
                              index
                              stringValue
                              intValue
                              booleanValue
                              floatValue
                              latitudeValue
                              longitudeValue
                              rangeFromValue
                              rangeToValue
                              isEditable
                              isInstanceProperty
                              isMandatory
                            }
                          }
                        }
                        parentEquipment {
                          id
                          name
                          equipmentType {
                            id
                            name
                            portDefinitions {
                              id
                              name
                              visibleLabel
                              portType {
                                id
                                name
                              }
                              bandwidth
                            }
                          }
                        }
                        link {
                          id
                          futureState
                          ports {
                            id
                            definition {
                              id
                              name
                              visibleLabel
                              portType {
                                linkPropertyTypes {
                                  id
                                  name
                                  type
                                  index
                                  stringValue
                                  intValue
                                  booleanValue
                                  floatValue
                                  latitudeValue
                                  longitudeValue
                                  rangeFromValue
                                  rangeToValue
                                  isEditable
                                  isInstanceProperty
                                  isMandatory
                                }
                                id
                              }
                            }
                            parentEquipment {
                              id
                              name
                              futureState
                              equipmentType {
                                id
                                name
                                portDefinitions {
                                  id
                                  name
                                  visibleLabel
                                  bandwidth
                                  portType {
                                    id
                                    name
                                  }
                                }
                              }
                              ...EquipmentBreadcrumbs_equipment
                            }
                          }
                          workOrder {
                            id
                            status
                          }
                          properties {
                            id
                            propertyType {
                              id
                              name
                              type
                              isEditable
                              isMandatory
                              isInstanceProperty
                              stringValue
                            }
                            stringValue
                            intValue
                            floatValue
                            booleanValue
                            latitudeValue
                            longitudeValue
                            rangeFromValue
                            rangeToValue
                            equipmentValue {
                              id
                              name
                            }
                            locationValue {
                              id
                              name
                            }
                            serviceValue {
                              id
                              name
                            }
                          }
                          services {
                            id
                            name
                          }
                        }
                        properties {
                          id
                          propertyType {
                            id
                            name
                            type
                            isEditable
                            isMandatory
                            isInstanceProperty
                            stringValue
                          }
                          stringValue
                          intValue
                          floatValue
                          booleanValue
                          latitudeValue
                          longitudeValue
                          rangeFromValue
                          rangeToValue
                          equipmentValue {
                            id
                            name
                          }
                          locationValue {
                            id
                            name
                          }
                          serviceValue {
                            id
                            name
                          }
                        }
                        serviceEndpoints {
                          role
                          service {
                            name
                            id
                          }
                          id
                        }
                      }
                      equipmentType {
                        portDefinitions {
                          id
                          name
                          visibleLabel
                          bandwidth
                        }
                        id
                      }
                    }
                    id
                  }
                }
                id
              }
            }
            id
          }
        }
        id
      }
      ports {
        id
        definition {
          id
          name
          index
          visibleLabel
          portType {
            id
            name
            propertyTypes {
              id
              name
              type
              index
              stringValue
              intValue
              booleanValue
              floatValue
              latitudeValue
              longitudeValue
              rangeFromValue
              rangeToValue
              isEditable
              isInstanceProperty
              isMandatory
            }
            linkPropertyTypes {
              id
              name
              type
              index
              stringValue
              intValue
              booleanValue
              floatValue
              latitudeValue
              longitudeValue
              rangeFromValue
              rangeToValue
              isEditable
              isInstanceProperty
              isMandatory
            }
          }
        }
        parentEquipment {
          id
          name
          equipmentType {
            id
            name
            portDefinitions {
              id
              name
              visibleLabel
              portType {
                id
                name
              }
              bandwidth
            }
          }
        }
        link {
          id
          futureState
          ports {
            id
            definition {
              id
              name
              visibleLabel
              portType {
                linkPropertyTypes {
                  id
                  name
                  type
                  index
                  stringValue
                  intValue
                  booleanValue
                  floatValue
                  latitudeValue
                  longitudeValue
                  rangeFromValue
                  rangeToValue
                  isEditable
                  isInstanceProperty
                  isMandatory
                }
                id
              }
            }
            parentEquipment {
              id
              name
              futureState
              equipmentType {
                id
                name
                portDefinitions {
                  id
                  name
                  visibleLabel
                  bandwidth
                  portType {
                    id
                    name
                  }
                }
              }
              ...EquipmentBreadcrumbs_equipment
            }
          }
          workOrder {
            id
            status
          }
          properties {
            id
            propertyType {
              id
              name
              type
              isEditable
              isMandatory
              isInstanceProperty
              stringValue
            }
            stringValue
            intValue
            floatValue
            booleanValue
            latitudeValue
            longitudeValue
            rangeFromValue
            rangeToValue
            equipmentValue {
              id
              name
            }
            locationValue {
              id
              name
            }
            serviceValue {
              id
              name
            }
          }
          services {
            id
            name
          }
        }
        properties {
          id
          propertyType {
            id
            name
            type
            isEditable
            isMandatory
            isInstanceProperty
            stringValue
          }
          stringValue
          intValue
          floatValue
          booleanValue
          latitudeValue
          longitudeValue
          rangeFromValue
          rangeToValue
          equipmentValue {
            id
            name
          }
          locationValue {
            id
            name
          }
          serviceValue {
            id
            name
          }
        }
        serviceEndpoints {
          role
          service {
            name
            id
          }
          id
        }
      }
    }
    id
  }
}

fragment EquipmentBreadcrumbs_equipment on Equipment {
  id
  name
  equipmentType {
    id
    name
  }
  locationHierarchy {
    id
    name
    locationType {
      name
      id
    }
  }
  positionHierarchy {
    id
    definition {
      id
      name
      visibleLabel
    }
    parentEquipment {
      id
      name
      equipmentType {
        id
        name
      }
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "equipmentId",
    "type": "ID!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "Variable",
    "name": "id",
    "variableName": "equipmentId"
  }
],
v2 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "id",
  "args": null,
  "storageKey": null
},
v3 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "name",
  "args": null,
  "storageKey": null
},
v4 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "visibleLabel",
  "args": null,
  "storageKey": null
},
v5 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "bandwidth",
  "args": null,
  "storageKey": null
},
v6 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "portDefinitions",
  "storageKey": null,
  "args": null,
  "concreteType": "EquipmentPortDefinition",
  "plural": true,
  "selections": [
    (v2/*: any*/),
    (v3/*: any*/),
    (v4/*: any*/),
    (v5/*: any*/)
  ]
},
v7 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "equipmentType",
  "storageKey": null,
  "args": null,
  "concreteType": "EquipmentType",
  "plural": false,
  "selections": [
    (v2/*: any*/),
    (v3/*: any*/),
    (v6/*: any*/)
  ]
},
v8 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "index",
  "args": null,
  "storageKey": null
},
v9 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "type",
  "args": null,
  "storageKey": null
},
v10 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "stringValue",
  "args": null,
  "storageKey": null
},
v11 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "intValue",
  "args": null,
  "storageKey": null
},
v12 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "booleanValue",
  "args": null,
  "storageKey": null
},
v13 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "floatValue",
  "args": null,
  "storageKey": null
},
v14 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "latitudeValue",
  "args": null,
  "storageKey": null
},
v15 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "longitudeValue",
  "args": null,
  "storageKey": null
},
v16 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeFromValue",
  "args": null,
  "storageKey": null
},
v17 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeToValue",
  "args": null,
  "storageKey": null
},
v18 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isEditable",
  "args": null,
  "storageKey": null
},
v19 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isInstanceProperty",
  "args": null,
  "storageKey": null
},
v20 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isMandatory",
  "args": null,
  "storageKey": null
},
v21 = [
  (v2/*: any*/),
  (v3/*: any*/),
  (v9/*: any*/),
  (v8/*: any*/),
  (v10/*: any*/),
  (v11/*: any*/),
  (v12/*: any*/),
  (v13/*: any*/),
  (v14/*: any*/),
  (v15/*: any*/),
  (v16/*: any*/),
  (v17/*: any*/),
  (v18/*: any*/),
  (v19/*: any*/),
  (v20/*: any*/)
],
v22 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "linkPropertyTypes",
  "storageKey": null,
  "args": null,
  "concreteType": "PropertyType",
  "plural": true,
  "selections": (v21/*: any*/)
},
v23 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "definition",
  "storageKey": null,
  "args": null,
  "concreteType": "EquipmentPortDefinition",
  "plural": false,
  "selections": [
    (v2/*: any*/),
    (v3/*: any*/),
    (v8/*: any*/),
    (v4/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "portType",
      "storageKey": null,
      "args": null,
      "concreteType": "EquipmentPortType",
      "plural": false,
      "selections": [
        (v2/*: any*/),
        (v3/*: any*/),
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "propertyTypes",
          "storageKey": null,
          "args": null,
          "concreteType": "PropertyType",
          "plural": true,
          "selections": (v21/*: any*/)
        },
        (v22/*: any*/)
      ]
    }
  ]
},
v24 = [
  (v2/*: any*/),
  (v3/*: any*/)
],
v25 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "portType",
  "storageKey": null,
  "args": null,
  "concreteType": "EquipmentPortType",
  "plural": false,
  "selections": (v24/*: any*/)
},
v26 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "parentEquipment",
  "storageKey": null,
  "args": null,
  "concreteType": "Equipment",
  "plural": false,
  "selections": [
    (v2/*: any*/),
    (v3/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "equipmentType",
      "storageKey": null,
      "args": null,
      "concreteType": "EquipmentType",
      "plural": false,
      "selections": [
        (v2/*: any*/),
        (v3/*: any*/),
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "portDefinitions",
          "storageKey": null,
          "args": null,
          "concreteType": "EquipmentPortDefinition",
          "plural": true,
          "selections": [
            (v2/*: any*/),
            (v3/*: any*/),
            (v4/*: any*/),
            (v25/*: any*/),
            (v5/*: any*/)
          ]
        }
      ]
    }
  ]
},
v27 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "futureState",
  "args": null,
  "storageKey": null
},
v28 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "equipmentType",
  "storageKey": null,
  "args": null,
  "concreteType": "EquipmentType",
  "plural": false,
  "selections": [
    (v2/*: any*/),
    (v3/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "portDefinitions",
      "storageKey": null,
      "args": null,
      "concreteType": "EquipmentPortDefinition",
      "plural": true,
      "selections": [
        (v2/*: any*/),
        (v3/*: any*/),
        (v4/*: any*/),
        (v5/*: any*/),
        (v25/*: any*/)
      ]
    }
  ]
},
v29 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "workOrder",
  "storageKey": null,
  "args": null,
  "concreteType": "WorkOrder",
  "plural": false,
  "selections": [
    (v2/*: any*/),
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "status",
      "args": null,
      "storageKey": null
    }
  ]
},
v30 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "properties",
  "storageKey": null,
  "args": null,
  "concreteType": "Property",
  "plural": true,
  "selections": [
    (v2/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "propertyType",
      "storageKey": null,
      "args": null,
      "concreteType": "PropertyType",
      "plural": false,
      "selections": [
        (v2/*: any*/),
        (v3/*: any*/),
        (v9/*: any*/),
        (v18/*: any*/),
        (v20/*: any*/),
        (v19/*: any*/),
        (v10/*: any*/)
      ]
    },
    (v10/*: any*/),
    (v11/*: any*/),
    (v13/*: any*/),
    (v12/*: any*/),
    (v14/*: any*/),
    (v15/*: any*/),
    (v16/*: any*/),
    (v17/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "equipmentValue",
      "storageKey": null,
      "args": null,
      "concreteType": "Equipment",
      "plural": false,
      "selections": (v24/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "locationValue",
      "storageKey": null,
      "args": null,
      "concreteType": "Location",
      "plural": false,
      "selections": (v24/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "serviceValue",
      "storageKey": null,
      "args": null,
      "concreteType": "Service",
      "plural": false,
      "selections": (v24/*: any*/)
    }
  ]
},
v31 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "services",
  "storageKey": null,
  "args": null,
  "concreteType": "Service",
  "plural": true,
  "selections": (v24/*: any*/)
},
v32 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "role",
  "args": null,
  "storageKey": null
},
v33 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "ports",
  "storageKey": null,
  "args": null,
  "concreteType": "EquipmentPort",
  "plural": true,
  "selections": [
    (v2/*: any*/),
    (v23/*: any*/),
    (v26/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "link",
      "storageKey": null,
      "args": null,
      "concreteType": "Link",
      "plural": false,
      "selections": [
        (v2/*: any*/),
        (v27/*: any*/),
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "ports",
          "storageKey": null,
          "args": null,
          "concreteType": "EquipmentPort",
          "plural": true,
          "selections": [
            (v2/*: any*/),
            {
              "kind": "LinkedField",
              "alias": null,
              "name": "definition",
              "storageKey": null,
              "args": null,
              "concreteType": "EquipmentPortDefinition",
              "plural": false,
              "selections": [
                (v2/*: any*/),
                (v3/*: any*/),
                (v4/*: any*/),
                {
                  "kind": "LinkedField",
                  "alias": null,
                  "name": "portType",
                  "storageKey": null,
                  "args": null,
                  "concreteType": "EquipmentPortType",
                  "plural": false,
                  "selections": [
                    (v22/*: any*/)
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
                (v2/*: any*/),
                (v3/*: any*/),
                (v27/*: any*/),
                (v28/*: any*/),
                {
                  "kind": "FragmentSpread",
                  "name": "EquipmentBreadcrumbs_equipment",
                  "args": null
                }
              ]
            }
          ]
        },
        (v29/*: any*/),
        (v30/*: any*/),
        (v31/*: any*/)
      ]
    },
    (v30/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "serviceEndpoints",
      "storageKey": null,
      "args": null,
      "concreteType": "ServiceEndpoint",
      "plural": true,
      "selections": [
        (v32/*: any*/),
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "service",
          "storageKey": null,
          "args": null,
          "concreteType": "Service",
          "plural": false,
          "selections": [
            (v3/*: any*/)
          ]
        }
      ]
    }
  ]
},
v34 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "equipmentType",
  "storageKey": null,
  "args": null,
  "concreteType": "EquipmentType",
  "plural": false,
  "selections": [
    (v6/*: any*/)
  ]
},
v35 = [
  (v3/*: any*/),
  (v2/*: any*/)
],
v36 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "ports",
  "storageKey": null,
  "args": null,
  "concreteType": "EquipmentPort",
  "plural": true,
  "selections": [
    (v2/*: any*/),
    (v23/*: any*/),
    (v26/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "link",
      "storageKey": null,
      "args": null,
      "concreteType": "Link",
      "plural": false,
      "selections": [
        (v2/*: any*/),
        (v27/*: any*/),
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "ports",
          "storageKey": null,
          "args": null,
          "concreteType": "EquipmentPort",
          "plural": true,
          "selections": [
            (v2/*: any*/),
            {
              "kind": "LinkedField",
              "alias": null,
              "name": "definition",
              "storageKey": null,
              "args": null,
              "concreteType": "EquipmentPortDefinition",
              "plural": false,
              "selections": [
                (v2/*: any*/),
                (v3/*: any*/),
                (v4/*: any*/),
                {
                  "kind": "LinkedField",
                  "alias": null,
                  "name": "portType",
                  "storageKey": null,
                  "args": null,
                  "concreteType": "EquipmentPortType",
                  "plural": false,
                  "selections": [
                    (v22/*: any*/),
                    (v2/*: any*/)
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
                (v2/*: any*/),
                (v3/*: any*/),
                (v27/*: any*/),
                (v28/*: any*/),
                {
                  "kind": "LinkedField",
                  "alias": null,
                  "name": "locationHierarchy",
                  "storageKey": null,
                  "args": null,
                  "concreteType": "Location",
                  "plural": true,
                  "selections": [
                    (v2/*: any*/),
                    (v3/*: any*/),
                    {
                      "kind": "LinkedField",
                      "alias": null,
                      "name": "locationType",
                      "storageKey": null,
                      "args": null,
                      "concreteType": "LocationType",
                      "plural": false,
                      "selections": (v35/*: any*/)
                    }
                  ]
                },
                {
                  "kind": "LinkedField",
                  "alias": null,
                  "name": "positionHierarchy",
                  "storageKey": null,
                  "args": null,
                  "concreteType": "EquipmentPosition",
                  "plural": true,
                  "selections": [
                    (v2/*: any*/),
                    {
                      "kind": "LinkedField",
                      "alias": null,
                      "name": "definition",
                      "storageKey": null,
                      "args": null,
                      "concreteType": "EquipmentPositionDefinition",
                      "plural": false,
                      "selections": [
                        (v2/*: any*/),
                        (v3/*: any*/),
                        (v4/*: any*/)
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
                        (v2/*: any*/),
                        (v3/*: any*/),
                        {
                          "kind": "LinkedField",
                          "alias": null,
                          "name": "equipmentType",
                          "storageKey": null,
                          "args": null,
                          "concreteType": "EquipmentType",
                          "plural": false,
                          "selections": (v24/*: any*/)
                        }
                      ]
                    }
                  ]
                }
              ]
            }
          ]
        },
        (v29/*: any*/),
        (v30/*: any*/),
        (v31/*: any*/)
      ]
    },
    (v30/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "serviceEndpoints",
      "storageKey": null,
      "args": null,
      "concreteType": "ServiceEndpoint",
      "plural": true,
      "selections": [
        (v32/*: any*/),
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "service",
          "storageKey": null,
          "args": null,
          "concreteType": "Service",
          "plural": false,
          "selections": (v35/*: any*/)
        },
        (v2/*: any*/)
      ]
    }
  ]
},
v37 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "equipmentType",
  "storageKey": null,
  "args": null,
  "concreteType": "EquipmentType",
  "plural": false,
  "selections": [
    (v6/*: any*/),
    (v2/*: any*/)
  ]
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "PortsConnectDialogQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "equipment",
        "name": "node",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": null,
        "plural": false,
        "selections": [
          {
            "kind": "InlineFragment",
            "type": "Equipment",
            "selections": [
              (v2/*: any*/),
              (v3/*: any*/),
              (v7/*: any*/),
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
                      (v2/*: any*/),
                      (v3/*: any*/),
                      (v33/*: any*/),
                      (v34/*: any*/),
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
                              (v2/*: any*/),
                              (v3/*: any*/),
                              (v33/*: any*/),
                              (v34/*: any*/),
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
                                      (v2/*: any*/),
                                      (v3/*: any*/),
                                      (v33/*: any*/),
                                      (v34/*: any*/),
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
                                              (v2/*: any*/),
                                              (v3/*: any*/),
                                              (v33/*: any*/),
                                              (v34/*: any*/)
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
              },
              (v33/*: any*/)
            ]
          }
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "PortsConnectDialogQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "equipment",
        "name": "node",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": null,
        "plural": false,
        "selections": [
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "__typename",
            "args": null,
            "storageKey": null
          },
          (v2/*: any*/),
          {
            "kind": "InlineFragment",
            "type": "Equipment",
            "selections": [
              (v3/*: any*/),
              (v7/*: any*/),
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
                      (v2/*: any*/),
                      (v3/*: any*/),
                      (v36/*: any*/),
                      (v37/*: any*/),
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
                              (v2/*: any*/),
                              (v3/*: any*/),
                              (v36/*: any*/),
                              (v37/*: any*/),
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
                                      (v2/*: any*/),
                                      (v3/*: any*/),
                                      (v36/*: any*/),
                                      (v37/*: any*/),
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
                                              (v2/*: any*/),
                                              (v3/*: any*/),
                                              (v36/*: any*/),
                                              (v37/*: any*/)
                                            ]
                                          },
                                          (v2/*: any*/)
                                        ]
                                      }
                                    ]
                                  },
                                  (v2/*: any*/)
                                ]
                              }
                            ]
                          },
                          (v2/*: any*/)
                        ]
                      }
                    ]
                  },
                  (v2/*: any*/)
                ]
              },
              (v36/*: any*/)
            ]
          }
        ]
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "PortsConnectDialogQuery",
    "id": null,
    "text": "query PortsConnectDialogQuery(\n  $equipmentId: ID!\n) {\n  equipment: node(id: $equipmentId) {\n    __typename\n    ... on Equipment {\n      id\n      name\n      equipmentType {\n        id\n        name\n        portDefinitions {\n          id\n          name\n          visibleLabel\n          bandwidth\n        }\n      }\n      positions {\n        attachedEquipment {\n          id\n          name\n          ports {\n            id\n            definition {\n              id\n              name\n              index\n              visibleLabel\n              portType {\n                id\n                name\n                propertyTypes {\n                  id\n                  name\n                  type\n                  index\n                  stringValue\n                  intValue\n                  booleanValue\n                  floatValue\n                  latitudeValue\n                  longitudeValue\n                  rangeFromValue\n                  rangeToValue\n                  isEditable\n                  isInstanceProperty\n                  isMandatory\n                }\n                linkPropertyTypes {\n                  id\n                  name\n                  type\n                  index\n                  stringValue\n                  intValue\n                  booleanValue\n                  floatValue\n                  latitudeValue\n                  longitudeValue\n                  rangeFromValue\n                  rangeToValue\n                  isEditable\n                  isInstanceProperty\n                  isMandatory\n                }\n              }\n            }\n            parentEquipment {\n              id\n              name\n              equipmentType {\n                id\n                name\n                portDefinitions {\n                  id\n                  name\n                  visibleLabel\n                  portType {\n                    id\n                    name\n                  }\n                  bandwidth\n                }\n              }\n            }\n            link {\n              id\n              futureState\n              ports {\n                id\n                definition {\n                  id\n                  name\n                  visibleLabel\n                  portType {\n                    linkPropertyTypes {\n                      id\n                      name\n                      type\n                      index\n                      stringValue\n                      intValue\n                      booleanValue\n                      floatValue\n                      latitudeValue\n                      longitudeValue\n                      rangeFromValue\n                      rangeToValue\n                      isEditable\n                      isInstanceProperty\n                      isMandatory\n                    }\n                    id\n                  }\n                }\n                parentEquipment {\n                  id\n                  name\n                  futureState\n                  equipmentType {\n                    id\n                    name\n                    portDefinitions {\n                      id\n                      name\n                      visibleLabel\n                      bandwidth\n                      portType {\n                        id\n                        name\n                      }\n                    }\n                  }\n                  ...EquipmentBreadcrumbs_equipment\n                }\n              }\n              workOrder {\n                id\n                status\n              }\n              properties {\n                id\n                propertyType {\n                  id\n                  name\n                  type\n                  isEditable\n                  isMandatory\n                  isInstanceProperty\n                  stringValue\n                }\n                stringValue\n                intValue\n                floatValue\n                booleanValue\n                latitudeValue\n                longitudeValue\n                rangeFromValue\n                rangeToValue\n                equipmentValue {\n                  id\n                  name\n                }\n                locationValue {\n                  id\n                  name\n                }\n                serviceValue {\n                  id\n                  name\n                }\n              }\n              services {\n                id\n                name\n              }\n            }\n            properties {\n              id\n              propertyType {\n                id\n                name\n                type\n                isEditable\n                isMandatory\n                isInstanceProperty\n                stringValue\n              }\n              stringValue\n              intValue\n              floatValue\n              booleanValue\n              latitudeValue\n              longitudeValue\n              rangeFromValue\n              rangeToValue\n              equipmentValue {\n                id\n                name\n              }\n              locationValue {\n                id\n                name\n              }\n              serviceValue {\n                id\n                name\n              }\n            }\n            serviceEndpoints {\n              role\n              service {\n                name\n                id\n              }\n              id\n            }\n          }\n          equipmentType {\n            portDefinitions {\n              id\n              name\n              visibleLabel\n              bandwidth\n            }\n            id\n          }\n          positions {\n            attachedEquipment {\n              id\n              name\n              ports {\n                id\n                definition {\n                  id\n                  name\n                  index\n                  visibleLabel\n                  portType {\n                    id\n                    name\n                    propertyTypes {\n                      id\n                      name\n                      type\n                      index\n                      stringValue\n                      intValue\n                      booleanValue\n                      floatValue\n                      latitudeValue\n                      longitudeValue\n                      rangeFromValue\n                      rangeToValue\n                      isEditable\n                      isInstanceProperty\n                      isMandatory\n                    }\n                    linkPropertyTypes {\n                      id\n                      name\n                      type\n                      index\n                      stringValue\n                      intValue\n                      booleanValue\n                      floatValue\n                      latitudeValue\n                      longitudeValue\n                      rangeFromValue\n                      rangeToValue\n                      isEditable\n                      isInstanceProperty\n                      isMandatory\n                    }\n                  }\n                }\n                parentEquipment {\n                  id\n                  name\n                  equipmentType {\n                    id\n                    name\n                    portDefinitions {\n                      id\n                      name\n                      visibleLabel\n                      portType {\n                        id\n                        name\n                      }\n                      bandwidth\n                    }\n                  }\n                }\n                link {\n                  id\n                  futureState\n                  ports {\n                    id\n                    definition {\n                      id\n                      name\n                      visibleLabel\n                      portType {\n                        linkPropertyTypes {\n                          id\n                          name\n                          type\n                          index\n                          stringValue\n                          intValue\n                          booleanValue\n                          floatValue\n                          latitudeValue\n                          longitudeValue\n                          rangeFromValue\n                          rangeToValue\n                          isEditable\n                          isInstanceProperty\n                          isMandatory\n                        }\n                        id\n                      }\n                    }\n                    parentEquipment {\n                      id\n                      name\n                      futureState\n                      equipmentType {\n                        id\n                        name\n                        portDefinitions {\n                          id\n                          name\n                          visibleLabel\n                          bandwidth\n                          portType {\n                            id\n                            name\n                          }\n                        }\n                      }\n                      ...EquipmentBreadcrumbs_equipment\n                    }\n                  }\n                  workOrder {\n                    id\n                    status\n                  }\n                  properties {\n                    id\n                    propertyType {\n                      id\n                      name\n                      type\n                      isEditable\n                      isMandatory\n                      isInstanceProperty\n                      stringValue\n                    }\n                    stringValue\n                    intValue\n                    floatValue\n                    booleanValue\n                    latitudeValue\n                    longitudeValue\n                    rangeFromValue\n                    rangeToValue\n                    equipmentValue {\n                      id\n                      name\n                    }\n                    locationValue {\n                      id\n                      name\n                    }\n                    serviceValue {\n                      id\n                      name\n                    }\n                  }\n                  services {\n                    id\n                    name\n                  }\n                }\n                properties {\n                  id\n                  propertyType {\n                    id\n                    name\n                    type\n                    isEditable\n                    isMandatory\n                    isInstanceProperty\n                    stringValue\n                  }\n                  stringValue\n                  intValue\n                  floatValue\n                  booleanValue\n                  latitudeValue\n                  longitudeValue\n                  rangeFromValue\n                  rangeToValue\n                  equipmentValue {\n                    id\n                    name\n                  }\n                  locationValue {\n                    id\n                    name\n                  }\n                  serviceValue {\n                    id\n                    name\n                  }\n                }\n                serviceEndpoints {\n                  role\n                  service {\n                    name\n                    id\n                  }\n                  id\n                }\n              }\n              equipmentType {\n                portDefinitions {\n                  id\n                  name\n                  visibleLabel\n                  bandwidth\n                }\n                id\n              }\n              positions {\n                attachedEquipment {\n                  id\n                  name\n                  ports {\n                    id\n                    definition {\n                      id\n                      name\n                      index\n                      visibleLabel\n                      portType {\n                        id\n                        name\n                        propertyTypes {\n                          id\n                          name\n                          type\n                          index\n                          stringValue\n                          intValue\n                          booleanValue\n                          floatValue\n                          latitudeValue\n                          longitudeValue\n                          rangeFromValue\n                          rangeToValue\n                          isEditable\n                          isInstanceProperty\n                          isMandatory\n                        }\n                        linkPropertyTypes {\n                          id\n                          name\n                          type\n                          index\n                          stringValue\n                          intValue\n                          booleanValue\n                          floatValue\n                          latitudeValue\n                          longitudeValue\n                          rangeFromValue\n                          rangeToValue\n                          isEditable\n                          isInstanceProperty\n                          isMandatory\n                        }\n                      }\n                    }\n                    parentEquipment {\n                      id\n                      name\n                      equipmentType {\n                        id\n                        name\n                        portDefinitions {\n                          id\n                          name\n                          visibleLabel\n                          portType {\n                            id\n                            name\n                          }\n                          bandwidth\n                        }\n                      }\n                    }\n                    link {\n                      id\n                      futureState\n                      ports {\n                        id\n                        definition {\n                          id\n                          name\n                          visibleLabel\n                          portType {\n                            linkPropertyTypes {\n                              id\n                              name\n                              type\n                              index\n                              stringValue\n                              intValue\n                              booleanValue\n                              floatValue\n                              latitudeValue\n                              longitudeValue\n                              rangeFromValue\n                              rangeToValue\n                              isEditable\n                              isInstanceProperty\n                              isMandatory\n                            }\n                            id\n                          }\n                        }\n                        parentEquipment {\n                          id\n                          name\n                          futureState\n                          equipmentType {\n                            id\n                            name\n                            portDefinitions {\n                              id\n                              name\n                              visibleLabel\n                              bandwidth\n                              portType {\n                                id\n                                name\n                              }\n                            }\n                          }\n                          ...EquipmentBreadcrumbs_equipment\n                        }\n                      }\n                      workOrder {\n                        id\n                        status\n                      }\n                      properties {\n                        id\n                        propertyType {\n                          id\n                          name\n                          type\n                          isEditable\n                          isMandatory\n                          isInstanceProperty\n                          stringValue\n                        }\n                        stringValue\n                        intValue\n                        floatValue\n                        booleanValue\n                        latitudeValue\n                        longitudeValue\n                        rangeFromValue\n                        rangeToValue\n                        equipmentValue {\n                          id\n                          name\n                        }\n                        locationValue {\n                          id\n                          name\n                        }\n                        serviceValue {\n                          id\n                          name\n                        }\n                      }\n                      services {\n                        id\n                        name\n                      }\n                    }\n                    properties {\n                      id\n                      propertyType {\n                        id\n                        name\n                        type\n                        isEditable\n                        isMandatory\n                        isInstanceProperty\n                        stringValue\n                      }\n                      stringValue\n                      intValue\n                      floatValue\n                      booleanValue\n                      latitudeValue\n                      longitudeValue\n                      rangeFromValue\n                      rangeToValue\n                      equipmentValue {\n                        id\n                        name\n                      }\n                      locationValue {\n                        id\n                        name\n                      }\n                      serviceValue {\n                        id\n                        name\n                      }\n                    }\n                    serviceEndpoints {\n                      role\n                      service {\n                        name\n                        id\n                      }\n                      id\n                    }\n                  }\n                  equipmentType {\n                    portDefinitions {\n                      id\n                      name\n                      visibleLabel\n                      bandwidth\n                    }\n                    id\n                  }\n                  positions {\n                    attachedEquipment {\n                      id\n                      name\n                      ports {\n                        id\n                        definition {\n                          id\n                          name\n                          index\n                          visibleLabel\n                          portType {\n                            id\n                            name\n                            propertyTypes {\n                              id\n                              name\n                              type\n                              index\n                              stringValue\n                              intValue\n                              booleanValue\n                              floatValue\n                              latitudeValue\n                              longitudeValue\n                              rangeFromValue\n                              rangeToValue\n                              isEditable\n                              isInstanceProperty\n                              isMandatory\n                            }\n                            linkPropertyTypes {\n                              id\n                              name\n                              type\n                              index\n                              stringValue\n                              intValue\n                              booleanValue\n                              floatValue\n                              latitudeValue\n                              longitudeValue\n                              rangeFromValue\n                              rangeToValue\n                              isEditable\n                              isInstanceProperty\n                              isMandatory\n                            }\n                          }\n                        }\n                        parentEquipment {\n                          id\n                          name\n                          equipmentType {\n                            id\n                            name\n                            portDefinitions {\n                              id\n                              name\n                              visibleLabel\n                              portType {\n                                id\n                                name\n                              }\n                              bandwidth\n                            }\n                          }\n                        }\n                        link {\n                          id\n                          futureState\n                          ports {\n                            id\n                            definition {\n                              id\n                              name\n                              visibleLabel\n                              portType {\n                                linkPropertyTypes {\n                                  id\n                                  name\n                                  type\n                                  index\n                                  stringValue\n                                  intValue\n                                  booleanValue\n                                  floatValue\n                                  latitudeValue\n                                  longitudeValue\n                                  rangeFromValue\n                                  rangeToValue\n                                  isEditable\n                                  isInstanceProperty\n                                  isMandatory\n                                }\n                                id\n                              }\n                            }\n                            parentEquipment {\n                              id\n                              name\n                              futureState\n                              equipmentType {\n                                id\n                                name\n                                portDefinitions {\n                                  id\n                                  name\n                                  visibleLabel\n                                  bandwidth\n                                  portType {\n                                    id\n                                    name\n                                  }\n                                }\n                              }\n                              ...EquipmentBreadcrumbs_equipment\n                            }\n                          }\n                          workOrder {\n                            id\n                            status\n                          }\n                          properties {\n                            id\n                            propertyType {\n                              id\n                              name\n                              type\n                              isEditable\n                              isMandatory\n                              isInstanceProperty\n                              stringValue\n                            }\n                            stringValue\n                            intValue\n                            floatValue\n                            booleanValue\n                            latitudeValue\n                            longitudeValue\n                            rangeFromValue\n                            rangeToValue\n                            equipmentValue {\n                              id\n                              name\n                            }\n                            locationValue {\n                              id\n                              name\n                            }\n                            serviceValue {\n                              id\n                              name\n                            }\n                          }\n                          services {\n                            id\n                            name\n                          }\n                        }\n                        properties {\n                          id\n                          propertyType {\n                            id\n                            name\n                            type\n                            isEditable\n                            isMandatory\n                            isInstanceProperty\n                            stringValue\n                          }\n                          stringValue\n                          intValue\n                          floatValue\n                          booleanValue\n                          latitudeValue\n                          longitudeValue\n                          rangeFromValue\n                          rangeToValue\n                          equipmentValue {\n                            id\n                            name\n                          }\n                          locationValue {\n                            id\n                            name\n                          }\n                          serviceValue {\n                            id\n                            name\n                          }\n                        }\n                        serviceEndpoints {\n                          role\n                          service {\n                            name\n                            id\n                          }\n                          id\n                        }\n                      }\n                      equipmentType {\n                        portDefinitions {\n                          id\n                          name\n                          visibleLabel\n                          bandwidth\n                        }\n                        id\n                      }\n                    }\n                    id\n                  }\n                }\n                id\n              }\n            }\n            id\n          }\n        }\n        id\n      }\n      ports {\n        id\n        definition {\n          id\n          name\n          index\n          visibleLabel\n          portType {\n            id\n            name\n            propertyTypes {\n              id\n              name\n              type\n              index\n              stringValue\n              intValue\n              booleanValue\n              floatValue\n              latitudeValue\n              longitudeValue\n              rangeFromValue\n              rangeToValue\n              isEditable\n              isInstanceProperty\n              isMandatory\n            }\n            linkPropertyTypes {\n              id\n              name\n              type\n              index\n              stringValue\n              intValue\n              booleanValue\n              floatValue\n              latitudeValue\n              longitudeValue\n              rangeFromValue\n              rangeToValue\n              isEditable\n              isInstanceProperty\n              isMandatory\n            }\n          }\n        }\n        parentEquipment {\n          id\n          name\n          equipmentType {\n            id\n            name\n            portDefinitions {\n              id\n              name\n              visibleLabel\n              portType {\n                id\n                name\n              }\n              bandwidth\n            }\n          }\n        }\n        link {\n          id\n          futureState\n          ports {\n            id\n            definition {\n              id\n              name\n              visibleLabel\n              portType {\n                linkPropertyTypes {\n                  id\n                  name\n                  type\n                  index\n                  stringValue\n                  intValue\n                  booleanValue\n                  floatValue\n                  latitudeValue\n                  longitudeValue\n                  rangeFromValue\n                  rangeToValue\n                  isEditable\n                  isInstanceProperty\n                  isMandatory\n                }\n                id\n              }\n            }\n            parentEquipment {\n              id\n              name\n              futureState\n              equipmentType {\n                id\n                name\n                portDefinitions {\n                  id\n                  name\n                  visibleLabel\n                  bandwidth\n                  portType {\n                    id\n                    name\n                  }\n                }\n              }\n              ...EquipmentBreadcrumbs_equipment\n            }\n          }\n          workOrder {\n            id\n            status\n          }\n          properties {\n            id\n            propertyType {\n              id\n              name\n              type\n              isEditable\n              isMandatory\n              isInstanceProperty\n              stringValue\n            }\n            stringValue\n            intValue\n            floatValue\n            booleanValue\n            latitudeValue\n            longitudeValue\n            rangeFromValue\n            rangeToValue\n            equipmentValue {\n              id\n              name\n            }\n            locationValue {\n              id\n              name\n            }\n            serviceValue {\n              id\n              name\n            }\n          }\n          services {\n            id\n            name\n          }\n        }\n        properties {\n          id\n          propertyType {\n            id\n            name\n            type\n            isEditable\n            isMandatory\n            isInstanceProperty\n            stringValue\n          }\n          stringValue\n          intValue\n          floatValue\n          booleanValue\n          latitudeValue\n          longitudeValue\n          rangeFromValue\n          rangeToValue\n          equipmentValue {\n            id\n            name\n          }\n          locationValue {\n            id\n            name\n          }\n          serviceValue {\n            id\n            name\n          }\n        }\n        serviceEndpoints {\n          role\n          service {\n            name\n            id\n          }\n          id\n        }\n      }\n    }\n    id\n  }\n}\n\nfragment EquipmentBreadcrumbs_equipment on Equipment {\n  id\n  name\n  equipmentType {\n    id\n    name\n  }\n  locationHierarchy {\n    id\n    name\n    locationType {\n      name\n      id\n    }\n  }\n  positionHierarchy {\n    id\n    definition {\n      id\n      name\n      visibleLabel\n    }\n    parentEquipment {\n      id\n      name\n      equipmentType {\n        id\n        name\n      }\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '9672bba6848ea6ee72e822943c65226a';
module.exports = node;
