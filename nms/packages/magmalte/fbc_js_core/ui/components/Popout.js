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
 *
 * @flow
 * @format
 */

import * as React from 'react';
import Popover from '@material-ui/core/Popover';
import {makeStyles} from '@material-ui/styles';

type Props = {
  content: React.Element<any> | string,
  children: React.Element<any> | string,
  onOpen?: () => void,
  onClose?: () => void,
  open?: boolean,
  contentClickTriggerClose?: boolean,
};

const useClasses = makeStyles(theme => ({
  root: {
    position: 'relative',
    display: 'inline-block',
  },
  menuPaper: {
    outline: 'none',
    overflowX: 'visible',
    overflowY: 'visible',
    position: 'absolute',
    '&:before, &:after': {
      content: '""',
      display: 'block',
      height: 0,
      left: '11px',
      position: 'absolute',
      width: 0,
    },
    '&:before': {
      borderLeft: '5px solid transparent',
      borderRight: '5px solid transparent',
      borderTop: `6px solid ${theme.palette.grey[100]}`,
      marginLeft: '-3px',
      bottom: '-6px',
      zIndex: 4,
    },
    '&:after': {
      borderLeft: '5px solid transparent',
      borderRight: '5px solid transparent',
      borderTop: '7px solid #fff',
      marginLeft: '-3px',
      bottom: '-5px',
      zIndex: 5,
    },
  },
  popover: {
    '& $menuPaper': {
      boxShadow: '0px 0px 4px 0px rgba(0, 0, 0, 0.15)',
    },
  },
  buttonContainer: {
    display: 'inline-block',
    position: 'relative',
  },
  buttonRelative: {
    display: 'inline-block',
    position: 'absolute',
    left: 0,
    top: '-14px',
    width: '100%',
  },
}));

export default function Popout(props: Props) {
  const {content, children, onOpen, onClose} = props;
  const classes = useClasses();
  const relativeRef = React.useRef();
  const [open, togglePopout] = React.useState(false);
  const handleClose = React.useCallback(
    () => (onClose ? onClose() : togglePopout(false)),
    [onClose, togglePopout],
  );

  const relativeRefPosition = relativeRef.current
    ? relativeRef.current.getBoundingClientRect()
    : null;

  return (
    <div className={classes.root}>
      <div
        className={classes.buttonContainer}
        onClick={() => {
          onOpen ? onOpen() : togglePopout(true);
        }}>
        {children}
      </div>
      {/* $FlowFixMe - Flow ref type definition is not up to date */}
      <div className={classes.buttonRelative} ref={relativeRef} />
      <Popover
        className={classes.popover}
        anchorReference="anchorPosition"
        anchorOrigin={{
          vertical: 'top',
          horizontal: 'left',
        }}
        anchorPosition={{
          top: relativeRefPosition?.y ?? 0,
          left: relativeRefPosition
            ? relativeRefPosition.x + relativeRefPosition.width / 2 - 14
            : 0,
        }}
        transformOrigin={{
          vertical: 'bottom',
          horizontal: 'left',
        }}
        PaperProps={{className: classes.menuPaper}}
        id="navigation-menu"
        open={props.open !== undefined ? props.open : open}
        onClose={handleClose}
        onClick={props.contentClickTriggerClose && handleClose}>
        {content}
      </Popover>
    </div>
  );
}
