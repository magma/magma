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

// $FlowFixMe migrated to typescript
import ApnContext from '../../components/context/ApnContext';
import Checkbox from '@material-ui/core/Checkbox';
import FormControl from '@material-ui/core/FormControl';
import List from '@material-ui/core/List';
import ListItemText from '@material-ui/core/ListItemText';
// $FlowFixMe migrated to typescript
import LoadingFiller from '../../components/LoadingFiller';
import MagmaV1API from '../../../generated/WebClient';
import MenuItem from '@material-ui/core/MenuItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
// $FlowFixMe migrated to typescript
import PolicyContext from '../../components/context/PolicyContext';
import React from 'react';
import Select from '@material-ui/core/Select';
// $FlowFixMe migrated to typescript
import nullthrows from '../../../shared/util/nullthrows';
import useMagmaAPI from '../../../api/useMagmaAPIFlow';

// $FlowFixMe migrated to typescript
import {AltFormField} from '../../components/FormField';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {EditSubscriberProps} from './SubscriberUtils';
import {makeStyles} from '@material-ui/styles';
import {useContext} from 'react';
import {useParams} from 'react-router-dom';

const useStyles = makeStyles(() => ({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '50%',
    fullWidth: true,
  },
}));

export default function EditSubscriberTrafficPolicy(
  props: EditSubscriberProps,
) {
  const classes = useStyles();
  const params = useParams();
  const apnCtx = useContext(ApnContext);
  const apns = Array.from(Object.keys(apnCtx.state || {}));
  const policyCtx = useContext(PolicyContext);

  const {isLoading: baseNamesLoading, response: baseNames} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdSubscriberConfigBaseNames,
    {
      networkId: nullthrows(params.networkId),
    },
  );

  if (baseNamesLoading) {
    return <LoadingFiller />;
  }
  return (
    <div>
      <List>
        <AltFormField label={'Active APNs'}>
          <FormControl className={classes.input}>
            <Select
              multiple
              id="activeApnTestId"
              value={props.subscriberState.active_apns ?? []}
              onChange={({target}) => {
                props.onSubscriberChange('active_apns', target.value);
              }}
              renderValue={selected => selected.join(', ')}
              input={<OutlinedInput />}>
              {apns.map((k: string, idx: number) => (
                <MenuItem key={idx} value={k}>
                  <Checkbox
                    checked={
                      props.subscriberState.active_apns != null
                        ? props.subscriberState.active_apns.indexOf(k) > -1
                        : false
                    }
                  />
                  <ListItemText primary={k} />
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </AltFormField>

        <AltFormField label={'Base Names'}>
          <FormControl className={classes.input}>
            <Select
              multiple
              value={props.subscriberState.active_base_names ?? []}
              onChange={({target}) => {
                props.onSubscriberChange('active_base_names', target.value);
              }}
              renderValue={selected => selected.join(', ')}
              input={<OutlinedInput />}>
              {(baseNames || []).map((k: string, idx: number) => (
                <MenuItem key={idx} value={k}>
                  <Checkbox
                    checked={
                      props.subscriberState.active_base_names != null
                        ? props.subscriberState.active_base_names.indexOf(k) >
                          -1
                        : false
                    }
                  />
                  <ListItemText primary={k} />
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </AltFormField>

        <AltFormField label={'Active Policies'}>
          <FormControl className={classes.input}>
            <Select
              multiple
              value={props.subscriberState.active_policies ?? []}
              onChange={({target}) => {
                props.onSubscriberChange('active_policies', target.value);
              }}
              renderValue={selected => selected.join(', ')}
              input={<OutlinedInput />}>
              {Object.keys(policyCtx.state).map((k: string, idx: number) => (
                <MenuItem key={idx} value={k}>
                  <Checkbox
                    checked={
                      props.subscriberState.active_policies != null
                        ? props.subscriberState.active_policies.indexOf(k) > -1
                        : false
                    }
                  />
                  <ListItemText primary={k} />
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </AltFormField>
      </List>
    </div>
  );
}
