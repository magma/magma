/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import * as React from 'react';
import AppContent from '@fbcnms/ui/components/layout/AppContent';
import AppContext from '@fbcnms/ui/context/AppContext';
import AppSideBar from '@fbcnms/ui/components/layout/AppSideBar.react';
import ApplicationMain from '@fbcnms/ui/components/ApplicationMain';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';

import nullthrows from '@fbcnms/util/nullthrows';
import {MagmaAPIUrls} from '../../common/MagmaAPI';
import {getProjectLinks} from '../../common/projects';
import {makeStyles} from '@material-ui/styles';
import {useAxios} from '@fbcnms/ui/hooks';
import {useContext} from 'react';

const useStyles = makeStyles(_theme => ({
  root: {
    display: 'flex',
  },
}));

type Props = {
  navItems: () => React.Node,
  navRoutes: () => React.Node,
};

function AdminMain(props: Props) {
  const classes = useStyles();
  const {tabs, user} = useContext(AppContext);

  return (
    <div className={classes.root}>
      <AppSideBar
        mainItems={props.navItems()}
        projects={getProjectLinks(tabs, user)}
        user={nullthrows(user)}
      />
      <AppContent>{props.navRoutes()}</AppContent>
    </div>
  );
}

export default (props: Props) => {
  const {error, isLoading, response} = useAxios({
    method: 'get',
    url: MagmaAPIUrls.networks(),
  });

  if (isLoading) {
    return <LoadingFiller />;
  }

  const networkIds = error || !response ? ['mpk_test'] : response.data.sort();

  const appContext = {
    ...window.CONFIG.appData,
    networkIds,
  };
  return (
    <ApplicationMain appContext={appContext}>
      <AdminMain {...props} />
    </ApplicationMain>
  );
};
