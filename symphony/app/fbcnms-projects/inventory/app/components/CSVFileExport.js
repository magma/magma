/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {FiltersQuery} from './comparison_view/ComparisonViewTypes';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';
import type {WithStyles} from '@material-ui/core';

import Button from '@fbcnms/ui/components/design-system/Button';
import CircularProgress from '@material-ui/core/CircularProgress';
import React, {useState} from 'react';
import axios from 'axios';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {withStyles} from '@material-ui/core/styles';

const styles = {
  exportButton: {
    paddingLeft: '16px',
    paddingRight: '16px',
    marginLeft: '10px',
  },
  exportButtonContainer: {
    display: 'flex',
  },
};
const PATH_PREFIX = '/graph/export';

type Props = {
  exportPath: string,
  title: string,
  filters: ?FiltersQuery,
} & WithStyles<typeof styles> &
  WithAlert;

const CSVFileExport = (props: Props) => {
  const {classes, title, exportPath} = props;
  const [isDownloading, setIsDownloading] = useState(false);

  const onClick = async () => {
    const path = PATH_PREFIX + exportPath;
    const fileName = exportPath.replace('/', '').replace(/\//g, '_') + '.csv';
    setIsDownloading(true);
    try {
      await axios
        .get(path, {
          params: {
            filters: JSON.stringify(props.filters),
          },
          responseType: 'blob',
        })
        .then(response => {
          setIsDownloading(false);
          const url = window.URL.createObjectURL(new Blob([response.data]));
          const link = document.createElement('a');
          link.href = url;
          link.setAttribute('download', fileName);
          link.click();
        });
    } catch (error) {
      props.alert(error.response?.data?.error || error);
      setIsDownloading(false);
    }
  };

  return (
    <div className={classes.exportButtonContainer}>
      <Button className={classes.exportButton} onClick={onClick}>
        {isDownloading ? (
          <CircularProgress size={20} color={'inherit'} />
        ) : (
          title
        )}
      </Button>
    </div>
  );
};

export default withStyles(styles)(withAlert(CSVFileExport));
