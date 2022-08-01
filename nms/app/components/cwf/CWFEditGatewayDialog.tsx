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

import type {CwfGateway, MagmadGatewayConfigs} from '../../../generated';

import Button from '@mui/material/Button';
import CWFGatewayConfigFields from './CWFGatewayConfigFields';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import GatewayCommandFields from '../GatewayCommandFields';
import MagmaDeviceFields from '../MagmaDeviceFields';
import React from 'react';
import Tab from '@mui/material/Tab';
import Tabs from '@mui/material/Tabs';

import MagmaAPI from '../../api/MagmaAPI';
import nullthrows from '../../../shared/util/nullthrows';
import {colors} from '../../theme/default';
import {getErrorMessage} from '../../util/ErrorUtils';
import {makeStyles} from '@mui/styles';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';
import {useParams} from 'react-router-dom';
import {useState} from 'react';

const useStyles = makeStyles({
  tabBar: {
    backgroundColor: colors.primary.brightGray,
  },
});

type Props = {
  gateway: CwfGateway;
  onCancel: () => void;
  onSave: (gateway: CwfGateway) => void;
};

export default function (props: Props) {
  const [tab, setTab] = useState(0);
  const [magmaConfigs, setMagmaConfigs] = useState(props.gateway.magmad);
  const [allowedGREPeers, setAllowedGREPeers] = useState(
    props.gateway.carrier_wifi.allowed_gre_peers,
  );
  const [ipdrExportDst, setIPDRExportDst] = useState(
    props.gateway.carrier_wifi.ipdr_export_dst,
  );

  const classes = useStyles();
  const params = useParams();
  const enqueueSnackbar = useEnqueueSnackbar();

  const gatewayID = nullthrows(params.gatewayID);
  const networkID = nullthrows(params.networkId);
  const onSave = async () => {
    try {
      await MagmaAPI.carrierWifiGateways.cwfNetworkIdGatewaysGatewayIdPut({
        networkId: networkID,
        gatewayId: gatewayID,
        gateway: {
          ...props.gateway,
          carrier_wifi: {
            allowed_gre_peers: allowedGREPeers,
            ipdr_export_dst: ipdrExportDst,
          },
          magmad: getMagmaConfigs(magmaConfigs),
        },
      });
      props.onSave(
        (
          await MagmaAPI.carrierWifiGateways.cwfNetworkIdGatewaysGatewayIdGet({
            networkId: networkID,
            gatewayId: gatewayID,
          })
        ).data,
      );
    } catch (e) {
      enqueueSnackbar(getErrorMessage(e), {
        variant: 'error',
      });
    }
  };

  let content;
  switch (tab) {
    case 0:
      content = (
        <CWFGatewayConfigFields
          allowedGREPeers={allowedGREPeers}
          onChange={setAllowedGREPeers}
          ipdrExportDst={ipdrExportDst}
          onIPDRChanged={setIPDRExportDst}
        />
      );
      break;
    case 1:
      content = (
        <MagmaDeviceFields
          configs={magmaConfigs}
          configChangeHandler={(fieldName, value) =>
            setMagmaConfigs({
              ...magmaConfigs,
              [fieldName]: value,
            })
          }
        />
      );
      break;
    case 2:
      content = (
        <GatewayCommandFields
          gatewayID={gatewayID}
          showRestartCommand={true}
          showRebootEnodebCommand={false}
          showPingCommand={true}
          showGenericCommand={true}
        />
      );
      break;

    default:
      // should never happen
      content = <div />;
  }

  return (
    <Dialog open={true} onClose={props.onCancel} maxWidth="md" scroll="body">
      <Tabs
        className={classes.tabBar}
        indicatorColor="primary"
        value={tab}
        onChange={(_, tab) => setTab(tab as number)}>
        <Tab label="Carrier Wifi" />
        <Tab label="Controller" />
        <Tab label="Command" />
      </Tabs>
      <DialogContent>{content}</DialogContent>
      <DialogActions>
        <Button onClick={props.onCancel}>Cancel</Button>
        <Button
          onClick={() => void onSave()}
          variant="contained"
          color="primary">
          Save
        </Button>
      </DialogActions>
    </Dialog>
  );
}

function getMagmaConfigs(
  magmaConfigs: MagmadGatewayConfigs,
): MagmadGatewayConfigs {
  //TODO [ts-migration]: The parseInt calls should be superfluous
  return {
    ...magmaConfigs,
    // @ts-ignore
    autoupgrade_poll_interval: parseInt(magmaConfigs.autoupgrade_poll_interval),
    // @ts-ignore
    checkin_interval: parseInt(magmaConfigs.checkin_interval),
    // @ts-ignore
    checkin_timeout: parseInt(magmaConfigs.checkin_timeout),
  };
}
