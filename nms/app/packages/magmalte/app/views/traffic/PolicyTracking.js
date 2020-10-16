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
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import MenuItem from '@material-ui/core/MenuItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';
import Select from '@material-ui/core/Select';
import Typography from '@material-ui/core/Typography';

import {AltFormField} from '../../components/FormField';
import {base64ToHex, decodeBase64} from '@fbcnms/util/strings';
import type {policy_rule} from '@fbcnms/magma-api';

type Props = {
  policyRule: policy_rule,
  onChange: policy_rule => void,
  descriptionClass: string,
  dialogClass: string,
  inputClass: string,
};

export default function PolicyTrackingEdit(props: Props) {
  return (
    <>
      <DialogContent
        data-testid="networkTrackingEdit"
        className={props.dialogClass}>
        <List>
          <Typography
            variant="caption"
            display="block"
            className={props.descriptionClass}
            gutterBottom>
            {'Tracking configuration for the policy'}
          </Typography>
          <ListItem dense disableGutters />
          <AltFormField disableGutters label={'Monitoring Key (base64)'}>
            <OutlinedInput
              required={true}
              className={props.inputClass}
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
            <Select
              fullWidth={true}
              className={props.inputClass}
              variant={'outlined'}
              value={props.policyRule.tracking_type || 'NO_TRACKING'}
              onChange={({target}) =>
                props.onChange({
                  ...props.policyRule,
                  // $FlowIgnore: value guaranteed to match the string literals
                  tracking_type: target.value,
                })
              }
              input={<OutlinedInput id="trackingType" />}>
              <MenuItem value={'ONLY_OCS'}>
                <ListItemText primary={'Only OCS'} />
              </MenuItem>
              <MenuItem value={'ONLY_PCRF'}>
                <ListItemText primary={'Only PCRF'} />
              </MenuItem>
              <MenuItem value={'OCS_AND_PCRF'}>
                <ListItemText primary={'OCS and PCRF'} />
              </MenuItem>
              <MenuItem value={'NO_TRACKING'}>
                <ListItemText primary={'No Tracking'} />
              </MenuItem>
            </Select>
          </AltFormField>
        </List>
      </DialogContent>
    </>
  );
}
