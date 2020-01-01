/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash e4a1d9e2441cd333833539d1b2226b67
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type ServicesView_service$ref = any;
export type FilterOperator = "CONTAINS" | "DATE_GREATER_THAN" | "DATE_LESS_THAN" | "IS" | "IS_NOT_ONE_OF" | "IS_ONE_OF" | "%future added value";
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "equipment" | "float" | "gps_location" | "int" | "location" | "range" | "service" | "string" | "%future added value";
export type ServiceFilterType = "EQUIPMENT_IN_SERVICE" | "LOCATION_INST" | "SERVICE_INST_CUSTOMER_NAME" | "SERVICE_INST_EXTERNAL_ID" | "SERVICE_INST_NAME" | "SERVICE_INST_PROPERTY" | "SERVICE_STATUS" | "SERVICE_TYPE" | "%future added value";
export type ServiceFilterInput = {|
  filterType: ServiceFilterType,
  operator: FilterOperator,
  stringValue?: ?string,
  propertyValue?: ?PropertyTypeInput,
  idSet?: ?$ReadOnlyArray<string>,
  maxDepth?: ?number,
|};
export type PropertyTypeInput = {|
  id?: ?string,
  name: string,
  type: PropertyKind,
  index?: ?number,
  category?: ?string,
  stringValue?: ?string,
  intValue?: ?number,
  booleanValue?: ?boolean,
  floatValue?: ?number,
  latitudeValue?: ?number,
  longitudeValue?: ?number,
  rangeFromValue?: ?number,
  rangeToValue?: ?number,
  isEditable?: ?boolean,
  isInstanceProperty?: ?boolean,
  isMandatory?: ?boolean,
  isDeleted?: ?boolean,
|};
export type ServiceComparisonViewQueryRendererSearchQueryVariables = {|
  limit?: ?number,
  filters: $ReadOnlyArray<ServiceFilterInput>,
|};
export type ServiceComparisonViewQueryRendererSearchQueryResponse = {|
  +serviceSearch: {|
    +services: $ReadOnlyArray<?{|
      +$fragmentRefs: ServicesView_service$ref
    |}>,
    +count: number,
  |}
|};
export type ServiceComparisonViewQueryRendererSearchQuery = {|
  variables: ServiceComparisonViewQueryRendererSearchQueryVariables,
  response: ServiceComparisonViewQueryRendererSearchQueryResponse,
|};
*/


