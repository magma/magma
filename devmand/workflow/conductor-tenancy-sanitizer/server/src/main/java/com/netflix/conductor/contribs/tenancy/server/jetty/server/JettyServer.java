/*
 * Copyright 2017 Netflix, Inc.
 * <p>
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * <p>
 * http://www.apache.org/licenses/LICENSE-2.0
 * <p>
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */
package com.netflix.conductor.contribs.tenancy.server.jetty.server;

import com.google.inject.servlet.GuiceFilter;
import java.util.EnumSet;
import javax.servlet.DispatcherType;
import org.eclipse.jetty.server.Server;
import org.eclipse.jetty.servlet.ServletContextHandler;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * @author Viren
 */
public class JettyServer {

    private static Logger logger = LoggerFactory.getLogger(JettyServer.class);

    private final int port;
    private final boolean join;

    private Server server;


    public JettyServer(int port, boolean join) {
        this.port = port;
        this.join = join;
    }

    public synchronized void start() throws Exception {

        if (server != null) {
            throw new IllegalStateException("Server is already running");
        }

        this.server = new Server(port);

        ServletContextHandler context = new ServletContextHandler();
        context.addFilter(GuiceFilter.class, "/*", EnumSet.allOf(DispatcherType.class));
        context.setWelcomeFiles(new String[]{"index.html"});

        server.setHandler(context);
        server.start();
        System.out.println("Started server on http://localhost:" + port + "/");


        if (join) {
            server.join();
        }

    }

    public synchronized void stop() throws Exception {
        if (server == null) {
            throw new IllegalStateException("Server is not running.  call #start() method to start the server");
        }
        server.stop();
        server = null;
    }

}
