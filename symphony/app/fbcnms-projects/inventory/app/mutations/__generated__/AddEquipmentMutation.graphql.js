/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash cfb6db113cb712a2295c1dda13303ca9
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type EquipmentTable_equipment$ref = any;
export type AddEquipmentInput = {|
  name: string,
  type: string,
  location?: ?string,
  parent?: ?string,
  positionDefinition?: ?string,
  properties?: ?$ReadOnlyArray<PropertyInput>,
  workOrder?: ?string,
  externalId?: ?string,
|};
export type PropertyInput = {|
  id?: ?string,
  propertyTypeID: string,
  stringValue?: ?string,
  intValue?: ?number,
  booleanValue?: ?boolean,
  floatValue?: ?number,
  latitudeValue?: ?number,
  longitudeValue?: ?number,
  rangeFromValue?: ?number,
  rangeToValue?: ?number,
  nodeIDValue?: ?string,
  isEditable?: ?boolean,
  isInstanceProperty?: ?boolean,
|};
export type AddEquipmentMutationVariables = {|
  input: AddEquipmentInput
|};
export type AddEquipmentMutationResponse = {|
  +addEquipment: {|
    +$fragmentRefs: EquipmentTable_equipment$ref
  |}
|};
export type AddEquipmentMutation = {|
  variables: AddEquipmentMutationVariables,
  response: AddEquipmentMutationResponse,
|};
*/


/*
mutation AddEquipmentMutation(
  $input: AddEquipmentInput!
) {
  addEquipment(input: $input) {
    ...EquipmentTable_equipment
    id
  }
}

fragment EquipmentTable_equipment on Equipment {
  id
  name
  futureState
  equipmentType {
    id
    name
  }
  workOrder {
    id
    status
  }
  device {
    up
  }
  services {
    id
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "input",
    "type": "AddEquipmentInput!",
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
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "AddEquipmentMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "addEquipment",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "Equipment",
        "plural": false,
        "selections": [
          {
            "kind": "FragmentSpread",
            "name": "EquipmentTable_equipment",
            "args": null
          }
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "AddEquipmentMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "addEquipment",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "Equipment",
        "plural": false,
        "selections": [
          (v2/*: any*/),
          (v3/*: any*/),
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "futureState",
            "args": null,
            "storageKey": null
          },
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
              (v3/*: any*/)
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
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "device",
            "storageKey": null,
            "args": null,
            "concreteType": "Device",
            "plural": false,
            "selections": [
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "up",
                "args": null,
                "storageKey": null
              }
            ]
          },
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "services",
            "storageKey": null,
            "args": null,
            "concreteType": "Service",
            "plural": true,
            "selections": [
              (v2/*: any*/)
            ]
          }
        ]
      }
    ]
  },
  "params": {
    "operationKind": "mutation",
    "name": "AddEquipmentMutation",
    "id": null,
    "text": "mutation AddEquipmentMutation(\n  $input: AddEquipmentInput!\n) {\n  addEquipment(input: $input) {\n    ...EquipmentTable_equipment\n    id\n  }\n}\n\nfragment EquipmentTable_equipment on Equipment {\n  id\n  name\n  futureState\n  equipmentType {\n    id\n    name\n  }\n  workOrder {\n    id\n    status\n  }\n  device {\n    up\n  }\n  services {\n    id\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '4e7f825e2543f45850548645d993cd19';
module.exports = node;
