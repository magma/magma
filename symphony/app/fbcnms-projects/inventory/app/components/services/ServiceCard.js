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
import type {
  EditServiceMutationResponse,
  EditServiceMutationVariables,
} from '../../mutations/__generated__/EditServiceMutation.graphql';
import type {Link} from '../../common/Equipment';
import type {MutationCallbacks} from '../../mutations/MutationCallbacks.js';
import type {WithStyles} from '@material-ui/core';

import AddCircleOutlineIcon from '@material-ui/icons/AddCircleOutline';
import AddLinkToServiceDialog from './AddLinkToServiceDialog';
import Card from '@fbcnms/ui/components/design-system/Card/Card';
import CardHeader from '@fbcnms/ui/components/design-system/Card/CardHeader';
import Dialog from '@material-ui/core/Dialog';
import EditServiceMutation from '../../mutations/EditServiceMutation';
import ExpandingPanel from '@fbcnms/ui/components/ExpandingPanel';
import Grid from '@material-ui/core/Grid';
import IconButton from '@material-ui/core/IconButton';
import React, {useState} from 'react';
import ServiceEquipmentTopology from './ServiceEquipmentTopology';
import ServiceHeader from './ServiceHeader';
import ServiceLinksView from './ServiceLinksView';
import Text from '@fbcnms/ui/components/design-system/Text';
import symphony from '@fbcnms/ui/theme/symphony';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {createFragmentContainer, graphql} from 'react-relay';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

type Props = {
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
  addLink: {
    marginRight: '8px',
  },
});

const ServiceCard = (props: Props) => {
  const {classes, service, history, match} = props;
  const [showAddLinkDialog, setShowAddLinkDialog] = useState(false);
  const [linksExpanded, setLinksExpanded] = useState(false);

  const onAddLink = (link: Link) => {
    const variables: EditServiceMutationVariables = {
      data: {
        id: service.id,
        name: service.name,
        externalId: service.externalId,
        customerId: service.customer?.id,
        upstreamServiceIds: [],
        properties: [],
        terminationPointIds: [],
        linkIds: [...service.links.map(l => l.id), link.id],
      },
    };
    const callbacks: MutationCallbacks<EditServiceMutationResponse> = {
      onCompleted: () => {
        setLinksExpanded(true);
      },
    };
    EditServiceMutation(variables, callbacks);
  };

  const navigateToMainPage = () => {
    ServerLogger.info(LogEvents.SERVICES_SEARCH_NAV_CLICKED, {
      source: 'service_card',
    });
    history.push(match.url);
  };
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
            <CardHeader className={classes.titleText}>Topology</CardHeader>
            <ServiceEquipmentTopology
              topology={service.topology}
              terminationPoints={service.terminationPoints}
            />
          </Card>
        </div>
      </Grid>
      <Grid item xs={6} sm={4} lg={4} xl={3} className={classes.sidePanel}>
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
          expansionPanelSummaryClassName={classes.expansionPanel}
          expanded={linksExpanded}
          onChange={expanded => setLinksExpanded(expanded)}
          rightContent={
            <IconButton
              className={classes.addLink}
              onClick={() => setShowAddLinkDialog(true)}>
              <AddCircleOutlineIcon />
            </IconButton>
          }>
          <ServiceLinksView links={service.links} />
        </ExpandingPanel>
        <div className={classes.separator} />
      </Grid>
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
    </Grid>
  );
};

export default withRouter(
  withStyles(styles)(
    createFragmentContainer(ServiceCard, {
      service: graphql`
        fragment ServiceCard_service on Service {
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
            id
            ...ServiceLinksView_links
          }
          terminationPoints {
            ...ServiceEquipmentTopology_terminationPoints
          }
          topology {
            ...ServiceEquipmentTopology_topology
          }
        }
      `,
    }),
  ),
);
