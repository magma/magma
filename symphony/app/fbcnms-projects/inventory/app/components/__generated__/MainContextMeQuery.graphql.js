/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash ad1818bc711a496cc206f1b4fcde578a
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type PermissionValue = "BY_CONDITION" | "NO" | "YES" | "%future added value";
export type MainContextMeQueryVariables = {||};
export type MainContextMeQueryResponse = {|
  +me: ?{|
    +user: ?{|
      +id: string,
      +authID: string,
      +email: string,
      +firstName: string,
      +lastName: string,
    |},
    +permissions: {|
      +canWrite: boolean,
      +adminPolicy: {|
        +access: {|
          +isAllowed: PermissionValue
        |}
      |},
      +inventoryPolicy: {|
        +read: {|
          +isAllowed: PermissionValue
        |},
        +location: {|
          +create: {|
            +isAllowed: PermissionValue,
            +locationTypeIds: ?$ReadOnlyArray<string>,
          |},
          +update: {|
            +isAllowed: PermissionValue,
            +locationTypeIds: ?$ReadOnlyArray<string>,
          |},
          +delete: {|
            +isAllowed: PermissionValue,
            +locationTypeIds: ?$ReadOnlyArray<string>,
          |},
        |},
        +equipment: {|
          +create: {|
            +isAllowed: PermissionValue
          |},
          +update: {|
            +isAllowed: PermissionValue
          |},
          +delete: {|
            +isAllowed: PermissionValue
          |},
        |},
        +equipmentType: {|
          +create: {|
            +isAllowed: PermissionValue
          |},
          +update: {|
            +isAllowed: PermissionValue
          |},
          +delete: {|
            +isAllowed: PermissionValue
          |},
        |},
        +locationType: {|
          +create: {|
            +isAllowed: PermissionValue
          |},
          +update: {|
            +isAllowed: PermissionValue
          |},
          +delete: {|
            +isAllowed: PermissionValue
          |},
        |},
        +portType: {|
          +create: {|
            +isAllowed: PermissionValue
          |},
          +update: {|
            +isAllowed: PermissionValue
          |},
          +delete: {|
            +isAllowed: PermissionValue
          |},
        |},
        +serviceType: {|
          +create: {|
            +isAllowed: PermissionValue
          |},
          +update: {|
            +isAllowed: PermissionValue
          |},
          +delete: {|
            +isAllowed: PermissionValue
          |},
        |},
      |},
      +workforcePolicy: {|
        +read: {|
          +isAllowed: PermissionValue,
          +projectTypeIds: ?$ReadOnlyArray<string>,
          +workOrderTypeIds: ?$ReadOnlyArray<string>,
        |},
        +templates: {|
          +create: {|
            +isAllowed: PermissionValue
          |},
          +update: {|
            +isAllowed: PermissionValue
          |},
          +delete: {|
            +isAllowed: PermissionValue
          |},
        |},
        +data: {|
          +create: {|
            +isAllowed: PermissionValue
          |},
          +update: {|
            +isAllowed: PermissionValue
          |},
          +delete: {|
            +isAllowed: PermissionValue
          |},
          +assign: {|
            +isAllowed: PermissionValue
          |},
          +transferOwnership: {|
            +isAllowed: PermissionValue
          |},
        |},
      |},
    |},
  |}
|};
export type MainContextMeQuery = {|
  variables: MainContextMeQueryVariables,
  response: MainContextMeQueryResponse,
|};
*/


