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
declare export opaque type TextCommentPost_comment$ref: FragmentReference;
declare export opaque type TextCommentPost_comment$fragmentType: TextCommentPost_comment$ref;
export type TextCommentPost_comment = {|
  +id: string,
  +authorName: string,
  +text: string,
  +createTime: any,
  +$refType: TextCommentPost_comment$ref,
|};
export type TextCommentPost_comment$data = TextCommentPost_comment;
export type TextCommentPost_comment$key = {
  +$data?: TextCommentPost_comment$data,
  +$fragmentRefs: TextCommentPost_comment$ref,
};
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "TextCommentPost_comment",
  "type": "Comment",
  "metadata": null,
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
      "name": "authorName",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "text",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "createTime",
      "args": null,
      "storageKey": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = 'bbba65892bc1bce4eda908c5b4722028';
module.exports = node;
