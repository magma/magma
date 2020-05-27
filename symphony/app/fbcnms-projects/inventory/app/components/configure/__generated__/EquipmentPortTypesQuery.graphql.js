/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash e88121962f3b9d96d60b8c4e0ccced25
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type AddEditEquipmentPortTypeCard_editingEquipmentPortType$ref = any;
type EquipmentPortTypeItem_equipmentPortType$ref = any;
export type EquipmentPortTypesQueryVariables = {||};
export type EquipmentPortTypesQueryResponse = {|
  +equipmentPortTypes: {|
    +edges: $ReadOnlyArray<{|
      +node: ?{|
        +id: string,
        +name: string,
        +$fragmentRefs: EquipmentPortTypeItem_equipmentPortType$ref & AddEditEquipmentPortTypeCard_editingEquipmentPortType$ref,
      |}
    |}>
  |}
|};
export type EquipmentPortTypesQuery = {|
  variables: EquipmentPortTypesQueryVariables,
  response: EquipmentPortTypesQueryResponse,
|};
*/


/*
query EquipmentPortTypesQuery {
  equipmentPortTypes(first: 500) {
    edges {
      node {
        ...EquipmentPortTypeItem_equipmentPortType
        ...AddEditEquipmentPortTypeCard_editingEquipmentPortType
        id
        name
        __typename
      }
      cursor
    }
    pageInfo {
      endCursor
      hasNextPage
    }
  }
}

fragment AddEditEquipmentPortTypeCard_editingEquipmentPortType on EquipmentPortType {
  id
  name
  numberOfPortDefinitions
  propertyTypes {
    id
    name
    type
    nodeType
    index
    stringValue
    intValue
    booleanValue
    floatValue
    latitudeValue
    longitudeValue
    isEditable
    isInstanceProperty
  }
  linkPropertyTypes {
    id
    name
    type
    nodeType
    index
    stringValue
    intValue
    booleanValue
    floatValue
    latitudeValue
    longitudeValue
    isEditable
    isInstanceProperty
  }
}

fragment DynamicPropertyTypesGrid_propertyTypes on PropertyType {
  ...PropertyTypeFormField_propertyType
  id
  index
}

fragment EquipmentPortTypeItem_equipmentPortType on EquipmentPortType {
  id
  name
  numberOfPortDefinitions
  propertyTypes {
    ...DynamicPropertyTypesGrid_propertyTypes
    id
  }
  linkPropertyTypes {
    ...DynamicPropertyTypesGrid_propertyTypes
    id
  }
}

fragment PropertyTypeFormField_propertyType on PropertyType {
  id
  name
  type
  nodeType
  index
  stringValue
  intValue
  booleanValue
  floatValue
  latitudeValue
  longitudeValue
  rangeFromValue
  rangeToValue
  isEditable
  isInstanceProperty
  isMandatory
  category
  isDeleted
}
*/

