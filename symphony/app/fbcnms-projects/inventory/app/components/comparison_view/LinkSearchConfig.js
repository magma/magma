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

import PowerSearchEquipmentTypeFilter from './PowerSearchEquipmentTypeFilter';
import PowerSearchExternalIDFilter from './PowerSearchExternalIDFilter';
import PowerSearchLinkFutureStateFilter from './PowerSearchLinkFutureStateFilter';
import PowerSearchLinkServiceNameFilter from './PowerSearchLinkServiceNameFilter';

const LinkCriteriaConfig: Array<EntityConfig> = [
  {
    type: 'link',
    label: 'Link',
    filters: [
      {
        key: 'link_future_status',
        name: 'link_future_status',
        entityType: 'link',
        label: 'Future State',
        component: PowerSearchLinkFutureStateFilter,
        defaultOperator: 'is_one_of',
      },
      {
        key: 'equipment_type',
        name: 'equipment_type',
        entityType: 'link',
        label: 'Equipment Type',
        component: PowerSearchEquipmentTypeFilter,
        defaultOperator: 'is_one_of',
      },
      {
        key: 'service_inst',
        name: 'service_inst',
        entityType: 'link',
        label: 'Used by Service',
        component: PowerSearchLinkServiceNameFilter,
        defaultOperator: 'contains',
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

export {LinkCriteriaConfig};