/*
query ServiceComparisonViewQueryRendererSearchQuery(
  $limit: Int
  $filters: [ServiceFilterInput!]!
) {
  serviceSearch(limit: $limit, filters: $filters) {
    services {
      ...ServicesView_service
      id
    }
    count
  }
}

fragment DynamicPropertiesGrid_properties on Property {
  ...PropertyFormField_property
  propertyType {
    id
    index
  }
}

fragment DynamicPropertiesGrid_propertyTypes on PropertyType {
  id
  name
  index
  isInstanceProperty
  type
  stringValue
  intValue
  booleanValue
  latitudeValue
  longitudeValue
  rangeFromValue
  rangeToValue
  floatValue
}

fragment PropertyFormField_property on Property {
  id
  propertyType {
    id
    name
    type
    isEditable
    isMandatory
    isInstanceProperty
    stringValue
  }
  stringValue
  intValue
  floatValue
  booleanValue
  latitudeValue
  longitudeValue
  rangeFromValue
  rangeToValue
  equipmentValue {
    id
    name
  }
  locationValue {
    id
    name
  }
  serviceValue {
    id
    name
  }
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

fragment ServicesView_service on Service {
  id
  name
  externalId
  status
  customer {
    id
    name
  }
  serviceType {
    id
    name
    propertyTypes {
      ...PropertyTypeFormField_propertyType
      ...DynamicPropertiesGrid_propertyTypes
      id
    }
  }
  properties {
    ...PropertyFormField_property
    ...DynamicPropertiesGrid_properties
    id
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "limit",
    "type": "Int",
    "defaultValue": null
  },
  {
    "kind": "LocalArgument",
    "name": "filters",
    "type": "[ServiceFilterInput!]!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "Variable",
    "name": "filters",
    "variableName": "filters"
  },
  {
    "kind": "Variable",
    "name": "limit",
    "variableName": "limit"
  }
],
v2 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "count",
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
  "name": "name",
  "args": null,
  "storageKey": null
},
v5 = [
  (v3/*: any*/),
  (v4/*: any*/)
],
v6 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "type",
  "args": null,
  "storageKey": null
},
v7 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "index",
  "args": null,
  "storageKey": null
},
v8 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "stringValue",
  "args": null,
  "storageKey": null
},
v9 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "intValue",
  "args": null,
  "storageKey": null
},
v10 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "booleanValue",
  "args": null,
  "storageKey": null
},
v11 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "floatValue",
  "args": null,
  "storageKey": null
},
v12 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "latitudeValue",
  "args": null,
  "storageKey": null
},
v13 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "longitudeValue",
  "args": null,
  "storageKey": null
},
v14 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeFromValue",
  "args": null,
  "storageKey": null
},
v15 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeToValue",
  "args": null,
  "storageKey": null
},
v16 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isEditable",
  "args": null,
  "storageKey": null
},
v17 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isInstanceProperty",
  "args": null,
  "storageKey": null
},
v18 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isMandatory",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "ServiceComparisonViewQueryRendererSearchQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "serviceSearch",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "ServiceSearchResult",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "services",
            "storageKey": null,
            "args": null,
            "concreteType": "Service",
            "plural": true,
            "selections": [
              {
                "kind": "FragmentSpread",
                "name": "ServicesView_service",
                "args": null
              }
            ]
          },
          (v2/*: any*/)
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "ServiceComparisonViewQueryRendererSearchQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "serviceSearch",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "ServiceSearchResult",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "services",
            "storageKey": null,
            "args": null,
            "concreteType": "Service",
            "plural": true,
            "selections": [
              (v3/*: any*/),
              (v4/*: any*/),
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
                "selections": (v5/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "serviceType",
                "storageKey": null,
                "args": null,
                "concreteType": "ServiceType",
                "plural": false,
                "selections": [
                  (v3/*: any*/),
                  (v4/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "propertyTypes",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "PropertyType",
                    "plural": true,
                    "selections": [
                      (v3/*: any*/),
                      (v4/*: any*/),
                      (v6/*: any*/),
                      (v7/*: any*/),
                      (v8/*: any*/),
                      (v9/*: any*/),
                      (v10/*: any*/),
                      (v11/*: any*/),
                      (v12/*: any*/),
                      (v13/*: any*/),
                      (v14/*: any*/),
                      (v15/*: any*/),
                      (v16/*: any*/),
                      (v17/*: any*/),
                      (v18/*: any*/)
                    ]
                  }
                ]
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "properties",
                "storageKey": null,
                "args": null,
                "concreteType": "Property",
                "plural": true,
                "selections": [
                  (v3/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "propertyType",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "PropertyType",
                    "plural": false,
                    "selections": [
                      (v3/*: any*/),
                      (v4/*: any*/),
                      (v6/*: any*/),
                      (v16/*: any*/),
                      (v18/*: any*/),
                      (v17/*: any*/),
                      (v8/*: any*/),
                      (v7/*: any*/)
                    ]
                  },
                  (v8/*: any*/),
                  (v9/*: any*/),
                  (v11/*: any*/),
                  (v10/*: any*/),
                  (v12/*: any*/),
                  (v13/*: any*/),
                  (v14/*: any*/),
                  (v15/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "equipmentValue",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "Equipment",
                    "plural": false,
                    "selections": (v5/*: any*/)
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "locationValue",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "Location",
                    "plural": false,
                    "selections": (v5/*: any*/)
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "serviceValue",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "Service",
                    "plural": false,
                    "selections": (v5/*: any*/)
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
    "name": "ServiceComparisonViewQueryRendererSearchQuery",
    "id": null,
    "text": "query ServiceComparisonViewQueryRendererSearchQuery(\n  $limit: Int\n  $filters: [ServiceFilterInput!]!\n) {\n  serviceSearch(limit: $limit, filters: $filters) {\n    services {\n      ...ServicesView_service\n      id\n    }\n    count\n  }\n}\n\nfragment DynamicPropertiesGrid_properties on Property {\n  ...PropertyFormField_property\n  propertyType {\n    id\n    index\n  }\n}\n\nfragment DynamicPropertiesGrid_propertyTypes on PropertyType {\n  id\n  name\n  index\n  isInstanceProperty\n  type\n  stringValue\n  intValue\n  booleanValue\n  latitudeValue\n  longitudeValue\n  rangeFromValue\n  rangeToValue\n  floatValue\n}\n\nfragment PropertyFormField_property on Property {\n  id\n  propertyType {\n    id\n    name\n    type\n    isEditable\n    isMandatory\n    isInstanceProperty\n    stringValue\n  }\n  stringValue\n  intValue\n  floatValue\n  booleanValue\n  latitudeValue\n  longitudeValue\n  rangeFromValue\n  rangeToValue\n  equipmentValue {\n    id\n    name\n  }\n  locationValue {\n    id\n    name\n  }\n  serviceValue {\n    id\n    name\n  }\n}\n\nfragment PropertyTypeFormField_propertyType on PropertyType {\n  id\n  name\n  type\n  index\n  stringValue\n  intValue\n  booleanValue\n  floatValue\n  latitudeValue\n  longitudeValue\n  rangeFromValue\n  rangeToValue\n  isEditable\n  isInstanceProperty\n  isMandatory\n}\n\nfragment ServicesView_service on Service {\n  id\n  name\n  externalId\n  status\n  customer {\n    id\n    name\n  }\n  serviceType {\n    id\n    name\n    propertyTypes {\n      ...PropertyTypeFormField_propertyType\n      ...DynamicPropertiesGrid_propertyTypes\n      id\n    }\n  }\n  properties {\n    ...PropertyFormField_property\n    ...DynamicPropertiesGrid_properties\n    id\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'ed58dc0dacf235000888b2776b58fd08';
module.exports = node;
