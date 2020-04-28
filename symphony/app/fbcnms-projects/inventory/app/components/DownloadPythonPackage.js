/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import Button from '@fbcnms/ui/components/design-system/Button';
import DateTimeFormat from '../common/DateTimeFormat.js';
import InventoryQueryRenderer from './InventoryQueryRenderer';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import axios from 'axios';
import {DocumentAPIUrls} from '../common/DocumentAPI';
import {graphql} from 'relay-runtime';
import {makeStyles} from '@material-ui/styles';

type PythonPackage = {
  version: string,
  whlFileKey: string,
};

const useStyles = makeStyles(_ => ({
  root: {
    marginLeft: '20px',
    maxWidth: '300px',
  },
}));

const DownloadPythonPackageQuery = graphql`
  query DownloadPythonPackageQuery {
    pythonPackages {
      version
      whlFileKey
      uploadTime
    }
  }
`;

const handleDownload = (key: string, fileName: string) => {
  const method = 'GET';
  const url = DocumentAPIUrls.download_url(key, fileName);

  axios
    .request({
      url,
      method,
      responseType: 'blob',
      headers: {
        'Is-Global': 'True',
      },
    })
    .then(({data}) => {
      const downloadUrl = window.URL.createObjectURL(new Blob([data]));
      const link = document.createElement('a');
      link.href = downloadUrl;
      link.setAttribute('download', fileName);
      if (document.body != null) {
        document.body.appendChild(link);
      }
      link.click();
      link.remove();
    });
};

const handleWhlDownload = (pythonPackage: PythonPackage) => {
  handleDownload(
    pythonPackage.whlFileKey,
    `pyinventory-${pythonPackage.version}-py3-none-any.whl`,
  );
};

const DownloadPythonPackage = () => {
  const classes = useStyles();
  return (
    <InventoryQueryRenderer
      query={DownloadPythonPackageQuery}
      variables={{}}
      render={props => {
        const pythonPackages = props.pythonPackages;

        if (pythonPackages == null) {
          return <div />;
        }

        return (
          <Table className={classes.root}>
            <TableHead>
              <TableRow>
                <TableCell>Version</TableCell>
                <TableCell>Release Date</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {pythonPackages.map(p => (
                <TableRow>
                  <TableCell component="th" scope="row">
                    <Button variant="text" onClick={() => handleWhlDownload(p)}>
                      {p.version}
                    </Button>
                  </TableCell>
                  <TableCell>{DateTimeFormat.dateOnly(p.uploadTime)}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        );
      }}
    />
  );
};

export default DownloadPythonPackage;
