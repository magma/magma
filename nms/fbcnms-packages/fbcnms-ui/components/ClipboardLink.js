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
import Button from '@material-ui/core/Button';
import IconButton from '@material-ui/core/IconButton';
import Tooltip from '@material-ui/core/Tooltip';

import copy from 'copy-to-clipboard';
import {useState} from 'react';

type Props = {
  ...React.ElementConfig<typeof Tooltip>,
  children: (props: {copyString: (content: string) => void}) =>
    | React.Element<typeof Button>
    | React.Element<typeof IconButton>,
};

const COPIED_MESSAGE = 'Copied to clipboard!';

/* Wrap a button with this component to copy a string to clipboard when
 * the button is clicked. After the button is clicked, a tooltip will
 * pop up saying the copying was successful. Tooltip custom props can be
 * passed into this component directly.
 */
export default function ClipboardLink(props: Props) {
  const [title, setTitle] = useState(props.title);
  const [showTooltip, setShowTooltip] = useState(false);

  if (props.title != null) {
    return (
      <Tooltip {...props} title={title} onClose={() => setTitle(props.title)}>
        {props.children({
          copyString: content => {
            copy(content);
            setTitle(COPIED_MESSAGE);
          },
        })}
      </Tooltip>
    );
  }
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
