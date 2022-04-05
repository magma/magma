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

import type {Organization} from './Organizations';
import type {SSOSelectedType} from '../../../fbc_js_core/types/auth';
import type {Tab} from '../../../fbc_js_core/types/tabs';

import Button from '../../../fbc_js_core/ui/components/design-system/Button';
import Checkbox from '@material-ui/core/Checkbox';
import FormControl from '@material-ui/core/FormControl';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import FormGroup from '@material-ui/core/FormGroup';
import Grid from '@material-ui/core/Grid';
import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';
import ListItemText from '@material-ui/core/ListItemText';
import LoadingFiller from '../../../fbc_js_core/ui/components/LoadingFiller';
import MenuItem from '@material-ui/core/MenuItem';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import SaveIcon from '@material-ui/icons/Save';
import Select from '@material-ui/core/Select';
import Text from '../../../fbc_js_core/ui/components/design-system/Text';
import TextField from '@material-ui/core/TextField';
import TypedSelect from '../../../fbc_js_core/ui/components/TypedSelect';
import axios from 'axios';
import renderList from '../../../fbc_js_core/util/renderList';
import symphony from '../../../fbc_js_core/ui/theme/symphony';
import {getProjectTabs as getAllProjectTabs} from '../../../fbc_js_core/projects/projects';
import {intersection} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useAxios, useRouter} from '../../../fbc_js_core/ui/hooks';
import {useCallback, useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '../../../fbc_js_core/ui/hooks/useSnackbar';

const useStyles = makeStyles(theme => ({
  header: {
    margin: '10px',
    display: 'flex',
    justifyContent: 'space-between',
  },
  leftIcon: {
    marginRight: theme.spacing(),
    color: symphony.palette.white,
  },
  root: {
    ...theme.mixins.gutters(),
    paddingTop: theme.spacing(2),
    paddingBottom: theme.spacing(2),
  },
  textField: {
    marginLeft: theme.spacing(),
    marginRight: theme.spacing(),
    width: 500,
  },
  networks: {
    flexDirection: 'row',
    flexWrap: 'nowrap',
    marginTop: '16px',
    marginBottom: '8px',
  },
  flexGrow: {
    flexGrow: 1,
  },
}));

type Props = {
  getProjectTabs?: () => Array<{id: Tab, name: string}>,
};

export default function OrganizationEdit(props: Props) {
  const {match} = useRouter();
  const classes = useStyles();
  const enqueueSnackbar = useEnqueueSnackbar();
  const [name, setName] = useState<string>('');
  const [csvCharset, setCsvCharset] = useState<string>('');
  const [tabs, setTabs] = useState<Set<string>>(new Set());
  const [shouldEnableAllNetworks, setShouldEnableAllNetworks] = useState(false);
  const [networkIds, setNetworkIds] = useState<Set<string>>(new Set());
  const [ssoSelectedType, setSSOSelectedType] = useState<SSOSelectedType>(
    'saml',
  );
  const [ssoIssuer, setSSOIssuer] = useState<string>('');
  const [ssoEntrypoint, setSSOEntrypoint] = useState<string>('');
  const [ssoCert, setSSOCert] = useState<string>('');
  const [ssoOidcClientID, setSSOOidcClientID] = useState<string>('');
  const [ssoOidcClientSecret, setSSOOidcClientSecret] = useState<string>('');
  const [
    ssoOidcConfigurationURL,
    setSSOOidcConfigurationURL,
  ] = useState<string>('');

  const orgRequest = useAxios<null, {organization: Organization}>({
    method: 'get',
    url: '/host/organization/async/' + match.params.name,
    onResponse: useCallback(res => {
      const {organization} = res.data;
      setName(organization.name);
      setTabs(new Set(organization.tabs));
      setNetworkIds(new Set(organization.networkIDs));
      setCsvCharset(organization.csvCharset);
      setSSOSelectedType(organization.ssoSelectedType);
      setSSOIssuer(organization.ssoIssuer);
      setSSOEntrypoint(organization.ssoEntrypoint);
      setSSOCert(organization.ssoCert);
      setSSOOidcClientID(organization.ssoOidcClientID);
      setSSOOidcClientSecret(organization.ssoOidcClientSecret);
      setSSOOidcConfigurationURL(organization.ssoOidcConfigurationURL);
    }, []),
  });

  const networksRequest = useAxios({
    method: 'get',
    url: '/host/networks/async',
  });

  useEffect(() => {
    if (orgRequest.response && networksRequest.response) {
      const allNetworks: string[] = networksRequest.response.data;
      const networkIDs = orgRequest.response.data.organization.networkIDs;
      if (intersection(allNetworks, networkIDs).length === allNetworks.length) {
        setShouldEnableAllNetworks(true);
      }
    }
  }, [orgRequest.response, networksRequest.response]);

  if (
    orgRequest.isLoading ||
    networksRequest.isLoading ||
    !orgRequest.response
  ) {
    return <LoadingFiller />;
  }

  const organization = orgRequest.response.data.organization;
  const allNetworks =
    networksRequest.error || !networksRequest.response
      ? []
      : networksRequest.response.data.sort();

  const onSave = () => {
    axios
      .put('/host/organization/async/' + match.params.name, {
        name,
        tabs: Array.from(tabs),
        networkIDs: shouldEnableAllNetworks
          ? allNetworks
          : Array.from(networkIds),
        csvCharset,
        ssoSelectedType,
        ssoIssuer,
        ssoEntrypoint,
        ssoCert,
        ssoOidcClientID,
        ssoOidcClientSecret,
        ssoOidcConfigurationURL,
      })
      .then(_res => {
        enqueueSnackbar('Updated organization successfully', {
          variant: 'success',
        });
      })
      .catch(error => {
        const message = error.response?.data?.error || error;
        enqueueSnackbar(`Unable to save organization: ${message}`, {
          variant: 'error',
        });
      });
  };

  const allTabs = props.getProjectTabs
    ? props.getProjectTabs()
    : getAllProjectTabs();

  return (
    <Grid container spacing={3}>
      <Grid item xs={12}>
        <Text className={classes.header} variant="h4">
          Organization: {organization.name}
        </Text>
      </Grid>
      <Grid item xs={3} />
      <Grid item xs={6}>
        <Paper className={classes.root} elevation={1}>
          <Text variant="h6">Basic Info</Text>
          <form noValidate autoComplete="off">
            <FormGroup>
              <TextField
                label="Name"
                className={classes.textField}
                value={name}
                onChange={evt => setName(evt.target.value)}
                margin="normal"
              />
            </FormGroup>
            <FormGroup>
              <FormControl className={classes.textField} margin="normal">
                <InputLabel htmlFor="tabs">Accessible Tabs</InputLabel>
                <Select
                  multiple
                  value={Array.from(tabs)}
                  onChange={({target}) => setTabs(new Set(target.value))}
                  renderValue={value => {
                    const selectedTabs = new Set(value);
                    return renderList(
                      allTabs
                        .filter(tab => selectedTabs.has(tab.id))
                        .map(x => x.name),
                    );
                  }}
                  input={<Input id="tabs" />}>
                  {allTabs.map(tab => (
                    <MenuItem key={tab.id} value={tab.id}>
                      <Checkbox checked={tabs.has(tab.id)} />
                      <ListItemText primary={tab.name} />
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </FormGroup>
            <FormGroup className={`${classes.networks} ${classes.textField}`}>
              <FormControlLabel
                control={
                  <Checkbox
                    checked={shouldEnableAllNetworks}
                    onChange={({target}) =>
                      setShouldEnableAllNetworks(target.checked)
                    }
                    color="primary"
                    margin="normal"
                  />
                }
                label="Enable All Networks"
              />
              {!shouldEnableAllNetworks && (
                <FormControl className={classes.flexGrow}>
                  <InputLabel htmlFor="network_ids">
                    Accessible Networks
                  </InputLabel>
                  <Select
                    multiple
                    value={Array.from(networkIds)}
                    onChange={({target}) =>
                      setNetworkIds(new Set(target.value))
                    }
                    renderValue={renderList}
                    input={<Input id="network_ids" />}>
                    {allNetworks.map(network => (
                      <MenuItem key={network} value={network}>
                        <Checkbox checked={networkIds.has(network)} />
                        <ListItemText primary={network} />
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              )}
            </FormGroup>
          </form>
        </Paper>
        <Paper className={classes.root} elevation={1}>
          <Text variant="h6">Advanced Settings</Text>
          <form noValidate autoComplete="off">
            <FormGroup>
              <TextField
                label="CSV Charset (default: utf-8)"
                className={classes.textField}
                value={csvCharset}
                onChange={evt => setCsvCharset(evt.target.value)}
                margin="normal"
              />
            </FormGroup>
          </form>
        </Paper>
        <Paper className={classes.root} elevation={1}>
          <Text variant="h6">Single Sign-On</Text>
          <form noValidate autoComplete="off">
            <FormGroup>
              <TypedSelect
                value={ssoSelectedType}
                className={classes.textField}
                items={{
                  none: 'Disabled',
                  oidc: 'OpenID Connect',
                  saml: 'SAML',
                }}
                onChange={value => setSSOSelectedType(value)}
                input={<Input id="ssoSelectedType" />}
              />
            </FormGroup>
            {ssoSelectedType === 'saml' ? (
              <>
                <FormGroup>
                  <TextField
                    label="Issuer"
                    className={classes.textField}
                    value={ssoIssuer}
                    onChange={evt => setSSOIssuer(evt.target.value)}
                    margin="normal"
                  />
                </FormGroup>
                <FormGroup>
                  <TextField
                    label="Entrypoint"
                    className={classes.textField}
                    value={ssoEntrypoint}
                    onChange={evt => setSSOEntrypoint(evt.target.value)}
                    margin="normal"
                  />
                </FormGroup>
                <FormGroup>
                  <TextField
                    label="Certificate"
                    multiline
                    rows="4"
                    value={ssoCert}
                    onChange={evt => setSSOCert(evt.target.value)}
                    className={classes.textField}
                    margin="normal"
                    variant="filled"
                  />
                </FormGroup>
              </>
            ) : null}
            {ssoSelectedType === 'oidc' ? (
              <>
                <FormGroup>
                  <TextField
                    label="Client ID"
                    className={classes.textField}
                    value={ssoOidcClientID}
                    onChange={evt => setSSOOidcClientID(evt.target.value)}
                    margin="normal"
                  />
                </FormGroup>
                <FormGroup>
                  <TextField
                    label="Client Secret"
                    className={classes.textField}
                    value={ssoOidcClientSecret}
                    onChange={evt => setSSOOidcClientSecret(evt.target.value)}
                    margin="normal"
                  />
                </FormGroup>
                <FormGroup>
                  <TextField
                    label="Configuration URL"
                    className={classes.textField}
                    value={ssoOidcConfigurationURL}
                    onChange={evt =>
                      setSSOOidcConfigurationURL(evt.target.value)
                    }
                    margin="normal"
                  />
                </FormGroup>
              </>
            ) : null}
          </form>
        </Paper>
        <Paper className={classes.root} elevation={1}>
          <form noValidate autoComplete="off">
            <Button onClick={() => onSave()}>
              <SaveIcon className={classes.leftIcon} />
              <Text color="light" weight="medium">
                Save
              </Text>
            </Button>
          </form>
        </Paper>
      </Grid>
    </Grid>
  );
}
