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
type PropertyFormField_property$ref = any;
import type { FragmentReference } from "relay-runtime";
declare export opaque type DynamicPropertiesGrid_properties$ref: FragmentReference;
declare export opaque type DynamicPropertiesGrid_properties$fragmentType: DynamicPropertiesGrid_properties$ref;
export type DynamicPropertiesGrid_properties = $ReadOnlyArray<{|
  +propertyType: {|
    +id: string,
    +index: ?number,
  |},
  +$fragmentRefs: PropertyFormField_property$ref,
  +$refType: DynamicPropertiesGrid_properties$ref,
|}>;
export type DynamicPropertiesGrid_properties$data = DynamicPropertiesGrid_properties;
export type DynamicPropertiesGrid_properties$key = $ReadOnlyArray<{
  +$data?: DynamicPropertiesGrid_properties$data,
  +$fragmentRefs: DynamicPropertiesGrid_properties$ref,
  ...
}>;
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "DynamicPropertiesGrid_properties",
  "type": "Property",
  "metadata": {
    "plural": true
  },
  "argumentDefinitions": [],
  "selections": [
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "propertyType",
      "storageKey": null,
      "args": null,
      "concreteType": "PropertyType",
      "plural": false,
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
          "name": "index",
          "args": null,
          "storageKey": null
        }
      ]
    },
    {
      "kind": "FragmentSpread",
      "name": "PropertyFormField_property",
      "args": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = '7b83547f439381a2346b2b5c487b5134';
module.exports = node;
