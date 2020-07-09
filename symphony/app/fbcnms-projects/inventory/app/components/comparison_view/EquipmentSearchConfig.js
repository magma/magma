/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {EntityConfig} from './ComparisonViewTypes';

import PowerSearchEquipmentNameFilter from './PowerSearchEquipmentNameFilter';
import PowerSearchEquipmentTypeFilter from './PowerSearchEquipmentTypeFilter';
import PowerSearchExternalIDFilter from './PowerSearchExternalIDFilter';

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
    type: 'locations',
    label: 'Location',
    filters: [
      {
        key: 'location_inst_external_id',
        name: 'location_inst_external_id',
        entityType: 'locations',
        label: 'Location External ID',
        component: PowerSearchExternalIDFilter,
        defaultOperator: 'contains',
      },
    ],
  },
  {
    type: 'properties',
    label: 'Properties',
    filters: [],
  },
];

export {EquipmentCriteriaConfig};
