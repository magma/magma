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

import CloseIcon from '@material-ui/icons/Close';
import DialogTitle from '@material-ui/core/DialogTitle';
import IconButton from '@material-ui/core/IconButton';
import React from 'react';
import Text from './Text';
import {colors} from '../default';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  closeButton: {
    color: colors.primary.white,
    padding: 0,
  },
}));

type Props = {
  label: string;
  onClose: () => void;
};

export default function CustomDialogTitle(props: Props) {
  const classes = useStyles(props);
  return (
    <DialogTitle>
      <Text variant="subtitle1">{props.label}</Text>
      <IconButton
        aria-label="close"
        className={classes.closeButton}
        onClick={props.onClose}>
        <CloseIcon />
      </IconButton>
    </DialogTitle>
  );
}
