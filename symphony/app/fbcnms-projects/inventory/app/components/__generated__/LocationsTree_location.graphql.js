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
declare export opaque type LocationsTree_location$ref: FragmentReference;
declare export opaque type LocationsTree_location$fragmentType: LocationsTree_location$ref;
export type LocationsTree_location = {
  +id: string,
  +externalId: ?string,
  +name: string,
  +locationType: {
    +id: string,
    +name: string,
  },
  +numChildren: number,
  +siteSurveyNeeded: boolean,
};
export type LocationsTree_location$data = LocationsTree_location;
export type LocationsTree_location$key = {
  +$data?: LocationsTree_location$data,
  +$fragmentRefs: LocationsTree_location$ref,
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
};
return {
  "kind": "Fragment",
  "name": "LocationsTree_location",
  "type": "Location",
  "metadata": {
    "mask": false
  },
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "externalId",
      "args": null,
      "storageKey": null
    },
    (v1/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "locationType",
      "storageKey": null,
      "args": null,
      "concreteType": "LocationType",
      "plural": false,
      "selections": [
        (v0/*: any*/),
        (v1/*: any*/)
      ]
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "numChildren",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "siteSurveyNeeded",
      "args": null,
      "storageKey": null
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = '5f66e0c34ff91e14429645285c63ea7f';
module.exports = node;
