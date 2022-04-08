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

import Autocomplete from '@material-ui/lab/Autocomplete';
import React from 'react';
import TextField from '@material-ui/core/TextField';
import {makeStyles} from '@material-ui/styles';
import type {SelectProps} from './ActionTable';

const useStyles = makeStyles(_ => ({
  inputRoot: {
    '&.MuiOutlinedInput-root': {
      padding: 0,
    },
  },
}));

/**
 * AutoCompleteEdit provides a text input field for editing.
 * Options are provided for autocomplete.
 */
export function AutoCompleteEdit(props: SelectProps) {
  const classes = useStyles();

  return (
    <Autocomplete
      disableClearable
      options={props.content}
      freeSolo
      value={props.value}
      classes={{
        inputRoot: classes.inputRoot,
      }}
      onChange={(_, newValue) => {
        props.onChange(newValue);
      }}
      inputValue={props.value}
      onInputChange={(_, newInputValue) => {
        props.onChange(newInputValue);
      }}
      renderInput={(params: {}) => <TextField {...params} variant="outlined" />}
    />
  );
}
