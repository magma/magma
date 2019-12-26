/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 45d00756c69f76c4ef35564f10c01d7e
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type EquipmentBreadcrumbs_equipment$ref = any;
export type FutureState = "INSTALL" | "REMOVE" | "%future added value";
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "equipment" | "float" | "gps_location" | "int" | "location" | "range" | "string" | "%future added value";
export type WorkOrderStatus = "DONE" | "PENDING" | "PLANNED" | "%future added value";
export type RemoveLinkMutationVariables = {|
  id: string,
  workOrderId?: ?string,
|};
export type RemoveLinkMutationResponse = {|
  +removeLink: ?{|
    +id: string,
    +futureState: ?FutureState,
    +ports: $ReadOnlyArray<?{|
      +id: string,
      +definition: {|
        +id: string,
        +name: string,
        +visibleLabel: ?string,
        +portType: ?{|
          +linkPropertyTypes: $ReadOnlyArray<?{|
            +id: string,
            +name: string,
            +type: PropertyKind,
            +index: ?number,
            +stringValue: ?string,
            +intValue: ?number,
            +booleanValue: ?boolean,
            +floatValue: ?number,
            +latitudeValue: ?number,
            +longitudeValue: ?number,
            +rangeFromValue: ?number,
            +rangeToValue: ?number,
            +isEditable: ?boolean,
            +isInstanceProperty: ?boolean,
            +isMandatory: ?boolean,
          |}>
        |},
      |},
      +parentEquipment: {|
        +id: string,
        +name: string,
        +futureState: ?FutureState,
        +equipmentType: {|
          +id: string,
          +name: string,
          +portDefinitions: $ReadOnlyArray<?{|
            +id: string,
            +name: string,
            +visibleLabel: ?string,
            +bandwidth: ?string,
            +portType: ?{|
              +id: string,
              +name: string,
            |},
          |}>,
        |},
        +$fragmentRefs: EquipmentBreadcrumbs_equipment$ref,
      |},
    |}>,
    +workOrder: ?{|
      +id: string,
      +status: WorkOrderStatus,
    |},
    +properties: $ReadOnlyArray<?{|
      +id: string,
      +propertyType: {|
        +id: string,
        +name: string,
        +type: PropertyKind,
        +isEditable: ?boolean,
        +isMandatory: ?boolean,
        +isInstanceProperty: ?boolean,
        +stringValue: ?string,
      |},
      +stringValue: ?string,
      +intValue: ?number,
      +floatValue: ?number,
      +booleanValue: ?boolean,
      +latitudeValue: ?number,
      +longitudeValue: ?number,
      +rangeFromValue: ?number,
      +rangeToValue: ?number,
      +equipmentValue: ?{|
        +id: string,
        +name: string,
      |},
      +locationValue: ?{|
        +id: string,
        +name: string,
      |},
    |}>,
    +services: $ReadOnlyArray<?{|
      +id: string,
      +name: string,
    |}>,
  |}
|};
export type RemoveLinkMutation = {|
  variables: RemoveLinkMutationVariables,
  response: RemoveLinkMutationResponse,
|};
*/


