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
import Popover from '@material-ui/core/Popover';
import {makeStyles} from '@material-ui/styles';

type Props = {
  className?: string;
  content: React.ReactElement<any> | string;
  children: React.ReactNode;
  onOpen?: () => void;
  onClose?: () => void;
  open?: boolean;
  contentClickTriggerClose?: boolean;
};

const useClasses = makeStyles(() => ({
  menuPaper: {
    '&.MuiPopover-paper': {
      borderRadius: 0,
    },
    outline: 'none',
    overflowX: 'visible',
    overflowY: 'visible',
    '&:after': {
      content: '""',
      display: 'block',
      height: 0,
      width: 0,
      position: 'absolute',
      left: '-14px',
      bottom: '30px',
      zIndex: 5,
      borderTop: '11px solid transparent',
      borderRight: '14px solid #fff',
      borderBottom: '11px solid transparent',
    },
  },
  popover: {
    '& $menuPaper': {
      boxShadow: '0px 8px 16px 0px rgba(0, 0, 0, 0.3)',
    },
  },
}));

export default function Popout(props: Props) {
  const classes = useClasses();
  const relativeRef = React.useRef<HTMLDivElement>(null);
  const [open, togglePopout] = React.useState(false);
  const {onClose} = props;
  const handleClose = React.useCallback(
    () => (onClose ? onClose() : togglePopout(false)),
    [onClose, togglePopout],
  );

  const relativeRefPosition = relativeRef.current
    ? relativeRef.current.getBoundingClientRect()
    : null;

  return (
    <>
      <div
        className={props.className}
        ref={relativeRef}
        onClick={() => {
          props.onOpen ? props.onOpen() : togglePopout(true);
        }}>
        {props.children}
      </div>
      <Popover
        className={classes.popover}
        anchorReference="anchorPosition"
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'left',
        }}
        anchorPosition={{
          top: relativeRefPosition ? relativeRefPosition.bottom + 15 : 0,
          // Hardcode the sidebar width otherwise the popover is misplaced
          // if opened during the sidebar animation
          left: 208,
        }}
        transformOrigin={{
          vertical: 'bottom',
          horizontal: 'left',
        }}
        PaperProps={{className: classes.menuPaper}}
        id="navigation-menu"
        open={props.open !== undefined ? props.open : open}
        onClose={handleClose}
        onClick={props.contentClickTriggerClose ? handleClose : undefined}
        onMouseOver={event => event.stopPropagation()}>
        {props.content}
      </Popover>
    </>
  );
}
