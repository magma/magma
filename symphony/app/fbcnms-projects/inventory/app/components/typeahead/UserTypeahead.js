/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {GroupMember} from '../admin/userManagement/utils/GroupMemberViewer';
import type {ShortUser} from '../../common/EntUtils';

import * as React from 'react';
import GroupMemberViewer from '../admin/userManagement/utils/GroupMemberViewer';
import Typeahead from '@fbcnms/ui/components/Typeahead';
import {useUserSearch} from '../admin/userManagement/utils/userSearch/UserSearchContext.js';

type Props = {
  className?: string,
  required?: boolean,
  headline?: string,
  selectedUser?: ?ShortUser,
  margin?: ?string,
  onUserSelection: (?ShortUser) => void,
};

const UserTypeahead = (props: Props) => {
  const {
    selectedUser,
    headline,
    required,
    className,
    margin,
    onUserSelection,
  } = props;

  const userSearch = useUserSearch();
  return (
    <div className={className}>
      <Typeahead
        margin={margin}
        required={!!required}
        // eslint-disable-next-line no-warning-comments
        // $FlowFixMe
        suggestions={userSearch.results.map(result => {
          // eslint-disable-next-line no-warning-comments
          // $FlowFixMe
          const member: GroupMember = result;
          return {
            entityId: member.user.id,
            entityType: 'user',
            name: member.user.authID,
            type: 'user',
            render: () => <GroupMemberViewer member={member} />,
          };
        })}
        onSuggestionsFetchRequested={userSearch.setSearchTerm}
        onEntitySelected={suggestion =>
          onUserSelection({
            id: suggestion.entityId,
            email: suggestion.name,
          })
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
};

export default UserTypeahead;
