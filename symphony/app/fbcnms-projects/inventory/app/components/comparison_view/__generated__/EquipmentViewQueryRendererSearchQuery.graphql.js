/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 4ffe84d7eb69a805485f15e8582a1552
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type PowerSearchEquipmentResultsTable_equipment$ref = any;
export type EquipmentFilterType = "EQUIPMENT_TYPE" | "EQUIP_INST_NAME" | "LOCATION_INST" | "PROPERTY" | "%future added value";
export type FilterOperator = "CONTAINS" | "DATE_GREATER_THAN" | "DATE_LESS_THAN" | "IS" | "IS_NOT_ONE_OF" | "IS_ONE_OF" | "%future added value";
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "equipment" | "float" | "gps_location" | "int" | "location" | "range" | "service" | "string" | "%future added value";
export type EquipmentFilterInput = {|
  filterType: EquipmentFilterType,
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
export type EquipmentViewQueryRendererSearchQueryVariables = {|
  limit?: ?number,
  filters: $ReadOnlyArray<EquipmentFilterInput>,
|};
export type EquipmentViewQueryRendererSearchQueryResponse = {|
  +equipmentSearch: {|
    +equipment: $ReadOnlyArray<?{|
      +$fragmentRefs: PowerSearchEquipmentResultsTable_equipment$ref
    |}>,
    +count: number,
  |}
|};
export type EquipmentViewQueryRendererSearchQuery = {|
  variables: EquipmentViewQueryRendererSearchQueryVariables,
  response: EquipmentViewQueryRendererSearchQueryResponse,
|};
*/


/*
query EquipmentViewQueryRendererSearchQuery(
  $limit: Int
  $filters: [EquipmentFilterInput!]!
) {
  equipmentSearch(limit: $limit, filters: $filters) {
    equipment {
      ...PowerSearchEquipmentResultsTable_equipment
      id
    }
    count
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

fragment PowerSearchEquipmentResultsTable_equipment on Equipment {
  id
  name
  futureState
  externalId
  equipmentType {
    id
    name
  }
  workOrder {
    id
    status
  }
  ...EquipmentBreadcrumbs_equipment
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
    "type": "[EquipmentFilterInput!]!",
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
v5 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "equipmentType",
  "storageKey": null,
  "args": null,
  "concreteType": "EquipmentType",
  "plural": false,
  "selections": [
    (v3/*: any*/),
    (v4/*: any*/)
  ]
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "EquipmentViewQueryRendererSearchQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "equipmentSearch",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "EquipmentSearchResult",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "equipment",
            "storageKey": null,
            "args": null,
            "concreteType": "Equipment",
            "plural": true,
            "selections": [
              {
                "kind": "FragmentSpread",
                "name": "PowerSearchEquipmentResultsTable_equipment",
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
    "name": "EquipmentViewQueryRendererSearchQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "equipmentSearch",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "EquipmentSearchResult",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "equipment",
            "storageKey": null,
            "args": null,
            "concreteType": "Equipment",
            "plural": true,
            "selections": [
              (v3/*: any*/),
              (v4/*: any*/),
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "futureState",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "externalId",
                "args": null,
                "storageKey": null
              },
              (v5/*: any*/),
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "workOrder",
                "storageKey": null,
                "args": null,
                "concreteType": "WorkOrder",
                "plural": false,
                "selections": [
                  (v3/*: any*/),
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "status",
                    "args": null,
                    "storageKey": null
                  }
                ]
              },
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
                  (v4/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "locationType",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "LocationType",
                    "plural": false,
                    "selections": [
                      (v4/*: any*/),
                      (v3/*: any*/)
                    ]
                  }
                ]
              },
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
                      (v4/*: any*/),
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
                      (v4/*: any*/),
                      (v5/*: any*/)
                    ]
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
    "name": "EquipmentViewQueryRendererSearchQuery",
    "id": null,
    "text": "query EquipmentViewQueryRendererSearchQuery(\n  $limit: Int\n  $filters: [EquipmentFilterInput!]!\n) {\n  equipmentSearch(limit: $limit, filters: $filters) {\n    equipment {\n      ...PowerSearchEquipmentResultsTable_equipment\n      id\n    }\n    count\n  }\n}\n\nfragment EquipmentBreadcrumbs_equipment on Equipment {\n  id\n  name\n  equipmentType {\n    id\n    name\n  }\n  locationHierarchy {\n    id\n    name\n    locationType {\n      name\n      id\n    }\n  }\n  positionHierarchy {\n    id\n    definition {\n      id\n      name\n      visibleLabel\n    }\n    parentEquipment {\n      id\n      name\n      equipmentType {\n        id\n        name\n      }\n    }\n  }\n}\n\nfragment PowerSearchEquipmentResultsTable_equipment on Equipment {\n  id\n  name\n  futureState\n  externalId\n  equipmentType {\n    id\n    name\n  }\n  workOrder {\n    id\n    status\n  }\n  ...EquipmentBreadcrumbs_equipment\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '67c29f4d532c39ff124d5f461d58abdb';
module.exports = node;
