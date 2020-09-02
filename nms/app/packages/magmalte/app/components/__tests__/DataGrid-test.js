import 'jest-dom/extend-expect';
import React from 'react';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import { MuiThemeProvider } from '@material-ui/core/styles';
import { cleanup, render, waitFor, fireEvent } from '@testing-library/react';
import DataGrid, { DataRows } from '../DataGrid';
import defaultTheme from '../../theme/default';

afterEach(cleanup);

const data: DataRows[] = [
  [
    {
      category: 'Total',
      value: 'eNodeBs',
      tooltip: 'Tooltip text',
    },
    {
      category: 'Severe Events',
      value: 'Value used as a tooltip',
    },
    {
      category: 'Max Latency',
      value: 100,
      unit: 'ms'
    },
  ],
];

const Wrapper = () => {
  return (
    <MuiThemeProvider theme={defaultTheme}>
      <MuiStylesThemeProvider theme={defaultTheme}>
        <DataGrid data={data} />
      </MuiStylesThemeProvider>
    </MuiThemeProvider>
  )
}

describe('<DataGrid />', () => {
  it('displays the passed tooltip', async () => {
    const { getByTestId, getByRole } = render(<Wrapper />);

    const dataEntryToHoverOver = getByTestId('Total');
    fireEvent.mouseOver(dataEntryToHoverOver);

    await waitFor(() => {
      expect(getByRole('tooltip')).toHaveTextContent(data[0][0].tooltip);
    });
  });

  it('defaults to the data entry value when the tooltip prop in not passed', async () => {
    const { getByTestId, getByRole } = render(<Wrapper />);

    const dataEntryToHoverOver = getByTestId('Severe Events');
    fireEvent.mouseOver(dataEntryToHoverOver);

    await waitFor(() => {
      expect(getByRole('tooltip')).toHaveTextContent(data[0][1].value);
    });
  });

  it('displays the data unit along with data value when unit prop is passed', async () => {
    const { getByTestId, getByRole } = render(<Wrapper />);

    const dataEntryToHoverOver = getByTestId('Max Latency');
    fireEvent.mouseOver(dataEntryToHoverOver);

    await waitFor(() => {
      const { value, unit } = data[0][2];
      expect(getByRole('tooltip')).toHaveTextContent(value + unit);
    });
  });
});
