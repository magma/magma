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
import Text from '@fbcnms/ui/components/design-system/Text';
import classNames from 'classnames';
import symphony from '@fbcnms/ui/theme/symphony';
import {
  WIDE_DIMENSION_HEIGHT_PX,
  WIDE_DIMENSION_WIDTH_PX,
} from '@fbcnms/ui/components/design-system/Experimental/FileUpload/FileUploadArea';
import {makeStyles} from '@material-ui/styles';

const PROGRESS_WIDTH = 70;

const useStyles = makeStyles(() => ({
  root: {
    height: WIDE_DIMENSION_HEIGHT_PX,
    width: WIDE_DIMENSION_WIDTH_PX,
    borderRadius: '4px',
    border: `1px solid ${symphony.palette.D100}`,
    position: 'relative',
    backgroundColor: symphony.palette.background,
  },
  name: {
    position: 'absolute',
    left: 8,
    right: 8,
    bottom: 8,
    zIndex: 2,
  },
  progressContainer: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
    padding: '40px',
  },
  progressBar: {
    width: `${PROGRESS_WIDTH}px`,
    height: '4px',
    borderRadius: '4px',
    backgroundColor: symphony.palette.D100,
  },
  actualProgress: {
    height: '4px',
    borderRadius: '4px 0px 0px 4px',
    backgroundColor: symphony.palette.B600,
    position: 'absolute',
    left: 0,
    top: 0,
  },
  progressBarContainer: {
    position: 'relative',
  },
}));

type Props = {
  fileName: string,
  progress: number,
  error?: ?string,
  className?: string,
};

const PendingFilePreview = ({
  fileName,
  progress,
  error: _error,
  className,
}: Props) => {
  const classes = useStyles();
  return (
    <div className={classNames(classes.root, className)}>
      <div className={classes.progressContainer}>
        <Text variant="body2" weight="medium">
          {progress}%
        </Text>
        <div className={classes.progressBarContainer}>
          <div className={classes.progressBar} />
          <div
            className={classes.actualProgress}
            style={{width: Math.ceil((progress / 100) * PROGRESS_WIDTH)}}
          />
        </div>
      </div>
      <Text
        className={classes.name}
        variant="caption"
        color="gray"
        weight="medium"
        useEllipsis={true}>
        {fileName}
      </Text>
    </div>
  );
};

export default PendingFilePreview;
