/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {
  AddEquipmentTypeMutationResponse,
  AddEquipmentTypeMutationVariables,
} from '../../mutations/__generated__/AddEquipmentTypeMutation.graphql';
import type {
  EditEquipmentTypeMutationResponse,
  EditEquipmentTypeMutationVariables,
} from '../../mutations/__generated__/EditEquipmentTypeMutation.graphql';
import type {EquipmentType} from '../../common/EquipmentType';
import type {MutationCallbacks} from '../../mutations/MutationCallbacks.js';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithSnackbarProps} from 'notistack';
import type {WithStyles} from '@material-ui/core';

import AddEquipmentTypeMutation from '../../mutations/AddEquipmentTypeMutation';
import Button from '@fbcnms/ui/components/design-system/Button';
import CardSection from '../CardSection';
import EditEquipmentTypeMutation from '../../mutations/EditEquipmentTypeMutation';
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import Grid from '@material-ui/core/Grid';
import PageFooter from '@fbcnms/ui/components/PageFooter';
import PortDefinitionsAddEditTable from './PortDefinitionsAddEditTable';
import PositionDefinitionsAddEditTable from './PositionDefinitionsAddEditTable';
import PropertyTypeTable from '../form/PropertyTypeTable';
import React from 'react';
import SectionedCard from '@fbcnms/ui/components/SectionedCard';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import nullthrows from '@fbcnms/util/nullthrows';
import update from 'immutability-helper';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {ConnectionHandler} from 'relay-runtime';
import {createFragmentContainer, graphql} from 'react-relay';
import {getGraphError} from '../../common/EntUtils';
import {getPropertyDefaultValue} from '../../common/PropertyType';
import {sortByIndex} from '../draggable/DraggableUtils';
import {withSnackbar} from 'notistack';
import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  header: {
    marginBottom: '21px',
    paddingBottom: '0px',
  },
  input: {
    width: '305px',
  },
  section: {
    marginBottom: '28px',
  },
  headerText: {
    fontSize: '20px',
    lineHeight: '24px',
    fontWeight: 500,
  },
  closeButton: {
    marginRight: theme.spacing(),
  },
  cards: {
    height: 'calc(100% - 60px)',
    padding: '8px 24px',
    overflowY: 'auto',
  },
});

type Props = WithSnackbarProps &
  WithStyles<typeof styles> &
  WithAlert & {
    editingEquipmentType?: ?EquipmentType,
    onClose: () => void,
    onSave: (equipmentType: any) => void,
  };

type State = {
  error: ?string,
  editingEquipmentType: EquipmentType,
  isSaving: boolean,
};

class AddEditEquipmentTypeCard extends React.Component<Props, State> {
  state = {
    error: null,
    editingEquipmentType: this.getEditingEquipmentType(),
    isSaving: false,
  };

  _nameInputRef = React.createRef();

  componentDidMount() {
    this._nameInputRef.current && this._nameInputRef.current.focus();
  }

