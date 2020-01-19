/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import Button from '../../components/design-system/Button';
import FilesUploadSnackbar from '../../components/design-system/Experimental/FilesUploadSnackbar';
import React, {useState} from 'react';
import {STORY_CATEGORIES} from '../storybookUtils';
import {makeStyles} from '@material-ui/styles';
import {storiesOf} from '@storybook/react';

const useStyles = makeStyles(_theme => ({
  root: {
    width: '100%',
  },
}));

const FilesUploadSnackbarRoot = () => {
  const classes = useStyles();
  const [isShown, setIsShown] = useState(true);
  return (
    <div className={classes.root}>
      <Button onClick={() => setIsShown(true)}>Upload</Button>
      {isShown && (
        <FilesUploadSnackbar
          onClose={() => setIsShown(false)}
          files={[
            {
              name: 'Blue_bird.jpg',
              status: 'done',
            },
            {
              name: 'Blue_bird2.jpg',
              status: 'error',
              errorMessage: 'We had a problem uploading this file',
            },
            {
              name: 'Blue_bird_singing_in_the_dead_of_night.jpg',
              status: 'uploading',
              status: 'error',
              errorMessage:
                'We had a problem uploading this file, long error message',
            },
            {
              name: 'Blue_bird4.jpg',
              status: 'uploading',
            },
            {
              name: 'Blue_bird5.jpg',
              status: 'uploading',
            },
            {
              name: 'Blue_bird6.jpg',
              status: 'uploading',
            },
          ]}
        />
      )}
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.EXPERIMENTAL}`, module).add(
  'FilesUploadSnackbar',
  () => <FilesUploadSnackbarRoot />,
);
