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
import {useGroupSearchContext} from './GroupSearchContext';

type Props = $ReadOnly<{|
  className?: ?String,
|}>;

const GroupSearchBox = (props: Props) => {
  const {className} = props;

  const groupSearch = useGroupSearchContext();

  return (
    <div className={className}>
      <TextInput
        variant="outlined"
        placeholder={`${fbt('Search groups...', '')}`}
        isProcessing={groupSearch.isSearchInProgress}
        fullWidth={true}
        value={groupSearch.searchTerm}
        onChange={e => groupSearch.setSearchTerm(e.target.value)}
        onEscPressed={() => groupSearch.clearSearch()}
        suffix={
          groupSearch.isEmptySearchTerm ? null : (
            <InputAffix onClick={groupSearch.clearSearch}>
              <CloseIcon color="gray" />
            </InputAffix>
          )
        }
      />
    </div>
  );
};

export default GroupSearchBox;
