/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ReaderFragment } from 'relay-runtime';
export type ActionsDataType = "string" | "stringArray" | "%future added value";
import type { FragmentReference } from "relay-runtime";
declare export opaque type TriggerFilterOperator_data$ref: FragmentReference;
declare export opaque type TriggerFilterOperator_data$fragmentType: TriggerFilterOperator_data$ref;
export type TriggerFilterOperator_data = {|
  +supportedOperators: $ReadOnlyArray<?{|
    +operatorID: string,
    +description: string,
    +dataType: ActionsDataType,
  |}>,
  +$refType: TriggerFilterOperator_data$ref,
|};
export type TriggerFilterOperator_data$data = TriggerFilterOperator_data;
export type TriggerFilterOperator_data$key = {
  +$data?: TriggerFilterOperator_data$data,
  +$fragmentRefs: TriggerFilterOperator_data$ref,
  ...
};
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "TriggerFilterOperator_data",
  "type": "ActionsFilter",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "supportedOperators",
      "storageKey": null,
      "args": null,
      "concreteType": "ActionsOperator",
      "plural": true,
      "selections": [
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "operatorID",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "description",
          "args": null,
          "storageKey": null
        },
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "dataType",
          "args": null,
          "storageKey": null
        }
      ]
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = '2c0b10a0f1bbd74ce1cb44bc791d1b9a';
module.exports = node;
