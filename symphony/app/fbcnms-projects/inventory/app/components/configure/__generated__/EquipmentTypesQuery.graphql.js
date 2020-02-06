/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 659750b55c93879d245ad940025aff3f
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type AddEditEquipmentTypeCard_editingEquipmentType$ref = any;
type EquipmentTypeItem_equipmentType$ref = any;
export type EquipmentTypesQueryVariables = {||};
export type EquipmentTypesQueryResponse = {|
  +equipmentTypes: {|
    +edges: $ReadOnlyArray<{|
      +node: ?{|
        +id: string,
        +name: string,
        +$fragmentRefs: EquipmentTypeItem_equipmentType$ref & AddEditEquipmentTypeCard_editingEquipmentType$ref,
      |}
    |}>
  |}
|};
export type EquipmentTypesQuery = {|
  variables: EquipmentTypesQueryVariables,
  response: EquipmentTypesQueryResponse,
|};
*/


/*
query EquipmentTypesQuery {
  equipmentTypes(first: 500) {
    edges {
      node {
        ...EquipmentTypeItem_equipmentType
        ...AddEditEquipmentTypeCard_editingEquipmentType
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

fragment AddEditEquipmentTypeCard_editingEquipmentType on EquipmentType {
  id
  name
  propertyTypes {
    id
    name
    type
    index
    stringValue
    intValue
    booleanValue
    floatValue
    latitudeValue
    longitudeValue
    isEditable
    isInstanceProperty
    isMandatory
  }
  positionDefinitions {
    id
    name
    index
    visibleLabel
  }
  portDefinitions {
    id
    name
    index
    visibleLabel
    portType {
      id
      name
    }
  }
  numberOfEquipment
}

fragment DynamicPropertyTypesGrid_propertyTypes on PropertyType {
  ...PropertyTypeFormField_propertyType
  id
  index
}

fragment EquipmentTypeItem_equipmentType on EquipmentType {
  id
  name
  propertyTypes {
    ...DynamicPropertyTypesGrid_propertyTypes
    id
  }
  positionDefinitions {
    ...PositionDefinitionsTable_positionDefinitions
    id
  }
  portDefinitions {
    ...PortDefinitionsTable_portDefinitions
    id
  }
  numberOfEquipment
}

fragment PortDefinitionsTable_portDefinitions on EquipmentPortDefinition {
  id
  name
  index
  visibleLabel
  portType {
    id
    name
  }
}

fragment PositionDefinitionsTable_positionDefinitions on EquipmentPositionDefinition {
  id
  name
  index
  visibleLabel
}

fragment PropertyTypeFormField_propertyType on PropertyType {
  id
  name
  type
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
v6 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "index",
  "args": null,
  "storageKey": null
},
v7 = {
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
    "name": "EquipmentTypesQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "equipmentTypes",
        "name": "__EquipmentTypes_equipmentTypes_connection",
        "storageKey": null,
        "args": null,
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
                  (v0/*: any*/),
                  (v1/*: any*/),
                  (v2/*: any*/),
                  {
                    "kind": "FragmentSpread",
                    "name": "EquipmentTypeItem_equipmentType",
                    "args": null
                  },
                  {
                    "kind": "FragmentSpread",
                    "name": "AddEditEquipmentTypeCard_editingEquipmentType",
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
    "name": "EquipmentTypesQuery",
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "equipmentTypes",
        "storageKey": "equipmentTypes(first:500)",
        "args": (v5/*: any*/),
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
                  (v0/*: any*/),
                  (v1/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "propertyTypes",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "PropertyType",
                    "plural": true,
                    "selections": [
                      (v0/*: any*/),
                      (v1/*: any*/),
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "type",
                        "args": null,
                        "storageKey": null
                      },
                      (v6/*: any*/),
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
                      (v0/*: any*/),
                      (v1/*: any*/),
                      (v6/*: any*/),
                      (v7/*: any*/)
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
                      (v0/*: any*/),
                      (v1/*: any*/),
                      (v6/*: any*/),
                      (v7/*: any*/),
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "portType",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "EquipmentPortType",
                        "plural": false,
                        "selections": [
                          (v0/*: any*/),
                          (v1/*: any*/)
                        ]
                      }
                    ]
                  },
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "numberOfEquipment",
                    "args": null,
                    "storageKey": null
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
        "name": "equipmentTypes",
        "args": (v5/*: any*/),
        "handle": "connection",
        "key": "EquipmentTypes_equipmentTypes",
        "filters": null
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "EquipmentTypesQuery",
    "id": null,
    "text": "query EquipmentTypesQuery {\n  equipmentTypes(first: 500) {\n    edges {\n      node {\n        ...EquipmentTypeItem_equipmentType\n        ...AddEditEquipmentTypeCard_editingEquipmentType\n        id\n        name\n        __typename\n      }\n      cursor\n    }\n    pageInfo {\n      endCursor\n      hasNextPage\n    }\n  }\n}\n\nfragment AddEditEquipmentTypeCard_editingEquipmentType on EquipmentType {\n  id\n  name\n  propertyTypes {\n    id\n    name\n    type\n    index\n    stringValue\n    intValue\n    booleanValue\n    floatValue\n    latitudeValue\n    longitudeValue\n    isEditable\n    isInstanceProperty\n    isMandatory\n  }\n  positionDefinitions {\n    id\n    name\n    index\n    visibleLabel\n  }\n  portDefinitions {\n    id\n    name\n    index\n    visibleLabel\n    portType {\n      id\n      name\n    }\n  }\n  numberOfEquipment\n}\n\nfragment DynamicPropertyTypesGrid_propertyTypes on PropertyType {\n  ...PropertyTypeFormField_propertyType\n  id\n  index\n}\n\nfragment EquipmentTypeItem_equipmentType on EquipmentType {\n  id\n  name\n  propertyTypes {\n    ...DynamicPropertyTypesGrid_propertyTypes\n    id\n  }\n  positionDefinitions {\n    ...PositionDefinitionsTable_positionDefinitions\n    id\n  }\n  portDefinitions {\n    ...PortDefinitionsTable_portDefinitions\n    id\n  }\n  numberOfEquipment\n}\n\nfragment PortDefinitionsTable_portDefinitions on EquipmentPortDefinition {\n  id\n  name\n  index\n  visibleLabel\n  portType {\n    id\n    name\n  }\n}\n\nfragment PositionDefinitionsTable_positionDefinitions on EquipmentPositionDefinition {\n  id\n  name\n  index\n  visibleLabel\n}\n\nfragment PropertyTypeFormField_propertyType on PropertyType {\n  id\n  name\n  type\n  index\n  stringValue\n  intValue\n  booleanValue\n  floatValue\n  latitudeValue\n  longitudeValue\n  rangeFromValue\n  rangeToValue\n  isEditable\n  isInstanceProperty\n  isMandatory\n}\n",
    "metadata": {
      "connection": [
        {
          "count": null,
          "cursor": null,
          "direction": "forward",
          "path": [
            "equipmentTypes"
          ]
        }
      ]
    }
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '22440e25d923db73f49ea30630898117';
module.exports = node;
