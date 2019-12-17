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
  AddServiceLinkMutationResponse,
  AddServiceLinkMutationVariables,
} from '../../mutations/__generated__/AddServiceLinkMutation.graphql';
import type {Link} from '../../common/Equipment';
import type {MutationCallbacks} from '../../mutations/MutationCallbacks.js';
import type {
  RemoveServiceLinkMutationResponse,
  RemoveServiceLinkMutationVariables,
} from '../../mutations/__generated__/RemoveServiceLinkMutation.graphql';
import type {Service} from '../../common/Service';

import AddCircleOutlineIcon from '@material-ui/icons/AddCircleOutline';
import AddServiceLinkMutation from '../../mutations/AddServiceLinkMutation';
import Button from '@fbcnms/ui/components/design-system/Button';
import Card from '@fbcnms/ui/components/design-system/Card/Card';
import ExpandingPanel from '@fbcnms/ui/components/ExpandingPanel';
import IconButton from '@material-ui/core/IconButton';
import React, {useState} from 'react';
import RemoveServiceLinkMutation from '../../mutations/RemoveServiceLinkMutation';
import ServiceLinksSubservicesMenu from './ServiceLinksSubservicesMenu';
import ServiceLinksView from './ServiceLinksView';
import Text from '@fbcnms/ui/components/design-system/Text';
import symphony from '@fbcnms/ui/theme/symphony';
import {makeStyles} from '@material-ui/styles';

type Props = {
  service: Service,
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
  addLink: {
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
});

/* $FlowFixMe - Flow doesn't support typing when using forwardRef on a
 * funcional component
 */
const ServicePanel = React.forwardRef((props: Props, ref) => {
  const classes = useStyles();
  const {service, onOpenDetailsPanel} = props;
  const [anchorEl, setAnchorEl] = useState<?HTMLElement>(null);
  const [showAddMenu, setShowAddMenu] = useState(false);
  const [linksExpanded, setLinksExpanded] = useState(false);

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
      <div className={classes.separator} />
      <ExpandingPanel
        title="Links & Subservices"
        defaultExpanded={false}
        expandedClassName={classes.expanded}
        className={classes.panel}
        expansionPanelSummaryClassName={classes.expansionPanel}
        detailsPaneClass={classes.detailsPanel}
        expanded={linksExpanded}
        onChange={expanded => setLinksExpanded(expanded)}
        rightContent={
          <IconButton
            className={classes.addLink}
            onClick={event => {
              setAnchorEl(event.currentTarget);
              setShowAddMenu(true);
            }}>
            <AddCircleOutlineIcon />
          </IconButton>
        }>
        <ServiceLinksView links={service.links} onDeleteLink={onDeleteLink} />
      </ExpandingPanel>
      <div className={classes.separator} />
      {showAddMenu ? (
        <ServiceLinksSubservicesMenu
          key={`${service.id}-menu`}
          service={service}
          anchorEl={anchorEl}
          onClose={() => setAnchorEl(null)}
          onAddLink={onAddLink}
        />
      ) : null}
    </div>
  );
});

export default ServicePanel;
