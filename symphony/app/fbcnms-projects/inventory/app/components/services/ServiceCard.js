/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ContextRouter} from 'react-router-dom';
import type {WithStyles} from '@material-ui/core';

import ExpandingPanel from '@fbcnms/ui/components/ExpandingPanel';
import FormField from '@fbcnms/ui/components/FormField';
import InventoryQueryRenderer from '../InventoryQueryRenderer';
import React from 'react';
import ServiceDetails from './ServiceDetails';
import ServiceHeader from './ServiceHeader';
import ServiceLinksView from './ServiceLinksView';
import ServiceNetworkMap from './ServiceNetworkMap';
import symphony from '@fbcnms/ui/theme/symphony';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {graphql} from 'react-relay';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

type Props = {
  serviceId: ?string,
} & WithStyles<typeof styles> &
  ContextRouter;

const styles = _ => ({
  root: {
    height: 'calc(100% - 80px)',
    display: 'flex',
    flexDirection: 'column',
    margin: '40px 32px',
  },
  contentRoot: {
    position: 'relative',
    flexGrow: 1,
    overflow: 'auto',
    backgroundColor: symphony.palette.white,
  },
  tabsContainer: {
    marginBottom: '16px',
    display: 'flex',
    flexDirection: 'column',
    flex: 1,
  },
  tabs: {
    backgroundColor: symphony.palette.white,
  },
  tabContainer: {
    width: 'auto',
  },
  detailsCard: {
    display: 'flex',
    flexWrap: 'wrap',
  },
  field: {
    width: '50%',
    marginBottom: '12px',
    paddingRight: '16px',
  },
  panel: {
    marginBottom: '16px',
    marginTop: '16px',
  },
  linksPanel: {
    width: '500px',
    marginRight: '16px',
  },
  topologyPanel: {
    flexGrow: 1,
  },
  topologyCard: {
    display: 'flex',
  },
});

const serviceQuery = graphql`
  query ServiceCardQuery($serviceId: ID!) {
    service(id: $serviceId) {
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
          ...PropertyTypeFormField_propertyType
          ...DynamicPropertiesGrid_propertyTypes
        }
      }
      properties {
        ...PropertyFormField_property
        ...DynamicPropertiesGrid_properties
      }
      links {
        ...ServiceLinksView_links
      }
    }
  }
`;

const ServiceCard = (props: Props) => {
  const {classes, serviceId, history, match} = props;

  const navigateToMainPage = () => {
    ServerLogger.info(LogEvents.SERVICES_SEARCH_NAV_CLICKED, {
      source: 'service_card',
    });
    history.push(match.url);
  };

  return (
    <InventoryQueryRenderer
      query={serviceQuery}
      variables={{
        serviceId,
      }}
      render={props => {
        const {service} = props;
        return (
          <div className={classes.root}>
            <ServiceHeader
              service={service}
              onBackClicked={navigateToMainPage}
              onServiceRemoved={navigateToMainPage}
            />
            <ExpandingPanel title="Details" className={classes.panel}>
              <div className={classes.detailsCard}>
                <div className={classes.field}>
                  <FormField label="Service ID" value={service.externalId} />
                </div>
                <div className={classes.field}>
                  <FormField label="Customer" value={service.customer?.name} />
                </div>
              </div>
            </ExpandingPanel>
            <ExpandingPanel title="Properties" className={classes.panel}>
              <ServiceDetails service={service} />
            </ExpandingPanel>
            <div className={classes.topologyCard}>
              <div className={classes.linksPanel}>
                <ExpandingPanel title="Links" className={classes.panel}>
                  <ServiceLinksView links={service.links} />
                </ExpandingPanel>
              </div>
              <div className={classes.topologyPanel}>
                <ExpandingPanel title="Topology" className={classes.panel}>
                  <ServiceNetworkMap serviceId={service.id} />
                </ExpandingPanel>
              </div>
            </div>
          </div>
        );
      }}
    />
  );
};

export default withRouter(withStyles(styles)(ServiceCard));
