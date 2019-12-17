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
type ServiceLinksView_links$ref = any;
export type ServiceStatus = "DISCONNECTED" | "IN_SERVICE" | "MAINTENANCE" | "PENDING" | "%future added value";
import type { FragmentReference } from "relay-runtime";
declare export opaque type ServicePanel_service$ref: FragmentReference;
declare export opaque type ServicePanel_service$fragmentType: ServicePanel_service$ref;
export type ServicePanel_service = {|
  +id: string,
  +name: string,
  +externalId: ?string,
  +status: ServiceStatus,
  +customer: ?{|
    +name: string
  |},
  +serviceType: {|
    +name: string
  |},
  +links: $ReadOnlyArray<?{|
    +id: string,
    +$fragmentRefs: ServiceLinksView_links$ref,
  |}>,
  +$refType: ServicePanel_service$ref,
|};
export type ServicePanel_service$data = ServicePanel_service;
export type ServicePanel_service$key = {
  +$data?: ServicePanel_service$data,
  +$fragmentRefs: ServicePanel_service$ref,
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
},
v2 = [
  (v1/*: any*/)
];
return {
  "kind": "Fragment",
  "name": "ServicePanel_service",
  "type": "Service",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
    (v1/*: any*/),
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "externalId",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "status",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "customer",
      "storageKey": null,
      "args": null,
      "concreteType": "Customer",
      "plural": false,
      "selections": (v2/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "serviceType",
      "storageKey": null,
      "args": null,
      "concreteType": "ServiceType",
      "plural": false,
      "selections": (v2/*: any*/)
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
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = 'eaf2a0ad573b59046b4aa8ddc78afcaa';
module.exports = node;
