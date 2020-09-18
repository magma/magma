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

import type {network_wifi_configs} from '@fbcnms/magma-api';

import Button from '@fbcnms/ui/components/design-system/Button';
import FormGroup from '@material-ui/core/FormGroup';
import KeyValueFields from '@fbcnms/magmalte/app/components/KeyValueFields';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React from 'react';
import TextField from '@material-ui/core/TextField';

import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';
import {
  additionalPropsToArray,
  additionalPropsToObject,
} from '../wifi/WifiUtils';
import {get} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  block: {
    display: 'block',
    marginRight: theme.spacing(),
    width: '245px',
  },
  formContainer: {
    paddingBottom: theme.spacing(2),
  },
  formGroup: {
    marginLeft: theme.spacing(2),
    marginBottom: theme.spacing(2),
  },
  keyValueFieldsInputValue: {
    width: '585px',
  },
  saveButton: {
    marginTop: theme.spacing(2),
  },
  textArea: {
    width: '600px',
  },
  textField: {
    marginRight: theme.spacing(),
    width: '245px',
  },
}));

export default function WifiNetworkConfig() {
  const classes = useStyles();
  const {match} = useRouter();
  const [config, setConfig] = useState<?network_wifi_configs>();
  const [additionalProps, setAdditionalProps] = useState<?Array<
    [string, string],
  >>();
  const enqueueSnackbar = useEnqueueSnackbar();

  const networkId = nullthrows(match.params.networkId);
  const {isLoading} = useMagmaAPI(
    MagmaV1API.getWifiByNetworkIdWifi,
    {networkId},
    useCallback(response => {
      setConfig({...response});
      setAdditionalProps(
        additionalPropsToArray(response.additional_props) || [],
      );
    }, []),
  );

  if (isLoading || !config) {
    return <LoadingFiller />;
  }

  const handleSave = () => {
    MagmaV1API.putWifiByNetworkIdWifi({
      networkId,
      config: {
        ...config,
        ping_num_packets: parseInt(config.ping_num_packets),
        ping_timeout_secs: parseInt(config.ping_timeout_secs),
        xwf_radius_auth_port: parseInt(config.xwf_radius_auth_port),
        xwf_radius_acct_port: parseInt(config.xwf_radius_acct_port),
        additional_props: additionalPropsToObject(additionalProps) || {},
      },
    })
      .then(() => enqueueSnackbar('Saved successfully', {variant: 'success'}))
      .catch(error =>
        enqueueSnackbar(get(error, 'response.data.message', error), {
          variant: 'error',
        }),
      );
  };

  const onConfigChangeHandler = (fieldName: string) => ({target}) =>
    setConfig({...config, [fieldName]: target.value});

  return (
    <div className={classes.formContainer}>
      <FormGroup row className={classes.formGroup}>
        <TextField
          required
          label="VPN Proto"
          margin="normal"
          className={classes.textField}
          value={config.mgmt_vpn_proto}
          onChange={onConfigChangeHandler('mgmt_vpn_proto')}
        />
        <TextField
          required
          label="VPN Remote"
          margin="normal"
          className={classes.textField}
          value={config.mgmt_vpn_remote}
          onChange={onConfigChangeHandler('mgmt_vpn_remote')}
        />
      </FormGroup>
      <FormGroup row className={classes.formGroup}>
        <TextField
          required
          label="Ping Host List"
          margin="normal"
          className={classes.textField}
          value={(config.ping_host_list || []).join(',')}
          onChange={({target}) =>
            setConfig({...config, ping_host_list: target.value.split(',')})
          }
        />
        <TextField
          required
          label="Ping Number of Packets"
          margin="normal"
          className={classes.textField}
          value={config.ping_num_packets}
          onChange={onConfigChangeHandler('ping_num_packets')}
        />
        <TextField
          required
          label="Ping Timeout (s)"
          margin="normal"
          className={classes.textField}
          value={config.ping_timeout_secs}
          onChange={onConfigChangeHandler('ping_timeout_secs')}
        />
      </FormGroup>
      <FormGroup row className={classes.formGroup}>
        <div>
          <TextField
            required
            label="XWF Radius Server"
            margin="normal"
            className={classes.block}
            value={config.xwf_radius_server}
            onChange={onConfigChangeHandler('xwf_radius_server')}
          />
          <TextField
            required
            label="XWF Radius Auth Port"
            margin="normal"
            className={classes.block}
            value={config.xwf_radius_auth_port}
            onChange={onConfigChangeHandler('xwf_radius_auth_port')}
          />
          <TextField
            required
            label="XWF Radius Acct Port"
            margin="normal"
            className={classes.block}
            value={config.xwf_radius_acct_port}
            onChange={onConfigChangeHandler('xwf_radius_acct_port')}
          />
        </div>
        <div>
          <TextField
            required
            label="XWF Radius Shared Secret"
            margin="normal"
            className={classes.block}
            value={config.xwf_radius_shared_secret}
            onChange={onConfigChangeHandler('xwf_radius_shared_secret')}
          />
          <TextField
            required
            label="XWF UAM Secret"
            margin="normal"
            className={classes.block}
            value={config.xwf_uam_secret}
            onChange={onConfigChangeHandler('xwf_uam_secret')}
          />
        </div>
      </FormGroup>
      <FormGroup row className={classes.formGroup}>
        <div>
          <TextField
            required
            label="XWF Partner Name"
            margin="normal"
            className={classes.block}
            value={config.xwf_partner_name}
            onChange={onConfigChangeHandler('xwf_partner_name')}
          />
          <TextField
            required
            label="XWF DHCP DNS 1"
            margin="normal"
            className={classes.block}
            value={config.xwf_dhcp_dns1}
            onChange={onConfigChangeHandler('xwf_dhcp_dns1')}
          />
          <TextField
            required
            label="XWF DHCP DNS 2"
            margin="normal"
            className={classes.block}
            value={config.xwf_dhcp_dns2}
            onChange={onConfigChangeHandler('xwf_dhcp_dns2')}
          />
        </div>
        <TextField
          multiline
          rowsMax="8"
          label="XWF Config"
          margin="normal"
          className={classes.textArea}
          value={config.xwf_config}
          onChange={onConfigChangeHandler('xwf_config')}
        />
      </FormGroup>
      <FormGroup className={classes.formGroup}>
        <KeyValueFields
          key_label="key"
          value_label="value"
          keyValuePairs={additionalProps || [['', '']]}
          onChange={setAdditionalProps}
          classes={{inputValue: classes.keyValueFieldsInputValue}}
        />
      </FormGroup>
      <FormGroup row className={classes.formGroup}>
        <Button className={classes.saveButton} onClick={handleSave}>
          Save
        </Button>
      </FormGroup>
    </div>
  );
}
