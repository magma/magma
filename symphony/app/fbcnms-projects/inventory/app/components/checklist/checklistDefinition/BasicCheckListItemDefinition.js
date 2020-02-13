/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import React, {useCallback} from 'react';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import fbt from 'fbt';
import {createFragmentContainer, graphql} from 'react-relay';
import {makeStyles} from '@material-ui/styles';
import type {BasicCheckListItemDefinition_item} from './__generated__/BasicCheckListItemDefinition_item.graphql';

type Props = {
  item: BasicCheckListItemDefinition_item,
  onChange: (updatedChecklistItem: BasicCheckListItemDefinition_item) => void,
};

const useStyles = makeStyles(() => ({
  container: {
    display: 'flex',
    flexDirection: 'row',
  },
  expandingPart: {
    flexGrow: 1,
    flexBasis: 0,
  },
}));

const BasicCheckListItemDefinition = (props: Props) => {
  const {item, onChange} = props;
  const classes = useStyles();

  const _updateOnChange = useCallback(
    newTitle => {
      const newItem = {
        ...item,
        title: newTitle,
      };
      onChange(newItem);
    },
    [item, onChange],
  );

  return (
    <div className={classes.container}>
      <FormField className={classes.expandingPart}>
        <TextInput
          type="string"
          placeholder={fbt(
            'What needs to be done?',
            'Placeholder for checkbox field title (user needs to type the title of the checkbox in this field).',
          )}
          value={item.title || ''}
          onChange={event => _updateOnChange(event.target.value)}
        />
      </FormField>
    </div>
  );
};

export default createFragmentContainer(BasicCheckListItemDefinition, {
  item: graphql`
    fragment BasicCheckListItemDefinition_item on CheckListItem {
      title
      checked
      ...CheckListItem_item
    }
  `,
});
