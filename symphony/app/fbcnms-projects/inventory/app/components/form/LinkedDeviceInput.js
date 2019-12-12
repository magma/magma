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
import {makeStyles} from '@material-ui/styles';

type Props = {
  label: string,
  value: string,
  onChange: (SyntheticInputEvent<HTMLInputElement>) => void,
};

const useStyles = makeStyles(() => ({
  input: {
    display: 'inline-flex',
    width: '100%',
  },
}));

export default function LinkedDeviceInput(props: Props) {
  const classes = useStyles();
  const {label, value, onChange} = props;
  return (
    <FormField label={label} hasSpacer={true}>
      <TextInput
        type="string"
        className={classes.input}
        value={value}
        onChange={onChange}
      />
    </FormField>
  );
}
