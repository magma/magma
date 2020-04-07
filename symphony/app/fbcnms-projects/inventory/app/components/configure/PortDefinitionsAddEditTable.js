/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {EquipmentPortType} from '../../common/EquipmentType';
import type {PortDefinitionsAddEditTable_portDefinitions} from './__generated__/PortDefinitionsAddEditTable_portDefinitions.graphql';
import type {WithStyles} from '@material-ui/core';

import Button from '@fbcnms/ui/components/design-system/Button';
import CardSection from '../CardSection';
import CircularProgress from '@material-ui/core/CircularProgress';
import DraggableTableRow from '../draggable/DraggableTableRow';
import DroppableTableBody from '../draggable/DroppableTableBody';
import IconButton from '@fbcnms/ui/components/design-system/IconButton';
import ListItemText from '@material-ui/core/ListItemText';
import MenuItem from '@material-ui/core/MenuItem';
import {DeleteIcon, PlusIcon} from '@fbcnms/ui/components/design-system/Icons';
// $FlowFixMe - it exists
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';
import RelayEnvironment from '../../common/RelayEnvironment.js';
import Select from '@material-ui/core/Select';
import Table from '@material-ui/core/Table';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import TextField from '@material-ui/core/TextField';
import inventoryTheme from '../../common/theme';
import update from 'immutability-helper';
import {fetchQuery, graphql} from 'react-relay';
import {reorder} from '../draggable/DraggableUtils';
import {sortLexicographically} from '@fbcnms/ui/utils/displayUtils';
import {withStyles} from '@material-ui/core/styles';

const styles = _theme => ({
  table: {
    marginBottom: '12px',
  },
  input: inventoryTheme.textField,
  cell: {
    ...inventoryTheme.textField,
    paddingLeft: '0px',
  },
  addButton: {
    marginBottom: '12px',
  },
});

graphql`
  fragment PortDefinitionsAddEditTable_portDefinitions on EquipmentPortDefinition
    @relay(plural: true) {
    id
    name
    index
    visibleLabel
    portType {
      id
      name
    }
  }
`;

const equipmentPortTypesQuery = graphql`
  query PortDefinitionsAddEditTable__equipmentPortTypesQuery {
    equipmentPortTypes(first: 500)
      @connection(key: "PortDefinitionsTable_equipmentPortTypes") {
      edges {
        node {
          id
          name
        }
      }
    }
  }
`;

type Props = {
  portDefinitions: PortDefinitionsAddEditTable_portDefinitions,
  onPortDefinitionsChanged?: ?(
    newPorts: PortDefinitionsAddEditTable_portDefinitions,
  ) => void,
} & WithStyles<typeof styles>;

type State = {
  equipmentPortTypes: ?Array<EquipmentPortType>,
};

class PortDefinitionsAddEditTable extends React.Component<Props, State> {
  state = {
    equipmentPortTypes: null,
  };

  componentDidMount() {
    this.getEquipmentPortTypes().then(equipmentPortTypes => {
      this.setState({
        equipmentPortTypes,
      });
    });
  }

