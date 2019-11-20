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
type DynamicPropertyTypesGrid_propertyTypes$ref = any;
import type { FragmentReference } from "relay-runtime";
declare export opaque type WorkOrderTypeItem_workOrderType$ref: FragmentReference;
declare export opaque type WorkOrderTypeItem_workOrderType$fragmentType: WorkOrderTypeItem_workOrderType$ref;
export type WorkOrderTypeItem_workOrderType = {|
  +id: string,
  +name: string,
  +propertyTypes: $ReadOnlyArray<?{|
    +$fragmentRefs: DynamicPropertyTypesGrid_propertyTypes$ref
  |}>,
  +numberOfWorkOrders: number,
  +$refType: WorkOrderTypeItem_workOrderType$ref,
|};
export type WorkOrderTypeItem_workOrderType$data = WorkOrderTypeItem_workOrderType;
export type WorkOrderTypeItem_workOrderType$key = {
  +$data?: WorkOrderTypeItem_workOrderType$data,
  +$fragmentRefs: WorkOrderTypeItem_workOrderType$ref,
};
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "WorkOrderTypeItem_workOrderType",
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
        {
          "kind": "FragmentSpread",
          "name": "DynamicPropertyTypesGrid_propertyTypes",
          "args": null
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
// prettier-ignore
(node/*: any*/).hash = '083798fec90c79fa0cc6d7a3e48ced97';
module.exports = node;
