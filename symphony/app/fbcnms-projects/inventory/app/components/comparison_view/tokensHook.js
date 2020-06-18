/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {Entry} from '@fbcnms/ui/components/Tokenizer';
import type {FilterValue} from './ComparisonViewTypes';

import WizardContext from '@fbcnms/ui/components/design-system/Wizard/WizardContext';
import {useContext} from 'react';

export function useTokens(value: FilterValue): Array<Entry> {
  const idSet = value.idSet ?? [];
  const wizardContext = useContext(WizardContext);
  return (
    idSet
      // eslint-disable-next-line no-warning-comments
      // $FlowFixMe
      .map(id => wizardContext.get(id))
      .filter(Boolean)
  );
}
