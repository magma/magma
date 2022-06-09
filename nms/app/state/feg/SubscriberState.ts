/*
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

import {FetchSubscriberState} from '../lte/SubscriberState';
import {NetworkId} from '../../../shared/types/network';
import {OptionsObject} from 'notistack';
import {SubscriberState} from '../../../generated-ts';
import {getServicedAccessNetworks} from '../../components/FEGServicingAccessGatewayKPIs';

/**
 * Props passed when fetching subscriber state.
 *
 * @param {NetworkId} networkId Id of the federation network.
 * @param {(msg, cfg,) => ?(string | number),} enqueueSnackbar Snackbar to display error.
 */
type FetchProps = {
  networkId: NetworkId;
  enqueueSnackbar?: (
    msg: string,
    cfg: OptionsObject,
  ) => string | number | null | undefined;
};

type FegSubscriberState = Record<NetworkId, Record<string, SubscriberState>>;

/**
 * Fetches and returns the subscriber session state of all the serviced
 * federated lte networks under by this federation network.
 *
 * @param {FetchProps} props an object containing the network id and snackbar to display error.
 * @returns {{[string]:{[string]: subscriber_state}}} returns an object containing the serviced
 *   network ids mapped to the each of their subscriber state. It returns an empty object and
 *   displays any error encountered on the snackbar when it fails to fetch the session state.
 */
export async function FetchFegSubscriberState(
  props: FetchProps,
): Promise<FegSubscriberState> {
  const {networkId, enqueueSnackbar} = props;
  const servicedAccessNetworks = await getServicedAccessNetworks(
    networkId,
    enqueueSnackbar,
  );
  const sessionState: FegSubscriberState = {};
  for (const servicedAccessNetwork of servicedAccessNetworks) {
    const servicedAccessNetworkId = servicedAccessNetwork.id;
    const state = await FetchSubscriberState({
      networkId: servicedAccessNetworkId,
      enqueueSnackbar,
    });
    // group session states under their network id
    sessionState[servicedAccessNetworkId] = state as Record<
      string,
      SubscriberState
    >;
  }
  return sessionState;
}
