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
import type {WorkOrderTypesListQuery_workOrderType} from './__generated__/WorkOrderTypesListQuery_workOrderType.graphql';

import Avatar from '@material-ui/core/Avatar';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemAvatar from '@material-ui/core/ListItemAvatar';
import ListItemText from '@material-ui/core/ListItemText';
import React from 'react';
import RelayEnvironment from '../../common/RelayEnvironment.js';
import WorkIcon from '@material-ui/icons/Work';
import {fetchQuery, graphql} from 'relay-runtime';
import {sortLexicographically} from '@fbcnms/ui/utils/displayUtils';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

const styles = _theme => ({
  avatar: {
    backgroundColor: '#e4f2ff',
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
});

type Props = ContextRouter & {
  onSelect: ?(workOrderTypeId: ?string) => void,
} & WithStyles<typeof styles>;

type State = {
  errorMessage: ?string,
  workOrderTypes: Array<WorkOrderTypesListQuery_workOrderType>,
  selectedWorkOrderTypeId: ?string,
  showDialog: boolean,
};

graphql`
  fragment WorkOrderTypesListQuery_workOrderType on WorkOrderType {
    id
    name
  }
`;

const workOrderTypesQuery = graphql`
  query WorkOrderTypesListQuery {
    workOrderTypes(first: 500)
      @connection(key: "WorkOrderTypesListQuery_workOrderTypes") {
      edges {
        node {
          ...WorkOrderTypesListQuery_workOrderType @relay(mask: false)
        }
      }
    }
  }
`;

class WorkOrderTypesList extends React.Component<Props, State> {
  state = {
    errorMessage: null,
    workOrderTypes: [],
    selectedWorkOrderTypeId: null,
    showDialog: false,
  };

  componentDidMount() {
    fetchQuery(RelayEnvironment, workOrderTypesQuery).then(response => {
      this.setState({
        workOrderTypes: response.workOrderTypes.edges.map(x => x.node),
      });
    });
  }

  render() {
    const {selectedWorkOrderTypeId} = this.state;
    const {classes} = this.props;
    const listItems = this.state.workOrderTypes
      .slice()
      .sort((workOrderTypeA, workOrderTypeB) =>
        sortLexicographically(workOrderTypeA.name, workOrderTypeB.name),
      )
      .map(workOrderType => (
        <ListItem
          className={classes.listItem}
          button
          key={workOrderType.id}
          selected={selectedWorkOrderTypeId === workOrderType.id}
          onClick={event => this.handleListItemClick(event, workOrderType)}>
          <ListItemAvatar className={classes.listAvatar}>
            <Avatar className={classes.avatar}>
              <WorkIcon />
            </Avatar>
          </ListItemAvatar>
          <ListItemText primary={workOrderType.name} />
        </ListItem>
      ));
    return <List>{listItems}</List>;
  }

  handleListItemClick = (event, selectedWorkOrderType) => {
    const selectedWorkOrderTypeId = selectedWorkOrderType?.id;
    this.setState(
      {selectedWorkOrderTypeId},
      () => this.props.onSelect && this.props.onSelect(selectedWorkOrderTypeId),
    );
  };
}

export default withStyles(styles)(withRouter(WorkOrderTypesList));
