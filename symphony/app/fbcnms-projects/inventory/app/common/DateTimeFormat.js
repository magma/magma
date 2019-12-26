/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {MomentInput} from 'moment';

import moment from 'moment';

const YEAR_FORMATTER = function(now) {
  return now.year() === this.year() ? 'MMMM D, h:mm A' : 'MMMM D YYYY, h:mm A';
};

const CALENDAR = {
  sameDay: function(now) {
    const diffMins = moment.duration(now.diff(this)).asMinutes();
    return diffMins < 1 ? '[Just now]' : '[Today,] h:mm A';
  },
  nextDay: '[Tomorrow]',
  nextWeek: 'dddd',
  lastDay: function(now) {
    const diffMins = moment.duration(now.diff(this)).asMinutes();
    return diffMins < 1 ? '[Just now]' : '[Yesterday,] h:mm A';
  },
  lastWeek: YEAR_FORMATTER,
  sameElse: YEAR_FORMATTER,
};

export default class DateTimeFormat {
  static commentTime(dateTimeValue: MomentInput): string {
    return moment(dateTimeValue).calendar(null, CALENDAR);
  }

  static dateTime = (
    dateTimeValue: ?string | ?number,
    fallback: string = '',
  ) => {
    // eslint-disable-next-line no-warning-comments
    // $FlowFixMe - Date.parse can handle number and nulls
    const dateTime = Date.parse(dateTimeValue);
    return Number.isNaN(dateTime)
      ? fallback
      : new Intl.DateTimeFormat('default', {
          hour: 'numeric',
          minute: 'numeric',
          year: 'numeric',
          month: 'numeric',
          day: 'numeric',
        }).format(dateTime);
  };
}
