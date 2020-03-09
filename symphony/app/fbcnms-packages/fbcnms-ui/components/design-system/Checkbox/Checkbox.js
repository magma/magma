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
import {useFormElementContext} from '../Form/FormElementContext';
import {useMemo} from 'react';

const useStyles = makeStyles(_theme => ({
  root: {
    width: '24px',
    height: '24px',
    '&:not($disabled)': {
      cursor: 'pointer',
      '&:hover': {
        '& $selection, & $noSelection': {
          fill: SymphonyTheme.palette.B700,
        },
      },
      '& $noSelection': {
        fill: SymphonyTheme.palette.D400,
      },
      '& $selection': {
        fill: SymphonyTheme.palette.primary,
      },
    },
  },
  disabled: {
    '& $noSelection, & $selection': {
      fill: SymphonyTheme.palette.D200,
    },
  },
  selection: {},
  noSelection: {},
}));

export type SelectionType = 'checked' | 'unchecked';

type Props = {
  className?: string,
  checked: boolean,
  indeterminate?: boolean,
  disabled?: ?boolean,
  onChange?: (selection: SelectionType) => void,
  onClick?: ?(SyntheticMouseEvent<Element>) => void,
};

const Checkbox = (props: Props) => {
  const {
    className,
    checked,
    indeterminate,
    onChange,
    onClick,
    disabled: propDisabled = false,
  } = props;
  const classes = useStyles();
  const CheckboxIcon = indeterminate
    ? IndeterminateCheckBoxIcon
    : checked
    ? CheckBoxIcon
    : CheckBoxOutlineBlankIcon;

  const {disabled: contextDisabled} = useFormElementContext();
  const disabled = useMemo(
    () => (propDisabled ? propDisabled : contextDisabled),
    [contextDisabled, propDisabled],
  );

  return (
    <div
      className={classNames(classes.root, className, {
        [classes.disabled]: disabled,
      })}
      onClick={e => {
        if (disabled) {
          return;
        }
        if (onChange) {
          onChange(
            indeterminate ? 'unchecked' : checked ? 'unchecked' : 'checked',
          );
        }
        if (onClick) {
          onClick(e);
        }
      }}>
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
