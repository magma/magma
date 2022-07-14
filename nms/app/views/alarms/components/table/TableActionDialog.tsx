/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
import * as React from 'react';
import Button from '@material-ui/core/Button';
import ClipboardLink from '../ClipboardLink';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles({
  paper: {
    minWidth: 360,
  },
  pre: {
    whiteSpace: 'pre-wrap',
    wordBreak: 'break-all',
  },
});

type Props<TRow> = {
  open: boolean;
  onClose: () => void;
  title: React.ReactNode;
  additionalContent?: React.ReactNode;
  row: TRow;
  showCopyButton?: boolean;
  showDeleteButton?: boolean;
  onDelete?: () => Promise<void>;
  RowViewer: React.ComponentType<{row?: TRow}>;
};

export default function TableActionDialog<TRow>(props: Props<TRow>) {
  const {
    open,
    onClose,
    title,
    additionalContent,
    row,
    showCopyButton,
    showDeleteButton,
    onDelete,
    RowViewer,
  } = props;
  const classes = useStyles();
  if (!row) {
    return null;
  }
  return (
    <Dialog
      PaperProps={{classes: {root: classes.paper}}}
      open={open}
      onClose={onClose}>
      <DialogTitle>{title}</DialogTitle>
      <DialogContent>
        <RowViewer row={row} />
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
                onClick={() => copyString(JSON.stringify(row) || '')}
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

TableActionDialog.defaultProps = {
  RowViewer: SimpleJsonViewer,
};

const useJsonStyles = makeStyles(() => ({
  pre: {
    whiteSpace: 'pre-wrap',
    wordBreak: 'break-all',
  },
}));

export function SimpleJsonViewer<TRow>({row}: {row: TRow}) {
  const classes = useJsonStyles();
  return <pre className={classes.pre}>{JSON.stringify(row, null, 2)}</pre>;
}
