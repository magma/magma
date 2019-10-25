/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {cwf_gateway} from '@fbcnms/magma-api';

import Button from '@material-ui/core/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import IconButton from '@material-ui/core/IconButton';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Typography from '@material-ui/core/Typography';

import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '../../common/useMagmaAPI';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {GatewayStatus} from '../GatewayUtils';
import {makeStyles} from '@material-ui/styles';
import {map} from 'lodash';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';

const useStyles = makeStyles(theme => ({
  header: {
    margin: '10px',
    display: 'flex',
    justifyContent: 'space-between',
  },
  paper: {
    margin: theme.spacing(3),
  },
}));

const FIVE_MINS = 5 * 60 * 1000;

function CWFGateways(props: WithAlert & {}) {
  const [gateways, setGateways] = useState<?(cwf_gateway[])>(null);
  const {match} = useRouter();
  const networkId = nullthrows(match.params.networkId);
  const classes = useStyles();

  const {isLoading} = useMagmaAPI(
    MagmaV1API.getCwfByNetworkIdGateways,
    {networkId},
    response => setGateways(map(response, g => g)),
  );

  if (!gateways || isLoading) {
    return <LoadingFiller />;
  }

  const deleteGateway = (gateway: cwf_gateway) => {
    props
      .confirm(`Are you sure you want to delete ${gateway.name}?`)
      .then(confirmed => {
        if (confirmed) {
          MagmaV1API.deleteCwfByNetworkIdGatewaysByGatewayId({
            networkId,
            gatewayId: gateway.id,
          }).then(() =>
            setGateways(gateways.filter(gw => gw.id != gateway.id)),
          );
        }
      });
  };

  const rows = gateways.map(gateway => (
    <TableRow key={gateway.id}>
      <TableCell>
        <GatewayStatus
          isGrey={!gateway.status?.checkin_time}
          isActive={
            Math.max(0, Date.now() - (gateway.status?.checkin_time || 0)) <
            FIVE_MINS
          }
        />
        {gateway.name}
      </TableCell>
      <TableCell>{gateway.device.hardware_id}</TableCell>
      <TableCell>
        <IconButton color="primary">
          <EditIcon />
        </IconButton>
        <IconButton color="primary" onClick={() => deleteGateway(gateway)}>
          <DeleteIcon />
        </IconButton>
      </TableCell>
    </TableRow>
  ));

  return (
    <div className={classes.paper}>
      <div className={classes.header}>
        <Typography variant="h5">Configure Gateways</Typography>
        <Button variant="contained" color="primary">
          Add Gateway
        </Button>
      </div>
      <Paper elevation={2}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Name</TableCell>
              <TableCell>Hardware UUID</TableCell>
              <TableCell />
            </TableRow>
          </TableHead>
          <TableBody>{rows}</TableBody>
        </Table>
      </Paper>
    </div>
  );
}

export default withAlert(CWFGateways);
