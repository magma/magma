/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import * as React from 'react';
import classNames from 'classnames';
import symphony from '../theme/symphony';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles({
  content: {
    paddingTop: '4px',
    paddingBottom: '4px',
  },
  root: {
    borderWidth: '1px',
    borderRadius: '4px',
    borderStyle: 'solid',
    margin: '4px 8px',
    padding: '4px 10px',
    alignItems: 'center',
  },
  border: {
    borderColor: symphony.palette.Y600,
  },
});
type errLine = {
  line: number,
  error: string,
  message: string,
};

type Props = {
  errors: Array<errLine>,
  skipped?: Array<errLine>,
};

const UploadErrorsList = (props: Props) => {
  const classes = useStyles();
  const {errors, skipped} = props;

  const allSkippedRows = errors.concat(skipped ?? []).sort((a, b) => {
    return a.line - b.line;
  });
  return (
    <div className={classNames(classes.root, classes.border)}>
      {allSkippedRows.map(e => (
        <div className={classes.content} key={e.line}>
          {e.line}
          {': '}
          {e.error}
          {': '}
          {e.message}
        </div>
      ))}
    </div>
  );
};

export default UploadErrorsList;
