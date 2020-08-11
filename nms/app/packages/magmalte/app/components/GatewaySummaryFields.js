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

import type {GatewayV1} from './GatewayUtils';

import Button from '@fbcnms/ui/components/design-system/Button';
import Check from '@material-ui/icons/Check';
import DeviceStatusCircle from '@fbcnms/ui/components/icons/DeviceStatusCircle';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import Divider from '@material-ui/core/Divider';
import Fade from '@material-ui/core/Fade';
import FormField from './FormField';
import Input from '@material-ui/core/Input';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import moment from 'moment';

import nullthrows from '@fbcnms/util/nullthrows';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';

type Props = {
  onClose: () => void,
  onSave: (gatewayID: string) => void,
  gateway: GatewayV1,
};

const useStyles = makeStyles(() => ({
  input: {
    width: '100%',
  },
  divider: {
    margin: '10px 0',
  },
}));

export default function GatewaySummaryFields(props: Props) {
  const classes = useStyles();
  const {match} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();
  const {gateway} = props;
  const [name, setName] = useState(gateway.name);
  const [showRebootCheck, setShowRebootCheck] = useState(false);
  const [showRestartCheck, setShowRestartCheck] = useState(false);

  const id = gateway.logicalID;
  const onSave = () => {
    MagmaV1API.putLteByNetworkIdGatewaysByGatewayIdName({
      networkId: nullthrows(match.params.networkId),
      gatewayId: id,
      name: JSON.stringify(`"${name}"`),
    })
      .then(() => props.onSave(id))
      .catch(error =>
        enqueueSnackbar(error.response.data.message, {variant: 'error'}),
      );
  };

  const handleRebootGateway = () => {
    MagmaV1API.postNetworksByNetworkIdGatewaysByGatewayIdCommandReboot({
      networkId: nullthrows(match.params.networkId),
      gatewayId: id,
    })
      .then(_resp => {
        enqueueSnackbar('Successfully initiated reboot', {variant: 'success'});
        setShowRebootCheck(true);
        setTimeout(() => setShowRebootCheck(false), 5000);
      })
      .catch(error =>
        enqueueSnackbar('Reboot failed: ' + error.response.data.message, {
          variant: 'error',
        }),
      );
  };

  const handleRestartServices = () => {
    MagmaV1API.postNetworksByNetworkIdGatewaysByGatewayIdCommandRestartServices(
      {
        networkId: nullthrows(match.params.networkId),
        gatewayId: id,
        services: [],
      },
    )
      .then(_resp => {
        enqueueSnackbar('Successfully initiated service restart', {
          variant: 'success',
        });
        setShowRestartCheck(true);
        setTimeout(() => setShowRestartCheck(false), 5000);
      })
      .catch(error =>
        enqueueSnackbar(
          'Restart services failed: ' + error.response.data.message,
          {variant: 'error'},
        ),
      );
  };

  return (
    <>
      <DialogContent>
        <FormField label="Name">
          <Input
            className={classes.input}
            value={name}
            onChange={({target}) => setName(target.value)}
            placeholder="E.g. Gateway 1234"
          />
        </FormField>
        <FormField label="Hardware UUID">{gateway.hardware_id}</FormField>
        <FormField label="Gateway ID">{gateway.logicalID}</FormField>
        <FormField label="Last Checkin">
          {moment(parseInt(gateway.lastCheckin, 10)).fromNow()}
        </FormField>
        <FormField label="Version">{gateway.version}</FormField>
        <FormField label="VPN IP">{gateway.vpnIP}</FormField>
        <FormField label="RF Transmitter">
          <DeviceStatusCircle
            isGrey={false}
            isActive={gateway.enodebRFTXEnabled}
          />
          {gateway.enodebRFTXEnabled ? '' : 'Not '}
          Allowed
          {'  '}
          <DeviceStatusCircle
            isGrey={gateway.isBackhaulDown}
            isActive={gateway.enodebConnected && gateway.enodebRFTXOn}
          />
          {gateway.enodebRFTXOn ? '' : 'Not '}
          Connected
        </FormField>
        <FormField label="GPS synchronized">
          <DeviceStatusCircle
            isGrey={gateway.isBackhaulDown}
            isActive={gateway.enodebConnected && gateway.gpsConnected}
          />
          {gateway.gpsConnected ? '' : 'Not '}
          Synced
        </FormField>
        <FormField label="MME">
          <DeviceStatusCircle
            isGrey={gateway.isBackhaulDown}
            isActive={gateway.enodebConnected && gateway.mmeConnected}
          />
          {gateway.mmeConnected ? '' : 'Not '}
          Connected
        </FormField>
        <Divider className={classes.divider} />
        <Text variant="subtitle1">Commands</Text>
        <FormField label="Reboot Gateway">
          <Button onClick={handleRebootGateway} variant="text">
            Reboot
          </Button>
          <Fade in={showRebootCheck} timeout={500}>
            <Check style={{verticalAlign: 'middle'}} htmlColor="green" />
          </Fade>
        </FormField>
        <FormField label="">
          <Button onClick={handleRestartServices} variant="text">
            Restart services
          </Button>
          <Fade in={showRestartCheck} timeout={500}>
            <Check style={{verticalAlign: 'middle'}} htmlColor="green" />
          </Fade>
        </FormField>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose} skin="regular">
          Cancel
        </Button>
        <Button onClick={onSave}>Save</Button>
      </DialogActions>
    </>
  );
}
