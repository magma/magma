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
 *
 */
import MagmaAPI from '../../api/MagmaAPI';
import getCwfAlerts from './cwfAlerts';
import getFegAlerts from './fegAlerts';
import getLteAlerts from './lteAlerts';

// $FlowFixMe migrated to typescript
import {CWF, FEG, FEG_LTE, LTE} from '../../shared/types/network';
import type {PromAlertConfig} from '../../generated-ts';
import type {Response} from 'express';

type OutputRequest<T> = {
  logIn: (user: T, callback: (err?: Error | null | undefined) => void) => void;
  logOut: () => void;
  logout: () => void;
  user: T;
  isAuthenticated: () => boolean;
  isUnauthenticated: () => boolean;
};
export type FBCNMSPassportRequest = OutputRequest<any>;
type Options = {loginUrl: string};
export type FBCNMSRequest = FBCNMSPassportRequest & {access: Options};

async function syncAlertsForNetwork(
  networkID: string,
  autoAlerts: {[name: string]: PromAlertConfig},
) {
  // Get currently configured alerts
  const alerts = await MagmaAPI.alerts.networksNetworkIdPrometheusAlertConfigGet(
    {
      networkId: networkID,
    },
  );

  const existingAlerts: {[name: string]: PromAlertConfig} = alerts.data.reduce(
    (map: {[name: string]: PromAlertConfig}, obj: PromAlertConfig) => (
      (map[obj.alert] = obj), map
    ),
    {},
  );

  const putAlerts: Array<PromAlertConfig> = [];
  const postAlerts: Array<PromAlertConfig> = [];
  for (const alertName in autoAlerts) {
    if (existingAlerts[alertName] !== undefined) {
      putAlerts.push(autoAlerts[alertName]);
    } else {
      postAlerts.push(autoAlerts[alertName]);
    }
  }

  const requests = [];
  for (const alert of postAlerts) {
    requests.push(
      MagmaAPI.alerts.networksNetworkIdPrometheusAlertConfigPost({
        networkId: networkID,
        alertConfig: alert,
      }),
    );
  }
  for (const alert of putAlerts) {
    requests.push(
      MagmaAPI.alerts.networksNetworkIdPrometheusAlertConfigAlertNamePut({
        networkId: networkID,
        alertName: alert.alert,
        alertConfig: alert,
      }),
    );
  }

  await Promise.all(requests).catch((error: {message: string}) => {
    throw error.message;
  });
}

async function syncAlerts(networkID: string, res: Response) {
  try {
    const type = await getNetworkType(networkID);
    if (type == null) {
      res.status(500).send(`Invalid network type`).end();
      return;
    }
    switch (type) {
      case CWF:
        await syncAlertsForNetwork(networkID, getCwfAlerts(networkID));
        break;
      case LTE:
        await syncAlertsForNetwork(networkID, getLteAlerts(networkID));
        break;
      case FEG_LTE:
        await syncAlertsForNetwork(networkID, getLteAlerts(networkID));
        break;
      case FEG:
        await syncAlertsForNetwork(networkID, getFegAlerts(networkID));
        break;
      default:
        res
          .status(400)
          .send(`Network type ${type} has no predefined alerts`)
          .end();
    }
    res.status(200).end();
  } catch (e) {
    const message = e instanceof Error ? e.message : 'unknown error';
    res.status(500).end(`Exception occurred ${message}`);
  }
}

async function getNetworkType(networkId: string): Promise<string | undefined> {
  const networkInfo = await MagmaAPI.networks.networksNetworkIdGet({networkId});
  return networkInfo.data.type;
}

export default syncAlerts;
