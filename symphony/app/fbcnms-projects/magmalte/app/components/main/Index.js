/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {Theme} from '@material-ui/core';

import AppContent from '@fbcnms/ui/components/layout/AppContent';
import AppContext from '@fbcnms/ui/context/AppContext';
import AppSideBar from '@fbcnms/ui/components/layout/AppSideBar';
import NetworkContext from '../context/NetworkContext';
import NetworkSelector from '../NetworkSelector';
import React, {useContext} from 'react';
import SectionLinks from '../layout/SectionLinks';
import SectionRoutes from '../layout/SectionRoutes';
import VersionTooltip from '../VersionTooltip';

import {getProjectLinks} from '../../common/projects';
import {makeStyles} from '@material-ui/styles';
import {shouldShowSettings} from '../Settings';
import {useRouter} from '@fbcnms/ui/hooks';

// These won't be considered networkIds
export const ROOT_PATHS = new Set<string>(['network']);

const useStyles = makeStyles((theme: Theme) => ({
  root: {
    display: 'flex',
  },
  toolbarIcon: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'flex-end',
    padding: '0 8px',
    ...theme.mixins.toolbar,
  },
}));

export default function Index() {
  const classes = useStyles();
  const {match} = useRouter();
  const {user, tabs, ssoEnabled} = useContext(AppContext);
  const networkId = ROOT_PATHS.has(match.params.networkId)
    ? null
    : match.params.networkId;

  return (
    <NetworkContext.Provider value={{networkId}}>
      <div className={classes.root}>
        <AppSideBar
          mainItems={[<SectionLinks key={1} />, <VersionTooltip key={2} />]}
          secondaryItems={[<NetworkSelector key={1} />]}
          projects={getProjectLinks(tabs, user)}
          showSettings={shouldShowSettings({
            isSuperUser: user.isSuperUser,
            ssoEnabled,
          })}
          user={user}
        />
        <AppContent>
          <SectionRoutes />
        </AppContent>
      </div>
    </NetworkContext.Provider>
  );
}
