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
import React from 'react';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import fbt from 'fbt';
import {createFragmentContainer, graphql} from 'react-relay';
import {makeStyles} from '@material-ui/styles';
import type {FreeTextCheckListItemDefinition_item} from './__generated__/FreeTextCheckListItemDefinition_item.graphql';

type Props = {
  item: FreeTextCheckListItemDefinition_item,
  onChange?: (
    updatedChecklistItem: FreeTextCheckListItemDefinition_item,
  ) => void,
};

const useStyles = makeStyles(() => ({
  container: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
  },
  expandingPart: {
    flexGrow: 1,
    flexBasis: 0,
    '&:not(:first-child)': {
      marginLeft: '8px',
    },
    '&:not(:last-child)': {
      marginRight: '8px',
    },
  },
}));

const FreeTextCheckListItemDefinition = (props: Props) => {
  const {item, onChange} = props;
  const classes = useStyles();

  const _updateItemValue = (
    updatedItem: FreeTextCheckListItemDefinition_item,
  ) => {
    if (!onChange) {
      return;
    }
    onChange(updatedItem);
  };

  const _updateTitle: (
    newTitle: string,
  ) => FreeTextCheckListItemDefinition_item = newTitle => {
    return {
      ...item,
      title: newTitle,
    };
  };
  const _updateHelpText: (
    newTitle: string,
  ) => FreeTextCheckListItemDefinition_item = newHelpText => {
    return {
      ...item,
      helpText: newHelpText,
    };
  };

  return (
    <div className={classes.container}>
      <FormField className={classes.expandingPart}>
        <TextInput
          type="string"
          placeholder={fbt(
            'What text should be filled?',
            'Place holder for free text field title (user needs to type the title of the textbox in this field).',
          )}
          value={item.title || ''}
          onChange={changeEvent =>
            _updateItemValue(_updateTitle(changeEvent.target.value))
          }
        />
      </FormField>
      <FormField className={classes.expandingPart}>
        <TextInput
          type="string"
          placeholder={fbt(
            'Optional hint text',
            'Place holder for free text field place holder (user needs to type the placeholder used for this text box in this field).',
          )}
          value={item.helpText || ''}
          onChange={changeEvent =>
            _updateItemValue(_updateHelpText(changeEvent.target.value))
          }
        />
      </FormField>
    </div>
  );
};

export default createFragmentContainer(FreeTextCheckListItemDefinition, {
  item: graphql`
    fragment FreeTextCheckListItemDefinition_item on CheckListItem {
      title
      helpText
      ...CheckListItem_item
    }
  `,
});
