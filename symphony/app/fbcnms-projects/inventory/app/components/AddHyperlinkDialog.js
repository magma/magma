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
import Button from '@fbcnms/ui/components/design-system/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormContext, {FormContextProvider} from '../common/FormContext';
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import Strings from '../common/CommonStrings';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import fbt from 'fbt';
import isUrl from 'is-url';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useState} from 'react';

const useStyles = makeStyles(() => ({
  field: {
    '&:not(:last-child)': {
      marginBottom: '8px',
    },
  },
}));

type Props = {
  isOpened: boolean,
  onAdd: (usr: string, displayName: ?string) => void,
  onClose: () => void,
  targetCategory?: ?string,
};

const AddHyperlinkDialog = (props: Props) => {
  const {isOpened, targetCategory} = props;
  const [url, setUrl] = useState('');
  const [displayName, setDisplayName] = useState('');

  const classes = useStyles();

  const onClose = useCallback(() => {
    if (props.onClose) {
      props.onClose();
    }
  }, [props]);
  const onSave = useCallback(() => {
    if (props.onAdd) {
      props.onAdd(url, displayName);
    }
    onClose();
  }, [onClose, props, url, displayName]);

  return (
    <FormContextProvider>
      <FormContext.Consumer>
        {form => {
          const urlValidationError = form.alerts.error.check({
            fieldId: 'url',
            fieldDisplayName: Strings.common.fields.url.label,
            value: url,
            required: true,
            checkCallback: value => {
              const rightFormat = value && isUrl(value);
              return rightFormat
                ? ''
                : `${fbt(
                    'URL must be in a valid URL format',
                    'URL format validation error message',
                  )}`;
            },
          });
          return (
            <Dialog
              maxWidth="sm"
              open={isOpened}
              onClose={props.onClose}
              fullWidth={true}>
              <DialogTitle>
                {!!targetCategory ? (
                  <Text>
                    <fbt desc="Adding url under given category dialog title">
                      Add a new URL under
                    </fbt>
                    {' ' + targetCategory}
                  </Text>
                ) : (
                  <Text>
                    <fbt desc="Adding url to some entity details">
                      Add a new URL
                    </fbt>
                  </Text>
                )}
              </DialogTitle>
              <DialogContent>
                <FormField
                  className={classes.field}
                  label={Strings.common.fields.url.label}
                  required={true}
                  hasError={!!urlValidationError}
                  errorText={urlValidationError}>
                  <TextInput
                    type="url"
                    placeholder={Strings.common.fields.url.placeholder}
                    value={url}
                    onChange={e => setUrl(e.target.value)}
                  />
                </FormField>
                <FormField
                  label={`${fbt(
                    'Display Name',
                    'Label for Display Name field',
                  )}`}
                  className={classes.field}>
                  <TextInput
                    value={displayName}
                    onChange={e => setDisplayName(e.target.value)}
                  />
                </FormField>
              </DialogContent>
              <DialogActions>
                <Button onClick={onClose} skin="regular">
                  {Strings.common.cancelButton}
                </Button>
                <Button onClick={onSave} disabled={form.alerts.error.detected}>
                  {Strings.documents.addLinkButton}
                </Button>
              </DialogActions>
            </Dialog>
          );
        }}
      </FormContext.Consumer>
    </FormContextProvider>
  );
};

export default AddHyperlinkDialog;
