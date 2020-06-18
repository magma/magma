/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash e2ebaf6750af599b5ad5530334994e0a
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type EquipmentBreadcrumbs_equipment$ref = any;
export type InventoryEntitiesTypeaheadQueryVariables = {|
  name: string
|};
export type InventoryEntitiesTypeaheadQueryResponse = {|
  +searchForNode: {|
    +edges: ?$ReadOnlyArray<{|
      +node: ?({|
        +__typename: "Location",
        +id: string,
        +externalId: ?string,
        +name: string,
        +locationType: {|
          +name: string
        |},
        +locationHierarchy: $ReadOnlyArray<{|
          +id: string,
          +name: string,
          +locationType: {|
            +name: string
          |},
        |}>,
      |} | {|
        +__typename: "Equipment",
        +id: string,
        +externalId: ?string,
        +name: string,
        +equipmentType: {|
          +name: string
        |},
        +$fragmentRefs: EquipmentBreadcrumbs_equipment$ref,
      |} | {|
        // This will never be '%other', but we need some
        // value in case none of the concrete values match.
        +__typename: "%other"
      |})
    |}>
  |}
|};
export type InventoryEntitiesTypeaheadQuery = {|
  variables: InventoryEntitiesTypeaheadQueryVariables,
  response: InventoryEntitiesTypeaheadQueryResponse,
|};
*/


