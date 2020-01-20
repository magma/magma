/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 11cce26e3186ad86421b4b8b69198fca
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type ServiceCard_service$ref = any;
export type RemoveServiceLinkMutationVariables = {|
  id: string,
  linkId: string,
|};
export type RemoveServiceLinkMutationResponse = {|
  +removeServiceLink: ?{|
    +$fragmentRefs: ServiceCard_service$ref
  |}
|};
export type RemoveServiceLinkMutation = {|
  variables: RemoveServiceLinkMutationVariables,
  response: RemoveServiceLinkMutationResponse,
|};
*/


/*
mutation RemoveServiceLinkMutation(
  $id: ID!
  $linkId: ID!
) {
  removeServiceLink(id: $id, linkId: $linkId) {
    ...ServiceCard_service
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

fragment ForceNetworkTopology_topology on NetworkTopology {
  nodes {
    __typename
    id
  }
  links {
    source {
      __typename
      id
    }
    target {
      __typename
      id
    }
  }
}

fragment ServiceCard_service on Service {
  id
  name
  ...ServiceDetailsPanel_service
  ...ServicePanel_service
  topology {
    ...ServiceEquipmentTopology_topology
  }
  endpoints {
    ...ServiceEquipmentTopology_endpoints
    id
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
      isMandatory
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
    serviceValue {
      id
      name
    }
  }
}

fragment ServiceEndpointsView_endpoints on ServiceEndpoint {
  id
  port {
    parentEquipment {
      name
      ...EquipmentBreadcrumbs_equipment
      id
    }
    definition {
      id
      name
    }
    id
  }
  role
}

fragment ServiceEquipmentTopology_endpoints on ServiceEndpoint {
  role
  port {
    parentEquipment {
      id
      positionHierarchy {
        parentEquipment {
          id
        }
        id
      }
    }
    id
  }
}

fragment ServiceEquipmentTopology_topology on NetworkTopology {
  nodes {
    __typename
    ... on Equipment {
      id
      name
    }
    id
  }
  ...ForceNetworkTopology_topology
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

fragment ServicePanel_service on Service {
  id
  name
  externalId
  status
  customer {
    name
    id
  }
  serviceType {
    name
    id
  }
  links {
    id
    ...ServiceLinksView_links
  }
  endpoints {
    ...ServiceEndpointsView_endpoints
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
  },
  {
    "kind": "LocalArgument",
    "name": "linkId",
    "type": "ID!",
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
    "name": "linkId",
    "variableName": "linkId"
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
  (v3/*: any*/),
  (v2/*: any*/)
],
v5 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "isInstanceProperty",
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
  "name": "floatValue",
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
  "name": "isMandatory",
  "args": null,
  "storageKey": null
},
v16 = [
  (v2/*: any*/),
  (v3/*: any*/)
],
v17 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "definition",
  "storageKey": null,
  "args": null,
  "concreteType": "EquipmentPortDefinition",
  "plural": false,
  "selections": (v16/*: any*/)
},
v18 = {
  "kind": "LinkedField",
  "alias": null,
  "name": "equipmentType",
  "storageKey": null,
  "args": null,
  "concreteType": "EquipmentType",
  "plural": false,
  "selections": (v16/*: any*/)
},
v19 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "__typename",
  "args": null,
  "storageKey": null
},
v20 = [
  (v19/*: any*/),
  (v2/*: any*/)
];
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "RemoveServiceLinkMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "removeServiceLink",
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
    "name": "RemoveServiceLinkMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "removeServiceLink",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "Service",
        "plural": false,
        "selections": [
          (v2/*: any*/),
          (v3/*: any*/),
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
            "selections": (v4/*: any*/)
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
              (v2/*: any*/),
              (v3/*: any*/),
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
                  (v5/*: any*/),
                  (v6/*: any*/),
                  (v7/*: any*/),
                  (v8/*: any*/),
                  (v9/*: any*/),
                  (v10/*: any*/),
                  (v11/*: any*/),
                  (v12/*: any*/),
                  (v13/*: any*/),
                  (v14/*: any*/),
                  (v15/*: any*/)
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
                  (v6/*: any*/),
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "isEditable",
                    "args": null,
                    "storageKey": null
                  },
                  (v5/*: any*/),
                  (v15/*: any*/),
                  (v7/*: any*/)
                ]
              },
              (v7/*: any*/),
              (v8/*: any*/),
              (v9/*: any*/),
              (v10/*: any*/),
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
                "selections": (v16/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "locationValue",
                "storageKey": null,
                "args": null,
                "concreteType": "Location",
                "plural": false,
                "selections": (v16/*: any*/)
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "serviceValue",
                "storageKey": null,
                "args": null,
                "concreteType": "Service",
                "plural": false,
                "selections": (v16/*: any*/)
              }
            ]
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
                    "selections": (v16/*: any*/)
                  },
                  (v17/*: any*/),
                  (v2/*: any*/)
                ]
              }
            ]
          },
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "endpoints",
            "storageKey": null,
            "args": null,
            "concreteType": "ServiceEndpoint",
            "plural": true,
            "selections": [
              (v2/*: any*/),
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "port",
                "storageKey": null,
                "args": null,
                "concreteType": "EquipmentPort",
                "plural": false,
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
                      (v3/*: any*/),
                      (v2/*: any*/),
                      (v18/*: any*/),
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
                            "selections": (v4/*: any*/)
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
                              (v2/*: any*/),
                              (v3/*: any*/),
                              (v18/*: any*/)
                            ]
                          }
                        ]
                      }
                    ]
                  },
                  (v17/*: any*/),
                  (v2/*: any*/)
                ]
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "role",
                "args": null,
                "storageKey": null
              }
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
                "concreteType": null,
                "plural": true,
                "selections": [
                  (v19/*: any*/),
                  (v2/*: any*/),
                  {
                    "kind": "InlineFragment",
                    "type": "Equipment",
                    "selections": [
                      (v3/*: any*/)
                    ]
                  }
                ]
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
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "source",
                    "storageKey": null,
                    "args": null,
                    "concreteType": null,
                    "plural": false,
                    "selections": (v20/*: any*/)
                  },
                  {
                    "kind": "LinkedField",
                    "alias": null,
                    "name": "target",
                    "storageKey": null,
                    "args": null,
                    "concreteType": null,
                    "plural": false,
                    "selections": (v20/*: any*/)
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
    "name": "RemoveServiceLinkMutation",
    "id": null,
    "text": "mutation RemoveServiceLinkMutation(\n  $id: ID!\n  $linkId: ID!\n) {\n  removeServiceLink(id: $id, linkId: $linkId) {\n    ...ServiceCard_service\n    id\n  }\n}\n\nfragment EquipmentBreadcrumbs_equipment on Equipment {\n  id\n  name\n  equipmentType {\n    id\n    name\n  }\n  locationHierarchy {\n    id\n    name\n    locationType {\n      name\n      id\n    }\n  }\n  positionHierarchy {\n    id\n    definition {\n      id\n      name\n      visibleLabel\n    }\n    parentEquipment {\n      id\n      name\n      equipmentType {\n        id\n        name\n      }\n    }\n  }\n}\n\nfragment ForceNetworkTopology_topology on NetworkTopology {\n  nodes {\n    __typename\n    id\n  }\n  links {\n    source {\n      __typename\n      id\n    }\n    target {\n      __typename\n      id\n    }\n  }\n}\n\nfragment ServiceCard_service on Service {\n  id\n  name\n  ...ServiceDetailsPanel_service\n  ...ServicePanel_service\n  topology {\n    ...ServiceEquipmentTopology_topology\n  }\n  endpoints {\n    ...ServiceEquipmentTopology_endpoints\n    id\n  }\n}\n\nfragment ServiceDetailsPanel_service on Service {\n  id\n  name\n  externalId\n  customer {\n    name\n    id\n  }\n  serviceType {\n    id\n    name\n    propertyTypes {\n      id\n      name\n      index\n      isInstanceProperty\n      type\n      stringValue\n      intValue\n      floatValue\n      booleanValue\n      latitudeValue\n      longitudeValue\n      rangeFromValue\n      rangeToValue\n      isMandatory\n    }\n  }\n  properties {\n    id\n    propertyType {\n      id\n      name\n      type\n      isEditable\n      isInstanceProperty\n      isMandatory\n      stringValue\n    }\n    stringValue\n    intValue\n    floatValue\n    booleanValue\n    latitudeValue\n    longitudeValue\n    rangeFromValue\n    rangeToValue\n    equipmentValue {\n      id\n      name\n    }\n    locationValue {\n      id\n      name\n    }\n    serviceValue {\n      id\n      name\n    }\n  }\n}\n\nfragment ServiceEndpointsView_endpoints on ServiceEndpoint {\n  id\n  port {\n    parentEquipment {\n      name\n      ...EquipmentBreadcrumbs_equipment\n      id\n    }\n    definition {\n      id\n      name\n    }\n    id\n  }\n  role\n}\n\nfragment ServiceEquipmentTopology_endpoints on ServiceEndpoint {\n  role\n  port {\n    parentEquipment {\n      id\n      positionHierarchy {\n        parentEquipment {\n          id\n        }\n        id\n      }\n    }\n    id\n  }\n}\n\nfragment ServiceEquipmentTopology_topology on NetworkTopology {\n  nodes {\n    __typename\n    ... on Equipment {\n      id\n      name\n    }\n    id\n  }\n  ...ForceNetworkTopology_topology\n}\n\nfragment ServiceLinksView_links on Link {\n  id\n  ports {\n    parentEquipment {\n      id\n      name\n    }\n    definition {\n      id\n      name\n    }\n    id\n  }\n}\n\nfragment ServicePanel_service on Service {\n  id\n  name\n  externalId\n  status\n  customer {\n    name\n    id\n  }\n  serviceType {\n    name\n    id\n  }\n  links {\n    id\n    ...ServiceLinksView_links\n  }\n  endpoints {\n    ...ServiceEndpointsView_endpoints\n    id\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '78ff6dbd88739e84dc5b8cfdff47db32';
module.exports = node;
