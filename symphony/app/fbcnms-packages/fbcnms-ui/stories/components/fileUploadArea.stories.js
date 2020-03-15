/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import FileUploadArea from '../../components/design-system/Experimental/FileUpload/FileUploadArea';
import React from 'react';
import {STORY_CATEGORIES} from '../storybookUtils';
import {makeStyles} from '@material-ui/styles';
import {storiesOf} from '@storybook/react';

const useStyles = makeStyles(_theme => ({
  root: {
    width: '100%',
  },
  card: {
    marginBottom: '16px',
  },
}));

const FileUploadAreaRoot = () => {
  const classes = useStyles();

  return (
    <div className={classes.root}>
      <FileUploadArea onFileChanged={files => alert(files[0].name)} />
    </div>
  );
};

storiesOf(`${STORY_CATEGORIES.EXPERIMENTAL}`, module).add(
  'FileUploadArea',
  () => <FileUploadAreaRoot />,
);
