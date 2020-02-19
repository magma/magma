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
import Table from './design-system/Table/Table';
import fbt from 'fbt';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  root: {
    margin: '4px 8px',
    alignItems: 'center',
  },
}));
type errLine = {
  key?: string,
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

  const errToMessage = (e: errLine): string => {
    return e.error == '' ? e.message : e.error + ': ' + e.message;
  };

  return (
    <div className={classes.root}>
      <Table
        data={allSkippedRows}
        columns={[
          {
            key: '0',
            title: fbt('Line', 'title of the number of the line'),
            render: row => row.line,
          },
          {
            key: '1',
            title: fbt('Issue', ' title of the error description'),
            render: row => errToMessage(row),
          },
        ]}
      />
    </div>
  );
};

export default UploadErrorsList;
