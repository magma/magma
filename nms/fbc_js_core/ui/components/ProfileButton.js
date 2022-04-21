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

import AppContext from '../context/AppContext';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import NetworkContext from '../../../app/components/context/NetworkContext';
import Popout from './Popout';
import PersonIcon from '@material-ui/icons/Person';
import React, {useContext, useEffect, useState} from 'react';
import Text from '../../../app/theme/design-system/Text';
import classNames from 'classnames';
import {Events, GeneralLogger} from '../utils/Logging';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '../hooks';
import {colors} from '../../../app/theme/default';
import Divider from '@material-ui/core/Divider';

const useStyles = makeStyles(theme => ({
  button: {
    display: 'flex',
    alignItems: 'center',
    width: '100%',
    padding: '15px 28px',
    cursor: 'pointer',
    outline: 'none',
    '&:hover $icon, &:hover $label, &$selected $icon, &$selected $label': {
      color: colors.primary.white,
    },
  },
  label: {
    '&&': {
      color: colors.primary.gullGray,
      whiteSpace: 'nowrap',
      paddingLeft: '16px',
    },
  },
  selected: {
    backgroundColor: colors.secondary.dodgerBlue,

    '& $icon': {
      color: colors.primary.white,
    },
  },
  icon: {
    color: colors.primary.gullGray,
    display: 'flex',
    justifyContent: 'center',
  },
  itemGutters: {
    '&&': {
      minWidth: '200px',
      padding: '6px 17px',
      '&:hover': {
        backgroundColor: colors.primary.concrete,
      },
    },
  },
  divider: {
    margin: '6px 17px',
  },
  profileList: {
    '&&': {
      padding: '10px 0',
    },
  },
  profileItemText: {
    fontSize: '14px',
    lineHeight: '20px',
  },
}));

type Props = {
  isMenuOpen: boolean,
  setMenuOpen: (isOpen: boolean) => void,
  expanded: boolean,
};

const ProfileButton = (props: Props) => {
  const {relativeUrl, history, location} = useRouter();
  const classes = useStyles();
  const {networkId: selectedNetworkId} = useContext(NetworkContext);
  const {
    user,
    hasAccountSettings,
    isFeatureEnabled,
    isOrganizations,
  } = useContext(AppContext);

  const getUrl = (path: string) =>
    selectedNetworkId != undefined || isOrganizations
      ? relativeUrl(path)
      : path;

  const adminUrl = getUrl('/admin');
  const settingsUrl = getUrl('/settings');

  const isSelected =
    location.pathname.includes(adminUrl) ||
    location.pathname.includes(settingsUrl);

  const hasAdministration = user.isSuperUser && !isOrganizations;
  const hasDocumentation = isFeatureEnabled('documents_site');

  return (
    <Popout
      className={classNames({
        [classes.button]: true,
        [classes.selected]: isSelected,
      })}
      open={props.isMenuOpen}
      content={
        <List component="nav" className={classes.profileList}>
          <ListItem classes={{gutters: classes.itemGutters}} disabled={true}>
            <Text className={classes.profileItemText}>{user.email}</Text>
          </ListItem>
          <Divider className={classes.divider} />
          {hasAccountSettings && (
            <ListItem
              classes={{gutters: classes.itemGutters}}
              button
              onClick={() => {
                GeneralLogger.info(Events.SETTINGS_CLICKED);
                props.setMenuOpen(false);
                history.push(settingsUrl);
              }}
              component="a">
              <Text className={classes.profileItemText}>Account Settings</Text>
            </ListItem>
          )}
          {hasAdministration && (
            <ListItem
              classes={{gutters: classes.itemGutters}}
              button
              onClick={() => {
                GeneralLogger.info(Events.ADMINISTRATION_CLICKED);
                props.setMenuOpen(false);
                history.push(adminUrl);
              }}
              component="a">
              <Text className={classes.profileItemText}>Administration</Text>
            </ListItem>
          )}
          {hasDocumentation && (
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
          {(hasAccountSettings || hasAdministration || hasDocumentation) && (
            <Divider className={classes.divider} />
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
      onOpen={() => props.setMenuOpen(true)}
      onClose={() => props.setMenuOpen(false)}>
      <PersonIcon data-testid="profileButton" className={classes.icon} />
      {props.expanded && (
        <Text className={classes.label} variant="body3">
          Account & Settings
        </Text>
      )}
    </Popout>
  );
};

export default ProfileButton;
