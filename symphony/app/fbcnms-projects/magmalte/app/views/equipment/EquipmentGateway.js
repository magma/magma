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

import CellWifiIcon from '@material-ui/icons/CellWifi';
import DataUsageIcon from '@material-ui/icons/DataUsage';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import MoreVertIcon from '@material-ui/icons/MoreVert';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Text from '@fbcnms/ui/components/design-system/Text';
import isGatewayHealthy from '../../components/GatewayUtils';
import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';
import {makeStyles} from '@material-ui/styles';
import {useHistory} from 'react-router-dom';
import {useRouter} from '@fbcnms/ui/hooks';
import type {gateway_id, lte_gateway} from '@fbcnms/magma-api';

export const DATE_TO_STRING_PARAMS = [
  'en-US',
  {
    year: 'numeric',
    month: 'numeric',
    day: 'numeric',
    hour: 'numeric',
    minute: 'numeric',
    second: 'numeric',
  },
];

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

export default function Gateway() {
  const classes = useStyles();
  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={3} alignItems="stretch">
        <Grid item xs={12}>
          <Text>
            <DataUsageIcon /> Gateway Check-Ins
          </Text>
          <Paper className={classes.paper}>Gateway Check in chart</Paper>
        </Grid>

        <Grid item xs={12}>
          <Text>
            <CellWifiIcon /> Gateways
          </Text>
          <GatewayTable />
        </Grid>
      </Grid>
    </div>
  );
}

function GatewayTable() {
  const {match} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);
  const {response, isLoading} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdGateways,
    {
      networkId: networkId,
    },
  );

  if (isLoading || !response) {
    return <LoadingFiller />;
  }
  const lte_gateways: {[string]: lte_gateway} = response;
  return (
    <TableContainer component={Paper}>
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>Name</TableCell>
            <TableCell>ID</TableCell>
            <TableCell>eNodeBs</TableCell>
            <TableCell>Subscribers</TableCell>
            <TableCell>Health</TableCell>
            <TableCell>Check In Time</TableCell>
            {/* placeholder column for "actionMenu" */}
            <TableCell />
          </TableRow>
        </TableHead>
        <TableBody>
          {Object.keys(lte_gateways)
            .map((gwId: string) => lte_gateways[gwId])
            .filter((g: lte_gateway) => g.cellular && g.id)
            .map((gateway: lte_gateway, rowIdx) => {
              return (
                <TableRow key={rowIdx} data-testid={'gatewayInfo-' + rowIdx}>
                  <TableCell key={'name-' + rowIdx} component="th" scope="row">
                    {gateway.name}
                  </TableCell>

                  <TableCell key={'id-' + rowIdx}>{gateway.id}</TableCell>

                  <TableCell key={'enbs-' + rowIdx}>
                    {gateway.connected_enodeb_serials
                      ? gateway.connected_enodeb_serials.length
                      : 0}
                  </TableCell>

                  <TableCell key={'subs-' + rowIdx}>0</TableCell>

                  <TableCell key={'health-' + rowIdx}>
                    {isGatewayHealthy(gateway) ? 'Good' : 'Bad'}
                  </TableCell>

                  <TableCell key={'created-' + rowIdx}>
                    {gateway.status &&
                    (gateway.status.checkin_time !== undefined ||
                      gateway.status.checkin_time === null)
                      ? new Date(
                          gateway.status.checkin_time,
                        ).toLocaleDateString(...DATE_TO_STRING_PARAMS)
                      : 0}
                  </TableCell>

                  <TableCell align="right" key={'action-' + rowIdx}>
                    <GatewayActionMenu gwId={gateway.id} />
                  </TableCell>
                </TableRow>
              );
            })}
        </TableBody>
      </Table>
    </TableContainer>
  );
}

function GatewayActionMenu({gwId}: {gwId: gateway_id}) {
  const history = useHistory();

  const [anchorEl, setAnchorEl] = React.useState(null);
  const {relativeUrl} = useRouter();
  const handleClick = event => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleView = () => {
    setAnchorEl(null);
    history.push(relativeUrl('/' + gwId));
  };

  return (
    <div>
      <IconButton onClick={handleClick}>
        <MoreVertIcon />
      </IconButton>
      <Menu
        id="simple-menu"
        anchorEl={anchorEl}
        keepMounted
        open={Boolean(anchorEl)}
        onClose={handleClose}>
        <MenuItem onClick={handleView}>View</MenuItem>
        <MenuItem onClick={handleClose}>Edit</MenuItem>
        <MenuItem onClick={handleClose}>Reboot</MenuItem>
        <MenuItem onClick={handleClose}>Deactivate</MenuItem>
        <MenuItem onClick={handleClose}>Remove</MenuItem>
      </Menu>
    </div>
  );
}
