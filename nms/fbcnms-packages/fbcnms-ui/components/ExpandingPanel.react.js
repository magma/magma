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
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import Typography from '@material-ui/core/Typography';
import classNames from 'classnames';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  expansionPanel: {
    padding: '24px',
    borderRadius: '4px',
    '&:before': {
      content: 'none',
    },
    boxShadow: '0px 1px 4px 0px rgba(0,0,0,0.17)',
  },
  expansionPanelSummary: {
    '&&': {
      padding: '0px',
      minHeight: 'auto',
    },
  },
  expandIcon: {
    padding: '0px',
  },
  summaryContent: {
    '&&': {
      margin: 0,
    },
  },
  panelTitle: {
    fontSize: '20px',
    color: theme.palette.blueGrayDark,
    lineHeight: '28px',
    fontWeight: 500,
  },
  panelDetails: {
    padding: 0,
    display: 'flex',
    flexDirection: 'column',
  },
}));

type Props = {
  title: string,
  children: React.Node,
  className?: string,
};

const ExpandingPanel = ({className, children, title}: Props) => {
  const classes = useStyles();
  return (
    <ExpansionPanel
      className={classNames(className, classes.expansionPanel)}
      defaultExpanded={true}>
      <ExpansionPanelSummary
        className={classes.expansionPanelSummary}
        classes={{
          expandIcon: classes.expandIcon,
          content: classes.summaryContent,
        }}
        expandIcon={<ExpandMoreIcon />}>
        <Typography className={classes.panelTitle}>{title}</Typography>
      </ExpansionPanelSummary>
      <ExpansionPanelDetails className={classes.panelDetails}>
        {children}
      </ExpansionPanelDetails>
    </ExpansionPanel>
  );
};

export default ExpandingPanel;