/*
query MainContextMeQuery {
  me {
    user {
      id
      authID
      email
      firstName
      lastName
    }
    permissions {
      canWrite
      adminPolicy {
        access {
          isAllowed
        }
      }
      inventoryPolicy {
        read {
          isAllowed
        }
        location {
          create {
            isAllowed
            locationTypeIds
          }
          update {
            isAllowed
            locationTypeIds
          }
          delete {
            isAllowed
            locationTypeIds
          }
        }
        equipment {
          create {
            isAllowed
          }
          update {
            isAllowed
          }
          delete {
            isAllowed
          }
        }
        equipmentType {
          create {
            isAllowed
          }
          update {
            isAllowed
          }
          delete {
            isAllowed
          }
        }
        locationType {
          create {
            isAllowed
          }
          update {
            isAllowed
          }
          delete {
            isAllowed
          }
        }
        portType {
          create {
            isAllowed
          }
          update {
            isAllowed
          }
          delete {
            isAllowed
          }
        }
        serviceType {
          create {
            isAllowed
          }
          update {
            isAllowed
          }
          delete {
            isAllowed
          }
        }
      }
      workforcePolicy {
        read {
          isAllowed
          projectTypeIds
          workOrderTypeIds
        }
        templates {
          create {
            isAllowed
          }
          update {
            isAllowed
          }
          delete {
            isAllowed
          }
        }
        data {
          create {
            isAllowed
          }
          update {
            isAllowed
          }
          delete {
            isAllowed
          }
          assign {
            isAllowed
          }
          transferOwnership {
            isAllowed
          }
        }
      }
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isAllowed",
  "args": null,
  "storageKey": null
},
v1 = [
  (v0/*: any*/)
],
v2 = [
  (v0/*: any*/),
  {
    "kind": "ScalarField",
    "alias": null,
    "name": "locationTypeIds",
    "args": null,
    "storageKey": null
  }
],
v3 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "create",
    "storageKey": null,
    "args": null,
    "concreteType": "BasicPermissionRule",
    "plural": false,
    "selections": (v1/*: any*/)
  },
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "update",
    "storageKey": null,
    "args": null,
    "concreteType": "BasicPermissionRule",
    "plural": false,
    "selections": (v1/*: any*/)
  },
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "delete",
    "storageKey": null,
    "args": null,
    "concreteType": "BasicPermissionRule",
    "plural": false,
    "selections": (v1/*: any*/)
  }
],
v4 = [
  {
    "kind": "LinkedField",
    "alias": null,
    "name": "me",
    "storageKey": null,
    "args": null,
    "concreteType": "Viewer",
    "plural": false,
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "user",
        "storageKey": null,
        "args": null,
        "concreteType": "User",
        "plural": false,
        "selections": [
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "id",
            "args": null,
            "storageKey": null
          },
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "authID",
            "args": null,
            "storageKey": null
          },
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "email",
            "args": null,
            "storageKey": null
          },
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "firstName",
            "args": null,
            "storageKey": null
          },
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "lastName",
            "args": null,
            "storageKey": null
          }
        ]
      },
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "permissions",
        "storageKey": null,
        "args": null,
        "concreteType": "PermissionSettings",
        "plural": false,
        "selections": [
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "canWrite",
            "args": null,
            "storageKey": null
          },
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "adminPolicy",
            "storageKey": null,
            "args": null,
            "concreteType": "AdministrativePolicy",
            "plural": false,
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "access",
                "storageKey": null,
                "args": null,
                "concreteType": "BasicPermissionRule",
                "plural": false,
                "selections": (v1/*: any*/)
              }
            ]
          },
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "inventoryPolicy",
            "storageKey": null,
            "args": null,
            "concreteType": "InventoryPolicy",
            "plural": false,
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "read",
                "storageKey": null,
                "args": null,
                "concreteType": "BasicPermissionRule",
                "plural": false,
                "selections": (v1/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "location",
                "storageKey": null,
                "args": null,
                "concreteType": "LocationCUD",
                "plural": false,
                "selections": [
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "create",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "LocationPermissionRule",
                    "plural": false,
                    "selections": (v2/*: any*/)
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "update",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "LocationPermissionRule",
                    "plural": false,
                    "selections": (v2/*: any*/)
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "delete",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "LocationPermissionRule",
                    "plural": false,
                    "selections": (v2/*: any*/)
                  }
                ]
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "equipment",
                "storageKey": null,
                "args": null,
                "concreteType": "CUD",
                "plural": false,
                "selections": (v3/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "equipmentType",
                "storageKey": null,
                "args": null,
                "concreteType": "CUD",
                "plural": false,
                "selections": (v3/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "locationType",
                "storageKey": null,
                "args": null,
                "concreteType": "CUD",
                "plural": false,
                "selections": (v3/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "portType",
                "storageKey": null,
                "args": null,
                "concreteType": "CUD",
                "plural": false,
                "selections": (v3/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "serviceType",
                "storageKey": null,
                "args": null,
                "concreteType": "CUD",
                "plural": false,
                "selections": (v3/*: any*/)
              }
            ]
          },
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "workforcePolicy",
            "storageKey": null,
            "args": null,
            "concreteType": "WorkforcePolicy",
            "plural": false,
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "read",
                "storageKey": null,
                "args": null,
                "concreteType": "WorkforcePermissionRule",
                "plural": false,
                "selections": [
                  (v0/*: any*/),
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "projectTypeIds",
                    "args": null,
                    "storageKey": null
                  },
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "workOrderTypeIds",
                    "args": null,
                    "storageKey": null
                  }
                ]
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "templates",
                "storageKey": null,
                "args": null,
                "concreteType": "CUD",
                "plural": false,
                "selections": (v3/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "data",
                "storageKey": null,
                "args": null,
                "concreteType": "WorkforceCUD",
                "plural": false,
                "selections": [
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "create",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "WorkforcePermissionRule",
                    "plural": false,
                    "selections": (v1/*: any*/)
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "update",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "WorkforcePermissionRule",
                    "plural": false,
                    "selections": (v1/*: any*/)
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "delete",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "WorkforcePermissionRule",
                    "plural": false,
                    "selections": (v1/*: any*/)
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "assign",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "WorkforcePermissionRule",
                    "plural": false,
                    "selections": (v1/*: any*/)
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "transferOwnership",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "WorkforcePermissionRule",
                    "plural": false,
                    "selections": (v1/*: any*/)
                  }
                ]
              }
            ]
          }
        ]
      }
    ]
  }
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "MainContextMeQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": [],
    "selections": (v4/*: any*/)
  },
  "operation": {
    "kind": "Operation",
    "name": "MainContextMeQuery",
    "argumentDefinitions": [],
    "selections": (v4/*: any*/)
  },
  "params": {
    "operationKind": "query",
    "name": "MainContextMeQuery",
    "id": null,
    "text": "query MainContextMeQuery {\n  me {\n    user {\n      id\n      authID\n      email\n      firstName\n      lastName\n    }\n    permissions {\n      canWrite\n      adminPolicy {\n        access {\n          isAllowed\n        }\n      }\n      inventoryPolicy {\n        read {\n          isAllowed\n        }\n        location {\n          create {\n            isAllowed\n            locationTypeIds\n          }\n          update {\n            isAllowed\n            locationTypeIds\n          }\n          delete {\n            isAllowed\n            locationTypeIds\n          }\n        }\n        equipment {\n          create {\n            isAllowed\n          }\n          update {\n            isAllowed\n          }\n          delete {\n            isAllowed\n          }\n        }\n        equipmentType {\n          create {\n            isAllowed\n          }\n          update {\n            isAllowed\n          }\n          delete {\n            isAllowed\n          }\n        }\n        locationType {\n          create {\n            isAllowed\n          }\n          update {\n            isAllowed\n          }\n          delete {\n            isAllowed\n          }\n        }\n        portType {\n          create {\n            isAllowed\n          }\n          update {\n            isAllowed\n          }\n          delete {\n            isAllowed\n          }\n        }\n        serviceType {\n          create {\n            isAllowed\n          }\n          update {\n            isAllowed\n          }\n          delete {\n            isAllowed\n          }\n        }\n      }\n      workforcePolicy {\n        read {\n          isAllowed\n          projectTypeIds\n          workOrderTypeIds\n        }\n        templates {\n          create {\n            isAllowed\n          }\n          update {\n            isAllowed\n          }\n          delete {\n            isAllowed\n          }\n        }\n        data {\n          create {\n            isAllowed\n          }\n          update {\n            isAllowed\n          }\n          delete {\n            isAllowed\n          }\n          assign {\n            isAllowed\n          }\n          transferOwnership {\n            isAllowed\n          }\n        }\n      }\n    }\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '5aea2c378611998a70ece0a7eb81c4c1';
module.exports = node;
