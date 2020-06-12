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
type CommentsActivitiesLog_comments$ref = any;
import type { FragmentReference } from "relay-runtime";
declare export opaque type CommentsActivitiesBox_comments$ref: FragmentReference;
declare export opaque type CommentsActivitiesBox_comments$fragmentType: CommentsActivitiesBox_comments$ref;
export type CommentsActivitiesBox_comments = $ReadOnlyArray<{|
  +$fragmentRefs: CommentsActivitiesLog_comments$ref,
  +$refType: CommentsActivitiesBox_comments$ref,
|}>;
export type CommentsActivitiesBox_comments$data = CommentsActivitiesBox_comments;
export type CommentsActivitiesBox_comments$key = $ReadOnlyArray<{
  +$data?: CommentsActivitiesBox_comments$data,
  +$fragmentRefs: CommentsActivitiesBox_comments$ref,
  ...
}>;
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "CommentsActivitiesBox_comments",
  "type": "Comment",
  "metadata": {
    "plural": true
  },
  "argumentDefinitions": [],
  "selections": [
    {
      "kind": "FragmentSpread",
      "name": "CommentsActivitiesLog_comments",
      "args": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = 'c118bd6378979ed6e9cdf189f1f347c8';
module.exports = node;
