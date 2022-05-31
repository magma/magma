/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

export type FiringAlarm = {
  annotations: Labels;
  endsAt: string;
  fingerprint: string;
  labels: Labels;
  receivers: Array<GettableReceiver>;
  startsAt: string;
  status: FiringAlarmStatus;
  updatedAt: string;
  generatorURL?: string;
};

type GettableReceiver = {
  name: string;
};

export type FiringAlarmStatus = {
  inhibitedBy: Array<string>;
  silencedBy: Array<string>;
  state: string;
};

export type AlertReceiver = {
  name: string;
  email_configs?: Array<ReceiverEmailConfig>;
  slack_configs?: Array<ReceiverSlackConfig>;
  webhook_configs?: Array<ReceiverWebhookConfig>;
  pagerduty_configs?: Array<ReceiverPagerDutyConfig>;
  pushover_configs?: Array<ReceiverPushoverConfig>;
};

// names of all the <type>_configs lists for a receiver
export type ReceiverConfigListName =
  | 'email_configs'
  | 'slack_configs'
  | 'webhook_configs'
  | 'pagerduty_configs'
  | 'pushover_configs';

export type ReceiverEmailConfig = {
  auth_identity?: string;
  auth_password?: string;
  auth_secret?: string;
  auth_username?: string;
  from: string;
  headers?: Record<string, string>;
  hello?: string;
  html?: string;
  smarthost: string;
  text?: string;
  to: string;
  send_resolved?: boolean;
  require_tls?: boolean;
};

export type ReceiverSlackConfig = {
  api_url: string;
  channel?: string;
  username?: string;
  color?: string;
  title?: string;
  pretext?: string;
  text?: string;
  fields?: Array<SlackField>;
  short_fields?: boolean;
  footer?: string;
  fallback?: string;
  callback_id?: string;
  icon_emoji?: string;
  icon_url?: string;
  image_url?: string;
  thumb_url?: string;
  link_names?: boolean;
  actions?: Array<SlackAction>;
};

export type SlackField = {
  title: string;
  value: string;
  short?: boolean;
};

export type SlackAction = {
  type: string;
  text: string;
  url: string;
  style?: string;
  name?: string;
  value?: string;
  confirm?: SlackConfirmField;
};

export type SlackConfirmField = {
  text: string;
  title: string;
  ok_text: string;
  dismiss_text: string;
};

export type ReceiverWebhookConfig = {
  http_config?: HTTPConfig;
  send_resolved?: boolean;
  url: string;
};

export type HTTPConfig = {
  basic_auth?: HTTPConfigBasicAuth;
  bearer_token?: string;
  proxy_url?: string;
};

export type HTTPConfigBasicAuth = {
  password: string;
  username: string;
};

export type ReceiverPagerDutyConfig = {
  send_resolved?: boolean;
  routing_key: string;
  service_key: string;
  url: string;
  client: string;
  client_url: string;
  description: string;
  severity: string;
};

export type ReceiverPushoverConfig = {
  send_resolved: boolean;
  user_key: string;
  token: string;
  title: string;
  message: string;
  url: string;
  priority: string;
  retry: string;
  expire: string;
  http_config: HTTPConfig;
};

/**
 * Prometheus alert rule configuration
 */
export type AlertConfig = {
  alert: string;
  expr: string;
  annotations?: Record<string, string>;
  for?: string;
  labels?: Labels;
  rawData?: AlertConfig;
  _isCustomAlertRule?: boolean;
};

export type AlertRoutingTree = {
  receiver: string;
  continue?: boolean;
  group_by?: Array<string>;
  group_interval?: string;
  group_wait?: string;
  match?: Record<string, string>;
  match_re?: Record<string, string>;
  repeat_interval?: string;
  routes?: Array<AlertRoutingTree>;
};

export type BulkAlertUpdateResponse = {
  errors: Record<string, string>;
  statuses: Record<string, string>;
};

export type AlertSuppressionMatcher = {
  name: string;
  value: string;
  isRegex: boolean;
};

export type AlertSuppressionState = {
  state: string;
};

export type AlertSuppression = {
  id: string;
  startsAt: string;
  endsAt: string;
  updatedAt: string;
  matchers: Array<AlertSuppressionMatcher>;
  createdBy: string;
  status: AlertSuppressionState;
  comment?: string;
};

export type Labels = Record<string, string>;

export type PrometheusLabelset = Record<string, string>;

export type AlertManagerGlobalConfig = {
  resolve_timeout: string;
  http_config: HTTPConfig;
  smtp_from: string;
  smtp_hello: string;
  smtp_smarthost: string;
  smtp_auth_username: string;
  smtp_auth_password: string;
  smtp_auth_secret: string;
  smtp_auth_identity: string;
  smtp_require_tls: boolean;
  slack_api_url: string;
  pagerduty_url: string;
  hipchat_api_url: string;
  hipchat_auth_token: string;
  opsgenie_api_url: string;
  opsgenie_api_key: string;
  wechat_api_url: string;
  wechat_api_secret: string;
  wechat_api_corp_id: string;
  victorops_api_url: string;
  victorops_api_key: string;
};

export type TenancyConfig = {
  restrictor_label: string;
  restrict_queries: boolean;
};
