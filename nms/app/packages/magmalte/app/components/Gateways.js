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
 * @flow
 * @format
 */

import type {GatewayV1} from './GatewayUtils';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {lte_gateway} from '@fbcnms/magma-api';

import AddGatewayDialog from './AddGatewayDialog';
import Button from '@fbcnms/ui/components/design-system/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import DeviceStatusCircle from '@fbcnms/ui/components/icons/DeviceStatusCircle';
import EditGatewayDialog from './EditGatewayDialog';
import EditIcon from '@material-ui/icons/Edit';
import IconButton from '@material-ui/core/IconButton';
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

import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {MAGMAD_DEFAULT_CONFIGS} from './AddGatewayDialog';
import {Route} from 'react-router-dom';
import {find} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useState} from 'react';
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

export default function Gateways() {
  const classes = useStyles();
  const {match, history, relativePath, relativeUrl} = useRouter();
  const [gateways, setGateways] = useState();
  const [lastFetchTime, setLastFetchTime] = useState(Date.now());

  const networkId = nullthrows(match.params.networkId);
  const {isLoading} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdGateways,
    {networkId: networkId},
    useCallback(
      response =>
        setGateways(
          Object.keys(response)
            .map(k => response[k])
            .filter(g => g.cellular && g.id)
            .map(_buildGatewayFromPayload),
        ),
      [],
    ),
    lastFetchTime,
  );

  if (isLoading || !gateways) {
    return <LoadingFiller />;
  }

  const onGatewayAdd = async ({
    gatewayID,
    name,
    description,
    hardwareID,
    challengeKey,
    tier,
  }) => {
    await MagmaV1API.postLteByNetworkIdGateways({
      networkId,
      gateway: {
        id: gatewayID,
        name,
        description,
        cellular: {
          epc: {nat_enabled: true, ip_block: '192.168.128.0/24'},
          ran: {pci: 260, transmit_enabled: false},
          non_eps_service: undefined,
        },
        magmad: MAGMAD_DEFAULT_CONFIGS,
        device: {
          hardware_id: hardwareID,
          key: {
            key: challengeKey,
            key_type: 'SOFTWARE_ECDSA_SHA256', // default key/challenge type
          },
        },
        connected_enodeb_serials: [],
        tier,
      },
    });

    history.push(relativeUrl(''));
    setLastFetchTime(Date.now());
  };

  const rows = gateways.map(gateway => (
    <GatewayRow
      key={gateway.logicalID}
      gateway={gateway}
      onSave={() => setLastFetchTime(Date.now())}
    />
  ));

  return (
    <div className={classes.paper}>
      <div className={classes.header}>
        <Text variant="h5">Configure Gateways</Text>
        <NestedRouteLink to="/new">
          <Button>Add Gateway</Button>
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
          <AddGatewayDialog
            onClose={() => history.push(relativeUrl(''))}
            onSave={onGatewayAdd}
          />
        )}
      />
    </div>
  );
}

