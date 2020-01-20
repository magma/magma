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
import type {EquipmentPortTypeItem_equipmentPortType} from './__generated__/EquipmentPortTypeItem_equipmentPortType.graphql';
import type {WithStyles} from '@material-ui/core';

import AddEditEquipmentPortTypeCard from './AddEditEquipmentPortTypeCard';
import Button from '@fbcnms/ui/components/design-system/Button';
import ConfigueTitle from '@fbcnms/ui/components/ConfigureTitle';
import EquipmentPortTypeItem from './EquipmentPortTypeItem';
import InventoryQueryRenderer from '../InventoryQueryRenderer';
import React from 'react';
import withInventoryErrorBoundary from '../../common/withInventoryErrorBoundary';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {graphql} from 'relay-runtime';
import {sortLexicographically} from '@fbcnms/ui/utils/displayUtils';
import {withRouter} from 'react-router-dom';
import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  header: {
    margin: '10px',
    display: 'flex',
    justifyContent: 'space-between',
  },
  root: {
    width: '100%',
    marginTop: '15px',
  },
  paper: {
    flexGrow: 1,
    overflowY: 'hidden',
  },
  content: {
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'flex-start',
  },
  listItem: {
    marginBottom: theme.spacing(),
  },
  addButton: {
    marginLeft: 'auto',
  },
  addButtonContainer: {
    display: 'flex',
  },
  typesList: {
    padding: '24px',
  },
  title: {
    marginLeft: '10px',
  },
  firstRow: {
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
});

type Props = ContextRouter & WithStyles<typeof styles>;

type State = {
  errorMessage: ?string,
  editingEquipmentPortType: ?EquipmentPortTypeItem_equipmentPortType,
  showAddEditCard: boolean,
};

const equipmentPortTypesQuery = graphql`
  query EquipmentPortTypesQuery {
    equipmentPortTypes(first: 50)
      @connection(key: "EquipmentPortTypes_equipmentPortTypes") {
      edges {
        node {
          ...EquipmentPortTypeItem_equipmentPortType
          ...AddEditEquipmentPortTypeCard_editingEquipmentPortType
          id
          name
        }
      }
    }
  }
`;

class EquipmentPortTypes extends React.Component<Props, State> {
  state = {
    errorMessage: null,
    editingEquipmentPortType: null,
    showAddEditCard: false,
  };

  render() {
    const {classes} = this.props;
    const {showAddEditCard, editingEquipmentPortType} = this.state;
    return (
      <InventoryQueryRenderer
        query={equipmentPortTypesQuery}
        variables={{}}
        render={props => {
          if (showAddEditCard) {
            return (
              <div className={classes.paper}>
                <AddEditEquipmentPortTypeCard
                  key={'new_port_type'}
                  open={showAddEditCard}
                  onClose={this.hideNewEquipmentPortTypeCard}
                  onSave={this.savePort}
                  editingEquipmentPortType={editingEquipmentPortType}
                />
              </div>
            );
          }

          const listItems = (props.equipmentPortTypes.edges ?? [])
            .map(edge => edge?.node)
            .filter(Boolean)
            .sort((eqTypeA, eqTypeB) =>
              sortLexicographically(eqTypeA.name, eqTypeB.name),
            )
            .map(equipmentPortType => (
              <div
                className={classes.listItem}
                key={`equipmentPortType_${equipmentPortType.id}`}>
                <EquipmentPortTypeItem
                  equipmentPortType={equipmentPortType}
                  onEdit={() =>
                    this.showAddEditEquipmentPortTypeCard(equipmentPortType)
                  }
                />
              </div>
            ));
          return (
            <div className={classes.typesList}>
              <div className={classes.firstRow}>
                <ConfigueTitle
                  className={classes.title}
                  title={'Port Types'}
                  subtitle={'Manage the types of ports in your inventory'}
                />
                <div className={classes.addButtonContainer}>
                  <Button
                    onClick={() => this.showAddEditEquipmentPortTypeCard(null)}>
                    Add Port Type
                  </Button>
                </div>
              </div>
              <div className={classes.root}>
                <div>{listItems}</div>
              </div>
            </div>
          );
        }}
      />
    );
  }

  showAddEditEquipmentPortTypeCard = (
    eqType: ?EquipmentPortTypeItem_equipmentPortType,
  ) => {
    ServerLogger.info(LogEvents.ADD_EQUIPMENT_TYPE_BUTTON_CLICKED);
    this.setState({editingEquipmentPortType: eqType, showAddEditCard: true});
  };

  hideNewEquipmentPortTypeCard = () =>
    this.setState({editingEquipmentPortType: null, showAddEditCard: false});
  savePort = () => {
    ServerLogger.info(LogEvents.SAVE_EQUIPMENT_TYPE_BUTTON_CLICKED);
    this.setState({
      editingEquipmentPortType: null,
      showAddEditCard: false,
    });
  };
}

export default withStyles(styles)(
  withRouter(withInventoryErrorBoundary(EquipmentPortTypes)),
);
