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
import type {ServiceTypesListQuery_serviceType} from './__generated__/ServiceTypesListQuery_serviceType.graphql';
import type {WithStyles} from '@material-ui/core';

import Avatar from '@material-ui/core/Avatar';
import Button from '@fbcnms/ui/components/design-system/Button';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import LinearScaleIcon from '@material-ui/icons/LinearScale';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemAvatar from '@material-ui/core/ListItemAvatar';
import ListItemText from '@material-ui/core/ListItemText';
import React from 'react';
import RelayEnvironment from '../../common/RelayEnvironment.js';
import Text from '@fbcnms/ui/components/design-system/Text';
import nullthrows from '@fbcnms/util/nullthrows';
import symphony from '@fbcnms/ui/theme/symphony';
import {fetchQuery, graphql} from 'relay-runtime';
import {sortLexicographically} from '@fbcnms/ui/utils/displayUtils';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

const styles = _ => ({
  avatar: {
    backgroundColor: symphony.palette.B50,
  },
  list: {
    paddingTop: 0,
    paddingBottom: 0,
  },
  listItem: {
    paddingLeft: '24px',
    paddingRight: '24px',
  },
  listAvatar: {
    minWidth: '52px',
  },
  dialogTitle: {
    padding: '24px',
    paddingBottom: '16px',
  },
  serviceTypesDialogContent: {
    padding: 0,
    height: '400px',
    overflowY: 'scroll',
  },
  dialogActions: {
    padding: '24px',
    bottom: 0,
    display: 'flex',
    justifyContent: 'flex-end',
    width: '100%',
    backgroundColor: symphony.palette.white,
  },
});

type Props = ContextRouter & {
  onSelect: (serviceTypeId: string) => void,
  onClose: () => void,
} & WithStyles<typeof styles>;

type State = {
  serviceTypes: Array<ServiceTypesListQuery_serviceType>,
  selectedServiceTypeId: ?string,
};

graphql`
  fragment ServiceTypesListQuery_serviceType on ServiceType {
    id
    name
  }
`;

const serviceTypesQuery = graphql`
  query ServiceTypesListQuery {
    serviceTypes(first: 500)
      @connection(key: "ServiceTypesListQuery_serviceTypes") {
      edges {
        node {
          ...ServiceTypesListQuery_serviceType @relay(mask: false)
        }
      }
    }
  }
`;

class ServiceTypesList extends React.Component<Props, State> {
  state = {
    serviceTypes: [],
    selectedServiceTypeId: null,
  };

  componentDidMount() {
    fetchQuery(RelayEnvironment, serviceTypesQuery).then(response => {
      this.setState({
        serviceTypes: response.serviceTypes.edges.map(x => x.node),
      });
    });
  }

  render() {
    const {selectedServiceTypeId} = this.state;
    const {classes, onSelect, onClose} = this.props;
    const listItems = this.state.serviceTypes
      .slice()
      .sort((serviceTypeA, serviceTypeB) =>
        sortLexicographically(serviceTypeA.name, serviceTypeB.name),
      )
      .map(serviceType => (
        <ListItem
          className={classes.listItem}
          button
          key={serviceType.id}
          selected={selectedServiceTypeId === serviceType.id}
          onClick={event => this.handleListItemClick(event, serviceType)}>
          <ListItemAvatar className={classes.listAvatar}>
            <Avatar className={classes.avatar}>
              <LinearScaleIcon />
            </Avatar>
          </ListItemAvatar>
          <ListItemText primary={serviceType.name} />
        </ListItem>
      ));
    return (
      <>
        <DialogTitle className={classes.dialogTitle}>
          <Text variant="h6">Select a Service type</Text>
        </DialogTitle>
        <DialogContent className={classes.serviceTypesDialogContent}>
          <List>{listItems}</List>
        </DialogContent>
        <DialogActions className={classes.dialogActions}>
          <Button onClick={onClose} skin="regular">
            Cancel
          </Button>
          <Button
            disabled={selectedServiceTypeId === null}
            onClick={() => onSelect(nullthrows(selectedServiceTypeId))}>
            Select
          </Button>
        </DialogActions>
      </>
    );
  }

  handleListItemClick = (event, selectedServiceType) => {
    const selectedServiceTypeId = selectedServiceType?.id;
    this.setState({selectedServiceTypeId});
  };
}

export default withStyles(styles)(withRouter(ServiceTypesList));
