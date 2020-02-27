/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {CheckListCategoryExpandingPanel_list} from './__generated__/CheckListCategoryExpandingPanel_list.graphql';

import * as React from 'react';
import AddIcon from '@fbcnms/ui/components/design-system/Icons/Actions/AddIcon';
import AppContext from '@fbcnms/ui/context/AppContext';
import Button from '@fbcnms/ui/components/design-system/Button';
import CheckListCategoryContext, {
  CheckListCategoryContextProvider,
} from './CheckListCategoryContext';
import CheckListCategoryTable from './CheckListCategoryTable';
import ExpandingPanel from '@fbcnms/ui/components/ExpandingPanel';
import FormValidationContext from '@fbcnms/ui/components/design-system/Form/FormValidationContext';
import fbt from 'fbt';
import {createFragmentContainer, graphql} from 'react-relay';
import {useContext, useMemo} from 'react';

type Props = {
  list: CheckListCategoryExpandingPanel_list,
  onListChanged?: (updatedList: CheckListCategoryExpandingPanel_list) => void,
};

const CheckListCategoryExpandingPanel = (props: Props) => {
  const appContext = useContext(AppContext);
  const formValidationContext = useContext(FormValidationContext);
  const categoriesEnabled = useMemo(
    () => appContext.isFeatureEnabled('checklistcategories'),
    [appContext],
  );
  if (!categoriesEnabled) {
    return null;
  }
  const hasCheckListCategories = props.list.length > 0;
  return (
    <CheckListCategoryContextProvider>
      <CheckListCategoryContext.Consumer>
        {categoryContext => (
          <ExpandingPanel
            allowExpandCollapse={hasCheckListCategories}
            title={fbt('Checklist Categories', 'Checklist section header')}
            rightContent={
              formValidationContext.editLock.detected ? null : (
                <Button
                  variant="text"
                  disabled={formValidationContext.editLock.detected}
                  onClick={() => categoryContext.call.addNewCategory()}>
                  <AddIcon color="primary" />
                </Button>
              )
            }>
            <CheckListCategoryTable {...props} />
          </ExpandingPanel>
        )}
      </CheckListCategoryContext.Consumer>
    </CheckListCategoryContextProvider>
  );
};

export default createFragmentContainer(CheckListCategoryExpandingPanel, {
  list: graphql`
    fragment CheckListCategoryExpandingPanel_list on CheckListCategory
      @relay(plural: true) {
      ...CheckListCategoryTable_list
    }
  `,
});
