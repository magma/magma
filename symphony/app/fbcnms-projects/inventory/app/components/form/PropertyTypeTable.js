/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {AppContextType} from '@fbcnms/ui/context/AppContext';
import type {FeatureID} from '@fbcnms/types/features';
import type {Property} from '../../common/Property';
import type {PropertyType} from '../../common/PropertyType';
import type {WithStyles} from '@material-ui/core';

import * as React from 'react';
import AppContext from '@fbcnms/ui/context/AppContext';
import Button from '@material-ui/core/Button';
import Checkbox from '@material-ui/core/Checkbox';
import DeleteIcon from '@material-ui/icons/Delete';
import DraggableTableRow from '../draggable/DraggableTableRow';
import DroppableTableBody from '../draggable/DroppableTableBody';
import IconButton from '@material-ui/core/IconButton';
import MenuItem from '@material-ui/core/MenuItem';
import PropertyValueInput from './PropertyValueInput';
import Table from '@material-ui/core/Table';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import TextField from '@material-ui/core/TextField';
import inventoryTheme from '../../common/theme';
import {removeItem, setItem, updateItem} from '@fbcnms/util/arrays';
import {reorder} from '../draggable/DraggableUtils';
import {withStyles} from '@material-ui/core/styles';

const styles = _theme => ({
  container: {
    maxWidth: '1366px',
    overflowX: 'auto',
  },
  root: {
    marginBottom: '12px',
    maxWidth: '100%',
  },
  input: {
    ...inventoryTheme.textField,
    marginTop: '0px',
    marginBottom: '0px',
    width: '100%',
  },
  cell: {
    ...inventoryTheme.textField,
    paddingLeft: '0px',
    width: 'unset',
  },
  addButton: {
    marginBottom: '12px',
  },
  selectMenu: {
    height: '14px',
  },
  actionsBar: {
    width: '20px',
  },
});

type PropertyTypeInfo = {|
  label: string,
  featureFlag?: FeatureID,
|};

const propertyTypeLabels: {[string]: PropertyTypeInfo} = {
  date: {label: 'Date'},
  datetime_local: {label: 'Date & Time'},
  int: {label: 'Number'},
  float: {label: 'Float'},
  string: {label: 'Text'},
  email: {label: 'Email'},
  gps_location: {label: 'Coordinates'},
  bool: {label: 'True or False'},
  range: {label: 'Range'},
  enum: {label: 'Multiple choice'},
  equipment: {label: 'Equipment'},
  location: {label: 'Location'},
  service: {label: 'Service', featureFlag: 'services'},
};

type Props = {
  propertyTypes: Array<PropertyType>,
  onPropertiesChanged: (newProperties: Array<PropertyType>) => void,
  supportMandatory?: boolean,
  supportDelete?: boolean,
} & WithStyles<typeof styles>;

class PropertyTypeTable extends React.Component<Props> {
  static contextType = AppContext;
  context: AppContextType;
  render() {
    const {classes} = this.props;
    const propertyTypes = this.props.propertyTypes;
    const {supportMandatory = true} = this.props;
    return (
      <div className={classes.container}>
        <Table component="div" className={classes.root}>
          <TableHead component="div">
            <TableRow component="div">
              <TableCell size="small" padding="none" component="div" />
              <TableCell component="div" className={classes.cell}>
                Name
              </TableCell>
              <TableCell component="div" className={classes.cell}>
                Property Type
              </TableCell>
              <TableCell component="div" className={classes.cell}>
                Default Value
              </TableCell>
              <TableCell
                padding="checkbox"
                component="div"
                className={classes.cell}>
                Fixed Value
              </TableCell>
              {supportMandatory && (
                <TableCell
                  padding="checkbox"
                  component="div"
                  className={classes.cell}>
                  Mandatory
                </TableCell>
              )}
              <TableCell component="div" />
            </TableRow>
          </TableHead>
          <DroppableTableBody onDragEnd={this._onDragEnd}>
            {propertyTypes
              .filter(property => !property.isDeleted)
              .map((property, i) => (
                <DraggableTableRow id={property.id} index={i} key={i}>
                  <TableCell
                    className={classes.cell}
                    component="div"
                    scope="row">
                    <TextField
                      autoFocus={true}
                      placeholder="Name"
                      variant="outlined"
                      className={classes.input}
                      value={property.name}
                      onChange={this._handleChange('name', i)}
                      onBlur={() => this._handleNameBlur(i)}
                      margin="dense"
                    />
                  </TableCell>
                  <TableCell
                    className={classes.cell}
                    component="div"
                    scope="row">
                    <TextField
                      select
                      variant="outlined"
                      className={classes.input}
                      value={property.type}
                      onChange={this._handleChange('type', i)}
                      SelectProps={{
                        classes: {selectMenu: classes.selectMenu},
                        MenuProps: {
                          className: classes.menu,
                        },
                      }}
                      margin="dense">
                      {Object.keys(propertyTypeLabels)
                        .filter(
                          type =>
                            !propertyTypeLabels[type].featureFlag ||
                            this.context.isFeatureEnabled(
                              propertyTypeLabels[type].featureFlag,
                            ),
                        )
                        .map(type => (
                          <MenuItem key={type} value={type}>
                            {propertyTypeLabels[type].label}
                          </MenuItem>
                        ))}
                    </TextField>
                  </TableCell>
                  <TableCell
                    className={classes.cell}
                    component="div"
                    scope="row">
                    <PropertyValueInput
                      label={null}
                      className={classes.input}
                      inputType="PropertyType"
                      property={property}
                      onChange={this._handlePropertyTypeChange(i)}
                      margin="dense"
                    />
                  </TableCell>
                  <TableCell padding="checkbox" component="div">
                    <Checkbox
                      checked={!property.isInstanceProperty}
                      onChange={this._handleChecked(i)}
                      color="primary"
                    />
                  </TableCell>
                  {supportMandatory && (
                    <TableCell padding="checkbox" component="div">
                      <Checkbox
                        checked={!!property.isMandatory}
                        onChange={this._handleIsMandatoryChecked(i)}
                        color="primary"
                      />
                    </TableCell>
                  )}
                  <TableCell
                    className={classes.actionsBar}
                    align="right"
                    component="div">
                    <IconButton
                      onClick={this._onRemovePropertyClicked(i, property)}
                      disabled={
                        !this.props.supportDelete &&
                        !property.id.includes('@tmp')
                      }>
                      <DeleteIcon />
                    </IconButton>
                  </TableCell>
                </DraggableTableRow>
              ))}
          </DroppableTableBody>
        </Table>
        <Button
          className={classes.addButton}
          color="primary"
          variant="outlined"
          onClick={this._onAddProperty}>
          Add Property
        </Button>
      </div>
    );
  }

