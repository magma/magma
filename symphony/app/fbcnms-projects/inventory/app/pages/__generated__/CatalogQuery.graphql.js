/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 6822f4c7294374a256559f7affed3685
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type EquipmentTypeItem_equipmentType$ref = any;
export type CatalogQueryVariables = {||};
export type CatalogQueryResponse = {|
  +equipmentTypes: {|
    +edges: ?$ReadOnlyArray<?{|
      +node: ?{|
        +$fragmentRefs: EquipmentTypeItem_equipmentType$ref
      |}
    |}>
  |}
|};
export type CatalogQuery = {|
  variables: CatalogQueryVariables,
  response: CatalogQueryResponse,
|};
*/


/*
query CatalogQuery {
  equipmentTypes(first: 50) {
    edges {
      node {
        ...EquipmentTypeItem_equipmentType
      }
    }
  }
}

fragment EquipmentTypeItem_equipmentType on EquipmentType {
  id
  name
  propertyTypes {
    ...PropertyTypeFormField_propertyType
    id
  }
  positionDefinitions {
    id
    name
    visibleLabel
  }
  portDefinitions {
    id
    name
    visibleLabel
    type
  }
  numberOfEquipment
}

fragment PropertyTypeFormField_propertyType on PropertyType {
  id
  name
  type
  stringValue
  intValue
  booleanValue
  floatValue
  latitudeValue
  longitudeValue
  isEditable
  isInstanceProperty
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "Literal",
    "name": "first",
    "value": 50,
    "type": "Int"
  }
],
v1 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "id",
  "args": null,
  "storageKey": null
},
v2 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "name",
  "args": null,
  "storageKey": null
},
v3 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "type",
  "args": null,
  "storageKey": null
},
v4 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "visibleLabel",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "CatalogQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "equipmentTypes",
        "storageKey": "equipmentTypes(first:50)",
        "args": (v0/*: any*/),
        "concreteType": "EquipmentTypeConnection",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "edges",
            "storageKey": null,
            "args": null,
            "concreteType": "EquipmentTypeEdge",
            "plural": true,
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "node",
                "storageKey": null,
                "args": null,
                "concreteType": "EquipmentType",
                "plural": false,
                "selections": [
                  {
                    "kind": "FragmentSpread",
                    "name": "EquipmentTypeItem_equipmentType",
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
    "name": "CatalogQuery",
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "equipmentTypes",
        "storageKey": "equipmentTypes(first:50)",
        "args": (v0/*: any*/),
        "concreteType": "EquipmentTypeConnection",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "edges",
            "storageKey": null,
            "args": null,
            "concreteType": "EquipmentTypeEdge",
            "plural": true,
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "node",
                "storageKey": null,
                "args": null,
                "concreteType": "EquipmentType",
                "plural": false,
                "selections": [
                  (v1/*: any*/),
                  (v2/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "propertyTypes",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "PropertyType",
                    "plural": true,
                    "selections": [
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "booleanValue",
                        "args": null,
                        "storageKey": null
                      },
                      (v1/*: any*/),
                      (v3/*: any*/),
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "stringValue",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "intValue",
                        "args": null,
                        "storageKey": null
                      },
                      (v2/*: any*/),
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "floatValue",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "latitudeValue",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "longitudeValue",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "isEditable",
                        "args": null,
                        "storageKey": null
                      },
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "isInstanceProperty",
                        "args": null,
                        "storageKey": null
                      }
                    ]
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "positionDefinitions",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "EquipmentPositionDefinition",
                    "plural": true,
                    "selections": [
                      (v1/*: any*/),
                      (v2/*: any*/),
                      (v4/*: any*/)
                    ]
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "portDefinitions",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "EquipmentPortDefinition",
                    "plural": true,
                    "selections": [
                      (v1/*: any*/),
                      (v2/*: any*/),
                      (v4/*: any*/),
                      (v3/*: any*/)
                    ]
                  },
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "numberOfEquipment",
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
  },
  "params": {
    "operationKind": "query",
    "name": "CatalogQuery",
    "id": null,
    "text": "query CatalogQuery {\n  equipmentTypes(first: 50) {\n    edges {\n      node {\n        ...EquipmentTypeItem_equipmentType\n      }\n    }\n  }\n}\n\nfragment EquipmentTypeItem_equipmentType on EquipmentType {\n  id\n  name\n  propertyTypes {\n    ...PropertyTypeFormField_propertyType\n    id\n  }\n  positionDefinitions {\n    id\n    name\n    visibleLabel\n  }\n  portDefinitions {\n    id\n    name\n    visibleLabel\n    type\n  }\n  numberOfEquipment\n}\n\nfragment PropertyTypeFormField_propertyType on PropertyType {\n  id\n  name\n  type\n  stringValue\n  intValue\n  booleanValue\n  floatValue\n  latitudeValue\n  longitudeValue\n  isEditable\n  isInstanceProperty\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '2d0bc55c6a32b6c2ea7f497c67485a3d';
module.exports = node;
