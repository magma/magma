/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 05c5a9ccbed6f494dd3facf3b5eaf25f
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type AddEditServiceTypeCard_editingServiceType$ref = any;
type ServiceTypeItem_serviceType$ref = any;
export type ServiceTypesQueryVariables = {||};
export type ServiceTypesQueryResponse = {|
  +serviceTypes: ?{|
    +edges: $ReadOnlyArray<{|
      +node: ?{|
        +id: string,
        +name: string,
        +isDeleted: boolean,
        +$fragmentRefs: ServiceTypeItem_serviceType$ref & AddEditServiceTypeCard_editingServiceType$ref,
      |}
    |}>
  |}
|};
export type ServiceTypesQuery = {|
  variables: ServiceTypesQueryVariables,
  response: ServiceTypesQueryResponse,
|};
*/


/*
query ServiceTypesQuery {
  serviceTypes(first: 500) {
    edges {
      node {
        ...ServiceTypeItem_serviceType
        ...AddEditServiceTypeCard_editingServiceType
        id
        name
        isDeleted
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

fragment AddEditServiceTypeCard_editingServiceType on ServiceType {
  id
  name
  numberOfServices
  discoveryMethod
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
    rangeFromValue
    rangeToValue
    isEditable
    isMandatory
    isInstanceProperty
  }
  endpointDefinitions {
    id
    index
    role
    name
    equipmentType {
      name
      id
    }
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

fragment ServiceEndpointDefinitionStaticTable_serviceEndpointDefinitions on ServiceEndpointDefinition {
  id
  name
  role
  index
  equipmentType {
    id
    name
  }
}

fragment ServiceTypeItem_serviceType on ServiceType {
  id
  name
  discoveryMethod
  propertyTypes {
    ...PropertyTypeFormField_propertyType
    id
  }
  endpointDefinitions {
    ...ServiceEndpointDefinitionStaticTable_serviceEndpointDefinitions
    id
  }
  numberOfServices
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
  "name": "isDeleted",
  "args": null,
  "storageKey": null
},
v3 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "__typename",
  "args": null,
  "storageKey": null
},
v4 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "cursor",
  "args": null,
  "storageKey": null
},
v5 = {
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
v6 = [
  {
    "kind": "Literal",
    "name": "first",
    "value": 500
  }
],
v7 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "index",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "ServiceTypesQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "serviceTypes",
        "name": "__ServiceTypes_serviceTypes_connection",
        "storageKey": null,
        "args": null,
        "concreteType": "ServiceTypeConnection",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "edges",
            "storageKey": null,
            "args": null,
            "concreteType": "ServiceTypeEdge",
            "plural": true,
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "node",
                "storageKey": null,
                "args": null,
                "concreteType": "ServiceType",
                "plural": false,
                "selections": [
                  (v0/*: any*/),
                  (v1/*: any*/),
                  (v2/*: any*/),
                  (v3/*: any*/),
                  {
                    "kind": "FragmentSpread",
                    "name": "ServiceTypeItem_serviceType",
                    "args": null
                  },
                  {
                    "kind": "FragmentSpread",
                    "name": "AddEditServiceTypeCard_editingServiceType",
                    "args": null
                  }
                ]
              },
              (v4/*: any*/)
            ]
          },
          (v5/*: any*/)
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "ServiceTypesQuery",
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "serviceTypes",
        "storageKey": "serviceTypes(first:500)",
        "args": (v6/*: any*/),
        "concreteType": "ServiceTypeConnection",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "edges",
            "storageKey": null,
            "args": null,
            "concreteType": "ServiceTypeEdge",
            "plural": true,
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "node",
                "storageKey": null,
                "args": null,
                "concreteType": "ServiceType",
                "plural": false,
                "selections": [
                  (v0/*: any*/),
                  (v1/*: any*/),
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "discoveryMethod",
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
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "nodeType",
                        "args": null,
                        "storageKey": null
                      },
                      (v7/*: any*/),
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
                      (v2/*: any*/)
                    ]
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "endpointDefinitions",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "ServiceEndpointDefinition",
                    "plural": true,
                    "selections": [
                      (v0/*: any*/),
                      (v1/*: any*/),
                      {
                        "kind": "ScalarField",
                        "alias": null,
                        "name": "role",
                        "args": null,
                        "storageKey": null
                      },
                      (v7/*: any*/),
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "equipmentType",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "EquipmentType",
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
                    "name": "numberOfServices",
                    "args": null,
                    "storageKey": null
                  },
                  (v2/*: any*/),
                  (v3/*: any*/)
                ]
              },
              (v4/*: any*/)
            ]
          },
          (v5/*: any*/)
        ]
      },
      {
        "kind": "LinkedHandle",
        "alias": null,
        "name": "serviceTypes",
        "args": (v6/*: any*/),
        "handle": "connection",
        "key": "ServiceTypes_serviceTypes",
        "filters": null
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "ServiceTypesQuery",
    "id": null,
    "text": "query ServiceTypesQuery {\n  serviceTypes(first: 500) {\n    edges {\n      node {\n        ...ServiceTypeItem_serviceType\n        ...AddEditServiceTypeCard_editingServiceType\n        id\n        name\n        isDeleted\n        __typename\n      }\n      cursor\n    }\n    pageInfo {\n      endCursor\n      hasNextPage\n    }\n  }\n}\n\nfragment AddEditServiceTypeCard_editingServiceType on ServiceType {\n  id\n  name\n  numberOfServices\n  discoveryMethod\n  propertyTypes {\n    id\n    name\n    type\n    nodeType\n    index\n    stringValue\n    intValue\n    booleanValue\n    floatValue\n    latitudeValue\n    longitudeValue\n    rangeFromValue\n    rangeToValue\n    isEditable\n    isMandatory\n    isInstanceProperty\n  }\n  endpointDefinitions {\n    id\n    index\n    role\n    name\n    equipmentType {\n      name\n      id\n    }\n  }\n}\n\nfragment PropertyTypeFormField_propertyType on PropertyType {\n  id\n  name\n  type\n  nodeType\n  index\n  stringValue\n  intValue\n  booleanValue\n  floatValue\n  latitudeValue\n  longitudeValue\n  rangeFromValue\n  rangeToValue\n  isEditable\n  isInstanceProperty\n  isMandatory\n  category\n  isDeleted\n}\n\nfragment ServiceEndpointDefinitionStaticTable_serviceEndpointDefinitions on ServiceEndpointDefinition {\n  id\n  name\n  role\n  index\n  equipmentType {\n    id\n    name\n  }\n}\n\nfragment ServiceTypeItem_serviceType on ServiceType {\n  id\n  name\n  discoveryMethod\n  propertyTypes {\n    ...PropertyTypeFormField_propertyType\n    id\n  }\n  endpointDefinitions {\n    ...ServiceEndpointDefinitionStaticTable_serviceEndpointDefinitions\n    id\n  }\n  numberOfServices\n}\n",
    "metadata": {
      "connection": [
        {
          "count": null,
          "cursor": null,
          "direction": "forward",
          "path": [
            "serviceTypes"
          ]
        }
      ]
    }
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '27909e6e8675e10f68f441661a634c35';
module.exports = node;
