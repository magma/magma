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

import BasePopoverTrigger from '../ContexualLayer/BasePopoverTrigger';
import React, {
  useEffect,
  useLayoutEffect,
  useMemo,
  useRef,
  useState,
} from 'react';
import SelectMenu from '../Select/SelectMenu';
import TextInput from '../Input/TextInput';
import TokenizerBasicPostFetchDecorator from './TokenizerBasicPostFetchDecorator';
import TokensList from './TokensList';
import classNames from 'classnames';
import {makeStyles} from '@material-ui/styles';
import {useFormElementContext} from '../Form/FormElementContext';

const useStyles = makeStyles(() => ({
  root: {
    display: 'flex',
    alignItems: 'center',
  },
  input: {
    width: '100%',
  },
  inputContainer: {
    paddingLeft: '0px',
    flexWrap: 'wrap',
  },
}));

export type TokenizerEntryType<T> = $ReadOnly<{|
  key: string,
  label: string,
  ...T,
|}>;

export type Entries<TEntry> = $ReadOnlyArray<TokenizerEntryType<TEntry>>;

export type NetworkDataSource<TEntry> = $ReadOnly<{|
  fetchNetwork: (query: string) => Promise<Entries<TEntry>>,
  postFetchDecorator?: (
    response: Entries<TEntry>,
    queryString: string,
    currentTokens: Entries<TEntry>,
  ) => Entries<TEntry>,
|}>;

export type TokenizerDisplayProps = $ReadOnly<{|
  disabled?: boolean,
  inputClassName?: string,
  menuClassName?: string,
|}>;

type Props<TEntry> = $ReadOnly<{|
  ...TokenizerDisplayProps,
  tokens: Entries<TEntry>,
  dataSource: NetworkDataSource<TEntry>,
  onTokensChange: (Entries<TEntry>) => void,
  queryString: string,
  onQueryStringChange: string => void,
|}>;

const Tokenizer = <TEntry>(props: Props<TEntry>) => {
  const {
    queryString,
    onQueryStringChange,
    dataSource: {
      fetchNetwork,
      postFetchDecorator = TokenizerBasicPostFetchDecorator,
    },
    onTokensChange,
    tokens,
    disabled: propDisabled = false,
    inputClassName,
    menuClassName,
  } = props;
  const classes = useStyles();
  const [searchEntries, setSearchEntries] = useState<Entries<TEntry>>([]);
  const inputRef = useRef(null);
  const popoverTriggerRef = useRef(null);

  useEffect(() => {
    fetchNetwork(queryString)
      .then(response => postFetchDecorator(response, queryString, tokens))
      .then(response => setSearchEntries(response));
  }, [queryString, fetchNetwork, postFetchDecorator, tokens]);

  const selectMenuOptions = useMemo(
    () => searchEntries.slice().map(entry => ({...entry, value: entry})),
    [searchEntries],
  );

  useLayoutEffect(
    () => popoverTriggerRef && popoverTriggerRef.current?.reposition(),
    [tokens],
  );

  const {disabled: contextDisabled} = useFormElementContext();
  const disabled = useMemo(
    () => (propDisabled ? propDisabled : contextDisabled),
    [contextDisabled, propDisabled],
  );

  return (
    <BasePopoverTrigger
      ref={popoverTriggerRef}
      popover={
        selectMenuOptions.length > 0 ? (
          <SelectMenu
            className={menuClassName}
            size="full"
            options={selectMenuOptions}
            onChange={token => {
              onTokensChange([...tokens, token]);
              onQueryStringChange('');
              inputRef && inputRef.current?.focus();
            }}
          />
        ) : null
      }
      position="below">
      {(onShow, _onHide, contextRef) => (
        <div className={classes.root} ref={contextRef}>
          <TextInput
            ref={inputRef}
            className={classNames(classes.input, inputClassName)}
            disabled={disabled}
            containerClassName={classes.inputContainer}
            prefix={
              <TokensList
                disabled={disabled}
                tokens={tokens}
                onTokensChange={newTokens => {
                  onTokensChange(newTokens);
                }}
              />
            }
            type="text"
            value={queryString}
            onChange={({target}) => {
              onShow();
              onQueryStringChange(target.value);
            }}
            onClick={onShow}
            onBackspacePressed={() => {
              queryString === '' &&
                onTokensChange(tokens.slice(0, tokens.length - 1));
            }}
          />
        </div>
      )}
    </BasePopoverTrigger>
  );
};

export default Tokenizer;
