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
import type {federation_gateway} from '@fbcnms/magma-api';

import Button from '@material-ui/core/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import FEGGatewayDialog from './FEGGatewayDialog';
import IconButton from '@material-ui/core/IconButton';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Text from '@fbcnms/ui/components/design-system/Text';

import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '../../common/useMagmaAPI';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {GatewayStatus} from '../GatewayUtils';
import {Route} from 'react-router-dom';
import {findIndex} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useState} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  header: {
    margin: '10px',
    display: 'flex',
    justifyContent: 'space-between',
  },
  greCell: {
    paddingBottom: '15px',
    paddingLeft: '75px',
    paddingRight: '15px',
    paddingTop: '15px',
  },
  paper: {
    margin: theme.spacing(3),
  },
  expandIconButton: {
    color: theme.palette.primary.dark,
    padding: '5px',
  },
  tableCell: {
    padding: '15px',
  },
  tableRow: {
    height: 'auto',
    whiteSpace: 'nowrap',
    verticalAlign: 'top',
  },
  gatewayName: {
    paddingRight: '10px',
  },
}));

const FIVE_MINS = 5 * 60 * 1000;

function CWFGateways(props: WithAlert & {}) {
  const [gateways, setGateways] = useState<?(federation_gateway[])>(null);
  const {match, history, relativePath, relativeUrl} = useRouter();
  const networkId = nullthrows(match.params.networkId);
  const classes = useStyles();

  const {isLoading} = useMagmaAPI(
    MagmaV1API.getFegByNetworkIdGateways,
    {networkId},
    useCallback(
      response =>
        setGateways(
          Object.keys(response)
            .map(k => response[k])
            .filter(g => g.id),
        ),
      [],
    ),
  );

  if (!gateways || isLoading) {
    return <LoadingFiller />;
  }

  const deleteGateway = (gateway: federation_gateway) => {
    props
      .confirm(`Are you sure you want to delete ${gateway.name}?`)
      .then(confirmed => {
        if (confirmed) {
          MagmaV1API.deleteFegByNetworkIdGatewaysByGatewayId({
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
        <span className={classes.gatewayName}>{gateway.name}</span>
        <GatewayStatus
          isGrey={!gateway.status?.checkin_time}
          isActive={
            Math.max(0, Date.now() - (gateway.status?.checkin_time || 0)) <
            FIVE_MINS
          }
        />
      </TableCell>
      <TableCell>{gateway.device.hardware_id}</TableCell>
      <TableCell>
        <IconButton
          color="primary"
          onClick={() => history.push(relativeUrl(`/edit/${gateway.id}`))}>
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
        <Text variant="h5">Configure Gateways</Text>
        <NestedRouteLink to="/new">
          <Button variant="contained" color="primary">
            Add Gateway
          </Button>
        </NestedRouteLink>
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
      <Route
        path={relativePath('/new')}
        render={() => (
          <FEGGatewayDialog
            onClose={() => history.push(relativeUrl(''))}
            onSave={gateway => {
              setGateways([...gateways, gateway]);
              history.push(relativeUrl(''));
            }}
          />
        )}
      />
      <Route
        path={relativePath('/edit/:gatewayID')}
        render={({match}) => (
          <FEGGatewayDialog
            editingGateway={nullthrows(
              gateways.find(gw => gw.id === match.params.gatewayID),
            )}
            onClose={() => history.push(relativeUrl(''))}
            onSave={gateway => {
              const newGateways = [...gateways];
              const i = findIndex(newGateways, g => g.id === gateway.id);
              newGateways[i] = gateway;
              setGateways(newGateways);
              history.push(relativeUrl(''));
            }}
          />
        )}
      />
    </div>
  );
}

export default withAlert(CWFGateways);
