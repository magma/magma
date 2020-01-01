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
import ErrorBoundary from '@fbcnms/ui/components/ErrorBoundary/ErrorBoundary';
import InventoryViewHeader, {DisplayOptions} from '../InventoryViewHeader';
import ProjectCard from './ProjectCard';
import ProjectComparisonViewQueryRenderer from './ProjectComparisonViewQueryRenderer';
import React, {useMemo, useState} from 'react';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {extractEntityIdFromUrl} from '../../common/RouterUtils';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles({
  root: {
    display: 'flex',
    flexDirection: 'column',
    height: '100%',
  },
  searchResults: {
    flexGrow: 1,
    paddingTop: '8px',
  },
  searchResultsTable: {
    paddingLeft: '24px',
    paddingRight: '24px',
  },
});

const ProjectComparisonView = () => {
  const classes = useStyles();
  const [dialogKey, setDialogKey] = useState(1);
  const [dialogOpen, setDialogOpen] = useState(false);
  const {match, history, location} = useRouter();
  const [resultsDisplayMode, setResultsDisplayMode] = useState(
    DisplayOptions.table,
  );

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
      <div className={classes.root}>
        <InventoryViewHeader
          title="Projects"
          onViewToggleClicked={setResultsDisplayMode}
          actionButtons={[
            {
              title: 'Add Project',
              action: () => {
                setDialogOpen(true);
                setDialogKey(dialogKey + 1);
                ServerLogger.info(LogEvents.ADD_PROJECT_BUTTON_CLICKED);
              },
            },
          ]}
        />
        <div className={classes.searchResults}>
          <ProjectComparisonViewQueryRenderer
            className={
              resultsDisplayMode === DisplayOptions.table
                ? classes.searchResultsTable
                : ''
            }
            limit={50}
            filters={[]}
            onProjectSelected={selectedProjectCardId =>
              navigateToProject(selectedProjectCardId)
            }
            displayMode={
              resultsDisplayMode === DisplayOptions.map
                ? DisplayOptions.map
                : DisplayOptions.table
            }
          />
        </div>
      </div>

      <AddProjectDialog
        key={`new_project_${dialogKey}`}
        open={dialogOpen}
        onClose={() => setDialogOpen(false)}
        onProjectTypeSelected={typeId => {
          navigateToAddProject(typeId);
          setDialogOpen(false);
        }}
      />
    </ErrorBoundary>
  );
};

export default ProjectComparisonView;
