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
import ErrorBoundary from '@fbcnms/ui/components/ErrorBoundary/ErrorBoundary';
import FormActionWithPermissions from '../../common/FormActionWithPermissions';
import InventoryView, {DisplayOptions} from '../InventoryViewContainer';
import ProjectCard from './ProjectCard';
import ProjectComparisonViewQueryRenderer from './ProjectComparisonViewQueryRenderer';
import React, {useMemo, useState} from 'react';
import fbt from 'fbt';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {extractEntityIdFromUrl} from '../../common/RouterUtils';

const ProjectComparisonView = () => {
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
  const header = {
    title: 'Projects',
    actionButtons: [
      <FormActionWithPermissions
        permissions={{
          entity: 'project',
          action: 'create',
          ignoreTypes: true,
        }}>
        <Button
          onClick={() => {
            setDialogOpen(true);
            setDialogKey(dialogKey + 1);
            ServerLogger.info(LogEvents.ADD_PROJECT_BUTTON_CLICKED);
          }}>
          <fbt desc="">Create Project</fbt>
        </Button>
      </FormActionWithPermissions>,
    ],
  };
  return (
    <ErrorBoundary>
      <InventoryView
        header={header}
        onViewToggleClicked={setResultsDisplayMode}
        permissions={{
          entity: 'project',
        }}>
        <ProjectComparisonViewQueryRenderer
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
        <AddProjectDialog
          key={`new_project_${dialogKey}`}
          open={dialogOpen}
          onClose={() => setDialogOpen(false)}
          onProjectTypeSelected={typeId => {
            navigateToAddProject(typeId);
            setDialogOpen(false);
          }}
        />
      </InventoryView>
    </ErrorBoundary>
  );
};

export default ProjectComparisonView;
