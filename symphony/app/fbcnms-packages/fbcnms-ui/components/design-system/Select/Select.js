/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ButtonProps} from '../Button';
import type {OptionProps} from './SelectMenu';

import * as React from 'react';
import ArrowDropDownIcon from '@material-ui/icons/ArrowDropDown';
import BasePopoverTrigger from '../ContexualLayer/BasePopoverTrigger';
import Button from '../Button';
import FormElementContext from '../Form/FormElementContext';
import SelectMenu from './SelectMenu';
import classNames from 'classnames';
import symphony from '@fbcnms/ui/theme/symphony';
import {makeStyles} from '@material-ui/styles';
import {useContext, useMemo} from 'react';

const useStyles = makeStyles({
  root: {
    justifyContent: 'flex-start',
    padding: '4px',
    border: `1px solid ${symphony.palette.D100}`,
    '&$disabled': {
      backgroundColor: symphony.palette.D50,
    },
  },
  disabled: {
    '&&': {
      '&&': {
        color: symphony.palette.disabled,
      },
    },
  },
  formValue: {
    ...symphony.typography.body2,
  },
  menu: {
    margin: '8px 0px',
  },
  label: {
    fontWeight: symphony.typography.body2.fontWeight,
  },
});

type Props<TValue> = {
  className?: string,
  label?: React.Node,
  options: Array<OptionProps<TValue>>,
  onChange: (value: TValue) => void | (() => void),
  selectedValue: ?TValue,
  ...ButtonProps,
};

const Select = <TValue>({
  label,
  className,
  ...selectMenuProps
}: Props<TValue>) => {
  const classes = useStyles();
  const {
    options,
    selectedValue,
    skin,
    variant,
    disabled: disabledProp,
  } = selectMenuProps;
  const {disabled: contextDisabled} = useContext(FormElementContext);
  const disabled = useMemo(
    () => (disabledProp ? disabledProp : contextDisabled),
    [disabledProp, contextDisabled],
  );
  return (
    <BasePopoverTrigger
      popover={<SelectMenu {...selectMenuProps} className={classes.menu} />}>
      {(onShow, _onHide, contextRef) => (
        <Button
          className={classNames(classes.root, className, {
            [classes.disabled]: disabled,
          })}
          ref={contextRef}
          onClick={onShow}
          skin={skin ?? 'regular'}
          variant={variant}
          disabled={disabled}
          rightIcon={ArrowDropDownIcon}
          rightIconClass={classNames({[classes.disabled]: disabled})}>
          <span className={classes.label}>{label}</span>
          {selectedValue != null && !!label ? ': ' : null}
          {selectedValue != null ? (
            <span
              className={
                classNames({
                  [classes.formValue]: !label,
                  [classes.disabled]: !label && disabled,
                }) || null
              }>
              {options.find(option => option.value === selectedValue)?.label ??
                ''}
            </span>
          ) : null}
        </Button>
      )}
    </BasePopoverTrigger>
  );
};

export default Select;
