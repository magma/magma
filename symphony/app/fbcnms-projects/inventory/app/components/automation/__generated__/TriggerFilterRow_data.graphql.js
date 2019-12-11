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
type TriggerFilterOperator_data$ref = any;
export type TriggerID = "magma_alert" | "%future added value";
import type { FragmentReference } from "relay-runtime";
declare export opaque type TriggerFilterRow_data$ref: FragmentReference;
declare export opaque type TriggerFilterRow_data$fragmentType: TriggerFilterRow_data$ref;
export type TriggerFilterRow_data = {|
  +triggerID: TriggerID,
  +supportedFilters: $ReadOnlyArray<?{|
    +filterID: string,
    +description: string,
    +supportedOperators: $ReadOnlyArray<?{|
      +operatorID: string
    |}>,
    +$fragmentRefs: TriggerFilterOperator_data$ref,
  |}>,
  +$refType: TriggerFilterRow_data$ref,
|};
export type TriggerFilterRow_data$data = TriggerFilterRow_data;
export type TriggerFilterRow_data$key = {
  +$data?: TriggerFilterRow_data$data,
  +$fragmentRefs: TriggerFilterRow_data$ref,
};
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "TriggerFilterRow_data",
  "type": "ActionsTrigger",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "triggerID",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "supportedFilters",
      "storageKey": null,
      "args": null,
      "concreteType": "ActionsFilter",
      "plural": true,
      "selections": [
        {
          "kind": "ScalarField",
          "alias": null,
          "name": "filterID",
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
            }
          ]
        },
        {
          "kind": "FragmentSpread",
          "name": "TriggerFilterOperator_data",
          "args": null
        }
      ]
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = '47c9db7480b461fd45c11dc90763638e';
module.exports = node;
