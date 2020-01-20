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
type EquipmentBreadcrumbs_equipment$ref = any;
import type { FragmentReference } from "relay-runtime";
declare export opaque type AvailableLinksTable_links$ref: FragmentReference;
declare export opaque type AvailableLinksTable_links$fragmentType: AvailableLinksTable_links$ref;
export type AvailableLinksTable_links = $ReadOnlyArray<{|
  +id: string,
  +ports: $ReadOnlyArray<?{|
    +parentEquipment: {|
      +id: string,
      +name: string,
      +positionHierarchy: $ReadOnlyArray<{|
        +parentEquipment: {|
          +id: string
        |}
      |}>,
      +$fragmentRefs: EquipmentBreadcrumbs_equipment$ref,
    |},
    +definition: {|
      +id: string,
      +name: string,
    |},
  |}>,
  +$refType: AvailableLinksTable_links$ref,
|}>;
export type AvailableLinksTable_links$data = AvailableLinksTable_links;
export type AvailableLinksTable_links$key = $ReadOnlyArray<{
  +$data?: AvailableLinksTable_links$data,
  +$fragmentRefs: AvailableLinksTable_links$ref,
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
  "name": "name",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Fragment",
  "name": "AvailableLinksTable_links",
  "type": "Link",
  "metadata": {
    "plural": true
  },
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "ports",
      "storageKey": null,
      "args": null,
      "concreteType": "EquipmentPort",
      "plural": true,
      "selections": [
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "parentEquipment",
          "storageKey": null,
          "args": null,
          "concreteType": "Equipment",
          "plural": false,
          "selections": [
            (v0/*: any*/),
            (v1/*: any*/),
            {
              "kind": "LinkedField",
              "alias": null,
              "name": "positionHierarchy",
              "storageKey": null,
              "args": null,
              "concreteType": "EquipmentPosition",
              "plural": true,
              "selections": [
                {
                  "kind": "LinkedField",
                  "alias": null,
                  "name": "parentEquipment",
                  "storageKey": null,
                  "args": null,
                  "concreteType": "Equipment",
                  "plural": false,
                  "selections": [
                    (v0/*: any*/)
                  ]
                }
              ]
            },
            {
              "kind": "FragmentSpread",
              "name": "EquipmentBreadcrumbs_equipment",
              "args": null
            }
          ]
        },
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "definition",
          "storageKey": null,
          "args": null,
          "concreteType": "EquipmentPortDefinition",
          "plural": false,
          "selections": [
            (v0/*: any*/),
            (v1/*: any*/)
          ]
        }
      ]
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = '6868986f0fbbbb699a7152390136e645';
module.exports = node;
