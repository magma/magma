/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {AddEditServiceTypeCard_editingServiceType} from './__generated__/AddEditServiceTypeCard_editingServiceType.graphql';
import type {
  AddServiceTypeMutationResponse,
  AddServiceTypeMutationVariables,
  ServiceTypeCreateData,
} from '../../mutations/__generated__/AddServiceTypeMutation.graphql';
import type {AppContextType} from '@fbcnms/ui/context/AppContext';
import type {
  EditServiceTypeMutationResponse,
  EditServiceTypeMutationVariables,
  ServiceTypeEditData,
} from '../../mutations/__generated__/EditServiceTypeMutation.graphql';
import type {MutationCallbacks} from '../../mutations/MutationCallbacks.js';
import type {PropertyType} from '../../common/PropertyType';
import type {
  ServiceEndpointDefinition,
  ServiceType,
} from '../../common/ServiceType';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithSnackbarProps} from 'notistack';
import type {WithStyles} from '@material-ui/core';

import AddServiceTypeMutation from '../../mutations/AddServiceTypeMutation';
import AppContext from '@fbcnms/ui/context/AppContext';
import Button from '@fbcnms/ui/components/design-system/Button';
import CardSection from '../CardSection';
import EditServiceTypeMutation from '../../mutations/EditServiceTypeMutation';
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import Grid from '@material-ui/core/Grid';
import MenuItem from '@material-ui/core/MenuItem';
import PageFooter from '@fbcnms/ui/components/PageFooter';
import PropertyTypeTable from '../form/PropertyTypeTable';
import React from 'react';
import SectionedCard from '@fbcnms/ui/components/SectionedCard';
import ServiceEndpointDefinitionTable from './ServiceEndpointDefinitionTable';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextField from '@material-ui/core/TextField';
import nullthrows from 'nullthrows';
import update from 'immutability-helper';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {ConnectionHandler} from 'relay-runtime';
import {FormContextProvider} from '../../common/FormContext';
import {createFragmentContainer, graphql} from 'react-relay';
import {discoveryMethods} from '../../common/Service';
import {generateTempId, removeTempIDs} from '../../common/EntUtils';
import {getGraphError} from '../../common/EntUtils';
import {getPropertyDefaultValue} from '../../common/PropertyType';
import {sortByIndex} from '../draggable/DraggableUtils';
import {withSnackbar} from 'notistack';
import {withStyles} from '@material-ui/core/styles';

const styles = _ => ({
  header: {
    marginBottom: '21px',
    paddingBottom: '0px',
  },
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '305px',
  },
  headerText: {
    fontSize: '20px',
    lineHeight: '24px',
    fontWeight: 500,
  },
  cards: {
    height: 'calc(100% - 60px)',
    padding: '8px 24px',
    overflowY: 'auto',
  },
  section: {
    marginBottom: '28px',
  },
  selectMenu: {
    height: '14px',
  },
});

type Props = {
  open: boolean,
  onClose: () => void,
  onSave: (serviceType: any) => void,
  editingServiceType?: AddEditServiceTypeCard_editingServiceType,
} & WithSnackbarProps &
  WithStyles<typeof styles> &
  WithAlert;

type State = {
  error: string,
  editingServiceType: ServiceType,
};

class AddEditServiceTypeCard extends React.Component<Props, State> {
  state = {
    error: '',
    editingServiceType: this.getEditingServiceType(),
  };

  static contextType = AppContext;
  context: AppContextType;

