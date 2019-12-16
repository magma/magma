/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ReaderFragment } from 'relay-runtime';
import type { FragmentReference } from "relay-runtime";
declare export opaque type ProjectsMap_projects$ref: FragmentReference;
declare export opaque type ProjectsMap_projects$fragmentType: ProjectsMap_projects$ref;
export type ProjectsMap_projects = $ReadOnlyArray<{|
  +id: string,
  +name: string,
  +location: ?{|
    +id: string,
    +name: string,
    +latitude: number,
    +longitude: number,
  |},
  +numberOfWorkOrders: number,
  +$refType: ProjectsMap_projects$ref,
|}>;
export type ProjectsMap_projects$data = ProjectsMap_projects;
export type ProjectsMap_projects$key = $ReadOnlyArray<{
  +$data?: ProjectsMap_projects$data,
  +$fragmentRefs: ProjectsMap_projects$ref,
}>;
*/


const node/*: ReaderFragment*/ = (function(){
var v0 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "id",
  "args": null,
  "storageKey": null
},
v1 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "name",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Fragment",
  "name": "ProjectsMap_projects",
  "type": "Project",
  "metadata": {
    "plural": true
  },
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
    (v1/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "location",
      "storageKey": null,
      "args": null,
      "concreteType": "Location",
      "plural": false,
      "selections": [
        (v0/*: any*/),
        (v1/*: any*/),
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "latitude",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "longitude",
          "args": null,
          "storageKey": null
        }
      ]
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "numberOfWorkOrders",
      "args": null,
      "storageKey": null
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = 'd43c7f541350f23d3936722943b7ca9b';
module.exports = node;
