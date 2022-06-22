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
 */

import type {FederationGateway} from '../../../generated-ts';
import type {WithAlert} from '../Alert/withAlert';

import Button from '@material-ui/core/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import DeviceStatusCircle from '../../theme/design-system/DeviceStatusCircle';
import EditIcon from '@material-ui/icons/Edit';
import FEGGatewayContext from '../context/FEGGatewayContext';
import FEGGatewayDialog from './FEGGatewayDialog';
import IconButton from '@material-ui/core/IconButton';
import NestedRouteLink from '../NestedRouteLink';
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
import nullthrows from '../../../shared/util/nullthrows';
import withAlert from '../Alert/withAlert';
import {HEALTHY_STATUS} from '../GatewayUtils';
import {Route, Routes, useNavigate, useParams} from 'react-router-dom';
import {Theme} from '@material-ui/core/styles';
import {colors} from '../../theme/default';
import {findIndex} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useContext, useState} from 'react';

const useStyles = makeStyles<Theme>(theme => ({
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

function EditDialog(props: {
  setGateways: (gateways: Array<FederationGateway>) => void;
  gateways: Array<FederationGateway>;
}) {
  const navigate = useNavigate();
  const params = useParams();

  return (
    <FEGGatewayDialog
      editingGateway={nullthrows(
        props.gateways.find(gw => gw.id === params.gatewayID),
      )}
      onClose={() => navigate('..')}
      onSave={gateway => {
        const newGateways = [...props.gateways];
        const i = findIndex(newGateways, g => g.id === gateway.id);
        newGateways[i] = gateway;
        props.setGateways(newGateways);
        navigate('..');
      }}
    />
  );
}

function CWFGateways(props: WithAlert) {
  const ctx = useContext(FEGGatewayContext);
  const [gateways, setGateways] = useState<Array<FederationGateway>>(
    Object.keys(ctx.state).map(gatewayId => ctx.state[gatewayId]),
  );
  const navigate = useNavigate();
  const classes = useStyles();
  const deleteGateway = (gateway: FederationGateway) => {
    void props
      .confirm(`Are you sure you want to delete ${gateway.name}?`)
      .then(confirmed => {
        if (confirmed) {
          void ctx
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
        <NestedRouteLink to="new">
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
      <Routes>
        <Route
          path="/new"
          element={
            <FEGGatewayDialog
              onClose={() => navigate('')}
              onSave={gateway => {
                setGateways([...gateways, gateway]);
                navigate('');
              }}
            />
          }
        />
        <Route
          path="edit/:gatewayID"
          element={<EditDialog gateways={gateways} setGateways={setGateways} />}
        />
      </Routes>
    </div>
  );
}

function GatewayRow(props: {
  gateway: FederationGateway;
  onDelete: (gateway: FederationGateway) => void;
  isPrimary: boolean;
}) {
  const classes = useStyles();
  const {gateway, onDelete, isPrimary} = props;
  const navigate = useNavigate();
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
          onClick={() => navigate(`edit/${gateway.id}`)}>
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
