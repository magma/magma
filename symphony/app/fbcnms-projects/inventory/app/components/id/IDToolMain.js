/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import AppContent from '@fbcnms/ui/components/layout/AppContent';
import ApplicationMain from '@fbcnms/ui/components/ApplicationMain';
import EntDetails from './EntDetails';
import EntSearchBar from './EntSearchBar';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import {AppContextProvider} from '@fbcnms/ui/context/AppContext';
import {Route} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  root: {
    display: 'flex',
  },
  content: {
    height: '100vh',
    display: 'flex',
    flexDirection: 'column',
  },
  placeholderContainer: {
    display: 'flex',
    flexGrow: 1,
    alignItems: 'center',
    justifyContent: 'center',
    color: theme.palette.blueGrayDark,
  },
}));

function IDToolMain() {
  const classes = useStyles();

  return (
    <div className={classes.root}>
      <AppContent>
        <div className={classes.content}>
          <EntSearchBar />
          <Route path="/id/:id" component={EntDetails} />
          <Route
            path="/id/"
            exact
            render={() => (
              <div className={classes.placeholderContainer}>
                <Text>Enter an ID above to see its fields and edges</Text>
              </div>
            )}
          />
        </div>
      </AppContent>
    </div>
  );
}

export default () => {
  return (
    <ApplicationMain>
      <AppContextProvider>
        <IDToolMain />
      </AppContextProvider>
    </ApplicationMain>
  );
};
