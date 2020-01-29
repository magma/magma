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
import {FileUploadStatuses} from '../../components/design-system/Experimental/FileUploadStatus';
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
              id: '0',
              name: 'Blue_bird.jpg',
              status: FileUploadStatuses.DONE,
            },
            {
              id: '1',
              name: 'Blue_bird2.jpg',
              status: FileUploadStatuses.ERROR,
              errorMessage: 'We had a problem uploading this file',
            },
            {
              id: '2',
              name: 'Blue_bird_singing_in_the_dead_of_night.jpg',
              status: FileUploadStatuses.ERROR,
              errorMessage:
                'We had a problem uploading this file, long error message',
            },
            {
              id: '3',
              name: 'Blue_bird.jpg',
              status: FileUploadStatuses.UPLOADING,
            },
            {
              id: '4',
              name: 'Blue_bird.jpg',
              status: FileUploadStatuses.UPLOADING,
            },
            {
              id: '5',
              name: 'Blue_bird.jpg',
              status: FileUploadStatuses.UPLOADING,
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
