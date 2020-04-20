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

import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {federation_gateway} from '@fbcnms/magma-api';

import Button from '@material-ui/core/Button';
import CircularProgress from '@material-ui/core/CircularProgress';
import DeleteIcon from '@material-ui/icons/Delete';
import DeviceStatusCircle from '@fbcnms/ui/components/icons/DeviceStatusCircle';
import EditIcon from '@material-ui/icons/Edit';
import FEGGatewayDialog from './FEGGatewayDialog';
import IconButton from '@material-ui/core/IconButton';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import StarIcon from '@material-ui/icons/Star';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Tooltip from '@material-ui/core/Tooltip';

import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import Text from '../../theme/design-system/Text';
import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {Route} from 'react-router-dom';
import {colors} from '../../theme/default';
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
    color: colors.primary.brightGray,
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
  star: {
    color: '#ffd700',
    width: '18px',
    verticalAlign: 'bottom',
  },
}));

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
  const {
    response: clusterStatus,
    isLoading: clusterStatusLoading,
  } = useMagmaAPI(MagmaV1API.getFegByNetworkIdClusterStatus, {networkId});

  if (!gateways || isLoading || clusterStatusLoading) {
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
    <GatewayRow
      key={gateway.id}
      gateway={gateway}
      onDelete={deleteGateway}
      isPrimary={clusterStatus?.active_gateway === gateway.id}
    />
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

function GatewayRow(props: {
  gateway: federation_gateway,
  onDelete: federation_gateway => void,
  isPrimary: boolean,
}) {
  const classes = useStyles();
  const {gateway, onDelete, isPrimary} = props;
  const {match, history, relativeUrl} = useRouter();
  const {isLoading, response} = useMagmaAPI(
    MagmaV1API.getFegByNetworkIdGatewaysByGatewayIdHealthStatus,
    {
      networkId: nullthrows(match.params.networkId),
      gatewayId: gateway.id,
    },
  );

  return (
    <TableRow key={gateway.id}>
      <TableCell>
        <span className={classes.gatewayName}>{gateway.name}</span>
        {isLoading ? (
          <CircularProgress size={20} />
        ) : (
          <DeviceStatusCircle
            isGrey={!response?.status}
            isActive={response?.status === 'HEALTHY'}
          />
        )}
        {isPrimary && (
          <Tooltip title="Primary FEG" placement="right">
            <StarIcon className={classes.star} />
          </Tooltip>
        )}
      </TableCell>
      <TableCell>{gateway.device.hardware_id}</TableCell>
      <TableCell>
        <IconButton
          color="primary"
          onClick={() => history.push(relativeUrl(`/edit/${gateway.id}`))}>
          <EditIcon />
        </IconButton>
        <IconButton color="primary" onClick={() => onDelete(gateway)}>
          <DeleteIcon />
        </IconButton>
      </TableCell>
    </TableRow>
  );
}

export default withAlert(CWFGateways);
