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
 * @flow strict-local
 * @format
 */

import * as React from 'react';
import Button from '@material-ui/core/Button';
import Menu from '@material-ui/core/Menu';

type Props = {
  id: string,
  children: React.Node,
  buttonContent: React.Node,
  className?: string,
};

export default function TopBarAnchoredMenu(props: Props) {
  const [anchorEl, setAnchorEl] = React.useState<?HTMLElement>(null);
  return (
    <>
      <Button
        aria-owns={anchorEl ? props.id : null}
        aria-haspopup="true"
        onClick={e => setAnchorEl(e.currentTarget)}
        className={props.className}
        color="inherit">
        {props.buttonContent}
      </Button>
      <Menu
        id={props.id}
        anchorEl={anchorEl}
        anchorOrigin={{vertical: 'top', horizontal: 'right'}}
        transformOrigin={{vertical: 'top', horizontal: 'right'}}
        open={!!anchorEl}
        onClose={() => setAnchorEl(null)}>
        {props.children}
      </Menu>
    </>
  );
}
