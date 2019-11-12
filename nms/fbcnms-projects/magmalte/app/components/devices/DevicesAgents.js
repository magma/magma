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

import Button from '@fbcnms/ui/components/design-system/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import DeviceStatusCircle from '@fbcnms/ui/components/icons/DeviceStatusCircle';
import DevicesEditAgentDialog from './DevicesEditAgentDialog';
import DevicesNewAgentDialog from './DevicesNewAgentDialog';
import EditIcon from '@material-ui/icons/Edit';
import IconButton from '@material-ui/core/IconButton';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import Paper from '@material-ui/core/Paper';
import React, {useCallback, useState} from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Text from '@fbcnms/ui/components/design-system/Text';
import {map} from 'lodash';

import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '../../common/useMagmaAPI';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {Route} from 'react-router-dom';
import {buildDevicesGatewayFromPayload} from './DevicesUtils';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';

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

function DevicesAgents(props: Props) {
  const {match, history, relativePath, relativeUrl} = useRouter();
  const [gateways, setGateways] = useState<?Array<DevicesGateway>>(null);
  const [errorMessage, setErrorMessage] = useState<?string>(null);
  const [editingGateway, setEditingGateway] = useState<?DevicesGateway>(null);
  const classes = useStyles();

  const {error, isLoading} = useMagmaAPI(
    MagmaV1API.getSymphonyByNetworkIdAgents,
    {networkId: nullthrows(match.params.networkId)},
    useCallback(response => {
      if (response != null) {
        setGateways(
          map(response, (agent, _) =>
            buildDevicesGatewayFromPayload(agent),
          ).sort((a, b) => a.id.localeCompare(b.id)),
        );
      }
    }, []),
  );

  if (error || isLoading || !gateways) {
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
    if (!gateway.id) {
      setErrorMessage('Error: cannot delete because id is empty');
    } else {
      props
        .confirm(`Are you sure you want to delete ${gateway.id}?`)
        .then(confirmed => {
          if (!confirmed) {
            return;
          }
          MagmaV1API.deleteSymphonyByNetworkIdAgentsByAgentId({
            networkId: nullthrows(match.params.networkId),
            agentId: gateway.id,
          }).then(() =>
            setGateways(gateways.filter(gw => gw.id != gateway.id)),
          );
          setErrorMessage(null);
        });
    }
  };

  const rows = gateways.map(gateway => (
    <TableRow key={gateway.id}>
      <TableCell>
        {status}
        <DeviceStatusCircle
          isGrey={gateway.status == null}
          isActive={!!gateway.up}
        />
        {gateway.id || 'Error: Missing ID'}
      </TableCell>
      <TableCell>
        {gateway.hardware_id}
        {gateway.devmand_config === undefined && (
          <Text color="error">missing devmand config</Text>
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
        <Text variant="h5">Configure Agents</Text>
        <NestedRouteLink to="/new">
          <Button>Add Agent</Button>
        </NestedRouteLink>
      </div>
      <Paper elevation={2}>
        {errorMessage && <div>{errorMessage.toString()}</div>}
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>ID</TableCell>
              <TableCell>Hardware UUID</TableCell>
              <TableCell />
            </TableRow>
          </TableHead>
          <TableBody>{rows}</TableBody>
        </Table>
      </Paper>
      {editingGateway && (
        <DevicesEditAgentDialog
          key={editingGateway.id}
          gateway={editingGateway}
          onClose={() => setEditingGateway(null)}
          onSave={onSave}
        />
      )}
      <Route
        path={relativePath('/new')}
        render={() => (
          <DevicesNewAgentDialog
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

export default withAlert(DevicesAgents);
