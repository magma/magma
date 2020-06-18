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
declare export opaque type ServiceLinksView_links$ref: FragmentReference;
declare export opaque type ServiceLinksView_links$fragmentType: ServiceLinksView_links$ref;
export type ServiceLinksView_links = $ReadOnlyArray<{|
  +id: string,
  +ports: $ReadOnlyArray<?{|
    +parentEquipment: {|
      +id: string,
      +name: string,
    |},
    +definition: {|
      +id: string,
      +name: string,
    |},
  |}>,
  +$refType: ServiceLinksView_links$ref,
|}>;
export type ServiceLinksView_links$data = ServiceLinksView_links;
export type ServiceLinksView_links$key = $ReadOnlyArray<{
  +$data?: ServiceLinksView_links$data,
  +$fragmentRefs: ServiceLinksView_links$ref,
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
v1 = [
  (v0/*: any*/),
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "name",
    "args": null,
    "storageKey": null
  }
];
return {
  "kind": "Fragment",
  "name": "ServiceLinksView_links",
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
          "selections": (v1/*: any*/)
        },
        {
          "kind": "LinkedField",
          "alias": null,
          "name": "definition",
          "storageKey": null,
          "args": null,
          "concreteType": "EquipmentPortDefinition",
          "plural": false,
          "selections": (v1/*: any*/)
        }
      ]
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = '0e01ac9e52351ad401bac91867ae8714';
module.exports = node;
