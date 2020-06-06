/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 0feed5d3c2967414efa18da5c07c2d67
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type CheckListItemEnumSelectionMode = "multiple" | "single" | "%future added value";
export type CheckListItemType = "cell_scan" | "enum" | "files" | "simple" | "string" | "wifi_scan" | "yes_no" | "%future added value";
export type FileType = "FILE" | "IMAGE" | "%future added value";
export type WorkOrderPriority = "HIGH" | "LOW" | "MEDIUM" | "NONE" | "URGENT" | "%future added value";
export type WorkOrderStatus = "DONE" | "PENDING" | "PLANNED" | "%future added value";
export type YesNoResponse = "NO" | "YES" | "%future added value";
export type AddWorkOrderInput = {|
  name: string,
  description?: ?string,
  workOrderTypeId: string,
  locationId?: ?string,
  projectId?: ?string,
  properties?: ?$ReadOnlyArray<PropertyInput>,
  checkList?: ?$ReadOnlyArray<CheckListItemInput>,
  ownerName?: ?string,
  ownerId?: ?string,
  checkListCategories?: ?$ReadOnlyArray<CheckListCategoryInput>,
  assignee?: ?string,
  assigneeId?: ?string,
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
  nodeIDValue?: ?string,
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
  enumSelectionMode?: ?CheckListItemEnumSelectionMode,
  selectedEnumValues?: ?string,
  stringValue?: ?string,
  checked?: ?boolean,
  files?: ?$ReadOnlyArray<FileInput>,
  yesNoResponse?: ?YesNoResponse,
|};
export type FileInput = {|
  id?: ?string,
  fileName: string,
  sizeInBytes?: ?number,
  modificationTime?: ?number,
  uploadTime?: ?number,
  fileType?: ?FileType,
  mimeType?: ?string,
  storeKey: string,
  annotation?: ?string,
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
    +owner: {|
      +id: string,
      +email: string,
    |},
    +creationDate: any,
    +installDate: ?any,
    +status: WorkOrderStatus,
    +assignedTo: ?{|
      +id: string,
      +email: string,
    |},
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
    +closeDate: ?any,
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
    owner {
      id
      email
    }
    creationDate
    installDate
    status
    assignedTo {
      id
      email
    }
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
    closeDate
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
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "email",
    "args": null,
    "storageKey": null
  }
],
v4 = [
  (v1/*: any*/),
  (v2/*: any*/)
],
v5 = [
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
        "kind": "LinkedField",
        "alias": null,
        "name": "owner",
        "storageKey": null,
        "args": null,
        "concreteType": "User",
        "plural": false,
        "selections": (v3/*: any*/)
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
        "kind": "LinkedField",
        "alias": null,
        "name": "assignedTo",
        "storageKey": null,
        "args": null,
        "concreteType": "User",
        "plural": false,
        "selections": (v3/*: any*/)
      },
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "location",
        "storageKey": null,
        "args": null,
        "concreteType": "Location",
        "plural": false,
        "selections": (v4/*: any*/)
      },
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "workOrderType",
        "storageKey": null,
        "args": null,
        "concreteType": "WorkOrderType",
        "plural": false,
        "selections": (v4/*: any*/)
      },
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "project",
        "storageKey": null,
        "args": null,
        "concreteType": "Project",
        "plural": false,
        "selections": (v4/*: any*/)
      },
      {
        "kind": "ScalarField",
        "alias": null,
        "name": "closeDate",
        "args": null,
        "storageKey": null
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
    "selections": (v5/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "AddWorkOrderMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": (v5/*: any*/)
  },
  "params": {
    "operationKind": "mutation",
    "name": "AddWorkOrderMutation",
    "id": null,
    "text": "mutation AddWorkOrderMutation(\n  $input: AddWorkOrderInput!\n) {\n  addWorkOrder(input: $input) {\n    id\n    name\n    description\n    owner {\n      id\n      email\n    }\n    creationDate\n    installDate\n    status\n    assignedTo {\n      id\n      email\n    }\n    location {\n      id\n      name\n    }\n    workOrderType {\n      id\n      name\n    }\n    project {\n      id\n      name\n    }\n    closeDate\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '59a17c820a60f80ca4a246e34ea79176';
module.exports = node;
