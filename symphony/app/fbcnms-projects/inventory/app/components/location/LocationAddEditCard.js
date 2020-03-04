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
  AddLocationMutationResponse,
  AddLocationMutationVariables,
} from '../../mutations/__generated__/AddLocationMutation.graphql';
import type {AppContextType} from '@fbcnms/ui/context/AppContext';
import type {
  EditLocationMutationResponse,
  EditLocationMutationVariables,
} from '../../mutations/__generated__/EditLocationMutation.graphql';
import type {Location} from '../../common/Location';
import type {LocationType} from '../../common/LocationType';
import type {MutationCallbacks} from '../../mutations/MutationCallbacks.js';
import type {Theme} from '@material-ui/core';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithSnackbarProps} from 'notistack';
import type {WithStyles} from '@material-ui/core';

import AddLocationMutation from '../../mutations/AddLocationMutation';
import AppContext from '@fbcnms/ui/context/AppContext';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import CardFooter from '@fbcnms/ui/components/CardFooter';
import CircularProgress from '@material-ui/core/CircularProgress';
import EditLocationMutation from '../../mutations/EditLocationMutation';
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import FormSaveCancelPanel from '@fbcnms/ui/components/design-system/Form/FormSaveCancelPanel';
import GPSPropertyValueInput from '../form/GPSPropertyValueInput';
import Grid from '@material-ui/core/Grid';
import NameInput from '@fbcnms/ui/components/design-system/Form/NameInput';
import PropertiesAddEditSection from '../form/PropertiesAddEditSection';
import React from 'react';
import RelayEnvironment from '../../common/RelayEnvironment.js';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import nullthrows from '@fbcnms/util/nullthrows';
import update from 'immutability-helper';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {ConnectionHandler, fetchQuery, graphql} from 'relay-runtime';
import {FormValidationContextProvider} from '@fbcnms/ui/components/design-system/Form/FormValidationContext';
import {getInitialPropertyFromType} from '../../common/PropertyType';
import {
  getNonInstancePropertyTypes,
  sortPropertiesByIndex,
  toPropertyInput,
} from '../../common/Property';
import {withSnackbar} from 'notistack';
import {withStyles} from '@material-ui/core/styles';

const styles = (theme: Theme) => ({
  root: {
    height: '100%',
  },
  header: {
    marginBottom: '16px',
  },
  input: {
    display: 'inline-flex',
  },
  loadingContainer: {
    minHeight: 500,
    paddingTop: 200,
    textAlign: 'center',
  },
  cancelButton: {
    marginRight: theme.spacing(),
  },
  row: {
    marginBottom: theme.spacing(2),
  },
});

type Props = WithSnackbarProps &
  WithStyles<typeof styles> &
  WithAlert & {
    editingLocationId?: ?string,
    parentId: ?string,
    type: ?LocationType,
    onCancel: () => void,
    onSave: (locationId: string) => void,
  };

type State = {
  editingLocation: ?Location,
  error: string,
  isSubmitting: boolean,
};

const locationAddEditCardQuery = graphql`
  query LocationAddEditCardQuery($locationId: ID!) {
    location: node(id: $locationId) {
      ... on Location {
        id
        name
        latitude
        longitude
        externalId
        locationType {
          id
          name
          mapType
          mapZoomLevel
          propertyTypes {
            id
            name
            index
            isInstanceProperty
            type
            stringValue
            intValue
            floatValue
            booleanValue
            latitudeValue
            longitudeValue
            rangeFromValue
            rangeToValue
            isMandatory
          }
        }
        equipments {
          id
          name
          equipmentType {
            id
            name
          }
        }
        properties {
          id
          propertyType {
            id
            name
            type
            index
            isEditable
            isInstanceProperty
            stringValue
            isMandatory
          }
          stringValue
          intValue
          booleanValue
          floatValue
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

const locationAddEditCard__locationTypeQuery = graphql`
  query LocationAddEditCard__locationTypeQuery($locationTypeId: ID!) {
    locationType: node(id: $locationTypeId) {
      ... on LocationType {
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
          isInstanceProperty
          isMandatory
        }
      }
    }
  }
