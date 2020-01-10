/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {AlertConfig} from './AlarmAPIType';

import * as React from 'react';
import Button from '@material-ui/core/Button';
import ClipboardLink from '@fbcnms/ui/components/ClipboardLink';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import {makeStyles} from '@material-ui/styles';
import type {RuleViewerProps} from './rules/RuleInterface';

const useStyles = makeStyles({
  paper: {
    minWidth: 360,
  },
  pre: {
    whiteSpace: 'pre-wrap',
    wordBreak: 'break-all',
  },
});

type Props<TRuleUnion> = {
  open: boolean,
  onClose: () => void,
  title: string,
  additionalContent?: React.Node,
  rule: TRuleUnion,
  showCopyButton?: boolean,
  showDeleteButton?: boolean,
  onDelete?: () => Promise<void>,
  RuleViewer: React.ComponentType<RuleViewerProps<TRuleUnion>>,
};

export default function AlertActionDialog<TRuleUnion>(
  props: Props<TRuleUnion>,
) {
  const {
    open,
    onClose,
    title,
    additionalContent,
    rule,
    showCopyButton,
    showDeleteButton,
    onDelete,
    RuleViewer,
  } = props;
  const classes = useStyles();

  return (
    <Dialog
      PaperProps={{classes: {root: classes.paper}}}
      open={open}
      onClose={onClose}>
      <DialogTitle>{title}</DialogTitle>
      <DialogContent>
        <RuleViewer rule={rule} />
        {additionalContent}
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose} color="primary">
          {showDeleteButton ? 'Cancel' : 'Close'}
        </Button>
        {showCopyButton && (
          <ClipboardLink>
            {({copyString}) => (
              <Button
                onClick={() => copyString(JSON.stringify(rule) || '')}
                color="primary"
                variant="contained">
                Copy
              </Button>
            )}
          </ClipboardLink>
        )}
        {showDeleteButton && (
          <Button onClick={onDelete} color="primary" variant="contained">
            Delete
          </Button>
        )}
      </DialogActions>
    </Dialog>
  );
}

AlertActionDialog.defaultProps = {
  RuleViewer: SimpleJsonViewer,
};

const useJsonStyles = makeStyles({
  pre: {
    whiteSpace: 'pre-wrap',
    wordBreak: 'break-all',
  },
});

function SimpleJsonViewer({rule}: {rule: AlertConfig}) {
  const classes = useJsonStyles();
  return <pre className={classes.pre}>{JSON.stringify(rule, null, 2)}</pre>;
}
