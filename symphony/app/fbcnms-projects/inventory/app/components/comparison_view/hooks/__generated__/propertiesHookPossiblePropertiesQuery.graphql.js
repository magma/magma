/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 413cd186d130ee3447909ec33968bf52
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type PropertyEntity = "EQUIPMENT" | "LINK" | "LOCATION" | "PORT" | "PROJECT" | "SERVICE" | "WORK_ORDER" | "%future added value";
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "float" | "gps_location" | "int" | "node" | "range" | "string" | "%future added value";
export type propertiesHookPossiblePropertiesQueryVariables = {|
  entityType: PropertyEntity
|};
export type propertiesHookPossiblePropertiesQueryResponse = {|
  +possibleProperties: $ReadOnlyArray<{|
    +name: string,
    +type: PropertyKind,
    +stringValue: ?string,
  |}>
|};
export type propertiesHookPossiblePropertiesQuery = {|
  variables: propertiesHookPossiblePropertiesQueryVariables,
  response: propertiesHookPossiblePropertiesQueryResponse,
|};
*/


/*
query propertiesHookPossiblePropertiesQuery(
  $entityType: PropertyEntity!
) {
  possibleProperties(entityType: $entityType) {
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
    "kind": "LocalArgument",
    "name": "entityType",
    "type": "PropertyEntity!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "Variable",
    "name": "entityType",
    "variableName": "entityType"
  }
],
v2 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "name",
  "args": null,
  "storageKey": null
},
v3 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "type",
  "args": null,
  "storageKey": null
},
v4 = {
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
    "name": "propertiesHookPossiblePropertiesQuery",
    "type": "Query",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "possibleProperties",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "PropertyType",
        "plural": true,
        "selections": [
          (v2/*: any*/),
          (v3/*: any*/),
          (v4/*: any*/)
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "propertiesHookPossiblePropertiesQuery",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "possibleProperties",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "PropertyType",
        "plural": true,
        "selections": [
          (v2/*: any*/),
          (v3/*: any*/),
          (v4/*: any*/),
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
    "name": "propertiesHookPossiblePropertiesQuery",
    "id": null,
    "text": "query propertiesHookPossiblePropertiesQuery(\n  $entityType: PropertyEntity!\n) {\n  possibleProperties(entityType: $entityType) {\n    name\n    type\n    stringValue\n    id\n  }\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'bb938da07f028b2a4f6fc67588a53b78';
module.exports = node;
