/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 58069ee91d98011daa20cebac01a2033
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type AddEditProjectTypeCard_editingProjectType$ref = any;
type ProjectTypeCard_projectType$ref = any;
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "float" | "gps_location" | "int" | "node" | "range" | "string" | "%future added value";
export type EditProjectTypeInput = {|
  id: string,
  name: string,
  description?: ?string,
  properties?: ?$ReadOnlyArray<PropertyTypeInput>,
  workOrders?: ?$ReadOnlyArray<WorkOrderDefinitionInput>,
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
export type WorkOrderDefinitionInput = {|
  id?: ?string,
  index?: ?number,
  type: string,
|};
export type EditProjectTypeMutationVariables = {|
  input: EditProjectTypeInput
|};
export type EditProjectTypeMutationResponse = {|
  +editProjectType: {|
    +$fragmentRefs: ProjectTypeCard_projectType$ref & AddEditProjectTypeCard_editingProjectType$ref
  |}
|};
export type EditProjectTypeMutation = {|
  variables: EditProjectTypeMutationVariables,
  response: EditProjectTypeMutationResponse,
|};
*/


/*
mutation EditProjectTypeMutation(
  $input: EditProjectTypeInput!
) {
  editProjectType(input: $input) {
    ...ProjectTypeCard_projectType
    ...AddEditProjectTypeCard_editingProjectType
    id
  }
}

fragment AddEditProjectTypeCard_editingProjectType on ProjectType {
  id
  name
  description
  workOrders {
    id
    type {
      id
      name
    }
  }
  properties {
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
    isMandatory
    isInstanceProperty
    isDeleted
  }
}

fragment ProjectTypeCard_projectType on ProjectType {
  id
  name
  description
  numberOfProjects
  workOrders {
    id
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "input",
    "type": "EditProjectTypeInput!",
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
    "name": "EditProjectTypeMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "editProjectType",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "ProjectType",
        "plural": false,
        "selections": [
          {
            "kind": "FragmentSpread",
            "name": "ProjectTypeCard_projectType",
            "args": null
          },
          {
            "kind": "FragmentSpread",
            "name": "AddEditProjectTypeCard_editingProjectType",
            "args": null
          }
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "EditProjectTypeMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "editProjectType",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "ProjectType",
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
            "name": "numberOfProjects",
            "args": null,
            "storageKey": null
          },
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "workOrders",
            "storageKey": null,
            "args": null,
            "concreteType": "WorkOrderDefinition",
            "plural": true,
            "selections": [
              (v2/*: any*/),
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "type",
                "storageKey": null,
                "args": null,
                "concreteType": "WorkOrderType",
                "plural": false,
                "selections": [
                  (v2/*: any*/),
                  (v3/*: any*/)
                ]
              }
            ]
          },
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "properties",
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
          }
        ]
      }
    ]
  },
  "params": {
    "operationKind": "mutation",
    "name": "EditProjectTypeMutation",
    "id": null,
    "text": "mutation EditProjectTypeMutation(\n  $input: EditProjectTypeInput!\n) {\n  editProjectType(input: $input) {\n    ...ProjectTypeCard_projectType\n    ...AddEditProjectTypeCard_editingProjectType\n    id\n  }\n}\n\nfragment AddEditProjectTypeCard_editingProjectType on ProjectType {\n  id\n  name\n  description\n  workOrders {\n    id\n    type {\n      id\n      name\n    }\n  }\n  properties {\n    id\n    name\n    type\n    nodeType\n    index\n    stringValue\n    intValue\n    booleanValue\n    floatValue\n    latitudeValue\n    longitudeValue\n    rangeFromValue\n    rangeToValue\n    isEditable\n    isMandatory\n    isInstanceProperty\n    isDeleted\n  }\n}\n\nfragment ProjectTypeCard_projectType on ProjectType {\n  id\n  name\n  description\n  numberOfProjects\n  workOrders {\n    id\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'daec04ab717918a0f3359b70b60a4dd9';
module.exports = node;
