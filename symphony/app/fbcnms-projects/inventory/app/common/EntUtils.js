/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

type EntWithID = $ReadOnly<{
  id?: ?string,
}>;

export const ENT_TEMP_ID_PREFIX = '@tmp';

export const removeTempID = (ent: EntWithID) => {
  if (ent.id && (ent.id.startsWith(ENT_TEMP_ID_PREFIX) || isNaN(ent.id))) {
    const {id: _, ...noIdEnt} = ent;
    return noIdEnt;
  }
  return ent;
};

export const removeTempIDs = (ents: Iterable<EntWithID>) => {
  return Array.prototype.map.call(ents, removeTempID);
};
