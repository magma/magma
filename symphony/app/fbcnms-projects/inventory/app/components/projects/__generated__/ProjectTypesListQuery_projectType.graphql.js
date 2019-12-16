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
declare export opaque type ProjectTypesListQuery_projectType$ref: FragmentReference;
declare export opaque type ProjectTypesListQuery_projectType$fragmentType: ProjectTypesListQuery_projectType$ref;
export type ProjectTypesListQuery_projectType = {|
  +id: string,
  +name: string,
  +$refType: ProjectTypesListQuery_projectType$ref,
|};
export type ProjectTypesListQuery_projectType$data = ProjectTypesListQuery_projectType;
export type ProjectTypesListQuery_projectType$key = {
  +$data?: ProjectTypesListQuery_projectType$data,
  +$fragmentRefs: ProjectTypesListQuery_projectType$ref,
};
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "ProjectTypesListQuery_projectType",
  "type": "ProjectType",
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
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = 'bfe406d4456143fb9fa4a1ad5470c960';
module.exports = node;
