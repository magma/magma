/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ChecklistCategoriesStateType} from '../ChecklistCategoriesMutateState';

import * as React from 'react';
import AddIcon from '@fbcnms/ui/components/design-system/Icons/Actions/AddIcon';
import AppContext from '@fbcnms/ui/context/AppContext';
import Button from '@fbcnms/ui/components/design-system/Button';
import CheckListCategoryTable from './CheckListCategoryTable';
import ChecklistCategoriesMutateDispatchContext from '../ChecklistCategoriesMutateDispatchContext';
import ExpandingPanel from '@fbcnms/ui/components/ExpandingPanel';
import fbt from 'fbt';
import {useContext, useMemo} from 'react';
import {useFormContext} from '../../../common/FormContext';

type Props = $ReadOnly<{|
  categories: ChecklistCategoriesStateType,
  isDefinitionsOnly?: boolean,
|}>;

const CheckListCategoryExpandingPanel = ({
  categories,
  isDefinitionsOnly,
}: Props) => {
  const appContext = useContext(AppContext);
  const dispatch = useContext(ChecklistCategoriesMutateDispatchContext);
  const form = useFormContext();
  const categoriesEnabled = useMemo(
    () => appContext.isFeatureEnabled('checklistcategories'),
    [appContext],
  );
  if (!categoriesEnabled) {
    return null;
  }
  const hasCheckListCategories = categories.length > 0;
  return (
    <ExpandingPanel
      allowExpandCollapse={hasCheckListCategories}
      title={fbt('Checklist Categories', 'Checklist section header')}
      rightContent={
        form.alerts.missingPermissions.detected ? null : (
          <Button
            variant="text"
            disabled={form.alerts.missingPermissions.detected}
            onClick={() => dispatch({type: 'ADD_CATEGORY'})}>
            <AddIcon color="primary" />
          </Button>
        )
      }>
      <CheckListCategoryTable
        categories={categories}
        isDefinitionsOnly={isDefinitionsOnly}
      />
    </ExpandingPanel>
  );
};

export default CheckListCategoryExpandingPanel;
