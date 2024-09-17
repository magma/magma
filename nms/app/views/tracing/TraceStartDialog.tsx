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
 */

import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '../../theme/design-system/DialogTitle';
import FormLabel from '@mui/material/FormLabel';
import Grid from '@mui/material/Grid';
import List from '@mui/material/List';
import OutlinedInput from '@mui/material/OutlinedInput';
import React from 'react';
import TraceContext from '../../context/TraceContext';
import TypedSelect from '../../components/TypedSelect';
import Typography from '@mui/material/Typography';
import {AltFormField} from '../../components/FormField';
import {colors, typography} from '../../theme/default';
import {getErrorMessage} from '../../util/ErrorUtils';
import {makeStyles} from '@mui/styles';
import {useContext} from 'react';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';
import {useState} from 'react';
import type {CallTraceConfig} from '../../../generated';

const DEFAULT_TRACE_CONFIG: CallTraceConfig = {
  gateway_id: '',
  timeout: 300,
  trace_id: '',
  trace_type: 'GATEWAY',
  capture_filters: '',
  display_filters: '',
};

const useStyles = makeStyles({
  topBar: {
    backgroundColor: colors.primary.mirage,
    padding: '20px 40px 20px 40px',
    color: colors.primary.white,
  },
  appBarBtn: {
    color: colors.primary.white,
    background: colors.primary.comet,
    fontFamily: typography.button.fontFamily,
    fontWeight: typography.button.fontWeight,
    fontSize: typography.button.fontSize,
    lineHeight: typography.button.lineHeight,
    letterSpacing: typography.button.letterSpacing,
    '&:hover': {
      background: colors.primary.mirage,
    },
  },
  addSubscriberDialog: {},
});

export default function CreateTraceButton() {
  const classes = useStyles();
  const [open, setOpen] = useState(false);

  return (
    <>
      <CreateTraceDialog open={open} onClose={() => setOpen(false)} />
      <Button onClick={() => setOpen(true)} className={classes.appBarBtn}>
        {'Start New Trace'}
      </Button>
    </>
  );
}

type DialogProps = {
  open: boolean;
  onClose: () => void;
};

function CreateTraceDialog(props: DialogProps) {
  const classes = useStyles();
  return (
    <Dialog
      data-testid="addSubscriberDialog"
      open={props.open}
      fullWidth={true}
      maxWidth="sm">
      <DialogTitle
        className={classes.topBar}
        onClose={props.onClose}
        label={'Start Call Trace'}
      />

      <CreateTraceDetails onClose={props.onClose} />
    </Dialog>
  );
}

type Props = {
  traceCfg?: CallTraceConfig;
  onClose: () => void;
};

