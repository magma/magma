/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {WorkOrderProjectTypesQueryResponse} from './__generated__/WorkOrderProjectTypesQuery.graphql';

import AddEditProjectTypeCard from './AddEditProjectTypeCard';
import Button from '@material-ui/core/Button';
import InventoryQueryRenderer from '../InventoryQueryRenderer';
import ProjectTypeCard from './ProjectTypeCard';
import React, {useState} from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {graphql} from 'relay-runtime';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  root: {
    flexGrow: 1,
    padding: '24px 16px',
  },
  header: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'flex-end',
    margin: '0px 8px',
    marginBottom: '24px',
  },
  titleContainer: {
    flexGrow: 1,
  },
  title: {
    fontSize: '20px',
    lineHeight: '24px',
    fontWeight: 500,
    color: theme.palette.blueGrayDark,
    display: 'block',
  },
  subtitle: {
    fontSize: '14px',
    lineHeight: '24px',
    color: '#73839e',
  },
  addButtonContainer: {
    display: 'flex',
  },
  addButton: {
    paddingLeft: '16px',
    paddingRight: '16px',
    marginLeft: 'auto',
  },
  typeCards: {
    display: 'flex',
    flexWrap: 'wrap',
    flexDirection: 'row',
  },
  typeCard: {
    padding: '8px',
    flexBasis: '16.66%', // 6 cards
  },
  '@media (max-width: 1950px)': {
    typeCard: {
      flexBasis: '20%', // 5 cards
    },
  },
  '@media (max-width: 1600px)': {
    typeCard: {
      flexBasis: '25%', // 4 cards
    },
  },
  '@media (max-width: 1024px)': {
    typeCard: {
      flexBasis: '33.33%', // 3 cards
    },
  },
  '@media (max-width: 650px)': {
    typeCard: {
      flexBasis: '100%', // 1 card
    },
  },
}));

const projectTypesQuery = graphql`
  query WorkOrderProjectTypesQuery {
    projectTypes(first: 50)
      @connection(key: "WorkOrderProjectTypesQuery_projectTypes") {
      edges {
        node {
          id
          ...ProjectTypeCard_projectType
          ...AddEditProjectTypeCard_editingProjectType
        }
      }
    }
    workOrderTypes {
      edges {
        node {
          ...ProjectTypeWorkOrderTemplatesPanel_workOrderTypes
        }
      }
    }
  }
`;

const WorkOrderProjectTypes = () => {
  const classes = useStyles();
  const [editingProjectType, setEditingProjectType] = useState(null);
  const [showAddEditCard, setShowAddEditCard] = useState(false);
  const hideAddEditCard = () => {
    setEditingProjectType(null);
    setShowAddEditCard(false);
  };
  return (
    <InventoryQueryRenderer
      query={projectTypesQuery}
      variables={{}}
      render={(props: WorkOrderProjectTypesQueryResponse) => {
        if (showAddEditCard || editingProjectType) {
          const workOrderTypes = props.workOrderTypes?.edges ?? [];
          return (
            <div className={classes.root}>
              <AddEditProjectTypeCard
                workOrderTypes={workOrderTypes
                  .map(e => e?.node)
                  .filter(Boolean)}
                editingProjectType={editingProjectType}
                onCancelClicked={hideAddEditCard}
                onProjectTypeSaved={hideAddEditCard}
              />
            </div>
          );
        }

        return (
          <div className={classes.root}>
            <div className={classes.header}>
              <div className={classes.titleContainer}>
                <Text className={classes.title}>Project Templates</Text>
                <Text className={classes.subtitle}>
                  Create and manage reusable project workflows
                </Text>
              </div>
              <div className={classes.addButtonContainer}>
                <Button
                  className={classes.addButton}
                  color="primary"
                  variant="contained"
                  onClick={() => {
                    ServerLogger.info(
                      LogEvents.ADD_PROJECT_TEMPLATE_BUTTON_CLICKED,
                    );
                    setShowAddEditCard(true);
                  }}>
                  Add Project Template
                </Button>
              </div>
            </div>
            <div className={classes.typeCards}>
              {(props.projectTypes?.edges ?? [])
                .map(edge => edge.node)
                .filter(Boolean)
                .map(projectType => (
                  <div key={projectType.id} className={classes.typeCard}>
                    <ProjectTypeCard
                      projectType={projectType}
                      onEditClicked={() => setEditingProjectType(projectType)}
                    />
                  </div>
                ))}
            </div>
          </div>
        );
      }}
    />
  );
};

export default WorkOrderProjectTypes;
