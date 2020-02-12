/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import * as React from 'react';
import ArrowBackIcon from '@material-ui/icons/ArrowBack';
import EditIcon from '@material-ui/icons/Edit';
import FormAction from '@fbcnms/ui/components/design-system/Form/FormAction';
import FormElementContext from '@fbcnms/ui/components/design-system/Form/FormElementContext';
import IconButton from '@material-ui/core/IconButton';
import {makeStyles} from '@material-ui/styles';

type Props = {
  isOnEdit: boolean,
  onChange: (isOnEdit: boolean) => void,
};

const useStyles = makeStyles({
  disabled: {
    opacity: 0.5,
    pointerEvents: 'none',
  },
});

const EditToggleButton = (props: Props) => {
  const {isOnEdit, onChange} = props;
  const classes = useStyles();

  return (
    <FormAction>
      <FormElementContext.Consumer>
        {context => (
          <IconButton
            onClick={() => onChange && onChange(!isOnEdit)}
            className={context.disabled ? classes.disabled : ''}
            color="primary">
            {isOnEdit ? <ArrowBackIcon /> : <EditIcon />}
          </IconButton>
        )}
      </FormElementContext.Consumer>
    </FormAction>
  );
};

export default EditToggleButton;
