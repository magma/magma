/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * A material-ui Typography component which also handles translation and key
 * extraction.
 * @flow
 * @format
 */

import * as React from 'react';
import Typography from '@material-ui/core/Typography';
import {Trans} from 'react-i18next';

type Props = {
  children: any,
} & TypographyProps &
  TransProps;

type TypographyProps = {};

type TransProps = {
  i18nKey?: string,
  count?: number,
  values?: Object,
};

export default function Text(props: Props) {
  const {i18nKey, count, values, children, ...typographyProps} = props;
  return (
    <Typography {...typographyProps}>
      <Trans i18nKey={i18nKey} count={count} values={values}>
        {children}
      </Trans>
    </Typography>
  );
}
