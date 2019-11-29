/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {Property} from '../../common/Property';
import type {PropertyType} from '../../common/PropertyType';
import type {WithStyles} from '@material-ui/core';

import * as React from 'react';
import EnumPropertySelectValueInput from './EnumPropertySelectValueInput';
import EnumPropertyValueInput from './EnumPropertyValueInput';
import EquipmentTypeahead from '../typeahead/EquipmentTypeahead';
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import GPSPropertyValueInput from './GPSPropertyValueInput';
import LocationTypeahead from '../typeahead/LocationTypeahead';
import MenuItem from '@material-ui/core/MenuItem';
import RangePropertyValueInput from './RangePropertyValueInput';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextField from '@material-ui/core/TextField';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import classNames from 'classnames';
import update from 'immutability-helper';
import {withStyles} from '@material-ui/core/styles';

type Props = {
  autoFocus: boolean,
  className: string,
  inputClassName?: ?string,
  label: ?string,
  inputType: 'Property' | 'PropertyType',
  property: Property | PropertyType,
  required: boolean,
  disabled: boolean,
  onChange: (Property | PropertyType) => void,
  onBlur?: () => void,
  onKeyDown?: ?(e: SyntheticKeyboardEvent<>) => void,
  margin: 'none' | 'dense' | 'normal',
  headlineVariant?: 'headline' | 'form',
  fullWidth?: boolean,
} & WithStyles<typeof styles>;

const styles = {
  input: {
    width: (props: Props): string => (props.fullWidth ? 'auto' : '300px'),
    display: 'flex',
    '&&': {
      margin: '0px',
    },
  },
  container: {
    display: 'flex',
    width: '280px',
  },
  toValue: {
    marginLeft: '6px',
  },
  selectMenu: {
    height: '32px',
    padding: '6px 8px',
    boxSizing: 'border-box',
    display: 'flex',
    alignItems: 'center',
  },
};

class PropertyValueInput extends React.Component<Props> {
  static defaultProps = {
    required: false,
    autoFocus: false,
    disabled: false,
    headlineVariant: 'headline',
    fullWidth: false,
  };

