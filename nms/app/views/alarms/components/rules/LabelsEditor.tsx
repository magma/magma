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
 * Edit rule labels
 */

import * as React from 'react';
import Button from '@mui/material/Button';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import CardHeader from '@mui/material/CardHeader';
import DeleteIcon from '@mui/icons-material/Delete';
import Grid from '@mui/material/Grid';
import IconButton from '@mui/material/IconButton';
import Typography from '@mui/material/Typography';
import {AltFormField} from '../../../../components/FormField';
import {OutlinedInput} from '@mui/material';
import type {Labels} from '../AlarmAPIType';

const filteredLabels = new Set(['networkID', 'severity']);

type Props = {
  labels: Record<string, string>;
  onChange: (newLabels: Labels) => void;
};

export default function LabelsEditor({labels, onChange}: Props) {
  /**
   * Use an array instead of an object because editing an object's key is not
   * possible in this context without causing weird issues.
   */
  const [labelsState, setLabelsState] = React.useState(
    convertLabelsToPairs(labels, filteredLabels),
  );

  // use this instead of using setLabelsState directly
  const updateLabels = React.useCallback(
    (newLabelsState: Array<[string, string]>) => {
      setLabelsState(newLabelsState);
      const newLabels = convertPairsToLabels(newLabelsState);
      onChange(newLabels);
    },
    [onChange, setLabelsState],
  );

  // update a single label by index
  const updateLabel = React.useCallback(
    (index: number, key: string, value: string) => {
      const labelsStateCopy = [...labelsState];
      const newLabel: [string, string] = [key, value];

      if (labelsStateCopy[index]) {
        // edit existing label
        labelsStateCopy[index] = newLabel;
      } else {
        console.error(`no label found at index: ${index}`);
      }
      updateLabels(labelsStateCopy);
    },
    [labelsState, updateLabels],
  );

  const handleKeyChange = React.useCallback(
    (index: number, newKey: string) => {
      updateLabel(index, newKey.replace(/\s/g, '_'), labelsState[index][1]);
    },
    [labelsState, updateLabel],
  );

  const handleValueChange = React.useCallback(
    (index: number, value: string) => {
      updateLabel(index, labelsState[index][0], value);
    },
    [labelsState, updateLabel],
  );

  const addNewLabel = React.useCallback(() => {
    updateLabels(labelsState.concat([['', '']]));
  }, [updateLabels, labelsState]);

  const removeLabel = React.useCallback(
    (index: number) => {
      updateLabels([
        ...labelsState.slice(0, index - 1),
        ...labelsState.slice(index + 1, labelsState.length),
      ]);
    },
    [labelsState, updateLabels],
  );

  return (
    <Card>
      <CardHeader
        title={
          <>
            <Typography variant="h5" gutterBottom>
              Labels
            </Typography>
            <Typography color="textSecondary" gutterBottom variant="body2">
              Add labels to attach data to this alert
            </Typography>
          </>
        }
      />
      <CardContent>
        <Grid container direction="column" spacing={2}>
          {labelsState &&
            labelsState.map(([key, value], index) => (
              <Grid container key={index} item spacing={1}>
                <Grid item xs={6}>
                  <AltFormField disableGutters label="Label Name">
                    <OutlinedInput
                      fullWidth={true}
                      name="description"
                      value={key}
                      id="label-name-input"
                      placeholder="Name"
                      onChange={e => handleKeyChange(index, e.target.value)}
                    />
                  </AltFormField>
                </Grid>
                <Grid item xs={5}>
                  <AltFormField disableGutters label="Value">
                    <OutlinedInput
                      fullWidth={true}
                      name="description"
                      value={value}
                      id="label-value-input"
                      placeholder="Value"
                      onChange={e => handleValueChange(index, e.target.value)}
                    />
                  </AltFormField>
                </Grid>
                <Grid item xs={1}>
                  <IconButton
                    title="Remove Label"
                    aria-label="Remove Label"
                    onClick={() => removeLabel(index)}
                    size="large">
                    <DeleteIcon />
                  </IconButton>
                </Grid>
              </Grid>
            ))}
          <Grid item>
            <Button
              variant="outlined"
              color="primary"
              size="small"
              onClick={addNewLabel}
              data-testid="add-new-label">
              Add new label
            </Button>
          </Grid>
        </Grid>
      </CardContent>
    </Card>
  );
}

// converts Labels to an array like [[key,value], [key,value]]
function convertLabelsToPairs(
  labels: Labels,
  filter: Set<string>,
): Array<[string, string]> {
  return Object.keys(labels)
    .filter(key => !filter.has(key))
    .map(key => [key, labels[key]]);
}

// converts n array like [[key,value], [key,value]] to Labels
function convertPairsToLabels(pairs: Array<[string, string]>): Labels {
  return pairs.reduce((map, [key, val]) => {
    if (key && key.trim() !== '') {
      map[key] = val;
    }
    return map;
  }, {} as Labels);
}
