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

import globalAxios from 'axios';
import https from 'https';
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
} from '../generated';
import {BaseAPI} from '../generated/base';

const config = new Configuration();

export function createBaseAPI(basePath: string) {
  return new BaseAPI(config, basePath);
}

/**
 * New APIs need to be added here
 */
export function setUpApi(basePath: string, httpsAgent?: https.Agent) {
  const axios = httpsAgent
    ? globalAxios.create({
        httpsAgent,
      })
    : globalAxios;

  return {
    apns: new APNsApi(config, basePath, axios),
    about: new AboutApi(config, basePath, axios),
    alerts: new AlertsApi(config, basePath, axios),
    callTracing: new CallTracingApi(config, basePath, axios),
    carrierWifiGateways: new CarrierWifiGatewaysApi(config, basePath, axios),
    carrierWifiNetworks: new CarrierWifiNetworksApi(config, basePath, axios),
    cbsds: new CbsdsApi(config, basePath, axios),
    commands: new CommandsApi(config, basePath, axios),
    default: new DefaultApi(config, basePath, axios),
    enodebs: new EnodeBsApi(config, basePath, axios),
    events: new EventsApi(config, basePath, axios),
    federatedLTENetworks: new FederatedLTENetworksApi(config, basePath, axios),
    federationGateways: new FederationGatewaysApi(config, basePath, axios),
    federationNetworks: new FederationNetworksApi(config, basePath, axios),
    gateways: new GatewaysApi(config, basePath, axios),
    lteGateways: new LTEGatewaysApi(config, basePath, axios),
    lteNetworks: new LTENetworksApi(config, basePath, axios),
    logs: new LogsApi(config, basePath, axios),
    metrics: new MetricsApi(config, basePath, axios),
    networkProbes: new NetworkProbesApi(config, basePath, axios),
    networks: new NetworksApi(config, basePath, axios),
    policies: new PoliciesApi(config, basePath, axios),
    ratingGroups: new RatingGroupsApi(config, basePath, axios),
    sms: new SMSApi(config, basePath, axios),
    subscribers: new SubscribersApi(config, basePath, axios),
    tenants: new TenantsApi(config, basePath, axios),
    upgrades: new UpgradesApi(config, basePath, axios),
    user: new UserApi(config, basePath, axios),
  };
}
