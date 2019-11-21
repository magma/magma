/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 67e513c8d7d9a9e625b0dadb992d96ea
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type ServiceEquipmentTopology_terminationPoints$ref = any;
type ServiceEquipmentTopology_topology$ref = any;
export type ServiceNetworkMapTabQueryVariables = {|
  serviceId: string
|};
export type ServiceNetworkMapTabQueryResponse = {|
  +service: ?{|
    +terminationPoints: $ReadOnlyArray<?{|
      +$fragmentRefs: ServiceEquipmentTopology_terminationPoints$ref
    |}>,
    +topology: {|
      +$fragmentRefs: ServiceEquipmentTopology_topology$ref
    |},
  |}
|};
export type ServiceNetworkMapTabQuery = {|
  variables: ServiceNetworkMapTabQueryVariables,
  response: ServiceNetworkMapTabQueryResponse,
|};
*/


/*
query ServiceNetworkMapTabQuery(
  $serviceId: ID!
) {
  service(id: $serviceId) {
    terminationPoints {
      ...ServiceEquipmentTopology_terminationPoints
      id
    }
    topology {
      ...ServiceEquipmentTopology_topology
    }
    id
  }
}

fragment ServiceEquipmentTopology_terminationPoints on Equipment {
  id
}

fragment ServiceEquipmentTopology_topology on NetworkTopology {
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
    "name": "serviceId",
    "type": "ID!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "Variable",
    "name": "id",
    "variableName": "serviceId"
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
    "name": "ServiceNetworkMapTabQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "service",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "Service",
        "plural": false,
        "selections": [
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
          }
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "ServiceNetworkMapTabQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "service",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "Service",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "terminationPoints",
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
          },
          (v2/*: any*/)
        ]
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "ServiceNetworkMapTabQuery",
    "id": null,
    "text": "query ServiceNetworkMapTabQuery(\n  $serviceId: ID!\n) {\n  service(id: $serviceId) {\n    terminationPoints {\n      ...ServiceEquipmentTopology_terminationPoints\n      id\n    }\n    topology {\n      ...ServiceEquipmentTopology_topology\n    }\n    id\n  }\n}\n\nfragment ServiceEquipmentTopology_terminationPoints on Equipment {\n  id\n}\n\nfragment ServiceEquipmentTopology_topology on NetworkTopology {\n  nodes {\n    id\n    name\n  }\n  links {\n    source\n    target\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'eca7ba62d713f7503fa43cdff70947f8';
module.exports = node;
