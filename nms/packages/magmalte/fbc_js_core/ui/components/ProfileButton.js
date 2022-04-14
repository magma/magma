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

import AppContext from '../../../fbc_js_core/ui/context/AppContext';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import NetworkContext from '../../../app/components/context/NetworkContext';
import Popout from '../../../fbc_js_core/ui/components/Popout';
import ProfileIcon from '../icons/ProfileIcon';
import React, {useContext, useState} from 'react';
import Text from './design-system/Text';
import classNames from 'classnames';
import {Events, GeneralLogger} from '../../../fbc_js_core/ui/utils/Logging';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '../hooks';

const useStyles = makeStyles(theme => ({
  accountButton: {
    backgroundColor: theme.palette.common.white,
    width: '36px',
    height: '36px',
    fontSize: '36px',
    cursor: 'pointer',
    borderRadius: '100%',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    '&:hover, &$openButton': {
      '& $accountButtonIcon': {
        fill: theme.palette.primary.main,
      },
    },
  },
  openButton: {},
  accountButtonIcon: {
    '&&': {
      fill: theme.palette.blueGrayDark,
      fontSize: '19px',
    },
  },
  accountButtonIconSelected: {
    '&&': {
      fill: theme.palette.primary.main,
    },
  },
  itemGutters: {
    '&&': {
      minWidth: '170px',
      borderRadius: '4px',
      padding: '8px 10px',
      '&:hover': {
        backgroundColor: 'rgba(145, 145, 145, 0.1)',
      },
    },
  },
  profileList: {
    '&&': {
      padding: '10px 5px',
    },
  },
  profileItemText: {
    fontSize: '12px',
    lineHeight: '16px',
  },
  selected: {
    backgroundColor: theme.palette.primary.main,
  },
  sidebarEntry: {
    display: 'flex',
    padding: '9px',
    justifyContent: 'center',
    width: '100%',
  },
}));

const ProfileButton = () => {
  const {relativeUrl, history, location} = useRouter();
  const classes = useStyles();
  const [isProfileMenuOpen, toggleProfileMenu] = useState(false);
  const {networkId: selectedNetworkId} = useContext(NetworkContext);
  const {user, ssoEnabled, isFeatureEnabled, isOrganizations} = useContext(
    AppContext,
  );
  const {email} = user;

  const getUrl = (path: string) =>
    (selectedNetworkId != undefined || isOrganizations) ? relativeUrl(path) : path;

  const adminUrl = getUrl('/admin');
  const settingsUrl = getUrl('/settings');

  const isSelected =
    location.pathname.includes(adminUrl) ||
    location.pathname.includes(settingsUrl);

  return (
    <Popout
      className={classNames({
        [classes.sidebarEntry]: true,
        [classes.selected]: isSelected,
      })}
      open={isProfileMenuOpen}
      content={
        <List component="nav" className={classes.profileList}>
          <ListItem classes={{gutters: classes.itemGutters}} disabled={true}>
            <Text className={classes.profileItemText}>{email}</Text>
          </ListItem>
          {!ssoEnabled && (
            <ListItem
              classes={{gutters: classes.itemGutters}}
              button
              onClick={() => {
                GeneralLogger.info(Events.SETTINGS_CLICKED);
                toggleProfileMenu(false);
                history.push(settingsUrl);
              }}
              component="a">
              <Text className={classes.profileItemText}>Account Settings</Text>
            </ListItem>
          )}
          {user.isSuperUser && !isOrganizations && (
            <ListItem
              classes={{gutters: classes.itemGutters}}
              button
              onClick={() => {
                GeneralLogger.info(Events.ADMINISTRATION_CLICKED);
                toggleProfileMenu(false);
                history.push(adminUrl);
              }}
              component="a">
              <Text className={classes.profileItemText}>Administration</Text>
            </ListItem>
          )}
          {isFeatureEnabled('documents_site') && (
            <ListItem
              classes={{gutters: classes.itemGutters}}
              button
              href={'/docs/docs/inventory-intro.html'}
              onClick={() =>
                GeneralLogger.info(Events.DOCUMENTATION_LINK_CLICKED)
              }
              component="a">
              <Text className={classes.profileItemText}>Documentation</Text>
            </ListItem>
          )}
          <ListItem
            classes={{gutters: classes.itemGutters}}
            button
            href="/user/logout"
            component="a">
            <Text className={classes.profileItemText}>Logout</Text>
          </ListItem>
        </List>
      }
      onOpen={() => toggleProfileMenu(true)}
      onClose={() => toggleProfileMenu(false)}>
      <div
        data-testid="profileButton"
        className={classNames({
          [classes.accountButton]: true,
          [classes.openButton]: isProfileMenuOpen,
        })}>
        <ProfileIcon
          className={classNames({
            [classes.accountButtonIcon]: true,
            [classes.accountButtonIconSelected]: isSelected,
          })}
        />
      </div>
    </Popout>
  );
};

export default ProfileButton;
