/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 0908d51a232e6d05f1d511f9398f4ebd
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "equipment" | "float" | "gps_location" | "int" | "location" | "range" | "service" | "string" | "%future added value";
export type ServiceComparisonViewQueryRendererPropertiesQueryVariables = {||};
export type ServiceComparisonViewQueryRendererPropertiesQueryResponse = {|
  +possibleProperties: $ReadOnlyArray<{|
    +name: string,
    +type: PropertyKind,
    +stringValue: ?string,
  |}>
|};
export type ServiceComparisonViewQueryRendererPropertiesQuery = {|
  variables: ServiceComparisonViewQueryRendererPropertiesQueryVariables,
  response: ServiceComparisonViewQueryRendererPropertiesQueryResponse,
|};
*/


/*
query ServiceComparisonViewQueryRendererPropertiesQuery {
  possibleProperties(entityType: SERVICE) {
    name
    type
    stringValue
    id
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "Literal",
    "name": "entityType",
    "value": "SERVICE"
  }
],
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
  "name": "type",
  "args": null,
  "storageKey": null
},
v3 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "stringValue",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "ServiceComparisonViewQueryRendererPropertiesQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "possibleProperties",
        "storageKey": "possibleProperties(entityType:\"SERVICE\")",
        "args": (v0/*: any*/),
        "concreteType": "PropertyType",
        "plural": true,
        "selections": [
          (v1/*: any*/),
          (v2/*: any*/),
          (v3/*: any*/)
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "ServiceComparisonViewQueryRendererPropertiesQuery",
    "argumentDefinitions": [],
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "possibleProperties",
        "storageKey": "possibleProperties(entityType:\"SERVICE\")",
        "args": (v0/*: any*/),
        "concreteType": "PropertyType",
        "plural": true,
        "selections": [
          (v1/*: any*/),
          (v2/*: any*/),
          (v3/*: any*/),
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "id",
            "args": null,
            "storageKey": null
          }
        ]
      }
    ]
  },
  "params": {
    "operationKind": "query",
    "name": "ServiceComparisonViewQueryRendererPropertiesQuery",
    "id": null,
    "text": "query ServiceComparisonViewQueryRendererPropertiesQuery {\n  possibleProperties(entityType: SERVICE) {\n    name\n    type\n    stringValue\n    id\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '5c530ac75d04a8db88c9284a7790ac3a';
module.exports = node;
