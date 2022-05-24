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
import type {EditProps} from './GatewayDetailConfigEdit';
// $FlowFixMe migrated to typescript
import type {EnodebInfo} from '../../components/lte/EnodebUtils';
import type {lte_gateway} from '../../../generated/MagmaAPIBindings';

import ActionTable from '../../components/ActionTable';
import AddEditGatewayButton from './GatewayDetailConfigEdit';
import Button from '@material-ui/core/Button';
import CardTitleRow from '../../components/layout/CardTitleRow';
import DataGrid from '../../components/DataGrid';
import EnodebContext from '../../components/context/EnodebContext';
import GatewayContext from '../../components/context/GatewayContext';
import Grid from '@material-ui/core/Grid';
import JsonEditor from '../../components/JsonEditor';
import React from 'react';
import SettingsIcon from '@material-ui/icons/Settings';
// $FlowFixMe migrated to typescript
import nullthrows from '../../../shared/util/nullthrows';

import {DynamicServices} from '../../components/GatewayUtils';
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useContext, useState} from 'react';
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';
import {useNavigate, useParams} from 'react-router-dom';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(3),
    flexGrow: 1,
  },
  appBarBtn: {
    color: colors.primary.white,
    background: colors.primary.comet,
    fontFamily: typography.button.fontFamily,
    fontWeight: typography.button.fontWeight,
    fontSize: typography.button.fontSize,
    lineHeight: typography.button.lineHeight,
    letterSpacing: typography.button.letterSpacing,

    '&:hover': {
      background: colors.primary.mirage,
    },
  },
}));

export function GatewayJsonConfig() {
  const navigate = useNavigate();
  const params = useParams();
  const [error, setError] = useState('');
  const gatewayId: string = nullthrows(params.gatewayId);
  const enqueueSnackbar = useEnqueueSnackbar();
  const ctx = useContext(GatewayContext);
  const gwInfo = ctx.state[gatewayId];
  const {['status']: _status, ...gwInfoJson} = gwInfo;
  return (
    <JsonEditor
      content={{
        ...gwInfoJson,
        connected_enodeb_serials: gwInfoJson.connected_enodeb_serials ?? [],
      }}
      error={error}
      onSave={async gateway => {
        try {
          await ctx.setState(gatewayId, gateway);
          enqueueSnackbar('Gateway saved successfully', {
            variant: 'success',
          });
          setError('');
          navigate(-1);
        } catch (e) {
          setError(e.response?.data?.message ?? e.message);
        }
      }}
    />
  );
}

export default function GatewayConfig() {
  const classes = useStyles();
  const navigate = useNavigate();
  const params = useParams();
  const gatewayId: string = nullthrows(params.gatewayId);
  const ctx = useContext(GatewayContext);
  const gwInfo = ctx.state[gatewayId];
  function ConfigFilter() {
    return (
      <Button className={classes.appBarBtn} onClick={() => navigate('json')}>
        Edit JSON
      </Button>
    );
  }

  function editFilter(editTab: EditProps) {
    return (
      <AddEditGatewayButton title={'Edit'} isLink={true} editProps={editTab} />
    );
  }

  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <Grid item xs={12}>
            <CardTitleRow
              icon={SettingsIcon}
              label="Config"
              filter={ConfigFilter}
            />
          </Grid>
          <Grid container spacing={4}>
            <Grid item xs={12} md={6} alignItems="center">
              <Grid container spacing={4}>
                <Grid item xs={12}>
                  <CardTitleRow
                    label="Gateway"
                    filter={() => editFilter({editTable: 'info'})}
                  />
                  <GatewayInfoConfig gwInfo={gwInfo} />
                </Grid>

                <Grid item xs={12}>
                  <CardTitleRow
                    label="Dynamic Services"
                    filter={() => editFilter({editTable: 'aggregation'})}
                  />
                  <GatewayDynamicServices gwInfo={gwInfo} />
                </Grid>

                {Object.keys(gwInfo.apn_resources || {}).length > 0 && (
                  <Grid item xs={12}>
                    <CardTitleRow
                      label="Apn Resources"
                      filter={() => editFilter({editTable: 'apnResources'})}
                    />
                    <ApnResourcesTable gwInfo={gwInfo} />
                  </Grid>
                )}
              </Grid>
            </Grid>
            <Grid item xs={12} md={6} alignItems="center">
              <Grid container spacing={4}>
                <Grid item xs={12}>
                  <CardTitleRow
                    label="EPC"
                    filter={() => editFilter({editTable: 'epc'})}
                  />
                  <GatewayEPC gwInfo={gwInfo} />
                </Grid>
                <Grid item xs={12}>
                  <CardTitleRow
                    label="Ran"
                    filter={() => editFilter({editTable: 'ran'})}
                  />
                  <GatewayRAN gwInfo={gwInfo} />
                </Grid>
                <Grid item xs={12}>
                  <CardTitleRow
                    label="Header Enrichment"
                    filter={() => editFilter({editTable: 'headerEnrichment'})}
                  />
                  <GatewayHE gwInfo={gwInfo} />
                </Grid>
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    </div>
  );
}

