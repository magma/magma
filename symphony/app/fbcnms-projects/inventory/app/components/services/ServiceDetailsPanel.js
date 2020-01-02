/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {EditServiceMutationResponse} from '../../mutations/__generated__/EditServiceMutation.graphql';
import type {MutationCallbacks} from '../../mutations/MutationCallbacks.js';
import type {Property} from '../../common/Property';
import type {Service} from '../../common/Service';

import ArrowBackIcon from '@material-ui/icons/ArrowBack';
import CustomerTypeahead from '../typeahead/CustomerTypeahead';
import EditServiceMutation from '../../mutations/EditServiceMutation';
import ExpandingPanel from '@fbcnms/ui/components/ExpandingPanel';
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import IconButton from '@material-ui/core/IconButton';
import NameInput from '@fbcnms/ui/components/design-system/Form/NameInput';
import PropertyValueInput from '../form/PropertyValueInput';
import React, {useRef, useState} from 'react';
import SideBar from '@fbcnms/ui/components/layout/SideBar';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import TextField from '@material-ui/core/TextField';
import symphony from '@fbcnms/ui/theme/symphony';
import update from 'immutability-helper';
import useStateWithCallback from 'use-state-with-callback';
import useVerticalScrollingEffect from '../../common/useVerticalScrollingEffect';

import {createFragmentContainer, graphql} from 'react-relay';
import {getInitialPropertyFromType} from '../../common/PropertyType';
import {
  getNonInstancePropertyTypes,
  sortPropertiesByIndex,
  toPropertyInput,
} from '../../common/Property';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';

type Props = {
  shown: boolean,
  service: Service,
  panelWidth?: number,
  onClose: () => void,
};

const useStyles = makeStyles({
  root: {
    height: '100%',
  },
  sideBar: {
    border: 'none',
    boxShadow: 'none',
    borderRadius: '0px',
    padding: '0px',
    overflowY: 'auto',
  },
  separator: {
    borderBottom: `1px solid ${symphony.palette.separator}`,
    marginTop: '8px',
  },
  expanded: {
    padding: '0px',
  },
  panel: {
    '&$expanded': {
      margin: '0px',
    },
    boxShadow: 'none',
    padding: '0px',
    background: 'transparent',
  },
  closeButton: {
    '&&': {
      backgroundColor: symphony.palette.D10,
      color: 'blue',
      margin: '32px 0px 0px 32px',
      padding: '2px',
      display: 'inline-block',
      '&:hover': {
        backgroundColor: symphony.palette.D100,
      },
    },
  },
  expansionPanel: {
    '&&': {
      padding: '24px 20px 16px 32px',
    },
  },
  topBar: {
    display: 'flex',
  },
  detailPane: {
    padding: '0px 32px',
  },
  input: {
    marginBottom: '20px',
  },
});

