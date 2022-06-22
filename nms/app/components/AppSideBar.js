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

import ProfileButton from './ProfileButton';
import React, {useContext, useState} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import SidebarItem from './SidebarItem';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Text from '../theme/design-system/Text';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import VersionContext from './context/VersionContext';
import classNames from 'classnames';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {colors} from '../theme/default';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  root: {
    width: '80px',
    overflowX: 'visible',
  },
  inner: {
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'space-between',
    backgroundColor: colors.primary.brightGray,
    boxShadow: '1px 0px 0px 0px rgba(0, 0, 0, 0.1)',
    height: '100vh',
    padding: '60px 0px 24px 0px',
    position: 'relative',
    zIndex: 100,
    overflowX: 'hidden',
    width: '80px',
    transition: 'width 500ms',
  },
  expanded: {
    width: '208px',
  },
  version: {
    display: 'block',
    padding: '13px 0 0 28px',
    color: colors.primary.gullGray,
    whiteSpace: 'nowrap',
  },
  versionHidden: {
    visibility: 'hidden',
  },
}));

type ItemConfig = {
  path: string,
  label: string,
  icon: any,
};

type Props = {
  items: Array<ItemConfig>,
};

const AppSideBar = (props: Props) => {
  const {items} = props;
  const classes = useStyles();
  const [expanded, setIsExpanded] = useState(false);
  const [isProfileMenuOpen, _setProfileMenuOpen] = useState(false);
  const {nmsVersion} = useContext(VersionContext);

  const setProfileMenuOpen = (isOpen: boolean) => {
    if (!isOpen) {
      setIsExpanded(false);
    }
    _setProfileMenuOpen(isOpen);
  };

  return (
    <div
      data-testid="app-sidebar"
      className={classes.root}
      onMouseOver={() => setIsExpanded(true)}
      onMouseLeave={() => {
        if (!isProfileMenuOpen) {
          setIsExpanded(false);
        }
      }}>
      <div
        className={classNames({
          [classes.inner]: true,
          [classes.expanded]: expanded,
        })}>
        <div>
          {items.map(({path, label, icon}) => (
            <SidebarItem
              key={label}
              path={path}
              label={label}
              icon={icon}
              expanded={expanded}
            />
          ))}
        </div>
        <div>
          <ProfileButton
            expanded={expanded}
            isMenuOpen={isProfileMenuOpen}
            setMenuOpen={setProfileMenuOpen}
          />
          <Text
            variant="body3"
            className={classNames({
              [classes.version]: true,
              [classes.versionHidden]: !expanded,
            })}>
            {'NMS: ' + nmsVersion}
          </Text>
        </div>
      </div>
    </div>
  );
};

export default AppSideBar;
