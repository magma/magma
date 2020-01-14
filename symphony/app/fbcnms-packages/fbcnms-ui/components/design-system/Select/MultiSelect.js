/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {OptionProps} from './SelectMenu';

import * as React from 'react';
import ArrowDropDownIcon from '@material-ui/icons/ArrowDropDown';
import BasePopoverTrigger from '../ContexualLayer/BasePopoverTrigger';
import Button from '../Button';
import MultiSelectMenu from './MultiSelectMenu';
import classNames from 'classnames';
import fbt from 'fbt';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles({
  root: {
    justifyContent: 'flex-start',
    padding: '4px',
  },
  value: {
    fontWeight: 500,
  },
  menu: {
    margin: '8px 0px',
  },
  label: {
    fontWeight: 400,
  },
});

type Props<TValue> = {
  className?: string,
  label: React.Node,
  options: Array<OptionProps<TValue>>,
  onChange: (option: OptionProps<TValue>) => void | (() => void),
  selectedValues: Array<OptionProps<TValue>>,
  onOptionsFetchRequested?: (searchTerm: string) => void,
  searchable?: boolean,
};

const MultiSelect = <TValue>({
  label,
  className,
  ...selectMenuProps
}: Props<TValue>) => {
  const classes = useStyles();
  const {selectedValues} = selectMenuProps;
  return (
    <BasePopoverTrigger
      popover={
        <MultiSelectMenu
          {...selectMenuProps}
          className={classes.menu}
          size="normal"
        />
      }>
      {(onShow, contextRef) => (
        <Button
          className={classNames(classes.root, className)}
          ref={contextRef}
          onClick={onShow}
          skin="regular"
          rightIcon={ArrowDropDownIcon}>
          <span className={classes.label}>{label}</span>
          {selectedValues.length > 0 ? ': ' : null}
          {selectedValues.length === 1 ? (
            <span className={classes.value} key={String(selectedValues[0])}>
              {selectedValues[0].label ?? ''}
            </span>
          ) : null}
          {selectedValues.length > 1 ? (
            <span className={classes.value}>
              <fbt desc="Amount of selected items">
                <fbt:param name="num_selected" number={true}>
                  {selectedValues.length}
                </fbt:param>
                Selected
              </fbt>
            </span>
          ) : null}
        </Button>
      )}
    </BasePopoverTrigger>
  );
};

export default MultiSelect;
