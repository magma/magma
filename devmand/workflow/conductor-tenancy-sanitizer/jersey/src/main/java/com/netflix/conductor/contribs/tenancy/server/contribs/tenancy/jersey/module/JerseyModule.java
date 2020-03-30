/**
 * Copyright 2016 Netflix, Inc.
 * <p>
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * <p>
 * http://www.apache.org/licenses/LICENSE-2.0
 * <p>
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package com.netflix.conductor.contribs.tenancy.server.contribs.tenancy.jersey.module;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.jaxrs.json.JacksonJsonProvider;
import com.google.inject.Provides;
import com.netflix.conductor.contribs.tenancy.server.contribs.tenancy.jersey.filter.auth.AuthFilter;
import com.netflix.conductor.contribs.tenancy.server.contribs.tenancy.jersey.filter.CORSFilter;
import com.netflix.conductor.contribs.tenancy.server.contribs.tenancy.jersey.resources.WorkflowResource;
import com.sun.jersey.api.core.PackagesResourceConfig;
import com.sun.jersey.api.core.ResourceConfig;
import com.sun.jersey.guice.JerseyServletModule;
import com.sun.jersey.guice.spi.container.servlet.GuiceContainer;
import java.util.HashMap;
import java.util.Map;
import javax.inject.Singleton;

public final class JerseyModule extends JerseyServletModule {

    @Override
    protected void configureServlets() {
        Map<String, String> jerseyParams = new HashMap<>();
        jerseyParams.put("com.sun.jersey.config.feature.FilterForwardOn404", "true");
        jerseyParams.put("com.sun.jersey.config.property.WebPageContentRegex", "/(((webjars|api-docs|swagger-ui/docs|manage)/.*)|(favicon\\.ico))");
        jerseyParams.put(PackagesResourceConfig.PROPERTY_PACKAGES,
                WorkflowResource.class.getPackage().getName() +
                        ";io.swagger.jaxrs.json;io.swagger.jaxrs.listing");
        jerseyParams.put(ResourceConfig.FEATURE_DISABLE_WADL, "false");
        serve("/api/*").with(GuiceContainer.class, jerseyParams);

        filter("/*").through(AuthFilter.class);
        filter("/*").through(CORSFilter.class);

    }

    @Provides
    @Singleton
    JacksonJsonProvider jacksonJsonProvider(ObjectMapper mapper) {
        return new JacksonJsonProvider(mapper);
    }

    @Override
    public boolean equals(Object obj) {
        return obj != null && getClass().equals(obj.getClass());
    }

    @Override
    public int hashCode() {
        return getClass().hashCode();
    }


}
