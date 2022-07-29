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

import Button from '@mui/material/Button';
import Card from '@mui/material/Card';
import CardActions from '@mui/material/CardActions';
import CardContent from '@mui/material/CardContent';
import FormLabel from '@mui/material/FormLabel';
import React from 'react';
import Text from '../../theme/design-system/Text';
import withStyles from '@mui/styles/withStyles';

import TextField from '@mui/material/TextField';
import {AltFormField} from '../../components/FormField';
import {colors} from '../../theme/default';

const ENTER_KEY = 13;
const styles = {
  capitalize: {
    textTransform: 'capitalize',
  },
  card: {
    maxWidth: '400px',
    margin: '24px auto 0',
    padding: '20px 0',
  },
  cardContent: {
    padding: '0 24px',
  },
  input: {
    display: 'inline-flex',
    width: '100%',
  },
  footer: {
    padding: '0 24px',
  },
  login: {
    marginTop: '10%',
  },
  title: {
    textAlign: 'center',
    display: 'flex',
    justifyContent: 'center',
    margin: '12px auto 0',
    flexDirection: 'column',
    maxWidth: '400px',
    alignItems: 'start',
  },
  formTitle: {
    marginBottom: '8px',
    textAlign: 'left',
    display: 'block',
    fontSize: '20px',
  },
  FormField: {
    paddingBottom: '0',
  },
  submitButton: {
    width: '100%',
    marginTop: '16px',
  },
  regular: {
    color: colors.primary.mirage,
  },
  gray: {
    color: colors.primary.gullGray,
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

const HOST_PORTAL_TITLE = 'Magma Host Portal';
const ORGANIZATION_PORTAL_TITLE = 'Magma Organization Portal';

class LoginForm extends React.Component<Props> {
  form: HTMLFormElement | null = null;

  render() {
    const {classes, csrfToken, ssoEnabled, ssoAction} = this.props;
    const error = this.props.error ? (
      <FormLabel error>{this.props.error}</FormLabel>
    ) : null;

    const params = new URLSearchParams(window.location.search);
    const to = params.get('to');

    const [organization] = window.location.hostname.split('.', 1);
    const hostPortal: boolean = organization === 'host';

    return (
      <div className={classes.login} data-testid="loginForm">
        <div className={classes.title}>
          {!hostPortal && (
            <Text className={classes.capitalize} variant="h4">
              {organization}
            </Text>
          )}
          <Text
            className={hostPortal ? classes.regular : classes.gray}
            variant={hostPortal ? 'h4' : 'subtitle1'}>
            {hostPortal ? HOST_PORTAL_TITLE : ORGANIZATION_PORTAL_TITLE}
          </Text>
        </div>
        <Card raised={true} className={classes.card}>
          {ssoEnabled ? (
            <>
              <CardContent className={classes.cardContent}>
                <Text className={classes.formTitle} variant="h6">
                  {`Log in to ${
                    hostPortal ? 'host' : organization
                  } user account`}
                </Text>
                {error}
              </CardContent>
              <CardActions className={classes.footer}>
                <Button
                  color="primary"
                  variant="contained"
                  className={classes.submitButton}
                  onClick={() => {
                    window.location.href =
                      (ssoAction || '') + window.location.search;
                  }}>
                  Login
                </Button>
              </CardActions>
            </>
          ) : (
            <>
              <form
                ref={ref => (this.form = ref)}
                method="post"
                action={this.props.action}>
                <input type="hidden" name="_csrf" value={csrfToken} />
                <input type="hidden" name="to" value={to || ''} />
                <CardContent className={classes.cardContent}>
                  <Text className={classes.formTitle} variant="h6">
                    {`Log in to ${
                      hostPortal ? 'host' : organization
                    } user account`}
                  </Text>
                  {error}
                  <AltFormField
                    disableGutters
                    label={'Email'}
                    className={classes.FormField}>
                    <TextField
                      variant="outlined"
                      name="email"
                      placeholder="Email"
                      className={classes.input}
                      onKeyUp={key =>
                        key.keyCode === ENTER_KEY && this.form!.submit()
                      }
                    />
                  </AltFormField>
                  <AltFormField disableGutters label={'Password'}>
                    <TextField
                      variant="outlined"
                      name="password"
                      placeholder="Password"
                      type="password"
                      className={classes.input}
                      onKeyUp={key =>
                        key.keyCode === ENTER_KEY && this.form!.submit()
                      }
                    />
                  </AltFormField>
                </CardContent>
                <CardActions className={classes.footer}>
                  <Button
                    color="primary"
                    variant="contained"
                    className={classes.submitButton}
                    onClick={() => this.form!.submit()}>
                    Login
                  </Button>
                </CardActions>
              </form>
            </>
          )}
        </Card>
      </div>
    );
  }
}

export default withStyles(styles)(LoginForm);
