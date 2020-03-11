/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 * @relayHash 02dec20f0bdba8d64639a9730708cafb
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
type AddEditLocationTypeCard_editingLocationType$ref = any;
type LocationTypeItem_locationType$ref = any;
export type PropertyKind = "bool" | "date" | "datetime_local" | "email" | "enum" | "equipment" | "float" | "gps_location" | "int" | "location" | "range" | "service" | "string" | "%future added value";
export type SurveyQuestionType = "BOOL" | "CELLULAR" | "COORDS" | "DATE" | "EMAIL" | "FLOAT" | "INTEGER" | "PHONE" | "PHOTO" | "TEXT" | "TEXTAREA" | "WIFI" | "%future added value";
export type AddLocationTypeInput = {|
  name: string,
  mapType?: ?string,
  mapZoomLevel?: ?number,
  isSite?: ?boolean,
  properties?: ?$ReadOnlyArray<PropertyTypeInput>,
  surveyTemplateCategories?: ?$ReadOnlyArray<SurveyTemplateCategoryInput>,
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
  isMandatory?: ?boolean,
  isDeleted?: ?boolean,
|};
export type SurveyTemplateCategoryInput = {|
  id?: ?string,
  categoryTitle: string,
  categoryDescription: string,
  surveyTemplateQuestions?: ?$ReadOnlyArray<?SurveyTemplateQuestionInput>,
|};
export type SurveyTemplateQuestionInput = {|
  id?: ?string,
  questionTitle: string,
  questionDescription: string,
  questionType: SurveyQuestionType,
  index: number,
|};
export type AddLocationTypeMutationVariables = {|
  input: AddLocationTypeInput
|};
export type AddLocationTypeMutationResponse = {|
  +addLocationType: {|
    +id: string,
    +name: string,
    +$fragmentRefs: LocationTypeItem_locationType$ref & AddEditLocationTypeCard_editingLocationType$ref,
  |}
|};
export type AddLocationTypeMutation = {|
  variables: AddLocationTypeMutationVariables,
  response: AddLocationTypeMutationResponse,
|};
*/


