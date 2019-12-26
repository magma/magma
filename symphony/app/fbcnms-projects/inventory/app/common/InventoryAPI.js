/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

const INVENTORY_PATH = '/inventory/inventory';
const LOCATION_SEARCH_PARAM = 'location';
const EQUIPMENT_SEARCH_PARAM = 'equipment';

export const InventoryAPIUrls = {
  location: (locationId: string) =>
    `${INVENTORY_PATH}?${LOCATION_SEARCH_PARAM}=${locationId}`,
  equipment: (equipmentId: string) =>
    `${INVENTORY_PATH}?${EQUIPMENT_SEARCH_PARAM}=${equipmentId}`,
  project: (projectId: string) =>
    `/workorders/projects/search?project=${projectId}`,
  workorder: (workorderId: ?string) =>
    `/workorders/search${!!workorderId ? `?workorder=${workorderId}` : ''}`,
};
