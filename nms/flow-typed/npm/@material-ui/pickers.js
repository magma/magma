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
 * @format
 * @flow strict-local
 */

type KeyboardDatePickerProps = {
  disableToolbar?: boolean,
  inputVariant?: string,
  format?: string,
  margin?: string,
  value?: string | Date,
  onChange?: Object => any,
  KeyboardButtonProps?: Object,
};

type DatePickerProps = {
  autoOk?: boolean,
  disableFuture?: boolean,
  variant?: string,
  inputVariant?: string,
  inputProps?: Object,
  maxDate?: ParsableDate,
  format?: string,
  value?: moment$Moment,
  onChange?: Object => any,
};

declare module '@material-ui/pickers/KeyboardDatePicker' {
  declare module.exports: React$ComponentType<KeyboardDatePickerProps>;
}

declare module '@material-ui/pickers/DateTimePicker' {
  declare module.exports: React$ComponentType<DatePickerProps>;
}

declare module '@material-ui/pickers/MuiPickersUtilsProvider' {
  declare module.exports: React$ComponentType<KeyboardDatePickerProps>;
}

declare module '@material-ui/pickers' {
  declare module.exports: {
    KeyboardDatePicker: $Exports<'@material-ui/pickers/KeyboardDatePicker'>,
    KeyboardTimePicker: $Exports<'@material-ui/pickers/KeyboardTimePicker'>,
    DateTimePicker: $Exports<'@material-ui/pickers/DateTimePicker'>,
    MuiPickersUtilsProvider: $Exports<
      '@material-ui/pickers/MuiPickersUtilsProvider',
    >,
  };
}
