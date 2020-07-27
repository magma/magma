/**
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

import * as React from 'react';
import CircularProgress from '@material-ui/core/CircularProgress';
import MenuItem from '@material-ui/core/MenuItem';
import TextField from '@material-ui/core/TextField';
import useRouter from '../../../hooks/useRouter';
import {useAlarmContext} from '../../AlarmContext';

type Props = {
  onChange: (receiverName: string) => void,
  receiver: ?string,
};

export default function SelectReceiver({
  onChange,
  receiver,
  ...fieldProps
}: Props) {
  const {apiUtil} = useAlarmContext();
  const {match} = useRouter();
  const {isLoading, error, response} = apiUtil.useAlarmsApi(
    apiUtil.getReceivers,
    {
      networkId: match.params.networkId,
    },
  );
  const handleChange = React.useCallback(
    (e: SyntheticInputEvent<HTMLInputElement>) => {
      onChange(e.target.value);
    },
    [onChange],
  );

  if (isLoading) {
    return <CircularProgress size={20} />;
  }

  return (
    <TextField
      {...fieldProps}
      select
      id="select-receiver"
      onChange={handleChange}
      inputProps={{'data-testid': 'select-receiver-input'}}
      value={receiver || ''}>
      <MenuItem value="" key={''}>
        None
      </MenuItem>
      {error && <MenuItem>Error: Could not load receivers</MenuItem>}
      {(response || []).map(receiver => (
        <MenuItem value={receiver.name} key={receiver.name}>
          {receiver.name}
        </MenuItem>
      ))}
    </TextField>
  );
}
