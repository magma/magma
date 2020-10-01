/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @flow strict-local
 * @format
 */

import type {WifiGateway} from './WifiUtils';

import Button from '@fbcnms/ui/components/design-system/Button';
import IconButton from '@material-ui/core/IconButton';
import LinearProgress from '@material-ui/core/LinearProgress';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import RefreshIcon from '@material-ui/icons/Refresh';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Text from '@fbcnms/ui/components/design-system/Text';
import Tooltip from '@material-ui/core/Tooltip';
import WifiDeviceDialog from './WifiDeviceDialog';
import WifiMeshDialog from './WifiMeshDialog';
import WifiMeshRow from './WifiMeshRow';
import nullthrows from '@fbcnms/util/nullthrows';

import {Route} from 'react-router-dom';
import {buildWifiGatewayFromPayloadV1} from './WifiUtils';
import {makeStyles} from '@material-ui/styles';
import {map} from 'lodash';
import {sortBy} from 'lodash';
import {useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  actionsColumn: {
    width: '160px',
  },
  header: {
    margin: '10px',
    display: 'flex',
    justifyContent: 'space-between',
  },
  infoColumn: {
    width: '400px',
  },
  paper: {
    margin: theme.spacing(3),
  },
}));

export default function WifiMeshesDevicesTable() {
  const classes = useStyles();
  const {match, relativePath, relativeUrl, history} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();

  const [isLoading, setIsLoading] = useState(false);
  const [meshes, setMeshes] = useState<Map<string, WifiGateway[]>>(new Map());
  const [lastRefreshTime, setLastRefreshTime] = useState(
    new Date().toLocaleString(),
  );

  const sortDevices = (d1, d2) =>
    d1.info.toLowerCase() > d2.info.toLowerCase() ? 1 : -1;

  const onCancelDialog = () => history.push(relativeUrl(''));

  useEffect(() => {
    setIsLoading(true);
    const networkId = nullthrows(match.params.networkId);
    Promise.all([
      MagmaV1API.getWifiByNetworkIdGateways({networkId}),
      MagmaV1API.getWifiByNetworkIdMeshes({networkId}),
    ])
      .then(([gatewaysResponse, meshesResponse]) => {
        const meshes = new Map();
        meshesResponse.forEach(meshID => meshes.set(meshID, []));

        const now = new Date().getTime();
        map(gatewaysResponse) // turn id->gateway map into gateway list
          .filter(gateway => gateway.device)
          .forEach(gatewayPayload => {
            const gateway = buildWifiGatewayFromPayloadV1(gatewayPayload, now);
            meshes.set(gateway.meshid, meshes.get(gateway.meshid) || []);
            nullthrows(meshes.get(gateway.meshid)).push(gateway);
          });

        meshes.forEach(gateways => gateways.sort(sortDevices));
        setIsLoading(false);
        setMeshes(meshes);
      })
      .catch((error, _) => {
        setIsLoading(false);
        enqueueSnackbar(error.toString(), {variant: 'error'});
      });
  }, [enqueueSnackbar, match.params.networkId, lastRefreshTime]);

  const onSaveDialog = () => {
    setLastRefreshTime(new Date().toLocaleString());
    onCancelDialog();
  };

  const meshIDs: Array<string> = sortBy(
    [...meshes.keys()], // sortBy can't sort a MapIterator
    [m => m.toLowerCase()],
  );

  const rows = meshIDs.map(meshID => (
    <WifiMeshRow
      enableDeviceEditing={true}
      key={meshID}
      gateways={meshes.get(meshID) || []}
      meshID={meshID}
      onDeleteMesh={() => setLastRefreshTime(new Date().toLocaleString())}
      onDeleteDevice={() => setLastRefreshTime(new Date().toLocaleString())}
    />
  ));

  return (
    <>
      <div className={classes.paper}>
        <div className={classes.header}>
          <Text variant="h5">Devices</Text>
          <div>
            <Tooltip title={'Last refreshed: ' + lastRefreshTime}>
              <span>
                <IconButton
                  color="inherit"
                  onClick={() =>
                    setLastRefreshTime(new Date().toLocaleString())
                  }
                  disabled={isLoading}>
                  <RefreshIcon />
                </IconButton>
              </span>
            </Tooltip>
            <NestedRouteLink to="/add_mesh/">
              <Button>New Mesh</Button>
            </NestedRouteLink>
          </div>
        </div>
        <Paper elevation={2}>
          {isLoading ? <LinearProgress /> : null}
          <Table>
            <TableHead>
              <TableRow>
                <TableCell className={classes.infoColumn}>Info</TableCell>
                <TableCell>ID</TableCell>
                <TableCell className={classes.actionsColumn}>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>{rows}</TableBody>
          </Table>
        </Paper>
        <Route
          path={relativePath('/add_mesh')}
          component={() => (
            <WifiMeshDialog onSave={onSaveDialog} onCancel={onCancelDialog} />
          )}
        />
        <Route
          path={relativePath('/edit_mesh/:meshID')}
          component={() => (
            <WifiMeshDialog onSave={onCancelDialog} onCancel={onCancelDialog} />
          )}
        />
        <Route
          path={relativePath('/add_device/:meshID')}
          component={() => (
            <WifiDeviceDialog
              title="Add"
              onSave={onSaveDialog}
              onCancel={onCancelDialog}
            />
          )}
        />
        <Route
          path={relativePath('/:meshID/edit_device/:deviceID')}
          component={() => (
            <WifiDeviceDialog
              title="Edit"
              onSave={onSaveDialog}
              onCancel={onCancelDialog}
            />
          )}
        />
      </div>
    </>
  );
}
