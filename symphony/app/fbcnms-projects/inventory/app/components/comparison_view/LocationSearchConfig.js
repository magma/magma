/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import PowerSearchExternalIDFilter from './PowerSearchExternalIDFilter';
import PowerSearchLocationHasEquipmentFilter from './PowerSearchLocationHasEquipmentFilter';
import PowerSearchLocationTypeFilter from './PowerSearchLocationTypeFilter';

import type {EntityConfig} from './ComparisonViewTypes';

const LocationCriteriaConfig: Array<EntityConfig> = [
  {
    type: 'location',
    label: 'Location',
    filters: [
      {
        key: 'location_type',
        name: 'location_type',
        entityType: 'location',
        label: 'Location Type',
        component: PowerSearchLocationTypeFilter,
        defaultOperator: 'is_one_of',
      },
      {
        key: 'location_inst_has_equipment',
        name: 'location_inst_has_equipment',
        entityType: 'location',
        label: 'Has Equipment',
        component: PowerSearchLocationHasEquipmentFilter,
        defaultOperator: 'is',
      },
    ],
  },
  {
    type: 'locations',
    label: 'Location Ancestor',
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

export {LocationCriteriaConfig};
