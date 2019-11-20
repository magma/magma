/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 1d0084ff72f93c0a6d7490d187516137
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type LocationEquipmentTopology_equipment$ref = any;
type LocationEquipmentTopology_topology$ref = any;
export type LocationNetworkMapTabQueryVariables = {|
  locationId: string
|};
export type LocationNetworkMapTabQueryResponse = {|
  +location: ?{|
    +equipments?: $ReadOnlyArray<?{|
      +$fragmentRefs: LocationEquipmentTopology_equipment$ref
    |}>,
    +topology?: {|
      +$fragmentRefs: LocationEquipmentTopology_topology$ref
    |},
  |}
|};
export type LocationNetworkMapTabQuery = {|
  variables: LocationNetworkMapTabQueryVariables,
  response: LocationNetworkMapTabQueryResponse,
|};
*/


/*
query LocationNetworkMapTabQuery(
  $locationId: ID!
) {
  location: node(id: $locationId) {
    __typename
    ... on Location {
      equipments {
        ...LocationEquipmentTopology_equipment
        id
      }
      topology {
        ...LocationEquipmentTopology_topology
      }
    }
    id
  }
}

fragment LocationEquipmentTopology_equipment on Equipment {
  id
}

fragment LocationEquipmentTopology_topology on NetworkTopology {
  nodes {
    id
    name
  }
  links {
    source
    target
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "locationId",
    "type": "ID!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "Variable",
    "name": "id",
    "variableName": "locationId"
  }
],
v2 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "id",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "LocationNetworkMapTabQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "location",
        "name": "node",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": null,
        "plural": false,
        "selections": [
          {
            "kind": "InlineFragment",
            "type": "Location",
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "equipments",
                "storageKey": null,
                "args": null,
                "concreteType": "Equipment",
                "plural": true,
                "selections": [
                  {
                    "kind": "FragmentSpread",
                    "name": "LocationEquipmentTopology_equipment",
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
                    "name": "LocationEquipmentTopology_topology",
                    "args": null
                  }
                ]
              }
            ]
          }
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "LocationNetworkMapTabQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "location",
        "name": "node",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": null,
        "plural": false,
        "selections": [
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "__typename",
            "args": null,
            "storageKey": null
          },
          (v2/*: any*/),
          {
            "kind": "InlineFragment",
            "type": "Location",
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "equipments",
                "storageKey": null,
                "args": null,
                "concreteType": "Equipment",
                "plural": true,
                "selections": [
                  (v2/*: any*/)
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
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "nodes",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "Equipment",
                    "plural": true,
                    "selections": [
                      (v2/*: any*/),
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "name",
                        "args": null,
                        "storageKey": null
                      }
                    ]
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "links",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "TopologyLink",
                    "plural": true,
                    "selections": [
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "source",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "target",
                        "args": null,
                        "storageKey": null
                      }
                    ]
                  }
                ]
              }
            ]
          }
        ]
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "LocationNetworkMapTabQuery",
    "id": null,
    "text": "query LocationNetworkMapTabQuery(\n  $locationId: ID!\n) {\n  location: node(id: $locationId) {\n    __typename\n    ... on Location {\n      equipments {\n        ...LocationEquipmentTopology_equipment\n        id\n      }\n      topology {\n        ...LocationEquipmentTopology_topology\n      }\n    }\n    id\n  }\n}\n\nfragment LocationEquipmentTopology_equipment on Equipment {\n  id\n}\n\nfragment LocationEquipmentTopology_topology on NetworkTopology {\n  nodes {\n    id\n    name\n  }\n  links {\n    source\n    target\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'f90e103b70337af4c1df59192964fdb0';
module.exports = node;
