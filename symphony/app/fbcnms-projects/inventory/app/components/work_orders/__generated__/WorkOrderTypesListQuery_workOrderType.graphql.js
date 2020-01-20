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
declare export opaque type WorkOrderTypesListQuery_workOrderType$ref: FragmentReference;
declare export opaque type WorkOrderTypesListQuery_workOrderType$fragmentType: WorkOrderTypesListQuery_workOrderType$ref;
export type WorkOrderTypesListQuery_workOrderType = {|
  +id: string,
  +name: string,
  +$refType: WorkOrderTypesListQuery_workOrderType$ref,
|};
export type WorkOrderTypesListQuery_workOrderType$data = WorkOrderTypesListQuery_workOrderType;
export type WorkOrderTypesListQuery_workOrderType$key = {
  +$data?: WorkOrderTypesListQuery_workOrderType$data,
  +$fragmentRefs: WorkOrderTypesListQuery_workOrderType$ref,
};
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "WorkOrderTypesListQuery_workOrderType",
  "type": "WorkOrderType",
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
(node/*: any*/).hash = 'f7a668cee825c02441b47bad8a2118ef';
module.exports = node;
