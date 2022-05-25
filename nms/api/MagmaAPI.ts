/**
 * Copyright 2022 The Magma Authors.
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

import {
  APNsApi,
  AboutApi,
  AlertsApi,
  CallTracingApi,
  CarrierWifiGatewaysApi,
  CarrierWifiNetworksApi,
  CbsdsApi,
  CommandsApi,
  Configuration,
  DefaultApi,
  EnodeBsApi,
  EventsApi,
  FederatedLTENetworksApi,
  FederationGatewaysApi,
  FederationNetworksApi,
  GatewaysApi,
  LTEGatewaysApi,
  LTENetworksApi,
  LogsApi,
  MetricsApi,
  NetworkProbesApi,
  NetworksApi,
  PoliciesApi,
  RatingGroupsApi,
  SMSApi,
  SubscribersApi,
  TenantsApi,
  UpgradesApi,
  UserApi,
} from '../generated-ts';

import {BaseAPI} from '../generated-ts/base';

const config = new Configuration();
const BASE_PATH = '/nms/apicontroller/magma/v1';

export const BASE_API = new BaseAPI(config, BASE_PATH);

/**
 * New API need to be added here
 */
export default {
  apns: new APNsApi(config, BASE_PATH),
  about: new AboutApi(config, BASE_PATH),
  alerts: new AlertsApi(config, BASE_PATH),
  callTracing: new CallTracingApi(config, BASE_PATH),
  carrierWifiGateways: new CarrierWifiGatewaysApi(config, BASE_PATH),
  carrierWifiNetworks: new CarrierWifiNetworksApi(config, BASE_PATH),
  cbsds: new CbsdsApi(config, BASE_PATH),
  commands: new CommandsApi(config, BASE_PATH),
  default: new DefaultApi(config, BASE_PATH),
  enodebs: new EnodeBsApi(config, BASE_PATH),
  events: new EventsApi(config, BASE_PATH),
  federatedLTENetworks: new FederatedLTENetworksApi(config, BASE_PATH),
  federationGateways: new FederationGatewaysApi(config, BASE_PATH),
  federationNetworks: new FederationNetworksApi(config, BASE_PATH),
  gateways: new GatewaysApi(config, BASE_PATH),
  lteGateways: new LTEGatewaysApi(config, BASE_PATH),
  lteNetworks: new LTENetworksApi(config, BASE_PATH),
  logs: new LogsApi(config, BASE_PATH),
  metrics: new MetricsApi(config, BASE_PATH),
  networkProbes: new NetworkProbesApi(config, BASE_PATH),
  networks: new NetworksApi(config, BASE_PATH),
  policies: new PoliciesApi(config, BASE_PATH),
  ratingGroups: new RatingGroupsApi(config, BASE_PATH),
  sms: new SMSApi(config, BASE_PATH),
  subscribers: new SubscribersApi(config, BASE_PATH),
  tenants: new TenantsApi(config, BASE_PATH),
  upgrades: new UpgradesApi(config, BASE_PATH),
  user: new UserApi(config, BASE_PATH),
};
