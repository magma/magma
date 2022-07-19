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
 */

import Button from '@material-ui/core/Button';
import Card from '@material-ui/core/Card';
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';
import FormLabel from '@material-ui/core/FormLabel';
import React from 'react';
import Text from '../../theme/design-system/Text';
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
    padding: '18px 16px 16px',
    float: 'right',
  },
  title: {
    marginBottom: '16px',
    textAlign: 'center',
    display: 'block',
  },
} as const;

type Props = {
  title: string;
  csrfToken: string;
  error?: string;
  classes: Record<string, string>;
  action: string;
  ssoEnabled?: boolean;
  ssoAction?: string;
};

class LoginForm extends React.Component<Props> {
  form: HTMLFormElement | null = null;

  render() {
    const {classes, csrfToken, ssoEnabled, ssoAction} = this.props;
    const error = this.props.error ? (
      <FormLabel error>{this.props.error}</FormLabel>
    ) : null;

    const params = new URLSearchParams(window.location.search);
    const to = params.get('to');

    if (ssoEnabled) {
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
              variant="contained"
              color="primary"
              onClick={() => {
                window.location.href =
                  (ssoAction || '') + window.location.search;
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
          <input type="hidden" name="to" value={to!} />
          <CardContent>
            <Text className={classes.title} variant="h5">
              {this.props.title}
            </Text>
            {error}
            <TextField
              name="email"
              label="Email"
              className={classes.input}
              onKeyUp={key => key.keyCode === ENTER_KEY && this.form?.submit()}
            />
            <TextField
              name="password"
              label="Password"
              type="password"
              className={classes.input}
              onKeyUp={key => key.keyCode === ENTER_KEY && this.form?.submit()}
            />
          </CardContent>
          <CardActions className={classes.footer}>
            <Button
              onClick={() => this.form?.submit()}
              variant="contained"
              color="primary">
              Login
            </Button>
          </CardActions>
        </form>
      </Card>
    );
  }
}

export default withStyles(styles)(LoginForm);
