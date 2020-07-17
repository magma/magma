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

export default function FddConfig(props: Props) {
  const classes = useStyles();

  return (
    <ListItem>
      <Grid container>
        <Grid item xs={6}>
          <Grid container>
            <Grid item xs={12}>
              EARFCNDL
            </Grid>
            <Grid item xs={12}>
              <OutlinedInput
                className={classes.input}
                data-testid="earfcndl"
                type="number"
                value={props.lteRanConfigs.fdd_config?.earfcndl}
                onChange={({target}) =>
                  props.setLteRanConfigs({
                    ...props.lteRanConfigs,
                    tdd_config: undefined,
                    fdd_config: {
                      earfcndl: parseInt(target.value),
                      earfcnul: props.lteRanConfigs.fdd_config?.earfcnul ?? 0,
                    },
                  })
                }
              />
            </Grid>
          </Grid>
        </Grid>
        <Grid item xs={6}>
          <Grid container>
            <Grid item xs={12}>
              EARFCNUL
            </Grid>
            <Grid item xs={12}>
              <OutlinedInput
                className={classes.input}
                type="number"
                data-testid="earfcnul"
                value={props.lteRanConfigs.fdd_config?.earfcnul}
                onChange={({target}) =>
                  props.setLteRanConfigs({
                    ...props.lteRanConfigs,
                    tdd_config: undefined,
                    fdd_config: {
                      earfcndl: props.lteRanConfigs.fdd_config?.earfcndl ?? 0,
                      earfcnul: parseInt(target.value),
                    },
                  })
                }
              />
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    </ListItem>
  );
}
