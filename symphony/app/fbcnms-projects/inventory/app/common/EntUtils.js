/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {FragmentReference} from 'relay-runtime';

import shortid from 'shortid';
import {camelCase, startCase, toUpper} from 'lodash';

type RefType = $ReadOnly<{|
  $refType: FragmentReference,
|}>;
export type OptionalRefTypeWrapper<T> = $ReadOnly<{|
  ...$Rest<T, RefType>,
|}>;

type EntWithID = $ReadOnly<{
  id?: ?string,
  ...
}>;

export type NamedNode = {id: string, name: string};

export type ShortUser = $ReadOnly<{
  id: string,
  email: string,
}>;

// http://github.com/golang/lint/blob/master/lint.go
const commonGoInitialisms = [
  'ACL',
  'API',
  'ASCII',
  'CPU',
  'CSS',
  'DNS',
  'EOF',
  'GUID',
  'HTML',
  'HTTP',
  'HTTPS',
  'ID',
  'IP',
  'JSON',
  'LHS',
  'QPS',
  'RAM',
  'RHS',
  'RPC',
  'SLA',
  'SMTP',
  'SQL',
  'SSH',
  'TCP',
  'TLS',
  'TTL',
  'UDP',
  'UI',
  'UID',
  'UUID',
  'URI',
  'URL',
  'UTF8',
  'VM',
  'XML',
  'XMPP',
  'XSRF',
  'XSS',
];

export const ENT_TEMP_ID_PREFIX = '@tmp';

export const generateTempId = () => {
  return `${ENT_TEMP_ID_PREFIX}${shortid.generate()}`;
};

export const isTempId = (id: string): boolean => {
  return id != null && (id.startsWith(ENT_TEMP_ID_PREFIX) || isNaN(id));
};

export const getGraphError = (error: Error): string => {
  if (error.hasOwnProperty('source')) {
    // eslint-disable-next-line no-warning-comments
    // $FlowFixMe verified there's sources T58630520
    return error.source.errors[0].message;
  }
  return error.message;
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

export const pascalCaseGoStyle = (word: string) => {
  return startCase(camelCase(word))
    .split(' ')
    .map(w => (commonGoInitialisms.includes(toUpper(w)) ? toUpper(w) : w))
    .join('');
};

export function haveDifferentValues<T_ENT>(entA: T_ENT, entB: T_ENT): boolean {
  if (
    entA == null ||
    entB == null ||
    typeof entA != 'object' ||
    typeof entB != 'object'
  ) {
    return entA != entB;
  }
  const propsToCompare = Object.keys(entA).filter(prop =>
    entA.hasOwnProperty(prop),
  );
  return !!propsToCompare.find(prop => entA[prop] != entB[prop]);
}

export type EntsMap<T: EntWithID> = Map<string, T>;
export function ent2EntsMap<T: EntWithID>(ents: Array<T>): EntsMap<T> {
  return new Map<string, T>(
    ents.filter(ent => ent.id != null).map(ent => [ent.id || '', ent]),
  );
}

export type KeyValueEnum<TValues> = {
  [key: TValues]: {
    key: TValues,
    value: string,
  },
};
