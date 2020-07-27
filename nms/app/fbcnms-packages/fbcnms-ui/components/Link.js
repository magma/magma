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

import * as React from 'react';

import Button from './design-system/Button';
import {makeStyles} from '@material-ui/styles';

type Props = {
  onClick: () => void,
  children: React.Node,
};

const useStyles = makeStyles(() => ({
  root: {
    textDecoration: 'underline',
  },
}));

// TODO(T38660666) - style according to design
export default function Link(props: Props) {
  const classes = useStyles();
  const {onClick, children} = props;
  return (
    <Button
      variant="text"
      useEllipsis={true}
      className={classes.root}
      onClick={onClick}>
      {children}
    </Button>
  );
}
