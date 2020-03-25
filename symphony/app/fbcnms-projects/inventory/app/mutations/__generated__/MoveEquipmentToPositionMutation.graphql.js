/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 77d821ccb140ff3554894ccdb10b435c
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type EquipmentPropertiesCard_position$ref = any;
export type MoveEquipmentToPositionMutationVariables = {|
  parent_equipment_id: string,
  position_definition_id: string,
  equipment_id: string,
|};
export type MoveEquipmentToPositionMutationResponse = {|
  +moveEquipmentToPosition: {|
    +$fragmentRefs: EquipmentPropertiesCard_position$ref
  |}
|};
export type MoveEquipmentToPositionMutation = {|
  variables: MoveEquipmentToPositionMutationVariables,
  response: MoveEquipmentToPositionMutationResponse,
|};
*/


/*
mutation MoveEquipmentToPositionMutation(
  $parent_equipment_id: ID!
  $position_definition_id: ID!
  $equipment_id: ID!
) {
  moveEquipmentToPosition(parentEquipmentId: $parent_equipment_id, positionDefinitionId: $position_definition_id, equipmentId: $equipment_id) {
    ...EquipmentPropertiesCard_position
    id
  }
}

fragment EquipmentPropertiesCard_position on EquipmentPosition {
  id
  definition {
    id
    name
    index
    visibleLabel
  }
  attachedEquipment {
    id
    name
    futureState
    workOrder {
      id
      status
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "parent_equipment_id",
    "type": "ID!",
    "defaultValue": null
  },
  {
    "kind": "LocalArgument",
    "name": "position_definition_id",
    "type": "ID!",
    "defaultValue": null
  },
  {
    "kind": "LocalArgument",
    "name": "equipment_id",
    "type": "ID!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "Variable",
    "name": "equipmentId",
    "variableName": "equipment_id"
  },
  {
    "kind": "Variable",
    "name": "parentEquipmentId",
    "variableName": "parent_equipment_id"
  },
  {
    "kind": "Variable",
    "name": "positionDefinitionId",
    "variableName": "position_definition_id"
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
    "name": "MoveEquipmentToPositionMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "moveEquipmentToPosition",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "EquipmentPosition",
        "plural": false,
        "selections": [
          {
            "kind": "FragmentSpread",
            "name": "EquipmentPropertiesCard_position",
            "args": null
          }
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "MoveEquipmentToPositionMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "moveEquipmentToPosition",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "EquipmentPosition",
        "plural": false,
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
                "name": "visibleLabel",
                "args": null,
                "storageKey": null
              }
            ]
          },
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
              }
            ]
          }
        ]
      }
    ]
  },
  "params": {
    "operationKind": "mutation",
    "name": "MoveEquipmentToPositionMutation",
    "id": null,
    "text": "mutation MoveEquipmentToPositionMutation(\n  $parent_equipment_id: ID!\n  $position_definition_id: ID!\n  $equipment_id: ID!\n) {\n  moveEquipmentToPosition(parentEquipmentId: $parent_equipment_id, positionDefinitionId: $position_definition_id, equipmentId: $equipment_id) {\n    ...EquipmentPropertiesCard_position\n    id\n  }\n}\n\nfragment EquipmentPropertiesCard_position on EquipmentPosition {\n  id\n  definition {\n    id\n    name\n    index\n    visibleLabel\n  }\n  attachedEquipment {\n    id\n    name\n    futureState\n    workOrder {\n      id\n      status\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '56b02f3e8af0b38080dd01ef75b4e682';
module.exports = node;
