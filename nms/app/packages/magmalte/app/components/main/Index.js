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

import MagmaV1API from '@fbcnms/magma-api/client/WebClient';

import {LteContextProvider} from '../lte/LteContext';
import {coalesceNetworkType} from '@fbcnms/types/network';
import type {NetworkType} from '@fbcnms/types/network';
import type {Theme} from '@material-ui/core';

import * as React from 'react';
import AppContent from '../layout/AppContent';
import AppContext from '@fbcnms/ui/context/AppContext';
import AppSideBar from '@fbcnms/ui/components/layout/AppSideBar';
import NetworkContext from '../context/NetworkContext';
import NetworkSelector from '../NetworkSelector';
import SectionLinks from '../layout/SectionLinks';
import SectionRoutes from '../layout/SectionRoutes';
import VersionTooltip from '../VersionTooltip';

import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import {getProjectLinks} from '@fbcnms/projects/projects';
import {makeStyles} from '@material-ui/styles';
import {shouldShowSettings} from '../Settings';
import {useContext, useEffect, useState} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';

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

export default function Index() {
  const classes = useStyles();
  const {match} = useRouter();
  const {user, tabs, ssoEnabled} = useContext(AppContext);
  const [networkType, setNetworkType] = useState<?NetworkType>(null);
  const networkId = ROOT_PATHS.has(match.params.networkId)
    ? null
    : match.params.networkId;

  useEffect(() => {
    const fetchNetworkType = async () => {
      if (networkId) {
        const networkType = await MagmaV1API.getNetworksByNetworkIdType({
          networkId,
        });
        setNetworkType(coalesceNetworkType(networkId, networkType));
      }
    };

    fetchNetworkType();
  }, [networkId]);

  if (networkId == null || networkType == null) {
    return <LoadingFiller />;
  }

  return (
    <NetworkContext.Provider value={{networkId, networkType}}>
      <LteContextProvider networkId={networkId} networkType={networkType}>
        <div className={classes.root}>
          <AppSideBar
            mainItems={[<SectionLinks key={1} />, <VersionTooltip key={2} />]}
            secondaryItems={[<NetworkSelector key={1} />]}
            projects={getProjectLinks(tabs, user)}
            showSettings={shouldShowSettings({
              isSuperUser: user.isSuperUser,
              ssoEnabled,
            })}
            user={user}
          />
          <AppContent>
            <SectionRoutes />
          </AppContent>
        </div>
      </LteContextProvider>
    </NetworkContext.Provider>
  );
}
