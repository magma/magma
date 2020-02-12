/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {Organization} from './Organizations';

import Button from '@fbcnms/ui/components/design-system/Button';
import Checkbox from '@material-ui/core/Checkbox';
import FormControl from '@material-ui/core/FormControl';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import FormGroup from '@material-ui/core/FormGroup';
import Grid from '@material-ui/core/Grid';
import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';
import ListItemText from '@material-ui/core/ListItemText';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MenuItem from '@material-ui/core/MenuItem';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import SaveIcon from '@material-ui/icons/Save';
import Select from '@material-ui/core/Select';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextField from '@material-ui/core/TextField';
import axios from 'axios';
import renderList from '@fbcnms/util/renderList';
import symphony from '@fbcnms/ui/theme/symphony';
import {getProjectTabs} from '@fbcnms/magmalte/app/common/projects';
import {intersection} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useAxios, useRouter} from '@fbcnms/ui/hooks';
import {useCallback, useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';

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

export default function OrganizationEdit() {
  const {match} = useRouter();
  const classes = useStyles();
  const enqueueSnackbar = useEnqueueSnackbar();
  const [name, setName] = useState<string>('');
  const [csvCharset, setCsvCharset] = useState<string>('');
  const [tabs, setTabs] = useState<Set<string>>(new Set());
  const [shouldEnableAllNetworks, setShouldEnableAllNetworks] = useState(false);
  const [networkIds, setNetworkIds] = useState<Set<string>>(new Set());
  const [ssoIssuer, setSSOIssuer] = useState<string>('');
  const [ssoEntrypoint, setSSOEntrypoint] = useState<string>('');
  const [ssoCert, setSSOCert] = useState<string>('');

  const orgRequest = useAxios<null, {organization: Organization}>({
    method: 'get',
    url: '/master/organization/async/' + match.params.name,
    onResponse: useCallback(res => {
      const {organization} = res.data;
      setName(organization.name);
      setTabs(new Set(organization.tabs));
      setNetworkIds(new Set(organization.networkIDs));
      setCsvCharset(organization.csvCharset);
      setSSOIssuer(organization.ssoIssuer);
      setSSOEntrypoint(organization.ssoEntrypoint);
      setSSOCert(organization.ssoCert);
    }, []),
  });

  const networksRequest = useAxios({
    method: 'get',
    url: '/master/networks/async',
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
      .put('/master/organization/async/' + match.params.name, {
        name,
        tabs: Array.from(tabs),
        networkIDs: shouldEnableAllNetworks
          ? allNetworks
          : Array.from(networkIds),
        csvCharset,
        ssoIssuer,
        ssoEntrypoint,
        ssoCert,
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

  const allTabs = getProjectTabs();
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
          <form className={classes.container} noValidate autoComplete="off">
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
          <form className={classes.container} noValidate autoComplete="off">
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
          <form className={classes.container} noValidate autoComplete="off">
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
