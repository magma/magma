/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import '@fbcnms/babel-register/polyfill';

import Button from '@fbcnms/ui/components/design-system/Button';
import Card from '@material-ui/core/Card';
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';
import Checkbox from '@material-ui/core/Checkbox';
import FormControl from '@material-ui/core/FormControl';
import FormGroup from '@material-ui/core/FormGroup';
import FormLabel from '@material-ui/core/FormLabel';
import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';
import ListItemText from '@material-ui/core/ListItemText';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import ReactDOM from 'react-dom';
import Select from '@material-ui/core/Select';
import TextField from '@material-ui/core/TextField';
import Typography from '@material-ui/core/Typography';
import axios from 'axios';

import {} from './common/axiosConfig';
import nullthrows from '@fbcnms/util/nullthrows';
import renderList from '@fbcnms/util/renderList';
import {getProjectTabs} from '@fbcnms/projects/projects';
import {makeStyles} from '@material-ui/styles';
import {useState} from 'react';

const useStyles = makeStyles(() => ({
  card: {
    maxWidth: '400px',
    margin: '10% auto 0',
  },
  input: {
    display: 'inline-flex',
    width: '100%',
    margin: '5px 0',
  },
  footer: {
    marginTop: '10px',
    float: 'right',
  },
}));

const ENTER_KEY = 13;
function Index() {
  const classes = useStyles();
  const [organization, setOrganization] = useState('');
  const [tabs, setTabs] = useState<Set<string>>(new Set());
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');

  const onClick = async () => {
    if (password !== confirmPassword) {
      setError('Passwords do not match');
      return;
    }

    try {
      await axios.post('/user/onboarding', {
        organization,
        tabs: Array.from(tabs),
        email,
        password,
      });

      const {protocol, host} = window.location;
      const to = encodeURIComponent(
        `/admin/organizations/detail/${organization}`,
      );
      window.location.href = `${protocol}//${organization}.${host}/user/login?to=${to}`;
    } catch (error) {
      setError(error.response?.data?.error || error);
    }
  };

  return (
    <Card raised={true} className={classes.card}>
      <CardContent>
        <Typography variant="h5" align="center" gutterBottom>
          Create new Organization and User
        </Typography>
        {error && <FormLabel error>{error}</FormLabel>}
        <TextField
          value={organization}
          label="Organization"
          className={classes.input}
          onChange={event => setOrganization(event.target.value)}
        />
        <FormGroup>
          <FormControl margin="normal">
            <InputLabel htmlFor="tabs">Accessible Tabs</InputLabel>
            <Select
              multiple
              value={Array.from(tabs)}
              onChange={({target}) => setTabs(new Set(target.value))}
              renderValue={renderList}
              input={<Input id="tabs" />}>
              {getProjectTabs().map(tab => (
                <MenuItem key={tab.id} value={tab.id}>
                  <Checkbox checked={tabs.has(tab.id)} />
                  <ListItemText primary={tab.name} />
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </FormGroup>
        <TextField
          value={email}
          label="Email"
          className={classes.input}
          onChange={event => setEmail(event.target.value)}
        />
        <TextField
          autoComplete="off"
          value={password}
          label="Password"
          type="password"
          className={classes.input}
          onChange={event => setPassword(event.target.value)}
        />
        <TextField
          autoComplete="off"
          value={confirmPassword}
          label="Confirm Password"
          type="password"
          className={classes.input}
          onChange={event => setConfirmPassword(event.target.value)}
          onKeyUp={key => key.keyCode === ENTER_KEY && onClick()}
        />
      </CardContent>
      <CardActions className={classes.footer}>
        <Button onClick={onClick}>Submit</Button>
      </CardActions>
    </Card>
  );
}

ReactDOM.render(<Index />, nullthrows(document.getElementById('root')));