`;

class LocationAddEditCard extends React.Component<Props, State> {
  static contextType = AppContext;
  context: AppContextType;

  componentDidMount() {
    this.getEditingLocation().then(editingLocation => {
      this.setState({
        editingLocation,
      });
    });
  }

  state = {
    editingLocation: null,
    error: '',
    isSubmitting: false,
  };

  render() {
    const {classes} = this.props;
    const {editingLocation} = this.state;
    const externalIDEnabled = this.context.isFeatureEnabled('external_id');
    if (!editingLocation) {
      return (
        <div className={classes.loadingContainer}>
          <CircularProgress size={50} />
        </div>
      );
    }
    const {latitude, longitude, properties} = editingLocation;
    return (
      <Card>
        <FormValidationContextProvider>
          <CardContent className={this.props.classes.root}>
            <div className={this.props.classes.header}>
              <Text variant="h5">{editingLocation.locationType.name}</Text>
            </div>
            <Grid container spacing={2} className={classes.row}>
              <Grid item xs={12} sm={12} lg={6} xl={4}>
                <NameInput
                  value={editingLocation.name}
                  onChange={this._onNameChanged}
                  inputClass={classes.input}
                />
              </Grid>
              {externalIDEnabled && (
                <Grid item xs={12} sm={12} lg={6} xl={4}>
                  <FormField
                    label="External ID"
                    hasSpacer
                    className={classes.externalIdFormField}>
                    <TextInput
                      className={classes.input}
                      type="string"
                      value={editingLocation.externalId ?? ''}
                      onChange={this._onExternalIdChanged}
                    />
                  </FormField>
                </Grid>
              )}
            </Grid>
            <Grid container spacing={0} className={classes.row}>
              <Grid item xs={12} sm={12} lg={6} xl={4}>
                <GPSPropertyValueInput
                  label="Location"
                  margin="normal"
                  fullWidth
                  className={classes.input}
                  value={{
                    latitude,
                    longitude,
                    accuracy: 0,
                    altitude: 0,
                    altitudeAccuracy: 0,
                  }}
                  onLatitudeChange={this._onLatitudeChanged}
                  onLongitudeChange={this._onLongitudeChanged}
                />
              </Grid>
            </Grid>
            {properties.length > 0 ? (
              <PropertiesAddEditSection
                properties={properties}
                onChange={index => this._propertyChangedHandler(index)}
              />
            ) : null}
          </CardContent>
          <CardFooter>
            <FormSaveCancelPanel
              onCancel={this.props.onCancel}
              onSave={this.onSave}
            />
          </CardFooter>
        </FormValidationContextProvider>
      </Card>
    );
  }

  onSave = () => {
    this.setState({isSubmitting: true});
    if (this.props.editingLocationId) {
      this._executeEdit();
    } else {
      this._executeAdd();
    }
  };

  _executeAdd = () => {
    const editingLocation = nullthrows(this.state.editingLocation);
    if (!editingLocation.name) {
      this.setState({error: 'Name cannot be empty'});
      return;
    }

    const variables: AddLocationMutationVariables = {
      input: {
        name: editingLocation.name,
        latitude: editingLocation.latitude ?? 0,
        longitude: editingLocation.longitude ?? 0,
        externalID: editingLocation.externalId,
        parent: this.props.parentId,
        type: nullthrows(this.props.type?.id),
        properties: toPropertyInput(editingLocation.properties),
      },
    };

    const callbacks: MutationCallbacks<AddLocationMutationResponse> = {
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
          this.props.onSave &&
            this.props.onSave(nullthrows(response.addLocation).id);
        }
        this.setState({isSubmitting: false});
      },
      onError: (error: Error) => {
        this.setState({error: error.message});
      },
    };

    const updater = store => {
      // $FlowFixMe (T62907961) Relay flow types
      const newNode = store.getRootField('addLocation');
      if (newNode === null) {
        return;
      }

      const parentId = this.props.parentId;
      if (!parentId) {
        // $FlowFixMe (T62907961) Relay flow types
        const rootQuery = store.getRoot();
        const locations = ConnectionHandler.getConnection(
          rootQuery,
          'LocationsTree_locations',
          {onlyTopLevel: true},
        );
        const edge = ConnectionHandler.createEdge(
          // $FlowFixMe (T62907961) Relay flow types
          store,
          // $FlowFixMe (T62907961) Relay flow types
          locations,
          newNode,
          'LocationsEdge',
        );
        // $FlowFixMe (T62907961) Relay flow types
        ConnectionHandler.insertEdgeAfter(locations, edge);
        return;
      }

      // $FlowFixMe (T62907961) Relay flow types
      const parentProxy = store.get(parentId);
      // $FlowFixMe (T62907961) Relay flow types
      const currNodes = parentProxy.getLinkedRecords('children');
      const parentLoaded =
        currNodes !== null &&
        // $FlowFixMe (T62907961) Relay flow types
        (currNodes.length === 0 || !!currNodes.find(node => node != undefined));
      if (parentLoaded) {
        // $FlowFixMe (T62907961) Relay flow types
        parentProxy.setLinkedRecords([...currNodes, newNode], 'children');
        // $FlowFixMe (T62907961) Relay flow types
        parentProxy.setValue(
          // $FlowFixMe (T62907961) Relay flow types
          parentProxy.getValue('numChildren') + 1,
          'numChildren',
        );
      }
    };

    AddLocationMutation(variables, callbacks, updater);
  };

  _executeEdit = () => {
    const editingLocation = nullthrows(this.state.editingLocation);
    if (!editingLocation.name) {
      this.setState({error: 'Name cannot be empty'});
      return;
    }

    const variables: EditLocationMutationVariables = {
      input: {
        id: nullthrows(this.props.editingLocationId),
        name: editingLocation.name,
        externalID: editingLocation.externalId,
        latitude: editingLocation.latitude ?? 0,
        longitude: editingLocation.longitude ?? 0,
        properties: toPropertyInput(editingLocation.properties),
      },
    };
    const callbacks: MutationCallbacks<EditLocationMutationResponse> = {
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
          this.props.onSave &&
            this.props.onSave(nullthrows(response.editLocation).id);
        }
        this.setState({isSubmitting: false});
      },
      onError: (error: Error) => {
        this.setState({error: error.message});
      },
    };

    EditLocationMutation(variables, callbacks);
  };

  async getEditingLocation(): any {
    let location = null;
    let locationType = null;
    if (this.props.editingLocationId) {
      const response = await fetchQuery(
        RelayEnvironment,
        locationAddEditCardQuery,
        {
          locationId: this.props.editingLocationId,
        },
      );
      location = response.location;
      locationType = location.locationType;
    } else {
      const response = await fetchQuery(
        RelayEnvironment,
        locationAddEditCard__locationTypeQuery,
        {
          locationTypeId: nullthrows(this.props.type?.id),
        },
      );
      locationType = response.locationType;
    }
    let initialProps = location?.properties ?? [];
    if (locationType && locationType.propertyTypes) {
      initialProps = [
        ...initialProps,
        ...getNonInstancePropertyTypes(
          initialProps,
          locationType.propertyTypes,
        ).map(propType => getInitialPropertyFromType(propType)),
      ];
      initialProps = initialProps.sort(sortPropertiesByIndex);
    }

    return {
      id: location?.id ?? 'Location@tmp',
      name: location?.name ?? '',
      locationType: locationType,
      externalId: location?.externalId,
      properties: initialProps,
      equipments: location?.equipments ?? [],
      children: location?.children ?? [],
      latitude: location?.latitude,
      longitude: location?.longitude,
    };
  }

  floatFieldChangedHandler = (field: 'latitude' | 'longitude') => event =>
    this.setState({
      error: '',
      editingLocation: {
        ...this.state.editingLocation,
        // $FlowFixMe Set state for each field
        [field]: parseFloat(event.target.value),
      },
    });

  fieldChangedHandler = (field: 'name' | 'externalId') => event =>
    this.setState({
      error: '',
      editingLocation: {
        ...nullthrows(this.state.editingLocation),
        // $FlowFixMe Set state for each field
        [field]: event.target.value,
      },
    });

  _onNameChanged = this.fieldChangedHandler('name');
  _onExternalIdChanged = this.fieldChangedHandler('externalId');
  _onLatitudeChanged = this.floatFieldChangedHandler('latitude');
  _onLongitudeChanged = this.floatFieldChangedHandler('longitude');

  _propertyChangedHandler = index => property =>
    this.setState(prevState => {
      return {
        error: '',
        editingLocation: update(prevState.editingLocation, {
          properties: {[index]: {$set: property}},
        }),
      };
    });
}

export default withStyles(styles)(withAlert(withSnackbar(LocationAddEditCard)));
