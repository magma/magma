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
export type CheckListItemType = "enum" | "simple" | "string" | "%future added value";
import type { FragmentReference } from "relay-runtime";
declare export opaque type CheckListCategoryExpandingPanel_list$ref: FragmentReference;
declare export opaque type CheckListCategoryExpandingPanel_list$fragmentType: CheckListCategoryExpandingPanel_list$ref;
export type CheckListCategoryExpandingPanel_list = $ReadOnlyArray<{|
  +id: string,
  +title: string,
  +description: ?string,
  +checkList: $ReadOnlyArray<{|
    +id: string,
    +index: ?number,
    +type: CheckListItemType,
    +title: string,
    +helpText: ?string,
    +checked: ?boolean,
    +enumValues: ?string,
    +stringValue: ?string,
  |}>,
  +$refType: CheckListCategoryExpandingPanel_list$ref,
|}>;
export type CheckListCategoryExpandingPanel_list$data = CheckListCategoryExpandingPanel_list;
export type CheckListCategoryExpandingPanel_list$key = $ReadOnlyArray<{
  +$data?: CheckListCategoryExpandingPanel_list$data,
  +$fragmentRefs: CheckListCategoryExpandingPanel_list$ref,
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
  "name": "title",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Fragment",
  "name": "CheckListCategoryExpandingPanel_list",
  "type": "CheckListCategory",
  "metadata": {
    "plural": true
  },
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
    (v1/*: any*/),
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "description",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "checkList",
      "storageKey": null,
      "args": null,
      "concreteType": "CheckListItem",
      "plural": true,
      "selections": [
        (v0/*: any*/),
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "index",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "type",
          "args": null,
          "storageKey": null
        },
        (v1/*: any*/),
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "helpText",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "checked",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "enumValues",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "stringValue",
          "args": null,
          "storageKey": null
        }
      ]
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = 'e332dc2f797c581eb5586dceda57c593';
module.exports = node;