  render() {
    const {classes, onClose} = this.props;
    const {editingServiceType} = this.state;
    const serviceEndpointsEnabled = this.context.isFeatureEnabled(
      'service_endpoints',
    );

    const propertyTypes = editingServiceType.propertyTypes
      .slice()
      .sort(sortByIndex);
    const endpointDefinitions = editingServiceType.endpointDefinitions
      .slice()
      .sort(sortByIndex);

    const isOnEdit = !!this.props.editingServiceType;
    return (
      <FormContextProvider
        permissions={{
          entity: 'serviceType',
          action: isOnEdit ? 'update' : 'create',
        }}>
        <div className={classes.cards}>
          <SectionedCard>
            <div className={classes.header}>
              <Text className={classes.headerText}>
                {isOnEdit ? 'Edit Service Type' : 'New Service Type'}
              </Text>
            </div>
            <Grid container spacing={1}>
              <Grid item xs={12} xl={7}>
                <FormField label="Name" required>
                  <TextField
                    name="name"
                    variant="outlined"
                    margin="dense"
                    className={classes.input}
                    value={editingServiceType.name}
                    onChange={this.nameChanged}
                  />
                </FormField>
              </Grid>
              {serviceEndpointsEnabled ? (
                <Grid item xs={12} xl={7}>
                  <FormField label="Discovery Method" className={classes.input}>
                    <TextField
                      select
                      variant="outlined"
                      disabled={this.props.editingServiceType != null}
                      className={classes.input}
                      value={
                        editingServiceType.discoveryMethod ??
                        discoveryMethods.MANUAL
                      }
                      onChange={this.discoveryMethodChanged}
                      SelectProps={{
                        classes: {selectMenu: classes.selectMenu},
                        MenuProps: {
                          className: classes.menu,
                        },
                      }}
                      margin="dense">
                      {Object.keys(discoveryMethods).map(discoveryMethod => {
                        return (
                          <MenuItem
                            key={discoveryMethod}
                            value={discoveryMethod}>
                            {discoveryMethods[discoveryMethod] ?? ''}
                          </MenuItem>
                        );
                      })}
                    </TextField>
                  </FormField>
                </Grid>
              ) : null}
            </Grid>
          </SectionedCard>
          <SectionedCard>
            <Grid container direction="column" spacing={3}>
              <Grid item xs={12} xl={7}>
                <CardSection
                  className={classes.section}
                  title="Service Endpoints">
                  <ServiceEndpointDefinitionTable
                    serviceEndpointDefinitions={endpointDefinitions}
                    onServiceEndpointDefinitionsChanged={
                      this._endpointChangedHandler
                    }
                  />
                </CardSection>
              </Grid>
            </Grid>
          </SectionedCard>
          <SectionedCard>
            <Grid container direction="column" spacing={3}>
              <Grid item xs={12} xl={7}>
                <CardSection className={classes.section} title="Properties">
                  <PropertyTypeTable
                    propertyTypes={propertyTypes}
                    onPropertiesChanged={this._propertyChangedHandler}
                  />
                </CardSection>
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
          <Button disabled={this.isSaveDisabled()} onClick={this.onSave}>
            Save
          </Button>
        </PageFooter>
      </FormContextProvider>
    );
  }

  isSaveDisabled() {
    const serviceType = this.state.editingServiceType;
    const endpointWithNames = serviceType.endpointDefinitions.filter(
      obj => obj.name,
    );
    return (
      !serviceType.name ||
      !serviceType.propertyTypes.every(propType => {
        return (
          propType.isInstanceProperty || !!getPropertyDefaultValue(propType)
        );
      }) ||
      !endpointWithNames.every(
        endpointType => endpointType.equipmentType && endpointType.name,
      ) ||
      (serviceType.discoveryMethod != 'MANUAL' &&
        (endpointWithNames.length == 1 || endpointWithNames.length > 5))
    );
  }

  onSave = () => {
    const {name} = this.state.editingServiceType;
    if (!name) {
      this.setState({error: 'Name cannot be empty'});
      return;
    }

    if (this.props.editingServiceType) {
      this.editServiceType();
    } else {
      this.addNewServiceType();
    }
  };

  deleteTempId = (propType: PropertyType) => {
    if (propType.id && isNaN(propType.id) && propType.id.includes('@tmp')) {
      return {
        ...propType,
        id: undefined,
      };
    }
    return {...propType};
  };

  handleEquipmentOnEndpoint = (endpointType: ServiceEndpointDefinition) => {
    const obj = {
      ...endpointType,
      equipmentTypeID: nullthrows(endpointType.equipmentType?.id),
    };
    const {equipmentType: _, ...relevant} = obj;
    return relevant;
  };

