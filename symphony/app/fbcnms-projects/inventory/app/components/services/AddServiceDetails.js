/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {AddServiceDetailsServiceTypeQuery} from './__generated__/AddServiceDetailsServiceTypeQuery.graphql';
import type {
  AddServiceMutationResponse,
  AddServiceMutationVariables,
} from '../../mutations/__generated__/AddServiceMutation.graphql';
import type {MutationCallbacks} from '../../mutations/MutationCallbacks.js';
import type {Service} from '../../common/Service';

import * as React from 'react';
import AddServiceMutation from '../../mutations/AddServiceMutation';
import Button from '@fbcnms/ui/components/design-system/Button';
import CustomerTypeahead from '../typeahead/CustomerTypeahead';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import Grid from '@material-ui/core/Grid';
import PropertyValueInput from '../form/PropertyValueInput';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import nullthrows from '@fbcnms/util/nullthrows';
import symphony from '@fbcnms/ui/theme/symphony';
import update from 'immutability-helper';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {getInitialPropertyFromType} from '../../common/PropertyType';
import {graphql, useLazyLoadQuery} from 'react-relay/hooks';
import {makeStyles} from '@material-ui/styles';
import {sortPropertiesByIndex} from '../../common/Property';
import {toPropertyInput} from '../../common/Property';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useState} from 'react';

const useStyles = makeStyles(_ => ({
  separator: {
    borderBottom: `1px solid ${symphony.palette.separator}`,
    margin: '0 0 24px 0px',
    paddingBottom: '24px',
  },
  input: {
    width: '100%',
    paddingBottom: '24px',
    marginLeft: '0px',
  },
  propInput: {
    width: '100%',
    paddingBottom: '0px',
    marginLeft: '0px',
  },
  detailInput: {
    display: 'inline-flex',
  },
  contentRoot: {
    marginLeft: '24px',
    marginRight: '24px',
  },
  loading: {
    display: 'flex',
    height: '100%',
    alignItems: 'center',
    justifyContent: 'center',
  },
  dialogTitle: {
    padding: '24px',
    paddingBottom: '16px',
  },
  serviceCreateDialogContent: {
    padding: 0,
    maxHeight: '500px',
    overflowY: 'auto',
  },
  dialogActions: {
    padding: '24px',
    bottom: 0,
    display: 'flex',
    justifyContent: 'flex-end',
    width: '100%',
    backgroundColor: symphony.palette.white,
  },
}));

const serviceTypeQuery = graphql`
  query AddServiceDetailsServiceTypeQuery($serviceTypeId: ID!) {
    serviceType(id: $serviceTypeId) {
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
      }
    }
  }
`;

type Props = {
  serviceTypeId: string,
  onBackClicked: () => void,
  onServiceCreated: (id: string) => void,
};

const AddServiceDetails = (props: Props) => {
  const {serviceTypeId, onBackClicked, onServiceCreated} = props;
  const [serviceState, setServiceState] = useState<?Service>(null);
  const classes = useStyles();
  const enqueueSnackbar = useEnqueueSnackbar();

  const data = useLazyLoadQuery<AddServiceDetailsServiceTypeQuery>(
    serviceTypeQuery,
    {
      serviceTypeId: serviceTypeId,
    },
  );

  const getService = () => {
    if (!serviceState) {
      const serviceType = data.serviceType;
      const initialProps = (serviceType.propertyTypes || [])
        .map(propType => getInitialPropertyFromType(propType))
        .sort(sortPropertiesByIndex);
      const service = {
        id: 'service@tmp',
        name: '',
        externalId: null,
        customer: null,
        serviceType: serviceType,
        properties: initialProps,
        upstream: [],
        downstream: [],
        terminationPoints: [],
        links: [],
      };
      setServiceState(service);
      return service;
    }
    return serviceState;
  };

  const service = getService();

  const propertyChangedHandler = index => property => {
    setServiceState(
      update(service, {
        properties: {[index]: {$set: property}},
      }),
    );
  };

  const isSaveDisabled = () => {
    return !service?.name;
  };

  const enqueueError = (message: string) => {
    enqueueSnackbar(message, {
      children: key => (
        <SnackbarItem id={key} message={message} variant="error" />
      ),
    });
  };

  const saveService = () => {
    const {name, externalId, customer, properties} = nullthrows(service);
    const serviceTypeId = nullthrows(service?.serviceType.id);
    const variables: AddServiceMutationVariables = {
      data: {
        name,
        externalId,
        serviceTypeId,
        customerId: customer?.id,
        properties: toPropertyInput(properties),
        upstreamServiceIds: [],
        terminationPointIds: [],
      },
    };

    const callbacks: MutationCallbacks<AddServiceMutationResponse> = {
      onCompleted: (response, errors) => {
        if (errors && errors[0]) {
          enqueueError(errors[0].message);
        } else {
          // navigate to main page
          onServiceCreated(nullthrows(response.addService?.id));
        }
      },
      onError: () => {
        enqueueError('Error saving service');
      },
    };
    ServerLogger.info(LogEvents.SAVE_SERVICE_BUTTON_CLICKED, {
      source: 'service_details',
    });
    AddServiceMutation(variables, callbacks);
  };

  return (
    <>
      <DialogTitle className={classes.dialogTitle}>
        <Text variant="h6">{service.serviceType.name}</Text>
      </DialogTitle>
      <DialogContent className={classes.serviceCreateDialogContent}>
        <div className={classes.contentRoot}>
          <div>
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <FormField label="Name" required>
                  <TextInput
                    name="name"
                    autoFocus={true}
                    type="string"
                    className={classes.input}
                    onChange={event =>
                      setServiceState({...service, name: event.target.value})
                    }
                  />
                </FormField>
              </Grid>
            </Grid>
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <FormField label="Service ID">
                  <TextInput
                    type="string"
                    className={classes.detailInput}
                    onChange={event =>
                      setServiceState({
                        ...service,
                        externalId: event.target.value,
                      })
                    }
                  />
                </FormField>
              </Grid>
              <Grid item xs={6}>
                <FormField label="Customer">
                  <CustomerTypeahead
                    className={classes.detailInput}
                    onCustomerSelection={customer =>
                      setServiceState({...service, customer: customer})
                    }
                    required={false}
                    margin="dense"
                  />
                </FormField>
              </Grid>
            </Grid>
            <div className={classes.separator} />
            {service.properties.length > 0 ? (
              <Grid container spacing={2}>
                {service.properties.map((property, index) => (
                  <Grid key={property.id} item xs={6}>
                    <PropertyValueInput
                      required={!!property.propertyType.isInstanceProperty}
                      disabled={!property.propertyType.isInstanceProperty}
                      label={property.propertyType.name}
                      className={classes.propInput}
                      margin="dense"
                      inputType="Property"
                      property={property}
                      onChange={propertyChangedHandler(index)}
                      headlineVariant="form"
                    />
                  </Grid>
                ))}
              </Grid>
            ) : null}
          </div>
        </div>
      </DialogContent>
      <DialogActions className={classes.dialogActions}>
        <Button onClick={onBackClicked} skin="regular">
          Back
        </Button>
        <Button disabled={isSaveDisabled()} onClick={saveService}>
          Create
        </Button>
      </DialogActions>
    </>
  );
};

export default AddServiceDetails;
