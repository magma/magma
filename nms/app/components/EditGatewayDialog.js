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
import type {lte_gateway} from '../../generated/MagmaAPIBindings';

import AppBar from '@material-ui/core/AppBar';
import Dialog from '@material-ui/core/Dialog';
import GatewayCellularFields from './GatewayCellularFields';
import GatewayCommandFields from './GatewayCommandFields';
import GatewayMagmadFields from './GatewayMagmadFields';
import GatewaySummaryFields from './GatewaySummaryFields';
import MagmaV1API from '../../generated/WebClient';
import React, {useState} from 'react';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';

import nullthrows from '../../fbc_js_core/util/nullthrows';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '../../fbc_js_core/ui/hooks/useSnackbar';
import {useParams} from 'react-router-dom';

const useStyles = makeStyles(() => ({
  appBar: {
    backgroundColor: '#f5f5f5',
    marginBottom: '20px',
  },
}));

type Props = {|
  onClose: () => void,
  onSave: lte_gateway => void,
  gateway: ?GatewayV1,
|};

export default function EditGatewayDialog({onClose, onSave, gateway}: Props) {
  const [tab, setTab] = useState(0);
  const classes = useStyles();
  const params = useParams();
  const enqueueSnackbar = useEnqueueSnackbar();

  const networkID = nullthrows(params.networkId);

  const wrappedOnSave = gatewayID => {
    MagmaV1API.getLteByNetworkIdGatewaysByGatewayId({
      networkId: networkID,
      gatewayId: gatewayID,
    })
      .then(onSave)
      .catch(e => {
        enqueueSnackbar(e?.response?.data?.message || e?.message || e, {
          variant: 'error',
        });
      });
  };

  if (!gateway) {
    return null;
  }

  let content;
  const childProps = {
    onClose,
    gateway,
    onSave: wrappedOnSave,
  };

  switch (tab) {
    case 0:
      content = <GatewaySummaryFields {...childProps} />;
      break;
    case 1:
      content = <GatewayCellularFields {...childProps} />;
      break;
    case 2:
      content = <GatewayMagmadFields {...childProps} />;
      break;
    case 3:
      content = (
        <GatewayCommandFields
          onClose={onClose}
          gatewayID={gateway.logicalID}
          showRestartCommand={true}
          showRebootEnodebCommand={true}
          showPingCommand={true}
          showGenericCommand={true}
        />
      );
      break;
  }
  return (
    <Dialog open={true} onClose={onClose} maxWidth="md" scroll="body">
      <AppBar position="static" className={classes.appBar}>
        <Tabs
          indicatorColor="primary"
          textColor="primary"
          value={tab}
          onChange={(event, tab) => setTab(tab)}>
          <Tab label="Summary" />
          <Tab label="LTE" />
          <Tab label="Magma" />
          <Tab label="Commands" />
        </Tabs>
      </AppBar>
      {content}
    </Dialog>
  );
}
