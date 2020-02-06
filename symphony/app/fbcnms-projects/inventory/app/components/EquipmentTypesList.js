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
import type {EquipmentType} from '../common/EquipmentType';
import type {WithStyles} from '@material-ui/core';

import Avatar from '@material-ui/core/Avatar';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import React from 'react';
import RelayEnvironment from '../common/RelayEnvironment.js';
import RouterIcon from '@material-ui/icons/Router';
import {fetchQuery, graphql} from 'relay-runtime';
import {sortLexicographically} from '@fbcnms/ui/utils/displayUtils';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

const styles = _theme => ({
  avatar: {
    marginRight: '8px',
  },
});

type Props = ContextRouter & {
  onSelect: ?(equipmentType: ?EquipmentType) => void,
} & WithStyles<typeof styles>;

type State = {
  errorMessage: ?string,
  equipmentTypes: Array<EquipmentType>,
  selectedEquipmentType: ?EquipmentType,
  showDialog: boolean,
};

const equipmentQuery = graphql`
  query EquipmentTypesListQuery {
    equipmentTypes {
      edges {
        node {
          id
          name
          propertyTypes {
            ...PropertyTypeFormField_propertyType @relay(mask: false)
          }
        }
      }
    }
  }
`;

class EquipmentTypesList extends React.Component<Props, State> {
  state = {
    errorMessage: null,
    equipmentTypes: [],
    selectedEquipmentType: null,
    showDialog: false,
  };

  componentDidMount() {
    fetchQuery(RelayEnvironment, equipmentQuery).then(response => {
      this.setState({
        equipmentTypes: response.equipmentTypes.edges.map(x => x.node),
      });
    });
  }

  render() {
    const {selectedEquipmentType} = this.state;
    const {classes} = this.props;
    const listItems = this.state.equipmentTypes
      .slice()
      .sort((equipmentTypeA, equipmentTypeB) =>
        sortLexicographically(equipmentTypeA.name, equipmentTypeB.name),
      )
      .map(equipmentType => (
        <ListItem
          dense
          button
          key={equipmentType.id}
          selected={
            selectedEquipmentType &&
            selectedEquipmentType.id === equipmentType.id
          }
          onClick={event => this.handleListItemClick(event, equipmentType)}>
          <Avatar className={classes.avatar}>
            <RouterIcon />
          </Avatar>
          <ListItemText primary={equipmentType.name} />
        </ListItem>
      ));

    return <List>{listItems}</List>;
  }

  handleListItemClick = (event, selectedEquipmentType) => {
    this.setState(
      {selectedEquipmentType},
      () => this.props.onSelect && this.props.onSelect(selectedEquipmentType),
    );
  };
}

export default withStyles(styles)(withRouter(EquipmentTypesList));
