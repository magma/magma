/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

jest.mock('@fbcnms/sequelize-models');

import OrganizationLocalStrategy from '@fbcnms/auth/strategies/OrganizationLocalStrategy';

import bodyParser from 'body-parser';
import express from 'express';
import fbcPassport from '../passport';
import passport from 'passport';
import request from 'supertest';
import userMiddleware from '../express';
import {USERS, USERS_EXPECTED} from '../test/UserModel';
import {User} from '@fbcnms/sequelize-models';

import {configureAccess} from '../access';

import type {FBCNMSRequest} from '../access';

function stripDatesMany(res) {
  res.body.users.map(u => {
    delete u['createdAt'];
    delete u['updatedAt'];
  });
}

function stripDates(res) {
  delete res.body.user['createdAt'];
  delete res.body.user['updatedAt'];
}

function mockOrgMiddleware(orgName: string) {
  return (req: FBCNMSRequest, _res, next) => {
    if (orgName) {
      // $FlowIgnore we know this is wrong, and that's okay for this test.
      req.organization = async () => {
        return {name: orgName};
      };
    }
    next();
  };
}

function mockUserMiddleware(where) {
  return async (req: FBCNMSRequest, _res, next) => {
    req.isAuthenticated = () => true;
    const user = await User.findOne({where});
    if (!user) {
      throw new Error('Could not find a user');
    }
    req.user = user;
    next();
  };
}

function getApp(orgName: string, loggedInEmail: ?string) {
  const app = express();
  app.use(bodyParser.json());
  app.use(bodyParser.urlencoded());
  fbcPassport.use();
  passport.use(new OrganizationLocalStrategy());
  app.use(passport.initialize());
  app.use(passport.session());
  app.use(configureAccess({loginUrl: '/user/login'}));
  app.use(mockOrgMiddleware(orgName));
  if (loggedInEmail) {
    app.use(mockUserMiddleware({email: loggedInEmail}));
  }
  app.use(
    '/user',
    userMiddleware({
      loginFailureUrl: '/failure',
      loginSuccessUrl: '/success',
    }),
  );
  return app;
}

describe('user tests', () => {
  beforeEach(async () => {
    USERS.forEach(async user => await User.create(user));
  });
  afterEach(async () => {
    await User.destroy({where: {}, truncate: true});
  });

  describe('login', () => {
    describe('with organization', () => {
      it('valid user can login', async () => {
        const app = getApp('validorg');
        await request(app)
          .post('/user/login')
          .type('form')
          .send({email: 'valid@123.com', password: 'password1234'})
          .expect(302)
          .expect('Location', '/success');
      });

      it('valid user can login (redirected)', async () => {
        const app = getApp('validorg');
        await request(app)
          .post('/user/login')
          .type('form')
          .send({
            email: 'valid@123.com',
            password: 'password1234',
            to: '/other/success',
          })
          .expect(302)
          .expect('Location', '/other/success');
      });

      it('valid user can login (non-relative redirect)', async () => {
        const app = getApp('validorg');
        await request(app)
          .post('/user/login')
          .type('form')
          .send({
            email: 'valid@123.com',
            password: 'password1234',
            to: 'http://evilsite.com/other/success',
          })
          .expect(302)
          .expect('Location', '/success');

        await request(app)
          .post('/user/login')
          .type('form')
          .send({
            email: 'valid@123.com',
            password: 'password1234',
            to: '//evilsite.com/other/success',
          })
          .expect(302)
          .expect('Location', '/success');
      });

      it('valid user, invalid org cant login', async () => {
        const app = getApp('invalidorg');
        await request(app)
          .post('/user/login')
          .type('form')
          .send({email: 'valid@123.com', password: 'password1234'})
          .expect(302)
          .expect('Location', '/failure');
      });

      it('invalid user cant login', async () => {
        const app = getApp('validorg');
        await request(app)
          .post('/user/login')
          .type('form')
          .send({email: 'invalid@123.com', password: 'password1234'})
          .expect(302)
          .expect('Location', '/failure');
      });
    });

    describe('no organization', () => {
      it('valid user can login', async () => {
        const app = getApp('');
        await request(app)
          .post('/user/login')
          .type('form')
          .send({email: 'noorg@123.com', password: 'password1234'})
          .expect(302)
          .expect('Location', '/success');
      });

      it('invalid user cant login', async () => {
        const app = getApp('');
        await request(app)
          .post('/user/login')
          .type('form')
          .send({email: 'invalid@123.com', password: 'password1234'})
          .expect(302)
          .expect('Location', '/failure');
      });
    });
  });

  describe('create user', () => {
    const validParams = {
      email: 'user@123.com',
      password: 'password1234',
      networkIDs: [],
      superUser: false,
      verificationType: 0,
      role: 0,
    };

    describe('as superuser', () => {
      const app = getApp('validorg', 'superuser@123.com');

      it('creates users successfully', async () => {
        const params = validParams;
        await request(app)
          .post('/user/async/')
          .send(params)
          .expect(201)
          .expect(stripDates)
          .expect({
            user: {
              isSuperUser: false,
              isReadOnlyUser: false,
              email: params.email,
              organization: 'validorg',
              networkIDs: params.networkIDs,
              role: 0,
              id: 5,
              tabs: [],
            },
          });
      });

      it('must supply email', async () => {
        const params = {
          ...validParams,
          email: '',
        };
        await request(app)
          .post('/user/async/')
          .send(params)
          .expect(400);
      });

      it('must be valid email', async () => {
        const params = {
          ...validParams,
          email: 'abc',
        };
        await request(app)
          .post('/user/async/')
          .send(params)
          .expect(400);
      });
    });
  });

  describe('list users', () => {
    it('lists users for organization', async () => {
      const app = getApp('validorg', 'superuser@123.com');
      await request(app)
        .get('/user/async/')
        .expect(200)
        .expect(stripDatesMany)
        .expect({
          users: [USERS_EXPECTED[0], USERS_EXPECTED[2]],
        });
    });
    it('list users (no organization)', async () => {
      const app = getApp('', 'superuser@123.com');
      await request(app)
        .get('/user/async/')
        .expect(200)
        .expect(stripDatesMany)
        .expect({
          users: USERS_EXPECTED,
        });
    });
  });

  describe('edit users', () => {
    const validUpdateParams = {
      networkIDs: [],
      password: 'mynewpassword',
      superUser: false,
      verificationType: 0,
      tabs: ['validtab'],
    };
    it('can update a user', async () => {
      const app = getApp('validorg', 'superuser@123.com');
      await request(app)
        .put('/user/async/1')
        .send(validUpdateParams)
        .expect(stripDates)
        .expect({
          user: {
            networkIDs: [],
            id: 1,
            email: 'valid@123.com',
            organization: 'validorg',
            role: 0,
            tabs: ['validtab'],
          },
        })
        .expect(200);
    });

    it('cannot edit another orgs user', async () => {
      const app = getApp('validorg', 'superuser@123.com');
      await request(app)
        .put('/user/async/2')
        .send(validUpdateParams)
        .expect(400)
        .expect({error: 'Error: User does not exist!'});
    });
  });

  describe('endpoints as normal user', () => {
    const app = getApp('validorg', 'valid@123.com');
    it('redirects restricted urls to login', async () => {
      await request(app)
        .get('/user/async/')
        .expect(302);
      await request(app)
        .post('/user/async/')
        .expect(302);
      await request(app)
        .put('/user/async/1')
        .expect(302);
      await request(app)
        .delete('/user/async/1/')
        .expect(302);
    });
  });
});
