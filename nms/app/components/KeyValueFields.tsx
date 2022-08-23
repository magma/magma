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

import AddCircleOutline from '@mui/icons-material/AddCircleOutline';
import Grid from '@mui/material/Grid';
import IconButton from '@mui/material/IconButton';
import OutlinedInput from '@mui/material/OutlinedInput';
import React from 'react';
import RemoveCircleOutline from '@mui/icons-material/RemoveCircleOutline';
import {AltFormField} from './FormField';
import {makeStyles} from '@mui/styles';

const useStyles = makeStyles({
  container: {
    display: 'block',
    margin: '10px 0 10px 0',
    whiteSpace: 'nowrap',
    width: '100%',
  },
  inputKey: {
    width: '245px',
    paddingRight: '10px',
  },
  inputValue: {
    width: '240px',
  },
  icon: {
    width: '40px',
    height: '40px',
    verticalAlign: 'bottom',
  },
});

type Props = {
  key_label: string;
  value_label: string;
  keyValuePairs: Array<[string, string]>;
  onChange: (keyValue: Array<[string, string]>) => void;
};

export default function KeyValueFields(props: Props) {
  const classes = useStyles();
  const onChange = (index: number, subIndex: number, value: string) => {
    const keyValuePairs = [...props.keyValuePairs];
    keyValuePairs[index] = [keyValuePairs[index][0], keyValuePairs[index][1]];
    keyValuePairs[index][subIndex] = value;
    props.onChange(keyValuePairs);
  };

  const removeField = (index: number) => {
    const keyValuePairs = [...props.keyValuePairs];
    keyValuePairs.splice(index, 1);
    props.onChange(keyValuePairs);
  };

  const addField = () => {
    props.onChange([...props.keyValuePairs, ['', '']]);
  };

  return (
    <AltFormField label={`${props.key_label} - ${props.value_label}`}>
      {props.keyValuePairs.map((pair, index) => (
        <div className={classes.container} key={index}>
          <Grid container spacing={2}>
            <Grid item>
              <OutlinedInput
                placeholder={props.key_label}
                value={pair[0]}
                onChange={({target}) => onChange(index, 0, target.value)}
              />
            </Grid>
            <Grid item>
              <OutlinedInput
                placeholder={props.value_label}
                value={pair[1]}
                onChange={({target}) => onChange(index, 1, target.value)}
              />
            </Grid>
            <Grid item>
              {props.keyValuePairs.length !== 1 && (
                <IconButton
                  onClick={() => removeField(index)}
                  className={classes.icon}
                  size="large">
                  <RemoveCircleOutline />
                </IconButton>
              )}
              {index === props.keyValuePairs.length - 1 && (
                <IconButton
                  onClick={addField}
                  className={classes.icon}
                  size="large">
                  <AddCircleOutline />
                </IconButton>
              )}
            </Grid>
          </Grid>
        </div>
      ))}
    </AltFormField>
  );
}
