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
declare export opaque type ProjectMoreActionsButton_project$ref: FragmentReference;
declare export opaque type ProjectMoreActionsButton_project$fragmentType: ProjectMoreActionsButton_project$ref;
export type ProjectMoreActionsButton_project = {|
  +id: string,
  +name: string,
  +numberOfWorkOrders: number,
  +type: {|
    +id: string
  |},
  +$refType: ProjectMoreActionsButton_project$ref,
|};
export type ProjectMoreActionsButton_project$data = ProjectMoreActionsButton_project;
export type ProjectMoreActionsButton_project$key = {
  +$data?: ProjectMoreActionsButton_project$data,
  +$fragmentRefs: ProjectMoreActionsButton_project$ref,
  ...
};
*/


const node/*: ReaderFragment*/ = (function(){
var v0 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "id",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Fragment",
  "name": "ProjectMoreActionsButton_project",
  "type": "Project",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "name",
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
      "name": "type",
      "storageKey": null,
      "args": null,
      "concreteType": "ProjectType",
      "plural": false,
      "selections": [
        (v0/*: any*/)
      ]
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = 'fab778d2924d1e0b30c094cc7dfa5572';
module.exports = node;
