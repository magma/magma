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
import Text from '../Text';
import classNames from 'classnames';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(_theme => ({
  root: {
    display: 'flex',
    alignItems: 'center',
    marginBottom: '24px',
  },
  titleText: {
    flexGrow: 1,
  },
}));

type Props = {
  className?: string,
  children: string,
  rightContent?: React.Node,
};

const CardHeader = (props: Props) => {
  const {children, className, rightContent} = props;
  const classes = useStyles();
  return (
    <div className={classNames(classes.root, className)}>
      <Text variant="h6" className={classes.titleText}>
        {children}
      </Text>
      {rightContent}
    </div>
  );
};

export default CardHeader;