  render() {
    const {classes, onClose} = this.props;
    const {editingEquipmentType} = this.state;
    const propertyTypes = editingEquipmentType.propertyTypes
      .slice()
      .sort(sortByIndex);
    const positionDefinitions = editingEquipmentType.positionDefinitions
      .slice()
      .sort(sortByIndex)
      .map(x => Object.freeze(x));
    const portDefinitions = editingEquipmentType.portDefinitions
      .slice()
      .sort(sortByIndex)
      .map(x => Object.freeze(x));

    return (
      <>
        <div className={classes.cards}>
          <SectionedCard>
            <div className={classes.header}>
              <Text className={classes.headerText}>
                {this.props.editingEquipmentType
                  ? 'Edit Equipment Type'
                  : 'New Equipment Type'}
              </Text>
            </div>
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <FormField label="Name" required>
                  <TextInput
                    name="name"
                    variant="outlined"
                    className={classes.input}
                    value={this.state.editingEquipmentType.name}
                    onChange={this.handleNameChanged}
                    ref={this._nameInputRef}
                  />
                </FormField>
              </Grid>
            </Grid>
          </SectionedCard>
          <SectionedCard>
            <Grid container direction="column" spacing={3}>
              <Grid item xs={12} xl={7}>
                <CardSection className={classes.section} title="Properties">
                  <PropertyTypeTable
                    propertyTypes={propertyTypes}
                    onPropertiesChanged={propertyTypes => {
                      this.setState(state => ({
                        editingEquipmentType: {
                          ...state.editingEquipmentType,
                          propertyTypes,
                        },
                      }));
                    }}
                  />
                </CardSection>
              </Grid>
            </Grid>
          </SectionedCard>
          <SectionedCard>
            <Grid container direction="column" spacing={3}>
              <Grid item xs={12} xl={7}>
                <PositionDefinitionsAddEditTable
                  positionDefinitions={positionDefinitions}
                  onPositionDefinitionsChanged={positionDefinitions =>
                    this.setState(state => ({
                      editingEquipmentType: update(state.editingEquipmentType, {
                        positionDefinitions: {$set: positionDefinitions},
                      }),
                    }))
                  }
                />
              </Grid>
            </Grid>
          </SectionedCard>
          <SectionedCard>
            <Grid container direction="column" spacing={3}>
              <Grid item xs={12} xl={7}>
                <PortDefinitionsAddEditTable
                  // $FlowFixMe mix between relay and hand typed. Please fix.
                  portDefinitions={portDefinitions}
                  onPortDefinitionsChanged={ports =>
                    this.setState(state => ({
                      editingEquipmentType: update(state.editingEquipmentType, {
                        portDefinitions: {$set: ports},
                      }),
                    }))
                  }
                />
              </Grid>
            </Grid>
          </SectionedCard>
        </div>
        <PageFooter>
          <Button
            className={classes.closeButton}
            skin="regular"
            onClick={onClose}>
            Cancel
          </Button>
          <Button onClick={this.onSave} disabled={this.isSaveDisabled()}>
            Save
          </Button>
        </PageFooter>
      </>
    );
  }
  isSaveDisabled() {
    return (
      !this.state.editingEquipmentType.name ||
      this.state.isSaving ||
      !this.state.editingEquipmentType.propertyTypes.every(property => {
        return (
          property.isInstanceProperty || !!getPropertyDefaultValue(property)
        );
      })
    );
  }
  onSave = () => {
    const {
      name,
      positionDefinitions,
      portDefinitions,
    } = this.state.editingEquipmentType;
    const {enqueueSnackbar} = this.props;

    let error = null;
    if (!name) {
      error = 'Name cannot be empty';
    }

    const hasDuplicateNames = (arr: Array<string>) =>
      arr.length !== new Set(arr).size;

    if (hasDuplicateNames(positionDefinitions.map(p => p.name))) {
      error = 'Cannot have duplicate position names';
    }

    if (hasDuplicateNames(portDefinitions.map(p => p.name))) {
      error = 'Cannot have duplicate port names';
    }

    if (error !== null) {
      enqueueSnackbar(error, {
        children: key => (
          <SnackbarItem id={key} message={nullthrows(error)} variant="error" />
        ),
      });
      return;
    }

    this.setState({isSaving: true});
    if (this.props.editingEquipmentType) {
      this.editEquipmentType();
    } else {
      this.addNewEquipmentType();
    }
  };

  buildAddMutationVariables = (): AddEquipmentTypeMutationVariables => {
    const {id: _, ...addVars} = this.buildEditMutationVariables().input;
    return {input: {...addVars}};
  };

  buildEditMutationVariables = (): EditEquipmentTypeMutationVariables => {
    const {
      id,
      name,
      positionDefinitions,
      propertyTypes,
      portDefinitions,
    } = this.state.editingEquipmentType;

    const deleteTempId = (definition: Object) => {
      const newDef = {...definition};
      if (newDef.id && newDef.id.includes('@tmp')) {
        newDef.id = undefined;
      }
      return newDef;
    };

    const variables: EditEquipmentTypeMutationVariables = {
      input: {
        id: id,
        name: name,
        properties: propertyTypes
          .filter(propType => !!propType.name)
          .map(deleteTempId),
        positions: positionDefinitions
          .filter(definition => !!definition.name)
          .map(deleteTempId),
        ports: portDefinitions
          .filter(port => !!port.name)
          .map(portDefinition => ({
            ...portDefinition,
            portTypeID: portDefinition.portType?.id,
          }))
          .map(portDefinition => {
            delete portDefinition.portType;
            return portDefinition;
          })
          .map(deleteTempId),
      },
    };

    return variables;
  };

