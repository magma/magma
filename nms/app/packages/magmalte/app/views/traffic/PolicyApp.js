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
import type {policy_rule} from '@fbcnms/magma-api';

type Props = {
  policyRule: policy_rule,
  onChange: policy_rule => void,
  descriptionClass: string,
  dialogClass: string,
  inputClass: string,
};

export default function PolicyAppEdit(props: Props) {
  const {policyRule} = props;

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
            {'App description and help'}
          </Typography>
          <ListItem dense disableGutters />
          <AltFormField disableGutters label={'App Name'}>
            <Select
              className={props.inputClass}
              fullWidth={true}
              variant={'outlined'}
              value={policyRule.app_name || 'NO_APP_NAME'}
              onChange={({target}) => {
                props.onChange({
                  ...policyRule,
                  // $FlowIgnore: value guaranteed to match the string literals
                  app_name: target.value,
                });
              }}
              input={<OutlinedInput id="appName" />}>
              <MenuItem value={'NO_APP_NAME'}>
                <ListItemText primary={'No App Name'} />
              </MenuItem>
              <MenuItem value={'FACEBOOK'}>
                <ListItemText primary={'Facebook'} />
              </MenuItem>
              <MenuItem value={'FACEBOOK_MESSENGER'}>
                <ListItemText primary={'Facebook Messenger'} />
              </MenuItem>
              <MenuItem value={'INSTAGRAM'}>
                <ListItemText primary={'Instagram'} />
              </MenuItem>
              <MenuItem value={'YOUTUBE'}>
                <ListItemText primary={'Youtube'} />
              </MenuItem>
              <MenuItem value={'GOOGLE'}>
                <ListItemText primary={'Google'} />
              </MenuItem>
              <MenuItem value={'GMAIL'}>
                <ListItemText primary={'Gmail'} />
              </MenuItem>
              <MenuItem value={'GOOGLE_DOCS'}>
                <ListItemText primary={'Google Docs'} />
              </MenuItem>
              <MenuItem value={'NETFLIX'}>
                <ListItemText primary={'Netflix'} />
              </MenuItem>
              <MenuItem value={'APPLE'}>
                <ListItemText primary={'Apple'} />
              </MenuItem>
              <MenuItem value={'MICROSOFT'}>
                <ListItemText primary={'Microsoft'} />
              </MenuItem>
              <MenuItem value={'REDDIT'}>
                <ListItemText primary={'Reddit'} />
              </MenuItem>
              <MenuItem value={'WHATSAPP'}>
                <ListItemText primary={'WhatsApp'} />
              </MenuItem>
              <MenuItem value={'GOOGLE_PLAY'}>
                <ListItemText primary={'Google Play'} />
              </MenuItem>
              <MenuItem value={'APPSTORE'}>
                <ListItemText primary={'App Store'} />
              </MenuItem>
              <MenuItem value={'AMAZON'}>
                <ListItemText primary={'Amazon'} />
              </MenuItem>
              <MenuItem value={'WECHAT'}>
                <ListItemText primary={'Wechat'} />
              </MenuItem>
              <MenuItem value={'TIKTOK'}>
                <ListItemText primary={'TikTok'} />
              </MenuItem>
              <MenuItem value={'TWITTER'}>
                <ListItemText primary={'Twitter'} />
              </MenuItem>
              <MenuItem value={'WIKIPEDIA'}>
                <ListItemText primary={'Wikipedia'} />
              </MenuItem>
              <MenuItem value={'GOOGLE_MAPS'}>
                <ListItemText primary={'Google Maps'} />
              </MenuItem>
              <MenuItem value={'YAHOO'}>
                <ListItemText primary={'Yahoo'} />
              </MenuItem>
              <MenuItem value={'IMO'}>
                <ListItemText primary={'IMO'} />
              </MenuItem>
            </Select>
          </AltFormField>
          <AltFormField disableGutters label={'App Service Type'}>
            <Select
              className={props.inputClass}
              fullWidth={true}
              variant={'outlined'}
              value={policyRule.app_service_type || 'NO_SERVICE_TYPE'}
              onChange={({target}) => {
                props.onChange({
                  ...policyRule,
                  // $FlowIgnore: value guaranteed to match the string literals
                  app_service_type: target.value,
                });
              }}
              input={<OutlinedInput id="appServiceType" />}>
              <MenuItem value={'NO_SERVICE_TYPE'}>
                <ListItemText primary={'No Service Type'} />
              </MenuItem>
              <MenuItem value={'CHAT'}>
                <ListItemText primary={'Chat'} />
              </MenuItem>
              <MenuItem value={'AUDIO'}>
                <ListItemText primary={'Audio'} />
              </MenuItem>
              <MenuItem value={'VIDEO'}>
                <ListItemText primary={'Video'} />
              </MenuItem>
            </Select>
          </AltFormField>
        </List>
      </DialogContent>
    </>
  );
}
