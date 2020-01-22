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
declare export opaque type CommentsLog_comments$ref: FragmentReference;
declare export opaque type CommentsLog_comments$fragmentType: CommentsLog_comments$ref;
export type CommentsLog_comments = $ReadOnlyArray<{|
  +id: string,
  +$fragmentRefs: TextCommentPost_comment$ref,
  +$refType: CommentsLog_comments$ref,
|}>;
export type CommentsLog_comments$data = CommentsLog_comments;
export type CommentsLog_comments$key = $ReadOnlyArray<{
  +$data?: CommentsLog_comments$data,
  +$fragmentRefs: CommentsLog_comments$ref,
  ...
}>;
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "CommentsLog_comments",
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
(node/*: any*/).hash = 'f895f92613d4a87fed5271af1064b73d';
module.exports = node;
