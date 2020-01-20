/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {EditUser} from '@fbcnms/ui/components/auth/EditUserDialog';
import type {Suggestion} from '@fbcnms/ui/components/Typeahead';

import * as React from 'react';
import Typeahead from '@fbcnms/ui/components/Typeahead';
import axios from 'axios';

type Props = {
  className?: string,
  required?: boolean,
  headline?: string,
  selectedUser?: ?string,
  margin?: ?string,
  onUserSelection: (projectId: ?string) => void,
};

type State = {
  suggestions: Array<Suggestion>,
  users: Array<EditUser>,
};

class UserTypeahead extends React.Component<Props, State> {
  state = {
    suggestions: [],
    users: [],
  };
  componentDidMount() {
    axios
      .get('/user/list/')
      .then(response => this.setState({users: response.data.users}));
  }

  fetchNewUserSuggestions = (searchTerm: string) => {
    const searchTermLC = searchTerm.toLowerCase();
    const users = this.state.users;
    const userEmails = users
      .map(e => e.email)
      .filter(e => e.toLowerCase().includes(searchTermLC));
    const suggestions = userEmails.map((e, i) => ({
      name: e,
      entityId: String(i),
      entityType: '',
      type: '',
    }));

    this.setState({
      suggestions,
    });
  };
  render() {
    const {selectedUser, headline, required, className} = this.props;
    const {suggestions} = this.state;
    return (
      <div className={className}>
        <Typeahead
          margin={this.props.margin}
          required={!!required}
          suggestions={suggestions}
          onSuggestionsFetchRequested={this.fetchNewUserSuggestions}
          onEntitySelected={suggestion =>
            this.props.onUserSelection(suggestion.name)
          }
          onEntriesRequested={() => {}}
          onSuggestionsClearRequested={() => this.props.onUserSelection('')}
          placeholder={headline}
          value={
            selectedUser
              ? {
                  name: selectedUser,
                  entityId: '1',
                  entityType: '',
                  type: '',
                }
              : null
          }
        />
      </div>
    );
  }
}

export default UserTypeahead;