/*
mutation RemoveLinkMutation(
  $id: ID!
  $workOrderId: ID
) {
  removeLink(id: $id, workOrderId: $workOrderId) {
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
    }
    services {
      id
      name
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
    "name": "id",
    "type": "ID!",
    "defaultValue": null
  },
  {
    "kind": "LocalArgument",
    "name": "workOrderId",
    "type": "ID",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "Variable",
    "name": "id",
    "variableName": "id"
  },
  {
    "kind": "Variable",
    "name": "workOrderId",
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
  "name": "futureState",
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
  "kind": "ScalarField",
  "alias": null,
  "name": "visibleLabel",
  "args": null,
  "storageKey": null
},
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
  "name": "stringValue",
  "args": null,
  "storageKey": null
},
v8 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "intValue",
  "args": null,
  "storageKey": null
},
v9 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "booleanValue",
  "args": null,
  "storageKey": null
},
v10 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "floatValue",
  "args": null,
  "storageKey": null
},
v11 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "latitudeValue",
  "args": null,
  "storageKey": null
},
v12 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "longitudeValue",
  "args": null,
  "storageKey": null
},
v13 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeFromValue",
  "args": null,
  "storageKey": null
},
v14 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeToValue",
  "args": null,
  "storageKey": null
},
v15 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isEditable",
  "args": null,
  "storageKey": null
},
v16 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isInstanceProperty",
  "args": null,
  "storageKey": null
},
v17 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isMandatory",
  "args": null,
  "storageKey": null
},
v18 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "linkPropertyTypes",
  "storageKey": null,
  "args": null,
  "concreteType": "PropertyType",
  "plural": true,
  "selections": [
    (v2/*: any*/),
    (v4/*: any*/),
    (v6/*: any*/),
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "index",
      "args": null,
      "storageKey": null
    },
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
    (v17/*: any*/)
  ]
},
v19 = [
  (v2/*: any*/),
  (v4/*: any*/)
],
v20 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "equipmentType",
  "storageKey": null,
  "args": null,
  "concreteType": "EquipmentType",
  "plural": false,
  "selections": [
    (v2/*: any*/),
    (v4/*: any*/),
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
        (v4/*: any*/),
        (v5/*: any*/),
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
          "selections": (v19/*: any*/)
        }
      ]
    }
  ]
},
v21 = {
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
v22 = {
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
        (v4/*: any*/),
        (v6/*: any*/),
        (v15/*: any*/),
        (v17/*: any*/),
        (v16/*: any*/),
        (v7/*: any*/)
      ]
    },
    (v7/*: any*/),
    (v8/*: any*/),
    (v10/*: any*/),
    (v9/*: any*/),
    (v11/*: any*/),
    (v12/*: any*/),
    (v13/*: any*/),
    (v14/*: any*/),
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "equipmentValue",
      "storageKey": null,
      "args": null,
      "concreteType": "Equipment",
      "plural": false,
      "selections": (v19/*: any*/)
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "locationValue",
      "storageKey": null,
      "args": null,
      "concreteType": "Location",
      "plural": false,
      "selections": (v19/*: any*/)
    }
  ]
},
v23 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "services",
  "storageKey": null,
  "args": null,
  "concreteType": "Service",
  "plural": true,
  "selections": (v19/*: any*/)
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "RemoveLinkMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "removeLink",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "Link",
        "plural": false,
        "selections": [
          (v2/*: any*/),
          (v3/*: any*/),
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
                  (v4/*: any*/),
                  (v5/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "portType",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "EquipmentPortType",
                    "plural": false,
                    "selections": [
                      (v18/*: any*/)
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
                  (v4/*: any*/),
                  (v3/*: any*/),
                  (v20/*: any*/),
                  {
                    "kind": "FragmentSpread",
                    "name": "EquipmentBreadcrumbs_equipment",
                    "args": null
                  }
                ]
              }
            ]
          },
          (v21/*: any*/),
          (v22/*: any*/),
          (v23/*: any*/)
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "RemoveLinkMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "removeLink",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "Link",
        "plural": false,
        "selections": [
          (v2/*: any*/),
          (v3/*: any*/),
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
                  (v4/*: any*/),
                  (v5/*: any*/),
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "portType",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "EquipmentPortType",
                    "plural": false,
                    "selections": [
                      (v18/*: any*/),
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
                  (v4/*: any*/),
                  (v3/*: any*/),
                  (v20/*: any*/),
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
                          (v2/*: any*/)
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
                          (v4/*: any*/),
                          (v5/*: any*/)
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
                          (v4/*: any*/),
                          {
                            "kind": "LinkedField",
                            "alias": null,
                            "name": "equipmentType",
                            "storageKey": null,
                            "args": null,
                            "concreteType": "EquipmentType",
                            "plural": false,
                            "selections": (v19/*: any*/)
                          }
                        ]
                      }
                    ]
                  }
                ]
              }
            ]
          },
          (v21/*: any*/),
          (v22/*: any*/),
          (v23/*: any*/)
        ]
      }
    ]
  },
  "params": {
    "operationKind": "mutation",
    "name": "RemoveLinkMutation",
    "id": null,
    "text": "mutation RemoveLinkMutation(\n  $id: ID!\n  $workOrderId: ID\n) {\n  removeLink(id: $id, workOrderId: $workOrderId) {\n    id\n    futureState\n    ports {\n      id\n      definition {\n        id\n        name\n        visibleLabel\n        portType {\n          linkPropertyTypes {\n            id\n            name\n            type\n            index\n            stringValue\n            intValue\n            booleanValue\n            floatValue\n            latitudeValue\n            longitudeValue\n            rangeFromValue\n            rangeToValue\n            isEditable\n            isInstanceProperty\n            isMandatory\n          }\n          id\n        }\n      }\n      parentEquipment {\n        id\n        name\n        futureState\n        equipmentType {\n          id\n          name\n          portDefinitions {\n            id\n            name\n            visibleLabel\n            bandwidth\n            portType {\n              id\n              name\n            }\n          }\n        }\n        ...EquipmentBreadcrumbs_equipment\n      }\n    }\n    workOrder {\n      id\n      status\n    }\n    properties {\n      id\n      propertyType {\n        id\n        name\n        type\n        isEditable\n        isMandatory\n        isInstanceProperty\n        stringValue\n      }\n      stringValue\n      intValue\n      floatValue\n      booleanValue\n      latitudeValue\n      longitudeValue\n      rangeFromValue\n      rangeToValue\n      equipmentValue {\n        id\n        name\n      }\n      locationValue {\n        id\n        name\n      }\n    }\n    services {\n      id\n      name\n    }\n  }\n}\n\nfragment EquipmentBreadcrumbs_equipment on Equipment {\n  id\n  name\n  equipmentType {\n    id\n    name\n  }\n  locationHierarchy {\n    id\n    name\n    locationType {\n      name\n      id\n    }\n  }\n  positionHierarchy {\n    id\n    definition {\n      id\n      name\n      visibleLabel\n    }\n    parentEquipment {\n      id\n      name\n      equipmentType {\n        id\n        name\n      }\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '779f6813617bbf5b040592f2e8c2c3e3';
module.exports = node;
