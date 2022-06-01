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
import IconButton from '@material-ui/core/IconButton';
import Tooltip, {TooltipProps} from '@material-ui/core/Tooltip';

import copy from 'copy-to-clipboard';
import {useState} from 'react';

type Props = {
  children: (props: {
    copyString: (content: string) => void;
  }) =>
    | React.ReactElement<typeof Button>
    | React.ReactElement<typeof IconButton>;
  // We set the title appropriately later
  title?: React.ReactNode;
} & Omit<TooltipProps, 'title'>;

const COPIED_MESSAGE = 'Copied to clipboard!';

/* Wrap a button with this component to copy a string to clipboard when
 * the button is clicked. After the button is clicked, a tooltip will
 * pop up saying the copying was successful. Tooltip custom props can be
 * passed into this component directly.
 */
export default function ClipboardLink({title, ...props}: Props) {
  if (title != null) {
    return <ClipboardLinkWithTitle {...props} title={title} />;
  }
  return <ClipboardLinkNoTitle {...props} />;
}

/* Since the logic and states are diffferent depending on whether a title for
 * the tooltip is passed in, we have 2 different components below for each
 * scenario.
 */

// If they pass in a title, we need to change that title briefly to
// COPIED_MESSAGE whenever the content is copied.
function ClipboardLinkWithTitle(
  props: Props & {title: NonNullable<React.ReactNode>},
) {
  const [currentTitle, setCurrentTitle] = useState<
    NonNullable<React.ReactNode>
  >(props.title);
  return (
    <Tooltip
      {...props}
      title={currentTitle}
      onClose={() => setCurrentTitle(props.title)}>
      {props.children({
        copyString: content => {
          copy(content);
          setCurrentTitle(COPIED_MESSAGE);
        },
      })}
    </Tooltip>
  );
}

// If they don't pass in a title, there should be no COPIED_MESSAGE tooltip
// shown until the content is copied.
function ClipboardLinkNoTitle(props: Props) {
  const [showTooltip, setShowTooltip] = useState(false);
  return (
    <Tooltip
      {...props}
      title={COPIED_MESSAGE}
      open={showTooltip}
      onClose={() => setShowTooltip(false)}>
      {props.children({
        copyString: content => {
          copy(content);
          setShowTooltip(true);
        },
      })}
    </Tooltip>
  );
}
