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
declare export opaque type LocationBreadcrumbsTitle_locationDetails$ref: FragmentReference;
declare export opaque type LocationBreadcrumbsTitle_locationDetails$fragmentType: LocationBreadcrumbsTitle_locationDetails$ref;
export type LocationBreadcrumbsTitle_locationDetails = {|
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
  +$refType: LocationBreadcrumbsTitle_locationDetails$ref,
|};
export type LocationBreadcrumbsTitle_locationDetails$data = LocationBreadcrumbsTitle_locationDetails;
export type LocationBreadcrumbsTitle_locationDetails$key = {
  +$data?: LocationBreadcrumbsTitle_locationDetails$data,
  +$fragmentRefs: LocationBreadcrumbsTitle_locationDetails$ref,
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
  "name": "LocationBreadcrumbsTitle_locationDetails",
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
(node/*: any*/).hash = '807c6b11117d3143cb9babd7a3239785';
module.exports = node;
