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

import type {GroupSearchContextQuery} from './__generated__/GroupSearchContextQuery.graphql';
import type {
  PermissionsPolicy,
  UserPermissionsGroup,
} from '../UserManagementUtils';

import RelayEnvironment from '../../../../../common/RelayEnvironment';
import createSearchContext from './SearchContext';
import {fetchQuery, graphql} from 'relay-runtime';
import {groupsResponse2Groups} from '../UserManagementUtils';

const groupSearchQuery = graphql`
  query GroupSearchContextQuery($filters: [UsersGroupFilterInput!]!) {
    usersGroupSearch(filters: $filters) {
      usersGroups {
        id
        name
        description
        status
        members {
          id
          authID
        }
      }
    }
  }
`;

const searchCallback = (searchTerm: string, _policy: ?PermissionsPolicy) =>
  fetchQuery<GroupSearchContextQuery>(RelayEnvironment, groupSearchQuery, {
    filters: [
      {
        filterType: 'GROUP_NAME',
        operator: 'CONTAINS',
        stringValue: searchTerm,
      },
    ],
  }).then(response => groupsResponse2Groups(response.usersGroupSearch));

const {
  SearchContext: GroupSearchContext,
  SearchContextProvider,
  useSearchContext,
  useSearch,
} = createSearchContext<UserPermissionsGroup>(searchCallback);

export const GroupSearchContextProvider = SearchContextProvider;
export const useGroupSearchContext = useSearchContext;
export const useGroupSearch = useSearch;
export default GroupSearchContext;
