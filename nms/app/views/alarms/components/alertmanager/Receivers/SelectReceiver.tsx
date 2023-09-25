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
 */

import * as React from 'react';
import Chip from '@mui/material/Chip';
import CircularProgress from '@mui/material/CircularProgress';
import MenuItem from '@mui/material/MenuItem';
import Select, {SelectChangeEvent} from '@mui/material/Select';
import {AltFormField} from '../../../../../components/FormField';
import {FormControl, OutlinedInput} from '@mui/material';
import {useAlarmContext} from '../../AlarmContext';
import {useParams} from 'react-router-dom';

type Props = {
  onChange: (receiverName: string) => void;
  receiver: string | null | undefined;
  fullWidth: boolean;
};

export default function SelectReceiver({
  onChange,
  receiver,
  ...fieldProps
}: Props) {
  const {apiUtil} = useAlarmContext();
  const params = useParams();
  const {isLoading, error, response} = apiUtil.useAlarmsApi(
    apiUtil.getReceivers,
    {
      networkId: params.networkId!,
    },
  );
  const handleChange = React.useCallback(
    (e: SelectChangeEvent<string>) => {
      onChange(e.target.value);
    },
    [onChange],
  );

  if (isLoading) {
    return <CircularProgress size={20} />;
  }

  return (
    <AltFormField disableGutters label="Audience">
      <FormControl fullWidth>
        <Select
          {...fieldProps}
          id="select-receiver"
          data-testid="select-receiver"
          onChange={handleChange}
          defaultValue="Select Team"
          input={
            <OutlinedInput
              inputProps={{'data-testid': 'select-receiver-input'}}
            />
          }
          renderValue={value => (
            <Chip
              key={value}
              label={value}
              variant="outlined"
              color="primary"
              size="small"
            />
          )}
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
        </Select>
      </FormControl>
    </AltFormField>
  );
}
