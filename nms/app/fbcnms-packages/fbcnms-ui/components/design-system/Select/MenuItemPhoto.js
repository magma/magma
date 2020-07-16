/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import React from 'react';
import {makeStyles} from '@material-ui/styles';

const PHOTO_SIZE = 36;

const useStyles = makeStyles(() => ({
  img: {
    width: `${PHOTO_SIZE}px`,
    padding: '2px 0px',
    height: `${PHOTO_SIZE}px`,
    borderRadius: '100%',
  },
}));

type Props = $ReadOnly<{|
  src: string,
|}>;

const MenuItemPhoto = ({src}: Props) => {
  const classes = useStyles();
  return <img className={classes.img} src={src} />;
};

export default MenuItemPhoto;
