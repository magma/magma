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
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import defaultTheme from '@fbcnms/ui/theme/default';
import {MemoryRouter} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {SnackbarProvider} from 'notistack';

import {act, render} from '@testing-library/react';

export function SymphonyWrapper({
  route,
  children,
}: {
  route?: string,
  children: React.Node,
}) {
  return (
    <MemoryRouter initialEntries={[route || '/']} initialIndex={0}>
      <MuiThemeProvider theme={defaultTheme}>
        <MuiStylesThemeProvider theme={defaultTheme}>
          <SnackbarProvider>{children}</SnackbarProvider>
        </MuiStylesThemeProvider>
      </MuiThemeProvider>
    </MemoryRouter>
  );
}

/*
 * Use this if a component asyncronously loads data when it is rendered.
 *
 * For example:
 *
 * function MyComponent(){
 *  const [data, setData] = useState(null);
 *  useEffect(() => {
 *    axios.get().then(response => {
 *      setData(response.data)
 *    })
 *  }, []);
 *  return <div>{data}</div>
 * }
 *
 * since the setData call happens asyncronously, react test renderer will
 * complain that you've modified state outside of an act() call.
 *
 * if your component needs to load data asyncronously on mount, replace:
 *
 * const {getByText} = render(<MyComponent/>);
 * with
 * const {getByText} = await renderAsync(<MyComponent/>);
 */
// eslint-disable-next-line flowtype/no-weak-types
export async function renderAsync(...renderArgs: Array<any>): Promise<any> {
  let result;
  await act(async () => {
    result = await render(...renderArgs);
  });
  return result;
}
