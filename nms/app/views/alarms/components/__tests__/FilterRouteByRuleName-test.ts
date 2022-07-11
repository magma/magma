import * as React from 'react';
import {AlertRoutingTree} from '../AlarmAPIType';
import {filterRouteByRuleName, filterUpdatedFilterRoutes} from '../hooks';

describe('Testing filterRouteByRuleName', () => {
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
    const routes = filterRouteByRuleName(response, undefined, undefined);
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

describe('Testint filterUpdatedFilterRoutes', () => {
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
  test('update only  High Disk Usage Alert from User1 to User2', () => {
    const ruleName = 'High Disk Usage Alert';
    const initialReceiver = 'User1';
    const receiver = 'User2';
    const initialState = response.routes.filter(
      route =>
        route.match?.alertname === ruleName && route.receiver == receiver,
    );
    const routes:
      | Array<AlertRoutingTree>
      | undefined = filterUpdatedFilterRoutes(
      response,
      initialReceiver,
      ruleName,
      receiver,
    );
    if (!!routes) {
      const endState = routes.filter(
        route =>
          route.match?.alertname === ruleName && route.receiver == receiver,
      );

      expect(endState.length).toBeGreaterThan(initialState.length);
      expect(response.routes).not.toBe(routes);
    } else {
      expect(routes).not.toBe(undefined);
    }
  });

  test('return exact routes if initial receiver undefined', () => {
    const ruleName = 'High Disk Usage Alert';
    const initialReceiver = undefined;
    const receiver = 'User2';

    const routes:
      | Array<AlertRoutingTree>
      | undefined = filterUpdatedFilterRoutes(
      response,
      initialReceiver,
      ruleName,
      receiver,
    );
    if (routes) {
      expect(response.routes).toBe(routes);
    }

    const routes2:
      | Array<AlertRoutingTree>
      | undefined = filterUpdatedFilterRoutes(
      response,
      null,
      ruleName,
      receiver,
    );
    if (routes2) {
      expect(response.routes).toBe(routes2);
    }

    const routes3:
      | Array<AlertRoutingTree>
      | undefined = filterUpdatedFilterRoutes(response, '', ruleName, receiver);
    if (routes3) {
      expect(response.routes).toBe(routes3);
    }
  });
  test('return exact routes if receiver undefined', () => {
    const ruleName = 'High Disk Usage Alert';
    const initialReceiver = 'User1';
    const receiverUndefined = undefined;
    const receiverNull = null;
    const receiverEmpty = '';

    const routesUndefined:
      | Array<AlertRoutingTree>
      | undefined = filterUpdatedFilterRoutes(
      response,
      initialReceiver,
      ruleName,
      receiverUndefined,
    );
    const routesNull:
      | Array<AlertRoutingTree>
      | undefined = filterUpdatedFilterRoutes(
      response,
      initialReceiver,
      ruleName,
      receiverNull,
    );
    const routesEmpty:
      | Array<AlertRoutingTree>
      | undefined = filterUpdatedFilterRoutes(
      response,
      initialReceiver,
      ruleName,
      receiverEmpty,
    );

    expect(response.routes).toBe(routesEmpty);
    expect(response.routes).toBe(routesNull);
    expect(response.routes).toBe(routesUndefined);
  });

  test('throws if response is undefined', () => {
    const _response = undefined;
    const ruleName = 'High Disk Usage Alert';
    const initialReceiver = 'User1';
    const receiver = 'User2';
    const error = () => {
      const routesUndefined:
        | Array<AlertRoutingTree>
        | undefined = filterUpdatedFilterRoutes(
        _response,
        initialReceiver,
        ruleName,
        receiver,
      );
    };
    expect(error).toThrow();
  });
});
