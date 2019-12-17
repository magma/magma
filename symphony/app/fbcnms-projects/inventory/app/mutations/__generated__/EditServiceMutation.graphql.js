/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 3b1cc4b4c02f7c40ef8454ebb70f83e5
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type ServiceCard_service$ref = any;
export type ServiceEditData = {|
  id: string,
  name: string,
  externalId?: ?string,
  customerId?: ?string,
  upstreamServiceIds: $ReadOnlyArray<string>,
  properties?: ?$ReadOnlyArray<?PropertyInput>,
  terminationPointIds: $ReadOnlyArray<string>,
|};
export type PropertyInput = {|
  id?: ?string,
  propertyTypeID: string,
  stringValue?: ?string,
  intValue?: ?number,
  booleanValue?: ?boolean,
  floatValue?: ?number,
  latitudeValue?: ?number,
  longitudeValue?: ?number,
  rangeFromValue?: ?number,
  rangeToValue?: ?number,
  equipmentIDValue?: ?string,
  locationIDValue?: ?string,
  isEditable?: ?boolean,
  isInstanceProperty?: ?boolean,
|};
export type EditServiceMutationVariables = {|
  data: ServiceEditData
|};
export type EditServiceMutationResponse = {|
  +editService: ?{|
    +$fragmentRefs: ServiceCard_service$ref
  |}
|};
export type EditServiceMutation = {|
  variables: EditServiceMutationVariables,
  response: EditServiceMutationResponse,
|};
*/


