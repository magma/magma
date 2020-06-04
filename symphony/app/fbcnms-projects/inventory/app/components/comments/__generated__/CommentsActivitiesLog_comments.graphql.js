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
      "kind": "FragmentSpread",
      "name": "TextCommentPost_comment",
      "args": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = 'e4b437292785a1a9a105a715bbccfce4';
module.exports = node;
