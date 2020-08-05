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
import type {EnodebInfo} from '../../components/lte/EnodebUtils';
import type {KPIRows} from '../../components/KPIGrid';
import type {lte_gateway} from '@fbcnms/magma-api';

import ActionTable from '../../components/ActionTable';
import AddEditGatewayButton from './GatewayDetailConfigEdit';
import Button from '@material-ui/core/Button';
import CardHeader from '@material-ui/core/CardHeader';
import Collapse from '@material-ui/core/Collapse';
import Divider from '@material-ui/core/Divider';
import EnodebContext from '../../components/context/EnodebContext';
import ExpandLess from '@material-ui/icons/ExpandLess';
import ExpandMore from '@material-ui/icons/ExpandMore';
import GatewayContext from '../../components/context/GatewayContext';
import Grid from '@material-ui/core/Grid';
import JsonEditor from '../../components/JsonEditor';
import KPIGrid from '../../components/KPIGrid';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import SettingsIcon from '@material-ui/icons/Settings';
import nullthrows from '@fbcnms/util/nullthrows';

import {CardTitleFilterRow} from '../../components/layout/CardTitleRow';
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useContext, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(3),
    flexGrow: 1,
  },
  list: {
    padding: 0,
  },
  kpiLabel: {
    color: colors.primary.comet,
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
  },
  kpiValue: {
    color: colors.primary.brightGray,
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    width: '100%',
  },
  kpiBox: {
    width: '100%',
    padding: 0,
    '& > div': {
      width: '100%',
    },
  },
  paper: {
    height: 100,
    padding: theme.spacing(10),
    textAlign: 'center',
  },
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '50%',
  },
  itemTitle: {
    color: colors.primary.comet,
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
  },
  itemValue: {
    color: colors.primary.brightGray,
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
  const {match} = useRouter();
  const [error, setError] = useState('');
  const gatewayId: string = nullthrows(match.params.gatewayId);
  const enqueueSnackbar = useEnqueueSnackbar();
  const ctx = useContext(GatewayContext);
  const gwInfo = ctx.state[gatewayId];

  return (
    <JsonEditor
      content={gwInfo}
      error={error}
      onSave={async gateway => {
        try {
          await ctx.setState(gatewayId, gateway);
          enqueueSnackbar('Gateway saved successfully', {
            variant: 'success',
          });
          setError('');
        } catch (e) {
          setError(e.response?.data?.message ?? e.message);
        }
      }}
    />
  );
}

export default function GatewayConfig() {
  const classes = useStyles();
  const {history, match, relativeUrl} = useRouter();
  const gatewayId: string = nullthrows(match.params.gatewayId);
  const ctx = useContext(GatewayContext);
  const gwInfo = ctx.state[gatewayId];

  function ConfigFilter() {
    return (
      <Button
        className={classes.appBarBtn}
        onClick={() => {
          history.push(relativeUrl('/json'));
        }}>
        Edit JSON
      </Button>
    );
  }

  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <Grid item xs={12}>
            <CardTitleFilterRow
              icon={SettingsIcon}
              label="Config"
              filter={ConfigFilter}
            />
          </Grid>
          <Grid container spacing={4}>
            <Grid item xs={12} md={6} alignItems="center">
              <Grid container spacing={4}>
                <Grid item xs={12}>
                  <CardTitleFilterRow label="Gateway" />
                  <AddEditGatewayButton
                    title={'Edit'}
                    isLink={true}
                    editProps={{
                      editTable: 'info',
                    }}
                  />
                  <GatewayInfoConfig gwInfo={gwInfo} />
                </Grid>
                <Grid item xs={12}>
                  <CardTitleFilterRow label="Aggregations" />
                  <AddEditGatewayButton
                    title={'Edit'}
                    isLink={true}
                    editProps={{
                      editTable: 'aggregation',
                    }}
                  />
                  <GatewayAggregation gwInfo={gwInfo} />
                </Grid>
              </Grid>
            </Grid>
            <Grid item xs={12} md={6} alignItems="center">
              <Grid container spacing={4}>
                <Grid item xs={12}>
                  <CardTitleFilterRow label="EPC" />
                  <AddEditGatewayButton
                    title={'Edit'}
                    isLink={true}
                    editProps={{
                      editTable: 'epc',
                    }}
                  />
                  <GatewayEPC gwInfo={gwInfo} />
                </Grid>
                <Grid item xs={12}>
                  <CardTitleFilterRow label="Ran" />
                  <AddEditGatewayButton
                    title={'Edit'}
                    isLink={true}
                    editProps={{
                      editTable: 'ran',
                    }}
                  />
                  <GatewayRAN gwInfo={gwInfo} />
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
  const kpiData: KPIRows[] = [
    [
      {
        category: 'Name',
        value: gwInfo.name,
        statusCircle: false,
      },
    ],
    [
      {
        category: 'Gateway ID',
        value: gwInfo.id,
        statusCircle: false,
      },
    ],
    [
      {
        category: 'Hardware UUID',
        value: gwInfo.device.hardware_id,
        statusCircle: false,
      },
    ],
    [
      {
        category: 'Version',
        value: gwInfo.status?.platform_info?.packages?.[0]?.version ?? 'null',
        statusCircle: false,
      },
    ],
    [
      {
        category: 'Description',
        value: gwInfo.description,
        statusCircle: false,
      },
    ],
  ];

  return <KPIGrid data={kpiData} />;
}

