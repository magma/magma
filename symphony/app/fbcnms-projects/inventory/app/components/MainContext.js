/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {
  MainContextMeQuery,
  MainContextMeQueryResponse,
} from './__generated__/MainContextMeQuery.graphql';
import type {SessionUser} from '@fbcnms/magmalte/app/common/UserModel';

import * as React from 'react';
import RelayEnvironment from '../common/RelayEnvironment';
import {DEACTIVATED_PAGE_PATH} from './DeactivatedPage';
import {PermissionValues} from './admin/userManagement/utils/UserManagementUtils';
import {fetchQuery, graphql} from 'relay-runtime';
import {useContext, useEffect, useState} from 'react';
import {useLocation} from 'react-router-dom';

export type MainContextValue = {
  initializing: boolean,
  integrationUserDefinition: SessionUser,
  ...MainContextMeQueryResponse,
};

const integrationUserDefinitionBuilder: (
  ?MainContextMeQueryResponse,
) => SessionUser = queryResponse => ({
  email: queryResponse?.me?.user?.email || '',
  isSuperUser:
    queryResponse?.me?.permissions.adminPolicy.access.isAllowed ===
    PermissionValues.YES,
});

const DEFUALT_VALUE = {
  initializing: true,
  integrationUserDefinition: integrationUserDefinitionBuilder(),
  me: null,
};

const MainContext = React.createContext<MainContextValue>(DEFUALT_VALUE);

export function useMainContext() {
  return useContext(MainContext);
}

const meQuery = graphql`
  query MainContextMeQuery {
    me {
      user {
        id
        authID
        email
        firstName
        lastName
      }
      permissions {
        canWrite
        adminPolicy {
          access {
            isAllowed
          }
        }
      }
    }
  }
`;

const getLoggedUserSettings = () => {
  return fetchQuery<MainContextMeQuery>(RelayEnvironment, meQuery, {});
};
type Props = $ReadOnly<{|
  children: React.Node,
|}>;

export function MainContextProvider(props: Props) {
  const [value, setValue] = useState(DEFUALT_VALUE);
  const location = useLocation();
  useEffect(() => {
    if (location.pathname === DEACTIVATED_PAGE_PATH) {
      setValue(currentValue => ({
        ...currentValue,
        initializing: false,
      }));
      return;
    }

    getLoggedUserSettings()
      .then(meValue =>
        setValue(currentValue => ({
          ...currentValue,
          integrationUserDefinition: integrationUserDefinitionBuilder(meValue),
          ...meValue,
        })),
      )
      .finally(() =>
        setValue(currentValue => ({
          ...currentValue,
          initializing: false,
        })),
      );
  }, [location.pathname]);
  return (
    <MainContext.Provider value={value}>{props.children}</MainContext.Provider>
  );
}

export default MainContext;
