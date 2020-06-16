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

import ActionTable from '../../components/ActionTable';
import EnodebThroughputChart from './EnodebThroughputChart';
import Grid from '@material-ui/core/Grid';
import React from 'react';
import SettingsInputAntennaIcon from '@material-ui/icons/SettingsInputAntenna';

import {colors} from '../../theme/default';
import {isEnodebHealthy} from '../../components/lte/EnodebUtils';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';

const CHART_TITLE = 'Total Throughput';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(3),
    flexGrow: 1,
  },
  topBar: {
    backgroundColor: colors.primary.mirage,
    padding: '20px 40px 20px 40px',
  },
  tabBar: {
    backgroundColor: colors.primary.brightGray,
    padding: '0 0 0 20px',
  },
  tabs: {
    color: colors.primary.white,
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
  },
  formControl: {
    margin: theme.spacing(1),
    minWidth: 120,
  },
}));

export default function Enodeb({enbInfo}: {enbInfo: {[string]: EnodebInfo}}) {
  const classes = useStyles();

  return (
    <div className={classes.dashboardRoot}>
      <Grid container justify="space-between" spacing={3}>
        <Grid item xs={12}>
          <EnodebThroughputChart
            title={CHART_TITLE}
            queries={[
              `sum(pdcp_user_plane_bytes_dl{service="enodebd"} + pdcp_user_plane_bytes_ul{service="enodebd"})/1000`,
            ]}
            legendLabels={['mbps']}
          />
        </Grid>
        <Grid item xs={12}>
          <EnodeTable enbInfo={enbInfo} />
        </Grid>
      </Grid>
    </div>
  );
}

type EnodebRowType = {
  name: string,
  id: string,
  sessionName: string,
  health: string,
  reportedTime: Date,
};

function EnodeTable({enbInfo}: {enbInfo: {[string]: EnodebInfo}}) {
  const {history, relativeUrl} = useRouter();
  const [currRow, setCurrRow] = useState<EnodebRowType>({});
  const enbRows: Array<EnodebRowType> = Object.keys(enbInfo).map(
    (serialNum: string) => {
      const enbInf = enbInfo[serialNum];
      return {
        name: enbInf.enb.name,
        id: serialNum,
        sessionName: enbInf.enb_state.fsm_state,
        health: isEnodebHealthy(enbInf) ? 'Good' : 'Bad',
        reportedTime: new Date(enbInf.enb_state.time_reported ?? 0),
      };
    },
  );

  return (
    <ActionTable
      titleIcon={SettingsInputAntennaIcon}
      title="EnodeBs"
      data={enbRows}
      columns={[
        {title: 'Name', field: 'name'},
        {title: 'Serial Number', field: 'id'},
        {title: 'Session State Name', field: 'sessionName'},
        {title: 'Health', field: 'health'},
        {title: 'Reported Time', field: 'reportedTime', type: 'datetime'},
      ]}
      handleCurrRow={(row: EnodebRowType) => setCurrRow(row)}
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
