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

import type {DataRows} from '../../components/DataGrid';
import type {TabOption} from '../../components/feg/FEGGatewayDialog';
import type {
  diameter_client_configs,
  federation_gateway,
  sctp_client_configs,
} from '../../../generated/MagmaAPIBindings';

import CardTitleRow from '../../components/layout/CardTitleRow';
import DataGrid from '../../components/DataGrid';
import EditGatewayButton from './FEGGatewayDetailConfigEdit';
import FEGGatewayContext from '../../components/context/FEGGatewayContext';
import Grid from '@material-ui/core/Grid';
import React from 'react';
import nullthrows from '../../../shared/util/nullthrows';

import {TAB_OPTIONS} from '../../components/feg/FEGGatewayDialog';
import {makeStyles} from '@material-ui/styles';
import {useContext} from 'react';
import {useParams} from 'react-router-dom';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(3),
    flexGrow: 1,
  },
}));

/**
 * Returns the configuration page of the selected federation
 * gateway. It provides information about the federation gateway
 * and its servers such as gx, gy, and the like.
 */
export default function FEGGatewayConfig() {
  const classes = useStyles();
  const params = useParams();
  const gatewayId: string = nullthrows(params.gatewayId);
  const ctx = useContext(FEGGatewayContext);
  const gwInfo: federation_gateway = ctx.state[gatewayId];

  function editFilter(tabOption: TabOption) {
    return (
      <EditGatewayButton
        title={'Edit'}
        tabOption={tabOption}
        editingGateway={gwInfo}
      />
    );
  }

  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <Grid container spacing={4}>
            <Grid item xs={12} md={6} alignItems="center">
              <Grid container spacing={4}>
                <Grid item xs={12}>
                  <CardTitleRow label="Gateway" />
                  <GatewayInfoConfig gwInfo={gwInfo} />
                </Grid>
                <Grid item xs={12}>
                  <CardTitleRow
                    label="Gx"
                    filter={() => editFilter(TAB_OPTIONS.GX)}
                  />
                  <GatewayDiameterServerConfig
                    serverConfig={gwInfo?.federation?.gx?.server || {}}
                    testID={'Gx'}
                  />
                </Grid>
                <Grid item xs={12}>
                  <CardTitleRow
                    label="CSFB"
                    filter={() => editFilter(TAB_OPTIONS.CSFB)}
                  />
                  <GatewaySctpServerConfig
                    serverConfig={gwInfo?.federation?.csfb?.client || {}}
                    testID={'CSFB'}
                  />
                </Grid>
              </Grid>
            </Grid>
            <Grid item xs={12} md={6} alignItems="center">
              <Grid container spacing={4}>
                <Grid item xs={12}>
                  <CardTitleRow
                    label="Gy"
                    filter={() => editFilter(TAB_OPTIONS.GY)}
                  />
                  <GatewayDiameterServerConfig
                    serverConfig={gwInfo?.federation?.gy?.server || {}}
                    testID={'Gy'}
                  />
                </Grid>
                <Grid item xs={12}>
                  <CardTitleRow
                    label="SWx"
                    filter={() => editFilter(TAB_OPTIONS.SWX)}
                  />
                  <GatewayDiameterServerConfig
                    serverConfig={gwInfo?.federation?.swx?.server || {}}
                    testID={'SWx'}
                  />
                </Grid>
                <Grid item xs={12}>
                  <CardTitleRow
                    label="S6a"
                    filter={() => editFilter(TAB_OPTIONS.S6A)}
                  />
                  <GatewayDiameterServerConfig
                    serverConfig={gwInfo?.federation?.s6a?.server || {}}
                    testID={'S6a'}
                  />
                </Grid>
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    </div>
  );
}

/**
 * Returns useful information about the federation gateway. It returns
 * its name, id, hardware uuid, version and description.
 * @param {federation_gateway} gwInfo The federation gateway that is being looked at.
 */
function GatewayInfoConfig({gwInfo}: {gwInfo: federation_gateway}) {
  const data: DataRows[] = [
    [
      {
        category: 'Name',
        value: gwInfo.name,
      },
    ],
    [
      {
        category: 'Gateway ID',
        value: gwInfo.id,
      },
    ],
    [
      {
        category: 'Hardware UUID',
        value: gwInfo.device?.hardware_id || '-',
      },
    ],
    [
      {
        category: 'Version',
        value: gwInfo.status?.platform_info?.packages?.[0]?.version ?? 'null',
      },
    ],
    [
      {
        category: 'Description',
        value: gwInfo.description,
      },
    ],
  ];

  return <DataGrid data={data} />;
}

/**
 * Returns useful information about the federation gateway's diameter based
 * server.
 * @param {diameter_client_configs} serverConfig Configuration object of the diameter based server.
 * @param {string} testId An id used to differentiate the various diameter servers.
 */
function GatewayDiameterServerConfig({
  serverConfig,
  testID,
}: {
  serverConfig: diameter_client_configs,
  testID: string,
}) {
  const data: DataRows[] = [
    [
      {
        category: 'Address',
        value: serverConfig?.address || '-',
      },
    ],
    [
      {
        category: 'Destination Host',
        value: serverConfig?.dest_host || '-',
      },
    ],
    [
      {
        category: 'Destination Realm',
        value: serverConfig?.dest_realm || '-',
      },
    ],
    [
      {
        category: 'Host',
        value: serverConfig?.host || '-',
      },
    ],
    [
      {
        category: 'Realm',
        value: serverConfig?.realm || '-',
      },
    ],
    [
      {
        category: 'Local Address',
        value: serverConfig?.local_address || '-',
      },
    ],
    [
      {
        category: 'Product Name',
        value: serverConfig?.product_name || '-',
      },
    ],
    [
      {
        category: 'Protocol',
        value: serverConfig?.protocol || '-',
      },
    ],
  ];

  return <DataGrid data={data} testID={testID} />;
}

/**
 * Returns useful information about the federation gateway's sctp based
 * server.
 * @param {sctp_client_configs} serverConfig Configuration object of the sctp based server.
 * @param {string} testId An id used to differentiate the various servers.
 */
function GatewaySctpServerConfig({
  serverConfig,
  testID,
}: {
  serverConfig: sctp_client_configs,
  testID: string,
}) {
  const data: DataRows[] = [
    [
      {
        category: 'Local Address',
        value: serverConfig?.local_address || '-',
      },
    ],
    [
      {
        category: 'Server Address',
        value: serverConfig?.server_address || '-',
      },
    ],
  ];

  return <DataGrid data={data} testID={testID} />;
}