const ServiceDetailsPanel = (props: Props) => {
  const classes = useStyles();
  const {shown, service, panelWidth, onClose} = props;
  const thisElement = useRef(null);
  const [dirtyValue, setDirtyValue] = useStateWithCallback(null, dirtyValue => {
    if (
      dirtyValue == 'equipment' ||
      dirtyValue == 'location' ||
      dirtyValue == 'service' ||
      dirtyValue == 'customer'
    ) {
      // don't wait for blur because there is no blur in those selection values
      editService();
    }
  });
  useVerticalScrollingEffect(thisElement);
  const enqueueSnackbar = useEnqueueSnackbar();
  let properties = service?.properties ?? [];
  if (service.serviceType.propertyTypes) {
    properties = [
      ...properties,
      ...getNonInstancePropertyTypes(
        properties,
        service.serviceType.propertyTypes,
      ).map(propType => getInitialPropertyFromType(propType)),
    ];
    properties = properties.sort(sortPropertiesByIndex);
  }

  const [editableService, setEditableService] = useState({
    id: service.id,
    name: service.name,
    externalId: service.externalId,
    customer: service.customer,
    properties: properties,
  });

  const getServiceInput = () => {
    return {
      data: {
        id: editableService.id,
        name: editableService.name,
        externalId: editableService.externalId,
        customerId: editableService.customer?.id,
        properties: toPropertyInput(editableService.properties),
        upstreamServiceIds: [],
      },
    };
  };

  const enqueueError = (message: string) => {
    enqueueSnackbar(message, {
      children: key => (
        <SnackbarItem id={key} message={message} variant="error" />
      ),
    });
  };

  const editService = () => {
    if (dirtyValue !== null) {
      const callbacks: MutationCallbacks<EditServiceMutationResponse> = {
        onCompleted: (response, errors) => {
          if (errors && errors[0]) {
            enqueueError(errors[0].message);
          }
        },
        onError: () => {
          enqueueError('Error saving service');
        },
      };

      EditServiceMutation(getServiceInput(), callbacks);
      setDirtyValue(null);
    }
  };

  const onChangeProperty = index => (property: Property) => {
    setEditableService(
      update(editableService, {
        properties: {
          [index]: {$set: property},
        },
      }),
    );
    setDirtyValue(property.propertyType.type);
  };
  const onChangeDetail = (key: 'name' | 'externalId' | 'customer', value) => {
    // $FlowFixMe Update specific value
    setEditableService(update(editableService, {[key]: {$set: value}}));
    setDirtyValue(key);
  };

  const backButton = (props: {onClose: () => void}) => (
    <div className={classes.topBar}>
      <IconButton className={classes.closeButton} onClick={props.onClose}>
        <ArrowBackIcon fontSize="small" color="primary" />
      </IconButton>
    </div>
  );
  return (
    <SideBar
      isShown={shown}
      top={0}
      width={panelWidth}
      onClose={onClose}
      className={classes.sideBar}
      backButton={backButton}>
      <div ref={thisElement} className={classes.scroller}>
        <ExpandingPanel
          title="Details"
          defaultExpanded={true}
          expandedClassName={classes.expanded}
          expansionPanelSummaryClassName={classes.expansionPanel}
          detailsPaneClass={classes.detailPane}
          className={classes.panel}>
          <div className={classes.input}>
            <NameInput
              value={editableService.name}
              onChange={event => onChangeDetail('name', event.target.value)}
              onBlur={editService}
              hasSpacer={false}
            />
          </div>
          <div className={classes.input}>
            <FormField label="Service ID">
              <TextField
                name="serviceId"
                variant="outlined"
                margin="dense"
                onChange={event =>
                  onChangeDetail('externalId', event.target.value)
                }
                value={editableService.externalId ?? ''}
                onBlur={editService}
              />
            </FormField>
          </div>
          <div className={classes.input}>
            <FormField label="Service Type">
              <TextField
                disabled
                name="type"
                variant="outlined"
                margin="dense"
                value={service.serviceType.name}
              />
            </FormField>
          </div>
          <div className={classes.input}>
            <FormField label="Customer">
              <CustomerTypeahead
                onCustomerSelection={customer => {
                  onChangeDetail('customer', customer);
                }}
                required={false}
                selectedCustomer={editableService.customer?.name}
                margin="dense"
              />
            </FormField>
          </div>
        </ExpandingPanel>
        <div className={classes.separator} />
        <ExpandingPanel
          title="Properties"
          defaultExpanded={true}
          expandedClassName={classes.expanded}
          expansionPanelSummaryClassName={classes.expansionPanel}
          detailsPaneClass={classes.detailPane}
          className={classes.panel}>
          {editableService.properties.map((property, index) => (
            <PropertyValueInput
              fullWidth
              required={!!property.propertyType.isMandatory}
              disabled={!property.propertyType.isInstanceProperty}
              label={property.propertyType.name}
              className={classes.input}
              margin="dense"
              inputType="Property"
              property={property}
              // $FlowFixMe pass property and not property type
              onChange={onChangeProperty(index)}
              onBlur={editService}
              headlineVariant="form"
            />
          ))}
        </ExpandingPanel>
      </div>
    </SideBar>
  );
};

export default createFragmentContainer(ServiceDetailsPanel, {
  service: graphql`
    fragment ServiceDetailsPanel_service on Service {
      id
      name
      externalId
      customer {
        name
      }
      serviceType {
        id
        name
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
      properties {
        id
        propertyType {
          id
          name
          type
          isEditable
          isInstanceProperty
          isMandatory
          stringValue
        }
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
  `,
});
