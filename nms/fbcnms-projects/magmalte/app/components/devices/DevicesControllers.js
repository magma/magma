/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {DevicesGateway} from './DevicesUtils';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';

import Button from '@material-ui/core/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import DeviceStatusCircle from '@fbcnms/ui/components/icons/DeviceStatusCircle';
import DevicesEditControllerDialog from './DevicesEditControllerDialog';
import EditIcon from '@material-ui/icons/Edit';
import IconButton from '@material-ui/core/IconButton';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import NewControllerDialog from './NewControllerDialog';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Typography from '@material-ui/core/Typography';

import nullthrows from '@fbcnms/util/nullthrows';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {MagmaAPIUrls} from '@fbcnms/magmalte/app/common/MagmaAPI';
import {Route} from 'react-router-dom';
import {buildDevicesGatewayFromPayload} from './DevicesUtils';
import {makeStyles} from '@material-ui/styles';
import {useAxios, useRouter} from '@fbcnms/ui/hooks';
import {useCallback, useState} from 'react';

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

type Props = WithAlert & {};

function DevicesControllers(props: Props) {
  const {match, history, relativePath, relativeUrl} = useRouter();
  const [gateways, setGateways] = useState<?Array<DevicesGateway>>(null);
  const [editingGateway, setEditingGateway] = useState<?DevicesGateway>(null);
  const classes = useStyles();

  const {error} = useAxios({
    method: 'get',
    url: MagmaAPIUrls.gateways(match, true),
    onResponse: useCallback(response => {
      const gateways = response.data
        .filter(g => g.record && g.config)
        .map(g => buildDevicesGatewayFromPayload(g))
        .sort((a, b) => a.id.localeCompare(b.id));
      setGateways(gateways);
    }, []),
  });

  if (error || !gateways) {
    return <LoadingFiller />;
  }

  const onSave = gatewayPayload => {
    const gateway = buildDevicesGatewayFromPayload(gatewayPayload);
    const newGateways = gateways.slice(0);
    if (editingGateway) {
      newGateways[newGateways.indexOf(editingGateway)] = gateway;
    } else {
      newGateways.push(gateway);
    }
    setGateways(newGateways);
    setEditingGateway(null);
  };

  const deleteGateway = gateway => {
    props
      .confirm(`Are you sure you want to delete ${gateway.id}?`)
      .then(confirmed => {
        if (!confirmed) {
          return;
        }
        MagmaV1API.deleteNetworksByNetworkIdGatewaysByGatewayId({
          networkId: nullthrows(match.params.networkId),
          gatewayId: gateway.id,
        }).then(() => setGateways(gateways.filter(gw => gw.id != gateway.id)));
      });
  };

  const rows = gateways.map(gateway => (
    <TableRow key={gateway.id}>
      <TableCell>
        {status}
        <DeviceStatusCircle
          isGrey={gateway.status == null}
          isActive={!!gateway.up}
        />
        {gateway.id}
      </TableCell>
      <TableCell>
        {gateway.hardware_id}
        {gateway.devmand_config === undefined && (
          <Typography color="error">missing devmand config</Typography>
        )}
      </TableCell>
      <TableCell>
        <IconButton color="primary" onClick={() => setEditingGateway(gateway)}>
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
        <Typography variant="h5">Configure Controllers</Typography>
        <NestedRouteLink to="/new">
          <Button variant="contained" color="primary">
            Add Controller
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
      {editingGateway && (
        <DevicesEditControllerDialog
          key={editingGateway.id}
          gateway={editingGateway}
          onClose={() => setEditingGateway(null)}
          onSave={onSave}
        />
      )}
      <Route
        path={relativePath('/new')}
        render={() => (
          <NewControllerDialog
            onClose={() => history.push(relativeUrl(''))}
            onSave={rawGateway => {
              setGateways([
                ...gateways,
                buildDevicesGatewayFromPayload(rawGateway),
              ]);
              history.push(relativeUrl(''));
            }}
          />
        )}
      />
    </div>
  );
}

export default withAlert(DevicesControllers);
