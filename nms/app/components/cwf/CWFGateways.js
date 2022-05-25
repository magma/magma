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

import type {WithAlert} from '../Alert/withAlert';
import type {cwf_gateway} from '../../../generated/MagmaAPIBindings';
import type {cwf_ha_pair} from '../../../generated/MagmaAPIBindings';

// $FlowFixMe migrated to typescript
import AddGatewayDialog from '../AddGatewayDialog';
import Button from '@material-ui/core/Button';
import CWFEditGatewayDialog from './CWFEditGatewayDialog';
import ChevronRight from '@material-ui/icons/ChevronRight';
import DeleteIcon from '@material-ui/icons/Delete';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import DeviceStatusCircle from '../../theme/design-system/DeviceStatusCircle';
import EditIcon from '@material-ui/icons/Edit';
import ExpandMore from '@material-ui/icons/ExpandMore';
import IconButton from '@material-ui/core/IconButton';
import MagmaV1API from '../../../generated/WebClient';
// $FlowFixMe migrated to typescript
import NestedRouteLink from '../NestedRouteLink';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import StarIcon from '@material-ui/icons/Star';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Tooltip from '@material-ui/core/Tooltip';

// $FlowFixMe migrated to typescript
import LoadingFiller from '../LoadingFiller';
import Text from '../../theme/design-system/Text';
// $FlowFixMe migrated to typescript
import nullthrows from '../../../shared/util/nullthrows';
import useMagmaAPI from '../../../api/useMagmaAPIFlow';
import withAlert from '../Alert/withAlert';
// $FlowFixMe migrated to typescript
import {MAGMAD_DEFAULT_CONFIGS} from '../AddGatewayDialog';
import {Route, Routes, useNavigate, useParams} from 'react-router-dom';
import {colors} from '../../theme/default';
import {findIndex} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {map} from 'lodash';
import {useCallback, useState} from 'react';
import {useInterval} from '../../../app/hooks';

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
  gatewayCell: {
    padding: '5px',
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
    color: colors.primary.brightGray,
    fontWeight: 'bolder',
    paddingRight: '10px',
  },
  star: {
    color: '#ffd700',
    width: '18px',
    verticalAlign: 'bottom',
  },
}));

const FIVE_MINS = 5 * 60 * 1000;
const REFRESH_INTERVAL = 2 * 60 * 1000;

function gatewayStatus(gateway: cwf_gateway): string {
  const gatewayHealthy =
    Math.max(0, Date.now() - (gateway.status?.checkin_time || 0)) < FIVE_MINS;
  let status = '';
  if (!gatewayHealthy) {
    const checkInTime = new Date(gateway.status?.checkin_time ?? 0);
    status = 'Last refreshed ' + checkInTime.toLocaleString();
  } else {
    if (gateway.carrier_wifi.allowed_gre_peers.length == 0) {
      status = 'Gateway is not functional. No GRE peers configured';
    }
  }
  return status;
}

