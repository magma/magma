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
type CommentsLog_comments$ref = any;
import type { FragmentReference } from "relay-runtime";
declare export opaque type CommentsBox_comments$ref: FragmentReference;
declare export opaque type CommentsBox_comments$fragmentType: CommentsBox_comments$ref;
export type CommentsBox_comments = $ReadOnlyArray<{|
  +$fragmentRefs: CommentsLog_comments$ref,
  +$refType: CommentsBox_comments$ref,
|}>;
export type CommentsBox_comments$data = CommentsBox_comments;
export type CommentsBox_comments$key = $ReadOnlyArray<{
  +$data?: CommentsBox_comments$data,
  +$fragmentRefs: CommentsBox_comments$ref,
}>;
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "CommentsBox_comments",
  "type": "Comment",
  "metadata": {
    "plural": true
  },
  "argumentDefinitions": [],
  "selections": [
    {
      "kind": "FragmentSpread",
      "name": "CommentsLog_comments",
      "args": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = 'aa03eae87ac4a993d91b61a591895c6f';
module.exports = node;
