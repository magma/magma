/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ElementRef} from 'react';

import Button from '../design-system/Button';
import Card from '@material-ui/core/Card';
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';
import FormLabel from '@material-ui/core/FormLabel';
import React from 'react';
import Text from '../design-system/Text';
import TextField from '@material-ui/core/TextField';
import {withStyles} from '@material-ui/core/styles';

const ENTER_KEY = 13;
const styles = {
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
  title: {
    marginBottom: '16px',
    textAlign: 'center',
    display: 'block',
  },
};

type Props = {
  title: string,
  csrfToken: string,
  error?: string,
  classes: {[string]: string},
  action: string,
  isSSO?: boolean,
  ssoAction?: string,
};

type State = {};

class LoginForm extends React.Component<Props, State> {
  form: ElementRef<any>;

  render() {
    const {classes, csrfToken, isSSO, ssoAction} = this.props;
    const error = this.props.error ? (
      <FormLabel error>{this.props.error}</FormLabel>
    ) : null;

    const params = new URLSearchParams(window.location.search);
    const to = params.get('to');

    if (isSSO) {
      return (
        <Card raised={true} className={classes.card}>
          <CardContent>
            <Text className={classes.title} variant="h5">
              {this.props.title}
            </Text>
            {error}
          </CardContent>
          <CardActions className={classes.footer}>
            <Button
              onClick={() => {
                window.location = (ssoAction || '') + window.location.search;
              }}>
              Sign In
            </Button>
          </CardActions>
        </Card>
      );
    }

    return (
      <Card raised={true} className={classes.card}>
        <form
          ref={ref => (this.form = ref)}
          method="post"
          action={this.props.action}>
          <input type="hidden" name="_csrf" value={csrfToken} />
          <input type="hidden" name="to" value={to} />
          <CardContent>
            <Text className={classes.title} variant="h5">
              {this.props.title}
            </Text>
            {error}
            <TextField
              name="email"
              label="Email"
              className={classes.input}
              onKeyUp={key => key.keyCode === ENTER_KEY && this.form.submit()}
            />
            <TextField
              name="password"
              label="Password"
              type="password"
              className={classes.input}
              onKeyUp={key => key.keyCode === ENTER_KEY && this.form.submit()}
            />
          </CardContent>
          <CardActions className={classes.footer}>
            <Button onClick={() => this.form.submit()}>Login</Button>
          </CardActions>
        </form>
      </Card>
    );
  }
}

export default withStyles(styles)(LoginForm);
