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
  AddServiceMutationResponse,
  AddServiceMutationVariables,
} from '../../mutations/__generated__/AddServiceMutation.graphql';
import type {ContextRouter} from 'react-router-dom';
import type {Link} from '../../common/Equipment';
import type {MutationCallbacks} from '../../mutations/MutationCallbacks.js';
import type {Service} from '../../common/Service';
import type {WithSnackbarProps} from 'notistack';
import type {WithStyles} from '@material-ui/core';

import AddLinkToServiceDialog from './AddLinkToServiceDialog';
import AddServiceMutation from '../../mutations/AddServiceMutation';
import Breadcrumbs from '@fbcnms/ui/components/Breadcrumbs';
import Button from '@fbcnms/ui/components/design-system/Button';
import CircularProgress from '@material-ui/core/CircularProgress';
import CustomerTypeahead from '../typeahead/CustomerTypeahead';
import Dialog from '@material-ui/core/Dialog';
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import Grid from '@material-ui/core/Grid';
import PropertiesAddEditSection from '../form/PropertiesAddEditSection';
import React, {useCallback, useEffect, useState} from 'react';
import RelayEnvironment from '../../common/RelayEnvironment.js';
import SectionedCard from '@fbcnms/ui/components/SectionedCard';
import ServiceLinksTable from './ServiceLinksTable';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextField from '@material-ui/core/TextField';
import nullthrows from '@fbcnms/util/nullthrows';
import update from 'immutability-helper';
import useRouter from '@fbcnms/ui/hooks/useRouter';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {fetchQuery, graphql} from 'relay-runtime';
import {getInitialPropertyFromType} from '../../common/PropertyType';
import {sortPropertiesByIndex, toPropertyInput} from '../../common/Property';
import {withRouter} from 'react-router-dom';
import {withSnackbar} from 'notistack';
import {withStyles} from '@material-ui/core/styles';

type Props = {
  serviceTypeId: ?string,
} & WithStyles<typeof styles> &
  ContextRouter &
  WithSnackbarProps;

const styles = theme => ({
  header: {
    marginBottom: '21px',
    paddingBottom: '0px',
  },
  root: {
    height: '100%',
    display: 'flex',
    flexDirection: 'column',
    padding: '20px 16px',
  },
  contentRoot: {
    position: 'relative',
    flexGrow: 1,
    overflow: 'auto',
  },
  cards: {
    height: 'calc(100% - 60px)',
    padding: '8px',
    overflowY: 'auto',
  },
  card: {
    display: 'flex',
    flexDirection: 'column',
  },
  input: {
    width: '250px',
    paddingBottom: '24px',
    margin: '5px',
    marginLeft: '0px',
  },
  detailInput: {
    display: 'inline-flex',
    margin: '5px 20px 5px 0px',
    width: '250px',
  },
  nameHeader: {
    display: 'flex',
    alignItems: 'center',
    marginBottom: '24px',
    marginRight: '8px',
  },
  breadcrumbs: {
    flexGrow: 1,
  },
  footer: {
    padding: '12px 16px',
    marginRight: '24px',
    marginLeft: '24px',
    boxShadow: '0px -1px 4px rgba(0, 0, 0, 0.1)',
  },
  separator: {
    borderBottom: `1px solid ${theme.palette.grey[100]}`,
    margin: '0 0 24px -24px',
    paddingBottom: '24px',
    width: 'calc(100% + 48px)',
  },
  separator: {
    borderBottom: `1px solid ${theme.palette.grey[100]}`,
    margin: '0 0 24px -24px',
    paddingBottom: '24px',
    width: 'calc(100% + 48px)',
  },
  dense: {
    paddingTop: '9px',
    paddingBottom: '9px',
    height: '14px',
  },
  headerText: {
    fontSize: '20px',
    lineHeight: '24px',
    fontWeight: 500,
  },
  section: {
    marginTop: '0px',
    overflow: 'visible',
  },
  details: {
    display: 'flex',
  },
  cancelButton: {
    marginRight: '8px',
  },
  footer: {
    padding: '12px 16px',
    boxShadow: '0px -1px 4px rgba(0, 0, 0, 0.1)',
  },
});

