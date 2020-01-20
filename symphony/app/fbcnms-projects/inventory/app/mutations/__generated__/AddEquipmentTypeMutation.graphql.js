/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash a83c6b90e175a93356f5144f563f9cfc
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type EquipmentTypeItem_equipmentType$ref = any;
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "equipment" | "float" | "gps_location" | "int" | "location" | "range" | "service" | "string" | "%future added value";
export type AddEquipmentTypeInput = {|
  name: string,
  category?: ?string,
  positions?: ?$ReadOnlyArray<EquipmentPositionInput>,
  ports?: ?$ReadOnlyArray<EquipmentPortInput>,
  properties?: ?$ReadOnlyArray<PropertyTypeInput>,
|};
export type EquipmentPositionInput = {|
  id?: ?string,
  name: string,
  index?: ?number,
  visibleLabel?: ?string,
|};
export type EquipmentPortInput = {|
  id?: ?string,
  name: string,
  index?: ?number,
  visibleLabel?: ?string,
  portTypeID?: ?string,
  bandwidth?: ?string,
|};
export type PropertyTypeInput = {|
  id?: ?string,
  name: string,
  type: PropertyKind,
  index?: ?number,
  category?: ?string,
  stringValue?: ?string,
  intValue?: ?number,
  booleanValue?: ?boolean,
  floatValue?: ?number,
  latitudeValue?: ?number,
  longitudeValue?: ?number,
  rangeFromValue?: ?number,
  rangeToValue?: ?number,
  isEditable?: ?boolean,
  isInstanceProperty?: ?boolean,
  isMandatory?: ?boolean,
  isDeleted?: ?boolean,
|};
export type AddEquipmentTypeMutationVariables = {|
  input: AddEquipmentTypeInput
|};
export type AddEquipmentTypeMutationResponse = {|
  +addEquipmentType: ?{|
    +id: string,
    +name: string,
    +$fragmentRefs: EquipmentTypeItem_equipmentType$ref,
  |}
|};
export type AddEquipmentTypeMutation = {|
  variables: AddEquipmentTypeMutationVariables,
  response: AddEquipmentTypeMutationResponse,
|};
*/


/*
mutation AddEquipmentTypeMutation(
  $input: AddEquipmentTypeInput!
) {
  addEquipmentType(input: $input) {
    id
    name
    ...EquipmentTypeItem_equipmentType
  }
}

fragment DynamicPropertyTypesGrid_propertyTypes on PropertyType {
  ...PropertyTypeFormField_propertyType
  id
  index
}

fragment EquipmentTypeItem_equipmentType on EquipmentType {
  id
  name
  propertyTypes {
    ...DynamicPropertyTypesGrid_propertyTypes
    id
  }
  positionDefinitions {
    ...PositionDefinitionsTable_positionDefinitions
    id
  }
  portDefinitions {
    ...PortDefinitionsTable_portDefinitions
    id
  }
  numberOfEquipment
}

fragment PortDefinitionsTable_portDefinitions on EquipmentPortDefinition {
  id
  name
  index
  visibleLabel
  portType {
    id
    name
  }
}

fragment PositionDefinitionsTable_positionDefinitions on EquipmentPositionDefinition {
  id
  name
  index
  visibleLabel
}

fragment PropertyTypeFormField_propertyType on PropertyType {
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
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "input",
    "type": "AddEquipmentTypeInput!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "Variable",
    "name": "input",
    "variableName": "input"
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
  "name": "index",
  "args": null,
  "storageKey": null
},
v5 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "visibleLabel",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "AddEquipmentTypeMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "addEquipmentType",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "EquipmentType",
        "plural": false,
        "selections": [
          (v2/*: any*/),
          (v3/*: any*/),
          {
            "kind": "FragmentSpread",
            "name": "EquipmentTypeItem_equipmentType",
            "args": null
          }
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "AddEquipmentTypeMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "addEquipmentType",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "EquipmentType",
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
            "selections": [
              (v2/*: any*/),
              (v3/*: any*/),
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "type",
                "args": null,
                "storageKey": null
              },
              (v4/*: any*/),
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "stringValue",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "intValue",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "booleanValue",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "floatValue",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "latitudeValue",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "longitudeValue",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "rangeFromValue",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "rangeToValue",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "isEditable",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "isInstanceProperty",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "isMandatory",
                "args": null,
                "storageKey": null
              }
            ]
          },
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "positionDefinitions",
            "storageKey": null,
            "args": null,
            "concreteType": "EquipmentPositionDefinition",
            "plural": true,
            "selections": [
              (v2/*: any*/),
              (v3/*: any*/),
              (v4/*: any*/),
              (v5/*: any*/)
            ]
          },
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
                  (v3/*: any*/)
                ]
              }
            ]
          },
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "numberOfEquipment",
            "args": null,
            "storageKey": null
          }
        ]
      }
    ]
  },
  "params": {
    "operationKind": "mutation",
    "name": "AddEquipmentTypeMutation",
    "id": null,
    "text": "mutation AddEquipmentTypeMutation(\n  $input: AddEquipmentTypeInput!\n) {\n  addEquipmentType(input: $input) {\n    id\n    name\n    ...EquipmentTypeItem_equipmentType\n  }\n}\n\nfragment DynamicPropertyTypesGrid_propertyTypes on PropertyType {\n  ...PropertyTypeFormField_propertyType\n  id\n  index\n}\n\nfragment EquipmentTypeItem_equipmentType on EquipmentType {\n  id\n  name\n  propertyTypes {\n    ...DynamicPropertyTypesGrid_propertyTypes\n    id\n  }\n  positionDefinitions {\n    ...PositionDefinitionsTable_positionDefinitions\n    id\n  }\n  portDefinitions {\n    ...PortDefinitionsTable_portDefinitions\n    id\n  }\n  numberOfEquipment\n}\n\nfragment PortDefinitionsTable_portDefinitions on EquipmentPortDefinition {\n  id\n  name\n  index\n  visibleLabel\n  portType {\n    id\n    name\n  }\n}\n\nfragment PositionDefinitionsTable_positionDefinitions on EquipmentPositionDefinition {\n  id\n  name\n  index\n  visibleLabel\n}\n\nfragment PropertyTypeFormField_propertyType on PropertyType {\n  id\n  name\n  type\n  index\n  stringValue\n  intValue\n  booleanValue\n  floatValue\n  latitudeValue\n  longitudeValue\n  rangeFromValue\n  rangeToValue\n  isEditable\n  isInstanceProperty\n  isMandatory\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'e3d9a9637e5a71aa50d1e5cc61468670';
module.exports = node;
