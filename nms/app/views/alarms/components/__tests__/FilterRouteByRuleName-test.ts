import * as React from 'react';
import {filterRouteByRuleName} from '../hooks';

describe('Testing routes', () => {
  const response = {
    receiver: 'test_tenant_base_route',
    match: {networkID: 'test'},
    routes: [
      {receiver: 'User1', match: {alertname: 'High Disk Usage Alert'}},
      {
        receiver: 'User1',
        match: {alertname: 'Certificate Expiring Soon'},
      },
      {
        receiver: 'User1',
        match: {alertname: 'Bootstrap Exception Alert'},
      },
      {
        receiver: 'User2',
        match: {alertname: 'Gateway Checkin Failure'},
      },
    ],
  };
  test('Remove only  High Disk Usage Alert from User1', () => {
    const ruleName = 'High Disk Usage Alert';
    const initialReceiver = 'User1';
    const routes = filterRouteByRuleName(response, initialReceiver, ruleName);
    expect(routes.length).toBe(3);
  });

  test('Remove only  Gateway Checkin Failure from User2', () => {
    const ruleName = 'Gateway Checkin Failure';
    const initialReceiver = 'User2';
    const routes = filterRouteByRuleName(response, initialReceiver, ruleName);
    expect(routes.length).toBe(3);
  });

  test('return all routes if initialReceiver is empty', () => {
    const ruleName = 'Gateway Checkin Failure';
    const initialReceiver = '';
    const routes = filterRouteByRuleName(response, initialReceiver, ruleName);
    expect(routes.length).toBe(4);
  });

  test('return all routes if ruleName is empty', () => {
    const ruleName = '';
    const initialReceiver = 'User2';
    const routes = filterRouteByRuleName(response, initialReceiver, ruleName);
    expect(routes.length).toBe(4);
  });

  test('return all routes if initialReceiver is undefined', () => {
    const ruleName = 'Gateway Checkin Failure';
    const initialReceiver = undefined;
    const routes = filterRouteByRuleName(response, initialReceiver, ruleName);
    expect(routes.length).toBe(4);
  });

  test('return all routes if ruleName is undefined', () => {
    const ruleName = undefined;
    const initialReceiver = 'User2';
    const routes = filterRouteByRuleName(response, initialReceiver, ruleName);
    expect(routes.length).toBe(4);
  });

  test('return all routes if ruleName and  initialReceiver is undefined', () => {
    const routes = filterRouteByRuleName(response);
    expect(routes.length).toBe(4);
  });

  test('throws if response is undefined', () => {
    const _response = undefined;
    const ruleName = 'High Disk Usage Alert';
    const initialReceiver = 'User1';
    const error = () => {
      const routes = filterRouteByRuleName(
        _response,
        initialReceiver,
        ruleName,
      );
    };
    expect(error).toThrow();
  });
});
