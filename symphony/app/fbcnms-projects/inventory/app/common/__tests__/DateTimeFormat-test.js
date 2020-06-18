/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import DateTimeFormat from '../DateTimeFormat';
import moment from 'moment';

const second = 1000;
const minute = 60 * second;

const timeFormat = ', [0-2]?[0-9]:[0-5][0-9] (AM|PM)';

const mockTimes = [
  new Date('2019-10-30T03:24:00'),
  new Date('2019-10-30T13:24:00'),
  new Date('2019-10-30T00:00:01'),
  new Date('2019-01-01T03:24:00'),
  new Date('2019-01-01T13:24:00'),
  new Date('2019-01-01T00:00:01'),
];

describe('commentTime', () => {
  const nowNotion = 'Just now';
  const todayNotion = 'Today';
  const yesterdayNotion = 'Yesterday';
  const originalGlobalDateNow = global.Date.now;

  afterAll(() => {
    global.Date.now = originalGlobalDateNow;
  });

  test(nowNotion, () => {
    mockTimes.forEach(mockTime => {
      global.Date.now = jest.fn(() => mockTime.getTime());

      expect(DateTimeFormat.commentTime(mockTime)).toBe(nowNotion);

      const afterOneSecond = mockTime - second;
      expect(DateTimeFormat.commentTime(afterOneSecond)).toBe(nowNotion);

      const afterLessThanMinute = mockTime - minute + second;
      expect(DateTimeFormat.commentTime(afterLessThanMinute)).toBe(nowNotion);
    });
  });

  test(todayNotion, () => {
    const todayFormat = new RegExp(`^${todayNotion}${timeFormat}$`);
    const isToday = new RegExp(`^${todayNotion}`);
    [mockTimes[0], mockTimes[1], mockTimes[3], mockTimes[4]].forEach(
      mockTime => {
        const mockDate = new Date(mockTime);
        global.Date.now = () => mockTime.getTime();

        const oneMinuteEarlier = mockTime - minute - second;
        expect(DateTimeFormat.commentTime(oneMinuteEarlier)).toMatch(
          todayFormat,
        );

        const beginingOfTheDay = mockDate.setHours(0, 0, 0, 0);
        expect(DateTimeFormat.commentTime(beginingOfTheDay)).toMatch(
          todayFormat,
        );

        const theDayBefore = beginingOfTheDay - 1;
        expect(DateTimeFormat.commentTime(theDayBefore)).not.toMatch(isToday);
      },
    );
  });

  test(yesterdayNotion, () => {
    const yesterdayFormat = new RegExp(`^${yesterdayNotion}${timeFormat}$`);
    const isYesterday = new RegExp(`^${yesterdayNotion}`);
    mockTimes.forEach(mockTime => {
      const mockDate = new Date(mockTime);
      global.Date.now = () => mockTime.getTime();

      const beginingOfTheDay = mockDate.setHours(0, 0, 0, 0);
      const theDayBefore = beginingOfTheDay - 1;
      const YesterdayDate = new Date(theDayBefore);
      const beginingOfYesterday = YesterdayDate.setHours(0, 0, 0, 0);
      expect(DateTimeFormat.commentTime(beginingOfYesterday)).toMatch(
        yesterdayFormat,
      );

      const twoDaysBack = beginingOfYesterday - 1;
      expect(DateTimeFormat.commentTime(twoDaysBack)).not.toMatch(isYesterday);
    });

    [mockTimes[0], mockTimes[1], mockTimes[3], mockTimes[4]].forEach(
      mockTime => {
        const mockDate = new Date(mockTime);
        global.Date.now = () => mockTime.getTime();

        const beginingOfTheDay = mockDate.setHours(0, 0, 0, 0);
        const theEndOfYesterday = beginingOfTheDay - 1;
        expect(DateTimeFormat.commentTime(theEndOfYesterday)).toMatch(
          yesterdayFormat,
        );
      },
    );

    [mockTimes[2], mockTimes[5]].forEach(mockTime => {
      global.Date.now = () => mockTime.getTime();

      const oneMinuteEarlier = mockTime - minute - second;
      expect(DateTimeFormat.commentTime(oneMinuteEarlier)).toMatch(
        yesterdayFormat,
      );
    });
  });

  const months = [
    'January',
    'February',
    'March',
    'April',
    'May',
    'June',
    'July',
    'August',
    'September',
    'October',
    'November',
    'December',
  ];
  const possibleMonths = months.join('|');

  const dateFormat = `(${possibleMonths}) ([1-9]|[1-2][0-9]|30|31)`;
  const fullDateFormat = new RegExp(`^${dateFormat}${timeFormat}$`);

  test('Before yesterday (this year)', () => {
    const now = new Date('2019-12-13T04:41:20');
    global.Date.now = jest.fn(() => now);

    const twoDaysAgo = moment(now).subtract(2, 'days');

    const thisYearPossibleMonths = Array.from(
      {length: twoDaysAgo.month()},
      (_, i) => i,
    ).reverse();

    const mockExactDates = [
      ...thisYearPossibleMonths.map(monthInd => ({
        date: twoDaysAgo.month(monthInd).toDate(),
        verify: months[monthInd],
      })),
      {
        date: twoDaysAgo.set({hour: 3, minute: 24}).toDate(),
        verify: '3:24 AM',
      },
      {
        date: twoDaysAgo.set({hour: 13, minute: 24}).toDate(),
        verify: '1:24 PM',
      },
    ];

    mockExactDates.forEach(mockDate => {
      expect(DateTimeFormat.commentTime(mockDate.date)).toMatch(fullDateFormat);
      expect(DateTimeFormat.commentTime(mockDate.date)).toContain(
        mockDate.verify,
      );
    });
  });

  const dateFormatWithYear = `${dateFormat} [1-2][0-9]{3}`;
  const fullDateFormatWithYear = new RegExp(
    `^${dateFormatWithYear}${timeFormat}$`,
  );

  test('Before this year (and yesterday)', () => {
    const now = new Date('2019-12-13T04:41:20');
    global.Date.now = jest.fn(() => now);
    const lastYear = moment(now).subtract({year: 1});

    const dateCheck = (month, date) => ({
      date: lastYear
        .clone()
        .set({month, date})
        .toDate(),
      verify: `${months[month]} ${date}`,
    });

    const assertDate = mockDate => {
      expect(DateTimeFormat.commentTime(mockDate.date)).toMatch(
        fullDateFormatWithYear,
      );
      expect(DateTimeFormat.commentTime(mockDate.date)).toContain(
        mockDate.verify,
      );
    };

    for (let i = 0; i < months.length; i++) {
      assertDate({
        date: lastYear
          .clone()
          .month(i)
          .toDate(),
        verify: months[i],
      });
    }

    assertDate({
      date: lastYear
        .clone()
        .hour(13)
        .minute(48)
        .toDate(),
      verify: '1:48 PM',
    });
    assertDate({
      date: lastYear
        .clone()
        .hour(3)
        .minute(24)
        .toDate(),
      verify: '3:24 AM',
    });
    [0, 2, 4, 6, 7, 9, 11].forEach(longMonth => {
      assertDate(dateCheck(longMonth, 30));
      assertDate(dateCheck(longMonth, 31));
    });
    [3, 5, 8, 10].forEach(shortMonth => {
      assertDate(dateCheck(shortMonth, 30));
    });
  });
});
