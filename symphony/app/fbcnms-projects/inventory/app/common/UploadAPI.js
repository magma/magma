/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

export const UploadAPIUrls = {
  locations: () => '/graph/import/location',
  equipment: () => '/graph/import/equipment',
  port_connect: () => '/graph/import/port_connect',
  port_definition: () => '/graph/import/port_def',
  position_definition: () => '/graph/import/position_def',
  exported_equipment: () => '/graph/import/export_equipment',
  exported_ports: () => '/graph/import/export_ports',
  exported_links: () => '/graph/import/export_links',
  exported_service: () => '/graph/import/export_service',
};
