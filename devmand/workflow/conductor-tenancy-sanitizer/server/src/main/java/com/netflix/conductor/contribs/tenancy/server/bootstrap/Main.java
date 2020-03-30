/*
 * Copyright 2017 Netflix, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */
package com.netflix.conductor.contribs.tenancy.server.bootstrap;

import com.google.inject.Guice;
import com.google.inject.Injector;
import com.netflix.conductor.contribs.tenancy.server.jetty.server.JettyServerProvider;
import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.lang.invoke.MethodHandles;
import java.net.URL;
import org.apache.log4j.PropertyConfigurator;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class Main {

    public static void main(String[] args) {
        if (args.length == 0) {
            URL log4j = Main.class.getResource("/log4j.properties");
            PropertyConfigurator.configure(log4j);
        } else {
            String log4jFileName = args[0];
            PropertyConfigurator.configureAndWatch(log4jFileName);
            LoggerFactory.getLogger(MethodHandles.lookup().lookupClass()).info("Configured log4j using {}", log4jFileName);
        }


        ModulesProvider modulesProvider = new ModulesProvider();
        Injector serverInjector = Guice.createInjector(modulesProvider.get());

        serverInjector.getInstance(JettyServerProvider.class).get().ifPresent(server -> {
            try {
                server.start();
            } catch (Exception ioe) {
                ioe.printStackTrace(System.err);
                System.exit(3);
            }
        });
    }
}
