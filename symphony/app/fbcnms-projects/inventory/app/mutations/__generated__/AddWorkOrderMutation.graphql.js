/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 7c90caee5746f98b2ad1b65ea1ca2203
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type CheckListItemType = "enum" | "simple" | "string" | "%future added value";
export type WorkOrderPriority = "HIGH" | "LOW" | "MEDIUM" | "NONE" | "URGENT" | "%future added value";
export type WorkOrderStatus = "DONE" | "PENDING" | "PLANNED" | "%future added value";
export type AddWorkOrderInput = {|
  name: string,
  description?: ?string,
  workOrderTypeId: string,
  locationId?: ?string,
  projectId?: ?string,
  properties?: ?$ReadOnlyArray<PropertyInput>,
  checkList?: ?$ReadOnlyArray<CheckListItemInput>,
  ownerName?: ?string,
  checkListCategories?: ?$ReadOnlyArray<CheckListCategoryInput>,
  assignee?: ?string,
  index?: ?number,
  status?: ?WorkOrderStatus,
  priority?: ?WorkOrderPriority,
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
  equipmentIDValue?: ?string,
  locationIDValue?: ?string,
  serviceIDValue?: ?string,
  isEditable?: ?boolean,
  isInstanceProperty?: ?boolean,
|};
export type CheckListItemInput = {|
  id?: ?string,
  title: string,
  type: CheckListItemType,
  index?: ?number,
  helpText?: ?string,
  enumValues?: ?string,
  stringValue?: ?string,
  checked?: ?boolean,
|};
export type CheckListCategoryInput = {|
  id?: ?string,
  title: string,
  description?: ?string,
  checkList?: ?$ReadOnlyArray<CheckListItemInput>,
|};
export type AddWorkOrderMutationVariables = {|
  input: AddWorkOrderInput
|};
export type AddWorkOrderMutationResponse = {|
  +addWorkOrder: {|
    +id: string,
    +name: string,
    +description: ?string,
    +ownerName: string,
    +creationDate: any,
    +installDate: ?any,
    +status: WorkOrderStatus,
    +assignee: ?string,
    +location: ?{|
      +id: string,
      +name: string,
    |},
    +workOrderType: {|
      +id: string,
      +name: string,
    |},
    +project: ?{|
      +id: string,
      +name: string,
    |},
  |}
|};
export type AddWorkOrderMutation = {|
  variables: AddWorkOrderMutationVariables,
  response: AddWorkOrderMutationResponse,
|};
*/


/*
mutation AddWorkOrderMutation(
  $input: AddWorkOrderInput!
) {
  addWorkOrder(input: $input) {
    id
    name
    description
    ownerName
    creationDate
    installDate
    status
    assignee
    location {
      id
      name
    }
    workOrderType {
      id
      name
    }
    project {
      id
      name
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "input",
    "type": "AddWorkOrderInput!",
    "defaultValue": null
  }
],
v1 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "id",
  "args": null,
  "storageKey": null
},
v2 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "name",
  "args": null,
  "storageKey": null
},
v3 = [
  (v1/*: any*/),
  (v2/*: any*/)
],
v4 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "addWorkOrder",
    "storageKey": null,
    "args": [
      {
        "kind": "Variable",
        "name": "input",
        "variableName": "input"
      }
    ],
    "concreteType": "WorkOrder",
    "plural": false,
    "selections": [
      (v1/*: any*/),
      (v2/*: any*/),
      {
        "kind": "ScalarField",
        "alias": null,
        "name": "description",
        "args": null,
        "storageKey": null
      },
      {
        "kind": "ScalarField",
        "alias": null,
        "name": "ownerName",
        "args": null,
        "storageKey": null
      },
      {
        "kind": "ScalarField",
        "alias": null,
        "name": "creationDate",
        "args": null,
        "storageKey": null
      },
      {
        "kind": "ScalarField",
        "alias": null,
        "name": "installDate",
        "args": null,
        "storageKey": null
      },
      {
        "kind": "ScalarField",
        "alias": null,
        "name": "status",
        "args": null,
        "storageKey": null
      },
      {
        "kind": "ScalarField",
        "alias": null,
        "name": "assignee",
        "args": null,
        "storageKey": null
      },
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "location",
        "storageKey": null,
        "args": null,
        "concreteType": "Location",
        "plural": false,
        "selections": (v3/*: any*/)
      },
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "workOrderType",
        "storageKey": null,
        "args": null,
        "concreteType": "WorkOrderType",
        "plural": false,
        "selections": (v3/*: any*/)
      },
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "project",
        "storageKey": null,
        "args": null,
        "concreteType": "Project",
        "plural": false,
        "selections": (v3/*: any*/)
      }
    ]
  }
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "AddWorkOrderMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v4/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "AddWorkOrderMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v4/*: any*/)
  },
  "params": {
    "operationKind": "mutation",
    "name": "AddWorkOrderMutation",
    "id": null,
    "text": "mutation AddWorkOrderMutation(\n  $input: AddWorkOrderInput!\n) {\n  addWorkOrder(input: $input) {\n    id\n    name\n    description\n    ownerName\n    creationDate\n    installDate\n    status\n    assignee\n    location {\n      id\n      name\n    }\n    workOrderType {\n      id\n      name\n    }\n    project {\n      id\n      name\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '59a17c820a60f80ca4a246e34ea79176';
module.exports = node;
