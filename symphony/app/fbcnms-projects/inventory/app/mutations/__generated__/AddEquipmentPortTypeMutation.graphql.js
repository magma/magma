/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 9be64b7c0280899298715826b0107562
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type EquipmentPortTypeItem_equipmentPortType$ref = any;
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "float" | "gps_location" | "int" | "node" | "range" | "string" | "%future added value";
export type AddEquipmentPortTypeInput = {|
  name: string,
  properties?: ?$ReadOnlyArray<PropertyTypeInput>,
  linkProperties?: ?$ReadOnlyArray<PropertyTypeInput>,
|};
export type PropertyTypeInput = {|
  id?: ?string,
  externalId?: ?string,
  name: string,
  type: PropertyKind,
  nodeType?: ?string,
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
export type AddEquipmentPortTypeMutationVariables = {|
  input: AddEquipmentPortTypeInput
|};
export type AddEquipmentPortTypeMutationResponse = {|
  +addEquipmentPortType: {|
    +id: string,
    +name: string,
    +$fragmentRefs: EquipmentPortTypeItem_equipmentPortType$ref,
  |}
|};
export type AddEquipmentPortTypeMutation = {|
  variables: AddEquipmentPortTypeMutationVariables,
  response: AddEquipmentPortTypeMutationResponse,
|};
*/


/*
mutation AddEquipmentPortTypeMutation(
  $input: AddEquipmentPortTypeInput!
) {
  addEquipmentPortType(input: $input) {
    id
    name
    ...EquipmentPortTypeItem_equipmentPortType
  }
}

fragment DynamicPropertyTypesGrid_propertyTypes on PropertyType {
  ...PropertyTypeFormField_propertyType
  id
  index
}

fragment EquipmentPortTypeItem_equipmentPortType on EquipmentPortType {
  id
  name
  numberOfPortDefinitions
  propertyTypes {
    ...DynamicPropertyTypesGrid_propertyTypes
    id
  }
  linkPropertyTypes {
    ...DynamicPropertyTypesGrid_propertyTypes
    id
  }
}

fragment PropertyTypeFormField_propertyType on PropertyType {
  id
  name
  type
  nodeType
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
  category
  isDeleted
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "input",
    "type": "AddEquipmentPortTypeInput!",
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
v4 = [
  (v2/*: any*/),
  (v3/*: any*/),
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "type",
    "args": null,
    "storageKey": null
  },
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "nodeType",
    "args": null,
    "storageKey": null
  },
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "index",
    "args": null,
    "storageKey": null
  },
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
  },
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "category",
    "args": null,
    "storageKey": null
  },
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "isDeleted",
    "args": null,
    "storageKey": null
  }
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "AddEquipmentPortTypeMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "addEquipmentPortType",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "EquipmentPortType",
        "plural": false,
        "selections": [
          (v2/*: any*/),
          (v3/*: any*/),
          {
            "kind": "FragmentSpread",
            "name": "EquipmentPortTypeItem_equipmentPortType",
            "args": null
          }
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "AddEquipmentPortTypeMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "addEquipmentPortType",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "EquipmentPortType",
        "plural": false,
        "selections": [
          (v2/*: any*/),
          (v3/*: any*/),
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "numberOfPortDefinitions",
            "args": null,
            "storageKey": null
          },
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "propertyTypes",
            "storageKey": null,
            "args": null,
            "concreteType": "PropertyType",
            "plural": true,
            "selections": (v4/*: any*/)
          },
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "linkPropertyTypes",
            "storageKey": null,
            "args": null,
            "concreteType": "PropertyType",
            "plural": true,
            "selections": (v4/*: any*/)
          }
        ]
      }
    ]
  },
  "params": {
    "operationKind": "mutation",
    "name": "AddEquipmentPortTypeMutation",
    "id": null,
    "text": "mutation AddEquipmentPortTypeMutation(\n  $input: AddEquipmentPortTypeInput!\n) {\n  addEquipmentPortType(input: $input) {\n    id\n    name\n    ...EquipmentPortTypeItem_equipmentPortType\n  }\n}\n\nfragment DynamicPropertyTypesGrid_propertyTypes on PropertyType {\n  ...PropertyTypeFormField_propertyType\n  id\n  index\n}\n\nfragment EquipmentPortTypeItem_equipmentPortType on EquipmentPortType {\n  id\n  name\n  numberOfPortDefinitions\n  propertyTypes {\n    ...DynamicPropertyTypesGrid_propertyTypes\n    id\n  }\n  linkPropertyTypes {\n    ...DynamicPropertyTypesGrid_propertyTypes\n    id\n  }\n}\n\nfragment PropertyTypeFormField_propertyType on PropertyType {\n  id\n  name\n  type\n  nodeType\n  index\n  stringValue\n  intValue\n  booleanValue\n  floatValue\n  latitudeValue\n  longitudeValue\n  rangeFromValue\n  rangeToValue\n  isEditable\n  isInstanceProperty\n  isMandatory\n  category\n  isDeleted\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '454f79d6f06b34bbd5fe4c8aae3f4743';
module.exports = node;