  render() {
    const {portDefinitions, classes} = this.props;
    if (portDefinitions.length === 0) {
      return null;
    }
    const {equipmentPortTypes} = this.state;
    if (!equipmentPortTypes) {
      return (
        <div className={classes.loadingContainer}>
          <CircularProgress size={50} />
        </div>
      );
    }
    return (
      <CardSection title="Ports">
        <Table component="div" className={classes.table}>
          <TableHead component="div">
            <TableRow component="div">
              <TableCell component="div" size="small" padding="checkbox" />
              <TableCell component="div" className={classes.cell}>
                Name
              </TableCell>
              <TableCell component="div" className={classes.cell}>
                Visible Label
              </TableCell>
              <TableCell component="div" className={classes.cell}>
                Type
              </TableCell>
              <TableCell component="div" />
            </TableRow>
          </TableHead>
          <DroppableTableBody onDragEnd={this._onDragEnd}>
            {portDefinitions.map((portDefinition, i) => (
              <DraggableTableRow key={i} id={portDefinition.id} index={i}>
                <TableCell className={classes.cell} component="div" scope="row">
                  {this.getEditablePortPropertyCell(
                    i,
                    portDefinition.name,
                    'name',
                    'Name',
                  )}
                </TableCell>
                <TableCell
                  component="div"
                  className={classes.cell}
                  align="left">
                  {this.getEditablePortPropertyCell(
                    i,
                    portDefinition.visibleLabel,
                    'visibleLabel',
                    'Visible Label',
                  )}
                </TableCell>
                <TableCell className={classes.cell} component="div" scope="row">
                  {portDefinition.id.includes('@tmp') ? (
                    <Select
                      value={portDefinition.portType?.id || ''}
                      input={<OutlinedInput margin="dense" />}
                      onChange={({target}) =>
                        this.onPortPropertyChanged(
                          'portType',
                          {id: target.value},
                          i,
                        )
                      }
                      MenuProps={{
                        className: classes.menu,
                      }}
                      margin="dense">
                      {equipmentPortTypes.map(equipmentPortType => (
                        <MenuItem value={equipmentPortType.id}>
                          <ListItemText>{equipmentPortType.name}</ListItemText>
                        </MenuItem>
                      ))}
                    </Select>
                  ) : (
                    portDefinition.portType?.name
                  )}
                </TableCell>
                <TableCell component="div" align="right">
                  <IconButton
                    onClick={this.onRemovePortClicked.bind(this, i)}
                    disabled={!portDefinition.id.includes('@tmp')}
                    icon={DeleteIcon}
                  />
                </TableCell>
              </DraggableTableRow>
            ))}
          </DroppableTableBody>
        </Table>
        <Button
          className={classes.addButton}
          color="primary"
          variant="text"
          leftIcon={PlusIcon}
          onClick={this.onAddPort}>
          Add Port
        </Button>
      </CardSection>
    );
  }

  getEditablePortPropertyCell(portIndex, value, name, placeholder) {
    const {classes} = this.props;
    return (
      <TextField
        className={classes.input}
        name={name}
        fullWidth={true}
        placeholder={placeholder}
        variant="outlined"
        value={value ? value : ''}
        onChange={({target}) =>
          this.onPortPropertyChanged(name, target.value, portIndex)
        }
        margin="dense"
      />
    );
  }

  onAddPort = () => {
    const {onPortDefinitionsChanged} = this.props;
    onPortDefinitionsChanged &&
      onPortDefinitionsChanged(
        update(this.props.portDefinitions, {
          $push: [this.getEditingPort()],
        }),
      );
  };

  onPortPropertyChanged = (propertyName, newValue, portIndex) => {
    const {onPortDefinitionsChanged} = this.props;
    onPortDefinitionsChanged &&
      onPortDefinitionsChanged(
        update(this.props.portDefinitions, {
          // $FlowFixMe Set state for each field
          [portIndex]: {[propertyName]: {$set: newValue}},
        }),
      );
  };

  onRemovePortClicked = portIndex => {
    const {onPortDefinitionsChanged} = this.props;
    onPortDefinitionsChanged &&
      onPortDefinitionsChanged(
        update(this.props.portDefinitions, {$splice: [[portIndex, 1]]}),
      );
  };

  _onDragEnd = result => {
    if (!result.destination) {
      return;
    }

    const items = reorder(
      this.props.portDefinitions,
      result.source.index,
      result.destination.index,
    );

    const newItems = [];
    items.map((portDefinition, i) => {
      newItems.push(update(portDefinition, {index: {$set: i}}));
    });

    this.props.onPortDefinitionsChanged &&
      this.props.onPortDefinitionsChanged(newItems);
  };

  getEditingPort(): $Shape<
    $ElementType<PortDefinitionsAddEditTable_portDefinitions, number>,
  > {
    const index = this.props.portDefinitions.length;
    return {
      id: `PortDefinition@tmp-${index}-${Date.now()}`,
      name: '',
      index: index,
      visibleLabel: '',
      portType: null,
    };
  }

  async getEquipmentPortTypes(): Promise<Array<EquipmentPortType>> {
    const response = await fetchQuery(
      RelayEnvironment,
      equipmentPortTypesQuery,
    );
    return response.equipmentPortTypes.edges
      .map(edge => edge.node)
      .filter(Boolean)
      .sort((portTypeA, portTypeB) =>
        sortLexicographically(portTypeA.name, portTypeB.name),
      );
  }
}

export default withStyles(styles)(PortDefinitionsAddEditTable);
