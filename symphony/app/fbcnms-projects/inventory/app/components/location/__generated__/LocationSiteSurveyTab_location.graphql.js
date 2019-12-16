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
type SiteSurveyPane_survey$ref = any;
import type { FragmentReference } from "relay-runtime";
declare export opaque type LocationSiteSurveyTab_location$ref: FragmentReference;
declare export opaque type LocationSiteSurveyTab_location$fragmentType: LocationSiteSurveyTab_location$ref;
export type LocationSiteSurveyTab_location = {|
  +id: string,
  +siteSurveyNeeded: boolean,
  +surveys: $ReadOnlyArray<?{|
    +id: string,
    +completionTimestamp: number,
    +name: string,
    +ownerName: ?string,
    +sourceFile: ?{|
      +id: string,
      +fileName: string,
      +storeKey: ?string,
    |},
    +$fragmentRefs: SiteSurveyPane_survey$ref,
  |}>,
  +$refType: LocationSiteSurveyTab_location$ref,
|};
export type LocationSiteSurveyTab_location$data = LocationSiteSurveyTab_location;
export type LocationSiteSurveyTab_location$key = {
  +$data?: LocationSiteSurveyTab_location$data,
  +$fragmentRefs: LocationSiteSurveyTab_location$ref,
};
*/


const node/*: ReaderFragment*/ = (function(){
var v0 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "id",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Fragment",
  "name": "LocationSiteSurveyTab_location",
  "type": "Location",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "siteSurveyNeeded",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "surveys",
      "storageKey": null,
      "args": null,
      "concreteType": "Survey",
      "plural": true,
      "selections": [
        (v0/*: any*/),
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "completionTimestamp",
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
          "name": "ownerName",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "sourceFile",
          "storageKey": null,
          "args": null,
          "concreteType": "File",
          "plural": false,
          "selections": [
            (v0/*: any*/),
            {
              "kind": "ScalarField",
              "alias": null,
              "name": "fileName",
              "args": null,
              "storageKey": null
            },
            {
              "kind": "ScalarField",
              "alias": null,
              "name": "storeKey",
              "args": null,
              "storageKey": null
            }
          ]
        },
        {
          "kind": "FragmentSpread",
          "name": "SiteSurveyPane_survey",
          "args": null
        }
      ]
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = 'd42a66c9f702b247082792e87c3b7a84';
module.exports = node;
