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
import type {
  network,
  network_epc_configs,
  network_ran_configs,
} from '@fbcnms/magma-api';

import AddEditNetworkButton from './NetworkEdit';
import CardTitleRow from '../../components/layout/CardTitleRow';
import Grid from '@material-ui/core/Grid';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import NetworkEpc from './NetworkEpc';
import NetworkInfo from './NetworkInfo';
import NetworkKPI from './NetworkKPIs';
import NetworkRanConfig from './NetworkRanConfig';
import React from 'react';
import TopBar from '../../components/TopBar';
import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';

import {NetworkCheck} from '@material-ui/icons';
import {Redirect, Route, Switch} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useState} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(5),
  },
}));

export default function NetworkDashboard() {
  const {relativePath, relativeUrl} = useRouter();

  return (
    <>
      <TopBar
        header="Network"
        tabs={[
          {
            label: 'Network',
            to: '/network',
            icon: NetworkCheck,
            filters: (
              <AddEditNetworkButton title={'Add Network'} isLink={false} />
            ),
          },
        ]}
      />

      <Switch>
        <Route
          path={relativePath('/network')}
          component={NetworkDashboardInternal}
        />
        <Redirect to={relativeUrl('/network')} />
      </Switch>
    </>
  );
}

function NetworkDashboardInternal() {
  const {match} = useRouter();
  const classes = useStyles();
  const networkId: string = nullthrows(match.params.networkId);

  const [networkInfo, setNetworkInfo] = useState<network>({});
  const [epcConfigs, setEpcConfigs] = useState<network_epc_configs>({});
  const [lteRanConfigs, setLteRanConfigs] = useState<network_ran_configs>({});

  const {isLoading: isInfoLoading} = useMagmaAPI(
    MagmaV1API.getNetworksByNetworkId,
    {
      networkId: networkId,
    },
    useCallback(networkInfo => {
      setNetworkInfo(networkInfo);
    }, []),
  );
  const {isLoading: isEpcLoading} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdCellularEpc,
    {
      networkId: networkId,
    },
    useCallback(epc => setEpcConfigs(epc), []),
  );

  const {isLoading: isRanLoading} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdCellularRan,
    {
      networkId: networkId,
    },
    useCallback(lteRanConfigs => setLteRanConfigs(lteRanConfigs), []),
  );

  const {response: lteGatwayResp, isLoading: isLteRespLoading} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdGateways,
    {
      networkId: networkId,
    },
  );

  const {response: enb, isLoading: isEnbRespLoading} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdEnodebs,
    {
      networkId: networkId,
    },
  );

  const {response: policyRules, isLoading: isPolicyLoading} = useMagmaAPI(
    MagmaV1API.getNetworksByNetworkIdPoliciesRules,
    {
      networkId: networkId,
    },
  );

  const {response: subscriber, isLoading: isSubscriberLoading} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdSubscribers,
    {
      networkId: networkId,
    },
  );

  const {response: apns, isLoading: isAPNsLoading} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdApns,
    {
      networkId: networkId,
    },
  );

  if (
    isEpcLoading ||
    isInfoLoading ||
    isRanLoading ||
    isLteRespLoading ||
    isEnbRespLoading ||
    isPolicyLoading ||
    isSubscriberLoading ||
    isAPNsLoading
  ) {
    return <LoadingFiller />;
  }
  const editProps = {
    networkInfo: networkInfo,
    lteRanConfigs: lteRanConfigs,
    epcConfigs: epcConfigs,
    onSaveNetworkInfo: setNetworkInfo,
    onSaveEpcConfigs: setEpcConfigs,
    onSaveLteRanConfigs: setLteRanConfigs,
  };

  function editNetwork() {
    return (
      <AddEditNetworkButton
        title={'Edit'}
        isLink={true}
        editProps={{
          editTable: 'info',
          ...editProps,
        }}
      />
    );
  }

  function editRAN() {
    return (
      <AddEditNetworkButton
        title={'Edit'}
        isLink={true}
        editProps={{
          editTable: 'ran',
          ...editProps,
        }}
      />
    );
  }

  function editEPC() {
    return (
      <AddEditNetworkButton
        title={'Edit'}
        isLink={true}
        editProps={{
          editTable: 'epc',
          ...editProps,
        }}
      />
    );
  }

  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <CardTitleRow label="Overview" />
          <NetworkKPI
            apns={apns}
            enb={enb}
            lteGatwayResp={lteGatwayResp}
            policyRules={policyRules}
            subscriber={subscriber}
          />
        </Grid>
        <Grid item xs={12} md={6}>
          <Grid container spacing={4}>
            <Grid item xs={12}>
              <CardTitleRow label="Network" filter={editNetwork} />
              <NetworkInfo networkInfo={networkInfo} />
            </Grid>
            <Grid item xs={12}>
              <CardTitleRow label="RAN" filter={editRAN} />
              <NetworkRanConfig lteRanConfigs={lteRanConfigs} />
            </Grid>
          </Grid>
        </Grid>
        <Grid item xs={12} md={6}>
          <Grid container spacing={4}>
            <Grid item xs={12}>
              <CardTitleRow label="EPC" filter={editEPC} />
              <NetworkEpc epcConfigs={epcConfigs} />
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    </div>
  );
}
