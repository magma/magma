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
 * @flow strict-local
 * @format
 */

import MagmaV1API from '../../server/magma/index';
import getCwfAlerts from './cwfAlerts';
import getFegAlerts from './fegAlerts';
import getLteAlerts from './lteAlerts';

// $FlowFixMe migrated to typescript
import {CWF, FEG, FEG_LTE, LTE} from '../../shared/types/network';
import type {ExpressResponse} from 'express';
import type {FBCNMSRequest} from '../../server/auth/access';
import type {
  network_type,
  prom_alert_config,
} from '../../generated/MagmaAPIBindings';

async function syncAlertsForNetwork(
  networkID: string,
  autoAlerts: {[string]: prom_alert_config},
) {
  // Get currently configured alerts
  const alerts = await MagmaV1API.getNetworksByNetworkIdPrometheusAlertConfig({
    networkId: networkID,
  });

  const existingAlerts = alerts.reduce(
    (map, obj) => ((map[obj.alert] = obj), map),
    {},
  );

  const putAlerts: prom_alert_config[] = [];
  const postAlerts: prom_alert_config[] = [];
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
      MagmaV1API.postNetworksByNetworkIdPrometheusAlertConfig({
        networkId: networkID,
        alertConfig: alert,
      }),
    );
  }
  for (const alert of putAlerts) {
    requests.push(
      MagmaV1API.putNetworksByNetworkIdPrometheusAlertConfigByAlertName({
        networkId: networkID,
        alertName: alert.alert,
        alertConfig: alert,
      }),
    );
  }

  await Promise.all(requests).catch(error => {
    throw error.message;
  });
}

async function syncAlerts(req: FBCNMSRequest, res: ExpressResponse) {
  try {
    const networkID = req.params.networkID;
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
    res.status(500).end('Exception occurred');
  }
}

async function getNetworkType(networkId: string): Promise<?network_type> {
  const networkInfo = await MagmaV1API.getNetworksByNetworkId({networkId});
  return networkInfo.type;
}

export default syncAlerts;