/*
mutation EditServiceMutation(
  $data: ServiceEditData!
) {
  editService(data: $data) {
    ...ServiceCard_service
    id
  }
}

fragment ServiceCard_service on Service {
  id
  name
  serviceType {
    name
    id
  }
  ...ServiceDetailsPanel_service
  links {
    id
    ...ServiceLinksView_links
  }
  terminationPoints {
    ...ServiceEquipmentTopology_terminationPoints
    id
  }
  topology {
    ...ServiceEquipmentTopology_topology
  }
}

fragment ServiceDetailsPanel_service on Service {
  id
  name
  externalId
  customer {
    name
    id
  }
  serviceType {
    id
    name
    propertyTypes {
      id
      name
      index
      isInstanceProperty
      type
      stringValue
      intValue
      floatValue
      booleanValue
      latitudeValue
      longitudeValue
      rangeFromValue
      rangeToValue
    }
  }
  properties {
    id
    propertyType {
      id
      name
      type
      isEditable
      isInstanceProperty
      isMandatory
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
}

fragment ServiceEquipmentTopology_terminationPoints on Equipment {
  id
}

fragment ServiceEquipmentTopology_topology on NetworkTopology {
  nodes {
    id
    name
  }
  links {
    source
    target
  }
}

fragment ServiceLinksView_links on Link {
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
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "data",
    "type": "ServiceEditData!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "Variable",
    "name": "data",
    "variableName": "data"
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
v4 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isInstanceProperty",
  "args": null,
  "storageKey": null
},
v5 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "type",
  "args": null,
  "storageKey": null
},
v6 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "stringValue",
  "args": null,
  "storageKey": null
},
v7 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "intValue",
  "args": null,
  "storageKey": null
},
v8 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "floatValue",
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
  "name": "latitudeValue",
  "args": null,
  "storageKey": null
},
v11 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "longitudeValue",
  "args": null,
  "storageKey": null
},
v12 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeFromValue",
  "args": null,
  "storageKey": null
},
v13 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "rangeToValue",
  "args": null,
  "storageKey": null
},
v14 = [
  (v2/*: any*/),
  (v3/*: any*/)
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "EditServiceMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "editService",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "Service",
        "plural": false,
        "selections": [
          {
            "kind": "FragmentSpread",
            "name": "ServiceCard_service",
            "args": null
          }
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "EditServiceMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "editService",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "Service",
        "plural": false,
        "selections": [
          (v2/*: any*/),
          (v3/*: any*/),
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
              (v2/*: any*/),
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "propertyTypes",
                "storageKey": null,
                "args": null,
                "concreteType": "PropertyType",
                "plural": true,
                "selections": [
                  (v2/*: any*/),
                  (v3/*: any*/),
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "index",
                    "args": null,
                    "storageKey": null
                  },
                  (v4/*: any*/),
                  (v5/*: any*/),
                  (v6/*: any*/),
                  (v7/*: any*/),
                  (v8/*: any*/),
                  (v9/*: any*/),
                  (v10/*: any*/),
                  (v11/*: any*/),
                  (v12/*: any*/),
                  (v13/*: any*/)
                ]
              }
            ]
          },
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "externalId",
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
            "selections": [
              (v3/*: any*/),
              (v2/*: any*/)
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
                  (v5/*: any*/),
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "isEditable",
                    "args": null,
                    "storageKey": null
                  },
                  (v4/*: any*/),
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "isMandatory",
                    "args": null,
                    "storageKey": null
                  },
                  (v6/*: any*/)
                ]
              },
              (v6/*: any*/),
              (v7/*: any*/),
              (v8/*: any*/),
              (v9/*: any*/),
              (v10/*: any*/),
              (v11/*: any*/),
              (v12/*: any*/),
              (v13/*: any*/),
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "equipmentValue",
                "storageKey": null,
                "args": null,
                "concreteType": "Equipment",
                "plural": false,
                "selections": (v14/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "locationValue",
                "storageKey": null,
                "args": null,
                "concreteType": "Location",
                "plural": false,
                "selections": (v14/*: any*/)
              }
            ]
          },
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "links",
            "storageKey": null,
            "args": null,
            "concreteType": "Link",
            "plural": true,
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
                    "selections": (v14/*: any*/)
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "definition",
                    "storageKey": null,
                    "args": null,
                    "concreteType": "EquipmentPortDefinition",
                    "plural": false,
                    "selections": (v14/*: any*/)
                  },
                  (v2/*: any*/)
                ]
              }
            ]
          },
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "terminationPoints",
            "storageKey": null,
            "args": null,
            "concreteType": "Equipment",
            "plural": true,
            "selections": [
              (v2/*: any*/)
            ]
          },
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "topology",
            "storageKey": null,
            "args": null,
            "concreteType": "NetworkTopology",
            "plural": false,
            "selections": [
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "nodes",
                "storageKey": null,
                "args": null,
                "concreteType": "Equipment",
                "plural": true,
                "selections": (v14/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "links",
                "storageKey": null,
                "args": null,
                "concreteType": "TopologyLink",
                "plural": true,
                "selections": [
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "source",
                    "args": null,
                    "storageKey": null
                  },
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "target",
                    "args": null,
                    "storageKey": null
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
    "operationKind": "mutation",
    "name": "EditServiceMutation",
    "id": null,
    "text": "mutation EditServiceMutation(\n  $data: ServiceEditData!\n) {\n  editService(data: $data) {\n    ...ServiceCard_service\n    id\n  }\n}\n\nfragment ServiceCard_service on Service {\n  id\n  name\n  serviceType {\n    name\n    id\n  }\n  ...ServiceDetailsPanel_service\n  links {\n    id\n    ...ServiceLinksView_links\n  }\n  terminationPoints {\n    ...ServiceEquipmentTopology_terminationPoints\n    id\n  }\n  topology {\n    ...ServiceEquipmentTopology_topology\n  }\n}\n\nfragment ServiceDetailsPanel_service on Service {\n  id\n  name\n  externalId\n  customer {\n    name\n    id\n  }\n  serviceType {\n    id\n    name\n    propertyTypes {\n      id\n      name\n      index\n      isInstanceProperty\n      type\n      stringValue\n      intValue\n      floatValue\n      booleanValue\n      latitudeValue\n      longitudeValue\n      rangeFromValue\n      rangeToValue\n    }\n  }\n  properties {\n    id\n    propertyType {\n      id\n      name\n      type\n      isEditable\n      isInstanceProperty\n      isMandatory\n      stringValue\n    }\n    stringValue\n    intValue\n    floatValue\n    booleanValue\n    latitudeValue\n    longitudeValue\n    rangeFromValue\n    rangeToValue\n    equipmentValue {\n      id\n      name\n    }\n    locationValue {\n      id\n      name\n    }\n  }\n}\n\nfragment ServiceEquipmentTopology_terminationPoints on Equipment {\n  id\n}\n\nfragment ServiceEquipmentTopology_topology on NetworkTopology {\n  nodes {\n    id\n    name\n  }\n  links {\n    source\n    target\n  }\n}\n\nfragment ServiceLinksView_links on Link {\n  id\n  ports {\n    parentEquipment {\n      id\n      name\n    }\n    definition {\n      id\n      name\n    }\n    id\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '0d262e10a1451069e21df0155a3ac18d';
module.exports = node;
