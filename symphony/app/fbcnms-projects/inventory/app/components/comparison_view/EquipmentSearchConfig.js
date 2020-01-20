/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {EntityConfig} from './ComparisonViewTypes';

import PowerSearchEquipmentNameFilter from './PowerSearchEquipmentNameFilter';
import PowerSearchEquipmentTypeFilter from './PowerSearchEquipmentTypeFilter';

const EquipmentCriteriaConfig: Array<EntityConfig> = [
  {
    type: 'equipment',
    label: 'Equipment',
    filters: [
      {
        key: 'equipment_name',
        name: 'equip_inst_name',
        entityType: 'equipment',
        label: 'Name',
        component: PowerSearchEquipmentNameFilter,
        defaultOperator: 'contains',
      },
      {
        key: 'equipment_type',
        name: 'equipment_type',
        entityType: 'equipment',
        label: 'Type',
        component: PowerSearchEquipmentTypeFilter,
        defaultOperator: 'is_one_of',
      },
    ],
  },
  {
    type: 'location_by_types',
    label: 'Location',
    filters: [],
  },
  {
    type: 'properties',
    label: 'Properties',
    filters: [],
  },
];

export {EquipmentCriteriaConfig};
