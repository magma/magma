/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {EntityConfig} from '../comparison_view/ComparisonViewTypes';

import PowerSearchServiceCustomerNameFilter from './PowerSearchServiceCustomerNameFilter';
import PowerSearchServiceDiscoveryMethodFilter from './PowerSearchServiceDiscoveryMethodFilter';
import PowerSearchServiceEquipmentInServiceFilter from './PowerSearchServiceEquipmentInServiceFilter';
import PowerSearchServiceExternalIDFilter from './PowerSearchServiceExternalIDFilter';
import PowerSearchServiceNameFilter from './PowerSearchServiceNameFilter';
import PowerSearchServiceStatusFilter from './PowerSearchServiceStatusFilter';
import PowerSearchServiceTypeFilter from './PowerSearchServiceTypeFilter';

const ServiceSearchConfig: Array<EntityConfig> = [
  {
    type: 'service',
    label: 'service',
    filters: [
      {
        key: 'service_name',
        name: 'service_inst_name',
        entityType: 'service',
        label: 'Name',
        component: PowerSearchServiceNameFilter,
        defaultOperator: 'contains',
      },
      {
        key: 'service_type',
        name: 'service_type',
        entityType: 'service',
        label: 'Type',
        component: PowerSearchServiceTypeFilter,
        defaultOperator: 'is_one_of',
      },
      {
        key: 'service_discovery_method',
        name: 'service_discovery_method',
        entityType: 'service',
        label: 'Discovery Method',
        component: PowerSearchServiceDiscoveryMethodFilter,
        defaultOperator: 'is_one_of',
      },
      {
        key: 'service_external_id',
        name: 'service_inst_external_id',
        entityType: 'service',
        label: 'Service ID',
        component: PowerSearchServiceExternalIDFilter,
        defaultOperator: 'is',
      },
      {
        key: 'service_status',
        name: 'service_status',
        entityType: 'service',
        label: 'Status',
        component: PowerSearchServiceStatusFilter,
        defaultOperator: 'is_one_of',
      },
      {
        key: 'equipment_in_service',
        name: 'equipment_in_service',
        entityType: 'service',
        label: 'Service Equipments',
        component: PowerSearchServiceEquipmentInServiceFilter,
        defaultOperator: 'contains',
      },
      {
        key: 'customer_name',
        name: 'service_inst_customer_name',
        entityType: 'service',
        label: 'Customer',
        component: PowerSearchServiceCustomerNameFilter,
        defaultOperator: 'contains',
      },
    ],
  },
  {
    type: 'location_by_types',
    label: 'Location (By consumer endpoints)',
    filters: [],
  },
  {
    type: 'properties',
    label: 'Properties',
    filters: [],
  },
];

export {ServiceSearchConfig};
