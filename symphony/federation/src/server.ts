/**
 * Copyright (c) 2004-present Facebook All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

import { ApolloServer, ServerInfo } from 'apollo-server'
import { ApolloGateway } from '@apollo/gateway'
import config from './config'

const gateway = new ApolloGateway(config.gateway)

const server = new ApolloServer({
    gateway,
    subscriptions: false,
    ...config.server,
})

server.listen(...config.listen).then(({ url }: ServerInfo) => {
    console.log(`🚀 Server ready at ${url}`)
})
