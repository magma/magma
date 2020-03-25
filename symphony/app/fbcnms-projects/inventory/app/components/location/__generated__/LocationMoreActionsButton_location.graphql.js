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
declare export opaque type LocationMoreActionsButton_location$ref: FragmentReference;
declare export opaque type LocationMoreActionsButton_location$fragmentType: LocationMoreActionsButton_location$ref;
export type LocationMoreActionsButton_location = {|
  +id: string,
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
  +$refType: LocationMoreActionsButton_location$ref,
|};
export type LocationMoreActionsButton_location$data = LocationMoreActionsButton_location;
export type LocationMoreActionsButton_location$key = {
  +$data?: LocationMoreActionsButton_location$data,
  +$fragmentRefs: LocationMoreActionsButton_location$ref,
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
  "name": "LocationMoreActionsButton_location",
  "type": "Location",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
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
(node/*: any*/).hash = 'c95f9da4658ca6a495b30d5d6809b583';
module.exports = node;
