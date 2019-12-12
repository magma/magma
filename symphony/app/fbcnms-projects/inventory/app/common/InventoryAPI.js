/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

export const InventoryAPIUrls = {
  location: (locationId: string) =>
    `/inventory/inventory?location=${locationId}`,
  equipment: (equipmentId: string) =>
    `/inventory/inventory?equipment=${equipmentId}`,
  project: (projectId: string) =>
    `/workorders/projects/search?project=${projectId}`,
};