function GatewayEPC({gwInfo}: {gwInfo: lte_gateway}) {
  const classes = useStyles();

  const [open, setOpen] = useState({
    ipAllocation: true,
    reservedIp: true,
  });
  const handleCollapse = (config: string) => {
    setOpen({
      ...open,
      [config]: !open[config],
    });
  };

  function ListItems(props) {
    return (
      <>
        <ListItem>
          <ListItemText primary={props.data} />
        </ListItem>
        <Divider />
      </>
    );
  }

  function ListNull() {
    return (
      <>
        <ListItem>
          <ListItemText primary="-" />
        </ListItem>
        <Divider />
      </>
    );
  }

  return (
    <List component={Paper} elevation={0} className={classes.list}>
      <ListItem button onClick={() => handleCollapse('ipAllocation')}>
        <CardHeader
          title="IP Allocation"
          className={classes.kpiBox}
          subheader={gwInfo.cellular.epc.nat_enabled ? 'NAT' : 'Custom'}
          titleTypographyProps={{
            variant: 'body3',
            className: classes.kpiLabel,
            title: 'IP Allocation',
          }}
          subheaderTypographyProps={{
            variant: 'body1',
            className: classes.kpiValue,
            title: gwInfo.cellular.epc.nat_enabled ? 'NAT' : 'Custom',
          }}
        />
        {open['ipAllocation'] ? <ExpandLess /> : <ExpandMore />}
      </ListItem>
      <Divider />
      <Collapse
        key="ipAllocation"
        in={open['ipAllocation']}
        timeout="auto"
        unmountOnExit>
        {gwInfo.cellular.epc.ip_block ? (
          <ListItems data={gwInfo.cellular.epc.ip_block} />
        ) : (
          <ListNull />
        )}
      </Collapse>
      <ListItem>
        <CardHeader
          title="Primary DNS"
          className={classes.kpiBox}
          subheader={gwInfo.cellular.epc.dns_primary ?? '-'}
          titleTypographyProps={{
            variant: 'body3',
            className: classes.kpiLabel,
            title: 'Primary DNS',
          }}
          subheaderTypographyProps={{
            variant: 'body1',
            className: classes.kpiValue,
            title: gwInfo.cellular.epc.dns_primary ?? '-',
          }}
        />
      </ListItem>
      <Divider />
      <ListItem>
        <CardHeader
          title="Secondary DNS"
          className={classes.kpiBox}
          subheader={gwInfo.cellular.epc.dns_secondary ?? '-'}
          titleTypographyProps={{
            variant: 'body3',
            className: classes.kpiLabel,
            title: 'Secondary DNS',
          }}
          subheaderTypographyProps={{
            variant: 'body1',
            className: classes.kpiValue,
            title: gwInfo.cellular.epc.dns_secondary ?? '-',
          }}
        />
      </ListItem>
    </List>
  );
}

function GatewayAggregation({gwInfo}: {gwInfo: lte_gateway}) {
  const logAggregation = !!gwInfo.magmad.dynamic_services?.includes(
    'td-agent-bit',
  );
  const eventAggregation = !!gwInfo.magmad?.dynamic_services?.includes(
    'eventd',
  );
  const aggregations: KPIRows[] = [
    [
      {
        category: 'Aggregation',
        value: logAggregation ? 'Enabled' : 'Disabled',
        statusCircle: false,
      },
      {
        category: 'Aggregation',
        value: eventAggregation ? 'Enabled' : 'Disabled',
        statusCircle: false,
      },
    ],
  ];

  return <KPIGrid data={aggregations} />;
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
  const [open, setOpen] = React.useState(true);
  const classes = useStyles();
  const enbCtx = useContext(EnodebContext);
  const enbInfo =
    gwInfo.connected_enodeb_serials?.reduce(
      (enbs: {[string]: EnodebInfo}, serial: string) => {
        if (enbCtx.state.enbInfo[serial] != null) {
          enbs[serial] = enbCtx.state.enbInfo[serial];
        }
        return enbs;
      },
      {},
    ) || {};

  const ran: KPIRows[] = [
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
  ];

  return (
    <>
      <KPIGrid data={ran} />
      <Divider />
      <List component={Paper} elevation={0} className={classes.list}>
        <ListItem button onClick={() => setOpen(!open)}>
          <CardHeader
            title="Registered eNodeBs"
            className={classes.kpiBox}
            subheader={gwInfo.connected_enodeb_serials?.length || 0}
            titleTypographyProps={{
              variant: 'body3',
              className: classes.kpiLabel,
              title: 'Registered eNodeBs',
            }}
            subheaderTypographyProps={{
              variant: 'body1',
              className: classes.kpiValue,
              title: gwInfo.connected_enodeb_serials?.length || 0,
            }}
          />
          {open ? <ExpandLess /> : <ExpandMore />}
        </ListItem>
        <Divider />
        <Collapse key="reservedIp" in={open} timeout="auto" unmountOnExit>
          <EnodebsTable gwInfo={gwInfo} enbInfo={enbInfo} />
        </Collapse>
      </List>
    </>
  );
}
