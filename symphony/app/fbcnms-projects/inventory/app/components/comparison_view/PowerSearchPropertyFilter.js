/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {FilterProps} from './ComparisonViewTypes';

import * as React from 'react';
import PowerSearchFilter from './PowerSearchFilter';
import PropertyValueInput from '../form/PropertyValueInput';
import nullthrows from '@fbcnms/util/nullthrows';
import {POWER_SEARCH_OPERATOR_ID} from './PowerSearchOperatorComponent';
import {getPropertyValue} from '../../common/Property';
import {makeStyles} from '@material-ui/styles';
import {useState} from 'react';

const useStyles = makeStyles(theme => ({
  inputRoot: {
    width: 'auto',
    backgroundColor: theme.palette.grey.A100,
    marginTop: '0px',
    marginBottom: '0px',
  },
  innerInput: {
    paddingTop: '6px',
    paddingBottom: '6px',
  },
}));

const ENTER_KEY_CODE = 13;

const PowerSearchPropertyFilter = (props: FilterProps) => {
  const {
    value,
    onInputBlurred,
    onValueChanged,
    onRemoveFilter,
    editMode,
    onNewInputBlurred,
  } = props;
  const classes = useStyles();
  const propertyValue = nullthrows(value.propertyValue);
  const [editedOperator, setEditedOperator] = useState(value.operator);

  const onChange = newValue => {
    const newFilterValue = {
      id: value.id,
      key: value.key,
      name: 'property',
      operator: editedOperator,
      propertyValue: {
        id: 'tmp@propertyType',
        name: newValue.propertyType
          ? newValue.propertyType.name
          : newValue.name,
        type: newValue.propertyType
          ? newValue.propertyType.type
          : newValue.type,
        index: 0,
        booleanValue: newValue.booleanValue,
        stringValue: newValue.stringValue,
        intValue: newValue.intValue,
        floatValue: newValue.floatValue,
        latitudeValue: newValue.latitudeValue,
        longitudeValue: newValue.longitudeValue,
        rangeFromValue: newValue.rangeFromValue,
        rangeToValue: newValue.rangeToValue,
      },
    };
    onValueChanged(newFilterValue);
    return newFilterValue;
  };

  return (
    <PowerSearchFilter
      name={propertyValue.name}
      operator={editedOperator}
      editMode={editMode}
      value={String(getPropertyValue(propertyValue))}
      onRemoveFilter={onRemoveFilter}
      onOperatorChange={setEditedOperator}
      input={
        <div>
          <PropertyValueInput
            className={classes.inputRoot}
            inputClassName={classes.innerInput}
            label={null}
            inputType="Property"
            property={propertyValue}
            onBlur={e => {
              if (e.relatedTarget.id === POWER_SEARCH_OPERATOR_ID) {
                return;
              }
              onInputBlurred();
            }}
            autoFocus={true}
            onKeyDown={e => {
              if (e.keyCode === ENTER_KEY_CODE) {
                onInputBlurred();
              }
            }}
            onChange={newValue => {
              if (
                propertyValue.type === 'enum' ||
                propertyValue.type === 'bool'
              ) {
                onNewInputBlurred(onChange(newValue));
              } else {
                onChange(newValue);
              }
            }}
            margin="none"
          />
        </div>
      }
    />
  );
};

export default PowerSearchPropertyFilter;
