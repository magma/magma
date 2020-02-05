/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

export type FeatureID =
  | 'lte_network_metrics'
  | 'sso_example_feature'
  | 'audit_log_example_feature'
  | 'audit_log_view'
  | 'third_party_devices'
  | 'network_topology'
  | 'site_survey'
  | 'alerts'
  | 'alert_receivers'
  | 'alert_routes'
  | 'alert_suppressions'
  | 'equipment_export'
  | 'import_exported_equipemnt'
  | 'import_exported_ports'
  | 'import_exported_links'
  | 'file_categories'
  | 'floor_plans'
  | 'work_order_map'
  | 'documents_site'
  | 'coverage_maps'
  | 'logs'
  | 'services'
  | 'planned_equipment'
  | 'multi_subject_reports'
  | 'equipment_live_status'
  | 'logged_out_alert'
  | 'external_id';
