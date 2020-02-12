/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import React, {useMemo, useState} from 'react';
import Table from '../../components/design-system/Table/Table';
import {STORY_CATEGORIES} from '../storybookUtils';
import {makeStyles} from '@material-ui/styles';
import {storiesOf} from '@storybook/react';

const DATA = [
  {
    firstName: 'Meghan',
    lastName: 'Bishop',
    birthDate: 'December 30, 2019',
    city: 'Tel Aviv',
  },
  {
    firstName: 'Sara',
    lastName: 'Porter',
    birthDate: 'June 28, 1990',
    city: 'Raanana',
  },
  {
    firstName: 'Dolev',
    lastName: 'Hadar',
    birthDate: 'Febuary 11, 1990',
    city: 'Tel Aviv',
  },
  {
    firstName: 'Walter',
    lastName: 'Jenning',
    birthDate: 'July 11, 2001',
    city: 'Ramat Gan',
  },
];

const useStyles = makeStyles(_theme => ({
  root: {
    width: '100%',
  },
  table: {
    marginBottom: '24px',
  },
}));

type DataType = {
  firstName: string,
  lastName: string,
  birthDate: string,
  city: string,
};

const TablesRoot = () => {
  const classes = useStyles();
  const [selectedIds, setSelectedIds] = useState([]);

  const sortData = (col, sortDirection) =>
    DATA.slice().sort(
      (d1: DataType, d2: DataType) =>
        d1[col].localeCompare(d2[col]) * (sortDirection === 'asc' ? -1 : 1),
    );
  const [sortDirection, setSortDirection] = useState('desc');
  const [sortColumn, setSortColumn] = useState('firstName');
  const sortedData = useMemo(() => sortData(sortColumn, sortDirection), [
    sortColumn,
    sortDirection,
  ]);

  return (
    <div className={classes.root}>
      <Table
        className={classes.table}
        data={DATA}
        columns={[
          {key: '0', title: 'First Name', render: row => row.firstName},
          {key: '1', title: 'Last Name', render: row => row.lastName},
          {key: '2', title: 'Birth Date', render: row => row.birthDate},
          {key: '3', title: 'City', render: row => row.city},
        ]}
      />
      <Table
        className={classes.table}
        showSelection
        selectedIds={selectedIds}
        onSelectionChanged={ids => setSelectedIds(ids)}
        data={DATA}
        columns={[
          {key: '0', title: 'First Name', render: row => row.firstName},
          {key: '1', title: 'Last Name', render: row => row.lastName},
          {key: '2', title: 'Birth Date', render: row => row.birthDate},
          {key: '3', title: 'City', render: row => row.city},
        ]}
      />
      <Table
        className={classes.table}
        data={sortedData}
        columns={[
          {
            key: 'firstName',
            title: 'First Name',
            render: row => row.firstName,
            sortable: true,
            sortDirection:
              sortColumn === 'firstName' ? sortDirection : undefined,
          },
          {
            key: 'lastName',
            title: 'Last Name',
            render: row => row.lastName,
            sortable: true,
            sortDirection:
              sortColumn === 'lastName' ? sortDirection : undefined,
          },
          {
            key: 'birthDate',
            title: 'Birth Date',
            render: row => row.birthDate,
          },
          {
            key: 'city',
            title: 'City',
            render: row => row.city,
            sortable: true,
            sortDirection: sortColumn === 'city' ? sortDirection : undefined,
          },
        ]}
        onSortClicked={col => {
          if (sortColumn === col) {
            setSortDirection(sortDirection === 'asc' ? 'desc' : 'asc');
          } else {
            setSortColumn(col);
            setSortDirection('desc');
          }
        }}
      />
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.COMPONENTS}`, module).add('Table', () => (
  <TablesRoot />
));