/*
mutation AddLocationTypeMutation(
  $input: AddLocationTypeInput!
) {
  addLocationType(input: $input) {
    id
    name
    ...LocationTypeItem_locationType
    ...AddEditLocationTypeCard_editingLocationType
  }
}

fragment AddEditLocationTypeCard_editingLocationType on LocationType {
  id
  name
  mapType
  mapZoomLevel
  numberOfLocations
  isSite
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
    isMandatory
    isInstanceProperty
  }
  surveyTemplateCategories {
    id
    categoryTitle
    categoryDescription
    surveyTemplateQuestions {
      id
      questionTitle
      questionDescription
      questionType
      index
    }
  }
}

fragment DynamicPropertyTypesGrid_propertyTypes on PropertyType {
  ...PropertyTypeFormField_propertyType
  id
  index
}

fragment LocationTypeItem_locationType on LocationType {
  id
  name
  index
  propertyTypes {
    ...DynamicPropertyTypesGrid_propertyTypes
    id
  }
  numberOfLocations
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
  isMandatory
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "kind": "LocalArgument",
    "name": "input",
    "type": "AddLocationTypeInput!",
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
},
v4 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "index",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Request",
  "fragment": {
    "kind": "Fragment",
    "name": "AddLocationTypeMutation",
    "type": "Mutation",
    "metadata": null,
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "addLocationType",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "LocationType",
        "plural": false,
        "selections": [
          (v2/*: any*/),
          (v3/*: any*/),
          {
            "kind": "FragmentSpread",
            "name": "LocationTypeItem_locationType",
            "args": null
          },
          {
            "kind": "FragmentSpread",
            "name": "AddEditLocationTypeCard_editingLocationType",
            "args": null
          }
        ]
      }
    ]
  },
  "operation": {
    "kind": "Operation",
    "name": "AddLocationTypeMutation",
    "argumentDefinitions": (v0/*: any*/),
    "selections": [
      {
        "kind": "LinkedField",
        "alias": null,
        "name": "addLocationType",
        "storageKey": null,
        "args": (v1/*: any*/),
        "concreteType": "LocationType",
        "plural": false,
        "selections": [
          (v2/*: any*/),
          (v3/*: any*/),
          (v4/*: any*/),
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
              (v4/*: any*/),
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
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "isMandatory",
                "args": null,
                "storageKey": null
              }
            ]
          },
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "numberOfLocations",
            "args": null,
            "storageKey": null
          },
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "mapType",
            "args": null,
            "storageKey": null
          },
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "mapZoomLevel",
            "args": null,
            "storageKey": null
          },
          {
            "kind": "ScalarField",
            "alias": null,
            "name": "isSite",
            "args": null,
            "storageKey": null
          },
          {
            "kind": "LinkedField",
            "alias": null,
            "name": "surveyTemplateCategories",
            "storageKey": null,
            "args": null,
            "concreteType": "SurveyTemplateCategory",
            "plural": true,
            "selections": [
              (v2/*: any*/),
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "categoryTitle",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "ScalarField",
                "alias": null,
                "name": "categoryDescription",
                "args": null,
                "storageKey": null
              },
              {
                "kind": "LinkedField",
                "alias": null,
                "name": "surveyTemplateQuestions",
                "storageKey": null,
                "args": null,
                "concreteType": "SurveyTemplateQuestion",
                "plural": true,
                "selections": [
                  (v2/*: any*/),
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "questionTitle",
                    "args": null,
                    "storageKey": null
                  },
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "questionDescription",
                    "args": null,
                    "storageKey": null
                  },
                  {
                    "kind": "ScalarField",
                    "alias": null,
                    "name": "questionType",
                    "args": null,
                    "storageKey": null
                  },
                  (v4/*: any*/)
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
    "name": "AddLocationTypeMutation",
    "id": null,
    "text": "mutation AddLocationTypeMutation(\n  $input: AddLocationTypeInput!\n) {\n  addLocationType(input: $input) {\n    id\n    name\n    ...LocationTypeItem_locationType\n    ...AddEditLocationTypeCard_editingLocationType\n  }\n}\n\nfragment AddEditLocationTypeCard_editingLocationType on LocationType {\n  id\n  name\n  mapType\n  mapZoomLevel\n  numberOfLocations\n  isSite\n  propertyTypes {\n    id\n    name\n    type\n    index\n    stringValue\n    intValue\n    booleanValue\n    floatValue\n    latitudeValue\n    longitudeValue\n    rangeFromValue\n    rangeToValue\n    isEditable\n    isMandatory\n    isInstanceProperty\n  }\n  surveyTemplateCategories {\n    id\n    categoryTitle\n    categoryDescription\n    surveyTemplateQuestions {\n      id\n      questionTitle\n      questionDescription\n      questionType\n      index\n    }\n  }\n}\n\nfragment DynamicPropertyTypesGrid_propertyTypes on PropertyType {\n  ...PropertyTypeFormField_propertyType\n  id\n  index\n}\n\nfragment LocationTypeItem_locationType on LocationType {\n  id\n  name\n  index\n  propertyTypes {\n    ...DynamicPropertyTypesGrid_propertyTypes\n    id\n  }\n  numberOfLocations\n}\n\nfragment PropertyTypeFormField_propertyType on PropertyType {\n  id\n  name\n  type\n  index\n  stringValue\n  intValue\n  booleanValue\n  floatValue\n  latitudeValue\n  longitudeValue\n  rangeFromValue\n  rangeToValue\n  isEditable\n  isInstanceProperty\n  isMandatory\n}\n",
    "metadata": {}
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '754cc54988ca2f679ee665da803fae10';
module.exports = node;
