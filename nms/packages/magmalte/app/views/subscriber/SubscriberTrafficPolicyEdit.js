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
import type {EditSubscriberProps} from './SubscriberAddDialog.js';

import ApnContext from '../../components/context/ApnContext';
import Checkbox from '@material-ui/core/Checkbox';
import FormControl from '@material-ui/core/FormControl';
import List from '@material-ui/core/List';
import ListItemText from '@material-ui/core/ListItemText';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import MenuItem from '@material-ui/core/MenuItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import PolicyContext from '../../components/context/PolicyContext';
import React from 'react';
import Select from '@material-ui/core/Select';
import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';

import {AltFormField} from '../../components/FormField';
import {makeStyles} from '@material-ui/styles';
import {useContext} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';

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
  const {match} = useRouter();
  const apnCtx = useContext(ApnContext);
  const apns = Array.from(Object.keys(apnCtx.state || {}));
  const policyCtx = useContext(PolicyContext);

  const {isLoading: baseNamesLoading, response: baseNames} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdSubscriberConfigBaseNames,
    {
      networkId: nullthrows(match.params.networkId),
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
