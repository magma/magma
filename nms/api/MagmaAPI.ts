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

import nullthrows from '../shared/util/nullthrows';
import {API_HOST} from '../config/config';

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

const BASE_PATH_FRONTEND = '/nms/apicontroller/magma/v1';
const config = new Configuration();

export const BASE_API = new BaseAPI(config, BASE_PATH_FRONTEND);

/**
 * New API need to be added here
 */
function setUpApi(basePath: string) {
  return {
    apns: new APNsApi(config, basePath),
    about: new AboutApi(config, basePath),
    alerts: new AlertsApi(config, basePath),
    callTracing: new CallTracingApi(config, basePath),
    carrierWifiGateways: new CarrierWifiGatewaysApi(config, basePath),
    carrierWifiNetworks: new CarrierWifiNetworksApi(config, basePath),
    cbsds: new CbsdsApi(config, basePath),
    commands: new CommandsApi(config, basePath),
    default: new DefaultApi(config, basePath),
    enodebs: new EnodeBsApi(config, basePath),
    events: new EventsApi(config, basePath),
    federatedLTENetworks: new FederatedLTENetworksApi(config, basePath),
    federationGateways: new FederationGatewaysApi(config, basePath),
    federationNetworks: new FederationNetworksApi(config, basePath),
    gateways: new GatewaysApi(config, basePath),
    lteGateways: new LTEGatewaysApi(config, basePath),
    lteNetworks: new LTENetworksApi(config, basePath),
    logs: new LogsApi(config, basePath),
    metrics: new MetricsApi(config, basePath),
    networkProbes: new NetworkProbesApi(config, basePath),
    networks: new NetworksApi(config, basePath),
    policies: new PoliciesApi(config, basePath),
    ratingGroups: new RatingGroupsApi(config, basePath),
    sms: new SMSApi(config, basePath),
    subscribers: new SubscribersApi(config, basePath),
    tenants: new TenantsApi(config, basePath),
    upgrades: new UpgradesApi(config, basePath),
    user: new UserApi(config, basePath),
  };
}
const orchestratorUrl = !/^https?\:\/\//.test(nullthrows(API_HOST))
  ? `https://${nullthrows(API_HOST)}/magma/v1`
  : `${nullthrows(API_HOST)}/magma/v1`;

export default setUpApi(BASE_PATH_FRONTEND);
export const OrchestratorAPI = setUpApi(orchestratorUrl);
