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
  AddEquipmentMutationResponse,
  AddEquipmentMutationVariables,
} from '../../mutations/__generated__/AddEquipmentMutation.graphql';
import type {AppContextType} from '@fbcnms/ui/context/AppContext';
import type {
  EditEquipmentMutationResponse,
  EditEquipmentMutationVariables,
} from '../../mutations/__generated__/EditEquipmentMutation.graphql';
import type {Equipment, EquipmentPosition} from '../../common/Equipment';
import type {EquipmentType} from '../../common/EquipmentType';
import type {MutationCallbacks} from '../../mutations/MutationCallbacks.js';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithSnackbarProps} from 'notistack';
import type {WithStyles} from '@material-ui/core';

import AddEquipmentMutation from '../../mutations/AddEquipmentMutation';
import AppContext from '@fbcnms/ui/context/AppContext';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import CardFooter from '@fbcnms/ui/components/CardFooter';
import CircularProgress from '@material-ui/core/CircularProgress';
import EditEquipmentMutation from '../../mutations/EditEquipmentMutation';
import FormLabel from '@material-ui/core/FormLabel';
import FormSaveCancelPanel from '@fbcnms/ui/components/design-system/Form/FormSaveCancelPanel';
import LinkedDeviceAddEditSection from '../form/LinkedDeviceAddEditSection';
import NameInput from '@fbcnms/ui/components/design-system/Form/NameInput';
import PropertiesAddEditSection from '../form/PropertiesAddEditSection';
import React from 'react';
import RelayEnvironment from '../../common/RelayEnvironment.js';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import Text from '@fbcnms/ui/components/design-system/Text';
import nullthrows from '@fbcnms/util/nullthrows';
import update from 'immutability-helper';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {FormContextProvider} from '../../common/FormContext';
import {fetchQuery, graphql} from 'relay-runtime';
import {getInitialPropertyFromType} from '../../common/PropertyType';
import {
  getNonInstancePropertyTypes,
  sortPropertiesByIndex,
  toPropertyInput,
} from '../../common/Property';
import {withSnackbar} from 'notistack';
import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  root: {
    height: '100%',
  },
  header: {
    minHeight: '50px',
  },
  loadingContainer: {
    minHeight: 500,
    paddingTop: 200,
    textAlign: 'center',
  },
  cancelButton: {
    marginRight: theme.spacing(),
  },
});

const equipmentAddEditCardQuery = graphql`
  query EquipmentAddEditCardQuery($equipmentId: ID!) {
    equipment: node(id: $equipmentId) {
      ... on Equipment {
        id
        name
        parentLocation {
          id
        }
        parentPosition {
          id
        }
        device {
          id
        }
        equipmentType {
          id
          name
          propertyTypes {
            id
            name
            index
            isInstanceProperty
            type
            isMandatory
            stringValue
            intValue
            floatValue
            booleanValue
            latitudeValue
            longitudeValue
            rangeFromValue
            rangeToValue
          }
        }
        properties {
          propertyType {
            id
            name
            index
            isInstanceProperty
            type
            stringValue
            isMandatory
          }
          id
          stringValue
          intValue
          floatValue
          booleanValue
          latitudeValue
          longitudeValue
          rangeFromValue
          rangeToValue
          equipmentValue {
            id
            name
          }
          locationValue {
            id
            name
          }
          serviceValue {
            id
            name
          }
        }
      }
    }
  }
`;

const equipmentAddEditCardQuery__equipmentTypeQuery = graphql`
  query EquipmentAddEditCardQuery__equipmentTypeQuery($equipmentTypeId: ID!) {
    equipmentType: node(id: $equipmentTypeId) {
      ... on EquipmentType {
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
          rangeFromValue
          rangeToValue
          isEditable
          isMandatory
          isInstanceProperty
        }
      }
    }
  }
`;

type Props = WithSnackbarProps &
  WithStyles<typeof styles> &
  WithAlert & {
    editingEquipmentId?: ?string,
    locationId: ?string,
    equipmentPosition: ?EquipmentPosition,
    workOrderId: ?string,
    type: ?EquipmentType,
    onCancel: () => void,
    onSave: () => void,
  };

type State = {
  editingEquipment: ?Equipment,
  error: string,
  isSubmitting: boolean,
};

class EquipmentAddEditCard extends React.Component<Props, State> {
  static contextType = AppContext;
  context: AppContextType;
  state = {
    editingEquipment: null,
    error: '',
    isSubmitting: false,
  };

  componentDidMount() {
    this.getEditingEquipment().then(editingEquipment => {
      this.setState({editingEquipment});
    });
  }

  render() {
    const {classes} = this.props;
    const {editingEquipment} = this.state;
    if (!editingEquipment) {
      return (
        <div className={classes.loadingContainer}>
          <CircularProgress size={50} />
        </div>
      );
    }
    const equipmentLiveStatusEnabled = this.context.isFeatureEnabled(
      'equipment_live_status',
    );
    return (
      <Card>
        <FormContextProvider>
          <CardContent className={this.props.classes.root}>
            {this.state.error && (
              <FormLabel error>{this.state.error}</FormLabel>
            )}
            <div className={this.props.classes.header}>
              <Text variant="h5">
                {editingEquipment?.equipmentType.name ?? this.props.type?.name}
              </Text>
            </div>
            <NameInput
              value={editingEquipment.name}
              onChange={this._onNameChanged}
              inputClass={classes.input}
            />
            {editingEquipment.properties.length > 0 ? (
              <PropertiesAddEditSection
                properties={editingEquipment.properties}
                onChange={index => this._propertyChangedHandler(index)}
              />
            ) : null}
            {this.props.editingEquipmentId && equipmentLiveStatusEnabled ? (
              <LinkedDeviceAddEditSection
                deviceID={editingEquipment.device?.id ?? ''}
                onChange={this._deviceIDChangedHandler}
              />
            ) : null}
          </CardContent>
          <CardFooter>
            <FormSaveCancelPanel
              isDisabled={this.state.isSubmitting}
              onCancel={this.props.onCancel}
              onSave={this.onSave}
            />
          </CardFooter>
        </FormContextProvider>
      </Card>
    );
  }

