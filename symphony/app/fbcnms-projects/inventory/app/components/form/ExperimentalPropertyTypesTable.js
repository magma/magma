/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {FeatureID} from '@fbcnms/types/features';
import type {PropertyType} from '../../common/PropertyType';

import * as React from 'react';
import Button from '@fbcnms/ui/components/design-system/Button';
import Checkbox from '@fbcnms/ui/components/design-system/Checkbox/Checkbox';
import DraggableTableRow from '../draggable/DraggableTableRow';
import DroppableTableBody from '../draggable/DroppableTableBody';
import FormAction from '@fbcnms/ui/components/design-system/Form/FormAction';
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import IconButton from '@fbcnms/ui/components/design-system/IconButton';
import PropertyTypeSelect from './PropertyTypeSelect';
import PropertyTypesTableDispatcher from './context/property_types/PropertyTypesTableDispatcher';
import PropertyValueInput from './PropertyValueInput';
import Table from '@material-ui/core/Table';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import fbt from 'fbt';
import inventoryTheme from '../../common/theme';
import {DeleteIcon, PlusIcon} from '@fbcnms/ui/components/design-system/Icons';
import {isTempId} from '../../common/EntUtils';
import {makeStyles} from '@material-ui/styles';
import {sortByIndex} from '../draggable/DraggableUtils';
import {useContext} from 'react';

const useStyles = makeStyles(() => ({
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
}));

export type PropertyTypeInfo = $ReadOnly<{|
  label: string,
  featureFlag?: FeatureID,
  isNode?: boolean,
|}>;

type Props = $ReadOnly<{|
  propertyTypes: Array<PropertyType>,
  supportMandatory?: boolean,
  supportDelete?: boolean,
|}>;

const ExperimentalPropertyTypesTable = ({
  propertyTypes,
  supportMandatory = true,
  supportDelete,
}: Props) => {
  const classes = useStyles();
  const dispatch = useContext(PropertyTypesTableDispatcher);

  return (
    <div className={classes.container}>
      <Table component="div" className={classes.root}>
        <TableHead component="div">
          <TableRow component="div">
            <TableCell size="small" padding="none" component="div" />
            <TableCell component="div" className={classes.cell}>
              <fbt desc="">Experimental Name</fbt>
            </TableCell>
            <TableCell component="div" className={classes.cell}>
              <fbt desc="">Property Type</fbt>
            </TableCell>
            <TableCell component="div" className={classes.cell}>
              <fbt desc="">Default Value</fbt>
            </TableCell>
            <TableCell
              padding="checkbox"
              component="div"
              className={classes.cell}>
              <fbt desc="">Fixed Value</fbt>
            </TableCell>
            {supportMandatory && (
              <TableCell
                padding="checkbox"
                component="div"
                className={classes.cell}>
                <fbt desc="">Mandatory</fbt>
              </TableCell>
            )}
            <TableCell component="div" />
          </TableRow>
        </TableHead>
        <DroppableTableBody
          onDragEnd={({source, destination}) =>
            dispatch({
              type: 'CHANGE_PROPERTY_TYPE_INDEX',
              sourceIndex: source.index,
              destinationIndex: destination.index,
            })
          }>
          {propertyTypes
            .slice()
            .sort(sortByIndex)
            .map((property, i) =>
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
                        onChange={({target}) =>
                          dispatch({
                            type: 'UPDATE_PROPERTY_TYPE_NAME',
                            id: property.id,
                            name: target.value,
                          })
                        }
                        onBlur={() =>
                          dispatch({
                            type: 'UPDATE_PROPERTY_TYPE_NAME',
                            id: property.id,
                            name: property.name.trim(),
                          })
                        }
                      />
                    </FormField>
                  </TableCell>
                  <TableCell
                    className={classes.cell}
                    component="div"
                    scope="row">
                    <FormField>
                      <PropertyTypeSelect propertyType={property} />
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
                      onChange={value =>
                        dispatch({
                          type: 'UPDATE_PROPERTY_TYPE',
                          value,
                        })
                      }
                      margin="dense"
                    />
                  </TableCell>
                  <TableCell padding="checkbox" component="div">
                    <FormField>
                      <Checkbox
                        checked={!property.isInstanceProperty}
                        onChange={checkedNewValue =>
                          dispatch({
                            type: 'UPDATE_PROPERTY_TYPE',
                            value: {
                              ...property,
                              isInstanceProperty: checkedNewValue !== 'checked',
                            },
                          })
                        }
                        title={null}
                      />
                    </FormField>
                  </TableCell>
                  {supportMandatory && (
                    <TableCell padding="checkbox" component="div">
                      <FormField>
                        <Checkbox
                          checked={!!property.isMandatory}
                          onChange={checkedNewValue =>
                            dispatch({
                              type: 'UPDATE_PROPERTY_TYPE',
                              value: {
                                ...property,
                                isMandatory: checkedNewValue === 'checked',
                              },
                            })
                          }
                          title={null}
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
                        onClick={() =>
                          dispatch({
                            type: 'REMOVE_PROPERTY_TYPE',
                            id: property.id,
                          })
                        }
                        disabled={!supportDelete && !isTempId(property.id)}
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
          onClick={() => dispatch({type: 'ADD_PROPERTY_TYPE'})}
          leftIcon={PlusIcon}>
          <fbt desc="">Add Property</fbt>
        </Button>
      </FormAction>
    </div>
  );
};

export default ExperimentalPropertyTypesTable;
