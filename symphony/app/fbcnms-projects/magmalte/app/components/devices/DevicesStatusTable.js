/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import type {symphony_device} from '@fbcnms/magma-api';

import Button from '@fbcnms/ui/components/design-system/Button';
import DevicesEditManagedDeviceDialog from './DevicesEditManagedDeviceDialog';
import DevicesManagedDeviceRow from './DevicesManagedDeviceRow';
import DevicesMetricsDialog from './DevicesMetricsDialog';
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
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';
import {Route} from 'react-router-dom';

import nullthrows from '@fbcnms/util/nullthrows';
import {augmentDevicesMap} from './DevicesUtils';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useState} from 'react';
import {useInterval, useRouter} from '@fbcnms/ui/hooks';

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
  const [devices, setDevices] = useState<?{[string]: symphony_device}>(null);

  const {isLoading: devicesIsLoading, error: devicesError} = useMagmaAPI(
    MagmaV1API.getSymphonyByNetworkIdDevices,
    {networkId: nullthrows(match.params.networkId)},
    useCallback(response => setDevices(response || {}), []),
  );

  useInterval(async () => {
    try {
      const response = await MagmaV1API.getSymphonyByNetworkIdDevices({
        networkId: nullthrows(match.params.networkId),
      });
      setDevices(response || {});
    } catch (err) {
      console.error(`Warning: cannot refresh'. ${err}`);
    }
  }, REFRESH_INTERVAL);

  let errorMessage = null;
  let fullDevices = {};

  if (devicesError) {
    errorMessage = devicesError.message;
  } else if (devices) {
    fullDevices = augmentDevicesMap(devices);
  }

  if (devicesError || devicesIsLoading) {
    return <LoadingFiller />;
  }

  const rows = Object.keys(fullDevices).map(id => (
    <DevicesManagedDeviceRow
      enableDeviceEditing={true}
      deviceID={id}
      onDeleteDevice={deletedDeviceID => {
        if (devices) {
          if (deletedDeviceID in devices) {
            const newDevices = {...devices};
            delete newDevices[deletedDeviceID];
            setDevices(newDevices);
          }
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
            onSave={(deviceID: string) => {
              const newDevices = {
                ...devices,
                [deviceID]: {
                  id: deviceID,
                  name: deviceID,
                  config: {},
                  managing_agent: '',
                  state: {},
                },
              };
              setDevices(newDevices);
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
      <Route
        path={relativePath('/metrics/:deviceID')}
        render={() => (
          <DevicesMetricsDialog onClose={() => history.push(relativeUrl(''))} />
        )}
      />
    </>
  );
}
