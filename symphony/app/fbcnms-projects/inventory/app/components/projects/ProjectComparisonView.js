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
import ListAltIcon from '@material-ui/icons/ListAlt';
import MapButtonGroup from '@fbcnms/ui/components/map/MapButtonGroup';
import MapIcon from '@material-ui/icons/Map';
import ProjectCard from './ProjectCard';
import ProjectComparisonViewQueryRenderer from './ProjectComparisonViewQueryRenderer';
import React, {useMemo, useState} from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import classNames from 'classnames';
import symphony from '@fbcnms/ui/theme/symphony';
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
    backgroundColor: symphony.palette.background,
  },
  root: {
    display: 'flex',
    flexDirection: 'column',
    backgroundColor: theme.palette.common.white,
    height: '100%',
  },
  searchResults: {
    flexGrow: 1,
    display: 'flex',
    flexDirection: 'column',
    backgroundColor: symphony.palette.background,
  },
  addProjectButton: {
    alignSelf: 'flex-end',
  },
  bar: {
    display: 'flex',
    flexDirection: 'row',
    boxShadow: '0px 2px 2px 0px rgba(0, 0, 0, 0.1)',
  },
  groupButtons: {
    display: 'flex',
    justifyContent: 'flex-end',
  },
  buttonContent: {
    paddingTop: '4px',
  },
  titleContainer: {
    margin: '32px',
    display: 'flex',
  },
  title: {
    flexGrow: 1,
    display: 'block',
  },
  searchBar: {
    flexGrow: 1,
  },
  comparisionViewTable: {
    margin: '0px 32px',
  },
}));

const ProjectComparisonView = () => {
  const classes = useStyles();
  const [dialogKey, setDialogKey] = useState(1);
  const [dialogOpen, setDialogOpen] = useState(false);
  const {match, history, location} = useRouter();
  const [resultsDisplayMode, setResultsDisplayMode] = useState('table');

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
      <div className={classes.cardRoot}>
        <div className={classes.root}>
          <div className={classes.bar}>
            <div className={classes.searchBar} />
            <MapButtonGroup
              onIconClicked={id => {
                setResultsDisplayMode(id === 'table' ? 'table' : 'map');
              }}
              buttons={[
                {
                  item: <ListAltIcon className={classes.buttonContent} />,
                  id: 'table',
                },
                {
                  item: <MapIcon className={classes.buttonContent} />,
                  id: 'map',
                },
              ]}
            />
          </div>
          <div className={classes.searchResults}>
            <div className={classes.titleContainer}>
              <Text className={classes.title} variant="h6">
                Projects
              </Text>
              <Button
                className={classes.addProjectButton}
                onClick={() => {
                  setDialogOpen(true);
                  setDialogKey(dialogKey + 1);
                  ServerLogger.info(LogEvents.ADD_PROJECT_BUTTON_CLICKED);
                }}>
                Add Project
              </Button>
            </div>
            <ProjectComparisonViewQueryRenderer
              className={classNames({
                [classes.comparisionViewTable]: resultsDisplayMode === 'table',
              })}
              limit={50}
              filters={[]}
              displayMode={'table'}
              onProjectSelected={selectedProjectCardId =>
                navigateToProject(selectedProjectCardId)
              }
              resultsDisplayMode={resultsDisplayMode}
            />
          </div>
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
