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
import Text from './design-system/Text';
import classNames from 'classnames';
import symphony from '../theme/symphony';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  root: {
    backgroundColor: symphony.palette.white,
    width: '234px',
    height: '100%',
    padding: '24px',
  },
  title: {
    display: 'block',
    marginBottom: '16px',
  },
  item: {
    borderRadius: '4px',
    cursor: 'pointer',
    padding: '10px 16px',
  },
  selectedItem: {
    backgroundColor: symphony.palette.background,
  },
  itemText: {
    lineHeight: '20px',
  },
}));

type NavigationItem = {
  key: string,
  label: string,
};

type Props = {
  className?: string,
  title: string,
  items: Array<NavigationItem>,
  selectedItemId: string,
  onItemClicked: (item: NavigationItem) => void,
};

const SideNavigationPanel = ({
  className,
  title,
  items,
  onItemClicked,
  selectedItemId,
}: Props) => {
  const classes = useStyles();
  return (
    <div className={classNames(classes.root, className)}>
      <Text className={classes.title} weight="medium">
        {title}
      </Text>
      {items.map(item => (
        <div
          className={classNames(classes.item, {
            [classes.selectedItem]: item.key === selectedItemId,
          })}
          onClick={() => onItemClicked(item)}>
          <Text
            className={classes.itemText}
            weight="medium"
            color={item.key === selectedItemId ? 'primary' : 'gray'}>
            {item.label}
          </Text>
        </div>
      ))}
    </div>
  );
};

export default SideNavigationPanel;
