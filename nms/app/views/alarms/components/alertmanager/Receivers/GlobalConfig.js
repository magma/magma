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
 *
 * Edit's alertmanager's "global" config section. This feature is NOT available
 * in Magma NMS since it's multitenant.
 */

import * as React from 'react';
import Checkbox from '@material-ui/core/Checkbox';
import CircularProgress from '@material-ui/core/CircularProgress';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Editor from '../../common/Editor';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import TextField from '@material-ui/core/TextField';
import Typography from '@material-ui/core/Typography';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import useForm from '../../../hooks/useForm';
import {makeStyles} from '@material-ui/styles';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useAlarmContext} from '../../AlarmContext';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useNetworkId} from '../../hooks';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useSnackbars} from '../../../../../hooks/useSnackbar';

// $FlowFixMe migrated to typescript
import type {AlertManagerGlobalConfig, HTTPConfig} from '../../AlarmAPIType';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {Props as EditorProps} from '../../common/Editor';

type Props = $Diff<
  EditorProps,
  {
    children: React.Node,
    title?: string,
    onSave: () => Promise<void> | void,
    isNew: boolean,
  },
>;

const useStyles = makeStyles(theme => ({
  root: {
    padding: theme.spacing(4),
  },
  loading: {
    display: 'flex',
    height: '100%',
    alignItems: 'center',
    justifyContent: 'center',
  },
}));