const node/*: ConcreteRequest*/ = (function(){
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
  "name": "cursor",
  "args": null,
  "storageKey": null
},
v4 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "pageInfo",
  "storageKey": null,
  "args": null,
  "concreteType": "PageInfo",
  "plural": false,
  "selections": [
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "endCursor",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "hasNextPage",
      "args": null,
      "storageKey": null
    }
  ]
},
v5 = [
  {
    "kind": "Literal",
    "name": "first",
    "value": 500
  }
],
v6 = [
  (v0/*: any*/),
  (v1/*: any*/),
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "type",
    "args": null,
    "storageKey": null
  },
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "nodeType",
    "args": null,
    "storageKey": null
  },
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
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "booleanValue",
    "args": null,
    "storageKey": null
  },
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
    "name": "rangeFromValue",
    "args": null,
    "storageKey": null
  },
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "rangeToValue",
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
  },
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "isMandatory",
    "args": null,
    "storageKey": null
  },
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "category",
    "args": null,
    "storageKey": null
  },
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "isDeleted",
    "args": null,
    "storageKey": null
  }
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "EquipmentPortTypesQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "equipmentPortTypes",
        "name": "__EquipmentPortTypes_equipmentPortTypes_connection",
        "storageKey": null,
        "args": null,
        "concreteType": "EquipmentPortTypeConnection",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "edges",
            "storageKey": null,
            "args": null,
            "concreteType": "EquipmentPortTypeEdge",
            "plural": true,
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "node",
                "storageKey": null,
                "args": null,
                "concreteType": "EquipmentPortType",
                "plural": false,
                "selections": [
                  (v0/*: any*/),
                  (v1/*: any*/),
                  (v2/*: any*/),
                  {
                    "kind": "FragmentSpread",
                    "name": "EquipmentPortTypeItem_equipmentPortType",
                    "args": null
                  },
                  {
                    "kind": "FragmentSpread",
                    "name": "AddEditEquipmentPortTypeCard_editingEquipmentPortType",
                    "args": null
                  }
                ]
              },
              (v3/*: any*/)
            ]
          },
          (v4/*: any*/)
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "EquipmentPortTypesQuery",
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "equipmentPortTypes",
        "storageKey": "equipmentPortTypes(first:500)",
        "args": (v5/*: any*/),
        "concreteType": "EquipmentPortTypeConnection",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "edges",
            "storageKey": null,
            "args": null,
            "concreteType": "EquipmentPortTypeEdge",
            "plural": true,
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "node",
                "storageKey": null,
                "args": null,
                "concreteType": "EquipmentPortType",
                "plural": false,
                "selections": [
                  (v0/*: any*/),
                  (v1/*: any*/),
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "numberOfPortDefinitions",
                    "args": null,
                    "storageKey": null
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "propertyTypes",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "PropertyType",
                    "plural": true,
                    "selections": (v6/*: any*/)
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "linkPropertyTypes",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "PropertyType",
                    "plural": true,
                    "selections": (v6/*: any*/)
                  },
                  (v2/*: any*/)
                ]
              },
              (v3/*: any*/)
            ]
          },
          (v4/*: any*/)
        ]
      },
      {
        "kind": "LinkedHandle",
        "alias": null,
        "name": "equipmentPortTypes",
        "args": (v5/*: any*/),
        "handle": "connection",
        "key": "EquipmentPortTypes_equipmentPortTypes",
        "filters": null
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "EquipmentPortTypesQuery",
    "id": null,
    "text": "query EquipmentPortTypesQuery {\n  equipmentPortTypes(first: 500) {\n    edges {\n      node {\n        ...EquipmentPortTypeItem_equipmentPortType\n        ...AddEditEquipmentPortTypeCard_editingEquipmentPortType\n        id\n        name\n        __typename\n      }\n      cursor\n    }\n    pageInfo {\n      endCursor\n      hasNextPage\n    }\n  }\n}\n\nfragment AddEditEquipmentPortTypeCard_editingEquipmentPortType on EquipmentPortType {\n  id\n  name\n  numberOfPortDefinitions\n  propertyTypes {\n    id\n    name\n    type\n    nodeType\n    index\n    stringValue\n    intValue\n    booleanValue\n    floatValue\n    latitudeValue\n    longitudeValue\n    isEditable\n    isInstanceProperty\n  }\n  linkPropertyTypes {\n    id\n    name\n    type\n    nodeType\n    index\n    stringValue\n    intValue\n    booleanValue\n    floatValue\n    latitudeValue\n    longitudeValue\n    isEditable\n    isInstanceProperty\n  }\n}\n\nfragment DynamicPropertyTypesGrid_propertyTypes on PropertyType {\n  ...PropertyTypeFormField_propertyType\n  id\n  index\n}\n\nfragment EquipmentPortTypeItem_equipmentPortType on EquipmentPortType {\n  id\n  name\n  numberOfPortDefinitions\n  propertyTypes {\n    ...DynamicPropertyTypesGrid_propertyTypes\n    id\n  }\n  linkPropertyTypes {\n    ...DynamicPropertyTypesGrid_propertyTypes\n    id\n  }\n}\n\nfragment PropertyTypeFormField_propertyType on PropertyType {\n  id\n  name\n  type\n  nodeType\n  index\n  stringValue\n  intValue\n  booleanValue\n  floatValue\n  latitudeValue\n  longitudeValue\n  rangeFromValue\n  rangeToValue\n  isEditable\n  isInstanceProperty\n  isMandatory\n  category\n  isDeleted\n}\n",
    "metadata": {
      "connection": [
        {
          "count": null,
          "cursor": null,
          "direction": "forward",
          "path": [
            "equipmentPortTypes"
          ]
        }
      ]
    }
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'bbfe2ae6972634e9f7fc0bdabbb613f7';
module.exports = node;
