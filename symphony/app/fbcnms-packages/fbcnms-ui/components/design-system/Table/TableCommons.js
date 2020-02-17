/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import {makeStyles} from '@material-ui/styles';

const styles = {
  cell: {
    padding: '4px 8px 4px 16px',
    minHeight: '48px',
    height: '48px',
  },
};

export const useTableCommonStyles = makeStyles<{}, typeof styles>(() => styles);
