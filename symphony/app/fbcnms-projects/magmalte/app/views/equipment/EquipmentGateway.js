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
import type {gateway_id, lte_gateway} from '@fbcnms/magma-api';

import ActionTable from '../../components/ActionTable';
import CellWifiIcon from '@material-ui/icons/CellWifi';
import EquipmentGatewayKPIs from './EquipmentGatewayKPIs';
import GatewayCheckinChart from './GatewayCheckinChart';
import Grid from '@material-ui/core/Grid';
import React, {useState} from 'react';
import isGatewayHealthy from '../../components/GatewayUtils';

import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(3),
    flexGrow: 1,
  },
  topBar: {
    backgroundColor: theme.palette.magmalte.background,
    padding: '20px 40px 20px 40px',
  },
  tabBar: {
    backgroundColor: theme.palette.magmalte.appbar,
    padding: '0 0 0 20px',
  },
  tabs: {
    color: 'white',
  },
  tab: {
    fontSize: '18px',
    textTransform: 'none',
  },
  tabLabel: {
    padding: '20px 0 20px 0',
  },
  tabIconLabel: {
    verticalAlign: 'middle',
    margin: '0 5px 3px 0',
  },
  // TODO: remove this when we actually fill out the grid sections
  contentPlaceholder: {
    padding: '50px 0',
  },
  paper: {
    height: 100,
    padding: theme.spacing(10),
    textAlign: 'center',
    color: theme.palette.text.secondary,
  },
  formControl: {
    margin: theme.spacing(1),
    minWidth: 120,
  },
}));

export default function Gateway({
  lte_gateways,
}: {
  lte_gateways: {[string]: lte_gateway},
}) {
  const classes = useStyles();

  return (
    <div className={classes.dashboardRoot}>
      <Grid container justify="space-between" spacing={3}>
        <Grid item xs={12}>
          <GatewayCheckinChart />
        </Grid>
        <Grid item xs={12}>
          <EquipmentGatewayKPIs lte_gateways={lte_gateways} />
        </Grid>
        <Grid item xs={12}>
          <GatewayTable lte_gateways={lte_gateways} />
        </Grid>
      </Grid>
    </div>
  );
}

type EquipmentGatewayRowType = {
  name: string,
  id: gateway_id,
  num_enodeb: number,
  num_subscribers: number,
  health: string,
  checkInTime: Date,
};

function GatewayTable({lte_gateways}: {lte_gateways: {[string]: lte_gateway}}) {
  const {history, relativeUrl} = useRouter();
  const [currRow, setCurrRow] = useState<EquipmentGatewayRowType>({});
  const lte_gateway_rows: Array<EquipmentGatewayRowType> = Object.keys(
    lte_gateways,
  )
    .map((gwId: string) => lte_gateways[gwId])
    .filter((g: lte_gateway) => g.cellular && g.id)
    .map((gateway: lte_gateway) => {
      let numEnodeBs = 0;
      if (gateway.connected_enodeb_serials) {
        numEnodeBs = gateway.connected_enodeb_serials.length;
      }

      let checkInTime = new Date(0);
      if (
        gateway.status &&
        (gateway.status.checkin_time !== undefined ||
          gateway.status.checkin_time === null)
      ) {
        checkInTime = new Date(gateway.status.checkin_time);
      }

      return {
        name: gateway.name,
        id: gateway.id,
        num_enodeb: numEnodeBs,
        num_subscribers: 0,
        health: isGatewayHealthy(gateway) ? 'Good' : 'Bad',
        checkInTime: checkInTime,
      };
    });

  return (
    <ActionTable
      titleIcon={CellWifiIcon}
      title={'Gateways'}
      data={lte_gateway_rows}
      columns={[
        {title: 'Name', field: 'name'},
        {title: 'ID', field: 'id'},
        {title: 'enodeBs', field: 'num_enodeb', type: 'numeric'},
        {title: 'Subscribers', field: 'num_subscribers', type: 'numeric'},
        {title: 'Health', field: 'health'},
        {title: 'Check In Time', field: 'checkInTime', type: 'datetime'},
      ]}
      handleCurrRow={(row: EquipmentGatewayRowType) => setCurrRow(row)}
      menuItems={[
        {
          name: 'View',
          handleFunc: () => {
            history.push(relativeUrl('/' + currRow.id));
          },
        },
        {name: 'Edit'},
        {name: 'Remove'},
        {name: 'Deactivate'},
        {name: 'Reboot'},
      ]}
      options={{
        actionsColumnIndex: -1,
        pageSizeOptions: [5, 10],
      }}
    />
  );
}
