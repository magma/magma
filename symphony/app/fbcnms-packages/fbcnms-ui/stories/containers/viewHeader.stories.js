/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import type {ActionButtonProps} from '../../components/design-system/View/ViewHeader';

import DeleteIcon from '../../components/design-system/Icons/Actions/DeleteIcon';
import ListAltIcon from '@material-ui/icons/ListAlt';
import MapIcon from '@material-ui/icons/Map';
import React from 'react';
import ViewHeader from '../../components/design-system/View/ViewHeader';
import {STORY_CATEGORIES} from '../storybookUtils';
import {makeStyles} from '@material-ui/styles';
import {storiesOf} from '@storybook/react';
import {useState} from 'react';

const useStyles = makeStyles(() => ({
  root: {
    width: '100%',
  },
  card: {
    marginBottom: '16px',
  },
}));

const ViewHeaderRoot = () => {
  const classes = useStyles();
  const [selectedButton, setSelectedButton] = useState('1');

  const actionButtons: Array<ActionButtonProps> = [
    {
      title: <DeleteIcon color="error" />,
      action: () => alert('Doing DELETE'),
      variant: 'text',
    },
    {
      title: 'Action B',
      action: () => alert('Doing action B'),
      variant: 'text',
      skin: 'gray',
    },
    {
      title: 'Action A',
      action: () => alert('Doing action A'),
    },
  ];

  const subTitle =
    'The Company is a secret group of multinational corporate alliances known only by those who work for them or oppose them. Its influence and power over individuals stretches to the White House, controlling every decision the country makes.';

  const viewOptions = {
    onItemClicked: setSelectedButton,
    selectedButtonId: selectedButton,
    buttons: [
      {
        item: <ListAltIcon />,
        id: '1',
      },
      {
        item: <MapIcon />,
        id: '2',
      },
    ],
  };
  return (
    <div className={classes.root}>
      <ViewHeader
        title="The Company"
        subtitle={subTitle}
        actionButtons={actionButtons}
        viewOptions={viewOptions}
      />
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.CONTAINERS}`, module).add('ViewHeader', () => (
  <ViewHeaderRoot />
));
