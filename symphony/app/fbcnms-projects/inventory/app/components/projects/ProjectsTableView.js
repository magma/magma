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
import React, {useMemo} from 'react';
import Table from '@fbcnms/ui/components/design-system/Table/Table';
import fbt from 'fbt';
import {createFragmentContainer, graphql} from 'react-relay';

type Props = {
  projects: ProjectsTableView_projects,
  onProjectSelected: string => void,
} & ContextRouter;

const ProjectsTableView = (props: Props) => {
  const {projects, onProjectSelected} = props;

  const data = useMemo(
    () => projects.map(project => ({...project, key: project.id})),
    [projects],
  );

  if (projects.length === 0) {
    return null;
  }

  return (
    <Table
      data={data}
      columns={[
        {
          key: 'name',
          title: 'Name',
          render: row => (
            <Button
              variant="text"
              useEllipsis={true}
              onClick={() => onProjectSelected(row.id)}>
              {row.name}
            </Button>
          ),
          getSortingValue: row => row.name,
        },
        {
          key: 'type',
          title: `${fbt('Template', '')}`,
          getSortingValue: row => row.type?.name,
          render: row => row.type?.name ?? '',
        },
        {
          key: 'location',
          title: 'Location',
          getSortingValue: row => row.location?.name,
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
          getSortingValue: row => row?.createdBy?.email,
          render: row => row?.createdBy?.email ?? '',
        },
      ]}
    />
  );
};

export default createFragmentContainer(ProjectsTableView, {
  projects: graphql`
    fragment ProjectsTableView_projects on Project @relay(plural: true) {
      id
      name
      createdBy {
        email
      }
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
