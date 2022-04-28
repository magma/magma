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

import NetworkSelector from '../../../../app/components/NetworkSelector';
import ProfileButton from '../ProfileButton';
import React, {useState} from 'react';
import SidebarItem from '../SidebarItem';
import classNames from 'classnames';
import {colors} from '../../../../app/theme/default';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '../../hooks';

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
    padding: '60px 0px 36px 0px',
    position: 'relative',
    zIndex: 100,
    overflowX: 'hidden',
    width: '80px',
    transition: 'width 500ms',
  },
  expanded: {
    width: '208px',
  },
}));

type ItemConfig = {
  path: string,
  label: string,
  icon: any,
};

type Props = {
  items: Array<ItemConfig>,
  showNetworkSwitch?: boolean,
};

const AppSideBar = (props: Props) => {
  const {items, showNetworkSwitch} = props;
  const classes = useStyles();
  const [expanded, setIsExpanded] = useState(false);
  const [isProfileMenuOpen, _setProfileMenuOpen] = useState(false);
  const [isNetworkMenuOpen, _setNetworkMenuOpen] = useState(false);
  const {relativeUrl} = useRouter();

  const setProfileMenuOpen = (isOpen: boolean) => {
    if (!isOpen) {
      setIsExpanded(false);
    }
    _setProfileMenuOpen(isOpen);
  };
  const setNetworkMenuOpen = (isOpen: boolean) => {
    if (!isOpen) {
      setIsExpanded(false);
    }
    _setNetworkMenuOpen(isOpen);
  };

  return (
    <div
      data-testid="app-sidebar"
      className={classes.root}
      onMouseOver={() => setIsExpanded(true)}
      onMouseLeave={() => {
        if (!isProfileMenuOpen && !isNetworkMenuOpen) {
          setIsExpanded(false);
        }
      }}>
      <div
        className={classNames({
          [classes.inner]: true,
          [classes.expanded]: expanded,
        })}>
        <div className={classes.mainItems}>
          {items.map(({path, label, icon}) => (
            <SidebarItem
              key={label}
              path={relativeUrl(path)}
              label={label}
              icon={icon}
              expanded={expanded}
            />
          ))}
        </div>
        <div className={classes.secondaryItems}>
          {showNetworkSwitch && (
            <NetworkSelector
              expanded={expanded}
              isMenuOpen={isNetworkMenuOpen}
              setMenuOpen={setNetworkMenuOpen}
            />
          )}
          <ProfileButton
            expanded={expanded}
            isMenuOpen={isProfileMenuOpen}
            setMenuOpen={setProfileMenuOpen}
          />
        </div>
      </div>
    </div>
  );
};

export default AppSideBar;