  editServiceType = () => {
    const {
      id,
      name,
      propertyTypes,
      endpointDefinitions,
    } = this.state.editingServiceType;

    const data: ServiceTypeEditData = {
      id,
      name,
      hasCustomer: false,
      // $FlowFixMe property input doesn't have an id
      properties: propertyTypes
        .filter(propType => !!propType.name)
        .map(this.deleteTempId),
      endpoints: removeTempIDs(
        endpointDefinitions
          .filter(endpoint => !!endpoint.name)
          .map(this.handleEquipmentOnEndpoint),
      ),
    };

    const variables: EditServiceTypeMutationVariables = {
      data,
    };
    const callbacks: MutationCallbacks<EditServiceTypeMutationResponse> = {
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
          this.props.onSave && this.props.onSave(response.editServiceType);
          this.setState({error: ''});
        }
      },
      onError: (error: Error) => {
        this.setState({error: getGraphError(error)});
      },
    };

    EditServiceTypeMutation(variables, callbacks);
  };

  addNewServiceType = () => {
    const {
      name,
      propertyTypes,
      endpointDefinitions,
    } = this.state.editingServiceType;
    let discoveryMethod = null;
    if (this.state.editingServiceType.discoveryMethod == 'INVENTORY') {
      discoveryMethod = 'INVENTORY';
    }
    const data: ServiceTypeCreateData = {
      name,
      hasCustomer: false,
      discoveryMethod,
      // $FlowFixMe property input doesn't have an id
      properties: propertyTypes
        .filter(propType => !!propType.name)
        .map(this.deleteTempId),
      endpoints: removeTempIDs(
        endpointDefinitions
          .filter(endpoint => !!endpoint.name)
          .map(this.handleEquipmentOnEndpoint),
      ),
    };

    const variables: AddServiceTypeMutationVariables = {
      data,
    };
    const callbacks: MutationCallbacks<AddServiceTypeMutationResponse> = {
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
          this.props.onSave && this.props.onSave(response.addServiceType);
          this.setState({error: ''});
        }
      },
      onError: (error: Error) => {
        this.setState({error: getGraphError(error)});
      },
    };
    const updater = store => {
      // $FlowFixMe (T62907961) Relay flow types
      const rootQuery = store.getRoot();
      // $FlowFixMe (T62907961) Relay flow types
      const newNode = store.getRootField('addServiceType');
      if (!newNode) {
        return;
      }
      const types = ConnectionHandler.getConnection(
        rootQuery,
        'ServiceTypes_serviceTypes',
      );
      const edge = ConnectionHandler.createEdge(
        // $FlowFixMe (T62907961) Relay flow types
        store,
        // $FlowFixMe (T62907961) Relay flow types
        types,
        newNode,
        'ServiceTypesEdge',
      );
      // $FlowFixMe (T62907961) Relay flow types
      ConnectionHandler.insertEdgeBefore(types, edge);
    };

    AddServiceTypeMutation(variables, callbacks, updater);
  };

  fieldChangedHandler = (field: string) => event => {
    this.setState({
      editingServiceType: {
        ...this.state.editingServiceType,
        [field]: event.target.value,
      },
    });
  };

  nameChanged = this.fieldChangedHandler('name');

  discoveryMethodChanged = this.fieldChangedHandler('discoveryMethod');

  _propertyChangedHandler = propertyTypes =>
    this.setState(prevState => {
      return {
        error: '',
        editingServiceType: update(prevState.editingServiceType, {
          propertyTypes: {$set: propertyTypes},
        }),
      };
    });

  _endpointChangedHandler = newEndpointDefinitions =>
    this.setState(prevState => {
      return {
        error: '',
        editingServiceType: {
          ...prevState.editingServiceType,
          endpointDefinitions: newEndpointDefinitions,
        },
      };
    });

  getEditingServiceType(): ServiceType {
    const {editingServiceType} = this.props;
    const propertyTypes = (editingServiceType?.propertyTypes ?? [])
      .filter(Boolean)
      .map(p => ({
        id: p.id,
        name: p.name,
        index: p.index || 0,
        type: p.type,
        nodeType: p.nodeType,
        booleanValue: p.booleanValue,
        stringValue: p.stringValue,
        intValue: p.intValue,
        floatValue: p.floatValue,
        latitudeValue: p.latitudeValue,
        longitudeValue: p.longitudeValue,
        isEditable: p.isEditable,
        isMandatory: p.isMandatory,
        isInstanceProperty: p.isInstanceProperty,
      }));
    const endpointDefinitions = (editingServiceType?.endpointDefinitions ?? [])
      .filter(Boolean)
      .map(ep => ({
        id: ep.id,
        name: ep.name,
        index: ep.index,
        role: ep.role,
        equipmentType: {
          id: ep.equipmentType?.id,
          name: ep.equipmentType?.name,
        },
      }));

    return {
      id: editingServiceType?.id ?? 'ServiceType@tmp0',
      name: editingServiceType?.name ?? '',
      numberOfServices: editingServiceType?.numberOfServices ?? 0,
      discoveryMethod: editingServiceType?.discoveryMethod ?? 'MANUAL',
      propertyTypes:
        propertyTypes.length > 0
          ? propertyTypes
          : [
              {
                id: 'PropertyType@tmp',
                name: '',
                type: 'string',
                nodeType: null,
                index: editingServiceType?.propertyTypes?.length ?? 0,
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
      endpointDefinitions:
        endpointDefinitions.length > 0
          ? endpointDefinitions
          : [
              {
                id: generateTempId(),
                name: '',
                index: editingServiceType?.endpointDefinitions?.length ?? 0,
                role: null,
                equipmentType: null,
              },
            ],
    };
  }
}

export default withStyles(styles)(
  withAlert(
    withSnackbar(
      createFragmentContainer(AddEditServiceTypeCard, {
        editingServiceType: graphql`
          fragment AddEditServiceTypeCard_editingServiceType on ServiceType {
            id
            name
            numberOfServices
            discoveryMethod
            propertyTypes {
              id
              name
              type
              nodeType
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
            endpointDefinitions {
              ...ServiceEndpointDefinitionTable_serviceEndpointDefinitions
                @relay(mask: false)
            }
          }
        `,
      }),
    ),
  ),
);