function EditDialog(props: {
  setGateways: (gateways: cwf_gateway[]) => void,
  gateways: cwf_gateway[],
}) {
  const navigate = useNavigate();
  const params = useParams();

  return (
    <CWFEditGatewayDialog
      gateway={nullthrows(
        props.gateways.find(gw => gw.id === params.gatewayID),
      )}
      onCancel={() => navigate('..')}
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

export function CWFGateways(props: WithAlert & {}) {
  const [gateways, setGateways] = useState<?(cwf_gateway[])>(null);
  const [haPairs, setHaPairs] = useState<?(cwf_ha_pair[])>(null);
  const params = useParams();
  const navigate = useNavigate();
  const [lastFetchTime, setLastFetchTime] = useState(Date.now());
  const networkId = nullthrows(params.networkId);
  const classes = useStyles();

  useMagmaAPI(
    MagmaV1API.getCwfByNetworkIdGateways,
    {networkId},
    useCallback(response => setGateways(map(response, g => g)), []),
    lastFetchTime,
  );

  useMagmaAPI(
    MagmaV1API.getCwfByNetworkIdHaPairs,
    {networkId},
    useCallback(response => setHaPairs(map(response, h => h)), []),
    lastFetchTime,
  );

  useInterval(() => setLastFetchTime(Date.now()), REFRESH_INTERVAL);

  if (!gateways || !haPairs) {
    return <LoadingFiller />;
  }

  const deleteGateway = (gateway: cwf_gateway) => {
    props
      .confirm(`Are you sure you want to delete ${gateway.name}?`)
      .then(confirmed => {
        if (confirmed) {
          MagmaV1API.deleteCwfByNetworkIdGatewaysByGatewayId({
            networkId,
            gatewayId: gateway.id,
          }).then(() =>
            setGateways(gateways.filter(gw => gw.id != gateway.id)),
          );
        }
      });
  };

  const addGateway = async ({
    gatewayID,
    name,
    description,
    hardwareID,
    challengeKey,
    tier,
  }) => {
    await MagmaV1API.postCwfByNetworkIdGateways({
      networkId,
      gateway: {
        carrier_wifi: {
          allowed_gre_peers: [],
        },
        description,
        device: {
          hardware_id: hardwareID,
          key: {
            key: challengeKey,
            key_type: 'SOFTWARE_ECDSA_SHA256', // default key type
          },
        },
        id: gatewayID,
        magmad: MAGMAD_DEFAULT_CONFIGS,
        name,
        tier,
      },
    });

    const gateway = await MagmaV1API.getCwfByNetworkIdGatewaysByGatewayId({
      networkId,
      gatewayId: gatewayID,
    });

    setGateways([...gateways, gateway]);
    navigate('');
  };

  const rows = gateways.map(gateway => (
    <GatewayRow
      key={gateway.id}
      gateway={gateway}
      haPairs={haPairs}
      onDelete={deleteGateway}
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
              <TableCell>Hardware UUID / GRE Key</TableCell>
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
            <AddGatewayDialog
              onClose={() => navigate('')}
              onSave={addGateway}
            />
          }
        />
        <Route
          path="/edit/:gatewayID"
          element={<EditDialog gateways={gateways} setGateways={setGateways} />}
        />
      </Routes>
    </div>
  );
}

function GatewayRow(props: {
  gateway: cwf_gateway,
  haPairs: cwf_ha_pair[],
  onDelete: cwf_gateway => void,
}) {
  const {gateway, haPairs, onDelete} = props;
  const [expanded, setExpanded] = useState<Set<string>>(new Set());
  const classes = useStyles();
  const navigate = useNavigate();

  const gatewayHaPair = haPairs.filter(haPair => {
    return (
      haPair.gateway_id_1 === gateway.id || haPair.gateway_id_2 === gateway.id
    );
  });

  const isPrimary =
    gatewayHaPair?.[0]?.state?.ha_pair_status?.active_gateway === gateway.id;
  const isGateway1 = gatewayHaPair?.[0]?.gateway_id_1 === gateway.id;

  const isNonHaGatewayHealthy =
    Math.max(0, Date.now() - (gateway.status?.checkin_time || 0)) < FIVE_MINS &&
    gateway.carrier_wifi.allowed_gre_peers.length > 0;
  const gatewayHealth = isGateway1
    ? gatewayHaPair[0]?.state?.gateway1_health?.status
    : gatewayHaPair?.[0]
    ? gatewayHaPair[0]?.state?.gateway2_health?.status
    : isNonHaGatewayHealthy
    ? 'HEALTHY'
    : 'UNHEALTHY';

  return (
    <>
      <TableRow key={gateway.id}>
        <Tooltip title={gatewayStatus(gateway)} placement={'bottom-start'}>
          <TableCell className={classes.gatewayCell}>
            <IconButton
              className={classes.expandIconButton}
              onClick={() => {
                const newExpanded = new Set(expanded);
                expanded.has(gateway.id)
                  ? newExpanded.delete(gateway.id)
                  : newExpanded.add(gateway.id);
                setExpanded(newExpanded);
              }}>
              {expanded.has(gateway.id) ? <ExpandMore /> : <ChevronRight />}
            </IconButton>

            <span className={classes.gatewayName}>{gateway.name}</span>
            <DeviceStatusCircle
              isGrey={!gateway.status?.checkin_time}
              isActive={gatewayHealth === 'HEALTHY'}
            />
            {isPrimary && (
              <Tooltip title="Primary CWAG" placement="right">
                <StarIcon className={classes.star} />
              </Tooltip>
            )}
          </TableCell>
        </Tooltip>

        <TableCell>{gateway.device?.hardware_id}</TableCell>
        <TableCell>
          <IconButton
            color="primary"
            onClick={() => navigate(`edit/${gateway.id}`)}>
            <EditIcon />
          </IconButton>
          <IconButton color="primary" onClick={() => onDelete(gateway)}>
            <DeleteIcon />
          </IconButton>
        </TableCell>
      </TableRow>
      {expanded.has(gateway.id) &&
        gateway.carrier_wifi.allowed_gre_peers.map((gre, i) => (
          <TableRow key={i} classeName={classes.tableRow}>
            <TableCell className={classes.greCell}>{gre.ip}</TableCell>
            <TableCell>{gre.key}</TableCell>
            <TableCell />
          </TableRow>
        ))}
    </>
  );
}

export default withAlert(CWFGateways);
