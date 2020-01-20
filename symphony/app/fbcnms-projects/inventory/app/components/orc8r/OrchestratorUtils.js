/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {DeviceConfig} from './PDUConfigRow';

import nullthrows from '@fbcnms/util/nullthrows';
import {find} from 'lodash';

export function deriveConfigsFromEquipment(equipment: any): DeviceConfig[] {
  return equipment.positions
    .filter(position => position?.attachedEquipment)
    .map(position => {
      const attachedEquipment = nullthrows(position?.attachedEquipment);
      const ipProperty = find(
        attachedEquipment.properties,
        property => property && property.name.match(/^ip$|ip\s?address/i),
      );

      return {
        id: attachedEquipment.id,
        name: attachedEquipment.name,
        ipAddress: ipProperty?.stringValue || '',
      };
    });
}
