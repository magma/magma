/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash fdc776485407d71b23ef4c0861e1c02c
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type WorkOrderDetailsPane_workOrder$ref = any;
export type WorkOrderDetailsPaneQueryRendererQueryVariables = {|
  workOrderId: string
|};
export type WorkOrderDetailsPaneQueryRendererQueryResponse = {|
  +workOrder: ?{|
    +$fragmentRefs: WorkOrderDetailsPane_workOrder$ref
  |}
|};
export type WorkOrderDetailsPaneQueryRendererQuery = {|
  variables: WorkOrderDetailsPaneQueryRendererQueryVariables,
  response: WorkOrderDetailsPaneQueryRendererQueryResponse,
|};
*/


/*
query WorkOrderDetailsPaneQueryRendererQuery(
  $workOrderId: ID!
) {
  workOrder: node(id: $workOrderId) {
    __typename
    ... on WorkOrder {
      ...WorkOrderDetailsPane_workOrder
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

fragment WorkOrderDetailsPaneEquipmentItem_equipment on Equipment {
  id
  name
  equipmentType {
    id
    name
  }
  parentLocation {
    id
    name
    locationType {
      id
      name
    }
  }
  parentPosition {
    id
    definition {
      name
      visibleLabel
      id
    }
    parentEquipment {
      id
      name
    }
  }
}

fragment WorkOrderDetailsPaneLinkItem_link on Link {
  id
  futureState
  ports {
    id
    definition {
      id
      name
      visibleLabel
      portType {
        linkPropertyTypes {
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
        id
      }
    }
    parentEquipment {
      id
      name
      futureState
      equipmentType {
        id
        name
        portDefinitions {
          id
          name
          visibleLabel
          bandwidth
          portType {
            id
            name
          }
        }
      }
      ...EquipmentBreadcrumbs_equipment
    }
    serviceEndpoints {
      role
      service {
        name
        id
      }
      id
    }
  }
  workOrder {
    id
    status
  }
  properties {
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
  services {
    id
    name
  }
}

fragment WorkOrderDetailsPane_workOrder on WorkOrder {
  id
  name
  equipmentToAdd {
    id
    ...WorkOrderDetailsPaneEquipmentItem_equipment
  }
  equipmentToRemove {
    id
    ...WorkOrderDetailsPaneEquipmentItem_equipment
  }
  linksToAdd {
    id
    ...WorkOrderDetailsPaneLinkItem_link
  }
  linksToRemove {
    id
    ...WorkOrderDetailsPaneLinkItem_link
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "workOrderId",
    "type": "ID!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "Variable",
    "name": "id",
    "variableName": "workOrderId"
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
  "name": "equipmentType",
  "storageKey": null,
  "args": null,
  "concreteType": "EquipmentType",
  "plural": false,
  "selections": (v4/*: any*/)
},
v6 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "visibleLabel",
  "args": null,
  "storageKey": null
},
v7 = [
  (v2/*: any*/),
  (v3/*: any*/),
  (v5/*: any*/),
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "parentLocation",
    "storageKey": null,
    "args": null,
    "concreteType": "Location",
    "plural": false,
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
        "selections": (v4/*: any*/)
      }
    ]
  },
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "parentPosition",
    "storageKey": null,
    "args": null,
    "concreteType": "EquipmentPosition",
    "plural": false,
    "selections": [
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
          (v3/*: any*/),
          (v6/*: any*/),
          (v2/*: any*/)
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
        "selections": (v4/*: any*/)
      }
    ]
  }
],
v8 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "futureState",
  "args": null,
  "storageKey": null
},
v9 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "type",
  "args": null,
  "storageKey": null
},
v10 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "stringValue",
  "args": null,
  "storageKey": null
},
v11 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "intValue",
  "args": null,
  "storageKey": null
},
v12 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "booleanValue",
  "args": null,
  "storageKey": null
},
v13 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "floatValue",
  "args": null,
  "storageKey": null
},
v14 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "latitudeValue",
  "args": null,
  "storageKey": null
},
v15 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "longitudeValue",
  "args": null,
  "storageKey": null
},
v16 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeFromValue",
  "args": null,
  "storageKey": null
},
v17 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeToValue",
  "args": null,
  "storageKey": null
},
v18 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isEditable",
  "args": null,
  "storageKey": null
},
v19 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isInstanceProperty",
  "args": null,
  "storageKey": null
},
v20 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isMandatory",
  "args": null,
  "storageKey": null
},
v21 = [
  (v3/*: any*/),
  (v2/*: any*/)
],
v22 = [
  (v2/*: any*/),
  (v8/*: any*/),
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "ports",
    "storageKey": null,
    "args": null,
    "concreteType": "EquipmentPort",
    "plural": true,
    "selections": [
      (v2/*: any*/),
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "definition",
        "storageKey": null,
        "args": null,
        "concreteType": "EquipmentPortDefinition",
        "plural": false,
        "selections": [
          (v2/*: any*/),
          (v3/*: any*/),
          (v6/*: any*/),
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "portType",
            "storageKey": null,
            "args": null,
            "concreteType": "EquipmentPortType",
            "plural": false,
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "linkPropertyTypes",
                "storageKey": null,
                "args": null,
                "concreteType": "PropertyType",
                "plural": true,
                "selections": [
                  (v2/*: any*/),
                  (v3/*: any*/),
                  (v9/*: any*/),
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "index",
                    "args": null,
                    "storageKey": null
                  },
                  (v10/*: any*/),
                  (v11/*: any*/),
                  (v12/*: any*/),
                  (v13/*: any*/),
                  (v14/*: any*/),
                  (v15/*: any*/),
                  (v16/*: any*/),
                  (v17/*: any*/),
                  (v18/*: any*/),
                  (v19/*: any*/),
                  (v20/*: any*/)
                ]
              },
              (v2/*: any*/)
            ]
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
          (v2/*: any*/),
          (v3/*: any*/),
          (v8/*: any*/),
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "equipmentType",
            "storageKey": null,
            "args": null,
            "concreteType": "EquipmentType",
            "plural": false,
            "selections": [
              (v2/*: any*/),
              (v3/*: any*/),
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "portDefinitions",
                "storageKey": null,
                "args": null,
                "concreteType": "EquipmentPortDefinition",
                "plural": true,
                "selections": [
                  (v2/*: any*/),
                  (v3/*: any*/),
                  (v6/*: any*/),
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "bandwidth",
                    "args": null,
                    "storageKey": null
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "portType",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "EquipmentPortType",
                    "plural": false,
                    "selections": (v4/*: any*/)
                  }
                ]
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
                "selections": (v21/*: any*/)
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
                  (v6/*: any*/)
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
                  (v2/*: any*/),
                  (v3/*: any*/),
                  (v5/*: any*/)
                ]
              }
            ]
          }
        ]
      },
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "serviceEndpoints",
        "storageKey": null,
        "args": null,
        "concreteType": "ServiceEndpoint",
        "plural": true,
        "selections": [
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "role",
            "args": null,
            "storageKey": null
          },
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "service",
            "storageKey": null,
            "args": null,
            "concreteType": "Service",
            "plural": false,
            "selections": (v21/*: any*/)
          },
          (v2/*: any*/)
        ]
      }
    ]
  },
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "workOrder",
    "storageKey": null,
    "args": null,
    "concreteType": "WorkOrder",
    "plural": false,
    "selections": [
      (v2/*: any*/),
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
    "name": "properties",
    "storageKey": null,
    "args": null,
    "concreteType": "Property",
    "plural": true,
    "selections": [
      (v2/*: any*/),
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "propertyType",
        "storageKey": null,
        "args": null,
        "concreteType": "PropertyType",
        "plural": false,
        "selections": [
          (v2/*: any*/),
          (v3/*: any*/),
          (v9/*: any*/),
          (v18/*: any*/),
          (v20/*: any*/),
          (v19/*: any*/),
          (v10/*: any*/)
        ]
      },
      (v10/*: any*/),
      (v11/*: any*/),
      (v13/*: any*/),
      (v12/*: any*/),
      (v14/*: any*/),
      (v15/*: any*/),
      (v16/*: any*/),
      (v17/*: any*/),
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "equipmentValue",
        "storageKey": null,
        "args": null,
        "concreteType": "Equipment",
        "plural": false,
        "selections": (v4/*: any*/)
      },
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "locationValue",
        "storageKey": null,
        "args": null,
        "concreteType": "Location",
        "plural": false,
        "selections": (v4/*: any*/)
      },
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "serviceValue",
        "storageKey": null,
        "args": null,
        "concreteType": "Service",
        "plural": false,
        "selections": (v4/*: any*/)
      }
    ]
  },
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "services",
    "storageKey": null,
    "args": null,
    "concreteType": "Service",
    "plural": true,
    "selections": (v4/*: any*/)
  }
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "WorkOrderDetailsPaneQueryRendererQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "workOrder",
        "name": "node",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": null,
        "plural": false,
        "selections": [
          {
            "kind": "InlineFragment",
            "type": "WorkOrder",
            "selections": [
              {
                "kind": "FragmentSpread",
                "name": "WorkOrderDetailsPane_workOrder",
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
    "name": "WorkOrderDetailsPaneQueryRendererQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": "workOrder",
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
            "type": "WorkOrder",
            "selections": [
              (v3/*: any*/),
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "equipmentToAdd",
                "storageKey": null,
                "args": null,
                "concreteType": "Equipment",
                "plural": true,
                "selections": (v7/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "equipmentToRemove",
                "storageKey": null,
                "args": null,
                "concreteType": "Equipment",
                "plural": true,
                "selections": (v7/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "linksToAdd",
                "storageKey": null,
                "args": null,
                "concreteType": "Link",
                "plural": true,
                "selections": (v22/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "linksToRemove",
                "storageKey": null,
                "args": null,
                "concreteType": "Link",
                "plural": true,
                "selections": (v22/*: any*/)
              }
            ]
          }
        ]
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "WorkOrderDetailsPaneQueryRendererQuery",
    "id": null,
    "text": "query WorkOrderDetailsPaneQueryRendererQuery(\n  $workOrderId: ID!\n) {\n  workOrder: node(id: $workOrderId) {\n    __typename\n    ... on WorkOrder {\n      ...WorkOrderDetailsPane_workOrder\n    }\n    id\n  }\n}\n\nfragment EquipmentBreadcrumbs_equipment on Equipment {\n  id\n  name\n  equipmentType {\n    id\n    name\n  }\n  locationHierarchy {\n    id\n    name\n    locationType {\n      name\n      id\n    }\n  }\n  positionHierarchy {\n    id\n    definition {\n      id\n      name\n      visibleLabel\n    }\n    parentEquipment {\n      id\n      name\n      equipmentType {\n        id\n        name\n      }\n    }\n  }\n}\n\nfragment WorkOrderDetailsPaneEquipmentItem_equipment on Equipment {\n  id\n  name\n  equipmentType {\n    id\n    name\n  }\n  parentLocation {\n    id\n    name\n    locationType {\n      id\n      name\n    }\n  }\n  parentPosition {\n    id\n    definition {\n      name\n      visibleLabel\n      id\n    }\n    parentEquipment {\n      id\n      name\n    }\n  }\n}\n\nfragment WorkOrderDetailsPaneLinkItem_link on Link {\n  id\n  futureState\n  ports {\n    id\n    definition {\n      id\n      name\n      visibleLabel\n      portType {\n        linkPropertyTypes {\n          id\n          name\n          type\n          index\n          stringValue\n          intValue\n          booleanValue\n          floatValue\n          latitudeValue\n          longitudeValue\n          rangeFromValue\n          rangeToValue\n          isEditable\n          isInstanceProperty\n          isMandatory\n        }\n        id\n      }\n    }\n    parentEquipment {\n      id\n      name\n      futureState\n      equipmentType {\n        id\n        name\n        portDefinitions {\n          id\n          name\n          visibleLabel\n          bandwidth\n          portType {\n            id\n            name\n          }\n        }\n      }\n      ...EquipmentBreadcrumbs_equipment\n    }\n    serviceEndpoints {\n      role\n      service {\n        name\n        id\n      }\n      id\n    }\n  }\n  workOrder {\n    id\n    status\n  }\n  properties {\n    id\n    propertyType {\n      id\n      name\n      type\n      isEditable\n      isMandatory\n      isInstanceProperty\n      stringValue\n    }\n    stringValue\n    intValue\n    floatValue\n    booleanValue\n    latitudeValue\n    longitudeValue\n    rangeFromValue\n    rangeToValue\n    equipmentValue {\n      id\n      name\n    }\n    locationValue {\n      id\n      name\n    }\n    serviceValue {\n      id\n      name\n    }\n  }\n  services {\n    id\n    name\n  }\n}\n\nfragment WorkOrderDetailsPane_workOrder on WorkOrder {\n  id\n  name\n  equipmentToAdd {\n    id\n    ...WorkOrderDetailsPaneEquipmentItem_equipment\n  }\n  equipmentToRemove {\n    id\n    ...WorkOrderDetailsPaneEquipmentItem_equipment\n  }\n  linksToAdd {\n    id\n    ...WorkOrderDetailsPaneLinkItem_link\n  }\n  linksToRemove {\n    id\n    ...WorkOrderDetailsPaneLinkItem_link\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'cc9a79bc9155445e523632a0125fb0fc';
module.exports = node;
