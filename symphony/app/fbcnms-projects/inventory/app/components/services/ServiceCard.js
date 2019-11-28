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

import Card from '@fbcnms/ui/components/design-system/Card/Card';
import CardHeader from '@fbcnms/ui/components/design-system/Card/CardHeader';
import ExpandingPanel from '@fbcnms/ui/components/ExpandingPanel';
import Grid from '@material-ui/core/Grid';
import InventoryQueryRenderer from '../InventoryQueryRenderer';
import React from 'react';
import ServiceHeader from './ServiceHeader';
import ServiceLinksView from './ServiceLinksView';
import ServiceNetworkMap from './ServiceNetworkMap';
import Text from '@fbcnms/ui/components/design-system/Text';
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
    height: '100%',
  },
  sidePanel: {
    display: 'flex',
    flexDirection: 'column',
    height: '100%',
    backgroundColor: symphony.palette.white,
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
    boxShadow: 'none',
    padding: '32px 32px 12px 32px',
  },
  field: {
    width: '50%',
    marginBottom: '12px',
    paddingRight: '16px',
  },
  expanded: {},
  panel: {
    '&$expanded': {
      margin: '0px 0px',
    },
    boxShadow: 'none',
  },
  linksPanel: {
    width: '500px',
    marginRight: '16px',
  },
  topologyPanel: {
    display: 'flex',
    flexDirection: 'column',
    height: '100%',
    padding: '32px',
  },
  separator: {
    borderBottom: `1px solid ${symphony.palette.separator}`,
    margin: 0,
  },
  detailValue: {
    color: symphony.palette.D500,
    display: 'block',
  },
  detail: {
    paddingBottom: '12px',
  },
  text: {
    display: 'block',
  },
  topologyCard: {
    flexGrow: 1,
  },
  titleText: {
    lineHeight: '28px',
  },
  expansionPanel: {
    '&&': {
      padding: '0px 20px 0px 32px',
    },
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
          <Grid container className={classes.root}>
            <Grid item xs={6} sm={8} lg={8} xl={9}>
              <div className={classes.topologyPanel}>
                <ServiceHeader
                  service={service}
                  onBackClicked={navigateToMainPage}
                  onServiceRemoved={navigateToMainPage}
                />
                <Card className={classes.topologyCard}>
                  <CardHeader className={classes.titleText}>
                    Topology
                  </CardHeader>
                  <ServiceNetworkMap serviceId={service.id} />
                </Card>
              </div>
            </Grid>
            <Grid
              item
              xs={6}
              sm={4}
              lg={4}
              xl={3}
              className={classes.sidePanel}>
              <Card className={classes.detailsCard}>
                <div className={classes.detail}>
                  <Text variant="h6" className={classes.text}>
                    {service.name}
                  </Text>
                  <Text
                    variant="subtitle2"
                    weight="regular"
                    className={classes.detailValue}>
                    {service.externalId}
                  </Text>
                </div>
                <div className={classes.detail}>
                  <Text variant="subtitle2" className={classes.text}>
                    Service Type
                  </Text>
                  <Text
                    variant="subtitle2"
                    weight="regular"
                    className={classes.detailValue}>
                    {service.serviceType.name}
                  </Text>
                </div>
                {service.customer && (
                  <div className={classes.detail}>
                    <Text variant="subtitle2" className={classes.text}>
                      Client
                    </Text>
                    <Text
                      variant="subtitle2"
                      weight="regular"
                      className={classes.detailValue}>
                      {service.customer.name}
                    </Text>
                  </div>
                )}
              </Card>
              <div className={classes.separator} />
              <ExpandingPanel
                title="Links"
                defaultExpanded={false}
                expandedClassName={classes.expanded}
                className={classes.panel}
                expansionPanelSummaryClassName={classes.expansionPanel}>
                <ServiceLinksView links={service.links} />
              </ExpandingPanel>
              <div className={classes.separator} />
            </Grid>
          </Grid>
        );
      }}
    />
  );
};

export default withRouter(withStyles(styles)(ServiceCard));
