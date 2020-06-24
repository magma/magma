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
declare export opaque type LocationMenu_location$ref: FragmentReference;
declare export opaque type LocationMenu_location$fragmentType: LocationMenu_location$ref;
export type LocationMenu_location = {|
  +id: string,
  +name: string,
  +locationType: {|
    +id: string
  |},
  +parentLocation: ?{|
    +id: string
  |},
  +children: $ReadOnlyArray<?{|
    +id: string
  |}>,
  +equipments: $ReadOnlyArray<?{|
    +id: string
  |}>,
  +images: $ReadOnlyArray<?{|
    +id: string
  |}>,
  +files: $ReadOnlyArray<?{|
    +id: string
  |}>,
  +surveys: $ReadOnlyArray<?{|
    +id: string
  |}>,
  +$refType: LocationMenu_location$ref,
|};
export type LocationMenu_location$data = LocationMenu_location;
export type LocationMenu_location$key = {
  +$data?: LocationMenu_location$data,
  +$fragmentRefs: LocationMenu_location$ref,
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
},
v1 = [
  (v0/*: any*/)
];
return {
  "kind": "Fragment",
  "name": "LocationMenu_location",
  "type": "Location",
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
      "kind": "LinkedField",
      "alias": null,
      "name": "locationType",
      "storageKey": null,
      "args": null,
      "concreteType": "LocationType",
      "plural": false,
      "selections": (v1/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "parentLocation",
      "storageKey": null,
      "args": null,
      "concreteType": "Location",
      "plural": false,
      "selections": (v1/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "children",
      "storageKey": null,
      "args": null,
      "concreteType": "Location",
      "plural": true,
      "selections": (v1/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "equipments",
      "storageKey": null,
      "args": null,
      "concreteType": "Equipment",
      "plural": true,
      "selections": (v1/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "images",
      "storageKey": null,
      "args": null,
      "concreteType": "File",
      "plural": true,
      "selections": (v1/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "files",
      "storageKey": null,
      "args": null,
      "concreteType": "File",
      "plural": true,
      "selections": (v1/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "surveys",
      "storageKey": null,
      "args": null,
      "concreteType": "Survey",
      "plural": true,
      "selections": (v1/*: any*/)
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = '4155d2d7cef7ceef7b79a4639755fec9';
module.exports = node;