  editEquipmentType = () => {
    const variables: EditEquipmentTypeMutationVariables = this.buildEditMutationVariables();
    const callbacks: MutationCallbacks<EditEquipmentTypeMutationResponse> = {
      onCompleted: (response, errors) => {
        this.setState({isSaving: false});
        if (errors && errors[0]) {
          this.props.enqueueSnackbar(errors[0].message, {
            children: key => (
              <SnackbarItem
                id={key}
                message={errors[0].message}
                variant="error"
              />
            ),
          });
        } else {
          this.props.onSave && this.props.onSave(response.editEquipmentType);
        }
      },

      onError: (error: Error) => {
        this.setState({error: getGraphError(error), isSaving: false});
      },
    };

    EditEquipmentTypeMutation(variables, callbacks);
  };

  addNewEquipmentType = () => {
    const variables: AddEquipmentTypeMutationVariables = this.buildAddMutationVariables();
    const callbacks: MutationCallbacks<AddEquipmentTypeMutationResponse> = {
      onCompleted: (response, errors) => {
        this.setState({isSaving: false});
        if (errors && errors[0]) {
          this.props.enqueueSnackbar(errors[0].message, {
            children: key => (
              <SnackbarItem
                id={key}
                message={errors[0].message}
                variant="error"
              />
            ),
          });
        } else {
          this.props.onSave && this.props.onSave(response.addEquipmentType);
        }
      },

      onError: (error: Error) => {
        this.setState({error: getGraphError(error), isSaving: false});
      },
    };
    const updater = store => {
      // $FlowFixMe (T62907961) Relay flow types
      const rootQuery = store.getRoot();
      // $FlowFixMe (T62907961) Relay flow types
      const newNode = store.getRootField('addEquipmentType');
      if (!newNode) {
        return;
      }
      const types = ConnectionHandler.getConnection(
        rootQuery,
        'EquipmentTypes_equipmentTypes',
      );
      const edge = ConnectionHandler.createEdge(
        // $FlowFixMe (T62907961) Relay flow types
        store,
        // $FlowFixMe (T62907961) Relay flow types
        types,
        newNode,
        'EquipmentTypesEdge',
      );
      // $FlowFixMe (T62907961) Relay flow types
      ConnectionHandler.insertEdgeBefore(types, edge);
    };
    AddEquipmentTypeMutation(variables, callbacks, updater);
  };

  fieldChangedHandler = (field: 'name') => event =>
    this.setState({
      error: null,
      editingEquipmentType: {
        ...this.state.editingEquipmentType,
        [field]: event.target.value,
      },
    });

  handleNameChanged = this.fieldChangedHandler('name');

  getEditingEquipmentType(): EquipmentType {
    const editingEquipmentType = this.props.editingEquipmentType;
    return {
      id: editingEquipmentType?.id ?? 'tmp',
      name: editingEquipmentType?.name ?? '',
      positionDefinitions: [
        ...(editingEquipmentType?.positionDefinitions ?? []),
        {
          id: 'PositionDefinition@tmp',
          name: '',
          visibleLabel: '',
          index: editingEquipmentType?.positionDefinitions.length ?? 0,
        },
      ],
      portDefinitions: [
        ...(editingEquipmentType?.portDefinitions ?? []),
        {
          id: 'PortDefinition@tmp',
          name: '',
          visibleLabel: '',
          portType: null,
          index: editingEquipmentType?.portDefinitions.length ?? 0,
        },
      ],
      propertyTypes: [
        ...(editingEquipmentType?.propertyTypes ?? []),
        {
          id: 'PropertyType@tmp',
          name: '',
          type: 'string',
          index: editingEquipmentType?.propertyTypes.length ?? 0,
          booleanValue: false,
          stringValue: null,
          intValue: null,
          floatValue: null,
          latitudeValue: null,
          longitudeValue: null,
          isEditable: true,
          isInstanceProperty: true,
        },
      ],
      numberOfEquipment: 0,
    };
  }
}

export default withStyles(styles)(
  withAlert(
    withSnackbar(
      createFragmentContainer(AddEditEquipmentTypeCard, {
        editingEquipmentType: graphql`
          fragment AddEditEquipmentTypeCard_editingEquipmentType on EquipmentType {
            id
            name
            propertyTypes {
              id
              name
              type
              index
              stringValue
              intValue
              booleanValue
              floatValue
              latitudeValue
              longitudeValue
              isEditable
              isInstanceProperty
              isMandatory
            }
            positionDefinitions {
              ...PositionDefinitionsAddEditTable_positionDefinition
                @relay(mask: false)
            }
            portDefinitions {
              ...PortDefinitionsAddEditTable_portDefinitions @relay(mask: false)
            }
            numberOfEquipment
          }
        `,
      }),
    ),
  ),
);
