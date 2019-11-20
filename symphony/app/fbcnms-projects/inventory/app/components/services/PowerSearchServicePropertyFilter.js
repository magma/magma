/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {FilterProps} from '../comparison_view/ComparisonViewTypes';

import * as React from 'react';
import PowerSearchFilter from '../comparison_view/PowerSearchFilter';
import PropertyValueInput from '../form/PropertyValueInput';
import nullthrows from '@fbcnms/util/nullthrows';
import {getPropertyValue} from '../../common/Property';
import {makeStyles} from '@material-ui/styles';

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

const SERVICE_PROPERTY_FILTER_NAME = 'service_inst_property';

const PowerSearchServicePropertyFilter = (props: FilterProps) => {
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
  const onChange = newValue => {
    const newFilterValue = {
      id: value.id,
      key: value.key,
      name: SERVICE_PROPERTY_FILTER_NAME,
      operator: 'is',
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
      operator={value.operator}
      editMode={editMode}
      value={String(getPropertyValue(propertyValue))}
      onRemoveFilter={onRemoveFilter}
      input={
        <div>
          <PropertyValueInput
            className={classes.inputRoot}
            inputClassName={classes.innerInput}
            label={null}
            inputType="Property"
            property={propertyValue}
            onBlur={onInputBlurred}
            onKeyDown={e => {
              if (e.keyCode === ENTER_KEY_CODE) {
                onInputBlurred();
              }
            }}
            onChange={newValue => {
              if (propertyValue.type === 'enum') {
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

export {PowerSearchServicePropertyFilter, SERVICE_PROPERTY_FILTER_NAME};
