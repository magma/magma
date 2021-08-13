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

import type {gateway_status, gateway_wifi_configs} from '@fbcnms/magma-api';

import Checkbox from '@material-ui/core/Checkbox';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import FormGroup from '@material-ui/core/FormGroup';
import KeyValueFields from '@fbcnms/magmalte/app/components/KeyValueFields';
import React from 'react';
import TextField from '@material-ui/core/TextField';

import nullthrows from '@fbcnms/util/nullthrows';
import {getAdditionalProp, setAdditionalProp} from './WifiUtils';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
}));

type Props = {
  macAddress: string,
  status: ?gateway_status,
  configs: gateway_wifi_configs,
  additionalProps: Array<[string, string]>,
  handleMACAddressChange?: string => void,
  configChangeHandler: (string, string | number) => void,
  additionalPropsChangeHandler: (Array<[string, string]>) => void,
};

export default function WifiDeviceFields(props: Props) {
  const classes = useStyles();
  const reboot_if_bootid = getAdditionalProp(
    props.additionalProps,
    'reboot_if_bootid',
  );

  const handleRequestReboot = ({target}) => {
    const keyValuePairs = props.additionalProps.slice(0);
    if (target.checked && props.status && props.status.meta) {
      // add the reboot directive
      setAdditionalProp(
        keyValuePairs,
        'reboot_if_bootid',
        props.status.meta.boot_id,
      );
    } else {
      // remove the reboot directive
      setAdditionalProp(keyValuePairs, 'reboot_if_bootid', undefined);
      // if there are no key/values, then add a dummy line for UI purposes
      if (keyValuePairs.length === 0) {
        keyValuePairs.push(['', '']);
      }
    }
    props.additionalPropsChangeHandler(keyValuePairs);
  };

  return (
    <>
      <FormGroup row>
        {props.handleMACAddressChange && (
          <TextField
            required
            className={classes.input}
            label="MAC Address"
            margin="normal"
            onChange={({target}) =>
              nullthrows(props.handleMACAddressChange)(target.value)
            }
            value={props.macAddress}
          />
        )}
        <TextField
          required
          className={classes.input}
          label="Info"
          margin="normal"
          onChange={({target}) =>
            props.configChangeHandler('info', target.value)
          }
          value={props.configs.info}
        />
        <TextField
          className={classes.input}
          label="Latitude"
          margin="normal"
          onChange={({target}) =>
            props.configChangeHandler('latitude', target.value)
          }
          value={props.configs.latitude}
        />
        <TextField
          className={classes.input}
          label="Longitude"
          margin="normal"
          onChange={({target}) =>
            props.configChangeHandler('longitude', target.value)
          }
          value={props.configs.longitude}
        />
        <TextField
          className={classes.input}
          label="Client Channel"
          margin="normal"
          onChange={({target}) =>
            props.configChangeHandler('client_channel', target.value)
          }
          value={props.configs.client_channel}
        />
        <FormControlLabel
          control={
            <Checkbox
              checked={props.configs.is_production}
              onChange={({target}) =>
                props.configChangeHandler('is_production', target.checked)
              }
              color="primary"
            />
          }
          label="Is Production"
        />

        <FormControlLabel
          control={
            <Checkbox
              disabled={props.status === null}
              checked={
                props.status !== null &&
                reboot_if_bootid !== null &&
                reboot_if_bootid === props.status?.meta?.boot_id
              }
              onChange={handleRequestReboot}
              color="primary"
            />
          }
          label="Reboot requested"
        />
      </FormGroup>
      <KeyValueFields
        key_label="key"
        value_label="value"
        keyValuePairs={props.additionalProps || [['', '']]}
        onChange={props.additionalPropsChangeHandler}
      />
    </>
  );
}
