/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type ServiceCard_service from './__generated__/ServiceCard_service.graphql';
import type {ContextRouter} from 'react-router-dom';
import type {WithStyles} from '@material-ui/core';

import Card from '@fbcnms/ui/components/design-system/Card/Card';
import CardHeader from '@fbcnms/ui/components/design-system/Card/CardHeader';

import Grid from '@material-ui/core/Grid';

import React, {useRef, useState} from 'react';
import ServiceDetailsPanel from './ServiceDetailsPanel';
import ServiceEquipmentTopology from './ServiceEquipmentTopology';
import ServiceHeader from './ServiceHeader';
import ServicePanel from './ServicePanel';
import symphony from '@fbcnms/ui/theme/symphony';
import {FormValidationContextProvider} from '@fbcnms/ui/components/design-system/Form/FormValidationContext';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {createFragmentContainer, graphql} from 'react-relay';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

type Props = {
  // $FlowFixMe (T62907961) Relay flow types
  service: ServiceCard_service,
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
  topologyPanel: {
    display: 'flex',
    flexDirection: 'column',
    height: '100%',
    padding: '32px',
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
  dialog: {
    width: '80%',
    maxWidth: '1280px',
    height: '90%',
    maxHeight: '800px',
  },
  detailsPanel: {
    padding: '0px',
  },
});

const ServiceCard = (props: Props) => {
  const {classes, service, history, match} = props;
  const [detailsPanelShown, setDetailsPanelShown] = useState(false);
  const panelRef = useRef();

  let panelWidth = undefined;
  const container = panelRef?.current;
  if (container) {
    panelWidth = container.clientWidth;
  }

  const navigateToMainPage = () => {
    ServerLogger.info(LogEvents.SERVICES_SEARCH_NAV_CLICKED, {
      source: 'service_card',
    });
    history.push(match.url);
  };
  return (
    <FormValidationContextProvider>
      <Grid container className={classes.root}>
        <Grid item xs={6} sm={8} lg={8} xl={9}>
          <div className={classes.topologyPanel}>
            <ServiceHeader
              service={service}
              onBackClicked={navigateToMainPage}
              onServiceRemoved={navigateToMainPage}
            />
            <Card className={classes.topologyCard}>
              <CardHeader className={classes.titleText}>Topology</CardHeader>
              <ServiceEquipmentTopology
                topology={service.topology}
                endpoints={service.endpoints}
              />
            </Card>
          </div>
        </Grid>
        <Grid item xs={6} sm={4} lg={4} xl={3} className={classes.sidePanel}>
          <ServicePanel
            service={service}
            onOpenDetailsPanel={() => setDetailsPanelShown(true)}
            ref={panelRef}
          />
          <ServiceDetailsPanel
            shown={detailsPanelShown}
            service={service}
            panelWidth={panelWidth}
            onClose={() => setDetailsPanelShown(false)}
          />
        </Grid>
      </Grid>
    </FormValidationContextProvider>
  );
};

export default withRouter(
  withStyles(styles)(
    createFragmentContainer(ServiceCard, {
      service: graphql`
        fragment ServiceCard_service on Service {
          id
          name
          ...ServiceDetailsPanel_service
          ...ServicePanel_service
          topology {
            ...ServiceEquipmentTopology_topology
          }
          endpoints {
            ...ServiceEquipmentTopology_endpoints
          }
        }
      `,
    }),
  ),
);
