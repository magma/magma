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
import type {policy_qos_profile, policy_rule} from '@fbcnms/magma-api';

type Props = {
  policyRule: policy_rule,
  onChange: policy_rule => void,
  qosProfiles: {[string]: policy_qos_profile},
  descriptionClass: string,
  dialogClass: string,
  inputClass: string,
};

export default function PolicyQosEdit(props: Props) {
  const {policyRule, qosProfiles} = props;

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
            {'Specify quality of service level'}
          </Typography>
          <ListItem dense disableGutters />
          <AltFormField disableGutters label={'QoS Profile'}>
            <Select
              className={props.inputClass}
              fullWidth={true}
              variant={'outlined'}
              value={policyRule?.qos_profile ?? ''}
              onChange={({target}) => {
                props.onChange({...policyRule, qos_profile: target.value});
              }}
              input={<OutlinedInput id="qosProfile" />}>
              {Object.keys(qosProfiles).map(profileID => (
                <MenuItem key={profileID} value={profileID}>
                  <ListItemText primary={profileID} />
                </MenuItem>
              ))}
            </Select>
          </AltFormField>
        </List>
      </DialogContent>
    </>
  );
}