const addServiceCard__serviceTypeQuery = graphql`
  query AddServiceCard__serviceTypeQuery($serviceTypeId: ID!) {
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

const AddServiceCard = (props: Props) => {
  const [service, setService] = useState<?Service>(null);
  const [showAddLinkDialog, setShowAddLinkDialog] = useState(false);
  const {history, match} = useRouter();

  const getService = useCallback(async (): Promise<Service> => {
    const response = await fetchQuery(
      RelayEnvironment,
      addServiceCard__serviceTypeQuery,
      {
        serviceTypeId: props.serviceTypeId,
      },
    );
    const serviceType = response.serviceType;
    const initialProps = (serviceType.propertyTypes || [])
      .map(propType => getInitialPropertyFromType(propType))
      .sort(sortPropertiesByIndex);
    return {
      id: 'service@tmp',
      name: '',
      externalId: null,
      customer: null,
      serviceType: serviceType,
      upstream: [],
      downstream: [],
      properties: initialProps,
      terminationPoints: [],
      links: [],
    };
  }, [props.serviceTypeId]);

  useEffect(() => {
    getService().then(service => {
      setService(service);
    });
  }, [getService]);

  const isSaveDisabled = () => {
    return !service?.name;
  };

  const saveService = () => {
    const {name, externalId, customer, properties, links} = nullthrows(service);
    const serviceTypeId = nullthrows(service?.serviceType.id);
    const variables: AddServiceMutationVariables = {
      data: {
        name,
        externalId,
        serviceTypeId,
        customerId: customer?.id,
        upstreamServiceIds: [],
        properties: toPropertyInput(properties),
        terminationPointIds: [],
        linkIds: links.map(l => l.id),
      },
    };

    const callbacks: MutationCallbacks<AddServiceMutationResponse> = {
      onCompleted: (response, errors) => {
        if (errors && errors[0]) {
          enqueueError(errors[0].message);
        } else {
          // navigate to main page
          history.push(match.url);
        }
      },
      onError: () => {
        enqueueError('Error saving work order');
      },
    };
    ServerLogger.info(LogEvents.SAVE_SERVICE_BUTTON_CLICKED, {
      source: 'service_details',
    });
    AddServiceMutation(variables, callbacks);
  };

  const enqueueError = (message: string) => {
    props.enqueueSnackbar(message, {
      children: key => (
        <SnackbarItem id={key} message={message} variant="error" />
      ),
    });
  };

  const setServiceDetail = (key: 'name' | 'externalId' | 'customer', value) => {
    // $FlowFixMe Set state for each field
    setService(update(service, {[key]: {$set: value}}));
  };

  const propertyChangedHandler = index => property => {
    setService(
      update(service, {
        properties: {[index]: {$set: property}},
      }),
    );
  };

  const navigateToMainPage = () => {
    ServerLogger.info(LogEvents.SERVICES_SEARCH_NAV_CLICKED, {
      source: 'service_details',
    });
    const {match} = props;
    props.history.push(match.url);
  };

  const onAddLink = (link: Link) => {
    if (service != null) {
      setService({...service, links: [...(service?.links || []), link]});
    }
  };

  const onDeleteLink = (link: Link) => {
    setService(
      update(service, {
        links: {$set: (service?.links || []).filter(l => l != link)},
      }),
    );
  };

  const {classes} = props;
  if (!service) {
    return (
      <div className={classes.root}>
        <CircularProgress />
      </div>
    );
  }
  return (
    <div className={classes.root}>
      <div className={classes.nameHeader}>
        <div className={classes.breadcrumbs}>
          <Breadcrumbs
            breadcrumbs={[
              {
                id: 'services',
                name: 'Services',
                onClick: () => navigateToMainPage(),
              },
              {
                id: `new_service_` + Date.now(),
                name: 'New Service',
              },
            ]}
            size="large"
          />
        </div>
        <Button
          className={classes.cancelButton}
          skin="regular"
          onClick={navigateToMainPage}>
          Cancel
        </Button>
        <Button disabled={isSaveDisabled()} onClick={() => saveService()}>
          Save
        </Button>
      </div>
      <div className={classes.contentRoot}>
        <div className={classes.cards}>
          <SectionedCard className={classes.section}>
            <div className={classes.header}>
              <Text className={classes.headerText}>Details</Text>
            </div>
            <div>
              <Grid container spacing={2}>
                <Grid item xs={6}>
                  <FormField label="Name" required>
                    <TextField
                      name="name"
                      variant="outlined"
                      margin="dense"
                      className={classes.input}
                      onChange={event =>
                        setServiceDetail('name', event.target.value)
                      }
                    />
                  </FormField>
                </Grid>
              </Grid>
              <div className={classes.details}>
                <FormField label="Service ID">
                  <TextField
                    className={classes.detailInput}
                    variant="outlined"
                    margin="dense"
                    onChange={event => {
                      setServiceDetail('externalId', event.target.value);
                    }}
                  />
                </FormField>
                <FormField label="Customer">
                  <CustomerTypeahead
                    className={classes.detailInput}
                    onCustomerSelection={customer =>
                      setServiceDetail('customer', customer)
                    }
                    required={false}
                    margin="dense"
                  />
                </FormField>
              </div>
            </div>
          </SectionedCard>
          {service.properties.length > 0 ? (
            <SectionedCard>
              <PropertiesAddEditSection
                properties={service.properties}
                onChange={index => propertyChangedHandler(index)}
              />
            </SectionedCard>
          ) : null}
          <SectionedCard>
            <ServiceLinksTable
              links={service.links}
              onDeleteLink={onDeleteLink}
            />
          </SectionedCard>
        </div>
      </div>
      <div className={classes.footer}>
        <Button onClick={() => setShowAddLinkDialog(true)}>Add Link</Button>
      </div>
      {showAddLinkDialog ? (
        <Dialog
          open={true}
          onClose={() => setShowAddLinkDialog(false)}
          maxWidth={false}
          fullWidth={true}>
          <AddLinkToServiceDialog
            service={service}
            onClose={() => setShowAddLinkDialog(false)}
            onAddLink={link => {
              onAddLink(link);
              setShowAddLinkDialog(false);
            }}
          />
        </Dialog>
      ) : null}
    </div>
  );
};

export default withSnackbar(withRouter(withStyles(styles)(AddServiceCard)));
