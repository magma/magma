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
import ArrowBackIcon from '@material-ui/icons/ArrowBack';
import EditIcon from '@material-ui/icons/Edit';
import IconButton from '@material-ui/core/IconButton';

type Props = {
  isOnEdit: boolean,
  onChange: (isOnEdit: boolean) => void,
};

const EditToggleButton = (props: Props) => {
  const {isOnEdit, onChange} = props;

  return (
    <IconButton onClick={() => onChange && onChange(!isOnEdit)} color="primary">
      {isOnEdit ? <ArrowBackIcon /> : <EditIcon />}
    </IconButton>
  );
};

export default EditToggleButton;
