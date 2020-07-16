/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import type {network_ran_configs} from '@fbcnms/magma-api';

import Grid from '@material-ui/core/Grid';
import ListItem from '@material-ui/core/ListItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';

import {colors} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '50%',
    fullWidth: true,
  },
  itemTitle: {
    color: colors.primary.comet,
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
  },
  itemValue: {
    color: colors.primary.brightGray,
  },
}));

type Props = {
  lteRanConfigs: network_ran_configs,
  setLteRanConfigs: network_ran_configs => void,
};

export default function TddConfig(props: Props) {
  const classes = useStyles();
  return (
    <>
      <ListItem>
        <Grid container>
          <Grid item xs={12}>
            EARFCNDL
          </Grid>
          <Grid item xs={12}>
            <OutlinedInput
              className={classes.input}
              type="number"
              data-testid="earfcndl"
              value={props.lteRanConfigs.tdd_config?.earfcndl}
              onChange={({target}) =>
                props.setLteRanConfigs({
                  ...props.lteRanConfigs,
                  fdd_config: undefined,
                  tdd_config: {
                    special_subframe_pattern:
                      props.lteRanConfigs.tdd_config
                        ?.special_subframe_pattern ?? 0,
                    subframe_assignment:
                      props.lteRanConfigs.tdd_config?.subframe_assignment ?? 0,
                    earfcndl: parseInt(target.value),
                  },
                })
              }
            />
          </Grid>
        </Grid>
      </ListItem>
      <ListItem>
        <Grid container>
          <Grid item xs={12}>
            Special Subframe Pattern
          </Grid>
          <Grid item xs={12}>
            <OutlinedInput
              className={classes.input}
              type="number"
              data-testid="specialSubframePattern"
              value={props.lteRanConfigs.tdd_config?.special_subframe_pattern}
              onChange={({target}) =>
                props.setLteRanConfigs({
                  ...props.lteRanConfigs,
                  fdd_config: undefined,
                  tdd_config: {
                    special_subframe_pattern: parseInt(target.value),
                    subframe_assignment:
                      props.lteRanConfigs.tdd_config?.subframe_assignment ?? 0,
                    earfcndl: props.lteRanConfigs.tdd_config?.earfcndl ?? 0,
                  },
                })
              }
            />
          </Grid>
        </Grid>
      </ListItem>
      <ListItem>
        <Grid container>
          <Grid item xs={12}>
            Subframe Assignment
          </Grid>
          <Grid item xs={12}>
            <OutlinedInput
              className={classes.input}
              type="number"
              data-testid="subframeAssignment"
              value={props.lteRanConfigs.tdd_config?.subframe_assignment}
              onChange={({target}) => {
                props.setLteRanConfigs({
                  ...props.lteRanConfigs,
                  fdd_config: undefined,
                  tdd_config: {
                    subframe_assignment: parseInt(target.value),
                    special_subframe_pattern:
                      props.lteRanConfigs.tdd_config
                        ?.special_subframe_pattern ?? 0,
                    earfcndl: props.lteRanConfigs.tdd_config?.earfcndl ?? 0,
                  },
                });
              }}
            />
          </Grid>
        </Grid>
      </ListItem>
    </>
  );
}
