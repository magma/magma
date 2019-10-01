/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {FeatureID} from '@fbcnms/types/features';

export type FeatureConfig = {
  id: FeatureID,
  title: string,
  enabledByDefault: boolean,
};

const {FeatureFlag} = require('@fbcnms/sequelize-models');

const arrayConfigs = [
  {
    id: 'sso_example_feature',
    title: 'SSO Example Feature',
    enabledByDefault: false,
  },
  {
    id: 'audit_log_example_feature',
    title: 'Audit Log Example Feature',
    enabledByDefault: true,
  },
  {
    id: 'audit_log_view',
    title: 'Audit Log View',
    enabledByDefault: false,
  },
  {
    id: 'third_party_devices',
    title: 'Third Party Devices',
    enabledByDefault: false,
  },
  {
    id: 'network_topology',
    title: 'Network Topology',
    enabledByDefault: false,
  },
  {
    id: 'upload_rural',
    title: 'Bulk Upload: Rural',
    enabledByDefault: false,
  },
  {
    id: 'upload_xwf',
    title: 'Bulk Upload: XWF',
    enabledByDefault: false,
  },
  {
    id: 'upload_ftth',
    title: 'Bulk Upload: FTTH',
    enabledByDefault: false,
  },
  {
    id: 'python_api',
    title: 'Download Puthon API',
    enabledByDefault: false,
  },
  {
    id: 'lte_network_metrics',
    title: 'LTE Network Metrics',
    enabledByDefault: true,
  },
  {
    id: 'site_survey',
    title: 'Site Survey',
    enabledByDefault: false,
  },
  {
    id: 'alerts',
    title: 'Alerts',
    enabledByDefault: true,
  },
  {
    id: 'logs',
    title: 'Logs',
    enabledByDefault: false,
  },
  {
    id: 'equipment_export',
    title: 'Equipment Export',
    enabledByDefault: true,
  },
  {
    id: 'magma_network_management',
    title: 'Magma Network Management',
    enabledByDefault: false,
  },
  {
    id: 'file_categories',
    title: 'File Categories (for IpT)',
    enabledByDefault: false,
  },
  {
    id: 'floor_plans',
    title: 'Floor Plans',
    enabledByDefault: false,
  },
  {
    id: 'import_exported_equipemnt',
    title: 'Imported Exported Equipment',
    enabledByDefault: false,
  },
  {
    id: 'work_order_map',
    title: 'Work order map',
    enabledByDefault: false,
  },
  {
    id: 'documents_site',
    title: 'Documents Site',
    enabledByDefault: true,
  },
];

const featureConfigs: {[FeatureID]: FeatureConfig} = {};
arrayConfigs.map(config => (featureConfigs[config.id] = config));

export async function isFeatureEnabled(
  featureId: FeatureID,
  organization: ?string,
) {
  if (organization) {
    const flag = await FeatureFlag.findOne({where: {organization, featureId}});
    if (flag) {
      return flag.enabled;
    }
  }

  return featureConfigs[featureId].enabledByDefault;
}

export async function getEnabledFeatures(
  organization: ?string,
): Promise<FeatureID[]> {
  const results = await Promise.all(
    arrayConfigs.map(async (config): Promise<?FeatureID> => {
      const enabled = await isFeatureEnabled(config.id, organization);
      return enabled ? config.id : null;
    }),
  );

  return results.filter(Boolean);
}

export default {...featureConfigs};