function GatewayInfoConfig({gwInfo}: {gwInfo: lte_gateway}) {
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

function GatewayEPC({gwInfo}: {gwInfo: lte_gateway}) {
  const collapse: DataRows[] = [
    [
      {
        category: 'IP Block',
        value: gwInfo.cellular.epc.ip_block ?? '-',
      },
      {
        category: 'IPv6 Block',
        value: gwInfo.cellular.epc.ipv6_block ?? '-',
      },
    ],
  ];

  const data: DataRows[] = [
    [
      {
        category: 'IP Allocation',
        value: gwInfo.cellular.epc.nat_enabled ? 'NAT' : 'Custom',
        collapse: <DataGrid data={collapse} />,
      },
    ],
    [
      {
        category: 'Primary DNS',
        value: gwInfo.cellular.epc.dns_primary ?? '-',
      },
    ],
    [
      {
        category: 'Secondary DNS',
        value: gwInfo.cellular.epc.dns_secondary ?? '-',
      },
    ],
  ];

  return <DataGrid data={data} />;
}

function GatewayDynamicServices({gwInfo}: {gwInfo: lte_gateway}) {
  const logAggregation = !!gwInfo.magmad.dynamic_services?.includes(
    DynamicServices.TD_AGENT_BIT,
  );
  const eventAggregation = !!gwInfo.magmad?.dynamic_services?.includes(
    DynamicServices.EVENTD,
  );
  const cpeMonitoring = !!gwInfo.magmad?.dynamic_services?.includes(
    DynamicServices.MONITORD,
  );
  const dynamicServices: DataRows[] = [
    [
      {
        category: 'Log Aggregation',
        value: logAggregation ? 'Enabled' : 'Disabled',
        statusCircle: true,
        status: logAggregation,
      },
      {
        category: 'Event Aggregation',
        value: eventAggregation ? 'Enabled' : 'Disabled',
        statusCircle: true,
        status: eventAggregation,
      },
      {
        category: 'CPE Monitoring',
        value: cpeMonitoring ? 'Enabled' : 'Disabled',
        statusCircle: true,
        status: cpeMonitoring,
      },
    ],
  ];

  return <DataGrid data={dynamicServices} />;
}

function EnodebsTable({enbInfo}: {enbInfo: {[string]: EnodebInfo}}) {
  type EnodebRowType = {
    name: string,
    id: string,
  };
  const enbRows: Array<EnodebRowType> = Object.keys(enbInfo).map(
    (serialNum: string) => {
      const enbInf = enbInfo[serialNum];
      return {
        name: enbInf.enb.name,
        id: serialNum,
      };
    },
  );

  return (
    <ActionTable
      title=""
      data={enbRows}
      columns={[{title: 'Serial Number', field: 'id'}]}
      menuItems={[
        {name: 'View'},
        {name: 'Edit'},
        {name: 'Remove'},
        {name: 'Deactivate'},
        {name: 'Reboot'},
      ]}
      options={{
        actionsColumnIndex: -1,
        pageSizeOptions: [5],
        toolbar: false,
        header: false,
        paging: false,
      }}
    />
  );
}

function GatewayRAN({gwInfo}: {gwInfo: lte_gateway}) {
  const enbCtx = useContext(EnodebContext);
  const enbInfo =
    gwInfo.connected_enodeb_serials?.reduce(
      (enbs: {[string]: EnodebInfo}, serial: string) => {
        if (enbCtx?.state?.enbInfo?.[serial] != null) {
          enbs[serial] = enbCtx.state.enbInfo[serial];
        }
        return enbs;
      },
      {},
    ) || {};
  const dhcpServiceStatus = gwInfo.cellular.dns?.dhcp_server_enabled ?? true;
  const ran: DataRows[] = [
    [
      {
        category: 'eNodeB DHCP Service',
        value: dhcpServiceStatus ? 'Enabled' : 'Disabled',
        statusCircle: true,
        status: dhcpServiceStatus,
      },
    ],
    [
      {
        category: 'PCI',
        value: gwInfo.cellular.ran.pci,
        statusCircle: false,
      },
      {
        category: 'eNodeB Transmit',
        value: gwInfo.cellular.ran.transmit_enabled ? 'Enabled' : 'Disabled',
        statusCircle: false,
      },
    ],
    [
      {
        category: 'Registered eNodeBs',
        value: gwInfo.connected_enodeb_serials?.length || 0,
        collapse: <EnodebsTable gwInfo={gwInfo} enbInfo={enbInfo} />,
      },
    ],
  ];

  return <DataGrid data={ran} />;
}

function ApnResourcesTable({gwInfo}: {gwInfo: lte_gateway}) {
  const apnResources = gwInfo.apn_resources || {};
  type ApnResourcesRowType = {
    name: string,
    id: string,
    vlanId: number | string,
  };
  const apnResourcesRows: Array<ApnResourcesRowType> = Object.keys(
    apnResources,
  ).map((apn: string) => {
    const apnRow = apnResources[apn];
    return {
      name: apn,
      id: apnRow.id,
      vlanId: apnRow.vlan_id ?? '-',
    };
  });

  return (
    <ActionTable
      title=""
      data={apnResourcesRows}
      columns={[
        {title: 'Name', field: 'name'},
        {title: 'Resource ID', field: 'id'},
        {title: 'VLAN ID', field: 'vlanId'},
      ]}
      options={{
        actionsColumnIndex: -1,
        pageSizeOptions: [5],
        toolbar: false,
      }}
    />
  );
}

function GatewayHE({gwInfo}: {gwInfo: lte_gateway}) {
  const heEnabled =
    gwInfo.cellular.he_config?.enable_header_enrichment ?? false;
  const encryptionEnabled =
    gwInfo.cellular.he_config?.enable_encryption ?? false;
  const EncryptionDetail = () => {
    const encryptionConfig: DataRows[] = [
      [
        {
          category: 'Encryption Key',
          value: gwInfo.cellular.he_config?.encryption_key || '',
          obscure: true,
        },
        {
          category: 'Encoding Type',
          value: gwInfo.cellular.he_config?.he_encoding_type || '',
        },
      ],
      [
        {
          category: 'Encryption Algorithm',
          value: gwInfo.cellular.he_config?.he_encryption_algorithm || '',
        },
        {
          category: 'Hash Function',
          value: gwInfo.cellular.he_config?.he_hash_function || '',
        },
      ],
    ];
    return <DataGrid data={encryptionConfig} />;
  };

  const heConfig: DataRows[] = [
    [
      {
        statusCircle: true,
        status: heEnabled,
        category: 'Header Enrichment',
        value: heEnabled ? 'Enabled' : 'Disabled',
        collapse: encryptionEnabled ? <EncryptionDetail /> : <></>,
      },
    ],
  ];
  return <DataGrid data={heConfig} />;
}
