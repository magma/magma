/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

export type FiringAlarm = {
  annotations: Labels,
  endsAt: string,
  fingerprint: string,
  labels: Labels,
  receivers: Array<GettableReceiver>,
  startsAt: string,
  status: FiringAlarmStatus,
  updatedAt: string,
  generatorURL?: string,
};

type GettableReceiver = {
  name: string,
};

export type FiringAlarmStatus = {
  inhibitedBy: Array<string>,
  silencedBy: Array<string>,
  state: string,
};

export type AlertReceiver = {
  name: string,
  email_configs?: Array<ReceiverEmailConfig>,
  slack_configs?: Array<ReceiverSlackConfig>,
  webhook_configs?: Array<ReceiverWebhookConfig>,
};

// names of all the <type>_configs lists for a receiver
export type ReceiverConfigListName =
  | 'email_configs'
  | 'slack_configs'
  | 'webhook_configs';

export type ReceiverEmailConfig = {
  auth_identity?: string,
  auth_password?: string,
  auth_secret?: string,
  auth_username?: string,
  from: string,
  headers?: {[string]: string},
  hello?: string,
  html?: string,
  send_resolved?: boolean,
  smarthost: string,
  text?: string,
  to: string,
};

export type ReceiverSlackConfig = {
  api_url: string,
  channel?: string,
  username?: string,
  color?: string,
  title?: string,
  pretext?: string,
  text?: string,
  fields?: Array<SlackField>,
  short_fields?: boolean,
  footer?: string,
  fallback?: string,
  callback_id?: string,
  icon_emoji?: string,
  icon_url?: string,
  image_url?: string,
  thumb_url?: string,
  link_names?: boolean,
  actions?: Array<SlackAction>,
};

export type SlackField = {
  title: string,
  value: string,
  short?: boolean,
};

export type SlackAction = {
  type: string,
  text: string,
  url: string,
  style?: string,
  name?: string,
  value?: string,
  confirm?: SlackConfirmField,
};

export type SlackConfirmField = {
  text: string,
  title: string,
  ok_text: string,
  dismiss_text: string,
};

export type ReceiverWebhookConfig = {
  http_config?: HTTPConfig,
  send_resolved?: boolean,
  url: string,
};

export type HTTPConfig = {
  basic_auth?: HTTPConfigBasicAuth,
  bearer_token?: string,
  proxy_url?: string,
};

export type HTTPConfigBasicAuth = {
  password: string,
  username: string,
};

/**
 * Prometheus alert rule configuration
 */
export type AlertConfig = {|
  alert: string,
  expr: string,
  annotations?: {[string]: string},
  for?: string,
  labels?: Labels,
  rawData?: AlertConfig,
  _isCustomAlertRule?: boolean,
|};

export type AlertRoutingTree = {
  receiver: string,
  continue?: boolean,
  group_by?: Array<string>,
  group_interval?: string,
  group_wait?: string,
  match?: {[string]: string},
  match_re?: {[string]: string},
  repeat_interval?: string,
  routes?: Array<AlertRoutingTree>,
};

export type BulkAlertUpdateResponse = {
  errors: {[string]: string},
  statuses: {[string]: string},
};

export type AlertSuppressionMatcher = {
  name: string,
  value: string,
  isRegex: boolean,
};

export type AlertSuppressionState = {
  state: string,
};

export type AlertSuppression = {
  id: string,
  startsAt: string,
  endsAt: string,
  updatedAt: string,
  matchers: Array<AlertSuppressionMatcher>,
  createdBy: string,
  status: AlertSuppressionState,
  comment?: string,
};

export type Labels = {
  [string]: string,
};

export type PrometheusLabelset = {
  [string]: string,
};
