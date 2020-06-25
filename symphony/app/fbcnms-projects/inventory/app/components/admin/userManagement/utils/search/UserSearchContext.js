/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 *
 * @flow strict-local
 * @format
 */

import type {User} from '../UserManagementUtils';
import type {UserSearchContextQuery} from './__generated__/UserSearchContextQuery.graphql';

import RelayEnvironment from '../../../../../common/RelayEnvironment';
import createSearchContext from './SearchContext';
import {USER_ROLES, USER_STATUSES} from '../UserManagementUtils';
import {fetchQuery, graphql} from 'relay-runtime';

const userSearchQuery = graphql`
  query UserSearchContextQuery($filters: [UserFilterInput!]!) {
    userSearch(filters: $filters) {
      users {
        ...UserManagementUtils_user @relay(mask: false)
      }
    }
  }
`;

const searchCallback = (searchTerm: string) =>
  fetchQuery<UserSearchContextQuery>(RelayEnvironment, userSearchQuery, {
    filters: [
      {
        filterType: 'USER_NAME',
        operator: 'CONTAINS',
        stringValue: searchTerm,
      },
      {
        filterType: 'USER_STATUS',
        operator: 'IS',
        statusValue: USER_STATUSES.ACTIVE.key,
      },
    ],
  }).then(response => {
    if (response?.userSearch == null) {
      return [];
    }
    return response.userSearch.users.filter(Boolean).map(user => {
      if (user.role === USER_ROLES.OWNER.key) {
        return {
          ...user,
          role: USER_ROLES.ADMIN.key,
        };
      }
      return user;
    });
  });

const {
  SearchContext: UserSearchContext,
  SearchContextProvider,
  useSearchContext,
  useSearch,
} = createSearchContext<User>(searchCallback);

export const UserSearchContextProvider = SearchContextProvider;
export const useUserSearchContext = useSearchContext;
export const useUserSearch = useSearch;
export default UserSearchContext;
