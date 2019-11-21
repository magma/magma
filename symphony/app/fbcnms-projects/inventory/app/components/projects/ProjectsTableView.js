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
import type {ProjectsTableView_projects} from './__generated__/ProjectsTableView_projects.graphql';

import Button from '@fbcnms/ui/components/design-system/Button';
import LocationLink from '../location/LocationLink';
import React, {useMemo, useState} from 'react';
import Table from '@fbcnms/ui/components/design-system/Table/Table';
import {createFragmentContainer, graphql} from 'react-relay';

type Props = {
  className?: string,
  projects: ProjectsTableView_projects,
  onProjectSelected: string => void,
} & ContextRouter;

const ProjectsTableView = (props: Props) => {
  const {projects, onProjectSelected, className} = props;

  const [sortDirection, setSortDirection] = useState('desc');
  const [sortColumn, setSortColumn] = useState('name');

  const sortedProjects = useMemo(
    () =>
      projects
        .slice()
        .sort(
          (p1, p2) =>
            p1[sortColumn].localeCompare(p2[sortColumn]) *
            (sortDirection === 'asc' ? -1 : 1),
        )
        .map(project => ({...project, key: project.id})),
    [projects, sortColumn, sortDirection],
  );

  if (projects.length === 0) {
    return null;
  }

  return (
    <div className={className}>
      <Table
        data={sortedProjects}
        onSortClicked={col => {
          if (sortColumn === col) {
            setSortDirection(sortDirection === 'asc' ? 'desc' : 'asc');
          } else {
            setSortColumn(col);
            setSortDirection('desc');
          }
        }}
        columns={[
          {
            key: 'name',
            title: 'Name',
            render: row => (
              <Button variant="text" onClick={() => onProjectSelected(row.id)}>
                {row.name}
              </Button>
            ),
            sortable: true,
            sortDirection: sortColumn === 'name' ? sortDirection : undefined,
          },
          {
            key: 'type',
            title: 'Type',
            render: row => row.type?.name ?? '',
          },
          {
            key: 'location',
            title: 'Location',
            render: row =>
              row.location ? (
                <LocationLink title={row.location.name} id={row.location.id} />
              ) : (
                ''
              ),
          },
          {
            key: 'owner',
            title: 'Owner',
            render: row => row.creator ?? '',
          },
        ]}
      />
    </div>
  );
};

export default createFragmentContainer(ProjectsTableView, {
  projects: graphql`
    fragment ProjectsTableView_projects on Project @relay(plural: true) {
      id
      name
      creator
      location {
        id
        name
      }
      type {
        id
        name
      }
    }
  `,
});
