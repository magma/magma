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

import FormControl from '@material-ui/core/FormControl';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Select from '@material-ui/core/Select';
import TextField from '@material-ui/core/TextField';

import {MagmadGatewayConfigs} from '../../generated-ts';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
}));

type Props = {
  configs: MagmadGatewayConfigs;
  configChangeHandler: <K extends keyof MagmadGatewayConfigs>(
    key: K,
    value: MagmadGatewayConfigs[K],
  ) => void;
};

export default function MagmaDeviceFields(props: Props) {
  const classes = useStyles();
  return (
    <>
      <FormControl className={classes.input}>
        <InputLabel htmlFor="autoupgradeEnabled">
          Autoupgrade Enabled
        </InputLabel>
        <Select
          inputProps={{id: 'autoupgradeEnabled'}}
          value={props.configs.autoupgrade_enabled ? 1 : 0}
          onChange={({target}) =>
            props.configChangeHandler('autoupgrade_enabled', !!target.value)
          }>
          <MenuItem value={1}>Enabled</MenuItem>
          <MenuItem value={0}>Disabled</MenuItem>
        </Select>
      </FormControl>
      <TextField
        label="Autoupgrade Poll Interval (seconds)"
        className={classes.input}
        value={props.configs.autoupgrade_poll_interval}
        onChange={({target}) =>
          props.configChangeHandler(
            'autoupgrade_poll_interval',
            parseInt(target.value),
          )
        }
        placeholder="E.g. 300"
      />
      <TextField
        label="Checkin Interval (seconds)"
        className={classes.input}
        value={props.configs.checkin_interval}
        onChange={({target}) =>
          props.configChangeHandler('checkin_interval', parseInt(target.value))
        }
        placeholder="E.g. 60"
      />
      <TextField
        label="Checkin Timeout (seconds)"
        className={classes.input}
        value={props.configs.checkin_timeout}
        onChange={({target}) =>
          props.configChangeHandler('checkin_timeout', parseInt(target.value))
        }
        placeholder="E.g. 5"
      />
    </>
  );
}
