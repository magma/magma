/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import PowerSearchEquipmentNameFilter from './PowerSearchEquipmentNameFilter';
import PowerSearchPortDefinitionFilter from './PowerSearchPortDefinitionFilter';
import PowerSearchPortHasLinkFilter from './PowerSearchPortHasLinkFilter';

import type {EntityConfig} from './ComparisonViewTypes';

const PortCriteriaConfig: Array<EntityConfig> = [
  {
    type: 'port',
    label: 'Port',
    filters: [
      {
        key: 'port_def',
        name: 'port_def',
        entityType: 'port',
        label: 'Port Name',
        component: PowerSearchPortDefinitionFilter,
        defaultOperator: 'is_one_of',
      },
      {
        key: 'port_inst_has_link',
        name: 'port_inst_has_link',
        entityType: 'port',
        label: 'Has Link',
        component: PowerSearchPortHasLinkFilter,
        defaultOperator: 'is',
      },
      {
        key: 'port_inst_equipment',
        name: 'port_inst_equipment',
        entityType: 'port',
        label: 'Equipment Name',
        component: PowerSearchEquipmentNameFilter,
        defaultOperator: 'contains',
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

export {PortCriteriaConfig};
