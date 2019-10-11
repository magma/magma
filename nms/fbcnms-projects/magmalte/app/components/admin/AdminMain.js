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
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '../../common/MagmaV1API';

import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '../../common/useMagmaAPI';
import {getProjectLinks} from '../../common/projects';
import {makeStyles} from '@material-ui/styles';
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
  const {error, isLoading, response} = useMagmaAPI(MagmaV1API.getNetworks, {});

  if (isLoading) {
    return <LoadingFiller />;
  }

  const networkIds = error || !response ? ['mpk_test'] : response.sort();
  const appContext = {
    ...window.CONFIG.appData,
    networkIds,
  };
  return (
    <AppContext.Provider value={appContext}>
      <AdminMain {...props} />
    </AppContext.Provider>
  );
};
