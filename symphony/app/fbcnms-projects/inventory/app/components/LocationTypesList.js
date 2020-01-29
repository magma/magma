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
import type {LocationType} from '../common/LocationType';
import type {WithStyles} from '@material-ui/core';

import Avatar from '@material-ui/core/Avatar';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import LocationOnIcon from '@material-ui/icons/LocationOn';
import React from 'react';
import RelayEnvironment from '../common/RelayEnvironment.js';
import {fetchQuery, graphql} from 'relay-runtime';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

type Props = ContextRouter & {
  onSelect: ?(locationType: ?LocationType) => void,
} & WithStyles<typeof styles>;

type State = {
  errorMessage: ?string,
  locationTypes: Array<LocationType>,
  selectedLocationType: ?LocationType,
  showDialog: boolean,
};

const locationQuery = graphql`
  query LocationTypesListQuery {
    locationTypes(first: 500) {
      edges {
        node {
          id
          name
        }
      }
    }
  }
`;

const styles = theme => ({
  gutters: {
    paddingLeft: theme.spacing(3),
    paddingRight: theme.spacing(3),
  },
  avatar: {
    marginRight: '8px',
  },
});

class LocationTypesList extends React.Component<Props, State> {
  state = {
    errorMessage: null,
    locationTypes: [],
    selectedLocationType: null,
    showDialog: false,
  };

  componentDidMount() {
    fetchQuery(RelayEnvironment, locationQuery).then(response => {
      this.setState({
        locationTypes: response.locationTypes.edges.map(x => x.node),
      });
    });
  }

  render() {
    const {classes} = this.props;
    const {selectedLocationType} = this.state;
    const listItems = this.state.locationTypes.map(locationType => (
      <ListItem
        classes={{gutters: classes.gutters}}
        dense
        button
        key={locationType.id}
        selected={
          selectedLocationType && selectedLocationType.id === locationType.id
        }
        onClick={event => this.handleListItemClick(event, locationType)}>
        <Avatar className={classes.avatar}>
          <LocationOnIcon />
        </Avatar>
        <ListItemText primary={locationType.name} />
      </ListItem>
    ));

    return <List>{listItems}</List>;
  }

  handleListItemClick = (event, selectedLocationType) => {
    this.setState(
      {selectedLocationType},
      () => this.props.onSelect && this.props.onSelect(selectedLocationType),
    );
  };
}

export default withStyles(styles)(withRouter(LocationTypesList));
