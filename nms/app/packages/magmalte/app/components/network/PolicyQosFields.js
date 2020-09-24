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
import type {policy_qos_profile} from '@fbcnms/magma-api';

import Checkbox from '@material-ui/core/Checkbox';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Grid from '@material-ui/core/Grid';
import Input from '@material-ui/core/Input';
import InputAdornment from '@material-ui/core/InputAdornment';
import React from 'react';
import TextField from '@material-ui/core/TextField';

import {makeStyles} from '@material-ui/styles';

const MAX_BW_SETTING = 1000000000; // 1 gbps

const useStyles = makeStyles(_ => ({
  input: {width: '100%'},
}));

type Props = {
  qosProfile: policy_qos_profile,
  setQosProfile: policy_qos_profile => void,
  qosEnabled: boolean,
  setIsQosEnabled: boolean => void,
};

export default function PolicyQosFields(props: Props) {
  const classes = useStyles();
  const {qosProfile, setQosProfile, qosEnabled, setIsQosEnabled} = props;

  const maxBwSetting = 1000000000;
  const err = bw =>
    bw > 0 && bw < MAX_BW_SETTING
      ? null
      : `value must be between 1-${maxBwSetting}`;
  const ulError = qosProfile?.max_req_bw_ul
    ? err(qosProfile.max_req_bw_ul)
    : null;
  const dlError = qosProfile?.max_req_bw_dl
    ? err(qosProfile.max_req_bw_dl)
    : null;
  return (
    <ExpansionPanel>
      <ExpansionPanelSummary
        expandIcon={<ExpandMoreIcon />}
        aria-label="Expand">
        <FormControlLabel
          checked={qosEnabled}
          onFocus={event => event.stopPropagation()}
          control={<Checkbox />}
          onChange={({target}) => setIsQosEnabled(target.checked)}
          label="Enable"
          input={<Input id="qos_enabled" />}
        />
      </ExpansionPanelSummary>
      <ExpansionPanelDetails>
        <Grid container spacing={2} justify="center">
          <Grid item xs={12}>
            <TextField
              className={classes.input}
              label="Profile ID"
              value={qosProfile?.id}
              onChange={({target}) =>
                setQosProfile({
                  ...qosProfile,
                  id: target.value,
                })
              }
            />
          </Grid>
          <Grid item xs={12}>
            <TextField
              error={ulError !== null}
              className={classes.input}
              label="Max UL B/W"
              type="number"
              value={qosProfile?.max_req_bw_ul ?? 1}
              helperText={ulError}
              onChange={({target}) =>
                setQosProfile({
                  ...qosProfile,
                  max_req_bw_ul: parseInt(target.value),
                })
              }
              InputProps={{
                endAdornment: (
                  <InputAdornment position="end">bps</InputAdornment>
                ),
              }}
            />
          </Grid>
          <Grid item xs={12}>
            <TextField
              className={classes.input}
              error={dlError !== null}
              label="Max DL B/W"
              type="number"
              value={qosProfile?.max_req_bw_dl ?? 1}
              helperText={dlError}
              onChange={({target}) =>
                setQosProfile({
                  ...qosProfile,
                  max_req_bw_dl: parseInt(target.value),
                })
              }
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
