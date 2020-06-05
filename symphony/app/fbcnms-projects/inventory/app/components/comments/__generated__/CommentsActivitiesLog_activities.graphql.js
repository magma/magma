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
export type ActivityField = "ASSIGNEE" | "CREATION_DATE" | "OWNER" | "PRIORITY" | "STATUS" | "%future added value";
import type { FragmentReference } from "relay-runtime";
declare export opaque type CommentsActivitiesLog_activities$ref: FragmentReference;
declare export opaque type CommentsActivitiesLog_activities$fragmentType: CommentsActivitiesLog_activities$ref;
export type CommentsActivitiesLog_activities = $ReadOnlyArray<{|
  +id: string,
  +author: ?{|
    +email: string
  |},
  +isCreate: boolean,
  +changedField: ActivityField,
  +newRelatedNode: ?({|
    +__typename: "User",
    +id: string,
    +email: string,
  |} | {|
    // This will never be '%other', but we need some
    // value in case none of the concrete values match.
    +__typename: "%other"
  |}),
  +oldRelatedNode: ?({|
    +__typename: "User",
    +id: string,
    +email: string,
  |} | {|
    // This will never be '%other', but we need some
    // value in case none of the concrete values match.
    +__typename: "%other"
  |}),
  +oldValue: ?string,
  +newValue: ?string,
  +createTime: any,
  +$refType: CommentsActivitiesLog_activities$ref,
|}>;
export type CommentsActivitiesLog_activities$data = CommentsActivitiesLog_activities;
export type CommentsActivitiesLog_activities$key = $ReadOnlyArray<{
  +$data?: CommentsActivitiesLog_activities$data,
  +$fragmentRefs: CommentsActivitiesLog_activities$ref,
  ...
}>;
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
  "name": "email",
  "args": null,
  "storageKey": null
},
v2 = [
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "__typename",
    "args": null,
    "storageKey": null
  },
  {
    "kind": "InlineFragment",
    "type": "User",
    "selections": [
      (v0/*: any*/),
      (v1/*: any*/)
    ]
  }
];
return {
  "kind": "Fragment",
  "name": "CommentsActivitiesLog_activities",
  "type": "Activity",
  "metadata": {
    "plural": true
  },
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "author",
      "storageKey": null,
      "args": null,
      "concreteType": "User",
      "plural": false,
      "selections": [
        (v1/*: any*/)
      ]
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "isCreate",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "changedField",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "newRelatedNode",
      "storageKey": null,
      "args": null,
      "concreteType": null,
      "plural": false,
      "selections": (v2/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "oldRelatedNode",
      "storageKey": null,
      "args": null,
      "concreteType": null,
      "plural": false,
      "selections": (v2/*: any*/)
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "oldValue",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "newValue",
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
})();
// prettier-ignore
(node/*: any*/).hash = '885c66100ffa7f6e09c28654684ec487';
module.exports = node;