function _buildGatewayFromPayload(gateway: lte_gateway): GatewayV1 {
  let enodebRFTXOn = false;
  let enodebConnected = false;
  let gpsConnected = false;
  let mmeConnected = false;
  let version = 'Not Reported';
  let vpnIP = 'Not Reported';
  let lastCheckin = 'Not Reported';
  let hardwareID = 'Not reported';
  let isBackhaulDown = true;
  const latLon = {lat: 0, lon: 0};
  const {status} = gateway;
  if (status) {
    vpnIP = status.platform_info?.vpn_ip || vpnIP;
    const packages = find(status.platform_info?.packages || [], {
      name: 'magma',
    });
    version = packages?.version || '';
    // if the last check-in time is more than 5 minutes
    // we treat it as backhaul is down
    const checkin = status.checkin_time;
    if (checkin != null) {
      const duration = Math.max(0, Date.now() - checkin);
      isBackhaulDown = duration > 1000 * 5 * 60;
      lastCheckin = checkin.toString();
    }

    const {meta} = status;
    if (meta) {
      if (!isBackhaulDown) {
        enodebRFTXOn = status.meta && status.meta.rf_tx_on;
      }

      latLon.lat = parseFloat(meta.gps_latitude);
      latLon.lon = parseFloat(meta.gps_longitude);
      gpsConnected = meta.gps_connected == '1';
      enodebConnected = meta.enodeb_connected == '1';
      mmeConnected = meta.mme_connected == '1';
    }

    if (status.hardware_id) {
      hardwareID = status.hardware_id;
    }
  }

  let autoupgradePollInterval,
    checkinInterval,
    checkinTimeout,
    autoupgradeEnabled;
  if (gateway.magmad) {
    const {magmad} = gateway;
    autoupgradePollInterval = magmad.autoupgrade_poll_interval;
    checkinInterval = magmad.checkin_interval;
    checkinTimeout = magmad.checkin_timeout;
    autoupgradeEnabled = magmad.autoupgrade_enabled;
  }

  let ipBlock, natEnabled;
  let pci, transmitEnabled;
  let control, csfbRAT, csfbMCC, csfbMNC, lac;
  if (gateway.cellular) {
    const {cellular} = gateway;
    ipBlock = cellular?.epc?.ip_block;
    natEnabled = cellular?.epc?.nat_enabled;
    pci = cellular?.ran?.pci;
    transmitEnabled = cellular?.ran?.transmit_enabled ?? false;

    const nonEPSService = cellular.non_eps_service || {};
    control = nonEPSService.non_eps_service_control;
    csfbRAT = nonEPSService.csfb_rat;
    csfbMNC = nonEPSService.csfb_mnc;
    csfbMCC = nonEPSService.csfb_mcc;
    lac = nonEPSService.lac;
  }

  return {
    hardware_id: hardwareID,
    name: gateway.name || 'N/A',
    logicalID: gateway.id,
    challengeType: gateway?.device?.key?.key_type || '',
    enodebRFTXEnabled: !!transmitEnabled,
    enodebRFTXOn: !!enodebRFTXOn,
    enodebConnected,
    gpsConnected,
    isBackhaulDown,
    lastCheckin,
    latLon,
    mmeConnected,
    version,
    vpnIP,
    autoupgradePollInterval,
    checkinInterval,
    checkinTimeout,
    tier: gateway.tier,
    autoupgradeEnabled: !!autoupgradeEnabled,
    attachedEnodebSerials: gateway.connected_enodeb_serials || [],
    ran: {
      pci,
      transmitEnabled: !!transmitEnabled,
    },
    epc: {
      ipBlock: ipBlock || '',
      natEnabled: natEnabled || false,
    },
    nonEPSService: {
      control: control || 0,
      csfbRAT: csfbRAT || 0,
      csfbMCC,
      csfbMNC,
      lac,
    },
    rawGateway: gateway,
  };
}

type Props = WithAlert & {onSave: () => void, gateway: GatewayV1};
function GatewayRowItem(props: Props) {
  const {match, history, relativePath, relativeUrl} = useRouter();
  const {gateway} = props;

  const deleteGateway = () => {
    props
      .confirm(`Are you sure you want to delete ${gateway.name}?`)
      .then(confirmed => {
        if (!confirmed) {
          return;
        }
        MagmaV1API.deleteLteByNetworkIdGatewaysByGatewayId({
          networkId: nullthrows(match.params.networkId),
          gatewayId: gateway.logicalID,
        }).then(props.onSave);
      });
  };

  return (
    <TableRow>
      <TableCell>
        <DeviceStatusCircle
          isGrey={!gateway.enodebRFTXOn}
          isActive={gateway.enodebRFTXOn === gateway.enodebRFTXEnabled}
        />
        {gateway.name}
      </TableCell>
      <TableCell>{gateway.hardware_id}</TableCell>
      <TableCell>
        <IconButton
          data-testid="edit-gateway-icon"
          color="primary"
          onClick={() =>
            history.push(relativeUrl(`/edit/${gateway.logicalID}`))
          }>
          <EditIcon />
        </IconButton>
        <IconButton
          data-testid="delete-gateway-icon"
          color="primary"
          onClick={deleteGateway}>
          <DeleteIcon />
        </IconButton>
      </TableCell>
      <Route
        path={relativePath(`/edit/${gateway.logicalID}`)}
        render={() => (
          <EditGatewayDialog
            gateway={gateway}
            onClose={() => history.push(relativeUrl(''))}
            onSave={() => {
              props.onSave();
              history.push(relativeUrl(''));
            }}
          />
        )}
      />
    </TableRow>
  );
}

const GatewayRow = withAlert(GatewayRowItem);
