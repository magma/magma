/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 *
 * Edit rule labels
 */

import * as React from 'react';
import Button from '@material-ui/core/Button';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import CardHeader from '@material-ui/core/CardHeader';
import DeleteIcon from '@material-ui/icons/Delete';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import SyncIcon from '@material-ui/icons/Sync';
import TextField from '@material-ui/core/TextField';
import Tooltip from '@material-ui/core/Tooltip';
import Typography from '@material-ui/core/Typography';
import type {Labels} from '../AlarmAPIType';

const filteredLabels = new Set(['networkID', 'severity']);

type Props = {
  labels: {[string]: string},
  onChange: (newLabels: Labels) => void,
};

export default function LabelsEditor({labels, onChange}: Props) {
  const [initialState] = React.useState(labels);
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
      const newLabel = [key, value];

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
      updateLabel(index, newKey, labelsState[index][1]);
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

  const resetLabels = React.useCallback(() => {
    updateLabels(convertLabelsToPairs(initialState, filteredLabels));
  }, [updateLabels, initialState]);

  const removeLabel = React.useCallback(
    index => {
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
              Add labels to the fired alert
            </Typography>
          </>
        }
        action={
          <Tooltip title="Reset labels">
            <IconButton
              aria-label="Reset labels"
              color="default"
              onClick={resetLabels}
              size="small">
              <SyncIcon />
            </IconButton>
          </Tooltip>
        }
      />
      <CardContent>
        <Grid container direction="column" spacing={2}>
          {labelsState &&
            labelsState.map(([key, value], index) => (
              <Grid container key={index} item spacing={1}>
                <Grid item xs={6}>
                  <TextField
                    placeholder="Name"
                    value={key}
                    fullWidth
                    onChange={e => handleKeyChange(index, e.target.value)}
                  />
                </Grid>
                <Grid item xs={5}>
                  <TextField
                    placeholder="Value"
                    value={value}
                    fullWidth
                    onChange={e => handleValueChange(index, e.target.value)}
                  />
                </Grid>
                <Grid item xs={1}>
                  <IconButton
                    title="Remove Label"
                    aria-label="Remove Label"
                    onClick={() => removeLabel(index)}>
                    <DeleteIcon />
                  </IconButton>
                </Grid>
              </Grid>
            ))}

          <Grid item>
            <Button onClick={addNewLabel} data-testid="add-new-label">
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
  }, {});
}
