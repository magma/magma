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
declare export opaque type LocationBreadcrumbsTitle_location$ref: FragmentReference;
declare export opaque type LocationBreadcrumbsTitle_location$fragmentType: LocationBreadcrumbsTitle_location$ref;
export type LocationBreadcrumbsTitle_location = {|
  +id: string,
  +name: string,
  +locationType: {|
    +name: string
  |},
  +locationHierarchy: $ReadOnlyArray<{|
    +id: string,
    +name: string,
    +locationType: {|
      +name: string
    |},
  |}>,
  +$refType: LocationBreadcrumbsTitle_location$ref,
|};
export type LocationBreadcrumbsTitle_location$data = LocationBreadcrumbsTitle_location;
export type LocationBreadcrumbsTitle_location$key = {
  +$data?: LocationBreadcrumbsTitle_location$data,
  +$fragmentRefs: LocationBreadcrumbsTitle_location$ref,
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
v1 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "name",
  "args": null,
  "storageKey": null
},
v2 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "locationType",
  "storageKey": null,
  "args": null,
  "concreteType": "LocationType",
  "plural": false,
  "selections": [
    (v1/*: any*/)
  ]
};
return {
  "kind": "Fragment",
  "name": "LocationBreadcrumbsTitle_location",
  "type": "Location",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
    (v1/*: any*/),
    (v2/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "locationHierarchy",
      "storageKey": null,
      "args": null,
      "concreteType": "Location",
      "plural": true,
      "selections": [
        (v0/*: any*/),
        (v1/*: any*/),
        (v2/*: any*/)
      ]
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = 'c642e01f89689b4dc8ba25d18d5ec3c4';
module.exports = node;
