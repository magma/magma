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
 *
 * @flow strict-local
 * @format
 */
import type {apn, policy_rule} from '@fbcnms/magma-api';

import ApnOverview from './ApnOverview';
import LibraryBooksIcon from '@material-ui/icons/LibraryBooks';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import PolicyOverview from './PolicyOverview';
import React from 'react';
import RssFeedIcon from '@material-ui/icons/RssFeed';
import TopBar from '../../components/TopBar';
import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';

import {ApnJsonConfig} from './ApnOverview';
import {PolicyJsonConfig} from './PolicyOverview';
import {Redirect, Route, Switch} from 'react-router-dom';
import {useCallback, useState} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';

export default function TrafficDashboard() {
  const {relativePath, relativeUrl, match} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);
  const [policies, setPolicies] = useState<{[string]: policy_rule}>({});
  const [apns, setApns] = useState<{[string]: apn}>({});
  const {isLoading: policyLoading} = useMagmaAPI(
    MagmaV1API.getNetworksByNetworkIdPoliciesRulesViewFull,
    {
      networkId: networkId,
    },
    useCallback(response => {
      setPolicies(response);
    }, []),
  );

  const {isLoading: apnLoading} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdApns,
    {
      networkId: networkId,
    },
    useCallback(response => {
      setApns(response);
    }, []),
  );
  if (policyLoading || apnLoading) {
    return <LoadingFiller />;
  }
  return (
    <>
      <TopBar
        header="Traffic"
        tabs={[
          {
            label: 'Policies',
            to: '/policy',
            icon: LibraryBooksIcon,
          },
          {
            label: 'APNs',
            to: '/apn',
            icon: RssFeedIcon,
          },
        ]}
      />

      <Switch>
        <Route
          path={relativePath('/policy/:policyId/json')}
          render={() => (
            <PolicyJsonConfig
              policies={policies}
              onSave={policy => setPolicies({...policies, [policy.id]: policy})}
            />
          )}
        />
        <Route
          path={relativePath('/policy/json')}
          render={() => (
            <PolicyJsonConfig
              policies={policies}
              onSave={policy => setPolicies({...policies, [policy.id]: policy})}
            />
          )}
        />
        <Route
          path={relativePath('/apn/:apnId/json')}
          render={() => (
            <ApnJsonConfig
              apns={apns}
              onSave={apn => setApns({...apns, [apn.apn_name]: apn})}
            />
          )}
        />
        <Route
          path={relativePath('/apn/json')}
          render={() => (
            <ApnJsonConfig
              apns={apns}
              onSave={apn => setApns({...apns, [apn.apn_name]: apn})}
            />
          )}
        />
        <Route
          path={relativePath('/policy')}
          render={() => (
            <PolicyOverview
              policies={policies}
              onDelete={policyId => {
                const {
                  [policyId]: _deletedPolicy,
                  ...updatedPolicies
                } = policies;
                setPolicies(updatedPolicies);
              }}
            />
          )}
        />
        <Route
          path={relativePath('/apn')}
          render={() => <ApnOverview apns={apns} />}
        />
        <Redirect to={relativeUrl('/policy')} />
      </Switch>
    </>
  );
}
