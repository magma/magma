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
import Text from '../../theme/design-system/Text';
import {AltFormField} from '../../components/FormField';
import {makeStyles} from '@material-ui/styles';
import {policyStyles} from './PolicyStyles';
import type {
  PolicyRule,
  PolicyRuleAppNameEnum,
  PolicyRuleAppServiceTypeEnum,
} from '../../../generated-ts';

const useStyles = makeStyles(() => policyStyles);

type Props = {
  policyRule: PolicyRule;
  onChange: (arg0: PolicyRule) => void;
};

export default function PolicyAppEdit(props: Props) {
  const {policyRule} = props;
  const classes = useStyles();
  const appList = [
    'NO_APP_NAME',
    'FACEBOOK',
    'FACEBOOK_MESSENGER',
    'INSTAGRAM',
    'YOUTUBE',
    'GOOGLE',
    'GMAIL',
    'GOOGLE_DOCS',
    'NETFLIX',
    'APPLE',
    'MICROSOFT',
    'REDDIT',
    'WHATSAPP',
    'GOOGLE_PLAY',
    'APPSTORE',
    'AMAZON',
    'WECHAT',
    'TIKTOK',
    'TWITTER',
    'WIKIPEDIA',
    'GOOGLE_MAPS',
    'YAHOO',
    'IMO',
  ];
  const serviceTypes = ['NO_SERVICE_TYPE', 'CHAT', 'AUDIO', 'VIDEO'];
  return (
    <div data-testid="appEdit">
      <Text weight="medium" variant="subtitle2" className={classes.description}>
        {'App description and help'}
      </Text>
      <ListItem dense disableGutters />
      <AltFormField disableGutters label={'App Name'}>
        <Select
          fullWidth={true}
          variant={'outlined'}
          value={policyRule.app_name || 'NO_APP_NAME'}
          onChange={({target}) => {
            props.onChange({
              ...policyRule,
              app_name: target.value as PolicyRuleAppNameEnum,
            });
          }}
          input={<OutlinedInput id="appName" />}>
          {appList.map(appName => (
            <MenuItem value={appName}>
              <ListItemText primary={appName} />
            </MenuItem>
          ))}
        </Select>
      </AltFormField>
      <AltFormField disableGutters label={'App Service Type'}>
        <Select
          fullWidth={true}
          variant={'outlined'}
          value={policyRule.app_service_type || 'NO_SERVICE_TYPE'}
          onChange={({target}) => {
            props.onChange({
              ...policyRule,
              app_service_type: target.value as PolicyRuleAppServiceTypeEnum,
            });
          }}
          input={<OutlinedInput id="appServiceType" />}>
          {serviceTypes.map(serviceType => (
            <MenuItem value={serviceType}>
              <ListItemText primary={serviceType} />
            </MenuItem>
          ))}
        </Select>
      </AltFormField>
    </div>
  );
}
