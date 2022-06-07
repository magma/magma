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
import type {call_trace_config} from '../../../generated/MagmaAPIBindings';

import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import DialogTitle from '../../theme/design-system/DialogTitle';
import FormLabel from '@material-ui/core/FormLabel';
import Grid from '@material-ui/core/Grid';
import List from '@material-ui/core/List';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';
// $FlowFixMe migrated to typescript
import TraceContext from '../../components/context/TraceContext';
import TypedSelect from '../../components/TypedSelect';
import Typography from '@material-ui/core/Typography';

// $FlowFixMe migrated to typescript
import {AltFormField} from '../../components/FormField';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useContext} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';
import {useState} from 'react';

const DEFAULT_TRACE_CONFIG: call_trace_config = {
  gateway_id: '',
  timeout: 300,
  trace_id: '',
  trace_type: 'GATEWAY',
  capture_filters: '',
  display_filters: '',
};

const useStyles = makeStyles(_ => ({
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
}));

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
  open: boolean,
  onClose: () => void,
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
  traceCfg?: call_trace_config,
  onClose: () => void,
};

function CreateTraceDetails(props: Props) {
  const classes = useStyles();
  const ctx = useContext(TraceContext);
  const [error, setError] = useState('');
  const [traceCfg, setTraceCfg] = useState<call_trace_config>({
    ...(props.traceCfg || DEFAULT_TRACE_CONFIG),
  });
  const enqueueSnackbar = useEnqueueSnackbar();

  const startTrace = async (cfg: call_trace_config) => {
    try {
      // $FlowFixMe[prop-missing]: Suppress type error, cannot refine type
      await ctx.setState?.(cfg.trace_id, cfg);
      props.onClose();
      enqueueSnackbar('Call trace started successfully', {
        variant: 'success',
      });
    } catch (e) {
      const errMsg = e.response?.data?.message ?? e.message ?? e;
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
                setTraceCfg({
                  ...traceCfg,
                  capture_filters: target.value,
                });
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
                setTraceCfg({
                  ...traceCfg,
                  display_filters: target.value,
                });
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
          onClick={async () => {
            await startTrace(traceCfg);
          }}>
          {'Start'}
        </Button>
      </DialogActions>
    </>
  );
}