  getTextInput = (): React.Node => {
    const {
      autoFocus,
      classes,
      disabled,
      onChange,
      onBlur,
      margin,
      required,
      className,
      inputClassName,
      onKeyDown,
      inputType,
      headlineVariant,
    } = this.props;
    const property = this.props.property;
    const propertyType = !!property.propertyType
      ? property.propertyType
      : property;
    const label = headlineVariant === 'form' ? null : this.props.label;
    const propInputType = propertyType.type;
    switch (propInputType) {
      case 'enum': {
        return inputType == 'Property' ? (
          <EnumPropertySelectValueInput
            label={label}
            className={className}
            inputClassName={classNames(classes.selectMenu, inputClassName)}
            margin={margin}
            property={property}
            onChange={onChange}
            autoFocus={true}
          />
        ) : (
          <EnumPropertyValueInput property={property} onChange={onChange} />
        );
      }
      case 'date':
      case 'email':
      case 'string':
        const coercedInputType: 'date' | 'email' | 'string' = propInputType;
        return (
          <TextInput
            autoFocus={autoFocus}
            required={required}
            disabled={disabled}
            id="property-value"
            label={label}
            variant="outlined"
            className={classNames(classes.input, className)}
            margin={margin}
            value={property.stringValue ?? ''}
            onBlur={() => onBlur && onBlur()}
            onKeyDown={e => onKeyDown && onKeyDown(e)}
            onChange={event =>
              onChange(
                update(property, {
                  stringValue: {$set: event.target.value},
                }),
              )
            }
            InputLabelProps={{
              shrink: !!(propInputType === 'date' || property.stringValue),
            }}
            inputProps={{className: inputClassName}}
            type={coercedInputType}
          />
        );
      case 'int':
        return (
          <TextInput
            autoFocus={autoFocus}
            required={required}
            disabled={disabled}
            id="property-value"
            label={label}
            variant="outlined"
            className={classNames(classes.input, className)}
            margin={margin}
            placeholder={'0'}
            {...(property.intValue ? {value: property.intValue} : {})}
            onBlur={() => onBlur && onBlur()}
            onKeyDown={e => onKeyDown && onKeyDown(e)}
            onChange={event =>
              onChange(
                update(property, {
                  intValue: {$set: parseInt(event.target.value)},
                }),
              )
            }
            inputProps={{className: inputClassName}}
            type="number"
          />
        );
      case 'float':
        return (
          <TextInput
            autoFocus={autoFocus}
            required={required}
            disabled={disabled}
            id="property-value"
            label={label}
            variant="outlined"
            className={classNames(classes.input, className)}
            margin={margin}
            value={property.floatValue ?? 0}
            onBlur={() => onBlur && onBlur()}
            onKeyDown={e => onKeyDown && onKeyDown(e)}
            onChange={event =>
              onChange(
                update(property, {
                  floatValue: {$set: parseFloat(event.target.value)},
                }),
              )
            }
            inputProps={{className: inputClassName}}
            type="number"
          />
        );
      case 'gps_location':
        return (
          <GPSPropertyValueInput
            required={required}
            disabled={disabled}
            id="property-value"
            label={label}
            className={classNames(classes.input, className)}
            margin={margin}
            value={{
              latitude: property.latitudeValue,
              longitude: property.longitudeValue,
            }}
            onLatitudeChange={event =>
              onChange(
                update(property, {
                  latitudeValue: {$set: parseFloat(event.target.value)},
                }),
              )
            }
            onLongitudeChange={event =>
              onChange(
                update(property, {
                  longitudeValue: {$set: parseFloat(event.target.value)},
                }),
              )
            }
          />
        );
      case 'bool':
        return (
          <TextField
            autoFocus={autoFocus}
            select
            id="property-value"
            variant="outlined"
            className={classNames(classes.input, className)}
            onBlur={() => onBlur && onBlur()}
            label={label}
            margin={margin}
            value={!!property.booleanValue ? 'True' : 'False'}
            onChange={event =>
              onChange(
                update(property, {
                  booleanValue: {
                    $set: event.target.value === 'True' ? true : false,
                  },
                }),
              )
            }
            SelectProps={{
              classes: {
                selectMenu: classes.selectMenu,
              },
            }}>
            {['True', 'False'].map(boolVal => (
              <MenuItem key={boolVal} value={boolVal}>
                {boolVal}
              </MenuItem>
            ))}
          </TextField>
        );
      case 'range':
        return (
          <RangePropertyValueInput
            required={required}
            disabled={disabled}
            id="property-value"
            label={label}
            className={classNames(classes.input, className)}
            margin={margin}
            onBlur={() => onBlur && onBlur()}
            value={{
              rangeFrom: property.rangeFromValue,
              rangeTo: property.rangeToValue,
            }}
            onRangeFromChange={event =>
              onChange(
                update(property, {
                  rangeFromValue: {$set: parseFloat(event.target.value)},
                }),
              )
            }
            onRangeToChange={event =>
              onChange(
                update(property, {
                  rangeToValue: {$set: parseFloat(event.target.value)},
                }),
              )
            }
          />
        );
      case 'equipment':
        return inputType == 'Property' ? (
          <EquipmentTypeahead
            margin="dense"
            // eslint-disable-next-line no-warning-comments
            // $FlowFixMe - need to fix this entire file as it receives either property or property type
            selectedEquipment={property.equipmentValue}
            onEquipmentSelection={equipment =>
              onChange(
                update(property, {
                  equipmentValue: {$set: equipment},
                }),
              )
            }
            headline={label}
          />
        ) : (
          <Text>-</Text>
        );
      case 'location':
        return inputType == 'Property' ? (
          <LocationTypeahead
            margin="dense"
            // eslint-disable-next-line no-warning-comments
            // $FlowFixMe - need to fix this entire file as it receives either property or property type
            selectedLocation={property.locationValue}
            onLocationSelection={location =>
              onChange(
                update(property, {
                  locationValue: {$set: location},
                }),
              )
            }
            headline={label}
          />
        ) : (
          <Text>-</Text>
        );
    }
    return null;
  };

  render() {
    const {property, headlineVariant} = this.props;
    const propertyType = !!property.propertyType
      ? property.propertyType
      : property;
    const input = this.getTextInput();
    return headlineVariant === 'form' ? (
      <FormField
        label={propertyType.name}
        hasSpacer={
          propertyType.type !== 'gps_location' && propertyType.type !== 'range'
        }>
        {input}
      </FormField>
    ) : (
      input
    );
  }
}

// eslint-disable-next-line no-warning-comments
// $FlowFixMe - styling based on props works, but flow doesn't recognize it.
export default withStyles(styles)(PropertyValueInput);
