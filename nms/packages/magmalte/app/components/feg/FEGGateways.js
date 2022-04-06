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

import type {WithAlert} from '../../../fbc_js_core/ui/components/Alert/withAlert';
import type {federation_gateway} from '../../../generated/MagmaAPIBindings';

import Button from '@material-ui/core/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import DeviceStatusCircle from '../../../fbc_js_core/ui/components/icons/DeviceStatusCircle';
import EditIcon from '@material-ui/icons/Edit';
import FEGGatewayContext from '../context/FEGGatewayContext';
import FEGGatewayDialog from './FEGGatewayDialog';
import IconButton from '@material-ui/core/IconButton';
import NestedRouteLink from '../../../fbc_js_core/ui/components/NestedRouteLink';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import StarIcon from '@material-ui/icons/Star';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Text from '../../theme/design-system/Text';
import Tooltip from '@material-ui/core/Tooltip';
import nullthrows from '../../../fbc_js_core/util/nullthrows';
import withAlert from '../../../fbc_js_core/ui/components/Alert/withAlert';

import {HEALTHY_STATUS} from '../GatewayUtils';
import {Route} from 'react-router-dom';
import {colors} from '../../theme/default';
import {findIndex} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useContext, useState} from 'react';
import {useRouter} from '../../../fbc_js_core/ui/hooks';

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
  const ctx = useContext(FEGGatewayContext);
  const [gateways, setGateways] = useState<federation_gateway[]>(
    Object.keys(ctx.state).map(gatewayId => ctx.state[gatewayId]),
  );
  const {history, relativePath, relativeUrl} = useRouter();
  const classes = useStyles();
  const deleteGateway = (gateway: federation_gateway) => {
    props
      .confirm(`Are you sure you want to delete ${gateway.name}?`)
      .then(confirmed => {
        if (confirmed) {
          ctx
            .setState(gateway.id)
            .then(() =>
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
      isPrimary={ctx.activeFegGatewayId === gateway.id}
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
  const {history, relativeUrl} = useRouter();
  const ctx = useContext(FEGGatewayContext);

  return (
    <TableRow key={gateway.id}>
      <TableCell>
        <span className={classes.gatewayName}>{gateway.name}</span>
        {
          <DeviceStatusCircle
            isGrey={!ctx.health[gateway.id]?.status}
            isActive={ctx.health[gateway.id]?.status === HEALTHY_STATUS}
          />
        }
        {isPrimary && (
          <Tooltip title="Primary FEG" placement="right">
            <StarIcon className={classes.star} />
          </Tooltip>
        )}
      </TableCell>
      <TableCell>{gateway.device?.hardware_id}</TableCell>
      <TableCell>
        <IconButton
          color="primary"
          onClick={() => history.push(relativeUrl(`/edit/${gateway.id}`))}>
          <EditIcon />
        </IconButton>
        <IconButton color="primary" onClick={() => onDelete(gateway)}>
          <DeleteIcon data-testid={`delete ${gateway.id}`} />
        </IconButton>
      </TableCell>
    </TableRow>
  );
}

export default withAlert(CWFGateways);