export default function GlobalConfig(props: Props) {
  const classes = useStyles();
  const {apiUtil} = useAlarmContext();
  const snackbars = useSnackbars();
  const [lastRefreshTime, _setLastRefreshTime] = React.useState(new Date());
  const networkId = useNetworkId();
  const {response, isLoading} = apiUtil.useAlarmsApi(
    apiUtil.getGlobalConfig,
    {networkId},
    lastRefreshTime.toLocaleString(),
  );

  const {formState, handleInputChange, updateFormState, setFormState} = useForm<
    $Shape<AlertManagerGlobalConfig>,
  >({
    initialState: {
      smtp_require_tls: true,
    },
  });

  const updateHttpConfigState = React.useCallback(
    (update: $Shape<HTTPConfig>) => {
      updateFormState({
        http_config: {
          ...(formState.http_config || {}),
          ...update,
        },
      });
    },
    [formState, updateFormState],
  );

  React.useEffect(() => {
    if (response) {
      setFormState(response);
    }
  }, [response, setFormState]);

  const handleSave = React.useCallback(async () => {
    try {
      const formStateCleaned = removeEmptys(formState);
      await apiUtil.editGlobalConfig({
        config: formStateCleaned,
        networkId,
      });
      snackbars.success('Successfully saved global config');
    } catch (error) {
      snackbars.error(
        `Unable to save global config: ${
          error.response ? error.response?.data?.message : error.message
        }`,
      );
    }
  }, [networkId, apiUtil, formState, snackbars]);

  if (isLoading) {
    return (
      <div className={classes.loading}>
        <CircularProgress />
      </div>
    );
  }
  return (
    <Grid className={classes.root} container>
      <Grid item xs={6}>
        <Paper elevation={1}>
          <Editor
            {...props}
            onSave={handleSave}
            title="Global Receiver Settings"
            description="Default settings which apply to all receivers."
            isNew={false}>
            <Grid container item direction="column" spacing={4}>
              <Grid item>
                <TextField
                  label="Resolve Timeout"
                  placeholder="Ex: 5s"
                  {...getIdProps('resolve_timeout')}
                  value={formState.resolve_timeout || ''}
                  onChange={handleInputChange(val => ({resolve_timeout: val}))}
                  fullWidth
                />
              </Grid>
              <ConfigSection title="Slack">
                <Grid item>
                  <TextField
                    label="API URL"
                    placeholder="Ex: https://hooks.slack.com/services/T0/B0/XXX"
                    {...getIdProps('slack_api_url')}
                    value={formState.slack_api_url || ''}
                    onChange={handleInputChange(val => ({
                      slack_api_url: val,
                    }))}
                    fullWidth
                  />
                </Grid>
              </ConfigSection>
              <ConfigSection title="Pagerduty">
                <Grid item>
                  <TextField
                    label="API URL"
                    placeholder="Ex: https://api.pagerduty.com"
                    {...getIdProps('pagerduty_url')}
                    value={formState.pagerduty_url || ''}
                    onChange={handleInputChange(val => ({pagerduty_url: val}))}
                    fullWidth
                  />
                </Grid>
              </ConfigSection>
              <ConfigSection title="SMTP">
                <Grid item>
                  <TextField
                    label="From"
                    placeholder="Ex: alert@terragraph.link"
                    {...getIdProps('smtp_from')}
                    value={formState.smtp_from || ''}
                    onChange={handleInputChange(val => ({smtp_from: val}))}
                    fullWidth
                  />
                </Grid>
                <Grid item>
                  <TextField
                    label="SMTP Hello"
                    placeholder="Ex: terragraph.link"
                    {...getIdProps('smtp_hello')}
                    value={formState.smtp_hello || ''}
                    onChange={handleInputChange(val => ({smtp_hello: val}))}
                    fullWidth
                  />
                </Grid>
                <Grid item>
                  <TextField
                    label="Smarthost"
                    {...getIdProps('smtp_smarthost')}
                    value={formState.smtp_smarthost || ''}
                    onChange={handleInputChange(val => ({smtp_smarthost: val}))}
                    fullWidth
                  />
                </Grid>
                <Grid item>
                  <TextField
                    label="Username"
                    {...getIdProps('smtp_auth_username')}
                    value={formState.smtp_auth_username || ''}
                    onChange={handleInputChange(val => ({
                      smtp_auth_username: val,
                    }))}
                    fullWidth
                  />
                </Grid>
                <Grid item>
                  <TextField
                    label="Password"
                    {...getIdProps('smtp_auth_password')}
                    value={formState.smtp_auth_password || ''}
                    onChange={handleInputChange(val => ({
                      smtp_auth_password: val,
                    }))}
                    fullWidth
                  />
                </Grid>
                <Grid item>
                  <TextField
                    label="Secret"
                    {...getIdProps('smtp_auth_secret')}
                    value={formState.smtp_auth_secret || ''}
                    onChange={handleInputChange(val => ({
                      smtp_auth_secret: val,
                    }))}
                    fullWidth
                  />
                </Grid>
                <Grid item>
                  <TextField
                    label="Identity"
                    {...getIdProps('smtp_auth_identity')}
                    value={formState.smtp_auth_identity || ''}
                    onChange={handleInputChange(val => ({
                      smtp_auth_identity: val,
                    }))}
                    fullWidth
                  />
                </Grid>
                <Grid item>
                  <FormControlLabel
                    control={
                      <Checkbox
                        checked={
                          // prevents uncontrolled->controlled error
                          typeof formState.smtp_require_tls === 'boolean'
                            ? formState.smtp_require_tls
                            : true
                        }
                        onChange={handleInputChange((_, event) => {
                          return {
                            smtp_require_tls: event.target.checked,
                          };
                        })}
                      />
                    }
                    label="Require TLS"
                  />
                </Grid>
              </ConfigSection>
              <ConfigSection title="Opsgenie">
                <Grid item>
                  <TextField
                    label="API URL"
                    placeholder="Ex: https://api.opsgenie.com/"
                    {...getIdProps('opsgenie_api_url')}
                    value={formState.opsgenie_api_url || ''}
                    onChange={handleInputChange(val => ({
                      opsgenie_api_url: val,
                    }))}
                    fullWidth
                  />
                </Grid>
                <Grid item>
                  <TextField
                    label="API Key"
                    placeholder="Ex: xxxxxxxx-xxxx-xxxx-xxxxx-xxxxxxxxxxxx"
                    {...getIdProps('opsgenie_api_key')}
                    value={formState.opsgenie_api_key || ''}
                    onChange={handleInputChange(val => ({
                      opsgenie_api_key: val,
                    }))}
                    fullWidth
                  />
                </Grid>
              </ConfigSection>
              <ConfigSection title="Hipchat">
                <Grid item>
                  <TextField
                    label="API URL"
                    placeholder="Ex: https://api.hipchat.com/v2"
                    {...getIdProps('hipchat_api_url')}
                    value={formState.hipchat_api_url || ''}
                    onChange={handleInputChange(val => ({
                      hipchat_api_url: val,
                    }))}
                    fullWidth
                  />
                </Grid>
                <Grid item>
                  <TextField
                    label="API Key"
                    placeholder="Ex: xxx-xxx-xxxx"
                    {...getIdProps('hipchat_auth_token')}
                    value={formState.hipchat_auth_token || ''}
                    onChange={handleInputChange(val => ({
                      hipchat_auth_token: val,
                    }))}
                    fullWidth
                  />
                </Grid>
              </ConfigSection>
              <ConfigSection title="WeChat">
                <Grid item>
                  <TextField
                    label="API URL"
                    placeholder="Ex: https://qyapi.weixin.qq.com/cgi-bin/"
                    {...getIdProps('wechat_api_url')}
                    value={formState.wechat_api_url || ''}
                    onChange={handleInputChange(val => ({
                      wechat_api_url: val,
                    }))}
                    fullWidth
                  />
                </Grid>
                <Grid item>
                  <TextField
                    label="API Key"
                    placeholder="Ex: xxxxx"
                    {...getIdProps('wechat_api_secret')}
                    value={formState.wechat_api_secret || ''}
                    onChange={handleInputChange(val => ({
                      wechat_api_secret: val,
                    }))}
                    fullWidth
                  />
                </Grid>
                <Grid item>
                  <TextField
                    label="Corp ID"
                    placeholder="Ex: xxxxx"
                    {...getIdProps('wechat_api_corp_id')}
                    value={formState.wechat_api_corp_id || ''}
                    onChange={handleInputChange(val => ({
                      wechat_api_corp_id: val,
                    }))}
                    fullWidth
                  />
                </Grid>
              </ConfigSection>
              <ConfigSection title="VictorOps">
                <Grid item>
                  <TextField
                    label="API URL"
                    placeholder="Ex: https://api.hipchat.com/v2"
                    {...getIdProps('victorops_api_url')}
                    value={formState.victorops_api_url || ''}
                    onChange={handleInputChange(val => ({
                      victorops_api_url: val,
                    }))}
                    fullWidth
                  />
                </Grid>
                <Grid item>
                  <TextField
                    label="API Key"
                    placeholder="Ex: xxx-xxx-xxxx"
                    {...getIdProps('victorops_api_key')}
                    value={formState.victorops_api_key || ''}
                    onChange={handleInputChange(val => ({
                      victorops_api_key: val,
                    }))}
                    fullWidth
                  />
                </Grid>
              </ConfigSection>
              <ConfigSection title="HTTP">
                <Grid item>
                  <TextField
                    {...getIdProps('http_config_bearer_token')}
                    label="Bearer Token"
                    value={formState.http_config?.bearer_token || ''}
                    onChange={e => {
                      updateHttpConfigState({
                        bearer_token: e.target.value,
                      });
                    }}
                    fullWidth
                  />
                </Grid>
                <Grid item>
                  <TextField
                    label="Proxy URL"
                    {...getIdProps('http_config_proxy_url')}
                    value={formState.http_config?.proxy_url || ''}
                    onChange={e => {
                      updateHttpConfigState({
                        proxy_url: e.target.value,
                      });
                    }}
                    fullWidth
                  />
                </Grid>
                <Grid item>
                  <TextField
                    label="Basic Auth Username"
                    {...getIdProps('basic_auth_username')}
                    value={formState.http_config?.basic_auth?.username || ''}
                    onChange={e => {
                      updateHttpConfigState({
                        basic_auth: {
                          username: e.target.value,
                          password:
                            formState.http_config?.basic_auth?.password || '',
                        },
                      });
                    }}
                    fullWidth
                  />
                </Grid>
                <Grid item>
                  <TextField
                    label="Basic Auth Password"
                    {...getIdProps('basic_auth_password')}
                    value={formState.http_config?.basic_auth?.password || ''}
                    onChange={e => {
                      updateHttpConfigState({
                        basic_auth: {
                          username:
                            formState.http_config?.basic_auth?.username || '',
                          password: e.target.value,
                        },
                      });
                    }}
                    fullWidth
                  />
                </Grid>
              </ConfigSection>
            </Grid>
          </Editor>
        </Paper>
      </Grid>
    </Grid>
  );
}

function ConfigSection({
  title,
  children,
}: {
  title: React.Node,
  children: React.Node,
}) {
  return (
    <Grid item container direction="column">
      <Grid item>
        <Typography variant="h6" color="textSecondary">
          {title}
        </Typography>
      </Grid>
      <Grid item container direction="column">
        {children}
      </Grid>
    </Grid>
  );
}

// Omit config keys that are set to empty strings
function removeEmptys(
  obj: $Shape<AlertManagerGlobalConfig>,
): $Shape<AlertManagerGlobalConfig> {
  const cleaned = {};
  for (const key in obj) {
    const val = obj[key];
    if (typeof val === 'string') {
      if (val.trim() !== '') {
        cleaned[key] = val;
      }
    } else if (typeof val === 'object') {
      cleaned[key] = removeEmptys(val);
    } else {
      cleaned[key] = val;
    }
  }
  return cleaned;
}

/**
 * Handles setting id and data-testid to the same value. id is needed for
 * accessibility purposes so the label's for value is set correctly. testid is
 * used for testing.
 */
function getIdProps(id: string) {
  return {
    id,
    // inputProps targets the raw HTMLInputElement so tests can assert its value
    inputProps: {'data-testid': id},
  };
}
