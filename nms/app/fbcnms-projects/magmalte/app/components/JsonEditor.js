/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import Button from '@material-ui/core/Button';
import FormLabel from '@material-ui/core/FormLabel';
import Grid from '@material-ui/core/Grid';
import React from 'react';
import ReactJson from 'react-json-view';
import SettingsIcon from '@material-ui/icons/Settings';
import Text from '@fbcnms/ui/components/design-system/Text';

import {colors, typography} from '../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useEffect, useState} from 'react';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(3),
    flexGrow: 1,
  },
  itemTitle: {
    color: colors.primary.comet,
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
  },
  itemValue: {
    color: colors.primary.brightGray,
  },
  appBarBtn: {
    color: colors.primary.white,
    background: colors.primary.comet,
    fontFamily: typography.button.fontFamily,
    fontWeight: typography.button.fontWeight,
    fontSize: typography.button.fontSize,
    lineHeight: typography.button.lineHeight,
    letterSpacing: typography.button.letterSpacing,

    '&:hover': {
      background: colors.primary.mirage,
    },
  },
  jsonTextarea: {
    fontFamily: 'monospace',
    height: '95%',
    border: 'none',
    margin: theme.spacing(2),
  },
  configBody: {
    display: 'flex',
    flexFlow: 'column',
    flexGrow: 1,
    overflowX: 'hidden',
  },
  appBarBtnSecondary: {
    color: colors.primary.brightGray,
  },
}));

type Props<T> = {
  content: T,
  error: string,
  onSave: T => Promise<void>,
};

export default function JsonEditor<T>(props: Props<T>) {
  const classes = useStyles();
  const [error, setError] = useState<string>(props.error);
  const [content, setContent] = useState<T>(props.content);

  useEffect(() => {
    setError(props.error);
  }, [props.error]);

  const handleChange = data => {
    setContent(data.updated_src);
  };

  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={3}>
        <Grid container item xs={12}>
          <Grid item xs={6}>
            <Text>
              <SettingsIcon /> JSON Config
            </Text>
          </Grid>
          <Grid container item justify="flex-end" xs={6}>
            <Grid item>
              <Button
                className={classes.appBarBtnSecondary}
                onClick={() => {
                  setContent(props.content);
                  setError('');
                }}>
                Cancel
              </Button>
            </Grid>
            <Grid item>
              <Button
                className={classes.appBarBtn}
                onClick={() => {
                  try {
                    props.onSave(content);
                  } catch (e) {
                    setError(e.message);
                  }
                }}>
                Save
              </Button>
            </Grid>
          </Grid>
        </Grid>

        <Grid
          container
          className={classes.configBody}
          alignItems="stretch"
          item
          xs={12}>
          {error !== '' && <FormLabel error>{error}</FormLabel>}
          <ReactJson
            src={content}
            enableClipboard={false}
            displayDataTypes={false}
            onAdd={handleChange}
            onEdit={handleChange}
            onDelete={handleChange}
          />
        </Grid>
      </Grid>
    </div>
  );
}
