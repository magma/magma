package com.netflix.conductor.contribs.tenancy.server.contribs.tenancy.jersey.config;

import javax.inject.Singleton;

@Singleton
public class ProxyConfiguration {

    /**
     * Base conductor URL.
     * Example: when api endpoints are at http://host:port/api/..., this method
     * should return http://host:port
      */
    public String getBaseConductorURL() {
        return "http://localhost:8080";
    }
}
