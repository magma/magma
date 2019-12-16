/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import CheckBoxIcon from '@material-ui/icons/CheckBox';
import CheckBoxOutlineBlankIcon from '@material-ui/icons/CheckBoxOutlineBlank';
import IndeterminateCheckBoxIcon from '@material-ui/icons/IndeterminateCheckBox';
import React from 'react';
import SymphonyTheme from '../../../theme/symphony';
import classNames from 'classnames';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(_theme => ({
  root: {
    width: '24px',
    height: '24px',
    cursor: 'pointer',
    '&:hover $selection, &:hover $noSelection': {
      fill: SymphonyTheme.palette.B700,
    },
  },
  selection: {
    fill: SymphonyTheme.palette.primary,
  },
  noSelection: {
    fill: SymphonyTheme.palette.D400,
  },
}));

export type SelectionType = 'checked' | 'unchecked';

type Props = {
  className?: string,
  checked: boolean,
  indeterminate?: boolean,
  onChange?: (selection: SelectionType) => void,
};

const Checkbox = (props: Props) => {
  const {className, checked, indeterminate, onChange} = props;
  const classes = useStyles();
  const CheckboxIcon = indeterminate
    ? IndeterminateCheckBoxIcon
    : checked
    ? CheckBoxIcon
    : CheckBoxOutlineBlankIcon;

  return (
    <div
      className={classNames(classes.root, className)}
      onClick={() =>
        onChange &&
        onChange(
          indeterminate ? 'unchecked' : checked ? 'unchecked' : 'checked',
        )
      }>
      <CheckboxIcon
        className={classNames({
          [classes.selection]: checked || indeterminate,
          [classes.noSelection]: !checked && !indeterminate,
        })}
      />
    </div>
  );
};

Checkbox.defaultProps = {
  checked: false,
  indeterminate: false,
};

export default Checkbox;
