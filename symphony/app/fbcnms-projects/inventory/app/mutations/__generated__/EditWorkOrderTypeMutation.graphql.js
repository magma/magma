/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 50764a9ae02a495cceb23c342f022d10
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type AddEditWorkOrderTypeCard_editingWorkOrderType$ref = any;
type WorkOrderTypeItem_workOrderType$ref = any;
export type CheckListItemType = "enum" | "simple" | "string" | "%future added value";
export type PropertyKind = "bool" | "date" | "email" | "enum" | "equipment" | "float" | "gps_location" | "int" | "location" | "range" | "string" | "%future added value";
export type EditWorkOrderTypeInput = {|
  id: string,
  name: string,
  description?: ?string,
  properties?: ?$ReadOnlyArray<?PropertyTypeInput>,
  checkList?: ?$ReadOnlyArray<?CheckListDefinitionInput>,
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
|};
export type CheckListDefinitionInput = {|
  id?: ?string,
  title: string,
  type: CheckListItemType,
  index?: ?number,
  enumValues?: ?string,
  helpText?: ?string,
|};
export type EditWorkOrderTypeMutationVariables = {|
  input: EditWorkOrderTypeInput
|};
export type EditWorkOrderTypeMutationResponse = {|
  +editWorkOrderType: ?{|
    +id: string,
    +name: string,
    +$fragmentRefs: WorkOrderTypeItem_workOrderType$ref & AddEditWorkOrderTypeCard_editingWorkOrderType$ref,
  |}
|};
export type EditWorkOrderTypeMutation = {|
  variables: EditWorkOrderTypeMutationVariables,
  response: EditWorkOrderTypeMutationResponse,
|};
*/


/*
mutation EditWorkOrderTypeMutation(
  $input: EditWorkOrderTypeInput!
) {
  editWorkOrderType(input: $input) {
    id
    name
    ...WorkOrderTypeItem_workOrderType
    ...AddEditWorkOrderTypeCard_editingWorkOrderType
  }
}

fragment AddEditWorkOrderTypeCard_editingWorkOrderType on WorkOrderType {
  id
  name
  description
  numberOfWorkOrders
  propertyTypes {
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
  }
}

fragment DynamicPropertyTypesGrid_propertyTypes on PropertyType {
  ...PropertyTypeFormField_propertyType
  id
  index
}

fragment PropertyTypeFormField_propertyType on PropertyType {
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
}

fragment WorkOrderTypeItem_workOrderType on WorkOrderType {
  id
  name
  propertyTypes {
    ...DynamicPropertyTypesGrid_propertyTypes
    id
  }
  numberOfWorkOrders
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "input",
    "type": "EditWorkOrderTypeInput!",
    "defaultValue": null
  }
],
v1 = [
  {
    "kind": "Variable",
    "name": "input",
    "variableName": "input"
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
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "EditWorkOrderTypeMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "editWorkOrderType",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "WorkOrderType",
        "plural": false,
        "selections": [
          (v2/*: any*/),
          (v3/*: any*/),
          {
            "kind": "FragmentSpread",
            "name": "WorkOrderTypeItem_workOrderType",
            "args": null
          },
          {
            "kind": "FragmentSpread",
            "name": "AddEditWorkOrderTypeCard_editingWorkOrderType",
            "args": null
          }
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "EditWorkOrderTypeMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "editWorkOrderType",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "WorkOrderType",
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
                "name": "type",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "index",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "stringValue",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "intValue",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "booleanValue",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "floatValue",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "latitudeValue",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "longitudeValue",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "rangeFromValue",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "rangeToValue",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "isEditable",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "isInstanceProperty",
                "args": null,
                "storageKey": null
              }
            ]
          },
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "numberOfWorkOrders",
            "args": null,
            "storageKey": null
          },
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "description",
            "args": null,
            "storageKey": null
          }
        ]
      }
    ]
  },
  "params": {
    "operationKind": "mutation",
    "name": "EditWorkOrderTypeMutation",
    "id": null,
    "text": "mutation EditWorkOrderTypeMutation(\n  $input: EditWorkOrderTypeInput!\n) {\n  editWorkOrderType(input: $input) {\n    id\n    name\n    ...WorkOrderTypeItem_workOrderType\n    ...AddEditWorkOrderTypeCard_editingWorkOrderType\n  }\n}\n\nfragment AddEditWorkOrderTypeCard_editingWorkOrderType on WorkOrderType {\n  id\n  name\n  description\n  numberOfWorkOrders\n  propertyTypes {\n    id\n    name\n    type\n    index\n    stringValue\n    intValue\n    booleanValue\n    floatValue\n    latitudeValue\n    longitudeValue\n    rangeFromValue\n    rangeToValue\n    isEditable\n    isInstanceProperty\n  }\n}\n\nfragment DynamicPropertyTypesGrid_propertyTypes on PropertyType {\n  ...PropertyTypeFormField_propertyType\n  id\n  index\n}\n\nfragment PropertyTypeFormField_propertyType on PropertyType {\n  id\n  name\n  type\n  index\n  stringValue\n  intValue\n  booleanValue\n  floatValue\n  latitudeValue\n  longitudeValue\n  rangeFromValue\n  rangeToValue\n  isEditable\n  isInstanceProperty\n}\n\nfragment WorkOrderTypeItem_workOrderType on WorkOrderType {\n  id\n  name\n  propertyTypes {\n    ...DynamicPropertyTypesGrid_propertyTypes\n    id\n  }\n  numberOfWorkOrders\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '84908125a2c8612e196f0c7590bd28da';
module.exports = node;
