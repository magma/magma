/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import FormValidationContext from '@fbcnms/ui/components/design-system/Form/FormValidationContext';
import NameInput from '@fbcnms/ui/components/design-system/Form/NameInput';
import React, {useContext} from 'react';
import Text from './design-system/Text';
import TextField from '@material-ui/core/TextField';
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
  const validationContext = useContext(FormValidationContext);
  return (
    <>
      <NameInput
        value={name}
        onChange={event => onNameChange && onNameChange(event.target.value)}
        inputClass={classes.nameField}
        title={title}
        placeholder={namePlaceholder || ''}
        disabled={validationContext.editLock.detected}
      />
      <Text className={classNames(classes.fieldName, classes.descriptionTitle)}>
        Description
      </Text>
      <TextField
        name="Description"
        InputProps={{
          classes: {
            root: classes.descriptionField,
            inputMultiline: classes.inputMultiline,
          },
        }}
        disabled={validationContext.editLock.detected}
        placeholder={descriptionPlaceholder}
        variant="outlined"
        multiline
        fullWidth
        rows="4"
        value={description ?? ''}
        onChange={event =>
          onDescriptionChange && onDescriptionChange(event.target.value)
        }
      />
    </>
  );
};

export default NameDescriptionSection;
