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

import FormControl from '@material-ui/core/FormControl';
import FormHelperText from '@material-ui/core/FormHelperText';
import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Select from '@material-ui/core/Select';

import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  formControl: {
    margin: theme.spacing(),
    minWidth: 120,
    width: 'calc(100% - 15px)',
  },
  selectEmpty: {
    marginTop: theme.spacing(2),
  },
}));

type Props = {
  onChange: (meshId: string) => void,
  meshes: Array<string>,
  selectedMeshID: string,
  disallowEmpty?: boolean,
  helperText?: string,
};

export default function WifiSelectMesh(props: Props) {
  const classes = useStyles();
  if (!props.meshes) {
    return null;
  }

  const meshes = [...props.meshes];
  meshes.sort((a, b) => (a.toLowerCase() > b.toLowerCase() ? 1 : -1));
  const meshItems = meshes.map(meshId => (
    <MenuItem value={meshId} key={meshId}>
      {meshId}
    </MenuItem>
  ));

  return (
    <FormControl className={classes.formControl}>
      <InputLabel htmlFor="meshid-helper">Mesh ID</InputLabel>
      <Select
        value={props.selectedMeshID}
        onChange={event => props.onChange(event.target.value)}
        input={<Input name="meshId" id="meshid-helper" />}>
        {props.disallowEmpty !== true && (
          <MenuItem value="">
            <em>All</em>
          </MenuItem>
        )}
        {meshItems}
      </Select>
      {props.helperText != null && (
        <FormHelperText>{props.helperText}</FormHelperText>
      )}
    </FormControl>
  );
}
