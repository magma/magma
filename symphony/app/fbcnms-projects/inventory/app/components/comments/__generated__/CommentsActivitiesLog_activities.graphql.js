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
type ActivityPost_activity$ref = any;
import type { FragmentReference } from "relay-runtime";
declare export opaque type CommentsActivitiesLog_activities$ref: FragmentReference;
declare export opaque type CommentsActivitiesLog_activities$fragmentType: CommentsActivitiesLog_activities$ref;
export type CommentsActivitiesLog_activities = $ReadOnlyArray<{|
  +id: string,
  +createTime: any,
  +$fragmentRefs: ActivityPost_activity$ref,
  +$refType: CommentsActivitiesLog_activities$ref,
|}>;
export type CommentsActivitiesLog_activities$data = CommentsActivitiesLog_activities;
export type CommentsActivitiesLog_activities$key = $ReadOnlyArray<{
  +$data?: CommentsActivitiesLog_activities$data,
  +$fragmentRefs: CommentsActivitiesLog_activities$ref,
  ...
}>;
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "CommentsActivitiesLog_activities",
  "type": "Activity",
  "metadata": {
    "plural": true
  },
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
      "name": "createTime",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "FragmentSpread",
      "name": "ActivityPost_activity",
      "args": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = 'dd73667e5b82cd9bacd860e1b733531e';
module.exports = node;
