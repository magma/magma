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

import {
  FEG,
  NetworkId,
  coalesceNetworkType,
} from '../../../shared/types/network';
import {FEGContextProvider} from '../feg/FEGContext';
import {LteContextProvider} from '../lte/LteContext';
import {VersionContextProvider} from '../context/VersionContext';
import type {NetworkType} from '../../../shared/types/network';
import type {Theme} from '@material-ui/core';

import * as React from 'react';
import AppContent from '../layout/AppContent';
import AppSideBar from '../AppSideBar';
import NetworkContext from '../context/NetworkContext';
import SectionRoutes from '../layout/SectionRoutes';
import {useEffect, useState} from 'react';

import LoadingFiller from '../LoadingFiller';
import MagmaAPI from '../../../api/MagmaAPI';
import useSections from '../layout/useSections';
import {makeStyles} from '@material-ui/styles';
import {useParams} from 'react-router-dom';

// These won't be considered networkIds
export const ROOT_PATHS = new Set<string>(['network']);

const useStyles = makeStyles((theme: Theme) => ({
  root: {
    display: 'flex',
  },
  toolbarIcon: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'flex-end',
    padding: '0 8px',
    ...theme.mixins.toolbar,
  },
}));

type Props = {
  networkId: NetworkId;
  networkType: NetworkType;
  children: React.ReactNode;
};

function Sidebar() {
  const [, sections] = useSections();
  return <AppSideBar items={[...sections]} />;
}

export default function Index() {
  const classes = useStyles();
  const params = useParams();
  const [networkType, setNetworkType] = useState<NetworkType | null>(null);
  const networkId = ROOT_PATHS.has(params.networkId!) ? null : params.networkId;

  useEffect(() => {
    const fetchNetworkType = async () => {
      if (networkId) {
        const networkType = (
          await MagmaAPI.networks.networksNetworkIdTypeGet({
            networkId,
          })
        ).data;
        setNetworkType(coalesceNetworkType(networkId, networkType));
      }
    };

    void fetchNetworkType();
  }, [networkId]);

  if (networkId == null || networkType == null) {
    return <LoadingFiller />;
  }

  return (
    <NetworkContextProvider {...{networkId, networkType}}>
      <div className={classes.root}>
        <Sidebar />
        <AppContent>
          <SectionRoutes />
        </AppContent>
      </div>
    </NetworkContextProvider>
  );
}

/**
 * Returns a Federation context provider if it is a federation network. It
 * otherwise returns a LTE context provider for a LTE or Federated LTE network.
 *
 * @param {network_id} network_id Id of the network
 * @param {network_type} network_type Type of the network
 */
function NetworkContextProvider(props: Props) {
  const {networkId, networkType} = props;

  return (
    <VersionContextProvider>
      <NetworkContext.Provider value={{networkId, networkType}}>
        {networkType === FEG ? (
          <FEGContextProvider networkId={networkId} networkType={networkType}>
            {props.children}
          </FEGContextProvider>
        ) : (
          <LteContextProvider networkId={networkId} networkType={networkType}>
            {props.children}
          </LteContextProvider>
        )}
      </NetworkContext.Provider>
    </VersionContextProvider>
  );
}
