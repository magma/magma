/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import useRouter from '@fbcnms/ui/hooks/useRouter';

import AddProjectCard from './AddProjectCard';
import AddProjectDialog from './AddProjectDialog';
import Button from '@fbcnms/ui/components/design-system/Button';
import Card from '@material-ui/core/Card';
import CardContent from '@material-ui/core/CardContent';
import CardFooter from '@fbcnms/ui/components/CardFooter';
import ErrorBoundary from '@fbcnms/ui/components/ErrorBoundary/ErrorBoundary';
import ProjectCard from './ProjectCard';
import ProjectComparisonViewQueryRenderer from './ProjectComparisonViewQueryRenderer';
import React, {useMemo, useState} from 'react';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {extractEntityIdFromUrl} from '../../common/RouterUtils';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  cardRoot: {
    height: '100%',
    display: 'flex',
    flexDirection: 'column',
    paddingLeft: '0px',
    paddingRight: '0px',
  },
  cardContent: {
    paddingLeft: '0px',
    paddingRight: '0px',
    paddingTop: '0px',
    flexGrow: 1,
    width: '100%',
  },
  root: {
    display: 'flex',
    flexDirection: 'column',
    backgroundColor: theme.palette.common.white,
    flexGrow: 1,
    height: '100%',
  },
  searchResults: {
    flexGrow: 1,
    display: 'flex',
    flexDirection: 'column',
  },
  bar: {
    display: 'flex',
    flexDirection: 'row',
    boxShadow: '0px 2px 2px 0px rgba(0, 0, 0, 0.1)',
  },
  searchBar: {
    flexGrow: 1,
  },
}));

const ProjectComparisonView = () => {
  const classes = useStyles();
  const [dialogKey, setDialogKey] = useState(1);
  const [dialogOpen, setDialogOpen] = useState(false);
  const {match, history, location} = useRouter();

  const selectedProjectTypeId = useMemo(
    () => extractEntityIdFromUrl('projectType', location.search),
    [location.search],
  );

  const selectedProjectCardId = useMemo(
    () => extractEntityIdFromUrl('project', location.search),
    [location.search],
  );

  function navigateToProject(selectedProjectCardId: ?string) {
    history.push(
      match.url +
        (selectedProjectCardId ? `?project=${selectedProjectCardId}` : ''),
    );
  }

  function navigateToAddProject(selectedProjectTypeId: ?string) {
    history.push(
      match.url +
        (selectedProjectTypeId ? `?projectType=${selectedProjectTypeId}` : ''),
    );
  }

  if (selectedProjectTypeId != null) {
    return (
      <ErrorBoundary>
        <AddProjectCard projectTypeId={selectedProjectTypeId} />
      </ErrorBoundary>
    );
  }
  if (selectedProjectCardId != null) {
    return (
      <ErrorBoundary>
        <ProjectCard
          projectId={selectedProjectCardId}
          onProjectExecuted={() => {}}
          onProjectRemoved={() => navigateToProject(null)}
        />
      </ErrorBoundary>
    );
  }
  return (
    <ErrorBoundary>
      <Card className={classes.cardRoot}>
        <CardContent className={classes.cardContent}>
          <div className={classes.root}>
            <div className={classes.searchResults}>
              <ProjectComparisonViewQueryRenderer
                limit={50}
                filters={[]}
                displayMode={'table'}
                onProjectSelected={selectedProjectCardId =>
                  navigateToProject(selectedProjectCardId)
                }
              />
            </div>
          </div>
        </CardContent>
        <CardFooter alignItems="left">
          <Button
            onClick={() => {
              setDialogOpen(true);
              setDialogKey(dialogKey + 1);
              ServerLogger.info(LogEvents.ADD_PROJECT_BUTTON_CLICKED);
            }}>
            New Project
          </Button>
          <AddProjectDialog
            key={`new_project_${dialogKey}`}
            open={dialogOpen}
            onClose={() => setDialogOpen(false)}
            onProjectTypeSelected={typeId => {
              navigateToAddProject(typeId);
              setDialogOpen(false);
            }}
          />
        </CardFooter>
      </Card>
    </ErrorBoundary>
  );
};

export default ProjectComparisonView;
