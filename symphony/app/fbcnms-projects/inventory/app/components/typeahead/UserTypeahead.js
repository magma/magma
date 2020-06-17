/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ShortUser} from '../../common/EntUtils';
import type {Suggestion} from '@fbcnms/ui/components/Typeahead';
import type {UserTypeahead_userQueryResponse} from './__generated__/UserTypeahead_userQuery.graphql';

import * as React from 'react';
import RelayEnvironment from '../../common/RelayEnvironment.js';
import Typeahead from '@fbcnms/ui/components/Typeahead';
import {debounce} from 'lodash';
import {fetchQuery, graphql} from 'relay-runtime';

const userTypeaheadQuery = graphql`
  query UserTypeahead_userQuery($filters: [UserFilterInput!]!) {
    userSearch(limit: 10, filters: $filters) {
      users {
        id
        email
      }
    }
  }
`;

const USER_SEARCH_DEBOUNCE_TIMEOUT_MS = 200;
const DEBOUNCE_CONFIG = {
  trailing: true,
  leading: true,
};

type Props = {
  className?: string,
  required?: boolean,
  headline?: string,
  selectedUser?: ?ShortUser,
  margin?: ?string,
  onUserSelection: (?ShortUser) => void,
};

type State = {
  userSuggestions: Array<Suggestion>,
};

class UserTypeahead extends React.Component<Props, State> {
  state = {
    userSuggestions: [],
  };

  debounceUserFetchSuggestions = debounce(
    (searchTerm: string) => this.fetchNewUserSuggestions(searchTerm),
    USER_SEARCH_DEBOUNCE_TIMEOUT_MS,
    DEBOUNCE_CONFIG,
  );

  fetchNewUserSuggestions = (searchTerm: string) => {
    fetchQuery(RelayEnvironment, userTypeaheadQuery, {
      filters: [
        {
          filterType: 'USER_NAME',
          operator: 'CONTAINS',
          stringValue: searchTerm,
        },
      ],
    }).then((response: ?UserTypeahead_userQueryResponse) => {
      if (!response || !response.userSearch) {
        return;
      }
      this.setState({
        userSuggestions: response.userSearch.users.filter(Boolean).map(e => ({
          name: e.email,
          entityId: e.id,
          entityType: '',
          type: 'user',
        })),
      });
    });
  };

  onUserSuggestionsFetchRequested = (searchTerm: string) => {
    this.debounceUserFetchSuggestions(searchTerm);
  };

  render() {
    const {
      selectedUser,
      headline,
      required,
      className,
      onUserSelection,
    } = this.props;
    const {userSuggestions} = this.state;
    return (
      <div className={className}>
        <Typeahead
          margin={this.props.margin}
          required={!!required}
          suggestions={userSuggestions}
          onSuggestionsFetchRequested={this.onUserSuggestionsFetchRequested}
          onEntitySelected={suggestion =>
            onUserSelection({id: suggestion.entityId, email: suggestion.name})
          }
          onEntriesRequested={() => {}}
          onSuggestionsClearRequested={() => onUserSelection(null)}
          placeholder={headline}
          value={
            selectedUser
              ? {
                  name: selectedUser.email,
                  entityId: selectedUser.id,
                  entityType: '',
                  type: 'user',
                }
              : null
          }
        />
      </div>
    );
  }
}

export default UserTypeahead;
