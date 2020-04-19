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
import CloseIcon from '@fbcnms/ui/components/design-system/Icons/Navigation/CloseIcon';
import InputAffix from '@fbcnms/ui/components/design-system/Input/InputAffix';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import fbt from 'fbt';
import {useUserSearchContext} from './UserSearchContext';

type Props = $ReadOnly<{|
  className?: ?String,
|}>;

const UserSearchBox = (props: Props) => {
  const {className} = props;

  const userSearch = useUserSearchContext();

  return (
    <div className={className}>
      <TextInput
        type="string"
        variant="outlined"
        placeholder={`${fbt('Search users...', '')}`}
        isProcessing={userSearch.isSearchInProgress}
        fullWidth={true}
        value={userSearch.searchTerm}
        onChange={e => userSearch.setSearchTerm(e.target.value)}
        suffix={
          userSearch.isEmptySearchTerm ? null : (
            <InputAffix onClick={userSearch.clearSearch}>
              <CloseIcon color="gray" />
            </InputAffix>
          )
        }
      />
    </div>
  );
};

export default UserSearchBox;
