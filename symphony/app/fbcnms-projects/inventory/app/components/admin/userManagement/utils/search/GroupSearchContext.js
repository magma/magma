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

import RelayEnvironment from '../../../../../common/RelayEnvironment';
import createSearchContext from './SearchContext';
import {fetchQuery, graphql} from 'relay-runtime';
import type {GroupSearchContextQuery} from './__generated__/GroupSearchContextQuery.graphql';
import type {UsersGroup} from '../../data/UsersGroups';

const groupSearchQuery = graphql`
  query GroupSearchContextQuery($filters: [UsersGroupFilterInput!]!) {
    usersGroups(first: 500, filterBy: $filters) {
      edges {
        node {
          ...UserManagementUtils_group @relay(mask: false)
        }
      }
    }
  }
`;

const searchCallback = (searchTerm: string) =>
  fetchQuery<GroupSearchContextQuery>(RelayEnvironment, groupSearchQuery, {
    filters: [
      {
        filterType: 'GROUP_NAME',
        operator: 'CONTAINS',
        stringValue: searchTerm,
      },
    ],
  }).then(response => {
    if (response?.usersGroups == null) {
      return [];
    }
    return response.usersGroups.edges.map(edge => edge.node).filter(Boolean);
  });

const {
  SearchContext: GroupSearchContext,
  SearchContextProvider,
  useSearchContext,
  useSearch,
} = createSearchContext<UsersGroup>(searchCallback);

export const GroupSearchContextProvider = SearchContextProvider;
export const useGroupSearchContext = useSearchContext;
export const useGroupSearch = useSearch;
export default GroupSearchContext;