/*
query InventoryEntitiesTypeaheadQuery(
  $name: String!
) {
  searchForNode(name: $name, first: 10) {
    edges {
      node {
        __typename
        ... on Location {
          id
          externalId
          name
          locationType {
            name
            id
          }
          locationHierarchy {
            id
            name
            locationType {
              name
              id
            }
          }
        }
        ... on Equipment {
          id
          externalId
          name
          equipmentType {
            name
            id
          }
          ...EquipmentBreadcrumbs_equipment
        }
        id
      }
    }
  }
}

fragment EquipmentBreadcrumbs_equipment on Equipment {
  id
  name
  equipmentType {
    id
    name
  }
  locationHierarchy {
    id
    name
    locationType {
      name
      id
    }
  }
  positionHierarchy {
    id
    definition {
      id
      name
      visibleLabel
    }
    parentEquipment {
      id
      name
      equipmentType {
        id
        name
      }
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "name",
    "type": "String!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "Literal",
    "name": "first",
    "value": 10
  },
  {
    "kind": "Variable",
    "name": "name",
    "variableName": "name"
  }
],
v2 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "__typename",
  "args": null,
  "storageKey": null
},
v3 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "id",
  "args": null,
  "storageKey": null
},
v4 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "externalId",
  "args": null,
  "storageKey": null
},
v5 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "name",
  "args": null,
  "storageKey": null
},
v6 = [
  (v5/*: any*/)
],
v7 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "locationType",
  "storageKey": null,
  "args": null,
  "concreteType": "LocationType",
  "plural": false,
  "selections": (v6/*: any*/)
},
v8 = [
  (v5/*: any*/),
  (v3/*: any*/)
],
v9 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "locationType",
  "storageKey": null,
  "args": null,
  "concreteType": "LocationType",
  "plural": false,
  "selections": (v8/*: any*/)
},
v10 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "locationHierarchy",
  "storageKey": null,
  "args": null,
  "concreteType": "Location",
  "plural": true,
  "selections": [
    (v3/*: any*/),
    (v5/*: any*/),
    (v9/*: any*/)
  ]
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "InventoryEntitiesTypeaheadQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "searchForNode",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "SearchNodesConnection",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "edges",
            "storageKey": null,
            "args": null,
            "concreteType": "SearchNodeEdge",
            "plural": true,
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "node",
                "storageKey": null,
                "args": null,
                "concreteType": null,
                "plural": false,
                "selections": [
                  (v2/*: any*/),
                  {
                    "kind": "InlineFragment",
                    "type": "Location",
                    "selections": [
                      (v3/*: any*/),
                      (v4/*: any*/),
                      (v5/*: any*/),
                      (v7/*: any*/),
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "locationHierarchy",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "Location",
                        "plural": true,
                        "selections": [
                          (v3/*: any*/),
                          (v5/*: any*/),
                          (v7/*: any*/)
                        ]
                      }
                    ]
                  },
                  {
                    "kind": "InlineFragment",
                    "type": "Equipment",
                    "selections": [
                      (v3/*: any*/),
                      (v4/*: any*/),
                      (v5/*: any*/),
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "equipmentType",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "EquipmentType",
                        "plural": false,
                        "selections": (v6/*: any*/)
                      },
                      {
                        "kind": "FragmentSpread",
                        "name": "EquipmentBreadcrumbs_equipment",
                        "args": null
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
  "operation": {
    "kind": "Operation",
    "name": "InventoryEntitiesTypeaheadQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "searchForNode",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "SearchNodesConnection",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "edges",
            "storageKey": null,
            "args": null,
            "concreteType": "SearchNodeEdge",
            "plural": true,
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "node",
                "storageKey": null,
                "args": null,
                "concreteType": null,
                "plural": false,
                "selections": [
                  (v2/*: any*/),
                  (v3/*: any*/),
                  {
                    "kind": "InlineFragment",
                    "type": "Location",
                    "selections": [
                      (v4/*: any*/),
                      (v5/*: any*/),
                      (v9/*: any*/),
                      (v10/*: any*/)
                    ]
                  },
                  {
                    "kind": "InlineFragment",
                    "type": "Equipment",
                    "selections": [
                      (v4/*: any*/),
                      (v5/*: any*/),
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "equipmentType",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "EquipmentType",
                        "plural": false,
                        "selections": (v8/*: any*/)
                      },
                      (v10/*: any*/),
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "positionHierarchy",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "EquipmentPosition",
                        "plural": true,
                        "selections": [
                          (v3/*: any*/),
                          {
                            "kind": "LinkedField",
                            "alias": null,
                            "name": "definition",
                            "storageKey": null,
                            "args": null,
                            "concreteType": "EquipmentPositionDefinition",
                            "plural": false,
                            "selections": [
                              (v3/*: any*/),
                              (v5/*: any*/),
                              {
                                "kind": "ScalarField",
                                "alias": null,
                                "name": "visibleLabel",
                                "args": null,
                                "storageKey": null
                              }
                            ]
                          },
                          {
                            "kind": "LinkedField",
                            "alias": null,
                            "name": "parentEquipment",
                            "storageKey": null,
                            "args": null,
                            "concreteType": "Equipment",
                            "plural": false,
                            "selections": [
                              (v3/*: any*/),
                              (v5/*: any*/),
                              {
                                "kind": "LinkedField",
                                "alias": null,
                                "name": "equipmentType",
                                "storageKey": null,
                                "args": null,
                                "concreteType": "EquipmentType",
                                "plural": false,
                                "selections": [
                                  (v3/*: any*/),
                                  (v5/*: any*/)
                                ]
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
          }
        ]
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "InventoryEntitiesTypeaheadQuery",
    "id": null,
    "text": "query InventoryEntitiesTypeaheadQuery(\n  $name: String!\n) {\n  searchForNode(name: $name, first: 10) {\n    edges {\n      node {\n        __typename\n        ... on Location {\n          id\n          externalId\n          name\n          locationType {\n            name\n            id\n          }\n          locationHierarchy {\n            id\n            name\n            locationType {\n              name\n              id\n            }\n          }\n        }\n        ... on Equipment {\n          id\n          externalId\n          name\n          equipmentType {\n            name\n            id\n          }\n          ...EquipmentBreadcrumbs_equipment\n        }\n        id\n      }\n    }\n  }\n}\n\nfragment EquipmentBreadcrumbs_equipment on Equipment {\n  id\n  name\n  equipmentType {\n    id\n    name\n  }\n  locationHierarchy {\n    id\n    name\n    locationType {\n      name\n      id\n    }\n  }\n  positionHierarchy {\n    id\n    definition {\n      id\n      name\n      visibleLabel\n    }\n    parentEquipment {\n      id\n      name\n      equipmentType {\n        id\n        name\n      }\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '7738e72176118873e296e7140a8d048b';
module.exports = node;
