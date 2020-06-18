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
type FileAttachment_file$ref = any;
import type { FragmentReference } from "relay-runtime";
declare export opaque type LocationFloorPlansTab_location$ref: FragmentReference;
declare export opaque type LocationFloorPlansTab_location$fragmentType: LocationFloorPlansTab_location$ref;
export type LocationFloorPlansTab_location = {|
  +id: string,
  +floorPlans: $ReadOnlyArray<?{|
    +id: string,
    +name: string,
    +image: {|
      +$fragmentRefs: FileAttachment_file$ref
    |},
  |}>,
  +$refType: LocationFloorPlansTab_location$ref,
|};
export type LocationFloorPlansTab_location$data = LocationFloorPlansTab_location;
export type LocationFloorPlansTab_location$key = {
  +$data?: LocationFloorPlansTab_location$data,
  +$fragmentRefs: LocationFloorPlansTab_location$ref,
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
};
return {
  "kind": "Fragment",
  "name": "LocationFloorPlansTab_location",
  "type": "Location",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "floorPlans",
      "storageKey": null,
      "args": null,
      "concreteType": "FloorPlan",
      "plural": true,
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
          "name": "image",
          "storageKey": null,
          "args": null,
          "concreteType": "File",
          "plural": false,
          "selections": [
            {
              "kind": "FragmentSpread",
              "name": "FileAttachment_file",
              "args": null
            }
          ]
        }
      ]
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = '33ba200160169ecb72ef20fca4c58fe5';
module.exports = node;
