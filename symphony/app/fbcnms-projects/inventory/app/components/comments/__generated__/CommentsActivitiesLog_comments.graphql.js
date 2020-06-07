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
type TextCommentPost_comment$ref = any;
import type { FragmentReference } from "relay-runtime";
declare export opaque type CommentsActivitiesLog_comments$ref: FragmentReference;
declare export opaque type CommentsActivitiesLog_comments$fragmentType: CommentsActivitiesLog_comments$ref;
export type CommentsActivitiesLog_comments = $ReadOnlyArray<{|
  +id: string,
  +createTime: any,
  +$fragmentRefs: TextCommentPost_comment$ref,
  +$refType: CommentsActivitiesLog_comments$ref,
|}>;
export type CommentsActivitiesLog_comments$data = CommentsActivitiesLog_comments;
export type CommentsActivitiesLog_comments$key = $ReadOnlyArray<{
  +$data?: CommentsActivitiesLog_comments$data,
  +$fragmentRefs: CommentsActivitiesLog_comments$ref,
  ...
}>;
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "CommentsActivitiesLog_comments",
  "type": "Comment",
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
      "name": "TextCommentPost_comment",
      "args": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = 'df495dc20009f04644722a4f796c7db9';
module.exports = node;
