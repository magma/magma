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
declare export opaque type CheckListCategoryTable_list$ref: FragmentReference;
declare export opaque type CheckListCategoryTable_list$fragmentType: CheckListCategoryTable_list$ref;
export type CheckListCategoryTable_list = $ReadOnlyArray<{|
  +id: string,
  +title: string,
  +description: ?string,
  +checkList: $ReadOnlyArray<{|
    +id: string
  |}>,
  +$refType: CheckListCategoryTable_list$ref,
|}>;
export type CheckListCategoryTable_list$data = CheckListCategoryTable_list;
export type CheckListCategoryTable_list$key = $ReadOnlyArray<{
  +$data?: CheckListCategoryTable_list$data,
  +$fragmentRefs: CheckListCategoryTable_list$ref,
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
};
return {
  "kind": "Fragment",
  "name": "CheckListCategoryTable_list",
  "type": "CheckListCategory",
  "metadata": {
    "plural": true
  },
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "title",
      "args": null,
      "storageKey": null
    },
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
        (v0/*: any*/)
      ]
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = '0919f0462cef49731594b133d7789ed4';
module.exports = node;
