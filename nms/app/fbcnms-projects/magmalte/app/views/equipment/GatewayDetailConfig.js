/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import type {EnodebInfo} from '../../components/lte/EnodebUtils';
import type {KPIRows} from '../../components/KPIGrid';
import type {lte_gateway} from '@fbcnms/magma-api';

import ActionTable from '../../components/ActionTable';
import Collapse from '@material-ui/core/Collapse';
import Divider from '@material-ui/core/Divider';
import ExpandLess from '@material-ui/icons/ExpandLess';
import ExpandMore from '@material-ui/icons/ExpandMore';
import Grid from '@material-ui/core/Grid';
import KPIGrid from '../../components/KPIGrid';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import SettingsIcon from '@material-ui/icons/Settings';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextField from '@material-ui/core/TextField';

import {colors} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useState} from 'react';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(3),
    flexGrow: 1,
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
}));

export default function GatewayConfig({
  gwInfo,
  enbInfo,
}: {
  gwInfo: lte_gateway,
  enbInfo: {[string]: EnodebInfo},
}) {
  const classes = useStyles();
  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={3}>
        <Grid container spacing={3} item xs={12}>
          <Text>
            <SettingsIcon /> Config
          </Text>
        </Grid>
        <Grid container spacing={3} item xs={6}>
          <Grid item xs={12}>
            <Grid container>
              <Grid item xs={6}>
                <Text>Gateway</Text>
              </Grid>
              <Grid container item xs={6} justify="flex-end">
                <Text>Edit</Text>
              </Grid>
            </Grid>
            <Grid item xs={12}>
              <GatewayInfoConfig gwInfo={gwInfo} />
            </Grid>
          </Grid>
          <Grid item xs={12}>
            <Grid container>
              <Grid item xs={6}>
                <Text>Aggregations</Text>
              </Grid>
              <Grid container item xs={6} justify="flex-end">
                <Text>Edit</Text>
              </Grid>
            </Grid>
            <Grid item xs={12}>
              <GatewayAggregation gwInfo={gwInfo} />
            </Grid>
          </Grid>
        </Grid>
        <Grid container spacing={3} item xs={6}>
          <Grid item xs={12}>
            <Grid container>
              <Grid item xs={6}>
                <Text>EPC</Text>
              </Grid>
              <Grid container item xs={6} justify="flex-end">
                <Text>Edit</Text>
              </Grid>
            </Grid>
            <Grid item xs={12}>
              <GatewayEPC gwInfo={gwInfo} />
            </Grid>
          </Grid>
          <Grid item xs={12}>
            <Grid container>
              <Grid item xs={6}>
                <Text>Ran</Text>
              </Grid>
              <Grid container item xs={6} justify="flex-end">
                <Text>Edit</Text>
              </Grid>
            </Grid>
            <Grid item xs={12}>
              <GatewayRAN gwInfo={gwInfo} enbInfo={enbInfo} />
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    </div>
  );
}

function GatewayInfoConfig({gwInfo}: {gwInfo: lte_gateway}) {
  const classes = useStyles();
  const typographyProps = {
    primaryTypographyProps: {
      variant: 'caption',
      className: classes.itemTitle,
    },
    secondaryTypographyProps: {
      variant: 'h6',
      className: classes.itemValue,
    },
  };
  return (
    <List component={Paper}>
      <ListItem>
        <ListItemText
          primary="Name"
          secondary={gwInfo.name}
          {...typographyProps}
        />
      </ListItem>
      <Divider />
      <ListItem>
        <ListItemText
          secondary={gwInfo.id}
          primary="Gateway ID"
          {...typographyProps}
        />
      </ListItem>
      <Divider />
      <ListItem>
        <ListItemText
          secondary={gwInfo.device.hardware_id}
          primary="Hardware UUID"
          {...typographyProps}
        />
      </ListItem>
      <Divider />
      <ListItem>
        <ListItemText
          secondary={
            gwInfo.status?.platform_info?.packages?.[0]?.version ?? 'null'
          }
          primary="Version"
          {...typographyProps}
        />
      </ListItem>
      <Divider />
      <ListItem>
        <ListItemText
          secondary={gwInfo.description}
          primary="Description"
          {...typographyProps}
        />
      </ListItem>
    </List>
  );
}

