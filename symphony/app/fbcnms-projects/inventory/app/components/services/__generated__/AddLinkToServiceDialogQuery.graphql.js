/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash b74e5f3e769bbd5ee46c1d78c014f147
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type AvailableLinksTable_links$ref = any;
export type FilterOperator = "CONTAINS" | "DATE_GREATER_OR_EQUAL_THAN" | "DATE_GREATER_THAN" | "DATE_LESS_OR_EQUAL_THAN" | "DATE_LESS_THAN" | "IS" | "IS_NOT_ONE_OF" | "IS_ONE_OF" | "%future added value";
export type LinkFilterType = "EQUIPMENT_INST" | "EQUIPMENT_TYPE" | "LINK_FUTURE_STATUS" | "LOCATION_INST" | "LOCATION_INST_EXTERNAL_ID" | "PROPERTY" | "SERVICE_INST" | "%future added value";
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "float" | "gps_location" | "int" | "node" | "range" | "string" | "%future added value";
export type LinkFilterInput = {|
  filterType: LinkFilterType,
  operator: FilterOperator,
  stringValue?: ?string,
  propertyValue?: ?PropertyTypeInput,
  idSet?: ?$ReadOnlyArray<string>,
  stringSet?: ?$ReadOnlyArray<string>,
  maxDepth?: ?number,
|};
export type PropertyTypeInput = {|
  id?: ?string,
  externalId?: ?string,
  name: string,
  type: PropertyKind,
  nodeType?: ?string,
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
export type AddLinkToServiceDialogQueryVariables = {|
  filters: $ReadOnlyArray<LinkFilterInput>
|};
export type AddLinkToServiceDialogQueryResponse = {|
  +links: {|
    +edges: $ReadOnlyArray<{|
      +node: ?{|
        +id: string,
        +ports: $ReadOnlyArray<?{|
          +parentEquipment: {|
            +id: string,
            +name: string,
          |},
          +definition: {|
            +id: string,
            +name: string,
          |},
        |}>,
        +$fragmentRefs: AvailableLinksTable_links$ref,
      |}
    |}>
  |}
|};
export type AddLinkToServiceDialogQuery = {|
  variables: AddLinkToServiceDialogQueryVariables,
  response: AddLinkToServiceDialogQueryResponse,
|};
*/


/*
query AddLinkToServiceDialogQuery(
  $filters: [LinkFilterInput!]!
) {
  links(filterBy: $filters, first: 50) {
    edges {
      node {
        id
        ports {
          parentEquipment {
            id
            name
          }
          definition {
            id
            name
          }
          id
        }
        ...AvailableLinksTable_links
      }
    }
  }
}

fragment AvailableLinksTable_links on Link {
  id
  ports {
    parentEquipment {
      id
      name
      positionHierarchy {
        parentEquipment {
          id
        }
        id
      }
      ...EquipmentBreadcrumbs_equipment
    }
    definition {
      id
      name
    }
    id
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
    "name": "filters",
    "type": "[LinkFilterInput!]!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "Variable",
    "name": "filterBy",
    "variableName": "filters"
  },
  {
    "kind": "Literal",
    "name": "first",
    "value": 50
  }
],
v2 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "id",
  "args": null,
  "storageKey": null
},
v3 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "name",
  "args": null,
  "storageKey": null
},
v4 = [
  (v2/*: any*/),
  (v3/*: any*/)
],
v5 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "definition",
  "storageKey": null,
  "args": null,
  "concreteType": "EquipmentPortDefinition",
  "plural": false,
  "selections": (v4/*: any*/)
},
v6 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "equipmentType",
  "storageKey": null,
  "args": null,
  "concreteType": "EquipmentType",
  "plural": false,
  "selections": (v4/*: any*/)
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "AddLinkToServiceDialogQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "links",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "LinkConnection",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "edges",
            "storageKey": null,
            "args": null,
            "concreteType": "LinkEdge",
            "plural": true,
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "node",
                "storageKey": null,
                "args": null,
                "concreteType": "Link",
                "plural": false,
                "selections": [
                  (v2/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "ports",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "EquipmentPort",
                    "plural": true,
                    "selections": [
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "parentEquipment",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "Equipment",
                        "plural": false,
                        "selections": (v4/*: any*/)
                      },
                      (v5/*: any*/)
                    ]
                  },
                  {
                    "kind": "FragmentSpread",
                    "name": "AvailableLinksTable_links",
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
    "name": "AddLinkToServiceDialogQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "links",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "LinkConnection",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "edges",
            "storageKey": null,
            "args": null,
            "concreteType": "LinkEdge",
            "plural": true,
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "node",
                "storageKey": null,
                "args": null,
                "concreteType": "Link",
                "plural": false,
                "selections": [
                  (v2/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "ports",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "EquipmentPort",
                    "plural": true,
                    "selections": [
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "parentEquipment",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "Equipment",
                        "plural": false,
                        "selections": [
                          (v2/*: any*/),
                          (v3/*: any*/),
                          {
                            "kind": "LinkedField",
                            "alias": null,
                            "name": "positionHierarchy",
                            "storageKey": null,
                            "args": null,
                            "concreteType": "EquipmentPosition",
                            "plural": true,
                            "selections": [
                              {
                                "kind": "LinkedField",
                                "alias": null,
                                "name": "parentEquipment",
                                "storageKey": null,
                                "args": null,
                                "concreteType": "Equipment",
                                "plural": false,
                                "selections": [
                                  (v2/*: any*/),
                                  (v3/*: any*/),
                                  (v6/*: any*/)
                                ]
                              },
                              (v2/*: any*/),
                              {
                                "kind": "LinkedField",
                                "alias": null,
                                "name": "definition",
                                "storageKey": null,
                                "args": null,
                                "concreteType": "EquipmentPositionDefinition",
                                "plural": false,
                                "selections": [
                                  (v2/*: any*/),
                                  (v3/*: any*/),
                                  {
                                    "kind": "ScalarField",
                                    "alias": null,
                                    "name": "visibleLabel",
                                    "args": null,
                                    "storageKey": null
                                  }
                                ]
                              }
                            ]
                          },
                          (v6/*: any*/),
                          {
                            "kind": "LinkedField",
                            "alias": null,
                            "name": "locationHierarchy",
                            "storageKey": null,
                            "args": null,
                            "concreteType": "Location",
                            "plural": true,
                            "selections": [
                              (v2/*: any*/),
                              (v3/*: any*/),
                              {
                                "kind": "LinkedField",
                                "alias": null,
                                "name": "locationType",
                                "storageKey": null,
                                "args": null,
                                "concreteType": "LocationType",
                                "plural": false,
                                "selections": [
                                  (v3/*: any*/),
                                  (v2/*: any*/)
                                ]
                              }
                            ]
                          }
                        ]
                      },
                      (v5/*: any*/),
                      (v2/*: any*/)
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
    "name": "AddLinkToServiceDialogQuery",
    "id": null,
    "text": "query AddLinkToServiceDialogQuery(\n  $filters: [LinkFilterInput!]!\n) {\n  links(filterBy: $filters, first: 50) {\n    edges {\n      node {\n        id\n        ports {\n          parentEquipment {\n            id\n            name\n          }\n          definition {\n            id\n            name\n          }\n          id\n        }\n        ...AvailableLinksTable_links\n      }\n    }\n  }\n}\n\nfragment AvailableLinksTable_links on Link {\n  id\n  ports {\n    parentEquipment {\n      id\n      name\n      positionHierarchy {\n        parentEquipment {\n          id\n        }\n        id\n      }\n      ...EquipmentBreadcrumbs_equipment\n    }\n    definition {\n      id\n      name\n    }\n    id\n  }\n}\n\nfragment EquipmentBreadcrumbs_equipment on Equipment {\n  id\n  name\n  equipmentType {\n    id\n    name\n  }\n  locationHierarchy {\n    id\n    name\n    locationType {\n      name\n      id\n    }\n  }\n  positionHierarchy {\n    id\n    definition {\n      id\n      name\n      visibleLabel\n    }\n    parentEquipment {\n      id\n      name\n      equipmentType {\n        id\n        name\n      }\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'fb869c69fa4c5764072912438439fdf6';
module.exports = node;
