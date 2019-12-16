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
declare export opaque type ProjectTypeCard_projectType$ref: FragmentReference;
declare export opaque type ProjectTypeCard_projectType$fragmentType: ProjectTypeCard_projectType$ref;
export type ProjectTypeCard_projectType = {|
  +id: string,
  +name: string,
  +description: ?string,
  +numberOfProjects: number,
  +workOrders: $ReadOnlyArray<?{|
    +id: string
  |}>,
  +$refType: ProjectTypeCard_projectType$ref,
|};
export type ProjectTypeCard_projectType$data = ProjectTypeCard_projectType;
export type ProjectTypeCard_projectType$key = {
  +$data?: ProjectTypeCard_projectType$data,
  +$fragmentRefs: ProjectTypeCard_projectType$ref,
};
*/


const node/*: ReaderFragment*/ = (function(){
var v0 = {
  "kind": "ScalarField",
  "alias": null,
  "name": "id",
  "args": null,
  "storageKey": null
};
return {
  "kind": "Fragment",
  "name": "ProjectTypeCard_projectType",
  "type": "ProjectType",
  "metadata": null,
  "argumentDefinitions": [],
  "selections": [
    (v0/*: any*/),
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
      "name": "description",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "ScalarField",
      "alias": null,
      "name": "numberOfProjects",
      "args": null,
      "storageKey": null
    },
    {
      "kind": "LinkedField",
      "alias": null,
      "name": "workOrders",
      "storageKey": null,
      "args": null,
      "concreteType": "WorkOrderDefinition",
      "plural": true,
      "selections": [
        (v0/*: any*/)
      ]
    }
  ]
};
})();
// prettier-ignore
(node/*: any*/).hash = 'a0ed06d279a9e96ad0fbb45c505ad5e8';
module.exports = node;
