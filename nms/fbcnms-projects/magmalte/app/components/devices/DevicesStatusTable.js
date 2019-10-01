/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {DevicesGatewayPayload} from './DevicesUtils';

import Button from '@material-ui/core/Button';
import DevicesDeviceDialog from './DevicesDeviceDialog';
import DevicesEditManagedDeviceDialog from './DevicesEditManagedDeviceDialog';
import DevicesManagedDeviceRow from './DevicesManagedDeviceRow';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Typography from '@material-ui/core/Typography';
import axios from 'axios';
import {Route} from 'react-router-dom';

import {MagmaAPIUrls} from '@fbcnms/magmalte/app/common/MagmaAPI';
import {
  buildDevicesGatewayFromPayload,
  mergeGatewaysDevices,
} from './DevicesUtils';
import {makeStyles} from '@material-ui/styles';
import {useAxios, useRouter} from '@fbcnms/ui/hooks';
import {useCallback, useEffect, useState} from 'react';

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
  const [rawGateways, setRawGateways] = useState<?(DevicesGatewayPayload[])>(
    null,
  );
  const [devices, setDevices] = useState<?(string[])>(null);

  const {isLoading: gatewaysIsLoading, error} = useAxios<
    null,
    DevicesGatewayPayload[],
  >({
    method: 'get',
    url: MagmaAPIUrls.gateways(match, true),
    onResponse: useCallback(res => {
      if (res.data) {
        setRawGateways(res.data.filter(device => device.record));
      }
    }, []),
  });

  const {isLoading: devicesIsLoading, error: devicesError} = useAxios<
    null,
    string[],
  >({
    method: 'get',
    url: MagmaAPIUrls.devices(match),
    onResponse: useCallback(res => {
      if (res.data) {
        setDevices(res.data);
      }
    }, []),
  });

  useEffect(() => {
    if (!rawGateways) {
      return;
    }

    const interval = setInterval(async () => {
      await Promise.all(
        rawGateways.map(async gateway => {
          try {
            const statusResponse = await axios.get(
              MagmaAPIUrls.gatewayStatus(match, gateway.gateway_id),
            );
            const newGateways = [...rawGateways];
            for (let i = 0; i < rawGateways.length; i++) {
              if (newGateways[i].gateway_id === gateway.gateway_id) {
                newGateways[i] = {
                  ...newGateways[i],
                  status: statusResponse.data,
                };
              }
            }
            setRawGateways(newGateways);
          } catch (err) {
            console.error(
              `Warning: cannot refresh gateway id '${gateway.gateway_id}'. ${err}`,
            );
          }
        }),
      );
    }, REFRESH_INTERVAL);
    return () => clearInterval(interval);
  }, [match, rawGateways]);

  let errorMessage = null;
  let fullDevices = {};

  if (error) {
    errorMessage = error.message;
  } else if (devicesError) {
    errorMessage = devicesError.message;
  } else if (rawGateways != null && devices) {
    const gateways = rawGateways.map(buildDevicesGatewayFromPayload);
    fullDevices = mergeGatewaysDevices(gateways, devices);
  }

  if (error || devicesError || gatewaysIsLoading || devicesIsLoading) {
    return <LoadingFiller />;
  }

  const rows = Object.keys(fullDevices).map(id => (
    <DevicesManagedDeviceRow
      enableDeviceEditing={true}
      deviceID={id}
      onDeleteDevice={_deletedDeviceID => {
        // delete doesn't actually work right now - wait until API V1
      }}
      key={id}
      device={fullDevices[id]}
    />
  ));

  return (
    <>
      <div className={classes.paper}>
        <div className={classes.header}>
          <Typography variant="h5">Devices</Typography>
          <NestedRouteLink to="/add_device/">
            <Button variant="contained" color="primary">
              New Device
            </Button>
          </NestedRouteLink>
        </div>
        {errorMessage && <Typography color="error">{errorMessage}</Typography>}
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
        path={relativePath('/new')}
        render={() => (
          <DevicesDeviceDialog
            onSave={device => {
              const existingGateways = rawGateways || [];
              setRawGateways([...existingGateways, device]);
              history.push(relativeUrl(''));
            }}
            onClose={() => history.push(relativeUrl(''))}
          />
        )}
      />
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
