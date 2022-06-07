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
 */

import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import MenuItem from '@material-ui/core/MenuItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';
import Select from '@material-ui/core/Select';
import Switch from '@material-ui/core/Switch';
import Text from '../../theme/design-system/Text';
import {AltFormField} from '../../components/FormField';
import {PolicyRule} from '../../../generated-ts';
import {makeStyles} from '@material-ui/styles';
import {policyStyles} from './PolicyStyles';

const useStyles = makeStyles(() => policyStyles);

type Props = {
  policyRule: PolicyRule;
  onChange: (arg0: PolicyRule) => void;
  inputClass: string;
};

export default function PolicyRedirectEdit(props: Props) {
  const classes = useStyles();
  const {policyRule} = props;
  // eslint-disable-next-line no-warning-comments
  // $FlowFixMe PolicyRule type will be updated to include field
  const redInfo = policyRule?.redirect || {
    server_address: '',
    address_type: 'IPv4',
    support: 'DISABLED',
  };

  const handleFieldChange = (field: string, value: number | string) => {
    props.onChange({
      ...policyRule,
      redirect: {...redInfo, [field]: value},
    });
  };

  return (
    <div data-testid="redirectEdit">
      <Text weight="medium" variant="subtitle2" className={classes.description}>
        {
          'If redirection is enabled, matching traffic can be redirected to a captive portal server'
        }
      </Text>
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
          fullWidth={true}
          className={props.inputClass}
          variant={'outlined'}
          // eslint-disable-next-line no-warning-comments
          // $FlowFixMe redirect_info type needed
          value={redInfo.address_type || 'IPv4'}
          onChange={({target}) => {
            handleFieldChange('address_type', target.value as string);
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
      <AltFormField disableGutters label={'Support'} isOptional>
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
      </AltFormField>
    </div>
  );
}
