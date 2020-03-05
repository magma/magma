/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 0a064347d2a50d75db6f2074e241dc62
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type AddEditWorkOrderTypeCard_editingWorkOrderType$ref = any;
export type CheckListItemType = "enum" | "simple" | "string" | "%future added value";
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "equipment" | "float" | "gps_location" | "int" | "location" | "range" | "service" | "string" | "%future added value";
export type EditWorkOrderTypeInput = {|
  id: string,
  name: string,
  description?: ?string,
  properties?: ?$ReadOnlyArray<?PropertyTypeInput>,
  checkList?: ?$ReadOnlyArray<?CheckListDefinitionInput>,
  checkListCategories?: ?$ReadOnlyArray<CheckListCategoryInput>,
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
export type CheckListDefinitionInput = {|
  id?: ?string,
  title: string,
  type: CheckListItemType,
  index?: ?number,
  enumValues?: ?string,
  helpText?: ?string,
|};
export type CheckListCategoryInput = {|
  id?: ?string,
  title: string,
  description?: ?string,
  checkList?: ?$ReadOnlyArray<CheckListItemInput>,
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
export type EditWorkOrderTypeMutationVariables = {|
  input: EditWorkOrderTypeInput
|};
export type EditWorkOrderTypeMutationResponse = {|
  +editWorkOrderType: {|
    +id: string,
    +name: string,
    +$fragmentRefs: AddEditWorkOrderTypeCard_editingWorkOrderType$ref,
  |}
|};
export type EditWorkOrderTypeMutation = {|
  variables: EditWorkOrderTypeMutationVariables,
  response: EditWorkOrderTypeMutationResponse,
|};
*/


/*
mutation EditWorkOrderTypeMutation(
  $input: EditWorkOrderTypeInput!
) {
  editWorkOrderType(input: $input) {
    id
    name
    ...AddEditWorkOrderTypeCard_editingWorkOrderType
  }
}

fragment AddEditWorkOrderTypeCard_editingWorkOrderType on WorkOrderType {
  id
  name
  description
  numberOfWorkOrders
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
    isMandatory
    isInstanceProperty
    isDeleted
  }
  checkListDefinitions {
    id
    title
    type
    index
    helpText
    enumValues
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "input",
    "type": "EditWorkOrderTypeInput!",
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
  "name": "type",
  "args": null,
  "storageKey": null
},
v5 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "index",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "EditWorkOrderTypeMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "editWorkOrderType",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "WorkOrderType",
        "plural": false,
        "selections": [
          (v2/*: any*/),
          (v3/*: any*/),
          {
            "kind": "FragmentSpread",
            "name": "AddEditWorkOrderTypeCard_editingWorkOrderType",
            "args": null
          }
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "EditWorkOrderTypeMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "editWorkOrderType",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "WorkOrderType",
        "plural": false,
        "selections": [
          (v2/*: any*/),
          (v3/*: any*/),
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
            "name": "numberOfWorkOrders",
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
            "selections": [
              (v2/*: any*/),
              (v3/*: any*/),
              (v4/*: any*/),
              (v5/*: any*/),
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
                "name": "isMandatory",
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
                "name": "isDeleted",
                "args": null,
                "storageKey": null
              }
            ]
          },
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "checkListDefinitions",
            "storageKey": null,
            "args": null,
            "concreteType": "CheckListItemDefinition",
            "plural": true,
            "selections": [
              (v2/*: any*/),
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "title",
                "args": null,
                "storageKey": null
              },
              (v4/*: any*/),
              (v5/*: any*/),
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "helpText",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "enumValues",
                "args": null,
                "storageKey": null
              }
            ]
          }
        ]
      }
    ]
  },
  "params": {
    "operationKind": "mutation",
    "name": "EditWorkOrderTypeMutation",
    "id": null,
    "text": "mutation EditWorkOrderTypeMutation(\n  $input: EditWorkOrderTypeInput!\n) {\n  editWorkOrderType(input: $input) {\n    id\n    name\n    ...AddEditWorkOrderTypeCard_editingWorkOrderType\n  }\n}\n\nfragment AddEditWorkOrderTypeCard_editingWorkOrderType on WorkOrderType {\n  id\n  name\n  description\n  numberOfWorkOrders\n  propertyTypes {\n    id\n    name\n    type\n    index\n    stringValue\n    intValue\n    booleanValue\n    floatValue\n    latitudeValue\n    longitudeValue\n    rangeFromValue\n    rangeToValue\n    isEditable\n    isMandatory\n    isInstanceProperty\n    isDeleted\n  }\n  checkListDefinitions {\n    id\n    title\n    type\n    index\n    helpText\n    enumValues\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '8fad9d18663268f869507ab8bd510f66';
module.exports = node;
