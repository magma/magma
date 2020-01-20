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
declare export opaque type ProjectTypeWorkOrderTemplatesPanel_workOrderTypes$ref: FragmentReference;
declare export opaque type ProjectTypeWorkOrderTemplatesPanel_workOrderTypes$fragmentType: ProjectTypeWorkOrderTemplatesPanel_workOrderTypes$ref;
export type ProjectTypeWorkOrderTemplatesPanel_workOrderTypes = $ReadOnlyArray<{|
  +id: string,
  +name: string,
  +$refType: ProjectTypeWorkOrderTemplatesPanel_workOrderTypes$ref,
|}>;
export type ProjectTypeWorkOrderTemplatesPanel_workOrderTypes$data = ProjectTypeWorkOrderTemplatesPanel_workOrderTypes;
export type ProjectTypeWorkOrderTemplatesPanel_workOrderTypes$key = $ReadOnlyArray<{
  +$data?: ProjectTypeWorkOrderTemplatesPanel_workOrderTypes$data,
  +$fragmentRefs: ProjectTypeWorkOrderTemplatesPanel_workOrderTypes$ref,
}>;
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "ProjectTypeWorkOrderTemplatesPanel_workOrderTypes",
  "type": "WorkOrderType",
  "metadata": {
    "plural": true
  },
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
(node/*: any*/).hash = 'b8387d4894b17459891a362d8757d82f';
module.exports = node;
