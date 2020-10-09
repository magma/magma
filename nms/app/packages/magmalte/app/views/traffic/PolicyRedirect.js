/*
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

import DialogContent from '@material-ui/core/DialogContent';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Grid from '@material-ui/core/Grid';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import MenuItem from '@material-ui/core/MenuItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';
import Select from '@material-ui/core/Select';
import Switch from '@material-ui/core/Switch';
import Typography from '@material-ui/core/Typography';

import {AltFormField} from '../../components/FormField';
import {makeStyles} from '@material-ui/styles';
import type {policy_rule} from '@fbcnms/magma-api';

const useStyles = makeStyles(() => ({
  title: {textAlign: 'center', margin: 'auto', marginLeft: '0px'},
  switch: {margin: 'auto 0px'},
}));

type Props = {
  policyRule: policy_rule,
  onChange: policy_rule => void,
  descriptionClass: string,
  dialogClass: string,
  inputClass: string,
};

export default function PolicyRedirectEdit(props: Props) {
  const classes = useStyles();
  const {policyRule} = props;
  // eslint-disable-next-line no-warning-comments
  // $FlowFixMe policy_rule type will be updated to include field
  const redInfo = policyRule?.redirect_information || {
    server_address: '',
    address_type: 'IPv4',
    support: 'DISABLED',
  };

  const handleFieldChange = (field: string, value: number | string) => {
    props.onChange({
      ...policyRule,
      redirect_information: {
        ...redInfo,
        [field]: value,
      },
    });
  };

  return (
    <>
      <DialogContent
        data-testid="networkRedirectEdit"
        className={props.dialogClass}>
        <List>
          <Typography
            variant="caption"
            display="block"
            className={props.descriptionClass}
            gutterBottom>
            {
              'If redirection is enabled, matching traffic can be redirected to a captive portal server'
            }
          </Typography>
          <ListItem dense disableGutters />
          <AltFormField disableGutters label={'Server Address'}>
            <OutlinedInput
              className={props.inputClass}
              required={true}
              data-testid="serverAddress"
              placeholder="Ex. 172.16.254.1 "
              fullWidth={true}
              // eslint-disable-next-line no-warning-comments
              // $FlowFixMe redirect_info type needed
              value={redInfo.server_address ?? ''}
              onChange={({target}) => {
                handleFieldChange('server_address', target.value);
              }}
            />
          </AltFormField>
          <AltFormField disableGutters label={'Address Type'}>
            <Select
              className={props.inputClass}
              fullWidth={true}
              variant={'outlined'}
              // eslint-disable-next-line no-warning-comments
              // $FlowFixMe redirect_info type needed
              value={redInfo.address_type || 'IPv4'}
              onChange={({target}) => {
                handleFieldChange('address_type', target.value);
              }}
              input={<OutlinedInput id="addressType" />}>
              <MenuItem value={'IPv4'}>
                <ListItemText primary={'IPv4'} />
              </MenuItem>
              <MenuItem value={'IPv6'}>
                <ListItemText primary={'IPv6'} />
              </MenuItem>
              <MenuItem value={'URL'}>
                <ListItemText primary={'URL'} />
              </MenuItem>
              <MenuItem value={'SIP URI'}>
                <ListItemText primary={'SIP URI'} />
              </MenuItem>
            </Select>
          </AltFormField>
          <Grid container justify="space-between" className={props.inputClass}>
            <Grid item className={classes.title}>
              <AltFormField disableGutters label={'Support'} isOptional />
            </Grid>
            <Grid item className={classes.switch}>
              <FormControlLabel
                control={
                  <Switch
                    color="primary"
                    // eslint-disable-next-line no-warning-comments
                    // $FlowFixMe redirect_info type needed
                    checked={redInfo.support === 'ENABLED'}
                    onChange={({target}) => {
                      handleFieldChange(
                        'support',
                        target.checked ? 'ENABLED' : 'DISABLED',
                      );
                    }}
                  />
                }
                label={redInfo.support === 'ENABLED' ? 'Enabled' : 'Disabled'}
                labelPlacement="start"
              />
            </Grid>
          </Grid>
        </List>
      </DialogContent>
    </>
  );
}
