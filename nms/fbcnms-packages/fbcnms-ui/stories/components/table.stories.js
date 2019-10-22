/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import React from 'react';
import Table from '../../components/design-system/Table/Table.react';
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
}));

const TablesRoot = () => {
  const classes = useStyles();
  return (
    <div className={classes.root}>
      <Table
        data={DATA}
        columns={[
          {title: 'First Name', render: row => row.firstName},
          {title: 'Last Name', render: row => row.lastName},
          {title: 'Birth Date', render: row => row.birthDate},
          {title: 'City', render: row => row.city},
        ]}
      />
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.COMPONENTS}`, module).add('Table', () => (
  <TablesRoot />
));
