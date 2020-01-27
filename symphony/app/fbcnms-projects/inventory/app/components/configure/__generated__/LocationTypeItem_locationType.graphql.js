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
declare export opaque type LocationTypeItem_locationType$ref: FragmentReference;
declare export opaque type LocationTypeItem_locationType$fragmentType: LocationTypeItem_locationType$ref;
export type LocationTypeItem_locationType = {|
  +id: string,
  +name: string,
  +index: ?number,
  +propertyTypes: $ReadOnlyArray<?{|
    +$fragmentRefs: DynamicPropertyTypesGrid_propertyTypes$ref
  |}>,
  +numberOfLocations: number,
  +$refType: LocationTypeItem_locationType$ref,
|};
export type LocationTypeItem_locationType$data = LocationTypeItem_locationType;
export type LocationTypeItem_locationType$key = {
  +$data?: LocationTypeItem_locationType$data,
  +$fragmentRefs: LocationTypeItem_locationType$ref,
  ...
};
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "LocationTypeItem_locationType",
  "type": "LocationType",
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
      "name": "index",
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
      "name": "numberOfLocations",
      "args": null,
      "storageKey": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = 'dbc429726327dcfb05c2a80d49cfa429';
module.exports = node;
