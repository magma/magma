/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash ccd4d1725e83faddbfe8638e56feebe8
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type ProjectsTableView_projects$ref = any;
export type AddProjectInput = {|
  name: string,
  description?: ?string,
  creator?: ?string,
  creatorId?: ?string,
  type: string,
  location?: ?string,
  properties?: ?$ReadOnlyArray<PropertyInput>,
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
export type AddProjectMutationVariables = {|
  input: AddProjectInput
|};
export type AddProjectMutationResponse = {|
  +createProject: {|
    +$fragmentRefs: ProjectsTableView_projects$ref
  |}
|};
export type AddProjectMutation = {|
  variables: AddProjectMutationVariables,
  response: AddProjectMutationResponse,
|};
*/


/*
mutation AddProjectMutation(
  $input: AddProjectInput!
) {
  createProject(input: $input) {
    ...ProjectsTableView_projects
    id
  }
}

fragment ProjectsTableView_projects on Project {
  id
  name
  creator
  location {
    id
    name
  }
  type {
    id
    name
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "input",
    "type": "AddProjectInput!",
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
  (v3/*: any*/)
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "AddProjectMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "createProject",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "Project",
        "plural": false,
        "selections": [
          {
            "kind": "FragmentSpread",
            "name": "ProjectsTableView_projects",
            "args": null
          }
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "AddProjectMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "createProject",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "Project",
        "plural": false,
        "selections": [
          (v2/*: any*/),
          (v3/*: any*/),
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "creator",
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
            "selections": (v4/*: any*/)
          },
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "type",
            "storageKey": null,
            "args": null,
            "concreteType": "ProjectType",
            "plural": false,
            "selections": (v4/*: any*/)
          }
        ]
      }
    ]
  },
  "params": {
    "operationKind": "mutation",
    "name": "AddProjectMutation",
    "id": null,
    "text": "mutation AddProjectMutation(\n  $input: AddProjectInput!\n) {\n  createProject(input: $input) {\n    ...ProjectsTableView_projects\n    id\n  }\n}\n\nfragment ProjectsTableView_projects on Project {\n  id\n  name\n  creator\n  location {\n    id\n    name\n  }\n  type {\n    id\n    name\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '4b62a168191f8fdb0ee345d1fd1212fc';
module.exports = node;
