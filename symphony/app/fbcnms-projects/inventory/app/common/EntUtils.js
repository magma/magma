/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import shortid from 'shortid';

type EntWithID = $ReadOnly<{
  id?: ?string,
}>;

export const ENT_TEMP_ID_PREFIX = '@tmp';

export const generateTempId = () => {
  return `${ENT_TEMP_ID_PREFIX}${shortid.generate()}`;
};

export const isTempId = (id: string): boolean => {
  return id != null && (id.startsWith(ENT_TEMP_ID_PREFIX) || isNaN(id));
};

export const removeTempID = (ent: EntWithID) => {
  if (ent.id && isTempId(ent.id)) {
    const {id: _, ...noIdEnt} = ent;
    return noIdEnt;
  }
  return ent;
};

export const removeTempIDs = (ents: Iterable<EntWithID>) => {
  return Array.prototype.map.call(ents, removeTempID);
};
