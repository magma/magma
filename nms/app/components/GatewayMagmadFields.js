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

import Button from '@material-ui/core/Button';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import FormControl from '@material-ui/core/FormControl';
import InputLabel from '@material-ui/core/InputLabel';
import MagmaV1API from '../../generated/WebClient';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Select from '@material-ui/core/Select';
import TextField from '@material-ui/core/TextField';

// $FlowFixMe migrated to typescript
import nullthrows from '../../shared/util/nullthrows';
import {makeStyles} from '@material-ui/styles';
import {toString} from './GatewayUtils';
import {useParams} from 'react-router-dom';
import {useState} from 'react';

const useStyles = makeStyles(() => ({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
}));

type Props = {
  onClose: () => void,
  onSave: (gatewayID: string) => void,
  gateway: GatewayV1,
};

export default function GatewayMagmadFields(props: Props) {
  const classes = useStyles();
  const params = useParams();
  const {gateway} = props;
  const [autoupgradeEnabled, setAutoupgradeEnabled] = useState(
    gateway.autoupgradeEnabled,
  );
  const [autoupgradePollInterval, setAutoupgradePollInterval] = useState(
    toString(gateway.autoupgradePollInterval),
  );
  const [checkinInterval, setCheckinInterval] = useState(
    toString(gateway.checkinInterval),
  );
  const [checkinTimeout, setCheckinTimeout] = useState(
    toString(gateway.checkinTimeout),
  );

  const onSave = () => {
    const magmad = {
      autoupgrade_enabled: autoupgradeEnabled,
      autoupgrade_poll_interval: parseInt(autoupgradePollInterval),
      checkin_interval: parseInt(checkinInterval),
      checkin_timeout: parseInt(checkinTimeout),
      tier: gateway.tier,
    };

    MagmaV1API.putLteByNetworkIdGatewaysByGatewayIdMagmad({
      networkId: nullthrows(params.networkId),
      gatewayId: gateway.logicalID,
      magmad,
    }).then(() => props.onSave(gateway.logicalID));
  };

  return (
    <>
      <DialogContent>
        <FormControl className={classes.input}>
          <InputLabel htmlFor="autoupgradeEnabled">
            Autoupgrade Enabled
          </InputLabel>
          <Select
            inputProps={{id: 'autoupgradeEnabled'}}
            value={autoupgradeEnabled ? 1 : 0}
            onChange={({target}) => setAutoupgradeEnabled(!!target.value)}>
            <MenuItem value={1}>Enabled</MenuItem>
            <MenuItem value={0}>Disabled</MenuItem>
          </Select>
        </FormControl>
        <TextField
          label="Autoupgrade Poll Interval (seconds)"
          className={classes.input}
          value={autoupgradePollInterval}
          onChange={({target}) => setAutoupgradePollInterval(target.value)}
          placeholder="E.g. 300"
        />
        <TextField
          label="Checkin Interval (seconds)"
          className={classes.input}
          value={checkinInterval}
          onChange={({target}) => setCheckinInterval(target.value)}
          placeholder="E.g. 60"
        />
        <TextField
          label="Checkin Timeout (seconds)"
          className={classes.input}
          value={checkinTimeout}
          onChange={({target}) => setCheckinTimeout(target.value)}
          placeholder="E.g. 5"
        />
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose}>Cancel</Button>
        <Button onClick={onSave} variant="contained" color="primary">
          Save
        </Button>
      </DialogActions>
    </>
  );
}
