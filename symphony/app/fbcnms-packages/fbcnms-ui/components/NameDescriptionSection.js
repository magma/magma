/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import FormField from './design-system/FormField/FormField';
import NameInput from '@fbcnms/ui/components/design-system/Form/NameInput';
import React from 'react';
import Text from './design-system/Text';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import classNames from 'classnames';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  fieldName: {
    fontSize: '14px',
    lineHeight: '24px',
    color: theme.palette.blueGrayDark,
    marginBottom: '4px',
    display: 'block',
  },
  nameField: {
    width: '50%',
  },
  descriptionField: {
    '&&': {
      padding: '6px 8px',
      lineHeight: '20px',
      fontSize: '14px',
    },
  },
  descriptionTitle: {
    marginTop: '20px',
  },
  inputMultiline: {
    padding: '6px 8px',
    lineHeight: '14px',
    fontSize: '14px',
    minHeight: 'inherit',
    height: 'auto',
    '&::placeholder': {
      fontSize: '14px',
      color: '#8895ad',
    },
  },
}));

type Props = {
  title?: string,
  name?: string,
  namePlaceholder?: ?string,
  description?: ?string,
  descriptionPlaceholder?: string,
  onNameChange?: string => void,
  onDescriptionChange?: string => void,
};

const NameDescriptionSection = ({
  title,
  name,
  namePlaceholder,
  description,
  descriptionPlaceholder,
  onNameChange,
  onDescriptionChange,
}: Props) => {
  const classes = useStyles();
  return (
    <>
      <NameInput
        value={name}
        onChange={event => onNameChange && onNameChange(event.target.value)}
        inputClass={classes.nameField}
        title={title}
        placeholder={namePlaceholder || ''}
      />
      <Text className={classNames(classes.fieldName, classes.descriptionTitle)}>
        Description
      </Text>
      <FormField>
        <TextInput
          type="multiline"
          placeholder={descriptionPlaceholder}
          rows={4}
          value={description ?? ''}
          onChange={event =>
            onDescriptionChange && onDescriptionChange(event.target.value)
          }
        />
      </FormField>
    </>
  );
};

export default NameDescriptionSection;
