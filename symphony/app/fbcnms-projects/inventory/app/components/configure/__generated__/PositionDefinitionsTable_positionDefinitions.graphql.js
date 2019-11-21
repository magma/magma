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
import type { FragmentReference } from "relay-runtime";
declare export opaque type PositionDefinitionsTable_positionDefinitions$ref: FragmentReference;
declare export opaque type PositionDefinitionsTable_positionDefinitions$fragmentType: PositionDefinitionsTable_positionDefinitions$ref;
export type PositionDefinitionsTable_positionDefinitions = $ReadOnlyArray<{|
  +id: string,
  +name: string,
  +index: ?number,
  +visibleLabel: ?string,
  +$refType: PositionDefinitionsTable_positionDefinitions$ref,
|}>;
export type PositionDefinitionsTable_positionDefinitions$data = PositionDefinitionsTable_positionDefinitions;
export type PositionDefinitionsTable_positionDefinitions$key = $ReadOnlyArray<{
  +$data?: PositionDefinitionsTable_positionDefinitions$data,
  +$fragmentRefs: PositionDefinitionsTable_positionDefinitions$ref,
}>;
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "PositionDefinitionsTable_positionDefinitions",
  "type": "EquipmentPositionDefinition",
  "metadata": {
    "plural": true
  },
  "argumentDefinitions": [],
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
      "name": "name",
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
      "name": "visibleLabel",
      "args": null,
      "storageKey": null
    }
  ]
};
// prettier-ignore
(node/*: any*/).hash = 'e380bf28ed2fa9bb1090ea936f6e7b25';
module.exports = node;
