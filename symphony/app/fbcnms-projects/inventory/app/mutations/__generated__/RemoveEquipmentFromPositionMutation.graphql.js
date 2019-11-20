/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash cf231258a59f59d07faf0f27f1b25f2f
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type EquipmentPropertiesCard_position$ref = any;
export type RemoveEquipmentFromPositionMutationVariables = {|
  position_id: string,
  work_order_id?: ?string,
|};
export type RemoveEquipmentFromPositionMutationResponse = {|
  +removeEquipmentFromPosition: ?{|
    +$fragmentRefs: EquipmentPropertiesCard_position$ref
  |}
|};
export type RemoveEquipmentFromPositionMutation = {|
  variables: RemoveEquipmentFromPositionMutationVariables,
  response: RemoveEquipmentFromPositionMutationResponse,
|};
*/


/*
mutation RemoveEquipmentFromPositionMutation(
  $position_id: ID!
  $work_order_id: ID
) {
  removeEquipmentFromPosition(positionId: $position_id, workOrderId: $work_order_id) {
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
    "name": "position_id",
    "type": "ID!",
    "defaultValue": null
  },
  {
    "kind": "LocalArgument",
    "name": "work_order_id",
    "type": "ID",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "Variable",
    "name": "positionId",
    "variableName": "position_id"
  },
  {
    "kind": "Variable",
    "name": "workOrderId",
    "variableName": "work_order_id"
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
    "name": "RemoveEquipmentFromPositionMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "removeEquipmentFromPosition",
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
    "name": "RemoveEquipmentFromPositionMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "removeEquipmentFromPosition",
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
    "name": "RemoveEquipmentFromPositionMutation",
    "id": null,
    "text": "mutation RemoveEquipmentFromPositionMutation(\n  $position_id: ID!\n  $work_order_id: ID\n) {\n  removeEquipmentFromPosition(positionId: $position_id, workOrderId: $work_order_id) {\n    ...EquipmentPropertiesCard_position\n    id\n  }\n}\n\nfragment EquipmentPropertiesCard_position on EquipmentPosition {\n  id\n  definition {\n    id\n    name\n    index\n    visibleLabel\n  }\n  attachedEquipment {\n    id\n    name\n    futureState\n    workOrder {\n      id\n      status\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'd186a2a7a210e4281b5091e52bc5138b';
module.exports = node;
