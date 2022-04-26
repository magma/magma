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
 * @flow
 * @format
 */

import * as React from 'react';
import Text from '../components/design-system/Text';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  root: {
    padding: '16px',
    boxSizing: 'border-box',
    width: '100%',
    height: '100%',
  },
  header: {
    marginBottom: '24px',
  },
}));

type Props = $ReadOnly<{|
  name: React.Node,
  children: React.Node,
|}>;

const Story = ({name, children}: Props) => {
  const classes = useStyles();
  return (
    <div className={classes.root}>
      <div className={classes.header}>
        <Text variant="h3">{name}</Text>
      </div>
      {children}
    </div>
  );
};

export default Story;
