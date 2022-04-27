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

import type {WithStyles} from '@material-ui/core';

import * as React from 'react';
import ErrorIcon from '@material-ui/icons/Error';
import Text from '../design-system/Text';
import {withStyles} from '@material-ui/core/styles';

const styles = {
  root: {
    padding: '8px',
    display: 'flex',
    alignItems: 'center',
  },
  errorIcon: {
    marginRight: '8px',
  },
};

type Props = {
  children: React.Node,
  onError?: ?(error: Error) => void,
} & WithStyles<typeof styles>;

type State = {
  error: ?Error,
};

class ErrorBoundary extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = {error: null};
  }

  componentDidCatch(error: Error) {
    this.setState({
      error: error,
    });
    this.props.onError && this.props.onError(error);
  }

  render() {
    const {classes} = this.props;
    if (this.state.error) {
      return (
        <div className={classes.root}>
          <ErrorIcon size="small" className={classes.errorIcon} />
          <Text variant="body1">Oops, something went wrong.</Text>
        </div>
      );
    }
    return this.props.children;
  }
}

export default withStyles(styles)(ErrorBoundary);
