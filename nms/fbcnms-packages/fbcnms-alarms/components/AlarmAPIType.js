/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
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
  slack_configs?: Array<ReceiverSlackConfig>,
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

/**
 * Prometheus alert rule configuration
 */
export type AlertConfig = {
  alert: string,
  expr: string,
  annotations?: {[string]: string},
  for?: string,
  labels?: {[string]: string},
  rawData?: AlertConfig,
  _isCustomAlertRule?: boolean,
};

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
