/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

export const extractEntityIdFromUrl = (
  entityType:
    | 'location'
    | 'equipment'
    | 'workorder'
    | 'workorderType'
    | 'project'
    | 'projectType'
    | 'service'
    | 'serviceType',
  searchParams: string,
): ?string => {
  const query = new URLSearchParams(searchParams);
  return query.get(entityType);
};
