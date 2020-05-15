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
import type {PropertyKind} from '../configure/__generated__/WorkOrderTypesQuery.graphql';
import type {PropertyType} from '../../common/PropertyType';
import type {WithStyles} from '@material-ui/core';

import * as React from 'react';
import AppContext from '@fbcnms/ui/context/AppContext';
import Button from '@fbcnms/ui/components/design-system/Button';
import Checkbox from '@fbcnms/ui/components/design-system/Checkbox/Checkbox';
import DraggableTableRow from '../draggable/DraggableTableRow';
import DroppableTableBody from '../draggable/DroppableTableBody';
import FormAction from '@fbcnms/ui/components/design-system/Form/FormAction';
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import IconButton from '@fbcnms/ui/components/design-system/IconButton';
import PropertyTypeSelect from './PropertyTypeSelect';
import PropertyValueInput from './PropertyValueInput';
import Table from '@material-ui/core/Table';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import inventoryTheme from '../../common/theme';
import {DeleteIcon, PlusIcon} from '@fbcnms/ui/components/design-system/Icons';
import {removeItem, setItem, updateItem} from '@fbcnms/util/arrays';
import {reorder} from '../draggable/DraggableUtils';
import {withStyles} from '@material-ui/core/styles';

const styles = () => ({
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
  selectMenu: {
    height: '14px',
  },
  actionsBar: {
    width: '20px',
  },
});

export type PropertyTypeInfo = $ReadOnly<{|
  kind: PropertyKind,
  label: string,
  featureFlag?: FeatureID,
|}>;

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
            {propertyTypes.map((property, i) =>
              property.isDeleted ? null : (
                <DraggableTableRow id={property.id} index={i} key={property.id}>
                  <TableCell
                    className={classes.cell}
                    component="div"
                    scope="row">
                    <FormField>
                      <TextInput
                        autoFocus={true}
                        placeholder="Name"
                        variant="outlined"
                        className={classes.input}
                        value={property.name}
                        onChange={this._handleNameChange(i)}
                        onBlur={() => this._handleNameBlur(i)}
                      />
                    </FormField>
                  </TableCell>
                  <TableCell
                    className={classes.cell}
                    component="div"
                    scope="row">
                    <FormField>
                      <PropertyTypeSelect
                        propertyType={property}
                        onPropertyTypeChange={this._handleTypeChange(i)}
                      />
                    </FormField>
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
                    <FormField>
                      <Checkbox
                        checked={!property.isInstanceProperty}
                        title={null}
                        onChange={this._handleChecked(i)}
                      />
                    </FormField>
                  </TableCell>
                  {supportMandatory && (
                    <TableCell padding="checkbox" component="div">
                      <FormField>
                        <Checkbox
                          checked={!!property.isMandatory}
                          title={null}
                          onChange={this._handleIsMandatoryChecked(i)}
                        />
                      </FormField>
                    </TableCell>
                  )}
                  <TableCell
                    className={classes.actionsBar}
                    align="right"
                    component="div">
                    <FormAction>
                      <IconButton
                        skin="primary"
                        onClick={this._onRemovePropertyClicked(i, property)}
                        disabled={
                          !this.props.supportDelete &&
                          !property.id.includes('@tmp')
                        }
                        icon={DeleteIcon}
                      />
                    </FormAction>
                  </TableCell>
                </DraggableTableRow>
              ),
            )}
          </DroppableTableBody>
        </Table>
        <FormAction>
          <Button
            color="primary"
            variant="text"
            onClick={this._onAddProperty}
            leftIcon={PlusIcon}>
            Add Property
          </Button>
        </FormAction>
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

  _handleNameChange = index => event => {
    this.props.onPropertiesChanged(
      updateItem<PropertyType, 'name'>(
        this.props.propertyTypes,
        index,
        'name',
        // $FlowFixMe: need to figure out how to cast string to PropertyKind
        event.target.value,
      ),
    );
  };

  _handleTypeChange = (index: number) => (value: PropertyType) => {
    this.props.onPropertiesChanged([
      ...this.props.propertyTypes.slice(0, index),
      value,
      ...this.props.propertyTypes.slice(index + 1),
    ]);
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

  _handleChecked = index => checkedNewValue => {
    this.props.onPropertiesChanged(
      updateItem<PropertyType, 'isInstanceProperty'>(
        this.props.propertyTypes,
        index,
        'isInstanceProperty',
        checkedNewValue !== 'checked',
      ),
    );
  };

  _handleIsMandatoryChecked = index => checkedNewValue => {
    this.props.onPropertiesChanged(
      updateItem<PropertyType, 'isMandatory'>(
        this.props.propertyTypes,
        index,
        'isMandatory',
        checkedNewValue === 'checked',
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
      nodeType: null,
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