  isSaveDisabled() {
    return this.state.isSubmitting || !this.state.editingEquipment?.name;
  }

  onSave = () => {
    this.setState({isSubmitting: true});
    if (this.props.editingEquipmentId) {
      this._executeEdit();
    } else {
      this._executeAdd();
    }
  };

  _executeAdd = () => {
    if (!this.state.editingEquipment?.name) {
      this.setState({error: 'Name cannot be empty'});
      return;
    }
    const variables: AddEquipmentMutationVariables = {
      input: {
        name: this.state.editingEquipment?.name ?? '',
        location: this.props.locationId,
        parent: this.props.equipmentPosition?.parentEquipment.id,
        positionDefinition: this.props.equipmentPosition?.definition.id,
        type: nullthrows(this.props.type?.id),
        workOrder: this.props.workOrderId,
        properties:
          (this.state.editingEquipment &&
            toPropertyInput(this.state.editingEquipment.properties)) ??
          [],
      },
    };

    const callbacks: MutationCallbacks<AddEquipmentMutationResponse> = {
      onCompleted: (response, errors) => {
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
          this.props.onSave && this.props.onSave();
        }

        this.setState({isSubmitting: false});
      },
      onError: (error: Error) => {
        this.setState({error: error.message});
      },
    };
    AddEquipmentMutation(variables, callbacks);
  };

  _executeEdit = () => {
    const editingEquipment = nullthrows(this.state.editingEquipment);
    if (!editingEquipment.name) {
      this.setState({error: 'Name cannot be empty'});
      return;
    }
    const variables: EditEquipmentMutationVariables = {
      input: {
        id: editingEquipment.id,
        name: editingEquipment.name,
        deviceID: editingEquipment.device?.id,
        properties: toPropertyInput(editingEquipment.properties),
      },
    };
    const callbacks: MutationCallbacks<EditEquipmentMutationResponse> = {
      onCompleted: (response, errors) => {
        if (errors && errors[0]) {
          this.setState({error: errors[0].message});
        } else {
          this.props.onSave && this.props.onSave();
        }
        this.setState({isSubmitting: false});
      },
      onError: (error: Error) => {
        this.setState({error: error.message});
      },
    };

    EditEquipmentMutation(variables, callbacks);
  };

  async getEditingEquipment(): Promise<Equipment> {
    let equipment = null;
    let equipmentType = null;
    if (this.props.editingEquipmentId) {
      const response = await fetchQuery(
        RelayEnvironment,
        equipmentAddEditCardQuery,
        {
          equipmentId: this.props.editingEquipmentId,
        },
      );
      equipment = response.equipment;
      equipmentType = equipment.equipmentType;
    } else {
      const response = await fetchQuery(
        RelayEnvironment,
        equipmentAddEditCardQuery__equipmentTypeQuery,
        {
          equipmentTypeId: nullthrows(this.props.type?.id),
        },
      );
      equipmentType = response.equipmentType;
    }

    let initialProps = equipment?.properties ?? [];
    if (equipmentType && equipmentType.propertyTypes) {
      initialProps = [
        ...initialProps,
        ...getNonInstancePropertyTypes(
          initialProps,
          equipmentType.propertyTypes,
        ).map(propType => getInitialPropertyFromType(propType)),
      ];
      initialProps = initialProps.sort(sortPropertiesByIndex);
    }

    return {
      id: equipment?.id ?? 'Equipment@tmp',
      name: equipment?.name ?? '',
      equipmentType: equipmentType,
      properties: initialProps,
      positions: equipment?.positions ?? [],
      ports: equipment?.ports ?? [],
      parentLocation: equipment?.parentLocation,
      parentPosition: equipment?.parentPosition,
      futureState: equipment?.futureState,
      device: equipment?.device,
      workOrder: equipment?.workOrder,
      locationHierarchy: equipment?.locationHierarchy ?? [],
      positionHierarchy: equipment?.positionHierarchy ?? [],
      services: equipment?.services ?? [],
    };
  }

  _fieldChangedHandler = (field: 'name') => event =>
    this.setState({
      error: '',
      editingEquipment: update(this.state.editingEquipment, {
        [field]: {$set: event.target.value},
      }),
    });

  _propertyChangedHandler = index => property =>
    this.setState(prevState => {
      return {
        error: '',
        editingEquipment: update(prevState.editingEquipment, {
          properties: {[index]: {$set: property}},
        }),
      };
    });

  _deviceIDChangedHandler = (deviceID: string) => {
    this.setState(prevState => {
      return {
        error: '',
        editingEquipment: update(prevState.editingEquipment, {
          device: {$set: {id: deviceID}},
        }),
      };
    });
  };

  _onNameChanged = this._fieldChangedHandler('name');
}

export default withStyles(styles)(
  withAlert(withSnackbar(EquipmentAddEditCard)),
);
