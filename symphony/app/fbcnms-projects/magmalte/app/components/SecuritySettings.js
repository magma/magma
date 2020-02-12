/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {ContextRouter} from 'react-router-dom';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithStyles} from '@material-ui/core';

import Button from '@fbcnms/ui/components/design-system/Button';
import FormGroup from '@material-ui/core/FormGroup';
import FormLabel from '@material-ui/core/FormLabel';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextField from '@material-ui/core/TextField';
import axios from 'axios';

import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  input: {},
  formContainer: {
    margin: theme.spacing(2),
    paddingBottom: theme.spacing(2),
  },
  paper: {
    margin: '10px',
  },
  formGroup: {
    marginBottom: theme.spacing(2),
  },
});

type Props = WithAlert & ContextRouter & WithStyles<typeof styles> & {};

type State = {
  error: string,
  currentPassword: string,
  newPassword: string,
  confirmPassword: string,
};

class SecuritySettings extends React.Component<Props, State> {
  state = {
    error: '',
    currentPassword: '',
    newPassword: '',
    confirmPassword: '',
  };

  render() {
    const {classes} = this.props;

    return (
      <div className={classes.formContainer}>
        <Text data-testid="change-password-title" variant="h5">
          Change Password
        </Text>
        {this.state.error ? (
          <FormLabel error>{this.state.error}</FormLabel>
        ) : null}
        <FormGroup row className={classes.formGroup}>
          <TextField
            required
            label="Current Password"
            type="password"
            value={this.state.currentPassword}
            onChange={this.onCurrentPasswordChange}
            className={classes.input}
          />
        </FormGroup>
        <FormGroup row className={classes.formGroup}>
          <TextField
            required
            label="New Password"
            type="password"
            value={this.state.newPassword}
            onChange={this.onNewPasswordChange}
            className={classes.input}
          />
        </FormGroup>
        <FormGroup row className={classes.formGroup}>
          <TextField
            required
            label="Confirm Password"
            type="password"
            value={this.state.confirmPassword}
            onChange={this.onConfirmPasswordChange}
            className={classes.input}
          />
        </FormGroup>
        <FormGroup row className={classes.formGroup}>
          <Button onClick={this.onSave}>Save</Button>
        </FormGroup>
      </div>
    );
  }

  onCurrentPasswordChange = ({target}) =>
    this.setState({currentPassword: target.value});
  onNewPasswordChange = ({target}) =>
    this.setState({newPassword: target.value});
  onConfirmPasswordChange = ({target}) =>
    this.setState({confirmPassword: target.value});

  onSave = async () => {
    if (
      !this.state.currentPassword ||
      !this.state.newPassword ||
      !this.state.confirmPassword
    ) {
      this.setState({error: 'Please complete all fields'});
      return;
    }

    if (this.state.newPassword !== this.state.confirmPassword) {
      this.setState({error: 'Passwords do not match'});
      return;
    }

    try {
      await axios.post('/user/change_password', {
        currentPassword: this.state.currentPassword,
        newPassword: this.state.newPassword,
      });

      this.props.alert('Success');
      this.setState({
        currentPassword: '',
        newPassword: '',
        confirmPassword: '',
        error: '',
      });
    } catch (error) {
      this.setState({error: error.response.data.error});
    }
  };
}

export default withStyles(styles)(withRouter(withAlert(SecuritySettings)));
