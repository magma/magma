/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import type {flow_qos} from '@fbcnms/magma-api';

import Checkbox from '@material-ui/core/Checkbox';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Grid from '@material-ui/core/Grid';
import Input from '@material-ui/core/Input';
import InputAdornment from '@material-ui/core/InputAdornment';
import React, {useState} from 'react';
import TextField from '@material-ui/core/TextField';

import {makeStyles} from '@material-ui/styles';
const MAX_BW_SETTING = 1000000000; // 1 gbps

const useStyles = makeStyles(_ => ({
  input: {width: '100%'},
}));

export type QosState = {
  enabled: boolean,
  qos: flow_qos,
};

type Props = {
  qos?: flow_qos,
  onChange: QosState => void,
};

export default function PolicyQosFields(props: Props) {
  const classes = useStyles();
  const [maxUlBw, setMaxUlBw] = useState(props.qos?.max_req_bw_ul ?? 1);
  const [maxDlBw, setMaxDlBw] = useState(props.qos?.max_req_bw_dl ?? 1);
  const [enabled, setEnabled] = useState(props.qos ? true : false);
  const maxBwSetting = 1000000000;
  const err = bw =>
    bw > 0 && bw < MAX_BW_SETTING
      ? null
      : `value must be between 1-${maxBwSetting}`;

  const ulError = err(maxUlBw);
  const dlError = err(maxDlBw);
  return (
    <ExpansionPanel>
      <ExpansionPanelSummary
        expandIcon={<ExpandMoreIcon />}
        aria-label="Expand">
        <FormControlLabel
          checked={enabled}
          onFocus={event => event.stopPropagation()}
          control={<Checkbox />}
          onChange={({target}) => {
            setEnabled(target.checked);
            props.onChange({
              enabled: target.checked,
              qos: {
                max_req_bw_ul: maxUlBw,
                max_req_bw_dl: maxDlBw,
              },
            });
          }}
          label="Enable"
          input={<Input id="qos_enabled" />}
        />
      </ExpansionPanelSummary>
      <ExpansionPanelDetails>
        <Grid container spacing={2} justify="center">
          <Grid item xs={6}>
            <TextField
              error={ulError !== null}
              className={classes.input}
              label="Max UL B/W"
              margin="normal"
              type="number"
              value={maxUlBw}
              helperText={ulError}
              onChange={({target}) => {
                const bw = parseInt(target.value);
                setMaxUlBw(bw);
                props.onChange({
                  enabled: enabled,
                  qos: {
                    max_req_bw_ul: bw,
                    max_req_bw_dl: maxDlBw,
                  },
                });
              }}
              InputProps={{
                endAdornment: (
                  <InputAdornment position="end">bps</InputAdornment>
                ),
              }}
            />
          </Grid>
          <Grid item xs={6}>
            <TextField
              className={classes.input}
              error={dlError !== null}
              label="Max DL B/W"
              margin="normal"
              type="number"
              value={maxDlBw}
              helperText={dlError}
              onChange={({target}) => {
                const bw = parseInt(target.value);
                setMaxDlBw(bw);
                props.onChange({
                  enabled: enabled,
                  qos: {
                    max_req_bw_ul: maxUlBw,
                    max_req_bw_dl: bw,
                  },
                });
              }}
              InputProps={{
                endAdornment: (
                  <InputAdornment position="end">bps</InputAdornment>
                ),
              }}
            />
          </Grid>
        </Grid>
      </ExpansionPanelDetails>
    </ExpansionPanel>
  );
}
