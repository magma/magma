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
import {makeStyles} from '@material-ui/styles';
import {policyStyles} from './PolicyStyles';
import type {PolicyQosProfile, PolicyRule} from '../../../generated-ts';

const useStyles = makeStyles(() => policyStyles);

type Props = {
  policyRule: PolicyRule;
  qosProfiles: Record<string, PolicyQosProfile>;
  onChange: (policyRule: PolicyRule) => void;
  isNetworkWide: boolean;
  setIsNetworkWide: (isNetworkWide: boolean) => void;
  inputClass: string;
};

export default function PolicyInfoEdit(props: Props) {
  const classes = useStyles();
  const {qosProfiles, policyRule} = props;
  return (
    <div data-testid="infoEdit">
      <Text weight="medium" variant="subtitle2" className={classes.description}>
        {'Basic policy rule fields'}
      </Text>
      <ListItem dense disableGutters />
      <AltFormField
        label={'Policy ID'}
        subLabel={'A unique identifier for the policy rule'}
        disableGutters>
        <OutlinedInput
          className={props.inputClass}
          fullWidth={true}
          data-testid="policyID"
          placeholder="Eg. policy_id"
          value={policyRule.id}
          onChange={({target}) => {
            props.onChange({...policyRule, id: target.value});
          }}
        />
      </AltFormField>
      <AltFormField
        label={'Priority Level'}
        subLabel={'Higher priority policies override lower priority ones'}
        disableGutters>
        <OutlinedInput
          type="number"
          className={props.inputClass}
          fullWidth={true}
          data-testid="policyPriority"
          placeholder="Value between 1 and 15"
          value={policyRule.priority}
          onChange={({target}) =>
            props.onChange({...policyRule, priority: parseInt(target.value)})
          }
        />
      </AltFormField>
      <AltFormField label={'Network Wide'} disableGutters>
        <Switch
          data-testid="networkWide"
          onChange={({target}) => props.setIsNetworkWide(target.checked)}
          checked={props.isNetworkWide}
        />
      </AltFormField>
      <AltFormField disableGutters label={'QoS Profile'}>
        <Select
          className={props.inputClass}
          fullWidth={true}
          variant={'outlined'}
          value={policyRule?.qos_profile ?? ''}
          onChange={({target}) => {
            props.onChange({
              ...policyRule,
              qos_profile: target.value as string,
            });
          }}
          input={<OutlinedInput id="qosProfile" />}>
          {Object.keys(qosProfiles).map(profileID => (
            <MenuItem key={profileID} value={profileID}>
              <ListItemText primary={profileID} />
            </MenuItem>
          ))}
        </Select>
      </AltFormField>
    </div>
  );
}
