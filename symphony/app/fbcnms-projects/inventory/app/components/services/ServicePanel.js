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
  AddServiceEndpointMutationResponse,
  AddServiceEndpointMutationVariables,
  ServiceEndpointRole,
} from '../../mutations/__generated__/AddServiceEndpointMutation.graphql';
import type {
  AddServiceLinkMutationResponse,
  AddServiceLinkMutationVariables,
} from '../../mutations/__generated__/AddServiceLinkMutation.graphql';
import type {EquipmentPort, Link} from '../../common/Equipment';
import type {MutationCallbacks} from '../../mutations/MutationCallbacks.js';
import type {
  RemoveServiceEndpointMutationResponse,
  RemoveServiceEndpointMutationVariables,
} from '../../mutations/__generated__/RemoveServiceEndpointMutation.graphql';
import type {
  RemoveServiceLinkMutationResponse,
  RemoveServiceLinkMutationVariables,
} from '../../mutations/__generated__/RemoveServiceLinkMutation.graphql';
import type {ServiceEndpoint, ServiceStatus} from '../../common/Service';
import type {ServicePanel_service} from './__generated__/ServicePanel_service.graphql';

import AddServiceEndpointMutation from '../../mutations/AddServiceEndpointMutation';
import AddServiceLinkMutation from '../../mutations/AddServiceLinkMutation';
import Button from '@fbcnms/ui/components/design-system/Button';
import Card from '@fbcnms/ui/components/design-system/Card/Card';
import EditServiceMutation from '../../mutations/EditServiceMutation';
import ExpandingPanel from '@fbcnms/ui/components/ExpandingPanel';
import React, {useState} from 'react';
import RemoveServiceEndpointMutation from '../../mutations/RemoveServiceEndpointMutation';
import RemoveServiceLinkMutation from '../../mutations/RemoveServiceLinkMutation';
import Select from '@fbcnms/ui/components/design-system/Select/Select';
import ServiceEndpointsMenu from './ServiceEndpointsMenu';
import ServiceEndpointsView from './ServiceEndpointsView';
import ServiceLinksSubservicesMenu from './ServiceLinksSubservicesMenu';
import ServiceLinksView from './ServiceLinksView';
import Text from '@fbcnms/ui/components/design-system/Text';
import symphony from '@fbcnms/ui/theme/symphony';
import {createFragmentContainer, graphql} from 'react-relay';
import {makeStyles} from '@material-ui/styles';
import {
  serviceStatusToColor,
  serviceStatusToVisibleNames,
} from '../../common/Service';

type Props = {
  service: ServicePanel_service,
  onOpenDetailsPanel: () => void,
};

