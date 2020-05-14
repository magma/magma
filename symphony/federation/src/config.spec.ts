/**
 * Copyright (c) 2004-present Facebook All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

import { Config } from './config'
import config from './config'
import { GatewayConfig } from '@apollo/gateway'
import { ApolloServerExpressConfig } from 'apollo-server-express'
import { ListenOptions } from 'net'

describe('Config', () => {
    it('should return undefined or empty on no config', () => {
        const Mock = jest.fn(() => ({
            has: jest.fn(),
            get: jest.fn(),
        }))
        const mock = new Mock()
        mock.has.mockImplementation(() => false)

        const config = new Config(mock)
        expect(config.gateway).toBeUndefined()
        expect(config.server).toBeUndefined()
        expect(config.listen).toHaveLength(0)
        expect(mock.has).toBeCalledTimes(3)
        expect(mock.get).not.toBeCalled()
    })

    it('should return config from test.json', () => {
        expect(config.gateway).toMatchObject<GatewayConfig>({ federationVersion: 1 })
        expect(config.server).toMatchObject<ApolloServerExpressConfig>({ introspection: false })
        expect(config.listen).toMatchObject<Array<ListenOptions>>([{ port: 8080 }])
    })
})
