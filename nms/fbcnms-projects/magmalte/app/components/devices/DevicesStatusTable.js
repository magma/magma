/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {symphony_agent} from '@fbcnms/magma-api';

import Button from '@fbcnms/ui/components/design-system/Button';
import DevicesEditManagedDeviceDialog from './DevicesEditManagedDeviceDialog';
import DevicesManagedDeviceRow from './DevicesManagedDeviceRow';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
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
import useMagmaAPI from '../../common/useMagmaAPI';
import {Route} from 'react-router-dom';
import {map} from 'lodash';

import nullthrows from '@fbcnms/util/nullthrows';
import {buildDevicesAgentFromPayload, mergeAgentsDevices} from './DevicesUtils';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useEffect, useState} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  actionsColumn: {
    textAlign: 'right',
    width: '160px',
  },
  infoColumn: {
    width: '600px',
  },
  header: {
    margin: '10px',
    display: 'flex',
    justifyContent: 'space-between',
  },
  paper: {
    margin: theme.spacing(3),
  },
  iconButton: {
    color: theme.palette.primary.dark,
    padding: '5px',
  },
  subrowCell: {
    paddingBottom: '15px',
    paddingLeft: '75px',
    paddingRight: '15px',
    paddingTop: '15px',
  },
}));

const REFRESH_INTERVAL = 10000;

export default function DevicesStatusTable() {
  const classes = useStyles();
  const {match, relativePath, relativeUrl, history} = useRouter();
  const [rawAgents, setRawAgents] = useState<?(symphony_agent[])>(null);
  const [devices, setDevices] = useState<?(string[])>(null);

  const {isLoading: agentsIsLoading, error} = useMagmaAPI(
    MagmaV1API.getSymphonyByNetworkIdAgents,
    {networkId: nullthrows(match.params.networkId)},
    useCallback(response => setRawAgents(map(response, agent => agent)), []),
  );

  const {isLoading: devicesIsLoading, error: devicesError} = useMagmaAPI(
    MagmaV1API.getSymphonyByNetworkIdDevices,
    {networkId: nullthrows(match.params.networkId)},
    useCallback(response => {
      if (response != null) {
        setDevices(Object.keys(response));
      }
    }, []),
  );

  useEffect(() => {
    if (!rawAgents) {
      return;
    }
    const interval = setInterval(async () => {
      try {
        const response = await MagmaV1API.getSymphonyByNetworkIdAgents({
          networkId: nullthrows(match.params.networkId),
        });
        setRawAgents(map(response, agent => agent));
      } catch (err) {
        console.error(`Warning: cannot refresh'. ${err}`);
      }
    }, REFRESH_INTERVAL);
    return () => clearInterval(interval);
  }, [match, rawAgents]);

  let errorMessage = null;
  let fullDevices = {};

  if (error) {
    errorMessage = error.message;
  } else if (devicesError) {
    errorMessage = devicesError.message;
  } else if (rawAgents != null && devices) {
    const agents = rawAgents.map(buildDevicesAgentFromPayload);
    fullDevices = mergeAgentsDevices(agents, devices);
  }

  if (error || devicesError || agentsIsLoading || devicesIsLoading) {
    return <LoadingFiller />;
  }

  const rows = Object.keys(fullDevices).map(id => (
    <DevicesManagedDeviceRow
      enableDeviceEditing={true}
      deviceID={id}
      onDeleteDevice={deletedDeviceID => {
        if (devices) {
          setDevices(devices.filter(deviceId => deviceId != deletedDeviceID));
        }
      }}
      key={id}
      device={fullDevices[id]}
    />
  ));

  return (
    <>
      <div className={classes.paper}>
        <div className={classes.header}>
          <Text variant="h5">Devices</Text>
          <NestedRouteLink to="/add_device/">
            <Button>New Device</Button>
          </NestedRouteLink>
        </div>
        {errorMessage && <Text color="error">{errorMessage}</Text>}
        <Paper>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Name</TableCell>
                <TableCell className={classes.infoColumn}>State</TableCell>
                <TableCell>Managing Agent</TableCell>
                <TableCell className={classes.actionsColumn}>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>{rows}</TableBody>
          </Table>
        </Paper>
      </div>
      <Route
        path={relativePath('/add_device')}
        render={() => (
          <DevicesEditManagedDeviceDialog
            title="Add New Device"
            onSave={deviceID => {
              setDevices([...(devices || []), deviceID]);
              history.push(relativeUrl(''));
            }}
            onCancel={() => history.push(relativeUrl(''))}
          />
        )}
      />
      <Route
        path={relativePath('/edit_device/:deviceID')}
        render={() => (
          <DevicesEditManagedDeviceDialog
            title="Edit Device Management Configs"
            onSave={() => history.push(relativeUrl(''))} // TODO update devices
            onCancel={() => history.push(relativeUrl(''))}
          />
        )}
      />
    </>
  );
}
