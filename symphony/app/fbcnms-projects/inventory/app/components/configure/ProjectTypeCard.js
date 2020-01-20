/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type ProjectTypeCard_projectType from './__generated__/ProjectTypeCard_projectType.graphql';

import ProjectTypeDeleteButton from './ProjectTypeDeleteButton';
import ProjectTypeWorkOrdersCount from './ProjectTypeWorkOrdersCount';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import Typography from '@material-ui/core/Typography';
import classNames from 'classnames';
import {createFragmentContainer, graphql} from 'react-relay';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  root: {
    height: '257px',
    backgroundColor: theme.palette.common.white,
    boxShadow: '0px 1px 4px 0px rgba(0,0,0,0.17)',
    padding: '24px 24px 16px 24px',
    overflow: 'hidden',
    borderRadius: '4px',
  },
  name: {
    fontSize: '20px',
    lineHeight: '28px',
    fontWeight: 500,
    color: theme.palette.blueGrayDark,
    marginBottom: '10px',
  },
  nameContainer: {
    display: 'flex',
    flexDirection: 'row',
    flexGrow: 1,
  },
  descriptionContainer: {
    overflow: 'hidden',
    flexGrow: 1,
  },
  description: {
    fontSize: '16px',
    lineHeight: '24px',
    color: '#8895ad',
    marginBottom: '8px',
    overflow: 'hidden',
  },
  deleteButton: {
    flexGrow: 1,
    display: 'flex',
    justifyContent: 'flex-end',
  },
  container: {
    display: 'flex',
    flexDirection: 'column',
    height: '100%',
  },
  buttonContainer: {
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'flex-end',
  },
  divider: {
    borderTop: '1px solid #edf0f9',
    margin: '16px 0px',
  },
  manageButton: {
    fontSize: '16px',
    lineHeight: '24px',
    color: theme.palette.primary.main,
    cursor: 'pointer',
  },
}));

type Props = {
  className?: string,
  projectType: ProjectTypeCard_projectType,
  onEditClicked: () => void,
};

const ProjectTypeCard = ({className, projectType, onEditClicked}: Props) => {
  const {name, description, numberOfProjects, workOrders} = projectType;
  const classes = useStyles();
  return (
    <div className={classNames(classes.root, className)}>
      <div className={classes.container}>
        <div className={classes.descriptionContainer}>
          <div className={classes.nameContainer}>
            <Text className={classes.name}>{name}</Text>
            {numberOfProjects === 0 && (
              <ProjectTypeDeleteButton
                className={classes.deleteButton}
                projectType={projectType}
              />
            )}
          </div>
          <Text className={classes.description}>{description}</Text>
        </div>
        <div className={classes.divider} />
        <div className={classes.buttonContainer}>
          <ProjectTypeWorkOrdersCount count={workOrders.length} />
          <Typography
            className={classes.manageButton}
            color="primary"
            onClick={onEditClicked}>
            Edit
          </Typography>
        </div>
      </div>
    </div>
  );
};

export default createFragmentContainer(ProjectTypeCard, {
  projectType: graphql`
    fragment ProjectTypeCard_projectType on ProjectType {
      id
      name
      description
      numberOfProjects
      workOrders {
        id
      }
    }
  `,
});
