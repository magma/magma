/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ServiceEndpointDefinition} from '../../common/ServiceType';
import type {WithStyles} from '@material-ui/core';

import Button from '@fbcnms/ui/components/design-system/Button';

import CircularProgress from '@material-ui/core/CircularProgress';
import DraggableTableRow from '../draggable/DraggableTableRow';
import DroppableTableBody from '../draggable/DroppableTableBody';
import IconButton from '@fbcnms/ui/components/design-system/IconButton';
import ListItemText from '@material-ui/core/ListItemText';
import MenuItem from '@material-ui/core/MenuItem';
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
import fbt from 'fbt';
import inventoryTheme from '../../common/theme';
import update from 'immutability-helper';
import {DeleteIcon, PlusIcon} from '@fbcnms/ui/components/design-system/Icons';
import {fetchQuery, graphql} from 'react-relay';
import {generateTempId, isTempId} from '../../common/EntUtils';
import {removeItem, updateItem} from '@fbcnms/util/arrays';
import {reorder} from '../draggable/DraggableUtils';
import {sortLexicographically} from '@fbcnms/ui/utils/displayUtils';
import {withStyles} from '@material-ui/core/styles';

const styles = _theme => ({
  container: {
    maxWidth: '1366px',
    overflowX: 'auto',
  },
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

type EquipmentTypeOption = {
  name: string,
  id: string,
};

graphql`
  fragment ServiceEndpointDefinitionTable_serviceEndpointDefinitions on ServiceEndpointDefinition
    @relay(plural: true) {
    id
    index
    role
    name
    equipmentType {
      name
      id
    }
  }
`;

const equipmentTypesQuery = graphql`
  query ServiceEndpointDefinitionTable_equipmentTypesQuery {
    equipmentTypes(first: 500)
      @connection(key: "ServiceEndpointDefinitionTable_equipmentTypes") {
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
  serviceEndpointDefinitions: Array<ServiceEndpointDefinition>,
  onServiceEndpointDefinitionsChanged?: ?(
    newEndpointTypes: Array<ServiceEndpointDefinition>,
  ) => void,
} & WithStyles<typeof styles>;

type State = {
  equipmentTypes: Array<EquipmentTypeOption>,
};

class ServiceEndpointDefinitionTable extends React.Component<Props, State> {
  state = {
    equipmentTypes: [],
  };

  componentDidMount() {
    this.getEquipmentTypes().then(equipmentTypes => {
      this.setState({
        equipmentTypes,
      });
    });
  }

  render() {
    const {serviceEndpointDefinitions, classes} = this.props;
    if (serviceEndpointDefinitions.length === 0) {
      return null;
    }
    const {equipmentTypes} = this.state;
    if (!equipmentTypes) {
      return (
        <div className={classes.loadingContainer}>
          <CircularProgress size={50} />
        </div>
      );
    }

    return (
      <div className={classes.container}>
        <Table component="div" className={classes.table}>
          <TableHead component="div">
            <TableRow component="div">
              <TableCell component="div" size="small" padding="checkbox" />
              <TableCell component="div" className={classes.cell}>
                <fbt desc="">Name</fbt>
              </TableCell>
              <TableCell component="div" className={classes.cell}>
                <fbt desc="">Endpoint Function</fbt>
              </TableCell>
              <TableCell component="div" className={classes.cell}>
                <fbt desc="">Equipment Type</fbt>
              </TableCell>
              <TableCell component="div" />
            </TableRow>
          </TableHead>
          <DroppableTableBody onDragEnd={this._onDragEnd}>
            {serviceEndpointDefinitions.map((serviceEndpointDefinition, i) => (
              <DraggableTableRow
                key={i}
                id={serviceEndpointDefinition.id}
                index={i}>
                <TableCell className={classes.cell} component="div" scope="row">
                  {this.getEditableCell(
                    i,
                    serviceEndpointDefinition.name,
                    'name',
                    'Name',
                    this.onNameChange,
                  )}
                </TableCell>
                <TableCell
                  component="div"
                  className={classes.cell}
                  align="left">
                  {this.getEditableCell(
                    i,
                    serviceEndpointDefinition.role,
                    'role',
                    'Role',
                    this.onRoleChange,
                  )}
                </TableCell>
                <TableCell className={classes.cell} component="div" scope="row">
                  {isTempId(serviceEndpointDefinition.id) ? (
                    <Select
                      value={serviceEndpointDefinition.equipmentType?.id || ''}
                      input={<OutlinedInput margin="dense" />}
                      onChange={({target}) =>
                        this.onEndpointEquipmentTypeChanged(target.value, i)
                      }
                      MenuProps={{
                        className: classes.menu,
                      }}
                      margin="dense">
                      {equipmentTypes.map(equipmentType => (
                        <MenuItem value={equipmentType.id}>
                          <ListItemText>{equipmentType.name}</ListItemText>
                        </MenuItem>
                      ))}
                    </Select>
                  ) : (
                    serviceEndpointDefinition.equipmentType?.name
                  )}
                </TableCell>
                <TableCell component="div" align="right">
                  <IconButton
                    onClick={() =>
                      this.onRemoveEndpointClicked(i, serviceEndpointDefinition)
                    }
                    disabled={!isTempId(serviceEndpointDefinition.id)}
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
          onClick={this.onAddEndpoint}>
          <fbt desc="">Add Endpoint</fbt>
        </Button>
      </div>
    );
  }

  _onDragEnd = result => {
    if (!result.destination) {
      return;
    }

    const items = reorder(
      this.props.serviceEndpointDefinitions,
      result.source.index,
      result.destination.index,
    );

    const newItems = [];
    items.map((endpoint, i) => {
      newItems.push(update(endpoint, {index: {$set: i}}));
    });

    this.props.onServiceEndpointDefinitionsChanged &&
      this.props.onServiceEndpointDefinitionsChanged(newItems);
  };

  getEditingEndpoint(): ServiceEndpointDefinition {
    const index = this.props.serviceEndpointDefinitions.length;
    return {
      id: generateTempId(),
      name: '',
      role: null,
      index: index,
      equipmentType: null,
    };
  }

  async getEquipmentTypes(): Promise<Array<EquipmentTypeOption>> {
    const response = await fetchQuery(RelayEnvironment, equipmentTypesQuery);
    return response.equipmentTypes.edges
      .map(edge => edge.node)
      .filter(Boolean)
      .sort((typ1, typ2) => sortLexicographically(typ1.name, typ2.name));
  }

  onAddEndpoint = () => {
    const {onServiceEndpointDefinitionsChanged} = this.props;
    onServiceEndpointDefinitionsChanged &&
      onServiceEndpointDefinitionsChanged([
        ...this.props.serviceEndpointDefinitions,
        this.getEditingEndpoint(),
      ]);
  };

  getEditableCell(
    index,
    value,
    name,
    placeholder,
    onFieldChange: (string, number) => void,
  ) {
    const {classes} = this.props;
    return (
      <TextField
        className={classes.input}
        name={name}
        fullWidth={true}
        placeholder={placeholder}
        variant="outlined"
        value={value ?? ''}
        onChange={({target}) => onFieldChange(target.value, index)}
        margin="dense"
      />
    );
  }

  onEndpointEquipmentTypeChanged = (equipmentID, index) => {
    const {onServiceEndpointDefinitionsChanged} = this.props;
    const equipmentObj = this.state.equipmentTypes.find(
      obj => obj.id == equipmentID,
    );

    equipmentObj &&
      onServiceEndpointDefinitionsChanged &&
      onServiceEndpointDefinitionsChanged(
        updateItem<ServiceEndpointDefinition, 'equipmentType'>(
          this.props.serviceEndpointDefinitions,
          index,
          'equipmentType',
          equipmentObj,
        ),
      );
  };

  onNameChange = (value, index) => {
    const {onServiceEndpointDefinitionsChanged} = this.props;
    onServiceEndpointDefinitionsChanged &&
      onServiceEndpointDefinitionsChanged(
        updateItem<ServiceEndpointDefinition, 'name'>(
          this.props.serviceEndpointDefinitions,
          index,
          'name',
          value,
        ),
      );
  };

  onRoleChange = (value, index) => {
    const {onServiceEndpointDefinitionsChanged} = this.props;
    onServiceEndpointDefinitionsChanged &&
      onServiceEndpointDefinitionsChanged(
        updateItem<ServiceEndpointDefinition, 'role'>(
          this.props.serviceEndpointDefinitions,
          index,
          'role',
          value,
        ),
      );
  };

  onRemoveEndpointClicked = (index, endpoint: ServiceEndpointDefinition) => {
    const {onServiceEndpointDefinitionsChanged} = this.props;
    if (isTempId(endpoint.id)) {
      onServiceEndpointDefinitionsChanged &&
        onServiceEndpointDefinitionsChanged(
          removeItem(this.props.serviceEndpointDefinitions, index),
        );
    }
  };
}

export default withStyles(styles)(ServiceEndpointDefinitionTable);