function CreateTraceDetails(props: Props) {
  const classes = useStyles();
  const ctx = useContext(TraceContext);
  const [error, setError] = useState('');
  const [traceCfg, setTraceCfg] = useState<CallTraceConfig>({
    ...(props.traceCfg || DEFAULT_TRACE_CONFIG),
  });
  const enqueueSnackbar = useEnqueueSnackbar();

  const startTrace = async (cfg: CallTraceConfig) => {
    try {
      // TODO[TS-migration] There is something seriously wrong with types here
      // @ts-ignore
      await ctx.setState?.(cfg.trace_id, cfg);
      props.onClose();
      enqueueSnackbar('Call trace started successfully', {
        variant: 'success',
      });
    } catch (e) {
      const errMsg = getErrorMessage(e);
      setError('error starting call trace: ' + errMsg);
    }
  };

  return (
    <>
      <DialogContent>
        <List className={classes.addSubscriberDialog}>
          {error !== '' && (
            <AltFormField label={''}>
              <FormLabel data-testid="configEditError" error>
                {error}
              </FormLabel>
            </AltFormField>
          )}
          <Grid container>
            <Grid item xs={12} sm={6}>
              <AltFormField label={'Trace ID'}>
                <OutlinedInput
                  data-testid="trace-id"
                  placeholder="Enter Trace ID"
                  fullWidth={true}
                  value={traceCfg.trace_id}
                  onChange={({target}) => {
                    setTraceCfg({...traceCfg, trace_id: target.value});
                  }}
                />
              </AltFormField>
            </Grid>
            <Grid item xs={12} sm={6}>
              <AltFormField label={'Timeout'}>
                <OutlinedInput
                  data-testid="timeout"
                  placeholder="Enter Trace Timeout (s)"
                  fullWidth={true}
                  value={traceCfg.timeout}
                  type="number"
                  onChange={({target}) => {
                    setTraceCfg({...traceCfg, timeout: parseInt(target.value)});
                  }}
                />
              </AltFormField>
            </Grid>
          </Grid>
          <AltFormField label={'Trace Type'}>
            <TypedSelect
              disabled={true}
              input={<OutlinedInput />}
              value={traceCfg.trace_type}
              fullWidth={true}
              items={{
                GATEWAY: 'Gateway',
              }}
              onChange={target => {
                setTraceCfg({...traceCfg, trace_type: target});
              }}
            />
          </AltFormField>
          <AltFormField label={'Gateway ID'}>
            <OutlinedInput
              data-testid="gateway-id"
              placeholder="Enter Gateway ID"
              fullWidth={true}
              value={traceCfg.gateway_id}
              onChange={({target}) => {
                setTraceCfg({...traceCfg, gateway_id: target.value});
              }}
            />
          </AltFormField>
          <AltFormField label={'TShark Custom Capture Filters'}>
            <OutlinedInput
              data-testid="capture-filters"
              placeholder="tcp and (port 80 or port 8080)"
              fullWidth={true}
              value={traceCfg.capture_filters}
              onChange={({target}) => {
                setTraceCfg({...traceCfg, capture_filters: target.value});
              }}
            />
          </AltFormField>
          <AltFormField label={'TShark Custom Display Filters'}>
            <OutlinedInput
              data-testid="display-filters"
              placeholder="http"
              fullWidth={true}
              value={traceCfg.display_filters}
              onChange={({target}) => {
                setTraceCfg({...traceCfg, display_filters: target.value});
              }}
            />
          </AltFormField>
          <AltFormField label={'Example Preset Filters'}>
            <Grid container>
              <Grid item xs={12} sm={4}>
                <Button
                  onClick={() => {
                    setTraceCfg({
                      ...traceCfg,
                      capture_filters: 'tcp and port 80',
                      display_filters: 'http',
                    });
                  }}
                  className={classes.appBarBtn}>
                  {'HTTP Messages'}
                </Button>
              </Grid>
              <Grid item xs={12} sm={4}>
                <Button
                  onClick={() => {
                    setTraceCfg({
                      ...traceCfg,
                      capture_filters: '',
                      display_filters: 'dns',
                    });
                  }}
                  className={classes.appBarBtn}>
                  {'DNS Messages'}
                </Button>
              </Grid>
              <Grid item xs={12} sm={4}>
                <Button
                  onClick={() => {
                    setTraceCfg({
                      ...traceCfg,
                      capture_filters: '',
                      display_filters: 'ip.src==192.168.0.0/16',
                    });
                  }}
                  className={classes.appBarBtn}>
                  {'Limit IP Src'}
                </Button>
              </Grid>
            </Grid>
          </AltFormField>
          <AltFormField label={'TShark Command'}>
            <Typography component="div" variant="caption">
              {'tshark -ni any -a filesize:4000 ' +
                (traceCfg.capture_filters || '')}
            </Typography>
          </AltFormField>
        </List>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose}> Cancel </Button>
        <Button
          data-testid="startTrace"
          onClick={() => void startTrace(traceCfg)}>
          {'Start'}
        </Button>
      </DialogActions>
    </>
  );
}
