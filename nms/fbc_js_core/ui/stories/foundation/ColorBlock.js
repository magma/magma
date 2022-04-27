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

import React from 'react';
import Text from '../../components/design-system/Text';
import classNames from 'classnames';
import {makeStyles} from '@material-ui/styles';

type Props = {
  color: string,
  name: string,
  code?: string,
  className?: string,
};

const useStyles = makeStyles(_theme => ({
  root: {
    display: 'flex',
    flexDirection: 'column',
    width: '160px',
  },
  block: {
    height: '160px',
    padding: '9px 12px',
    boxSizing: 'border-box',
  },
  nameContainer: {
    display: 'flex',
    flexDirection: 'row',
    padding: '8px',
    backgroundColor: 'white',
    alignItems: 'center',
  },
  name: {
    flexGrow: 1,
  },
}));

const ColorBlock = (props: Props) => {
  const {color, name, code, className} = props;
  const classes = useStyles();
  return (
    <div className={classNames(classes.root, className)}>
      <div className={classes.block} style={{backgroundColor: color}} />
      <div className={classes.nameContainer}>
        <Text className={classes.name} weight="medium">
          {name}
        </Text>
        <Text weight="medium">{code ?? color}</Text>
      </div>
    </div>
  );
};

export default ColorBlock;
