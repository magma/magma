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
type ServiceDetailsPanel_service$ref = any;
type ServiceEquipmentTopology_terminationPoints$ref = any;
type ServiceEquipmentTopology_topology$ref = any;
type ServiceLinksView_links$ref = any;
import type { FragmentReference } from "relay-runtime";
declare export opaque type ServiceCard_service$ref: FragmentReference;
declare export opaque type ServiceCard_service$fragmentType: ServiceCard_service$ref;
export type ServiceCard_service = {|
  +id: string,
  +name: string,
  +serviceType: {|
    +name: string
  |},
  +links: $ReadOnlyArray<?{|
    +id: string,
    +$fragmentRefs: ServiceLinksView_links$ref,
  |}>,
  +terminationPoints: $ReadOnlyArray<?{|
    +$fragmentRefs: ServiceEquipmentTopology_terminationPoints$ref
  |}>,
  +topology: {|
    +$fragmentRefs: ServiceEquipmentTopology_topology$ref
  |},
  +$fragmentRefs: ServiceDetailsPanel_service$ref,
  +$refType: ServiceCard_service$ref,
|};
export type ServiceCard_service$data = ServiceCard_service;
export type ServiceCard_service$key = {
  +$data?: ServiceCard_service$data,
  +$fragmentRefs: ServiceCard_service$ref,
};
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
  "name": "ServiceCard_service",
  "type": "Service",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
    (v1/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "serviceType",
      "storageKey": null,
      "args": null,
      "concreteType": "ServiceType",
      "plural": false,
      "selections": [
        (v1/*: any*/)
      ]
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "links",
      "storageKey": null,
      "args": null,
      "concreteType": "Link",
      "plural": true,
      "selections": [
        (v0/*: any*/),
        {
          "kind": "FragmentSpread",
          "name": "ServiceLinksView_links",
          "args": null
        }
      ]
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "terminationPoints",
      "storageKey": null,
      "args": null,
      "concreteType": "Equipment",
      "plural": true,
      "selections": [
        {
          "kind": "FragmentSpread",
          "name": "ServiceEquipmentTopology_terminationPoints",
          "args": null
        }
      ]
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "topology",
      "storageKey": null,
      "args": null,
      "concreteType": "NetworkTopology",
      "plural": false,
      "selections": [
        {
          "kind": "FragmentSpread",
          "name": "ServiceEquipmentTopology_topology",
          "args": null
        }
      ]
    },
    {
      "kind": "FragmentSpread",
      "name": "ServiceDetailsPanel_service",
      "args": null
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = '1f502fbd72a8f98d41147f74d27395f3';
module.exports = node;
