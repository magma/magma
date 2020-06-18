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
import Autosuggest from 'react-autosuggest';
import CancelIcon from '@material-ui/icons/Cancel';
import FormAction from './design-system/Form/FormAction';
import InputAffix from './design-system/Input/InputAffix';
import Text from './design-system/Text';
import TextInput from './design-system/Input/TextInput';
import Tooltip from '@material-ui/core/Tooltip';
import emptyFunction from '@fbcnms/util/emptyFunction';
import symphony from '../theme/symphony';
import useFollowElement from './useFollowElement';
import {blue05} from '../theme/colors';
import {makeStyles, useTheme} from '@material-ui/styles';
import {useFormElementContext} from '@fbcnms/ui/components/design-system/Form/FormElementContext';
import {useMemo, useRef, useState} from 'react';

const autoSuggestStyles = theme => ({
  container: {
    position: 'relative',
  },
  input: {
    width: '100%',
    padding: '0px',
    ...symphony.typography.subtitle1,
    color: theme.typography.subtitle1.color,
  },
  suggestionsContainer: {
    display: 'none',
  },
  suggestionsContainerOpen: {
    backgroundColor: theme.palette.common.white,
    borderRadius: '2px',
    display: 'block',
    ...symphony.typography.subtitle1,
    position: 'fixed',
    boxShadow: theme.shadows[2],
    zIndex: 5,
    transition: 'top 100ms ease-out 0s',
  },
  suggestionsList: {
    margin: 0,
    padding: 0,
    listStyleType: 'none',
  },
  suggestion: {
    color: symphony.palette.D900,
    cursor: 'pointer',
    padding: '10px 20px',
  },
  suggestionHighlighted: {
    backgroundColor: symphony.palette.background,
  },
});

const useStyles = makeStyles(theme => ({
  container: {
    width: '100%',
  },
  suggestionRoot: {
    display: 'flex',
    alignItems: 'center',
  },
  suggestionType: {
    marginLeft: theme.spacing(),
    lineHeight: '20px',
  },
  outlinedInput: {
    '&&': {
      backgroundColor: blue05,
      color: theme.palette.text.primary,
    },
  },
  smallSuggest: {
    paddingTop: '9px',
    paddingBottom: '9px',
    paddingLeft: '14px',
    paddingRight: '14px',
    height: '14px',
  },
  cancelIcon: {
    color: symphony.palette.D300,
  },
}));

type Props = {
  margin?: ?string,
  required: ?boolean,
  suggestions: Array<Suggestion>,
  onEntitySelected: Suggestion => void,
  onSuggestionsFetchRequested: (searchTerm: string) => void,
  onSuggestionsClearRequested?: () => void,
  placeholder?: ?string,
  value?: ?Suggestion,
  disabled?: boolean,
};

export type Suggestion = {
  entityId: string,
  entityType: string,
  name: string,
  type: string,
  render?: () => React.Node,
};

const Typeahead = (props: Props) => {
  const {
    onSuggestionsFetchRequested,
    onSuggestionsClearRequested,
    onEntitySelected,
    suggestions,
    placeholder,
    required,
    value,
    margin,
    disabled: propDisabled = false,
  } = props;
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedSuggestion, setSelectedSuggestion] = useState(value);
  const classes = useStyles();
  const theme = useTheme();

  const inputContainer = useRef(null);
  const dropdownContainer = useRef(null);
  useFollowElement(dropdownContainer, inputContainer);

  const {disabled: contextDisabled} = useFormElementContext();
  const disabled = useMemo(() => propDisabled || contextDisabled, [
    propDisabled,
    contextDisabled,
  ]);

  return (
    <div className={classes.container}>
      {selectedSuggestion && onSuggestionsClearRequested ? (
        <Tooltip
          arrow
          interactive
          placement="top"
          title={
            <Text variant="caption" color="light">
              {selectedSuggestion.entityId}
            </Text>
          }>
          <div>
            <TextInput
              type="string"
              required={!!required}
              variant="outlined"
              placeholder={placeholder ?? ''}
              fullWidth={true}
              disabled={disabled || selectedSuggestion != null}
              value={selectedSuggestion ? selectedSuggestion.name : ''}
              onChange={emptyFunction}
              suffix={
                searchTerm === '' ? (
                  <FormAction disabled={disabled}>
                    <InputAffix
                      onClick={() => {
                        setSearchTerm('');
                        setSelectedSuggestion(null);
                        onSuggestionsClearRequested &&
                          onSuggestionsClearRequested();
                      }}>
                      <CancelIcon className={classes.cancelIcon} />
                    </InputAffix>
                  </FormAction>
                ) : null
              }
            />
          </div>
        </Tooltip>
      ) : (
        <Autosuggest
          suggestions={suggestions}
          getSuggestionValue={(suggestion: Suggestion) => suggestion.name}
          onSuggestionsFetchRequested={({value}) => {
            onSuggestionsFetchRequested(value);
          }}
          renderSuggestionsContainer={({containerProps, children}) => (
            <div
              {...containerProps}
              ref={refInput => {
                dropdownContainer.current = refInput;
                containerProps.ref && containerProps.ref(refInput);
              }}>
              {children}
            </div>
          )}
          onSuggestionsClearRequested={emptyFunction}
          renderSuggestion={(suggestion: Suggestion) => (
            <div className={classes.suggestionRoot}>
              {suggestion.render != null ? (
                suggestion.render()
              ) : (
                <div>
                  <Text>{suggestion.name}</Text>
                  <Text
                    className={classes.suggestionType}
                    variant="caption"
                    color="gray">
                    {suggestion.type}
                  </Text>
                </div>
              )}
            </div>
          )}
          onSuggestionSelected={(e, data) => {
            const suggestion: Suggestion = data.suggestion;
            setSearchTerm('');
            setSelectedSuggestion(suggestion);
            onEntitySelected(suggestion);
          }}
          theme={autoSuggestStyles(theme)}
          renderInputComponent={inputProps => (
            <div ref={inputContainer}>
              <TextInput
                type="string"
                placeholder={placeholder ?? 'Search...'}
                {...inputProps}
              />
            </div>
          )}
          inputProps={{
            style: {},
            required: !!required,
            value: searchTerm,
            margin,
            onChange: (_e, {newValue}) => setSearchTerm(newValue),
          }}
          highlightFirstSuggestion={true}
        />
      )}
    </div>
  );
};

export default Typeahead;