function GatewayEPC({gwInfo}: {gwInfo: lte_gateway}) {
  const classes = useStyles();
  const typographyProps = {
    primaryTypographyProps: {
      variant: 'caption',
      className: classes.itemTitle,
    },
    secondaryTypographyProps: {
      variant: 'h6',
      className: classes.itemValue,
    },
  };

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
  return (
    <List component={Paper}>
      <ListItem button onClick={() => handleCollapse('ipAllocation')}>
        <ListItemText
          primary={'IP Allocation'}
          secondary={gwInfo.cellular.epc.nat_enabled ? 'NAT' : 'Custom'}
          {...typographyProps}
        />
        {open['ipAllocation'] ? <ExpandLess /> : <ExpandMore />}
      </ListItem>
      <Divider />
      <Collapse
        key="ipAllocation"
        in={open['ipAllocation']}
        timeout="auto"
        unmountOnExit>
        <ListItem>
          <ListItemText
            secondary={gwInfo.cellular.epc.ip_block}
            {...typographyProps}
          />
        </ListItem>
        <Divider />
      </Collapse>
      <ListItem button onClick={() => handleCollapse('reservedIp')}>
        <ListItemText
          primary={'Primary DNS'}
          secondary={gwInfo.cellular.epc.dns_primary ?? '-'}
          {...typographyProps}
        />
      </ListItem>
      <Divider />
      <ListItem>
        <ListItemText
          primary={'Secondary DNS'}
          secondary={gwInfo.cellular.epc.dns_secondary ?? '-'}
          {...typographyProps}
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
      columns={[
        {title: 'Name', field: 'name'},
        {title: 'Serial Number', field: 'id'},
      ]}
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

function GatewayRAN({
  gwInfo,
  enbInfo,
}: {
  gwInfo: lte_gateway,
  enbInfo: {[string]: EnodebInfo},
}) {
  const [open, setOpen] = React.useState(true);
  const classes = useStyles();
  const typographyProps = {
    primaryTypographyProps: {
      variant: 'caption',
      className: classes.itemTitle,
      readOnly: true,
    },
    secondaryTypographyProps: {
      variant: 'h6',
      className: classes.itemValue,
      readOnly: true,
    },
  };

  return (
    <List component={Paper}>
      <ListItem>
        <Grid container>
          <Grid item xs={6}>
            <TextField
              fullWidth={true}
              value={gwInfo.cellular.ran.pci}
              label="PCI"
              InputProps={{disableUnderline: true, readOnly: true}}
            />
          </Grid>
          <Grid item xs={6}>
            <TextField
              fullWidth={true}
              value={
                gwInfo.cellular.ran.transmit_enabled ? 'Enabled' : 'Disabled'
              }
              label="eNodeB Transmit"
              InputProps={{disableUnderline: true, readOnly: true}}
            />
          </Grid>
        </Grid>
      </ListItem>
      <Divider />
      <ListItem button onClick={() => setOpen(!open)}>
        <ListItemText
          primary={'Registered eNodeBs'}
          secondary={gwInfo.connected_enodeb_serials?.length || 0}
          {...typographyProps}
        />
        {open ? <ExpandLess /> : <ExpandMore />}
      </ListItem>
      <Collapse key="reservedIp" in={open} timeout="auto" unmountOnExit>
        <Divider />
        <Grid item xs={12}>
          <EnodebsTable gwInfo={gwInfo} enbInfo={enbInfo} />
        </Grid>
      </Collapse>
    </List>
  );
}