  _handlePropertyTypeChange = (index: number) => (
    property: PropertyType | Property,
  ) => {
    if (property.propertyType) {
      // Filter out properties, we are just dealing with propertyTypes
      return;
    }
    this.props.onPropertiesChanged(
      setItem(this.props.propertyTypes, index, property),
    );
  };

  _handleChange = (changedProp: 'name' | 'type', index) => event => {
    this.props.onPropertiesChanged(
      updateItem<PropertyType, typeof changedProp>(
        this.props.propertyTypes,
        index,
        changedProp,
        // $FlowFixMe: need to figure out how to cast string to PropertyKind
        event.target.value,
      ),
    );
  };

  _handleNameBlur = index => {
    const name = this.props.propertyTypes[index]?.name;
    const trimmedName = name && name.trim();
    if (name === trimmedName) {
      return;
    }

    this.props.onPropertiesChanged(
      updateItem<PropertyType, 'name'>(
        this.props.propertyTypes,
        index,
        'name',
        trimmedName,
      ),
    );
  };

  _handleChecked = index => event => {
    this.props.onPropertiesChanged(
      updateItem<PropertyType, 'isInstanceProperty'>(
        this.props.propertyTypes,
        index,
        'isInstanceProperty',
        !event.target.checked,
      ),
    );
  };

  _handleIsMandatoryChecked = index => event => {
    this.props.onPropertiesChanged(
      updateItem<PropertyType, 'isMandatory'>(
        this.props.propertyTypes,
        index,
        'isMandatory',
        event.target.checked,
      ),
    );
  };

  _onAddProperty = () => {
    this.props.onPropertiesChanged([
      ...this.props.propertyTypes,
      this.getInitialProperty(),
    ]);
  };

  _onRemovePropertyClicked = (index, property: PropertyType) => _event => {
    if (property.id?.includes('@tmp')) {
      this.props.onPropertiesChanged(
        removeItem(this.props.propertyTypes, index),
      );
    } else {
      this.props.onPropertiesChanged(
        updateItem<PropertyType, 'isDeleted'>(
          this.props.propertyTypes,
          index,
          'isDeleted',
          true,
        ),
      );
    }
  };

  _onDragEnd = result => {
    if (!result.destination) {
      return;
    }

    const items = reorder(
      this.props.propertyTypes,
      result.source.index,
      result.destination.index,
    );

    const newItems = items.map((property, i) => ({...property, index: i}));
    this.props.onPropertiesChanged(newItems);
  };

  getInitialProperty(): PropertyType {
    return {
      id: `PropertyType@tmp-${this.props.propertyTypes.length}-${Date.now()}`,
      name: '',
      index: this.props.propertyTypes.length,
      type: 'string',
      booleanValue: false,
      stringValue: null,
      intValue: null,
      floatValue: null,
      latitudeValue: null,
      longitudeValue: null,
      rangeFromValue: null,
      rangeToValue: null,
      isEditable: true,
      isInstanceProperty: true,
    };
  }
}

export default withStyles(styles)(PropertyTypeTable);
