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
declare export opaque type PositionDefinitionsAddEditTable_positionDefinition$ref: FragmentReference;
declare export opaque type PositionDefinitionsAddEditTable_positionDefinition$fragmentType: PositionDefinitionsAddEditTable_positionDefinition$ref;
export type PositionDefinitionsAddEditTable_positionDefinition = {
  +id: string,
  +name: string,
  +index: ?number,
  +visibleLabel: ?string,
};
export type PositionDefinitionsAddEditTable_positionDefinition$data = PositionDefinitionsAddEditTable_positionDefinition;
export type PositionDefinitionsAddEditTable_positionDefinition$key = {
  +$data?: PositionDefinitionsAddEditTable_positionDefinition$data,
  +$fragmentRefs: PositionDefinitionsAddEditTable_positionDefinition$ref,
};
*/


const node/*: ReaderFragment*/ = {
  "kind": "Fragment",
  "name": "PositionDefinitionsAddEditTable_positionDefinition",
  "type": "EquipmentPositionDefinition",
  "metadata": {
    "mask": false
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
(node/*: any*/).hash = '3952fd6597286104bfc0889a1d16bb1a';
module.exports = node;