const useStyles = makeStyles({
  root: {
    overflowY: 'auto',
    height: '100%',
  },
  contentRoot: {
    position: 'relative',
    flexGrow: 1,
    overflow: 'auto',
    backgroundColor: symphony.palette.white,
  },
  detailsCard: {
    boxShadow: 'none',
    padding: '32px 32px 12px 32px',
    position: 'relative',
  },
  expanded: {},
  panel: {
    '&$expanded': {
      margin: '0px 0px',
      padding: '24px 0px 18px 0px',
    },
    boxShadow: 'none',
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
  expansionPanel: {
    '&&': {
      padding: '0px 20px 0px 32px',
    },
  },
  addButton: {
    marginRight: '8px',
    '&:hover': {
      backgroundColor: 'transparent',
    },
  },
  dialog: {
    width: '80%',
    maxWidth: '1280px',
    height: '90%',
    maxHeight: '800px',
  },
  edit: {
    position: 'absolute',
    bottom: '24px',
    right: '24px',
  },
  editText: {
    color: symphony.palette.B500,
  },
  select: {
    marginBottom: '24px',
  },
  detailsPanel: {
    padding: '0px',
  },
});

/* $FlowFixMe - Flow doesn't support typing when using forwardRef on a
 * funcional component
 */
const ServicePanel = React.forwardRef((props: Props, ref) => {
  const classes = useStyles();
  const {service, onOpenDetailsPanel} = props;
  const [endpointsExpanded, setEndpointsExpanded] = useState(false);
  const [linksExpanded, setLinksExpanded] = useState(false);

  const onAddEndpoint = (port: EquipmentPort, role: ServiceEndpointRole) => {
    const variables: AddServiceEndpointMutationVariables = {
      input: {
        id: service.id,
        portId: port.id,
        role: role,
      },
    };
    const callbacks: MutationCallbacks<AddServiceEndpointMutationResponse> = {
      onCompleted: () => {
        setEndpointsExpanded(true);
      },
    };
    AddServiceEndpointMutation(variables, callbacks);
  };

  const onAddLink = (link: Link) => {
    const variables: AddServiceLinkMutationVariables = {
      id: service.id,
      linkId: link.id,
    };
    const callbacks: MutationCallbacks<AddServiceLinkMutationResponse> = {
      onCompleted: () => {
        setLinksExpanded(true);
      },
    };
    AddServiceLinkMutation(variables, callbacks);
  };

  const onDeleteLink = (link: Link) => {
    const variables: RemoveServiceLinkMutationVariables = {
      id: service.id,
      linkId: link.id,
    };
    const callbacks: MutationCallbacks<RemoveServiceLinkMutationResponse> = {
      onCompleted: () => {
        setLinksExpanded(true);
      },
    };
    RemoveServiceLinkMutation(variables, callbacks);
  };

  const onDeleteEndpoint = (endpoint: ServiceEndpoint) => {
    const variables: RemoveServiceEndpointMutationVariables = {
      serviceEndpointId: endpoint.id,
    };
    const callbacks: MutationCallbacks<RemoveServiceEndpointMutationResponse> = {
      onCompleted: () => {
        setEndpointsExpanded(true);
      },
    };
    RemoveServiceEndpointMutation(variables, callbacks);
  };

  const onStatusChange = (status: ServiceStatus) => {
    EditServiceMutation({
      data: {
        id: service.id,
        status: status,
      },
    });
  };

  const getValidServiceStatus = (type: string): ServiceStatus => {
    if (
      type === 'DISCONNECTED' ||
      type === 'IN_SERVICE' ||
      type === 'MAINTENANCE' ||
      type === 'PENDING'
    ) {
      return type;
    }

    return 'PENDING';
  };

  return (
    <div className={classes.root} ref={ref}>
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
        <Select
          className={classes.select}
          label="Status"
          options={Object.entries(serviceStatusToVisibleNames).map(entry => {
            // $FlowFixMe - Flow doesn't value type well from object
            return {value: entry[0], label: entry[1]};
          })}
          selectedValue={service.status}
          onChange={value => onStatusChange(getValidServiceStatus(value))}
          skin={serviceStatusToColor[getValidServiceStatus(service.status)]}
        />
        <div className={classes.detail}>
          <Text variant="subtitle2" className={classes.text}>
            Service Type
          </Text>
          <Text variant="body2" className={classes.detailValue}>
            {service.serviceType.name}
          </Text>
        </div>
        {service.customer && (
          <div className={classes.detail}>
            <Text variant="subtitle2" className={classes.text}>
              Client
            </Text>
            <Text variant="body2" className={classes.detailValue}>
              {service.customer.name}
            </Text>
          </div>
        )}
        <div className={classes.edit}>
          <Button variant="text" onClick={onOpenDetailsPanel}>
            <Text variant="body2" className={classes.editText}>
              View & Edit Details
            </Text>
          </Button>
        </div>
      </Card>
      <>
        <div className={classes.separator} />
        <ExpandingPanel
          title="Endpoints"
          defaultExpanded={false}
          expandedClassName={classes.expanded}
          className={classes.panel}
          expansionPanelSummaryClassName={classes.expansionPanel}
          detailsPaneClass={classes.detailsPanel}
          expanded={endpointsExpanded}
          onChange={expanded => setEndpointsExpanded(expanded)}
          rightContent={
            <ServiceEndpointsMenu
              service={{id: service.id, name: service.name}}
              onAddEndpoint={onAddEndpoint}
            />
          }>
          <ServiceEndpointsView
            endpoints={service.endpoints}
            onDeleteEndpoint={onDeleteEndpoint}
          />
        </ExpandingPanel>
      </>
      <div className={classes.separator} />
      <ExpandingPanel
        title="Links"
        defaultExpanded={false}
        expandedClassName={classes.expanded}
        className={classes.panel}
        expansionPanelSummaryClassName={classes.expansionPanel}
        detailsPaneClass={classes.detailsPanel}
        expanded={linksExpanded}
        onChange={expanded => setLinksExpanded(expanded)}
        rightContent={
          <ServiceLinksSubservicesMenu
            service={{id: service.id, name: service.name}}
            onAddLink={onAddLink}
          />
        }>
        <ServiceLinksView links={service.links} onDeleteLink={onDeleteLink} />
      </ExpandingPanel>
      <div className={classes.separator} />
    </div>
  );
});

export default createFragmentContainer(ServicePanel, {
  service: graphql`
    fragment ServicePanel_service on Service {
      id
      name
      externalId
      status
      customer {
        name
      }
      serviceType {
        name
      }
      links {
        id
        ...ServiceLinksView_links
      }
      endpoints {
        ...ServiceEndpointsView_endpoints
      }
    }
  `,
});
