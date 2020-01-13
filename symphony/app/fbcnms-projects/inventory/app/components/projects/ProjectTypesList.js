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
import type {ProjectTypesListQuery_projectType} from './__generated__/ProjectTypesListQuery_projectType.graphql';
import type {WithStyles} from '@material-ui/core';

import Avatar from '@material-ui/core/Avatar';
import List from '@material-ui/core/List';
import ListIcon from '@material-ui/icons/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemAvatar from '@material-ui/core/ListItemAvatar';
import ListItemText from '@material-ui/core/ListItemText';
import React from 'react';
import RelayEnvironment from '../../common/RelayEnvironment.js';
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
  onSelect: ?(projectTypeId: ?string) => void,
} & WithStyles<typeof styles>;

type State = {
  errorMessage: ?string,
  projectTypes: Array<ProjectTypesListQuery_projectType>,
  selectedProjectTypeId: ?string,
  showDialog: boolean,
};

graphql`
  fragment ProjectTypesListQuery_projectType on ProjectType {
    id
    name
  }
`;

const projectTypesQuery = graphql`
  query ProjectTypesListQuery {
    projectTypes(first: 50) {
      edges {
        node {
          ...ProjectTypesListQuery_projectType @relay(mask: false)
        }
      }
    }
  }
`;

class ProjectTypesList extends React.Component<Props, State> {
  state = {
    errorMessage: null,
    projectTypes: [],
    selectedProjectTypeId: null,
    showDialog: false,
  };

  componentDidMount() {
    fetchQuery(RelayEnvironment, projectTypesQuery).then(response => {
      this.setState({
        projectTypes: response.projectTypes.edges.map(x => x.node),
      });
    });
  }

  render() {
    const {selectedProjectTypeId} = this.state;
    const {classes} = this.props;
    const listItems = this.state.projectTypes
      .slice()
      .sort((projectTypeA, projectTypeB) =>
        sortLexicographically(projectTypeA.name, projectTypeB.name),
      )
      .map(projectType => (
        <ListItem
          className={classes.listItem}
          button
          key={projectType.id}
          selected={selectedProjectTypeId === projectType.id}
          onClick={event => this.handleListItemClick(event, projectType)}>
          <ListItemAvatar className={classes.listAvatar}>
            <Avatar className={classes.avatar}>
              <ListIcon />
            </Avatar>
          </ListItemAvatar>
          <ListItemText primary={projectType.name} />
        </ListItem>
      ));
    return <List>{listItems}</List>;
  }

  handleListItemClick = (event, selectedProjectType) => {
    const selectedProjectTypeId = selectedProjectType?.id;
    this.setState(
      {selectedProjectTypeId},
      () => this.props.onSelect && this.props.onSelect(selectedProjectTypeId),
    );
  };
}

export default withStyles(styles)(withRouter(ProjectTypesList));
