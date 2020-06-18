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
  +$refType: ProjectMoreActionsButton_project$ref,
|};
export type ProjectMoreActionsButton_project$data = ProjectMoreActionsButton_project;
export type ProjectMoreActionsButton_project$key = {
  +$data?: ProjectMoreActionsButton_project$data,
  +$fragmentRefs: ProjectMoreActionsButton_project$ref,
  ...
};
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "ProjectMoreActionsButton_project",
  "type": "Project",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "id",
      "args": null,
      "storageKey": null
    },
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
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = 'a087beefd40f47f0a1bc08f83e8c0667';
module.exports = node;
