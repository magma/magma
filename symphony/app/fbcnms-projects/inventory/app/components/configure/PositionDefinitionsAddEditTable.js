/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {PositionDefinitionsAddEditTable_positionDefinition} from './__generated__/PositionDefinitionsAddEditTable_positionDefinition.graphql';

import Button from '@fbcnms/ui/components/design-system/Button';
import CardSection from '../CardSection';
import DraggableTableRow from '../draggable/DraggableTableRow';
import DroppableTableBody from '../draggable/DroppableTableBody';
import FormAction from '@fbcnms/ui/components/design-system/Form/FormAction';
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import IconButton from '@fbcnms/ui/components/design-system/IconButton';
import React, {useCallback} from 'react';
import Table from '@material-ui/core/Table';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import fbt from 'fbt';
import inventoryTheme from '../../common/theme';
import update from 'immutability-helper';
import {DeleteIcon, PlusIcon} from '@fbcnms/ui/components/design-system/Icons';
import {graphql} from 'react-relay';
import {makeStyles} from '@material-ui/styles';
import {reorder} from '../draggable/DraggableUtils';

const useStyles = makeStyles(_theme => ({
  table: {
    marginBottom: '12px',
  },
  input: {
    ...inventoryTheme.textField,
    marginTop: '0px',
    marginBottom: '0px',
  },
  cell: {
    ...inventoryTheme.textField,
    paddingLeft: '0px',
  },
  addButton: {
    marginBottom: '12px',
  },
}));

graphql`
  fragment PositionDefinitionsAddEditTable_positionDefinition on EquipmentPositionDefinition
    @relay(mask: false) {
    id
    name
    index
    visibleLabel
  }
`;

type Props = {
  positionDefinitions: Array<PositionDefinitionsAddEditTable_positionDefinition>,
  onPositionDefinitionsChanged?: ?(
    newPositioDef: Array<PositionDefinitionsAddEditTable_positionDefinition>,
  ) => void,
};

const PositionDefinitionsAddEditTable = (props: Props) => {
  const classes = useStyles();

  const getEditingPositionDefinition = useCallback((): PositionDefinitionsAddEditTable_positionDefinition => {
    const index = props.positionDefinitions.length;
    return {
      id: 'PositionDefinition@tmp' + index,
      name: '',
      index: index,
      visibleLabel: '',
    };
  }, [props]);

  const onAddPosition = useCallback(() => {
    const {onPositionDefinitionsChanged} = props;
    onPositionDefinitionsChanged &&
      onPositionDefinitionsChanged(
        update(props.positionDefinitions, {
          $push: [getEditingPositionDefinition()],
        }),
      );
  }, [getEditingPositionDefinition, props]);

  const onPositionDefinitionsChanged = useCallback(
    (newValue, keyName, positionIndex) => {
      const {onPositionDefinitionsChanged} = props;
      onPositionDefinitionsChanged &&
        onPositionDefinitionsChanged(
          update(props.positionDefinitions, {
            // $FlowFixMe Set state for each field
            [positionIndex]: {[keyName]: {$set: newValue}},
          }),
        );
    },
    [props],
  );

  const onRemovePositionClicked = useCallback(
    positionIndex => {
      const {onPositionDefinitionsChanged} = props;
      onPositionDefinitionsChanged &&
        onPositionDefinitionsChanged(
          props.positionDefinitions.length === 1
            ? [getEditingPositionDefinition()]
            : update(props.positionDefinitions, {
                $splice: [[positionIndex, 1]],
              }),
        );
    },
    [getEditingPositionDefinition, props],
  );

  const onDragEnd = useCallback(
    result => {
      if (!result.destination) {
        return;
      }
      const items = reorder(
        props.positionDefinitions,
        result.source.index,
        result.destination.index,
      );
      const newItems = items.map((positionDefinition, i) =>
        update(positionDefinition, {index: {$set: i}}),
      );
      props.onPositionDefinitionsChanged &&
        props.onPositionDefinitionsChanged(newItems);
    },
    [props],
  );

  const {positionDefinitions} = props;
  if (positionDefinitions.length === 0) {
    return null;
  }

  return (
    <CardSection title="Positions">
      <Table component="div" className={classes.table}>
        <TableHead component="div">
          <TableRow component="div">
            <TableCell component="div" size="small" padding="none" />
            <TableCell component="div" className={classes.cell}>
              Name
            </TableCell>
            <TableCell component="div" className={classes.cell}>
              Visible Label
            </TableCell>
            <TableCell component="div" />
          </TableRow>
        </TableHead>
        <DroppableTableBody onDragEnd={onDragEnd}>
          {positionDefinitions.map((definition, i) => (
            <DraggableTableRow key={i} id={definition.id} index={i}>
              <TableCell className={classes.cell} component="div" scope="row">
                <FormField>
                  <TextInput
                    placeholder={`${fbt('Name', '')}`}
                    variant="outlined"
                    className={classes.input}
                    value={definition.name ?? ''}
                    onChange={({target}) =>
                      onPositionDefinitionsChanged(target.value, 'name', i)
                    }
                  />
                </FormField>
              </TableCell>
              <TableCell component="div" className={classes.cell} align="left">
                <FormField>
                  <TextInput
                    placeholder={`${fbt('Visible Label', '')}`}
                    variant="outlined"
                    className={classes.input}
                    value={definition.visibleLabel ?? ''}
                    onChange={({target}) =>
                      onPositionDefinitionsChanged(
                        target.value,
                        'visibleLabel',
                        i,
                      )
                    }
                  />
                </FormField>
              </TableCell>
              <TableCell component="div" align="right">
                <IconButton
                  onClick={onRemovePositionClicked.bind(this, i)}
                  disabled={!definition.id.includes('@tmp')}
                  icon={DeleteIcon}
                />
              </TableCell>
            </DraggableTableRow>
          ))}
        </DroppableTableBody>
      </Table>
      <FormAction>
        <Button
          className={classes.addButton}
          color="primary"
          variant="text"
          leftIcon={PlusIcon}
          onClick={onAddPosition}>
          Add Position
        </Button>
      </FormAction>
    </CardSection>
  );
};

export default PositionDefinitionsAddEditTable;
