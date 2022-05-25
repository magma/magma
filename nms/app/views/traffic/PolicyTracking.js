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
import type {policy_rule} from '../../../generated/MagmaAPIBindings';

import ListItem from '@material-ui/core/ListItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Text from '../../theme/design-system/Text';
import TypedSelect from '../../components/TypedSelect';

import {AltFormField} from '../../components/FormField';
// $FlowFixMe[cannot-resolve-module]
import {base64ToHex, decodeBase64} from '../../util/strings';
import {makeStyles} from '@material-ui/styles';
import {policyStyles} from './PolicyStyles';

const useStyles = makeStyles(() => policyStyles);

type Props = {
  policyRule: policy_rule,
  onChange: policy_rule => void,
  inputClass: string,
};

export default function PolicyTrackingEdit(props: Props) {
  const classes = useStyles();

  return (
    <div data-testid="trackingEdit">
      <Text weight="medium" variant="subtitle2" className={classes.description}>
        {'Tracking configuration for the policy'}
      </Text>
      <ListItem dense disableGutters />
      <AltFormField disableGutters label={'Monitoring Key (base64)'}>
        <OutlinedInput
          className={props.inputClass}
          required={true}
          data-testid="monitoringKey"
          placeholder="Enter Monitoring Key"
          fullWidth={true}
          value={props.policyRule.monitoring_key ?? ''}
          onChange={({target}) => {
            props.onChange({
              ...props.policyRule,
              monitoring_key: target.value,
            });
          }}
        />
      </AltFormField>
      <AltFormField disableGutters label={'Monitoring Key (hex)'}>
        <OutlinedInput
          className={props.inputClass}
          data-testid="monitoringKey"
          placeholder="Enter Monitoring Key"
          fullWidth={true}
          disabled={true}
          value={base64ToHex(props.policyRule.monitoring_key ?? '')}
        />
      </AltFormField>
      <AltFormField disableGutters label={'Monitoring Key (ascii)'}>
        <OutlinedInput
          className={props.inputClass}
          data-testid="monitoringKey"
          placeholder="Enter Monitoring Key"
          fullWidth={true}
          disabled={true}
          value={decodeBase64(props.policyRule.monitoring_key ?? '')}
        />
      </AltFormField>
      <AltFormField disableGutters label={'Rating Group'}>
        <OutlinedInput
          required={true}
          className={props.inputClass}
          data-testid="ratingGroup"
          placeholder="0"
          fullWidth={true}
          value={props.policyRule.rating_group}
          type={'number'}
          onChange={({target}) =>
            props.onChange({
              ...props.policyRule,
              rating_group: parseInt(target.value),
            })
          }
        />
      </AltFormField>
      <AltFormField disableGutters label={'Tracking Type'}>
        <TypedSelect
          className={props.inputClass}
          input={<OutlinedInput />}
          value={props.policyRule?.tracking_type ?? 'NO_TRACKING'}
          items={{
            NO_TRACKING: 'No Tracking',
            ONLY_OCS: 'Only OCS',
            ONLY_PCRF: 'Only PCRF',
            OCS_AND_PCRF: 'OCS and PCRF',
          }}
          onChange={value => {
            props.onChange({
              ...props.policyRule,
              // $FlowIgnore: value guaranteed to match the string literals
              tracking_type: value,
            });
          }}
        />
      </AltFormField>
    </div>
  );
}
