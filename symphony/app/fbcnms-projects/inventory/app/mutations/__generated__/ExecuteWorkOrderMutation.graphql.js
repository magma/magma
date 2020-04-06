/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 4083779e13415f314473b6566ef13ca3
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type EquipmentPortsTable_link$ref = any;
type EquipmentTable_equipment$ref = any;
export type ExecuteWorkOrderMutationVariables = {|
  id: string
|};
export type ExecuteWorkOrderMutationResponse = {|
  +executeWorkOrder: {|
    +equipmentAdded: $ReadOnlyArray<{|
      +$fragmentRefs: EquipmentTable_equipment$ref
    |}>,
    +equipmentRemoved: $ReadOnlyArray<string>,
    +linkAdded: $ReadOnlyArray<{|
      +$fragmentRefs: EquipmentPortsTable_link$ref
    |}>,
    +linkRemoved: $ReadOnlyArray<string>,
  |}
|};
export type ExecuteWorkOrderMutation = {|
  variables: ExecuteWorkOrderMutationVariables,
  response: ExecuteWorkOrderMutationResponse,
|};
*/


/*
mutation ExecuteWorkOrderMutation(
  $id: ID!
) {
  executeWorkOrder(id: $id) {
    equipmentAdded {
      ...EquipmentTable_equipment
      id
    }
    equipmentRemoved
    linkAdded {
      ...EquipmentPortsTable_link
      id
    }
    linkRemoved
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

fragment EquipmentPortsTable_link on Link {
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

fragment EquipmentTable_equipment on Equipment {
  id
  name
  futureState
  equipmentType {
    id
    name
  }
  workOrder {
    id
    status
  }
  device {
    up
  }
  services {
    id
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "id",
    "type": "ID!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "Variable",
    "name": "id",
    "variableName": "id"
  }
],
v2 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "equipmentRemoved",
  "args": null,
  "storageKey": null
},
v3 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "linkRemoved",
  "args": null,
  "storageKey": null
},
v4 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "id",
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
v6 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "futureState",
  "args": null,
  "storageKey": null
},
v7 = [
  (v4/*: any*/),
  (v5/*: any*/)
],
v8 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "equipmentType",
  "storageKey": null,
  "args": null,
  "concreteType": "EquipmentType",
  "plural": false,
  "selections": (v7/*: any*/)
},
v9 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "workOrder",
  "storageKey": null,
  "args": null,
  "concreteType": "WorkOrder",
  "plural": false,
  "selections": [
    (v4/*: any*/),
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "status",
      "args": null,
      "storageKey": null
    }
  ]
},
v10 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "visibleLabel",
  "args": null,
  "storageKey": null
},
v11 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "type",
  "args": null,
  "storageKey": null
},
v12 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "stringValue",
  "args": null,
  "storageKey": null
},
v13 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "intValue",
  "args": null,
  "storageKey": null
},
v14 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "booleanValue",
  "args": null,
  "storageKey": null
},
v15 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "floatValue",
  "args": null,
  "storageKey": null
},
v16 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "latitudeValue",
  "args": null,
  "storageKey": null
},
v17 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "longitudeValue",
  "args": null,
  "storageKey": null
},
v18 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeFromValue",
  "args": null,
  "storageKey": null
},
v19 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeToValue",
  "args": null,
  "storageKey": null
},
v20 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isEditable",
  "args": null,
  "storageKey": null
},
v21 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isInstanceProperty",
  "args": null,
  "storageKey": null
},
v22 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isMandatory",
  "args": null,
  "storageKey": null
},
v23 = [
  (v5/*: any*/),
  (v4/*: any*/)
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "ExecuteWorkOrderMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "executeWorkOrder",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "WorkOrderExecutionResult",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "equipmentAdded",
            "storageKey": null,
            "args": null,
            "concreteType": "Equipment",
            "plural": true,
            "selections": [
              {
                "kind": "FragmentSpread",
                "name": "EquipmentTable_equipment",
                "args": null
              }
            ]
          },
          (v2/*: any*/),
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "linkAdded",
            "storageKey": null,
            "args": null,
            "concreteType": "Link",
            "plural": true,
            "selections": [
              {
                "kind": "FragmentSpread",
                "name": "EquipmentPortsTable_link",
                "args": null
              }
            ]
          },
          (v3/*: any*/)
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "ExecuteWorkOrderMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "executeWorkOrder",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "WorkOrderExecutionResult",
        "plural": false,
        "selections": [
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "equipmentAdded",
            "storageKey": null,
            "args": null,
            "concreteType": "Equipment",
            "plural": true,
            "selections": [
              (v4/*: any*/),
              (v5/*: any*/),
              (v6/*: any*/),
              (v8/*: any*/),
              (v9/*: any*/),
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "device",
                "storageKey": null,
                "args": null,
                "concreteType": "Device",
                "plural": false,
                "selections": [
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "up",
                    "args": null,
                    "storageKey": null
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
                "selections": [
                  (v4/*: any*/)
                ]
              }
            ]
          },
          (v2/*: any*/),
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "linkAdded",
            "storageKey": null,
            "args": null,
            "concreteType": "Link",
            "plural": true,
            "selections": [
              (v4/*: any*/),
              (v6/*: any*/),
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "ports",
                "storageKey": null,
                "args": null,
                "concreteType": "EquipmentPort",
                "plural": true,
                "selections": [
                  (v4/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "definition",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "EquipmentPortDefinition",
                    "plural": false,
                    "selections": [
                      (v4/*: any*/),
                      (v5/*: any*/),
                      (v10/*: any*/),
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
                              (v4/*: any*/),
                              (v5/*: any*/),
                              (v11/*: any*/),
                              {
                                "kind": "ScalarField",
                                "alias": null,
                                "name": "index",
                                "args": null,
                                "storageKey": null
                              },
                              (v12/*: any*/),
                              (v13/*: any*/),
                              (v14/*: any*/),
                              (v15/*: any*/),
                              (v16/*: any*/),
                              (v17/*: any*/),
                              (v18/*: any*/),
                              (v19/*: any*/),
                              (v20/*: any*/),
                              (v21/*: any*/),
                              (v22/*: any*/)
                            ]
                          },
                          (v4/*: any*/)
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
                      (v4/*: any*/),
                      (v5/*: any*/),
                      (v6/*: any*/),
                      {
                        "kind": "LinkedField",
                        "alias": null,
                        "name": "equipmentType",
                        "storageKey": null,
                        "args": null,
                        "concreteType": "EquipmentType",
                        "plural": false,
                        "selections": [
                          (v4/*: any*/),
                          (v5/*: any*/),
                          {
                            "kind": "LinkedField",
                            "alias": null,
                            "name": "portDefinitions",
                            "storageKey": null,
                            "args": null,
                            "concreteType": "EquipmentPortDefinition",
                            "plural": true,
                            "selections": [
                              (v4/*: any*/),
                              (v5/*: any*/),
                              (v10/*: any*/),
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
                                "selections": (v7/*: any*/)
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
                          (v4/*: any*/),
                          (v5/*: any*/),
                          {
                            "kind": "LinkedField",
                            "alias": null,
                            "name": "locationType",
                            "storageKey": null,
                            "args": null,
                            "concreteType": "LocationType",
                            "plural": false,
                            "selections": (v23/*: any*/)
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
                          (v4/*: any*/),
                          {
                            "kind": "LinkedField",
                            "alias": null,
                            "name": "definition",
                            "storageKey": null,
                            "args": null,
                            "concreteType": "EquipmentPositionDefinition",
                            "plural": false,
                            "selections": [
                              (v4/*: any*/),
                              (v5/*: any*/),
                              (v10/*: any*/)
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
                              (v4/*: any*/),
                              (v5/*: any*/),
                              (v8/*: any*/)
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
                        "selections": (v23/*: any*/)
                      },
                      (v4/*: any*/)
                    ]
                  }
                ]
              },
              (v9/*: any*/),
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "properties",
                "storageKey": null,
                "args": null,
                "concreteType": "Property",
                "plural": true,
                "selections": [
                  (v4/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "propertyType",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "PropertyType",
                    "plural": false,
                    "selections": [
                      (v4/*: any*/),
                      (v5/*: any*/),
                      (v11/*: any*/),
                      (v20/*: any*/),
                      (v22/*: any*/),
                      (v21/*: any*/),
                      (v12/*: any*/)
                    ]
                  },
                  (v12/*: any*/),
                  (v13/*: any*/),
                  (v15/*: any*/),
                  (v14/*: any*/),
                  (v16/*: any*/),
                  (v17/*: any*/),
                  (v18/*: any*/),
                  (v19/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "equipmentValue",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "Equipment",
                    "plural": false,
                    "selections": (v7/*: any*/)
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "locationValue",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "Location",
                    "plural": false,
                    "selections": (v7/*: any*/)
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "serviceValue",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "Service",
                    "plural": false,
                    "selections": (v7/*: any*/)
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
                "selections": (v7/*: any*/)
              }
            ]
          },
          (v3/*: any*/),
          (v4/*: any*/)
        ]
      }
    ]
  },
  "params": {
    "operationKind": "mutation",
    "name": "ExecuteWorkOrderMutation",
    "id": null,
    "text": "mutation ExecuteWorkOrderMutation(\n  $id: ID!\n) {\n  executeWorkOrder(id: $id) {\n    equipmentAdded {\n      ...EquipmentTable_equipment\n      id\n    }\n    equipmentRemoved\n    linkAdded {\n      ...EquipmentPortsTable_link\n      id\n    }\n    linkRemoved\n    id\n  }\n}\n\nfragment EquipmentBreadcrumbs_equipment on Equipment {\n  id\n  name\n  equipmentType {\n    id\n    name\n  }\n  locationHierarchy {\n    id\n    name\n    locationType {\n      name\n      id\n    }\n  }\n  positionHierarchy {\n    id\n    definition {\n      id\n      name\n      visibleLabel\n    }\n    parentEquipment {\n      id\n      name\n      equipmentType {\n        id\n        name\n      }\n    }\n  }\n}\n\nfragment EquipmentPortsTable_link on Link {\n  id\n  futureState\n  ports {\n    id\n    definition {\n      id\n      name\n      visibleLabel\n      portType {\n        linkPropertyTypes {\n          id\n          name\n          type\n          index\n          stringValue\n          intValue\n          booleanValue\n          floatValue\n          latitudeValue\n          longitudeValue\n          rangeFromValue\n          rangeToValue\n          isEditable\n          isInstanceProperty\n          isMandatory\n        }\n        id\n      }\n    }\n    parentEquipment {\n      id\n      name\n      futureState\n      equipmentType {\n        id\n        name\n        portDefinitions {\n          id\n          name\n          visibleLabel\n          bandwidth\n          portType {\n            id\n            name\n          }\n        }\n      }\n      ...EquipmentBreadcrumbs_equipment\n    }\n    serviceEndpoints {\n      role\n      service {\n        name\n        id\n      }\n      id\n    }\n  }\n  workOrder {\n    id\n    status\n  }\n  properties {\n    id\n    propertyType {\n      id\n      name\n      type\n      isEditable\n      isMandatory\n      isInstanceProperty\n      stringValue\n    }\n    stringValue\n    intValue\n    floatValue\n    booleanValue\n    latitudeValue\n    longitudeValue\n    rangeFromValue\n    rangeToValue\n    equipmentValue {\n      id\n      name\n    }\n    locationValue {\n      id\n      name\n    }\n    serviceValue {\n      id\n      name\n    }\n  }\n  services {\n    id\n    name\n  }\n}\n\nfragment EquipmentTable_equipment on Equipment {\n  id\n  name\n  futureState\n  equipmentType {\n    id\n    name\n  }\n  workOrder {\n    id\n    status\n  }\n  device {\n    up\n  }\n  services {\n    id\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'b2fc04d3444e653e33dc106aa1bcd291';
module.exports = node;
